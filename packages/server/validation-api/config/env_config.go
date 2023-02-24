package config

type Config struct {
	ApiPort string `env:"PORT"`

	ReacherApiPath string `env:"REACHER_API_PATH"`

	Postgres struct {
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
	Neo4j struct {
		Target                          string `env:"NEO4J_TARGET,required"`
		User                            string `env:"NEO4J_AUTH_USER,required,unset"`
		Pwd                             string `env:"NEO4J_AUTH_PWD,required,unset"`
		Realm                           string `env:"NEO4J_AUTH_REALM"`
		MaxConnectionPoolSize           int    `env:"NEO4J_MAX_CONN_POOL_SIZE" envDefault:"100"`
		ConnectionAcquisitionTimeoutSec int    `env:"NEO4J_CONN_ACQUISITION_TIMEOUT_SEC" envDefault:"60"`
		LogLevel                        string `env:"NEO4J_LOG_LEVEL" envDefault:"WARNING"`
	}

	Smarty struct {
		AuthId    string `env:"SMARTY_AUTH_ID,required"`
		AuthToken string `env:"SMARTY_AUTH_TOKEN,required"`
	}
}
