package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	cmngrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	issuegrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/issue"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestIssueService_UpsertIssue_CreateIssue(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(nil, aggregateStore)
	require.Nil(t, err)

	issueClient := issuegrpc.NewIssueGrpcServiceClient(grpcConnection)
	timeNow := utils.Now()
	tenant := "ziggy"
	response, err := issueClient.UpsertIssue(ctx, &issuegrpc.UpsertIssueGrpcRequest{
		Tenant:      tenant,
		Subject:     "This is subject",
		Description: "This is description",
		Status:      "open",
		Priority:    "high",
		CreatedAt:   timestamppb.New(timeNow),
		SourceFields: &cmngrpc.SourceFields{
			Source:    "openline",
			AppSource: "unit-test",
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
	require.Equal(t, "unit-test", eventData.AppSource)
	require.Equal(t, "openline", eventData.Source)
	require.Equal(t, timeNow, eventData.CreatedAt)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, "This is subject", eventData.Subject)
	require.Equal(t, "This is description", eventData.Description)
	require.Equal(t, "open", eventData.Status)
	require.Equal(t, "high", eventData.Priority)
	require.Equal(t, "456", eventData.ReportedByOrganizationId)
	require.Equal(t, "ABC", eventData.SubmittedByOrganizationId)
	require.Equal(t, "DEF", eventData.SubmittedByUserId)
}

func TestIssueService_UpsertIssue_UpdateIssue(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(nil, aggregateStore)
	require.Nil(t, err)

	issueClient := issuegrpc.NewIssueGrpcServiceClient(grpcConnection)

	timeNow := utils.Now()
	issueId := uuid.New().String()
	tenant := "ziggy"

	// prepare aggregate
	issueAggregate := aggregate.NewIssueAggregateWithTenantAndID(tenant, issueId)
	createEvent := eventstore.NewBaseEvent(issueAggregate, event.IssueCreateV1)
	preconfiguredEventData := event.IssueCreateEvent{
		Tenant:  tenant,
		Subject: "This is subject",
	}
	err = createEvent.SetJsonData(&preconfiguredEventData)
	require.Nil(t, err)
	issueAggregate.UncommittedEvents = []eventstore.Event{
		createEvent,
	}
	err = aggregateStore.Save(ctx, issueAggregate)
	require.Nil(t, err)

	response, err := issueClient.UpsertIssue(ctx, &issuegrpc.UpsertIssueGrpcRequest{
		Tenant:      tenant,
		Id:          issueId,
		Subject:     "New subject",
		Description: "New description",
		Status:      "closed",
		Priority:    "low",
		SourceFields: &cmngrpc.SourceFields{
			Source:    "openline",
			AppSource: "unit-test",
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
	require.Equal(t, "New subject", eventData.Subject)
	require.Equal(t, "New description", eventData.Description)
	require.Equal(t, "closed", eventData.Status)
	require.Equal(t, "low", eventData.Priority)
}
