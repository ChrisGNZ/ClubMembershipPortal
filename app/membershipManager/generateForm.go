package membershipManager

import (
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/model/Forms"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func RenderForm(appCtx *appContextConfig.Application, formName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.LogInfo("Invoked RenderForm for : " + formName)
		html, err := Forms.GenerateHTML(appCtx, formName)
		if err != nil {
			appCtx.LogInfo("Error calling Forms.GenerateHTML(" + formName + ") : " + err.Error())
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred.")
			return
		}
		appCtx.LogInfo(fmt.Sprint("html size: ", len(html)))
		ctx.HTML(http.StatusOK, "membershipApplicationForm.html", gin.H{"formHTML": template.HTML(html), "recaptchaSiteKey": appCtx.Config.RecaptchaSiteKey})
	}
}
