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
		linkresult, userId, memberID, err := Users.AddUpdateUserFromLoginSessionAndMembershipAppForm(appCtx.DBconn, login.SessionID, HeaderID)
		if err != nil || linkresult != "OK" {
			appCtx.LogInfo(fmt.Sprint("AddUpdateUserFromLoginSessionAndMembershipAppForm() returned error: ", result, " ", err))
			ctx.String(http.StatusBadGateway, "An unexpected server error occurred")
			return
		}

		nawMembershipStatus, fullOrHalfYearStatus, calculatedFee, err := Users.CalculateNewMembershipFee(appCtx.DBconn, userId)
		ctx.HTML(http.StatusOK, "membershipApplicationForm2.html",
			gin.H{
				"HeaderID":             HeaderID,
				"UserID":               userId,
				"MemberID":             memberID,
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
		login := Logins.LogVisit(ctx, appCtx)
		appCtx.LogInfo("Login Info : SessionID = " + login.SessionID + ", Username = " + login.Username + ", Email = " + login.Email + ", Name: " + login.GivenName + " " + login.FamilyName)

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
				"membershipOfficer": "membership@sawg.org.nz",
			})
	}
}
