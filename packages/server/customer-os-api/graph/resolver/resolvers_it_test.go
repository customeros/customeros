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
		Tenant: tenantName,
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

func TestQueryResolver_Users(t *testing.T) {
	defer setupTestCase()(t)
	otherTenant := "other"
	createTenant(driver, tenantName)
	createTenant(driver, otherTenant)
	createUser(driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Email:     "test@openline.ai",
	})
	createUser(driver, otherTenant, entity.UserEntity{
		FirstName: "otherFirst",
		LastName:  "otherLast",
		Email:     "otherEmail",
	})

	rawResponse, err := c.RawPost(getQuery("get_users"))
	assertRawResponseSuccess(t, rawResponse, err)

	var users struct {
		Users model.UserPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &users)
	require.Nil(t, err)
	require.NotNil(t, users)
	require.Equal(t, 1, users.Users.TotalPages)
	require.Equal(t, int64(1), users.Users.TotalElements)
	require.Equal(t, "first", users.Users.Content[0].FirstName)
	require.Equal(t, "last", users.Users.Content[0].LastName)
	require.Equal(t, "test@openline.ai", users.Users.Content[0].Email)
	require.NotNil(t, users.Users.Content[0].CreatedAt)
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

	rawResponse, err := c.RawPost(getQuery("get_contact_by_phone"), client.Var("e164", "+1234567890"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		ContactByPhone model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, contactId1, contact.ContactByPhone.ID)
}

func TestMutationResolver_CreateUser(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	createTenant(driver, "other")

	rawResponse, err := c.RawPost(getQuery("create_user"))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		CreateUser model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, "first", user.CreateUser.FirstName)
	require.Equal(t, "last", user.CreateUser.LastName)
	require.Equal(t, "user@openline.ai", user.CreateUser.Email)
	require.NotNil(t, user.CreateUser.CreatedAt)
	require.NotNil(t, user.CreateUser.ID)
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

	require.Equal(t, 5, len(contact.CreateContact.CustomFields))

	boolField := contact.CreateContact.CustomFields[0]
	require.NotNil(t, boolField.GetID())
	require.Equal(t, "boolField", boolField.Name)
	require.Equal(t, model.CustomFieldDataTypeBool, boolField.Datatype)
	require.Equal(t, true, boolField.Value.RealValue())

	decimalField := contact.CreateContact.CustomFields[1]
	require.NotNil(t, decimalField.GetID())
	require.Equal(t, "decimalField", decimalField.Name)
	require.Equal(t, model.CustomFieldDataTypeDecimal, decimalField.Datatype)
	require.Equal(t, 0.001, decimalField.Value.RealValue())

	integerField := contact.CreateContact.CustomFields[2]
	require.NotNil(t, integerField.GetID())
	require.Equal(t, "integerField", integerField.Name)
	require.Equal(t, model.CustomFieldDataTypeInteger, integerField.Datatype)
	// issue in decoding, int converted to float 64
	require.Equal(t, float64(123), integerField.Value.RealValue())

	textField := contact.CreateContact.CustomFields[3]
	require.NotNil(t, textField.GetID())
	require.Equal(t, "textField", textField.Name)
	require.Equal(t, model.CustomFieldDataTypeText, textField.Datatype)
	require.Equal(t, "value1", textField.Value.RealValue())

	timeField := contact.CreateContact.CustomFields[4]
	require.NotNil(t, timeField.GetID())
	require.Equal(t, "timeField", timeField.Name)
	require.Equal(t, model.CustomFieldDataTypeDatetime, timeField.Datatype)
	require.Equal(t, "2022-11-13T20:21:56.732Z", timeField.Value.RealValue())

	require.Equal(t, 1, len(contact.CreateContact.Emails))
	require.NotNil(t, contact.CreateContact.Emails[0].ID)
	require.Equal(t, "contact@abc.com", contact.CreateContact.Emails[0].Email)
	require.Equal(t, "WORK", contact.CreateContact.Emails[0].Label.String())
	require.Equal(t, false, contact.CreateContact.Emails[0].Primary)

	require.Equal(t, 1, len(contact.CreateContact.PhoneNumbers))
	require.NotNil(t, contact.CreateContact.PhoneNumbers[0].ID)
	require.Equal(t, "+1234567890", contact.CreateContact.PhoneNumbers[0].E164)
	require.Equal(t, "MOBILE", contact.CreateContact.PhoneNumbers[0].Label.String())
	require.Equal(t, true, contact.CreateContact.PhoneNumbers[0].Primary)

	require.Equal(t, 1, len(contact.CreateContact.Companies))
	require.Equal(t, "abc", contact.CreateContact.Companies[0].CompanyName)
	require.Equal(t, "CTO", *contact.CreateContact.Companies[0].JobTitle)

	require.Equal(t, 0, len(contact.CreateContact.Groups))

	require.Equal(t, 2, getCountOfNodes(driver, "Tenant"))
	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 0, getCountOfNodes(driver, "ContactGroup"))
	require.Equal(t, 5, getCountOfNodes(driver, "CustomField"))
	require.Equal(t, 1, getCountOfNodes(driver, "TextField"))
	require.Equal(t, 1, getCountOfNodes(driver, "IntField"))
	require.Equal(t, 1, getCountOfNodes(driver, "FloatField"))
	require.Equal(t, 1, getCountOfNodes(driver, "BoolField"))
	require.Equal(t, 1, getCountOfNodes(driver, "TimeField"))
	require.Equal(t, 1, getCountOfNodes(driver, "Email"))
	require.Equal(t, 1, getCountOfNodes(driver, "PhoneNumber"))
}

