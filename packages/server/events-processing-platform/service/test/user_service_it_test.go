package servicet

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	jobrolepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/job_role"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	events2 "github.com/openline-ai/openline-customer-os/packages/server/events/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestUserService_UpsertUser(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to events processing platform: %v", err)
	}
	userClient := userpb.NewUserGrpcServiceClient(grpcConnection)
	timeNow := utils.Now()
	userId, _ := uuid.NewUUID()

	response, err := userClient.UpsertUser(ctx, &userpb.UpsertUserGrpcRequest{
		Id:              userId.String(),
		Tenant:          "ziggy",
		FirstName:       "Bob",
		LastName:        "Dole",
		Name:            "Bob Dole",
		Internal:        true,
		Bot:             true,
		ProfilePhotoUrl: "https://www.google.com",
		Timezone:        "America/Los_Angeles",
		CreatedAt:       timestamppb.New(timeNow),
		UpdatedAt:       timestamppb.New(timeNow),
		SourceFields: &commonpb.SourceFields{
			AppSource: "event-processing-platform",
			Source:    "N/A",
		},
	})

	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, userId.String(), response.Id)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[aggregate.NewUserAggregateWithTenantAndID("ziggy", userId.String()).ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, events.UserCreateV1, eventList[0].EventType)
	var eventData events.UserCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	fmt.Printf("Got an envent %s\n", string(eventList[0].GetData()))
	if err := eventList[0].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, "Bob", eventData.FirstName)
	require.Equal(t, "Dole", eventData.LastName)
	require.Equal(t, "Bob Dole", eventData.Name)
	require.Equal(t, "event-processing-platform", eventData.SourceFields.AppSource)
	require.Equal(t, "N/A", eventData.SourceFields.Source)
	require.Equal(t, "N/A", eventData.SourceFields.SourceOfTruth)
	require.Equal(t, timeNow, eventData.CreatedAt)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, "ziggy", eventData.Tenant)
	require.True(t, eventData.Internal)
	require.True(t, eventData.Bot)
	require.Equal(t, "https://www.google.com", eventData.ProfilePhotoUrl)
	require.Equal(t, "America/Los_Angeles", eventData.Timezone)
}

func TestUserService_UpsertUserAndLinkJobRole(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to events processing platform: %v", err)
	}
	userClient := userpb.NewUserGrpcServiceClient(grpcConnection)
	jobRoleClient := jobrolepb.NewJobRoleGrpcServiceClient(grpcConnection)

	timeNow := utils.Now()
	userId, _ := uuid.NewUUID()

	createUserResponse, err := userClient.UpsertUser(ctx, &userpb.UpsertUserGrpcRequest{
		Id:        userId.String(),
		Tenant:    "ziggy",
		FirstName: "Bob",
		LastName:  "Dole",
		Name:      "Bob Dole",
		CreatedAt: timestamppb.New(timeNow),
		UpdatedAt: timestamppb.New(timeNow),
		SourceFields: &commonpb.SourceFields{
			AppSource: "event-processing-platform",
			Source:    "N/A",
		},
	})

	require.Nil(t, err)
	require.NotNil(t, createUserResponse)
	require.Equal(t, userId.String(), createUserResponse.Id)

	timeStarted := utils.Now().AddDate(0, -6, 0)
	timeEnded := utils.Now().AddDate(0, 6, 0)
	description := "I clean things"
	createJobRoleResponse, err := jobRoleClient.CreateJobRole(ctx, &jobrolepb.CreateJobRoleGrpcRequest{
		Tenant:        "ziggy",
		JobTitle:      "Chief Janitor",
		Description:   &description,
		Source:        "N/A",
		SourceOfTruth: "N/A",
		AppSource:     "event-processing-platform",
		CreatedAt:     timestamppb.New(timeNow),
		StartedAt:     timestamppb.New(timeStarted),
		EndedAt:       timestamppb.New(timeEnded),
	})
	if err != nil {
		t.Fatalf("Failed to create job role: %v", err)
	}
	require.Nil(t, err)
	require.NotNil(t, createJobRoleResponse)

	linkJobRoleResponse, err := userClient.LinkJobRoleToUser(ctx, &userpb.LinkJobRoleToUserGrpcRequest{
		UserId:    createUserResponse.Id,
		JobRoleId: createJobRoleResponse.Id,
		Tenant:    "ziggy",
		AppSource: "event-processing-platform",
	})

	if err != nil {
		t.Fatalf("Failed to link job role to user: %v", err)
	}
	require.Nil(t, err)
	require.NotNil(t, linkJobRoleResponse)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 2, len(eventsMap))
	eventList := eventsMap[aggregate.NewUserAggregateWithTenantAndID("ziggy", userId.String()).ID]
	require.Equal(t, 2, len(eventList))
	require.Equal(t, events.UserCreateV1, eventList[0].EventType)
	require.Equal(t, events.UserJobRoleLinkV1, eventList[1].EventType)
}

func TestUserService_LinkEmail(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"

	aggregateStore := eventstore.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	userClient := userpb.NewUserGrpcServiceClient(grpcConnection)

	userId := uuid.New().String()
	emailId := uuid.New().String()

	// Grpc call
	response, err := userClient.LinkEmailToUser(ctx, &userpb.LinkEmailToUserGrpcRequest{
		Tenant:         tenant,
		UserId:         userId,
		LoggedInUserId: userId,
		EmailId:        emailId,
		Primary:        true,
		Label:          "work",
	})
	require.Nil(t, err)

	// Assert response
	require.NotNil(t, response)
	require.NotEmpty(t, response.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[aggregate.NewUserAggregateWithTenantAndID("ziggy", userId).ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, events.UserEmailLinkV1, eventList[0].EventType)
	var eventData events.UserLinkEmailEvent

	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err)

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, emailId, eventData.EmailId)
	require.Equal(t, "work", eventData.Label)
	require.True(t, eventData.Primary)
}

