package grpc

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"
	"log"
	"net"

	comlog "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/server"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
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
	appLogger := logger.NewExtendedAppLogger(&comlog.Config{
		LogLevel: "debug",
		DevMode:  false,
		Encoder:  "console",
	})
	appLogger.InitLogger()
	appLogger.WithName("event-processing-platform")

	myServer := server.NewServer(&config.Config{
		Utils: config.Utils{
			RetriesOnOptimisticLockException: 3,
		},
	}, appLogger)

	myServer.GrpcServer = grpcServer
	myServer.Repositories = repositories
	myServer.AggregateStore = aggregateStore
	bufferService := eventbuffer.NewEventBufferStoreService(myServer.Repositories.PostgresRepositories.EventBufferRepository, appLogger)
	myServer.CommandHandlers = command.NewCommandHandlers(appLogger, &config.Config{}, aggregateStore, bufferService)
	myServer.Services = service.InitServices(&config.Config{}, repositories, aggregateStore, myServer.CommandHandlers, appLogger, bufferService)

	server.RegisterGrpcServices(myServer.GrpcServer, myServer.Services)

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
