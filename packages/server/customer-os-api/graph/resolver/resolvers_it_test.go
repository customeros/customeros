package resolver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
	"os"
	"reflect"
	"sort"
	"testing"
)

var (
	neo4jContainer testcontainers.Container
	driver         *neo4j.DriverWithContext

	postgresContainer testcontainers.Container
	postgresGormDB    *gorm.DB
	postgresSqlDB     *sql.DB
	c                 *client.Client
	cAdmin            *client.Client
	cAdminWithTenant  *client.Client
)

const tenantName = "openline"
const testUserId = "test-user-id"
const testContactId = "test-contact-id"

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
			logrus.Fatal("Error during container termination")
		}
	}(postgresContainer, context.Background())

	prepareClient()

	os.Exit(m.Run())
}

func tearDownTestCase(ctx context.Context) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jt.CleanupAllData(ctx, driver)
	}
}

func prepareClient() {
	commonServices := commonService.InitServices(postgresGormDB, driver)
	serviceContainer := service.InitServices(driver, commonServices, nil)
	graphResolver := NewResolver(serviceContainer, nil)
	loader := dataloader.NewDataLoader(serviceContainer)
	customCtx := &common.CustomContext{
		Tenant: tenantName,
		UserId: testUserId,
		Role:   model.RoleUser,
	}

	customAdminCtx := &common.CustomContext{
		Role: model.RoleAdmin,
	}

	customAdminWTenantCtx := &common.CustomContext{
		Tenant: tenantName,
		Role:   model.RoleAdmin,
	}
	schemaConfig := generated.Config{Resolvers: graphResolver}
	schemaConfig.Directives.HasRole = cosHandler.GetRoleChecker()
	schemaConfig.Directives.HasTenant = cosHandler.GetTenantChecker()
	server := handler.NewDefaultServer(generated.NewExecutableSchema(schemaConfig))
	dataloaderServer := dataloader.Middleware(loader, server)
	handler := common.WithContext(customCtx, dataloaderServer)
	c = client.New(handler)
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
		logrus.Errorf("Error in response: %v", string(response.Errors))
	}
	require.NotNil(t, response.Data)
	require.Nil(t, response.Errors)
}

func assertNeo4jLabels(ctx context.Context, t *testing.T, driver *neo4j.DriverWithContext, expectedLabels []string) {
	actualLabels := neo4jt.GetAllLabels(ctx, driver)
	sort.Strings(expectedLabels)
	sort.Strings(actualLabels)
	if !reflect.DeepEqual(actualLabels, expectedLabels) {
		t.Errorf("Expected labels: %v, \nActual labels: %v", expectedLabels, actualLabels)
	}
}

func TestMutationResolver_FieldSetMergeToContact_AllowMultipleFieldSetWithSameNameOnDifferentContacts(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId1 := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		Prefix:    "MR",
		FirstName: "first",
		LastName:  "last",
	})
	contactId2 := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		Prefix:    "MR",
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
	require.NotNil(t, fieldSet1.FieldSetMergeToContact.CreatedAt)
	require.Equal(t, "some name", fieldSet2.FieldSetMergeToContact.Name)
	require.NotNil(t, fieldSet2.FieldSetMergeToContact.CreatedAt)
	require.Equal(t, model.DataSourceOpenline, fieldSet1.FieldSetMergeToContact.Source)
	require.Equal(t, model.DataSourceOpenline, fieldSet2.FieldSetMergeToContact.Source)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "FieldSet"))
}

func TestMutationResolver_MergeCustomFieldToFieldSet(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	fieldSetId := neo4jt.CreateDefaultFieldSet(ctx, driver, contactId)

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
	require.Equal(t, model.DataSourceOpenline, textField.CustomFieldMergeToFieldSet.Source)
	require.Equal(t, model.CustomFieldDataTypeText, textField.CustomFieldMergeToFieldSet.Datatype)
	require.NotNil(t, textField.CustomFieldMergeToFieldSet.ID)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "FieldSet"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "CustomField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "TextField"))
}

func TestMutationResolver_CustomFieldUpdateInFieldSet(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	fieldSetId := neo4jt.CreateDefaultFieldSet(ctx, driver, contactId)
	fieldId := neo4jt.CreateDefaultCustomFieldInSet(ctx, driver, fieldSetId)

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
	require.Equal(t, model.DataSourceOpenline, textField.CustomFieldUpdateInFieldSet.Source)
	require.Equal(t, model.CustomFieldDataTypeText, textField.CustomFieldUpdateInFieldSet.Datatype)
	require.Equal(t, fieldId, textField.CustomFieldUpdateInFieldSet.ID)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "FieldSet"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "CustomField"))
}

func TestMutationResolver_CustomFieldDeleteFromFieldSetByID(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	fieldSetId := neo4jt.CreateDefaultFieldSet(ctx, driver, contactId)
	fieldId := neo4jt.CreateDefaultCustomFieldInSet(ctx, driver, fieldSetId)

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

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "FieldSet"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "CustomField"))
}

func TestMutationResolver_FieldSetDeleteFromContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	fieldSetId := neo4jt.CreateDefaultFieldSet(ctx, driver, contactId)
	neo4jt.CreateDefaultCustomFieldInSet(ctx, driver, fieldSetId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "FieldSet"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "CustomField"))

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

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "FieldSet"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "CustomField"))
}

