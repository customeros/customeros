package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
)

func (i *Loaders) GetEmailsForContact(ctx context.Context, contactId string) (*entity.EmailEntities, error) {
	thunk := i.EmailsForContact.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.EmailEntities)
	return &resultObj, nil
}

func (i *Loaders) GetEmailsForOrganization(ctx context.Context, organizationId string) (*entity.EmailEntities, error) {
	thunk := i.EmailsForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.EmailEntities)
	return &resultObj, nil
}

func (b *emailBatcher) getEmailsForContacts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailDataLoader.getEmailsForContacts", opentracing.ChildOf(tracing.ExtractSpanCtx(ctx)))
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	emailEntitiesPtr, err := b.emailService.GetAllForEntityTypeByIds(ctx, entity.CONTACT, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get emails for contacts")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	emailEntitiesGrouped := make(map[string]entity.EmailEntities)
	for _, val := range *emailEntitiesPtr {
		if list, ok := emailEntitiesGrouped[val.DataloaderKey]; ok {
			emailEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			emailEntitiesGrouped[val.DataloaderKey] = entity.EmailEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range emailEntitiesGrouped {
		ix, ok := keyOrder[contactId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.EmailEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.EmailEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("output - results_length", len(results)))

	return results
}

func (b *emailBatcher) getEmailsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailDataLoader.getEmailsForOrganizations", opentracing.ChildOf(tracing.ExtractSpanCtx(ctx)))
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	emailEntitiesPtr, err := b.emailService.GetAllForEntityTypeByIds(ctx, entity.ORGANIZATION, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get emails for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	emailEntitiesGrouped := make(map[string]entity.EmailEntities)
	for _, val := range *emailEntitiesPtr {
		if list, ok := emailEntitiesGrouped[val.DataloaderKey]; ok {
			emailEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			emailEntitiesGrouped[val.DataloaderKey] = entity.EmailEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range emailEntitiesGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.EmailEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.EmailEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("output - results_length", len(results)))

	return results
}
