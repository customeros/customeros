package dataloader

import (
	"context"
	"errors"
	"reflect"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
)

func (i *Loaders) GetCountryForPhoneNumber(ctx context.Context, phoneNumberId string) (*neo4j_entity.CountryEntity, error) {
	thunk := i.CountryForPhoneNumber.Load(ctx, dataloader.StringKey(phoneNumberId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	resultObj := result.(*neo4j_entity.CountryEntity)
	return resultObj, nil
}

func (b *countryBatcher) getCountriesForPhoneNumbers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "CountryDataLoader.getCountriesForPhoneNumbers")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	countryEntities, err := b.countryService.GetCountriesForPhoneNumbers(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get countries for phoneNumbers")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	countryEntityByPhoneNumberId := make(map[string]neo4j_entity.CountryEntity)
	for _, val := range *countryEntities {
		countryEntityByPhoneNumberId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for phoneNumberId := range countryEntityByPhoneNumberId {
		if ix, ok := keyOrder[phoneNumberId]; ok {
			val := countryEntityByPhoneNumberId[phoneNumberId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, phoneNumberId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4j_entity.CountryEntity{}), true); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
