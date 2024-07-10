package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	issuepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/issue"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestIssueService_UpsertIssue_CreateIssue(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)

	issueClient := issuepb.NewIssueGrpcServiceClient(grpcConnection)
	timeNow := utils.Now()
	tenant := "ziggy"
	response, err := issueClient.UpsertIssue(ctx, &issuepb.UpsertIssueGrpcRequest{
		Tenant:      tenant,
		GroupId:     &tenant,
		Subject:     "This is subject",
		Description: "This is description",
		Status:      "open",
		Priority:    "high",
		CreatedAt:   timestamppb.New(timeNow),
		SourceFields: &commonpb.SourceFields{
			Source:    "openline",
			AppSource: "event-processing-platform",
		},
		ReportedByOrganizationId:  utils.StringPtr("456"),
		SubmittedByOrganizationId: utils.StringPtr("ABC"),
		SubmittedByUserId:         utils.StringPtr("DEF"),
	})
	require.Nil(t, err, "Failed to create issue")
	require.NotNil(t, response)

	issueId := response.Id

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	issueAggregate := aggregate.NewIssueAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[issueAggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, event.IssueCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(aggregate.IssueAggregateType)+"-"+tenant+"-"+issueId, eventList[0].GetAggregateID())

	var eventData event.IssueCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err)

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "event-processing-platform", eventData.AppSource)
	require.Equal(t, "openline", eventData.Source)
	require.Equal(t, timeNow, eventData.CreatedAt)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, "This is subject", eventData.Subject)
	require.Equal(t, "This is description", eventData.Description)
	require.Equal(t, tenant, eventData.GroupId)
	require.Equal(t, "open", eventData.Status)
	require.Equal(t, "high", eventData.Priority)
	require.Equal(t, "456", eventData.ReportedByOrganizationId)
	require.Equal(t, "ABC", eventData.SubmittedByOrganizationId)
	require.Equal(t, "DEF", eventData.SubmittedByUserId)
}

func TestIssueService_UpsertIssue_UpdateIssue(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)

	issueClient := issuepb.NewIssueGrpcServiceClient(grpcConnection)

	timeNow := utils.Now()
	issueId := uuid.New().String()
	tenant := "ziggy"

	// prepare aggregate
	issueAggregate := aggregate.NewIssueAggregateWithTenantAndID(tenant, issueId)
	createEvent := eventstore.NewBaseEvent(issueAggregate, event.IssueCreateV1)
	preconfiguredEventData := event.IssueCreateEvent{
		Tenant:  tenant,
		GroupId: "This is groupId",
		Subject: "This is subject",
	}
	err = createEvent.SetJsonData(&preconfiguredEventData)
	require.Nil(t, err)
	issueAggregate.UncommittedEvents = []eventstore.Event{
		createEvent,
	}
	err = aggregateStore.Save(ctx, issueAggregate)
	require.Nil(t, err)

	response, err := issueClient.UpsertIssue(ctx, &issuepb.UpsertIssueGrpcRequest{
		Tenant:      tenant,
		Id:          issueId,
		GroupId:     &tenant,
		Subject:     "New subject",
		Description: "New description",
		Status:      "closed",
		Priority:    "low",
		SourceFields: &commonpb.SourceFields{
			Source:    "openline",
			AppSource: "event-processing-platform",
		},
		UpdatedAt: timestamppb.New(timeNow),
	})
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, issueId, response.Id)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[aggregate.NewIssueAggregateWithTenantAndID(tenant, response.Id).ID]
	require.Equal(t, 2, len(eventList))

	require.Equal(t, event.IssueCreateV1, eventList[0].GetEventType())
	require.Equal(t, event.IssueUpdateV1, eventList[1].GetEventType())
	require.Equal(t, string(aggregate.IssueAggregateType)+"-"+tenant+"-"+issueId, eventList[1].GetAggregateID())

	var eventData event.IssueUpdateEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err)

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "openline", eventData.Source)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, tenant, eventData.GroupId)
	require.Equal(t, "New subject", eventData.Subject)
	require.Equal(t, "New description", eventData.Description)
	require.Equal(t, "closed", eventData.Status)
	require.Equal(t, "low", eventData.Priority)
}
