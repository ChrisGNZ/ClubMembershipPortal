SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO
drop table if exists [WebForms]
go
CREATE TABLE [dbo].[WebForms](
                                 [ID] [int] IDENTITY(1,1) NOT NULL,
                                 [FormName] [varchar](35) NOT NULL,
                                 [FormDescription] [varchar](max) NOT NULL,
                                 [DefaultEmailReportRecipient] [varchar](max) NULL
                                     PRIMARY KEY CLUSTERED ( [ID] ASC  )  ON [PRIMARY],
                                     CONSTRAINT [FormNameUC] UNIQUE NONCLUSTERED ([FormName] ASC)
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[WebForms] ADD  DEFAULT ('') FOR [DefaultEmailReportRecipient]
GO
-----------------------------------------------------------------------------------------------------------------------
CREATE TABLE [dbo].[WebFormQuestionTemplates](
                                                 [QuestionTemplateName] [varchar](35) NOT NULL,
                                                 [QuestionTemplateText] [nvarchar](max) NULL,
                                                 PRIMARY KEY CLUSTERED ( [QuestionTemplateName] ASC  )
) ON [PRIMARY]
GO
-----------------------------------------------------------------------------------------------------------------------
CREATE TABLE [dbo].[WebFormQuestions](
                                         [ID] [int] IDENTITY(1,1) NOT NULL,
                                         [FormID] [int] NOT NULL REFERENCES [dbo].[WebForms] ([ID]),
                                         [InputFieldName] [varchar](35) NOT NULL,
                                         [EntityName] [varchar](75) NOT NULL,
                                         [QuestionText] [nvarchar](max) NOT NULL,
                                         [QuestionType] [varchar](35) NOT NULL,
                                         [TemplateName] [varchar](35) NOT NULL,
                                         [Seq] [int] NULL,
                                         [SubSeq] [int] NULL,
                                         [AnswerRequired] [char](1) NOT NULL DEFAULT('N'),
                                         [TextRows] [int] NULL,
                                         [RadioOption] [nvarchar](max) NULL,
                                         [QuestionLabel] [nvarchar](max) NULL,
                                         [QuestionExtraText] [nvarchar](max) NULL,
                                         [Hidden] [char](1) NOT NULL DEFAULT('N'),
                                         PRIMARY KEY CLUSTERED
                                             (
                                              [ID] ASC
                                                 )
) ON [PRIMARY]
GO
-----------------------------------------------------------------------------------------------------------------------
CREATE TABLE [dbo].[WebFormResponseHeader](
                                              [ID] [int] IDENTITY(1,1) NOT NULL,
                                              [FormID] [int] NOT NULL,
                                              [UserID] [int] NULL,
                                              [Created] [datetime] NOT NULL,
                                              [DateSubmitted] [date] NOT NULL,
                                              [TimeSubmitted] [datetime] NOT NULL,
                                              [ClientIPAddress] [varchar](75) NULL,
                                              [RecaptchaV3Score] [money] NULL,
                                              LoginSessionID nvarchar(255) NULL,
                                              [EmailReportSent] [datetime] NULL,
                                              [EmailedTo] [varchar](100) NULL,
                                              PRIMARY KEY CLUSTERED([ID] ASC)
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[WebFormResponseHeader] ADD  DEFAULT (getdate()) FOR [Created]
GO
ALTER TABLE [dbo].[WebFormResponseHeader] ADD  DEFAULT (CONVERT([date],getdate())) FOR [DateSubmitted]
GO
ALTER TABLE [dbo].[WebFormResponseHeader] ADD  DEFAULT (getdate()) FOR [TimeSubmitted]
GO
ALTER TABLE [dbo].[WebFormResponseHeader]  WITH CHECK ADD FOREIGN KEY([FormID])  REFERENCES [dbo].[WebForms] ([ID])
GO
-----------------------------------------------------------------------------------------------------------------------
CREATE TABLE [dbo].[WebFormResponseDetail](
                                              [ID] [int] IDENTITY(1,1) NOT NULL,
                                              [HeaderID] [int] NOT NULL,
                                              [QuestionID] [int] NOT NULL,
                                              [QuestionResponse] [nvarchar](max) NULL,
                                              PRIMARY KEY CLUSTERED ( [ID] ASC )
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[WebFormResponseDetail] ADD  DEFAULT ('') FOR [QuestionResponse]
GO
ALTER TABLE [dbo].[WebFormResponseDetail]  WITH CHECK ADD FOREIGN KEY([HeaderID]) REFERENCES [dbo].[WebFormResponseHeader] ([ID])
GO
ALTER TABLE [dbo].[WebFormResponseDetail]  WITH CHECK ADD FOREIGN KEY([QuestionID]) REFERENCES [dbo].[WebFormQuestions] ([ID])
GO
-----------------------------------------------------------------------------------------------------------------------
create table AuditTrail(
                           ID int primary key identity(1,1),
                           LogTime datetime default(getdate()),
                           TableName nvarchar(255) default(''),
                           DataID int null,
                           NewValue nvarchar(255) default(''),
                           OldValue nvarchar(255) default(''),
                           EventDescription  nvarchar(255) default(''),
                           SessionID int null,
                           UserID int null,
                           MemberID int null
)
create index idx_audittrail_time on AuditTrail(LogTime)
create index idx_audittrail_timetable on AuditTrail(LogTime,TableName)
-----------------------------------------------------------------------------------------------------------------------
create table LoginSessionLog(
                                ID int primary key identity(1,1),
                                SessionID nvarchar(255) default(''),
                                Username nvarchar(255) default(''),
                                Nickname nvarchar(255) default(''),
                                Picture nvarchar(255) default(''),
                                UserId nvarchar(255) default(''),
                                Email nvarchar(255) default(''),
                                EmailVerified nvarchar(255) default(''),
                                GivenName nvarchar(255) default(''),
                                FamilyName nvarchar(255) default(''),
                                ClientIP nvarchar(255) default(''),
                                SessionStart datetime default(getdate()),
                                SessionLastAccess datetime default(getdate())
)
-----------------------------------------------------------------------------------------------------------------------
create table [Users](
                        ID int primary key identity(1,1),
                        CreateDate datetime default(getdate()),
                        Username nvarchar(255) default(''),
                        AuthUsername nvarchar(255) default(''),
                        Nickname nvarchar(255) default(''),
                        Picture nvarchar(255) default(''),
                        AuthUserId nvarchar(255) default(''),
                        Email nvarchar(255) default(''),
                        EmailVerified nvarchar(255) default(''),
                        GivenName nvarchar(255) default(''),
                        FamilyName nvarchar(255) default('')
)
-----------------------------------------------------------------------------------------------------------------------
create table [Roles](
                        [ID] int primary key identity(1,1),
                        [RoleName] varchar(35) not null,
                        [RoleDescription] varchar(255) not null default('')
)
-----------------------------------------------------------------------------------------------------------------------
create table [UsersRoles](
                             [ID] int primary key identity(1,1),
                             [UserID] int not null references [Users](ID),
                             [RoleID] int not null references [Roles](ID)
)
-----------------------------------------------------------------------------------------------------------------------
create table [dbo].Members(
                              ID int primary key identity(1,1), /* this is an internal database ID rather than the SAWG Membership # , as I don't know what the full story is about SAWG membership #s and other clubs might not have membership numbers */
                              MembershipStatus nvarchar(50) not null default(''), /* status of the member might include:  New Application Submitted, Accepted */
                              ClubMembershipNumber varchar(35) null, /* for clubs that do keep membership numbers for each of their members... can be left blank if not needed */
                              AuthIdentifier nvarchar(255) null,  /* Auth0 uses an email address, but other auth providers might use a different identifier */
                              FirstName nvarchar(50) not null default(''),
                              LastName nvarchar(50) not null default(''),
                              Email nvarchar(100) not null default(''),   /* keep in mind some members, eg married couples, might share the same email, so allow for this */
                              EmailVerified datetime null,
                              [ClubTitle] nvarchar(255) not null default(''), /* eg "President", "Treasurer", "Membership Officer", and so on.. leave blank for normal members.  Can be shown in badge label printing */
                              [Address] nvarchar(255) not null default(''),  /* i don't see any point in having structured address fields, eg street, building, suburb, etc */
                              Postcode nvarchar(12) not null default(''),
                              PreferredPhone nvarchar(35) not null default(''),
                              SecondaryPhone nvarchar(35) not null default(''),
                              EmergencyContact nvarchar(255) not null default(''), /* SAWG form calls this "Spouse/Partner" but this might be better? */
                              Occupation nvarchar(50) not null default(''),
                              Retired char(1) not null default('N'), /* store retired status separately from occupation */
                              NAWMembershipNumber int null,   /* i think it is safe to assume that NAW will always only use a numeric membership number */
                              AllowSharingOfMembershipDetails char(1) not null default('N'),  /* if 'Y' then this member has given us permission to share their details  with other members */
                              YouthMember char(1) not null default('N'), /* If 'Y' then this member is a youth who qualifies for reduced rate */
                              LifeMember char(1) not null default('N'),  /* If 'Y' then this member is a Life Member and exempt from membership fees */
                              YearOfJoining int null,  /* what year first joined the club */
                              CalculatedJoiningFee money null,  /* estimated fee for new  member */
                              SubmittedJoiningFee money null, /* fee for new member as entered on the application form */
                              LastUpdated datetime not null default(getdate())
)
-----------------------------------------------------------------------------------------------------------------------
create table [dbo].MemberUserLogin(
                                      ID int primary key identity(1,1),
                                      MemberID int not null references Members(ID),
                                      UserID int not null references Users(ID)
)
-----------------------------------------------------------------------------------------------------------------------
-----------------------------------------------------------------------------------------------------------------------