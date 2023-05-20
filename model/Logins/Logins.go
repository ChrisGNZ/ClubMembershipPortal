package Logins

import (
	"ClubMembershipPortal/appContextConfig"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strings"
	"time"
)

func LogVisit(ctx *gin.Context, appCtx *appContextConfig.Application) LoginStruct {
	session := sessions.Default(ctx)
	profile := session.Get("profile") //https://auth0.com/docs/manage-users/user-accounts/user-profiles/normalized-user-profile-schema

	sessionID := fmt.Sprint(session.Get("sessionID"))
	appCtx.LogInfo("sessionID from session.Get = " + sessionID)
	var err error
	if sessionID == "" || sessionID == "<nil>" {
		sessionID = fmt.Sprint(generateRandomSessionID())
		session.Set("sessionID", sessionID)
		if err := session.Save(); err != nil {
			appCtx.LogInfo("Error saving session: " + err.Error())
		}
	}
	login := LoginStruct{}
	login.SessionID = sessionID
	//again, these are all based on what Auth0 can return.  When other auth providers are added, we might add more fields
	login.Username = strings.Replace(fmt.Sprint(profile.(map[string]interface{})["name"]), "<nil>", "", 1)
	login.Nickname = strings.Replace(fmt.Sprint(profile.(map[string]interface{})["nickname"]), "<nil>", "", 1)
	login.Picture = strings.Replace(fmt.Sprint(profile.(map[string]interface{})["picture"]), "<nil>", "", 1)
	login.UserId = strings.Replace(fmt.Sprint(profile.(map[string]interface{})["user_id"]), "<nil>", "", 1)
	login.Email = strings.Replace(fmt.Sprint(profile.(map[string]interface{})["email"]), "<nil>", "", 1)
	login.EmailVerified = strings.Replace(fmt.Sprint(profile.(map[string]interface{})["email_verified"]), "<nil>", "", 1)
	login.GivenName = strings.Replace(fmt.Sprint(profile.(map[string]interface{})["given_name"]), "<nil>", "", 1)
	login.FamilyName = strings.Replace(fmt.Sprint(profile.(map[string]interface{})["family_name"]), "<nil>", "", 1)
	login.ClientIP = ctx.ClientIP()

	result, err := logSession(appCtx.DBconn, login)
	if err != nil {
		appCtx.LogInfo("An error occurred calling Logins.LogSession() : " + result + ", " + err.Error())
		return LoginStruct{}
	}
	appCtx.LogInfo("Logins.LogSession() returned: " + result)
	return login
}

// --------------------------------------------------------------------------------------------------------------------
func generateRandomSessionID() int {

	rand.Seed(time.Now().UnixNano())
	min := 1000000
	max := 9999999
	return rand.Intn(max-min+1) + min

}
