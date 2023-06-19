-- noinspection SqlNoDataSourceInspectionForFile

go
GO
create or alter   procedure [dbo].[LoginsLogSession] @SessionID nvarchar(255)='', @Username nvarchar(255)='', @Nickname nvarchar(255)='', @Picture nvarchar(255)='',
                                                     @UserId nvarchar(255)='', @Email nvarchar(255)='', @EmailVerified nvarchar(255)='', @GivenName nvarchar(255)='', @FamilyName nvarchar(255)='', @ClientIP nvarchar(255)=''
as
    set nocount on
declare @existingID int

select @existingID = ID
from LoginSessionLog l
where ((l.GivenName = @GivenName and l.FamilyName = @FamilyName) or l.Username=@Username or l.UserId=@UserId or l.Email=@Email) and @SessionID=l.SessionID and @ClientIP=ClientIP

    if @existingID is null
        begin

            -- add new logins entry
            declare @outputIDTable table (logID int)

            insert into LoginSessionLog(SessionID, Username, Nickname, Picture, UserId, Email, EmailVerified, GivenName, FamilyName, ClientIP )
            output inserted.ID into @outputIDTable
            select @SessionID, @Username, @Nickname, @Picture, @UserId, @Email, @EmailVerified, @GivenName, @FamilyName, @ClientIP

            declare @logId int
            select @logID = logID from @outputIDTable
            insert into AuditTrail(TableName,DataID,EventDescription,SessionID,NewValue) select 'LoginSessionLog',@logID,'Added new LoginSessionLog',@SessionID,'Username: '+isnull(@Username,'')+', Email: '+isnull(@email,'')
        end else begin
        update LoginSessionLog set SessionLastAccess = getdate() where ID = @existingID
        insert into AuditTrail(TableName,DataID,EventDescription,SessionID,NewValue) select 'LoginSessionLog',@existingID,'Updated LoginSessionLog',@SessionID,'Username: '+isnull(@Username,'')+', Email: '+isnull(@email,'')
    end

select 'OK' as [Result], [MembershipStatus], [UsersUserID], [MemberId], '|'+STRING_AGG([Role],'|')+'|' as [Roles]
from (
	select distinct isnull(m.MembershipStatus,'') as [MembershipStatus], 
		isnull(u.ID,0) as [UsersUserID], isnull(mul.MemberID,0) as [MemberId], isnull(r.RoleName,'') as [Role]
	from LoginSessionLog lsl
			 left join Users u on u.AuthUsername=lsl.Username
			 left join MemberUserLogin mul on mul.UserID=u.ID
			 left join Members m on m.ID=mul.MemberID
			 left join UsersRoles ur on ur.UserID=u.ID
			 left join Roles r on r.ID=ur.RoleID
	where  lsl.SessionID=@SessionID
	) s
group by [MembershipStatus], [UsersUserID], [MemberId] 

return
GO













-----------------------------------------------------------------------------------------------------------------------
go
create or alter procedure UserAddOrUpdate   @Username nvarchar(255)='', @Nickname nvarchar(255)='', @Picture nvarchar(255)='',
                                            @AuthUserId nvarchar(255)='', @Email nvarchar(255)='', @EmailVerified nvarchar(255)='',
                                            @GivenName nvarchar(255)='', @FamilyName nvarchar(255)=''
as
    set nocount on
declare @ID int

select @ID = ID
from Users
where ((GivenName = @GivenName and FamilyName = @FamilyName) or Username=@Username or AuthUserId=@AuthUserId or Email=@Email)

    if @ID is null
        begin
            declare @outputTable table (ID int)
            insert into Users(Username, Nickname, Picture, AuthUserId, Email, EmailVerified, GivenName, FamilyName )
            output inserted.ID into @outputTable
            select  @Username, @Nickname, @Picture, @AuthUserId, @Email, @EmailVerified, @GivenName, @FamilyName

            select @ID = ID from @outputTable

            insert into AuditTrail(TableName,DataID,EventDescription,UserID,NewValue) select 'Users',@ID,'Added new User',@ID,'Username: '+isnull(@Username,'')+', Email: '+isnull(@email,'')
        end else begin
        insert into AuditTrail(TableName,DataID,EventDescription,UserID,NewValue) select 'Users',@ID,'Updated existing User',@ID,'Username: '+isnull(@Username,'')+', Email: '+isnull(@email,'')
    end

