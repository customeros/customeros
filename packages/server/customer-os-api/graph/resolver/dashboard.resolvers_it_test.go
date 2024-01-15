package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryResolver_Search_Organization_By_Name(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 2")
	neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 3, neo4jtest.GetCountOfRelationships(ctx, driver, "ORGANIZATION_BELONGS_TO_TENANT"))

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
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Website: "org 1",
	})
	neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Website: "org 2",
	})
	neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 3, neo4jtest.GetCountOfRelationships(ctx, driver, "ORGANIZATION_BELONGS_TO_TENANT"))

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
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

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
		Source: neo4jentity.DataSourceOpenline,
		Region: "NY",
	})
	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 2",
		Source: neo4jentity.DataSourceOpenline,
		Region: "TX",
	})

	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId1, locationId1)
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId2, locationId2)

	require.Equal(t, 7, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 7, neo4jtest.GetCountOfRelationships(ctx, driver, "ORGANIZATION_BELONGS_TO_TENANT"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

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
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateTenantOrganization(ctx, driver, tenantName, "org excluded")
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 2")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 3")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 4")

	locationId1 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 1",
		Source: neo4jentity.DataSourceOpenline,
		Region: "NY",
	})
	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 2",
		Source: neo4jentity.DataSourceOpenline,
		Region: "TX",
	})

	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId1, locationId1)
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId2, locationId2)

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 5, neo4jtest.GetCountOfRelationships(ctx, driver, "ORGANIZATION_BELONGS_TO_TENANT"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

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
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateTenantOrganization(ctx, driver, tenantName, "org excluded")

	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 2")
	organizationId3 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 3")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org 4")

	locationId1 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 1",
		Source: neo4jentity.DataSourceOpenline,
		Region: "NY",
	})
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId1, locationId1)
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId2, locationId1)

	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:   "LOCATION 2",
		Source: neo4jentity.DataSourceOpenline,
		Region: "TX",
	})
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId3, locationId2)

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 5, neo4jtest.GetCountOfRelationships(ctx, driver, "ORGANIZATION_BELONGS_TO_TENANT"))
	require.Equal(t, 3, neo4jtest.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

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
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	userId1 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	userId2 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)

	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1 for owner 1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 2 for owner 1")
	organizationId3 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1 for owner 2")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org without owner")

	neo4jt.UserOwnsOrganization(ctx, driver, userId1, organizationId1)
	neo4jt.UserOwnsOrganization(ctx, driver, userId1, organizationId2)
	neo4jt.UserOwnsOrganization(ctx, driver, userId2, organizationId3)

	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 3, neo4jtest.GetCountOfRelationships(ctx, driver, "OWNS"))

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
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	userId1 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	userId2 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)

	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1 for owner 1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 2 for owner 1")
	organizationId3 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1 for owner 2")
	organizationId4 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org without owner")

	neo4jt.UserOwnsOrganization(ctx, driver, userId1, organizationId1)
	neo4jt.UserOwnsOrganization(ctx, driver, userId1, organizationId2)
	neo4jt.UserOwnsOrganization(ctx, driver, userId2, organizationId3)

	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 3, neo4jtest.GetCountOfRelationships(ctx, driver, "OWNS"))

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
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

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

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})

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
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

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

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})

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
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organizationId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org1",
		RenewalSummary: entity.RenewalSummary{
			ArrForecast: utils.ToPtr[float64](200),
		},
	})
	organizationId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org2",
	})
	organizationId3 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org3",
		RenewalSummary: entity.RenewalSummary{
			ArrForecast: utils.ToPtr[float64](100.5),
		},
	})
	organizationId4 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org4",
		RenewalSummary: entity.RenewalSummary{
			ArrForecast: utils.ToPtr[float64](300),
		},
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 4})

	rawResponse := callGraphQL(t, "dashboard_view/organization/dashboard_view_organization_sort",
		map[string]interface{}{
			"page":    1,
			"limit":   10,
			"sortBy":  "FORECAST_ARR",
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
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organizationId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org1",
		RenewalSummary: entity.RenewalSummary{
			RenewalLikelihoodOrder: utils.Int64Ptr(30),
		},
	})
	organizationId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org2",
	})
	organizationId3 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org3",
		RenewalSummary: entity.RenewalSummary{
			RenewalLikelihoodOrder: utils.Int64Ptr(40),
		},
	})
	organizationId4 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org4",
		RenewalSummary: entity.RenewalSummary{
			RenewalLikelihoodOrder: utils.Int64Ptr(20),
		},
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 4})

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

	require.Equal(t, organizationId4, organizationsPageStruct.DashboardView_Organizations.Content[0].ID)
	require.Equal(t, organizationId1, organizationsPageStruct.DashboardView_Organizations.Content[1].ID)
	require.Equal(t, organizationId3, organizationsPageStruct.DashboardView_Organizations.Content[2].ID)
	require.Equal(t, organizationId2, organizationsPageStruct.DashboardView_Organizations.Content[3].ID)
}

