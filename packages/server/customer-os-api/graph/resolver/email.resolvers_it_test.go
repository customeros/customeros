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
	neo4jt.CreateTenant(driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("merge_email_to_contact"),
		client.Var("contactId", contactId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var email struct {
		EmailMergeToContact model.Email
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &email)
	require.Nil(t, err)
	require.NotNil(t, email)

	require.NotNil(t, email.EmailMergeToContact.ID)
	require.Equal(t, true, email.EmailMergeToContact.Primary)
	require.Equal(t, "test@gmail.com", email.EmailMergeToContact.Email)
	require.Equal(t, model.EmailLabelWork, *email.EmailMergeToContact.Label)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Email_"+tenantName))
	require.Equal(t, 3, neo4jt.GetTotalCountOfNodes(driver))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "EMAILED_AT"))

	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Email", "Email_" + tenantName})
}
