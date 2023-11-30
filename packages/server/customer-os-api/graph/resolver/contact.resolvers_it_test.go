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
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	contactgrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	emailgrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	phonenumbergrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
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
	createdContactId := uuid.New().String()

	calledCreateContact, calledCreateEmail, calledCreatePhoneNumber := false, false, false

	contactServiceCallbacks := events_platform.MockContactServiceCallbacks{
		CreateContact: func(context context.Context, contact *contactgrpc.UpsertContactGrpcRequest) (*contactgrpc.ContactIdGrpcResponse, error) {
			require.Equal(t, "", contact.FirstName)
			require.Equal(t, "", contact.LastName)
			require.Equal(t, "", contact.Prefix)
			require.Equal(t, "", contact.Name)
			require.Equal(t, "", contact.Description)
			require.Equal(t, "", contact.Timezone)
			require.Equal(t, "", contact.ProfilePhotoUrl)
			require.Equal(t, "openline", contact.Tenant)
			require.Equal(t, string(entity.DataSourceOpenline), contact.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, contact.SourceFields.AppSource)
			require.Equal(t, testUserId, contact.LoggedInUserId)
			calledCreateContact = true
			neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
				Id: createdContactId,
			})
			return &contactgrpc.ContactIdGrpcResponse{
				Id: createdContactId,
			}, nil
		},
	}
	emailServiceCallbacks := events_platform.MockEmailServiceCallbacks{
		UpsertEmail: func(ctx context.Context, data *emailgrpc.UpsertEmailGrpcRequest) (*emailgrpc.EmailIdGrpcResponse, error) {
			calledCreateEmail = true
			return &emailgrpc.EmailIdGrpcResponse{
				Id: uuid.New().String(),
			}, nil
		},
	}
	phoneNumberServiceCallbacks := events_platform.MockPhoneNumberServiceCallbacks{
		UpsertPhoneNumber: func(ctx context.Context, data *phonenumbergrpc.UpsertPhoneNumberGrpcRequest) (*phonenumbergrpc.PhoneNumberIdGrpcResponse, error) {
			calledCreatePhoneNumber = true
			return &phonenumbergrpc.PhoneNumberIdGrpcResponse{
				Id: uuid.New().String(),
			}, nil
		},
	}
	events_platform.SetContactCallbacks(&contactServiceCallbacks)
	events_platform.SetEmailCallbacks(&emailServiceCallbacks)
	events_platform.SetPhoneNumberCallbacks(&phoneNumberServiceCallbacks)

	rawResponse := callGraphQL(t, "contact/create_contact_min", map[string]interface{}{})

	var contactStruct struct {
		Contact_Create model.Contact
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)
	require.Equal(t, createdContactId, contactStruct.Contact_Create.ID)
	require.True(t, calledCreateContact)
	require.False(t, calledCreateEmail)
	require.False(t, calledCreatePhoneNumber)
}

