package Members

import (
	"database/sql"
	"errors"
)

// --------------------------------------------------------------------------------------------------------------------
func GetMember(db *sql.DB, memberID int64) (MemberInfo, error) {
	sqlstr := ` exec MemberGet  @MemberID=?  `
	rows, err := db.Query(sqlstr, memberID)

	if err != nil {
		return MemberInfo{}, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	m := MemberInfo{}
	if rows.Next() {
		err = rows.Scan(&m.ID, &m.MembershipStatus, &m.ClubMembershipNumber, &m.AuthIdentifier, &m.FirstName,
			&m.LastName, &m.Email, &m.ClubTitle, &m.Address, &m.Postcode, &m.PreferredPhone, &m.SecondaryPhone,
			&m.EmergencyContact, &m.Occupation, &m.Retired, &m.NAWMembershipNumber, &m.AllowSharingOfMembershipDetails,
			&m.YouthMember, &m.LifeMember, &m.YearOfJoining, &m.CalculatedJoiningFee, &m.SubmittedJoiningFee)
		if err != nil {
			return MemberInfo{}, err
		}
	} else {
		if err != nil {
			return MemberInfo{}, errors.New("Not found")
		}
	}
	return m, nil
}

// --------------------------------------------------------------------------------------------------------------------
func UpdateMembershipStatus(db *sql.DB, userID int64, memberID int64, newStatus string) (string, MemberInfo, error) {
	sqlstr := ` exec MemberUpdateStatus  @MemberID=?, @UserID=?, @Status=?  `
	rows, err := db.Query(sqlstr, memberID, userID, newStatus)

	if err != nil {
		return "Error calling db.Query()", MemberInfo{}, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	result := ""
	if rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return result, MemberInfo{}, err
		}
	}
	memberInfo := MemberInfo{}
	if result == "OK" {
		memberInfo, err = GetMember(db, memberID)
		if err != nil {
			return "error calling GetMember()", MemberInfo{}, err
		}
	}
	return result, memberInfo, nil
}

// --------------------------------------------------------------------------------------------------------------------
func GetUserMemberListing(db *sql.DB) ([]UserMemberListing, error) {

	userMembers := []UserMemberListing{}

	sqlstr := ` exec MemberListing;  `
	rows, err := db.Query(sqlstr)
	if err != nil {
		return userMembers, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	u := UserMemberListing{}
	for rows.Next() {
		err := rows.Scan(&u.Key, &u.UserID, &u.MemberID, &u.Username, &u.AuthUsername, &u.EmailVerified, &u.ClubMembershipNumber, &u.FirstName, &u.LastName, &u.Email, &u.Address, &u.PreferredPhone, &u.EmergencyContact, &u.NAWMembershipNumber, &u.YouthMember, &u.LifeMember, &u.MembershipStatus)
		if err != nil {
			return userMembers, err
		}
		userMembers = append(userMembers, u)
	}

	return userMembers, nil
}

// --------------------------------------------------------------------------------------------------------------------
func ValidateMemberUserSessionLookupKey(db *sql.DB, key string) (string, string, int64, error) {

	sqlstr := ` exec MemberLookupKey @key = ?  `
	rows, err := db.Query(sqlstr, key)
	if err != nil {
		return "Error calling db.Query", "", 0, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	result := "No value returned from db.Query()"
	keyType := ""
	var keyValue int64 = 0
	if rows.Next() {
		err := rows.Scan(&result, &keyType, &keyValue)
		if err != nil {
			return "Error calling rows.Scan", "", 0, err
		}
	}

	return result, keyType, keyValue, nil
}
