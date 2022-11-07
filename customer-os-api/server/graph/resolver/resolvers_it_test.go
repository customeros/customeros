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

func setupTestCase() func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		cleanupAllData(driver)
	}
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

func assertRawResponseSuccess(t *testing.T, response *client.Response, err error) {
	require.Nil(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Data)
	require.Nil(t, response.Errors)
}

func TestQueryGetTenantUsers(t *testing.T) {
	defer setupTestCase()(t)
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
	assertRawResponseSuccess(t, rawResponse, err)

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
	defer setupTestCase()(t)
	createTenant(driver, "openline")
	createTenant(driver, "other")

	rawResponse, err := c.RawPost(getQuery("create_tenant_user"))
	assertRawResponseSuccess(t, rawResponse, err)

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
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	createTenant(driver, "otherTenant")

	rawResponse, err := c.RawPost(getQuery("create_contact"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		CreateContact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "MR", contact.CreateContact.Title.String())
	require.Equal(t, "first", contact.CreateContact.FirstName)
	require.Equal(t, "last", contact.CreateContact.LastName)
	require.Equal(t, "customer", *contact.CreateContact.ContactType)
	require.Equal(t, "Some notes...", *contact.CreateContact.Notes)
	require.Equal(t, "Some label", *contact.CreateContact.Label)

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

	//require.Equal(t, 2, getCountOfNodes(driver, "Tenant"))
	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 0, getCountOfNodes(driver, "ContactGroup"))
	require.Equal(t, 2, getCountOfNodes(driver, "TextCustomField"))
	require.Equal(t, 1, getCountOfNodes(driver, "Email"))
	require.Equal(t, 1, getCountOfNodes(driver, "PhoneNumber"))
}

func TestUpdateContact(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactId := createContact(driver, tenantName, entity.ContactEntity{
		Title:       model.PersonTitleMr.String(),
		FirstName:   "first",
		LastName:    "last",
		Label:       "label",
		ContactType: "type",
		Notes:       "notes",
	})

	rawResponse, err := c.RawPost(getQuery("update_contact"), client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		UpdateContact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "DR", contact.UpdateContact.Title.String())
	require.Equal(t, "updated first", contact.UpdateContact.FirstName)
	require.Equal(t, "updated last", contact.UpdateContact.LastName)
	require.Equal(t, "updated type", *contact.UpdateContact.ContactType)
	require.Equal(t, "updated notes", *contact.UpdateContact.Notes)
	require.Equal(t, "updated label", *contact.UpdateContact.Label)
}

func TestMergeFieldsSetToContact_AllowMultipleFieldsSetWithSameNameOnDifferentContacts(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactId1 := createContact(driver, tenantName, entity.ContactEntity{
		Title:     model.PersonTitleMr.String(),
		FirstName: "first",
		LastName:  "last",
	})
	contactId2 := createContact(driver, tenantName, entity.ContactEntity{
		Title:     model.PersonTitleMr.String(),
		FirstName: "first",
		LastName:  "last",
	})

	rawResponse1, err := c.RawPost(getQuery("merge_fields_set_to_contact"), client.Var("contactId", contactId1))
	rawResponse2, err := c.RawPost(getQuery("merge_fields_set_to_contact"), client.Var("contactId", contactId2))
	assertRawResponseSuccess(t, rawResponse1, err)
	assertRawResponseSuccess(t, rawResponse2, err)

	var fieldsSet1 struct {
		MergeFieldsSetToContact model.FieldsSet
	}
	var fieldsSet2 struct {
		MergeFieldsSetToContact model.FieldsSet
	}

	err = decode.Decode(rawResponse1.Data.(map[string]interface{}), &fieldsSet1)
	require.Nil(t, err)
	err = decode.Decode(rawResponse2.Data.(map[string]interface{}), &fieldsSet2)
	require.Nil(t, err)
	require.NotNil(t, fieldsSet1)
	require.NotNil(t, fieldsSet2)

	require.NotNil(t, fieldsSet1.MergeFieldsSetToContact.ID)
	require.NotNil(t, fieldsSet2.MergeFieldsSetToContact.ID)
	require.NotEqual(t, fieldsSet1.MergeFieldsSetToContact.ID, fieldsSet2.MergeFieldsSetToContact.ID)
	require.Equal(t, "some name", fieldsSet1.MergeFieldsSetToContact.Name)
	require.Equal(t, "some type", fieldsSet1.MergeFieldsSetToContact.Type)
	require.NotNil(t, fieldsSet1.MergeFieldsSetToContact.Added)
	require.Equal(t, "some name", fieldsSet2.MergeFieldsSetToContact.Name)
	require.Equal(t, "some type", fieldsSet2.MergeFieldsSetToContact.Type)
	require.NotNil(t, fieldsSet2.MergeFieldsSetToContact.Added)

	require.Equal(t, 2, getCountOfNodes(driver, "FieldsSet"))
}
