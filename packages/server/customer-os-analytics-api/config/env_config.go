package config

type Config struct {
	GraphQL struct {
		PlaygroundEnabled bool `env:"GRAPHQL_PLAYGROUND_ENABLED" envDefault:"false"`
	}
	Db struct {
		Host            string `env:"DB_HOST,required"`
		Port            string `env:"DB_PORT" envDefault:"5432"`
		Pwd             string `env:"DB_PWD,required,unset"`
		Name            string `env:"DB_NAME,required"`
		User            string `env:"DB_USER,required"`
		MaxConn         int    `env:"DB_MAX_CONN"`
		MaxIdleConn     int    `env:"DB_MAX_IDLE_CONN"`
		ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME"`
	}
}