func TestQueryResolver_Sort_Organizations_ByRenewalDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	daysFromNow10 := utils.Now().AddDate(0, 0, 10)
	daysFromNow20 := utils.Now().AddDate(0, 0, 20)

	organizationId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org1",
		RenewalSummary: entity.RenewalSummary{
			NextRenewalAt: utils.TimePtr(daysFromNow10),
		},
	})
	organizationId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org2",
	})
	organizationId3 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org3",
		RenewalSummary: entity.RenewalSummary{
			NextRenewalAt: utils.TimePtr(daysFromNow20),
		},
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 3})

	rawResponse := callGraphQL(t, "dashboard_view/organization/dashboard_view_organization_sort",
		map[string]interface{}{
			"page":    1,
			"limit":   10,
			"sortBy":  "RENEWAL_DATE",
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
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

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

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 7})

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

func insertContractWithActiveRenewalOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, orgId string, contract entity.ContractEntity, opportunity entity.OpportunityEntity) string {
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, contract)
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, opportunity)
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)
	return contractId
}

func insertContractWithOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, orgId string, contract entity.ContractEntity, opportunity entity.OpportunityEntity) string {
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, contract)
	neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, opportunity)
	return contractId
}

func insertServiceLineItem(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, quantity int64, startedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  id,
		Billed:    billedType,
		Price:     price,
		Quantity:  quantity,
		StartedAt: startedAt,
	})
	return id
}

func insertServiceLineItemEnded(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, quantity int64, startedAt, endedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  id,
		Billed:    billedType,
		Price:     price,
		Quantity:  quantity,
		StartedAt: startedAt,
		EndedAt:   &endedAt,
	})
	return id
}

func insertServiceLineItemCanceled(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, quantity int64, startedAt, endedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:         id,
		ParentID:   id,
		Billed:     billedType,
		Price:      price,
		Quantity:   quantity,
		IsCanceled: true,
		StartedAt:  startedAt,
		EndedAt:    &endedAt,
	})
	return id
}

func insertServiceLineItemWithParent(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, quantity int64, previousBilledType entity.BilledType, previousPrice float64, previousQuantity int64, startedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:               id,
		ParentID:         parentId,
		Billed:           billedType,
		Price:            price,
		Quantity:         quantity,
		PreviousBilled:   previousBilledType,
		PreviousPrice:    previousPrice,
		PreviousQuantity: previousQuantity,
		StartedAt:        startedAt,
	})
}

func insertServiceLineItemEndedWithParent(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, quantity int64, previousBilledType entity.BilledType, previousPrice float64, previousQuantity int64, startedAt, endedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:               id,
		ParentID:         parentId,
		Billed:           billedType,
		Price:            price,
		Quantity:         quantity,
		PreviousBilled:   previousBilledType,
		PreviousPrice:    previousPrice,
		PreviousQuantity: previousQuantity,
		StartedAt:        startedAt,
		EndedAt:          &endedAt,
	})
}

func insertServiceLineItemCanceledWithParent(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, quantity int64, previousBilledType entity.BilledType, previousPrice float64, previousQuantity int64, startedAt, endedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:               id,
		ParentID:         parentId,
		Billed:           billedType,
		Price:            price,
		Quantity:         quantity,
		PreviousBilled:   previousBilledType,
		PreviousPrice:    previousPrice,
		PreviousQuantity: previousQuantity,
		IsCanceled:       true,
		StartedAt:        startedAt,
		EndedAt:          &endedAt,
	})
}

