package config

type AzureOAuthConfig struct {
	ClientId     string `env:"AZURE_OAUTH_CLIENT_ID" envDefault:""`
	ClientSecret string `env:"AZURE_OAUTH_CLIENT_SECRET" envDefault:""`
}
