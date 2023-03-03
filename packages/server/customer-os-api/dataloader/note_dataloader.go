package dataloader

import (
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"golang.org/x/net/context"
	"time"
)

const noteContextTimeout = 10 * time.Second

func (i *Loaders) GetNotesForTicket(ctx context.Context, ticketId string) (*entity.NoteEntities, error) {
	thunk := i.NotesForTicket.Load(ctx, dataloader.StringKey(ticketId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.NoteEntities)
	return &resultObj, nil
}

func (b *noteBatcher) getNotesForTickets(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, noteContextTimeout)
	defer cancel()

	noteEntitiesPtr, err := b.noteService.GetNotesForTickets(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get notes for tickets")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	noteEntitiesByTicketId := make(map[string]entity.NoteEntities)
	for _, val := range *noteEntitiesPtr {
		if list, ok := noteEntitiesByTicketId[val.DataloaderKey]; ok {
			noteEntitiesByTicketId[val.DataloaderKey] = append(list, val)
		} else {
			noteEntitiesByTicketId[val.DataloaderKey] = entity.NoteEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range noteEntitiesByTicketId {
		if ix, ok := keyOrder[contactId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.NoteEntities{}, Error: nil}
	}

	return results
}
