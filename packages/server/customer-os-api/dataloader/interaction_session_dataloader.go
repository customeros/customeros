package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"reflect"
	"time"
)

const interactionSessionContextTimeout = 10 * time.Second

func (i *Loaders) GetInteractionSessionForInteractionEvent(ctx context.Context, interactionEventId string) (*entity.InteractionSessionEntity, error) {
	thunk := i.InteractionSessionForInteractionEvent.Load(ctx, dataloader.StringKey(interactionEventId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	resultObj := result.(*entity.InteractionSessionEntity)
	return resultObj, nil
}

func (b *interactionSessionBatcher) getInteractionSessionsForInteractionEvents(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, interactionSessionContextTimeout)
	defer cancel()

	interactionSessionEntities, err := b.interactionSessionService.GetInteractionEventsForInteractionSessions(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get interaction sessions for interaction events")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	interactionSessionEntityByInteractionEventId := make(map[string]entity.InteractionSessionEntity)
	for _, val := range *interactionSessionEntities {
		interactionSessionEntityByInteractionEventId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for interactionEventId, _ := range interactionSessionEntityByInteractionEventId {
		if ix, ok := keyOrder[interactionEventId]; ok {
			val := interactionSessionEntityByInteractionEventId[interactionEventId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, interactionEventId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(entity.InteractionSessionEntity{}), true); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}
