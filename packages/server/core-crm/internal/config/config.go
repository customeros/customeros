package config

type AppConfig struct {
	APIPort     string `env:"API_PORT,required" envDefault:"8080"`
	Environment string `env:"ENVIRONMENT,required" envDefault:"dev"`
}

type NATSConfig struct {
	Node1 string `env:"NATS_NODE_1_URL,required"`
	Node2 string `env:"NATS_NODE_2_URL"`
	Node3 string `env:"NATS_NODE_3_URL"`
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