select ID, Username, Nickname, Picture, AuthUserId, Email, EmailVerified, GivenName, FamilyName
from Users
where ID = @ID
go
-----------------------------------------------------------------------------------------------------------------------
go
create or alter procedure UserAddUpdateFromLoginSessionAndMembershipAppForm @LoginSessionID NVARCHAR(255), @FormHeaderID int
as
    set nocount on

    /*
    the general assumption is that a brand new user will use the Membership Application Form as opposed to an existing user who would use the Membership Maintenance form.

    so the assumed workflow is:

    1) user authenticates via Auth0 and reaches the home page or even directly to the new membership application form
    2) retrieve the auth details via the Auth0 profile, generate a session ID, and save to the LoginSessionLog table (and save the session id to the session cookie)
    3) user fills out the new membership application form form, which includes a mandatory email field, which gets saved - and we get a HeaderID in return
    4) now we can link the HeaderID and the LoginSessionLog session id to add or update the Users table
    5) the other form fields are then used to insert or update the Members table
    6) update the MemberUserLogin table to link the [Users] record to the [Members] record

    */
--declare @LoginSessionID NVARCHAR(255), @FormHeaderID int; select @LoginSessionID=8893817, @FormHeaderID=29
    if not exists(select 1 from LoginSessionLog where SessionID=@LoginSessionID)
        begin
            select 'Error: Invalid or missing LoginSessionID' as [Result], -1 as [UserID]
            return
        end

    if not exists(select 1 from WebFormResponseDetail det join WebFormQuestions q on det.QuestionID=q.ID and q.InputFieldName='Email' and det.HeaderID=@FormHeaderID)
        begin
            select 'Error: Invalid or missing Form Email field' as [Result], -1 as [UserID]
            return
        end

declare @LogID int
select @LogID = max(ID) from LoginSessionLog where SessionID=@LoginSessionID

declare @Username nvarchar(255), @FirstName nvarchar(255), @LastName nvarchar(255), @email nvarchar(255)
select @FirstName	= det.QuestionResponse from WebFormResponseDetail det join WebFormQuestions q on det.QuestionID=q.ID and q.InputFieldName='FirstName' and det.HeaderID=@FormHeaderID
select @LastName	= det.QuestionResponse from WebFormResponseDetail det join WebFormQuestions q on det.QuestionID=q.ID and q.InputFieldName='LastName' and det.HeaderID=@FormHeaderID
select @email		= det.QuestionResponse from WebFormResponseDetail det join WebFormQuestions q on det.QuestionID=q.ID and q.InputFieldName='Email' and det.HeaderID=@FormHeaderID

select @Username = coalesce(l.Username,@email,@FirstName+' '+@LastName,'') from LoginSessionLog l where l.ID=@LogID
    if @Username = ''
        begin
            select 'Error: Invalid Username' as [Result], -1 as [UserID]
            return
        end


declare @outputTable table (ID int)
------------------------------------------------------------------------------------------------------------------------
declare @userID int
------------------------------------------------------------------------------------------------------------------------

    if exists(select 1 from [Users] where Username = @Username)
        begin
            -- user already exists! - update the fields

            select @userID = ID from [Users] where [Username] = @Username

            update u set
                         [AuthUsername] = coalesce(nullif(l.Username,''),[AuthUsername],''),
                         [Nickname] = coalesce(l.Nickname,u.[Nickname],''),
                         [Picture] = coalesce(l.Picture,u.[Picture],''),
                         [AuthUserId] = coalesce(l.UserId,u.[AuthUserId],''),
                         [Email] = coalesce(@email,l.Email,u.[Email],''),
                         [EmailVerified] = coalesce(l.EmailVerified,u.[EmailVerified],''),
                         [GivenName] = coalesce(@FirstName,l.GivenName,u.[GivenName],''),
                         [FamilyName] = coalesce(@LastName,l.FamilyName,u.[FamilyName],'')
            from [Users] u
                     join LoginSessionLog l on l.ID=@LogID
            where u.ID = @userID

            if @@ROWCOUNT > 0 insert into AuditTrail(TableName,DataID,EventDescription,SessionID,UserID,NewValue) select 'Users',@userID,'Updated existing User',@LoginSessionID,@userID,'Username: '+isnull(@Username,'')+', Email: '+isnull(@email,'')

        end else begin

        -- add new user
        insert into [Users](Username,AuthUsername,Nickname,Picture,AuthUserId,Email,EmailVerified,GivenName,FamilyName)
        output inserted.ID into @outputTable
        select @Username as [Username],
               coalesce(l.Username,'') as [AuthUsername],
               coalesce(l.Nickname,'') as [Nickname],
               coalesce(l.Picture,'') as [Picture],
               coalesce(l.UserId,'') as [AuthUserId],
               coalesce(@email,l.Email,'') as [Email],
               coalesce(l.EmailVerified,'') as [EmailVerified],
               coalesce(@FirstName,l.GivenName,'') as [GivenName],
               coalesce(@LastName,l.FamilyName,'') as [FamilyName]
        from LoginSessionLog l
        where l.ID=@LogID

        select @userID = ID from @outputTable
        insert into AuditTrail(TableName,DataID,EventDescription,SessionID,UserID,NewValue) select 'Users',@userID,'Added new User',@LoginSessionID,@userID,'Username: '+isnull(@Username,'')+', Email: '+isnull(@email,'')

    end

    --ok, now that the  User record has been updated/inserted, time to update/insert the Member record

