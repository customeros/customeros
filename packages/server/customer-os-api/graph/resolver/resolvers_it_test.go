package resolver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/test"
	"log"
	"os"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

var (
	neo4jContainer testcontainers.Container
	driver         *neo4j.DriverWithContext

	postgresContainer        testcontainers.Container
	postgresGormDB           *gorm.DB
	postgresSqlDB            *sql.DB
	c                        *client.Client
	cOwner                   *client.Client
	cCustomerOsPlatformOwner *client.Client
	cAdmin                   *client.Client
	cAdminWithTenant         *client.Client

	//services
	customerOsApiServices *service.Services
)

const tenantName = "openline"
const testUserId = "test-user-id"
const testUserEmail = "test-user-email"
const testContactId = "test-contact-id"
const testPlayerId = "test-player-id"

func TestMain(m *testing.M) {
	neo4jContainer, driver = neo4jt.InitTestNeo4jDB()
	defer func(dbContainer testcontainers.Container, driver neo4j.DriverWithContext, ctx context.Context) {
		neo4jt.CloseDriver(driver)
		neo4jt.TerminateNeo4j(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	postgresContainer, postgresGormDB, postgresSqlDB = neo4jt.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		neo4jt.TerminatePostgres(postgresContainer, ctx)
	}(postgresContainer, context.Background())

	// Start RabbitMQ container
	//_, cleanup, err := neo4jt.SetupRabbitMQTestContainer()
	//if err != nil {
	//	log.Fatalf("Failed to setup RabbitMQ test container: %v", err)
	//}
	//defer cleanup()

	prepareClient()

	os.Exit(m.Run())
}

func tearDownTestCase(ctx context.Context) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jtest.CleanupAllData(ctx, driver)
	}
}

func prepareClient() {
	appLogger := logger.NewAppLogger(&logger.Config{
		DevMode: true,
	})
	appLogger.InitLogger()

	testDialFactory := events_platform.NewTestDialFactory()
	gRPCconn, _ := testDialFactory.GetEventsProcessingPlatformConn()

	grpcClient := grpc_client.InitClients(gRPCconn)
	commonServices := commonService.InitServices(&commonConfig.GlobalConfig{}, postgresGormDB, driver, "neo4j", grpcClient, appLogger)
	customerOsApiServices = service.InitServices(appLogger, driver, &config.Config{}, commonServices, grpcClient, postgresGormDB)
	graphResolver := NewResolver(appLogger, customerOsApiServices, customerOsApiServices.CommonServices.GrpcClients, &config.Config{})
	loader := dataloader.NewDataLoader(customerOsApiServices)
	customCtx := &common.CustomContext{
		Tenant:     tenantName,
		UserId:     testUserId,
		UserEmail:  testUserEmail,
		IdentityId: testPlayerId,
		Roles:      []string{model.RoleUser.String()},
	}

	customOwnerCtx := &common.CustomContext{
		Tenant:     tenantName,
		UserId:     testUserId,
		UserEmail:  testUserEmail,
		IdentityId: testPlayerId,
		Roles:      []string{model.RoleUser.String(), model.RoleOwner.String()},
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
