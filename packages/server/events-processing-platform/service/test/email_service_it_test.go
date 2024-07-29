package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email/event"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestEmailService_UpsertEmail(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to events processing platform: %v", err)
	}
	emailClient := emailpb.NewEmailGrpcServiceClient(grpcConnection)
	timeNow := utils.Now()
	emailId, _ := uuid.NewUUID()
	response, err := emailClient.UpsertEmail(ctx, &emailpb.UpsertEmailGrpcRequest{
		Tenant:   "ziggy",
		RawEmail: "test@openline.ai",
		SourceFields: &commonpb.SourceFields{
			AppSource:     "event-processing-platform",
			Source:        "N/A",
			SourceOfTruth: "N/A",
		},
		CreatedAt: timestamppb.New(timeNow),
		UpdatedAt: timestamppb.New(timeNow),
		Id:        emailId.String(),
	})
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, emailId.String(), response.Id)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[email.NewEmailAggregateWithTenantAndID("ziggy", emailId.String()).ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, event.EmailCreateV1, eventList[0].GetEventType())
	var eventData event.EmailCreateEvent
	if err := eventList[0].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, "test@openline.ai", eventData.RawEmail)
	require.Equal(t, "event-processing-platform", eventData.SourceFields.AppSource)
	require.Equal(t, "N/A", eventData.SourceFields.Source)
	require.Equal(t, "N/A", eventData.SourceFields.SourceOfTruth)
	require.Equal(t, timeNow, eventData.CreatedAt)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, "ziggy", eventData.Tenant)

}
