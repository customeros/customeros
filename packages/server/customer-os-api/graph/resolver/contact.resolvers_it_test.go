package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

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
		Contact_ByEmail model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, contactId1, contact.Contact_ByEmail.ID)
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
		Contact_ByPhone model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, contactId1, contact.Contact_ByPhone.ID)
}

func TestMutationResolver_ContactCreate(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	createTenant(driver, "otherTenant")
	contactTypeId := createContactType(driver, tenantName, "CUSTOMER")

	rawResponse, err := c.RawPost(getQuery("create_contact"),
		client.Var("contactTypeId", contactTypeId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact_Create model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "MR", contact.Contact_Create.Title.String())
	require.Equal(t, "first", contact.Contact_Create.FirstName)
	require.Equal(t, "last", contact.Contact_Create.LastName)
	require.Equal(t, contactTypeId, contact.Contact_Create.ContactType.ID)
	require.Equal(t, "CUSTOMER", contact.Contact_Create.ContactType.Name)
	require.Equal(t, "Some notes...", *contact.Contact_Create.Notes)
	require.Equal(t, "Some label", *contact.Contact_Create.Label)

	require.Equal(t, 5, len(contact.Contact_Create.CustomFields))

	boolField := contact.Contact_Create.CustomFields[0]
	require.NotNil(t, boolField.GetID())
	require.Equal(t, "boolField", boolField.Name)
	require.Equal(t, model.CustomFieldDataTypeBool, boolField.Datatype)
	require.Equal(t, true, boolField.Value.RealValue())

	decimalField := contact.Contact_Create.CustomFields[1]
	require.NotNil(t, decimalField.GetID())
	require.Equal(t, "decimalField", decimalField.Name)
	require.Equal(t, model.CustomFieldDataTypeDecimal, decimalField.Datatype)
	require.Equal(t, 0.001, decimalField.Value.RealValue())

	integerField := contact.Contact_Create.CustomFields[2]
	require.NotNil(t, integerField.GetID())
	require.Equal(t, "integerField", integerField.Name)
	require.Equal(t, model.CustomFieldDataTypeInteger, integerField.Datatype)
	// issue in decoding, int converted to float 64
	require.Equal(t, float64(123), integerField.Value.RealValue())

	textField := contact.Contact_Create.CustomFields[3]
	require.NotNil(t, textField.GetID())
	require.Equal(t, "textField", textField.Name)
	require.Equal(t, model.CustomFieldDataTypeText, textField.Datatype)
	require.Equal(t, "value1", textField.Value.RealValue())

	timeField := contact.Contact_Create.CustomFields[4]
	require.NotNil(t, timeField.GetID())
	require.Equal(t, "timeField", timeField.Name)
	require.Equal(t, model.CustomFieldDataTypeDatetime, timeField.Datatype)
	require.Equal(t, "2022-11-13T20:21:56.732Z", timeField.Value.RealValue())

	require.Equal(t, 1, len(contact.Contact_Create.Emails))
	require.NotNil(t, contact.Contact_Create.Emails[0].ID)
	require.Equal(t, "contact@abc.com", contact.Contact_Create.Emails[0].Email)
	require.Equal(t, "WORK", contact.Contact_Create.Emails[0].Label.String())
	require.Equal(t, false, contact.Contact_Create.Emails[0].Primary)

	require.Equal(t, 1, len(contact.Contact_Create.PhoneNumbers))
	require.NotNil(t, contact.Contact_Create.PhoneNumbers[0].ID)
	require.Equal(t, "+1234567890", contact.Contact_Create.PhoneNumbers[0].E164)
	require.Equal(t, "MOBILE", contact.Contact_Create.PhoneNumbers[0].Label.String())
	require.Equal(t, true, contact.Contact_Create.PhoneNumbers[0].Primary)

	require.Equal(t, 0, len(contact.Contact_Create.Groups))

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
	require.Equal(t, 1, getCountOfNodes(driver, "ContactType"))
}

func TestMutationResolver_ContactCreate_WithEntityDefinition(t *testing.T) {
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
		Contact_Create model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)

	require.Equal(t, entityDefinitionId, contact.Contact_Create.Definition.ID)
	require.Equal(t, 2, len(contact.Contact_Create.CustomFields))
	require.Equal(t, "field1", contact.Contact_Create.CustomFields[0].Name)
	require.Equal(t, "TEXT", contact.Contact_Create.CustomFields[0].Datatype.String())
	require.Equal(t, "value1", contact.Contact_Create.CustomFields[0].Value.RealValue())
	require.Equal(t, fieldDefinitionId, contact.Contact_Create.CustomFields[0].Definition.ID)
	require.NotNil(t, contact.Contact_Create.CustomFields[0].GetID())
	require.Equal(t, "field2", contact.Contact_Create.CustomFields[1].Name)
	require.Equal(t, "TEXT", contact.Contact_Create.CustomFields[1].Datatype.String())
	require.Equal(t, "value2", contact.Contact_Create.CustomFields[1].Value.RealValue())
	require.NotNil(t, contact.Contact_Create.CustomFields[1].GetID())

}

