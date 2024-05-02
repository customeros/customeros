package resolver

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

// SLI on first nanosecond of the billing cycle is considered to be included in the cycle invoice
func TestMutationResolver_InvoiceSimulate_OnCycle_PostPaidFalse_1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	januaryFirst := utils.FirstTimeOfMonth(2024, 1)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{
		InvoicingPostpaid: false,
	})
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		BillingCycleInMonths:  1,
		InvoicingEnabled:      true,
		InvoicingStartDate:    &januaryFirst,
		OrganizationLegalName: "Test Organization",
	})

	rawResponse := callGraphQL(t, "invoice/simulate_invoice", map[string]interface{}{
		"contractId": contractId,
		"serviceLines": []model.InvoiceSimulateServiceLineInput{
			{
				Key:            "1",
				Description:    "S1",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          1,
				Quantity:       1,
				ServiceStarted: januaryFirst,
			},
		},
	})

	var invoiceStruct struct {
		Invoice_Simulate []*model.InvoiceSimulate
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, 1, len(invoiceStruct.Invoice_Simulate))

	invoice := invoiceStruct.Invoice_Simulate[0]

	require.Equal(t, 1, len(invoice.InvoiceLineItems))
	require.Equal(t, "Test Organization", *invoice.Customer.Name)

	asserInvoice(t, invoice, "2024-01-01T00:00:00Z", "2024-01-31T00:00:00Z", false, false, 1)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[0], "1", "S1", 1, 1, 1)
}

// 1 new SLI prorated
// TODO suppress test until offcycle invoices are re-enabled
//func TestMutationResolver_InvoiceSimulate_OnCycle_PostPaidFalse_2(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//
//	januaryFirst := utils.FirstTimeOfMonth(2024, 1)
//	januaryMid := utils.MiddleTimeOfMonth(2024, 1)
//
//	neo4jtest.CreateTenant(ctx, driver, tenantName)
//	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{
//		InvoicingPostpaid: false,
//	})
//	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
//	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
//	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
//		BillingCycleInMonths: 1,
//		InvoicingEnabled:     true,
//		InvoicingStartDate:   &januaryFirst,
//	})
//
//	rawResponse := callGraphQL(t, "invoice/simulate_invoice", map[string]interface{}{
//		"contractId": contractId,
//		"serviceLines": []model.InvoiceSimulateServiceLineInput{
//			{
//				Key:            "1",
//				Description:    "S1",
//				BillingCycle:   model.BilledTypeMonthly,
//				Price:          29,
//				Quantity:       1,
//				ServiceStarted: januaryMid,
//			},
//		},
//	})
//
//	var invoiceStruct struct {
//		Invoice_Simulate []*model.InvoiceSimulate
//	}
//
//	require.Nil(t, rawResponse.Errors)
//	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
//	require.Nil(t, err)
//
//	require.Equal(t, 2, len(invoiceStruct.Invoice_Simulate))
//
//	proratedInvoice := invoiceStruct.Invoice_Simulate[0]
//	require.Equal(t, 1, len(proratedInvoice.InvoiceLineItems))
//	asserInvoice(t, proratedInvoice, "2024-01-17T00:00:00Z", "2024-01-31T00:00:00Z", true, false, 14.3)
//	asserInvoiceLineItem(t, proratedInvoice.InvoiceLineItems[0], "1", "S1", 29, 1, 14.3)
//
//	onCycleInvoice := invoiceStruct.Invoice_Simulate[1]
//	require.Equal(t, 1, len(onCycleInvoice.InvoiceLineItems))
//	asserInvoice(t, onCycleInvoice, "2024-02-01T00:00:00Z", "2024-02-29T00:00:00Z", false, false, 29)
//	asserInvoiceLineItem(t, onCycleInvoice.InvoiceLineItems[0], "1", "S1", 29, 1, 29)
//}

