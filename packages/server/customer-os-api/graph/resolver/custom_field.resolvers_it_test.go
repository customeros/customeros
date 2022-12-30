package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_CustomFieldsMergeAndUpdateInContact(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	entityDefinitionId := neo4jt.CreateEntityDefinition(driver, tenantName, model.EntityDefinitionExtensionContact.String())
	fieldDefinitionId := neo4jt.AddFieldDefinitionToEntity(driver, entityDefinitionId)
	setDefinitionId := neo4jt.AddSetDefinitionToEntity(driver, entityDefinitionId)
	fieldInSetDefinitionId := neo4jt.AddFieldDefinitionToSet(driver, setDefinitionId)
	neo4jt.LinkEntityDefinitionToContact(driver, entityDefinitionId, contactId)
	fieldInContactId := neo4jt.CreateDefaultCustomFieldInContact(driver, contactId)
	fieldSetId := neo4jt.CreateDefaultFieldSet(driver, contactId)
	fieldInSetId := neo4jt.CreateDefaultCustomFieldInSet(driver, fieldSetId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "CustomField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "FieldSet"))

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "EntityDefinition"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "CustomFieldDefinition"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "FieldSetDefinition"))

	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "IS_DEFINED_BY"))

	rawResponse, err := c.RawPost(getQuery("update_custom_fields_and_filed_sets_in_contact"),
		client.Var("contactId", contactId),
		client.Var("customFieldId", fieldInContactId),
		client.Var("fieldSetId", fieldSetId),
		client.Var("customFieldInSetId", fieldInSetId),
		client.Var("fieldDefinitionId", fieldDefinitionId),
		client.Var("setDefinitionId", setDefinitionId),
		client.Var("fieldInSetDefinitionId", fieldInSetDefinitionId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "CustomField"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "CustomField_"+tenantName))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "FieldSet"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "FieldSet_"+tenantName))

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "EntityDefinition"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "CustomFieldDefinition"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "FieldSetDefinition"))

	require.Equal(t, 4, neo4jt.GetCountOfRelationships(driver, "IS_DEFINED_BY"))

	var contact struct {
		CustomFieldsMergeAndUpdateInContact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)

	updatedContact := contact.CustomFieldsMergeAndUpdateInContact
	require.Equal(t, entityDefinitionId, updatedContact.Definition.ID)
	require.Equal(t, 2, len(updatedContact.CustomFields))
	if updatedContact.CustomFields[0].ID == fieldInContactId {
		checkCustomField(t, *updatedContact.CustomFields[0], "field1", "value1", nil)
		checkCustomField(t, *updatedContact.CustomFields[1], "field2", "value2", &fieldDefinitionId)
	} else {
		checkCustomField(t, *updatedContact.CustomFields[1], "field1", "value1", nil)
		checkCustomField(t, *updatedContact.CustomFields[0], "field2", "value2", &fieldDefinitionId)
	}

	require.Equal(t, 2, len(updatedContact.FieldSets))
	require.ElementsMatch(t, []string{"set1", "set2"}, []string{updatedContact.FieldSets[0].Name, updatedContact.FieldSets[1].Name})

	if updatedContact.FieldSets[0].Definition != nil {
		require.Equal(t, fieldInSetId, updatedContact.FieldSets[1].CustomFields[0].ID)
		checkCustomField(t, *updatedContact.FieldSets[1].CustomFields[0], "field3", "value3", nil)
		checkCustomField(t, *updatedContact.FieldSets[0].CustomFields[0], "field4", "value4", &fieldInSetDefinitionId)
	} else {
		require.Equal(t, fieldInSetId, updatedContact.FieldSets[0].CustomFields[0].ID)
		checkCustomField(t, *updatedContact.FieldSets[0].CustomFields[0], "field3", "value3", nil)
		checkCustomField(t, *updatedContact.FieldSets[1].CustomFields[0], "field4", "value4", &fieldInSetDefinitionId)
	}

	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "EntityDefinition", "CustomFieldDefinition",
		"FieldSetDefinition", "TextField", "CustomField", "CustomField_" + tenantName, "FieldSet", "FieldSet_" + tenantName})
}

func checkCustomField(t *testing.T, customField model.CustomField, name, value string, fieldDefinitionId *string) {
	require.Equal(t, name, customField.Name)
	require.Equal(t, value, customField.Value.RealValue())
	if fieldDefinitionId != nil {
		require.Equal(t, *fieldDefinitionId, customField.Definition.ID)
	}
}
