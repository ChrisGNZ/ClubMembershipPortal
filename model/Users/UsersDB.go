package Users

import (
	"database/sql"
	"errors"
)

// --------------------------------------------------------------------------------------------------------------------
func AddUpdateUserFromLoginSessionAndMembershipAppForm(db *sql.DB, LoginSessionID string, FormHeaderID int) (string, int, int, error) {

	sqlstr := ` exec UserAddUpdateFromLoginSessionAndMembershipAppForm @LoginSessionID=?, @FormHeaderID=? `
	rows, err := db.Query(sqlstr, LoginSessionID, FormHeaderID)
	if err != nil {
		return "", 0, 0, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	result := "Unexpected Error (no response received from database server)"
	var UserID int
	var MemberID int
	if rows.Next() {
		err = rows.Scan(&result, &UserID, &MemberID)
		if err != nil {
			return "", 0, 0, err
		}
	}
	return result, UserID, MemberID, nil
}

// --------------------------------------------------------------------------------------------------------------------
func CalculateNewMembershipFee(db *sql.DB, userId int) (string, string, float64, error) {

	sqlstr := ` exec UserCalculateNewMembershipFee @UserID=?  `
	rows, err := db.Query(sqlstr, userId)
	if err != nil {
		return "", "", 0, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	nawMembershipStatus := ""
	fullOrHalfYearStatus := ""
	var calculatedFee float64
	if rows.Next() {
		err = rows.Scan(&nawMembershipStatus, &fullOrHalfYearStatus, &calculatedFee)
		if err != nil {
			return "", "", 0, err
		}
	} else {
		return "", "", 0, errors.New("No rows returned from database")
	}

	return nawMembershipStatus, fullOrHalfYearStatus, calculatedFee, nil
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
