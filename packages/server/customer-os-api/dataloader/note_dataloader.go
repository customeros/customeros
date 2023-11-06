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

func (i *Loaders) GetNotesForMeeting(ctx context.Context, meetingId string) (*entity.NoteEntities, error) {
	thunk := i.NotesForMeeting.Load(ctx, dataloader.StringKey(meetingId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.NoteEntities)
	return &resultObj, nil
}

func (b *noteBatcher) getNotesForMeetings(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteDataLoader.getNotesForMeetings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	noteEntitiesPtr, err := b.noteService.GetNotesForMeetings(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get noted entities for notes")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	notesForMeetings := make(map[string]entity.NoteEntities)
	for _, val := range *noteEntitiesPtr {
		if list, ok := notesForMeetings[val.GetDataloaderKey()]; ok {
			notesForMeetings[val.GetDataloaderKey()] = append(list, val)
		} else {
			notesForMeetings[val.GetDataloaderKey()] = entity.NoteEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range notesForMeetings {
		if ix, ok := keyOrder[contactId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.NoteEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.NoteEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