func TestMutationResolver_ContactCreate_WithOwner(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	userId := createUser(driver, tenantName, entity.UserEntity{
		FirstName: "Agent",
		LastName:  "Smith",
	})

	rawResponse, err := c.RawPost(getQuery("create_contact_with_owner"),
		client.Var("ownerId", userId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact_Create model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "MR", contact.Contact_Create.Title.String())
	require.Equal(t, "first", contact.Contact_Create.FirstName)
	require.Equal(t, "last", contact.Contact_Create.LastName)
	require.Equal(t, userId, contact.Contact_Create.Owner.ID)
	require.Equal(t, "Agent", contact.Contact_Create.Owner.FirstName)
	require.Equal(t, "Smith", contact.Contact_Create.Owner.LastName)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, getCountOfNodes(driver, "User"))
	require.Equal(t, 1, getCountOfNodes(driver, "Tenant"))
	require.Equal(t, 1, getCountOfRelationships(driver, "OWNS"))
}

func TestMutationResolver_UpdateContact(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	origOwnerId := createDefaultUser(driver, tenantName)
	newOwnerId := createDefaultUser(driver, tenantName)
	contactId := createContact(driver, tenantName, entity.ContactEntity{
		Title:     model.PersonTitleMr.String(),
		FirstName: "first",
		LastName:  "last",
		Label:     "label",
		Notes:     "notes",
	})
	contactTypeIdOrig := createContactType(driver, tenantName, "ORIG")
	contactTypeIdUpdate := createContactType(driver, tenantName, "UPDATED")

	setContactTypeForContact(driver, contactId, contactTypeIdOrig)
	userOwnsContact(driver, origOwnerId, contactId)

	rawResponse, err := c.RawPost(getQuery("update_contact"),
		client.Var("contactId", contactId),
		client.Var("contactTypeId", contactTypeIdUpdate),
		client.Var("ownerId", newOwnerId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact_Update model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "DR", contact.Contact_Update.Title.String())
	require.Equal(t, "updated first", contact.Contact_Update.FirstName)
	require.Equal(t, "updated last", contact.Contact_Update.LastName)
	require.Equal(t, "updated notes", *contact.Contact_Update.Notes)
	require.Equal(t, "updated label", *contact.Contact_Update.Label)
	require.Equal(t, contactTypeIdUpdate, contact.Contact_Update.ContactType.ID)
	require.Equal(t, "UPDATED", contact.Contact_Update.ContactType.Name)
	require.Equal(t, newOwnerId, contact.Contact_Update.Owner.ID)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, getCountOfNodes(driver, "ContactType"))
	require.Equal(t, 2, getCountOfNodes(driver, "User"))
	require.Equal(t, 1, getCountOfRelationships(driver, "IS_OF_TYPE"))
	require.Equal(t, 1, getCountOfRelationships(driver, "OWNS"))
}

