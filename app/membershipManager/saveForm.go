package membershipManager

import (
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/model/Forms"
	"ClubMembershipPortal/model/Logins"
	"ClubMembershipPortal/model/Members"
	"ClubMembershipPortal/model/Users"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func SaveForm(appCtx *appContextConfig.Application, formName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.LogInfo("Invoked SaveForm for : " + formName)
		login, membershipStatus, userID, memberID, roles := Logins.LogVisit(ctx, appCtx)
		appCtx.LogInfo("Login Info : SessionID = " + login.SessionID + ", Username = " + login.Username +
			", Email = " + login.Email + ", Name: " + login.GivenName + " " + login.FamilyName + " " +
			membershipStatus + " " + strconv.FormatInt(userID, 10) + " " + strconv.FormatInt(memberID, 10) +
			", Roles = " + roles)

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
		linkresult, userId, memberID, err := Users.AddUpdateUserFromLoginSessionAndMembershipAppForm(appCtx.DBconn, login.SessionID, HeaderID)
		if err != nil || linkresult != "OK" {
			appCtx.LogInfo(fmt.Sprint("AddUpdateUserFromLoginSessionAndMembershipAppForm() returned error: ", result, " ", err))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}

		isFullYear, nawMembershipStatus, fullOrHalfYearStatus, calculatedFee, err := Users.CalculateNewMembershipFee(appCtx.DBconn, userId)
		ctx.HTML(http.StatusOK, "membershipApplicationForm2.html",
			gin.H{
				"HeaderID":             HeaderID,
				"UserID":               userId,
				"MemberID":             memberID,
				"isFullYear":           isFullYear,
				"nawMembershipStatus":  nawMembershipStatus,  //{{ .nawMembershipStatus }}You are a member of NAW (# 2905)
				"fullOrHalfYearStatus": fullOrHalfYearStatus, //{{ .fullOrHalfYearStatus }}This is for a Full Year
				"calculatedFee":        calculatedFee,        //{{ .calculatedFee }}
			})
	}
}

// --------------------------------------------------------------------------------------------------------------------
func SaveFormFee(appCtx *appContextConfig.Application, formName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.LogInfo("Invoked SaveForm for : " + formName)
		login, membershipStatus, userID, memberID, roles := Logins.LogVisit(ctx, appCtx)
		appCtx.LogInfo("Login Info : SessionID = " + login.SessionID + ", Username = " + login.Username +
			", Email = " + login.Email + ", Name: " + login.GivenName + " " + login.FamilyName +
			", membershipStatus = " + membershipStatus + ", userID = " + strconv.FormatInt(userID, 10) +
			", memberID = " + strconv.FormatInt(memberID, 10) +
			", roles = " + roles)

		//Call to ParseForm makes form fields available.
		err := ctx.Request.ParseForm()
		if err != nil {
			// Handle error here via logging and then return
			appCtx.LogInfo(fmt.Sprint("Error calling : ctx.Request.ParseForm()", err.Error()))
			ctx.String(http.StatusBadGateway, "Error parsing form")
			return
		}

		headerIDstr := ctx.Request.FormValue("formheader")
		if headerIDstr == "" {
			appCtx.LogInfo(fmt.Sprint("2nd part of New Membership Form did not return a header ID"))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}

		var HeaderID int
		id, err := strconv.ParseInt(headerIDstr, 10, 64)
		if err != nil {
			appCtx.LogInfo(fmt.Sprint("2nd part of New Membership Form did not return a numeric header ID"))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}
		HeaderID = int(id)

		hdr, err := Forms.GetResponseHeader(appCtx.DBconn, HeaderID)
		if err != nil || hdr.ID == 0 {
			appCtx.LogInfo(fmt.Sprint("2nd part of New Membership Form did not return a valid header ID"))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}

		userIDstr := ctx.Request.FormValue("userid")
		if userIDstr == "" {
			appCtx.LogInfo(fmt.Sprint("2nd part of New Membership Form did not return a User ID"))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}
		var UserID int
		id, err = strconv.ParseInt(userIDstr, 10, 64)
		if err != nil {
			appCtx.LogInfo(fmt.Sprint("2nd part of New Membership Form did not return a numeric User ID"))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}
		UserID = int(id)

		memberIDstr := ctx.Request.FormValue("memberid")
		if memberIDstr == "" {
			appCtx.LogInfo(fmt.Sprint("2nd part of New Membership Form did not return a Member ID"))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}

		MemberID, err := strconv.ParseInt(memberIDstr, 10, 64)
		if err != nil {
			appCtx.LogInfo(fmt.Sprint("2nd part of New Membership Form did not return a numeric Member ID"))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}

		youth := ctx.Request.FormValue("youth")
		fee := ctx.Request.FormValue("fee")
		calculatedFee := ctx.Request.FormValue("calculated")

		result, err := Forms.SaveWebFormResponseDetail(appCtx.DBconn, HeaderID, "YouthMember", youth)
		if err != nil {
			appCtx.LogInfo(fmt.Sprint("SaveWebFormResponseDetail() returned error: ", result.ResultMessage, err))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}
		result, err = Forms.SaveWebFormResponseDetail(appCtx.DBconn, HeaderID, "SubmittedFee", fee)
		if err != nil {
			appCtx.LogInfo(fmt.Sprint("SaveWebFormResponseDetail() returned error: ", result.ResultMessage, err))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}

		member, err := Members.GetMember(appCtx.DBconn, MemberID)
		if err != nil {
			appCtx.LogInfo(fmt.Sprint("Members.GetMember() returned error: ", err))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}

		ctx.HTML(http.StatusOK, "membershipApplicationForm3.html",
			gin.H{
				"HeaderID":          HeaderID,
				"UserID":            UserID,
				"MemberID":          MemberID,
				"Youth":             youth,         //Youth Membership appplied for?
				"SubmittedFee":      fee,           //user submitted Fee
				"calculatedFee":     calculatedFee, //algorithmically calculated Fee
				"email":             member.Email,
				"phone":             member.PreferredPhone,
				"membershipOfficer": "memsec@sawg.org.nz",
			})
	}
}

