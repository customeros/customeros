package servicet

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	orgaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	utils2 "github.com/openline-ai/openline-customer-os/packages/server/events/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestContractService_CreateContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	orgId := "Org123"

	aggregateStore := eventstoret.NewTestAggregateStore()
	organizationAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, orgId)
	aggregateStore.Save(ctx, organizationAggregate)
	grpcConnection, _ := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	contractClient := contractpb.NewContractGrpcServiceClient(grpcConnection)
	timeNow := utils.Now()
	response, err := contractClient.CreateContract(ctx, &contractpb.CreateContractGrpcRequest{
		Tenant:               tenant,
		Name:                 "New Contract",
		ContractUrl:          "http://contract.url",
		OrganizationId:       orgId,
		CreatedByUserId:      "User123",
		ServiceStartedAt:     timestamppb.New(timeNow),
		SignedAt:             timestamppb.New(timeNow),
		LengthInMonths:       int64(1),
		BillingCycleInMonths: int64(2),
		Approved:             true,
		AutoRenew:            true,
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
	require.Equal(t, "http://contract.url", eventData.ContractUrl)
	require.Equal(t, orgId, eventData.OrganizationId)
	require.Equal(t, "User123", eventData.CreatedByUserId)
	require.Equal(t, utils.ToDate(timeNow), *eventData.ServiceStartedAt)
	require.Equal(t, utils.ToDate(timeNow), *eventData.SignedAt)
	require.Equal(t, int64(1), eventData.LengthInMonths)
	require.Equal(t, int64(2), eventData.BillingCycleInMonths)
	require.Equal(t, "ExternalSystemID", eventData.ExternalSystem.ExternalSystemId)
	require.Equal(t, "http://external.url", eventData.ExternalSystem.ExternalUrl)
	require.Equal(t, "ExternalID", eventData.ExternalSystem.ExternalId)
	require.Equal(t, "ExternalIDSecond", eventData.ExternalSystem.ExternalIdSecond)
	require.Equal(t, "ExternalSource", eventData.ExternalSystem.ExternalSource)
	require.True(t, timeNow.Equal(*eventData.ExternalSystem.SyncDate))
	require.True(t, eventData.Approved)
	require.True(t, eventData.AutoRenew)
}

func TestContractService_CreateContract_ServiceStartedInFuture(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup
	tenant := "ziggy"
	orgId := "Org123"

	aggregateStore := eventstoret.NewTestAggregateStore()
	organizationAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, orgId)
	aggregateStore.Save(ctx, organizationAggregate)
	grpcConnection, _ := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	contractClient := contractpb.NewContractGrpcServiceClient(grpcConnection)

	// Create a future date
	futureDate := utils.Now().AddDate(0, 1, 0) // 1 month into the future

	// Call CreateContract with future ServiceStartedAt
	response, err := contractClient.CreateContract(ctx, &contractpb.CreateContractGrpcRequest{
		Tenant:           tenant,
		OrganizationId:   orgId,
		ServiceStartedAt: timestamppb.New(futureDate),
	})
	require.Nil(t, err, "Failed to create contract with future ServiceStartedAt")

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 2, len(eventsMap))
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[contractAggregate.ID]
	require.Equal(t, 1, len(eventList))
	var eventData event.ContractCreateEvent
	err = eventList[0].GetJsonData(&eventData)
}

func TestContractService_CreateContract_ServiceStartedNil(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup
	tenant := "ziggy"
	orgId := "Org123"

	aggregateStore := eventstoret.NewTestAggregateStore()
	organizationAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, orgId)
	aggregateStore.Save(ctx, organizationAggregate)
	grpcConnection, _ := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	contractClient := contractpb.NewContractGrpcServiceClient(grpcConnection)

	// Call CreateContract with future ServiceStartedAt
	response, err := contractClient.CreateContract(ctx, &contractpb.CreateContractGrpcRequest{
		Tenant:         tenant,
		OrganizationId: orgId,
	})
	require.Nil(t, err, "Failed to create contract with future ServiceStartedAt")

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 2, len(eventsMap))
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[contractAggregate.ID]
	require.Equal(t, 1, len(eventList))
	var eventData event.ContractCreateEvent
	err = eventList[0].GetJsonData(&eventData)
}

