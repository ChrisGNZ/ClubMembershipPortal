package membershipManager

import (
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/model/Logins"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ListAllUsersAndMembers(appCtx *appContextConfig.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		login, membershipStatus, userID, memberID, roles := Logins.LogVisit(ctx, appCtx)
		appCtx.LogInfo(fmt.Sprint("Invoked membershipManager.ListAllUsersAndMembers for : ", login.Username, ", membershipStatus: ", membershipStatus, ", userID: ", userID, ", memberID: ", memberID))
		if !strings.Contains(roles, "Membership Manager") {
			appCtx.LogInfo("ListAllUsersAndMembers() : Error: User: " + login.Username + " does not have the 'Membership Manager' role")
			ctx.String(http.StatusBadGateway, "Access Denied.")
			return
		}
		ctx.HTML(http.StatusOK, "membershipListing.html",
			gin.H{
				"loggedinusername": login.Username,
				"firstname":        login.GivenName,
				"lastname":         login.FamilyName,
				"avatar":           login.Picture,
				"membershipstatus": membershipStatus,
				"userid":           userID,
				"memberid":         memberID,
				"roles":            roles,
			})
		//use https://datatables.net/extensions/rowreorder/examples/initialisation/responsive.html
	}
}
