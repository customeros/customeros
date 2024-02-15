package resolver

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
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
	yesterday := timeNow.Add(-24 * time.Hour)
	tomorrow := timeNow.Add(24 * time.Hour)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		CreatedAt:        timeNow,
		UpdatedAt:        timeNow,
		DryRun:           false,
		Number:           "1",
		Currency:         "RON",
		PeriodStartDate:  yesterday,
		PeriodEndDate:    tomorrow,
		DueDate:          timeNow,
		Amount:           100,
		Vat:              19,
		TotalAmount:      119,
		RepositoryFileId: "ABC",
		Note:             "Note",
	})

	neo4jtest.CreateInvoiceLine(ctx, driver, tenantName, invoiceId, neo4jentity.InvoiceLineEntity{
		CreatedAt:   timeNow,
		Name:        "SLI 1",
		Price:       100,
		Quantity:    1,
		Amount:      100,
		Vat:         19,
		TotalAmount: 119,
	})

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
	require.False(t, invoice.DryRun)
	require.False(t, invoice.Postpaid)
	require.False(t, invoice.OffCycle)
	require.Equal(t, "1", invoice.InvoiceNumber)
	require.Equal(t, fmt.Sprintf("%s/%s", constants.UrlInvoices, invoice.Metadata.ID), invoice.InvoiceURL)
	require.Equal(t, "RON", invoice.Currency)
	require.Equal(t, yesterday, invoice.InvoicePeriodStart)
	require.Equal(t, tomorrow, invoice.InvoicePeriodEnd)
	require.Equal(t, timeNow, invoice.Due)
	require.Equal(t, 100.0, invoice.Subtotal)
	require.Equal(t, 19.0, invoice.TaxDue)
	require.Equal(t, 119.0, invoice.AmountDue)
	require.Equal(t, "ABC", invoice.RepositoryFileID)
	require.Equal(t, "Note", *invoice.Note)
	require.False(t, invoice.Paid)
	require.Equal(t, 119.0, invoice.AmountRemaining)
	require.Equal(t, 0.0, invoice.AmountPaid)

	require.Equal(t, 1, len(invoice.InvoiceLineItems))
	require.Equal(t, "SLI 1", invoice.InvoiceLineItems[0].Description)
	require.Equal(t, 100.0, invoice.InvoiceLineItems[0].Price)
	require.Equal(t, 1, invoice.InvoiceLineItems[0].Quantity)
	require.Equal(t, 100.0, invoice.InvoiceLineItems[0].Subtotal)
	require.Equal(t, 19.0, invoice.InvoiceLineItems[0].TaxDue)
	require.Equal(t, 119.0, invoice.InvoiceLineItems[0].Total)

	require.Equal(t, organizationId, invoice.Organization.ID)
}

func TestInvoiceResolver_Invoices(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	timeNow := utils.Now()
	yesterday := timeNow.Add(-24 * time.Hour)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})

	invoice1Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Number:    "1",
		DryRun:    false,
	})
	neo4jtest.CreateInvoiceLine(ctx, driver, tenantName, invoice1Id, neo4jentity.InvoiceLineEntity{
		Name: "SLI 1",
	})

	invoice2Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		CreatedAt: yesterday,
		UpdatedAt: yesterday,
		Number:    "2",
		DryRun:    false,
	})
	neo4jtest.CreateInvoiceLine(ctx, driver, tenantName, invoice2Id, neo4jentity.InvoiceLineEntity{
		Name: "SLI 2",
	})

	invoice3Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		CreatedAt: yesterday,
		UpdatedAt: yesterday,
		Number:    "11",
		DryRun:    false,
	})
	neo4jtest.CreateInvoiceLine(ctx, driver, tenantName, invoice3Id, neo4jentity.InvoiceLineEntity{
		Name: "SLI 3",
	})

	invoice4Id := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		CreatedAt: yesterday,
		UpdatedAt: yesterday,
		Number:    "11",
		DryRun:    true,
	})
	neo4jtest.CreateInvoiceLine(ctx, driver, tenantName, invoice4Id, neo4jentity.InvoiceLineEntity{
		Name: "SLI 4",
	})

	rawResponse := callGraphQL(t, "invoice/get_invoices", map[string]interface{}{
		"page":  0,
		"limit": 10,
	})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoices model.InvoicesPage
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, int64(2), invoiceStruct.Invoices.TotalElements)
	require.Equal(t, 2, len(invoiceStruct.Invoices.Content))

	require.Equal(t, invoice1Id, invoiceStruct.Invoices.Content[0].ID)
	require.Equal(t, "1", invoiceStruct.Invoices.Content[0].Number)
	require.Equal(t, "SLI 1", invoiceStruct.Invoices.Content[0].InvoiceLines[0].Description)
	require.Equal(t, "11", invoiceStruct.Invoices.Content[1].Number)
	require.Equal(t, "SLI 3", invoiceStruct.Invoices.Content[1].InvoiceLines[0].Description)
}

