package servicet

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	event "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	orgaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/opportunity"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestOpportunityService_CreateOpportunity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	orgId := "Org123"

	aggregateStore := eventstoret.NewTestAggregateStore()
	organizationAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, orgId)
	aggregateStore.Save(ctx, organizationAggregate)

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to connect to processing platform")
	opportunityClient := opportunitypb.NewOpportunityGrpcServiceClient(grpcConnection)
	timeNow := utils.Now()
	response, err := opportunityClient.CreateOpportunity(ctx, &opportunitypb.CreateOpportunityGrpcRequest{
		Tenant:             tenant,
		Name:               "New Opportunity",
		MaxAmount:          10000,
		InternalType:       opportunitypb.OpportunityInternalType_NBO,
		ExternalType:       "TypeA",
		InternalStage:      opportunitypb.OpportunityInternalStage_OPEN,
		ExternalStage:      "Stage1",
		EstimatedCloseDate: timestamppb.New(timeNow),
		OwnerUserId:        "OwnerUser123",
		GeneralNotes:       "Some general notes",
		NextSteps:          "Next steps to be taken",
		SourceFields: &commonpb.SourceFields{
			Source:    "openline",
			AppSource: "event-processing-platform",
		},
		OrganizationId: orgId,
		CreatedAt:      timestamppb.New(timeNow),
		UpdatedAt:      timestamppb.New(timeNow),
		ExternalSystemFields: &commonpb.ExternalSystemFields{
			ExternalSystemId: "ExternalSystemID",
			ExternalUrl:      "http://external.url",
			ExternalId:       "ExternalID",
			ExternalIdSecond: "ExternalIDSecond",
			ExternalSource:   "ExternalSource",
			SyncDate:         timestamppb.New(timeNow),
		},
	})
	require.Nil(t, err, "Failed to create opportunity")

	require.NotNil(t, response)
	opportunityId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 2, len(eventsMap))
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[opportunityAggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, opportunityevent.OpportunityCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(aggregate.OpportunityAggregateType)+"-"+tenant+"-"+opportunityId, eventList[0].GetAggregateID())

	var eventData event.OpportunityCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "New Opportunity", eventData.Name)
	require.Equal(t, 10000.0, eventData.MaxAmount)
	require.Equal(t, neo4jenum.OpportunityInternalTypeNBO.String(), eventData.InternalType)
	require.Equal(t, "TypeA", eventData.ExternalType)
	require.Equal(t, neo4jenum.OpportunityInternalStageOpen.String(), eventData.InternalStage)
	require.Equal(t, "Stage1", eventData.ExternalStage)
	require.Equal(t, timeNow, eventData.EstimatedClosedAt.UTC())
	require.Equal(t, "OwnerUser123", eventData.OwnerUserId)
	require.Equal(t, "openline", eventData.Source.Source)
	require.Equal(t, "event-processing-platform", eventData.Source.AppSource)
	require.Equal(t, "Org123", eventData.OrganizationId)
	require.Equal(t, "Some general notes", eventData.GeneralNotes)
	require.Equal(t, "Next steps to be taken", eventData.NextSteps)
	require.True(t, timeNow.Equal(eventData.CreatedAt.UTC()))
	require.True(t, timeNow.Equal(eventData.UpdatedAt.UTC()))
	require.Equal(t, "ExternalSystemID", eventData.ExternalSystem.ExternalSystemId)
	require.Equal(t, "http://external.url", eventData.ExternalSystem.ExternalUrl)
	require.Equal(t, "ExternalID", eventData.ExternalSystem.ExternalId)
	require.Equal(t, "ExternalIDSecond", eventData.ExternalSystem.ExternalIdSecond)
	require.Equal(t, "ExternalSource", eventData.ExternalSystem.ExternalSource)
	require.True(t, timeNow.Equal(*eventData.ExternalSystem.SyncDate))
}

func TestOpportunityService_CreateOpportunity_MissingOrganizationId(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	orgId := ""

	aggregateStore := eventstoret.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to connect to processing platform")
	opportunityClient := opportunitypb.NewOpportunityGrpcServiceClient(grpcConnection)
	_, err = opportunityClient.CreateOpportunity(ctx, &opportunitypb.CreateOpportunityGrpcRequest{
		Tenant:         tenant,
		Name:           "New Opportunity",
		OrganizationId: orgId,
	})
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
	require.Contains(t, st.Message(), "missing required field: organizationId")
}

