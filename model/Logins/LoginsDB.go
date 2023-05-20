package Logins

import "database/sql"

// --------------------------------------------------------------------------------------------------------------------
func logSession(db *sql.DB, login LoginStruct) (string, error) {

	result := "An unexpected server error occurred"

	sqlstr := ` exec LoginsLogSession @SessionID=?, @Username=?, @Nickname=?, @Picture=?, @UserId=?, @Email=?, @EmailVerified=?, @GivenName=?, @FamilyName=?, @ClientIP=?; `
	rows, err := db.Query(sqlstr, login.SessionID, login.Username, login.Nickname, login.Picture, login.UserId, login.Email, login.EmailVerified, login.GivenName, login.FamilyName, login.ClientIP)

	if err != nil {
		return result, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// --------------------------------------------------------------------------------------------------------------------
