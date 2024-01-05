package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_PhoneNumberMergeToContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create a default contact
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.CreateCountry(ctx, driver, "US", "USA", "United States", "1")

	rawResponse, err := c.RawPost(getQuery("phone_number/merge_phone_number_to_contact"),
		client.Var("contactId", contactId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumberStruct struct {
		PhoneNumberMergeToContact model.PhoneNumber
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	createdPhoneNumber := phoneNumberStruct.PhoneNumberMergeToContact
	// Check that the fields of the phoneNumber struct have the expected values
	require.NotNil(t, createdPhoneNumber.CreatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPhoneNumber.CreatedAt)
	require.NotNil(t, createdPhoneNumber.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPhoneNumber.UpdatedAt)

	require.NotNil(t, createdPhoneNumber.ID, "PhoneNumber ID is nil")
	require.Equal(t, true, createdPhoneNumber.Primary, "PhoneNumber Primary field is not true")
	require.Nil(t, createdPhoneNumber.Validated, "PhoneNumber Validated field is not nil")
	require.Nil(t, createdPhoneNumber.E164)
	require.Equal(t, "+1234567890", *createdPhoneNumber.RawPhoneNumber, "PhoneNumber E164 field is not expected value")
	if createdPhoneNumber.Label == nil {
		t.Errorf("PhoneNumber Label field is nil")
	} else {
		require.Equal(t, model.PhoneNumberLabelWork, *createdPhoneNumber.Label, "PhoneNumber Label field is not expected value")
	}
	require.Equal(t, model.DataSourceOpenline, createdPhoneNumber.Source, "PhoneNumber Source field is not expected value")
	require.NotNil(t, createdPhoneNumber.Country)
	require.Equal(t, "US", createdPhoneNumber.Country.CodeA2)
	require.Equal(t, "USA", createdPhoneNumber.Country.CodeA3)
	require.Equal(t, "United States", createdPhoneNumber.Country.Name)
	require.Equal(t, "1", createdPhoneNumber.Country.PhoneCode)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Country"), "Incorrect number of Country nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"), "Incorrect number of Contact nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"), "Incorrect number of PhoneNumber nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber_"+tenantName), "Incorrect number of PhoneNumber_%s nodes in Neo4j", tenantName)
	require.Equal(t, 4, neo4jt.GetTotalCountOfNodes(ctx, driver), "Incorrect total number of nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of PHONE_ASSOCIATED_WITH relationships in Neo4j")

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Country", "Contact", "Contact_" + tenantName, "PhoneNumber", "PhoneNumber_" + tenantName})
}

func TestMutationResolver_PhoneNumberUpdateInContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create a default contact and phone number
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId, "+1234567890", false, "WORK")

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/update_phone_number_for_contact"),
		client.Var("contactId", contactId),
		client.Var("phoneNumberId", phoneNumberId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the phone number struct
	var phoneNumberStruct struct {
		PhoneNumberUpdateInContact model.PhoneNumber
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	phoneNumber := phoneNumberStruct.PhoneNumberUpdateInContact

	// Check that the fields of the phone struct have the expected values
	require.Equal(t, phoneNumberId, phoneNumber.ID, "Phone number ID is nil")
	require.Equal(t, true, phoneNumber.Primary, "Phone number Primary field is not true")
	require.Equal(t, "+1234567890", *phoneNumber.RawPhoneNumber, "Phone number expected not to be changed")
	require.Equal(t, "+1234567890", *phoneNumber.E164, "Phone number expected not to be changed")
	require.NotNil(t, phoneNumber.UpdatedAt, "Missing updatedAt field")
	if phoneNumber.Label == nil {
		t.Errorf("Phone number Label field is nil")
	} else {
		require.Equal(t, model.PhoneNumberLabelHome, *phoneNumber.Label, "Phone number label field is not expected value")
	}

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"), "Incorrect number of PhoneNumber nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")
}

func TestMutationResolver_PhoneNumberUpdateInContact_ReplacePhoneNumber(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create a default contact and phone number
	neo4jt.CreateCountry(ctx, driver, "US", "USA", "United States", "1")
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId, "+1234567890", false, "WORK")

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/replace_phone_number_for_contact"),
		client.Var("contactId", contactId),
		client.Var("phoneNumberId", phoneNumberId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the phone number struct
	var phoneNumberStruct struct {
		PhoneNumberUpdateInContact model.PhoneNumber
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	phoneNumber := phoneNumberStruct.PhoneNumberUpdateInContact

	// Check that the fields of the phoneNumberStruct struct have the expected values
	require.NotEqual(t, phoneNumberId, phoneNumber.ID, "Expected new phone number ID to be generated")
	require.Equal(t, true, phoneNumber.Primary, "Phone number primary field is not true")
	require.Equal(t, "+987654321", *phoneNumber.RawPhoneNumber)
	require.Nil(t, phoneNumber.E164)
	require.Nil(t, phoneNumber.Validated, "New phone number is not nil")
	require.NotNil(t, phoneNumber.CreatedAt, "Missing createdAt field")
	require.NotNil(t, phoneNumber.UpdatedAt, "Missing updatedAt field")
	if phoneNumber.Label == nil {
		t.Errorf("Phone number Label field is nil")
	} else {
		require.Equal(t, model.PhoneNumberLabelHome, *phoneNumber.Label, "Phone number label field is not expected value")
	}
	require.NotNil(t, phoneNumber.Country)
	require.Equal(t, "US", phoneNumber.Country.CodeA2)
	require.Equal(t, "USA", phoneNumber.Country.CodeA3)
	require.Equal(t, "United States", phoneNumber.Country.Name)
	require.Equal(t, "1", phoneNumber.Country.PhoneCode)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Country"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"), "Expected 2 PhoneNumber nodes, original one and new")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "LINKED_TO"))

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Country", "Contact", "Contact_" + tenantName, "PhoneNumber", "PhoneNumber_" + tenantName})
}

