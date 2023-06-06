package main

import (
	"ClubMembershipPortal/app/home"
	"ClubMembershipPortal/app/login"
	"ClubMembershipPortal/app/membershipManager"
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/middleware"
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

// New registers the routes and returns the router.
func NewRouter(appCtx *appContextConfig.Application) *gin.Engine {
	router := gin.Default()

	//------------------------------------------------------------------------------------------------------------------
	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	//------------------------------------------------------------------------------------------------------------------
	// Setup the static content folders and load the html templates
	router.Static("/public", "web/static")
	router.StaticFile("/favicon.ico", "web/static/image/favicon.ico")
	router.LoadHTMLGlob("web/template/**/*")

	//------------------------------------------------------------------------------------------------------------------
	// Setup the auth handlers (currently configure for use with Auth0)
	router.GET("/login", login.Handler(appCtx))
	router.GET("/callback", login.Callback(appCtx))
	//router.GET("/logout", logout.Handler(appCtx))

	//------------------------------------------------------------------------------------------------------------------
	// HOME and application
	//------------------------------------------------------------------------------------------------------------------
	router.GET("/", middleware.IsAuthenticated, home.HomeHandler(appCtx))

	router.GET("/newmembershipappform", middleware.IsAuthenticated, membershipManager.RenderForm(appCtx, "Membership Application Form"))
	router.POST("/savemembershipappform", middleware.IsAuthenticated, membershipManager.SaveForm(appCtx, "Membership Application Form"))
	router.POST("/savemembershipappform3", middleware.IsAuthenticated, membershipManager.SaveFormFee(appCtx, "Membership Application Form"))

	//------------------------------------------------------------------------------------------------------------------
	// troubleshooting
	router.NoRoute(func(c *gin.Context) {
		appCtx.LogInfo(fmt.Sprint("404 for: ", c.Request.RequestURI))
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "404 page not found"})
	})

	router.NoMethod(func(c *gin.Context) {
		appCtx.LogInfo(fmt.Sprint("405 for: ", c.Request.RequestURI))
		c.JSON(http.StatusMethodNotAllowed, gin.H{"code": "METHOD_NOT_ALLOWED", "message": "405 method not allowed"})
	})

	//------------------------------------------------------------------------------------------------------------------
	return router
}
