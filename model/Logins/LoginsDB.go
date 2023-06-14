package Logins

import "database/sql"

// --------------------------------------------------------------------------------------------------------------------
func logSession(db *sql.DB, login LoginStruct) (string, string, int64, int64, string, error) {

	result := "An unexpected server error occurred"

	sqlstr := ` exec LoginsLogSession @SessionID=?, @Username=?, @Nickname=?, @Picture=?, @UserId=?, @Email=?, @EmailVerified=?, @GivenName=?, @FamilyName=?, @ClientIP=?; `
	rows, err := db.Query(sqlstr, login.SessionID, login.Username, login.Nickname, login.Picture, login.UserId, login.Email, login.EmailVerified, login.GivenName, login.FamilyName, login.ClientIP)

	if err != nil {
		return result, "", 0, 0, "", err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	var MembershipStatus string
	var UserID int64
	var MemberID int64
	var Roles string
	if rows.Next() {
		err = rows.Scan(&result, &MembershipStatus, &UserID, &MemberID, &Roles)
		if err != nil {
			return result, "", 0, 0, "", err
		}
	}
	return result, MembershipStatus, UserID, MemberID, Roles, nil
}

// --------------------------------------------------------------------------------------------------------------------
