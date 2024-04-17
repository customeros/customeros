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

func TestMutationResolver_InvoiceSimulate_OnCycle_PostPaidFalse_FirstInvoiceForContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	januaryFirst := utils.FirstTimeOfMonth(2024, 1)
	februaryFirst := utils.FirstTimeOfMonth(2024, 2)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingEnabled:   true,
		InvoicingStartDate: &januaryFirst,
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
				BillingCycle:   model.BilledTypeQuarterly,
				Price:          3,
				Quantity:       3,
				ServiceStarted: januaryFirst,
			},
			{
				Key:            "3",
				Description:    "S3",
				BillingCycle:   model.BilledTypeAnnually,
				Price:          12,
				Quantity:       12,
				ServiceStarted: januaryFirst,
			},
			{
				Key:            "4",
				Description:    "S4 - excluded - in the future",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          1,
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

	invoice := invoiceStruct.Invoice_Simulate[0]

	require.Equal(t, 3, len(invoice.InvoiceLineItems))

	asserInvoiceLineItem(t, invoice.InvoiceLineItems[0], "1", "S1", 1, 1, 1)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[1], "2", "S2", 1, 3, 3)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[2], "3", "S3", 1, 12, 12)
}

func TestMutationResolver_InvoiceSimulate_OnCycle_PostPaidFalse_SecondInvoiceForContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	januaryFirst := utils.FirstTimeOfMonth(2024, 1)
	februaryFirst := utils.FirstTimeOfMonth(2024, 2)
	marchFirst := utils.FirstTimeOfMonth(2024, 3)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingEnabled:   true,
		InvoicingStartDate: &januaryFirst,
		NextInvoiceDate:    &februaryFirst,
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
				BillingCycle:   model.BilledTypeQuarterly,
				Price:          3,
				Quantity:       3,
				ServiceStarted: januaryFirst,
			},
			{
				Key:            "3",
				Description:    "S3",
				BillingCycle:   model.BilledTypeAnnually,
				Price:          12,
				Quantity:       12,
				ServiceStarted: januaryFirst,
			},
			{
				Key:            "4",
				Description:    "S4",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          1,
				Quantity:       1,
				ServiceStarted: februaryFirst,
			},
			{
				Key:            "5",
				Description:    "S5 - excluded - in the future",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          1,
				Quantity:       1,
				ServiceStarted: marchFirst,
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

	require.Equal(t, 4, len(invoice.InvoiceLineItems))

	asserInvoiceLineItem(t, invoice.InvoiceLineItems[0], "1", "S1", 1, 1, 1)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[1], "2", "S2", 1, 3, 3)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[2], "3", "S3", 1, 12, 12)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[3], "4", "S4", 1, 1, 1)
}

func TestMutationResolver_InvoiceSimulate_OnCycle_PostPaidTrue_FirstInvoiceForContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	januaryFirst := utils.FirstTimeOfMonth(2024, 1)
	januaryMid := utils.MiddleTimeOfMonth(2024, 1)
	februaryFirst := utils.FirstTimeOfMonth(2024, 2)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{
		InvoicingPostpaid: true,
	})
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingEnabled:   true,
		InvoicingStartDate: &januaryFirst,
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
				BillingCycle:   model.BilledTypeQuarterly,
				Price:          3,
				Quantity:       3,
				ServiceStarted: januaryFirst,
			},
			{
				Key:            "3",
				Description:    "S3",
				BillingCycle:   model.BilledTypeAnnually,
				Price:          12,
				Quantity:       12,
				ServiceStarted: januaryFirst,
			},
			{
				Key:            "4",
				Description:    "S4",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          1,
				Quantity:       1,
				ServiceStarted: januaryMid,
			},
			{
				Key:            "5",
				Description:    "S5 - excluded - in future",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          1,
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

	invoice := invoiceStruct.Invoice_Simulate[0]

	require.Equal(t, 4, len(invoice.InvoiceLineItems))

	asserInvoiceLineItem(t, invoice.InvoiceLineItems[0], "1", "S1", 1, 1, 1)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[1], "2", "S2", 1, 3, 3)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[2], "3", "S3", 1, 12, 12)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[3], "4", "S4", 1, 1, 1)
}

