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

func (i *Loaders) GetNotedEntitiesForNote(ctx context.Context, noteId string) (*entity.NotedEntities, error) {
	thunk := i.NotedEntitiesForNote.Load(ctx, dataloader.StringKey(noteId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.NotedEntities)
	return &resultObj, nil
}

func (b *notedEntityBatcher) getNotedEntitiesForNotes(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NotedEntityDataLoader.getNotedEntitiesForNotes")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	noteEntitiesPtr, err := b.noteService.GetNotedEntities(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get noted entities for notes")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	notedEntitiesByNoteId := make(map[string]entity.NotedEntities)
	for _, val := range *noteEntitiesPtr {
		if list, ok := notedEntitiesByNoteId[val.GetDataloaderKey()]; ok {
			notedEntitiesByNoteId[val.GetDataloaderKey()] = append(list, val)
		} else {
			notedEntitiesByNoteId[val.GetDataloaderKey()] = entity.NotedEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range notedEntitiesByNoteId {
		if ix, ok := keyOrder[contactId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.NotedEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.NotedEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
