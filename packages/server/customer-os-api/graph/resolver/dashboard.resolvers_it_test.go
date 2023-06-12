package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_Search_Contact_By_Name(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId1 := neo4jt.CreateContactWith(ctx, driver, tenantName, "a", "b")
	neo4jt.CreateContactWith(ctx, driver, tenantName, "c", "d")
	contactId3 := neo4jt.CreateContactWith(ctx, driver, tenantName, "e", "f")
	neo4jt.CreateContactWith(ctx, driver, tenantName, "a", "d")

	locationId1 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 1",
		Source: entity.DataSourceOpenline,
		Region: "NY",
	})
	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 2",
		Source: entity.DataSourceOpenline,
		Region: "TX",
	})

	neo4jt.ContactAssociatedWithLocation(ctx, driver, contactId1, locationId1)
	neo4jt.ContactAssociatedWithLocation(ctx, driver, contactId3, locationId2)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(ctx, driver, "CONTACT_BELONGS_TO_TENANT"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	require.Equal(t, int64(2), assert_Search_Contact_By_Name(t, "a").TotalElements)
	require.Equal(t, int64(1), assert_Search_Contact_By_Name(t, "b").TotalElements)
	require.Equal(t, int64(1), assert_Search_Contact_By_Name(t, "c").TotalElements)
	require.Equal(t, int64(2), assert_Search_Contact_By_Name(t, "d").TotalElements)
	require.Equal(t, int64(1), assert_Search_Contact_By_Name(t, "e").TotalElements)
	require.Equal(t, int64(0), assert_Search_Contact_By_Name(t, "TEST").TotalElements)
}

func assert_Search_Contact_By_Name(t *testing.T, searchTerm string) model.ContactsPage {
	rawResponse, err := c.RawPost(getQuery("/dashboard_view/contact/dashboard_view_contact_filter_by_name"),
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", searchTerm),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var responseRaw struct {
		DashboardView_Contacts model.ContactsPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &responseRaw)
	require.Nil(t, err)
	require.NotNil(t, responseRaw)

	return responseRaw.DashboardView_Contacts
}

func TestQueryResolver_Search_Contact_By_Regions(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId1 := neo4jt.CreateContactWith(ctx, driver, tenantName, "a", "b")
	contactId2 := neo4jt.CreateContactWith(ctx, driver, tenantName, "c", "d")
	neo4jt.CreateContactWith(ctx, driver, tenantName, "a", "f")
	neo4jt.CreateContactWith(ctx, driver, tenantName, "g", "h")

	locationId1 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 1",
		Source: entity.DataSourceOpenline,
		Region: "NY",
	})
	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 2",
		Source: entity.DataSourceOpenline,
		Region: "TX",
	})

	neo4jt.ContactAssociatedWithLocation(ctx, driver, contactId1, locationId1)
	neo4jt.ContactAssociatedWithLocation(ctx, driver, contactId2, locationId2)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(ctx, driver, "CONTACT_BELONGS_TO_TENANT"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	regionTX := "TX"

	require.Equal(t, int64(1), assert_Search_Contact_By_Regions(t, "NY", nil).TotalElements)
	require.Equal(t, int64(1), assert_Search_Contact_By_Regions(t, "TX", nil).TotalElements)
	require.Equal(t, int64(0), assert_Search_Contact_By_Regions(t, "TEST", nil).TotalElements)
	require.Equal(t, int64(2), assert_Search_Contact_By_Regions(t, "NY", &regionTX).TotalElements)
}

func assert_Search_Contact_By_Regions(t *testing.T, region1 string, region2 *string) model.ContactsPage {
	query := "/dashboard_view/contact/dashboard_view_contact_filter_by_region"
	options := []client.Option{client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("region1", region1),
	}

	if region2 != nil {
		query = "/dashboard_view/contact/dashboard_view_contact_filter_by_regions"
		options = append(options, client.Var("region2", *region2))
	}

	rawResponse, err := c.RawPost(getQuery(query), options...)
	assertRawResponseSuccess(t, rawResponse, err)

	var responseRaw struct {
		DashboardView_Contacts model.ContactsPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &responseRaw)
	require.Nil(t, err)
	require.NotNil(t, responseRaw)

	return responseRaw.DashboardView_Contacts
}

