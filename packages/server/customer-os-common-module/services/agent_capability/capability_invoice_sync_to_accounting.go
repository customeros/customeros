package agent_capability

import (
	"context"
	"fmt"
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresrepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type SyncInvoiceToAccountingCapability struct {
	postgresRepositories *postgresrepository.Repositories
	neo4jRepositories    *neo4jrepository.Repositories
	invoiceService       interfaces.InvoiceService
	quickbooksService    interfaces.QuickbooksService
	contractService      interfaces.ContractService
}

func NewSyncInvoiceToAccountingCapability(
	postgresRepositories *postgresrepository.Repositories,
	neo4jRepositories *neo4jrepository.Repositories,
	invoiceService interfaces.InvoiceService,
	quickbooksService interfaces.QuickbooksService,
	contractService interfaces.ContractService) *SyncInvoiceToAccountingCapability {
	return &SyncInvoiceToAccountingCapability{
		neo4jRepositories:    neo4jRepositories,
		postgresRepositories: postgresRepositories,
		invoiceService:       invoiceService,
		quickbooksService:    quickbooksService,
		contractService:      contractService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[SyncInvoiceToAccountingInput, SyncInvoiceToAccountingOutput, SyncInvoiceToAccountingConfig] = (*SyncInvoiceToAccountingCapability)(nil)
)

func (c *SyncInvoiceToAccountingCapability) Type() enum.AgentCapability {
	return enum.CapabilitySyncInvoiceToAccounting
}

func (c *SyncInvoiceToAccountingCapability) Name() string {
	return "Sync invoice to accounting system"
}

func (c *SyncInvoiceToAccountingCapability) NewInput() SyncInvoiceToAccountingInput {
	return SyncInvoiceToAccountingInput{}
}

func (c *SyncInvoiceToAccountingCapability) NewConfig() SyncInvoiceToAccountingConfig {
	return SyncInvoiceToAccountingConfig{}
}

func (c *SyncInvoiceToAccountingCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *SyncInvoiceToAccountingCapability) DefaultActive() bool {
	return false
}

func (c *SyncInvoiceToAccountingCapability) ValidateInput(input SyncInvoiceToAccountingInput) error {
	if input.InvoiceID == "" {
		return errors.New("InvoiceID required")
	}
	return nil
}

func (c *SyncInvoiceToAccountingCapability) ValidateConfig(SyncInvoiceToAccountingConfig) error {
	return nil
}

type SyncInvoiceToAccountingInput struct {
	InvoiceID string `json:"invoiceId"`
}

type SyncInvoiceToAccountingOutput struct{}

type SyncInvoiceToAccountingConfig struct {
	Quickbooks              ConfigSingleBoolValue `json:"quickbooks"`
	AccountingMethodAccrual ConfigSingleBoolValue `json:"accountingMethodAccrual"`
}

func (c *SyncInvoiceToAccountingCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[SyncInvoiceToAccountingInput, SyncInvoiceToAccountingConfig]) (enum.CapabilityExecutionStatus, SyncInvoiceToAccountingOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncInvoiceToAccountingCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := SyncInvoiceToAccountingOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionCompleted, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionCompleted, result, err
	}

	invoice, err := c.invoiceService.GetById(ctx, nil, executionContainer.InputData.InvoiceID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionCompleted, result, err
	}
	if invoice == nil {
		err = errors.New("invoice not found")
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionCompleted, result, err
	}
	if invoice.DryRun {
		span.LogFields(log.String("skip", "Dry run"))
		return enum.CapabilityExecutionCompleted, result, nil
	}

	if invoice.QuickbooksInvoiceId == "" {
		err = c.syncInvoiceToQuickbooks(ctx, *invoice)
		if err != nil {
			tracing.TraceErr(span, err)
			return enum.CapabilityExecutionRetry, result, err
		}
	}

	// re-fetch invoice to get updated quickbooks invoice id
	invoice, err = c.invoiceService.GetById(ctx, nil, executionContainer.InputData.InvoiceID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionRetry, result, err
	}

	if executionContainer.ConfigData.AccountingMethodAccrual.Value && invoice.QuickbooksJournalEntryId == "" {
		err = c.syncInvoiceToQuickbooksJournalEntry(ctx, *invoice)
		if err != nil {
			tracing.TraceErr(span, err)
			return enum.CapabilityExecutionRetry, result, err
		}
	}

	// re-fetch invoice to get updated quickbooks journal id
	invoice, err = c.invoiceService.GetById(ctx, nil, executionContainer.InputData.InvoiceID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionRetry, result, err
	}

	if invoice.Status == neo4jenum.InvoiceStatusPaid {
		err = c.syncPaidInvoiceToQuickbooks(ctx, *invoice)
		if err != nil {
			tracing.TraceErr(span, err)
			return enum.CapabilityExecutionRetry, result, err
		}
	} else if invoice.Status == neo4jenum.InvoiceStatusVoid {
		err = c.syncVoidInvoiceToQuickbooks(ctx, *invoice)
		if err != nil {
			tracing.TraceErr(span, err)
			return enum.CapabilityExecutionRetry, result, err
		}
		if executionContainer.ConfigData.AccountingMethodAccrual.Value && invoice.QuickbooksJournalEntryId != "" {
			err = c.quickbooksService.ZeroJournalEntry(ctx, invoice.QuickbooksJournalEntryId)
			if err != nil {
				tracing.TraceErr(span, err)
				return enum.CapabilityExecutionRetry, result, err
			}
		}
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}

func (c *SyncInvoiceToAccountingCapability) syncInvoiceToQuickbooks(ctx context.Context, invoice neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncInvoiceToAccountingCapability.syncInvoiceToQuickbooks")
	defer span.Finish()
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tenant := common.GetTenantFromContext(ctx)

	quickbooksSettingsEntity, err := c.postgresRepositories.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if quickbooksSettingsEntity == nil {
		err = errors.New("Quickbooks settings not found")
		tracing.TraceErr(span, err)
		return err
	}

	invoiceLines, err := c.invoiceService.GetInvoiceLinesForInvoices(ctx, []string{invoice.Id})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if invoiceLines == nil || len(*invoiceLines) == 0 {
		span.LogKV("skip", "Invoice lines not found")
		return nil
	}

	// validate all invoice line have skuId and the skuId has been pushed to quickbooks
	for _, invoiceLine := range *invoiceLines {
		if invoiceLine.SkuId == "" {
			span.LogKV("skip", fmt.Sprintf("SkuId not found for invoice line %s", invoiceLine.Id))
			return nil
		}

		sku, err := c.postgresRepositories.SkuRepository.Get(ctx, tenant, invoiceLine.SkuId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if sku == nil {
			err = errors.New(fmt.Sprintf("Sku not found for invoice line %s", invoiceLine.Id))
			tracing.TraceErr(span, err)
			return err
		}

		if sku.QuickbooksId == "" {
			// Save sku to quickbooks
			qbProduct, err := c.quickbooksService.SaveProduct(ctx, "", sku.Name, sku.Archived, sku.Price)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error saving product to quickbooks"))
				return err
			}
			if qbProduct == nil || qbProduct.Product == nil {
				tracing.LogObjectAsJson(span, "qbProduct", qbProduct)
				err := errors.New("Quickbooks product could not be saved")
				tracing.TraceErr(span, err)
				return err
			}
			sku.QuickbooksId = qbProduct.Product.Id
			_, err = c.postgresRepositories.SkuRepository.Save(ctx, sku)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error saving sku in db"))
				return err
			}
		}
	}

	// save invoice to quickbooks
	quickbooksInvoiceLines := make([]interfaces.QuickbooksInvoiceLine, 0)
	for _, invoiceLine := range *invoiceLines {
		sku, err := c.postgresRepositories.SkuRepository.Get(ctx, tenant, invoiceLine.SkuId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		quickbooksInvoiceLine := interfaces.QuickbooksInvoiceLine{
			DetailType: "SalesItemLineDetail",
			Amount:     invoiceLine.Amount,
		}
		quickbooksInvoiceLine.SalesItemLineDetail.ItemRef.Value = sku.QuickbooksId
		quickbooksInvoiceLine.SalesItemLineDetail.ServiceDate = invoice.PeriodStartDate.Format("2006-01-02")

		quickbooksInvoiceLines = append(quickbooksInvoiceLines, quickbooksInvoiceLine)
	}

	contractEntity, err := c.contractService.GetContractForInvoice(ctx, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if contractEntity.QuickbooksCustomerId == "" {
		quickbooksSaveCustomerResponse, err := c.quickbooksService.SaveCustomer(ctx, "", contractEntity.OrganizationLegalName)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if quickbooksSaveCustomerResponse.Customer == nil {
			err := errors.New("Quickbooks customer not found")
			tracing.TraceErr(span, err)
			return err
		}

		contractEntity.QuickbooksCustomerId = quickbooksSaveCustomerResponse.Customer.Id

		err = c.neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonmodel.NodeLabelContract, contractEntity.Id, string(neo4jentity.ContractPropertyQuickbooksCustomerId), contractEntity.QuickbooksCustomerId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	savedInvoiced, err := c.quickbooksService.SaveInvoice(ctx, contractEntity.QuickbooksCustomerId, invoice.Number, invoice.IssuedDate, invoice.DueDate, contractEntity.InvoiceEmail, quickbooksInvoiceLines)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if savedInvoiced == nil {
		err := errors.New("Invoice not saved in quickbooks")
		tracing.TraceErr(span, err)
		return err
	}

	err = c.neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonmodel.NodeLabelInvoice, invoice.Id, string(neo4jentity.InvoicePropertyQuickbooksInvoiceId), savedInvoiced.Invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (c *SyncInvoiceToAccountingCapability) syncInvoiceToQuickbooksJournalEntry(ctx context.Context, invoice neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncInvoiceToAccountingCapability.syncInvoiceToQuickbooksJournalEntry")
	defer span.Finish()
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tenant := common.GetTenantFromContext(ctx)

	quickbooksSettingsEntity, err := c.postgresRepositories.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if quickbooksSettingsEntity == nil {
		err = errors.New("Quickbooks settings not found")
		tracing.TraceErr(span, err)
		return err
	}

	invoiceLines, err := c.invoiceService.GetInvoiceLinesForInvoices(ctx, []string{invoice.Id})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if invoiceLines == nil || len(*invoiceLines) == 0 {
		span.LogKV("skip", "Invoice lines not found")
		return nil
	}

	contractEntity, err := c.contractService.GetContractForInvoice(ctx, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// prepare amounts by income account
	journalLineAmountsByIncomeAccount := make(map[string]float64)
	for _, invoiceLine := range *invoiceLines {
		sku, err := c.postgresRepositories.SkuRepository.Get(ctx, tenant, invoiceLine.SkuId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if sku == nil {
			err = errors.New(fmt.Sprintf("Sku not found for invoice line %s", invoiceLine.Id))
			tracing.TraceErr(span, err)
			return err
		}
		if sku.QuickbooksId == "" {
			err = errors.New(fmt.Sprintf("QuickbooksId not found for sku %s", sku.ID))
			tracing.TraceErr(span, err)
			return err
		}
		quickbooksProduct, err := c.quickbooksService.GetProduct(ctx, sku.QuickbooksId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if quickbooksProduct == nil {
			err = errors.New(fmt.Sprintf("Quickbooks product not found for sku %s", sku.ID))
			tracing.TraceErr(span, err)
			return err
		}

		if _, ok := journalLineAmountsByIncomeAccount[quickbooksProduct.Item.IncomeAccountRef.Value]; !ok {
			journalLineAmountsByIncomeAccount[quickbooksProduct.Item.IncomeAccountRef.Value] = invoiceLine.Amount
		} else {
			journalLineAmountsByIncomeAccount[quickbooksProduct.Item.IncomeAccountRef.Value] += invoiceLine.Amount
		}
	}

	// prepare debit journal line item
	journalLineItems := make([]interfaces.QuickbooksJournalEntryLine, 0)
	debitJournalLineItem := interfaces.QuickbooksJournalEntryLine{
		DetailType:  "JournalEntryLineDetail",
		Amount:      invoice.Amount,
		Description: invoice.Number,
	}
	debitJournalLineItem.JournalEntryLineDetail.PostingType = "Debit"
	debitJournalLineItem.JournalEntryLineDetail.Entity.EntityRef.Value = contractEntity.QuickbooksCustomerId
	debitJournalLineItem.JournalEntryLineDetail.AccountRef.Name = "Accounts receivable (A/R)"
	debtorsAccountId, err := c.quickbooksService.GetAccountIdByName(ctx, "Accounts receivable (A/R)")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	debitJournalLineItem.JournalEntryLineDetail.AccountRef.Value = debtorsAccountId
	journalLineItems = append(journalLineItems, debitJournalLineItem)

	// prepare credit journal line items
	for incomeAccount, amount := range journalLineAmountsByIncomeAccount {
		journalLineItem := interfaces.QuickbooksJournalEntryLine{
			DetailType:  "JournalEntryLineDetail",
			Amount:      amount,
			Description: invoice.Number,
		}
		journalLineItem.JournalEntryLineDetail.PostingType = "Credit"
		journalLineItem.JournalEntryLineDetail.Entity.EntityRef.Value = contractEntity.QuickbooksCustomerId
		journalLineItem.JournalEntryLineDetail.AccountRef.Value = incomeAccount
		journalLineItems = append(journalLineItems, journalLineItem)
	}

	savedJournalEntry, err := c.quickbooksService.SaveJournalEntry(ctx, invoice.PeriodEndDate, journalLineItems)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if savedJournalEntry == nil {
		err := errors.New("Journal entry not saved in quickbooks")
		tracing.TraceErr(span, err)
		return err
	}

	err = c.neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonmodel.NodeLabelInvoice, invoice.Id, string(neo4jentity.InvoicePropertyQuickbooksJournalEntryId), savedJournalEntry.JournalEntry.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	quickbooksPayment, err := c.quickbooksService.SavePaymentLinkingJournalEntryToInvoice(ctx, contractEntity.QuickbooksCustomerId, invoice.QuickbooksInvoiceId, savedJournalEntry.JournalEntry.Id, invoice.IssuedDate, invoice.Amount)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error saving payment linking journal entry to invoice"))
		return err
	}
	if quickbooksPayment != nil {
		err = c.neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonmodel.NodeLabelInvoice, invoice.Id, string(neo4jentity.InvoicePropertyQuickbooksPaymentId), quickbooksPayment.Payment.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (c *SyncInvoiceToAccountingCapability) syncPaidInvoiceToQuickbooks(ctx context.Context, invoice neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncInvoiceToAccountingCapability.syncPaidInvoiceToQuickbooks")
	defer span.Finish()
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tenant := common.GetTenantFromContext(ctx)

	quickbooksSettingsEntity, err := c.postgresRepositories.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if quickbooksSettingsEntity == nil {
		err = errors.New("Quickbooks settings not found")
		tracing.TraceErr(span, err)
		return err
	}

	contractEntity, err := c.contractService.GetContractForInvoice(ctx, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	paymentResponse, err := c.quickbooksService.PayInvoice(ctx, contractEntity.QuickbooksCustomerId, invoice.QuickbooksInvoiceId, invoice.TotalAmount)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if paymentResponse == nil {
		err = errors.New("Invoice not paid in quickbooks")
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (c *SyncInvoiceToAccountingCapability) syncVoidInvoiceToQuickbooks(ctx context.Context, invoice neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncInvoiceToAccountingCapability.syncVoidInvoiceToQuickbooks")
	defer span.Finish()
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tenant := common.GetTenantFromContext(ctx)

	quickbooksSettingsEntity, err := c.postgresRepositories.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if quickbooksSettingsEntity == nil {
		err = errors.New("Quickbooks settings not found")
		tracing.TraceErr(span, err)
		return err
	}

	voidedResponse, err := c.quickbooksService.VoidInvoice(ctx, invoice.QuickbooksInvoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if voidedResponse == nil {
		err = errors.New("Invoice not voided in quickbooks")
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
