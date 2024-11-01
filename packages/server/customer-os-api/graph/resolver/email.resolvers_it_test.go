package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_GetEmail_WithParentOwners(t *testing.T) {
	ctx := context.Background()
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
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "test org1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "test org2")
	userId1 := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		FirstName: "a",
		LastName:  "b",
	})
	userId2 := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		FirstName: "c",
		LastName:  "d",
	})

	emailId := neo4jt.AddEmailTo(ctx, driver, commonModel.USER, tenantName, userId1, "test@openline.com", false, "WORK")
	neo4jt.AddEmailTo(ctx, driver, commonModel.USER, tenantName, userId2, "test@openline.com", false, "WORK")
	neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId1, "test@openline.com", false, "WORK")
	neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId2, "test@openline.com", false, "WORK")
	neo4jt.AddEmailTo(ctx, driver, commonModel.ORGANIZATION, tenantName, organizationId1, "test@openline.com", false, "WORK")
	neo4jt.AddEmailTo(ctx, driver, commonModel.ORGANIZATION, tenantName, organizationId2, "test@openline.com", false, "WORK")

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))

	rawResponse, err := c.RawPost(getQuery("email/get_email_with_parent_owners_via_organization_query"),
		client.Var("organizationId", organizationId1))
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationStruct struct {
		Organization model.Organization
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	require.Equal(t, 1, len(organizationStruct.Organization.Emails))

	email := organizationStruct.Organization.Emails[0]

	require.Equal(t, emailId, email.ID)
	require.Equal(t, 2, len(email.Users))
	require.Equal(t, 2, len(email.Contacts))
	require.Equal(t, 2, len(email.Organizations))
	require.Equal(t, userId1, email.Users[0].ID)
	require.Equal(t, userId2, email.Users[1].ID)
	require.Equal(t, contactId1, email.Contacts[0].ID)
	require.Equal(t, contactId2, email.Contacts[1].ID)
	require.Equal(t, organizationId1, email.Organizations[0].ID)
	require.Equal(t, organizationId2, email.Organizations[1].ID)
}

func TestQueryResolver_GetEmail_ById(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	emailId := neo4jtest.CreateEmail(ctx, driver, tenantName, neo4jentity.EmailEntity{
		Email:     "test@openline.ai",
		RawEmail:  "testRaw@openline.ai",
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
		IsRisky:   utils.BoolPtr(true),
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Email": 1, "Email_" + tenantName: 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, driver, map[string]int{"EMAIL_ADDRESS_BELONGS_TO_TENANT": 1})

	// Make the RawPost request and check for errors
	rawResponse := callGraphQL(t, "email/get_email", map[string]interface{}{"emailId": emailId})

	// Unmarshal the response data into the email struct
	var emailStruct struct {
		Email model.Email
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	email := emailStruct.Email

	require.Equal(t, emailId, email.ID)
	test.AssertRecentTime(t, email.UpdatedAt)
	test.AssertRecentTime(t, email.CreatedAt)
	require.Equal(t, "test@openline.ai", *email.Email)
	require.Equal(t, "testRaw@openline.ai", *email.RawEmail)
	require.Equal(t, true, *email.EmailValidationDetails.IsRisky)

	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Email", "Email_" + tenantName})
}
