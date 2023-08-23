package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/event_store"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	contactProto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	emailProto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryResolver_ContactByEmail(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	otherTenant := "other"
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, otherTenant)
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, otherTenant)
	neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId1, "test@test.com", true, "MAIN")
	neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, otherTenant, contactId2, "test@test.com", true, "MAIN")

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_by_email"), client.Var("email", "test@test.com"))
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
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	otherTenant := "other"
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, otherTenant)
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, otherTenant)
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId1, "+1234567890", false, "OTHER")
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId2, "+1234567890", true, "MAIN")

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_by_phone"), client.Var("e164", "+1234567890"))
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
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	rawResponse, err := c.RawPost(getQuery("contact/create_contact_min"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact_Create model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "", *contact.Contact_Create.Prefix)
	require.Equal(t, "", *contact.Contact_Create.Name)
	require.Equal(t, "", *contact.Contact_Create.FirstName)
	require.Equal(t, "", *contact.Contact_Create.LastName)
	require.Equal(t, "", *contact.Contact_Create.Description)
	require.Equal(t, "", *contact.Contact_Create.Timezone)
	require.Equal(t, "", *contact.Contact_Create.ProfilePhotoURL)
	require.Equal(t, model.DataSourceOpenline, contact.Contact_Create.Source)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 3, neo4jt.GetTotalCountOfNodes(ctx, driver))

	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "User", "User_" + tenantName, "Contact", "Contact_" + tenantName})
}

func TestMutationResolver_ContactCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, "otherTenant")
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	rawResponse, err := c.RawPost(getQuery("contact/create_contact"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact_Create model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "MR", *contact.Contact_Create.Prefix)
	require.Equal(t, "first", *contact.Contact_Create.FirstName)
	require.Equal(t, "last", *contact.Contact_Create.LastName)
	require.Equal(t, "full name", *contact.Contact_Create.Name)
	require.Equal(t, "Some description", *contact.Contact_Create.Description)
	require.Equal(t, "America/Los_Angeles", *contact.Contact_Create.Timezone)
	require.Equal(t, model.DataSourceOpenline, contact.Contact_Create.Source)

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
	require.Equal(t, "contact@abc.com", *contact.Contact_Create.Emails[0].RawEmail)
	require.Equal(t, "contact@abc.com", *contact.Contact_Create.Emails[0].Email)
	require.Equal(t, "WORK", contact.Contact_Create.Emails[0].Label.String())
	require.Equal(t, false, contact.Contact_Create.Emails[0].Primary)
	require.Equal(t, model.DataSourceOpenline, contact.Contact_Create.Emails[0].Source)

	require.Equal(t, 1, len(contact.Contact_Create.PhoneNumbers))
	require.NotNil(t, contact.Contact_Create.PhoneNumbers[0].ID)
	require.Equal(t, "+1234567890", *contact.Contact_Create.PhoneNumbers[0].RawPhoneNumber)
	require.Nil(t, contact.Contact_Create.PhoneNumbers[0].E164)
	require.Equal(t, "MOBILE", contact.Contact_Create.PhoneNumbers[0].Label.String())
	require.Equal(t, true, contact.Contact_Create.PhoneNumbers[0].Primary)
	require.Equal(t, model.DataSourceOpenline, contact.Contact_Create.PhoneNumbers[0].Source)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "CustomField"))
	require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "CustomField_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "TextField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "IntField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "FloatField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "BoolField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "TimeField"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber_"+tenantName))
	require.Equal(t, 11, neo4jt.GetTotalCountOfNodes(ctx, driver))

	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "User", "User_" + tenantName, "Contact", "Contact_" + tenantName,
		"Email", "Email_" + tenantName, "PhoneNumber", "PhoneNumber_" + tenantName,
		"CustomField", "BoolField", "TextField", "FloatField", "TimeField", "IntField", "CustomField_" + tenantName})
}

func TestMutationResolver_CustomerContactCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, "otherTenant")
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	timeNow := time.Now().UTC()

	createdContactId, _ := uuid.NewUUID()
	createdEmailId, _ := uuid.NewUUID()

	calledCreateContact, calledCreateEmail, calledLinkEmailToContact := false, false, false

	contactServiceCallbacks := event_store.MockContactServiceCallbacks{
		CreateContact: func(context context.Context, contact *contactProto.CreateContactGrpcRequest) (*contactProto.CreateContactGrpcResponse, error) {
			require.Equal(t, "Bob", contact.FirstName)
			require.Equal(t, "Smith", contact.LastName)
			require.Equal(t, "Mr.", contact.Prefix)
			require.Equal(t, "This is a person", contact.Description)
			require.Equal(t, "unit-test", contact.AppSource)
			require.Equal(t, timeNow.Unix(), contact.CreatedAt.Seconds)
			require.Equal(t, "openline", contact.Tenant)
			calledCreateContact = true
			return &contactProto.CreateContactGrpcResponse{
				Id: createdContactId.String(),
			}, nil
		},
		LinkEmailToContact: func(context context.Context, link *contactProto.LinkEmailToContactGrpcRequest) (*contactProto.ContactIdGrpcResponse, error) {
			require.Equal(t, createdContactId.String(), link.ContactId)
			require.Equal(t, createdEmailId.String(), link.EmailId)
			require.Equal(t, true, link.Primary)
			require.Equal(t, "WORK", link.Label)
			require.Equal(t, "openline", link.Tenant)
			calledLinkEmailToContact = true
			return &contactProto.ContactIdGrpcResponse{
				Id: createdContactId.String(),
			}, nil
		},
	}
	event_store.SetContactCallbacks(&contactServiceCallbacks)

	emailServiceCallbacks := event_store.MockEmailServiceCallbacks{
		UpsertEmail: func(ctx context.Context, data *emailProto.UpsertEmailGrpcRequest) (*emailProto.EmailIdGrpcResponse, error) {
			require.Equal(t, "contact@abc.com", data.RawEmail)
			require.Equal(t, "openline", data.Tenant)
			calledCreateEmail = true
			return &emailProto.EmailIdGrpcResponse{
				Id: createdEmailId.String(),
			}, nil
		},
	}
	event_store.SetEmailCallbacks(&emailServiceCallbacks)

	rawResponse, err := c.RawPost(getQuery("contact/customer_create_contact"),
		client.Var("firstName", "Bob"),
		client.Var("lastName", "Smith"),
		client.Var("prefix", "Mr."),
		client.Var("description", "This is a person"),
		client.Var("appSource", "unit-test"),
		client.Var("createdAt", timeNow),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Customer_contact_Create model.CustomerContact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.Equal(t, createdContactId.String(), contact.Customer_contact_Create.ID)
	require.Equal(t, createdEmailId.String(), contact.Customer_contact_Create.Email.ID)
	require.True(t, calledCreateContact)
	require.True(t, calledCreateEmail)
	require.True(t, calledLinkEmailToContact)
}

