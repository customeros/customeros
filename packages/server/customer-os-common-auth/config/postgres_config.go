package config

type PostgresConfig struct {
	Host            string `env:"POSTGRES_HOST,required"`
	Port            string `env:"POSTGRES_PORT,required"`
	User            string `env:"POSTGRES_USER,required,unset"`
	Db              string `env:"POSTGRES_DB,required"`
	Password        string `env:"POSTGRES_PASSWORD,required,unset"`
	MaxConn         int    `env:"POSTGRES_DB_MAX_CONN"`
	MaxIdleConn     int    `env:"POSTGRES_DB_MAX_IDLE_CONN"`
	ConnMaxLifetime int    `env:"POSTGRES_DB_CONN_MAX_LIFETIME"`
	LogLevel        string `env:"POSTGRES_LOG_LEVEL" envDefault:"WARN"`
}
