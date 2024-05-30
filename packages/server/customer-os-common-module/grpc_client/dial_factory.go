package grpc_client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

type DialFactory interface {
	GetEventsProcessingPlatformConn() (*grpc.ClientConn, error)
	Close(conn *grpc.ClientConn)
}

type DialFactoryImpl struct {
	conf                         *config.GrpcClientConfig
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

	var conn *grpc.ClientConn
	var err error

	if dfi.conf.EventsProcessingPlatformCertificate != "" {
		decodedCertificate, err := base64.StdEncoding.DecodeString(dfi.conf.EventsProcessingPlatformCertificate)
		if err != nil {
			log.Fatalf("Failed to decode certificate: %v", err)
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(decodedCertificate) {
			log.Fatalf("Failed to append CA certificate to pool")
		}

		creds := credentials.NewTLS(&tls.Config{
			RootCAs:    certPool,
			ServerName: dfi.conf.EventsProcessingPlatformServername,
		})

		// Dial the gRPC server with TLS
		conn, err = grpc.Dial(
			dfi.conf.EventsProcessingPlatformUrl,
			grpc.WithTransportCredentials(creds),
			grpc.WithUnaryInterceptor(
				interceptor.ApiKeyEnricher(dfi.conf.EventsProcessingPlatformApiKey),
			),
		)

		if err != nil {
			log.Fatalf("Failed to connect to gRPC server: %v", err)
		}
	} else {
		conn, err = grpc.Dial(dfi.conf.EventsProcessingPlatformUrl, grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(
				interceptor.ApiKeyEnricher(dfi.conf.EventsProcessingPlatformApiKey),
			))
		if err != nil {
			log.Fatalf("Failed to connect to gRPC server: %v", err)
		}
	}

	dfi.eventsProcessingPlatformConn = conn
	return conn, err
}

func NewDialFactory(conf *config.GrpcClientConfig) DialFactory {
	dfi := new(DialFactoryImpl)
	dfi.conf = conf
	return *dfi
}
