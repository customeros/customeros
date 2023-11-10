package servicet

import (
	"context"
	"fmt"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	orgaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestContractService_CreateContract(t *testing.T) {
	ctx := context.TODO()
	// setup test environment
	tenant := "ziggy"
	orgId := "Org123"

	aggregateStore := eventstoret.NewTestAggregateStore()
	organizationAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, orgId)
	aggregateStore.Save(ctx, organizationAggregate)

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to connect to processing platform")
	contractClient := contractpb.NewContractGrpcServiceClient(grpcConnection)
	timeNow := time.Now()
	response, err := contractClient.CreateContract(ctx, &contractpb.CreateContractGrpcRequest{
		Tenant:           tenant,
		Name:             "New Contract",
		OrganizationId:   orgId,
		CreatedByUserId:  "User123",
		ServiceStartedAt: timestamppb.New(timeNow),
		SignedAt:         timestamppb.New(timeNow),
		RenewalCycle:     contractpb.RenewalCycle_MONTHLY_RENEWAL,
		Status:           contractpb.ContractStatus_LIVE,
		ExternalSystemFields: &commonpb.ExternalSystemFields{
			ExternalSystemId: "ExternalSystemID",
			ExternalUrl:      "http://external.url",
			ExternalId:       "ExternalID",
			ExternalIdSecond: "ExternalIDSecond",
			ExternalSource:   "ExternalSource",
			SyncDate:         timestamppb.New(timeNow),
		},
	})
	require.Nil(t, err, "Failed to create contract")

	require.NotNil(t, response)
	contractId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 2, len(eventsMap))
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[contractAggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, event.ContractCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(aggregate.ContractAggregateType)+"-"+tenant+"-"+contractId, eventList[0].GetAggregateID())

	var eventData event.ContractCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "New Contract", eventData.Name)
	require.Equal(t, orgId, eventData.OrganizationId)
	require.Equal(t, "User123", eventData.CreatedByUserId)
	require.True(t, timeNow.Equal(*eventData.ServiceStartedAt))
	require.True(t, timeNow.Equal(*eventData.SignedAt))
	require.Equal(t, model.MonthlyRenewal.String(), eventData.RenewalCycle)
	require.Equal(t, model.Live.String(), eventData.Status)
	require.Equal(t, "ExternalSystemID", eventData.ExternalSystem.ExternalSystemId)
	require.Equal(t, "http://external.url", eventData.ExternalSystem.ExternalUrl)
	require.Equal(t, "ExternalID", eventData.ExternalSystem.ExternalId)
	require.Equal(t, "ExternalIDSecond", eventData.ExternalSystem.ExternalIdSecond)
	require.Equal(t, "ExternalSource", eventData.ExternalSystem.ExternalSource)
	require.True(t, timeNow.Equal(*eventData.ExternalSystem.SyncDate))
}

func TestCreateContract_MissingOrganizationId(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	orgId := ""

	aggregateStore := eventstoret.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to connect to processing platform")
	contractClient := contractpb.NewContractGrpcServiceClient(grpcConnection)
	_, err = contractClient.CreateContract(ctx, &contractpb.CreateContractGrpcRequest{
		Tenant:         tenant,
		Name:           "New Contract",
		OrganizationId: orgId,
	})
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
	require.Contains(t, st.Message(), "missing required field: organizationId")
}

func TestCreateContract_OrganizationAggregateDoesNotExists(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	orgId := "org123"

	aggregateStore := eventstoret.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to connect to processing platform")

	contractClient := contractpb.NewContractGrpcServiceClient(grpcConnection)
	_, err = contractClient.CreateContract(ctx, &contractpb.CreateContractGrpcRequest{
		Tenant:         tenant,
		Name:           "New Contract",
		OrganizationId: orgId,
	})
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
	require.Contains(t, st.Message(), fmt.Sprintf("organization with ID %s not found", orgId))
}