func TestMutationResolver_PhoneNumberRemoveFromContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId, "+1234567890", false, "WORK")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))

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

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "PhoneNumber", "PhoneNumber_" + tenantName, "Contact", "Contact_" + tenantName})
}

func TestMutationResolver_PhoneNumberMergeToOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "Edgeless Systems")

	organizationServiceCallbacks := events_platform.MockOrganizationServiceCallbacks{
		RefreshLastTouchpoint: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: organizationId,
			}, nil
		},
	}
	events_platform.SetOrganizationCallbacks(&organizationServiceCallbacks)

	rawResponse, err := c.RawPost(getQuery("phone_number/merge_phone_number_to_organization"),
		client.Var("organizationId", organizationId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumberStruct struct {
		PhoneNumberMergeToOrganization model.PhoneNumber
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	createdPhoneNumber := phoneNumberStruct.PhoneNumberMergeToOrganization
	// Check that the fields of the phoneNumber struct have the expected values
	require.NotNil(t, createdPhoneNumber.CreatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPhoneNumber.CreatedAt)
	require.NotNil(t, createdPhoneNumber.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPhoneNumber.UpdatedAt)

	require.NotNil(t, createdPhoneNumber.ID, "PhoneNumber ID is nil")
	require.Equal(t, true, createdPhoneNumber.Primary, "PhoneNumber Primary field is not true")
	require.Nil(t, createdPhoneNumber.Validated, "PhoneNumber Validated field is not nil")
	require.Nil(t, createdPhoneNumber.E164)
	require.Equal(t, "+1234567890", *createdPhoneNumber.RawPhoneNumber, "PhoneNumber E164 field is not expected value")
	if createdPhoneNumber.Label == nil {
		t.Errorf("PhoneNumber Label field is nil")
	} else {
		require.Equal(t, model.PhoneNumberLabelWork, *createdPhoneNumber.Label, "PhoneNumber Label field is not expected value")
	}
	require.Equal(t, model.DataSourceOpenline, createdPhoneNumber.Source, "PhoneNumber Source field is not expected value")

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"), "Incorrect number of Contact nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"), "Incorrect number of PhoneNumber nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber_"+tenantName), "Incorrect number of PhoneNumber_%s nodes in Neo4j", tenantName)
	require.Equal(t, 3, neo4jt.GetTotalCountOfNodes(ctx, driver), "Incorrect total number of nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of PHONE_ASSOCIATED_WITH relationships in Neo4j")

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Organization", "Organization_" + tenantName, "PhoneNumber", "PhoneNumber_" + tenantName})
}

func TestMutationResolver_PhoneNumberUpdateInOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "Edgeless Systems")
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, organizationId, "+1234567890", false, "WORK")

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/update_phone_number_for_organization"),
		client.Var("organizationId", organizationId),
		client.Var("phoneNumberId", phoneNumberId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the phone number struct
	var phoneNumberStruct struct {
		PhoneNumberUpdateInOrganization model.PhoneNumber
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	phoneNumber := phoneNumberStruct.PhoneNumberUpdateInOrganization

	// Check that the fields of the phone struct have the expected values
	require.Equal(t, phoneNumberId, phoneNumber.ID, "Phone number ID is nil")
	require.Equal(t, true, phoneNumber.Primary, "Phone number Primary field is not true")
	require.Equal(t, "+1234567890", *phoneNumber.RawPhoneNumber, "Phone number expected not to be changed")
	require.Equal(t, "+1234567890", *phoneNumber.E164, "Phone number expected not to be changed")
	require.NotNil(t, phoneNumber.UpdatedAt, "Missing updatedAt field")
	if phoneNumber.Label == nil {
		t.Errorf("Phone number Label field is nil")
	} else {
		require.Equal(t, model.PhoneNumberLabelHome, *phoneNumber.Label, "Phone number label field is not expected value")
	}

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"), "Incorrect number of PhoneNumber nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")
}

func TestMutationResolver_PhoneNumberUpdateInOrganization_ReplacePhoneNumber(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "Edgeless Systems")
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, organizationId, "+1234567890", false, "WORK")

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/replace_phone_number_for_organization"),
		client.Var("organizationId", organizationId),
		client.Var("phoneNumberId", phoneNumberId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the phone number struct
	var phoneNumberStruct struct {
		PhoneNumberUpdateInOrganization model.PhoneNumber
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	phoneNumber := phoneNumberStruct.PhoneNumberUpdateInOrganization

	// Check that the fields of the phoneNumberStruct struct have the expected values
	require.NotEqual(t, phoneNumberId, phoneNumber.ID, "Expected new phone number ID to be generated")
	require.Equal(t, true, phoneNumber.Primary, "Phone number primary field is not true")
	require.Equal(t, "+987654321", *phoneNumber.RawPhoneNumber)
	require.Nil(t, phoneNumber.E164)
	require.Nil(t, phoneNumber.Validated, "New phone number is not validated yet")
	require.NotNil(t, phoneNumber.CreatedAt, "Missing createdAt field")
	require.NotNil(t, phoneNumber.UpdatedAt, "Missing updatedAt field")
	if phoneNumber.Label == nil {
		t.Errorf("Phone number Label field is nil")
	} else {
		require.Equal(t, model.PhoneNumberLabelHome, *phoneNumber.Label, "Phone number label field is not expected value")
	}

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"), "Expected 2 PhoneNumber nodes, original one and new")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Organization", "Organization_" + tenantName, "PhoneNumber", "PhoneNumber_" + tenantName})
}

