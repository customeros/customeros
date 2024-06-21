package config

type GoogleOAuthConfig struct {
	ClientId     string `env:"GOOGLE_OAUTH_CLIENT_ID" envDefault:""`
	ClientSecret string `env:"GOOGLE_OAUTH_CLIENT_SECRET" envDefault:""`
}