func TestMutationResolver_ContactCreate_WithCustomFields(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	entityTemplateId := neo4jt.CreateEntityTemplate(ctx, driver, tenantName, model.EntityTemplateExtensionContact.String())
	fieldTemplateId := neo4jt.AddFieldTemplateToEntity(ctx, driver, entityTemplateId)
	setTemplateId := neo4jt.AddSetTemplateToEntity(ctx, driver, entityTemplateId)
	fieldInSetTemplateId := neo4jt.AddFieldTemplateToSet(ctx, driver, setTemplateId)

	rawResponse, err := c.RawPost(getQuery("contact/create_contact_with_custom_fields"),
		client.Var("entityTemplateId", entityTemplateId),
		client.Var("fieldTemplateId", fieldTemplateId),
		client.Var("setTemplateId", setTemplateId),
		client.Var("fieldInSetTemplateId", fieldInSetTemplateId))
	assertRawResponseSuccess(t, rawResponse, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "CustomField"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "CustomField_"+tenantName))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "TextField"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "EntityTemplate"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "CustomFieldTemplate"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "FieldSetTemplate"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "FieldSet"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "FieldSet_"+tenantName))
	require.Equal(t, 13, neo4jt.GetTotalCountOfNodes(ctx, driver))

	var contact struct {
		Contact_Create model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)

	createdContact := contact.Contact_Create
	require.Equal(t, model.DataSourceOpenline, createdContact.Source)
	require.Equal(t, entityTemplateId, createdContact.Template.ID)
	require.Equal(t, 2, len(createdContact.CustomFields))
	require.Equal(t, "field1", createdContact.CustomFields[0].Name)
	require.Equal(t, "TEXT", createdContact.CustomFields[0].Datatype.String())
	require.Equal(t, "value1", createdContact.CustomFields[0].Value.RealValue())
	require.Equal(t, model.DataSourceOpenline, createdContact.CustomFields[0].Source)
	require.Equal(t, fieldTemplateId, createdContact.CustomFields[0].Template.ID)
	require.NotNil(t, createdContact.CustomFields[0].ID)
	require.NotNil(t, createdContact.CustomFields[0].CreatedAt)
	require.Equal(t, "field2", createdContact.CustomFields[1].Name)
	require.Equal(t, "TEXT", createdContact.CustomFields[1].Datatype.String())
	require.Equal(t, "value2", createdContact.CustomFields[1].Value.RealValue())
	require.Equal(t, model.DataSourceOpenline, createdContact.CustomFields[1].Source)
	require.NotNil(t, createdContact.CustomFields[1].ID)
	require.NotNil(t, createdContact.CustomFields[1].CreatedAt)
	require.Equal(t, 2, len(createdContact.FieldSets))
	var set1, set2 *model.FieldSet
	if createdContact.FieldSets[0].Name == "set1" {
		set1 = createdContact.FieldSets[0]
		set2 = createdContact.FieldSets[1]
	} else {
		set1 = createdContact.FieldSets[1]
		set2 = createdContact.FieldSets[0]
	}
	require.NotNil(t, set1.ID)
	require.NotNil(t, set1.CreatedAt)
	require.Equal(t, "set1", set1.Name)
	require.Equal(t, 2, len(set1.CustomFields))
	require.NotNil(t, set1.CustomFields[0].CreatedAt)
	require.Equal(t, "field3InSet", set1.CustomFields[0].Name)
	require.Equal(t, "value3", set1.CustomFields[0].Value.RealValue())
	require.Equal(t, model.DataSourceOpenline, set1.CustomFields[0].Source)
	require.Equal(t, "TEXT", set1.CustomFields[0].Datatype.String())
	require.Equal(t, fieldInSetTemplateId, set1.CustomFields[0].Template.ID)
	require.NotNil(t, set1.CustomFields[1].CreatedAt)
	require.Equal(t, "field4InSet", set1.CustomFields[1].Name)
	require.Equal(t, "value4", set1.CustomFields[1].Value.RealValue())
	require.Equal(t, model.DataSourceOpenline, set1.CustomFields[1].Source)
	require.Equal(t, "TEXT", set1.CustomFields[1].Datatype.String())
	require.Nil(t, set1.CustomFields[1].Template)
	require.Equal(t, model.DataSourceOpenline, set1.Source)
	require.NotNil(t, set2.ID)
	require.NotNil(t, set2.CreatedAt)
	require.Equal(t, "set2", set2.Name)
	require.Equal(t, model.DataSourceOpenline, set2.Source)

	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "User", "User_" + tenantName, "Contact", "Contact_" + tenantName,
		"CustomFieldTemplate", "EntityTemplate", "FieldSet", "FieldSet_" + tenantName, "FieldSetTemplate",
		"CustomField", "TextField", "CustomField_" + tenantName})
}

