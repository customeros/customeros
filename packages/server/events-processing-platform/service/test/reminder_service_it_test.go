package servicet

import (
	"context"
	"testing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	reminderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
	"github.com/stretchr/testify/require"
)

func TestReminderService_CreateReminder(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to get grpc connection")

	reminderClient := reminderpb.NewReminderGrpcServiceClient(grpcConnection)

	response, err := reminderClient.CreateReminder(ctx, &reminderpb.CreateReminderGrpcRequest{
		Tenant:         tenant,
		Content:        "New Reminder",
		LoggedInUserId: "user",
		OrganizationId: "org",
		Dismissed:      false,
		DueDate:        utils.ConvertTimeToTimestampPtr(utils.NowPtr()),
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
			Source:    "source",
		},
		CreatedAt: utils.ConvertTimeToTimestampPtr(utils.NowPtr()),
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	reminderId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	reminderAggregate := reminder.NewReminderAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[reminderAggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, reminder.ReminderCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(reminder.ReminderAggregateType)+"-"+tenant+"-"+reminderId, eventList[0].GetAggregateID())

	var eventData reminder.ReminderCreateEvent

	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "New Reminder", eventData.Content)
	test.AssertRecentTime(t, eventData.CreatedAt)
}

func TestReminderService_UpdateReminder(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to get grpc connection")

	reminderClient := reminderpb.NewReminderGrpcServiceClient(grpcConnection)

	createResponse, err := reminderClient.CreateReminder(ctx, &reminderpb.CreateReminderGrpcRequest{
		Tenant:         tenant,
		Content:        "New Reminder",
		LoggedInUserId: "user",
		OrganizationId: "org",
		Dismissed:      false,
		DueDate:        utils.ConvertTimeToTimestampPtr(utils.NowPtr()),
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
			Source:    "source",
		},
		// CreatedAt: utils.ConvertTimeToTimestampPtr(utils.NowPtr()),
	})
	require.Nil(t, err)
	require.NotNil(t, createResponse)

	reminderId := createResponse.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	reminderAggregate := reminder.NewReminderAggregateWithTenantAndID(tenant, createResponse.Id)
	eventList := eventsMap[reminderAggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, reminder.ReminderCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(reminder.ReminderAggregateType)+"-"+tenant+"-"+reminderId, eventList[0].GetAggregateID())

	updateResponse, err := reminderClient.UpdateReminder(ctx, &reminderpb.UpdateReminderGrpcRequest{
		Tenant:         tenant,
		ReminderId:     reminderId,
		Content:        "Updated Reminder",
		AppSource:      "app",
		LoggedInUserId: "user",
		FieldsMask: []reminderpb.ReminderFieldMask{
			reminderpb.ReminderFieldMask_REMINDER_PROPERTY_CONTENT,
		},
		// Dismissed: false,
	})
	require.Nil(t, err)
	require.NotNil(t, updateResponse)

	eventsMap = aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList = eventsMap[reminderAggregate.ID]
	require.Equal(t, 2, len(eventList))
	require.Equal(t, reminder.ReminderUpdateV1, eventList[1].GetEventType())
	require.Equal(t, string(reminder.ReminderAggregateType)+"-"+tenant+"-"+reminderId, eventList[1].GetAggregateID())

	var eventData reminder.ReminderUpdateEvent

	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "Updated Reminder", eventData.Content)
}
