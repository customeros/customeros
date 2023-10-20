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

func (i *Loaders) GetMentionedEntitiesForNote(ctx context.Context, noteId string) (*entity.MentionedEntities, error) {
	thunk := i.MentionedEntitiesForNote.Load(ctx, dataloader.StringKey(noteId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.MentionedEntities)
	return &resultObj, nil
}

func (b *mentionedEntityBatcher) getMentionedEntitiesForNotes(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MentionedEntityDataLoader.getMentionedEntitiesForNotes", opentracing.ChildOf(tracing.ExtractSpanCtx(ctx)))
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	noteEntitiesPtr, err := b.noteService.GetMentionedEntities(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get mentioned entities for notes")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	mentionedEntitiesByNoteId := make(map[string]entity.MentionedEntities)
	for _, val := range *noteEntitiesPtr {
		if list, ok := mentionedEntitiesByNoteId[val.GetDataloaderKey()]; ok {
			mentionedEntitiesByNoteId[val.GetDataloaderKey()] = append(list, val)
		} else {
			mentionedEntitiesByNoteId[val.GetDataloaderKey()] = entity.MentionedEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range mentionedEntitiesByNoteId {
		if ix, ok := keyOrder[contactId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.MentionedEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.MentionedEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