func TestMutationResolver_ContactCreate_WithOwner(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	userId := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "Agent",
		LastName:  "Smith",
	})

	rawResponse, err := c.RawPost(getQuery("contact/create_contact_with_owner"),
		client.Var("ownerId", userId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact_Create model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	createdContact := contact.Contact_Create
	require.Equal(t, "", *createdContact.Prefix)
	require.Equal(t, "first", *createdContact.FirstName)
	require.Equal(t, "last", *createdContact.LastName)
	require.Equal(t, userId, createdContact.Owner.ID)
	require.Equal(t, "Agent", createdContact.Owner.FirstName)
	require.Equal(t, "Smith", createdContact.Owner.LastName)
	require.Equal(t, model.DataSourceOpenline, createdContact.Source)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 3, neo4jt.GetTotalCountOfNodes(ctx, driver))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "OWNS"))

	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "User", "User_" + tenantName})
}

func TestMutationResolver_ContactCreate_WithExternalReference(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	neo4jt.CreateHubspotExternalSystem(ctx, driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("contact/create_contact_with_external_reference"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact_Create model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.NotNil(t, contact.Contact_Create.ID)
	require.Equal(t, model.DataSourceOpenline, contact.Contact_Create.Source)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "ExternalSystem"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "ExternalSystem_"+tenantName))
	require.Equal(t, 4, neo4jt.GetTotalCountOfNodes(ctx, driver))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "IS_LINKED_WITH"))

	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "User", "User_" + tenantName, "Contact", "Contact_" + tenantName, "ExternalSystem", "ExternalSystem_" + tenantName})
}

func TestMutationResolver_ContactUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	origOwnerId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	newOwnerId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	contactId := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		Prefix:        "MR",
		FirstName:     "first",
		LastName:      "last",
		Description:   "description",
		Source:        entity.DataSourceHubspot,
		SourceOfTruth: entity.DataSourceHubspot,
	})

	neo4jt.UserOwnsContact(ctx, driver, origOwnerId, contactId)

	rawResponse, err := c.RawPost(getQuery("contact/update_contact"),
		client.Var("contactId", contactId),
		client.Var("ownerId", newOwnerId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactStruct struct {
		Contact_Update model.Contact
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)

	updatedContact := contactStruct.Contact_Update

	require.Equal(t, "DR", *updatedContact.Prefix)
	require.Equal(t, "updated name", *updatedContact.Name)
	require.Equal(t, "updated first", *updatedContact.FirstName)
	require.Equal(t, "updated last", *updatedContact.LastName)
	require.Equal(t, "updated description", *updatedContact.Description)
	require.Equal(t, "updated timezone", *updatedContact.Timezone)
	require.Equal(t, "test", *updatedContact.AppSource)
	require.Equal(t, model.DataSourceOpenline, updatedContact.SourceOfTruth)
	require.Equal(t, newOwnerId, updatedContact.Owner.ID)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "OWNS"))

	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "User", "User_" + tenantName})
}

func TestMutationResolver_ContactUpdate_ClearTitle(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		Prefix:    "MR",
		FirstName: "first",
		LastName:  "last",
	})

	rawResponse, err := c.RawPost(getQuery("contact/update_contact_clear_title"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact_Update model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	updatedContact := contact.Contact_Update
	require.Equal(t, "", *updatedContact.Prefix)
	require.Equal(t, "updated first", *updatedContact.FirstName)
	require.Equal(t, "updated last", *updatedContact.LastName)
}

func TestQueryResolver_Contact_WithJobRoles_ById(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name:        "name1",
		Description: "description1",
		Website:     "website1",
		Industry:    "industry1",
		IsPublic:    true,
	})
	neo4jt.AddDomainToOrg(ctx, driver, organizationId1, "domain1")
	organizationId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name:        "name2",
		Description: "description2",
		Website:     "website2",
		Industry:    "industry2",
		IsPublic:    false,
	})
	neo4jt.AddDomainToOrg(ctx, driver, organizationId2, "domain2")
	role1 := neo4jt.ContactWorksForOrganization(ctx, driver, contactId, organizationId1, "CTO", false)
	role2 := neo4jt.ContactWorksForOrganization(ctx, driver, contactId, organizationId2, "CEO", true)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "JobRole"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "WORKS_AS"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_job_roles_by_id"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var searchedContact struct {
		Contact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedContact)
	require.Nil(t, err)
	require.Equal(t, contactId, searchedContact.Contact.ID)

	roles := searchedContact.Contact.JobRoles
	require.Equal(t, 2, len(roles))
	var cto, ceo *model.JobRole
	ceo = roles[0]
	cto = roles[1]
	require.Equal(t, role1, cto.ID)
	require.Equal(t, "CTO", *cto.JobTitle)
	require.Equal(t, false, cto.Primary)
	require.Equal(t, organizationId1, cto.Organization.ID)
	require.Equal(t, "name1", cto.Organization.Name)
	require.Equal(t, "description1", *cto.Organization.Description)
	require.Equal(t, []string{"domain1"}, cto.Organization.Domains)
	require.Equal(t, "website1", *cto.Organization.Website)
	require.Equal(t, "industry1", *cto.Organization.Industry)
	require.Equal(t, true, *cto.Organization.IsPublic)
	require.NotNil(t, cto.Organization.CreatedAt)

	require.Equal(t, role2, ceo.ID)
	require.Equal(t, "CEO", *ceo.JobTitle)
	require.Equal(t, true, ceo.Primary)
	require.Equal(t, organizationId2, ceo.Organization.ID)
	require.Equal(t, "name2", ceo.Organization.Name)
	require.Equal(t, "description2", *ceo.Organization.Description)
	require.Equal(t, []string{"domain2"}, ceo.Organization.Domains)
	require.Equal(t, "website2", *ceo.Organization.Website)
	require.Equal(t, "industry2", *ceo.Organization.Industry)
	require.Equal(t, false, *ceo.Organization.IsPublic)
	require.NotNil(t, ceo.Organization.CreatedAt)
}

