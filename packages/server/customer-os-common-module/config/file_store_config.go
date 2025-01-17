package config

type FileStoreConfig struct {
	FileStoreJwtSecret             string `env:"FILE_STORE_JWT_SECRET"`
	CloudflareImageUploadAccountId string `env:"CLOUDFLARE_IMAGE_UPLOAD_ACCOUNT_ID" envDefault:""`
	CloudflareImageUploadApiKey    string `env:"CLOUDFLARE_IMAGE_UPLOAD_API_KEY" envDefault:""`
	CloudflareImageUploadSignKey   string `env:"CLOUDFLARE_IMAGE_UPLOAD_SIGN_KEY" envDefault:""`
	MaxFileSizeMB                  int64  `env:"MAX_FILE_SIZE_MB"`
	AWS                            struct {
		Region string `env:"AWS_S3_REGION"`
		Bucket string `env:"AWS_S3_BUCKET"`
	}
}
