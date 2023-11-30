package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphPhoneNumberEventHandler_OnPhoneNumberCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	phoneNumberEventHandler := &GraphPhoneNumberEventHandler{
		Repositories: testDatabase.Repositories,
	}
	phoneNumberId, _ := uuid.NewUUID()
	phoneNumberAggregate := aggregate.NewPhoneNumberAggregateWithTenantAndID(tenantName, phoneNumberId.String())
	phoneNumber := "+0123456789"
	curTime := time.Now().UTC()
	event, err := events.NewPhoneNumberCreateEvent(phoneNumberAggregate, tenantName, phoneNumber, cmnmod.Source{
		Source:        constants.SourceOpenline,
		SourceOfTruth: constants.SourceOpenline,
		AppSource:     "test",
	}, curTime, curTime)
	require.Nil(t, err)
	err = phoneNumberEventHandler.OnPhoneNumberCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName), "Incorrect number of PhoneNumber_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "PHONE_NUMBER_BELONGS_TO_TENANT"), "Incorrect number of PHONE_NUMBER_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, phoneNumberId.String(), utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, phoneNumber, utils.GetStringPropOrEmpty(props, "rawPhoneNumber"))
	require.Equal(t, "test", utils.GetStringPropOrEmpty(props, "appSource"))
}