func TestQueryResolver_Contact_WithNotes_ById(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	noteId1 := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "note1", utils.Now())
	noteId2 := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "note2", utils.Now())
	neo4jt.NoteCreatedByUser(ctx, driver, noteId1, userId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "NOTED"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "CREATED"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_notes_by_id"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var searchedContact struct {
		Contact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedContact)
	require.Nil(t, err)
	require.Equal(t, contactId, searchedContact.Contact.ID)

	notes := searchedContact.Contact.Notes.Content
	require.Equal(t, 2, len(notes))
	var noteWithUser, noteWithoutUser *model.Note
	if noteId1 == notes[0].ID {
		noteWithUser = notes[0]
		noteWithoutUser = notes[1]
	} else {
		noteWithUser = notes[1]
		noteWithoutUser = notes[0]
	}
	require.Equal(t, noteId1, noteWithUser.ID)
	require.Equal(t, "note1", noteWithUser.HTML)
	require.NotNil(t, noteWithUser.CreatedAt)
	require.NotNil(t, noteWithUser.CreatedBy)
	require.Equal(t, userId, noteWithUser.CreatedBy.ID)
	require.Equal(t, "first", noteWithUser.CreatedBy.FirstName)
	require.Equal(t, "last", noteWithUser.CreatedBy.LastName)

	require.Equal(t, noteId2, noteWithoutUser.ID)
	require.Equal(t, "note2", noteWithoutUser.HTML)
	require.NotNil(t, noteWithoutUser.CreatedAt)
	require.Nil(t, noteWithoutUser.CreatedBy)
}

func TestQueryResolver_Contact_WithNotes_ById_Time_Range(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	noteId1 := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "note1", utils.Now())
	noteId2 := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "note2", utils.Now())
	neo4jt.NoteCreatedByUser(ctx, driver, noteId1, userId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "NOTED"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "CREATED"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_notes_by_id_time_range"),
		client.Var("contactId", contactId),
		client.Var("start", time.Now().Add(-1*time.Hour)),
		client.Var("end", time.Now().Add(1*time.Hour)))

	assertRawResponseSuccess(t, rawResponse, err)

	var searchedContact struct {
		Contact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedContact)
	require.Nil(t, err)
	require.Equal(t, contactId, searchedContact.Contact.ID)

	notes := searchedContact.Contact.NotesByTime
	require.Equal(t, 2, len(notes))
	var noteWithUser, noteWithoutUser *model.Note
	if noteId1 == notes[0].ID {
		noteWithUser = notes[0]
		noteWithoutUser = notes[1]
	} else {
		noteWithUser = notes[1]
		noteWithoutUser = notes[0]
	}
	require.Equal(t, noteId1, noteWithUser.ID)
	require.Equal(t, "note1", noteWithUser.HTML)
	require.NotNil(t, noteWithUser.CreatedAt)
	require.NotNil(t, noteWithUser.CreatedBy)
	require.Equal(t, userId, noteWithUser.CreatedBy.ID)
	require.Equal(t, "first", noteWithUser.CreatedBy.FirstName)
	require.Equal(t, "last", noteWithUser.CreatedBy.LastName)

	require.Equal(t, noteId2, noteWithoutUser.ID)
	require.Equal(t, "note2", noteWithoutUser.HTML)
	require.NotNil(t, noteWithoutUser.CreatedAt)
	require.Nil(t, noteWithoutUser.CreatedBy)

	// test with time range that does not include any notes
	rawResponse, err = c.RawPost(getQuery("contact/get_contact_with_notes_by_id_time_range"),
		client.Var("contactId", contactId),
		client.Var("start", time.Now().Add(-2*time.Hour)),
		client.Var("end", time.Now().Add(-1*time.Hour)))

	assertRawResponseSuccess(t, rawResponse, err)

	searchedContact.Contact = model.Contact{}
	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedContact)
	require.Nil(t, err)
	require.Equal(t, contactId, searchedContact.Contact.ID)

	notes = searchedContact.Contact.NotesByTime
	require.Equal(t, 0, len(notes))

	rawResponse, err = c.RawPost(getQuery("contact/get_contact_with_notes_by_id_time_range"),
		client.Var("contactId", contactId),
		client.Var("start", time.Now().Add(1*time.Hour)),
		client.Var("end", time.Now().Add(2*time.Hour)))

	assertRawResponseSuccess(t, rawResponse, err)

	searchedContact.Contact = model.Contact{}
	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedContact)
	require.Nil(t, err)
	require.Equal(t, contactId, searchedContact.Contact.ID)

	notes = searchedContact.Contact.NotesByTime
	require.Equal(t, 0, len(notes))
}

