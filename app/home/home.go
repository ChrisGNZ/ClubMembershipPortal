package home

import (
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/model/Logins"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomeHandler(appCtx *appContextConfig.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.SysLog.Err("Loading home: /")

		/*
			username := model.GetUsernameAndLogUserFromCtx(ctx, appCtx)
			if username == "" {
				appCtx.SysLog.Err("home.Handler() model.GetUsernameAndLogUserFromCtx() returned empty username")
				ctx.HTML(http.StatusOK, "LoginRequired.html", nil)
				return
			}
		*/
		Logins.LogVisit(ctx, appCtx)
		ctx.HTML(http.StatusOK, "home.html", gin.H{"loggedinusername": ""})

	}
}
