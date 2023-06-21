go
create  or alter  procedure [dbo].[FormsCreateResponseHeader] @FormName varchar(35), @ClientIPAddress varchar(75), @RecaptchaV3Score money, @LoginSessionID nvarchar(255)
as
    set nocount on

declare @webFormID int
select @webFormID = wf.ID from WebForms wf where wf.FormName = @FormName

    if @webFormID is null
        begin
            insert into AuditTrail(TableName,SessionID,EventDescription,NewValue) select 'WebFormResponseHeader',@LoginSessionID,'Missing or Invalid Form Name',@FormName

            select 'Missing or Invalid Form Name' as [Result], 0 as [ID]
            return
        end

declare @userID int, @username nvarchar(255)
select @username = Username from LoginSessionLog where SessionID = @LoginSessionID

select @userID = u.ID from Users u where u.AuthUsername = @username


declare @outputIDTable table (HeaderID int)

insert into [WebFormResponseHeader](FormID,ClientIPAddress,RecaptchaV3Score,LoginSessionID,UserID) output inserted.ID into @outputIDTable
select @webFormID, @ClientIPAddress, @RecaptchaV3Score,@LoginSessionID,@userID

declare @headerID int
select @headerID = HeaderID from @outputIDTable

    if @headerID is null
        begin
            insert into AuditTrail(TableName,SessionID,EventDescription,NewValue) select 'WebFormResponseHeader',@LoginSessionID,'Error inserting into [WebFormResponseHeader]',@FormName

            select 'Error inserting into [WebFormResponseHeader]' as [Result], 0 as [ID]
            return
        end

insert into AuditTrail(TableName,DataID,EventDescription,SessionID,NewValue) select 'WebFormResponseHeader',@headerID,'Added new form response header',@LoginSessionID,@FormName

select 'OK' as [Result], @headerID as [ID]
go

GO
-----------------------------------------------------------------------------------------------------------------------
go
create or alter  procedure [dbo].[FormsGetQuestions] @FormName varchar(35)
as
    set nocount on
select wfq.ID, wfq.FormID, wfq.InputFieldName, wfq.EntityName, wfq.QuestionText, wfq.QuestionType, wfq.TemplateName,
       isnull(wfq.Seq,0) as [Seq],
       isnull(wfq.SubSeq,0) as [SubSeq],
       wfq.AnswerRequired,
       isnull(wfq.TextRows,0) as [TextRows],
       isnull(wfq.RadioOption,'') as [RadioOption],
       isnull(wqt.QuestionTemplateText,'') as [Template1],
       isnull(wqt2.QuestionTemplateText,'') as [Template2],
       isnull(wqt3.QuestionTemplateText,'') as [Template3],
       isnull(wfq.QuestionLabel,'') as [QuestionLabel],
       isnull(wfq.[QuestionExtraText],'') as [QuestionExtraText]
from WebFormQuestions wfq
         join WebForms wf on wfq.FormID=wf.ID
         left join WebFormQuestionTemplates wqt on (wqt.QuestionTemplateName =  wfq.TemplateName )
         left join WebFormQuestionTemplates wqt2 on wqt2.QuestionTemplateName = case when wfq.TemplateName = 'RadioButton' then 'RadioButtonHeader' end
         left join WebFormQuestionTemplates wqt3 on wqt3.QuestionTemplateName = case when wfq.TemplateName = 'RadioButton' then 'RadioButtonDetail' end
where wf.FormName = @FormName and isnull(wfq.[Hidden],'N')='N'
order by wfq.seq,wfq.SubSeq,wfq.ID
GO
-----------------------------------------------------------------------------------------------------------------------
GO
create or alter  procedure [dbo].[FormsGetResponseHeader] @HeaderID int
as
    set nocount on

select hdr.ID,FormID,  wf.FormName, wf.FormDescription,
    hdr.DateSubmitted, hdr.TimeSubmitted
     ,isnull(hdr.ClientIPAddress,'') as [ClientIPAddress]
     ,isnull(hdr.RecaptchaV3Score,-1) as [RecaptchaV3Score]
     ,isnull(hdr.LoginSessionID,'') as [LoginSessionID]
     ,isnull(wf.DefaultEmailReportRecipient,'') as [DefaultEmailReportRecipient]
     ,isnull(convert(varchar(10),[EmailReportSent],103),'') as [EmailReportSent]
     ,isnull(convert(varchar(19),[EmailedTo],126),'') as [EmailedTo]
from [WebFormResponseHeader] hdr
         join [WebForms] wf on hdr.FormID = wf.ID
where hdr.ID = @HeaderID
GO
-----------------------------------------------------------------------------------------------------------------------
go
create or alter procedure [dbo].[FormsGetResponseDetails] @HeaderID int
as
    set nocount on
select rd.HeaderID, rd.QuestionID, isnull(q.Seq,rd.QuestionID) as [Seq], q.EntityName, q.QuestionText, rd.QuestionResponse from WebFormResponseDetail rd
                                                                                                                                    join WebFormQuestions q on q.ID=rd.QuestionID
where rd.HeaderID = @HeaderID
order by q.Seq, q.ID
GO
-----------------------------------------------------------------------------------------------------------------------
go
create  or alter procedure FormsSaveResponseDetail @HeaderID int, @InputFieldName varchar(35), @QuestionResponse nvarchar(max)
as
    set nocount on

declare @QuestionID int, @FormID int
select @QuestionID = q.ID, @FormID = rh.FormID
from [WebFormQuestions] q
         join [WebFormResponseHeader] rh on rh.ID = @HeaderID
where q.FormID = rh.FormID and q.InputFieldName = @InputFieldName

    if @QuestionID is null or @FormID is null
        begin
            insert into AuditTrail(TableName,DataID,EventDescription,NewValue) select 'WebFormResponseHeader',@HeaderID,'Error searching for Question ID',@InputFieldName

            select 'Error searching for Question ID for Form ID # '+isnull(convert(varchar,@FormID),'')
                       +', @HeaderID # '+isnull(convert(varchar,@HeaderID),'')
                       +',  @InputFieldName: '+isnull(convert(varchar,@InputFieldName),'')  as [Result], 0 as [ID]
            return
        end

declare @outputIDTable table (DetailID int)

insert into WebFormResponseDetail(HeaderID,QuestionID,QuestionResponse) output inserted.ID into @outputIDTable
select @HeaderID, @QuestionID, @QuestionResponse

declare @detailID int
select @detailID = DetailID from @outputIDTable

    if @@ROWCOUNT = 0 or @detailID is null
        begin

            insert into AuditTrail(TableName,DataID,EventDescription,NewValue) select 'WebFormResponseHeader',@HeaderID,'Error inserting answer for Question', isnull(@InputFieldName,'')+' : '+ isnull(@QuestionResponse,'')

            select 'Error inserting answer for Question ID # '+convert(varchar,@QuestionID)
                       +', Form ID # '+convert(varchar,@FormID)+', @HeaderID # '+convert(varchar,@HeaderID)
                       +',  @InputFieldName: '+convert(varchar,@InputFieldName)  as [Result], 0 as [ID]
            return
        end

select 'OK' as [Result], @detailID as [ID]
go
-----------------------------------------------------------------------------------------------------------------------
-----------------------------------------------------------------------------------------------------------------------
