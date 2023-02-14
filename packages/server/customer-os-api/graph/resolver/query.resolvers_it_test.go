package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
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

func TestQueryResolver_GetData_Multiple_Records(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	contact1Id := neo4jt.CreateContactWith(driver, tenantName, "c", "1")
	neo4jt.CreateContactWith(driver, tenantName, "c", "2")

	organization1Id := neo4jt.CreateOrganization(driver, tenantName, "org 1")
	neo4jt.CreateOrganization(driver, tenantName, "org 2")

	neo4jt.LinkContactWithOrganization(driver, contact1Id, organization1Id)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "CONTACT_BELONGS_TO_TENANT"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "ORGANIZATION_BELONGS_TO_TENANT"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "CONTACT_OF"))

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
	require.Equal(t, int64(3), pagedData.TotalElements)

	require.Equal(t, "org 1", pagedData.Content[0].Organization.Name)
	require.Equal(t, "c", *pagedData.Content[0].Contact.FirstName)
	require.Equal(t, "1", *pagedData.Content[0].Contact.LastName)

	require.Nil(t, pagedData.Content[1].Contact)
	require.Equal(t, "org 2", pagedData.Content[1].Organization.Name)

	require.Nil(t, pagedData.Content[2].Organization)
	require.Equal(t, "c", *pagedData.Content[2].Contact.FirstName)
	require.Equal(t, "2", *pagedData.Content[2].Contact.LastName)
}

func TestQueryResolver_GetData_Search_Organization(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	neo4jt.CreateContactWith(driver, tenantName, "c", "1")
	neo4jt.CreateContactWith(driver, tenantName, "c", "2")

	neo4jt.CreateOrganization(driver, tenantName, "org 1")
	neo4jt.CreateOrganization(driver, tenantName, "org 2")

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "CONTACT_BELONGS_TO_TENANT"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "ORGANIZATION_BELONGS_TO_TENANT"))

	rawResponse, err := c.RawPost(getQuery("/dashboard_view/dashboard_view_with_filters"),
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", "org 1"),
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

func TestQueryResolver_GetData_Search_Contact(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	neo4jt.CreateContactWith(driver, tenantName, "c", "1")
	neo4jt.CreateContactWith(driver, tenantName, "c", "2")
	neo4jt.CreateContactWith(driver, tenantName, "c", "3")

	neo4jt.CreateOrganization(driver, tenantName, "org 1")
	neo4jt.CreateOrganization(driver, tenantName, "org 2")

	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(driver, "CONTACT_BELONGS_TO_TENANT"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "ORGANIZATION_BELONGS_TO_TENANT"))

	rawResponse, err := c.RawPost(getQuery("/dashboard_view/dashboard_view_with_filters"),
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", "3"),
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
	require.Equal(t, "3", *pagedData.Content[0].Contact.LastName)
}

func TestQueryResolver_GetData_Search_Contact_And_Organization(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	neo4jt.CreateContactWith(driver, tenantName, "c", "1")
	neo4jt.CreateContactWith(driver, tenantName, "c", "2")

	neo4jt.CreateOrganization(driver, tenantName, "org 1")
	neo4jt.CreateOrganization(driver, tenantName, "org 2")

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "CONTACT_BELONGS_TO_TENANT"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "ORGANIZATION_BELONGS_TO_TENANT"))

	rawResponse, err := c.RawPost(getQuery("/dashboard_view/dashboard_view_with_filters"),
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", "2"),
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
	require.Equal(t, int64(2), pagedData.TotalElements)

	require.Nil(t, pagedData.Content[0].Contact)
	require.Equal(t, "org 2", pagedData.Content[0].Organization.Name)

	require.Nil(t, pagedData.Content[1].Organization)
	require.Equal(t, "c", *pagedData.Content[1].Contact.FirstName)
	require.Equal(t, "2", *pagedData.Content[1].Contact.LastName)
}

