package resolver

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestInvoiceResolver_Invoice(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	timeNow := utils.Now()
	today := utils.ToDate(timeNow)
	yesterday := timeNow.Add(-24 * time.Hour)
	tomorrow := timeNow.Add(24 * time.Hour)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})
	sliId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Price:    10.5,
		Quantity: 2,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		CreatedAt:        timeNow,
		UpdatedAt:        timeNow,
		DryRun:           true,
		Preview:          true,
		OffCycle:         true,
		Postpaid:         true,
		Number:           "1",
		Currency:         "RON",
		PeriodStartDate:  yesterday,
		PeriodEndDate:    tomorrow,
		DueDate:          timeNow,
		IssuedDate:       today,
		Amount:           100,
		Vat:              19,
		TotalAmount:      119,
		RepositoryFileId: "ABC",
		Note:             "Note",
		PaymentDetails: neo4jentity.PaymentDetails{
			PaymentLink: "Link",
		},
	})
	invoiceLineId := neo4jtest.CreateInvoiceLine(ctx, driver, tenantName, invoiceId, neo4jentity.InvoiceLineEntity{
		CreatedAt:   timeNow,
		Name:        "SLI 1",
		Price:       100,
		Quantity:    1,
		Amount:      100,
		Vat:         19,
		TotalAmount: 119,
	})
	neo4jtest.LinkNodes(ctx, driver, invoiceLineId, sliId, "INVOICED")

	rawResponse := callGraphQL(t, "invoice/get_invoice", map[string]interface{}{"id": invoiceId})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoice model.Invoice
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	invoice := invoiceStruct.Invoice
	require.Equal(t, invoiceId, invoice.Metadata.ID)
	require.Equal(t, timeNow, invoice.Metadata.Created)
	require.Equal(t, timeNow, invoice.Metadata.LastUpdated)
	require.True(t, invoice.DryRun)
	require.True(t, invoice.Postpaid)
	require.True(t, invoice.OffCycle)
	require.True(t, invoice.Preview)
	require.Equal(t, "1", invoice.InvoiceNumber)
	require.Equal(t, fmt.Sprintf(constants.UrlFileStoreFileDownloadUrlTemplate, invoice.RepositoryFileID), invoice.InvoiceURL)
	require.Equal(t, "RON", invoice.Currency)
	require.Equal(t, utils.ToDate(yesterday), invoice.InvoicePeriodStart)
	require.Equal(t, utils.ToDate(tomorrow), invoice.InvoicePeriodEnd)
	require.Equal(t, utils.ToDate(timeNow), invoice.Due)
	require.Equal(t, today, invoice.Issued)
	require.Equal(t, 100.0, invoice.Subtotal)
	require.Equal(t, 19.0, invoice.TaxDue)
	require.Equal(t, 119.0, invoice.AmountDue)
	require.Equal(t, "ABC", invoice.RepositoryFileID)
	require.Equal(t, "Note", *invoice.Note)
	require.False(t, invoice.Paid)
	require.Equal(t, 119.0, invoice.AmountRemaining)
	require.Equal(t, 0.0, invoice.AmountPaid)
	require.Equal(t, "Link", *invoice.PaymentLink)

	require.Equal(t, 1, len(invoice.InvoiceLineItems))
	require.Equal(t, "SLI 1", invoice.InvoiceLineItems[0].Description)
	require.Equal(t, 100.0, invoice.InvoiceLineItems[0].Price)
	require.Equal(t, int64(1), invoice.InvoiceLineItems[0].Quantity)
	require.Equal(t, 100.0, invoice.InvoiceLineItems[0].Subtotal)
	require.Equal(t, 19.0, invoice.InvoiceLineItems[0].TaxDue)
	require.Equal(t, 119.0, invoice.InvoiceLineItems[0].Total)
	require.Equal(t, sliId, invoice.InvoiceLineItems[0].ContractLineItem.Metadata.ID)
	require.Equal(t, 10.5, invoice.InvoiceLineItems[0].ContractLineItem.Price)
	require.Equal(t, int64(2), invoice.InvoiceLineItems[0].ContractLineItem.Quantity)

	require.Equal(t, organizationId, invoice.Organization.ID)
	require.Equal(t, contractId, invoice.Contract.ID)
}

