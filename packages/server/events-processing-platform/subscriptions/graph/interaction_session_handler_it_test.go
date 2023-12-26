package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGraphInteractionSessionEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	externalSystemId := "sf"
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, externalSystemId)

	// prepare event handler
	interactionSessionEventHandler := &InteractionSessionEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	now := utils.Now()
	interactionSessionId := uuid.New().String()
	interactionSessionAggregate := aggregate.NewInteractionSessionAggregateWithTenantAndID(tenantName, interactionSessionId)
	createEvent, err := event.NewInteractionSessionCreateEvent(interactionSessionAggregate, model.InteractionSessionDataFields{
		Channel:     "test channel",
		ChannelData: "test channel data",
		Identifier:  "test identifier",
		Name:        "test name",
		Status:      "test status",
		Type:        "test type",
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
	err = interactionSessionEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"ExternalSystem": 1, "ExternalSystem_" + tenantName: 1,
		"InteractionSession": 1, "InteractionSession_" + tenantName: 1,
		"Tenant": 1,
	})
	neo4jt.AssertNeo4jLabels(ctx, t, testDatabase.Driver, []string{"InteractionSession", "InteractionSession_" + tenantName, "Tenant", "ExternalSystem", "ExternalSystem_" + tenantName})
	neo4jt.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"IS_LINKED_WITH": 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, interactionSessionId, "IS_LINKED_WITH", externalSystemId)

	interactionSessionDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "InteractionSession_"+tenantName, interactionSessionId)
	require.Nil(t, err)
	require.NotNil(t, interactionSessionDbNode)

	// verify interaction session
	interactionEvent := graph_db.MapDbNodeToInteractionSessionEntity(*interactionSessionDbNode)
	require.Equal(t, interactionSessionId, interactionEvent.Id)
	require.Equal(t, "test channel", interactionEvent.Channel)
	require.Equal(t, "test channel data", interactionEvent.ChannelData)
	require.Equal(t, "test identifier", interactionEvent.Identifier)
	require.Equal(t, "test name", interactionEvent.Name)
	require.Equal(t, "test status", interactionEvent.Status)
	require.Equal(t, "test type", interactionEvent.Type)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), interactionEvent.Source)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), interactionEvent.SourceOfTruth)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, interactionEvent.AppSource)
	require.Equal(t, now, interactionEvent.CreatedAt)
	require.Equal(t, now, interactionEvent.UpdatedAt)

}