func TestMutationResolver_CreateContact_WithEntityDefinition(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	entityDefinitionId := createEntityDefinition(driver, tenantName, model.EntityDefinitionExtensionContact.String())
	fieldDefinitionId := addFieldDefinitionToEntity(driver, entityDefinitionId)

	rawResponse, err := c.RawPost(getQuery("create_contact_with_entity_definition"),
		client.Var("entityDefinitionId", entityDefinitionId),
		client.Var("fieldDefinitionId", fieldDefinitionId))
	assertRawResponseSuccess(t, rawResponse, err)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 0, getCountOfNodes(driver, "ContactGroup"))
	require.Equal(t, 0, getCountOfNodes(driver, "Company"))
	require.Equal(t, 2, getCountOfNodes(driver, "CustomField"))
	require.Equal(t, 2, getCountOfNodes(driver, "TextField"))
	require.Equal(t, 0, getCountOfNodes(driver, "Email"))
	require.Equal(t, 0, getCountOfNodes(driver, "PhoneNumber"))
	require.Equal(t, 1, getCountOfNodes(driver, "EntityDefinition"))
	require.Equal(t, 1, getCountOfNodes(driver, "CustomFieldDefinition"))
	require.Equal(t, 0, getCountOfNodes(driver, "FieldSetDefinition"))

	var contact struct {
		CreateContact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)

	require.Equal(t, entityDefinitionId, contact.CreateContact.Definition.ID)
	require.Equal(t, 2, len(contact.CreateContact.CustomFields))
	require.Equal(t, "field1", contact.CreateContact.CustomFields[0].Name)
	require.Equal(t, "TEXT", contact.CreateContact.CustomFields[0].Datatype.String())
	require.Equal(t, "value1", contact.CreateContact.CustomFields[0].Value.RealValue())
	require.Equal(t, fieldDefinitionId, contact.CreateContact.CustomFields[0].Definition.ID)
	require.NotNil(t, contact.CreateContact.CustomFields[0].GetID())
	require.Equal(t, "field2", contact.CreateContact.CustomFields[1].Name)
	require.Equal(t, "TEXT", contact.CreateContact.CustomFields[1].Datatype.String())
	require.Equal(t, "value2", contact.CreateContact.CustomFields[1].Value.RealValue())
	require.NotNil(t, contact.CreateContact.CustomFields[1].GetID())

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

	rawResponse1, err := c.RawPost(getQuery("merge_field_set_to_contact"), client.Var("contactId", contactId1))
	rawResponse2, err := c.RawPost(getQuery("merge_field_set_to_contact"), client.Var("contactId", contactId2))
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
	require.NotNil(t, fieldSet1.MergeFieldSetToContact.Added)
	require.Equal(t, "some name", fieldSet2.MergeFieldSetToContact.Name)
	require.NotNil(t, fieldSet2.MergeFieldSetToContact.Added)

	require.Equal(t, 2, getCountOfNodes(driver, "FieldSet"))
}