// 1 upsell SLI prorated
// TODO suppress test until offcycle invoices are re-enabled
//func TestMutationResolver_InvoiceSimulate_OnCycle_PostPaidFalse_3(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//
//	decemberFirst := utils.FirstTimeOfMonth(2023, 12)
//	januaryMid := utils.MiddleTimeOfMonth(2024, 1)
//	februaryFirst := utils.FirstTimeOfMonth(2024, 2)
//
//	neo4jtest.CreateTenant(ctx, driver, tenantName)
//	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{
//		InvoicingPostpaid: false,
//	})
//	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
//	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
//	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
//		BillingCycleInMonths: 1,
//		InvoicingEnabled:     true,
//		InvoicingStartDate:   &decemberFirst,
//		NextInvoiceDate:      &februaryFirst,
//	})
//	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
//		Name:      "S1",
//		Billed:    neo4jenum.BilledTypeMonthly,
//		Price:     5,
//		Quantity:  1,
//		StartedAt: decemberFirst,
//	})
//
//	rawResponse := callGraphQL(t, "invoice/simulate_invoice", map[string]interface{}{
//		"contractId": contractId,
//		"serviceLines": []model.InvoiceSimulateServiceLineInput{
//			{
//				Key:               "1",
//				ServiceLineItemID: &serviceLineItemId,
//				ParentID:          &serviceLineItemId,
//				Description:       "S1",
//				BillingCycle:      model.BilledTypeMonthly,
//				Price:             29,
//				Quantity:          1,
//				ServiceStarted:    januaryMid,
//			},
//		},
//	})
//
//	var invoiceStruct struct {
//		Invoice_Simulate []*model.InvoiceSimulate
//	}
//
//	require.Nil(t, rawResponse.Errors)
//	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
//	require.Nil(t, err)
//
//	require.Equal(t, 2, len(invoiceStruct.Invoice_Simulate))
//
//	proratedInvoice := invoiceStruct.Invoice_Simulate[0]
//	require.Equal(t, 1, len(proratedInvoice.InvoiceLineItems))
//	asserInvoice(t, proratedInvoice, "2024-01-17T00:00:00Z", "2024-01-31T00:00:00Z", true, false, 14.3)
//	asserInvoiceLineItem(t, proratedInvoice.InvoiceLineItems[0], "1", "S1", 29, 1, 14.3)
//
//	onCycleInvoice := invoiceStruct.Invoice_Simulate[1]
//	require.Equal(t, 1, len(onCycleInvoice.InvoiceLineItems))
//	asserInvoice(t, onCycleInvoice, "2024-02-01T00:00:00Z", "2024-02-29T00:00:00Z", false, false, 29)
//	asserInvoiceLineItem(t, onCycleInvoice.InvoiceLineItems[0], "1", "S1", 29, 1, 29)
//}

// 1 downgrade SLI not prorated
func TestMutationResolver_InvoiceSimulate_OnCycle_PostPaidFalse_4(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	decemberFirst := utils.FirstTimeOfMonth(2023, 12)
	januaryMid := utils.MiddleTimeOfMonth(2024, 1)
	februaryFirst := utils.FirstTimeOfMonth(2024, 2)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{
		InvoicingPostpaid: false,
	})
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		BillingCycleInMonths: 1,
		InvoicingEnabled:     true,
		InvoicingStartDate:   &decemberFirst,
		NextInvoiceDate:      &februaryFirst,
	})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:      "S1",
		Billed:    neo4jenum.BilledTypeMonthly,
		Price:     5,
		Quantity:  1,
		StartedAt: decemberFirst,
	})

	rawResponse := callGraphQL(t, "invoice/simulate_invoice", map[string]interface{}{
		"contractId": contractId,
		"serviceLines": []model.InvoiceSimulateServiceLineInput{
			{
				Key:               "1",
				ServiceLineItemID: &serviceLineItemId,
				ParentID:          &serviceLineItemId,
				Description:       "S1",
				BillingCycle:      model.BilledTypeMonthly,
				Price:             2,
				Quantity:          1,
				ServiceStarted:    januaryMid,
			},
		},
	})

	var invoiceStruct struct {
		Invoice_Simulate []*model.InvoiceSimulate
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, 1, len(invoiceStruct.Invoice_Simulate))

	onCycleInvoice := invoiceStruct.Invoice_Simulate[0]
	require.Equal(t, 1, len(onCycleInvoice.InvoiceLineItems))
	asserInvoice(t, onCycleInvoice, "2024-02-01T00:00:00Z", "2024-02-29T00:00:00Z", false, false, 2)
	asserInvoiceLineItem(t, onCycleInvoice.InvoiceLineItems[0], "1", "S1", 2, 1, 2)
}

