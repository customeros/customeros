package graph

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	contractmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	organizationcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	organizationEvents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestOpportunityEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	userIdOwner := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{})
	userIdCreator := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{})
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, "sf")
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "User": 2, "ExternalSystem": 1, "Opportunity": 0})

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create an OpportunityCreateEvent
	opportunityId := uuid.New().String()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	timeNow := utils.Now()
	opportunityData := model.OpportunityDataFields{
		Name:              "New Opportunity",
		Amount:            10000,
		InternalType:      model.NBO,
		ExternalType:      "TypeA",
		InternalStage:     model.OPEN,
		ExternalStage:     "Stage1",
		EstimatedClosedAt: &timeNow,
		OwnerUserId:       userIdOwner,
		CreatedByUserId:   userIdCreator,
		GeneralNotes:      "Some general notes about the opportunity",
		NextSteps:         "Next steps to proceed with the opportunity",
		OrganizationId:    orgId,
	}
	createEvent, err := event.NewOpportunityCreateEvent(
		opportunityAggregate,
		opportunityData,
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
	require.Nil(t, err, "failed to create opportunity create event")

	// EXECUTE
	err = opportunityEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute opportunity create event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization":   1,
		"User":           2,
		"ExternalSystem": 1,
		"Opportunity":    1, "Opportunity_" + tenantName: 1})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, userIdOwner, "OWNS", opportunityId)
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, opportunityId, "CREATED_BY", userIdCreator)
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, orgId, "HAS_OPPORTUNITY", opportunityId)
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, opportunityId, "IS_LINKED_WITH", "sf")

	opportunityDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Opportunity_"+tenantName, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(*opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), opportunity.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, opportunity.AppSource)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), opportunity.SourceOfTruth)
	require.Equal(t, timeNow, opportunity.CreatedAt)
	test.AssertRecentTime(t, opportunity.UpdatedAt)
	require.Equal(t, timeNow, *opportunity.EstimatedClosedAt)
	require.Equal(t, opportunityData.Name, opportunity.Name)
	require.Equal(t, opportunityData.Amount, opportunity.Amount)
	require.Equal(t, string(opportunityData.InternalType.StringValue()), opportunity.InternalType)
	require.Equal(t, opportunityData.ExternalType, opportunity.ExternalType)
	require.Equal(t, string(opportunityData.InternalStage.StringValue()), opportunity.InternalStage)
	require.Equal(t, opportunityData.ExternalStage, opportunity.ExternalStage)
	require.Equal(t, opportunityData.GeneralNotes, opportunity.GeneralNotes)
	require.Equal(t, opportunityData.NextSteps, opportunity.NextSteps)
}

func TestOpportunityEventHandler_OnCreateRenewal(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract": 1, "Opportunity": 0})

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	// Create an OpportunityCreateEvent
	opportunityId := uuid.New().String()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	timeNow := utils.Now()
	createEvent, err := event.NewOpportunityCreateRenewalEvent(
		opportunityAggregate,
		contractId,
		string(model.RenewalLikelihoodStringLow),
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
	)
	require.Nil(t, err, "failed to create opportunity create renewal event")

	// EXECUTE
	err = opportunityEventHandler.OnCreateRenewal(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute opportunity create event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":    1,
		"Opportunity": 1, "Opportunity_" + tenantName: 1})
	neo4jt.AssertRelationships(ctx, t, testDatabase.Driver, contractId, []string{"HAS_OPPORTUNITY", "ACTIVE_RENEWAL"}, opportunityId)

	opportunityDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Opportunity_"+tenantName, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(*opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), opportunity.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, opportunity.AppSource)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), opportunity.SourceOfTruth)
	require.Equal(t, timeNow, opportunity.CreatedAt)
	test.AssertRecentTime(t, opportunity.UpdatedAt)
	require.Nil(t, opportunity.EstimatedClosedAt)
	require.Equal(t, "", opportunity.Name)
	require.Equal(t, "", opportunity.ExternalType)
	require.Equal(t, "", opportunity.ExternalStage)
	require.Equal(t, string(model.OpportunityInternalTypeStringRenewal), opportunity.InternalType)
	require.Equal(t, string(model.OpportunityInternalStageStringOpen), opportunity.InternalStage)
	require.Equal(t, string(model.RenewalLikelihoodStringLow), opportunity.RenewalDetails.RenewalLikelihood)
}

