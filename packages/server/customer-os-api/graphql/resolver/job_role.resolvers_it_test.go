package resolver

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	neo4jtest "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	neo4jt "github.com/customeros/customeros/packages/server/customer-os-api/test/neo4j"
	"github.com/customeros/customeros/packages/server/customer-os-api/utils/decode"
)

func TestMutationResolver_JobRoleDelete(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "LLC LLC")
	roleId := neo4jt.ContactWorksForOrganization(ctx, driver, contactId, organizationId, "CTO", false)

	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "JobRole"))

	rawResponse, err := c.RawPost(getQuery("job_role/delete_job_role"),
		client.Var("contactId", contactId),
		client.Var("roleId", roleId))
	assertRawResponseSuccess(t, rawResponse, err)

	var resultStruct struct {
		JobRole_Delete model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &resultStruct)
	require.Nil(t, err)
	require.NotNil(t, resultStruct)
	require.Equal(t, true, resultStruct.JobRole_Delete.Result)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 0, neo4jtest.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 0, neo4jtest.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
	require.Equal(t, 0, neo4jtest.GetCountOfNodes(ctx, driver, "JobRole"))

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "Organization", "Organization_" + tenantName})
}
