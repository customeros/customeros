package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryResolver_Search_Organization_By_Name(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 2")
	neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(ctx, driver, "ORGANIZATION_BELONGS_TO_TENANT"))

	require.Equal(t, int64(1), assert_Search_Organization_By_Name(t, "org 1", false).TotalElements)
	require.Equal(t, int64(2), assert_Search_Organization_By_Name(t, "org 1", true).TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_Name(t, "org 2", false).TotalElements)
	require.Equal(t, int64(2), assert_Search_Organization_By_Name(t, "org 2", true).TotalElements)
	require.Equal(t, int64(2), assert_Search_Organization_By_Name(t, "org", false).TotalElements)
	require.Equal(t, int64(3), assert_Search_Organization_By_Name(t, "org", true).TotalElements)
}

func assert_Search_Organization_By_Name(t *testing.T, searchTerm string, includeEmpty bool) model.OrganizationPage {
	rawResponse, err := c.RawPost(getQuery("/dashboard_view/organization/dashboard_view_organization_filter_by_name"),
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", searchTerm),
		client.Var("includeEmpty", includeEmpty),
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

func TestQueryResolver_Search_Organization_By_Website(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Website: "org 1",
	})
	neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Website: "org 2",
	})
	neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(ctx, driver, "ORGANIZATION_BELONGS_TO_TENANT"))

	require.Equal(t, int64(1), assert_Search_Organization_By_Website(t, "org 1", false).TotalElements)
	require.Equal(t, int64(2), assert_Search_Organization_By_Website(t, "org 1", true).TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_Website(t, "org 2", false).TotalElements)
	require.Equal(t, int64(2), assert_Search_Organization_By_Website(t, "org 2", true).TotalElements)
	require.Equal(t, int64(2), assert_Search_Organization_By_Website(t, "org", false).TotalElements)
	require.Equal(t, int64(3), assert_Search_Organization_By_Website(t, "org", true).TotalElements)
}

func assert_Search_Organization_By_Website(t *testing.T, searchTerm string, includeEmpty bool) model.OrganizationPage {
	rawResponse, err := c.RawPost(getQuery("/dashboard_view/organization/dashboard_view_organization_filter_by_website"),
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", searchTerm),
		client.Var("includeEmpty", includeEmpty),
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

func TestQueryResolver_Search_Organization_By_ORGANIZATION_Filter(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateTenantOrganization(ctx, driver, tenantName, "org excluded")
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 2")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 3")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 4")
	neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		ReferenceId: "100/200",
	})
	neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		CustomerOsId: "C-123-ABC",
	})

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

	require.Equal(t, 7, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 7, neo4jt.GetCountOfRelationships(ctx, driver, "ORGANIZATION_BELONGS_TO_TENANT"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	require.Equal(t, int64(1), assert_Search_Organization_By_ORGANIZATION(t, "org 1").TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_ORGANIZATION(t, "org 2").TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_ORGANIZATION(t, "org 3").TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_ORGANIZATION(t, "org 4").TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_ORGANIZATION(t, "100").TotalElements)
	require.Equal(t, int64(1), assert_Search_Organization_By_ORGANIZATION(t, "ABC").TotalElements)
	require.Equal(t, int64(4), assert_Search_Organization_By_ORGANIZATION(t, "org").TotalElements)
	require.Equal(t, int64(0), assert_Search_Organization_By_ORGANIZATION(t, "org excluded").TotalElements)
}

func assert_Search_Organization_By_ORGANIZATION(t *testing.T, searchTerm string) model.OrganizationPage {
	rawResponse, err := c.RawPost(getQuery("/dashboard_view/organization/dashboard_view_organization_filter_by_organization"),
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

func TestQueryResolver_Search_Organizations_By_Owner_In_IncludeEmptyFalse(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	userId1 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	userId2 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)

	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1 for owner 1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 2 for owner 1")
	organizationId3 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1 for owner 2")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org without owner")

	neo4jt.UserOwnsOrganization(ctx, driver, userId1, organizationId1)
	neo4jt.UserOwnsOrganization(ctx, driver, userId1, organizationId2)
	neo4jt.UserOwnsOrganization(ctx, driver, userId2, organizationId3)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(ctx, driver, "OWNS"))

	rawResponse := callGraphQL(t, "dashboard_view/organization/dashboard_view_organization_filter_by_owner", map[string]interface{}{"ownerIdList": []string{userId1}, "ownerIdEmpty": false, "page": 1, "limit": 10})

	var organizationsPageStruct struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationsPageStruct)
	require.Nil(t, err)

	require.Equal(t, int64(4), organizationsPageStruct.DashboardView_Organizations.TotalAvailable)
	require.Equal(t, int64(2), organizationsPageStruct.DashboardView_Organizations.TotalElements)
	require.Equal(t, 2, len(organizationsPageStruct.DashboardView_Organizations.Content))
	require.ElementsMatch(t, []string{organizationId1, organizationId2},
		[]string{organizationsPageStruct.DashboardView_Organizations.Content[0].ID, organizationsPageStruct.DashboardView_Organizations.Content[1].ID})
}

