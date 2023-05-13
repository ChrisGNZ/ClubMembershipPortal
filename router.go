package main

import (
	"ClubMembershipPortal/app/home"
	"ClubMembershipPortal/app/login"
	"ClubMembershipPortal/app/membershipManager"
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/middleware"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
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

	router.GET("/newmembershipappform", middleware.IsAuthenticated, membershipManager.RenderForm(appCtx))
	router.POST("/savemembershipappform", middleware.IsAuthenticated, membershipManager.SaveForm(appCtx))
	return router
}
