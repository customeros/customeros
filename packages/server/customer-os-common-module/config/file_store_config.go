package config

type FileStoreConfig struct {
	FileStoreJwtSecret             string `env:"FILE_STORE_JWT_SECRET"`
	CloudflareImageUploadAccountId string `env:"CLOUDFLARE_IMAGE_UPLOAD_ACCOUNT_ID"`
	CloudflareImageUploadApiKey    string `env:"CLOUDFLARE_IMAGE_UPLOAD_API_KEY"`
	CloudflareImageUploadSignKey   string `env:"CLOUDFLARE_IMAGE_UPLOAD_SIGN_KEY"`
	MaxFileSizeMB                  int64  `env:"MAX_FILE_SIZE_MB"`
	AWS                            AwsConfig
}
