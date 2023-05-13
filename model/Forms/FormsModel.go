package Forms

type WebFormQuestion struct {
	ID             int
	FormID         int
	InputFieldName string
	EntityName     string
	QuestionText   string
	QuestionType   string
	TemplateName   string
	Seq            int
	SubSeq         int
	AnswerRequired string
	TextRows       int
	RadioOption    string
	Template1      string
	Template2      string
	Template3      string
}

// --------------------------------------------------------------------------------------------------------------------

type WebFormHeader struct {
	ID                          int
	FormID                      int
	FormName                    string
	FormDescription             string
	DateSubmitted               string
	TimeSubmitted               string
	ClientIPAddress             string
	RecaptchaV3Score            string
	DefaultEmailReportRecipient string
	EmailReportSent             string
	EmailedTo                   string
}

type WebFormDetail struct {
	HeaderID         int
	QuestionID       int
	Seq              int
	EntityName       string
	QuestionText     string
	QuestionResponse string
}

type Result struct {
	ResultMessage string
	Id            int
}

// --------------------------------------------------------------------------------------------------------------------
