package events_listeners

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
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
	spans, ctx := telemetry.StartListenerSpan(ctx, "SkuUpdateListener.Handle")
	defer spans.Finish()
	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	skuId := event.Event.EntityId
	spans.TagEntity(skuId)

	return l.handle(ctx, skuId)
}

func (l *SkuUpdateListener) handle(ctx context.Context, skuId string) error {
	spans, ctx := telemetry.StartListenerSpan(ctx, "SkuUpdateListener.handle")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	quickbooksSettingsEntity, err := l.dependencies.CommonServices.PostgresRepositories.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if quickbooksSettingsEntity == nil {
		spans.LogKV("error", "Quickbooks settings not found")
		return nil
	}

	skuEntity, err := l.dependencies.CommonServices.PostgresRepositories.SkuRepository.Get(ctx, tenant, skuId)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	qbProduct, err := l.dependencies.CommonServices.QuickbooksService.SaveProduct(ctx, skuEntity.QuickbooksId, skuEntity.Name, skuEntity.Archived, skuEntity.Price)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if qbProduct == nil || qbProduct.Product == nil {
		spans.LogObjectAsJson("qbProduct", qbProduct)
		err := errors.New("Quickbooks product could not be saved")
		spans.TraceError(err)
		return err
	}

	if skuEntity.QuickbooksId == "" {
		skuEntity.QuickbooksId = qbProduct.Product.Id
		_, err = l.dependencies.CommonServices.PostgresRepositories.SkuRepository.Save(ctx, skuEntity)
		if err != nil {
			spans.TraceError(err)
			return err
		}
	} else {
		spans.LogKV("info", "Quickbooks ID already exists")
	}

	return nil
}