func TestQueryResolver_Contact_WithTags_ById(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	tagId1 := neo4jt.CreateTag(ctx, driver, tenantName, "tag1")
	tagId2 := neo4jt.CreateTag(ctx, driver, tenantName, "tag2")
	tagId3 := neo4jt.CreateTag(ctx, driver, tenantName, "tag3")
	neo4jt.TagContact(ctx, driver, contactId, tagId1)
	neo4jt.TagContact(ctx, driver, contactId, tagId2)
	neo4jt.TagContact(ctx, driver, contactId2, tagId3)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "Tag"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(ctx, driver, "TAGGED"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_tags_by_id"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactStruct struct {
		Contact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	contact := contactStruct.Contact

	require.Nil(t, err)
	require.Equal(t, contactId, contact.ID)

	tags := contact.Tags
	require.Equal(t, 2, len(tags))
	require.Equal(t, tagId1, tags[0].ID)
	require.Equal(t, "tag1", tags[0].Name)
	require.Equal(t, tagId2, tags[1].ID)
	require.Equal(t, "tag2", tags[1].Name)
}

func TestQueryResolver_Contact_WithLocations_ById(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	locationId1 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:         "WORK",
		Source:       entity.DataSourceOpenline,
		AppSource:    "test",
		Country:      "testCountry",
		Region:       "testRegion",
		Locality:     "testLocality",
		Address:      "testAddress",
		Address2:     "testAddress2",
		Zip:          "testZip",
		AddressType:  "testAddressType",
		HouseNumber:  "testHouseNumber",
		PostalCode:   "testPostalCode",
		PlusFour:     "testPlusFour",
		Commercial:   true,
		Predirection: "testPredirection",
		District:     "testDistrict",
		Street:       "testStreet",
		RawAddress:   "testRawAddress",
		UtcOffset:    1,
		TimeZone:     "paris",
		Latitude:     utils.ToPtr(float64(0.001)),
		Longitude:    utils.ToPtr(float64(-2.002)),
	})
	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:      "UNKNOWN",
		Source:    entity.DataSourceOpenline,
		AppSource: "test",
	})
	neo4jt.ContactAssociatedWithLocation(ctx, driver, contactId, locationId1)
	neo4jt.ContactAssociatedWithLocation(ctx, driver, contactId, locationId2)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_locations_by_id"),
		client.Var("contactId", contactId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var contactStruct struct {
		Contact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)

	contact := contactStruct.Contact
	require.NotNil(t, contact)
	require.Equal(t, 2, len(contact.Locations))

	var locationWithAddressDtls, locationWithoutAddressDtls *model.Location
	if contact.Locations[0].ID == locationId1 {
		locationWithAddressDtls = contact.Locations[0]
		locationWithoutAddressDtls = contact.Locations[1]
	} else {
		locationWithAddressDtls = contact.Locations[1]
		locationWithoutAddressDtls = contact.Locations[0]
	}

	require.Equal(t, locationId1, locationWithAddressDtls.ID)
	require.Equal(t, "WORK", *locationWithAddressDtls.Name)
	require.NotNil(t, locationWithAddressDtls.CreatedAt)
	require.NotNil(t, locationWithAddressDtls.UpdatedAt)
	require.Equal(t, "test", locationWithAddressDtls.AppSource)
	require.Equal(t, model.DataSourceOpenline, locationWithAddressDtls.Source)
	require.Equal(t, "testCountry", *locationWithAddressDtls.Country)
	require.Equal(t, "testLocality", *locationWithAddressDtls.Locality)
	require.Equal(t, "testRegion", *locationWithAddressDtls.Region)
	require.Equal(t, "testAddress", *locationWithAddressDtls.Address)
	require.Equal(t, "testAddress2", *locationWithAddressDtls.Address2)
	require.Equal(t, "testZip", *locationWithAddressDtls.Zip)
	require.Equal(t, "testAddressType", *locationWithAddressDtls.AddressType)
	require.Equal(t, "testHouseNumber", *locationWithAddressDtls.HouseNumber)
	require.Equal(t, "testPostalCode", *locationWithAddressDtls.PostalCode)
	require.Equal(t, "testPlusFour", *locationWithAddressDtls.PlusFour)
	require.Equal(t, true, *locationWithAddressDtls.Commercial)
	require.Equal(t, "testPredirection", *locationWithAddressDtls.Predirection)
	require.Equal(t, "testDistrict", *locationWithAddressDtls.District)
	require.Equal(t, "testStreet", *locationWithAddressDtls.Street)
	require.Equal(t, "testRawAddress", *locationWithAddressDtls.RawAddress)
	require.Equal(t, "paris", *locationWithAddressDtls.TimeZone)
	require.Equal(t, int64(1), *locationWithAddressDtls.UtcOffset)
	require.Equal(t, float64(0.001), *locationWithAddressDtls.Latitude)
	require.Equal(t, float64(-2.002), *locationWithAddressDtls.Longitude)

	require.Equal(t, locationId2, locationWithoutAddressDtls.ID)
	require.Equal(t, "UNKNOWN", *locationWithoutAddressDtls.Name)
	require.NotNil(t, locationWithoutAddressDtls.CreatedAt)
	require.NotNil(t, locationWithoutAddressDtls.UpdatedAt)
	require.Equal(t, "test", locationWithoutAddressDtls.AppSource)
	require.Equal(t, model.DataSourceOpenline, locationWithoutAddressDtls.Source)
	require.Equal(t, "", *locationWithoutAddressDtls.Country)
	require.Equal(t, "", *locationWithoutAddressDtls.Region)
	require.Equal(t, "", *locationWithoutAddressDtls.Locality)
	require.Equal(t, "", *locationWithoutAddressDtls.Address)
	require.Equal(t, "", *locationWithoutAddressDtls.Address2)
	require.Equal(t, "", *locationWithoutAddressDtls.Zip)
	require.False(t, *locationWithoutAddressDtls.Commercial)
	require.Nil(t, locationWithoutAddressDtls.Latitude)
	require.Nil(t, locationWithoutAddressDtls.Longitude)
}

func TestQueryResolver_Contacts_SortByTitleAscFirstNameAscLastNameDesc(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contact1 := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		Prefix:    "MR",
		FirstName: "contact",
		LastName:  "1",
	})
	contact2 := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		Prefix:    "DR",
		FirstName: "contact",
		LastName:  "9",
	})
	contact3 := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		Prefix:    "",
		FirstName: "contact",
		LastName:  "222",
	})
	contact4 := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		Prefix:    "MR",
		FirstName: "other contact",
		LastName:  "A",
	})

	rawResponse, err := c.RawPost(getQuery("contact/get_contacts_with_sorting"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contacts struct {
		Contacts model.ContactsPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contacts)
	require.Nil(t, err)
	require.NotNil(t, contacts.Contacts)
	require.Equal(t, 4, len(contacts.Contacts.Content))
	require.Equal(t, contact3, contacts.Contacts.Content[0].ID)
	require.Equal(t, contact2, contacts.Contacts.Content[1].ID)
	require.Equal(t, contact1, contacts.Contacts.Content[2].ID)
	require.Equal(t, contact4, contacts.Contacts.Content[3].ID)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
}

func TestQueryResolver_Contact_BasicFilters_FindContactWithLetterAInName(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactFoundByFirstName := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		Prefix:    "MR",
		Name:      "contact1",
		FirstName: "aa",
		LastName:  "bb",
	})
	contactFoundByLastName := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		Prefix:    "MR",
		FirstName: "bb",
		LastName:  "AA",
	})
	contactFilteredOut := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		Prefix:    "MR",
		FirstName: "bb",
		LastName:  "BB",
	})

	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contacts_basic_filters"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactsStruct struct {
		Contacts model.ContactsPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactsStruct)
	require.Nil(t, err)
	require.NotNil(t, contactsStruct.Contacts)

	contacts := contactsStruct.Contacts.Content

	require.Equal(t, 2, len(contacts))
	require.Equal(t, contactFoundByFirstName, contacts[0].ID)
	require.Equal(t, "contact1", *contacts[0].Name)
	require.Equal(t, "aa", *contacts[0].FirstName)
	require.Equal(t, "bb", *contacts[0].LastName)
	require.Equal(t, contactFoundByLastName, contacts[1].ID)
	require.Equal(t, 1, contactsStruct.Contacts.TotalPages)
	require.Equal(t, int64(2), contactsStruct.Contacts.TotalElements)
	// suppress unused warnings
	require.NotNil(t, contactFilteredOut)
}

