package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryResolver_ContactByEmail(t *testing.T) {
	defer tearDownTestCase()(t)
	otherTenant := "other"
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, otherTenant)
	contactId1 := neo4jt.CreateDefaultContact(driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(driver, otherTenant)
	neo4jt.AddEmailToContact(driver, contactId1, "test@test.com", true, "MAIN")
	neo4jt.AddEmailToContact(driver, contactId2, "test@test.com", true, "MAIN")

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
	defer tearDownTestCase()(t)
	otherTenant := "other"
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, otherTenant)
	contactId1 := neo4jt.CreateDefaultContact(driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(driver, otherTenant)
	neo4jt.AddPhoneNumberToContact(driver, contactId1, "+1234567890", false, "OTHER")
	neo4jt.AddPhoneNumberToContact(driver, contactId2, "+1234567890", true, "MAIN")

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

func TestMutationResolver_ContactCreate_Min(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("create_contact_min"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact_Create model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "", contact.Contact_Create.Title.String())
	require.Equal(t, "", *contact.Contact_Create.FirstName)
	require.Equal(t, "", *contact.Contact_Create.LastName)
	require.Equal(t, "", *contact.Contact_Create.Notes)
	require.Equal(t, "", *contact.Contact_Create.Label)
	require.Equal(t, false, contact.Contact_Create.Readonly)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetTotalCountOfNodes(driver))
}

func TestMutationResolver_ContactCreate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, "otherTenant")
	contactTypeId := neo4jt.CreateContactType(driver, tenantName, "CUSTOMER")

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
	require.Equal(t, "first", *contact.Contact_Create.FirstName)
	require.Equal(t, "last", *contact.Contact_Create.LastName)
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

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "ContactGroup"))
	require.Equal(t, 5, neo4jt.GetCountOfNodes(driver, "CustomField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "TextField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "IntField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "FloatField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "BoolField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "TimeField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "ContactType"))
}

func TestMutationResolver_ContactCreate_WithCustomFields(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	entityDefinitionId := neo4jt.CreateEntityDefinition(driver, tenantName, model.EntityDefinitionExtensionContact.String())
	fieldDefinitionId := neo4jt.AddFieldDefinitionToEntity(driver, entityDefinitionId)
	setDefinitionId := neo4jt.AddSetDefinitionToEntity(driver, entityDefinitionId)
	fieldInSetDefinitionId := neo4jt.AddFieldDefinitionToSet(driver, setDefinitionId)

	rawResponse, err := c.RawPost(getQuery("create_contact_with_custom_fields"),
		client.Var("entityDefinitionId", entityDefinitionId),
		client.Var("fieldDefinitionId", fieldDefinitionId),
		client.Var("setDefinitionId", setDefinitionId),
		client.Var("fieldInSetDefinitionId", fieldInSetDefinitionId))
	assertRawResponseSuccess(t, rawResponse, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "ContactGroup"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Company"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "CustomField"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "TextField"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Email"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "EntityDefinition"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "CustomFieldDefinition"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "FieldSetDefinition"))

	var contact struct {
		Contact_Create model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)

	createdContact := contact.Contact_Create
	require.Equal(t, entityDefinitionId, createdContact.Definition.ID)
	require.Equal(t, 2, len(createdContact.CustomFields))
	require.Equal(t, "field1", createdContact.CustomFields[0].Name)
	require.Equal(t, "TEXT", createdContact.CustomFields[0].Datatype.String())
	require.Equal(t, "value1", createdContact.CustomFields[0].Value.RealValue())
	require.Equal(t, "", *createdContact.CustomFields[0].Source)
	require.Equal(t, fieldDefinitionId, createdContact.CustomFields[0].Definition.ID)
	require.NotNil(t, createdContact.CustomFields[0].GetID())
	require.Equal(t, "field2", createdContact.CustomFields[1].Name)
	require.Equal(t, "TEXT", createdContact.CustomFields[1].Datatype.String())
	require.Equal(t, "value2", createdContact.CustomFields[1].Value.RealValue())
	require.Equal(t, "hubspot", *createdContact.CustomFields[1].Source)
	require.NotNil(t, createdContact.CustomFields[1].GetID())
	require.Equal(t, 2, len(createdContact.FieldSets))
	require.NotNil(t, createdContact.FieldSets[0].ID)
	require.NotNil(t, createdContact.FieldSets[0].Added)
	require.Equal(t, "set1", createdContact.FieldSets[0].Name)
	require.Equal(t, 2, len(createdContact.FieldSets[0].CustomFields))
	require.Equal(t, "field3InSet", createdContact.FieldSets[0].CustomFields[0].Name)
	require.Equal(t, "value3", createdContact.FieldSets[0].CustomFields[0].Value.RealValue())
	require.Equal(t, "", *createdContact.FieldSets[0].CustomFields[0].Source)
	require.Equal(t, "TEXT", createdContact.FieldSets[0].CustomFields[0].Datatype.String())
	require.Equal(t, fieldInSetDefinitionId, createdContact.FieldSets[0].CustomFields[0].Definition.ID)
	require.Equal(t, "field4InSet", createdContact.FieldSets[0].CustomFields[1].Name)
	require.Equal(t, "value4", createdContact.FieldSets[0].CustomFields[1].Value.RealValue())
	require.Equal(t, "zendesk", *createdContact.FieldSets[0].CustomFields[1].Source)
	require.Equal(t, "TEXT", createdContact.FieldSets[0].CustomFields[1].Datatype.String())
	require.Nil(t, createdContact.FieldSets[0].CustomFields[1].Definition)
	require.NotNil(t, createdContact.FieldSets[1].ID)
	require.NotNil(t, createdContact.FieldSets[1].Added)
	require.Equal(t, "set2", createdContact.FieldSets[1].Name)
}

