package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphLocationEventHandler_OnLocationCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)

	locationId, _ := uuid.NewUUID()
	locationAggregate := aggregate.NewLocationAggregateWithTenantAndID(tenantName, locationId.String())
	name := "test_location_name"
	rawAddress := "test_location_raw_address"
	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()
	country := "US"
	region := "test_region"
	district := "test_district"
	locality := "test_locality"
	street := "test_street"
	address1 := "test_address1"
	address2 := "test_address2"
	zip := "test_zip"
	addressType := "test_address_type"
	houseNumber := "test_house_number"
	postalCode := "test_postal_code"
	plusFour := "test_plus_four"
	commercial := false
	predirection := "test_prediction"
	timeZone := "test_timezone"
	var utcOffset = 1
	var latitude float64 = 1
	var longitude float64 = 2
	event, err := events.NewLocationCreateEvent(
		locationAggregate,
		name,
		rawAddress,
		cmnmod.Source{
			Source:        constants.SourceOpenline,
			SourceOfTruth: constants.SourceOpenline,
			AppSource:     constants.SourceOpenline,
		},
		createdAt,
		updatedAt,
		models.LocationAddress{
			Country:      country,
			Region:       region,
			District:     district,
			Locality:     locality,
			Street:       street,
			Address1:     address1,
			Address2:     address2,
			Zip:          zip,
			AddressType:  addressType,
			HouseNumber:  houseNumber,
			PostalCode:   postalCode,
			PlusFour:     plusFour,
			Commercial:   commercial,
			Predirection: predirection,
			Latitude:     &latitude,
			Longitude:    &longitude,
			TimeZone:     timeZone,
			UtcOffset:    utcOffset,
		})

	require.Nil(t, err)
	locationEventHandler := &LocationEventHandler{
		Repositories: testDatabase.Repositories,
	}
	err = locationEventHandler.OnLocationCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "Location"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "Location_"+tenantName), "Incorrect number of Location_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "LOCATION_BELONGS_TO_TENANT"), "Incorrect number of LOCATION_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Location_"+tenantName, locationId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, locationId.String(), utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, rawAddress, utils.GetStringPropOrEmpty(props, "rawAddress"))
	require.Equal(t, name, utils.GetStringPropOrEmpty(props, "name"))
	require.Equal(t, country, utils.GetStringPropOrEmpty(props, "country"))
	require.Equal(t, region, utils.GetStringPropOrEmpty(props, "region"))
	require.Equal(t, district, utils.GetStringPropOrEmpty(props, "district"))
	require.Equal(t, locality, utils.GetStringPropOrEmpty(props, "locality"))
	require.Equal(t, street, utils.GetStringPropOrEmpty(props, "street"))
	require.Equal(t, address1, utils.GetStringPropOrEmpty(props, "address"))
	require.Equal(t, address2, utils.GetStringPropOrEmpty(props, "address2"))
	require.Equal(t, addressType, utils.GetStringPropOrEmpty(props, "addressType"))
	require.Equal(t, houseNumber, utils.GetStringPropOrEmpty(props, "houseNumber"))
	require.Equal(t, postalCode, utils.GetStringPropOrEmpty(props, "postalCode"))
	require.Equal(t, commercial, utils.GetBoolPropOrFalse(props, "commercial"))
	require.Equal(t, predirection, utils.GetStringPropOrEmpty(props, "predirection"))
	require.Equal(t, &latitude, utils.GetFloatPropOrNil(props, "latitude"))
	require.Equal(t, &longitude, utils.GetFloatPropOrNil(props, "longitude"))
	require.Equal(t, timeZone, utils.GetStringPropOrEmpty(props, "timeZone"))
	require.Equal(t, int64(utcOffset), utils.GetInt64PropOrZero(props, "utcOffset"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(props, "source"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(props, "sourceOfTruth"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(props, "appSource"))
	require.Equal(t, &createdAt, utils.GetTimePropOrNil(props, "createdAt"))
	require.Equal(t, &updatedAt, utils.GetTimePropOrNil(props, "updatedAt"))
}

func TestGraphLocationEventHandler_OnLocationValidated(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	name := "test_location_name"
	updatedAt := time.Now().UTC()
	country := "US"
	region := "test_region"
	locality := "test_locality"
	street := "test_street"
	address1 := "test_address1"
	address2 := "test_address2"
	zip := "test_zip"
	addressType := "test_address_type"
	rawAddress := "test_location_raw_address"
	houseNumber := "test_house_number"
	postalCode := "test_postal_code"
	plusFour := "test_plus_four"
	commercial := false
	predirection := "test_prediction"
	district := "test_district"
	timeZone := "test_timezone"
	var utcOffset = 1
	var latitude float64 = 1
	var longitude float64 = 2
	locationId := neo4jt.CreateLocation(ctx, testDatabase.Driver, tenantName, entity.LocationEntity{
		Name:          name,
		UpdatedAt:     updatedAt,
		Country:       country,
		Region:        region,
		Locality:      locality,
		Address:       address1,
		Address2:      address2,
		Zip:           zip,
		AddressType:   addressType,
		HouseNumber:   houseNumber,
		PostalCode:    postalCode,
		PlusFour:      plusFour,
		Commercial:    commercial,
		Predirection:  predirection,
		District:      district,
		Street:        street,
		RawAddress:    rawAddress,
		Latitude:      &latitude,
		Longitude:     &longitude,
		TimeZone:      timeZone,
		UtcOffset:     int64(utcOffset),
		SourceOfTruth: constants.SourceOpenline,
		Source:        constants.SourceOpenline,
		AppSource:     constants.SourceOpenline,
	})

	dbNodeAfterCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Location_"+tenantName, locationId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterCreate)

	propsAfterCreate := utils.GetPropsFromNode(*dbNodeAfterCreate)
	require.Equal(t, rawAddress, utils.GetStringPropOrEmpty(propsAfterCreate, "rawAddress"))
	require.Equal(t, name, utils.GetStringPropOrEmpty(propsAfterCreate, "name"))
	require.Equal(t, country, utils.GetStringPropOrEmpty(propsAfterCreate, "country"))
	require.Equal(t, region, utils.GetStringPropOrEmpty(propsAfterCreate, "region"))
	require.Equal(t, district, utils.GetStringPropOrEmpty(propsAfterCreate, "district"))
	require.Equal(t, locality, utils.GetStringPropOrEmpty(propsAfterCreate, "locality"))
	require.Equal(t, street, utils.GetStringPropOrEmpty(propsAfterCreate, "street"))
	require.Equal(t, zip, utils.GetStringPropOrEmpty(propsAfterCreate, "zip"))
	require.Equal(t, plusFour, utils.GetStringPropOrEmpty(propsAfterCreate, "plusFour"))
	require.Equal(t, address1, utils.GetStringPropOrEmpty(propsAfterCreate, "address"))
	require.Equal(t, address2, utils.GetStringPropOrEmpty(propsAfterCreate, "address2"))
	require.Equal(t, addressType, utils.GetStringPropOrEmpty(propsAfterCreate, "addressType"))
	require.Equal(t, houseNumber, utils.GetStringPropOrEmpty(propsAfterCreate, "houseNumber"))
	require.Equal(t, postalCode, utils.GetStringPropOrEmpty(propsAfterCreate, "postalCode"))
	require.Equal(t, commercial, utils.GetBoolPropOrFalse(propsAfterCreate, "commercial"))
	require.Equal(t, predirection, utils.GetStringPropOrEmpty(propsAfterCreate, "predirection"))
	require.Equal(t, &latitude, utils.GetFloatPropOrNil(propsAfterCreate, "latitude"))
	require.Equal(t, &longitude, utils.GetFloatPropOrNil(propsAfterCreate, "longitude"))
	require.Equal(t, timeZone, utils.GetStringPropOrEmpty(propsAfterCreate, "timeZone"))
	require.Equal(t, int64(utcOffset), utils.GetInt64PropOrZero(propsAfterCreate, "utcOffset"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(propsAfterCreate, "source"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(propsAfterCreate, "sourceOfTruth"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(propsAfterCreate, "appSource"))
	require.Equal(t, &updatedAt, utils.GetTimePropOrNil(propsAfterCreate, "updatedAt"))

	locationAggregate := aggregate.NewLocationAggregateWithTenantAndID(tenantName, locationId)
	neo4jt.CreateCountry(ctx, testDatabase.Driver, "US", "USA", "United States", "1")
	event, err := events.NewLocationValidatedEvent(locationAggregate, rawAddress, "US", models.LocationAddress{
		Country:      country,
		Region:       region,
		District:     district,
		Locality:     locality,
		Street:       street,
		Address1:     address1,
		Address2:     address2,
		Zip:          zip,
		AddressType:  addressType,
		HouseNumber:  houseNumber,
		PostalCode:   postalCode,
		PlusFour:     plusFour,
		Commercial:   commercial,
		Predirection: predirection,
		Latitude:     &latitude,
		Longitude:    &longitude,
		TimeZone:     timeZone,
		UtcOffset:    utcOffset,
	})
	require.Nil(t, err)

	locationEventHandler := &LocationEventHandler{
		Repositories: testDatabase.Repositories,
	}
	err = locationEventHandler.OnLocationValidated(context.Background(), event)
	require.Nil(t, err)

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Location_"+tenantName, locationId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, "", utils.GetStringPropOrEmpty(props, "validationError"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "validated"))
	require.NotEqual(t, &updatedAt, utils.GetTimePropOrNil(props, "updatedAt"))
	require.Equal(t, commercial, utils.GetBoolPropOrFalse(props, "commercial"))
	require.Equal(t, country, utils.GetStringPropOrEmpty(props, "country"))
	require.Equal(t, region, utils.GetStringPropOrEmpty(props, "region"))
	require.Equal(t, district, utils.GetStringPropOrEmpty(props, "district"))
	require.Equal(t, locality, utils.GetStringPropOrEmpty(props, "locality"))
	require.Equal(t, street, utils.GetStringPropOrEmpty(props, "street"))
	require.Equal(t, address1, utils.GetStringPropOrEmpty(props, "address"))
	require.Equal(t, address2, utils.GetStringPropOrEmpty(props, "address2"))
	require.Equal(t, zip, utils.GetStringPropOrEmpty(props, "zip"))
	require.Equal(t, addressType, utils.GetStringPropOrEmpty(props, "addressType"))
	require.Equal(t, houseNumber, utils.GetStringPropOrEmpty(props, "houseNumber"))
	require.Equal(t, postalCode, utils.GetStringPropOrEmpty(props, "postalCode"))
	require.Equal(t, plusFour, utils.GetStringPropOrEmpty(props, "plusFour"))
	require.Equal(t, predirection, utils.GetStringPropOrEmpty(props, "predirection"))
	require.Equal(t, &latitude, utils.GetFloatPropOrNil(props, "latitude"))
	require.Equal(t, &longitude, utils.GetFloatPropOrNil(props, "longitude"))
	require.Equal(t, timeZone, utils.GetStringPropOrEmpty(props, "timeZone"))
	require.Equal(t, int64(utcOffset), utils.GetInt64PropOrZero(props, "utcOffset"))
}

func TestGraphLocationEventHandler_OnLocationValidationFailed(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	name := "test_location_name"
	updatedAt := time.Now().UTC()
	country := "US"
	region := "test_region"
	locality := "test_locality"
	street := "test_street"
	address1 := "test_address1"
	address2 := "test_address2"
	zip := "test_zip"
	addressType := "test_address_type"
	rawAddress := "test_location_raw_address"
	houseNumber := "test_house_number"
	postalCode := "test_postal_code"
	plusFour := "test_plus_four"
	commercial := false
	predirection := "test_prediction"
	district := "test_district"
	timeZone := "test_timezone"
	var utcOffset = 1
	var latitude float64 = 1
	var longitude float64 = 2
	locationId := neo4jt.CreateLocation(ctx, testDatabase.Driver, tenantName, entity.LocationEntity{
		Name:          name,
		UpdatedAt:     updatedAt,
		Country:       country,
		Region:        region,
		Locality:      locality,
		Address:       address1,
		Address2:      address2,
		Zip:           zip,
		AddressType:   addressType,
		HouseNumber:   houseNumber,
		PostalCode:    postalCode,
		PlusFour:      plusFour,
		Commercial:    commercial,
		Predirection:  predirection,
		District:      district,
		Street:        street,
		RawAddress:    rawAddress,
		Latitude:      &latitude,
		Longitude:     &longitude,
		TimeZone:      timeZone,
		UtcOffset:     int64(utcOffset),
		SourceOfTruth: constants.SourceOpenline,
		Source:        constants.SourceOpenline,
		AppSource:     constants.SourceOpenline,
	})

	dbNodeAfterCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Location_"+tenantName, locationId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterCreate)

	propsAfterCreate := utils.GetPropsFromNode(*dbNodeAfterCreate)
	require.Equal(t, rawAddress, utils.GetStringPropOrEmpty(propsAfterCreate, "rawAddress"))
	require.Equal(t, name, utils.GetStringPropOrEmpty(propsAfterCreate, "name"))
	require.Equal(t, country, utils.GetStringPropOrEmpty(propsAfterCreate, "country"))
	require.Equal(t, region, utils.GetStringPropOrEmpty(propsAfterCreate, "region"))
	require.Equal(t, district, utils.GetStringPropOrEmpty(propsAfterCreate, "district"))
	require.Equal(t, locality, utils.GetStringPropOrEmpty(propsAfterCreate, "locality"))
	require.Equal(t, street, utils.GetStringPropOrEmpty(propsAfterCreate, "street"))
	require.Equal(t, zip, utils.GetStringPropOrEmpty(propsAfterCreate, "zip"))
	require.Equal(t, plusFour, utils.GetStringPropOrEmpty(propsAfterCreate, "plusFour"))
	require.Equal(t, address1, utils.GetStringPropOrEmpty(propsAfterCreate, "address"))
	require.Equal(t, address2, utils.GetStringPropOrEmpty(propsAfterCreate, "address2"))
	require.Equal(t, addressType, utils.GetStringPropOrEmpty(propsAfterCreate, "addressType"))
	require.Equal(t, houseNumber, utils.GetStringPropOrEmpty(propsAfterCreate, "houseNumber"))
	require.Equal(t, postalCode, utils.GetStringPropOrEmpty(propsAfterCreate, "postalCode"))
	require.Equal(t, commercial, utils.GetBoolPropOrFalse(propsAfterCreate, "commercial"))
	require.Equal(t, predirection, utils.GetStringPropOrEmpty(propsAfterCreate, "predirection"))
	require.Equal(t, &latitude, utils.GetFloatPropOrNil(propsAfterCreate, "latitude"))
	require.Equal(t, &longitude, utils.GetFloatPropOrNil(propsAfterCreate, "longitude"))
	require.Equal(t, timeZone, utils.GetStringPropOrEmpty(propsAfterCreate, "timeZone"))
	require.Equal(t, int64(utcOffset), utils.GetInt64PropOrZero(propsAfterCreate, "utcOffset"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(propsAfterCreate, "source"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(propsAfterCreate, "sourceOfTruth"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(propsAfterCreate, "appSource"))
	require.Equal(t, &updatedAt, utils.GetTimePropOrNil(propsAfterCreate, "updatedAt"))

	locationAggregate := aggregate.NewLocationAggregateWithTenantAndID(tenantName, locationId)
	neo4jt.CreateCountry(ctx, testDatabase.Driver, "US", "USA", "United States", "1")
	validationError := "Phone validation failed with this custom message!"
	event, err := events.NewLocationFailedValidationEvent(locationAggregate, rawAddress, "US", validationError)
	require.Nil(t, err)

	locationEventHandler := &LocationEventHandler{
		Repositories: testDatabase.Repositories,
	}
	err = locationEventHandler.OnLocationValidationFailed(context.Background(), event)
	require.Nil(t, err)

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Location_"+tenantName, locationId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, validationError, utils.GetStringPropOrEmpty(props, "validationError"))
	require.Equal(t, false, utils.GetBoolPropOrFalse(props, "validated"))
	require.NotEqual(t, updatedAt, utils.GetTimePropOrNil(props, "updatedAt"))
}

func TestGraphLocationEventHandler_OnLocationUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	name := "test_location_name"
	updatedAt := time.Now().UTC()
	country := "US"
	region := "test_region"
	locality := "test_locality"
	street := "test_street"
	address1 := "test_address1"
	address2 := "test_address2"
	zip := "test_zip"
	addressType := "test_address_type"
	rawAddress := "test_location_raw_address"
	houseNumber := "test_house_number"
	postalCode := "test_postal_code"
	plusFour := "test_plus_four"
	commercial := false
	predirection := "test_prediction"
	district := "test_district"
	timeZone := "test_timezone"
	var utcOffset = 1
	var latitude float64 = 1
	var longitude float64 = 2
	locationId := neo4jt.CreateLocation(ctx, testDatabase.Driver, tenantName, entity.LocationEntity{
		Name:          name,
		UpdatedAt:     updatedAt,
		Country:       country,
		Region:        region,
		Locality:      locality,
		Address:       address1,
		Address2:      address2,
		Zip:           zip,
		AddressType:   addressType,
		HouseNumber:   houseNumber,
		PostalCode:    postalCode,
		PlusFour:      plusFour,
		Commercial:    commercial,
		Predirection:  predirection,
		District:      district,
		Street:        street,
		RawAddress:    rawAddress,
		Latitude:      &latitude,
		Longitude:     &longitude,
		TimeZone:      timeZone,
		UtcOffset:     int64(utcOffset),
		SourceOfTruth: constants.SourceOpenline,
		Source:        constants.SourceOpenline,
		AppSource:     constants.SourceOpenline,
	})

	dbNodeAfterCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Location_"+tenantName, locationId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterCreate)

	propsAfterCreate := utils.GetPropsFromNode(*dbNodeAfterCreate)
	require.Equal(t, rawAddress, utils.GetStringPropOrEmpty(propsAfterCreate, "rawAddress"))
	require.Equal(t, name, utils.GetStringPropOrEmpty(propsAfterCreate, "name"))
	require.Equal(t, country, utils.GetStringPropOrEmpty(propsAfterCreate, "country"))
	require.Equal(t, region, utils.GetStringPropOrEmpty(propsAfterCreate, "region"))
	require.Equal(t, district, utils.GetStringPropOrEmpty(propsAfterCreate, "district"))
	require.Equal(t, locality, utils.GetStringPropOrEmpty(propsAfterCreate, "locality"))
	require.Equal(t, street, utils.GetStringPropOrEmpty(propsAfterCreate, "street"))
	require.Equal(t, zip, utils.GetStringPropOrEmpty(propsAfterCreate, "zip"))
	require.Equal(t, plusFour, utils.GetStringPropOrEmpty(propsAfterCreate, "plusFour"))
	require.Equal(t, address1, utils.GetStringPropOrEmpty(propsAfterCreate, "address"))
	require.Equal(t, address2, utils.GetStringPropOrEmpty(propsAfterCreate, "address2"))
	require.Equal(t, addressType, utils.GetStringPropOrEmpty(propsAfterCreate, "addressType"))
	require.Equal(t, houseNumber, utils.GetStringPropOrEmpty(propsAfterCreate, "houseNumber"))
	require.Equal(t, postalCode, utils.GetStringPropOrEmpty(propsAfterCreate, "postalCode"))
	require.Equal(t, commercial, utils.GetBoolPropOrFalse(propsAfterCreate, "commercial"))
	require.Equal(t, predirection, utils.GetStringPropOrEmpty(propsAfterCreate, "predirection"))
	require.Equal(t, &latitude, utils.GetFloatPropOrNil(propsAfterCreate, "latitude"))
	require.Equal(t, &longitude, utils.GetFloatPropOrNil(propsAfterCreate, "longitude"))
	require.Equal(t, timeZone, utils.GetStringPropOrEmpty(propsAfterCreate, "timeZone"))
	require.Equal(t, int64(utcOffset), utils.GetInt64PropOrZero(propsAfterCreate, "utcOffset"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(propsAfterCreate, "source"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(propsAfterCreate, "sourceOfTruth"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(propsAfterCreate, "appSource"))
	require.Equal(t, &updatedAt, utils.GetTimePropOrNil(propsAfterCreate, "updatedAt"))

	locationAggregate := aggregate.NewLocationAggregateWithTenantAndID(tenantName, locationId)
	neo4jt.CreateCountry(ctx, testDatabase.Driver, "US", "USA", "United States", "1")
	updatedAtUpdate := time.Now().UTC()
	locationAddressLatitudeUpdate := 1.1
	locationAddressLongitudeUpdate := 2.2
	locationAddressCountryUpdate := "locationAddressCountryUpdate"
	locationAddressRegionUpdate := "locationAddressRegionUpdate"
	locationAddressDistrictUpdate := "locationAddressDistrictUpdate"
	locationAddressLocalityUpdate := "locationAddressLocalityUpdate"
	locationAddressStreetUpdate := "locationAddressStreetUpdate"
	locationAddressAddress1Update := "locationAddressAddress1Update"
	locationAddressAddress2Update := "locationAddressAddress2Update"
	locationAddressZipUpdate := "locationAddressZipUpdate"
	locationAddressAddressType := "locationAddressAddressType"
	locationAddressHouseNumber := "locationAddressHouseNumber"
	locationAddressPostalCodeUpdate := "locationAddressPostalCodeUpdate"
	locationAddressPlusFourUpdate := "locationAddressPlusFourUpdate"
	locationAddressCommercialUpdate := true
	locationAddressPredirectionUpdate := "locationAddressPredirectionUpdate"
	locationAddressTimeZoneUpdate := "locationAddressTimeZoneUpdate"
	locationAddressUtcOffsetUpdate := 1
	locationUpdateEvent, err := events.NewLocationUpdateEvent(locationAggregate, name, rawAddress, constants.SourceOpenline, updatedAtUpdate, models.LocationAddress{
		Country:      locationAddressCountryUpdate,
		Region:       locationAddressRegionUpdate,
		District:     locationAddressDistrictUpdate,
		Locality:     locationAddressLocalityUpdate,
		Street:       locationAddressStreetUpdate,
		Address1:     locationAddressAddress1Update,
		Address2:     locationAddressAddress2Update,
		Zip:          locationAddressZipUpdate,
		AddressType:  locationAddressAddressType,
		HouseNumber:  locationAddressHouseNumber,
		PostalCode:   locationAddressPostalCodeUpdate,
		PlusFour:     locationAddressPlusFourUpdate,
		Commercial:   locationAddressCommercialUpdate,
		Predirection: locationAddressPredirectionUpdate,
		Latitude:     &locationAddressLatitudeUpdate,
		Longitude:    &locationAddressLongitudeUpdate,
		TimeZone:     locationAddressTimeZoneUpdate,
		UtcOffset:    locationAddressUtcOffsetUpdate,
	})
	require.Nil(t, err)

	locationEventHandler := &LocationEventHandler{
		Repositories: testDatabase.Repositories,
	}
	err = locationEventHandler.OnLocationUpdate(context.Background(), locationUpdateEvent)
	require.Nil(t, err)

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Location_"+tenantName, locationId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	locationUpdateProps := utils.GetPropsFromNode(*dbNode)
	require.Equal(t, 27, len(locationUpdateProps))

	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(locationUpdateProps, "appSource"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(locationUpdateProps, "source"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(locationUpdateProps, "sourceOfTruth"))
	require.Equal(t, locationAddressStreetUpdate, utils.GetStringPropOrEmpty(locationUpdateProps, "street"))
	require.Equal(t, locationAddressZipUpdate, utils.GetStringPropOrEmpty(locationUpdateProps, "zip"))
	require.Equal(t, locationAddressAddressType, utils.GetStringPropOrEmpty(locationUpdateProps, "addressType"))
	require.Equal(t, locationAddressLongitudeUpdate, utils.GetFloatPropOrZero(locationUpdateProps, "longitude"))
	require.Equal(t, locationAddressPlusFourUpdate, utils.GetStringPropOrEmpty(locationUpdateProps, "plusFour"))
	require.Equal(t, locationAddressAddress2Update, utils.GetStringPropOrEmpty(locationUpdateProps, "address2"))
	require.Equal(t, locationAddressLocalityUpdate, utils.GetStringPropOrEmpty(locationUpdateProps, "locality"))
	require.Equal(t, locationAddressRegionUpdate, utils.GetStringPropOrEmpty(locationUpdateProps, "region"))
	require.Equal(t, locationAddressCommercialUpdate, utils.GetBoolPropOrFalse(locationUpdateProps, "commercial"))
	require.Equal(t, locationAddressLatitudeUpdate, utils.GetFloatPropOrZero(locationUpdateProps, "latitude"))
	require.Equal(t, locationAddressPredirectionUpdate, utils.GetStringPropOrEmpty(locationUpdateProps, "predirection"))
	require.Equal(t, locationAddressAddress1Update, utils.GetStringPropOrEmpty(locationUpdateProps, "address"))
	require.Equal(t, locationAddressPostalCodeUpdate, utils.GetStringPropOrEmpty(locationUpdateProps, "postalCode"))
	require.Equal(t, locationAddressHouseNumber, utils.GetStringPropOrEmpty(locationUpdateProps, "houseNumber"))
	require.Equal(t, locationAddressUtcOffsetUpdate, int(utils.GetInt64PropOrZero(locationUpdateProps, "utcOffset")))
	require.Equal(t, name, utils.GetStringPropOrEmpty(locationUpdateProps, "name"))
	require.Equal(t, locationAddressCountryUpdate, utils.GetStringPropOrEmpty(locationUpdateProps, "country"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(locationUpdateProps, "syncedWithEventStore"))
	creationTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	require.Equal(t, creationTime, utils.GetTimePropOrNow(locationUpdateProps, "createdAt"))
	require.Equal(t, locationAddressTimeZoneUpdate, utils.GetStringPropOrEmpty(locationUpdateProps, "timeZone"))
	require.Equal(t, rawAddress, utils.GetStringPropOrEmpty(locationUpdateProps, "rawAddress"))
	require.Less(t, updatedAt, utils.GetTimePropOrNow(locationUpdateProps, "updatedAt"))
	require.Equal(t, locationAddressDistrictUpdate, utils.GetStringPropOrEmpty(locationUpdateProps, "district"))
}
