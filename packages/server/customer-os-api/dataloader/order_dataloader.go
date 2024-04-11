package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
)

func (i *Loaders) GetOrdersForOrganization(ctx context.Context, organizationId string) (*entity.OrderEntities, error) {
	thunk := i.OrdersForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.OrderEntities)
	return &resultObj, nil
}

func (b *orderBatcher) getOrdersForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderDataLoader.getOrdersForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	orderEntitiesPtr, err := b.orderService.GetAllForOrganizations(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get orders for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	orderEntitiesGrouped := make(map[string]entity.OrderEntities)
	for _, val := range *orderEntitiesPtr {
		if list, ok := orderEntitiesGrouped[val.DataloaderKey]; ok {
			orderEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			orderEntitiesGrouped[val.DataloaderKey] = entity.OrderEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range orderEntitiesGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.OrderEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.OrderEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