func TestMutationResolver_MergeCustomFieldToFieldSet(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)
	fieldSetId := createDefaultFieldSet(driver, contactId)

	rawResponse, err := c.RawPost(getQuery("merge_text_field_to_field_set"),
		client.Var("contactId", contactId), client.Var("fieldSetId", fieldSetId))
	assertRawResponseSuccess(t, rawResponse, err)

	var textField struct {
		MergeCustomFieldToFieldSet model.CustomField
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &textField)
	require.Nil(t, err)

	require.Equal(t, "some name", textField.MergeCustomFieldToFieldSet.Name)
	require.Equal(t, "some value", textField.MergeCustomFieldToFieldSet.Value.RealValue())
	require.Equal(t, model.CustomFieldDataTypeText, textField.MergeCustomFieldToFieldSet.Datatype)
	require.NotNil(t, textField.MergeCustomFieldToFieldSet.ID)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, getCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 1, getCountOfNodes(driver, "CustomField"))
	require.Equal(t, 1, getCountOfNodes(driver, "TextField"))
}

func TestMutationResolver_UpdateCustomFieldInFieldSet(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)
	fieldSetId := createDefaultFieldSet(driver, contactId)
	fieldId := createDefaultCustomFieldInSet(driver, fieldSetId)

	rawResponse, err := c.RawPost(getQuery("update_text_field_in_field_set"),
		client.Var("contactId", contactId),
		client.Var("fieldSetId", fieldSetId),
		client.Var("fieldId", fieldId))
	assertRawResponseSuccess(t, rawResponse, err)

	var textField struct {
		UpdateCustomFieldInFieldSet model.CustomField
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &textField)
	require.Nil(t, err)

	require.Equal(t, "new name", textField.UpdateCustomFieldInFieldSet.Name)
	require.Equal(t, "new value", textField.UpdateCustomFieldInFieldSet.Value.RealValue())
	require.Equal(t, model.CustomFieldDataTypeText, textField.UpdateCustomFieldInFieldSet.Datatype)
	require.Equal(t, fieldId, textField.UpdateCustomFieldInFieldSet.ID)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, getCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 1, getCountOfNodes(driver, "CustomField"))
}

func TestMutationResolver_RemoveCustomFieldFromFieldSetByID(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)
	fieldSetId := createDefaultFieldSet(driver, contactId)
	fieldId := createDefaultCustomFieldInSet(driver, fieldSetId)

	rawResponse, err := c.RawPost(getQuery("remove_text_field_from_field_set"),
		client.Var("contactId", contactId),
		client.Var("fieldSetId", fieldSetId),
		client.Var("fieldId", fieldId))
	assertRawResponseSuccess(t, rawResponse, err)

	var textField struct {
		RemoveCustomFieldFromFieldSetByID model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &textField)
	require.Nil(t, err)

	require.Equal(t, true, textField.RemoveCustomFieldFromFieldSetByID.Result)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, getCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 0, getCountOfNodes(driver, "CustomField"))
}

func TestMutationResolver_RemoveFieldSetFromContact(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)
	fieldSetId := createDefaultFieldSet(driver, contactId)
	createDefaultCustomFieldInSet(driver, fieldSetId)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, getCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 1, getCountOfNodes(driver, "CustomField"))

	rawResponse, err := c.RawPost(getQuery("remove_field_set_from_contact"),
		client.Var("contactId", contactId),
		client.Var("fieldSetId", fieldSetId))
	assertRawResponseSuccess(t, rawResponse, err)

	var textField struct {
		RemoveFieldSetFromContact model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &textField)
	require.Nil(t, err)

	require.Equal(t, true, textField.RemoveFieldSetFromContact.Result)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 0, getCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 0, getCountOfNodes(driver, "CustomField"))
}

