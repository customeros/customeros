package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_CustomFieldsMergeAndUpdateInContact(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)
	entityDefinitionId := createEntityDefinition(driver, tenantName, model.EntityDefinitionExtensionContact.String())
	fieldDefinitionId := addFieldDefinitionToEntity(driver, entityDefinitionId)
	setDefinitionId := addSetDefinitionToEntity(driver, entityDefinitionId)
	fieldInSetDefinitionId := addFieldDefinitionToSet(driver, setDefinitionId)
	linkEntityDefinitionToContact(driver, entityDefinitionId, contactId)
	fieldInContactId := createDefaultCustomFieldInContact(driver, contactId)
	fieldSetId := createDefaultFieldSet(driver, contactId)
	fieldInSetId := createDefaultCustomFieldInSet(driver, fieldSetId)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, getCountOfNodes(driver, "CustomField"))
	require.Equal(t, 1, getCountOfNodes(driver, "FieldSet"))

	require.Equal(t, 1, getCountOfNodes(driver, "EntityDefinition"))
	require.Equal(t, 2, getCountOfNodes(driver, "CustomFieldDefinition"))
	require.Equal(t, 1, getCountOfNodes(driver, "FieldSetDefinition"))

	require.Equal(t, 1, getCountOfRelationships(driver, "IS_DEFINED_BY"))

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

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 4, getCountOfNodes(driver, "CustomField"))
	require.Equal(t, 2, getCountOfNodes(driver, "FieldSet"))

	require.Equal(t, 1, getCountOfNodes(driver, "EntityDefinition"))
	require.Equal(t, 2, getCountOfNodes(driver, "CustomFieldDefinition"))
	require.Equal(t, 1, getCountOfNodes(driver, "FieldSetDefinition"))

	require.Equal(t, 4, getCountOfRelationships(driver, "IS_DEFINED_BY"))

	var contact struct {
		CustomFieldsMergeAndUpdateInContact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)

	updatedContact := contact.CustomFieldsMergeAndUpdateInContact
	require.Equal(t, entityDefinitionId, updatedContact.Definition.ID)
	require.Equal(t, 2, len(updatedContact.CustomFields))
	require.Equal(t, fieldInContactId, updatedContact.CustomFields[0].ID)
	require.Equal(t, "field1", updatedContact.CustomFields[0].Name)
	require.Equal(t, "value1", updatedContact.CustomFields[0].Value.RealValue())
	require.Equal(t, "field2", updatedContact.CustomFields[1].Name)
	require.Equal(t, "value2", updatedContact.CustomFields[1].Value.RealValue())
	require.Equal(t, fieldDefinitionId, updatedContact.CustomFields[1].Definition.ID)

	require.Equal(t, 2, len(updatedContact.FieldSets))
	require.Equal(t, fieldSetId, updatedContact.FieldSets[0].ID)
	require.Equal(t, "set1", updatedContact.FieldSets[0].Name)
	require.Equal(t, "set2", updatedContact.FieldSets[1].Name)
	require.Equal(t, setDefinitionId, updatedContact.FieldSets[1].Definition.ID)

	require.Equal(t, fieldInSetId, updatedContact.FieldSets[0].CustomFields[0].ID)
	require.Equal(t, "field3", updatedContact.FieldSets[0].CustomFields[0].Name)
	require.Equal(t, "value3", updatedContact.FieldSets[0].CustomFields[0].Value.RealValue())

	require.Equal(t, "field4", updatedContact.FieldSets[1].CustomFields[0].Name)
	require.Equal(t, "value4", updatedContact.FieldSets[1].CustomFields[0].Value.RealValue())
	require.Equal(t, fieldInSetDefinitionId, updatedContact.FieldSets[1].CustomFields[0].Definition.ID)
}
