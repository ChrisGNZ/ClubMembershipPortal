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

		/*
			OK, trying to figure out the best way to go here, so lets try something and see if it works...
			so lets get the membershipstatus back from the members table

			1) if it is blank (either a member record doesn't exist or no status set for some reason) then:

				The user needs to either fill out a New Membership Application Form or if they are an existing club
				member who is registering online for the first time, then they need to be able to identify themselves
				as an existing member of the club...

			 	Let's show two buttons:

					[New Membership Application Form]
					(Click here if you would like to join the South Auckland Woodturners Guild)

					[Existing Membership Registration]
					(If you are an existing member of the Guild then click here to register online)

			2) if the membership status is: "Active Member" then that means you are already fully enrolled (if new) or
				already registered (if existing).  Here we can do stuff like show your current subscription status
				and your transaction history

			3) if the membership status is: "Awaiting Processing" then that means your New Membership or Existing registration
				is still being processed.

		*/

		templateName := "home.html"
		switch membershipStatus {
		case "Active Member":
			templateName = "homeActiveMember.html"
		case "Awaiting Processing":
			templateName = "homeAwaitingProcessing.html"
		}

		ctx.HTML(http.StatusOK, templateName,
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
