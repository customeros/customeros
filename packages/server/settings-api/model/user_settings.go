package model

type UserSettings struct {
	ID                          string
	TenantName                  string
	UserName                    string
	GoogleOAuthAllScopesEnabled bool
	GoogleOAuthUserAccessToken  string
}