func TestMutationResolver_CreateEntityDefinition(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	createTenant(driver, "other")

	rawResponse, err := c.RawPost(getQuery("create_entity_definition"))
	assertRawResponseSuccess(t, rawResponse, err)

	var entityDefinition struct {
		CreateEntityDefinition model.EntityDefinition
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &entityDefinition)
	actual := entityDefinition.CreateEntityDefinition
	require.Nil(t, err)
	require.NotNil(t, actual)
	require.NotNil(t, actual.ID)
	require.NotNil(t, actual.Added)
	require.Equal(t, "the entity definition name", actual.Name)
	require.Equal(t, 1, actual.Version)
	require.Nil(t, actual.Extends)

	require.Equal(t, 2, len(actual.FieldSets))

	set := actual.FieldSets[0]
	require.NotNil(t, set.ID)
	require.Equal(t, "set 1", set.Name)
	require.Equal(t, 1, set.Order)
	require.Equal(t, 2, len(set.CustomFields))

	field := set.CustomFields[0]
	require.NotNil(t, field)
	require.Equal(t, "field 3", field.Name)
	require.Equal(t, 1, field.Order)
	require.Equal(t, true, field.Mandatory)
	require.Equal(t, model.CustomFieldDefinitionTypeText, field.Type)
	require.Nil(t, field.Min)
	require.Nil(t, field.Max)
	require.Nil(t, field.Length)

	field = set.CustomFields[1]
	require.NotNil(t, field)
	require.Equal(t, "field 4", field.Name)
	require.Equal(t, 2, field.Order)
	require.Equal(t, false, field.Mandatory)
	require.Equal(t, model.CustomFieldDefinitionTypeText, field.Type)
	require.Equal(t, 10, *field.Min)
	require.Equal(t, 990, *field.Max)
	require.Equal(t, 2550, *field.Length)

	set = actual.FieldSets[1]
	require.NotNil(t, set.ID)
	require.Equal(t, "set 2", set.Name)
	require.Equal(t, 2, set.Order)
	require.Equal(t, 0, len(set.CustomFields))

	field = actual.CustomFields[0]
	require.NotNil(t, field)
	require.Equal(t, "field 1", field.Name)
	require.Equal(t, 1, field.Order)
	require.Equal(t, true, field.Mandatory)
	require.Equal(t, model.CustomFieldDefinitionTypeText, field.Type)
	require.Nil(t, field.Min)
	require.Nil(t, field.Max)
	require.Nil(t, field.Length)

	field = actual.CustomFields[1]
	require.NotNil(t, field)
	require.Equal(t, "field 2", field.Name)
	require.Equal(t, 2, field.Order)
	require.Equal(t, false, field.Mandatory)
	require.Equal(t, model.CustomFieldDefinitionTypeText, field.Type)
	require.Equal(t, 1, *field.Min)
	require.Equal(t, 99, *field.Max)
	require.Equal(t, 255, *field.Length)

	require.Equal(t, 1, getCountOfNodes(driver, "EntityDefinition"))
	require.Equal(t, 2, getCountOfNodes(driver, "FieldSetDefinition"))
	require.Equal(t, 4, getCountOfNodes(driver, "CustomFieldDefinition"))
}

func TestMutationResolver_CreateConversation_AutogenerateID(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	userId := createDefaultUser(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("create_conversation"),
		client.Var("contactId", contactId),
		client.Var("userId", userId))
	assertRawResponseSuccess(t, rawResponse, err)

	var conversation struct {
		CreateConversation model.Conversation
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &conversation)
	require.Nil(t, err)
	require.NotNil(t, conversation)
	require.NotNil(t, conversation.CreateConversation.ID)
	require.NotNil(t, conversation.CreateConversation.StartedAt)
}

func TestMutationResolver_CreateConversation_WithGivenID(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	conversationId := "Some conversation ID"
	userId := createDefaultUser(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("create_conversation_with_id"),
		client.Var("contactId", contactId),
		client.Var("userId", userId),
		client.Var("conversationId", conversationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var conversation struct {
		CreateConversation model.Conversation
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &conversation)
	require.Nil(t, err)
	require.NotNil(t, conversation)
	require.NotNil(t, conversation.CreateConversation.StartedAt)
	require.Equal(t, conversationId, conversation.CreateConversation.ID)
}

func TestMutationResolver_ContactTypeCreate(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	createTenant(driver, "otherTenantName")

	rawResponse, err := c.RawPost(getQuery("create_contact_type"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactType struct {
		ContactType_Create model.ContactType
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactType)
	require.Nil(t, err)
	require.NotNil(t, contactType)
	require.NotNil(t, contactType.ContactType_Create.ID)
	require.Equal(t, "the contact type", contactType.ContactType_Create.Name)

	require.Equal(t, 1, getCountOfNodes(driver, "ContactType"))
}

func TestMutationResolver_ContactTypeUpdate(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactTypeId := createContactType(driver, tenantName, "original type")

	rawResponse, err := c.RawPost(getQuery("update_contact_type"),
		client.Var("contactTypeId", contactTypeId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactType struct {
		ContactType_Update model.ContactType
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactType)
	require.Nil(t, err)
	require.NotNil(t, contactType)
	require.Equal(t, contactTypeId, contactType.ContactType_Update.ID)
	require.Equal(t, "updated type", contactType.ContactType_Update.Name)

	require.Equal(t, 1, getCountOfNodes(driver, "ContactType"))
}
