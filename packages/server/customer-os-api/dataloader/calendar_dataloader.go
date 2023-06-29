package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"reflect"
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
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	calendarEntitiesPtr, err := b.calendarService.GetAllForUsers(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get calendars for users")}}
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
		return []*dataloader.Result{{nil, err}}
	}

	return results
}
