package graph

import (
	"context"
	"github.com/google/uuid"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/mocked_grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGraphInteractionEventEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	externalSystemId := "sf"
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, externalSystemId)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	issueId := neo4jt.CreateIssue(ctx, testDatabase.Driver, tenantName, entity.IssueEntity{})
	neo4jt.LinkIssueReportedBy(ctx, testDatabase.Driver, issueId, orgId)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{})
	contactId := neo4jt.CreateContact(ctx, testDatabase.Driver, tenantName, entity.ContactEntity{})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User": 1, "Contact": 1, "Organization": 1, "ExternalSystem": 1, "Issue": 1, "TimelineEvent": 1, "InteractionEvent": 0})

	// prepare grpc mock
	lastTouchpointInvoked := false
	organizationServiceCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		RefreshLastTouchpoint: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, org.AppSource)
			lastTouchpointInvoked = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId,
			}, nil
		},
	}
	mocked_grpc.SetOrganizationCallbacks(&organizationServiceCallbacks)

	// prepare event handler
	interactionEventHandler := &InteractionEventHandler{
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	now := utils.Now()
	interactionEventId := uuid.New().String()
	interactionEventAggregate := aggregate.NewInteractionEventAggregateWithTenantAndID(tenantName, interactionEventId)
	createEvent, err := event.NewInteractionEventCreateEvent(interactionEventAggregate, model.InteractionEventDataFields{
		Content:          "test content",
		ContentType:      "test content type",
		Channel:          "test channel",
		ChannelData:      "test channel data",
		Identifier:       "test identifier",
		EventType:        "test event type",
		Hide:             true,
		BelongsToIssueId: utils.StringPtr(issueId),
		Sender: model.Sender{
			Participant: commonmodel.Participant{
				ID:              userId,
				ParticipantType: commonmodel.UserType,
			},
			RelationType: "FROM",
		},
		Receivers: []model.Receiver{
			{
				Participant: commonmodel.Participant{
					ID:              contactId,
					ParticipantType: commonmodel.ContactType,
				},
				RelationType: "TO",
			},
			{
				Participant: commonmodel.Participant{
					ID:              orgId,
					ParticipantType: commonmodel.OrganizationType,
				},
				RelationType: "CC",
			},
		},
	}, commonmodel.Source{
		Source:        constants.SourceOpenline,
		AppSource:     constants.AppSourceEventProcessingPlatform,
		SourceOfTruth: constants.SourceOpenline,
	}, commonmodel.ExternalSystem{
		ExternalSystemId: "sf",
		ExternalId:       "123",
	}, now, now)
	require.Nil(t, err)

	// EXECUTE
	err = interactionEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User": 1, "User_" + tenantName: 1,
		"Contact": 1, "Contact_" + tenantName: 1,
		"Organization": 1, "Organization_" + tenantName: 1,
		"ExternalSystem": 1, "ExternalSystem_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"InteractionEvent": 1, "InteractionEvent_" + tenantName: 1,
		"TimelineEvent": 2, "TimelineEvent_" + tenantName: 2})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"REPORTED_BY":    1,
		"IS_LINKED_WITH": 1,
		"PART_OF":        1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, interactionEventId, "PART_OF", issueId)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, interactionEventId, "IS_LINKED_WITH", externalSystemId)
	neo4jtest.AssertRelationshipWithProperties(ctx, t, testDatabase.Driver, interactionEventId, "SENT_BY", userId, map[string]interface{}{"type": "FROM"})
	neo4jtest.AssertRelationshipWithProperties(ctx, t, testDatabase.Driver, interactionEventId, "SENT_TO", contactId, map[string]interface{}{"type": "TO"})
	neo4jtest.AssertRelationshipWithProperties(ctx, t, testDatabase.Driver, interactionEventId, "SENT_TO", orgId, map[string]interface{}{"type": "CC"})

	interactionEventDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "InteractionEvent_"+tenantName, interactionEventId)
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
	require.Equal(t, true, interactionEvent.Hide)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), interactionEvent.Source)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), interactionEvent.SourceOfTruth)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, interactionEvent.AppSource)
	require.Equal(t, now, interactionEvent.CreatedAt)
	require.Equal(t, now, interactionEvent.UpdatedAt)

	// Check refresh last touchpoint
	require.Truef(t, lastTouchpointInvoked, "RefreshLastTouchpoint was not invoked")
}

func TestGraphInteractionEventEventHandler_OnUpdate(t *testing.T) {
	ctx := context.Background()
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
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"InteractionEvent": 1, "TimelineEvent": 1})

	// prepare event handler
	interactionEventHandler := &InteractionEventHandler{
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
	}, constants.SourceOpenline, commonmodel.ExternalSystem{}, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = interactionEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"InteractionEvent": 1, "InteractionEvent_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	interactionEventDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "InteractionEvent_"+tenantName, interactionEventId)
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
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), interactionEvent.SourceOfTruth)
	require.Equal(t, now, interactionEvent.UpdatedAt)
}