func TestMutationResolver_PhoneNumberRemoveFromOrganizationByID(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "Edgeless Systems")
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, organizationId, "+1234567890", false, "WORK")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/remove_phone_number_from_organization_by_id"),
		client.Var("organizationId", organizationId),
		client.Var("phoneNumberId", phoneNumberId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumberStruct struct {
		PhoneNumberRemoveFromOrganizationById model.Result
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	require.Equal(t, true, phoneNumberStruct.PhoneNumberRemoveFromOrganizationById.Result)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "PhoneNumber", "PhoneNumber_" + tenantName, "Organization", "Organization_" + tenantName})
}

func TestMutationResolver_PhoneNumberMergeToUser(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("phone_number/merge_phone_number_to_user"),
		client.Var("userId", userId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumberStruct struct {
		PhoneNumberMergeToUser model.PhoneNumber
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	createdPhoneNumber := phoneNumberStruct.PhoneNumberMergeToUser
	// Check that the fields of the phoneNumber struct have the expected values
	require.NotNil(t, createdPhoneNumber.CreatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPhoneNumber.CreatedAt)
	require.NotNil(t, createdPhoneNumber.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPhoneNumber.UpdatedAt)

	require.NotNil(t, createdPhoneNumber.ID, "PhoneNumber ID is nil")
	require.Equal(t, true, createdPhoneNumber.Primary, "PhoneNumber Primary field is not true")
	require.Nil(t, createdPhoneNumber.Validated, "PhoneNumber Validated field is not nil")
	require.Nil(t, createdPhoneNumber.E164)
	require.Equal(t, "+1234567890", *createdPhoneNumber.RawPhoneNumber, "PhoneNumber E164 field is not expected value")
	if createdPhoneNumber.Label == nil {
		t.Errorf("PhoneNumber Label field is nil")
	} else {
		require.Equal(t, model.PhoneNumberLabelWork, *createdPhoneNumber.Label, "PhoneNumber Label field is not expected value")
	}
	require.Equal(t, model.DataSourceOpenline, createdPhoneNumber.Source, "PhoneNumber Source field is not expected value")

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"), "Incorrect number of Contact nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"), "Incorrect number of PhoneNumber nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber_"+tenantName), "Incorrect number of PhoneNumber_%s nodes in Neo4j", tenantName)
	require.Equal(t, 3, neo4jt.GetTotalCountOfNodes(ctx, driver), "Incorrect total number of nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of PHONE_ASSOCIATED_WITH relationships in Neo4j")

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "User", "User_" + tenantName, "PhoneNumber", "PhoneNumber_" + tenantName})
}

