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

func TestQueryResolver_TenantUsers(t *testing.T) {
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

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantUsers)
	require.Nil(t, err)
	require.NotNil(t, tenantUsers)
	require.Equal(t, 1, tenantUsers.TenantUsers.TotalPages)
	require.Equal(t, int64(1), tenantUsers.TenantUsers.TotalElements)
	require.Equal(t, "first", tenantUsers.TenantUsers.Content[0].FirstName)
	require.Equal(t, "last", tenantUsers.TenantUsers.Content[0].LastName)
	require.Equal(t, "test@openline.ai", tenantUsers.TenantUsers.Content[0].Email)
	require.NotNil(t, tenantUsers.TenantUsers.Content[0].CreatedAt)
}

func TestQueryResolver_ContactByEmail(t *testing.T) {
	defer setupTestCase()(t)
	otherTenant := "other"
	createTenant(driver, tenantName)
	createTenant(driver, otherTenant)
	contactId1 := createDefaultContact(driver, tenantName)
	contactId2 := createDefaultContact(driver, otherTenant)
	addEmailToContact(driver, contactId1, "test@test.com", true, "MAIN")
	addEmailToContact(driver, contactId2, "test@test.com", true, "MAIN")

	rawResponse, err := c.RawPost(getQuery("get_contact_by_email"), client.Var("email", "test@test.com"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		ContactByEmail model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, contactId1, contact.ContactByEmail.ID)
}

func TestQueryResolver_ContactByPhone(t *testing.T) {
	defer setupTestCase()(t)
	otherTenant := "other"
	createTenant(driver, tenantName)
	createTenant(driver, otherTenant)
	contactId1 := createDefaultContact(driver, tenantName)
	contactId2 := createDefaultContact(driver, otherTenant)
	addPhoneNumberToContact(driver, contactId1, "+1234567890", false, "OTHER")
	addPhoneNumberToContact(driver, contactId2, "+1234567890", true, "MAIN")

	rawResponse, err := c.RawPost(getQuery("get_contact_by_phone"), client.Var("number", "+1234567890"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		ContactByPhone model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, contactId1, contact.ContactByPhone.ID)
}

func TestMutationResolver_CreateTenantUser(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, "openline")
	createTenant(driver, "other")

	rawResponse, err := c.RawPost(getQuery("create_tenant_user"))
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantUser struct {
		CreateTenantUser model.TenantUser
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantUser)
	require.Nil(t, err)
	require.NotNil(t, tenantUser)
	require.Equal(t, "first", tenantUser.CreateTenantUser.FirstName)
	require.Equal(t, "last", tenantUser.CreateTenantUser.LastName)
	require.Equal(t, "user@openline.ai", tenantUser.CreateTenantUser.Email)
	require.NotNil(t, tenantUser.CreateTenantUser.CreatedAt)
	require.NotNil(t, tenantUser.CreateTenantUser.ID)
}

func TestMutationResolver_CreateContact(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	createTenant(driver, "otherTenant")

	rawResponse, err := c.RawPost(getQuery("create_contact"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		CreateContact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
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

func TestMutationResolver_UpdateContact(t *testing.T) {
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

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "DR", contact.UpdateContact.Title.String())
	require.Equal(t, "updated first", contact.UpdateContact.FirstName)
	require.Equal(t, "updated last", contact.UpdateContact.LastName)
	require.Equal(t, "updated type", *contact.UpdateContact.ContactType)
	require.Equal(t, "updated notes", *contact.UpdateContact.Notes)
	require.Equal(t, "updated label", *contact.UpdateContact.Label)
}

func TestMutationResolver_MergeFieldSetToContact_AllowMultipleFieldSetWithSameNameOnDifferentContacts(t *testing.T) {
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

	var fieldSet1 struct {
		MergeFieldSetToContact model.FieldSet
	}
	var fieldSet2 struct {
		MergeFieldSetToContact model.FieldSet
	}

	err = decode.Decode(rawResponse1.Data.(map[string]any), &fieldSet1)
	require.Nil(t, err)
	err = decode.Decode(rawResponse2.Data.(map[string]any), &fieldSet2)
	require.Nil(t, err)
	require.NotNil(t, fieldSet1)
	require.NotNil(t, fieldSet2)

	require.NotNil(t, fieldSet1.MergeFieldSetToContact.ID)
	require.NotNil(t, fieldSet2.MergeFieldSetToContact.ID)
	require.NotEqual(t, fieldSet1.MergeFieldSetToContact.ID, fieldSet2.MergeFieldSetToContact.ID)
	require.Equal(t, "some name", fieldSet1.MergeFieldSetToContact.Name)
	require.Equal(t, "some type", fieldSet1.MergeFieldSetToContact.Type)
	require.NotNil(t, fieldSet1.MergeFieldSetToContact.Added)
	require.Equal(t, "some name", fieldSet2.MergeFieldSetToContact.Name)
	require.Equal(t, "some type", fieldSet2.MergeFieldSetToContact.Type)
	require.NotNil(t, fieldSet2.MergeFieldSetToContact.Added)

	require.Equal(t, 2, getCountOfNodes(driver, "FieldSet"))
}

func TestMutationResolver_MergeTextCustomFieldToFieldSet(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)
	fieldSetId := createDefaultFieldSet(driver, contactId)

	rawResponse, err := c.RawPost(getQuery("merge_text_field_to_fields_set"),
		client.Var("contactId", contactId), client.Var("fieldSetId", fieldSetId))
	assertRawResponseSuccess(t, rawResponse, err)

	var textField struct {
		MergeTextCustomFieldToFieldSet model.TextCustomField
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &textField)
	require.Nil(t, err)

	require.Equal(t, "some name", textField.MergeTextCustomFieldToFieldSet.Name)
	require.Equal(t, "some value", textField.MergeTextCustomFieldToFieldSet.Value)
	require.NotNil(t, textField.MergeTextCustomFieldToFieldSet.ID)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, getCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 1, getCountOfNodes(driver, "TextCustomField"))
}

func TestMutationResolver_UpdateTextCustomFieldInFieldSet(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)
	fieldSetId := createDefaultFieldSet(driver, contactId)
	fieldId := createDefaultTextFieldInSet(driver, fieldSetId)

	rawResponse, err := c.RawPost(getQuery("update_text_field_in_fields_set"),
		client.Var("contactId", contactId),
		client.Var("fieldSetId", fieldSetId),
		client.Var("fieldId", fieldId))
	assertRawResponseSuccess(t, rawResponse, err)

	var textField struct {
		UpdateTextCustomFieldInFieldSet model.TextCustomField
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &textField)
	require.Nil(t, err)

	require.Equal(t, "new name", textField.UpdateTextCustomFieldInFieldSet.Name)
	require.Equal(t, "new value", textField.UpdateTextCustomFieldInFieldSet.Value)
	require.Equal(t, fieldId, textField.UpdateTextCustomFieldInFieldSet.ID)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, getCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 1, getCountOfNodes(driver, "TextCustomField"))
}

