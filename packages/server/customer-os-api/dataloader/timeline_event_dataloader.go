package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventDataLoader.getTimelineEventsForTimelineEventIds", opentracing.ChildOf(tracing.ExtractSpanCtx(ctx)))
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	timelineEventsPtr, err := b.timelineEventService.GetTimelineEventsWithIds(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
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

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
