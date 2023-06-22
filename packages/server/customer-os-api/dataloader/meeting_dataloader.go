package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"reflect"
)

func (i *Loaders) GetMeetingForInteractionEvent(ctx context.Context, interactionEventId string) (*entity.MeetingEntity, error) {
	thunk := i.MeetingForInteractionEvent.Load(ctx, dataloader.StringKey(interactionEventId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	resultObj := result.(*entity.MeetingEntity)
	return resultObj, nil
}

func (b *meetingBatcher) getMeetingsForInteractionEvents(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	meetingEntities, err := b.meetingService.GetMeetingsForInteractionEvents(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get interaction sessions for interaction events")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	meetingEntityByInteractionEventId := make(map[string]entity.MeetingEntity)
	for _, val := range *meetingEntities {
		meetingEntityByInteractionEventId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for interactionEventId, _ := range meetingEntityByInteractionEventId {
		if ix, ok := keyOrder[interactionEventId]; ok {
			val := meetingEntityByInteractionEventId[interactionEventId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, interactionEventId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(entity.MeetingEntity{}), true); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}