func TestMutationResolver_EntityTemplateCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, "other")

	rawResponse, err := c.RawPost(getQuery("create_entity_template"))
	assertRawResponseSuccess(t, rawResponse, err)

	var entityTemplate struct {
		EntityTemplateCreate model.EntityTemplate
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &entityTemplate)
	actual := entityTemplate.EntityTemplateCreate
	require.Nil(t, err)
	require.NotNil(t, actual)
	require.NotNil(t, actual.ID)
	require.NotNil(t, actual.CreatedAt)
	require.Equal(t, "the entity template name", actual.Name)
	require.Equal(t, 1, actual.Version)
	require.Nil(t, actual.Extends)

	require.Equal(t, 2, len(actual.FieldSetTemplates))

	set := actual.FieldSetTemplates[0]
	require.NotNil(t, set.ID)
	require.NotNil(t, set.CreatedAt)
	require.Equal(t, "set 1", set.Name)
	require.Equal(t, 1, set.Order)
	require.Equal(t, 2, len(set.CustomFieldTemplates))

	field := set.CustomFieldTemplates[0]
	require.NotNil(t, field)
	require.NotNil(t, field.CreatedAt)
	require.Equal(t, "field 3", field.Name)
	require.Equal(t, 1, field.Order)
	require.Equal(t, true, field.Mandatory)
	require.Equal(t, model.CustomFieldTemplateTypeText, field.Type)
	require.Nil(t, field.Min)
	require.Nil(t, field.Max)
	require.Nil(t, field.Length)

	field = set.CustomFieldTemplates[1]
	require.NotNil(t, field)
	require.NotNil(t, field.CreatedAt)
	require.Equal(t, "field 4", field.Name)
	require.Equal(t, 2, field.Order)
	require.Equal(t, false, field.Mandatory)
	require.Equal(t, model.CustomFieldTemplateTypeText, field.Type)
	require.Equal(t, 10, *field.Min)
	require.Equal(t, 990, *field.Max)
	require.Equal(t, 2550, *field.Length)

	set = actual.FieldSetTemplates[1]
	require.NotNil(t, set.ID)
	require.NotNil(t, set.CreatedAt)
	require.Equal(t, "set 2", set.Name)
	require.Equal(t, 2, set.Order)
	require.Equal(t, 0, len(set.CustomFieldTemplates))

	field = actual.CustomFieldTemplates[0]
	require.NotNil(t, field)
	require.NotNil(t, field.CreatedAt)
	require.Equal(t, "field 1", field.Name)
	require.Equal(t, 1, field.Order)
	require.Equal(t, true, field.Mandatory)
	require.Equal(t, model.CustomFieldTemplateTypeText, field.Type)
	require.Nil(t, field.Min)
	require.Nil(t, field.Max)
	require.Nil(t, field.Length)

	field = actual.CustomFieldTemplates[1]
	require.NotNil(t, field)
	require.NotNil(t, field.CreatedAt)
	require.Equal(t, "field 2", field.Name)
	require.Equal(t, 2, field.Order)
	require.Equal(t, false, field.Mandatory)
	require.Equal(t, model.CustomFieldTemplateTypeText, field.Type)
	require.Equal(t, 1, *field.Min)
	require.Equal(t, 99, *field.Max)
	require.Equal(t, 255, *field.Length)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "EntityTemplate"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "EntityTemplate_"+tenantName))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "FieldSetTemplate"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "FieldSetTemplate_"+tenantName))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "CustomFieldTemplate"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "CustomFieldTemplate_"+tenantName))

	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "EntityTemplate", "EntityTemplate_" + tenantName,
		"FieldSetTemplate", "FieldSetTemplate_" + tenantName, "CustomFieldTemplate", "CustomFieldTemplate_" + tenantName})
}

func TestQueryResolver_EntityTemplates_FilterExtendsProperty(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateEntityTemplate(ctx, driver, tenantName, "")
	id2 := neo4jt.CreateEntityTemplate(ctx, driver, tenantName, model.EntityTemplateExtensionContact.String())
	id3 := neo4jt.CreateEntityTemplate(ctx, driver, tenantName, model.EntityTemplateExtensionContact.String())

	rawResponse, err := c.RawPost(getQuery("get_entity_templates_filter_by_extends"),
		client.Var("extends", model.EntityTemplateExtensionContact.String()))
	assertRawResponseSuccess(t, rawResponse, err)

	var entityTemplate struct {
		EntityTemplates []model.EntityTemplate
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &entityTemplate)
	require.Nil(t, err)
	require.NotNil(t, entityTemplate.EntityTemplates)
	require.Equal(t, 2, len(entityTemplate.EntityTemplates))
	require.Equal(t, "CONTACT", entityTemplate.EntityTemplates[0].Extends.String())
	require.Equal(t, "CONTACT", entityTemplate.EntityTemplates[1].Extends.String())
	require.ElementsMatch(t, []string{id2, id3}, []string{entityTemplate.EntityTemplates[0].ID, entityTemplate.EntityTemplates[1].ID})

	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "EntityTemplate"))
}
