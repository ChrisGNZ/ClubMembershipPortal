package Users

import "database/sql"

// --------------------------------------------------------------------------------------------------------------------
func AddUpdateUserFromLoginSessionAndMembershipAppForm(db *sql.DB, LoginSessionID string, FormHeaderID int) (string, int, error) {

	sqlstr := ` exec UserAddUpdateFromLoginSessionAndMembershipAppForm @LoginSessionID=?, @FormHeaderID=? `
	rows, err := db.Query(sqlstr, LoginSessionID, FormHeaderID)
	if err != nil {
		return "", 0, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	result := "Unexpected Error (no response received from database server)"
	var UserID int
	if rows.Next() {
		err = rows.Scan(&result, &UserID)
		if err != nil {
			return "", 0, err
		}
	}
	return result, UserID, nil
}

// --------------------------------------------------------------------------------------------------------------------
func AddUpdateUser(db *sql.DB, Username string, Nickname string, Picture string, AuthUserId string, Email string, EmailVerified string, GivenName string, FamilyName string) (int64, string, string, string, string, string, string, string, string, error) {

	sqlstr := ` exec UserAddOrUpdate  @Username=?, @Nickname=?, @Picture=?, @UserId=?, @Email=?, @EmailVerified=?, @GivenName=?, @FamilyName=?  `
	rows, err := db.Query(sqlstr, Username, Nickname, Picture, AuthUserId, Email, EmailVerified, GivenName, FamilyName)

	if err != nil {
		return 0, "", "", "", "", "", "", "", "", err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	var ID int64
	if rows.Next() {
		err = rows.Scan(&ID, &Username, &Nickname, &Picture, &AuthUserId, &Email, &EmailVerified, &GivenName, &FamilyName)
		if err != nil {
			return 0, "", "", "", "", "", "", "", "", err
		}
	}
	return ID, Username, Nickname, Picture, AuthUserId, Email, EmailVerified, GivenName, FamilyName, nil
}

// --------------------------------------------------------------------------------------------------------------------
