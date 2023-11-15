package graph

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
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
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, aggregateStore),
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
		opportunityCommands: opportunitycmdhandler.NewCommandHandlers(testLogger, aggregateStore),
	}

	// Create an OpportunityCreateEvent
	opportunityId := uuid.New().String()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	timeNow := utils.Now()
	createEvent, err := event.NewOpportunityCreateRenewalEvent(
		opportunityAggregate,
		contractId,
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