func TestOpportunityEventHandler_OnUpdateNextCycleDate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalStage: string(model.OpportunityInternalStageStringOpen),
	})
	updatedAt := utils.Now()
	renewedAt := updatedAt.AddDate(0, 6, 0) // 6 months later

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an OpportunityUpdateNextCycleDateEvent
	updateEvent, err := event.NewOpportunityUpdateNextCycleDateEvent(
		aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId),
		updatedAt,
		&renewedAt,
	)
	require.Nil(t, err)

	// Execute the event handler
	err = opportunityEventHandler.OnUpdateNextCycleDate(ctx, updateEvent)
	require.Nil(t, err)

	// Assert Neo4j Node
	opportunityDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, constants.NodeLabel_Opportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// Validate that the opportunity next cycle date is updated in the repository
	opportunity := graph_db.MapDbNodeToOpportunityEntity(*opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, updatedAt, opportunity.UpdatedAt)
	require.Equal(t, renewedAt, *opportunity.RenewalDetails.RenewedAt)
}

func TestOpportunityEventHandler_OnUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Name:      "test opportunity",
		Amount:    10000,
		MaxAmount: 20000,
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{constants.NodeLabel_Opportunity: 1})

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	now := utils.Now()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	updateEvent, err := event.NewOpportunityUpdateEvent(opportunityAggregate,
		model.OpportunityDataFields{
			Name:      "updated opportunity",
			Amount:    30000,
			MaxAmount: 40000,
		},
		constants.SourceOpenline,
		commonmodel.ExternalSystem{},
		now,
		nil)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = opportunityEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{constants.NodeLabel_Opportunity: 1, constants.NodeLabel_Opportunity + "_" + tenantName: 1})

	opportunityDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, constants.NodeLabel_Opportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(*opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, now, opportunity.UpdatedAt)
	require.Equal(t, "updated opportunity", opportunity.Name)
	require.Equal(t, float64(30000), opportunity.Amount)
	require.Equal(t, float64(40000), opportunity.MaxAmount)
}

func TestOpportunityEventHandler_OnUpdate_OnlyAmountIsChangedByFieldsMask(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Name:   "test opportunity",
		Amount: 10000,
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{constants.NodeLabel_Opportunity: 1})

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	now := utils.Now()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	updateEvent, err := event.NewOpportunityUpdateEvent(opportunityAggregate,
		model.OpportunityDataFields{
			Name:   "updated opportunity",
			Amount: 20000,
		},
		constants.SourceOpenline,
		commonmodel.ExternalSystem{},
		now,
		[]string{"amount"})
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = opportunityEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{constants.NodeLabel_Opportunity: 1, constants.NodeLabel_Opportunity + "_" + tenantName: 1})

	opportunityDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Opportunity_"+tenantName, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(*opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, now, opportunity.UpdatedAt)
	require.Equal(t, "test opportunity", opportunity.Name)
	require.Equal(t, float64(20000), opportunity.Amount)
}

func TestOpportunityEventHandler_OnUpdateRenewal_AmountAndRenewalChangedByUser(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Name:         "test opportunity",
		Amount:       10000,
		InternalType: string(model.OpportunityInternalTypeStringRenewal),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: "HIGH",
		},
	})
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	contractId := neo4jt.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, entity.ContractEntity{
		RenewalCycle: string(contractmodel.MonthlyRenewalCycleString),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{constants.NodeLabel_Opportunity: 1})

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:                  testLogger,
		repositories:         testDatabase.Repositories,
		opportunityCommands:  opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
		organizationCommands: organizationcmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore, nil),
	}

	now := utils.Now()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	updateEvent, err := event.NewOpportunityUpdateRenewalEvent(opportunityAggregate,
		"MEDIUM",
		"some comments",
		"user-123",
		"openline",
		float64(10),
		now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = opportunityEventHandler.OnUpdateRenewal(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{constants.NodeLabel_Opportunity: 1, constants.NodeLabel_Opportunity + "_" + tenantName: 1})

	opportunityDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, constants.NodeLabel_Opportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(*opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, now, opportunity.UpdatedAt)
	require.Equal(t, "MEDIUM", opportunity.RenewalDetails.RenewalLikelihood)
	require.Equal(t, "user-123", opportunity.RenewalDetails.RenewalUpdatedByUserId)
	require.Equal(t, now, *opportunity.RenewalDetails.RenewalUpdatedByUserAt)
	require.Equal(t, float64(10), opportunity.Amount)
	require.Equal(t, "some comments", opportunity.Comments)

	// Check event to refresh organization arr is generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))

	generatedEvent1 := eventList[0]
	require.Equal(t, organizationEvents.OrganizationRefreshArrV1, generatedEvent1.EventType)
	var eventData1 organizationEvents.OrganizationRefreshArrEvent
	err = generatedEvent1.GetJsonData(&eventData1)
	require.Nil(t, err)
	require.Equal(t, tenantName, eventData1.Tenant)
}

