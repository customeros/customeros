package servicet

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	masterplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMasterPlanService_CreateMasterPlan(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to get grpc connection")
	masterPlanClient := masterplanpb.NewMasterPlanGrpcServiceClient(grpcConnection)

	response, err := masterPlanClient.CreateMasterPlan(ctx, &masterplanpb.CreateMasterPlanGrpcRequest{
		Tenant: tenant,
		Name:   "New Master Plan",
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
			Source:    "source",
		},
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	masterPlanId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[masterPlanAggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, event.MasterPlanCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(aggregate.MasterPlanAggregateType)+"-"+tenant+"-"+masterPlanId, eventList[0].GetAggregateID())

	var eventData event.MasterPlanCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "New Master Plan", eventData.Name)
	test.AssertRecentTime(t, eventData.CreatedAt)
	require.Equal(t, "app", eventData.SourceFields.AppSource)
	require.Equal(t, "source", eventData.SourceFields.Source)
	require.Equal(t, "source", eventData.SourceFields.SourceOfTruth)
}

func TestMasterPlanService_CreateMasterPlanMilestone(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	masterPlanId := "master-plan-id"

	aggregateStore := eventstoret.NewTestAggregateStore()
	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenant, masterPlanId)
	createEvent, _ := event.NewMasterPlanCreateEvent(masterPlanAggregate, "", commonmodel.Source{}, utils.Now())
	masterPlanAggregate.UncommittedEvents = append(masterPlanAggregate.UncommittedEvents, createEvent)
	aggregateStore.Save(ctx, masterPlanAggregate)

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	masterPlanClient := masterplanpb.NewMasterPlanGrpcServiceClient(grpcConnection)

	response, err := masterPlanClient.CreateMasterPlanMilestone(ctx, &masterplanpb.CreateMasterPlanMilestoneGrpcRequest{
		Tenant:       tenant,
		MasterPlanId: masterPlanId,
		Name:         "New Milestone",
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
			Source:    "source",
		},
		DurationHours: 1,
		Order:         2,
		Items:         []string{"item1", "item2"},
		Optional:      true,
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	milestoneId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[masterPlanAggregate.ID]
	require.Equal(t, 2, len(eventList))
	require.Equal(t, event.MasterPlanCreateV1, eventList[0].GetEventType())
	require.Equal(t, event.MasterPlanMilestoneCreateV1, eventList[1].GetEventType())
	require.Equal(t, string(aggregate.MasterPlanAggregateType)+"-"+tenant+"-"+masterPlanId, eventList[1].GetAggregateID())

	var eventData event.MasterPlanMilestoneCreateEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, milestoneId, eventData.MilestoneId)
	require.Equal(t, "New Milestone", eventData.Name)
	test.AssertRecentTime(t, eventData.CreatedAt)
	require.Equal(t, "app", eventData.SourceFields.AppSource)
	require.Equal(t, "source", eventData.SourceFields.Source)
	require.Equal(t, "source", eventData.SourceFields.SourceOfTruth)
	require.Equal(t, int64(1), eventData.DurationHours)
	require.Equal(t, int64(2), eventData.Order)
	require.Equal(t, []string{"item1", "item2"}, eventData.Items)
	require.Equal(t, true, eventData.Optional)
}