func TestMutationResolver_InvoiceSimulate_OnCycle_PostPaidTrue_1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	januaryFirst := utils.FirstTimeOfMonth(2024, 1)
	februaryFirst := utils.FirstTimeOfMonth(2024, 2)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{
		InvoicingPostpaid: true,
	})
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		BillingCycleInMonths: 1,
		InvoicingEnabled:     true,
		InvoicingStartDate:   &januaryFirst,
	})

	rawResponse := callGraphQL(t, "invoice/simulate_invoice", map[string]interface{}{
		"contractId": contractId,
		"serviceLines": []model.InvoiceSimulateServiceLineInput{
			{
				Key:            "1",
				Description:    "S1",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          1,
				Quantity:       1,
				ServiceStarted: februaryFirst.Add(-1),
			},
			{
				Key:            "2",
				Description:    "S2",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          2,
				Quantity:       1,
				ServiceStarted: januaryFirst,
			},
			{
				Key:            "3",
				Description:    "S3",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          3,
				Quantity:       1,
				ServiceStarted: februaryFirst,
			},
		},
	})

	var invoiceStruct struct {
		Invoice_Simulate []*model.InvoiceSimulate
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, 1, len(invoiceStruct.Invoice_Simulate))

	onCycleInvoice := invoiceStruct.Invoice_Simulate[0]
	require.Equal(t, 2, len(onCycleInvoice.InvoiceLineItems))
	asserInvoice(t, onCycleInvoice, "2024-01-01T00:00:00Z", "2024-01-31T00:00:00Z", false, true, 3)
	asserInvoiceLineItem(t, onCycleInvoice.InvoiceLineItems[0], "1", "S1", 1, 1, 1)
	asserInvoiceLineItem(t, onCycleInvoice.InvoiceLineItems[1], "2", "S2", 2, 1, 2)
}

func TestMutationResolver_InvoiceSimulate_OnCycle_PostPaidFalse_CanceledSLI_NotInvoiced(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	januaryFirst := utils.FirstTimeOfMonth(2024, 1)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{
		InvoicingPostpaid: false,
	})
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		BillingCycleInMonths:  1,
		InvoicingEnabled:      true,
		InvoicingStartDate:    &januaryFirst,
		OrganizationLegalName: "Test Organization",
	})

	rawResponse := callGraphQL(t, "invoice/simulate_invoice", map[string]interface{}{
		"contractId": contractId,
		"serviceLines": []model.InvoiceSimulateServiceLineInput{
			{
				Key:            "1",
				Description:    "S1",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          1,
				Quantity:       1,
				ServiceStarted: januaryFirst,
			},
			{
				Key:            "2",
				Description:    "S2",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          2,
				Quantity:       2,
				ServiceStarted: januaryFirst,
				CloseVersion:   utils.BoolPtr(true),
			},
		},
	})

	var invoiceStruct struct {
		Invoice_Simulate []*model.InvoiceSimulate
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, 1, len(invoiceStruct.Invoice_Simulate))

	invoice := invoiceStruct.Invoice_Simulate[0]

	require.Equal(t, 1, len(invoice.InvoiceLineItems))
	require.Equal(t, "Test Organization", *invoice.Customer.Name)

	asserInvoice(t, invoice, "2024-01-01T00:00:00Z", "2024-01-31T00:00:00Z", false, false, 1)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[0], "1", "S1", 1, 1, 1)
}

func asserInvoice(t *testing.T, invoice *model.InvoiceSimulate, periodStart, periodEnd string, offCycle, postpaid bool, total float64) {
	require.Equal(t, periodStart, invoice.InvoicePeriodStart.Format("2006-01-02T15:04:05Z"))
	require.Equal(t, periodEnd, invoice.InvoicePeriodEnd.Format("2006-01-02T15:04:05Z"))
	require.Equal(t, offCycle, invoice.OffCycle)
	require.Equal(t, postpaid, invoice.Postpaid)
	require.Equal(t, total, invoice.Total)
}

func asserInvoiceLineItem(t *testing.T, invoiceLineItem *model.InvoiceLineSimulate, key, description string, price float64, quantity int, total float64) {
	//require.Equal(t, key, invoiceLineItem.Key) //TODO put back when key is added to the response
	require.Equal(t, description, invoiceLineItem.Description)
	require.Equal(t, price, invoiceLineItem.Price)
	require.Equal(t, int64(quantity), invoiceLineItem.Quantity)
	require.Equal(t, total, invoiceLineItem.Total)
}
