package auth

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"os"
)

/*
user profile received from auth0 when an example user logs in
{
  "sub": "auth0|1234",
  "nickname": "chris",
  "name": "chris@fantail.net.nz",
  "picture": "https://s.gravatar.com/avatar/1234?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fch.png",
  "updated_at": "2022-05-25T06:40:20.048Z"
}
*/

type AuthenticatorConfig struct {
	AUTH0_DOMAIN        string // The URL of our Auth0 Tenant Domain.
	AUTH0_CLIENT_ID     string // Our Auth0 Application's Client ID.
	AUTH0_CLIENT_SECRET string // Our Auth0 Application's Client Secret.
	AUTH0_CALLBACK_URL  string
	AUTH0_LOGGEDOUT_URL string //redirect to here after logout was processed by Auth0
}

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

func LoadAuthConfigFromEnv() AuthenticatorConfig {
	ac := AuthenticatorConfig{}
	ac.AUTH0_DOMAIN = os.Getenv("AUTH0_DOMAIN")
	ac.AUTH0_CLIENT_ID = os.Getenv("AUTH0_CLIENT_ID")
	ac.AUTH0_CLIENT_SECRET = os.Getenv("AUTH0_CLIENT_SECRET")
	ac.AUTH0_CALLBACK_URL = os.Getenv("AUTH0_CALLBACK_URL")
	ac.AUTH0_LOGGEDOUT_URL = os.Getenv("AUTH0_LOGGEDOUT_URL")
	return ac
}

// New instantiates the *Authenticator.
func New(authConfig AuthenticatorConfig) (*Authenticator, error) {

	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+authConfig.AUTH0_DOMAIN+"/",
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     authConfig.AUTH0_CLIENT_ID,
		ClientSecret: authConfig.AUTH0_CLIENT_SECRET,
		RedirectURL:  authConfig.AUTH0_CALLBACK_URL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