func TestQueryResolver_Search_Contact_By_Name_And_Regions(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateTenantOrganization(ctx, driver, tenantName, "org excluded")

	contactId1 := neo4jt.CreateContactWith(ctx, driver, tenantName, "a", "b")
	contactId2 := neo4jt.CreateContactWith(ctx, driver, tenantName, "c", "d")
	contactId3 := neo4jt.CreateContactWith(ctx, driver, tenantName, "a", "e")
	neo4jt.CreateContactWith(ctx, driver, tenantName, "f", "g")

	locationId1 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 1",
		Source: entity.DataSourceOpenline,
		Region: "NY",
	})
	neo4jt.ContactAssociatedWithLocation(ctx, driver, contactId1, locationId1)
	neo4jt.ContactAssociatedWithLocation(ctx, driver, contactId2, locationId1)

	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 2",
		Source: entity.DataSourceOpenline,
		Region: "TX",
	})
	neo4jt.ContactAssociatedWithLocation(ctx, driver, contactId3, locationId2)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(ctx, driver, "CONTACT_BELONGS_TO_TENANT"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	regionTX := "TX"

	require.Equal(t, int64(0), assert_Search_Contact_By_Name_And_Regions(t, "NY", nil, "TEST").TotalElements)
	require.Equal(t, int64(1), assert_Search_Contact_By_Name_And_Regions(t, "NY", nil, "a").TotalElements)
	require.Equal(t, int64(1), assert_Search_Contact_By_Name_And_Regions(t, "NY", nil, "b").TotalElements)
	require.Equal(t, int64(1), assert_Search_Contact_By_Name_And_Regions(t, "NY", nil, "c").TotalElements)
	require.Equal(t, int64(0), assert_Search_Contact_By_Name_And_Regions(t, "NY", nil, "e").TotalElements)

	require.Equal(t, int64(0), assert_Search_Contact_By_Name_And_Regions(t, "TEST", nil, "a").TotalElements)
	require.Equal(t, int64(0), assert_Search_Contact_By_Name_And_Regions(t, "TEST", nil, "b").TotalElements)
	require.Equal(t, int64(0), assert_Search_Contact_By_Name_And_Regions(t, "TEST", nil, "c").TotalElements)
	require.Equal(t, int64(0), assert_Search_Contact_By_Name_And_Regions(t, "TEST", nil, "d").TotalElements)
	require.Equal(t, int64(0), assert_Search_Contact_By_Name_And_Regions(t, "TEST", nil, "e").TotalElements)
	require.Equal(t, int64(0), assert_Search_Contact_By_Name_And_Regions(t, "TEST", nil, "f").TotalElements)

	require.Equal(t, int64(2), assert_Search_Contact_By_Name_And_Regions(t, "NY", &regionTX, "a").TotalElements)
	require.Equal(t, int64(1), assert_Search_Contact_By_Name_And_Regions(t, "NY", &regionTX, "b").TotalElements)
	require.Equal(t, int64(1), assert_Search_Contact_By_Name_And_Regions(t, "NY", &regionTX, "c").TotalElements)
	require.Equal(t, int64(1), assert_Search_Contact_By_Name_And_Regions(t, "NY", &regionTX, "d").TotalElements)
	require.Equal(t, int64(1), assert_Search_Contact_By_Name_And_Regions(t, "NY", &regionTX, "e").TotalElements)
	require.Equal(t, int64(0), assert_Search_Contact_By_Name_And_Regions(t, "NY", &regionTX, "f").TotalElements)
}

func assert_Search_Contact_By_Name_And_Regions(t *testing.T, region1 string, region2 *string, searchTerm string) model.ContactsPage {
	query := "/dashboard_view/contact/dashboard_view_contact_filter_by_name_and_region"
	options := []client.Option{client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", searchTerm),
		client.Var("region1", region1),
	}

	if region2 != nil {
		query = "/dashboard_view/contact/dashboard_view_contact_filter_by_name_and_regions"
		options = append(options, client.Var("region2", *region2))
	}

	rawResponse, err := c.RawPost(getQuery(query), options...)
	assertRawResponseSuccess(t, rawResponse, err)

	var responseRaw struct {
		DashboardView_Contacts model.ContactsPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &responseRaw)
	require.Nil(t, err)
	require.NotNil(t, responseRaw)

	return responseRaw.DashboardView_Contacts
}

