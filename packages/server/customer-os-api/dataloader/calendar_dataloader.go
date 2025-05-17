package dataloader

import (
	"context"
	"reflect"

	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"
)

func (i *Loaders) GetCalendarsForUser(ctx context.Context, userId string) (*entity.CalendarEntities, error) {
	thunk := i.CalendarsForUser.Load(ctx, dataloader.StringKey(userId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.CalendarEntities)
	return &resultObj, nil
}

func (b *calendarBatcher) getCalendarsForUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "CalendarDataLoader.getCalendarsForUsers")
	defer spans.Finish()
	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	calendarEntitiesPtr, err := b.calendarService.GetAllForUsers(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	calendarEntitiesGroupedByUserId := make(map[string]entity.CalendarEntities)
	for _, val := range *calendarEntitiesPtr {
		if list, ok := calendarEntitiesGroupedByUserId[val.DataloaderKey]; ok {
			calendarEntitiesGroupedByUserId[val.DataloaderKey] = append(list, val)
		} else {
			calendarEntitiesGroupedByUserId[val.DataloaderKey] = entity.CalendarEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for userId, record := range calendarEntitiesGroupedByUserId {
		ix, ok := keyOrder[userId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, userId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.CalendarEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.CalendarEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
