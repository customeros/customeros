package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
)

func (i *Loaders) GetExternalSystemsForEntity(ctx context.Context, entityId string) (*entity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForEntity.Load(ctx, dataloader.StringKey(entityId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ExternalSystemEntities)
	return &resultObj, nil
}

func (b *externalSystemBatcher) getExternalSystemsForEntities(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemDataLoader.getExternalSystemsForEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ExternalSystemsPtr, err := b.externalSystemService.GetExternalSystemsForEntities(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get external systems for entities")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	ExternalSystemsByEntityId := make(map[string]entity.ExternalSystemEntities)
	for _, val := range *ExternalSystemsPtr {
		if list, ok := ExternalSystemsByEntityId[val.DataloaderKey]; ok {
			ExternalSystemsByEntityId[val.DataloaderKey] = append(list, val)
		} else {
			ExternalSystemsByEntityId[val.DataloaderKey] = entity.ExternalSystemEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for entityId, record := range ExternalSystemsByEntityId {
		if ix, ok := keyOrder[entityId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, entityId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.ExternalSystemEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.ExternalSystemEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("output - results_length", len(results)))

	return results
}