func TestQueryResolver_Search_Organization_By_Name(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateTenantOrganization(ctx, driver, tenantName, "org excluded")
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 2")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 3")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 4")

	locationId1 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 1",
		Source: entity.DataSourceOpenline,
		Region: "NY",
	})
	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 2",
		Source: entity.DataSourceOpenline,
		Region: "TX",
	})

	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId1, locationId1)
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId2, locationId2)

	require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 5, neo4jt.GetCountOfRelationships(ctx, driver, "ORGANIZATION_BELONGS_TO_TENANT"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	require.Equal(t, int64(1), assert_Search_Organization_By_Name(t, "org 1").TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_Name(t, "org 2").TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_Name(t, "org 3").TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_Name(t, "org 4").TotalElements)
	require.Equal(t, int64(4), assert_Search_Organization_By_Name(t, "org").TotalElements)
	require.Equal(t, int64(0), assert_Search_Organization_By_Name(t, "org excluded").TotalElements)
}

func assert_Search_Organization_By_Name(t *testing.T, searchTerm string) model.OrganizationPage {
	rawResponse, err := c.RawPost(getQuery("/dashboard_view/organization/dashboard_view_organization_filter_by_name"),
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", searchTerm),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var responseRaw struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &responseRaw)
	require.Nil(t, err)
	require.NotNil(t, responseRaw)

	return responseRaw.DashboardView_Organizations
}

func TestQueryResolver_Search_Organization_By_Regions(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateTenantOrganization(ctx, driver, tenantName, "org excluded")
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 2")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 3")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 4")

	locationId1 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 1",
		Source: entity.DataSourceOpenline,
		Region: "NY",
	})
	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 2",
		Source: entity.DataSourceOpenline,
		Region: "TX",
	})

	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId1, locationId1)
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId2, locationId2)

	require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 5, neo4jt.GetCountOfRelationships(ctx, driver, "ORGANIZATION_BELONGS_TO_TENANT"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	testRegion := "TEST"
	region2 := "TX"

	require.Equal(t, int64(0), assert_Search_Organization_By_Regions(t, "TEST", nil).TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_Regions(t, "NY", nil).TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_Regions(t, "TX", nil).TotalElements)
	require.Equal(t, int64(2), assert_Search_Organization_By_Regions(t, "NY", &region2).TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_Regions(t, "NY", &testRegion).TotalElements)
	require.Equal(t, int64(0), assert_Search_Organization_By_Regions(t, "org", nil).TotalElements)
}

func assert_Search_Organization_By_Regions(t *testing.T, region1 string, region2 *string) model.OrganizationPage {
	query := "/dashboard_view/organization/dashboard_view_organization_filter_by_region"
	options := []client.Option{client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("region1", region1),
	}

	if region2 != nil {
		query = "/dashboard_view/organization/dashboard_view_organization_filter_by_regions"
		options = append(options, client.Var("region2", *region2))
	}

	rawResponse, err := c.RawPost(getQuery(query), options...)
	assertRawResponseSuccess(t, rawResponse, err)

	var responseRaw struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &responseRaw)
	require.Nil(t, err)
	require.NotNil(t, responseRaw)

	return responseRaw.DashboardView_Organizations
}

func TestQueryResolver_Search_Organization_By_Name_And_Regions(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateTenantOrganization(ctx, driver, tenantName, "org excluded")

	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 2")
	organizationId3 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 3")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 4")

	locationId1 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 1",
		Source: entity.DataSourceOpenline,
		Region: "NY",
	})
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId1, locationId1)
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId2, locationId1)

	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 2",
		Source: entity.DataSourceOpenline,
		Region: "TX",
	})
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId3, locationId2)

	require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 5, neo4jt.GetCountOfRelationships(ctx, driver, "ORGANIZATION_BELONGS_TO_TENANT"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	regionTX := "TX"

	require.Equal(t, int64(0), assert_Search_Organization_By_Name_And_Regions(t, "NY", nil, "TEST").TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_Name_And_Regions(t, "NY", nil, "org 1").TotalElements)
	require.Equal(t, int64(2), assert_Search_Organization_By_Name_And_Regions(t, "NY", nil, "org").TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_Name_And_Regions(t, "TX", nil, "org").TotalElements)

	require.Equal(t, int64(0), assert_Search_Organization_By_Name_And_Regions(t, "TEST", nil, "TEST").TotalElements)
	require.Equal(t, int64(0), assert_Search_Organization_By_Name_And_Regions(t, "TEST", nil, "org 1").TotalElements)
	require.Equal(t, int64(0), assert_Search_Organization_By_Name_And_Regions(t, "TEST", nil, "org").TotalElements)

	require.Equal(t, int64(0), assert_Search_Organization_By_Name_And_Regions(t, "NY", &regionTX, "TEST").TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_Name_And_Regions(t, "NY", &regionTX, "org 1").TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_Name_And_Regions(t, "NY", &regionTX, "org 3").TotalElements)
	require.Equal(t, int64(3), assert_Search_Organization_By_Name_And_Regions(t, "NY", &regionTX, "org").TotalElements)
}