func TestQueryResolver_Search_Organizations_By_Owner_In_IncludeEmptyTrue(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	userId1 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	userId2 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)

	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1 for owner 1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 2 for owner 1")
	organizationId3 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1 for owner 2")
	organizationId4 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org without owner")

	neo4jt.UserOwnsOrganization(ctx, driver, userId1, organizationId1)
	neo4jt.UserOwnsOrganization(ctx, driver, userId1, organizationId2)
	neo4jt.UserOwnsOrganization(ctx, driver, userId2, organizationId3)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(ctx, driver, "OWNS"))

	rawResponse := callGraphQL(t, "dashboard_view/organization/dashboard_view_organization_filter_by_owner", map[string]interface{}{"ownerIdList": []string{userId1}, "ownerIdEmpty": true, "page": 1, "limit": 10})

	var organizationsPageStruct struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationsPageStruct)
	require.Nil(t, err)

	require.Equal(t, int64(4), organizationsPageStruct.DashboardView_Organizations.TotalAvailable)
	require.Equal(t, int64(3), organizationsPageStruct.DashboardView_Organizations.TotalElements)
	require.Equal(t, 3, len(organizationsPageStruct.DashboardView_Organizations.Content))
	require.ElementsMatch(t, []string{organizationId1, organizationId2, organizationId4},
		[]string{organizationsPageStruct.DashboardView_Organizations.Content[0].ID, organizationsPageStruct.DashboardView_Organizations.Content[1].ID, organizationsPageStruct.DashboardView_Organizations.Content[2].ID})
}

func TestQueryResolver_Sort_Organizations_ByLastTouchpointAt(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := time.Now().UTC()
	secAgo60 := now.Add(-60 * time.Second)

	organizationId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name:             "org1",
		LastTouchpointAt: &secAgo60,
	})
	organizationId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name:             "org2",
		LastTouchpointAt: &now,
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})

	rawResponse := callGraphQL(t, "dashboard_view/organization/dashboard_view_organization_sort",
		map[string]interface{}{
			"page":    1,
			"limit":   10,
			"sortBy":  "LAST_TOUCHPOINT_AT",
			"sortDir": "ASC",
		})

	var organizationsPageStruct struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationsPageStruct)
	require.Nil(t, err)

	require.Equal(t, int64(2), organizationsPageStruct.DashboardView_Organizations.TotalAvailable)
	require.Equal(t, int64(2), organizationsPageStruct.DashboardView_Organizations.TotalElements)

	require.Equal(t, organizationId1, organizationsPageStruct.DashboardView_Organizations.Content[0].ID)
	require.Equal(t, organizationId2, organizationsPageStruct.DashboardView_Organizations.Content[1].ID)
}

