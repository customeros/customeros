package config

type Config struct {
	DB struct {
		Host     string `env:"DB_HOST,required"`
		Port     int    `env:"DB_PORT,required"`
		User     string `env:"DB_USER,required,unset"`
		Name     string `env:"DB_NAME,required"`
		Password string `env:"DB_PASSWORD,required,unset"`
	}
	Neo4jDb struct {
		Target                string `env:"NEO4J_TARGET,required"`
		User                  string `env:"NEO4J_AUTH_USER,required,unset"`
		Pwd                   string `env:"NEO4J_AUTH_PWD,required,unset"`
		Realm                 string `env:"NEO4J_AUTH_REALM"`
		MaxConnectionPoolSize int    `env:"NEO4J_MAX_CONN_POOL_SIZE" envDefault:"100"`
		LogLevel              string `env:"NEO4J_LOG_LEVEL" envDefault:"WARNING"`
	}
	Service struct {
		ServerPort       int    `env:"MESSAGE_STORE_SERVER_PORT,required"`
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
	}
	Identity struct {
		DefaultUserId string `env:"DEFAULT_USER_ID" envDefault:"AgentSmith"`
	}
}
