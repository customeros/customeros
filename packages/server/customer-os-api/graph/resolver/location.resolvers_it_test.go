package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_LocationUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	locationId := neo4jt.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{})

	rawResponse, err := c.RawPost(getQuery("location/update_location"),
		client.Var("locationId", locationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var locationStruct struct {
		Location_Update model.Location
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &locationStruct)
	require.Nil(t, err)

	updatedLocation := locationStruct.Location_Update

	require.Equal(t, locationId, updatedLocation.ID)
	test.AssertRecentTime(t, updatedLocation.UpdatedAt)
	require.Equal(t, model.DataSourceOpenline, updatedLocation.SourceOfTruth)
	require.Equal(t, "name", *updatedLocation.Name)
	require.Equal(t, "rawAddress", *updatedLocation.RawAddress)
	require.Equal(t, "country", *updatedLocation.Country)
	require.Equal(t, "region", *updatedLocation.Region)
	require.Equal(t, "district", *updatedLocation.District)
	require.Equal(t, "locality", *updatedLocation.Locality)
	require.Equal(t, "street", *updatedLocation.Street)
	require.Equal(t, "address", *updatedLocation.Address)
	require.Equal(t, "address2", *updatedLocation.Address2)
	require.Equal(t, "zip", *updatedLocation.Zip)
	require.Equal(t, "addressType", *updatedLocation.AddressType)
	require.Equal(t, "houseNumber", *updatedLocation.HouseNumber)
	require.Equal(t, "postalCode", *updatedLocation.PostalCode)
	require.Equal(t, "plusFour", *updatedLocation.PlusFour)
	require.Equal(t, true, *updatedLocation.Commercial)
	require.Equal(t, "predirection", *updatedLocation.Predirection)
	require.Equal(t, 1.0, *updatedLocation.Latitude)
	require.Equal(t, -2.0, *updatedLocation.Longitude)
	require.Equal(t, "timeZone", *updatedLocation.TimeZone)
	require.Equal(t, 3.0, *updatedLocation.UtcOffset)

	// Check the number of nodes in the Neo4j database
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Location_"+tenantName))
}

func TestMutationResolver_LocationRemoveFromOrganization_UniqueRelation(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	locationId := neo4jt.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{})
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org")
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId, locationId)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Location": 1, "Organization": 1, "Location_" + tenantName: 1, "Organization_" + tenantName: 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, driver, map[string]int{"ASSOCIATED_WITH": 1})

	rawResponse := callGraphQL(t, "location/remove_location_from_organization", map[string]interface{}{
		"locationId":     locationId,
		"organizationId": organizationId,
	})

	var organizationStruct struct {
		Location_RemoveFromOrganization model.Organization
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)

	org := organizationStruct.Location_RemoveFromOrganization
	require.Equal(t, organizationId, org.ID)
	require.Equal(t, "org", org.Name)

	// Check the number of nodes in the Neo4j database
	neo4jtest.AssertNeo4jRelationCount(ctx, t, driver, map[string]int{"ASSOCIATED_WITH": 0})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Location":                   0,
		"Location_" + tenantName:     0,
		"Organization":               1,
		"Organization_" + tenantName: 1})
}

func TestMutationResolver_LocationRemoveFromOrganization_SharedLocation(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	locationId := neo4jt.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{})
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org2")
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId1, locationId)
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId2, locationId)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Location":                   1,
		"Location_" + tenantName:     1,
		"Organization":               2,
		"Organization_" + tenantName: 2})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, driver, map[string]int{"ASSOCIATED_WITH": 2})

	rawResponse := callGraphQL(t, "location/remove_location_from_organization", map[string]interface{}{
		"locationId":     locationId,
		"organizationId": organizationId1,
	})

	var organizationStruct struct {
		Location_RemoveFromOrganization model.Organization
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)

	org := organizationStruct.Location_RemoveFromOrganization
	require.Equal(t, organizationId1, org.ID)
	require.Equal(t, "org1", org.Name)

	// Check the number of nodes in the Neo4j database
	neo4jtest.AssertNeo4jRelationCount(ctx, t, driver, map[string]int{"ASSOCIATED_WITH": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Location":                   1,
		"Location_" + tenantName:     1,
		"Organization":               2,
		"Organization_" + tenantName: 2})
}
