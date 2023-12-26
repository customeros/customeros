package servicet

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/event"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	ispb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_session"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestInteractionSessionService_UpsertInteractionSession_Create(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, _ := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	interactionSessionClient := ispb.NewInteractionSessionGrpcServiceClient(grpcConnection)
	now := utils.Now()

	response, err := interactionSessionClient.UpsertInteractionSession(ctx, &ispb.UpsertInteractionSessionGrpcRequest{
		Tenant:      tenant,
		Name:        "Name-123",
		Channel:     "Channel-123",
		ChannelData: "ChannelData-123",
		Identifier:  "Identifier-123",
		Type:        "Type-123",
		Status:      "Status-123",
		ExternalSystemFields: &commonpb.ExternalSystemFields{
			ExternalSystemId: "ExternalSystemID",
			ExternalUrl:      "http://external.url",
			ExternalId:       "ExternalID",
			ExternalIdSecond: "ExternalIDSecond",
			ExternalSource:   "ExternalSource",
			SyncDate:         timestamppb.New(now),
		},
		SourceFields: &commonpb.SourceFields{
			Source:    "Source",
			AppSource: "AppSource",
		},
	})
	require.Nil(t, err)

	require.NotNil(t, response)
	interactionSessionId := response.Id

	// Verify the aggregate store
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	interactionSessionAggregate := aggregate.NewInteractionSessionAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[interactionSessionAggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, event.InteractionSessionCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(aggregate.InteractionSessionAggregateType)+"-"+tenant+"-"+interactionSessionId, eventList[0].GetAggregateID())

	var eventData event.InteractionSessionCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err)

	// Assertions to validate the interaction session create event data
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "Name-123", eventData.Name)
	require.Equal(t, "Channel-123", eventData.Channel)
	require.Equal(t, "ChannelData-123", eventData.ChannelData)
	require.Equal(t, "Identifier-123", eventData.Identifier)
	require.Equal(t, "Type-123", eventData.Type)
	require.Equal(t, "Status-123", eventData.Status)
	require.Equal(t, "Source", eventData.Source)
	require.Equal(t, "AppSource", eventData.AppSource)
	require.Equal(t, "ExternalSystemID", eventData.ExternalSystem.ExternalSystemId)
	require.Equal(t, "http://external.url", eventData.ExternalSystem.ExternalUrl)
	require.Equal(t, "ExternalID", eventData.ExternalSystem.ExternalId)
	require.Equal(t, "ExternalIDSecond", eventData.ExternalSystem.ExternalIdSecond)
	require.Equal(t, "ExternalSource", eventData.ExternalSystem.ExternalSource)
	require.True(t, now.Equal(*eventData.ExternalSystem.SyncDate))
}
