package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func (i *Loaders) GetTimelineEventForTimelineEventId(ctx context.Context, timelineEventId string) (*entity.TimelineEvent, error) {
	thunk := i.TimelineEventForTimelineEventId.Load(ctx, dataloader.StringKey(timelineEventId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*entity.TimelineEvent), nil
}

func (b *timelineEventBatcher) getTimelineEventsForTimelineEventIds(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	timelineEventsPtr, err := b.timelineEventService.GetTimelineEventsWithIds(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get timeline events for timeline event ids")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	lastTouchpointTimelineEvent := make(map[string]entity.TimelineEvent)
	for _, val := range *timelineEventsPtr {
		lastTouchpointTimelineEvent[val.GetDataloaderKey()] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for timelineEventid, _ := range lastTouchpointTimelineEvent {
		if ix, ok := keyOrder[timelineEventid]; ok {
			val := lastTouchpointTimelineEvent[timelineEventid]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, timelineEventid)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	return results
}
