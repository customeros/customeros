package servicet

import (
	"context"
	contactAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/event"
	emailAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	emailEvents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestContactService_CreateContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to emailEvents processing platform: %v", err)
	}
	contactClient := contactpb.NewContactGrpcServiceClient(grpcConnection)
	timeNow := utils.Now()
	response, err := contactClient.UpsertContact(ctx, &contactpb.UpsertContactGrpcRequest{
		Tenant:          "ziggy",
		FirstName:       "Bob",
		LastName:        "Smith",
		Prefix:          "Mr.",
		Description:     "This is a contact description",
		Timezone:        "America/Los_Angeles",
		ProfilePhotoUrl: "https://www.google.com",
		AppSource:       "unit-test",
		Source:          "N/A",
		SourceOfTruth:   "N/A",
		CreatedAt:       timestamppb.New(timeNow),
	})
	if err != nil {
		t.Errorf("Failed to create contact: %v", err)
	}
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[contactAggregate.NewContactAggregateWithTenantAndID("ziggy", response.Id).ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, event.ContactCreateV1, eventList[0].GetEventType())
	var eventData event.ContactCreateEvent
	if err := eventList[0].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, "Bob", eventData.FirstName)
	require.Equal(t, "Smith", eventData.LastName)
	require.Equal(t, "Mr.", eventData.Prefix)
	require.Equal(t, "America/Los_Angeles", eventData.Timezone)
	require.Equal(t, "https://www.google.com", eventData.ProfilePhotoUrl)
	require.Equal(t, "unit-test", eventData.AppSource)
	require.Equal(t, "N/A", eventData.Source)
	require.Equal(t, "N/A", eventData.SourceOfTruth)
	require.Equal(t, timeNow, eventData.CreatedAt)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, "ziggy", eventData.Tenant)

}

func TestContactService_CreateContactWithEmail(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to emailEvents processing platform: %v", err)
	}
	contactClient := contactpb.NewContactGrpcServiceClient(grpcConnection)
	emailClient := emailpb.NewEmailGrpcServiceClient(grpcConnection)

	timeNow := utils.Now()
	responseContact, err := contactClient.UpsertContact(ctx, &contactpb.UpsertContactGrpcRequest{
		Tenant:          "ziggy",
		FirstName:       "Bob",
		LastName:        "Smith",
		Prefix:          "Mr.",
		Description:     "This is a contact description",
		Timezone:        "America/Los_Angeles",
		ProfilePhotoUrl: "https://www.google.com?id=123",
		AppSource:       "unit-test",
		Source:          "N/A",
		SourceOfTruth:   "N/A",
		CreatedAt:       timestamppb.New(timeNow),
	})
	if err != nil {
		t.Errorf("Failed to create contact: %v", err)
	}
	require.NotNil(t, responseContact)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	contactEventList := eventsMap[contactAggregate.NewContactAggregateWithTenantAndID("ziggy", responseContact.Id).ID]
	require.Equal(t, 1, len(contactEventList))
	require.Equal(t, event.ContactCreateV1, contactEventList[0].GetEventType())
	var createEventData event.ContactCreateEvent
	if err := contactEventList[0].GetJsonData(&createEventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}

	responseEmail, err := emailClient.UpsertEmail(ctx, &emailpb.UpsertEmailGrpcRequest{
		Tenant:        "ziggy",
		RawEmail:      "test@openline.ai",
		AppSource:     "unit-test",
		Source:        "N/A",
		SourceOfTruth: "N/A",
		CreatedAt:     timestamppb.New(timeNow),
		UpdatedAt:     timestamppb.New(timeNow),
		Id:            "",
	})
	if err != nil {
		t.Errorf("Failed to create email: %v", err)
	}
	require.Nil(t, err)
	require.NotNil(t, responseEmail)

	emailEventList := eventsMap[emailAggregate.NewEmailAggregateWithTenantAndID("ziggy", responseEmail.Id).ID]
	require.Equal(t, 1, len(emailEventList))
	require.Equal(t, emailEvents.EmailCreateV1, emailEventList[0].GetEventType())
	var eventData emailEvents.EmailCreateEvent
	if err := emailEventList[0].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, "test@openline.ai", eventData.RawEmail)

	responseLinkEmail, err := contactClient.LinkEmailToContact(ctx, &contactpb.LinkEmailToContactGrpcRequest{
		Tenant:    "ziggy",
		ContactId: responseContact.Id,
		EmailId:   responseEmail.Id,
		Primary:   true,
		Label:     "WORK",
		AppSource: "unit-test",
	})
	if err != nil {
		t.Errorf("Failed to link email to contact: %v", err)
	}
	require.Nil(t, err)
	require.NotNil(t, responseLinkEmail)

	contactEventList = eventsMap[contactAggregate.NewContactAggregateWithTenantAndID("ziggy", responseContact.Id).ID]

	require.Equal(t, 2, len(contactEventList))
	require.Equal(t, event.ContactEmailLinkV1, contactEventList[1].GetEventType())
	var linkEmailToContact event.ContactLinkEmailEvent
	if err := contactEventList[1].GetJsonData(&linkEmailToContact); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, responseEmail.Id, linkEmailToContact.EmailId)
	require.Equal(t, "ziggy", linkEmailToContact.Tenant)
	require.Equal(t, "WORK", linkEmailToContact.Label)
	require.Equal(t, true, linkEmailToContact.Primary)

}
