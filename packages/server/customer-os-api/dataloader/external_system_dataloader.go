package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
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
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	ExternalSystemsPtr, err := b.externalSystemService.GetExternalSystemsForEntities(ctx, ids)
	if err != nil {
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
		return []*dataloader.Result{{nil, err}}
	}

	return results
}
