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
