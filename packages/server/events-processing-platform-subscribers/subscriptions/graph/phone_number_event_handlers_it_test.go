package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/mocked_grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	commonEvents "github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphPhoneNumberEventHandler_OnPhoneNumberCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	phoneNumberEventHandler := &PhoneNumberEventHandler{
		repositories: testDatabase.Repositories,
	}
	phoneNumberId, _ := uuid.NewUUID()
	phoneNumberAggregate := aggregate.NewPhoneNumberAggregateWithTenantAndID(tenantName, phoneNumberId.String())
	phoneNumber := "+0123456789"
	curTime := utils.Now()
	event, err := events.NewPhoneNumberCreateEvent(phoneNumberAggregate, tenantName, phoneNumber, commonEvents.Source{
		Source:        constants.SourceOpenline,
		SourceOfTruth: constants.SourceOpenline,
		AppSource:     "test",
	}, curTime, curTime)
	require.Nil(t, err)
	err = phoneNumberEventHandler.OnPhoneNumberCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, neo4jutil.NodeLabelPhoneNumber))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, neo4jutil.NodeLabelPhoneNumber+"_"+tenantName), "Incorrect number of PhoneNumber_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "PHONE_NUMBER_BELONGS_TO_TENANT"), "Incorrect number of PHONE_NUMBER_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId.String())
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

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	validated := false
	e164 := "+0123456789"
	phoneNumberId := neo4jtest.CreatePhoneNumber(ctx, testDatabase.Driver, tenantName, neo4jentity.PhoneNumberEntity{
		E164:           e164,
		Validated:      &validated,
		RawPhoneNumber: e164,
		Source:         constants.SourceOpenline,
		SourceOfTruth:  constants.SourceOpenline,
		AppSource:      constants.SourceOpenline,
	})

	dbNodeAfterCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
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

	phoneNumberEventHandler := &PhoneNumberEventHandler{
		repositories: testDatabase.Repositories,
	}
	err = phoneNumberEventHandler.OnPhoneNumberValidated(context.Background(), event)
	require.Nil(t, err)

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
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

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	validated := false
	e164 := "+0123456789"
	phoneNumberId := neo4jtest.CreatePhoneNumber(ctx, testDatabase.Driver, tenantName, neo4jentity.PhoneNumberEntity{
		E164:           e164,
		Validated:      &validated,
		RawPhoneNumber: e164,
		Source:         constants.SourceOpenline,
		SourceOfTruth:  constants.SourceOpenline,
		AppSource:      constants.SourceOpenline,
	})

	dbNodeAfterCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
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

	phoneNumberEventHandler := &PhoneNumberEventHandler{
		repositories: testDatabase.Repositories,
	}
	err = phoneNumberEventHandler.OnPhoneNumberValidationFailed(context.Background(), event)
	require.Nil(t, err)

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
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

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	validated := false
	e164 := "+0123456789"
	phoneNumberId := neo4jtest.CreatePhoneNumber(ctx, testDatabase.Driver, tenantName, neo4jentity.PhoneNumberEntity{
		E164:           e164,
		Validated:      &validated,
		RawPhoneNumber: e164,
		Source:         constants.SourceOpenline,
		SourceOfTruth:  constants.SourceOpenline,
		AppSource:      constants.SourceOpenline,
	})

	dbNodeAfterCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterCreate)
	propsAfterCreate := utils.GetPropsFromNode(*dbNodeAfterCreate)
	require.Equal(t, false, utils.GetBoolPropOrFalse(propsAfterCreate, "validated"))
	creationTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	creationUpdatedAt := utils.GetTimePropOrNil(propsAfterCreate, "updatedAt")
	require.Equal(t, &creationTime, creationUpdatedAt)

	phoneNumberAggregate := aggregate.NewPhoneNumberAggregateWithTenantAndID(tenantName, phoneNumberId)
	neo4jt.CreateCountry(ctx, testDatabase.Driver, "US", "USA", "United States", "1")
	updatedAtUpdate := utils.Now()

	phoneNumberEventHandler := &PhoneNumberEventHandler{
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// prepare grpc mock
	calledPhoneNumberValidateRequest := false
	phoneNumberCallbacks := mocked_grpc.MockPhoneNumberServiceCallbacks{
		RequestPhoneNumberValidation: func(context context.Context, op *phonenumberpb.RequestPhoneNumberValidationGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, phoneNumberId, op.Id)
			calledPhoneNumberValidateRequest = true
			return &phonenumberpb.PhoneNumberIdGrpcResponse{}, nil
		},
	}
	mocked_grpc.SetPhoneNumberCallbacks(&phoneNumberCallbacks)

	event, err := events.NewPhoneNumberUpdateEvent(phoneNumberAggregate, tenantName, constants.SourceOpenline, "+998877", updatedAtUpdate)
	require.Nil(t, err)

	err = phoneNumberEventHandler.OnPhoneNumberUpdate(context.Background(), event)
	require.Nil(t, err)

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	phoneUpdateProps := utils.GetPropsFromNode(*dbNode)
	require.Equal(t, 9, len(phoneUpdateProps))

	require.Less(t, *creationUpdatedAt, utils.GetTimePropOrNow(phoneUpdateProps, "updatedAt"))
	require.Equal(t, creationTime, utils.GetTimePropOrNow(phoneUpdateProps, "createdAt"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(phoneUpdateProps, "syncedWithEventStore"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(phoneUpdateProps, "source"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(phoneUpdateProps, "sourceOfTruth"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(phoneUpdateProps, "appSource"))
	require.Equal(t, "+998877", utils.GetStringPropOrEmpty(phoneUpdateProps, "rawPhoneNumber"))
	require.Equal(t, "", utils.GetStringPropOrEmpty(phoneUpdateProps, "e164"))
	require.Equal(t, false, utils.GetBoolPropOrFalse(phoneUpdateProps, "validated"))

	require.True(t, calledPhoneNumberValidateRequest)
}