func TestMutationResolver_ContactCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	createdContactId := uuid.New().String()
	createdEmailId := uuid.New().String()
	createdPhoneNumberId := uuid.New().String()

	calledCreateContact, calledCreateEmail, calledCreatePhoneNumber, calledLinkEmail, calledLinkPhoneNumber := false, false, false, false, false

	contactServiceCallbacks := events_platform.MockContactServiceCallbacks{
		CreateContact: func(context context.Context, contact *contactgrpc.UpsertContactGrpcRequest) (*contactgrpc.ContactIdGrpcResponse, error) {
			require.Equal(t, "MR", contact.Prefix)
			require.Equal(t, "first", contact.FirstName)
			require.Equal(t, "last", contact.LastName)
			require.Equal(t, "full name", contact.Name)
			require.Equal(t, "Some description", contact.Description)
			require.Equal(t, "America/Los_Angeles", contact.Timezone)
			require.Equal(t, "http://www.abc.com", contact.ProfilePhotoUrl)
			require.Equal(t, tenantName, contact.Tenant)
			require.Equal(t, string(entity.DataSourceOpenline), contact.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, contact.SourceFields.AppSource)
			require.Equal(t, testUserId, contact.LoggedInUserId)
			calledCreateContact = true
			neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
				Id: createdContactId,
			})
			return &contactgrpc.ContactIdGrpcResponse{
				Id: createdContactId,
			}, nil
		},
		LinkEmailToContact: func(context context.Context, link *contactgrpc.LinkEmailToContactGrpcRequest) (*contactgrpc.ContactIdGrpcResponse, error) {
			require.Equal(t, createdContactId, link.ContactId)
			require.Equal(t, createdEmailId, link.EmailId)
			require.Equal(t, true, link.Primary)
			require.Equal(t, "WORK", link.Label)
			require.Equal(t, tenantName, link.Tenant)
			require.Equal(t, testUserId, link.LoggedInUserId)
			calledLinkEmail = true
			neo4jt.LinkEmail(ctx, driver, createdContactId, createdEmailId, link.Primary, link.Label)
			return &contactgrpc.ContactIdGrpcResponse{
				Id: createdContactId,
			}, nil
		},
		LinkPhoneNumberToContact: func(context context.Context, link *contactgrpc.LinkPhoneNumberToContactGrpcRequest) (*contactgrpc.ContactIdGrpcResponse, error) {
			require.Equal(t, createdContactId, link.ContactId)
			require.Equal(t, createdPhoneNumberId, link.PhoneNumberId)
			require.Equal(t, true, link.Primary)
			require.Equal(t, "MOBILE", link.Label)
			require.Equal(t, tenantName, link.Tenant)
			require.Equal(t, testUserId, link.LoggedInUserId)
			calledLinkPhoneNumber = true
			neo4jt.LinkPhoneNumber(ctx, driver, createdContactId, createdPhoneNumberId, link.Primary, link.Label)
			return &contactgrpc.ContactIdGrpcResponse{
				Id: createdContactId,
			}, nil
		},
	}
	emailServiceCallbacks := events_platform.MockEmailServiceCallbacks{
		UpsertEmail: func(ctx context.Context, data *emailgrpc.UpsertEmailGrpcRequest) (*emailgrpc.EmailIdGrpcResponse, error) {
			require.Equal(t, "contact@abc.com", data.RawEmail)
			require.Equal(t, tenantName, data.Tenant)
			require.Equal(t, testUserId, data.LoggedInUserId)
			require.Equal(t, string(entity.DataSourceOpenline), data.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, data.SourceFields.AppSource)
			calledCreateEmail = true
			neo4jt.CreateEmail(ctx, driver, tenantName, entity.EmailEntity{
				Id:    createdEmailId,
				Email: "contact@abc.com",
			})
			return &emailgrpc.EmailIdGrpcResponse{
				Id: createdEmailId,
			}, nil
		},
	}
	phoneNumberServiceCallbacks := events_platform.MockPhoneNumberServiceCallbacks{
		UpsertPhoneNumber: func(ctx context.Context, data *phonenumbergrpc.UpsertPhoneNumberGrpcRequest) (*phonenumbergrpc.PhoneNumberIdGrpcResponse, error) {
			require.Equal(t, "+1234567890", data.PhoneNumber)
			require.Equal(t, tenantName, data.Tenant)
			require.Equal(t, testUserId, data.LoggedInUserId)
			require.Equal(t, string(entity.DataSourceOpenline), data.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, data.SourceFields.AppSource)
			calledCreatePhoneNumber = true
			neo4jt.CreatePhoneNumber(ctx, driver, tenantName, entity.PhoneNumberEntity{
				Id: createdPhoneNumberId,
			})
			return &phonenumbergrpc.PhoneNumberIdGrpcResponse{
				Id: createdPhoneNumberId,
			}, nil
		},
	}
	events_platform.SetContactCallbacks(&contactServiceCallbacks)
	events_platform.SetEmailCallbacks(&emailServiceCallbacks)
	events_platform.SetPhoneNumberCallbacks(&phoneNumberServiceCallbacks)

	rawResponse := callGraphQL(t, "contact/create_contact", map[string]interface{}{})

	var contactStruct struct {
		Contact_Create model.Contact
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)
	contact := contactStruct.Contact_Create
	require.Equal(t, createdContactId, contact.ID)
	require.Equal(t, 1, len(contact.Emails))
	require.Equal(t, createdEmailId, contact.Emails[0].ID)
	require.Equal(t, 1, len(contact.PhoneNumbers))
	require.Equal(t, createdPhoneNumberId, contact.PhoneNumbers[0].ID)
	require.True(t, calledCreateContact)
	require.True(t, calledCreateEmail)
	require.True(t, calledCreatePhoneNumber)
	require.True(t, calledLinkEmail)
	require.True(t, calledLinkPhoneNumber)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Contact":     1,
		"Email":       1,
		"PhoneNumber": 1,
	})
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

	contactServiceCallbacks := events_platform.MockContactServiceCallbacks{
		CreateContact: func(context context.Context, contact *contactgrpc.UpsertContactGrpcRequest) (*contactgrpc.ContactIdGrpcResponse, error) {
			require.Equal(t, "Bob", contact.FirstName)
			require.Equal(t, "Smith", contact.LastName)
			require.Equal(t, "Mr.", contact.Prefix)
			require.Equal(t, "This is a person", contact.Description)
			require.Equal(t, "unit-test", contact.SourceFields.AppSource)
			require.Equal(t, timeNow.Unix(), contact.CreatedAt.Seconds)
			require.Equal(t, "openline", contact.Tenant)
			calledCreateContact = true
			return &contactgrpc.ContactIdGrpcResponse{
				Id: createdContactId.String(),
			}, nil
		},
		LinkEmailToContact: func(context context.Context, link *contactgrpc.LinkEmailToContactGrpcRequest) (*contactgrpc.ContactIdGrpcResponse, error) {
			require.Equal(t, createdContactId.String(), link.ContactId)
			require.Equal(t, createdEmailId.String(), link.EmailId)
			require.Equal(t, true, link.Primary)
			require.Equal(t, "WORK", link.Label)
			require.Equal(t, "openline", link.Tenant)
			calledLinkEmailToContact = true
			return &contactgrpc.ContactIdGrpcResponse{
				Id: createdContactId.String(),
			}, nil
		},
	}
	events_platform.SetContactCallbacks(&contactServiceCallbacks)

	emailServiceCallbacks := events_platform.MockEmailServiceCallbacks{
		UpsertEmail: func(ctx context.Context, data *emailgrpc.UpsertEmailGrpcRequest) (*emailgrpc.EmailIdGrpcResponse, error) {
			require.Equal(t, "contact@abc.com", data.RawEmail)
			require.Equal(t, "openline", data.Tenant)
			calledCreateEmail = true
			return &emailgrpc.EmailIdGrpcResponse{
				Id: createdEmailId.String(),
			}, nil
		},
	}
	events_platform.SetEmailCallbacks(&emailServiceCallbacks)

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

