package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_GetData_EmptyDB(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("/dashboard_view/dashboard_view_no_filters"),
		client.Var("page", 1),
		client.Var("limit", 10),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var response struct {
		PageResponse model.DashboardViewItemPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &response)
	require.Nil(t, err)
	require.NotNil(t, response)
	pagedData := response.PageResponse
	require.Equal(t, 0, pagedData.TotalPages)
	require.Equal(t, int64(0), pagedData.TotalElements)
}

func TestQueryResolver_GetData_One_Organization(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	neo4jt.CreateOrganization(driver, tenantName, "org 1")

	rawResponse, err := c.RawPost(getQuery("/dashboard_view/dashboard_view_no_filters"),
		client.Var("page", 1),
		client.Var("limit", 10),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var response struct {
		DashboardView model.DashboardViewItemPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &response)
	require.Nil(t, err)
	require.NotNil(t, response)
	pagedData := response.DashboardView
	require.Equal(t, 1, pagedData.TotalPages)
	require.Equal(t, int64(1), pagedData.TotalElements)

	require.Nil(t, pagedData.Content[0].Contact)
	require.Equal(t, "org 1", pagedData.Content[0].Organization.Name)
}

func TestQueryResolver_GetData_One_Contact(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	neo4jt.CreateContactWith(driver, tenantName, "c", "1")

	rawResponse, err := c.RawPost(getQuery("/dashboard_view/dashboard_view_no_filters"),
		client.Var("page", 1),
		client.Var("limit", 10),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var response struct {
		DashboardView model.DashboardViewItemPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &response)
	require.Nil(t, err)
	require.NotNil(t, response)
	pagedData := response.DashboardView
	require.Equal(t, 1, pagedData.TotalPages)
	require.Equal(t, int64(1), pagedData.TotalElements)

	require.Nil(t, pagedData.Content[0].Organization)
	require.Equal(t, "c", *pagedData.Content[0].Contact.FirstName)
	require.Equal(t, "1", *pagedData.Content[0].Contact.LastName)
}

func TestQueryResolver_GetData_One_Contact_Linked_To_One_Organization(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	contactId := neo4jt.CreateContactWith(driver, tenantName, "c", "1")
	organizationId := neo4jt.CreateOrganization(driver, tenantName, "org 1")
	neo4jt.LinkContactWithOrganization(driver, contactId, organizationId)

	rawResponse, err := c.RawPost(getQuery("/dashboard_view/dashboard_view_no_filters"),
		client.Var("page", 1),
		client.Var("limit", 10),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var response struct {
		DashboardView model.DashboardViewItemPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &response)
	require.Nil(t, err)
	require.NotNil(t, response)
	pagedData := response.DashboardView
	require.Equal(t, 1, pagedData.TotalPages)
	require.Equal(t, int64(1), pagedData.TotalElements)

	require.Equal(t, "org 1", pagedData.Content[0].Organization.Name)
	require.Equal(t, "c", *pagedData.Content[0].Contact.FirstName)
	require.Equal(t, "1", *pagedData.Content[0].Contact.LastName)
}

//func TestQueryResolver_GetData_Dataset1(t *testing.T) {
//	defer tearDownTestCase()(t)
//	neo4jt.CreateTenant(driver, tenantName)
//
//	contact1Id := neo4jt.CreateContactWith(driver, tenantName, "c", "1")
//	neo4jt.CreateContactWith(driver, tenantName, "c", "2")
//
//	organization1Id := neo4jt.CreateOrganization(driver, tenantName, "org 1")
//	neo4jt.CreateOrganization(driver, tenantName, "org 2")
//
//	neo4jt.LinkContactWithOrganization(driver, contact1Id, organization1Id)
//	//neo4jt.ContactWorksForOrganization(driver, contact1Id, organization1Id, "employee", true)
//
//	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Contact"))
//	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Organization"))
//	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tenant"))
//	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "CONTACT_BELONGS_TO_TENANT"))
//	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "ORGANIZATION_BELONGS_TO_TENANT"))
//	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "CONTACT_OF"))
//
//	require.Equal(t, int64(1), neo4jt.Q1(driver, tenantName))
//	require.Equal(t, int64(1), neo4jt.Q2(driver, tenantName))
//	require.Equal(t, int64(1), neo4jt.Q3(driver, tenantName))
//	require.Equal(t, int64(1), neo4jt.Q4(driver, tenantName))
//
//	rawResponse, err := c.RawPost(getQuery("/dashboard_view/dashboard_view_no_filters"),
//		client.Var("page", 1),
//		client.Var("limit", 10),
//	)
//	assertRawResponseSuccess(t, rawResponse, err)
//
//	var response struct {
//		DashboardView model.DashboardViewItemPage
//	}
//
//	err = decode.Decode(rawResponse.Data.(map[string]any), &response)
//	require.Nil(t, err)
//	require.NotNil(t, response)
//	pagedData := response.DashboardView
//	require.Equal(t, 1, pagedData.TotalPages)
//	require.Equal(t, int64(3), pagedData.TotalElements)
//
//	require.Equal(t, "org 1", pagedData.Content[0].Organization.Name)
//	require.Equal(t, "c", *pagedData.Content[0].Contact.FirstName)
//	require.Equal(t, "1", *pagedData.Content[0].Contact.LastName)
//
//	require.Nil(t, pagedData.Content[1].Contact)
//	require.Equal(t, "org 2", pagedData.Content[1].Organization.Name)
//
//	require.Nil(t, pagedData.Content[2].Organization)
//	require.Equal(t, "c", pagedData.Content[1].Contact.FirstName)
//	require.Equal(t, "2", pagedData.Content[1].Contact.LastName)
//
//}

//func TestQueryResolver_GetData_EmptyDB(t *testing.T) {
//	defer tearDownTestCase()(t)
//	neo4jt.CreateTenant(driver, tenantName)
//	neo4jt.CreateOrganization(driver, tenantName, "org 1")
//	neo4jt.CreateOrganization(driver, tenantName, "org 2")
//
//	neo4jt.CreateContactWith(driver, tenantName, "c", "1")
//	neo4jt.CreateContactWith(driver, tenantName, "c", "2")
//
//	require.Equal(t, 5, neo4jt.GetCountOfNodes(driver, "Organization"))
//
//	rawResponse, err := c.RawPost(getQuery("get_organizations"),
//		client.Var("page", 1),
//		client.Var("limit", 3),
//	)
//	assertRawResponseSuccess(t, rawResponse, err)
//
//	var organizations struct {
//		Organizations model.OrganizationPage
//	}
//
//	err = decode.Decode(rawResponse.Data.(map[string]any), &organizations)
//	require.Nil(t, err)
//	require.NotNil(t, organizations)
//	pagedOrganizations := organizations.Organizations
//	require.Equal(t, 2, pagedOrganizations.TotalPages)
//	require.Equal(t, int64(4), pagedOrganizations.TotalElements)
//	require.Equal(t, "OPENLINE", pagedOrganizations.Content[0].Name)
//	require.Equal(t, "OpEnLiNe", pagedOrganizations.Content[1].Name)
//	require.Equal(t, "some other open organization", pagedOrganizations.Content[2].Name)
//}
//