func TestQueryResolver_Search_Organization_ByOnboardingStatus(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateTenantOrganization(ctx, driver, tenantName, "org excluded")

	today := utils.Now()
	yesterday := today.AddDate(0, 0, -1)

	orgNA := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		OnboardingDetails: entity.OnboardingDetails{
			Status:    entity.OnboardingStatusNotApplicable,
			UpdatedAt: nil,
		},
	})
	orgSuccess := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		OnboardingDetails: entity.OnboardingDetails{
			Status:       entity.OnboardingStatusSuccessful,
			UpdatedAt:    &today,
			SortingOrder: utils.Int64Ptr(60),
		},
	})
	orgStuckYesterday := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		OnboardingDetails: entity.OnboardingDetails{
			Status:       entity.OnboardingStatusStuck,
			UpdatedAt:    &yesterday,
			SortingOrder: utils.Int64Ptr(20),
		},
	})
	orgDone := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		OnboardingDetails: entity.OnboardingDetails{
			Status:       entity.OnboardingStatusDone,
			UpdatedAt:    &today,
			SortingOrder: utils.Int64Ptr(50),
		},
	})
	orgOnTrack := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		OnboardingDetails: entity.OnboardingDetails{
			Status:       entity.OnboardingStatusOnTrack,
			UpdatedAt:    &yesterday,
			SortingOrder: utils.Int64Ptr(40),
		},
	})
	orgStuckToday := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		OnboardingDetails: entity.OnboardingDetails{
			Status:       entity.OnboardingStatusStuck,
			UpdatedAt:    &today,
			SortingOrder: utils.Int64Ptr(20),
		},
	})
	orgLate := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		OnboardingDetails: entity.OnboardingDetails{
			Status:       entity.OnboardingStatusLate,
			UpdatedAt:    &yesterday,
			SortingOrder: utils.Int64Ptr(30),
		},
	})
	orgNotStarted := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		OnboardingDetails: entity.OnboardingDetails{
			Status:       entity.OnboardingStatusNotStarted,
			UpdatedAt:    &yesterday,
			SortingOrder: utils.Int64Ptr(10),
		},
	})

	require.Equal(t, 9, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	assert_Search_Organization_ByOnboardingStatus(t, []string{"DONE"}, []string{orgDone})
	assert_Search_Organization_ByOnboardingStatus(t, []string{"LATE"}, []string{orgLate})
	assert_Search_Organization_ByOnboardingStatus(t, []string{"STUCK"}, []string{orgStuckYesterday, orgStuckToday})
	assert_Search_Organization_ByOnboardingStatus(t, []string{"NOT_STARTED"}, []string{orgNotStarted})
	assert_Search_Organization_ByOnboardingStatus(t, []string{"SUCCESSFUL"}, []string{orgSuccess})
	assert_Search_Organization_ByOnboardingStatus(t, []string{"DONE", "LATE"}, []string{orgLate, orgDone})
	assert_Search_Organization_ByOnboardingStatus(t, []string{"DONE", "LATE", "NOT_APPLICABLE"}, []string{orgLate, orgDone, orgNA})
	assert_Search_Organization_ByOnboardingStatus(t, []string{}, []string{orgNotStarted, orgStuckYesterday, orgStuckToday, orgLate, orgOnTrack, orgDone, orgSuccess, orgNA})
}

func assert_Search_Organization_ByOnboardingStatus(t *testing.T, searchStatuses []string, expectedOrgs []string) {
	query := "/dashboard_view/organization/dashboard_view_organization_filter_by_onboarding_status"
	options := []client.Option{
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", searchStatuses),
	}

	rawResponse, err := c.RawPost(getQuery(query), options...)
	assertRawResponseSuccess(t, rawResponse, err)

	var responseRaw struct {
		DashboardView_Organizations model.OrganizationPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &responseRaw)
	require.Nil(t, err)
	require.NotNil(t, responseRaw)

	require.Equal(t, int64(len(expectedOrgs)), responseRaw.DashboardView_Organizations.TotalElements)
	for i, org := range responseRaw.DashboardView_Organizations.Content {
		require.Equal(t, expectedOrgs[i], org.ID)
	}
}

