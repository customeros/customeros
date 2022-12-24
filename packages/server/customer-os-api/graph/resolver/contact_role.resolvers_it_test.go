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
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "WORKS"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "HAS_ROLE"))
}

//func TestMutationResolver_ContactUpdateCompanyPosition_SameCompanyNewPosition(t *testing.T) {
//	defer tearDownTestCase()(t)
//	neo4jt.CreateTenant(driver, tenantName)
//	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
//	companyId := neo4jt.CreateCompany(driver, tenantName, "LLC LLC")
//	positionId := neo4jt.ContactWorksForCompany(driver, contactId, companyId, "CTO", true)
//
//	rawResponse, err := c.RawPost(getQuery("update_company_position_same_company"),
//		client.Var("contactId", contactId),
//		client.Var("companyId", companyId),
//		client.Var("companyPositionId", positionId))
//	assertRawResponseSuccess(t, rawResponse, err)
//
//	var companyPosition struct {
//		Contact_UpdateCompanyPosition model.ContactRole
//	}
//
//	err = decode.Decode(rawResponse.Data.(map[string]any), &companyPosition)
//	require.Nil(t, err)
//
//	require.NotNil(t, companyPosition.Contact_UpdateCompanyPosition.ID)
//	require.Equal(t, companyId, companyPosition.Contact_UpdateCompanyPosition.Company.ID)
//	require.Equal(t, "LLC LLC", companyPosition.Contact_UpdateCompanyPosition.Company.Name)
//	require.Equal(t, "CEO", *companyPosition.Contact_UpdateCompanyPosition.JobTitle)
//
//	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
//	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Company"))
//	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "WORKS_AT"))
//}
//
//func TestMutationResolver_ContactUpdateCompanyPosition_InOtherExistingCompany(t *testing.T) {
//	defer tearDownTestCase()(t)
//	neo4jt.CreateTenant(driver, tenantName)
//	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
//	companyId := neo4jt.CreateCompany(driver, tenantName, "Current Company")
//	otherCompanyId := neo4jt.CreateCompany(driver, tenantName, "Other Company")
//	positionId := neo4jt.ContactWorksForCompany(driver, contactId, companyId, "CTO", false)
//
//	rawResponse, err := c.RawPost(getQuery("update_company_position_other_company"),
//		client.Var("contactId", contactId),
//		client.Var("companyId", otherCompanyId),
//		client.Var("companyPositionId", positionId))
//	assertRawResponseSuccess(t, rawResponse, err)
//
//	var companyPosition struct {
//		Contact_UpdateCompanyPosition model.ContactRole
//	}
//
//	err = decode.Decode(rawResponse.Data.(map[string]any), &companyPosition)
//	require.Nil(t, err)
//
//	require.NotNil(t, companyPosition.Contact_UpdateCompanyPosition.ID)
//	require.Equal(t, otherCompanyId, companyPosition.Contact_UpdateCompanyPosition.Company.ID)
//	require.Equal(t, "Other Company", companyPosition.Contact_UpdateCompanyPosition.Company.Name)
//	require.Equal(t, "CEO", *companyPosition.Contact_UpdateCompanyPosition.JobTitle)
//
//	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
//	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Company"))
//	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "WORKS_AT"))
//}
//
//func TestMutationResolver_ContactUpdateCompanyPosition_InNewCompany(t *testing.T) {
//	defer tearDownTestCase()(t)
//	neo4jt.CreateTenant(driver, tenantName)
//	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
//	companyId := neo4jt.CreateCompany(driver, tenantName, "LLC LLC")
//	positionId := neo4jt.ContactWorksForCompany(driver, contactId, companyId, "CTO", false)
//
//	rawResponse, err := c.RawPost(getQuery("update_company_position_new_company"),
//		client.Var("contactId", contactId),
//		client.Var("companyPositionId", positionId))
//	assertRawResponseSuccess(t, rawResponse, err)
//
//	var companyPosition struct {
//		Contact_UpdateCompanyPosition model.ContactRole
//	}
//
//	err = decode.Decode(rawResponse.Data.(map[string]any), &companyPosition)
//	require.Nil(t, err)
//
//	require.NotNil(t, companyPosition.Contact_UpdateCompanyPosition.ID)
//	require.NotEqual(t, companyId, companyPosition.Contact_UpdateCompanyPosition.Company.ID)
//	require.Equal(t, "new company", companyPosition.Contact_UpdateCompanyPosition.Company.Name)
//	require.Equal(t, "CEO", *companyPosition.Contact_UpdateCompanyPosition.JobTitle)
//
//	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
//	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Company"))
//	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "WORKS_AT"))
//}