func TestQueryResolver_Contact_WithTimelineEvents(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.CreateDefaultUser(ctx, driver, tenantName)

	now := time.Now().UTC()
	secAgo1 := now.Add(time.Duration(-1) * time.Second)
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo30 := now.Add(time.Duration(-30) * time.Second)
	secAgo40 := now.Add(time.Duration(-40) * time.Second)
	secAgo50 := now.Add(time.Duration(-50) * time.Second)
	secAgo55 := now.Add(time.Duration(-55) * time.Second)
	secAgo60 := now.Add(time.Duration(-60) * time.Second)
	minAgo5 := now.Add(time.Duration(-5) * time.Minute)

	// prepare page views
	pageViewId1 := neo4jt.CreatePageView(ctx, driver, contactId, entity.PageViewEntity{
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
	pageViewId2 := neo4jt.CreatePageView(ctx, driver, contactId, entity.PageViewEntity{
		StartedAt:      secAgo10,
		EndedAt:        now,
		TrackerName:    "tracker2",
		SessionId:      "session2",
		Application:    "application2",
		PageTitle:      "page2",
		PageUrl:        "http://app-2.ai",
		OrderInSession: 2,
		EngagedTime:    20,
	})
	neo4jt.CreatePageView(ctx, driver, contactId2, entity.PageViewEntity{})

	voiceSession := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "CALL", "ACTIVE", "VOICE", now, false)

	analysis1 := neo4jt.CreateAnalysis(ctx, driver, tenantName, "This is a summary of the conversation", "text/plain", "SUMMARY", secAgo55)
	neo4jt.AnalysisDescribes(ctx, driver, tenantName, analysis1, voiceSession, string(repository.LINKED_WITH_INTERACTION_SESSION))

	// prepare meeting
	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "meeting-name", secAgo60)
	neo4jt.MeetingAttendedBy(ctx, driver, meetingId, contactId)

	// prepare contact notes
	contactNoteId := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 1", secAgo30)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 2", minAgo5)

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 1", "application/json", &channel, secAgo40)
	interactionEventId2 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 2", "application/json", &channel, secAgo50)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "email1", false, "WORK")
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId, "+1234", false, "WORK")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, phoneNumberId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId2, phoneNumberId, "")
	neo4jt.InteractionSessionAttendedBy(ctx, driver, tenantName, voiceSession, phoneNumberId, "")

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "PageView"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Analysis"))
	require.Equal(t, 9, neo4jt.GetCountOfNodes(ctx, driver, "TimelineEvent"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_timeline_events"),
		client.Var("contactId", contactId),
		client.Var("from", now),
		client.Var("size", 8))
	assertRawResponseSuccess(t, rawResponse, err)

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])

	timelineEvents := contact.(map[string]interface{})["timelineEvents"].([]interface{})
	require.Equal(t, 8, len(timelineEvents))

	timelineEvent1 := timelineEvents[0].(map[string]interface{})
	require.Equal(t, "PageView", timelineEvent1["__typename"].(string))
	require.Equal(t, pageViewId1, timelineEvent1["id"].(string))
	require.NotNil(t, timelineEvent1["startedAt"].(string))
	require.NotNil(t, timelineEvent1["endedAt"].(string))
	require.Equal(t, "session1", timelineEvent1["sessionId"].(string))
	require.Equal(t, "application1", timelineEvent1["application"].(string))
	require.Equal(t, "page1", timelineEvent1["pageTitle"].(string))
	require.Equal(t, "http://app-1.ai", timelineEvent1["pageUrl"].(string))
	require.Equal(t, float64(1), timelineEvent1["orderInSession"].(float64))
	require.Equal(t, float64(10), timelineEvent1["engagedTime"].(float64))
	require.Equal(t, "test", timelineEvent1["appSource"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["sourceOfTruth"].(string))

	timelineEvent2 := timelineEvents[1].(map[string]interface{})
	require.Equal(t, "PageView", timelineEvent2["__typename"].(string))
	require.Equal(t, pageViewId2, timelineEvent2["id"].(string))
	require.NotNil(t, timelineEvent2["startedAt"].(string))
	require.NotNil(t, timelineEvent2["endedAt"].(string))
	require.Equal(t, "session2", timelineEvent2["sessionId"].(string))
	require.Equal(t, "application2", timelineEvent2["application"].(string))
	require.Equal(t, "page2", timelineEvent2["pageTitle"].(string))
	require.Equal(t, "http://app-2.ai", timelineEvent2["pageUrl"].(string))
	require.Equal(t, float64(2), timelineEvent2["orderInSession"].(float64))
	require.Equal(t, float64(20), timelineEvent2["engagedTime"].(float64))

	timelineEvent4 := timelineEvents[3].(map[string]interface{})
	require.Equal(t, "Note", timelineEvent4["__typename"].(string))
	require.Equal(t, contactNoteId, timelineEvent4["id"].(string))
	require.NotNil(t, timelineEvent4["createdAt"].(string))
	require.Equal(t, "contact note 1", timelineEvent4["html"].(string))

	timelineEvent5 := timelineEvents[4].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent5["__typename"].(string))
	require.Equal(t, interactionEventId1, timelineEvent5["id"].(string))
	require.NotNil(t, timelineEvent5["createdAt"].(string))
	require.Equal(t, "IE text 1", timelineEvent5["content"].(string))
	require.Equal(t, "application/json", timelineEvent5["contentType"].(string))
	require.Equal(t, "EMAIL", timelineEvent5["channel"].(string))

	timelineEvent6 := timelineEvents[5].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent6["__typename"].(string))
	require.Equal(t, interactionEventId2, timelineEvent6["id"].(string))
	require.NotNil(t, timelineEvent6["createdAt"].(string))
	require.Equal(t, "IE text 2", timelineEvent6["content"].(string))
	require.Equal(t, "application/json", timelineEvent6["contentType"].(string))
	require.Equal(t, "EMAIL", timelineEvent6["channel"].(string))

	timelineEvent7 := timelineEvents[6].(map[string]interface{})
	require.Equal(t, "Meeting", timelineEvent7["__typename"].(string))
	require.Equal(t, meetingId, timelineEvent7["id"].(string))
	require.NotNil(t, timelineEvent7["createdAt"].(string))
	require.Equal(t, "meeting-name", timelineEvent7["name"].(string))
}

