package config

type Config struct {
	ApiPort       string `env:"PORT"`
	ApiServiceUrl string `env:"SERVICE_URL"`
	MaxFileSizeMB int64  `env:"MAX_FILE_SIZE_MB"`

	AWS struct {
		Region string `env:"AWS_S3_REGION,required"`
		Bucket string `env:"AWS_S3_BUCKET,required"`
	}

	Postgres struct {
		Host            string `env:"POSTGRES_HOST,required"`
		Port            string `env:"POSTGRES_PORT,required"`
		User            string `env:"POSTGRES_USER,required,unset"`
		Db              string `env:"POSTGRES_DB,required"`
		Password        string `env:"POSTGRES_PASSWORD,required,unset"`
		MaxConn         int    `env:"POSTGRES_DB_MAX_CONN"`
		MaxIdleConn     int    `env:"POSTGRES_DB_MAX_IDLE_CONN"`
		ConnMaxLifetime int    `env:"POSTGRES_DB_CONN_MAX_LIFETIME"`
	}

	Service struct {
		CustomerOsAPI                  string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey               string `env:"CUSTOMER_OS_API_KEY,required"`
		FileStoreAPIJwtSecret          string `env:"FILE_STORE_API_JWT_SECRET,required"`
		CloudflareImageUploadAccountId string `env:"CLOUDFLARE_IMAGE_UPLOAD_ACCOUNT_ID" envDefault:""`
		CloudflareImageUploadApiKey    string `env:"CLOUDFLARE_IMAGE_UPLOAD_API_KEY" envDefault:""`
		CloudflareImageUploadSignKey   string `env:"CLOUDFLARE_IMAGE_UPLOAD_SIGN_KEY" envDefault:""`
	}

	Neo4j struct {
		Target                string `env:"NEO4J_TARGET,required"`
		User                  string `env:"NEO4J_AUTH_USER,required,unset"`
		Pwd                   string `env:"NEO4J_AUTH_PWD,required,unset"`
		Realm                 string `env:"NEO4J_AUTH_REALM"`
		MaxConnectionPoolSize int    `env:"NEO4J_MAX_CONN_POOL_SIZE" envDefault:"100"`
		LogLevel              string `env:"NEO4J_LOG_LEVEL" envDefault:"WARNING"`
	}
}
