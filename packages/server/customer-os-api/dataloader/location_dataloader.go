package dataloader

import (
	"context"
	"errors"
	"fmt"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"reflect"
	"time"
)

const locationContextTimeout = 10 * time.Second

func (i *Loaders) GetLocationsForContact(ctx context.Context, contactId string) (*entity.LocationEntities, error) {
	thunk := i.LocationsForContact.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.LocationEntities)
	return &resultObj, nil
}

func (i *Loaders) GetLocationsForOrganization(ctx context.Context, organizationId string) (*entity.LocationEntities, error) {
	thunk := i.LocationsForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.LocationEntities)
	return &resultObj, nil
}

func (b *locationBatcher) getLocationsForContacts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, locationContextTimeout)
	defer cancel()

	locationEntitiesPtr, err := b.locationService.GetAllForContacts(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get locations for contacts")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	locationEntitiesGrouped := make(map[string]entity.LocationEntities)
	for _, val := range *locationEntitiesPtr {
		if list, ok := locationEntitiesGrouped[val.DataloaderKey]; ok {
			locationEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			locationEntitiesGrouped[val.DataloaderKey] = entity.LocationEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range locationEntitiesGrouped {
		ix, ok := keyOrder[contactId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.LocationEntities{}, Error: nil}
	}

	if err = assertLocationEntitiesType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func (b *locationBatcher) getLocationsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, locationContextTimeout)
	defer cancel()

	locationEntitiesPtr, err := b.locationService.GetAllForOrganizations(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get locations for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	locationEntitiesGrouped := make(map[string]entity.LocationEntities)
	for _, val := range *locationEntitiesPtr {
		if list, ok := locationEntitiesGrouped[val.DataloaderKey]; ok {
			locationEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			locationEntitiesGrouped[val.DataloaderKey] = entity.LocationEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range locationEntitiesGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.LocationEntities{}, Error: nil}
	}

	if err = assertLocationEntitiesType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func assertLocationEntitiesType(results []*dataloader.Result) error {
	for _, res := range results {
		if _, ok := res.Data.(entity.LocationEntities); !ok {
			return errors.New(fmt.Sprintf("Not expected type :%v", reflect.TypeOf(res.Data)))
		}
	}
	return nil
}
