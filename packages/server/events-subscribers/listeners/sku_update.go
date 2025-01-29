package listeners

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type SkuUpdateListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewSkuUpdateListener(logger logger.Logger, deps *model.DependencyContainer) events.EventListener {
	return &SkuUpdateListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.HideContact](), // subscribed event
			events.QueueEvents,                     // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *SkuUpdateListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SkuUpdateListener.Handle")
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

	_, err = l.validateMessage(ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	skuId := event.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, skuId)

	return l.onSkuUpdate(ctx, skuId)
}

func (l *SkuUpdateListener) onSkuUpdate(ctx context.Context, skuId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SkuUpdateListener.onSkuUpdate")
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

	skuEntity, err := l.dependencies.CommonServices.PostgresRepositories.SkuRepository.Get(ctx, tenant, skuId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	qbProduct, err := l.dependencies.CommonServices.QuickbooksService.SaveProduct(ctx, skuEntity.QuickbooksId, skuEntity.Name, skuEntity.Archived)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if qbProduct == nil || qbProduct.Product == nil {
		tracing.LogObjectAsJson(span, "qbProduct", qbProduct)
		err := errors.New("Quickbooks product could not be saved")
		tracing.TraceErr(span, err)
		return err
	}

	if skuEntity != nil && skuEntity.QuickbooksId == "" {
		skuEntity.QuickbooksId = qbProduct.Product.Id
		_, err = l.dependencies.CommonServices.PostgresRepositories.SkuRepository.Save(ctx, skuEntity)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		span.LogKV("info", "Quickbooks ID already exists")
	}

	return nil
}

func (l *SkuUpdateListener) validateMessage(ctx context.Context, event *dto.Event) (*dto.HideContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SkuUpdateListener.validateMessage")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message, ok := event.Event.Data.(*dto.HideContact)
	if !ok {
		err := fmt.Errorf("expected RequestEnrichContact, got %T", event.Event.Data)
		tracing.TraceErr(span, err)
		return nil, err
	}

	if event.Event.EntityId == "" {
		err := errors.New("EntityId not set on event")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return message, nil
}
