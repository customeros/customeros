package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_PhoneNumberMergeToContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	// Create a default contact
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	phoneNumberId := uuid.New().String()
	neo4jtest.CreateCountry(ctx, driver, neo4jentity.CountryEntity{
		Name:      "United States",
		CodeA2:    "US",
		CodeA3:    "USA",
		PhoneCode: "1",
	})

	phoneNumberServiceCalled := false
	contactServiceCalled := false

	phoneNumberServiceCallbacks := events_platform.MockPhoneNumberServiceCallbacks{
		UpsertPhoneNumber: func(ctx context.Context, phoneNumber *phonenumberpb.UpsertPhoneNumberGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
			require.Equal(t, tenantName, phoneNumber.Tenant)
			require.NotNil(t, phoneNumber)
			phoneNumberServiceCalled = true
			neo4jtest.CreatePhoneNumber(ctx, driver, tenantName, neo4jentity.PhoneNumberEntity{
				Id:             phoneNumberId,
				RawPhoneNumber: "+1234567890",
				CreatedAt:      utils.Now(),
				UpdatedAt:      utils.Now(),
			})
			return &phonenumberpb.PhoneNumberIdGrpcResponse{
				Id: phoneNumberId,
			}, nil
		},
	}
	events_platform.SetPhoneNumberCallbacks(&phoneNumberServiceCallbacks)

	contactServiceCallbacks := events_platform.MockContactServiceCallbacks{
		LinkPhoneNumberToContact: func(context context.Context, contact *contactpb.LinkPhoneNumberToContactGrpcRequest) (*contactpb.ContactIdGrpcResponse, error) {
			require.Equal(t, tenantName, contact.Tenant)
			require.Equal(t, contactId, contact.ContactId)
			require.Equal(t, phoneNumberId, contact.PhoneNumberId)
			contactServiceCalled = true
			return &contactpb.ContactIdGrpcResponse{
				Id: contactId,
			}, nil
		},
	}
	events_platform.SetContactCallbacks(&contactServiceCallbacks)

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/merge_phone_number_to_contact"),
		client.Var("contactId", contactId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumberStruct struct {
		PhoneNumberMergeToContact model.PhoneNumber
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	require.True(t, phoneNumberServiceCalled, "Phone number service was not called")
	require.True(t, contactServiceCalled, "Contact service was not called")
}

func TestMutationResolver_PhoneNumberRemoveFromContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId, "+1234567890", false, "WORK")

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "PhoneNumber_"+tenantName))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "HAS"))

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/remove_phone_number_from_contact"),
		client.Var("contactId", contactId),
		client.Var("e164", "+1234567890"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumberStruct struct {
		PhoneNumberRemoveFromContactByE164 model.Result
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	require.Equal(t, true, phoneNumberStruct.PhoneNumberRemoveFromContactByE164.Result)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "PhoneNumber_"+tenantName))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 0, neo4jtest.GetCountOfRelationships(ctx, driver, "HAS"))
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "PhoneNumber", "PhoneNumber_" + tenantName, "Contact", "Contact_" + tenantName})
}

