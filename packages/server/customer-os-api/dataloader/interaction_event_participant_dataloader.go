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

func (i *Loaders) GetSentByParticipantsForInteractionEvent(ctx context.Context, contactId string) (*entity.InteractionEventParticipants, error) {
	thunk := i.SentByParticipantsForInteractionEvent.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.InteractionEventParticipants)
	return &resultObj, nil
}

func (i *Loaders) GetSentToParticipantsForInteractionEvent(ctx context.Context, organizationId string) (*entity.InteractionEventParticipants, error) {
	thunk := i.SentToParticipantsForInteractionEvent.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.InteractionEventParticipants)
	return &resultObj, nil
}

func (b *interactionEventParticipantBatcher) getSentByParticipantsForInteractionEvents(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventParticipantDataLoader.getSentByParticipantsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	participantEntitiesPtr, err := b.interactionEventService.GetSentByParticipantsForInteractionEvents(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get interaction event participants")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	participantEntitiesGrouped := make(map[string]entity.InteractionEventParticipants)
	for _, val := range *participantEntitiesPtr {
		if list, ok := participantEntitiesGrouped[val.GetDataloaderKey()]; ok {
			participantEntitiesGrouped[val.GetDataloaderKey()] = append(list, val)
		} else {
			participantEntitiesGrouped[val.GetDataloaderKey()] = entity.InteractionEventParticipants{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range participantEntitiesGrouped {
		ix, ok := keyOrder[contactId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.InteractionEventParticipants{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.InteractionEventParticipants{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *interactionEventParticipantBatcher) getSentToParticipantsForInteractionEvents(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventParticipantDataLoader.getSentToParticipantsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	participantEntitiesPtr, err := b.interactionEventService.GetSentToParticipantsForInteractionEvents(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get interaction event participants")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	participantEntitiesGrouped := make(map[string]entity.InteractionEventParticipants)
	for _, val := range *participantEntitiesPtr {
		if list, ok := participantEntitiesGrouped[val.GetDataloaderKey()]; ok {
			participantEntitiesGrouped[val.GetDataloaderKey()] = append(list, val)
		} else {
			participantEntitiesGrouped[val.GetDataloaderKey()] = entity.InteractionEventParticipants{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range participantEntitiesGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.InteractionEventParticipants{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.InteractionEventParticipants{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
