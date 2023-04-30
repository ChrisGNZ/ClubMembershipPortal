package main

import (
	"ClubMembershipPortal/appContextConfig"
	"ClubMembershipPortal/lib/auth"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	appConfig, err := appContextConfig.LoadEnvFile()
	if err != nil {
		log.Fatalf("Failed to load the env variables: %v", err)
	}

	auth, err := auth.New(appConfig.AuthConfig)
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	appCtx, status, err := appContextConfig.CreateApplicationContext(appConfig, auth)
	if err != nil {
		fmt.Println(status)
		fmt.Println(err)
		os.Exit(501)
	}

	rtr := NewRouter(appCtx)

	log.Print("Server listening on http://localhost:" + appCtx.Config.HttpServerPort)
	if err := http.ListenAndServe("0.0.0.0:"+appCtx.Config.HttpServerPort, rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
