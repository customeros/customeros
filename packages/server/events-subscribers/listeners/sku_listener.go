package listeners

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/events-subscribers/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func OnSkuUpdate(ctx context.Context, dependencies *model.DependencyContainer, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnSkuUpdate")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message := input.(*dto.Event)
	skuId := message.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, skuId)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Missing tenant in context")
		tracing.TraceErr(span, err)
		return err
	}

	quickbooksSettingsEntity, err := dependencies.CommonServices.PostgresRepositories.QuickbooksSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if quickbooksSettingsEntity == nil {
		span.LogFields(log.String("error", "Quickbooks settings not found"))
		return nil
	}

	skuEntity, err := dependencies.CommonServices.PostgresRepositories.SkuRepository.Get(ctx, tenant, skuId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	qbProduct, err := dependencies.CommonServices.QuickbooksService.SaveProduct(ctx, skuEntity.QuickbooksId, skuEntity.Name, skuEntity.Archived)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if qbProduct == nil || qbProduct.Product == nil {
		span.LogFields(log.Object("qbProduct", qbProduct))
		err := errors.New("Quickbooks product could not be saved")
		tracing.TraceErr(span, err)
		return err
	}

	if skuEntity != nil && skuEntity.QuickbooksId == "" {
		skuEntity.QuickbooksId = qbProduct.Product.Id
		_, err = dependencies.CommonServices.PostgresRepositories.SkuRepository.Save(ctx, skuEntity)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		span.LogFields(log.String("info", "Quickbooks ID already exists"))
	}

	return nil
}
