package listeners

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/events-subscribers/model"
	"github.com/opentracing/opentracing-go"
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

	if skuEntity.QuickbooksId == "" {
		skuEntity.QuickbooksId = qbProduct.Product.Id
		_, err = dependencies.CommonServices.PostgresRepositories.SkuRepository.Save(ctx, skuEntity)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
