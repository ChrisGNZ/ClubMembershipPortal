package membershipManager

import (
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/model/Forms"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SaveForm(appCtx *appContextConfig.Application, formName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.LogInfo("Invoked SaveForm for : " + formName)

		//Call to ParseForm makes form fields available.
		err := ctx.Request.ParseForm()
		if err != nil {
			// Handle error here via logging and then return
			appCtx.LogInfo(fmt.Sprint("Error calling : ctx.Request.ParseForm()", err.Error()))
			ctx.String(http.StatusBadGateway, "Error parsing form")
			return
		}
		clientIP := ctx.ClientIP()
		appCtx.LogInfo(fmt.Sprint("Client IP = ", clientIP))

		result, err := Forms.CreateWebFormResponseHeader(appCtx.DBconn, formName, clientIP, -1)
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

		ctx.HTML(http.StatusOK, "membershipApplicationFormPage2.html", gin.H{"HeaderID": HeaderID})
	}
}
