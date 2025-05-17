package dataloader

import (
	"context"
	"errors"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
	"reflect"
)

func (i *Loaders) GetPhoneNumbersForOrganization(ctx context.Context, organizationId string) (*neo4jentity.PhoneNumberEntities, error) {
	thunk := i.PhoneNumbersForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.PhoneNumberEntities)
	return &resultObj, nil
}

func (i *Loaders) GetPhoneNumbersForUser(ctx context.Context, userId string) (*neo4jentity.PhoneNumberEntities, error) {
	thunk := i.PhoneNumbersForUser.Load(ctx, dataloader.StringKey(userId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.PhoneNumberEntities)
	return &resultObj, nil
}

func (i *Loaders) GetPhoneNumbersForContact(ctx context.Context, contactId string) (*neo4jentity.PhoneNumberEntities, error) {
	thunk := i.PhoneNumbersForContact.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.PhoneNumberEntities)
	return &resultObj, nil
}

func (b *phoneNumberBatcher) getPhoneNumbersForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "PhoneNumberDataLoader.getPhoneNumbersForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	phoneNumberEntitiesPtr, err := b.phoneNumberService.GetAllForEntityTypeByIds(ctx, commonModel.ORGANIZATION, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get phone numbers for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	phoneNumberEntitiesGrouped := make(map[string]neo4jentity.PhoneNumberEntities)
	for _, val := range *phoneNumberEntitiesPtr {
		if list, ok := phoneNumberEntitiesGrouped[val.DataloaderKey]; ok {
			phoneNumberEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			phoneNumberEntitiesGrouped[val.DataloaderKey] = neo4jentity.PhoneNumberEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range phoneNumberEntitiesGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.PhoneNumberEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.PhoneNumberEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *phoneNumberBatcher) getPhoneNumbersForUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "PhoneNumberDataLoader.getPhoneNumbersForUsers")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	phoneNumberEntitiesPtr, err := b.phoneNumberService.GetAllForEntityTypeByIds(ctx, commonModel.USER, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get phone numbers for users")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	phoneNumberEntitiesGrouped := make(map[string]neo4jentity.PhoneNumberEntities)
	for _, val := range *phoneNumberEntitiesPtr {
		if list, ok := phoneNumberEntitiesGrouped[val.DataloaderKey]; ok {
			phoneNumberEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			phoneNumberEntitiesGrouped[val.DataloaderKey] = neo4jentity.PhoneNumberEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for userId, record := range phoneNumberEntitiesGrouped {
		ix, ok := keyOrder[userId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, userId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.PhoneNumberEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.PhoneNumberEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *phoneNumberBatcher) getPhoneNumbersForContacts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "PhoneNumberDataLoader.getPhoneNumbersForContacts")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	phoneNumberEntitiesPtr, err := b.phoneNumberService.GetAllForEntityTypeByIds(ctx, commonModel.CONTACT, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get phone numbers for users")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	phoneNumberEntitiesGrouped := make(map[string]neo4jentity.PhoneNumberEntities)
	for _, val := range *phoneNumberEntitiesPtr {
		if list, ok := phoneNumberEntitiesGrouped[val.DataloaderKey]; ok {
			phoneNumberEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			phoneNumberEntitiesGrouped[val.DataloaderKey] = neo4jentity.PhoneNumberEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range phoneNumberEntitiesGrouped {
		ix, ok := keyOrder[contactId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.PhoneNumberEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.PhoneNumberEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
