package resolver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service/container"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
	"os"
	"testing"
)

var (
	neo4jContainer testcontainers.Container
	driver         *neo4j.Driver

	postgresContainer testcontainers.Container
	postgresGormDB    *gorm.DB
	postgresSqlDB     *sql.DB
	c                 *client.Client
)

const tenantName = "openline"

func TestMain(m *testing.M) {
	neo4jContainer, driver = neo4jt.InitTestNeo4jDB()
	defer func(dbContainer testcontainers.Container, driver neo4j.Driver, ctx context.Context) {
		neo4jt.Close(driver, "Driver")
		neo4jt.Terminate(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	postgresContainer, postgresGormDB, postgresSqlDB = postgres.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		err := postgresContainer.Terminate(ctx)
		if err != nil {
			logrus.Fatal("Error during container termination")
		}
	}(postgresContainer, context.Background())

	prepareClient()

	os.Exit(m.Run())
}

func tearDownTestCase() func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jt.CleanupAllData(driver)
	}
}

func prepareClient() {
	repositoryContainer := commonRepository.InitCommonRepositories(postgresGormDB)
	serviceContainer := container.InitServices(driver)
	graphResolver := NewResolver(serviceContainer, repositoryContainer)
	customCtx := &common.CustomContext{
		Tenant: tenantName,
	}
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graphResolver}))
	h := common.WithContext(customCtx, server)
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

func TestMutationResolver_FieldSetMergeToContact_AllowMultipleFieldSetWithSameNameOnDifferentContacts(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId1 := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     model.PersonTitleMr.String(),
		FirstName: "first",
		LastName:  "last",
	})
	contactId2 := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     model.PersonTitleMr.String(),
		FirstName: "first",
		LastName:  "last",
	})

	rawResponse1, err := c.RawPost(getQuery("merge_field_set_to_contact"), client.Var("contactId", contactId1))
	rawResponse2, err := c.RawPost(getQuery("merge_field_set_to_contact"), client.Var("contactId", contactId2))
	assertRawResponseSuccess(t, rawResponse1, err)
	assertRawResponseSuccess(t, rawResponse2, err)

	var fieldSet1 struct {
		FieldSetMergeToContact model.FieldSet
	}
	var fieldSet2 struct {
		FieldSetMergeToContact model.FieldSet
	}

	err = decode.Decode(rawResponse1.Data.(map[string]any), &fieldSet1)
	require.Nil(t, err)
	err = decode.Decode(rawResponse2.Data.(map[string]any), &fieldSet2)
	require.Nil(t, err)
	require.NotNil(t, fieldSet1)
	require.NotNil(t, fieldSet2)

	require.NotNil(t, fieldSet1.FieldSetMergeToContact.ID)
	require.NotNil(t, fieldSet2.FieldSetMergeToContact.ID)
	require.NotEqual(t, fieldSet1.FieldSetMergeToContact.ID, fieldSet2.FieldSetMergeToContact.ID)
	require.Equal(t, "some name", fieldSet1.FieldSetMergeToContact.Name)
	require.NotNil(t, fieldSet1.FieldSetMergeToContact.Added)
	require.Equal(t, "some name", fieldSet2.FieldSetMergeToContact.Name)
	require.NotNil(t, fieldSet2.FieldSetMergeToContact.Added)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "FieldSet"))
}

func TestMutationResolver_MergeCustomFieldToFieldSet(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	fieldSetId := neo4jt.CreateDefaultFieldSet(driver, contactId)

	rawResponse, err := c.RawPost(getQuery("merge_custom_field_to_field_set"),
		client.Var("contactId", contactId), client.Var("fieldSetId", fieldSetId))
	assertRawResponseSuccess(t, rawResponse, err)

	var textField struct {
		CustomFieldMergeToFieldSet model.CustomField
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &textField)
	require.Nil(t, err)

	require.Equal(t, "some name", textField.CustomFieldMergeToFieldSet.Name)
	require.Equal(t, "some value", textField.CustomFieldMergeToFieldSet.Value.RealValue())
	require.Equal(t, "manual", *textField.CustomFieldMergeToFieldSet.Source)
	require.Equal(t, model.CustomFieldDataTypeText, textField.CustomFieldMergeToFieldSet.Datatype)
	require.NotNil(t, textField.CustomFieldMergeToFieldSet.ID)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "CustomField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "TextField"))
}

func TestMutationResolver_CustomFieldUpdateInFieldSet(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	fieldSetId := neo4jt.CreateDefaultFieldSet(driver, contactId)
	fieldId := neo4jt.CreateDefaultCustomFieldInSet(driver, fieldSetId)

	rawResponse, err := c.RawPost(getQuery("update_custom_field_in_field_set"),
		client.Var("contactId", contactId),
		client.Var("fieldSetId", fieldSetId),
		client.Var("fieldId", fieldId))
	assertRawResponseSuccess(t, rawResponse, err)

	var textField struct {
		CustomFieldUpdateInFieldSet model.CustomField
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &textField)
	require.Nil(t, err)

	require.Equal(t, "new name", textField.CustomFieldUpdateInFieldSet.Name)
	require.Equal(t, "new value", textField.CustomFieldUpdateInFieldSet.Value.RealValue())
	require.Equal(t, "new source", *textField.CustomFieldUpdateInFieldSet.Source)
	require.Equal(t, model.CustomFieldDataTypeText, textField.CustomFieldUpdateInFieldSet.Datatype)
	require.Equal(t, fieldId, textField.CustomFieldUpdateInFieldSet.ID)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "CustomField"))
}

