/*
create procedure MemberMatchExisting @FormHeaderID int
as
set nocount on
*/

/*
this sp differs from UserAddUpdateFromLoginSessionAndMembershipAppForm in that the user has declared that they are an existing member of the club
and have given us some basic details.  We will try to match at least 2 of those fields (or is just matching email good enough??)
*/
declare @FormHeaderID int = 72
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

declare @matches table(MinMatchedMemberID int, MaxMatchedMemberID)
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

-- otherwise keep going, this time just match only on email (if that is considered good enough)

