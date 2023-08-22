package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
)

func (i *Loaders) GetAttachmentsForInteractionEvent(ctx context.Context, interactionEventId string) (*entity.AttachmentEntities, error) {
	thunk := i.AttachmentsForInteractionEvent.Load(ctx, dataloader.StringKey(interactionEventId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.AttachmentEntities)
	return &resultObj, nil
}

func (i *Loaders) GetAttachmentsForInteractionSession(ctx context.Context, interactionSessionId string) (*entity.AttachmentEntities, error) {
	thunk := i.AttachmentsForInteractionSession.Load(ctx, dataloader.StringKey(interactionSessionId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.AttachmentEntities)
	return &resultObj, nil
}

func (i *Loaders) GetAttachmentsForMeeting(ctx context.Context, meetingId string) (*entity.AttachmentEntities, error) {
	thunk := i.AttachmentsForMeeting.Load(ctx, dataloader.StringKey(meetingId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.AttachmentEntities)
	return &resultObj, nil
}

func (b *attachmentBatcher) getAttachmentsForInteractionEvents(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentDataLoader.getAttachmentsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	attachmentEntitiesPtr, err := b.attachmentService.GetAttachmentsForNode(ctx, repository.LINKED_WITH_INTERACTION_EVENT, nil, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get attachments for interaction events")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	attachmentEntitiesByInteractionEventId := make(map[string]entity.AttachmentEntities)
	for _, val := range *attachmentEntitiesPtr {
		if list, ok := attachmentEntitiesByInteractionEventId[val.DataloaderKey]; ok {
			attachmentEntitiesByInteractionEventId[val.DataloaderKey] = append(list, val)
		} else {
			attachmentEntitiesByInteractionEventId[val.DataloaderKey] = entity.AttachmentEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range attachmentEntitiesByInteractionEventId {
		if ix, ok := keyOrder[organizationId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.AttachmentEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.AttachmentEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("output - results_length", len(results)))

	return results
}

func (b *attachmentBatcher) getAttachmentsForInteractionSessions(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentDataLoader.getAttachmentsForInteractionSessions")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	attachmentEntitiesPtr, err := b.attachmentService.GetAttachmentsForNode(ctx, repository.LINKED_WITH_INTERACTION_SESSION, nil, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get attachments for interaction sessions")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	attachmentEntitiesByInteractionSessionId := make(map[string]entity.AttachmentEntities)
	for _, val := range *attachmentEntitiesPtr {
		if list, ok := attachmentEntitiesByInteractionSessionId[val.DataloaderKey]; ok {
			attachmentEntitiesByInteractionSessionId[val.DataloaderKey] = append(list, val)
		} else {
			attachmentEntitiesByInteractionSessionId[val.DataloaderKey] = entity.AttachmentEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range attachmentEntitiesByInteractionSessionId {
		if ix, ok := keyOrder[organizationId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.AttachmentEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.AttachmentEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("output - results_length", len(results)))

	return results
}

func (b *attachmentBatcher) getAttachmentsForMeetings(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentDataLoader.getAttachmentsForMeetings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	attachmentEntitiesPtr, err := b.attachmentService.GetAttachmentsForNode(ctx, repository.LINKED_WITH_MEETING, nil, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get attachments for interaction sessions")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	attachmentEntitiesByMeetingId := make(map[string]entity.AttachmentEntities)
	for _, val := range *attachmentEntitiesPtr {
		if list, ok := attachmentEntitiesByMeetingId[val.DataloaderKey]; ok {
			attachmentEntitiesByMeetingId[val.DataloaderKey] = append(list, val)
		} else {
			attachmentEntitiesByMeetingId[val.DataloaderKey] = entity.AttachmentEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range attachmentEntitiesByMeetingId {
		if ix, ok := keyOrder[organizationId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.AttachmentEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.AttachmentEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("output - results_length", len(results)))

	return results
}
