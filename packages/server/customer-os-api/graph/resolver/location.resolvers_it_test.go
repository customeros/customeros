package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_LocationUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	locationId := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{})

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
	test.AssertTimeRecentlyChanged(t, updatedLocation.UpdatedAt)
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
	require.Equal(t, int64(3), *updatedLocation.UtcOffset)

	// Check the number of nodes in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Location_"+tenantName))
}