func TestMutationResolver_ContactCreate_WithOwner(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	userId := neo4jt.CreateUser(driver, tenantName, entity.UserEntity{
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
	require.Equal(t, "", contact.Contact_Create.Title.String())
	require.Equal(t, "first", *contact.Contact_Create.FirstName)
	require.Equal(t, "last", *contact.Contact_Create.LastName)
	require.Equal(t, userId, contact.Contact_Create.Owner.ID)
	require.Equal(t, "Agent", contact.Contact_Create.Owner.FirstName)
	require.Equal(t, "Smith", contact.Contact_Create.Owner.LastName)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "OWNS"))
}

func TestMutationResolver_ContactCreate_WithExternalReference(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateHubspotExternalSystem(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("create_contact_with_external_reference"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact_Create model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.NotNil(t, contact.Contact_Create.ID)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "ExternalSystem"))
	require.Equal(t, 3, neo4jt.GetTotalCountOfNodes(driver))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "IS_LINKED_WITH"))
}

func TestMutationResolver_UpdateContact(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	origOwnerId := neo4jt.CreateDefaultUser(driver, tenantName)
	newOwnerId := neo4jt.CreateDefaultUser(driver, tenantName)
	contactId := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     model.PersonTitleMr.String(),
		FirstName: "first",
		LastName:  "last",
		Label:     "label",
		Notes:     "notes",
	})
	contactTypeIdOrig := neo4jt.CreateContactType(driver, tenantName, "ORIG")
	contactTypeIdUpdate := neo4jt.CreateContactType(driver, tenantName, "UPDATED")

	neo4jt.SetContactTypeForContact(driver, contactId, contactTypeIdOrig)
	neo4jt.UserOwnsContact(driver, origOwnerId, contactId)

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
	require.Equal(t, "updated first", *contact.Contact_Update.FirstName)
	require.Equal(t, "updated last", *contact.Contact_Update.LastName)
	require.Equal(t, "updated notes", *contact.Contact_Update.Notes)
	require.Equal(t, "updated label", *contact.Contact_Update.Label)
	require.Equal(t, contactTypeIdUpdate, contact.Contact_Update.ContactType.ID)
	require.Equal(t, "UPDATED", contact.Contact_Update.ContactType.Name)
	require.Equal(t, newOwnerId, contact.Contact_Update.Owner.ID)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "ContactType"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "IS_OF_TYPE"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "OWNS"))
}

func TestQueryResolver_Contact(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	companyId1 := neo4jt.CreateCompany(driver, tenantName, "ABC")
	companyId2 := neo4jt.CreateCompany(driver, tenantName, "XYZ")
	positionId1 := neo4jt.ContactWorksForCompany(driver, contactId, companyId1, "CTO")
	positionId2 := neo4jt.ContactWorksForCompany(driver, contactId, companyId2, "CEO")

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

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Company"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "WORKS_AT"))
}

