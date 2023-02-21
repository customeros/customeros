package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestMutationResolver_PhoneNumberMergeToContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	// Create a default contact
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("phone_number/merge_phone_number_to_contact"),
		client.Var("contactId", contactId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumber struct {
		PhoneNumberMergeToContact model.PhoneNumber
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumber)
	require.Nil(t, err, "Error unmarshalling response data")

	createdPhoneNumber := phoneNumber.PhoneNumberMergeToContact
	// Check that the fields of the phoneNumber struct have the expected values
	require.NotNil(t, createdPhoneNumber.CreatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPhoneNumber.CreatedAt)
	require.NotNil(t, createdPhoneNumber.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPhoneNumber.UpdatedAt)

	require.NotNil(t, createdPhoneNumber.ID, "PhoneNumber ID is nil")
	require.Equal(t, true, createdPhoneNumber.Primary, "PhoneNumber Primary field is not true")
	require.Equal(t, false, *createdPhoneNumber.Validated, "PhoneNumber Validated field is not false")
	require.Nil(t, createdPhoneNumber.E164)
	require.Equal(t, "+1234567890", *createdPhoneNumber.RawPhoneNumber, "PhoneNumber E164 field is not expected value")
	if createdPhoneNumber.Label == nil {
		t.Errorf("PhoneNumber Label field is nil")
	} else {
		require.Equal(t, model.PhoneNumberLabelWork, *createdPhoneNumber.Label, "PhoneNumber Label field is not expected value")
	}
	require.Equal(t, model.DataSourceOpenline, createdPhoneNumber.Source, "PhoneNumber Source field is not expected value")

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"), "Incorrect number of Contact nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"), "Incorrect number of PhoneNumber nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber_"+tenantName), "Incorrect number of PhoneNumber_%s nodes in Neo4j", tenantName)
	require.Equal(t, 3, neo4jt.GetTotalCountOfNodes(ctx, driver), "Incorrect total number of nodes in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"), "Incorrect number of PHONE_ASSOCIATED_WITH relationships in Neo4j")

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "PhoneNumber", "PhoneNumber_" + tenantName})
}
