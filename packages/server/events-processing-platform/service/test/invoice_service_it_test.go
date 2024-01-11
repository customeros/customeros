package servicet

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInvoiceService_NewInvoice(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	now := utils.Now()

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to get grpc connection")
	invoiceServiceClient := invoicepb.NewInvoiceServiceClient(grpcConnection)

	response, err := invoiceServiceClient.NewInvoice(ctx, &invoicepb.NewInvoiceRequest{
		Tenant:         tenant,
		OrganizationId: "1",
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
			Source:    "source",
		},
		CreatedAt: utils.ConvertTimeToTimestampPtr(&now),
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	invoiceId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[invoiceAggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, invoice.InvoiceNewV1, eventList[0].GetEventType())
	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[0].GetAggregateID())

	var eventData invoice.InvoiceNewEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "1", eventData.OrganizationId)
	test.AssertRecentTime(t, eventData.CreatedAt)
	require.Equal(t, "app", eventData.SourceFields.AppSource)
	require.Equal(t, "source", eventData.SourceFields.Source)
	require.Equal(t, "source", eventData.SourceFields.SourceOfTruth)
}