func TestMutationResolver_PhoneNumberUpdateInUser(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, userId, "+1234567890", false, "WORK")

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/update_phone_number_for_user"),
		client.Var("userId", userId),
		client.Var("phoneNumberId", phoneNumberId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the phone number struct
	var phoneNumberStruct struct {
		PhoneNumberUpdateInUser model.PhoneNumber
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	phoneNumber := phoneNumberStruct.PhoneNumberUpdateInUser

	// Check that the fields of the phone struct have the expected values
	require.Equal(t, phoneNumberId, phoneNumber.ID, "Phone number ID is nil")
	require.Equal(t, true, phoneNumber.Primary, "Phone number Primary field is not true")
	require.Equal(t, "+1234567890", *phoneNumber.RawPhoneNumber, "Phone number expected not to be changed")
	require.Equal(t, "+1234567890", *phoneNumber.E164, "Phone number expected not to be changed")
	require.NotNil(t, phoneNumber.UpdatedAt, "Missing updatedAt field")
	if phoneNumber.Label == nil {
		t.Errorf("Phone number Label field is nil")
	} else {
		require.Equal(t, model.PhoneNumberLabelHome, *phoneNumber.Label, "Phone number label field is not expected value")
	}

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"), "Incorrect number of PhoneNumber nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")
}

func TestMutationResolver_PhoneNumberUpdateInUser_ReplacePhoneNumber(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, userId, "+1234567890", false, "WORK")

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/replace_phone_number_for_user"),
		client.Var("userId", userId),
		client.Var("phoneNumberId", phoneNumberId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the phone number struct
	var phoneNumberStruct struct {
		PhoneNumberUpdateInUser model.PhoneNumber
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	phoneNumber := phoneNumberStruct.PhoneNumberUpdateInUser

	// Check that the fields of the phoneNumberStruct struct have the expected values
	require.NotEqual(t, phoneNumberId, phoneNumber.ID, "Expected new phone number ID to be generated")
	require.Equal(t, true, phoneNumber.Primary, "Phone number primary field is not true")
	require.Equal(t, "+987654321", *phoneNumber.RawPhoneNumber)
	require.Nil(t, phoneNumber.E164)
	require.Nil(t, phoneNumber.Validated, "New phone number is not validated yet")
	require.NotNil(t, phoneNumber.CreatedAt, "Missing createdAt field")
	require.NotNil(t, phoneNumber.UpdatedAt, "Missing updatedAt field")
	if phoneNumber.Label == nil {
		t.Errorf("Phone number Label field is nil")
	} else {
		require.Equal(t, model.PhoneNumberLabelHome, *phoneNumber.Label, "Phone number label field is not expected value")
	}

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"), "Expected 2 PhoneNumber nodes, original one and new")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "User", "User_" + tenantName, "PhoneNumber", "PhoneNumber_" + tenantName})
}

func TestQueryResolver_GetPhoneNumber_WithParentOwners(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId1 := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		FirstName: "a",
		LastName:  "b",
	})
	contactId2 := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{
		FirstName: "c",
		LastName:  "d",
	})
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "test org1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "test org2")
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "a",
		LastName:  "b",
	})
	userId2 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "c",
		LastName:  "d",
	})

	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, userId1, "+12345", false, "WORK")
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, userId2, "+12345", false, "WORK")
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId1, "+12345", false, "WORK")
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId2, "+12345", false, "WORK")
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, organizationId1, "+12345", false, "WORK")
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, organizationId2, "+12345", false, "WORK")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))

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

	neo4jt.CreateTenant(ctx, driver, tenantName)

	phoneNumberId := neo4jt.CreatePhoneNumber(ctx, driver, tenantName, entity.PhoneNumberEntity{
		E164:           "+123456789",
		RawPhoneNumber: "+ 123 456 789",
		CreatedAt:      utils.Now(),
		UpdatedAt:      utils.Now(),
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"PhoneNumber": 1, "PhoneNumber_" + tenantName: 1})
	assertNeo4jRelationCount(ctx, t, driver, map[string]int{"PHONE_NUMBER_BELONGS_TO_TENANT": 1})

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
