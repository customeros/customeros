package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContractEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	userIdCreator := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{})
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, "sf")
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "User": 1, "ExternalSystem": 1, "Contract": 0})

	// Prepare the event handler
	contractEventHandler := &ContractEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create a ContractCreateEvent
	contractId := uuid.New().String()
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	timeNow := utils.Now()
	createEvent, err := event.NewContractCreateEvent(
		contractAggregate,
		model.ContractDataFields{
			Name:             "New Contract",
			ContractUrl:      "http://contract.url",
			OrganizationId:   orgId,
			CreatedByUserId:  userIdCreator,
			ServiceStartedAt: &timeNow,
			SignedAt:         &timeNow,
			RenewalCycle:     model.MonthlyRenewal,
			Status:           model.Live,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		commonmodel.ExternalSystem{
			ExternalSystemId: "sf",
			ExternalId:       "ext-id-1",
		},
		timeNow,
		timeNow,
	)
	require.Nil(t, err, "failed to create contract create event")

	// Execute
	err = contractEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute contract create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1,
		"User":         1,
		"Contract":     1, "Contract_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "CREATED_BY", userIdCreator)
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, orgId, "HAS_CONTRACT", contractId)
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "IS_LINKED_WITH", "sf")

	contractDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// Verify contract
	contract := graph_db.MapDbNodeToContractEntity(*contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, "New Contract", contract.Name)
	require.Equal(t, "http://contract.url", contract.ContractUrl)
	require.Equal(t, model.Live.String(), contract.Status)
	require.Equal(t, model.MonthlyRenewal.String(), contract.RenewalCycle)
	require.True(t, timeNow.Equal(contract.CreatedAt.UTC()))
	require.True(t, timeNow.Equal(contract.UpdatedAt.UTC()))
	require.True(t, timeNow.Equal(*contract.ServiceStartedAt))
	require.True(t, timeNow.Equal(*contract.SignedAt))
	require.Nil(t, contract.EndedAt)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), contract.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, contract.AppSource)

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

func TestContractEventHandler_OnUpdate_FrequencySet(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Name:        "test contract",
		ContractUrl: "http://contract.url",
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1})

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}
	now := utils.Now()
	yesterday := now.AddDate(0, 0, -1)
	daysAgo2 := now.AddDate(0, 0, -2)
	tomorrow := now.AddDate(0, 0, 1)
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	updateEvent, err := event.NewContractUpdateEvent(contractAggregate,
		model.ContractDataFields{
			Name:             "test contract updated",
			ContractUrl:      "http://contract.url/updated",
			ServiceStartedAt: &yesterday,
			SignedAt:         &daysAgo2,
			EndedAt:          &tomorrow,
			RenewalCycle:     model.MonthlyRenewal,
			Status:           model.Live,
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1, "Contract_" + tenantName: 1})

	contractDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := graph_db.MapDbNodeToContractEntity(*contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, "test contract updated", contract.Name)
	require.Equal(t, "http://contract.url/updated", contract.ContractUrl)
	require.Equal(t, model.Live.String(), contract.Status)
	require.Equal(t, model.MonthlyRenewal.String(), contract.RenewalCycle)
	require.True(t, now.Equal(contract.UpdatedAt))
	require.True(t, yesterday.Equal(*contract.ServiceStartedAt))
	require.True(t, daysAgo2.Equal(*contract.SignedAt))
	require.True(t, tomorrow.Equal(*contract.EndedAt))
	require.Equal(t, entity.DataSource(constants.SourceOpenline), contract.SourceOfTruth)

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

func TestContractEventHandler_OnUpdate_FrequencyNotChanged(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		RenewalCycle: string(model.MonthlyRenewalCycleString),
	})

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	updateEvent, err := event.NewContractUpdateEvent(contractAggregate,
		model.ContractDataFields{
			Name: "test contract updated",
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		utils.Now())
	require.Nil(t, err)

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	// Check no create renewal opportunity command was generated
	require.Equal(t, 0, len(aggregateStore.GetEventMap()))
}

func TestContractEventHandler_OnUpdate_CurrentSourceOpenline_UpdateSourceNonOpenline(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	now := utils.Now()
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Name:             "test contract",
		ContractUrl:      "http://contract.url",
		Status:           "DRAFT",
		RenewalCycle:     "ANNUALLY",
		ServiceStartedAt: &now,
		SignedAt:         &now,
		EndedAt:          &now,
		SourceOfTruth:    constants.SourceOpenline,
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1})

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}
	yesterday := now.AddDate(0, 0, -1)
	daysAgo2 := now.AddDate(0, 0, -2)
	tomorrow := now.AddDate(0, 0, 1)
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	updateEvent, err := event.NewContractUpdateEvent(contractAggregate,
		model.ContractDataFields{
			Name:             "test contract updated",
			ContractUrl:      "http://contract.url/updated",
			ServiceStartedAt: &yesterday,
			SignedAt:         &daysAgo2,
			EndedAt:          &tomorrow,
			RenewalCycle:     model.MonthlyRenewal,
			Status:           model.Live,
		},
		commonmodel.ExternalSystem{},
		"hubspot",
		now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1, "Contract_" + tenantName: 1})

	contractDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := graph_db.MapDbNodeToContractEntity(*contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, "test contract", contract.Name)
	require.Equal(t, "http://contract.url", contract.ContractUrl)
	require.Equal(t, model.Draft.String(), contract.Status)
	require.Equal(t, model.AnnuallyRenewal.String(), contract.RenewalCycle)
	require.True(t, now.Equal(contract.UpdatedAt))
	require.True(t, now.Equal(*contract.ServiceStartedAt))
	require.True(t, now.Equal(*contract.SignedAt))
	require.True(t, now.Equal(*contract.EndedAt))
	require.Equal(t, entity.DataSource(constants.SourceOpenline), contract.SourceOfTruth)
}
