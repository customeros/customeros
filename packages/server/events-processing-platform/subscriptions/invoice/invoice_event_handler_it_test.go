package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/mocked_grpc"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func prepareInvoiceEventHandler() *InvoiceEventHandler {
	return &InvoiceEventHandler{
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
		log:          testLogger,
	}
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_AnnualPrice_MonthlyInvoicing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     1200,
		Billed:    neo4jenum.BilledTypeAnnually,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleMonthlyBilling,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(100), inv.Amount)
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_AnnualPrice_QuarterlyInvoicing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     1200,
		Billed:    neo4jenum.BilledTypeAnnually,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleQuarterlyBilling,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(300), inv.Amount)
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_AnnualPrice_AnnualInvoicing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  2,
		Price:     1200,
		Billed:    neo4jenum.BilledTypeAnnually,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleAnnuallyBilling,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(2400), inv.Amount)
			require.Equal(t, float64(2400), inv.Total)
			require.Equal(t, float64(0), inv.Vat)
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_QuarterlyPrice_MonthlyInvoicing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     300,
		Billed:    neo4jenum.BilledTypeQuarterly,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleMonthlyBilling,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(100), inv.Amount)
			require.Equal(t, float64(100), inv.Total)
			require.Equal(t, float64(0), inv.Vat)
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_QuarterlyPrice_QuarterlyInvoicing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     300,
		Billed:    neo4jenum.BilledTypeQuarterly,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleQuarterlyBilling,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(300), inv.Amount)
			require.Equal(t, float64(300), inv.Total)
			require.Equal(t, float64(0), inv.Vat)
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_QuarterlyPrice_AnnualInvoicing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     300,
		Billed:    neo4jenum.BilledTypeQuarterly,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleAnnuallyBilling,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(1200), inv.Amount)
			require.Equal(t, float64(1200), inv.Total)
			require.Equal(t, float64(0), inv.Vat)
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_MonthlyPrice_MonthlyInvoicing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     100,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleMonthlyBilling,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(100), inv.Amount)
			require.Equal(t, float64(100), inv.Total)
			require.Equal(t, float64(0), inv.Vat)
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_MonthlyPrice_QuarterlyInvoicing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     100,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleQuarterlyBilling,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(300), inv.Amount)
			require.Equal(t, float64(300), inv.Total)
			require.Equal(t, float64(0), inv.Vat)
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_MonthlyPrice_AnnualInvoicing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     100,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleAnnuallyBilling,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(1200), inv.Amount)
			require.Equal(t, float64(1200), inv.Total)
			require.Equal(t, float64(0), inv.Vat)
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_MultipleServiceLineItems(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     100,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod.AddDate(0, 0, -1),
		Quantity:  2,
		Price:     10,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  3,
		Price:     3,
		Billed:    neo4jenum.BilledTypeOnce,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleMonthlyBilling,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(129), inv.Amount)
			require.Equal(t, float64(129), inv.Total)
			require.Equal(t, float64(0), inv.Vat)
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_Postpaid_MultipleServiceLineItems(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)
	secondsBeforeInvoiceEndPeriodEOD := invoiceEndPeriod.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceEndPeriod,
		Quantity:  1,
		Price:     100,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod.AddDate(0, 0, 10),
		Quantity:  2,
		Price:     10,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: secondsBeforeInvoiceEndPeriodEOD,
		Quantity:  3,
		Price:     3,
		Billed:    neo4jenum.BilledTypeOnce,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleMonthlyBilling,
		Postpaid:        true,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(129), inv.Total)
			require.Equal(t, 3, len(inv.InvoiceLines))
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_Postpaid_SkipUsageSLIs(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     100,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  9999,
		Price:     9999,
		Billed:    neo4jenum.BilledTypeUsage,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  9999,
		Price:     9999,
		Billed:    neo4jenum.BilledTypeNone,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleMonthlyBilling,
		Postpaid:        true,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(100), inv.Total)
			require.Equal(t, 1, len(inv.InvoiceLines))
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_Postpaid_SliEndedBeforeEndOfInvoicingCycle(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     100,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		EndedAt:   &invoiceEndPeriod,
		Quantity:  9999,
		Price:     9999,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleMonthlyBilling,
		Postpaid:        true,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(100), inv.Total)
			require.Equal(t, 1, len(inv.InvoiceLines))
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_Prepaid_SliEndedBeforeStartOfInvoicingCycle(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     100,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod.AddDate(0, 0, -10),
		EndedAt:   utils.ToPtr(invoiceStartPeriod.Add(-1 * time.Second)),
		Quantity:  9999,
		Price:     9999,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleMonthlyBilling,
		Postpaid:        false,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(100), inv.Total)
			require.Equal(t, 1, len(inv.InvoiceLines))
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_Prepaid_SliNotActiveOnStartOfInvoicingCycle(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     100,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod.Add(1 * time.Minute),
		EndedAt:   utils.ToPtr(invoiceStartPeriod.Add(10 * time.Minute)),
		Quantity:  9999,
		Price:     9999,
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleMonthlyBilling,
		Postpaid:        false,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(100), inv.Total)
			require.Equal(t, 1, len(inv.InvoiceLines))
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_TaxRate_Provided(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     100,
		VatRate:   float64(20),
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  2,
		Price:     100,
		VatRate:   float64(10),
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleMonthlyBilling,
		Postpaid:        false,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, float64(340), inv.Total)
			require.Equal(t, float64(300), inv.Amount)
			require.Equal(t, float64(40), inv.Vat)
			require.Equal(t, 2, len(inv.InvoiceLines))
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}

func TestInvoiceEventHandler_OnInvoiceFillRequestedV1_CycleInvoice_Check2DecimalsRounding(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// test data
	invoiceStartPeriod := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	invoiceEndPeriod := time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		StartedAt: invoiceStartPeriod,
		Quantity:  1,
		Price:     0.33333,
		VatRate:   float64(10),
		Billed:    neo4jenum.BilledTypeMonthly,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		PeriodStartDate: invoiceStartPeriod,
		PeriodEndDate:   invoiceEndPeriod,
		BillingCycle:    neo4jenum.BillingCycleMonthlyBilling,
		Postpaid:        false,
	})

	// prepare event handler
	invoiceEventHandler := prepareInvoiceEventHandler()
	// prepare aggregate
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	// prepare event
	invoiceFillRequestedEvent, _ := invoice.NewInvoiceFillRequestedEvent(invoiceAggregate, contractId)

	// prepare grpc mock
	calledFillInvoice := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		FillInvoice: func(context context.Context, inv *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, 0.36, inv.Total)
			require.Equal(t, 0.33, inv.Amount)
			require.Equal(t, 0.03, inv.Vat)
			require.Equal(t, 1, len(inv.InvoiceLines))
			calledFillInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// EXECUTE
	err := invoiceEventHandler.onInvoiceFillRequestedV1(context.Background(), invoiceFillRequestedEvent)
	require.Nil(t, err, "invoicing failed")

	// VERIFY
	require.True(t, calledFillInvoice)
}
