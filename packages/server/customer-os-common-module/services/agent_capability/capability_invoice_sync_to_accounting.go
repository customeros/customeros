package agent_capability

import (
	"context"
	"fmt"
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
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
}

func NewSyncInvoiceToAccountingCapability(
	postgresRepositories *postgresrepository.Repositories,
	neo4jRepositories *neo4jrepository.Repositories,
	invoiceService interfaces.InvoiceService,
	quickbooksService interfaces.QuickbooksService) *SyncInvoiceToAccountingCapability {
	return &SyncInvoiceToAccountingCapability{
		neo4jRepositories:    neo4jRepositories,
		postgresRepositories: postgresRepositories,
		invoiceService:       invoiceService,
		quickbooksService:    quickbooksService,
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
	Quickbooks ConfigSingleBoolValue `json:"quickbooks"`
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
		return enum.CapabilityExecutionCompleted, result, nil
	}

	if invoice.QuickbooksInvoiceId != "" {
		err = c.syncInvoiceToQuickbooks(ctx, *invoice)
		if err != nil {
			tracing.TraceErr(span, err)
			return enum.CapabilityExecutionCompleted, result, err
		}
	} else {
		if invoice.Status == neo4jenum.InvoiceStatusPaid {
			err = c.syncPaidInvoiceToQuickbooks(ctx, *invoice)
			if err != nil {
				tracing.TraceErr(span, err)
				return enum.CapabilityExecutionCompleted, result, err
			}
		} else if invoice.Status == neo4jenum.InvoiceStatusVoid {
			err = c.syncVoidInvoiceToQuickbooks(ctx, *invoice)
			if err != nil {
				tracing.TraceErr(span, err)
				return enum.CapabilityExecutionCompleted, result, err
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
			span.LogFields(log.String("skip", fmt.Sprintf("QuickbooksId not found for sku %s", sku.ID)))
			return nil
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

		quickbooksInvoiceLines = append(quickbooksInvoiceLines, quickbooksInvoiceLine)
	}

	organizationNode, err := c.neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if organizationNode == nil {
		err := errors.New("Organization not found for invoice")
		tracing.TraceErr(span, err)
		return err
	}
	organization := mapper.MapDbNodeToOrganizationEntity(organizationNode)

	if organization.QuickbooksCustomerId == "" {
		quickbooksSaveCustomerResponse, err := c.quickbooksService.SaveCustomer(ctx, "", organization.Name)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if quickbooksSaveCustomerResponse.Customer == nil {
			err := errors.New("Quickbooks customer not found")
			tracing.TraceErr(span, err)
			return err
		}

		organization.QuickbooksCustomerId = quickbooksSaveCustomerResponse.Customer.Id

		err = c.neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonmodel.NodeLabelOrganization, organization.ID, string(neo4jentity.OrganizationPropertyQuickbooksCustomerId), organization.QuickbooksCustomerId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	savedInvoiced, err := c.quickbooksService.SaveInvoice(ctx, organization.QuickbooksCustomerId, invoice.PeriodStartDate, quickbooksInvoiceLines)
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

	organizationNode, err := c.neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if organizationNode == nil {
		err := errors.New("Organization not found for invoice")
		tracing.TraceErr(span, err)
		return err
	}
	organization := mapper.MapDbNodeToOrganizationEntity(organizationNode)

	paymentResponse, err := c.quickbooksService.PayInvoice(ctx, organization.QuickbooksCustomerId, invoice.QuickbooksInvoiceId, invoice.TotalAmount)
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