func TestQueryResolver_Sort_Renewals_ByRenewalDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	daysFromNow10 := utils.Now().AddDate(0, 0, 10)
	daysFromNow20 := utils.Now().AddDate(0, 0, 20)

	_ = neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org1",
	})
	_ = neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org2",
	})
	organizationId3 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org3",
	})

	contractStartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contractId1 := insertContractWithActiveRenewalOpportunity(ctx, driver, organizationId3, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, entity.OpportunityEntity{RenewedAt: daysFromNow20})
	insertServiceLineItem(ctx, driver, contractId1, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	contractId2 := insertContractWithActiveRenewalOpportunity(ctx, driver, organizationId3, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
	}, entity.OpportunityEntity{RenewedAt: daysFromNow10})
	insertServiceLineItem(ctx, driver, contractId2, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 3, "Contract": 2, "Opportunity": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_view_renewals_sort",
		map[string]interface{}{
			"page":    1,
			"limit":   10,
			"sortBy":  "RENEWAL_DATE",
			"sortDir": "ASC",
		})

	var renewalsPageStruct struct {
		DashboardView_Renewals model.RenewalsPage
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &renewalsPageStruct)
	require.Nil(t, err)

	require.Equal(t, 1, renewalsPageStruct.DashboardView_Renewals.TotalPages)
	require.Equal(t, int64(2), renewalsPageStruct.DashboardView_Renewals.TotalAvailable)
	require.Equal(t, int64(2), renewalsPageStruct.DashboardView_Renewals.TotalElements)

	require.Equal(t, organizationId3, renewalsPageStruct.DashboardView_Renewals.Content[0].Organization.ID)
	require.Equal(t, contractId2, renewalsPageStruct.DashboardView_Renewals.Content[0].Contract.ID)
	require.Equal(t, contractId1, renewalsPageStruct.DashboardView_Renewals.Content[1].Contract.ID)
}

func TestQueryResolver_Sort_Renewals_ByForecastAmountASC(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	daysFromNow10 := utils.Now().AddDate(0, 0, 10)
	daysFromNow20 := utils.Now().AddDate(0, 0, 20)

	_ = neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org1",
	})
	_ = neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org2",
	})
	organizationId3 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org3",
	})

	contractStartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contractId1 := insertContractWithActiveRenewalOpportunity(ctx, driver, organizationId3, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
		RenewalCycle:     "MONTHLY",
	}, entity.OpportunityEntity{RenewedAt: daysFromNow20, MaxAmount: 100})
	insertServiceLineItem(ctx, driver, contractId1, entity.BilledTypeAnnually, 3, 2, sli1StartedAt)

	contractId2 := insertContractWithActiveRenewalOpportunity(ctx, driver, organizationId3, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
		RenewalCycle:     "MONTHLY",
	}, entity.OpportunityEntity{RenewedAt: daysFromNow10, MaxAmount: 200})
	insertServiceLineItem(ctx, driver, contractId2, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 3, "Contract": 2, "Opportunity": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_view_renewals_sort",
		map[string]interface{}{
			"page":    1,
			"limit":   10,
			"sortBy":  "FORECAST_ARR",
			"sortDir": "ASC",
		})

	var renewalsPageStruct struct {
		DashboardView_Renewals model.RenewalsPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &renewalsPageStruct)
	require.Nil(t, err)

	require.Equal(t, 1, renewalsPageStruct.DashboardView_Renewals.TotalPages)
	require.Equal(t, int64(2), renewalsPageStruct.DashboardView_Renewals.TotalAvailable)
	require.Equal(t, int64(2), renewalsPageStruct.DashboardView_Renewals.TotalElements)

	require.Equal(t, organizationId3, renewalsPageStruct.DashboardView_Renewals.Content[0].Organization.ID)
	require.Equal(t, contractId1, renewalsPageStruct.DashboardView_Renewals.Content[0].Contract.ID)
	require.Equal(t, contractId2, renewalsPageStruct.DashboardView_Renewals.Content[1].Contract.ID)
}

func TestQueryResolver_Sort_Renewals_ByForecastAmountDESC(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	daysFromNow10 := utils.Now().AddDate(0, 0, 10)
	daysFromNow20 := utils.Now().AddDate(0, 0, 20)

	_ = neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org1",
	})
	_ = neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org2",
	})
	organizationId3 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org3",
	})

	contractStartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contractId1 := insertContractWithActiveRenewalOpportunity(ctx, driver, organizationId3, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
		RenewalCycle:     "MONTHLY",
	}, entity.OpportunityEntity{RenewedAt: daysFromNow20, MaxAmount: 100})
	insertServiceLineItem(ctx, driver, contractId1, entity.BilledTypeAnnually, 3, 2, sli1StartedAt)

	contractId2 := insertContractWithActiveRenewalOpportunity(ctx, driver, organizationId3, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
		RenewalCycle:     "MONTHLY",
	}, entity.OpportunityEntity{RenewedAt: daysFromNow10, MaxAmount: 200})
	insertServiceLineItem(ctx, driver, contractId2, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 3, "Contract": 2, "Opportunity": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_view_renewals_sort",
		map[string]interface{}{
			"page":    1,
			"limit":   10,
			"sortBy":  "FORECAST_ARR",
			"sortDir": "DESC",
		})

	var renewalsPageStruct struct {
		DashboardView_Renewals model.RenewalsPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &renewalsPageStruct)
	require.Nil(t, err)

	require.Equal(t, 1, renewalsPageStruct.DashboardView_Renewals.TotalPages)
	require.Equal(t, int64(2), renewalsPageStruct.DashboardView_Renewals.TotalAvailable)
	require.Equal(t, int64(2), renewalsPageStruct.DashboardView_Renewals.TotalElements)

	require.Equal(t, organizationId3, renewalsPageStruct.DashboardView_Renewals.Content[0].Organization.ID)
	require.Equal(t, contractId1, renewalsPageStruct.DashboardView_Renewals.Content[1].Contract.ID)
	require.Equal(t, contractId2, renewalsPageStruct.DashboardView_Renewals.Content[0].Contract.ID)
}