func TestOpportunityEventHandler_OnUpdateRenewal_OnlyCommentsChangedByUser_DoNotUpdatePreviousUpdatedByUser(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Name:         "test opportunity",
		Amount:       10000,
		Comments:     "no comments",
		InternalType: string(model.OpportunityInternalTypeStringRenewal),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood:      "HIGH",
			RenewalUpdatedByUserId: "orig-user",
		},
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{constants.NodeLabel_Opportunity: 1})

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	now := utils.Now()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	updateEvent, err := event.NewOpportunityUpdateRenewalEvent(opportunityAggregate,
		"HIGH",
		"some comments",
		"user-123",
		"openline",
		float64(10000),
		now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = opportunityEventHandler.OnUpdateRenewal(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{constants.NodeLabel_Opportunity: 1, constants.NodeLabel_Opportunity + "_" + tenantName: 1})

	opportunityDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, constants.NodeLabel_Opportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(*opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, now, opportunity.UpdatedAt)
	require.Equal(t, "HIGH", opportunity.RenewalDetails.RenewalLikelihood)
	require.Equal(t, "orig-user", opportunity.RenewalDetails.RenewalUpdatedByUserId)
	require.Nil(t, opportunity.RenewalDetails.RenewalUpdatedByUserAt)
	require.Equal(t, float64(10000), opportunity.Amount)
	require.Equal(t, "some comments", opportunity.Comments)

	// Check no events were generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 0, len(eventsMap))
}

func TestOpportunityEventHandler_OnUpdateRenewal_LikelihoodChangedByUser_GenerateEventToRecalculateArr(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		RenewalCycle: string(contractmodel.MonthlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Amount:        10000,
		InternalType:  string(model.OpportunityInternalTypeStringRenewal),
		InternalStage: string(model.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood:      "HIGH",
			RenewalUpdatedByUserId: "orig-user",
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		constants.NodeLabel_Opportunity: 1,
		constants.NodeLabel_Contract:    1})

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:                 testLogger,
		repositories:        testDatabase.Repositories,
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore),
	}

	now := utils.Now()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	updateEvent, err := event.NewOpportunityUpdateRenewalEvent(opportunityAggregate,
		"MEDIUM",
		"Updated likelihood",
		"user-123",
		"openline",
		float64(10000),
		now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = opportunityEventHandler.OnUpdateRenewal(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{constants.NodeLabel_Opportunity: 1, constants.NodeLabel_Opportunity + "_" + tenantName: 1})

	opportunityDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, constants.NodeLabel_Opportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(*opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, now, opportunity.UpdatedAt)
	require.Equal(t, "MEDIUM", opportunity.RenewalDetails.RenewalLikelihood)
	require.Equal(t, "user-123", opportunity.RenewalDetails.RenewalUpdatedByUserId)
	require.Equal(t, now, *opportunity.RenewalDetails.RenewalUpdatedByUserAt)
	require.Equal(t, float64(10000), opportunity.Amount)
	require.Equal(t, "Updated likelihood", opportunity.Comments)

	// Check no events were generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	var eventList []eventstore.Event
	for _, value := range eventsMap {
		eventList = value
	}
	require.Equal(t, 1, len(eventList))

	generatedEvent1 := eventList[0]
	require.Equal(t, event.OpportunityUpdateV1, generatedEvent1.EventType)
	var eventData1 event.OpportunityUpdateEvent
	err = generatedEvent1.GetJsonData(&eventData1)
	require.Nil(t, err)
	require.Equal(t, tenantName, eventData1.Tenant)
	require.Equal(t, float64(0), eventData1.Amount)
	require.Equal(t, float64(0), eventData1.MaxAmount)
}