func TestMutationResolver_ContactUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		Prefix:          "MR",
		FirstName:       "first",
		LastName:        "last",
		Description:     "description",
		ProfilePhotoUrl: "original url",
		Source:          entity.DataSourceHubspot,
		SourceOfTruth:   entity.DataSourceHubspot,
	})

	contactServiceCallbacks := events_platform.MockContactServiceCallbacks{
		CreateContact: func(context context.Context, contact *contactgrpc.UpsertContactGrpcRequest) (*contactgrpc.ContactIdGrpcResponse, error) {
			require.Equal(t, contactId, contact.Id)
			require.Equal(t, "updated first", contact.FirstName)
			require.Equal(t, "updated last", contact.LastName)
			require.Equal(t, "DR", contact.Prefix)
			require.Equal(t, "updated name", contact.Name)
			require.Equal(t, "updated description", contact.Description)
			require.Equal(t, "updated timezone", contact.Timezone)
			require.Equal(t, "http://updated.com", contact.ProfilePhotoUrl)
			require.Equal(t, constants.AppSourceCustomerOsApi, contact.SourceFields.AppSource)
			require.Equal(t, string(entity.DataSourceOpenline), contact.SourceFields.Source)
			require.Equal(t, tenantName, contact.Tenant)
			require.Equal(t, testUserId, contact.LoggedInUserId)
			return &contactgrpc.ContactIdGrpcResponse{
				Id: contactId,
			}, nil
		},
	}
	events_platform.SetContactCallbacks(&contactServiceCallbacks)

	rawResponse, err := c.RawPost(getQuery("contact/update_contact"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactStruct struct {
		Contact_Update model.Contact
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)
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
	noteId1 := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "note1", "text/plain", utils.Now())
	noteId2 := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "note2", "text/plain", utils.Now())
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
	require.Equal(t, "note1", *noteWithUser.Content)
	require.NotNil(t, noteWithUser.CreatedAt)
	require.NotNil(t, noteWithUser.CreatedBy)
	require.Equal(t, userId, noteWithUser.CreatedBy.ID)
	require.Equal(t, "first", noteWithUser.CreatedBy.FirstName)
	require.Equal(t, "last", noteWithUser.CreatedBy.LastName)

	require.Equal(t, noteId2, noteWithoutUser.ID)
	require.Equal(t, "note2", *noteWithoutUser.Content)
	require.NotNil(t, noteWithoutUser.CreatedAt)
	require.Nil(t, noteWithoutUser.CreatedBy)
}

