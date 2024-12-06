package config

type PostgresAsyncConfig struct {
	Host            string `env:"POSTGRES_ASYNC_HOST,required"`
	Port            string `env:"POSTGRES_ASYNC_PORT,required"`
	User            string `env:"POSTGRES_ASYNC_USER,required,unset"`
	Db              string `env:"POSTGRES_ASYNC_DB,required"`
	Password        string `env:"POSTGRES_ASYNC_PASSWORD,required,unset"`
	MaxConn         int    `env:"POSTGRES_ASYNC_DB_MAX_CONN"`
	MaxIdleConn     int    `env:"POSTGRES_ASYNC_DB_MAX_IDLE_CONN"`
	ConnMaxLifetime int    `env:"POSTGRES_ASYNC_DB_CONN_MAX_LIFETIME"`
	LogLevel        string `env:"POSTGRES_LOG_LEVEL" envDefault:"WARN"`
}
