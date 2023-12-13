package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	contractmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	opportunitymodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServiceLineItemEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Id: contractId,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.MonthlyBilled,
			Quantity:   10,
			Price:      100.50,
			Name:       "Test service line item",
			ContractId: contractId,
			ParentId:   serviceLineItemId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreate(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentId)
	require.Equal(t, model.MonthlyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, float64(100.50), serviceLineItem.Price)
	require.Equal(t, "Test service line item", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, timeNow, serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)
}

func TestServiceLineItemEventHandler_OnUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Prepare test data in Neo4j
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{})
	serviceLineItemId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Billed: model.MonthlyBilled.String(),
	})
	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Name:     "Updated Service Line Item",
			Price:    200.0,
			Quantity: 20,
			Billed:   model.AnnuallyBilled,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		updatedAt,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdate(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
	})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, model.AnnuallyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, int64(20), serviceLineItem.Quantity)
	require.Equal(t, float64(200.0), serviceLineItem.Price)
	require.Equal(t, "Updated Service Line Item", serviceLineItem.Name)
}

func TestServiceLineItemEventHandler_OnDeleteUnnamed(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Prepare test data in Neo4j
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		RenewalCycle: string(contractmodel.AnnuallyRenewalCycleString),
	})
	serviceLineItemId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Billed: model.MonthlyBilled.String(),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemDeleteEvent
	deleteEvent, err := event.NewServiceLineItemDeleteEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
	)
	require.Nil(t, err)

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = deleteEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnDelete(ctx, deleteEvent)
	require.Nil(t, err)

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 0, "ServiceLineItem_" + tenantName: 0,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// check event was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))
	require.Equal(t, opportunityevent.OpportunityUpdateV1, eventList[0].GetEventType())

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemRemoved, action.Type)
	require.Equal(t, "logged-in user removed Unnamed service from Unnamed contract", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Unnamed service","comment":"service line item removed is Unnamed service from Unnamed contract by logged-in user"}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnDelete(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Prepare test data in Neo4j
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Name:         "Contract 1",
		RenewalCycle: string(contractmodel.AnnuallyRenewalCycleString),
	})
	serviceLineItemId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Name:   "Service 1",
		Billed: model.MonthlyBilled.String(),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemDeleteEvent
	deleteEvent, err := event.NewServiceLineItemDeleteEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
	)
	require.Nil(t, err)

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = deleteEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnDelete(ctx, deleteEvent)
	require.Nil(t, err)

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 0, "ServiceLineItem_" + tenantName: 0,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// check event was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))
	require.Equal(t, opportunityevent.OpportunityUpdateV1, eventList[0].GetEventType())

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemRemoved, action.Type)
	require.Equal(t, "logged-in user removed Service 1 from Contract 1", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","comment":"service line item removed is Service 1 from Contract 1 by logged-in user"}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnClose(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Prepare test data in Neo4j
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		RenewalCycle: string(contractmodel.AnnuallyRenewalCycleString),
	})
	serviceLineItemId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Billed: model.MonthlyBilled.String(),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	now := utils.Now()
	// Create a ServiceLineItemCloseEvent
	closeEvent, err := event.NewServiceLineItemCloseEvent(aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId), now, now, true)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnClose(ctx, closeEvent)
	require.Nil(t, err)

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
	})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, now, serviceLineItem.UpdatedAt)
	require.Equal(t, now, *serviceLineItem.EndedAt)
	require.True(t, serviceLineItem.IsCanceled)

	// check event was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))
	require.Equal(t, opportunityevent.OpportunityUpdateV1, eventList[0].GetEventType())
}