func TestInvoiceResolver_Invoice_ByNumber(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		CreatedAt:       utils.Now(),
		UpdatedAt:       utils.Now(),
		PeriodStartDate: utils.Now(),
		PeriodEndDate:   utils.Now(),
		DueDate:         utils.Now(),
		IssuedDate:      utils.Now(),
		DryRun:          true,
		Number:          "INV-001",
		TotalAmount:     119.01,
	})

	rawResponse := callGraphQL(t, "invoice/get_invoice_by_number", map[string]interface{}{"number": "INV-001"})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoice_ByNumber model.Invoice
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	invoice := invoiceStruct.Invoice_ByNumber
	require.Equal(t, invoiceId, invoice.Metadata.ID)
	require.True(t, invoice.DryRun)
	require.Equal(t, "INV-001", invoice.InvoiceNumber)
	require.Equal(t, 119.01, invoice.AmountDue)
	require.Equal(t, 0, len(invoice.InvoiceLineItems))
}

func TestInvoiceResolver_Invoices_Contract_Name(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contract1Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		Name: "A",
	})
	invoice1Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract1Id, neo4jentity.InvoiceEntity{})

	contract2Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		Name: "AA",
	})
	invoice2Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract2Id, neo4jentity.InvoiceEntity{})

	contract3Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		Name: "B",
	})
	invoice3Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract3Id, neo4jentity.InvoiceEntity{})

	assertInvoicesFilter(t, "CONTRACT_NAME", "A", []string{invoice1Id, invoice2Id}, int64(3))
	assertInvoicesFilter(t, "CONTRACT_NAME", "B", []string{invoice3Id}, int64(3))
}

func TestInvoiceResolver_Invoices_Contract_BillingCycle(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contract1Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		BillingCycleInMonths: 1,
	})
	invoice1Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract1Id, neo4jentity.InvoiceEntity{})

	contract2Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		BillingCycleInMonths: 3,
	})
	invoice2Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract2Id, neo4jentity.InvoiceEntity{})

	contract3Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		BillingCycleInMonths: 12,
	})
	invoice3Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract3Id, neo4jentity.InvoiceEntity{})

	contract4Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		BillingCycleInMonths: 0,
	})
	invoice4Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract4Id, neo4jentity.InvoiceEntity{})

	assertInvoicesFilter(t, "CONTRACT_BILLING_CYCLE", []string{"NONE"}, []string{invoice4Id}, int64(4))
	assertInvoicesFilter(t, "CONTRACT_BILLING_CYCLE", []string{"MONTHLY"}, []string{invoice1Id}, int64(4))
	assertInvoicesFilter(t, "CONTRACT_BILLING_CYCLE", []string{"QUARTERLY"}, []string{invoice2Id}, int64(4))
	assertInvoicesFilter(t, "CONTRACT_BILLING_CYCLE", []string{"ANNUALLY"}, []string{invoice3Id}, int64(4))
}

func assertInvoicesFilter(t *testing.T, propertyName string, propertyValue any, expectedInvoiceOrderIds []string, totalAvailable int64) {
	rawResponse := callGraphQL(t, "invoice/get_invoices_filter", map[string]interface{}{
		"page":     0,
		"limit":    10,
		"property": propertyName,
		"value":    propertyValue,
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, totalAvailable, invoiceStruct.Invoices.TotalAvailable)
	require.Equal(t, int64(len(expectedInvoiceOrderIds)), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, len(expectedInvoiceOrderIds), len(invoiceStruct.Invoices.Content))

	for _, invoice := range invoiceStruct.Invoices.Content {
		require.Contains(t, expectedInvoiceOrderIds, invoice.Metadata.ID)
	}
}

func TestInvoiceResolver_Invoices_Contract_Ended_False(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()
	tenWeeksAgo := now.Add(-10 * 7 * 24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contract1Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		EndedAt:              &tenWeeksAgo,
		BillingCycleInMonths: 1,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract1Id, neo4jentity.InvoiceEntity{})

	contract2Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		BillingCycleInMonths: 3,
	})
	invoice2Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract2Id, neo4jentity.InvoiceEntity{})

	rawResponse := callGraphQL(t, "invoice/get_invoices_filter_contract_ended", map[string]interface{}{
		"page":          0,
		"limit":         10,
		"contractEnded": false,
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(1), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, int64(2), invoiceStruct.Invoices.TotalAvailable)
	require.Equal(t, 1, len(invoiceStruct.Invoices.Content))

	require.Equal(t, invoice2Id, invoiceStruct.Invoices.Content[0].Metadata.ID)
}

