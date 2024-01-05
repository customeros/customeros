package resolver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	commonAuthService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
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
	services *service.Services
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
		neo4jt.Terminate(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	postgresContainer, postgresGormDB, postgresSqlDB = postgres.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		err := postgresContainer.Terminate(ctx)
		if err != nil {
			log.Fatal("Error during container termination")
		}
	}(postgresContainer, context.Background())

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

	commonServices := commonService.InitServices(postgresGormDB, driver)
	commonAuthServices := commonAuthService.InitServices(nil, postgresGormDB)
	testDialFactory := events_platform.NewTestDialFactory()
	gRPCconn, _ := testDialFactory.GetEventsProcessingPlatformConn()
	services = service.InitServices(appLogger, driver, &config.Config{}, commonServices, commonAuthServices, grpc_client.InitClients(gRPCconn))
	graphResolver := NewResolver(appLogger, services, grpc_client.InitClients(gRPCconn))
	loader := dataloader.NewDataLoader(services)
	customCtx := &common.CustomContext{
		Tenant:     tenantName,
		UserId:     testUserId,
		UserEmail:  testUserEmail,
		IdentityId: testPlayerId,
		Roles:      []model.Role{model.RoleUser},
	}

	customOwnerCtx := &common.CustomContext{
		Tenant:     tenantName,
		UserId:     testUserId,
		UserEmail:  testUserEmail,
		IdentityId: testPlayerId,
		Roles:      []model.Role{model.RoleUser, model.RoleOwner},
	}
	customCustomerOsPlatformOwnerCtx := &common.CustomContext{
		Tenant:    tenantName,
		UserId:    testUserId,
		UserEmail: testUserEmail,
		Roles:     []model.Role{model.RoleUser, model.RoleCustomerOsPlatformOwner},
	}
	customAdminCtx := &common.CustomContext{
		Roles: []model.Role{model.RoleAdmin},
	}

	customAdminWTenantCtx := &common.CustomContext{
		Tenant: tenantName,
		Roles:  []model.Role{model.RoleAdmin},
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

// Deprecated, use neo4jtest.AssertNeo4jNodeCount instead
func assertNeo4jNodeCount(ctx context.Context, t *testing.T, driver *neo4j.DriverWithContext, nodes map[string]int) {
	for name, expectedCount := range nodes {
		actualCount := neo4jt.GetCountOfNodes(ctx, driver, name)
		require.Equal(t, expectedCount, actualCount, "Unexpected count for node: "+name)
	}
}

// Deprecated, use neo4jtest.AssertNeo4jRelationCount instead
func assertNeo4jRelationCount(ctx context.Context, t *testing.T, driver *neo4j.DriverWithContext, relations map[string]int) {
	for name, expectedCount := range relations {
		actualCount := neo4jt.GetCountOfRelationships(ctx, driver, name)
		require.Equal(t, expectedCount, actualCount, "Unexpected count for relationship: "+name)
	}
}

// Deprecated, use neo4jtest.AssertRelationship instead
func assertRelationship(ctx context.Context, t *testing.T, driver *neo4j.DriverWithContext, fromNodeId, relationshipType, toNodeId string) {
	rel, err := neo4jt.GetRelationship(ctx, driver, fromNodeId, toNodeId)
	require.Nil(t, err)
	require.NotNil(t, rel)
	require.Equal(t, relationshipType, rel.Type)
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

	var rr struct {
		GraphQlErrorResponse GraphQlErrorResponse
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &rr)
	require.Nil(t, err)

	return rr.GraphQlErrorResponse
}

type GraphQlErrorResponse struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}
