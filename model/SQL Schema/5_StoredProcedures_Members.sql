go
create or alter procedure UserCalculateNewMembershipFee @UserID int   -- returns nawMembershipStatus, fullOrHalfYearStatus, calculatedFee
as
    set nocount on
declare @nawMembershipStatus varchar(255), @fullOrHalfYearStatus varchar(255), @calculatedFee money, @isFullYear char(1)

    if not exists(select 1 from [Users] where ID = @UserID)
        begin
            select '' as [isFullYear],'Invalid User ID!  User ID '+isnull(convert(varchar,@UserID),'NULL')+' not found.' as [nawMembershipStatus], 'ERROR' as [fullOrHalfYearStatus], 0.00 as [calculatedFee]
            return
        end

declare @joiningDate date
select @joiningDate = CreateDate from [Users] where ID = @UserID
    if @joiningDate is null
        begin

            select '' as [isFullYear], 'Invalid User Record. Invalid CreateDate field for User ID '+isnull(convert(varchar,@UserID),'NULL')+'.' as [nawMembershipStatus], 'ERROR' as [fullOrHalfYearStatus], 0.00 as [calculatedFee]
            return
        end

    if not exists(select 1 from [Members] m join [MemberUserLogin] mul on mul.MemberID = m.ID join [Users] u on u.ID=mul.UserID where u.ID = @UserID)
        begin
            select '' as [isFullYear], 'Missing Member record!!  User ID '+isnull(convert(varchar,@UserID),'NULL')+' is not linked to a Member record' as [nawMembershipStatus], 'ERROR' as [fullOrHalfYearStatus], 0.00 as [calculatedFee]
            return
        end

    ---------------------------------------------------------------------------------------------------------------
    ---------------------------------------------------------------------------------------------------------------
    --if not exists(select 1 from [Users] where
    --first determene if this is a full year or half year subscription
    --if new member joins after October 1 and before March 31 then this is a half year subscription
    -- otherwise, between April 1 and August 31 means a full year subscription
    ---------------------------------------------------------------------------------------------------------------
    ---------------------------------------------------------------------------------------------------------------
declare @year char(4), @FullYearMarkBegin date, @FullYearMarkEnd date
select @year = datepart(year,@joiningDate), @FullYearMarkBegin = convert(date,@year+'0401',112), @FullYearMarkEnd = convert(date,@year+'0831',112)  --these hardcoded MMDD values should probably go into some sort of lookup table...

    if @joiningDate between @FullYearMarkBegin and @FullYearMarkEnd
        begin
            select @fullOrHalfYearStatus = 'This is for a Full-Year Membership', @calculatedFee = 80.00   --this hardcoded value should probably go into some sort of lookup table...
            select @isFullYear = 'Y'
        end else begin
        select @fullOrHalfYearStatus = 'This is for a Half-Year Membership', @calculatedFee = 40.00   --this hardcoded value should probably go into some sort of lookup table...
        select @isFullYear = ''
    end

declare @NAWnumber int
select @NAWnumber = isnull(NAWMembershipNumber,0) from [Members] m join [MemberUserLogin] mul on mul.MemberID = m.ID join [Users] u on u.ID=mul.UserID where u.ID = @UserID

    if @NAWnumber = 0
        begin
            select @nawMembershipStatus = 'You are not a member of NAW.'
            select @calculatedFee = @calculatedFee + 12.00  --this hardcoded value should probably go into some sort of lookup table...
        end else begin
        select @nawMembershipStatus = 'You are a member of NAW (# '+convert(varchar,@NAWnumber)+').'
    end
    ---------------------------------------------------------------------------------------------------------------
insert into AuditTrail(TableName,DataID,EventDescription, UserID,MemberID,NewValue)
select 'Members',m.ID,'UserCalculateNewMembershipFee',u.ID,m.ID, 'Caculated fee: '+convert(varchar,@calculatedFee)
from [Members] m join [MemberUserLogin] mul on mul.MemberID = m.ID join [Users] u on u.ID=mul.UserID where u.ID = @UserID

select @isFullYear as [isFullYear], @nawMembershipStatus as [nawMembershipStatus], @fullOrHalfYearStatus as[fullOrHalfYearStatus], @calculatedFee as [calculatedFee]
go
-----------------------------------------------------------------------------------------------------------------------
go
create or alter procedure MemberGet @MemberID int
as
    set nocount on
select ID,  MembershipStatus,
       isnull(ClubMembershipNumber,'') as [ClubMembershipNumber],
       isnull(AuthIdentifier, '') as [AuthIdentifier],
       FirstName, LastName, Email, ClubTitle,[Address],
       Postcode, PreferredPhone, SecondaryPhone,
       EmergencyContact, Occupation, Retired,
       isnull(NAWMembershipNumber,0) as [NAWMembershipNumber],
       AllowSharingOfMembershipDetails,
       YouthMember, LifeMember,
       isnull(YearOfJoining,0) as [YearOfJoining],
       isnull(CalculatedJoiningFee,0) as [CalculatedJoiningFee],
       isnull(SubmittedJoiningFee,0) as [SubmittedJoiningFee]
from [Members]
where ID = @MemberID
go
-----------------------------------------------------------------------------------------------------------------------
-----------------------------------------------------------------------------------------------------------------------
-----------------------------------------------------------------------------------------------------------------------


create or alter procedure MemberMatchExisting @FormHeaderID int
as
    set nocount on

    /*
    this sp differs from UserAddUpdateFromLoginSessionAndMembershipAppForm in that the user has declared that they are an existing member of the club
    and have given us some basic details.  We will try to match at least 2 of those fields (or is just matching email good enough??)
    */
-- test declare @FormHeaderID int = 72
---------------------------------------------------------------------------------------------------------------------------------------------
declare @userID int, @BadgeNumber varchar(255), @BadgeName nvarchar(255), @Email nvarchar(255), @PreferredPhone varchar(35)

select @userID			= UserID from WebFormResponseHeader where ID = @FormHeaderID
select @BadgeNumber		= isnull(d.QuestionResponse,'') from WebFormResponseDetail d join WebFormQuestions q on d.QuestionID=q.ID where  d.HeaderID = @FormHeaderID and q.EntityName='BadgeNumber'
select @BadgeName		= isnull(d.QuestionResponse,'') from WebFormResponseDetail d join WebFormQuestions q on d.QuestionID=q.ID where  d.HeaderID = @FormHeaderID and q.EntityName='BadgeName'
select @Email			= isnull(d.QuestionResponse,'') from WebFormResponseDetail d join WebFormQuestions q on d.QuestionID=q.ID where  d.HeaderID = @FormHeaderID and q.EntityName='Email'
select @PreferredPhone	= isnull(d.QuestionResponse,'') from WebFormResponseDetail d join WebFormQuestions q on d.QuestionID=q.ID where  d.HeaderID = @FormHeaderID and q.EntityName='PreferredPhone'

select @PreferredPhone = dbo.StripNonNumeric(@PreferredPhone)

-- (housekeeping) if the user table does not yet have an email and we were given one, then update the user table
update Users set Email = @Email where ID = @userID and isnull(Email,'') = '' and @Email != ''

---------------------------------------------------------------------------------------------------------------------------------------------

declare @matches table(MatchedMemberID int )
declare @matchesCount int
insert into @matches(MatchedMemberID) select m.ID from Members m where m.Email = @Email and dbo.StripNonNumeric(PreferredPhone) = @PreferredPhone
select @matchesCount = count(*) from @matches

    if @matchesCount = 1
        begin
            -- we have exactly one match on both email and phone # : this is a very strong match, so return a valid result
            select 'OK' as [MatchStatus], MatchedMemberID as [MemberID] from @matches
            return
        end

    if @matchesCount > 1
        begin
            -- sometimes some clubs allow family memberships where multiple members of a family can share the same email and phone - in this case manual matching is needed.
            select 'Manual Matching Required as more than one matching email and phone found' as [MatchStatus], 0 as [MemberID]
            return
        end
    ---------------------------------------------------------------------------------------------------------------------------------------------
-- otherwise keep going, this time just match only on email (if that is considered good enough)
insert into @matches(MatchedMemberID) select m.ID from Members m where m.Email = @Email
select @matchesCount = count(*) from @matches
    if @matchesCount = 1
        begin
            -- we have exactly one match on email, and this might be good enough...
            select 'OK' as [MatchStatus], MatchedMemberID as [MemberID] from @matches
            return
        end

    if @matchesCount > 1
        begin
            -- sometimes some clubs allow family memberships where multiple members of a family can share the same email  - in this case manual matching is needed.
            select 'Manual Matching Required as more than one matching email found' as [MatchStatus], 0 as [MemberID]
            return
        end
    ---------------------------------------------------------------------------------------------------------------------------------------------
-- otherwise lets give up at this point
-- maybe future TO DO:
-- try matching on badge name and number (this would require some sort of fuzzy matching for the name)
select 'Manual Matching Required' as [MatchStatus], 0 as [MemberID]
    return
go
grant execute on MemberMatchExisting to TestPortalUser
