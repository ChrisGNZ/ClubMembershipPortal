package Forms

import (
	"ClubMembershipPortal/appContextConfig"
	"strings"
)

func GenerateHTML(appCtx *appContextConfig.Application, formName string) (string, error) {
	/* so this is a very simple loop to generate HTML form elements using the form fields read from the db */
	questions, err := GetFormQuestions(appCtx.DBconn, formName)
	if err != nil {
		appCtx.LogInfo("Error calling: Forms.GetFormQuestions(" + formName + "): " + err.Error())
		return "An unexpected server error occurred in Forms.GenerateHTML()", err
	}

	html := ""
	currentRadioButtonFieldName := ""
	currentRadioButtonContainer := ""
	currentRadioButtonChildren := ""

	for _, q := range questions {

		rpl := ""

		//are we in the middle of a radio button?
		if currentRadioButtonChildren != "" && q.TemplateName != "RadioButton" { // no, we no longer in Radio mode
			//load up the individual buttons into the overal button container
			rpl = strings.Replace(currentRadioButtonContainer, "{{ .radioButtonDetail }}", currentRadioButtonChildren, -1)
			html += rpl
			rpl = ""
			currentRadioButtonChildren = ""
		}

		if q.TemplateName == "StdText" || q.TemplateName == "StdTextArea" {
			rpl = strings.Replace(q.Template1, "{{ .questionType }}", q.QuestionType, -1)
			rpl = strings.Replace(rpl, "{{ .InputFieldName }}", q.InputFieldName, -1)
			rpl = strings.Replace(rpl, "{{ .questionID }}", q.InputFieldName, -1)
			rpl = strings.Replace(rpl, "{{ .questionID }}", q.InputFieldName, -1)
			rpl = strings.Replace(rpl, "{{ .questionValue }}", "", -1)
			if q.AnswerRequired == "Y" {
				rpl = strings.Replace(rpl, "{{ .required }}", "required", -1)
			} else {
				rpl = strings.Replace(rpl, "{{ .required }}", "", -1)
			}
			rpl = strings.Replace(rpl, "{{ .questionLabel }}", q.QuestionText, -1)
			html += rpl
			continue // continue to the next item in the FOR loop
		}

		if q.TemplateName == "RadioButton" { //we have either encountered a new radio, or are already in an existing radio button sequence
			if q.InputFieldName != currentRadioButtonFieldName { //must be a new radio
				currentRadioButtonContainer = strings.Replace(q.Template2, "{{ .InputFieldName }}", q.QuestionText, -1)
				currentRadioButtonFieldName = q.InputFieldName
				currentRadioButtonChildren = ""
			}
			radioBtn := q.Template3
			radioBtn = strings.Replace(radioBtn, "{{ .radioButtonName }}", q.InputFieldName, -1)
			radioBtn = strings.Replace(radioBtn, "{{ .radioButtonID }}", q.InputFieldName, -1)
			radioBtn = strings.Replace(radioBtn, "{{ .radioButtonID }}", q.InputFieldName, -1)
			radioBtn = strings.Replace(radioBtn, "{{ .radioButtonValue }}", q.RadioOption, -1)
			currentRadioButtonChildren += radioBtn
		}

	}
	return html, nil
}