func TestQueryResolver_Contact_WithTimelineEvents_FilterByType(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	now := time.Now().UTC()
	secAgo1 := now.Add(time.Duration(-1) * time.Second)

	neo4jt.CreatePageView(ctx, driver, contactId, entity.PageViewEntity{
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

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "TimelineEvent"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PageView"))

	types := []model.TimelineEventType{}
	types = append(types, model.TimelineEventTypePageView)

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_timeline_filter_by_timeline_event_type"),
		client.Var("contactId", contactId),
		client.Var("from", now),
		client.Var("types", types))
	assertRawResponseSuccess(t, rawResponse, err)

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])

	timelineEvents := contact.(map[string]interface{})["timelineEvents"].([]interface{})
	require.Equal(t, 0, len(timelineEvents))
	//timelineEvent1 := timelineEvents[0].(map[string]interface{})
	//require.Equal(t, "PageView", timelineEvent1["__typename"].(string))
	//require.Equal(t, actionId1, timelineEvent1["id"].(string))
}

func TestQueryResolver_Contact_WithTimelineEventsTotalCount(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.CreateDefaultUser(ctx, driver, tenantName)

	now := time.Now().UTC()

	// prepare page views
	neo4jt.CreatePageView(ctx, driver, contactId, entity.PageViewEntity{
		StartedAt: now,
	})
	neo4jt.CreatePageView(ctx, driver, contactId, entity.PageViewEntity{
		StartedAt: now,
	})

	// prepare contact notes
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 1", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 2", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 3", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 4", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 5", now)

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text", "application/json", &channel, now)
	interactionEventId2 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text", "application/json", &channel, now)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "email1", false, "WORK")
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId, "+1234", false, "WORK")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId1, phoneNumberId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId, "")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "PageView"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Conversation"))
	require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 10, neo4jt.GetCountOfNodes(ctx, driver, "TimelineEvent"))

	rawResponse := callGraphQL(t, "contact/get_contact_with_timeline_events_total_count", map[string]interface{}{"contactId": contactId})

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])
	require.Equal(t, float64(10), contact.(map[string]interface{})["timelineEventsTotalCount"].(float64))
}

func TestQueryResolver_Contact_WithOrganizations_ById(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "organization1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "organization2")
	organizationId3 := neo4jt.CreateOrganization(ctx, driver, tenantName, "organization3")
	organizationId0 := neo4jt.CreateOrganization(ctx, driver, tenantName, "organization0")
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId, organizationId1)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId, organizationId2)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId, organizationId3)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId2, organizationId0)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(ctx, driver, "ROLE_IN"))

	rawResponse := callGraphQL(t, "contact/get_contact_with_organizations_by_id",
		map[string]interface{}{"contactId": contactId, "limit": 2, "page": 1})

	var searchedContact struct {
		Contact model.Contact
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &searchedContact)
	require.Nil(t, err)
	require.Equal(t, contactId, searchedContact.Contact.ID)
	require.Equal(t, 2, searchedContact.Contact.Organizations.TotalPages)
	require.Equal(t, int64(3), searchedContact.Contact.Organizations.TotalElements)

	organizations := searchedContact.Contact.Organizations.Content
	require.Equal(t, 2, len(organizations))
	require.Equal(t, organizationId1, organizations[0].ID)
	require.Equal(t, organizationId2, organizations[1].ID)
}

func TestMutationResolver_ContactAddTagByID(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	tagId1 := neo4jt.CreateTag(ctx, driver, tenantName, "tag1")
	tagId2 := neo4jt.CreateTag(ctx, driver, tenantName, "tag2")
	neo4jt.TagContact(ctx, driver, contactId, tagId1)
	time.Sleep(100 * time.Millisecond)

	rawResponse := callGraphQL(t, "contact/add_tag_to_contact", map[string]interface{}{"contactId": contactId, "tagId": tagId2})

	var contactStruct struct {
		Contact_AddTagById model.Contact
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)
	require.NotNil(t, contactStruct)
	tags := contactStruct.Contact_AddTagById.Tags
	require.Equal(t, contactId, contactStruct.Contact_AddTagById.ID)
	require.NotNil(t, contactStruct.Contact_AddTagById.UpdatedAt)
	require.Equal(t, 2, len(tags))
	require.Equal(t, tagId1, tags[0].ID)
	require.Equal(t, "tag1", tags[0].Name)
	require.Equal(t, tagId2, tags[1].ID)
	require.Equal(t, "tag2", tags[1].Name)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Tag"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "TAGGED"))
}

func TestMutationResolver_ContactRemoveTagByID(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	tagId1 := neo4jt.CreateTag(ctx, driver, tenantName, "tag1")
	tagId2 := neo4jt.CreateTag(ctx, driver, tenantName, "tag2")
	neo4jt.TagContact(ctx, driver, contactId, tagId1)
	neo4jt.TagContact(ctx, driver, contactId, tagId2)

	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "TAGGED"))

	rawResponse := callGraphQL(t, "contact/remove_tag_from_contact", map[string]interface{}{"contactId": contactId, "tagId": tagId2})

	var contactStruct struct {
		Contact_RemoveTagById model.Contact
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)
	require.NotNil(t, contactStruct)
	tags := contactStruct.Contact_RemoveTagById.Tags
	require.Equal(t, contactId, contactStruct.Contact_RemoveTagById.ID)
	require.NotNil(t, contactStruct.Contact_RemoveTagById.UpdatedAt)
	require.Equal(t, 1, len(tags))
	require.Equal(t, tagId1, tags[0].ID)
	require.Equal(t, "tag1", tags[0].Name)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Tag"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "TAGGED"))
}

