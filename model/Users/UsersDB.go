package Users

import "database/sql"

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