func TestQueryResolver_Contact_WithNotes_ById_Time_Range(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	noteId1 := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "note1", "text/plain", utils.Now())
	noteId2 := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "note2", "text/plain", utils.Now())
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
	require.Equal(t, "note1", *noteWithUser.Content)
	require.NotNil(t, noteWithUser.CreatedAt)
	require.NotNil(t, noteWithUser.CreatedBy)
	require.Equal(t, userId, noteWithUser.CreatedBy.ID)
	require.Equal(t, "first", noteWithUser.CreatedBy.FirstName)
	require.Equal(t, "last", noteWithUser.CreatedBy.LastName)

	require.Equal(t, noteId2, noteWithoutUser.ID)
	require.Equal(t, "note2", *noteWithoutUser.Content)
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
	secAgo40 := now.Add(time.Duration(-40) * time.Second)
	secAgo50 := now.Add(time.Duration(-50) * time.Second)
	secAgo55 := now.Add(time.Duration(-55) * time.Second)
	secAgo60 := now.Add(time.Duration(-60) * time.Second)

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
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Analysis"))
	require.Equal(t, 6, neo4jt.GetCountOfNodes(ctx, driver, "TimelineEvent"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_timeline_events"),
		client.Var("contactId", contactId),
		client.Var("from", now),
		client.Var("size", 8))
	assertRawResponseSuccess(t, rawResponse, err)

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])

	timelineEvents := contact.(map[string]interface{})["timelineEvents"].([]interface{})
	require.Equal(t, 5, len(timelineEvents))

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

	timelineEvent3 := timelineEvents[2].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent3["__typename"].(string))
	require.Equal(t, interactionEventId1, timelineEvent3["id"].(string))
	require.NotNil(t, timelineEvent3["createdAt"].(string))
	require.Equal(t, "IE text 1", timelineEvent3["content"].(string))
	require.Equal(t, "application/json", timelineEvent3["contentType"].(string))
	require.Equal(t, "EMAIL", timelineEvent3["channel"].(string))

	timelineEvent4 := timelineEvents[3].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent4["__typename"].(string))
	require.Equal(t, interactionEventId2, timelineEvent4["id"].(string))
	require.NotNil(t, timelineEvent4["createdAt"].(string))
	require.Equal(t, "IE text 2", timelineEvent4["content"].(string))
	require.Equal(t, "application/json", timelineEvent4["contentType"].(string))
	require.Equal(t, "EMAIL", timelineEvent4["channel"].(string))

	timelineEvent5 := timelineEvents[4].(map[string]interface{})
	require.Equal(t, "Meeting", timelineEvent5["__typename"].(string))
	require.Equal(t, meetingId, timelineEvent5["id"].(string))
	require.NotNil(t, timelineEvent5["createdAt"].(string))
	require.Equal(t, "meeting-name", timelineEvent5["name"].(string))
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
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 1", "text/plain", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 2", "text/plain", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 3", "text/plain", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 4", "text/plain", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 5", "text/plain", now)

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
	require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 9, neo4jt.GetCountOfNodes(ctx, driver, "TimelineEvent"))

	rawResponse := callGraphQL(t, "contact/get_contact_with_timeline_events_total_count", map[string]interface{}{"contactId": contactId})

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])
	require.Equal(t, float64(4), contact.(map[string]interface{})["timelineEventsTotalCount"].(float64))
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

	organizationServiceCallbacks := events_platform.MockOrganizationServiceCallbacks{
		RefreshLastTouchpoint: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId1,
			}, nil
		},
	}
	events_platform.SetOrganizationCallbacks(&organizationServiceCallbacks)

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
