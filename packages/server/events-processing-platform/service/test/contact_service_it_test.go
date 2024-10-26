package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact"
	contactevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/contact/event"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestContactService_AddSocial(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenantName := "ziggy"
	contactId := uuid.New().String()
	now := utils.Now()

	// setup aggregate and create initial event
	aggregateStore := eventstoret.NewTestAggregateStore()
	contactAggregate := contact.NewContactAggregateWithTenantAndID(tenantName, contactId)
	newEvent, _ := contactevent.NewContactCreateEvent(contactAggregate, contactevent.ContactDataFields{}, commonmodel.Source{}, commonmodel.ExternalSystem{}, now, now)
	contactAggregate.UncommittedEvents = append(contactAggregate.UncommittedEvents, newEvent)
	aggregateStore.Save(ctx, contactAggregate)

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to get grpc connection")
	contactServiceClient := contactpb.NewContactGrpcServiceClient(grpcConnection)

	response, err := contactServiceClient.AddSocial(ctx, &contactpb.ContactAddSocialGrpcRequest{
		Tenant:   tenantName,
		SocialId: uuid.New().String(),
		SourceFields: &commonpb.SourceFields{
			Source:    "test",
			AppSource: "test-app",
		},
		CreatedAt: timestamppb.New(now),
		ContactId: contactId,
		Url:       "https://www.google.com",
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	eventList := eventsMap[contactAggregate.ID]
	require.Equal(t, 2, len(eventList))
	require.Equal(t, contactevent.ContactAddSocialV1, eventList[1].GetEventType())
	require.Equal(t, string(contact.ContactAggregateType)+"-"+tenantName+"-"+contactId, eventList[1].GetAggregateID())

	var eventData contactevent.ContactAddSocialEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, tenantName, eventData.Tenant)
	require.Equal(t, now, eventData.CreatedAt)
	require.Equal(t, "test", eventData.Source.Source)
	require.Equal(t, "test-app", eventData.Source.AppSource)
	require.Equal(t, "https://www.google.com", eventData.Url)
	require.NotEmpty(t, eventData.SocialId)
}