func TestQueryResolver_Contacts_SortByTitleAscFirstNameAscLastNameDesc(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	contact1 := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "contact",
		LastName:  "1",
	})
	contact2 := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     "DR",
		FirstName: "contact",
		LastName:  "9",
	})
	contact3 := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     "",
		FirstName: "contact",
		LastName:  "222",
	})
	contact4 := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
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
	require.Equal(t, contact3, contacts.Contacts.Content[0].ID)
	require.Equal(t, contact2, contacts.Contacts.Content[1].ID)
	require.Equal(t, contact1, contacts.Contacts.Content[2].ID)
	require.Equal(t, contact4, contacts.Contacts.Content[3].ID)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "Contact"))
}

func TestQueryResolver_Contact_BasicFilters_FindContactWithLetterAInName(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	contactFoundByFirstName := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "aa",
		LastName:  "bb",
	})
	contactFoundByLastName := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "bb",
		LastName:  "AA",
	})
	contactFilteredOut := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "bb",
		LastName:  "BB",
	})

	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "Contact"))

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

func TestQueryResolver_Contact_WithConversations(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	user1 := neo4jt.CreateDefaultUser(driver, tenantName)
	user2 := neo4jt.CreateDefaultUser(driver, tenantName)
	contact1 := neo4jt.CreateDefaultContact(driver, tenantName)
	contact2 := neo4jt.CreateDefaultContact(driver, tenantName)
	contact3 := neo4jt.CreateDefaultContact(driver, tenantName)

	conv1_1 := neo4jt.CreateConversation(driver, user1, contact1)
	conv1_2 := neo4jt.CreateConversation(driver, user1, contact2)
	conv2_1 := neo4jt.CreateConversation(driver, user2, contact1)
	conv2_3 := neo4jt.CreateConversation(driver, user2, contact3)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "User"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "Conversation"))

	rawResponse, err := c.RawPost(getQuery("get_contact_with_conversations"),
		client.Var("contactId", contact1))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, contact1, contact.Contact.ID)
	require.Equal(t, 1, contact.Contact.Conversations.TotalPages)
	require.Equal(t, int64(2), contact.Contact.Conversations.TotalElements)
	require.Equal(t, 2, len(contact.Contact.Conversations.Content))
	conversations := contact.Contact.Conversations.Content
	require.ElementsMatch(t, []string{conv1_1, conv2_1}, []string{conversations[0].ID, conversations[1].ID})
	require.ElementsMatch(t, []string{user1, user2}, []string{conversations[0].User.ID, conversations[1].User.ID})
	require.ElementsMatch(t, []string{user1, user2}, []string{conversations[0].UserID, conversations[1].UserID})
	require.Equal(t, contact1, conversations[0].Contact.ID)
	require.Equal(t, contact1, conversations[1].Contact.ID)
	require.Equal(t, contact1, conversations[0].ContactID)
	require.Equal(t, contact1, conversations[1].ContactID)

	require.NotNil(t, conv1_2)
	require.NotNil(t, conv2_3)
}

