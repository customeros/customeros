package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_PhoneNumberMergeToContact(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	// Create a default contact
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("merge_phone_number_to_contact"),
		client.Var("contactId", contactId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumber struct {
		PhoneNumberMergeToContact model.PhoneNumber
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumber)
	require.Nil(t, err, "Error unmarshalling response data")

	// Check that the fields of the phoneNumber struct have the expected values
	require.NotNil(t, phoneNumber.PhoneNumberMergeToContact.ID, "PhoneNumber ID is nil")
	require.Equal(t, true, phoneNumber.PhoneNumberMergeToContact.Primary, "PhoneNumber Primary field is not true")
	require.Equal(t, "+1234567890", phoneNumber.PhoneNumberMergeToContact.E164, "PhoneNumber E164 field is not expected value")
	if phoneNumber.PhoneNumberMergeToContact.Label == nil {
		t.Errorf("PhoneNumber Label field is nil")
	} else {
		require.Equal(t, model.PhoneNumberLabelWork, *phoneNumber.PhoneNumberMergeToContact.Label, "PhoneNumber Label field is not expected value")
	}
	require.Equal(t, model.DataSourceOpenline, phoneNumber.PhoneNumberMergeToContact.Source, "PhoneNumber Source field is not expected value")

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"), "Incorrect number of Contact nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "PhoneNumber"), "Incorrect number of PhoneNumber nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "PhoneNumber_"+tenantName), "Incorrect number of PhoneNumber_%s nodes in Neo4j", tenantName)
	require.Equal(t, 3, neo4jt.GetTotalCountOfNodes(driver), "Incorrect total number of nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "PHONE_ASSOCIATED_WITH"), "Incorrect number of PHONE_ASSOCIATED_WITH relationships in Neo4j")

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "PhoneNumber", "PhoneNumber_" + tenantName})
}
