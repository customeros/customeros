package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

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
			Source:    constants.SourceOpenline,
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
	require.Equal(t, constants.SourceOpenline, eventData.Source)
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
			Source:    constants.SourceOpenline,
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