func TestGraphPhoneNumberEventHandler_OnPhoneNumberValidated(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	validated := false
	e164 := "+0123456789"
	phoneNumberId := neo4jt.CreatePhoneNumber(ctx, testDatabase.Driver, tenantName, entity.PhoneNumberEntity{
		E164:           e164,
		Validated:      &validated,
		RawPhoneNumber: e164,
		Source:         constants.SourceOpenline,
		SourceOfTruth:  constants.SourceOpenline,
		AppSource:      constants.SourceOpenline,
	})

	dbNodeAfterCreate, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterCreate)
	propsAfterCreate := utils.GetPropsFromNode(*dbNodeAfterCreate)
	require.Equal(t, false, utils.GetBoolPropOrFalse(propsAfterCreate, "validated"))
	creationTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	require.Equal(t, &creationTime, utils.GetTimePropOrNil(propsAfterCreate, "updatedAt"))

	phoneNumberAggregate := aggregate.NewPhoneNumberAggregateWithTenantAndID(tenantName, phoneNumberId)
	neo4jt.CreateCountry(ctx, testDatabase.Driver, "US", "USA", "United States", "1")
	event, err := events.NewPhoneNumberValidatedEvent(phoneNumberAggregate, tenantName, e164, e164, "US")
	require.Nil(t, err)

	phoneNumberEventHandler := &GraphPhoneNumberEventHandler{
		Repositories: testDatabase.Repositories,
	}
	err = phoneNumberEventHandler.OnPhoneNumberValidated(context.Background(), event)
	require.Nil(t, err)

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, e164, utils.GetStringPropOrEmpty(props, "e164"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "validated"))
	require.NotEqual(t, &creationTime, utils.GetTimePropOrNil(props, "updatedAt"))
	require.Equal(t, e164, utils.GetStringPropOrEmpty(props, "rawPhoneNumber"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(props, "source"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(props, "sourceOfTruth"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(props, "appSource"))
}

func TestGraphPhoneNumberEventHandler_OnPhoneNumberValidationFailed(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	validated := false
	e164 := "+0123456789"
	phoneNumberId := neo4jt.CreatePhoneNumber(ctx, testDatabase.Driver, tenantName, entity.PhoneNumberEntity{
		E164:           e164,
		Validated:      &validated,
		RawPhoneNumber: e164,
		Source:         constants.SourceOpenline,
		SourceOfTruth:  constants.SourceOpenline,
		AppSource:      constants.SourceOpenline,
	})

	dbNodeAfterCreate, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterCreate)
	propsAfterCreate := utils.GetPropsFromNode(*dbNodeAfterCreate)
	require.Equal(t, false, utils.GetBoolPropOrFalse(propsAfterCreate, "validated"))
	creationTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	require.Equal(t, &creationTime, utils.GetTimePropOrNil(propsAfterCreate, "updatedAt"))

	phoneNumberAggregate := aggregate.NewPhoneNumberAggregateWithTenantAndID(tenantName, phoneNumberId)
	neo4jt.CreateCountry(ctx, testDatabase.Driver, "US", "USA", "United States", "1")
	validationError := "Phone validation failed with this custom message!"
	event, err := events.NewPhoneNumberFailedValidationEvent(phoneNumberAggregate, tenantName, e164, "US", validationError)
	require.Nil(t, err)

	phoneNumberEventHandler := &GraphPhoneNumberEventHandler{
		Repositories: testDatabase.Repositories,
	}
	err = phoneNumberEventHandler.OnPhoneNumberValidationFailed(context.Background(), event)
	require.Nil(t, err)

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, e164, utils.GetStringPropOrEmpty(props, "e164"))
	require.Equal(t, false, utils.GetBoolPropOrFalse(props, "validated"))
	require.Equal(t, validationError, utils.GetStringPropOrEmpty(props, "validationError"))
	require.NotEqual(t, &creationTime, utils.GetTimePropOrNil(props, "updatedAt"))
	require.Equal(t, e164, utils.GetStringPropOrEmpty(props, "rawPhoneNumber"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(props, "source"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(props, "sourceOfTruth"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(props, "appSource"))
}

func TestGraphPhoneNumberEventHandler_OnPhoneNumberUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	validated := false
	e164 := "+0123456789"
	phoneNumberId := neo4jt.CreatePhoneNumber(ctx, testDatabase.Driver, tenantName, entity.PhoneNumberEntity{
		E164:           e164,
		Validated:      &validated,
		RawPhoneNumber: e164,
		Source:         constants.SourceOpenline,
		SourceOfTruth:  constants.SourceOpenline,
		AppSource:      constants.SourceOpenline,
	})

	dbNodeAfterCreate, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterCreate)
	propsAfterCreate := utils.GetPropsFromNode(*dbNodeAfterCreate)
	require.Equal(t, false, utils.GetBoolPropOrFalse(propsAfterCreate, "validated"))
	creationTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	creationUpdatedAt := utils.GetTimePropOrNil(propsAfterCreate, "updatedAt")
	require.Equal(t, &creationTime, creationUpdatedAt)

	phoneNumberAggregate := aggregate.NewPhoneNumberAggregateWithTenantAndID(tenantName, phoneNumberId)
	neo4jt.CreateCountry(ctx, testDatabase.Driver, "US", "USA", "United States", "1")
	updatedAtUpdate := time.Now().UTC()
	event, err := events.NewPhoneNumberUpdateEvent(phoneNumberAggregate, tenantName, constants.SourceWebscrape, updatedAtUpdate)
	require.Nil(t, err)

	phoneNumberEventHandler := &GraphPhoneNumberEventHandler{
		Repositories: testDatabase.Repositories,
	}
	err = phoneNumberEventHandler.OnPhoneNumberUpdate(context.Background(), event)
	require.Nil(t, err)

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	phoneUpdateProps := utils.GetPropsFromNode(*dbNode)
	require.Equal(t, 10, len(phoneUpdateProps))

	require.Less(t, *creationUpdatedAt, utils.GetTimePropOrNow(phoneUpdateProps, "updatedAt"))
	require.Equal(t, creationTime, utils.GetTimePropOrNow(phoneUpdateProps, "createdAt"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(phoneUpdateProps, "syncedWithEventStore"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(phoneUpdateProps, "source"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(phoneUpdateProps, "sourceOfTruth"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(phoneUpdateProps, "appSource"))
	require.Equal(t, e164, utils.GetStringPropOrEmpty(phoneUpdateProps, "rawPhoneNumber"))
	require.Equal(t, e164, utils.GetStringPropOrEmpty(phoneUpdateProps, "e164"))
	require.Equal(t, false, utils.GetBoolPropOrFalse(phoneUpdateProps, "validated"))
}