func TestUserService_LinkEmail_IgnoreDuplicates(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	userId := uuid.New().String()
	emailId := uuid.New().String()

	aggregateStore := eventstore.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	userClient := userpb.NewUserGrpcServiceClient(grpcConnection)

	// Grpc call
	response1, err := userClient.LinkEmailToUser(ctx, &userpb.LinkEmailToUserGrpcRequest{
		Tenant:    tenant,
		UserId:    userId,
		EmailId:   emailId,
		Primary:   false,
		Label:     "work",
		AppSource: events2.AppSourceIntegrationApp,
	})
	require.Nil(t, err)
	// Second grpc call
	response2, err := userClient.LinkEmailToUser(ctx, &userpb.LinkEmailToUserGrpcRequest{
		Tenant:    tenant,
		UserId:    userId,
		EmailId:   emailId,
		Primary:   true,
		Label:     "work",
		AppSource: events2.AppSourceIntegrationApp,
	})
	require.Nil(t, err)

	// Assert response
	require.NotNil(t, response1)
	require.NotEmpty(t, response1.Id)
	require.NotNil(t, response2)
	require.NotEmpty(t, response2.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[aggregate.NewUserAggregateWithTenantAndID("ziggy", userId).ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, events.UserEmailLinkV1, eventList[0].EventType)
	var eventData events.UserLinkEmailEvent

	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err)

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, emailId, eventData.EmailId)
	require.Equal(t, "work", eventData.Label)
	require.False(t, eventData.Primary)
}

func TestUserService_LinkPhoneNumber(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"

	aggregateStore := eventstore.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	userClient := userpb.NewUserGrpcServiceClient(grpcConnection)

	userId := uuid.New().String()
	phoneNumberId := uuid.New().String()

	// Grpc call
	response, err := userClient.LinkPhoneNumberToUser(ctx, &userpb.LinkPhoneNumberToUserGrpcRequest{
		Tenant:         tenant,
		UserId:         userId,
		LoggedInUserId: userId,
		PhoneNumberId:  phoneNumberId,
		Primary:        true,
		Label:          "work",
	})
	require.Nil(t, err)

	// Assert response
	require.NotNil(t, response)
	require.NotEmpty(t, response.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[aggregate.NewUserAggregateWithTenantAndID("ziggy", userId).ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, events.UserPhoneNumberLinkV1, eventList[0].EventType)
	var eventData events.UserLinkPhoneNumberEvent

	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err)

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, phoneNumberId, eventData.PhoneNumberId)
	require.Equal(t, "work", eventData.Label)
	require.True(t, eventData.Primary)
}

func TestUserService_LinkPhoneNumber_IgnoreDuplicates(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	userId := uuid.New().String()
	phoneNumberId := uuid.New().String()

	aggregateStore := eventstore.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	userClient := userpb.NewUserGrpcServiceClient(grpcConnection)

	// Grpc call
	response1, err := userClient.LinkPhoneNumberToUser(ctx, &userpb.LinkPhoneNumberToUserGrpcRequest{
		Tenant:        tenant,
		UserId:        userId,
		PhoneNumberId: phoneNumberId,
		Primary:       false,
		Label:         "work",
		AppSource:     events2.AppSourceIntegrationApp,
	})
	require.Nil(t, err)
	// Second grpc call same label different primary flag
	response2, err := userClient.LinkPhoneNumberToUser(ctx, &userpb.LinkPhoneNumberToUserGrpcRequest{
		Tenant:        tenant,
		UserId:        userId,
		PhoneNumberId: phoneNumberId,
		Primary:       true,
		Label:         "work",
		AppSource:     events2.AppSourceIntegrationApp,
	})
	require.Nil(t, err)
	// Third grpc call different label
	response3, err := userClient.LinkPhoneNumberToUser(ctx, &userpb.LinkPhoneNumberToUserGrpcRequest{
		Tenant:        tenant,
		UserId:        userId,
		PhoneNumberId: phoneNumberId,
		Primary:       true,
		Label:         "home",
		AppSource:     events2.AppSourceIntegrationApp,
	})
	require.Nil(t, err)

	// Assert response
	require.NotEmpty(t, response1.Id)
	require.NotEmpty(t, response2.Id)
	require.NotEmpty(t, response3.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[aggregate.NewUserAggregateWithTenantAndID("ziggy", userId).ID]
	require.Equal(t, 2, len(eventList))

	require.Equal(t, events.UserPhoneNumberLinkV1, eventList[0].EventType)
	var eventData1 events.UserLinkPhoneNumberEvent
	var eventData2 events.UserLinkPhoneNumberEvent

	err = eventList[0].GetJsonData(&eventData1)
	require.Nil(t, err)
	err = eventList[1].GetJsonData(&eventData2)
	require.Nil(t, err)

	require.Equal(t, tenant, eventData1.Tenant)
	require.Equal(t, phoneNumberId, eventData1.PhoneNumberId)
	require.Equal(t, "work", eventData1.Label)
	require.False(t, eventData1.Primary)

	require.Equal(t, tenant, eventData2.Tenant)
	require.Equal(t, phoneNumberId, eventData2.PhoneNumberId)
	require.Equal(t, "home", eventData2.Label)
	require.True(t, eventData2.Primary)
}