func TestInvoiceResolver_Invoices_Contract_Ended_True(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()
	tenWeeksAgo := now.Add(-10 * 7 * 24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contract1Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		EndedAt:              &tenWeeksAgo,
		BillingCycleInMonths: 1,
	})
	invoice1Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract1Id, neo4jentity.InvoiceEntity{})

	contract2Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		BillingCycleInMonths: 3,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract2Id, neo4jentity.InvoiceEntity{})

	rawResponse := callGraphQL(t, "invoice/get_invoices_filter_contract_ended", map[string]interface{}{
		"page":          0,
		"limit":         10,
		"contractEnded": true,
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(1), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, int64(2), invoiceStruct.Invoices.TotalAvailable)
	require.Equal(t, 1, len(invoiceStruct.Invoices.Content))

	require.Equal(t, invoice1Id, invoiceStruct.Invoices.Content[0].Metadata.ID)
}

func TestInvoiceResolver_Invoices_Status(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})

	invoice1Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Status: neo4jenum.InvoiceStatusDue,
	})
	invoice2Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Status: neo4jenum.InvoiceStatusPaid,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Status: neo4jenum.InvoiceStatusVoid,
	})

	rawResponse := callGraphQL(t, "invoice/get_invoices_filter_status", map[string]interface{}{
		"page":          0,
		"limit":         10,
		"invoiceStatus": []string{"DUE", "PAID"},
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(2), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, 2, len(invoiceStruct.Invoices.Content))

	require.Contains(t, []string{invoice1Id, invoice2Id}, invoiceStruct.Invoices.Content[0].Metadata.ID)
	require.Contains(t, []string{invoice1Id, invoice2Id}, invoiceStruct.Invoices.Content[1].Metadata.ID)
}

func TestInvoiceResolver_Invoices_IssueDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	today := utils.ToDate(utils.Now())
	tenWeeksAgo := today.Add(-10 * 7 * 24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	inTenWeeks := today.Add(10 * 7 * 24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})

	invoice1Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		IssuedDate: yesterday,
	})
	invoice2Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		IssuedDate: tomorrow,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		IssuedDate: tenWeeksAgo,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		IssuedDate: inTenWeeks,
	})

	rawResponse := callGraphQL(t, "invoice/get_invoices_filter_issue_date", map[string]interface{}{
		"page":             0,
		"limit":            10,
		"invoiceIssueDate": []time.Time{yesterday, tomorrow},
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(2), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, 2, len(invoiceStruct.Invoices.Content))

	require.Contains(t, []string{invoice1Id, invoice2Id}, invoiceStruct.Invoices.Content[0].Metadata.ID)
	require.Contains(t, []string{invoice1Id, invoice2Id}, invoiceStruct.Invoices.Content[1].Metadata.ID)
}

func TestInvoiceResolver_Invoices_Preview_True(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})

	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Preview: false,
	})

	invoice2Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Preview: true,
	})

	rawResponse := callGraphQL(t, "invoice/get_invoices_filter_preview", map[string]interface{}{
		"page":    0,
		"limit":   10,
		"preview": true,
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(1), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, int64(1), invoiceStruct.Invoices.TotalAvailable)
	require.Equal(t, 1, len(invoiceStruct.Invoices.Content))

	require.Equal(t, invoice2Id, invoiceStruct.Invoices.Content[0].Metadata.ID)
}

