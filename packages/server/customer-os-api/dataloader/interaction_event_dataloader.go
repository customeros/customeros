package dataloader

import (
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"golang.org/x/net/context"
	"time"
)

const interactionEventContextTimeout = 10 * time.Second

func (i *Loaders) GetInteractionEventsForInteractionSession(ctx context.Context, interactionSessionId string) (*entity.InteractionEventEntities, error) {
	thunk := i.InteractionEventsForInteractionSession.Load(ctx, dataloader.StringKey(interactionSessionId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.InteractionEventEntities)
	return &resultObj, nil
}

func (i *Loaders) GetInteractionEventsForMeeting(ctx context.Context, meetingId string) (*entity.InteractionEventEntities, error) {
	thunk := i.InteractionEventsForMeeting.Load(ctx, dataloader.StringKey(meetingId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.InteractionEventEntities)
	return &resultObj, nil
}

func (i *Loaders) GetInteractionEventsForInteractionEvent(ctx context.Context, interactionEventId string) (*entity.InteractionEventEntities, error) {
	thunk := i.ReplyToInteractionEventForInteractionEvent.Load(ctx, dataloader.StringKey(interactionEventId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.InteractionEventEntities)
	return &resultObj, nil
}

func (b *interactionEventBatcher) getInteractionEventsForInteractionSessions(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, interactionEventContextTimeout)
	defer cancel()

	interactionEventEntitiesPtr, err := b.interactionEventService.GetInteractionEventsForInteractionSessions(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get interaction events for interaction sessions")}}
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
	for interactionSessionId, record := range interactionEventEntitiesByInteractionSessionId {
		if ix, ok := keyOrder[interactionSessionId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, interactionSessionId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.InteractionEventEntities{}, Error: nil}
	}

	return results
}

func (b *interactionEventBatcher) getInteractionEventsForMeetings(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, interactionEventContextTimeout)
	defer cancel()

	interactionEventEntitiesPtr, err := b.interactionEventService.GetInteractionEventsForMeetings(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get interaction events for meetings")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	interactionEventEntitiesByMeetingId := make(map[string]entity.InteractionEventEntities)
	for _, val := range *interactionEventEntitiesPtr {
		if list, ok := interactionEventEntitiesByMeetingId[val.DataloaderKey]; ok {
			interactionEventEntitiesByMeetingId[val.DataloaderKey] = append(list, val)
		} else {
			interactionEventEntitiesByMeetingId[val.DataloaderKey] = entity.InteractionEventEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for meetingsId, record := range interactionEventEntitiesByMeetingId {
		if ix, ok := keyOrder[meetingsId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, meetingsId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.InteractionEventEntities{}, Error: nil}
	}

	return results
}

func (b *interactionEventBatcher) getReplyToInteractionEventsForInteractionEvents(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, participantContextTimeout)
	defer cancel()

	interactionEventEntitiesPtr, err := b.interactionEventService.GetReplyToInteractionsEventForInteractionEvents(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get interaction event participants")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	interactionEventEntitiesGrouped := make(map[string]entity.InteractionEventEntities)
	for _, val := range *interactionEventEntitiesPtr {
		if list, ok := interactionEventEntitiesGrouped[val.DataloaderKey]; ok {
			interactionEventEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			interactionEventEntitiesGrouped[val.DataloaderKey] = entity.InteractionEventEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range interactionEventEntitiesGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.InteractionEventEntities{}, Error: nil}
	}

	return results
}
