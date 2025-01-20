package resolver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commontest "github.com/customeros/customeros/packages/server/customer-os-common-module/test"
	neo4jtest "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/test"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/customer-os-api/config"
	"github.com/customeros/customeros/packages/server/customer-os-api/dataloader"
	cosHandler "github.com/customeros/customeros/packages/server/customer-os-api/graphql"
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/generated"
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	cosapiservices "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-api/test/grpc/events_platform"
)

var (
	neo4jContainer testcontainers.Container
	driver         *neo4j.DriverWithContext

	postgresContainer testcontainers.Container
	gormDB            *gorm.DB
	sqlDB             *sql.DB

	rabbitMqContainer testcontainers.Container
	rabbitMqUrl       string

	c                        *client.Client
	cOwner                   *client.Client
	cCustomerOsPlatformOwner *client.Client
	cAdmin                   *client.Client
	cAdminWithTenant         *client.Client

	// services
	customerOsApiServices *cosapiservices.Services
)

const (
	tenantName    = "openline"
	testUserId    = "test-user-id"
	testUserEmail = "test-user-email"
)

func TestMain(m *testing.M) {
	// Start Neo4j container
	neo4jContainer, driver = commontest.InitTestNeo4jDB()
	defer func(dbContainer testcontainers.Container, driver neo4j.DriverWithContext, ctx context.Context) {
		commontest.CloseDriver(driver)
		commontest.TerminateNeo4j(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	// Start Postgres container
	postgresContainer, gormDB, sqlDB = commontest.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		commontest.TerminatePostgres(postgresContainer, ctx)
	}(postgresContainer, context.Background())

	// Start RabbitMQ container
	rabbitMqContainer, rabbitMqUrl = commontest.InitTestRabbitMQ()
	defer func(rabbitMqContainer testcontainers.Container, ctx context.Context) {
		commontest.TerminateRabbitMQ(rabbitMqContainer, ctx)
	}(rabbitMqContainer, context.Background())

	prepareClient()

	os.Exit(m.Run())
}

func tearDownTestCase(ctx context.Context) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %s, cleaning neo4j DB", tb.Name())
		neo4jtest.CleanupAllData(ctx, driver)
	}
}

func prepareClient() {
	appLogger := logger.NewAppLogger(&logger.Config{
		DevMode: true,
	})
	appLogger.InitLogger()

	postgresDB := &commonConfig.PostgresDB{
		GormDB:      gormDB,
		AsyncGormDB: gormDB,
	}

	testDialFactory := events_platform.NewTestDialFactory()
	gRPCconn, _ := testDialFactory.GetEventsProcessingPlatformConn()

	grpcClient := grpc_client.InitClients(gRPCconn)

	customerOsApiServices = cosapiservices.InitServices(
		appLogger,
		driver,
		postgresDB,
		&config.Config{Common: &commonConfig.CommonConfig{
			Infrastructure: commonConfig.InfrastructureConfig{
				RabbitMQConfig: commonConfig.RabbitMQConfig{
					Url: rabbitMqUrl,
				},
			},
		},
		},
		grpcClient,
	)

	graphResolver := NewResolver(
		appLogger,
		customerOsApiServices,
		grpcClient,
		&config.Config{},
	)

	loader := dataloader.NewDataLoader(customerOsApiServices)
	customCtx := &common.CustomContext{
		Tenant:    tenantName,
		UserId:    testUserId,
		UserEmail: testUserEmail,
		Roles:     []string{model.RoleUser.String()},
	}

	customOwnerCtx := &common.CustomContext{
		Tenant:    tenantName,
		UserId:    testUserId,
		UserEmail: testUserEmail,
		Roles:     []string{model.RoleUser.String(), model.RoleOwner.String()},
	}

	customCustomerOsPlatformOwnerCtx := &common.CustomContext{
		Tenant:    tenantName,
		UserId:    testUserId,
		UserEmail: testUserEmail,
		Roles:     []string{model.RoleUser.String(), model.RolePlatformOwner.String()},
	}

	customAdminCtx := &common.CustomContext{
		Roles: []string{model.RoleAdmin.String()},
	}

	customAdminWTenantCtx := &common.CustomContext{
		Tenant: tenantName,
		Roles:  []string{model.RoleAdmin.String()},
	}

	schemaConfig := generated.Config{Resolvers: graphResolver}
	schemaConfig.Directives.HasRole = cosHandler.GetRoleChecker()
	schemaConfig.Directives.HasTenant = cosHandler.GetTenantChecker()
	server := handler.NewDefaultServer(generated.NewExecutableSchema(schemaConfig))
	dataloaderServer := dataloader.Middleware(loader, server)
	handler := common.WithContext(customCtx, dataloaderServer)
	c = client.New(handler)
	cOwner = client.New(common.WithContext(customOwnerCtx, dataloaderServer))
	cCustomerOsPlatformOwner = client.New(common.WithContext(customCustomerOsPlatformOwnerCtx, dataloaderServer))
	cAdmin = client.New(common.WithContext(customAdminCtx, dataloaderServer))
	cAdminWithTenant = client.New(common.WithContext(customAdminWTenantCtx, dataloaderServer))
}

func getQuery(fileName string) string {
	b, err := os.ReadFile(fmt.Sprintf("test_queries/%s.txt", fileName))
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}

func assertRawResponseSuccess(t *testing.T, response *client.Response, err error) {
	require.Nil(t, err)
	require.NotNil(t, response)
	if response.Errors != nil {
		log.Println(fmt.Sprintf("Error in response: %v", string(response.Errors)))
	}
	require.NotNil(t, response.Data)
	require.Nil(t, response.Errors)
}

func assertRawResponseError(t *testing.T, response *client.Response, err error) {
	require.Nil(t, err)
	require.NotNil(t, response)
	if response.Errors != nil {
		log.Println(fmt.Sprintf("Error in response: %v", string(response.Errors)))
	}
	require.NotNil(t, response.Data)
	require.NotNil(t, response.Errors)
}

func callGraphQL(t *testing.T, queryLocation string, vars map[string]interface{}) (rawResponse *client.Response) {
	// Transform map into var args of options
	options := make([]client.Option, 0, len(vars))
	for key, value := range vars {
		options = append(options, client.Var(key, value))
	}

	// Call RawPost with options
	rawResponse, err := c.RawPost(getQuery(queryLocation), options...)
	require.Nil(t, err)
	assertRawResponseSuccess(t, rawResponse, err)
	return
}

func callGraphQLExpectError(t *testing.T, queryLocation string, vars map[string]interface{}) (response GraphQlErrorResponse) {
	// Transform map into var args of options
	options := make([]client.Option, 0, len(vars))
	for key, value := range vars {
		options = append(options, client.Var(key, value))
	}

	// Call RawPost with options
	rawResponse, err := c.RawPost(getQuery(queryLocation), options...)
	require.Nil(t, err)
	assertRawResponseError(t, rawResponse, err)

	var rr []GraphQlErrorResponse

	err = json.Unmarshal(rawResponse.Errors, &rr)
	require.Nil(t, err)

	return rr[0]
}

type GraphQlErrorResponse struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}
