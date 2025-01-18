package config

type SmartyConfig struct {
	AuthId    string `env:"SMARTY_AUTH_ID"`
	AuthToken string `env:"SMARTY_AUTH_TOKEN"`
}
