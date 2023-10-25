package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/model"
	orgaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	orgcmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	orgevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGraphInteractionEventEventHandler_OnCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	// prepare neo4j data
	externalSystemId := "sf"
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, externalSystemId)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	issueId := neo4jt.CreateIssue(ctx, testDatabase.Driver, tenantName, entity.IssueEntity{})
	neo4jt.LinkIssueReportedBy(ctx, testDatabase.Driver, issueId, orgId)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "ExternalSystem": 1, "Issue": 1, "TimelineEvent": 1, "InteractionEvent": 0})

	// prepare event handler
	interactionEventHandler := &GraphInteractionEventHandler{
		repositories:         testDatabase.Repositories,
		organizationCommands: orgcmdhnd.NewOrganizationCommands(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()
	interactionEventId := uuid.New().String()
	interactionEventAggregate := aggregate.NewInteractionEventAggregateWithTenantAndID(tenantName, interactionEventId)
	createEvent, err := event.NewInteractionEventCreateEvent(interactionEventAggregate, model.InteractionEventDataFields{
		Content:       "test content",
		ContentType:   "test content type",
		Channel:       "test channel",
		ChannelData:   "test channel data",
		Identifier:    "test identifier",
		EventType:     "test event type",
		PartOfIssueId: utils.StringPtr(issueId),
	}, cmnmod.Source{
		Source:        constants.SourceOpenline,
		AppSource:     constants.AppSourceEventProcessingPlatform,
		SourceOfTruth: constants.SourceOpenline,
	}, cmnmod.ExternalSystem{
		ExternalSystemId: "sf",
		ExternalId:       "123",
	}, now, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = interactionEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"ExternalSystem": 1, "ExternalSystem_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"InteractionEvent": 1, "InteractionEvent_" + tenantName: 1,
		"TimelineEvent": 2, "TimelineEvent_" + tenantName: 2})
	neo4jt.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"REPORTED_BY":    1,
		"IS_LINKED_WITH": 1,
		"PART_OF":        1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, interactionEventId, "PART_OF", issueId)
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, interactionEventId, "IS_LINKED_WITH", externalSystemId)

	interactionEventDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "InteractionEvent_"+tenantName, interactionEventId)
	require.Nil(t, err)
	require.NotNil(t, interactionEventDbNode)

	// verify interaction event
	interactionEvent := graph_db.MapDbNodeToInteractionEventEntity(*interactionEventDbNode)
	require.Equal(t, interactionEventId, interactionEvent.Id)
	require.Equal(t, "test content", interactionEvent.Content)
	require.Equal(t, "test content type", interactionEvent.ContentType)
	require.Equal(t, "test channel", interactionEvent.Channel)
	require.Equal(t, "test channel data", interactionEvent.ChannelData)
	require.Equal(t, "test identifier", interactionEvent.Identifier)
	require.Equal(t, "test event type", interactionEvent.EventType)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), interactionEvent.Source)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), interactionEvent.SourceOfTruth)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, interactionEvent.AppSource)
	require.Equal(t, now, interactionEvent.CreatedAt)
	require.Equal(t, now, interactionEvent.UpdatedAt)

	// Check refresh last touchpoint event was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[orgAggregate.GetID()]
	require.Equal(t, 1, len(eventList))
	generatedEvent := eventList[0]
	require.Equal(t, orgevents.OrganizationRefreshLastTouchpointV1, generatedEvent.EventType)
	var eventData orgevents.OrganizationRefreshLastTouchpointEvent
	err = generatedEvent.GetJsonData(&eventData)
	require.Nil(t, err)
	require.Equal(t, tenantName, eventData.Tenant)
}