func TestInvoiceResolver_Invoices_Preview_False(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})

	invoice1Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Preview: false,
	})

	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Preview: true,
	})

	rawResponse := callGraphQL(t, "invoice/get_invoices_filter_preview", map[string]interface{}{
		"page":    0,
		"limit":   10,
		"preview": false,
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(1), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, int64(1), invoiceStruct.Invoices.TotalAvailable)
	require.Equal(t, 1, len(invoiceStruct.Invoices.Content))

	require.Equal(t, invoice1Id, invoiceStruct.Invoices.Content[0].Metadata.ID)
}

func TestInvoiceResolver_Invoices_DryRun_True(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})

	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		DryRun: false,
	})

	invoice2Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		DryRun: true,
	})

	rawResponse := callGraphQL(t, "invoice/get_invoices_filter_dry_run", map[string]interface{}{
		"page":   0,
		"limit":  10,
		"dryRun": true,
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(1), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, int64(1), invoiceStruct.Invoices.TotalAvailable)
	require.Equal(t, 1, len(invoiceStruct.Invoices.Content))

	require.Equal(t, invoice2Id, invoiceStruct.Invoices.Content[0].Metadata.ID)
}

func TestInvoiceResolver_Invoices_DryRun_False(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})

	invoice1Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		DryRun: false,
	})

	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		DryRun: true,
	})

	rawResponse := callGraphQL(t, "invoice/get_invoices_filter_dry_run", map[string]interface{}{
		"page":   0,
		"limit":  10,
		"dryRun": false,
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(1), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, int64(1), invoiceStruct.Invoices.TotalAvailable)
	require.Equal(t, 1, len(invoiceStruct.Invoices.Content))

	require.Equal(t, invoice1Id, invoiceStruct.Invoices.Content[0].Metadata.ID)
}

func TestInvoiceResolver_Invoices_Number(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})

	invoice1Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Number: "1",
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Number: "11",
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Number: "2",
	})

	rawResponse := callGraphQL(t, "invoice/get_invoices_filter_number", map[string]interface{}{
		"page":          0,
		"limit":         10,
		"invoiceNumber": "1",
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(1), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, int64(3), invoiceStruct.Invoices.TotalAvailable)
	require.Equal(t, 1, len(invoiceStruct.Invoices.Content))

	require.Equal(t, invoice1Id, invoiceStruct.Invoices.Content[0].Metadata.ID)
}

func TestInvoiceResolver_Invoices_Exclude_INITIALIZED_EMPTY_Statuses(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})

	invoice1Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Status: neo4jenum.InvoiceStatusInitialized,
	})
	invoice2Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Status: neo4jenum.InvoiceStatusEmpty,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Status: neo4jenum.InvoiceStatusDue,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Status: neo4jenum.InvoiceStatusPaid,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Status: neo4jenum.InvoiceStatusVoid,
	})

	rawResponse := callGraphQL(t, "invoice/get_invoices_exclude_initialized_status", map[string]interface{}{
		"page":  0,
		"limit": 10,
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(3), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, int64(3), invoiceStruct.Invoices.TotalAvailable)
	require.Equal(t, 3, len(invoiceStruct.Invoices.Content))

	for _, invoice := range invoiceStruct.Invoices.Content {
		require.NotEqual(t, invoice1Id, invoice.Metadata.ID)
		require.NotEqual(t, invoice2Id, invoice.Metadata.ID)
	}
}

