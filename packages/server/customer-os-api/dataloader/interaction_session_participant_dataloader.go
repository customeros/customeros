package dataloader

import (
	"context"
	"errors"
	"fmt"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"reflect"
	"time"
)

const participantSessionContextTimeout = 10 * time.Second

func (i *Loaders) GetAttendedByParticipantsForInteractionSession(ctx context.Context, contactId string) (*entity.InteractionSessionParticipants, error) {
	thunk := i.AttendedByParticipantsForInteractionSession.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.InteractionSessionParticipants)
	return &resultObj, nil
}

func (b *interactionSessionParticipantBatcher) getAttendedByParticipantsForInteractionSessions(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, participantSessionContextTimeout)
	defer cancel()

	participantEntitiesPtr, err := b.interactionSessionService.GetAttendedByParticipantsForInteractionSessions(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get interaction event participants")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	participantEntitiesGrouped := make(map[string]entity.InteractionSessionParticipants)
	for _, val := range *participantEntitiesPtr {
		if list, ok := participantEntitiesGrouped[val.GetDataloaderKey()]; ok {
			participantEntitiesGrouped[val.GetDataloaderKey()] = append(list, val)
		} else {
			participantEntitiesGrouped[val.GetDataloaderKey()] = entity.InteractionSessionParticipants{val}
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
		results[ix] = &dataloader.Result{Data: entity.InteractionSessionParticipants{}, Error: nil}
	}

	if err = assertInteractionSessionParticipantsType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func assertInteractionSessionParticipantsType(results []*dataloader.Result) error {
	for _, res := range results {
		if _, ok := res.Data.(entity.InteractionSessionParticipants); !ok {
			return errors.New(fmt.Sprintf("Not expected type :%v", reflect.TypeOf(res.Data)))
		}
	}
	return nil
}
