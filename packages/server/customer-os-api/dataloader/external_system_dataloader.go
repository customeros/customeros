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

func (i *Loaders) GetExternalSystemsForComment(ctx context.Context, commentId string) (*entity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForComment.Load(ctx, dataloader.StringKey(commentId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForIssue(ctx context.Context, issueId string) (*entity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForIssue.Load(ctx, dataloader.StringKey(issueId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForOrganization(ctx context.Context, organizationId string) (*entity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForLogEntry(ctx context.Context, logEntryId string) (*entity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForLogEntry.Load(ctx, dataloader.StringKey(logEntryId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForMeeting(ctx context.Context, meetingId string) (*entity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForMeeting.Load(ctx, dataloader.StringKey(meetingId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForInteractionEvent(ctx context.Context, ieId string) (*entity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForInteractionEvent.Load(ctx, dataloader.StringKey(ieId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ExternalSystemEntities)
	return &resultObj, nil
}

func (b *externalSystemBatcher) getExternalSystemsForComments(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemDataLoader.getExternalSystemsForComments")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	return b.getExternalSystemsFor(ctx, keys, entity.COMMENT, span)
}

func (b *externalSystemBatcher) getExternalSystemsForIssues(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemDataLoader.getExternalSystemsForIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	return b.getExternalSystemsFor(ctx, keys, entity.ISSUE, span)
}

func (b *externalSystemBatcher) getExternalSystemsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemDataLoader.getExternalSystemsForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	return b.getExternalSystemsFor(ctx, keys, entity.ORGANIZATION, span)
}

func (b *externalSystemBatcher) getExternalSystemsForLogEntries(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemDataLoader.getExternalSystemsForLogEntries")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	return b.getExternalSystemsFor(ctx, keys, entity.LOG_ENTRY, span)
}

func (b *externalSystemBatcher) getExternalSystemsForMeetings(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemDataLoader.getExternalSystemsForMeetings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	return b.getExternalSystemsFor(ctx, keys, entity.MEETING, span)
}

func (b *externalSystemBatcher) getExternalSystemsForInteractionEvents(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemDataLoader.getExternalSystemsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	return b.getExternalSystemsFor(ctx, keys, entity.INTERACTION_EVENT, span)
}

func (b *externalSystemBatcher) getExternalSystemsFor(ctx context.Context, keys dataloader.Keys, entityType entity.EntityType, span opentracing.Span) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ExternalSystemsPtr, err := b.externalSystemService.GetExternalSystemsFor(ctx, ids, entityType)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get external systems for entities")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	ExternalSystemsByEntityId := make(map[string]entity.ExternalSystemEntities)
	for _, val := range *ExternalSystemsPtr {
		if list, ok := ExternalSystemsByEntityId[val.DataloaderKey]; ok {
			ExternalSystemsByEntityId[val.DataloaderKey] = append(list, val)
		} else {
			ExternalSystemsByEntityId[val.DataloaderKey] = entity.ExternalSystemEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for entityId, record := range ExternalSystemsByEntityId {
		if ix, ok := keyOrder[entityId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, entityId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.ExternalSystemEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.ExternalSystemEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("output - results_length", len(results)))

	return results
}

func (i *Loaders) GetExternalSystemsForContract(ctx context.Context, contractId string) (*entity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForContract.Load(ctx, dataloader.StringKey(contractId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForServiceLineItem(ctx context.Context, serviceLineItemId string) (*entity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForServiceLineItem.Load(ctx, dataloader.StringKey(serviceLineItemId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ExternalSystemEntities)
	return &resultObj, nil
}
