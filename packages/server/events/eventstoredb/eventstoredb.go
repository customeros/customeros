package eventstoredb

import (
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"time"
)

func NewEventStoreDB(cfg EventStoreConfig, log logger.Logger) (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString(cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	(*settings).KeepAliveTimeout = time.Millisecond * time.Duration(cfg.KeepAliveTimeout)
	(*settings).KeepAliveInterval = time.Millisecond * time.Duration(cfg.KeepAliveInterval)
	(*settings).SkipCertificateVerification = !cfg.TlsVerifyCert
	(*settings).DisableTLS = cfg.TlsDisable
	if cfg.ConnectionUser != "" && cfg.ConnectionPassword != "" {
		(*settings).Username = cfg.ConnectionUser
		(*settings).Password = cfg.ConnectionPassword
	}

	log.Infof("Obtaining connection to EventStoreDB...")
	//log.Infof("EventStoreDB connection settings: {%+v}", *settings)

	esdbClient, err := esdb.NewClient(settings)

	serverVersion, _ := esdbClient.GetServerVersion()
	log.Infof("EventStore version details: {%+v}", *serverVersion)

	return esdbClient, err
}
