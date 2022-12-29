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

	contactId := neo4jt.CreateDefaultContact(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("merge_phone_number_to_contact"),
		client.Var("contactId", contactId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var phoneNumber struct {
		PhoneNumberMergeToContact model.PhoneNumber
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &phoneNumber)
	require.Nil(t, err)
	require.NotNil(t, phoneNumber)

	require.NotNil(t, phoneNumber.PhoneNumberMergeToContact.ID)
	require.Equal(t, true, phoneNumber.PhoneNumberMergeToContact.Primary)
	require.Equal(t, "+1234567890", phoneNumber.PhoneNumberMergeToContact.E164)
	require.Equal(t, model.PhoneNumberLabelWork, *phoneNumber.PhoneNumberMergeToContact.Label)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "PhoneNumber_"+tenantName))
	require.Equal(t, 3, neo4jt.GetTotalCountOfNodes(driver))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "CALLED_AT"))

	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "PhoneNumber", "PhoneNumber_" + tenantName})
}