func TestMutationResolver_PhoneNumberMergeToOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	// Create a default organization
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	phoneNumberId := uuid.New().String()
	neo4jtest.CreateCountry(ctx, driver, neo4jentity.CountryEntity{
		Name:      "United States",
		CodeA2:    "US",
		CodeA3:    "USA",
		PhoneCode: "1",
	})

	phoneNumberServiceCalled := false
	organizationServiceCalled := false

	phoneNumberServiceCallbacks := events_platform.MockPhoneNumberServiceCallbacks{
		UpsertPhoneNumber: func(ctx context.Context, phoneNumber *phonenumberpb.UpsertPhoneNumberGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
			require.Equal(t, tenantName, phoneNumber.Tenant)
			require.NotNil(t, phoneNumber)
			phoneNumberServiceCalled = true
			neo4jtest.CreatePhoneNumber(ctx, driver, tenantName, neo4jentity.PhoneNumberEntity{
				Id:             phoneNumberId,
				RawPhoneNumber: "+1234567890",
				CreatedAt:      utils.Now(),
				UpdatedAt:      utils.Now(),
			})
			return &phonenumberpb.PhoneNumberIdGrpcResponse{
				Id: phoneNumberId,
			}, nil
		},
	}
	events_platform.SetPhoneNumberCallbacks(&phoneNumberServiceCallbacks)

	organizationServiceCallbacks := events_platform.MockOrganizationServiceCallbacks{
		LinkPhoneNumberToOrganization: func(context context.Context, organization *organizationpb.LinkPhoneNumberToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, organization.Tenant)
			require.Equal(t, organizationId, organization.OrganizationId)
			require.Equal(t, phoneNumberId, organization.PhoneNumberId)
			organizationServiceCalled = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: organizationId,
			}, nil
		},
	}
	events_platform.SetOrganizationCallbacks(&organizationServiceCallbacks)

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/merge_phone_number_to_organization"),
		client.Var("organizationId", organizationId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumberStruct struct {
		PhoneNumberMergeToOrganization model.PhoneNumber
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	require.True(t, phoneNumberServiceCalled, "Phone number service was not called")
	require.True(t, organizationServiceCalled, "Organization service was not called")
}

func TestMutationResolver_PhoneNumberMergeToUser(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	// Create a default user
	userId := neo4jtest.CreateDefaultUser(ctx, driver, tenantName)
	phoneNumberId := uuid.New().String()
	neo4jtest.CreateCountry(ctx, driver, neo4jentity.CountryEntity{
		Name:      "United States",
		CodeA2:    "US",
		CodeA3:    "USA",
		PhoneCode: "1",
	})

	phoneNumberServiceCalled := false
	userServiceCalled := false

	phoneNumberServiceCallbacks := events_platform.MockPhoneNumberServiceCallbacks{
		UpsertPhoneNumber: func(ctx context.Context, phoneNumber *phonenumberpb.UpsertPhoneNumberGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
			require.Equal(t, tenantName, phoneNumber.Tenant)
			require.NotNil(t, phoneNumber)
			phoneNumberServiceCalled = true
			neo4jtest.CreatePhoneNumber(ctx, driver, tenantName, neo4jentity.PhoneNumberEntity{
				Id:             phoneNumberId,
				RawPhoneNumber: "+1234567890",
				CreatedAt:      utils.Now(),
				UpdatedAt:      utils.Now(),
			})
			return &phonenumberpb.PhoneNumberIdGrpcResponse{
				Id: phoneNumberId,
			}, nil
		},
	}
	events_platform.SetPhoneNumberCallbacks(&phoneNumberServiceCallbacks)

	userServiceCallbacks := events_platform.MockUserServiceCallbacks{
		LinkPhoneNumberToUser: func(context context.Context, user *userpb.LinkPhoneNumberToUserGrpcRequest) (*userpb.UserIdGrpcResponse, error) {
			require.Equal(t, tenantName, user.Tenant)
			require.Equal(t, userId, user.UserId)
			require.Equal(t, phoneNumberId, user.PhoneNumberId)
			userServiceCalled = true
			return &userpb.UserIdGrpcResponse{
				Id: userId,
			}, nil
		},
	}
	events_platform.SetUserCallbacks(&userServiceCallbacks)

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/merge_phone_number_to_user"),
		client.Var("userId", userId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumberStruct struct {
		PhoneNumberMergeToUser model.PhoneNumber
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	require.True(t, phoneNumberServiceCalled, "Phone number service was not called")
	require.True(t, userServiceCalled, "User service was not called")
}

func TestQueryResolver_GetPhoneNumber_WithParentOwners(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId1 := neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{
		FirstName: "a",
		LastName:  "b",
	})
	contactId2 := neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{
		FirstName: "c",
		LastName:  "d",
	})
	organizationId1 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "test org1"})
	organizationId2 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "test org2"})
	userId1 := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		FirstName: "a",
		LastName:  "b",
	})
	userId2 := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		FirstName: "c",
		LastName:  "d",
	})

	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, userId1, "+12345", false, "WORK")
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, userId2, "+12345", false, "WORK")
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId1, "+12345", false, "WORK")
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId2, "+12345", false, "WORK")
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, organizationId1, "+12345", false, "WORK")
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, organizationId2, "+12345", false, "WORK")

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))

	rawResponse, err := c.RawPost(getQuery("phone_number/get_phone_number_with_parent_owners_via_organization_query"),
		client.Var("organizationId", organizationId1))
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationStruct struct {
		Organization model.Organization
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	require.Equal(t, 1, len(organizationStruct.Organization.PhoneNumbers))

	phoneNumber := organizationStruct.Organization.PhoneNumbers[0]

	require.Equal(t, phoneNumberId, phoneNumber.ID)
	require.Equal(t, 2, len(phoneNumber.Users))
	require.Equal(t, 2, len(phoneNumber.Contacts))
	require.Equal(t, 2, len(phoneNumber.Organizations))
	require.Equal(t, userId1, phoneNumber.Users[0].ID)
	require.Equal(t, userId2, phoneNumber.Users[1].ID)
	require.Equal(t, contactId1, phoneNumber.Contacts[0].ID)
	require.Equal(t, contactId2, phoneNumber.Contacts[1].ID)
	require.Equal(t, organizationId1, phoneNumber.Organizations[0].ID)
	require.Equal(t, organizationId2, phoneNumber.Organizations[1].ID)
}

func TestQueryResolver_GetPhoneNumber_ById(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	phoneNumberId := neo4jtest.CreatePhoneNumber(ctx, driver, tenantName, neo4jentity.PhoneNumberEntity{
		E164:           "+123456789",
		RawPhoneNumber: "+ 123 456 789",
		CreatedAt:      utils.Now(),
		UpdatedAt:      utils.Now(),
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"PhoneNumber": 1, "PhoneNumber_" + tenantName: 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, driver, map[string]int{"PHONE_NUMBER_BELONGS_TO_TENANT": 1})

	// Make the RawPost request and check for errors
	rawResponse := callGraphQL(t, "phone_number/get_phone_number", map[string]interface{}{"phoneNumberId": phoneNumberId})

	// Unmarshal the response data into the phone number struct
	var phoneNumberStruct struct {
		PhoneNumber model.PhoneNumber
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	phoneNumber := phoneNumberStruct.PhoneNumber

	require.Equal(t, phoneNumberId, phoneNumber.ID)
	test.AssertRecentTime(t, phoneNumber.UpdatedAt)
	test.AssertRecentTime(t, phoneNumber.CreatedAt)
	require.Equal(t, "+123456789", *phoneNumber.E164)
	require.Equal(t, "+ 123 456 789", *phoneNumber.RawPhoneNumber)

	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "PhoneNumber", "PhoneNumber_" + tenantName})
}