func TestQueryResolver_Contact(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)
	companyId1 := createCompany(driver, tenantName, "ABC")
	companyId2 := createCompany(driver, tenantName, "XYZ")
	positionId1 := contactWorksForCompany(driver, contactId, companyId1, "CTO")
	positionId2 := contactWorksForCompany(driver, contactId, companyId2, "CEO")

	rawResponse, err := c.RawPost(getQuery("get_contact_by_id"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var searchedContact struct {
		Contact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedContact)
	require.Nil(t, err)
	require.Equal(t, contactId, searchedContact.Contact.ID)

	companyPositions := searchedContact.Contact.CompanyPositions
	require.Equal(t, 2, len(companyPositions))
	require.Equal(t, positionId1, companyPositions[0].ID)
	require.Equal(t, "CTO", *companyPositions[0].JobTitle)
	require.Equal(t, companyId1, companyPositions[0].Company.ID)
	require.Equal(t, "ABC", companyPositions[0].Company.Name)
	require.Equal(t, positionId2, companyPositions[1].ID)
	require.Equal(t, "CEO", *companyPositions[1].JobTitle)
	require.Equal(t, companyId2, companyPositions[1].Company.ID)
	require.Equal(t, "XYZ", companyPositions[1].Company.Name)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, getCountOfNodes(driver, "Company"))
	require.Equal(t, 2, getCountOfRelationships(driver, "WORKS_AT"))
}

func TestQueryResolver_Contacts_SortByTitleAscFirstNameAscLastNameDesc(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)

	contact1 := createContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "contact",
		LastName:  "1",
	})
	contact2 := createContact(driver, tenantName, entity.ContactEntity{
		Title:     "DR",
		FirstName: "contact",
		LastName:  "9",
	})
	contact3 := createContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "contact",
		LastName:  "222",
	})
	contact4 := createContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "other contact",
		LastName:  "A",
	})

	rawResponse, err := c.RawPost(getQuery("get_contacts_with_sorting"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contacts struct {
		Contacts model.ContactGroupPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contacts)
	require.Nil(t, err)
	require.NotNil(t, contacts.Contacts)
	require.Equal(t, 4, len(contacts.Contacts.Content))
	require.Equal(t, contact2, contacts.Contacts.Content[0].ID)
	require.Equal(t, contact3, contacts.Contacts.Content[1].ID)
	require.Equal(t, contact1, contacts.Contacts.Content[2].ID)
	require.Equal(t, contact4, contacts.Contacts.Content[3].ID)

	require.Equal(t, 4, getCountOfNodes(driver, "Contact"))
}

func TestQueryResolver_Contact_BasicFilters_FindContactWithLetterAInName(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)

	contactFoundByFirstName := createContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "aa",
		LastName:  "bb",
	})
	contactFoundByLastName := createContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "bb",
		LastName:  "AA",
	})
	contactFilteredOut := createContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "bb",
		LastName:  "BB",
	})

	require.Equal(t, 3, getCountOfNodes(driver, "Contact"))

	rawResponse, err := c.RawPost(getQuery("get_contacts_basic_filters"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contacts struct {
		Contacts model.ContactsPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contacts)
	require.Nil(t, err)
	require.NotNil(t, contacts.Contacts)
	require.Equal(t, 2, len(contacts.Contacts.Content))
	require.Equal(t, contactFoundByFirstName, contacts.Contacts.Content[0].ID)
	require.Equal(t, contactFoundByLastName, contacts.Contacts.Content[1].ID)
	require.Equal(t, 1, contacts.Contacts.TotalPages)
	require.Equal(t, int64(2), contacts.Contacts.TotalElements)
	// suppress unused warnings
	require.NotNil(t, contactFilteredOut)
}
