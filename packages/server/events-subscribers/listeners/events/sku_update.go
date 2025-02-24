package events_listeners

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
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

func NewSkuUpdateListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &SkuUpdateListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.SkuUpdate](), // subscribed event
			events.QueueEvents,                   // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *SkuUpdateListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SkuUpdateListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	skuId := event.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, skuId)

	return l.handle(ctx, skuId)
}

func (l *SkuUpdateListener) handle(ctx context.Context, skuId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SkuUpdateListener.handle")
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

	qbProduct, err := l.dependencies.CommonServices.QuickbooksService.SaveProduct(ctx, skuEntity.QuickbooksId, skuEntity.Name, skuEntity.Archived, skuEntity.Price)
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

	if skuEntity.QuickbooksId == "" {
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
