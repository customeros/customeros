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

func TestMutationResolver_ContactRoleCreate_WithCompany(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	companyId := neo4jt.CreateCompany(driver, tenantName, "LLC LLC")

	rawResponse, err := c.RawPost(getQuery("create_contact_role_for_company"),
		client.Var("contactId", contactId),
		client.Var("companyId", companyId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactRole struct {
		ContactRole_Create model.ContactRole
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactRole)
	require.Nil(t, err)

	createdRole := contactRole.ContactRole_Create

	require.NotNil(t, createdRole.ID)
	require.Equal(t, companyId, createdRole.Company.ID)
	require.Equal(t, "CEO", *createdRole.JobTitle)
	require.Equal(t, true, createdRole.Primary)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Company"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Role"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Role_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "WORKS"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "HAS_ROLE"))
}

func TestMutationResolver_ContactRoleCreate_WithoutCompany(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("create_contact_role_without_company"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactRole struct {
		ContactRole_Create model.ContactRole
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactRole)
	require.Nil(t, err)

	createdRole := contactRole.ContactRole_Create

	require.NotNil(t, createdRole.ID)
	require.Equal(t, "CEO", *createdRole.JobTitle)
	require.Equal(t, true, createdRole.Primary)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Company"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Role"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Role_"+tenantName))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "WORKS"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "HAS_ROLE"))
}

func TestMutationResolver_ContactRoleUpdate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	companyId := neo4jt.CreateCompany(driver, tenantName, "LLC LLC")
	roleId := neo4jt.ContactWorksForCompany(driver, contactId, companyId, "CTO", false)

	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "HAS_ROLE"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "WORKS"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Role"))

	rawResponse, err := c.RawPost(getQuery("update_contact_role"),
		client.Var("contactId", contactId),
		client.Var("roleId", roleId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactRole struct {
		ContactRole_Update model.ContactRole
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactRole)
	require.Nil(t, err)

	updatedRole := contactRole.ContactRole_Update

	require.NotNil(t, updatedRole)
	require.Equal(t, true, updatedRole.Primary)
	require.Equal(t, "CEO", *updatedRole.JobTitle)
	require.Equal(t, companyId, updatedRole.Company.ID)
}

func TestMutationResolver_ContactRoleUpdate_ChangeCompany(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	companyId := neo4jt.CreateCompany(driver, tenantName, "LLC LLC")
	newCompanyId := neo4jt.CreateCompany(driver, tenantName, "NEW CO")
	roleId := neo4jt.ContactWorksForCompany(driver, contactId, companyId, "CTO", false)

	rawResponse, err := c.RawPost(getQuery("update_contact_role_change_company"),
		client.Var("contactId", contactId),
		client.Var("roleId", roleId),
		client.Var("companyId", newCompanyId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactRole struct {
		ContactRole_Update model.ContactRole
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactRole)
	require.Nil(t, err)

	updatedRole := contactRole.ContactRole_Update

	require.NotNil(t, updatedRole)
	require.Equal(t, true, updatedRole.Primary)
	require.Equal(t, "CEO", *updatedRole.JobTitle)
	require.Equal(t, newCompanyId, updatedRole.Company.ID)

	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "HAS_ROLE"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "WORKS"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Role"))
}