------------------------------------------------------------------------------------------------------------------------
declare @memberID int
------------------------------------------------------------------------------------------------------------------------
select @memberID = MemberID from MemberUserLogin where UserID = @userID

    drop table if exists #mbr
select det.HeaderID
     ,max(case when q.EntityName='FirstName' then det.QuestionResponse else '' end) as [FirstName]
     ,max(case when q.EntityName='LastName' then det.QuestionResponse else '' end) as [LastName]
     ,max(case when q.EntityName='Address' then det.QuestionResponse else '' end) as [Address]
     ,max(case when q.EntityName='Postcode' then det.QuestionResponse else '' end) as [Postcode]
     ,max(case when q.EntityName='PreferredPhone' then det.QuestionResponse else '' end) as [PreferredPhone]
     ,max(case when q.EntityName='Email' then det.QuestionResponse else '' end) as [Email]
     ,max(case when q.EntityName='EmergencyContact' then det.QuestionResponse else '' end) as [EmergencyContact]
     ,max(case when q.EntityName='Occupation' then det.QuestionResponse else '' end) as [Occupation]
     ,max(case when q.EntityName='NAWMembershipNumber' then det.QuestionResponse else '' end) as [NAWMembershipNumber]
     ,max(case when q.EntityName='AllowSharingOfMembershipDetails' then left(det.QuestionResponse,1) else '' end) as [AllowSharingOfMembershipDetails]
into #mbr
from WebFormResponseDetail det
         join WebFormQuestions q on det.QuestionID = q.ID
where det.HeaderID = @FormHeaderID
group by det.HeaderID

    if exists (select 1 from [Members] where ID = @memberID)
        begin
            --update existing member

            update m set FirstName=mbr.FirstName, LastName=mbr.LastName, Email=mbr.Email, [Address]=mbr.[Address],
                         PreferredPhone=mbr.PreferredPhone, EmergencyContact=mbr.EmergencyContact, Occupation=mbr.Occupation,
                         NAWMembershipNumber=mbr.NAWMembershipNumber, AllowSharingOfMembershipDetails=mbr.AllowSharingOfMembershipDetails
            from [Members] m
                     cross join #mbr mbr
            where m.ID = @memberID

            if @@ROWCOUNT > 0 insert into AuditTrail(TableName,DataID,EventDescription,SessionID,UserID,MemberID,NewValue) select 'Members',@memberID,'Updated Existing Member',@LoginSessionID,@userID,@memberID,'Username: '+isnull(@Username,'')+', Email: '+isnull(@email,'')

        end else begin

        -- add new member

        insert into [Members]( FirstName, LastName, Email, [Address], PreferredPhone, EmergencyContact, Occupation, NAWMembershipNumber, AllowSharingOfMembershipDetails)
        output inserted.ID into @outputTable
        select FirstName, LastName, Email, [Address],  PreferredPhone, EmergencyContact, Occupation, NAWMembershipNumber, AllowSharingOfMembershipDetails
        from #mbr

        select @memberID = ID from @outputTable
        insert into MemberUserLogin(MemberID,UserID) select @memberID, @userID
        insert into AuditTrail(TableName,DataID,EventDescription,SessionID,UserID,MemberID,NewValue) select 'Members',@memberID,'Added new Member',@LoginSessionID,@userID,@memberID,'Username: '+isnull(@Username,'')+', Email: '+isnull(@email,'')

        --misc details
        update [Members] set YearOfJoining = datepart(year,getdate()) where ID = @memberID
    end

    ------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
select 'OK' as [Result], @userID as [UserID], @memberID as [MemberID]
    return
go
-----------------------------------------------------------------------------------------------------------------------
