package events_processing_client

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client/interceptor"
	"google.golang.org/grpc"
)

type DialFactory interface {
	GetEventsProcessingPlatformConn() (*grpc.ClientConn, error)
	Close(conn *grpc.ClientConn)
}

type DialFactoryImpl struct {
	conf                         *config.Config
	eventsProcessingPlatformConn *grpc.ClientConn
	log                          logger.Logger
}

func (dfi DialFactoryImpl) Close(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		dfi.log.Printf("Error closing connection: %v", err)
	}
}

func (dfi DialFactoryImpl) GetEventsProcessingPlatformConn() (*grpc.ClientConn, error) {
	if dfi.eventsProcessingPlatformConn != nil {
		return dfi.eventsProcessingPlatformConn, nil
	}
	// TODO: alexb investigate for required dial options
	conn, err := grpc.Dial(dfi.conf.EventsProcessing.EventsProcessingPlatformUrl, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			interceptor.ApiKeyEnricher(dfi.conf.EventsProcessing.EventsProcessingPlatformApiKey),
		))
	dfi.eventsProcessingPlatformConn = conn
	return conn, err
}

func NewDialFactory(conf *config.Config, log logger.Logger) DialFactory {
	dfi := new(DialFactoryImpl)
	dfi.conf = conf
	dfi.log = log
	return *dfi
}