func assert_Search_Organization_By_Name_And_Regions(t *testing.T, region1 string, region2 *string, searchTerm string) model.OrganizationPage {
	query := "/dashboard_view/organization/dashboard_view_organization_filter_by_name_and_region"
	options := []client.Option{client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", searchTerm),
		client.Var("region1", region1),
	}

	if region2 != nil {
		query = "/dashboard_view/organization/dashboard_view_organization_filter_by_name_and_regions"
		options = append(options, client.Var("region2", *region2))
	}

	rawResponse, err := c.RawPost(getQuery(query), options...)
	assertRawResponseSuccess(t, rawResponse, err)

	var responseRaw struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &responseRaw)
	require.Nil(t, err)
	require.NotNil(t, responseRaw)

	return responseRaw.DashboardView_Organizations
}

func TestQueryResolver_DashboardViewPortfolioOrganizations(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	userId1 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	userId2 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)

	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org for portfolio 1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "second org for portfolio 1")
	organizationId3 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org for portfolio 2")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org without owner")
	neo4jt.UserOwnsOrganization(ctx, driver, userId1, organizationId1)
	neo4jt.UserOwnsOrganization(ctx, driver, userId1, organizationId2)
	neo4jt.UserOwnsOrganization(ctx, driver, userId2, organizationId3)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(ctx, driver, "OWNS"))

	rawResponse := callGraphQL(t, "dashboard_view/organization/dashboard_view_organization_filter_by_owner", map[string]interface{}{"ownerId": userId1, "page": 1, "limit": 10})

	var organizationsPageStruct struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationsPageStruct)
	require.Nil(t, err)

	require.Equal(t, int64(2), organizationsPageStruct.DashboardView_Organizations.TotalElements)
	require.Equal(t, 2, len(organizationsPageStruct.DashboardView_Organizations.Content))
	require.Equal(t, organizationId1, organizationsPageStruct.DashboardView_Organizations.Content[0].ID)
	require.Equal(t, organizationId2, organizationsPageStruct.DashboardView_Organizations.Content[1].ID)
}

func TestQueryResolver_DashboardViewRelationshipOrganizations(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateOrganizationRelationship(ctx, driver, entity.Customer.String())
	neo4jt.CreateOrganizationRelationship(ctx, driver, entity.Investor.String())

	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "customer org")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "second customer org")
	organizationId3 := neo4jt.CreateOrganization(ctx, driver, tenantName, "investor org")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org without relationship")
	neo4jt.CreateOrganizationRelationshipStages(ctx, driver, tenantName, entity.Customer.String(), []string{"A", "B"})
	neo4jt.LinkOrganizationWithRelationshipAndStage(ctx, driver, organizationId1, entity.Customer.String(), "A")
	neo4jt.LinkOrganizationWithRelationshipAndStage(ctx, driver, organizationId2, entity.Customer.String(), "B")
	neo4jt.LinkOrganizationWithRelationship(ctx, driver, organizationId3, entity.Investor.String())

	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationRelationship"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationRelationshipStage"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(ctx, driver, "IS"))

	rawResponse := callGraphQL(t, "dashboard_view/organization/dashboard_view_by_relationship",
		map[string]interface{}{"relationship": model.OrganizationRelationshipCustomer.String(), "page": 1, "limit": 10})

	var organizationsPageStruct struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationsPageStruct)
	require.Nil(t, err)

	require.Equal(t, int64(2), organizationsPageStruct.DashboardView_Organizations.TotalElements)
	require.Equal(t, 2, len(organizationsPageStruct.DashboardView_Organizations.Content))
	require.Equal(t, organizationId1, organizationsPageStruct.DashboardView_Organizations.Content[0].ID)
	require.Equal(t, organizationId2, organizationsPageStruct.DashboardView_Organizations.Content[1].ID)
}