func TestMutationResolver_RemoveTextCustomFieldFromFieldSetByID(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)
	fieldSetId := createDefaultFieldSet(driver, contactId)
	fieldId := createDefaultTextFieldInSet(driver, fieldSetId)

	rawResponse, err := c.RawPost(getQuery("remove_text_field_from_fields_set"),
		client.Var("contactId", contactId),
		client.Var("fieldSetId", fieldSetId),
		client.Var("fieldId", fieldId))
	assertRawResponseSuccess(t, rawResponse, err)

	var textField struct {
		RemoveTextCustomFieldFromFieldSetByID model.BooleanResult
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &textField)
	require.Nil(t, err)

	require.Equal(t, true, textField.RemoveTextCustomFieldFromFieldSetByID.Result)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, getCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 0, getCountOfNodes(driver, "TextCustomField"))
}

func TestMutationResolver_RemoveFieldSetFromContact(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)
	fieldSetId := createDefaultFieldSet(driver, contactId)
	createDefaultTextFieldInSet(driver, fieldSetId)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, getCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 1, getCountOfNodes(driver, "TextCustomField"))

	rawResponse, err := c.RawPost(getQuery("remove_fields_set_from_contact"),
		client.Var("contactId", contactId),
		client.Var("fieldSetId", fieldSetId))
	assertRawResponseSuccess(t, rawResponse, err)

	var textField struct {
		RemoveFieldSetFromContact model.BooleanResult
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &textField)
	require.Nil(t, err)

	require.Equal(t, true, textField.RemoveFieldSetFromContact.Result)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 0, getCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 0, getCountOfNodes(driver, "TextCustomField"))
}
