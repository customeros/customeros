package config

type Neo4jConfig struct {
	Target                          string `env:"NEO4J_TARGET,required"`
	User                            string `env:"NEO4J_AUTH_USER,required,unset"`
	Pwd                             string `env:"NEO4J_AUTH_PWD,required,unset"`
	Realm                           string `env:"NEO4J_AUTH_REALM"`
	MaxConnectionPoolSize           int    `env:"NEO4J_MAX_CONN_POOL_SIZE" envDefault:"100"`
	SocketConnectTimeout            int    `env:"NEO4J_SOCKET_CONN_TIMEOUT_SEC" envDefault:"5"`
	SocketKeepalive                 bool   `env:"NEO4J_SOCKET_KEEP_ALIVE" envDefault:"true"`
	ConnectionAcquisitionTimeoutSec int    `env:"NEO4J_CONN_ACQUISITION_TIMEOUT_SEC" envDefault:"60"`
	LogLevel                        string `env:"NEO4J_LOG_LEVEL" envDefault:"WARNING"`
	Database                        string `env:"NEO4J_DATABASE" envDefault:"neo4j"`
}