func TestMutationResolver_CustomFieldDeleteFromFieldSetByID(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	fieldSetId := neo4jt.CreateDefaultFieldSet(driver, contactId)
	fieldId := neo4jt.CreateDefaultCustomFieldInSet(driver, fieldSetId)

	rawResponse, err := c.RawPost(getQuery("delete_custom_field_from_field_set"),
		client.Var("contactId", contactId),
		client.Var("fieldSetId", fieldSetId),
		client.Var("fieldId", fieldId))
	assertRawResponseSuccess(t, rawResponse, err)

	var textField struct {
		CustomFieldDeleteFromFieldSetByID model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &textField)
	require.Nil(t, err)

	require.Equal(t, true, textField.CustomFieldDeleteFromFieldSetByID.Result)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "CustomField"))
}

func TestMutationResolver_FieldSetDeleteFromContact(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	fieldSetId := neo4jt.CreateDefaultFieldSet(driver, contactId)
	neo4jt.CreateDefaultCustomFieldInSet(driver, fieldSetId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "CustomField"))

	rawResponse, err := c.RawPost(getQuery("delete_field_set_from_contact"),
		client.Var("contactId", contactId),
		client.Var("fieldSetId", fieldSetId))
	assertRawResponseSuccess(t, rawResponse, err)

	var textField struct {
		FieldSetDeleteFromContact model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &textField)
	require.Nil(t, err)

	require.Equal(t, true, textField.FieldSetDeleteFromContact.Result)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "CustomField"))
}

func TestMutationResolver_EntityDefinitionCreate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, "other")

	rawResponse, err := c.RawPost(getQuery("create_entity_definition"))
	assertRawResponseSuccess(t, rawResponse, err)

	var entityDefinition struct {
		EntityDefinitionCreate model.EntityDefinition
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &entityDefinition)
	actual := entityDefinition.EntityDefinitionCreate
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

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "EntityDefinition"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "FieldSetDefinition"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "CustomFieldDefinition"))
}

func TestQueryResolver_EntityDefinitions_FilterExtendsProperty(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	neo4jt.CreateEntityDefinition(driver, tenantName, "")
	id2 := neo4jt.CreateEntityDefinition(driver, tenantName, model.EntityDefinitionExtensionContact.String())
	id3 := neo4jt.CreateEntityDefinition(driver, tenantName, model.EntityDefinitionExtensionContact.String())

	rawResponse, err := c.RawPost(getQuery("get_entity_definitions_filter_by_extends"),
		client.Var("extends", model.EntityDefinitionExtensionContact.String()))
	assertRawResponseSuccess(t, rawResponse, err)

	var entityDefinition struct {
		EntityDefinitions []model.EntityDefinition
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &entityDefinition)
	require.Nil(t, err)
	require.NotNil(t, entityDefinition.EntityDefinitions)
	require.Equal(t, 2, len(entityDefinition.EntityDefinitions))
	require.Equal(t, "CONTACT", entityDefinition.EntityDefinitions[0].Extends.String())
	require.Equal(t, "CONTACT", entityDefinition.EntityDefinitions[1].Extends.String())
	require.ElementsMatch(t, []string{id2, id3}, []string{entityDefinition.EntityDefinitions[0].ID, entityDefinition.EntityDefinitions[1].ID})

	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "EntityDefinition"))
}

func TestMutationResolver_ContactTypeCreate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, "otherTenantName")

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

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "ContactType"))
}

func TestMutationResolver_ContactTypeUpdate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactTypeId := neo4jt.CreateContactType(driver, tenantName, "original type")

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

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "ContactType"))
}

func TestMutationResolver_ContactTypeDelete(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactTypeId := neo4jt.CreateContactType(driver, tenantName, "the type")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "ContactType"))

	rawResponse, err := c.RawPost(getQuery("delete_contact_type"),
		client.Var("contactTypeId", contactTypeId))
	assertRawResponseSuccess(t, rawResponse, err)

	var result struct {
		ContactType_Delete model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &result)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, true, result.ContactType_Delete.Result)

	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "ContactType"))
}

func TestQueryResolver_ContactTypes(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, "other")
	contactTypeId1 := neo4jt.CreateContactType(driver, tenantName, "first")
	contactTypeId2 := neo4jt.CreateContactType(driver, tenantName, "second")
	neo4jt.CreateContactType(driver, "other", "contact type for other tenant")

	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "ContactType"))

	rawResponse, err := c.RawPost(getQuery("get_contact_types"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactType struct {
		ContactTypes []model.ContactType
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactType)
	require.Nil(t, err)
	require.NotNil(t, contactType)
	require.Equal(t, 2, len(contactType.ContactTypes))
	require.Equal(t, contactTypeId1, contactType.ContactTypes[0].ID)
	require.Equal(t, "first", contactType.ContactTypes[0].Name)
	require.Equal(t, contactTypeId2, contactType.ContactTypes[1].ID)
	require.Equal(t, "second", contactType.ContactTypes[1].Name)
}