func TestOpportunityService_CreateOpportunity_OrganizationAggregateDoesNotExists(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	orgId := "org123"

	aggregateStore := eventstoret.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to connect to processing platform")

	opportunityClient := opportunitypb.NewOpportunityGrpcServiceClient(grpcConnection)
	_, err = opportunityClient.CreateOpportunity(ctx, &opportunitypb.CreateOpportunityGrpcRequest{
		Tenant:         tenant,
		Name:           "New Opportunity",
		OrganizationId: orgId,
	})
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
	require.Contains(t, st.Message(), fmt.Sprintf("organization with ID %s not found", orgId))
}

func TestOpportunityService_UpdateRenewalOpportunity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	opportunityId := uuid.New().String()

	// Setup opportunity store and create initial opportunity event
	aggregateStore := eventstoret.NewTestAggregateStore()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenant, opportunityId)
	createEvent, err := event.NewOpportunityCreateEvent(opportunityAggregate, model.OpportunityDataFields{
		OrganizationId: "org-123",
	}, commonmodel.Source{}, commonmodel.ExternalSystem{}, utils.Now(), utils.Now())
	require.Nil(t, err)
	opportunityAggregate.UncommittedEvents = append(opportunityAggregate.UncommittedEvents, createEvent)
	aggregateStore.Save(ctx, opportunityAggregate)

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	opportunityClient := opportunitypb.NewOpportunityGrpcServiceClient(grpcConnection)

	// Update opportunity request
	response, err := opportunityClient.UpdateRenewalOpportunity(ctx, &opportunitypb.UpdateRenewalOpportunityGrpcRequest{
		Tenant:            tenant,
		Id:                opportunityId,
		Amount:            10000,
		OwnerUserId:       "user-id-123",
		RenewalLikelihood: opportunitypb.RenewalLikelihood_MEDIUM_RENEWAL,
		RenewalApproved:   true,
		SourceFields: &commonpb.SourceFields{
			Source:    "openline",
			AppSource: "event-processing-platform",
		},
		LoggedInUserId: "user-id-123",
		Comments:       "Some comments",
	})
	require.Nil(t, err, "Failed to update opportunity")

	// Assert response
	require.NotNil(t, response)
	require.Equal(t, opportunityId, response.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[opportunityAggregate.ID]
	require.Equal(t, 2, len(eventList))

	require.Equal(t, opportunityevent.OpportunityUpdateRenewalV1, eventList[1].GetEventType())
	require.Equal(t, string(aggregate.OpportunityAggregateType)+"-"+tenant+"-"+opportunityId, eventList[0].GetAggregateID())

	var eventData event.OpportunityUpdateRenewalEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, 10000.0, eventData.Amount)
	test.AssertRecentTime(t, eventData.UpdatedAt)
	require.Equal(t, "user-id-123", eventData.UpdatedByUserId)
	require.Equal(t, "user-id-123", eventData.OwnerUserId)
	require.Equal(t, "Some comments", eventData.Comments)
	require.Equal(t, neo4jenum.RenewalLikelihoodMedium.String(), eventData.RenewalLikelihood)
	require.Equal(t, "openline", eventData.Source)
	require.True(t, eventData.RenewalApproved)
}

func TestOpportunityService_CloseLooseOpportunity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	opportunityId := uuid.New().String()

	// Setup opportunity store and create initial opportunity event
	aggregateStore := eventstoret.NewTestAggregateStore()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenant, opportunityId)
	createEvent, err := event.NewOpportunityCreateEvent(opportunityAggregate, model.OpportunityDataFields{
		OrganizationId: "org-123",
	}, commonmodel.Source{}, commonmodel.ExternalSystem{}, utils.Now(), utils.Now())
	require.Nil(t, err)
	opportunityAggregate.UncommittedEvents = append(opportunityAggregate.UncommittedEvents, createEvent)
	aggregateStore.Save(ctx, opportunityAggregate)

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	opportunityClient := opportunitypb.NewOpportunityGrpcServiceClient(grpcConnection)

	// Update opportunity request
	response, err := opportunityClient.CloseLooseOpportunity(ctx, &opportunitypb.CloseLooseOpportunityGrpcRequest{
		Tenant:         tenant,
		Id:             opportunityId,
		AppSource:      "event-processing-platform",
		LoggedInUserId: "user-id-123",
	})
	require.Nil(t, err, "Failed to close loose opportunity")

	// Assert response
	require.NotNil(t, response)
	require.Equal(t, opportunityId, response.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[opportunityAggregate.ID]
	require.Equal(t, 2, len(eventList))

	require.Equal(t, opportunityevent.OpportunityCloseLooseV1, eventList[1].GetEventType())
	require.Equal(t, string(aggregate.OpportunityAggregateType)+"-"+tenant+"-"+opportunityId, eventList[0].GetAggregateID())

	var eventData opportunityevent.OpportunityCloseLooseEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
	test.AssertRecentTime(t, eventData.ClosedAt)
}