func TestMasterPlanService_UpdateMasterPlan(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	masterPlanId := "master-plan-id"

	// prepare master plan aggregate
	aggregateStore := eventstoret.NewTestAggregateStore()
	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenant, masterPlanId)
	createEvent, _ := event.NewMasterPlanCreateEvent(masterPlanAggregate, "", commonmodel.Source{}, utils.Now())
	masterPlanAggregate.UncommittedEvents = append(masterPlanAggregate.UncommittedEvents, createEvent)
	aggregateStore.Save(ctx, masterPlanAggregate)

	// prepare connection to grpc server
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	masterPlanClient := masterplanpb.NewMasterPlanGrpcServiceClient(grpcConnection)

	// Execute the command
	response, err := masterPlanClient.UpdateMasterPlan(ctx, &masterplanpb.UpdateMasterPlanGrpcRequest{
		Tenant:         tenant,
		MasterPlanId:   masterPlanId,
		Name:           "New Plan Name",
		Retired:        true,
		AppSource:      "app",
		LoggedInUserId: "user-id",
		FieldsMask:     []masterplanpb.MasterPlanFieldMask{masterplanpb.MasterPlanFieldMask_MASTER_PLAN_PROPERTY_NAME, masterplanpb.MasterPlanFieldMask_MASTER_PLAN_PROPERTY_RETIRED},
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	// verify
	require.Equal(t, masterPlanId, response.Id)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[masterPlanAggregate.ID]
	require.Equal(t, 2, len(eventList))
	require.Equal(t, event.MasterPlanCreateV1, eventList[0].GetEventType())
	require.Equal(t, event.MasterPlanUpdateV1, eventList[1].GetEventType())
	require.Equal(t, string(aggregate.MasterPlanAggregateType)+"-"+tenant+"-"+masterPlanId, eventList[1].GetAggregateID())

	var eventData event.MasterPlanUpdateEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "New Plan Name", eventData.Name)
	require.Equal(t, true, eventData.Retired)
	test.AssertRecentTime(t, eventData.UpdatedAt)
	require.ElementsMatch(t, []string{event.FieldMaskName, event.FieldMaskRetired}, eventData.FieldsMask)
}

func TestMasterPlanService_UpdateMasterPlanMilestone(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	masterPlanId := "master-plan-id"
	milestoneId := "milestone-id"

	// prepare master plan aggregate
	aggregateStore := eventstoret.NewTestAggregateStore()
	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenant, masterPlanId)
	createEvent, _ := event.NewMasterPlanCreateEvent(masterPlanAggregate, "", commonmodel.Source{}, utils.Now())
	masterPlanAggregate.UncommittedEvents = append(masterPlanAggregate.UncommittedEvents, createEvent)
	aggregateStore.Save(ctx, masterPlanAggregate)

	// prepare connection to grpc server
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	masterPlanClient := masterplanpb.NewMasterPlanGrpcServiceClient(grpcConnection)

	// execute the command
	response, err := masterPlanClient.UpdateMasterPlanMilestone(ctx, &masterplanpb.UpdateMasterPlanMilestoneGrpcRequest{
		Tenant:                tenant,
		MasterPlanId:          masterPlanId,
		MasterPlanMilestoneId: milestoneId,
		Name:                  "Updated Milestone",
		DurationHours:         1,
		Order:                 2,
		Items:                 []string{"item1", "item2"},
		Optional:              true,
		Retired:               true,
		AppSource:             "app",
		FieldsMask: []masterplanpb.MasterPlanMilestoneFieldMask{masterplanpb.MasterPlanMilestoneFieldMask_MASTER_PLAN_MILESTONE_PROPERTY_NAME,
			masterplanpb.MasterPlanMilestoneFieldMask_MASTER_PLAN_MILESTONE_PROPERTY_ITEMS},
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	// verify
	require.Equal(t, milestoneId, response.Id)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[masterPlanAggregate.ID]
	require.Equal(t, 2, len(eventList))
	require.Equal(t, event.MasterPlanCreateV1, eventList[0].GetEventType())
	require.Equal(t, event.MasterPlanMilestoneUpdateV1, eventList[1].GetEventType())
	require.Equal(t, string(aggregate.MasterPlanAggregateType)+"-"+tenant+"-"+masterPlanId, eventList[1].GetAggregateID())

	var eventData event.MasterPlanMilestoneUpdateEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, milestoneId, eventData.MilestoneId)
	require.Equal(t, "Updated Milestone", eventData.Name)
	require.Equal(t, int64(0), eventData.DurationHours)
	require.Equal(t, int64(0), eventData.Order)
	require.Equal(t, []string{"item1", "item2"}, eventData.Items)
	require.Equal(t, false, eventData.Optional)
	require.Equal(t, false, eventData.Retired)
	test.AssertRecentTime(t, eventData.UpdatedAt)
	require.ElementsMatch(t, []string{event.FieldMaskName, event.FieldMaskItems}, eventData.FieldsMask)
}
