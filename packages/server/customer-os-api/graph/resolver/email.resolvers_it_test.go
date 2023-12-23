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
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_EmailMergeToContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create a default contact
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("email/merge_email_to_contact"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var emailStruct struct {
		EmailMergeToContact model.Email
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	email := emailStruct.EmailMergeToContact

	// Check that the fields of the email struct have the expected values
	require.NotNil(t, email.ID, "Email ID is nil")
	require.NotNil(t, email.CreatedAt, "Missing createdAt field")
	require.NotNil(t, email.UpdatedAt, "Missing updatedAt field")
	require.Equal(t, true, email.Primary, "Email Primary field is not true")
	require.Equal(t, "test@gmail.com", *email.Email)
	require.Equal(t, "test@gmail.com", *email.RawEmail)
	require.Nil(t, email.EmailValidationDetails.Validated)
	if email.Label == nil {
		t.Errorf("Email Label field is nil")
	} else {
		require.Equal(t, model.EmailLabelWork, *email.Label, "Email Label field is not expected value")
	}
	require.Equal(t, model.DataSourceOpenline, email.Source, "Email Source field is not expected value")
	require.Equal(t, model.DataSourceOpenline, email.SourceOfTruth, "Email Source of truth field is not expected value")
	require.Equal(t, "test", email.AppSource, "Email App source field is not expected value")

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"), "Incorrect number of Contact nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"), "Incorrect number of Email nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName), "Incorrect number of Email_%s nodes in Neo4j", tenantName)
	require.Equal(t, 3, neo4jt.GetTotalCountOfNodes(ctx, driver), "Incorrect total number of nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "Email", "Email_" + tenantName})
}

func TestMutationResolver_EmailMergeToContact_SecondEmail(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create a default contact
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "original@email.com", true, "")

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("email/merge_second_email_to_contact"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var emailStruct struct {
		EmailMergeToContact model.Email
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	email := emailStruct.EmailMergeToContact

	// Check that the fields of the email struct have the expected values
	require.NotNil(t, email.ID, "Email ID is nil")
	require.NotEqual(t, emailId, email.ID)
	require.NotNil(t, email.CreatedAt, "Missing createdAt field")
	require.NotNil(t, email.UpdatedAt, "Missing updatedAt field")
	require.Equal(t, true, email.Primary, "Email Primary field is not true")
	require.Nil(t, email.Email)
	require.Nil(t, email.RawEmail)
	require.Nil(t, email.EmailValidationDetails.Validated)
	if email.Label == nil {
		t.Errorf("Email Label field is nil")
	} else {
		require.Equal(t, model.EmailLabelWork, *email.Label, "Email Label field is not expected value")
	}
	require.Equal(t, model.DataSourceOpenline, email.Source, "Email Source field is not expected value")
	require.Equal(t, model.DataSourceOpenline, email.SourceOfTruth, "Email Source of truth field is not expected value")
	require.Equal(t, "test", email.AppSource, "Email App source field is not expected value")

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"), "Incorrect number of Contact nodes in Neo4j")
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Email"), "Incorrect number of Email nodes in Neo4j")
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName), "Incorrect number of Email_%s nodes in Neo4j", tenantName)
	require.Equal(t, 4, neo4jt.GetTotalCountOfNodes(ctx, driver), "Incorrect total number of nodes in Neo4j")
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "Email", "Email_" + tenantName})
}

