package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMutationResolver_JobRoleCreate_WithOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "LLC LLC")

	rawResponse, err := c.RawPost(getQuery("job_role/create_job_role_for_organization"),
		client.Var("contactId", contactId),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var jobRoleStruct struct {
		JobRole_Create model.JobRole
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &jobRoleStruct)
	require.Nil(t, err)

	createdRole := jobRoleStruct.JobRole_Create

	require.NotNil(t, createdRole.ID)
	require.NotNil(t, createdRole.CreatedAt)
	require.NotNil(t, createdRole.UpdatedAt)
	expectedStartedAt := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	expectedEndedAt := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	require.Equal(t, expectedStartedAt, *createdRole.StartedAt)
	require.Equal(t, expectedEndedAt, *createdRole.EndedAt)
	require.Equal(t, organizationId, createdRole.Organization.ID)
	require.Equal(t, "CEO", *createdRole.JobTitle)
	require.Equal(t, "CEO of the company", *createdRole.Description)
	require.Equal(t, "Acme Inc.", *createdRole.Company)
	require.Equal(t, true, createdRole.Primary)
	require.Equal(t, model.DataSourceOpenline, createdRole.Source)
	require.Equal(t, "Hubspot", createdRole.AppSource)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "JobRole"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "JobRole_"+tenantName))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "WORKS_AS"))

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName,
		"Organization", "Organization_" + tenantName, "JobRole", "JobRole_" + tenantName})
}

func TestMutationResolver_JobRoleCreate_WithoutOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("job_role/create_job_role_without_organization"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var jobRoleStruct struct {
		JobRole_Create model.JobRole
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &jobRoleStruct)
	require.Nil(t, err)

	createdRole := jobRoleStruct.JobRole_Create

	require.NotNil(t, createdRole.ID)
	require.NotNil(t, createdRole.CreatedAt)
	require.NotNil(t, createdRole.UpdatedAt)
	test.AssertRecentTime(t, utils.Now())
	require.Nil(t, createdRole.EndedAt)
	require.Equal(t, "CEO", *createdRole.JobTitle)
	require.Equal(t, true, createdRole.Primary)
	require.Equal(t, model.DataSourceOpenline, createdRole.Source)
	require.Equal(t, "customer-os-api", createdRole.AppSource)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 0, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "JobRole"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "JobRole_"+tenantName))
	require.Equal(t, 0, neo4jtest.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "WORKS_AS"))

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "JobRole", "JobRole_" + tenantName})
}

func TestMutationResolver_JobRoleUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "LLC LLC")
	roleId := neo4jt.ContactWorksForOrganization(ctx, driver, contactId, organizationId, "CTO", false)

	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "JobRole"))

	rawResponse, err := c.RawPost(getQuery("job_role/update_job_role"),
		client.Var("contactId", contactId),
		client.Var("roleId", roleId))
	assertRawResponseSuccess(t, rawResponse, err)

	var jobRoleStruct struct {
		JobRole_Update model.JobRole
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &jobRoleStruct)
	require.Nil(t, err)

	updatedRole := jobRoleStruct.JobRole_Update

	require.NotNil(t, updatedRole)
	require.NotNil(t, updatedRole.UpdatedAt)
	require.Nil(t, updatedRole.EndedAt)
	require.Nil(t, updatedRole.StartedAt)
	require.Equal(t, true, updatedRole.Primary)
	require.Equal(t, "CEO", *updatedRole.JobTitle)
	require.Equal(t, "CEO of Apple", *updatedRole.Description)
	require.Equal(t, "Apple", *updatedRole.Company)
	require.Equal(t, organizationId, updatedRole.Organization.ID)

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName,
		"Organization", "Organization_" + tenantName, "JobRole"})
}

func TestMutationResolver_JobRoleUpdate_ChangeOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "LLC LLC")
	newOrganizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "NEW CO")
	roleId := neo4jt.ContactWorksForOrganization(ctx, driver, contactId, organizationId, "CTO", false)

	rawResponse, err := c.RawPost(getQuery("job_role/update_job_role_change_organization"),
		client.Var("contactId", contactId),
		client.Var("roleId", roleId),
		client.Var("organizationId", newOrganizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var jobRoleStruct struct {
		JobRole_Update model.JobRole
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &jobRoleStruct)
	require.Nil(t, err)

	updatedRole := jobRoleStruct.JobRole_Update

	require.NotNil(t, updatedRole)
	require.NotNil(t, updatedRole.UpdatedAt)
	require.Equal(t, true, updatedRole.Primary)
	require.Equal(t, "CEO", *updatedRole.JobTitle)
	require.Equal(t, newOrganizationId, updatedRole.Organization.ID)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "JobRole"))

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName,
		"Organization", "Organization_" + tenantName, "JobRole"})
}

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
