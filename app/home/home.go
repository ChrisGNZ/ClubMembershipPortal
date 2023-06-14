package home

import (
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/model/Logins"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func HomeHandler(appCtx *appContextConfig.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.SysLog.Err("Loading home: /")

		login, membershipStatus, userID, memberID, roles := Logins.LogVisit(ctx, appCtx)
		ctx.SetCookie("UserID", strconv.FormatInt(userID, 10), 60*60*24, "/", "", true, true)
		ctx.SetCookie("MemberID", strconv.FormatInt(memberID, 10), 60*60*24, "/", "", true, true)

		ctx.HTML(http.StatusOK, "home.html",
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

	}
}
