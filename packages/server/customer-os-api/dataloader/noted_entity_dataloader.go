package dataloader

import (
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"golang.org/x/net/context"
	"time"
)

const notedEntityContextTimeout = 10 * time.Second

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
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, notedEntityContextTimeout)
	defer cancel()

	noteEntitiesPtr, err := b.noteService.GetNotedEntities(ctx, ids)
	if err != nil {
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

	return results
}
