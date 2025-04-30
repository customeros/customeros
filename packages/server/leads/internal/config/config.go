package config

type AppConfig struct {
	APIPort     string `env:"API_PORT,required" envDefault:"8080"`
	APIKey      string `env:"API_KEY,required"`
	Environment string `env:"ENVIRONMENT,required" envDefault:"dev"`
}

type NATSConfig struct {
	Node1 string `env:"NATS_NODE_1_URL,required"`
	Node2 string `env:"NATS_NODE_2_URL"`
	Node3 string `env:"NATS_NODE_3_URL"`
}

type LeadsDatabaseConfig struct {
	Host            string `env:"LEADS_POSTGRES_HOST,required"`
	WritePort       string `env:"LEADS_POSTGRES_WRITE_PORT,required"`
	ReadPort        string `env:"LEADS_POSTGRES_READ_PORT,required"`
	User            string `env:"LEADS_POSTGRES_USER,required"`
	DBName          string `env:"LEADS_POSTGRES_DB_NAME,required"`
	Password        string `env:"LEADS_POSTGRES_PASSWORD,required"`
	MaxConn         int    `env:"LEADS_POSTGRES_DB_MAX_CONN"`
	MaxIdleConn     int    `env:"LEADS_POSTGRES_DB_MAX_IDLE_CONN"`
	ConnMaxLifetime int    `env:"LEADS_POSTGRES_DB_CONN_MAX_LIFETIME"`
	LogLevel        string `env:"LEADS_POSTGRES_LOG_LEVEL" envDefault:"WARN"`
}

type DataWarehouseConfig struct {
	Host            string `env:"WAREHOUSE_DB_HOST,required"`
	WritePort       string `env:"WAREHOUSE_DB_WRITE_PORT,required"`
	ReadPort        string `env:"WAREHOUSE_DB_READ_PORT,required"`
	User            string `env:"WAREHOUSE_DB_USER,required"`
	DBName          string `env:"WAREHOUSE_DB_NAME,required"`
	Password        string `env:"WAREHOUSE_DB_PASSWORD,required"`
	MaxConn         int    `env:"WAREHOUSE_DB_MAX_CONN"`
	MaxIdleConn     int    `env:"WAREHOUSE_DB_MAX_IDLE_CONN"`
	ConnMaxLifetime int    `env:"WAREHOUSE_DB_CONN_MAX_LIFETIME"`
	LogLevel        string `env:"WAREHOUSE_DB_LOG_LEVEL" envDefault:"WARN"`
}

type IPDataConfig struct {
	ApiUrl             string `env:"IPDATA_API_URL"`
	ApiKey             string `env:"IPDATA_API_KEY"`
	IpDataCacheTtlDays int    `env:"IPDATA_CACHE_TTL_DAYS" envDefault:"90"`
}

type SnitcherConfig struct {
	Url    string `env:"SNITCHER_API_URL" required:"true" envDefault:"https://app.snitcher.com/api"`
	ApiKey string `env:"SNITCHER_API_KEY" `
}
