package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/comment"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestCommentService_UpsertComment_CreateComment(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to connect to processing platform")
	commentClient := commentpb.NewCommentGrpcServiceClient(grpcConnection)
	timeNow := utils.Now()
	tenant := "ziggy"
	response, err := commentClient.UpsertComment(ctx, &commentpb.UpsertCommentGrpcRequest{
		Tenant:      tenant,
		Content:     "This is a log entry",
		ContentType: "text/plain",
		CreatedAt:   timestamppb.New(timeNow),
		SourceFields: &commonpb.SourceFields{
			Source:    "openline",
			AppSource: "unit-test",
		},
		AuthorUserId:     utils.StringPtr("123"),
		CommentedIssueId: utils.StringPtr("456"),
	})
	require.Nil(t, err, "Failed to create log entry")

	require.NotNil(t, response)
	commentId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	commentAggregate := aggregate.NewCommentAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[commentAggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, event.CommentCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(aggregate.CommentAggregateType)+"-"+tenant+"-"+commentId, eventList[0].GetAggregateID())

	var eventData event.CommentCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "unit-test", eventData.AppSource)
	require.Equal(t, "openline", eventData.Source)
	require.Equal(t, timeNow, eventData.CreatedAt)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, "This is a log entry", eventData.Content)
	require.Equal(t, "text/plain", eventData.ContentType)
	require.Equal(t, "123", eventData.AuthorUserId)
	require.Equal(t, "456", eventData.CommentedIssueId)
}

func TestCommentService_UpsertComment_UpdateComment(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)

	commentClient := commentpb.NewCommentGrpcServiceClient(grpcConnection)

	timeNow := utils.Now()
	commentId := uuid.New().String()
	tenant := "ziggy"

	// prepare aggregate
	commentAggregate := aggregate.NewCommentAggregateWithTenantAndID(tenant, commentId)
	updateEvent := eventstore.NewBaseEvent(commentAggregate, event.CommentCreateV1)
	preconfiguredEventData := event.CommentCreateEvent{
		Source: "openline",
	}
	err = updateEvent.SetJsonData(&preconfiguredEventData)
	require.Nil(t, err)
	commentAggregate.UncommittedEvents = []eventstore.Event{
		updateEvent,
	}
	err = aggregateStore.Save(ctx, commentAggregate)
	require.Nil(t, err)

	response, err := commentClient.UpsertComment(ctx, &commentpb.UpsertCommentGrpcRequest{
		Tenant:           tenant,
		Id:               commentId,
		Content:          "This is a log entry",
		ContentType:      "text/plain",
		UpdatedAt:        timestamppb.New(timeNow),
		AuthorUserId:     utils.StringPtr("123"),
		CommentedIssueId: utils.StringPtr("456"),
	})
	require.Nil(t, err, "Failed to create log entry")

	require.NotNil(t, response)
	require.Equal(t, commentId, response.Id)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[aggregate.NewCommentAggregateWithTenantAndID(tenant, response.Id).ID]
	require.Equal(t, 2, len(eventList))

	require.Equal(t, event.CommentCreateV1, eventList[0].GetEventType())
	require.Equal(t, event.CommentUpdateV1, eventList[1].GetEventType())
	require.Equal(t, string(aggregate.CommentAggregateType)+"-"+tenant+"-"+commentId, eventList[1].GetAggregateID())

	var eventData event.CommentUpdateEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "openline", eventData.Source)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, "This is a log entry", eventData.Content)
	require.Equal(t, "text/plain", eventData.ContentType)
}
