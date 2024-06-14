package model

import "time"

type OAuthToken struct {
	AccessToken       string    `json:"accessToken"`
	RefreshToken      string    `json:"refreshToken"`
	ExpiresAt         time.Time `json:"expiresAt"`
	Scope             string    `json:"scope"`
	ProviderAccountId string    `json:"providerAccountId"`
	IdToken           string    `json:"idToken"`
}

type SignInRequest struct {
	Provider           string     `json:"provider"`
	Tenant             string     `json:"tenant"`
	LoggedInEmail      string     `json:"loggedInEmail"`
	OAuthTokenForEmail string     `json:"oAuthTokenForEmail"`
	OAuthTokenType     string     `json:"oAuthTokenType"`
	OAuthToken         OAuthToken `json:"oAuthToken"`
}

type UpdateUserRequest struct {
	Tenant string `json:"tenant"`
	Email  string `json:"email"`
	UserId string `json:"userId"`
}

type RevokeRequest struct {
	Tenant   string `json:"tenant"`
	Provider string `json:"provider"`
	Email    string `json:"email"`
}
