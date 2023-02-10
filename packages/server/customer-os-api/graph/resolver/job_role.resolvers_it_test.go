package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_JobRoleCreate_WithOrganization(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	organizationId := neo4jt.CreateOrganization(driver, tenantName, "LLC LLC")

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
	require.Equal(t, organizationId, createdRole.Organization.ID)
	require.Equal(t, "CEO", *createdRole.JobTitle)
	require.Equal(t, true, createdRole.Primary)
	require.Equal(t, int64(2), createdRole.ResponsibilityLevel)
	require.Equal(t, model.DataSourceOpenline, createdRole.Source)
	require.Equal(t, model.DataSourceOpenline, createdRole.SourceOfTruth)
	require.Equal(t, "Hubspot", createdRole.AppSource)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "JobRole"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "JobRole_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "ROLE_IN"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "WORKS_AS"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName,
		"Organization", "Organization_" + tenantName, "JobRole", "JobRole_" + tenantName})
}

func TestMutationResolver_JobRoleCreate_WithoutOrganization(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)

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
	require.Equal(t, "CEO", *createdRole.JobTitle)
	require.Equal(t, true, createdRole.Primary)
	require.Equal(t, int64(0), createdRole.ResponsibilityLevel)
	require.Equal(t, model.DataSourceOpenline, createdRole.Source)
	require.Equal(t, model.DataSourceOpenline, createdRole.SourceOfTruth)
	require.Equal(t, "customer-os-api", createdRole.AppSource)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "JobRole"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "JobRole_"+tenantName))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "ROLE_IN"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "WORKS_AS"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "JobRole", "JobRole_" + tenantName})
}

func TestMutationResolver_JobRoleUpdate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	organizationId := neo4jt.CreateOrganization(driver, tenantName, "LLC LLC")
	roleId := neo4jt.ContactWorksForOrganization(driver, contactId, organizationId, "CTO", false)

	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "WORKS_AS"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "ROLE_IN"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "JobRole"))

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
	require.Equal(t, true, updatedRole.Primary)
	require.Equal(t, int64(1), updatedRole.ResponsibilityLevel)
	require.Equal(t, "CEO", *updatedRole.JobTitle)
	require.Equal(t, organizationId, updatedRole.Organization.ID)
	require.Equal(t, model.DataSourceOpenline, updatedRole.SourceOfTruth)

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName,
		"Organization", "Organization_" + tenantName, "JobRole"})
}

func TestMutationResolver_JobRoleUpdate_ChangeOrganization(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	organizationId := neo4jt.CreateOrganization(driver, tenantName, "LLC LLC")
	newOrganizationId := neo4jt.CreateOrganization(driver, tenantName, "NEW CO")
	roleId := neo4jt.ContactWorksForOrganization(driver, contactId, organizationId, "CTO", false)

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
	require.Equal(t, int64(0), updatedRole.ResponsibilityLevel)
	require.Equal(t, "CEO", *updatedRole.JobTitle)
	require.Equal(t, newOrganizationId, updatedRole.Organization.ID)
	require.Equal(t, model.DataSourceOpenline, updatedRole.SourceOfTruth)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "WORKS_AS"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "ROLE_IN"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "JobRole"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName,
		"Organization", "Organization_" + tenantName, "JobRole"})
}

func TestMutationResolver_JobRoleDelete(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	organizationId := neo4jt.CreateOrganization(driver, tenantName, "LLC LLC")
	roleId := neo4jt.ContactWorksForOrganization(driver, contactId, organizationId, "CTO", false)

	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "WORKS_AS"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "ROLE_IN"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "JobRole"))

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
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "WORKS_AS"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "ROLE_IN"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "JobRole"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "Organization", "Organization_" + tenantName})
}
