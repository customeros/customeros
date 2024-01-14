package servicet

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/currency"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	currencypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/currency"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInvoiceService_CreateCurrency(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to get grpc connection")
	currencyServiceClient := currencypb.NewCurrencyGrpcServiceClient(grpcConnection)

	response, err := currencyServiceClient.CreateCurrency(ctx, &currencypb.CreateCurrencyRequest{
		Name:           "USD",
		Symbol:         "$",
		LoggedInUserId: "user",
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
			Source:    "source",
		},
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	currencyId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	invoiceAggregate := currency.NewCurrencyAggregateWithID(response.Id)
	eventList := eventsMap[invoiceAggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, currency.CurrencyCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(currency.CurrencyAggregateType)+"-"+currencyId, eventList[0].GetAggregateID())

	var eventData currency.CurrencyCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, "USD", eventData.Name)
	require.Equal(t, "$", eventData.Symbol)
	test.AssertRecentTime(t, eventData.CreatedAt)
	require.Equal(t, "app", eventData.SourceFields.AppSource)
	require.Equal(t, "source", eventData.SourceFields.Source)
	require.Equal(t, "source", eventData.SourceFields.SourceOfTruth)
}
