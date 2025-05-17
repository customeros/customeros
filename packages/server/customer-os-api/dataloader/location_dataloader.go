package dataloader

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
	"reflect"
)

func (i *Loaders) GetLocationsForContact(ctx context.Context, contactId string) (*neo4jentity.LocationEntities, error) {
	thunk := i.LocationsForContact.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.LocationEntities)
	return &resultObj, nil
}

func (i *Loaders) GetLocationsForOrganization(ctx context.Context, organizationId string) (*neo4jentity.LocationEntities, error) {
	thunk := i.LocationsForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.LocationEntities)
	return &resultObj, nil
}

func (b *locationBatcher) getLocationsForContacts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "LocationDataLoader.getLocationsForContacts")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	locationEntitiesPtr, err := b.locationCommonService.GetAllForContacts(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get locations for contacts")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	locationEntitiesGrouped := make(map[string]neo4jentity.LocationEntities)
	for _, val := range *locationEntitiesPtr {
		if list, ok := locationEntitiesGrouped[val.DataloaderKey]; ok {
			locationEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			locationEntitiesGrouped[val.DataloaderKey] = neo4jentity.LocationEntities{val}
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
		results[ix] = &dataloader.Result{Data: neo4jentity.LocationEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.LocationEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *locationBatcher) getLocationsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "LocationDataLoader.getLocationsForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	locationEntitiesPtr, err := b.locationCommonService.GetAllForOrganizations(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get locations for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	locationEntitiesGrouped := make(map[string]neo4jentity.LocationEntities)
	for _, val := range *locationEntitiesPtr {
		if list, ok := locationEntitiesGrouped[val.DataloaderKey]; ok {
			locationEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			locationEntitiesGrouped[val.DataloaderKey] = neo4jentity.LocationEntities{val}
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
		results[ix] = &dataloader.Result{Data: neo4jentity.LocationEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.LocationEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
