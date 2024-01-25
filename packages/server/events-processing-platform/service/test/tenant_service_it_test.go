package servicet

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	tenant "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTenantService_AddBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenantName := "ziggy"
	now := utils.Now()

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to get grpc connection")
	tenantServiceClient := tenantpb.NewTenantGrpcServiceClient(grpcConnection)

	response, err := tenantServiceClient.AddBillingProfile(ctx, &tenantpb.AddBillingProfileRequest{
		Tenant: tenantName,
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
			Source:    "source",
		},
		CreatedAt:                         utils.ConvertTimeToTimestampPtr(&now),
		Email:                             "email",
		Phone:                             "phone",
		AddressLine1:                      "addressLine1",
		AddressLine2:                      "addressLine2",
		AddressLine3:                      "addressLine3",
		Locality:                          "locality",
		Country:                           "country",
		Zip:                               "zip",
		LegalName:                         "legalName",
		DomesticPaymentsBankInfo:          "domesticPaymentsBankInfo",
		DomesticPaymentsBankName:          "domesticPaymentsBankName",
		DomesticPaymentsAccountNumber:     "domesticPaymentsAccountNumber",
		DomesticPaymentsSortCode:          "domesticPaymentsSortCode",
		InternationalPaymentsBankInfo:     "internationalPaymentsBankInfo",
		InternationalPaymentsSwiftBic:     "internationalPaymentsSwiftBic",
		InternationalPaymentsBankName:     "internationalPaymentsBankName",
		InternationalPaymentsBankAddress:  "internationalPaymentsBankAddress",
		InternationalPaymentsInstructions: "internationalPaymentsInstructions",
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	billingProfileId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	tenantAggregate := tenant.NewTenantAggregate(tenantName)
	eventList := eventsMap[tenantAggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, event.TenantAddBillingProfileV1, eventList[0].GetEventType())
	require.Equal(t, string(tenant.TenantAggregateType)+"-"+tenantName, eventList[0].GetAggregateID())

	var eventData event.CreateTenantBillingProfileEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, tenantName, eventData.Tenant)
	require.Equal(t, now, eventData.CreatedAt)
	require.Equal(t, billingProfileId, eventData.Id)
	require.Equal(t, "email", eventData.Email)
	require.Equal(t, "phone", eventData.Phone)
	require.Equal(t, "addressLine1", eventData.AddressLine1)
	require.Equal(t, "addressLine2", eventData.AddressLine2)
	require.Equal(t, "addressLine3", eventData.AddressLine3)
	require.Equal(t, "locality", eventData.Locality)
	require.Equal(t, "country", eventData.Country)
	require.Equal(t, "zip", eventData.Zip)
	require.Equal(t, "legalName", eventData.LegalName)
	require.Equal(t, "domesticPaymentsBankInfo", eventData.DomesticPaymentsBankInfo)
	require.Equal(t, "domesticPaymentsBankName", eventData.DomesticPaymentsBankName)
	require.Equal(t, "domesticPaymentsAccountNumber", eventData.DomesticPaymentsAccountNumber)
	require.Equal(t, "domesticPaymentsSortCode", eventData.DomesticPaymentsSortCode)
	require.Equal(t, "internationalPaymentsBankInfo", eventData.InternationalPaymentsBankInfo)
	require.Equal(t, "internationalPaymentsSwiftBic", eventData.InternationalPaymentsSwiftBic)
	require.Equal(t, "internationalPaymentsBankName", eventData.InternationalPaymentsBankName)
	require.Equal(t, "internationalPaymentsBankAddress", eventData.InternationalPaymentsBankAddress)
	require.Equal(t, "internationalPaymentsInstructions", eventData.InternationalPaymentsInstructions)
	require.Equal(t, "app", eventData.SourceFields.AppSource)
	require.Equal(t, "source", eventData.SourceFields.Source)
	require.Equal(t, "source", eventData.SourceFields.SourceOfTruth)
}
