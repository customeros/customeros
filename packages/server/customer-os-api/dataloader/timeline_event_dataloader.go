package dataloader

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/graph-gophers/dataloader"
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

func (i *Loaders) GetInboundCommsCountForOrganization(ctx context.Context, organizationId string) (int64, error) {
	thunk := i.InboundCommsCountForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, nil
	}
	return *result.(*int64), nil
}

func (i *Loaders) GetOutboundCommsCountForOrganization(ctx context.Context, organizationId string) (int64, error) {
	thunk := i.OutboundCommsCountForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, nil
	}
	return *result.(*int64), nil
}

func (b *timelineEventBatcher) getTimelineEventsForTimelineEventIds(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TimelineEventDataLoader.getTimelineEventsForTimelineEventIds")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	timelineEventsPtr, err := b.timelineEventService.GetTimelineEventsWithIds(ctx, ids)
	if err != nil {
		spans.TraceError(err)
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

	spans.LogKV("result.length", len(results))

	return results
}

func (b *timelineEventBatcher) getInboundCommsCountForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TimelineEventBatcher.getInboundCommsCountForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	countsPerOrg, err := b.timelineEventService.GetInboundCommsCountCountByOrganizations(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get inbound comms count for organization")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for orgId, _ := range countsPerOrg {
		if ix, ok := keyOrder[orgId]; ok {
			val := countsPerOrg[orgId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, orgId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: 0, Error: nil}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *timelineEventBatcher) getOutboundCommsCountForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TimelineEventBatcher.getInboundCommsCountForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	countsPerOrg, err := b.timelineEventService.GetOutboundCommsCountCountByOrganizations(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get outbound comms count for organization")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for orgId, _ := range countsPerOrg {
		if ix, ok := keyOrder[orgId]; ok {
			val := countsPerOrg[orgId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, orgId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: 0, Error: nil}
	}

	spans.LogKV("result.length", len(results))

	return results
}
