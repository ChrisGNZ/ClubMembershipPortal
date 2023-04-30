package login

import (
	"ClubMembershipPortal/appContextConfig"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func Logout(appCtx *appContextConfig.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logoutUrl, err := url.Parse("https://" + appCtx.Config.AuthConfig.AUTH0_DOMAIN + "/v2/logout")
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		returnTo, err := url.Parse(appCtx.Config.AuthConfig.AUTH0_LOGGEDOUT_URL)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		parameters := url.Values{}
		parameters.Add("returnTo", returnTo.String())
		parameters.Add("client_id", appCtx.Config.AuthConfig.AUTH0_CLIENT_ID)
		logoutUrl.RawQuery = parameters.Encode()

		ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
	}
}
