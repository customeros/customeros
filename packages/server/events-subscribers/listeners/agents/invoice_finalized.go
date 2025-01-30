package agent_listeners

import (
	"context"
	"errors"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	common_model "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type InvoiceFinalizedListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewInvoiceFinalizedListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &InvoiceFinalizedListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.InvoiceFinalized](), // subscribed event
			events.QueueAgents,                          // listening on Agents queue
		),
		dependencies: deps,
	}
}

// Add all Agent types subscribed to this event here
func (h *InvoicePaidListener) subscribedAgents() []enum.AgentType {
	return []enum.AgentType{}
}

// check if organization exists in quickbooks, if not create it
// check if all line items have skuid and the skuid has been pushed to quickbooks
// save invoice to quickbooks
func (l *InvoiceFinalizedListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceFinalizedListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Missing tenant in context")
		tracing.TraceErr(span, err)
		return err
	}

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	invoiceId := event.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	return l.handle(ctx, invoiceId)
}

func (l *InvoiceFinalizedListener) handle(ctx context.Context, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceFinalizedListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	quickbooksSettingsEntity, err := l.dependencies.CommonServices.PostgresRepositories.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if quickbooksSettingsEntity == nil {
		span.LogKV("error", "Quickbooks settings not found")
		return nil
	}

	invoice, err := l.dependencies.CommonServices.InvoiceService.GetById(ctx, nil, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if invoice == nil {
		err := errors.New("Invoice not found")
		tracing.TraceErr(span, err)
		return err
	}

	invoiceLines, err := l.dependencies.CommonServices.InvoiceService.GetInvoiceLinesForInvoices(ctx, []string{invoiceId})
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

		sku, err := l.dependencies.CommonServices.PostgresRepositories.SkuRepository.Get(ctx, tenant, invoiceLine.SkuId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if sku == nil {
			err := errors.New(fmt.Sprintf("Sku not found for invoice line %s", invoiceLine.Id))
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
		sku, err := l.dependencies.CommonServices.PostgresRepositories.SkuRepository.Get(ctx, tenant, invoiceLine.SkuId)
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

	organizationNode, err := l.dependencies.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoiceId)
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
		quickbooksSaveCustomerResponse, err := l.dependencies.CommonServices.QuickbooksService.SaveCustomer(ctx, "", organization.Name)
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

		err = l.dependencies.Neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, common_model.NodeLabelOrganization, organization.ID, string(neo4j_entity.OrganizationPropertyQuickbooksCustomerId), organization.QuickbooksCustomerId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	savedInvoiced, err := l.dependencies.CommonServices.QuickbooksService.SaveInvoice(ctx, organization.QuickbooksCustomerId, invoice.IssuedDate, quickbooksInvoiceLines)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if savedInvoiced == nil {
		err := errors.New("Invoice not saved in quickbooks")
		tracing.TraceErr(span, err)
		return err
	}

	err = l.dependencies.Neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, common_model.NodeLabelInvoice, invoice.Id, string(neo4j_entity.InvoicePropertyQuickbooksInvoiceId), savedInvoiced.Invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
