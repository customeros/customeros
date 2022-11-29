package config

type Config struct {
	Neo4jDb struct {
		Target                string `env:"NEO4J_TARGET,required"`
		User                  string `env:"NEO4J_AUTH_USER,required,unset"`
		Pwd                   string `env:"NEO4J_AUTH_PWD,required,unset"`
		Realm                 string `env:"NEO4J_AUTH_REALM"`
		MaxConnectionPoolSize int    `env:"NEO4J_MAX_CONN_POOL_SIZE" envDefault:"100"`
		LogLevel              string `env:"NEO4J_LOG_LEVEL" envDefault:"WARNING"`
	}
	PostgresDb struct {
		Host            string `env:"DB_HOST,required"`
		Port            int    `env:"DB_PORT,required"`
		Pwd             string `env:"DB_PWD,required,unset"`
		Name            string `env:"DB_NAME,required"`
		User            string `env:"DB_USER,required"`
		MaxConn         int    `env:"DB_MAX_CONN"`
		MaxIdleConn     int    `env:"DB_MAX_IDLE_CONN"`
		ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME"`
	}
	TimeoutAfterTaskRun int `env:"TIMEOUT_AFTER_TASK_RUN_SEC" envDefault:"60"`
}
