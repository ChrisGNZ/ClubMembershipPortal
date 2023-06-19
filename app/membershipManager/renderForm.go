package membershipManager

import (
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/model/Forms"
	"ClubMembershipPortal/model/Logins"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func RenderForm(appCtx *appContextConfig.Application, formName string, templateName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.LogInfo("Invoked RenderForm for : " + formName)

		login, membershipStatus, userID, memberID, roles := Logins.LogVisit(ctx, appCtx)

		html, err := Forms.GenerateHTML(appCtx, formName)
		if err != nil {
			appCtx.LogInfo("Error calling Forms.GenerateHTML(" + formName + ") : " + err.Error())
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred.")
			return
		}
		appCtx.LogInfo(fmt.Sprint("html size: ", len(html)))
		//ctx.HTML(http.StatusOK, templateName, gin.H{"formHTML": template.HTML(html), "recaptchaSiteKey": appCtx.Config.RecaptchaSiteKey})

		ctx.HTML(http.StatusOK, templateName,
			gin.H{
				"formHTML":         template.HTML(html),
				"recaptchaSiteKey": appCtx.Config.RecaptchaSiteKey,
				"loggedinusername": login.Username,
				"firstname":        login.GivenName,
				"lastname":         login.FamilyName,
				"avatar":           login.Picture,
				"membershipstatus": membershipStatus,
				"userid":           userID,
				"memberid":         memberID,
				"roles":            roles,
			})

	}
}
