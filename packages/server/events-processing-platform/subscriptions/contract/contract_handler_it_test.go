package contract

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	opportunityaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	opportunitymodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	servicelineitemmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/mocked_grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestContractEventHandler_UpdateRenewalNextCycleDate_CreateRenewalOpportunityWhenMissing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
	})

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityNextCycleDate(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check create renewal opportunity command was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))

	generatedEvent1 := eventList[0]
	require.Equal(t, opportunityevent.OpportunityCreateRenewalV1, generatedEvent1.EventType)
	var eventData1 opportunityevent.OpportunityCreateRenewalEvent
	err = generatedEvent1.GetJsonData(&eventData1)
	require.Nil(t, err)
	require.Equal(t, tenantName, eventData1.Tenant)
	require.Equal(t, contractId, eventData1.ContractId)
	require.Equal(t, constants.SourceOpenline, eventData1.Source.Source)
	test.AssertRecentTime(t, eventData1.CreatedAt)
	test.AssertRecentTime(t, eventData1.UpdatedAt)
}

func TestContractEventHandler_UpdateRenewalNextCycleDate_MonthlyContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate := false
	opportunityServiceRefreshCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunityNextCycleDate: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.OpportunityId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.AppSource)
			require.Equal(t, startOfNextMonth(utils.Now()), *utils.TimestampProtoToTimePtr(op.RenewedAt))
			calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityServiceRefreshCallbacks)

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityNextCycleDate(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))
}

func TestContractEventHandler_UpdateRenewalNextCycleDate_QuarterlyContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	yesterday := utils.Now().AddDate(0, 0, -1)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: &yesterday,
		RenewalCycle:     string(model.QuarterlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate := false
	opportunityServiceRefreshCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunityNextCycleDate: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.OpportunityId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.AppSource)
			in1Quarter := yesterday.AddDate(0, 3, 0)
			require.Equal(t, in1Quarter, *utils.TimestampProtoToTimePtr(op.RenewedAt))
			calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityServiceRefreshCallbacks)

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityNextCycleDate(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))
}

func TestContractEventHandler_UpdateRenewalNextCycleDate_AnnualContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.AnnuallyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// prepare grpc client
	calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate := false
	opportunityServiceRefreshCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunityNextCycleDate: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.OpportunityId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.AppSource)
			require.Equal(t, startOfNextYear(utils.Now()), *utils.TimestampProtoToTimePtr(op.RenewedAt))
			calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityServiceRefreshCallbacks)

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityNextCycleDate(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))
}

func TestContractEventHandler_UpdateRenewalNextCycleDate_MultiAnnualContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	yesterday := utils.Now().AddDate(0, 0, -1)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: &yesterday,
		RenewalCycle:     string(model.AnnuallyRenewalCycleString),
		RenewalPeriods:   utils.Int64Ptr(10),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// prepare grpc client
	calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate := false
	opportunityServiceRefreshCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunityNextCycleDate: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.OpportunityId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.AppSource)
			in10Years := yesterday.AddDate(10, 0, 0)
			require.Equal(t, in10Years, *utils.TimestampProtoToTimePtr(op.RenewedAt))
			calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityServiceRefreshCallbacks)

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityNextCycleDate(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))
}

func prepareRenewalOpportunity(t *testing.T, tenant, opportunityId string, aggregateStore *eventstoret.TestAggregateStore) {
	// prepare aggregate
	opportunityAggregate := opportunityaggregate.NewOpportunityAggregateWithTenantAndID(tenant, opportunityId)
	createEvent := eventstore.NewBaseEvent(opportunityAggregate, opportunityevent.OpportunityCreateRenewalV1)
	preconfiguredEventData := opportunityevent.OpportunityCreateRenewalEvent{
		Tenant:        tenant,
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	}
	err := createEvent.SetJsonData(&preconfiguredEventData)
	require.Nil(t, err)
	opportunityAggregate.UncommittedEvents = []eventstore.Event{
		createEvent,
	}
	err = aggregateStore.Save(context.Background(), opportunityAggregate)
	require.Nil(t, err)
}

func startOfNextMonth(current time.Time) time.Time {
	year, month, _ := current.Date()
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, current.Location())

	// Handle December to January transition
	if month == time.December {
		firstOfNextMonth = time.Date(year+1, time.January, 1, 0, 0, 0, 0, current.Location())
	}
	return firstOfNextMonth
}

func startOfNextYear(current time.Time) time.Time {
	year, _, _ := current.Date()
	firstOfNextYear := time.Date(year+1, time.January, 1, 0, 0, 0, 0, current.Location())
	return firstOfNextYear
}

