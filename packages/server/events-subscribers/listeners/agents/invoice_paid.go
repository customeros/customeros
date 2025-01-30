package agent_listeners

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type InvoicePaidListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewInvoicePaidListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &InvoicePaidListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.InvoicePaid](), // subscribed event
			events.QueueAgents,                     // listening on Agents queue
		),
		dependencies: deps,
	}
}

// Add all Agent types subscribed to this event here
func (h *InvoicePaidListener) subscribedAgents() []enum.AgentType {
	return []enum.AgentType{}
}

func (l *InvoicePaidListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicePaidListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	invoiceId := event.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	return l.handle(ctx, invoiceId)
}

func (l *InvoicePaidListener) handle(ctx context.Context, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicePaidListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Missing tenant in context")
		tracing.TraceErr(span, err)
		return err
	}

	quickbooksSettingsEntity, err := l.dependencies.CommonServices.PostgresRepositories.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if quickbooksSettingsEntity == nil {
		span.LogFields(log.String("error", "Quickbooks settings not found"))
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

	paymentResponse, err := l.dependencies.CommonServices.QuickbooksService.PayInvoice(ctx, organization.QuickbooksCustomerId, invoice.QuickbooksInvoiceId, invoice.TotalAmount)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if paymentResponse == nil {
		err := errors.New("Invoice not paid in quickbooks")
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
