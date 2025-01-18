package resolver

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	"github.com/customeros/customeros/packages/server/customer-os-api/test"
	neo4jt "github.com/customeros/customeros/packages/server/customer-os-api/test/neo4j"
	"github.com/customeros/customeros/packages/server/customer-os-api/utils/decode"
)

func TestMutationResolver_PhoneNumberMergeToContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	// Create a default contact
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jtest.CreateCountry(ctx, driver, neo4jentity.CountryEntity{
		Name:      "United States",
		CodeA2:    "US",
		CodeA3:    "USA",
		PhoneCode: "1",
	})

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/merge_phone_number_to_contact"),
		client.Var("contactId", contactId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumberStruct struct {
		PhoneNumberMergeToContact model.PhoneNumber
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err)
	assert.NotNil(t, phoneNumberStruct.PhoneNumberMergeToContact.ID)
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
	neo4jtest.CreateCountry(ctx, driver, neo4jentity.CountryEntity{
		Name:      "United States",
		CodeA2:    "US",
		CodeA3:    "USA",
		PhoneCode: "1",
	})

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("phone_number/merge_phone_number_to_organization"),
		client.Var("organizationId", organizationId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumberStruct struct {
		PhoneNumberMergeToOrganization model.PhoneNumber
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumberStruct)
	require.Nil(t, err)
	assert.NotNil(t, phoneNumberStruct.PhoneNumberMergeToOrganization.ID)
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