func TestMutationResolver_EmailUpdateInContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create a default contact and email
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "original@email.com", true, "")

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("email/update_email_for_contact"),
		client.Var("contactId", contactId),
		client.Var("emailId", emailId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var emailStruct struct {
		EmailUpdateInContact model.Email
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	email := emailStruct.EmailUpdateInContact

	// Check that the fields of the email struct have the expected values
	require.Equal(t, emailId, email.ID, "Email ID is nil")
	require.Equal(t, true, email.Primary, "Email Primary field is not true")
	require.Equal(t, "original@email.com", *email.RawEmail, "Email address expected not to be changed")
	require.Equal(t, "original@email.com", *email.Email, "Email address expected not to be changed")
	require.NotNil(t, email.UpdatedAt, "Missing updatedAt field")
	if email.Label == nil {
		t.Errorf("Email Label field is nil")
	} else {
		require.Equal(t, model.EmailLabelPersonal, *email.Label, "Email Label field is not expected value")
	}

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"), "Incorrect number of Email nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")
}

func TestMutationResolver_EmailUpdateInContact_ReplaceEmail(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create a default contact and email
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "original@email.com", true, "")

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("email/replace_email_for_contact"),
		client.Var("contactId", contactId),
		client.Var("emailId", emailId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the emailStruct struct
	var emailStruct struct {
		EmailUpdateInContact model.Email
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	email := emailStruct.EmailUpdateInContact

	// Check that the fields of the emailStruct struct have the expected values
	require.NotEqual(t, emailId, email.ID, "Expected new email id to be generated")
	require.Equal(t, true, email.Primary, "Email Primary field is not true")
	require.Equal(t, "new@email.com", *email.RawEmail)
	require.Equal(t, "new@email.com", *email.Email)
	require.Nil(t, email.EmailValidationDetails.Validated)
	require.NotNil(t, email.CreatedAt, "Missing createdAt field")
	require.NotNil(t, email.UpdatedAt, "Missing updatedAt field")
	if email.Label == nil {
		t.Errorf("Email Label field is nil")
	} else {
		require.Equal(t, model.EmailLabelPersonal, *email.Label, "Email Label field is not expected value")
	}

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Email"), "Expected 2 email nodes, original one and new")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "Email", "Email_" + tenantName})
}

func TestMutationResolver_EmailUpdateInUser(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create a default contact and email
	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId, "original@email.com", true, "")

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("email/update_email_for_user"),
		client.Var("userId", userId),
		client.Var("emailId", emailId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var emailStruct struct {
		EmailUpdateInUser model.Email
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	email := emailStruct.EmailUpdateInUser

	// Check that the fields of the email struct have the expected values
	require.Equal(t, emailId, email.ID, "Email ID is nil")
	require.Equal(t, true, email.Primary, "Email Primary field is not true")
	require.Equal(t, "original@email.com", *email.Email, "Email address expected not to be changed")
	require.Equal(t, "original@email.com", *email.RawEmail, "Email address expected not to be changed")
	require.NotNil(t, email.UpdatedAt, "Missing updatedAt field")
	if email.Label == nil {
		t.Errorf("Email Label field is nil")
	} else {
		require.Equal(t, model.EmailLabelPersonal, *email.Label, "Email Label field is not expected value")
	}

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"), "Incorrect number of Email nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")
}

func TestMutationResolver_EmailDelete(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId, "original@email.com", true, "")
	neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "original@email.com", true, "")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("email/delete_email"),
		client.Var("emailId", emailId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var emailStruct struct {
		EmailDelete model.Result
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	require.Equal(t, true, emailStruct.EmailDelete.Result)

	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "User", "User_" + tenantName})
}

func TestMutationResolver_EmailRemoveFromUser(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create user and email
	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId, "original@email.com", true, "")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("email/remove_email_from_user"),
		client.Var("userId", userId),
		client.Var("email", "original@email.com"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var emailStruct struct {
		EmailRemoveFromUser model.Result
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	require.Equal(t, true, emailStruct.EmailRemoveFromUser.Result)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Email", "Email_" + tenantName, "User", "User_" + tenantName})
}

func TestMutationResolver_EmailRemoveFromUserById(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create user and email
	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId, "original@email.com", true, "")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("email/remove_email_from_user_by_id"),
		client.Var("userId", userId),
		client.Var("emailId", emailId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var emailStruct struct {
		EmailRemoveFromUserById model.Result
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	require.Equal(t, true, emailStruct.EmailRemoveFromUserById.Result)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Email", "Email_" + tenantName, "User", "User_" + tenantName})
}

func TestMutationResolver_EmailMergeToOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create organization
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "Edgeless Systems")

	organizationServiceCallbacks := events_platform.MockOrganizationServiceCallbacks{
		RefreshLastTouchpoint: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: organizationId,
			}, nil
		},
	}
	events_platform.SetOrganizationCallbacks(&organizationServiceCallbacks)

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("email/merge_email_to_organization"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var emailStruct struct {
		EmailMergeToOrganization model.Email
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	email := emailStruct.EmailMergeToOrganization

	// Check that the fields of the email struct have the expected values
	require.NotNil(t, email.ID, "Email ID is nil")
	require.NotNil(t, email.CreatedAt, "Missing createdAt field")
	require.NotNil(t, email.UpdatedAt, "Missing updatedAt field")
	require.Equal(t, true, email.Primary, "Email Primary field is not true")
	require.Equal(t, "test@gmail.com", *email.Email)
	require.Equal(t, "test@gmail.com", *email.RawEmail)
	require.Nil(t, email.EmailValidationDetails.Validated)
	if email.Label == nil {
		t.Errorf("Email Label field is nil")
	} else {
		require.Equal(t, model.EmailLabelWork, *email.Label, "Email Label field is not expected value")
	}
	require.Equal(t, model.DataSourceOpenline, email.Source, "Email Source field is not expected value")
	require.Equal(t, model.DataSourceOpenline, email.SourceOfTruth, "Email Source of truth field is not expected value")
	require.Equal(t, "test", email.AppSource, "Email App source field is not expected value")

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"), "Incorrect number of Organization nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"), "Incorrect number of Email nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName), "Incorrect number of Email_%s nodes in Neo4j", tenantName)
	require.Equal(t, 3, neo4jt.GetTotalCountOfNodes(ctx, driver), "Incorrect total number of nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Organization", "Organization_" + tenantName,
		"Email", "Email_" + tenantName})
}

