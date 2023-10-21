package grpc

import (
	"context"
	common_logger "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	server "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
)

type TestDialFactoryImpl struct {
	eventsProcessingPlatformConn *grpc.ClientConn
}

func (dfi TestDialFactoryImpl) Close(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}

func (dfi TestDialFactoryImpl) GetEventsProcessingPlatformConn(repositories *repository.Repositories, aggregateStore eventstore.AggregateStore) (*grpc.ClientConn, error) {
	if dfi.eventsProcessingPlatformConn != nil {
		return dfi.eventsProcessingPlatformConn, nil
	}
	listener := bufconn.Listen(1024 * 1024)

	grpcServer := grpc.NewServer()
	appLogger := logger.NewExtendedAppLogger(&common_logger.Config{
		LogLevel: "debug",
		DevMode:  false,
		Encoder:  "console",
	})
	appLogger.InitLogger()
	appLogger.WithName("unit-test")

	myServer := server.NewServer(&config.Config{}, appLogger)
	myServer.SetRepository(repositories)
	myServer.SetCommands(commands.InitCommandHandlers(appLogger, &config.Config{}, aggregateStore, repositories))

	server.RegisterGrpcServices(myServer, grpcServer)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return grpc.DialContext(context.Background(), "", grpc.WithInsecure(),
		grpc.WithContextDialer(func(ctx context.Context, target string) (net.Conn, error) {
			return listener.Dial()
		}))
}
