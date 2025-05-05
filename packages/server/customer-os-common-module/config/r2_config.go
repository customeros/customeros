package config

type R2StorageConfig struct {
	AccountID       string `env:"CLOUDFLARE_R2_ACCOUNT_ID"`
	AccessKeyID     string `env:"CLOUDFLARE_R2_ACCESS_KEY_ID"`
	AccessKeySecret string `env:"CLOUDFLARE_R2_ACCESS_KEY_SECRET"`
}