func TestGraphInteractionEventEventHandler_OnUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	interactionEventId := neo4jt.CreateInteractionEvent(ctx, testDatabase.Driver, tenantName, entity.InteractionEventEntity{
		Content:     "test content",
		Channel:     "test channel",
		Identifier:  "test identifier",
		EventType:   "test event type",
		ContentType: "test content type",
		ChannelData: "test channel data",
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"InteractionEvent": 1, "TimelineEvent": 1})

	// prepare event handler
	interactionEventHandler := &GraphInteractionEventHandler{
		repositories: testDatabase.Repositories,
	}
	now := utils.Now()
	interactionEventAggregate := aggregate.NewInteractionEventAggregateWithTenantAndID(tenantName, interactionEventId)
	updateEvent, err := event.NewInteractionEventUpdateEvent(interactionEventAggregate, model.InteractionEventDataFields{
		Content:     "test content updated",
		ContentType: "test content type updated",
		Channel:     "test channel updated",
		ChannelData: "test channel data updated",
		Identifier:  "test identifier updated",
		EventType:   "test event type updated",
	}, constants.SourceOpenline, cmnmod.ExternalSystem{}, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = interactionEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"InteractionEvent": 1, "InteractionEvent_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	interactionEventDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "InteractionEvent_"+tenantName, interactionEventId)
	require.Nil(t, err)
	require.NotNil(t, interactionEventDbNode)

	// verify interaction event
	interactionEvent := graph_db.MapDbNodeToInteractionEventEntity(*interactionEventDbNode)
	require.Equal(t, interactionEventId, interactionEvent.Id)
	require.Equal(t, "test content updated", interactionEvent.Content)
	require.Equal(t, "test content type updated", interactionEvent.ContentType)
	require.Equal(t, "test channel updated", interactionEvent.Channel)
	require.Equal(t, "test channel data updated", interactionEvent.ChannelData)
	require.Equal(t, "test identifier updated", interactionEvent.Identifier)
	require.Equal(t, "test event type updated", interactionEvent.EventType)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), interactionEvent.SourceOfTruth)
	require.Equal(t, now, interactionEvent.UpdatedAt)
}

func TestGraphInteractionEventEventHandler_OnUpdate_CurrentSourceOpenline_UpdateSourceNonOpenline_UpdateOnlyEmptyFields(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	interactionEventId := neo4jt.CreateInteractionEvent(ctx, testDatabase.Driver, tenantName, entity.InteractionEventEntity{
		Content:       "test content",
		Channel:       "test channel",
		Identifier:    "test identifier",
		EventType:     "test event type",
		ContentType:   "test content type",
		ChannelData:   "test channel data",
		SourceOfTruth: constants.SourceOpenline,
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"InteractionEvent": 1, "TimelineEvent": 1})

	// prepare event handler
	interactionEventHandler := &GraphInteractionEventHandler{
		repositories: testDatabase.Repositories,
	}
	now := utils.Now()
	interactionEventAggregate := aggregate.NewInteractionEventAggregateWithTenantAndID(tenantName, interactionEventId)
	updateEvent, err := event.NewInteractionEventUpdateEvent(interactionEventAggregate, model.InteractionEventDataFields{
		Content:     "test content",
		Channel:     "test channel",
		Identifier:  "test identifier",
		EventType:   "test event type",
		ContentType: "test content type",
		ChannelData: "test channel data",
	}, "hubspot", cmnmod.ExternalSystem{}, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = interactionEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"InteractionEvent": 1, "InteractionEvent_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	interactionEventDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "InteractionEvent_"+tenantName, interactionEventId)
	require.Nil(t, err)
	require.NotNil(t, interactionEventDbNode)

	// verify interaction event
	interactionEvent := graph_db.MapDbNodeToInteractionEventEntity(*interactionEventDbNode)
	require.Equal(t, interactionEventId, interactionEvent.Id)
	require.Equal(t, "test content", interactionEvent.Content)
	require.Equal(t, "test content type", interactionEvent.ContentType)
	require.Equal(t, "test channel", interactionEvent.Channel)
	require.Equal(t, "test channel data", interactionEvent.ChannelData)
	require.Equal(t, "test identifier", interactionEvent.Identifier)
	require.Equal(t, "test event type", interactionEvent.EventType)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), interactionEvent.SourceOfTruth)
	require.Equal(t, now, interactionEvent.UpdatedAt)
}