func TestQueryResolver_Sort_Renewals_ByRenewalLikelihood(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	daysFromNow10 := utils.Now().AddDate(0, 0, 10)
	daysFromNow20 := utils.Now().AddDate(0, 0, 20)

	_ = neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org1",
	})
	_ = neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org2",
	})
	organizationId3 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org3",
	})

	contractStartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contractId1 := insertContractWithActiveRenewalOpportunity(ctx, driver, organizationId3, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
		RenewalCycle:     "MONTHLY",
	}, entity.OpportunityEntity{RenewedAt: daysFromNow20, MaxAmount: 100, RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh})
	insertServiceLineItem(ctx, driver, contractId1, entity.BilledTypeAnnually, 3, 2, sli1StartedAt)

	contractId2 := insertContractWithActiveRenewalOpportunity(ctx, driver, organizationId3, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
		RenewalCycle:     "MONTHLY",
	}, entity.OpportunityEntity{RenewedAt: daysFromNow10, MaxAmount: 200, RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium})
	insertServiceLineItem(ctx, driver, contractId2, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	contractId3 := insertContractWithActiveRenewalOpportunity(ctx, driver, organizationId3, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
		RenewalCycle:     "MONTHLY",
	}, entity.OpportunityEntity{RenewedAt: daysFromNow10, MaxAmount: 200, RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow})
	insertServiceLineItem(ctx, driver, contractId3, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 3, "Contract": 3, "Opportunity": 3})
	require.Equal(t, int64(2), assert_Search_Renewals_By_Likelihood(t, []string{"HIGH_RENEWAL", "MEDIUM_RENEWAL"}).TotalElements)
}

func assert_Search_Renewals_By_Likelihood(t *testing.T, searchTerm []string) model.RenewalsPage {
	rawResponse, err := c.RawPost(getQuery("dashboard_view/dashboard_view_renewals_filter_by_renewal_likelihood"),
		client.Var("page", 1),
		client.Var("limit", 10),
		client.Var("searchTerm", searchTerm),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var responseRaw struct {
		DashboardView_Renewals model.RenewalsPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &responseRaw)
	require.Nil(t, err)
	require.NotNil(t, responseRaw)

	return responseRaw.DashboardView_Renewals
}

func TestQueryResolver_Search_Renewals_By_Owner_In_IncludeEmptyFalse(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	userId1 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	userId2 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	org := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Name: "org",
	})

	contractId1 := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, org, entity.ContractEntity{})
	opportunityId1 := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId1, entity.OpportunityEntity{})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId1, opportunityId1)
	neo4jt.OpportunityOwnedBy(ctx, driver, opportunityId1, userId1)

	contractId2 := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, org, entity.ContractEntity{})
	opportunityId2 := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId2, entity.OpportunityEntity{})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId2, opportunityId2)
	neo4jt.OpportunityOwnedBy(ctx, driver, opportunityId2, userId2)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "OWNS"))

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_view_renewals_filter_by_owner",
		map[string]interface{}{
			"ownerIdList":  []string{userId1},
			"ownerIdEmpty": false,
			"page":         1,
			"limit":        10})

	var renewalsPageStruct struct {
		DashboardView_Renewals model.RenewalsPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &renewalsPageStruct)
	require.Nil(t, err)

	require.Equal(t, int64(2), renewalsPageStruct.DashboardView_Renewals.TotalAvailable)
	require.Equal(t, int64(1), renewalsPageStruct.DashboardView_Renewals.TotalElements)
	require.Equal(t, 1, len(renewalsPageStruct.DashboardView_Renewals.Content))
	require.ElementsMatch(t, []string{opportunityId1},
		[]string{renewalsPageStruct.DashboardView_Renewals.Content[0].Opportunity.ID})
}