func TestInvoiceResolver_Invoices_Sorting(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contract1Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		EndedAt: &yesterday,
		Name:    "A",
	})
	contract2Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		EndedAt: &now,
		Name:    "B",
	})
	contract3Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		EndedAt: &tomorrow,
		Name:    "C",
	})

	invoice1Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract1Id, neo4jentity.InvoiceEntity{
		CreatedAt: yesterday,
		DueDate:   yesterday,
		Number:    "1",
		Status:    neo4jenum.InvoiceStatusDue,
	})
	invoice2Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract2Id, neo4jentity.InvoiceEntity{
		CreatedAt: now,
		DueDate:   now,
		Number:    "2",
		Status:    neo4jenum.InvoiceStatusPaid,
	})
	invoice3Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract3Id, neo4jentity.InvoiceEntity{
		CreatedAt: tomorrow,
		DueDate:   tomorrow,
		Number:    "3",
		Status:    neo4jenum.InvoiceStatusVoid,
	})

	assertInvoicesSorted(t, "CONTRACT_NAME", "ASC", []string{invoice1Id, invoice2Id, invoice3Id})
	assertInvoicesSorted(t, "CONTRACT_NAME", "DESC", []string{invoice3Id, invoice2Id, invoice1Id})
	assertInvoicesSorted(t, "CONTRACT_ENDED_AT", "ASC", []string{invoice1Id, invoice2Id, invoice3Id})
	assertInvoicesSorted(t, "CONTRACT_ENDED_AT", "DESC", []string{invoice3Id, invoice2Id, invoice1Id})
	assertInvoicesSorted(t, "INVOICE_CREATED_AT", "ASC", []string{invoice1Id, invoice2Id, invoice3Id})
	assertInvoicesSorted(t, "INVOICE_CREATED_AT", "DESC", []string{invoice3Id, invoice2Id, invoice1Id})
	assertInvoicesSorted(t, "INVOICE_DUE_DATE", "ASC", []string{invoice1Id, invoice2Id, invoice3Id})
	assertInvoicesSorted(t, "INVOICE_DUE_DATE", "DESC", []string{invoice3Id, invoice2Id, invoice1Id})
	assertInvoicesSorted(t, "INVOICE_STATUS", "ASC", []string{invoice1Id, invoice2Id, invoice3Id})
	assertInvoicesSorted(t, "INVOICE_STATUS", "DESC", []string{invoice3Id, invoice2Id, invoice1Id})
}

func assertInvoicesSorted(t *testing.T, sortBy string, sortDirection string, expectedInvoiceOrderIds []string) {
	rawResponse := callGraphQL(t, "invoice/get_invoices_sort", map[string]interface{}{
		"page":          0,
		"limit":         10,
		"sortBy":        sortBy,
		"sortDirection": sortDirection,
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(3), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, int64(3), invoiceStruct.Invoices.TotalAvailable)
	require.Equal(t, 3, len(invoiceStruct.Invoices.Content))

	for i, invoice := range invoiceStruct.Invoices.Content {
		require.Equal(t, expectedInvoiceOrderIds[i], invoice.Metadata.ID)
	}
}

func TestInvoiceResolver_InvoicesForOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organization1Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contract1Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organization1Id, neo4jentity.ContractEntity{})
	contract2Id := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organization1Id, neo4jentity.ContractEntity{})
	invoice1Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract1Id, neo4jentity.InvoiceEntity{
		Number: "1",
	})
	invoice2Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contract2Id, neo4jentity.InvoiceEntity{
		Number: "2",
	})

	organization2Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId3 := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organization2Id, neo4jentity.ContractEntity{})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId3, neo4jentity.InvoiceEntity{
		Number: "3",
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		model2.NodeLabelOrganization: 2,
		model2.NodeLabelContract:     3,
		model2.NodeLabelInvoice:      3,
	})

	rawResponse := callGraphQL(t, "invoice/get_invoices_for_organization", map[string]interface{}{
		"page":           0,
		"limit":          10,
		"organizationId": organization1Id,
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(2), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, int64(2), invoiceStruct.Invoices.TotalAvailable)
	require.Equal(t, 2, len(invoiceStruct.Invoices.Content))

	require.ElementsMatch(t, []string{invoice1Id, invoice2Id}, []string{invoiceStruct.Invoices.Content[0].Metadata.ID, invoiceStruct.Invoices.Content[1].Metadata.ID})
	require.ElementsMatch(t, []string{"1", "2"}, []string{invoiceStruct.Invoices.Content[0].InvoiceNumber, invoiceStruct.Invoices.Content[1].InvoiceNumber})
}

