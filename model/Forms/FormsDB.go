package Forms

import (
	"database/sql"
)

// --------------------------------------------------------------------------------------------------------------------
func GetFormQuestions(db *sql.DB, FormName string) ([]WebFormQuestion, error) {
	questions := []WebFormQuestion{}

	sqlstr := `exec [FormsGetQuestions] @FormName=?`
	rows, err := db.Query(sqlstr, FormName)
	if err != nil {
		return questions, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	q := WebFormQuestion{}
	for rows.Next() {
		err = rows.Scan(&q.ID, &q.FormID, &q.InputFieldName, &q.EntityName, &q.QuestionText, &q.QuestionType,
			&q.TemplateName, &q.Seq, &q.SubSeq, &q.AnswerRequired, &q.TextRows, &q.RadioOption,
			&q.Template1, &q.Template2, &q.Template3, &q.QuestionLabel, &q.QuestionExtraText)
		if err != nil {
			return questions, err
		}
		questions = append(questions, q)
	}
	return questions, nil
}

// --------------------------------------------------------------------------------------------------------------------
func CreateWebFormResponseHeader(db *sql.DB, FormName string, ClientIPAddress string, RecaptchaV3Score float64, LoginSessionID string) (Result, error) {

	result := Result{}
	sqlstr := `exec [FormsCreateResponseHeader] @FormName=?, @ClientIPAddress=?, @RecaptchaV3Score=?, @LoginSessionID=?; `
	rows, err := db.Query(sqlstr, FormName, ClientIPAddress, RecaptchaV3Score, LoginSessionID)
	if err != nil {
		return result, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if rows.Next() {
		err = rows.Scan(&result.ResultMessage, &result.Id)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// --------------------------------------------------------------------------------------------------------------------
func SaveWebFormResponseDetail(db *sql.DB, HeaderID int, InputFieldName string, QuestionResponse string) (Result, error) {

	result := Result{}
	sqlstr := `exec [FormsSaveResponseDetail] @HeaderID=?, @InputFieldName=?, @QuestionResponse=?`
	rows, err := db.Query(sqlstr, HeaderID, InputFieldName, QuestionResponse)
	if err != nil {
		return result, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if rows.Next() {
		err = rows.Scan(&result.ResultMessage, &result.Id)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// --------------------------------------------------------------------------------------------------------------------
func GetResponseHeader(db *sql.DB, HeaderID int) (WebFormHeader, error) {
	frmHdr := WebFormHeader{}
	sqlstr := `exec [FormsGetResponseHeader] @HeaderID=? `
	rows, err := db.Query(sqlstr, HeaderID)
	if err != nil {
		return frmHdr, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if rows.Next() {
		err = rows.Scan(&frmHdr.ID, &frmHdr.FormID,
			&frmHdr.FormName, &frmHdr.FormDescription,
			&frmHdr.DateSubmitted, &frmHdr.TimeSubmitted,
			&frmHdr.ClientIPAddress, &frmHdr.RecaptchaV3Score,
			&frmHdr.LoginSessionID,
			&frmHdr.DefaultEmailReportRecipient,
			&frmHdr.EmailReportSent, &frmHdr.EmailedTo)
		if err != nil {
			return frmHdr, err
		}
	}
	return frmHdr, nil
}

// --------------------------------------------------------------------------------------------------------------------
func MatchExistingMembership(db *sql.DB, HeaderID int) (string, int64, error) {
	sqlstr := `exec [MemberMatchExisting] @FormHeaderID=? `
	rows, err := db.Query(sqlstr, HeaderID)
	if err != nil {
		return "Error calling db.Query()", 0, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	matchStatus := ""
	var memberID int64

	if rows.Next() {
		err = rows.Scan(&matchStatus, &memberID)
		if err != nil {
			return "Error calling rows.Scan()", 0, err
		}
	}
	return matchStatus, memberID, nil
}

// --------------------------------------------------------------------------------------------------------------------
func GetResponseDetails(db *sql.DB, HeaderID int) ([]WebFormDetail, error) {

	details := []WebFormDetail{}
	sqlstr := `exec [FormsGetResponseDetails] @HeaderID=? `
	rows, err := db.Query(sqlstr, HeaderID)
	if err != nil {
		return details, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	d := WebFormDetail{}
	for rows.Next() {
		err = rows.Scan(&d.HeaderID, &d.QuestionID, &d.Seq,
			&d.EntityName, &d.QuestionText, &d.QuestionResponse)
		if err != nil {
			return details, err
		}
		details = append(details, d)
	}
	return details, nil
}

// --------------------------------------------------------------------------------------------------------------------
// if the option to email new forms to a person such as a membership official is set, mark the email as sent
/*
func MarkEmailSuccess(db *sql.DB, HeaderID int, EmailedTo string) error {
	sqlstr := `exec MarkEmailSuccess @HeaderID=?, @EmailedTo=? `
	rows, err := db.Query(sqlstr, HeaderID, EmailedTo)
	if err != nil {
		return err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	status := ""
	for rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			return err
		}
	}
	if status == "OK" {
		return nil
	} else {
		return errors.New("An unexpected server error occurred")
	}
}
*/
