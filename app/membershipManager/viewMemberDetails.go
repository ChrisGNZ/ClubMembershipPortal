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

func ViewMemberDetails(appCtx *appContextConfig.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		login, membershipStatus, userID, memberID, roles := Logins.LogVisit(ctx, appCtx)
		appCtx.LogInfo(fmt.Sprint("Invoked membershipManager.ViewMemberDetails for : ", login.Username, ", membershipStatus: ", membershipStatus, ", userID: ", userID, ", memberID: ", memberID))
		if !strings.Contains(roles, "Membership Manager") {
			appCtx.LogInfo("ViewMemberDetails() : Error: User: " + login.Username + " does not have the 'Membership Manager' role")
			ctx.String(http.StatusBadGateway, "Access Denied.")
			return
		}

		keystr := ctx.Query("key")
		//result, keyType, keyValue, err := Members.ValidateMemberUserSessionLookupKey(appCtx.DBconn, keystr)
		result, _, _, err := Members.ValidateMemberUserSessionLookupKey(appCtx.DBconn, keystr)
		if err != nil {
			appCtx.LogInfo("Error calling Members.ValidateMemberUserSessionLookupKey(" + keystr + "): " + result + " : " + err.Error())
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}

		/* TO DO
		details, err := Members.GetMemberUserSessionLookupKeyDetails(appCtx.DBconn, keyType, keyValue)
		if err != nil {
			appCtx.LogInfo(fmt.Sprint("Error calling Members.ValidateMemberUserSessionLookupKeyDetails(", keyType, ", ", keyValue, "): ", err.Error()))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}
		*/
		ctx.String(http.StatusOK, "TO DO: show membership details for key: "+keystr)
		/*
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
		*/

		//use https://datatables.net/extensions/rowreorder/examples/initialisation/responsive.html
	}
}