func TestInvoiceResolver_SimulateInvoice(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	timeNow := utils.Now()
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{})

	calledSimulateInvoice := false
	invoiceServiceCallbacks := events_platform.MockInvoiceServiceCallbacks{
		SimulateInvoice: func(context context.Context, request *invoicepb.SimulateInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, request.Tenant)
			require.Equal(t, testUserId, request.LoggedInUserId)
			require.Equal(t, contractId, request.ContractId)
			require.Equal(t, 2, len(request.DryRunServiceLineItems))

			require.Equal(t, "1", request.DryRunServiceLineItems[0].ServiceLineItemId)
			require.Equal(t, "SLI 1", request.DryRunServiceLineItems[0].Name)
			require.Equal(t, commonpb.BilledType_MONTHLY_BILLED, request.DryRunServiceLineItems[0].Billed)
			require.Equal(t, 100.0, request.DryRunServiceLineItems[0].Price)
			require.Equal(t, int64(1), request.DryRunServiceLineItems[0].Quantity)

			require.Equal(t, "", request.DryRunServiceLineItems[1].ServiceLineItemId)
			require.Equal(t, "New SLI", request.DryRunServiceLineItems[1].Name)
			require.Equal(t, commonpb.BilledType_NONE_BILLED, request.DryRunServiceLineItems[1].Billed)
			require.Equal(t, 10.0, request.DryRunServiceLineItems[1].Price)
			require.Equal(t, int64(5), request.DryRunServiceLineItems[1].Quantity)

			require.Equal(t, constants.AppSourceCustomerOsApi, request.SourceFields.AppSource)
			calledSimulateInvoice = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	events_platform.SetInvoiceCallbacks(&invoiceServiceCallbacks)

	rawResponse := callGraphQL(t, "invoice/simulate_invoice", map[string]interface{}{
		"invoice": map[string]interface{}{
			"contractId":      contractId,
			"periodStartDate": timeNow,
			"invoiceLines": []map[string]interface{}{
				{
					"serviceLineItemId": "1",
					"name":              "SLI 1",
					"billed":            "MONTHLY",
					"price":             100,
					"quantity":          1,
				},
				{
					"serviceLineItemId": "",
					"billed":            "NONE",
					"name":              "New SLI",
					"price":             10,
					"quantity":          5,
				},
			},
		},
	})
	require.Nil(t, rawResponse.Errors)

	require.True(t, calledSimulateInvoice)

	var invoiceStruct struct {
		Invoice_Simulate string
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	require.Equal(t, invoiceId, invoiceStruct.Invoice_Simulate)
}

func TestInvoiceResolver_InvoicesForOrganization(t *testing.T) {
	ctx := context.TODO()
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
	contrac3tId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organization2Id, neo4jentity.ContractEntity{})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contrac3tId, neo4jentity.InvoiceEntity{
		Number: "3",
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		neo4jutil.NodeLabelOrganization: 2,
		neo4jutil.NodeLabelContract:     3,
		neo4jutil.NodeLabelInvoice:      3,
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
	require.Equal(t, 2, len(invoiceStruct.Invoices.Content))

	require.ElementsMatch(t, []string{invoice1Id, invoice2Id}, []string{invoiceStruct.Invoices.Content[0].ID, invoiceStruct.Invoices.Content[1].ID})
	require.ElementsMatch(t, []string{"1", "2"}, []string{invoiceStruct.Invoices.Content[0].Number, invoiceStruct.Invoices.Content[1].Number})
}

func TestInvoiceResolver_NextDryRunForContract(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	nextInvoiceDate := neo4jtest.FirstTimeOfMonth(2023, 6)
	periodStartExpected := neo4jtest.FirstTimeOfMonth(2023, 6)
	periodEndExpected := utils.LastDayOfMonth(2023, 6)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{
		BillingCycle:    neo4jenum.BillingCycleMonthlyBilling,
		Currency:        neo4jenum.CurrencyAUD,
		InvoiceNote:     "abc",
		NextInvoiceDate: &nextInvoiceDate,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{})

	calledNextDryRun := false
	invoiceServiceCallbacks := events_platform.MockInvoiceServiceCallbacks{
		NewInvoiceForContract: func(context context.Context, request *invoicepb.NewInvoiceForContractRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, request.Tenant)
			require.Equal(t, testUserId, request.LoggedInUserId)
			require.Equal(t, contractId, request.ContractId)
			require.Equal(t, commonpb.BillingCycle_MONTHLY_BILLING, request.BillingCycle)
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
		neo4jutil.NodeLabelOrganization: 1,
		neo4jutil.NodeLabelContract:     1,
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
