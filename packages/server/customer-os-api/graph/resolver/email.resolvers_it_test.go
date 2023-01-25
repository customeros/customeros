package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_EmailMergeToContact(t *testing.T) {
	defer tearDownTestCase()(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(driver, tenantName)

	// Create a default contact
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("merge_email_to_contact"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var email struct {
		EmailMergeToContact model.Email
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &email)
	require.Nil(t, err, "Error unmarshalling response data")

	e := email.EmailMergeToContact

	// Check that the fields of the email struct have the expected values
	require.NotNil(t, e.ID, "Email ID is nil")
	require.NotNil(t, e.CreatedAt, "Missing createdAt field")
	require.NotNil(t, e.UpdatedAt, "Missing updatedAt field")
	require.Equal(t, true, e.Primary, "Email Primary field is not true")
	require.Equal(t, "test@gmail.com", e.Email, "Email Email field is not expected value")
	if e.Label == nil {
		t.Errorf("Email Label field is nil")
	} else {
		require.Equal(t, model.EmailLabelWork, *e.Label, "Email Label field is not expected value")
	}
	require.Equal(t, model.DataSourceOpenline, e.Source, "Email Source field is not expected value")
	require.Equal(t, model.DataSourceOpenline, e.SourceOfTruth, "Email Source of truth field is not expected value")
	require.Equal(t, "test", e.AppSource, "Email App source field is not expected value")

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"), "Incorrect number of Contact nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Email"), "Incorrect number of Email nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Email_"+tenantName), "Incorrect number of Email_%s nodes in Neo4j", tenantName)
	require.Equal(t, 3, neo4jt.GetTotalCountOfNodes(driver), "Incorrect total number of nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "EMAILED_AT"), "Incorrect number of EMAILED_AT relationships in Neo4j")

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Email", "Email_" + tenantName})
}

func TestMutationResolver_EmailUpdateInContact(t *testing.T) {
	defer tearDownTestCase()(t)

	// Create a tenant in the Neo4j database
	neo4jt.CreateTenant(driver, tenantName)

	// Create a default contact and email
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	emailId := neo4jt.AddEmailToContact(driver, contactId, "original@email.com", true, "")

	// Make the RawPost request and check for errors
	rawResponse, err := c.RawPost(getQuery("update_email_for_contact"),
		client.Var("contactId", contactId),
		client.Var("emailId", emailId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the email struct
	var email struct {
		EmailUpdateInContact model.Email
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &email)
	require.Nil(t, err, "Error unmarshalling response data")

	e := email.EmailUpdateInContact

	// Check that the fields of the email struct have the expected values
	require.Equal(t, emailId, e.ID, "Email ID is nil")
	require.Equal(t, true, e.Primary, "Email Primary field is not true")
	require.Equal(t, "new@email.com", e.Email, "Email Email field is not expected value")
	require.NotNil(t, e.UpdatedAt, "Missing updatedAt field")
	if e.Label == nil {
		t.Errorf("Email Label field is nil")
	} else {
		require.Equal(t, model.EmailLabelHome, *e.Label, "Email Label field is not expected value")
	}

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Email"), "Incorrect number of Email nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "EMAILED_AT"), "Incorrect number of EMAILED_AT relationships in Neo4j")
}
