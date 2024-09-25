package eventstoredb

type EventStoreConfig struct {
	ConnectionString   string `env:"EVENT_STORE_CONNECTION_STRING" validate:"required"`
	TlsDisable         bool   `env:"EVENT_STORE_CONNECTION_TLS_DISABLED" envDefault:"false"`
	TlsVerifyCert      bool   `env:"EVENT_STORE_CONNECTION_TLS_VERIFY_CERT" envDefault:"true"`
	KeepAliveTimeout   int    `env:"EVENT_STORE_CONNECTION_KEEP_ALIVE_TIMEOUT_MS" envDefault:"10000"`
	KeepAliveInterval  int    `env:"EVENT_STORE_CONNECTION_KEEP_ALIVE_INTERVAL_MS" envDefault:"10000"`
	ConnectionUser     string `env:"EVENT_STORE_CONNECTION_USERNAME"`
	ConnectionPassword string `env:"EVENT_STORE_CONNECTION_PASSWORD,unset"`
}
