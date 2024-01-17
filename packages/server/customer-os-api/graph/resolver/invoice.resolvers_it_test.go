package resolver

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestQueryResolver_Invoice(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	timeNow := utils.Now()
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContract(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoice(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		CreatedAt:        timeNow,
		UpdatedAt:        timeNow,
		DryRun:           false,
		Number:           "1",
		Currency:         "RON",
		Date:             timeNow,
		DueDate:          timeNow,
		Amount:           100,
		Vat:              19,
		Total:            119,
		RepositoryFileId: "ABC",
	})

	neo4jtest.CreateInvoiceLine(ctx, driver, tenantName, invoiceId, neo4jentity.InvoiceLineEntity{
		CreatedAt: timeNow,
		Name:      "SLI 1",
		Price:     100,
		Quantity:  1,
		Amount:    100,
		Vat:       19,
		Total:     119,
	})

	rawResponse := callGraphQL(t, "invoice/get_invoice", map[string]interface{}{"id": invoiceId})
	require.Nil(t, rawResponse.Errors)

	var invoiceStruct struct {
		Invoice model.Invoice
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &invoiceStruct)
	require.Nil(t, err)

	invoice := invoiceStruct.Invoice
	require.Equal(t, invoiceId, invoice.ID)
	require.Equal(t, timeNow, invoice.CreatedAt)
	require.Equal(t, timeNow, invoice.UpdatedAt)
	require.Equal(t, false, invoice.DryRun)
	require.Equal(t, "1", invoice.Number)
	require.Equal(t, "RON", invoice.Currency)
	require.Equal(t, timeNow, invoice.Date)
	require.Equal(t, timeNow, invoice.DueDate)
	require.Equal(t, 100.0, invoice.Amount)
	require.Equal(t, 19.0, invoice.Vat)
	require.Equal(t, 119.0, invoice.Total)
	require.Equal(t, "ABC", invoice.RepositoryFileID)

	require.Equal(t, 1, len(invoice.InvoiceLines))
	require.Equal(t, "SLI 1", invoice.InvoiceLines[0].Name)
	require.Equal(t, 100.0, invoice.InvoiceLines[0].Price)
	require.Equal(t, 1, invoice.InvoiceLines[0].Quantity)
	require.Equal(t, 100.0, invoice.InvoiceLines[0].Amount)
	require.Equal(t, 19.0, invoice.InvoiceLines[0].Vat)
	require.Equal(t, 119.0, invoice.InvoiceLines[0].Total)
}
