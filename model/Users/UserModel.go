package Users

import (
	"ClubMembershipPortal/appContextConfig"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/*
// ---------------------------------------------------------------------------------------------------------------------
type UserInfo struct {
	ID            int64
	UserName      string
	DisplayName   string
	Email         string
	Phone         string
	DateCreated   string
	UserStatus    string
	LastAccess    string
	IsUserManager string
	IsSysAdmin    string
}
*/
// ---------------------------------------------------------------------------------------------------------------------

func GetUsernameAndLogUserFromCtx(ctx *gin.Context, appCtx *appContextConfig.Application) string {
	session := sessions.Default(ctx)
	profile := session.Get("profile") //https://auth0.com/docs/manage-users/user-accounts/user-profiles/normalized-user-profile-schema

	username := fmt.Sprint(profile.(map[string]interface{})["name"])
	appCtx.SysLog.Info("Loaded Session Profile for Username: " + username)

	return username
}
