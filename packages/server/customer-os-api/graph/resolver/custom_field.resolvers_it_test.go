package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_CustomFieldsMergeAndUpdateInContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	entityTemplateId := neo4jt.CreateEntityTemplate(ctx, driver, tenantName, model.EntityTemplateExtensionContact.String())
	fieldTemplateId := neo4jt.AddFieldTemplateToEntity(ctx, driver, entityTemplateId)
	setTemplateId := neo4jt.AddSetTemplateToEntity(ctx, driver, entityTemplateId)
	fieldInSetTemplateId := neo4jt.AddFieldTemplateToSet(ctx, driver, setTemplateId)
	neo4jt.LinkEntityTemplateToContact(ctx, driver, entityTemplateId, contactId)
	fieldInContactId := neo4jt.CreateDefaultCustomFieldInContact(ctx, driver, contactId)
	fieldSetId := neo4jt.CreateDefaultFieldSet(ctx, driver, contactId)
	fieldInSetId := neo4jt.CreateDefaultCustomFieldInSet(ctx, driver, fieldSetId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "CustomField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "FieldSet"))

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "EntityTemplate"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "CustomFieldTemplate"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "FieldSetTemplate"))

	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "IS_DEFINED_BY"))

	rawResponse, err := c.RawPost(getQuery("update_custom_fields_and_filed_sets_in_contact"),
		client.Var("contactId", contactId),
		client.Var("customFieldId", fieldInContactId),
		client.Var("fieldSetId", fieldSetId),
		client.Var("customFieldInSetId", fieldInSetId),
		client.Var("fieldTemplateId", fieldTemplateId),
		client.Var("setTemplateId", setTemplateId),
		client.Var("fieldInSetTemplateId", fieldInSetTemplateId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "CustomField"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "CustomField_"+tenantName))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "FieldSet"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "FieldSet_"+tenantName))

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "EntityTemplate"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "CustomFieldTemplate"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "FieldSetTemplate"))

	require.Equal(t, 4, neo4jt.GetCountOfRelationships(ctx, driver, "IS_DEFINED_BY"))

	var contact struct {
		CustomFieldsMergeAndUpdateInContact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)

	updatedContact := contact.CustomFieldsMergeAndUpdateInContact
	require.Equal(t, entityTemplateId, updatedContact.Template.ID)
	require.Equal(t, 2, len(updatedContact.CustomFields))
	if updatedContact.CustomFields[0].ID == fieldInContactId {
		checkCustomField(t, *updatedContact.CustomFields[0], "field1", "value1", nil)
		checkCustomField(t, *updatedContact.CustomFields[1], "field2", "value2", &fieldTemplateId)
	} else {
		checkCustomField(t, *updatedContact.CustomFields[1], "field1", "value1", nil)
		checkCustomField(t, *updatedContact.CustomFields[0], "field2", "value2", &fieldTemplateId)
	}

	require.Equal(t, 2, len(updatedContact.FieldSets))
	require.Equal(t, model.DataSourceOpenline, updatedContact.FieldSets[0].Source)
	require.Equal(t, model.DataSourceOpenline, updatedContact.FieldSets[1].Source)
	require.ElementsMatch(t, []string{"set1", "set2"}, []string{updatedContact.FieldSets[0].Name, updatedContact.FieldSets[1].Name})

	if updatedContact.FieldSets[0].Template != nil {
		require.Equal(t, fieldInSetId, updatedContact.FieldSets[1].CustomFields[0].ID)
		checkCustomField(t, *updatedContact.FieldSets[1].CustomFields[0], "field3", "value3", nil)
		checkCustomField(t, *updatedContact.FieldSets[0].CustomFields[0], "field4", "value4", &fieldInSetTemplateId)
	} else {
		require.Equal(t, fieldInSetId, updatedContact.FieldSets[0].CustomFields[0].ID)
		checkCustomField(t, *updatedContact.FieldSets[0].CustomFields[0], "field3", "value3", nil)
		checkCustomField(t, *updatedContact.FieldSets[1].CustomFields[0], "field4", "value4", &fieldInSetTemplateId)
	}

	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "EntityTemplate", "CustomFieldTemplate",
		"FieldSetTemplate", "TextField", "CustomField", "CustomField_" + tenantName, "FieldSet", "FieldSet_" + tenantName})
}

func checkCustomField(t *testing.T, customField model.CustomField, name, value string, fieldTemplateId *string) {
	require.Equal(t, name, customField.Name)
	require.Equal(t, value, customField.Value.RealValue())
	require.Equal(t, model.DataSourceOpenline, customField.Source)
	if fieldTemplateId != nil {
		require.Equal(t, *fieldTemplateId, customField.Template.ID)
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

	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "EntityTemplate", "EntityTemplate_" + tenantName,
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