func TestMutationResolver_EmailUpdateInOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create organization and email
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "Edgeless Systems")
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.ORGANIZATION, tenantName, organizationId, "original@email.com", true, "")

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("email/update_email_for_organization"),
		client.Var("organizationId", organizationId),
		client.Var("emailId", emailId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var emailStruct struct {
		EmailUpdateInOrganization model.Email
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	email := emailStruct.EmailUpdateInOrganization

	// Check that the fields of the email struct have the expected values
	require.Equal(t, emailId, email.ID, "Email ID is missing")
	require.Equal(t, true, email.Primary, "Email Primary field is not true")
	require.Equal(t, "original@email.com", *email.Email, "Email address expected not to be changed")
	require.Equal(t, "original@email.com", *email.RawEmail, "Email address expected not to be changed")
	require.NotNil(t, email.UpdatedAt, "Missing updatedAt field")
	if email.Label == nil {
		t.Errorf("Email Label field is nil")
	} else {
		require.Equal(t, model.EmailLabelWork, *email.Label, "Email Label field is not expected value")
	}

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"), "Incorrect number of Email nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")
}

func TestMutationResolver_EmailRemoveFromOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create organization and email
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "Edgeless Systems")
	neo4jt.AddEmailTo(ctx, driver, entity.ORGANIZATION, tenantName, organizationId, "original@email.com", true, "")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("email/remove_email_from_organization"),
		client.Var("organizationId", organizationId),
		client.Var("email", "original@email.com"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var emailStruct struct {
		EmailRemoveFromOrganization model.Result
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	require.Equal(t, true, emailStruct.EmailRemoveFromOrganization.Result)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Email", "Email_" + tenantName, "Organization", "Organization_" + tenantName})
}

func TestMutationResolver_EmailRemoveFromOrganizationById(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create organization and email
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "Edgeless Systems")
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.ORGANIZATION, tenantName, organizationId, "original@email.com", true, "")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("email/remove_email_from_organization_by_id"),
		client.Var("organizationId", organizationId),
		client.Var("emailId", emailId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var emailStruct struct {
		EmailRemoveFromOrganizationById model.Result
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &emailStruct)
	require.Nil(t, err, "Error unmarshalling response data")

	require.Equal(t, true, emailStruct.EmailRemoveFromOrganizationById.Result)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Email", "Email_" + tenantName, "Organization", "Organization_" + tenantName})
}

func TestQueryResolver_GetEmail_WithParentOwners(t *testing.T) {
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

	emailId := neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId1, "test@openline.com", false, "WORK")
	neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId2, "test@openline.com", false, "WORK")
	neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId1, "test@openline.com", false, "WORK")
	neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId2, "test@openline.com", false, "WORK")
	neo4jt.AddEmailTo(ctx, driver, entity.ORGANIZATION, tenantName, organizationId1, "test@openline.com", false, "WORK")
	neo4jt.AddEmailTo(ctx, driver, entity.ORGANIZATION, tenantName, organizationId2, "test@openline.com", false, "WORK")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))

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
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	emailId := neo4jt.CreateEmail(ctx, driver, tenantName, entity.EmailEntity{
		Email:       "test@openline.ai",
		RawEmail:    "testRaw@openline.ai",
		IsReachable: utils.StringPtr("reachable"),
		CreatedAt:   utils.Now(),
		UpdatedAt:   utils.Now(),
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Email": 1, "Email_" + tenantName: 1})
	assertNeo4jRelationCount(ctx, t, driver, map[string]int{"EMAIL_ADDRESS_BELONGS_TO_TENANT": 1})

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
	require.Equal(t, "reachable", *email.EmailValidationDetails.IsReachable)

	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Email", "Email_" + tenantName})
}
