package grpc_client

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client/interceptor"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"google.golang.org/grpc"
	"log"
)

type DialFactory interface {
	GetEventsProcessingPlatformConn() (*grpc.ClientConn, error)
	Close(conn *grpc.ClientConn)
}

type DialFactoryImpl struct {
	conf                         *config.Config
	eventsProcessingPlatformConn *grpc.ClientConn
}

func (dfi DialFactoryImpl) Close(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}

func (dfi DialFactoryImpl) GetEventsProcessingPlatformConn() (*grpc.ClientConn, error) {
	if dfi.eventsProcessingPlatformConn != nil {
		return dfi.eventsProcessingPlatformConn, nil
	}
	conn, err := grpc.Dial(dfi.conf.Services.EventsProcessingPlatformUrl, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			interceptor.ApiKeyEnricher(dfi.conf.Services.EventsProcessingPlatformApiKey),
		))

	dfi.eventsProcessingPlatformConn = conn
	return conn, err
}

func NewDialFactory(conf *config.Config) DialFactory {
	dfi := new(DialFactoryImpl)
	dfi.conf = conf
	return *dfi
}
