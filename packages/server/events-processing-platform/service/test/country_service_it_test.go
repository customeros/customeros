package servicet

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/country"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	countrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/country"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCountryService_CreateCountry(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to get grpc connection")
	countryServiceClient := countrypb.NewCountryGrpcServiceClient(grpcConnection)

	response, err := countryServiceClient.CreateCountry(ctx, &countrypb.CreateCountryRequest{
		Name:           "ABC",
		CodeA2:         "1",
		CodeA3:         "2",
		PhoneCode:      "3",
		LoggedInUserId: "user",
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
			Source:    "source",
		},
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	countryId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	invoiceAggregate := country.NewCountryAggregateWithID(response.Id)
	eventList := eventsMap[invoiceAggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, country.CountryCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(country.CountryAggregateType)+"-"+countryId, eventList[0].GetAggregateID())

	var eventData country.CountryCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, "ABC", eventData.Name)
	require.Equal(t, "1", eventData.CodeA2)
	require.Equal(t, "2", eventData.CodeA3)
	require.Equal(t, "3", eventData.PhoneCode)
	test.AssertRecentTime(t, eventData.CreatedAt)
	require.Equal(t, "app", eventData.SourceFields.AppSource)
	require.Equal(t, "source", eventData.SourceFields.Source)
	require.Equal(t, "source", eventData.SourceFields.SourceOfTruth)
}
