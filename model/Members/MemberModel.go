package Members

// ID	ClubMembershipNumber	AuthIdentifier	FirstName	LastName	Email	Address	Postcode	PreferredPhone	SecondaryPhone	EmergencyContact	Occupation	Retired	NAWMembershipNumber	AllowSharingOfMembershipDetails	YouthMember	LifeMember	CalculatedJoiningFee	YearOfJoining	MembershipStatus	SubmittedJoiningFee	ClubTitle
// ---------------------------------------------------------------------------------------------------------------------
type MemberInfo struct {
	ID                              int64
	MembershipStatus                string
	ClubMembershipNumber            string
	AuthIdentifier                  string
	FirstName                       string
	LastName                        string
	Email                           string
	ClubTitle                       string
	Address                         string
	Postcode                        string
	PreferredPhone                  string
	SecondaryPhone                  string
	EmergencyContact                string
	Occupation                      string
	Retired                         string
	NAWMembershipNumber             string
	AllowSharingOfMembershipDetails string
	YouthMember                     string
	LifeMember                      string
	YearOfJoining                   int64
	CalculatedJoiningFee            float64
	SubmittedJoiningFee             float64
}
