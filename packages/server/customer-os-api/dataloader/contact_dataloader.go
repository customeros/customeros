package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
)

func (i *Loaders) GetContactsForEmail(ctx context.Context, emailId string) (*entity.ContactEntities, error) {
	thunk := i.ContactsForEmail.Load(ctx, dataloader.StringKey(emailId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ContactEntities)
	return &resultObj, nil
}

func (i *Loaders) GetContactsForPhoneNumber(ctx context.Context, phoneNumberId string) (*entity.ContactEntities, error) {
	thunk := i.ContactsForPhoneNumber.Load(ctx, dataloader.StringKey(phoneNumberId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ContactEntities)
	return &resultObj, nil
}

func (i *Loaders) GetContactForJobRole(ctx context.Context, jobRoleId string) (*entity.ContactEntity, error) {
	thunk := i.ContactForJobRole.Load(ctx, dataloader.StringKey(jobRoleId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*entity.ContactEntity), nil
}

func (b *contactBatcher) getContactsForEmails(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactDataLoader.getContactsForEmails")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	contactEntitiesPtr, err := b.contactService.GetContactsForEmails(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get contacts for emails")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	contactEntitiesByEmailId := make(map[string]entity.ContactEntities)
	for _, val := range *contactEntitiesPtr {
		if list, ok := contactEntitiesByEmailId[val.DataloaderKey]; ok {
			contactEntitiesByEmailId[val.DataloaderKey] = append(list, val)
		} else {
			contactEntitiesByEmailId[val.DataloaderKey] = entity.ContactEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for emailId, record := range contactEntitiesByEmailId {
		if ix, ok := keyOrder[emailId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, emailId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.ContactEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.ContactEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Object("output - results_length", len(results)))

	return results
}

func (b *contactBatcher) getContactsForPhoneNumbers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactDataLoader.getContactsForPhoneNumbers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	contactEntitiesPtr, err := b.contactService.GetContactsForPhoneNumbers(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get contacts for phone numbers")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	contactEntitiesByPhoneNumberId := make(map[string]entity.ContactEntities)
	for _, val := range *contactEntitiesPtr {
		if list, ok := contactEntitiesByPhoneNumberId[val.DataloaderKey]; ok {
			contactEntitiesByPhoneNumberId[val.DataloaderKey] = append(list, val)
		} else {
			contactEntitiesByPhoneNumberId[val.DataloaderKey] = entity.ContactEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for phoneNumberId, record := range contactEntitiesByPhoneNumberId {
		if ix, ok := keyOrder[phoneNumberId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, phoneNumberId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.ContactEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.ContactEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Object("output - results_length", len(results)))

	return results
}

func (b *contactBatcher) getContactsForJobRoles(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactDataLoader.getContactsForJobRoles")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	contactEntities, err := b.contactService.GetContactsForJobRoles(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get contacts for job roles")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	contactEntityByJobRoleId := make(map[string]entity.ContactEntity)
	for _, val := range *contactEntities {
		contactEntityByJobRoleId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for jobRoleId, _ := range contactEntityByJobRoleId {
		if ix, ok := keyOrder[jobRoleId]; ok {
			val := contactEntityByJobRoleId[jobRoleId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, jobRoleId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(entity.ContactEntity{}), true); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Object("output - results_length", len(results)))

	return results
}