func TestInvoiceResolver_NextDryRunForContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	nextInvoiceDate := utils.FirstTimeOfMonth(2023, 6)
	periodStartExpected := utils.FirstTimeOfMonth(2023, 6)
	periodEndExpected := utils.LastDayOfMonth(2023, 6)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		BillingCycleInMonths: 1,
		Currency:             neo4jenum.CurrencyAUD,
		InvoiceNote:          "abc",
		NextInvoiceDate:      &nextInvoiceDate,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{})

	calledNextDryRun := false
	invoiceServiceCallbacks := events_platform.MockInvoiceServiceCallbacks{
		NewInvoiceForContract: func(context context.Context, request *invoicepb.NewInvoiceForContractRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, request.Tenant)
			require.Equal(t, testUserId, request.LoggedInUserId)
			require.Equal(t, contractId, request.ContractId)
			require.Equal(t, int64(1), request.BillingCycleInMonths)
			require.Equal(t, neo4jenum.CurrencyAUD.String(), request.Currency)
			require.Equal(t, utils.ConvertTimeToTimestampPtr(&periodStartExpected), request.InvoicePeriodStart)
			require.Equal(t, utils.ConvertTimeToTimestampPtr(&periodEndExpected), request.InvoicePeriodEnd)
			require.Equal(t, constants.AppSourceCustomerOsApi, request.SourceFields.AppSource)
			calledNextDryRun = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	events_platform.SetInvoiceCallbacks(&invoiceServiceCallbacks)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		model2.NodeLabelOrganization: 1,
		model2.NodeLabelContract:     1,
	})

	rawResponse := callGraphQL(t, "invoice/next_dry_run_for_contract", map[string]interface{}{
		"page":       0,
		"limit":      10,
		"contractId": contractId,
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoice_NextDryRunForContract string
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.True(t, calledNextDryRun)

	require.Equal(t, invoiceId, invoiceStruct.Invoice_NextDryRunForContract)
}

func TestMutationResolver_InvoiceUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{})

	rawResponse := callGraphQL(t, "invoice/update_invoice", map[string]interface{}{
		"invoiceId": invoiceId,
	})

	var invoiceStruct struct {
		Invoice_Update model.Invoice
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)
	invoice := invoiceStruct.Invoice_Update
	require.Equal(t, invoiceId, invoice.Metadata.ID)
}

func TestMutationResolver_InvoicePay(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Status: neo4jenum.InvoiceStatusDue,
	})

	rawResponse := callGraphQL(t, "invoice/pay_invoice", map[string]interface{}{
		"invoiceId": invoiceId,
	})

	var invoiceStruct struct {
		Invoice_Pay model.Invoice
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)
	invoice := invoiceStruct.Invoice_Pay
	require.Equal(t, invoiceId, invoice.Metadata.ID)
}

func TestMutationResolver_InvoiceVoid(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{})

	calledVoidInvoice := false

	invoiceServiceCallbacks := events_platform.MockInvoiceServiceCallbacks{
		VoidInvoice: func(context context.Context, invoice *invoicepb.VoidInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, invoice.Tenant)
			require.Equal(t, invoiceId, invoice.InvoiceId)
			require.Equal(t, testUserId, invoice.LoggedInUserId)
			require.Equal(t, constants.AppSourceCustomerOsApi, invoice.AppSource)

			calledVoidInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	events_platform.SetInvoiceCallbacks(&invoiceServiceCallbacks)

	rawResponse := callGraphQL(t, "invoice/void_invoice", map[string]interface{}{
		"invoiceId": invoiceId,
	})

	var invoiceStruct struct {
		Invoice_Void model.Invoice
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)
	invoice := invoiceStruct.Invoice_Void
	require.Equal(t, invoiceId, invoice.Metadata.ID)

	require.True(t, calledVoidInvoice)
}
