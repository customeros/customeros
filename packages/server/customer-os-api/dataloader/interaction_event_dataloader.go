package dataloader

import (
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"golang.org/x/net/context"
	"time"
)

const interactionEventContextTimeout = 10 * time.Second

func (i *Loaders) GetInteractionsEventsForInteractionsSession(ctx context.Context, ticketId string) (*entity.InteractionEventEntities, error) {
	thunk := i.InteractionEventsForInteractionSession.Load(ctx, dataloader.StringKey(ticketId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.InteractionEventEntities)
	return &resultObj, nil
}

func (b *interactionEventBatcher) getInteractionsEventsForInteractionsSessions(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, interactionEventContextTimeout)
	defer cancel()

	interactionEventEntitiesPtr, err := b.interactionEventService.GetInteractionEventsForInteractionSessions(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get notes for tickets")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	interactionEventEntitiesByInteractionSessionId := make(map[string]entity.InteractionEventEntities)
	for _, val := range *interactionEventEntitiesPtr {
		if list, ok := interactionEventEntitiesByInteractionSessionId[val.DataloaderKey]; ok {
			interactionEventEntitiesByInteractionSessionId[val.DataloaderKey] = append(list, val)
		} else {
			interactionEventEntitiesByInteractionSessionId[val.DataloaderKey] = entity.InteractionEventEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range interactionEventEntitiesByInteractionSessionId {
		if ix, ok := keyOrder[contactId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.NoteEntities{}, Error: nil}
	}

	return results
}