func TestServiceLineItemEventHandler_OnUpdatePriceIncreaseRetroactively_TimelineEvent(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Prepare test data in Neo4j
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{})
	serviceLineItemId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Name:   "Service 1",
		Billed: model.MonthlyBilled.String(),
		Price:  150.0,
	})
	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Price:    200.0,
			Comments: "this is the reason for change",
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		updatedAt,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdate(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, model.MonthlyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, float64(200.0), serviceLineItem.Price)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user retroactively increased the price for Service 1 from 150.00/month to 200.00/month", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","quantity":0,"price":200,"previousPrice":150,"billedType":"MONTHLY","comment":"price changed is 150.00 for service Service 1","reasonForChange":"this is the reason for change"}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdatePriceIncreasePerUseRetroactively_TimelineEvent(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Prepare test data in Neo4j
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{})
	serviceLineItemId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Name:   "Service 1",
		Billed: model.UsageBilled.String(),
		Price:  150.0,
	})
	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Price:    200.0,
			Billed:   model.UsageBilled,
			Comments: "test reason for change",
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		updatedAt,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdate(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, model.UsageBilled.String(), serviceLineItem.Billed)
	require.Equal(t, float64(200.0), serviceLineItem.Price)
	require.Equal(t, "test reason for change", serviceLineItem.Comments)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user retroactively increased the price for Service 1 from 150.0000 to 200.0000", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","quantity":0,"price":200,"previousPrice":150,"billedType":"USAGE","comment":"price changed is 150.00 for service Service 1","reasonForChange":"test reason for change"}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdatePriceDecreaseRetroactively_TimelineEvent(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Prepare test data in Neo4j
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{})
	serviceLineItemId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Name:   "Service 1",
		Billed: model.AnnuallyBilled.String(),
		Price:  150.0,
	})
	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Name:     "Service 1",
			Price:    50.0,
			Billed:   model.AnnuallyBilled,
			Comments: "Reason for change is x",
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		updatedAt,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdate(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, model.AnnuallyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, float64(50.0), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user retroactively decreased the price for Service 1 from 150.00/year to 50.00/year", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","quantity":0,"price":50,"previousPrice":150,"billedType":"ANNUALLY","comment":"price changed is 150.00 for service Service 1","reasonForChange":"Reason for change is x"}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdatePriceDecreaseOnceRetroactively_TimelineEvent(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Prepare test data in Neo4j
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{})
	serviceLineItemId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Name:   "Service 1",
		Billed: model.OnceBilled.String(),
		Price:  150.0,
	})
	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Name:   "Service 1",
			Price:  50.0,
			Billed: model.OnceBilled,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		updatedAt,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdate(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, model.OnceBilled.String(), serviceLineItem.Billed)
	require.Equal(t, float64(50.0), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user retroactively decreased the price for Service 1 from 150.00 to 50.00", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","quantity":0,"price":50,"previousPrice":150,"billedType":"ONCE","comment":"price changed is 150.00 for service Service 1","reasonForChange":""}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdateQuantityIncreaseRetroactively_TimelineEvent(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Prepare test data in Neo4j
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{})
	serviceLineItemId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Name:     "Service 1",
		Quantity: 15,
	})
	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Quantity: 20,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		updatedAt,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdate(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, int64(20), serviceLineItem.Quantity)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemQuantityUpdated, action.Type)
	require.Equal(t, "logged-in user retroactively increased the quantity of Service 1 from 15 to 20", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","quantity":20,"previousQuantity":15,"price":0,"billedType":"","comment":"quantity changed is 15 for service Service 1","reasonForChange":""}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdateQuantityDecreaseRetroactively_TimelineEvent(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Prepare test data in Neo4j
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{})
	serviceLineItemId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Name:     "Service 1",
		Quantity: 400,
	})
	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Quantity: 350,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		updatedAt,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)
	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdate(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, model.MonthlyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, int64(350), serviceLineItem.Quantity)

	// verify actionat
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemQuantityUpdated, action.Type)
	require.Equal(t, "logged-in user retroactively decreased the quantity of Service 1 from 400 to 350", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","quantity":350,"previousQuantity":400,"price":0,"billedType":"","comment":"quantity changed is 400 for service Service 1","reasonForChange":""}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdateBilledType_TimelineEvent(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Prepare test data in Neo4j
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{})
	serviceLineItemId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Name:   "Service 1",
		Price:  20,
		Billed: model.AnnuallyBilled.String(),
	})
	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Name:   "Service 1",
			Price:  20,
			Billed: model.MonthlyBilled,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		updatedAt,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)
	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdate(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, model.MonthlyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, "Service 1", serviceLineItem.Name)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemBilledTypeUpdated, action.Type)
	require.Equal(t, "logged-in user changed the billing cycle for Service 1 from 20.00/year to 20.00/month", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":20,"quantity":0,"billedType":"MONTHLY","previousBilledType":"ANNUALLY","comment":"billed type changed is ANNUALLY for service Service 1","reasonForChange":""}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateRecurringMonthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.MonthlyBilled,
			Quantity:   10,
			Price:      170.25,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreate(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentId)
	require.Equal(t, model.MonthlyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, float64(170.25), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, timeNow, serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemBilledTypeRecurringCreated, action.Type)
	require.Equal(t, "logged-in user added a recurring service to Contract 1: Service 1 at 10 x 170.25/month", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":170.25,"quantity":10,"billedType":"MONTHLY","previousBilledType":"","comment":"billed type is MONTHLY for service Service 1","reasonForChange":""}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateRecurringAnnually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.AnnuallyBilled,
			Quantity:   10,
			Price:      170.25,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreate(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentId)
	require.Equal(t, model.AnnuallyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, float64(170.25), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, timeNow, serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemBilledTypeRecurringCreated, action.Type)
	require.Equal(t, "logged-in user added a recurring service to Contract 1: Service 1 at 10 x 170.25/year", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":170.25,"quantity":10,"billedType":"ANNUALLY","previousBilledType":"","comment":"billed type is ANNUALLY for service Service 1","reasonForChange":""}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateRecurringQuarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.QuarterlyBilled,
			Quantity:   10,
			Price:      170.25,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreate(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentId)
	require.Equal(t, model.QuarterlyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, float64(170.25), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, timeNow, serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemBilledTypeRecurringCreated, action.Type)
	require.Equal(t, "logged-in user added a recurring service to Contract 1: Service 1 at 10 x 170.25/quarter", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":170.25,"quantity":10,"billedType":"QUARTERLY","previousBilledType":"","comment":"billed type is QUARTERLY for service Service 1","reasonForChange":""}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateOnce(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.OnceBilled,
			Price:      170.25,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreate(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentId)
	require.Equal(t, model.OnceBilled.String(), serviceLineItem.Billed)
	require.Equal(t, float64(170.25), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, timeNow, serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemBilledTypeOnceCreated, action.Type)
	require.Equal(t, "logged-in user added an one time service to Contract 1: Service 1 at 170.25", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":170.25,"quantity":0,"billedType":"ONCE","previousBilledType":"","comment":"billed type is ONCE for service Service 1","reasonForChange":""}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreatePerUse(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.UsageBilled,
			Price:      170.25,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreate(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentId)
	require.Equal(t, model.UsageBilled.String(), serviceLineItem.Billed)
	require.Equal(t, float64(170.25), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, timeNow, serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemBilledTypeUsageCreated, action.Type)
	require.Equal(t, "logged-in user added a per use service to Contract 1: Service 1 at 170.2500", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":170.25,"quantity":0,"billedType":"USAGE","previousBilledType":"","comment":"billed type is USAGE for service Service 1","reasonForChange":""}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateNewVersionForNonRetroactiveQuantityDecrease(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	serviceLineItemParentId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Name:     "Service Parent",
		Quantity: 400,
	})
	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}
	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.MonthlyBilled,
			Quantity:   10,
			Name:       "Test service line item",
			ContractId: contractId,
			ParentId:   serviceLineItemParentId,
			Comments:   "reason for what change?",
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreate(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 2, "ServiceLineItem_" + tenantName: 2,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, serviceLineItemParentId, serviceLineItem.ParentId)
	require.Equal(t, model.MonthlyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, "Test service line item", serviceLineItem.Name)
	require.Equal(t, "reason for what change?", serviceLineItem.Comments)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, timeNow, serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemQuantityUpdated, action.Type)
	require.Equal(t, "logged-in user decreased the quantity of Test service line item from 400 to 10", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Test service line item","quantity":10,"previousQuantity":400,"price":0,"billedType":"MONTHLY","comment":"quantity is 10 for service Test service line item","reasonForChange":"reason for what change?"}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateNewVersionForNonRetroactivePriceIncrease(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	serviceLineItemParentId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Name:  "Service Parent",
		Price: 1500.56,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.MonthlyBilled,
			Price:      850.75,
			Name:       "Test service line item",
			ContractId: contractId,
			ParentId:   serviceLineItemParentId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreate(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 2, "ServiceLineItem_" + tenantName: 2,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, serviceLineItemParentId, serviceLineItem.ParentId)
	require.Equal(t, model.MonthlyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, float64(850.75), serviceLineItem.Price)
	require.Equal(t, "Test service line item", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, timeNow, serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user decreased the price for Test service line item from 1500.56/month to 850.75/month", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Test service line item","quantity":0,"price":850.75,"previousPrice":1500.56,"billedType":"MONTHLY","comment":"price is 850.75 for service Test service line item","reasonForChange":""}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateNewVersionForNonRetroactivePriceIncreaseNoNameService(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	serviceLineItemParentId := neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price: 1500.56,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.MonthlyBilled,
			Price:      850.75,
			ContractId: contractId,
			ParentId:   serviceLineItemParentId,
			Comments:   "This is a reason for change",
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreate(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 2, "ServiceLineItem_" + tenantName: 2,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, serviceLineItemParentId, serviceLineItem.ParentId)
	require.Equal(t, model.MonthlyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, float64(850.75), serviceLineItem.Price)
	require.Equal(t, "This is a reason for change", serviceLineItem.Comments)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, timeNow, serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user decreased the price for Unnamed service from 1500.56/month to 850.75/month", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Unnamed service","quantity":0,"price":850.75,"previousPrice":1500.56,"billedType":"MONTHLY","comment":"price is 850.75 for service Unnamed service","reasonForChange":"This is a reason for change"}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdateBilledTypeNonRetroactiveForExistingSLI(t *testing.T) {
	//in this test we will test the following scenario:
	// 1. create a service line item with billed type annually
	// 2. create a service line item with billed type quarterly & use as parent id the id of the service line item created in step 1
	// 3. Figure out the action message
	// 4. This is a new version of an existing service line item where only the billed type is changed(for retroactive changes on billed type we have a different test)
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Setup test environment
	serviceLineItemId1 := "service-line-item-id-1"
	serviceLineItemId2 := "service-line-item-id-2"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Id:       serviceLineItemId1,
		Name:     "service 1",
		Billed:   model.AnnuallyBilled.String(),
		Price:    170.25,
		Quantity: 10,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId2)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.QuarterlyBilled,
			Quantity:   10,
			Price:      170.25,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId1,
			Comments:   "Reason for change",
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreate(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 2, "ServiceLineItem_" + tenantName: 2,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId2)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId2)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId2, serviceLineItem.Id)
	require.Equal(t, serviceLineItemId1, serviceLineItem.ParentId)
	require.Equal(t, model.QuarterlyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, float64(170.25), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, "Reason for change", serviceLineItem.Comments)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, timeNow, serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemBilledTypeUpdated, action.Type)
	require.Equal(t, "logged-in user changed the billing cycle for Service 1 from 170.25/year to 170.25/quarter", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":170.25,"quantity":10,"billedType":"QUARTERLY","previousBilledType":"ANNUALLY","comment":"billed type is QUARTERLY for service Service 1","reasonForChange":"Reason for change"}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdatePriceAndBilledTypeNonRetroactiveForExistingSLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Setup test environment
	serviceLineItemId1 := "service-line-item-id-1"
	serviceLineItemId2 := "service-line-item-id-2"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Id:       serviceLineItemId1,
		Name:     "service 1",
		Billed:   model.AnnuallyBilled.String(),
		Price:    170.25,
		Quantity: 10,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId2)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.MonthlyBilled,
			Quantity:   10,
			Price:      100,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId1,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreate(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 2, "ServiceLineItem_" + tenantName: 2,
		"TimelineEvent": 2, "TimelineEvent_" + tenantName: 2,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId2)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId2)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId2, serviceLineItem.Id)
	require.Equal(t, serviceLineItemId1, serviceLineItem.ParentId)
	require.Equal(t, model.MonthlyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, float64(100), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, timeNow, serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, entity.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user decreased the price for Service 1 from 170.25/year to 100.00/month", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","quantity":10,"price":100,"previousPrice":170.25,"billedType":"MONTHLY","comment":"price is 100.00 for service Service 1","reasonForChange":""}`, action.Metadata)
}