func TestQueryResolver_Contact_WithActions(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(driver, tenantName)
	userId := neo4jt.CreateDefaultUser(driver, tenantName)
	conversationId := neo4jt.CreateConversation(driver, userId, contactId)

	now := time.Now().UTC()
	secAgo1 := now.Add(time.Duration(-1) * time.Second)
	secAgo30 := now.Add(time.Duration(-30) * time.Second)
	secAgo60 := now.Add(time.Duration(-60) * time.Second)
	from := now.Add(time.Duration(-10) * time.Minute)

	messageId := neo4jt.AddMessageToConversation(driver, conversationId, mapper.MapMessageChannelFromModel(model.MessageChannelChat), secAgo60)

	pageViewId1 := neo4jt.CreatePageView(driver, contactId, entity.PageViewEntity{
		StartedAt:      secAgo1,
		EndedAt:        now,
		TrackerName:    "tracker1",
		SessionId:      "session1",
		Application:    "application1",
		PageTitle:      "page1",
		PageUrl:        "http://app-1.ai",
		OrderInSession: 1,
		EngagedTime:    10,
	})

	pageViewId2 := neo4jt.CreatePageView(driver, contactId, entity.PageViewEntity{
		StartedAt:      secAgo30,
		EndedAt:        now,
		TrackerName:    "tracker2",
		SessionId:      "session2",
		Application:    "application2",
		PageTitle:      "page2",
		PageUrl:        "http://app-2.ai",
		OrderInSession: 2,
		EngagedTime:    20,
	})

	neo4jt.CreatePageView(driver, contactId2, entity.PageViewEntity{})

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "User"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "Action"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "PageView"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Conversation"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Message"))

	rawResponse, err := c.RawPost(getQuery("get_contact_with_actions"),
		client.Var("contactId", contactId),
		client.Var("from", from),
		client.Var("to", now))
	assertRawResponseSuccess(t, rawResponse, err)

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])

	actions := contact.(map[string]interface{})["actions"].([]interface{})
	require.Equal(t, 3, len(actions))
	action1 := actions[0].(map[string]interface{})
	require.Equal(t, "PageViewAction", action1["__typename"].(string))
	require.Equal(t, pageViewId1, action1["id"].(string))
	require.NotNil(t, action1["startedAt"].(string))
	require.NotNil(t, action1["endedAt"].(string))
	require.Equal(t, "session1", action1["sessionId"].(string))
	require.Equal(t, "application1", action1["application"].(string))
	require.Equal(t, "page1", action1["pageTitle"].(string))
	require.Equal(t, "http://app-1.ai", action1["pageUrl"].(string))
	require.Equal(t, float64(1), action1["orderInSession"].(float64))
	require.Equal(t, float64(10), action1["engagedTime"].(float64))

	action2 := actions[1].(map[string]interface{})
	require.Equal(t, "PageViewAction", action2["__typename"].(string))
	require.Equal(t, pageViewId2, action2["id"].(string))
	require.NotNil(t, action2["startedAt"].(string))
	require.NotNil(t, action2["endedAt"].(string))
	require.Equal(t, "session2", action2["sessionId"].(string))
	require.Equal(t, "application2", action2["application"].(string))
	require.Equal(t, "page2", action2["pageTitle"].(string))
	require.Equal(t, "http://app-2.ai", action2["pageUrl"].(string))
	require.Equal(t, float64(2), action2["orderInSession"].(float64))
	require.Equal(t, float64(20), action2["engagedTime"].(float64))

	action3 := actions[2].(map[string]interface{})
	require.Equal(t, "MessageAction", action3["__typename"].(string))
	require.Equal(t, messageId, action3["id"].(string))
	require.Equal(t, conversationId, action3["conversationId"].(string))
	require.Equal(t, "CHAT", action3["channel"].(string))
	require.NotNil(t, action3["startedAt"].(string))
}

func TestQueryResolver_Contact_WithActions_FilterByActionType(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	userId := neo4jt.CreateDefaultUser(driver, tenantName)
	conversationId := neo4jt.CreateConversation(driver, userId, contactId)

	now := time.Now().UTC()
	secAgo1 := now.Add(time.Duration(-1) * time.Second)
	from := now.Add(time.Duration(-10) * time.Minute)

	actionId1 := neo4jt.CreatePageView(driver, contactId, entity.PageViewEntity{
		StartedAt:      secAgo1,
		EndedAt:        now,
		TrackerName:    "tracker1",
		SessionId:      "session1",
		Application:    "application1",
		PageTitle:      "page1",
		PageUrl:        "http://app-1.ai",
		OrderInSession: 1,
		EngagedTime:    10,
	})
	neo4jt.AddMessageToConversation(driver, conversationId, mapper.MapMessageChannelFromModel(model.MessageChannelChat), secAgo1)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Action"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "PageView"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Message"))

	types := []model.ActionType{}
	types = append(types, model.ActionTypePageView)

	rawResponse, err := c.RawPost(getQuery("get_contact_with_actions_filter_by_action_type"),
		client.Var("contactId", contactId),
		client.Var("from", from),
		client.Var("to", now),
		client.Var("types", types))
	assertRawResponseSuccess(t, rawResponse, err)

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])

	actions := contact.(map[string]interface{})["actions"].([]interface{})
	require.Equal(t, 1, len(actions))
	action1 := actions[0].(map[string]interface{})
	require.Equal(t, "PageViewAction", action1["__typename"].(string))
	require.Equal(t, actionId1, action1["id"].(string))
}
