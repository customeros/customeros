package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/neo4j"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGraphInteractionSessionEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	externalSystemId := "sf"
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
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
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		SourceOfTruth: constants.SourceOpenline,
	}, commonmodel.ExternalSystem{
		ExternalSystemId: "sf",
		ExternalId:       "123",
	}, now, now)
	require.Nil(t, err)

	// EXECUTE
	err = interactionSessionEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"ExternalSystem": 1, "ExternalSystem_" + tenantName: 1,
		"InteractionSession": 1, "InteractionSession_" + tenantName: 1,
		"Tenant": 1,
	})
	neo4jtest.AssertNeo4jLabels(ctx, t, testDatabase.Driver, []string{"InteractionSession", "InteractionSession_" + tenantName, "Tenant", "ExternalSystem", "ExternalSystem_" + tenantName})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"IS_LINKED_WITH": 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, interactionSessionId, "IS_LINKED_WITH", externalSystemId)

	interactionSessionDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "InteractionSession_"+tenantName, interactionSessionId)
	require.Nil(t, err)
	require.NotNil(t, interactionSessionDbNode)

	// verify interaction session
	interactionEvent := neo4jmapper.MapDbNodeToInteractionSessionEntity(interactionSessionDbNode)
	require.Equal(t, interactionSessionId, interactionEvent.Id)
	require.Equal(t, "test channel", interactionEvent.Channel)
	require.Equal(t, "test channel data", interactionEvent.ChannelData)
	require.Equal(t, "test identifier", interactionEvent.Identifier)
	require.Equal(t, "test name", interactionEvent.Name)
	require.Equal(t, "test status", interactionEvent.Status)
	require.Equal(t, "test type", interactionEvent.Type)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), interactionEvent.Source)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), interactionEvent.SourceOfTruth)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, interactionEvent.AppSource)
	require.Equal(t, now, interactionEvent.CreatedAt)
	test.AssertRecentTime(t, interactionEvent.UpdatedAt)

}