func TestMutationResolver_ContactAddOrganizationByID(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	orgId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org1")
	orgId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org2")
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId, orgId1)

	rawResponse := callGraphQL(t, "contact/add_organization_to_contact", map[string]interface{}{"contactId": contactId, "organizationId": orgId2})

	var contactStruct struct {
		Contact_AddOrganizationById model.Contact
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)
	require.NotNil(t, contactStruct)
	organizations := contactStruct.Contact_AddOrganizationById.Organizations.Content
	require.Equal(t, contactId, contactStruct.Contact_AddOrganizationById.ID)
	require.NotNil(t, contactStruct.Contact_AddOrganizationById.UpdatedAt)
	require.Equal(t, 2, len(organizations))
	require.ElementsMatch(t, []string{orgId1, orgId2}, []string{organizations[0].ID, organizations[1].ID})
	require.ElementsMatch(t, []string{"org1", "org2"}, []string{organizations[0].Name, organizations[1].Name})

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
}

func TestMutationResolver_ContactRemoveOrganizationByID(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	orgId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org1")
	orgId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org2")
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId, orgId1)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId, orgId2)

	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ROLE_IN"))

	rawResponse := callGraphQL(t, "contact/remove_organization_from_contact", map[string]interface{}{"contactId": contactId, "organizationId": orgId2})

	var contactStruct struct {
		Contact_RemoveOrganizationById model.Contact
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)
	require.NotNil(t, contactStruct)
	organizations := contactStruct.Contact_RemoveOrganizationById.Organizations.Content
	require.Equal(t, contactId, contactStruct.Contact_RemoveOrganizationById.ID)
	require.NotNil(t, contactStruct.Contact_RemoveOrganizationById.UpdatedAt)
	require.Equal(t, 1, len(organizations))
	require.Equal(t, orgId1, organizations[0].ID)
	require.Equal(t, "org1", organizations[0].Name)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "JobRole"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
}

func TestMutationResolver_ContactAddNewLocation(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	rawResponse := callGraphQL(t, "contact/add_new_location_to_contact", map[string]interface{}{"contactId": contactId})

	var locationStruct struct {
		Contact_AddNewLocation model.Location
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &locationStruct)
	require.Nil(t, err)
	require.NotNil(t, locationStruct)
	location := locationStruct.Contact_AddNewLocation
	require.NotNil(t, location.ID)
	require.NotNil(t, location.CreatedAt)
	require.NotNil(t, location.UpdatedAt)
	require.Equal(t, constants.AppSourceCustomerOsApi, location.AppSource)
	require.Equal(t, model.DataSourceOpenline, location.Source)
	require.Equal(t, model.DataSourceOpenline, location.SourceOfTruth)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "LOCATION_BELONGS_TO_TENANT"))
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Location", "Location_" + tenantName, "Contact", "Contact_" + tenantName})
}

func TestMutationResolver_ContactAddSocial(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	rawResponse := callGraphQL(t, "contact/add_social_to_contact", map[string]interface{}{"contactId": contactId})

	var socialStruct struct {
		Contact_AddSocial model.Social
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &socialStruct)
	require.Nil(t, err)
	require.NotNil(t, socialStruct)
	social := socialStruct.Contact_AddSocial
	require.NotNil(t, social.ID)
	require.NotNil(t, social.CreatedAt)
	require.NotNil(t, social.UpdatedAt)
	test.AssertRecentTime(t, social.CreatedAt)
	test.AssertRecentTime(t, social.UpdatedAt)
	require.Equal(t, constants.AppSourceCustomerOsApi, social.AppSource)
	require.Equal(t, model.DataSourceOpenline, social.Source)
	require.Equal(t, model.DataSourceOpenline, social.SourceOfTruth)
	require.Equal(t, "social url", social.URL)
	require.Equal(t, "social platform", *social.PlatformName)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Social"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Social", "Social_" + tenantName, "Contact", "Contact_" + tenantName})
}

func TestQueryResolver_Contact_WithSocials(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	socialId1 := neo4jt.CreateSocial(ctx, driver, tenantName, entity.SocialEntity{
		PlatformName: "p1",
		Url:          "url1",
	})
	socialId2 := neo4jt.CreateSocial(ctx, driver, tenantName, entity.SocialEntity{
		PlatformName: "p2",
		Url:          "url2",
	})
	neo4jt.LinkSocialWithEntity(ctx, driver, contactId, socialId1)
	neo4jt.LinkSocialWithEntity(ctx, driver, contactId, socialId2)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Social"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))

	rawResponse := callGraphQL(t, "contact/get_contact_with_socials", map[string]interface{}{"contactId": contactId})

	var contactStruct struct {
		Contact model.Contact
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)

	contact := contactStruct.Contact
	require.NotNil(t, contact)
	require.Equal(t, 2, len(contact.Socials))

	require.Equal(t, socialId1, contact.Socials[0].ID)
	require.Equal(t, "p1", *contact.Socials[0].PlatformName)
	require.Equal(t, "url1", contact.Socials[0].URL)
	require.NotNil(t, contact.Socials[0].CreatedAt)
	require.NotNil(t, contact.Socials[0].UpdatedAt)
	require.Equal(t, "test", contact.Socials[0].AppSource)

	require.Equal(t, socialId2, contact.Socials[1].ID)
	require.Equal(t, "p2", *contact.Socials[1].PlatformName)
	require.Equal(t, "url2", contact.Socials[1].URL)
	require.NotNil(t, contact.Socials[1].CreatedAt)
	require.NotNil(t, contact.Socials[1].UpdatedAt)
	require.Equal(t, "test", contact.Socials[1].AppSource)
}
