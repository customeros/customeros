package config

type AppConfig struct {
	APIPort           string `env:"PORT,required" envDefault:"12222"`
	TrackingPublicUrl string `env:"TRACKING_PUBLIC_URL" envDefault:"https://custosmetrics.com"`
}

type MailstackDatabaseConfig struct {
	Host            string `env:"POSTGRES_HOST,required"`
	Port            string `env:"POSTGRES_PORT,required"`
	User            string `env:"POSTGRES_USER,required"`
	DBName          string `env:"POSTGRES_DB_NAME,required"`
	Password        string `env:"POSTGRES_PASSWORD,required"`
	MaxConn         int    `env:"POSTGRES_DB_MAX_CONN"`
	MaxIdleConn     int    `env:"POSTGRES_DB_MAX_IDLE_CONN"`
	ConnMaxLifetime int    `env:"POSTGRES_DB_CONN_MAX_LIFETIME"`
	LogLevel        string `env:"POSTGRES_LOG_LEVEL" envDefault:"WARN"`
	SSLMode         string `env:"POSTGRES_SSL_MODE"`
}
