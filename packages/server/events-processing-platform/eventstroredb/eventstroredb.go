package eventstroredb

import (
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
)

func NewEventStoreDB(cfg EventStoreConfig) (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString(cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	return esdb.NewClient(settings)
}
