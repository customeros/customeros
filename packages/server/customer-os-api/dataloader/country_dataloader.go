package dataloader

import (
	"context"
	"errors"
	"reflect"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryDataLoader.getCountriesForPhoneNumbers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	countryEntities, err := b.countryService.GetCountriesForPhoneNumbers(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("result.length", len(results)))

	return results
}
