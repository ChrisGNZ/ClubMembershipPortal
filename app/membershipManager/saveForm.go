package membershipManager

import (
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/model/Forms"
	"ClubMembershipPortal/model/Logins"
	"ClubMembershipPortal/model/Users"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SaveForm(appCtx *appContextConfig.Application, formName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.LogInfo("Invoked SaveForm for : " + formName)
		login := Logins.LogVisit(ctx, appCtx)
		appCtx.LogInfo("Login Info : SessionID = " + login.SessionID + ", Username = " + login.Username + ", Email = " + login.Email + ", Name: " + login.GivenName + " " + login.FamilyName)
		//displayName := Users.AddUpdateUser(appCtx.DBconn, login.Username, login.Nickname, login.Picture, login.UserId, login.Email, login.EmailVerified, login.GivenName, login.FamilyName)

		//Call to ParseForm makes form fields available.
		err := ctx.Request.ParseForm()
		if err != nil {
			// Handle error here via logging and then return
			appCtx.LogInfo(fmt.Sprint("Error calling : ctx.Request.ParseForm()", err.Error()))
			ctx.String(http.StatusBadGateway, "Error parsing form")
			return
		}

		result, err := Forms.CreateWebFormResponseHeader(appCtx.DBconn, formName, login.ClientIP, -1, login.SessionID)
		if err != nil {
			appCtx.LogInfo(fmt.Sprint("CreateWebFormResponseHeader() error: ", result.ResultMessage, err))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}
		HeaderID := result.Id

		if HeaderID == 0 {
			appCtx.LogInfo(fmt.Sprint("CreateWebFormResponseHeader returned HeaderID = zero, cannot continue."))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}

		for key, value := range ctx.Request.PostForm {

			val := ""
			if len(value) > 0 {
				val = value[0]
			}

			appCtx.LogInfo(fmt.Sprint("calling SaveWebFormResponseDetail(", HeaderID, key, val))
			result, err := Forms.SaveWebFormResponseDetail(appCtx.DBconn, HeaderID, key, val)
			if err != nil {
				appCtx.LogInfo(fmt.Sprint("SaveWebFormResponseDetail() returned error: ", result.ResultMessage, err))
				ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
				return
			}
		}

		//now link login to user, if not already done so
		appCtx.LogInfo(fmt.Sprint("About to call AddUpdateUserFromLoginSessionAndMembershipAppForm(SessionID = ", login.SessionID, ", HeaderID = ", HeaderID))
		linkresult, userId, err := Users.AddUpdateUserFromLoginSessionAndMembershipAppForm(appCtx.DBconn, login.SessionID, HeaderID)
		if err != nil || linkresult != "OK" {
			appCtx.LogInfo(fmt.Sprint("AddUpdateUserFromLoginSessionAndMembershipAppForm() returned error: ", result, " ", err))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}
		ctx.HTML(http.StatusOK, "membershipApplicationForm2.html", gin.H{"HeaderID": HeaderID, "UserID": userId})
	}
}
