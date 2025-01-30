package agent_listeners

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type InvoiceVoidedListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewInvoiceVoidedListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &InvoiceVoidedListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.InvoiceVoided](), // subscribed event
			events.QueueAgents,                       // listening on Agents queue
		),
		dependencies: deps,
	}
}

// Add all Agent types subscribed to this event here
func (h *InvoiceVoidedListener) subscribedAgents() []enum.AgentType {
	return []enum.AgentType{}
}

func (l *InvoiceVoidedListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceVoidedListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = l.validateMessage(ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	invoiceId := event.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	return l.handle(ctx, invoiceId)
}

func (l *InvoiceVoidedListener) handle(ctx context.Context, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceVoidedListener.handle")
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

	voidedResponse, err := l.dependencies.CommonServices.QuickbooksService.VoidInvoice(ctx, invoice.QuickbooksInvoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if voidedResponse == nil {
		err := errors.New("Invoice not voided in quickbooks")
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (l *InvoiceVoidedListener) validateMessage(ctx context.Context, event *dto.Event) (*dto.InvoiceVoided, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceVoidedListener.validateMessage")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message, ok := event.Event.Data.(*dto.InvoiceVoided)
	if !ok {
		err := fmt.Errorf("expected InvoiceVoided, got %T", event.Event.Data)
		tracing.TraceErr(span, err)
		return nil, err
	}

	return message, nil
}