func TestContractService_CreateContract_MissingOrganizationId(t *testing.T) {
	ctx := context.Background()
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

func TestContractService_CreateContract_OrganizationAggregateDoesNotExists(t *testing.T) {
	ctx := context.Background()
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

func TestContractService_UpdateContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	tenant := "ziggy"
	contractId := uuid.New().String()

	// Setup aggregate store and create initial contract
	aggregateStore := eventstoret.NewTestAggregateStore()
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenant, contractId)
	aggregateStore.Save(ctx, contractAggregate)

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to connect to processing platform")
	contractClient := contractpb.NewContractGrpcServiceClient(grpcConnection)

	// Create update request
	timeNow := utils.Now()
	updateRequest := &contractpb.UpdateContractGrpcRequest{
		Tenant:               tenant,
		Id:                   contractId,
		Name:                 "Updated Contract",
		ContractUrl:          "http://new.contract.url",
		UpdatedAt:            timestamppb.New(timeNow),
		ServiceStartedAt:     timestamppb.New(timeNow),
		SignedAt:             timestamppb.New(timeNow),
		EndedAt:              timestamppb.New(timeNow.AddDate(0, 1, 0)),
		NextInvoiceDate:      timestamppb.New(timeNow),
		LengthInMonths:       int64(1),
		BillingCycleInMonths: int64(2),
		Currency:             "USD",
		FieldsMask: []contractpb.ContractFieldMask{
			contractpb.ContractFieldMask_CONTRACT_FIELD_NAME,
			contractpb.ContractFieldMask_CONTRACT_FIELD_CONTRACT_URL,
			contractpb.ContractFieldMask_CONTRACT_FIELD_LENGTH_IN_MONTHS,
			contractpb.ContractFieldMask_CONTRACT_FIELD_BILLING_CYCLE_IN_MONTHS,
			contractpb.ContractFieldMask_CONTRACT_FIELD_SERVICE_STARTED_AT,
			contractpb.ContractFieldMask_CONTRACT_FIELD_SIGNED_AT,
			contractpb.ContractFieldMask_CONTRACT_FIELD_ENDED_AT,
			contractpb.ContractFieldMask_CONTRACT_FIELD_CURRENCY,
		},
		SourceFields: &commonpb.SourceFields{
			Source:    utils2.SourceOpenline,
			AppSource: "event-processing-platform",
		},
		ExternalSystemFields: &commonpb.ExternalSystemFields{
			ExternalSystemId: "ExternalSystemID",
			ExternalUrl:      "http://external.url",
			ExternalId:       "ExternalID",
			ExternalIdSecond: "ExternalIDSecond",
			ExternalSource:   "ExternalSource",
			SyncDate:         timestamppb.New(timeNow),
		},
	}

	// Execute update contract request
	response, err := contractClient.UpdateContract(ctx, updateRequest)
	require.Nil(t, err, "Failed to update contract")

	// Assert response
	require.NotNil(t, response)
	require.Equal(t, contractId, response.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	contractEvents := eventsMap[contractAggregate.ID]
	require.Equal(t, 1, len(contractEvents))

	require.Equal(t, event.ContractUpdateV1, contractEvents[0].GetEventType())

	var eventData event.ContractUpdateEvent
	err = contractEvents[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assert event data
	require.Equal(t, "Updated Contract", eventData.Name)
	require.Equal(t, "http://new.contract.url", eventData.ContractUrl)
	require.Equal(t, int64(1), eventData.LengthInMonths)
	require.Equal(t, int64(2), eventData.BillingCycleInMonths)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, utils.ToDate(timeNow), *eventData.ServiceStartedAt)
	require.Equal(t, utils.ToDate(timeNow), *eventData.SignedAt)
	require.Equal(t, utils.ToDate(timeNow).AddDate(0, 1, 0), *eventData.EndedAt)
	require.Equal(t, utils2.SourceOpenline, eventData.Source)
	require.Equal(t, "ExternalSystemID", eventData.ExternalSystem.ExternalSystemId)
	require.Equal(t, "USD", eventData.Currency)
	require.Nil(t, eventData.NextInvoiceDate) // next invoice date was not mentioned in fields mask, hence it should be nil
}

func TestContractService_UpdateContract_OnlySelectedFieldsModified(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	tenant := "ziggy"
	contractId := uuid.New().String()

	// Setup aggregate store and create initial contract
	aggregateStore := eventstoret.NewTestAggregateStore()
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenant, contractId)
	aggregateStore.Save(ctx, contractAggregate)

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to connect to processing platform")
	contractClient := contractpb.NewContractGrpcServiceClient(grpcConnection)

	// Create update request
	updateRequest := &contractpb.UpdateContractGrpcRequest{
		Tenant:          tenant,
		Id:              contractId,
		InvoiceEmailTo:  "to@gmail.com",
		InvoiceEmailCc:  []string{"cc1@gmail.com", "cc2@gmail.com"},
		InvoiceEmailBcc: []string{"bcc1@gmail.com", "bcc2@gmail.com"},
		SourceFields: &commonpb.SourceFields{
			Source:    utils2.SourceOpenline,
			AppSource: "event-processing-platform",
		},
		FieldsMask: []contractpb.ContractFieldMask{
			contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_EMAIL_TO,
			contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_EMAIL_CC,
			contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_EMAIL_BCC,
		},
	}

	// Execute update contract request
	response, err := contractClient.UpdateContract(ctx, updateRequest)
	require.Nil(t, err, "Failed to update contract")

	// Assert response
	require.NotNil(t, response)
	require.Equal(t, contractId, response.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	contractEvents := eventsMap[contractAggregate.ID]
	require.Equal(t, 1, len(contractEvents))

	require.Equal(t, event.ContractUpdateV1, contractEvents[0].GetEventType())

	var eventData event.ContractUpdateEvent
	err = contractEvents[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assert event data
	require.Equal(t, []string{
		event.FieldMaskInvoiceEmail,
		event.FieldMaskInvoiceEmailCC,
		event.FieldMaskInvoiceEmailBCC}, eventData.FieldsMask)
	require.Equal(t, "to@gmail.com", eventData.InvoiceEmail)
	require.Equal(t, []string{"cc1@gmail.com", "cc2@gmail.com"}, eventData.InvoiceEmailCC)
	require.Equal(t, []string{"bcc1@gmail.com", "bcc2@gmail.com"}, eventData.InvoiceEmailBCC)
	require.Nil(t, eventData.NextInvoiceDate) // next invoice date was not mentioned in fields mask, hence it should be nil
}

func TestContractService_SoftDeleteContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	tenant := "ziggy"
	contractId := uuid.New().String()

	// Setup aggregate store and create initial contract
	aggregateStore := eventstoret.NewTestAggregateStore()
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenant, contractId)
	aggregateStore.Save(ctx, contractAggregate)

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	contractClient := contractpb.NewContractGrpcServiceClient(grpcConnection)

	// Create update request
	deleteRequest := &contractpb.SoftDeleteContractGrpcRequest{
		Tenant: tenant,
		Id:     contractId,
	}

	// Execute update contract request
	response, err := contractClient.SoftDeleteContract(ctx, deleteRequest)
	require.Nil(t, err, "Failed to delete contract")

	// Assert response
	require.NotNil(t, response)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	contractEvents := eventsMap[contractAggregate.ID]
	require.Equal(t, 1, len(contractEvents))

	require.Equal(t, event.ContractDeleteV1, contractEvents[0].GetEventType())

	var eventData event.ContractDeleteEvent
	err = contractEvents[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assert event data
	require.NotNil(t, eventData.UpdatedAt)
	require.Equal(t, tenant, eventData.Tenant)
}
