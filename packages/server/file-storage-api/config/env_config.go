package config

type Config struct {
	ApiPort       string `env:"PORT"`
	ApiBaseUrl    string `env:"BASE_URL"`
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
}
