package membershipManager

import (
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/model/Forms"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func RenderForm(appCtx *appContextConfig.Application, formName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		html, err := Forms.GenerateHTML(appCtx, formName)
		if err != nil {
			appCtx.LogInfo("Error calling Forms.GenerateHTML(" + formName + ") : " + err.Error())
			ctx.String(http.StatusBadGateway)
		}
		ctx.HTML(http.StatusOK, "ClientReferralForm.html", gin.H{"formHTML": template.HTML(html), "recaptchaSiteKey": appCtx.Config.RecaptchaSiteKey})
	}
}