func TestContractEventHandler_UpdateRenewalArrForecast_OnlyOnceBilled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.AnnuallyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringHigh),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(10),
		Billed:    string(servicelineitemmodel.OnceBilledString),
		CreatedAt: utils.Now(),
	})

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityArr(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check create renewal opportunity command was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 2, len(eventList))

	updateNextCycleDateEvent := eventList[1]
	require.Equal(t, opportunityevent.OpportunityUpdateV1, updateNextCycleDateEvent.EventType)
	var eventData opportunityevent.OpportunityUpdateEvent
	err = updateNextCycleDateEvent.GetJsonData(&eventData)
	require.Nil(t, err)
	require.Equal(t, tenantName, eventData.Tenant)
	require.Equal(t, float64(0), eventData.Amount)
	require.Equal(t, float64(0), eventData.MaxAmount)
	require.Equal(t, []string{opportunitymodel.FieldMaskAmount, opportunitymodel.FieldMaskMaxAmount}, eventData.FieldsMask)
}

func TestContractEventHandler_UpdateRenewalArrForecast_MultipleServices(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringHigh),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(1000000),
		Billed:    string(servicelineitemmodel.OnceBilledString),
		CreatedAt: utils.Now(),
	})
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(10),
		Quantity:  int64(5),
		Billed:    string(servicelineitemmodel.MonthlyBilledString),
		CreatedAt: utils.Now(),
	})
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(1000),
		Quantity:  int64(2),
		Billed:    string(servicelineitemmodel.AnnuallyBilledString),
		CreatedAt: utils.Now(),
	})

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityArr(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check create renewal opportunity command was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 2, len(eventList))

	updateNextCycleDateEvent := eventList[1]
	require.Equal(t, opportunityevent.OpportunityUpdateV1, updateNextCycleDateEvent.EventType)
	var eventData opportunityevent.OpportunityUpdateEvent
	err = updateNextCycleDateEvent.GetJsonData(&eventData)
	require.Nil(t, err)
	require.Equal(t, tenantName, eventData.Tenant)
	require.Equal(t, float64(2600), eventData.Amount)
	require.Equal(t, float64(2600), eventData.MaxAmount)
	require.Equal(t, []string{opportunitymodel.FieldMaskAmount, opportunitymodel.FieldMaskMaxAmount}, eventData.FieldsMask)
}

func TestContractEventHandler_UpdateRenewalArrForecast_MediumLikelihood(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringMedium),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(1000),
		Quantity:  int64(4),
		Billed:    string(servicelineitemmodel.AnnuallyBilledString),
		CreatedAt: utils.Now(),
	})

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityArr(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check create renewal opportunity command was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 2, len(eventList))

	updateNextCycleDateEvent := eventList[1]
	require.Equal(t, opportunityevent.OpportunityUpdateV1, updateNextCycleDateEvent.EventType)
	var eventData opportunityevent.OpportunityUpdateEvent
	err = updateNextCycleDateEvent.GetJsonData(&eventData)
	require.Nil(t, err)
	require.Equal(t, tenantName, eventData.Tenant)
	require.Equal(t, float64(4000), eventData.MaxAmount)
	require.Equal(t, float64(2000), eventData.Amount)
	require.Equal(t, []string{opportunitymodel.FieldMaskAmount, opportunitymodel.FieldMaskMaxAmount}, eventData.FieldsMask)
}

func TestContractEventHandler_UpdateRenewalArrForecast_ContractEndsBeforeNextRenewal(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()
	in5Minutes := utils.Now().Add(time.Minute * 5)
	in10Minutes := utils.Now().Add(time.Minute * 10)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.AnnuallyRenewalCycleString),
		EndedAt:          utils.TimePtr(in5Minutes),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringMedium),
			RenewedAt:         utils.TimePtr(in10Minutes),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(1000),
		Quantity:  int64(1),
		Billed:    string(servicelineitemmodel.AnnuallyBilledString),
		CreatedAt: utils.Now(),
	})

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityArr(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check create renewal opportunity command was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 2, len(eventList))

	updateNextCycleDateEvent := eventList[1]
	require.Equal(t, opportunityevent.OpportunityUpdateV1, updateNextCycleDateEvent.EventType)
	var eventData opportunityevent.OpportunityUpdateEvent
	err = updateNextCycleDateEvent.GetJsonData(&eventData)
	require.Nil(t, err)
	require.Equal(t, tenantName, eventData.Tenant)
	require.Equal(t, float64(0), eventData.MaxAmount)
	require.Equal(t, float64(0), eventData.Amount)
	require.Equal(t, []string{opportunitymodel.FieldMaskAmount, opportunitymodel.FieldMaskMaxAmount}, eventData.FieldsMask)
}

