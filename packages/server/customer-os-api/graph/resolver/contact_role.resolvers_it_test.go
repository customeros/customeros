package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_ContactRoleDelete(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	companyId := neo4jt.CreateCompany(driver, tenantName, "LLC LLC")
	roleId := neo4jt.ContactWorksForCompany(driver, contactId, companyId, "CTO", false)

	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "HAS_ROLE"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "WORKS"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Role"))

	rawResponse, err := c.RawPost(getQuery("delete_contact_role"),
		client.Var("contactId", contactId),
		client.Var("roleId", roleId))
	assertRawResponseSuccess(t, rawResponse, err)

	var result struct {
		ContactRole_Delete model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &result)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, true, result.ContactRole_Delete.Result)

	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "HAS_ROLE"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "WORKS"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Role"))
}
