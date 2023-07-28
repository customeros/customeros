package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"reflect"
)

func (i *Loaders) GetActionItemsForInteractionEvent(ctx context.Context, interactionEventId string) (*entity.ActionItemEntities, error) {
	thunk := i.ActionItemsForInteractionEvent.Load(ctx, dataloader.StringKey(interactionEventId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ActionItemEntities)
	return &resultObj, nil
}

func (b *actionItemBatcher) getActionItemsForInteractionEvents(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	actionItemsForNodes, err := b.actionItemService.GetActionItemsForNodes(ctx, repository.LINKED_WITH_INTERACTION_EVENT, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get attachments for interaction events")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	entitiesByInteractionEventId := make(map[string]entity.ActionItemEntities)
	for _, val := range *actionItemsForNodes {
		if list, ok := entitiesByInteractionEventId[val.DataloaderKey]; ok {
			entitiesByInteractionEventId[val.DataloaderKey] = append(list, val)
		} else {
			entitiesByInteractionEventId[val.DataloaderKey] = entity.ActionItemEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for id, record := range entitiesByInteractionEventId {
		if ix, ok := keyOrder[id]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, id)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.ActionItemEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.ActionItemEntities{})); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}
