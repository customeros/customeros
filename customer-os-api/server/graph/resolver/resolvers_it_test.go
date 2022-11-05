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

const tenantName = "openline"

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
	b, err := os.ReadFile(fmt.Sprintf("test_queries/%s.txt", fileName))
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

	var tenantUsers struct {
		TenantUsers model.TenantUsersPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &tenantUsers)
	require.Nil(t, err)
	require.NotNil(t, tenantUsers)
	require.Equal(t, 1, tenantUsers.TenantUsers.TotalPages)
	require.Equal(t, int64(1), tenantUsers.TenantUsers.TotalElements)
	require.Equal(t, "first", tenantUsers.TenantUsers.Content[0].FirstName)
	require.Equal(t, "last", tenantUsers.TenantUsers.Content[0].LastName)
	require.Equal(t, "test@openline.ai", tenantUsers.TenantUsers.Content[0].Email)
	require.NotNil(t, tenantUsers.TenantUsers.Content[0].CreatedAt)
}

func TestCreateTenantUser(t *testing.T) {
	createTenant(driver, "openline")
	createTenant(driver, "other")

	rawResponse, err := c.RawPost(getQuery("create_tenant_user"))
	assertRawResponseNotNil(t, rawResponse, err)

	var tenantUser struct {
		CreateTenantUser model.TenantUser
	}

	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &tenantUser)
	require.Nil(t, err)
	require.NotNil(t, tenantUser)
	require.Equal(t, "first", tenantUser.CreateTenantUser.FirstName)
	require.Equal(t, "last", tenantUser.CreateTenantUser.LastName)
	require.Equal(t, "user@openline.ai", tenantUser.CreateTenantUser.Email)
	require.NotNil(t, tenantUser.CreateTenantUser.CreatedAt)
	require.NotNil(t, tenantUser.CreateTenantUser.ID)
}

func TestCreateContact(t *testing.T) {
	createTenant(driver, tenantName)
	createTenant(driver, "otherTenant")

	rawResponse, err := c.RawPost(getQuery("create_contact"))
	assertRawResponseNotNil(t, rawResponse, err)

	var contact struct {
		CreateContact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "first", contact.CreateContact.FirstName)
	require.Equal(t, "last", contact.CreateContact.LastName)
	require.Equal(t, "MR", contact.CreateContact.Title.String())
	require.Equal(t, "customer", *contact.CreateContact.ContactType)
	require.Equal(t, "Some notes...", *contact.CreateContact.Notes)

	require.Equal(t, 2, len(contact.CreateContact.TextCustomFields))
	require.Equal(t, "field1", contact.CreateContact.TextCustomFields[0].Name)
	require.Equal(t, "value1", contact.CreateContact.TextCustomFields[0].Value)
	require.NotNil(t, contact.CreateContact.TextCustomFields[0].ID)
	require.Equal(t, "field2", contact.CreateContact.TextCustomFields[1].Name)
	require.Equal(t, "value2", contact.CreateContact.TextCustomFields[1].Value)
	require.NotNil(t, contact.CreateContact.TextCustomFields[1].ID)

	require.Equal(t, 1, len(contact.CreateContact.Emails))
	require.NotNil(t, contact.CreateContact.Emails[0].ID)
	require.Equal(t, "contact@abc.com", contact.CreateContact.Emails[0].Email)
	require.Equal(t, "WORK", contact.CreateContact.Emails[0].Label.String())
	require.Equal(t, false, contact.CreateContact.Emails[0].Primary)

	require.Equal(t, 1, len(contact.CreateContact.PhoneNumbers))
	require.NotNil(t, contact.CreateContact.PhoneNumbers[0].ID)
	require.Equal(t, "+1234567890", contact.CreateContact.PhoneNumbers[0].Number)
	require.Equal(t, "MOBILE", contact.CreateContact.PhoneNumbers[0].Label.String())
	require.Equal(t, true, contact.CreateContact.PhoneNumbers[0].Primary)

	require.Equal(t, 1, len(contact.CreateContact.CompanyPositions))
	require.Equal(t, "abc", contact.CreateContact.CompanyPositions[0].CompanyName)
	require.Equal(t, "CTO", *contact.CreateContact.CompanyPositions[0].JobTitle)

	require.Equal(t, 0, len(contact.CreateContact.Groups))
}