func TestQueryResolver_GetData_Search_By_Email(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	//c email associated with org - wrong email
	contactId1 := neo4jt.CreateContactWith(driver, tenantName, "c", "1")
	neo4jt.AddEmailTo(driver, entity.CONTACT, tenantName, contactId1, "wrong email", true, "WORK")
	organizationId1 := neo4jt.CreateOrganization(driver, tenantName, "org 1")
	neo4jt.LinkContactWithOrganization(driver, contactId1, organizationId1)

	//c email associated with org - good email
	contactId2 := neo4jt.CreateContactWith(driver, tenantName, "c", "2")
	neo4jt.AddEmailTo(driver, entity.CONTACT, tenantName, contactId2, "good email", true, "WORK")
	organizationId2 := neo4jt.CreateOrganization(driver, tenantName, "org 2")
	neo4jt.LinkContactWithOrganization(driver, contactId2, organizationId2)

	//c associated with org email - wrong email
	contactId3 := neo4jt.CreateContactWith(driver, tenantName, "c", "3")
	organizationId3 := neo4jt.CreateOrganization(driver, tenantName, "org 3")
	neo4jt.AddEmailTo(driver, entity.ORGANIZATION, tenantName, organizationId3, "wrong email", true, "WORK")
	neo4jt.LinkContactWithOrganization(driver, contactId3, organizationId3)

	//c associated with org email - good email
	contactId4 := neo4jt.CreateContactWith(driver, tenantName, "c", "4")
	organizationId4 := neo4jt.CreateOrganization(driver, tenantName, "org 4")
	neo4jt.AddEmailTo(driver, entity.ORGANIZATION, tenantName, organizationId4, "good email", true, "WORK")
	neo4jt.LinkContactWithOrganization(driver, contactId4, organizationId4)

	//c not associated - wrong email
	contactId5 := neo4jt.CreateContactWith(driver, tenantName, "c", "5")
	neo4jt.AddEmailTo(driver, entity.CONTACT, tenantName, contactId5, "wrong email", true, "WORK")

	//c not associated - good email
	contactId6 := neo4jt.CreateContactWith(driver, tenantName, "c", "6")
	neo4jt.AddEmailTo(driver, entity.CONTACT, tenantName, contactId6, "good email", true, "WORK")

	//org 3 not associated - wrong email
	organizationId7 := neo4jt.CreateOrganization(driver, tenantName, "org 7")
	neo4jt.AddEmailTo(driver, entity.ORGANIZATION, tenantName, organizationId7, "wrong email", true, "WORK")

	//org not associated - good email
	organizationId8 := neo4jt.CreateOrganization(driver, tenantName, "org 8")
	neo4jt.AddEmailTo(driver, entity.ORGANIZATION, tenantName, organizationId8, "good email", true, "WORK")

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Email"))
	require.Equal(t, 6, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 6, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 6, neo4jt.GetCountOfRelationships(driver, "CONTACT_BELONGS_TO_TENANT"))
	require.Equal(t, 6, neo4jt.GetCountOfRelationships(driver, "ORGANIZATION_BELONGS_TO_TENANT"))
	require.Equal(t, 8, neo4jt.GetCountOfRelationships(driver, "HAS"))

	rawResponse, err := c.RawPost(getQuery("/dashboard_view/dashboard_view_with_filters"),
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", "good"),
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
	require.Equal(t, int64(4), pagedData.TotalElements)

	//c associated with org email
	require.Equal(t, "org 4", pagedData.Content[0].Organization.Name)
	require.Equal(t, "c", *pagedData.Content[0].Contact.FirstName)
	require.Equal(t, "4", *pagedData.Content[0].Contact.LastName)

	//c email associated with org
	require.Equal(t, "org 2", pagedData.Content[1].Organization.Name)
	require.Equal(t, "c", *pagedData.Content[1].Contact.FirstName)
	require.Equal(t, "2", *pagedData.Content[1].Contact.LastName)

	//org not associated
	require.Nil(t, pagedData.Content[2].Contact)
	require.Equal(t, "org 8", pagedData.Content[2].Organization.Name)

	//c not associated
	require.Nil(t, pagedData.Content[3].Organization)
	require.Equal(t, "c", *pagedData.Content[3].Contact.FirstName)
	require.Equal(t, "6", *pagedData.Content[3].Contact.LastName)
}

func TestQueryResolver_GetData_Search_By_Place(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	contactId1 := neo4jt.CreateContactWith(driver, tenantName, "c", "1")
	contactId2 := neo4jt.CreateContactWith(driver, tenantName, "c", "2")
	contactId3 := neo4jt.CreateContactWith(driver, tenantName, "c", "3")

	organizationId1 := neo4jt.CreateOrganization(driver, tenantName, "org 1")
	organizationId2 := neo4jt.CreateOrganization(driver, tenantName, "org 2")
	organizationId3 := neo4jt.CreateOrganization(driver, tenantName, "org 3")

	locationId1 := neo4jt.CreateLocation(driver, tenantName, entity.LocationEntity{
		Name:      "LOCATION 1",
		Source:    entity.DataSourceOpenline,
		AppSource: "test",
		Country:   "testCountry1",
		Region:    "testState1",
		Locality:  "testCity1",
		Address:   "testAddress1",
		Address2:  "testAddress21",
		Zip:       "testZip1",
	})
	locationId2 := neo4jt.CreateLocation(driver, tenantName, entity.LocationEntity{
		Name:      "LOCATION 2",
		Source:    entity.DataSourceOpenline,
		AppSource: "test",
		Country:   "testCountry2",
		Region:    "testState2",
		Locality:  "testCity2",
		Address:   "testAddress2",
		Address2:  "testAddress22",
		Zip:       "testZip2",
	})
	locationId3 := neo4jt.CreateLocation(driver, tenantName, entity.LocationEntity{
		Name:      "LOCATION 3",
		Source:    entity.DataSourceOpenline,
		AppSource: "test",
	})

	neo4jt.ContactAssociatedWithLocation(driver, contactId1, locationId1)
	neo4jt.ContactAssociatedWithLocation(driver, contactId2, locationId2)
	neo4jt.ContactAssociatedWithLocation(driver, contactId3, locationId3)

	neo4jt.OrganizationAssociatedWithLocation(driver, organizationId1, locationId1)
	neo4jt.OrganizationAssociatedWithLocation(driver, organizationId2, locationId2)
	neo4jt.OrganizationAssociatedWithLocation(driver, organizationId3, locationId3)

	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "Location"))
	require.Equal(t, 6, neo4jt.GetCountOfRelationships(driver, "ASSOCIATED_WITH"))

	//region search by country
	rawResponseCountry, err := c.RawPost(getQuery("/dashboard_view/dashboard_view_with_filters"),
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", "testCountry2"),
	)
	assertRawResponseSuccess(t, rawResponseCountry, err)

	var responseCountry struct {
		DashboardView model.DashboardViewItemPage
	}

	err = decode.Decode(rawResponseCountry.Data.(map[string]any), &responseCountry)
	require.Nil(t, err)
	require.NotNil(t, responseCountry)
	pagedDataCountry := responseCountry.DashboardView
	require.Equal(t, 1, pagedDataCountry.TotalPages)
	require.Equal(t, int64(2), pagedDataCountry.TotalElements)

	require.Equal(t, "org 2", pagedDataCountry.Content[0].Organization.Name)
	require.Nil(t, pagedDataCountry.Content[0].Contact)

	require.Nil(t, pagedDataCountry.Content[1].Organization)
	require.Equal(t, "c", *pagedDataCountry.Content[1].Contact.FirstName)
	require.Equal(t, "2", *pagedDataCountry.Content[1].Contact.LastName)
	//endregion

	//region search by state
	rawResponseState, err := c.RawPost(getQuery("/dashboard_view/dashboard_view_with_filters"),
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", "testState2"),
	)
	assertRawResponseSuccess(t, rawResponseState, err)

	var responseState struct {
		DashboardView model.DashboardViewItemPage
	}

	err = decode.Decode(rawResponseState.Data.(map[string]any), &responseState)
	require.Nil(t, err)
	require.NotNil(t, responseState)
	pagedDataState := responseState.DashboardView
	require.Equal(t, 1, pagedDataState.TotalPages)
	require.Equal(t, int64(2), pagedDataState.TotalElements)

	require.Equal(t, "org 2", pagedDataState.Content[0].Organization.Name)
	require.Nil(t, pagedDataState.Content[0].Contact)

	require.Nil(t, pagedDataState.Content[1].Organization)
	require.Equal(t, "c", *pagedDataState.Content[1].Contact.FirstName)
	require.Equal(t, "2", *pagedDataState.Content[1].Contact.LastName)
	//endregion

	//region search by city
	rawResponseCity, err := c.RawPost(getQuery("/dashboard_view/dashboard_view_with_filters"),
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", "testCity2"),
	)
	assertRawResponseSuccess(t, rawResponseCity, err)

	var responseCity struct {
		DashboardView model.DashboardViewItemPage
	}

	err = decode.Decode(rawResponseCity.Data.(map[string]any), &responseCity)
	require.Nil(t, err)
	require.NotNil(t, responseCity)
	pagedDataCity := responseCity.DashboardView
	require.Equal(t, 1, pagedDataCity.TotalPages)
	require.Equal(t, int64(2), pagedDataCity.TotalElements)

	require.Equal(t, "org 2", pagedDataCity.Content[0].Organization.Name)
	require.Nil(t, pagedDataCity.Content[0].Contact)

	require.Nil(t, pagedDataCity.Content[1].Organization)
	require.Equal(t, "c", *pagedDataCity.Content[1].Contact.FirstName)
	require.Equal(t, "2", *pagedDataCity.Content[1].Contact.LastName)
	//endregion
}
