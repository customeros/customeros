package config

type AppConfig struct {
	Environment string `env:"ENVIRONMENT,required" envDefault:"dev"`
}

type NATSConfig struct {
	Node1 string `env:"NATS_NODE_1_URL,required"`
	Node2 string `env:"NATS_NODE_2_URL"`
	Node3 string `env:"NATS_NODE_3_URL"`
}

type InboxDatabaseConfig struct {
	Host            string `env:"MAILSTACK_POSTGRES_HOST,required"`
	Port            string `env:"MAILSTACK_POSTGRES_PORT,required"`
	User            string `env:"MAILSTACK_POSTGRES_USER,required"`
	DBName          string `env:"MAILSTACK_POSTGRES_DB_NAME,required"`
	Password        string `env:"MAILSTACK_POSTGRES_PASSWORD,required"`
	MaxConn         int    `env:"MAILSTACK_POSTGRES_DB_MAX_CONN"`
	MaxIdleConn     int    `env:"MAILSTACK_POSTGRES_DB_MAX_IDLE_CONN"`
	ConnMaxLifetime int    `env:"MAILSTACK_POSTGRES_DB_CONN_MAX_LIFETIME"`
	LogLevel        string `env:"MAILSTACK_POSTGRES_LOG_LEVEL" envDefault:"WARN"`
}

type DataWarehouseConfig struct {
	Host            string `env:"WAREHOUSE_DB_HOST,required"`
	Port            string `env:"WAREHOUSE_DB_PORT,required"`
	User            string `env:"WAREHOUSE_DB_USER,required"`
	DBName          string `env:"WAREHOUSE_DB_NAME,required"`
	Password        string `env:"WAREHOUSE_DB_PASSWORD,required"`
	MaxConn         int    `env:"WAREHOUSE_DB_MAX_CONN"`
	MaxIdleConn     int    `env:"WAREHOUSE_DB_MAX_IDLE_CONN"`
	ConnMaxLifetime int    `env:"WAREHOUSE_DB_CONN_MAX_LIFETIME"`
	LogLevel        string `env:"WAREHOUSE_DB_LOG_LEVEL" envDefault:"WARN"`
}