func TestQueryResolver_Sort_Organizations_ByLastTouchpointType(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	created := "CREATED"
	updated := "UPDATED"

	organizationId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name:               "org1",
		LastTouchpointType: &updated,
	})
	organizationId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name:               "org2",
		LastTouchpointType: &created,
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})

	rawResponse := callGraphQL(t, "dashboard_view/organization/dashboard_view_organization_sort",
		map[string]interface{}{
			"page":    1,
			"limit":   10,
			"sortBy":  "LAST_TOUCHPOINT_TYPE",
			"sortDir": "ASC",
		})

	var organizationsPageStruct struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationsPageStruct)
	require.Nil(t, err)

	require.Equal(t, int64(2), organizationsPageStruct.DashboardView_Organizations.TotalAvailable)
	require.Equal(t, int64(2), organizationsPageStruct.DashboardView_Organizations.TotalElements)

	require.Equal(t, organizationId2, organizationsPageStruct.DashboardView_Organizations.Content[0].ID)
	require.Equal(t, organizationId1, organizationsPageStruct.DashboardView_Organizations.Content[1].ID)
}

func TestQueryResolver_Sort_Organizations_ByForecastAmount(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org1",
		RenewalForecast: entity.RenewalForecast{
			Amount: utils.ToPtr[float64](200),
		},
	})
	organizationId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org2",
	})
	organizationId3 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org3",
		RenewalForecast: entity.RenewalForecast{
			Amount: utils.ToPtr[float64](100.5),
		},
	})
	organizationId4 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org4",
		RenewalForecast: entity.RenewalForecast{
			Amount: utils.ToPtr[float64](300),
		},
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 4})

	rawResponse := callGraphQL(t, "dashboard_view/organization/dashboard_view_organization_sort",
		map[string]interface{}{
			"page":    1,
			"limit":   10,
			"sortBy":  "FORECAST_AMOUNT",
			"sortDir": "ASC",
		})

	var organizationsPageStruct struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationsPageStruct)
	require.Nil(t, err)

	require.Equal(t, int64(4), organizationsPageStruct.DashboardView_Organizations.TotalAvailable)
	require.Equal(t, int64(4), organizationsPageStruct.DashboardView_Organizations.TotalElements)

	require.Equal(t, organizationId3, organizationsPageStruct.DashboardView_Organizations.Content[0].ID)
	require.Equal(t, organizationId1, organizationsPageStruct.DashboardView_Organizations.Content[1].ID)
	require.Equal(t, organizationId4, organizationsPageStruct.DashboardView_Organizations.Content[2].ID)
	require.Equal(t, organizationId2, organizationsPageStruct.DashboardView_Organizations.Content[3].ID)
}

func TestQueryResolver_Sort_Organizations_ByRenewalLikelihood(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org1",
		RenewalLikelihood: entity.RenewalLikelihood{
			RenewalLikelihood: string(entity.RenewalLikelihoodProbabilityMedium),
		},
	})
	organizationId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org2",
	})
	organizationId3 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org3",
		RenewalLikelihood: entity.RenewalLikelihood{
			RenewalLikelihood: string(entity.RenewalLikelihoodProbabilityHigh),
		},
	})
	organizationId4 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org4",
		RenewalLikelihood: entity.RenewalLikelihood{
			RenewalLikelihood: string(entity.RenewalLikelihoodProbabilityLow),
		},
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 4})

	rawResponse := callGraphQL(t, "dashboard_view/organization/dashboard_view_organization_sort",
		map[string]interface{}{
			"page":    1,
			"limit":   10,
			"sortBy":  "RENEWAL_LIKELIHOOD",
			"sortDir": "ASC",
		})

	var organizationsPageStruct struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationsPageStruct)
	require.Nil(t, err)

	require.Equal(t, int64(4), organizationsPageStruct.DashboardView_Organizations.TotalAvailable)
	require.Equal(t, int64(4), organizationsPageStruct.DashboardView_Organizations.TotalElements)

	require.Equal(t, organizationId3, organizationsPageStruct.DashboardView_Organizations.Content[0].ID)
	require.Equal(t, organizationId1, organizationsPageStruct.DashboardView_Organizations.Content[1].ID)
	require.Equal(t, organizationId4, organizationsPageStruct.DashboardView_Organizations.Content[2].ID)
	require.Equal(t, organizationId2, organizationsPageStruct.DashboardView_Organizations.Content[3].ID)
}