// --------------------------------------------------------------------------------------------------------------------
func RegisterExistingMembership(appCtx *appContextConfig.Application, formName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.LogInfo("Invoked SaveForm for : " + formName)
		login, membershipStatus, userID, memberID, roles := Logins.LogVisit(ctx, appCtx)
		appCtx.LogInfo("Login Info : SessionID = " + login.SessionID + ", Username = " + login.Username +
			", Email = " + login.Email + ", Name: " + login.GivenName + " " + login.FamilyName + " " +
			membershipStatus + " " + strconv.FormatInt(userID, 10) + " " + strconv.FormatInt(memberID, 10) +
			", Roles = " + roles)

		err := ctx.Request.ParseForm()
		if err != nil {
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
				appCtx.LogInfo(fmt.Sprint("Forms.SaveWebFormResponseDetail() returned error: ", result.ResultMessage, err))
				ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
				return
			}
		}

		//now we need the logic to try and either automagically register a user who has an exact match on email address and a close? match on name?
		matchStatus, matchedMembershipID, err := Forms.MatchExistingMembership(appCtx.DBconn, HeaderID)
		if err != nil {
			appCtx.LogInfo(fmt.Sprint("Forms.MatchExistingMembership() returned error: ", result.ResultMessage, err))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}

		/*
			if matchedMembershipID is not zero then that means we have successfully matched this user with an existing membership record!
			--> now we need to update the linkage between user & member then display a logged-in member page

			else if matchedMembershipID is zero then manual matching is required
			1) add this user to the list of members waiting for manual matching
			2) display a message to the user advising them that they will be contacted shortly
			3) do we need to validate their email?  maybe we should
		*/
		if matchedMembershipID != 0 {
			//matched to a member!  So now we can update that Member's status to be "Active"
			// update [Members].membershipstatus and [MemberUserLogin] tables
			updateResult, updatedMemberInfo, err := Members.UpdateMembershipStatus(appCtx.DBconn, userID, memberID, "Active")
			if updateResult != "OK" || err != nil {
				appCtx.LogInfo(fmt.Sprint("Forms.MatchExistingMembership() error: ", updateResult, ", err: ", err))
				ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
				return
			}
			ctx.HTML(http.StatusOK, "homeActiveMember.html",
				gin.H{
					"loggedinusername": login.Username,
					"firstname":        login.GivenName,
					"lastname":         login.FamilyName,
					"avatar":           login.Picture,
					"membershipstatus": updatedMemberInfo.MembershipStatus,
					"userid":           userID,
					"memberid":         memberID,
					"roles":            roles,
				})
			return
		}

		//otherwise, we are NOT matched to a member!
		// we need to
		//     1) add this user to the list of members waiting for manual matching
		//     2) display a message to the user advising them that they will be contacted shortly
		//     3) do we need to validate their email?  maybe we should
		ctx.String(http.StatusOK, fmt.Sprint("matchStatus: ", matchStatus, ", Membership ID: ", matchedMembershipID))
		/*
			appCtx.LogInfo(fmt.Sprint("About to call AddUpdateUserFromLoginSessionAndMembershipAppForm(SessionID = ", login.SessionID, ", HeaderID = ", HeaderID))
			linkresult, userId, memberID, err := Users.AddUpdateUserFromLoginSessionAndMembershipAppForm(appCtx.DBconn, login.SessionID, HeaderID)
			if err != nil || linkresult != "OK" {
				appCtx.LogInfo(fmt.Sprint("AddUpdateUserFromLoginSessionAndMembershipAppForm() returned error: ", result, " ", err))
				ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
				return
			}

			isFullYear, nawMembershipStatus, fullOrHalfYearStatus, calculatedFee, err := Users.CalculateNewMembershipFee(appCtx.DBconn, userId)
			ctx.HTML(http.StatusOK, "membershipApplicationForm2.html",
				gin.H{
					"HeaderID":             HeaderID,
					"UserID":               userId,
					"MemberID":             memberID,
					"isFullYear":           isFullYear,
					"nawMembershipStatus":  nawMembershipStatus,  //{{ .nawMembershipStatus }}You are a member of NAW (# 2905)
					"fullOrHalfYearStatus": fullOrHalfYearStatus, //{{ .fullOrHalfYearStatus }}This is for a Full Year
					"calculatedFee":        calculatedFee,        //{{ .calculatedFee }}
				})

		*/
	}
}