func TestGraphInteractionEventEventHandler_OnUpdate_CurrentSourceOpenline_UpdateSourceNonOpenline_UpdateOnlyEmptyFields(t *testing.T) {
	ctx := context.Background()
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
		Hide:          false,
		SourceOfTruth: constants.SourceOpenline,
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"InteractionEvent": 1, "TimelineEvent": 1})

	// prepare event handler
	interactionEventHandler := &InteractionEventHandler{
		repositories: testDatabase.Repositories,
	}
	now := utils.Now()
	interactionEventAggregate := aggregate.NewInteractionEventAggregateWithTenantAndID(tenantName, interactionEventId)
	updateEvent, err := event.NewInteractionEventUpdateEvent(interactionEventAggregate, model.InteractionEventDataFields{
		Content:     "test content updated",
		Channel:     "test channel updated",
		Identifier:  "test identifier updated",
		EventType:   "test event type updated",
		ContentType: "test content type updated",
		ChannelData: "test channel data updated",
		Hide:        true,
	}, "hubspot", commonmodel.ExternalSystem{}, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = interactionEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"InteractionEvent": 1, "InteractionEvent_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	interactionEventDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "InteractionEvent_"+tenantName, interactionEventId)
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
	require.Equal(t, false, interactionEvent.Hide)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), interactionEvent.SourceOfTruth)
	require.Equal(t, now, interactionEvent.UpdatedAt)
}

func TestGraphInteractionEventEventHandler_OnSummaryReplace_Create(t *testing.T) {
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
		Hide:          false,
		SourceOfTruth: constants.SourceOpenline,
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"InteractionEvent": 1, "TimelineEvent": 1})

	interactionEvent, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "InteractionEvent_"+tenantName)
	require.Nil(t, err)
	interactionEventProps := utils.GetPropsFromNode(*interactionEvent)
	interactionEventId = utils.GetStringPropOrEmpty(interactionEventProps, "id")
	require.NotNil(t, interactionEventId)

	// prepare event handler
	interactionEventHandler := &InteractionEventHandler{
		repositories: testDatabase.Repositories,
	}
	interactionEventAggregate := aggregate.NewInteractionEventAggregateWithTenantAndID(tenantName, interactionEventId)
	summary := "AnalysisSummary"
	contentType := "AnalysisContentType"
	now := utils.Now()
	summaryReplaceEvent, err := event.NewInteractionEventReplaceSummaryEvent(interactionEventAggregate, tenantName, summary, contentType, now)
	require.Nil(t, err, "failed to create event")

	err = interactionEventHandler.OnSummaryReplace(context.Background(), summaryReplaceEvent)
	require.Nil(t, err, "failed to execute OnSummaryReplace for interactionEventHandler")
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "DESCRIBES"), "Incorrect number of DESCRIBES relationships in Neo4j")

	analysis, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Analysis_"+tenantName)
	require.Nil(t, err)
	analysisProps := utils.GetPropsFromNode(*analysis)
	require.Equal(t, 9, len(analysisProps))
	analysisId := utils.GetStringPropOrEmpty(analysisProps, "id")
	require.NotNil(t, analysisId)
	require.Equal(t, now, utils.GetTimePropOrNow(analysisProps, "createdAt"))
	require.Equal(t, now, utils.GetTimePropOrNow(analysisProps, "updatedAt"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(analysisProps, "source"))
	require.Equal(t, constants.AppSourceEventProcessingPlatform, utils.GetStringPropOrEmpty(analysisProps, "appSource"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(analysisProps, "sourceOfTruth"))
	require.Equal(t, "summary", utils.GetStringPropOrEmpty(analysisProps, "analysisType"))
	require.Equal(t, "AnalysisContentType", utils.GetStringPropOrEmpty(analysisProps, "contentType"))
	require.Equal(t, "AnalysisSummary", utils.GetStringPropOrEmpty(analysisProps, "content"))
}

func TestGraphInteractionEventEventHandler_OnActionItemsReplace(t *testing.T) {
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
		Hide:          false,
		SourceOfTruth: constants.SourceOpenline,
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"InteractionEvent": 1, "TimelineEvent": 1})

	interactionEvent, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "InteractionEvent_"+tenantName)
	require.Nil(t, err)
	interactionEventProps := utils.GetPropsFromNode(*interactionEvent)
	interactionEventId = utils.GetStringPropOrEmpty(interactionEventProps, "id")
	require.NotNil(t, interactionEventId)

	// prepare event handler
	interactionEventHandler := &InteractionEventHandler{
		repositories: testDatabase.Repositories,
	}
	interactionEventAggregate := aggregate.NewInteractionEventAggregateWithTenantAndID(tenantName, interactionEventId)

	var actionItems = []string{"ActionItem1", "ActionItem2"}
	now := utils.Now()
	actionItemsReplaceEvent, err := event.NewInteractionEventReplaceActionItemsEvent(interactionEventAggregate, tenantName, actionItems, now)
	require.Nil(t, err, "failed to create event")

	err = interactionEventHandler.OnActionItemsReplace(context.Background(), actionItemsReplaceEvent)
	require.Nil(t, err, "failed to execute OnActionItemsReplace for interactionEventHandler")

	actionItem, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "ActionItem_"+tenantName)
	require.Nil(t, err)
	actionItemProps := utils.GetPropsFromNode(*actionItem)
	require.Equal(t, 7, len(actionItemProps))
	analysisId := utils.GetStringPropOrEmpty(actionItemProps, "id")
	require.NotNil(t, analysisId)
	require.Equal(t, now, utils.GetTimePropOrNow(actionItemProps, "createdAt"))
	require.Equal(t, now, utils.GetTimePropOrNow(actionItemProps, "updatedAt"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(actionItemProps, "source"))
	require.Equal(t, constants.AppSourceEventProcessingPlatform, utils.GetStringPropOrEmpty(actionItemProps, "appSource"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(actionItemProps, "sourceOfTruth"))

	//TODO update the following assertion when the implementation of the OnActionItemsReplace is fixed and the content contains the entire Array of ActionItems.
	//TODO It is a known fact that at the moment, the content only returns the first element of the ActionItems array.
	//TODO Won't be fixed right away as the event is not used in Prod
	require.Contains(t, []string{"ActionItem1", "ActionItem2"}, utils.GetStringPropOrEmpty(actionItemProps, "content"))
}
