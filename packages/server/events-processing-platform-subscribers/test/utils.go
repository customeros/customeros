package test

import (
	"context"
	"testing"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	comlog "github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonServices "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	neo4jt "github.com/customeros/customeros/packages/server/customer-os-common-module/test"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	neo4jtest "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/test"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/logger"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/service"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/test/mocked_grpc"
)

type TestDatabase struct {
	Neo4jContainer    testcontainers.Container
	Driver            *neo4j.DriverWithContext
	postgresContainer testcontainers.Container
	GormDB            *gorm.DB
	CommonServices    *commonServices.CommonServices
	Services          *service.Services
	GrpcClients       *grpc_client.Clients
}

func SetupTestLogger() logger.Logger {
	testLogger := logger.NewExtendedAppLogger(&comlog.Config{
		DevMode: true,
	})
	testLogger.InitLogger()
	return testLogger
}

func SetupTestDatabase() (TestDatabase, func()) {
	testDBs := TestDatabase{}
	testDBs.Neo4jContainer, testDBs.Driver = neo4jtest.InitTestNeo4jDB()
	testDBs.postgresContainer, testDBs.GormDB, _ = neo4jt.InitTestDB()

	rabbitMqContainer, rabbitMqUrl := neo4jt.InitTestRabbitMQ()
	testDialFactory := mocked_grpc.NewMockedTestDialFactory()
	grpcConn, _ := testDialFactory.GetEventsProcessingPlatformConn()
	testDBs.GrpcClients = grpc_client.InitClients(grpcConn)

	// Setup config
	cfg := &commonConfig.CommonConfig{
		Infrastructure: commonConfig.InfrastructureConfig{
			RabbitMQConfig: commonConfig.RabbitMQConfig{
				Url: rabbitMqUrl,
			},
		},
	}

	// Initialize repositories
	neo4jRepositories := neo4j_repository.InitNeo4jRepositories(testDBs.Driver, "neo4j")
	postgresRepositories := &postgres_repository.Repositories{}

	testDBs.CommonServices = commonServices.InitCommonServices(
		SetupTestLogger(),
		neo4jRepositories,
		postgresRepositories,
		cfg,
		testDBs.GrpcClients,
	)

	testDBs.Services = &service.Services{
		CommonServices: testDBs.CommonServices,
	}

	shutdown := func() {
		neo4jtest.CloseDriver(*testDBs.Driver)
		neo4jtest.Terminate(testDBs.Neo4jContainer, context.Background())
		neo4jt.TerminatePostgres(testDBs.postgresContainer, context.Background())
		neo4jt.TerminateRabbitMQ(rabbitMqContainer, context.Background())
	}

	return testDBs, shutdown
}

func SetupMockedTestGrpcClient() *grpc_client.Clients {
	testDialFactory := mocked_grpc.NewMockedTestDialFactory()
	grpcConn, _ := testDialFactory.GetEventsProcessingPlatformConn()
	return grpc_client.InitClients(grpcConn)
}

func AssertRecentTime(t *testing.T, checkTime time.Time) {
	x := 5 // Set the time difference to 5 seconds
	diff := time.Since(checkTime)
	require.True(t, diff <= time.Duration(x)*time.Second, "The time is within the last %d seconds.", x)
}