func TestContractEventHandler_UpdateRenewalArrForecast_ContractEndsIn6Months_ProrateAmount(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()
	in6Months := utils.Now().AddDate(0, 6, 0)
	in1Month := utils.Now().AddDate(0, 1, 0)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
		EndedAt:          utils.TimePtr(in6Months),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringHigh),
			RenewedAt:         utils.TimePtr(in1Month),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(1000),
		Quantity:  int64(1),
		Billed:    string(servicelineitemmodel.AnnuallyBilledString),
		CreatedAt: utils.Now(),
	})

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityArr(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check create renewal opportunity command was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 2, len(eventList))

	updateNextCycleDateEvent := eventList[1]
	require.Equal(t, opportunityevent.OpportunityUpdateV1, updateNextCycleDateEvent.EventType)
	var eventData opportunityevent.OpportunityUpdateEvent
	err = updateNextCycleDateEvent.GetJsonData(&eventData)
	require.Nil(t, err)
	require.Equal(t, tenantName, eventData.Tenant)
	require.Equal(t, float64(500), eventData.MaxAmount)
	require.Equal(t, float64(500), eventData.Amount)
	require.Equal(t, []string{opportunitymodel.FieldMaskAmount, opportunitymodel.FieldMaskMaxAmount}, eventData.FieldsMask)
}

func TestContractEventHandler_UpdateRenewalArrForecast_ContractEndsInMoreThan12Months_FullAmount(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()
	in13Months := utils.Now().AddDate(0, 13, 0)
	in12Months := utils.Now().AddDate(1, 0, 0)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
		EndedAt:          utils.TimePtr(in13Months),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringHigh),
			RenewedAt:         utils.TimePtr(in12Months),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(1000),
		Quantity:  int64(1),
		Billed:    string(servicelineitemmodel.AnnuallyBilledString),
		CreatedAt: utils.Now(),
	})

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityArr(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check create renewal opportunity command was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 2, len(eventList))

	updateNextCycleDateEvent := eventList[1]
	require.Equal(t, opportunityevent.OpportunityUpdateV1, updateNextCycleDateEvent.EventType)
	var eventData opportunityevent.OpportunityUpdateEvent
	err = updateNextCycleDateEvent.GetJsonData(&eventData)
	require.Nil(t, err)
	require.Equal(t, tenantName, eventData.Tenant)
	require.Equal(t, float64(1000), eventData.MaxAmount)
	require.Equal(t, float64(1000), eventData.Amount)
	require.Equal(t, []string{opportunitymodel.FieldMaskAmount, opportunitymodel.FieldMaskMaxAmount}, eventData.FieldsMask)
}

func TestContractEventHandler_UpdateActiveRenewalOpportunityLikelihood_EndedContract_UpdateLikelihood(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	tomorrow := utils.Now().AddDate(0, 0, 1)
	afterTomorrow := utils.Now().AddDate(0, 0, 2)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		EndedAt:      &tomorrow,
		RenewalCycle: string(model.AnnuallyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringLow),
			RenewedAt:         &afterTomorrow,
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// prepare grpc mock
	calledEventsPlatformToUpdateRenewalOpportunity := false
	opportunityServiceRefreshCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunity: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Equal(t, opportunitypb.RenewalLikelihood_ZERO_RENEWAL, op.RenewalLikelihood)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD}, op.FieldsMask)
			calledEventsPlatformToUpdateRenewalOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityServiceRefreshCallbacks)

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityLikelihood(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunity)
}

func TestContractEventHandler_UpdateActiveRenewalOpportunityLikelihood_EndedContract_LikelihoodAlreadyZero(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	tomorrow := utils.Now().AddDate(0, 0, 1)
	afterTomorrow := utils.Now().AddDate(0, 0, 2)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		EndedAt:      &tomorrow,
		RenewalCycle: string(model.AnnuallyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringZero),
			RenewedAt:         &afterTomorrow,
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityLikelihood(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check create renewal opportunity command was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))

	updateNextCycleDateEvent := eventList[0]
	require.Equal(t, opportunityevent.OpportunityCreateRenewalV1, updateNextCycleDateEvent.EventType)
}

func TestContractEventHandler_UpdateActiveRenewalOpportunityLikelihood_ReinitiatedContract_UpdateLikelihood(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	afterTomorrow := utils.Now().AddDate(0, 0, 2)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		RenewalCycle: string(model.AnnuallyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringZero),
			RenewedAt:         &afterTomorrow,
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// prepare grpc mock
	calledEventsPlatformToUpdateRenewalOpportunity := false
	opportunityServiceRefreshCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunity: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Equal(t, opportunitypb.RenewalLikelihood_MEDIUM_RENEWAL, op.RenewalLikelihood)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD}, op.FieldsMask)
			calledEventsPlatformToUpdateRenewalOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityServiceRefreshCallbacks)

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityLikelihood(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunity)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))
}

func TestContractEventHandler_UpdateActiveRenewalOpportunityLikelihood_ReinitiatedContract_LikelihoodNotZero(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	afterTomorrow := utils.Now().AddDate(0, 0, 2)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		RenewalCycle: string(model.AnnuallyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringHigh),
			RenewedAt:         &afterTomorrow,
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	prepareRenewalOpportunity(t, tenantName, opportunityId, aggregateStore)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore), testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityLikelihood(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check create renewal opportunity command was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))

	updateNextCycleDateEvent := eventList[0]
	require.Equal(t, opportunityevent.OpportunityCreateRenewalV1, updateNextCycleDateEvent.EventType)
}