func TestQueryResolver_Sort_Organizations_ByRenewalCycleNext(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	daysFromNow10 := time.Now().AddDate(0, 0, 10)
	daysFromNow20 := time.Now().AddDate(0, 0, 20)

	organizationId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org1",
		BillingDetails: entity.BillingDetails{
			RenewalCycleNext: utils.TimePtr(daysFromNow10),
		},
	})
	organizationId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org2",
	})
	organizationId3 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org3",
		BillingDetails: entity.BillingDetails{
			RenewalCycleNext: utils.TimePtr(daysFromNow20),
		},
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 3})

	rawResponse := callGraphQL(t, "dashboard_view/organization/dashboard_view_organization_sort",
		map[string]interface{}{
			"page":    1,
			"limit":   10,
			"sortBy":  "RENEWAL_CYCLE_NEXT",
			"sortDir": "ASC",
		})

	var organizationsPageStruct struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationsPageStruct)
	require.Nil(t, err)

	require.Equal(t, int64(3), organizationsPageStruct.DashboardView_Organizations.TotalAvailable)
	require.Equal(t, int64(3), organizationsPageStruct.DashboardView_Organizations.TotalElements)

	require.Equal(t, organizationId1, organizationsPageStruct.DashboardView_Organizations.Content[0].ID)
	require.Equal(t, organizationId3, organizationsPageStruct.DashboardView_Organizations.Content[1].ID)
	require.Equal(t, organizationId2, organizationsPageStruct.DashboardView_Organizations.Content[2].ID)
}

func TestQueryResolver_Sort_Organizations_ByOrganizationName_WithOrganizationHierarchy(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	independentOrgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "C-Org",
	})
	parent1OrgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "B-Parent",
	})
	sub1_2OrgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "X-Sub",
	})
	sub1_1OrgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "A-Sub",
	})
	parent2OrgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "O-Parent",
	})
	sub2_1OrgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "K-Sub",
	})
	sub2_2OrgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "Y-Sub",
	})
	neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, parent1OrgId, sub1_1OrgId, "")
	neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, parent1OrgId, sub1_2OrgId, "")
	neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, parent2OrgId, sub2_1OrgId, "")
	neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, parent2OrgId, sub2_2OrgId, "")

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 7})

	rawResponse := callGraphQL(t, "dashboard_view/organization/dashboard_view_organization_sort",
		map[string]interface{}{
			"page":    1,
			"limit":   10,
			"sortBy":  "ORGANIZATION",
			"sortDir": "ASC",
		})

	var organizationsPageStruct struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationsPageStruct)
	require.Nil(t, err)

	require.Equal(t, int64(7), organizationsPageStruct.DashboardView_Organizations.TotalElements)

	require.Equal(t, "B-Parent", organizationsPageStruct.DashboardView_Organizations.Content[0].Name)
	require.Equal(t, "A-Sub", organizationsPageStruct.DashboardView_Organizations.Content[1].Name)
	require.Equal(t, "X-Sub", organizationsPageStruct.DashboardView_Organizations.Content[2].Name)
	require.Equal(t, "C-Org", organizationsPageStruct.DashboardView_Organizations.Content[3].Name)
	require.Equal(t, "O-Parent", organizationsPageStruct.DashboardView_Organizations.Content[4].Name)
	require.Equal(t, "K-Sub", organizationsPageStruct.DashboardView_Organizations.Content[5].Name)
	require.Equal(t, "Y-Sub", organizationsPageStruct.DashboardView_Organizations.Content[6].Name)

	require.Equal(t, parent1OrgId, organizationsPageStruct.DashboardView_Organizations.Content[0].ID)
	require.Equal(t, sub1_1OrgId, organizationsPageStruct.DashboardView_Organizations.Content[1].ID)
	require.Equal(t, sub1_2OrgId, organizationsPageStruct.DashboardView_Organizations.Content[2].ID)
	require.Equal(t, independentOrgId, organizationsPageStruct.DashboardView_Organizations.Content[3].ID)
	require.Equal(t, parent2OrgId, organizationsPageStruct.DashboardView_Organizations.Content[4].ID)
	require.Equal(t, sub2_1OrgId, organizationsPageStruct.DashboardView_Organizations.Content[5].ID)
	require.Equal(t, sub2_2OrgId, organizationsPageStruct.DashboardView_Organizations.Content[6].ID)
}
