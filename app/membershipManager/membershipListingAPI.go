package membershipManager

import (
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/model/Logins"
	"ClubMembershipPortal/model/Members"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type UsersAndMembersJSONData struct {
	Result string
	Data   []Members.UserMemberListing `json:"data"`
}

func UsersAndMembersJSON(appCtx *appContextConfig.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		data := UsersAndMembersJSONData{}

		login, membershipStatus, userID, memberID, roles := Logins.LogVisit(ctx, appCtx)
		appCtx.LogInfo(fmt.Sprint("Invoked membershipManager.ListAllUsersAndMembers for : ", login.Username, ", membershipStatus: ", membershipStatus, ", userID: ", userID, ", memberID: ", memberID))
		if !strings.Contains(roles, "Membership Manager") {
			appCtx.LogInfo("ListAllUsersAndMembers() : Error: User: " + login.Username + " does not have the 'Membership Manager' role")
			data.Result = "API Access Denied."
			ctx.JSON(http.StatusBadRequest, data)
			return
		}

		userMembers, err := Members.GetUserMemberListing(appCtx.DBconn)
		if err != nil {
			appCtx.LogInfo("Members.GetUserMemberListing() returned error: " + err.Error())
			data.Result = "An unexpected server error occurred"
			ctx.JSON(http.StatusBadRequest, data)
			return
		}

		data.Result = "OK"
		data.Data = userMembers
		ctx.JSON(http.StatusOK, data)
	}
}
