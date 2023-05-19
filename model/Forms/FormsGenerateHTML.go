package Forms

import (
	"ClubMembershipPortal/appContextConfig"
	"fmt"
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
	state := "start"
	radioName := ""
	for _, q := range questions {

		appCtx.LogInfo(fmt.Sprint("InputFieldName: ", q.InputFieldName, ", QuestionType: ", q.QuestionType, ", TemplateName: ", q.TemplateName, ", State: ", state))
		if state == "start" {
			if q.TemplateName == "StdText" || q.TemplateName == "StdTextArea" {
				state = "std"
				html += stdField(q)
			} else if q.TemplateName == "RadioButton" {
				state = "radio"
				radioName = q.InputFieldName
				html += startRadio(q)
			}
		} else if state == "std" {
			if q.TemplateName == "RadioButton" {
				state = "radio"
				radioName = q.InputFieldName
				html += startRadio(q)
			} else {
				html += stdField(q)
			}
		} else if state == "radio" {
			if q.TemplateName == "RadioButton" { //continue with radio
				if radioName != q.InputFieldName {
					html += finishRadio()
					radioName = q.InputFieldName
					html += startRadio(q)
				} else {
					html += continueRadio(q)
				}
			} else {
				html += finishRadio()
				state = "std"
				html += stdField(q)
			}
		}
	} //end for range loop
	if state == "radio" {
		html += finishRadio()
	}

	return html, nil
} //end func

func startRadio(q WebFormQuestion) string {
	rpl := strings.Replace(q.Template2, "{{ .InputFieldName }}", q.QuestionText, -1)
	rpl = strings.Replace(rpl, "{{ .questionExtraText }}", q.QuestionExtraText, -1)
	rpl += continueRadio(q)
	return rpl
}

func continueRadio(q WebFormQuestion) string {
	radioBtn := q.Template3
	radioBtn = strings.Replace(radioBtn, "{{ .radioButtonName }}", q.InputFieldName, -1)
	radioBtn = strings.Replace(radioBtn, "{{ .radioButtonID }}", q.InputFieldName, -1)
	radioBtn = strings.Replace(radioBtn, "{{ .radioButtonID }}", q.InputFieldName, -1)
	radioBtn = strings.Replace(radioBtn, "{{ .radioButtonValue }}", q.RadioOption, -1)
	return radioBtn
}

func finishRadio() string {
	rpl := `		</div>
    </div>
</div>`
	return rpl
}

func stdField(q WebFormQuestion) string {
	rpl := strings.Replace(q.Template1, "{{ .questionType }}", q.QuestionType, -1)
	rpl = strings.Replace(rpl, "{{ .InputFieldName }}", q.InputFieldName, -1)
	rpl = strings.Replace(rpl, "{{ .questionID }}", q.InputFieldName, -1)
	rpl = strings.Replace(rpl, "{{ .questionID }}", q.InputFieldName, -1)
	rpl = strings.Replace(rpl, "{{ .questionExtraText }}", q.QuestionExtraText, -1)
	rpl = strings.Replace(rpl, "{{ .questionValue }}", "", -1)
	if q.AnswerRequired == "Y" {
		rpl = strings.Replace(rpl, "{{ .required }}", "required", -1)
	} else {
		rpl = strings.Replace(rpl, "{{ .required }}", "", -1)
	}
	rpl = strings.Replace(rpl, "{{ .questionLabel }}", q.QuestionText, -1)
	return rpl
}

/*


		currentRadioButtonFieldName := ""
		currentRadioButtonContainer := ""
		currentRadioButtonChildren := ""


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
*/