func TestMutationResolver_InvoiceSimulate_OnCycle_PostPaidTrue_SecondInvoiceForContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	januaryFirst := utils.FirstTimeOfMonth(2024, 1)
	februaryFirst := utils.FirstTimeOfMonth(2024, 2)
	februaryMid := utils.MiddleTimeOfMonth(2024, 2)
	marchFirst := utils.FirstTimeOfMonth(2024, 3)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{
		InvoicingPostpaid: true,
	})
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingEnabled:   true,
		InvoicingStartDate: &januaryFirst,
		NextInvoiceDate:    &februaryFirst,
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
				BillingCycle:   model.BilledTypeQuarterly,
				Price:          3,
				Quantity:       3,
				ServiceStarted: januaryFirst,
			},
			{
				Key:            "3",
				Description:    "S3",
				BillingCycle:   model.BilledTypeAnnually,
				Price:          12,
				Quantity:       12,
				ServiceStarted: januaryFirst,
			},
			{
				Key:            "4",
				Description:    "S4",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          1,
				Quantity:       1,
				ServiceStarted: februaryFirst,
			},
			{
				Key:            "5",
				Description:    "S5",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          1,
				Quantity:       1,
				ServiceStarted: februaryMid,
			},
			{
				Key:            "6",
				Description:    "S6 - excluded - in the future",
				BillingCycle:   model.BilledTypeMonthly,
				Price:          1,
				Quantity:       1,
				ServiceStarted: marchFirst,
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

	require.Equal(t, 5, len(invoice.InvoiceLineItems))

	asserInvoiceLineItem(t, invoice.InvoiceLineItems[0], "1", "S1", 1, 1, 1)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[1], "2", "S2", 1, 3, 3)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[2], "3", "S3", 1, 12, 12)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[3], "4", "S4", 1, 1, 1)
	asserInvoiceLineItem(t, invoice.InvoiceLineItems[4], "5", "S5", 1, 1, 1)
}

//func TestMutationResolver_InvoiceSimulate_OffCycle_PostPaidFalse_FirstOffCycleInvoice(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//
//	januaryFirst := utils.FirstTimeOfMonth(2024, 1)
//	januaryMid := utils.MiddleTimeOfMonth(2024, 1)
//
//	neo4jtest.CreateTenant(ctx, driver, tenantName)
//	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
//
//	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{
//		InvoicingPostpaid: false,
//	})
//	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
//
//	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
//		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
//		InvoicingEnabled:   true,
//		InvoicingStartDate: &januaryFirst,
//	})
//	sli1Id := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
//		Name:      "S1",
//		Billed:    neo4jenum.BilledTypeMonthly,
//		Price:     1,
//		Quantity:  1,
//		StartedAt: januaryFirst,
//	})
//
//	rawResponse := callGraphQL(t, "invoice/simulate_invoice", map[string]interface{}{
//		"contractId": contractId,
//		"serviceLines": []model.InvoiceSimulateServiceLineInput{
//			{
//				Key:               "1",
//				ServiceLineItemID: &sli1Id,
//				ParentID:          &sli1Id,
//				Description:       "S1",
//				BillingCycle:      model.BilledTypeMonthly,
//				Price:             10,
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
//	require.Equal(t, 1, len(invoiceStruct.Invoice_Simulate))
//
//	invoice := invoiceStruct.Invoice_Simulate[0]
//
//	require.Equal(t, 1, len(invoice.InvoiceLineItems))
//
//	//valoarea este diferenta la numarul de zile ramase din perioada curenta
//	asserInvoiceLineItem(t, invoice.InvoiceLineItems[0], "1", "S1", 9, 1, 1)
//}

func asserInvoiceLineItem(t *testing.T, invoiceLineItem *model.InvoiceLineSimulate, key, description string, price float64, quantity int, total float64) {
	require.Equal(t, key, invoiceLineItem.Key)
	require.Equal(t, description, invoiceLineItem.Description)
	require.Equal(t, price, invoiceLineItem.Price)
	require.Equal(t, int64(quantity), invoiceLineItem.Quantity)
	require.Equal(t, total, invoiceLineItem.Total)
}
