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
			if strings.ToLower(q.TemplateName) == strings.ToLower("StdText") || strings.ToLower(q.TemplateName) == strings.ToLower("StdTextArea") {
				state = "std"
				html += stdField(q)
			} else if strings.ToLower(q.TemplateName) == strings.ToLower("RadioButton") {
				state = "radio"
				radioName = q.InputFieldName
				html += startRadio(q)
			}
		} else if state == "std" {
			if strings.ToLower(q.TemplateName) == strings.ToLower("RadioButton") {
				state = "radio"
				radioName = q.InputFieldName
				html += startRadio(q)
			} else {
				html += stdField(q)
			}
		} else if state == "radio" {
			if strings.ToLower(q.TemplateName) == strings.ToLower("RadioButton") { //continue with radio
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
		//appCtx.LogInfo(fmt.Sprint("len(html) = ", len(html)))
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
	if strings.ToLower(q.AnswerRequired) == strings.ToLower("Y") {
		rpl = strings.Replace(rpl, "{{ .required }}", "required", -1)
	} else {
		rpl = strings.Replace(rpl, "{{ .required }}", "", -1)
	}
	rpl = strings.Replace(rpl, "{{ .questionLabel }}", q.QuestionText, -1)
	return rpl
}
