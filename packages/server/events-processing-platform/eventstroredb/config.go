package eventstroredb

type EventStoreConfig struct {
	ConnectionString string `env:"EVENT_STORE_CONNECTION_STRING" validate:"required"`
}
