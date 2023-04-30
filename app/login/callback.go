package login

import (
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/model"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handler for our callback.
func Callback(appCtx *appContextConfig.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			appCtx.SysLog.Err("callback.Handler() : Invalid state parameter.")
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		appCtx.SysLog.Info("callback.Handler() : code = " + ctx.Query("code"))
		token, err := appCtx.Auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			appCtx.SysLog.Err("callback.Handler() : Failed to convert an authorization code into a token.")
			ctx.String(http.StatusUnauthorized, "Failed to convert an authorization code into a token.")
			return
		}

		idToken, err := appCtx.Auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			appCtx.SysLog.Err("callback.Handler() : Failed to verify ID Token.")
			ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			appCtx.SysLog.Err("callback.Handler() : idToken.Claims(&profile) returned error: " + err.Error())
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		if err := session.Save(); err != nil {
			appCtx.SysLog.Err("callback.Handler() : session.Save() returned error: " + err.Error())
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		username := model.GetUsernameAndLogUserFromCtx(ctx, appCtx)
		// Redirect to logged in page.
		appCtx.SysLog.Info(fmt.Sprint("callback.Handler() : Success! state = ", ctx.Query("state"),
			", code = ", ctx.Query("code"), ", token.AccessToken = ", token.AccessToken,
			", session username = ", username))
		ctx.Redirect(http.StatusTemporaryRedirect, "/user")
	}
}
