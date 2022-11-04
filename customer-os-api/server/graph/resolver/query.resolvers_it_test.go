package resolver

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/integration_tests"
	"github.com/openline-ai/openline-customer-os/customer-os-api/service/container"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"os"
	"testing"
)

var (
	neo4jContainer testcontainers.Container
	driver         *neo4j.Driver
	c              *client.Client
)

var tenantUsers struct {
	TenantUsers model.TenantUsersPage
}

func TestMain(m *testing.M) {
	neo4jContainer, driver = integration_tests.InitTestNeo4jDB()
	defer func(dbContainer testcontainers.Container, driver neo4j.Driver, ctx context.Context) {
		integration_tests.Close(driver, "Driver")
		integration_tests.Terminate(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	prepareClient()

	os.Exit(m.Run())
}

func prepareClient() {
	serviceContainer := container.InitServices(driver)
	graphResolver := NewResolver(serviceContainer)
	customCtx := &common.CustomContext{
		Tenant: "openline",
	}
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graphResolver}))
	h := common.CreateContext(customCtx, server)
	c = client.New(h)
}

func getQuery(fileName string) string {
	b, err := os.ReadFile(fmt.Sprintf("test_query/%s.txt", fileName))
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}

func assertRawResponseNotNil(t *testing.T, response *client.Response, err error) {
	require.Nil(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Data)
}

func TestQueryGetTenantUsers(t *testing.T) {
	tenant := "openline"
	otherTenant := "other"
	createTenant(driver, tenant)
	createTenant(driver, otherTenant)
	createTenantUser(driver, tenant, entity.TenantUserEntity{
		FirstName: "first",
		LastName:  "last",
		Email:     "test@openline.ai",
	})
	createTenantUser(driver, otherTenant, entity.TenantUserEntity{
		FirstName: "otherFirst",
		LastName:  "otherLast",
		Email:     "otherEmail",
	})

	rawResponse, err := c.RawPost(getQuery("get_tenant_users"))
	assertRawResponseNotNil(t, rawResponse, err)

	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &tenantUsers)
	require.Nil(t, err)
	require.NotNil(t, tenantUsers)
	require.Equal(t, 1, tenantUsers.TenantUsers.TotalPages)
	require.Equal(t, int64(1), tenantUsers.TenantUsers.TotalElements)
	require.Equal(t, "first", tenantUsers.TenantUsers.Content[0].FirstName)
	require.Equal(t, "last", tenantUsers.TenantUsers.Content[0].LastName)
	require.Equal(t, "test@openline.ai", tenantUsers.TenantUsers.Content[0].Email)
	require.NotNil(t, "test@openline.ai", tenantUsers.TenantUsers.Content[0].CreatedAt)
}
