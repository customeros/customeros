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

func (i *Loaders) GetCreatedByParticipantsForMeeting(ctx context.Context, contactId string) (*entity.MeetingParticipants, error) {
	thunk := i.CreatedByParticipantsForMeeting.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.MeetingParticipants)
	return &resultObj, nil
}

func (i *Loaders) GetAttendedByParticipantsForMeeting(ctx context.Context, contactId string) (*entity.MeetingParticipants, error) {
	thunk := i.AttendedByParticipantsForMeeting.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.MeetingParticipants)
	return &resultObj, nil
}

func (b *meetingParticipantBatcher) getCreatedByParticipantsForMeeting(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingParticipantDataLoader.getCreatedByParticipantsForMeeting")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	participantEntitiesPtr, err := b.meetingService.GetParticipantsForMeetings(ctx, ids, entity.CREATED_BY)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get meeting created participants")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	participantEntitiesGrouped := make(map[string]entity.MeetingParticipants)
	for _, val := range *participantEntitiesPtr {
		if list, ok := participantEntitiesGrouped[val.GetDataloaderKey()]; ok {
			participantEntitiesGrouped[val.GetDataloaderKey()] = append(list, val)
		} else {
			participantEntitiesGrouped[val.GetDataloaderKey()] = entity.MeetingParticipants{val}
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
		results[ix] = &dataloader.Result{Data: entity.MeetingParticipants{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.MeetingParticipants{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *meetingParticipantBatcher) getAttendedByParticipantsForMeeting(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingParticipantDataLoader.getAttendedByParticipantsForMeeting")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	participantEntitiesPtr, err := b.meetingService.GetParticipantsForMeetings(ctx, ids, entity.ATTENDED_BY)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get meeting attended participants")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	participantEntitiesGrouped := make(map[string]entity.MeetingParticipants)
	for _, val := range *participantEntitiesPtr {
		if list, ok := participantEntitiesGrouped[val.GetDataloaderKey()]; ok {
			participantEntitiesGrouped[val.GetDataloaderKey()] = append(list, val)
		} else {
			participantEntitiesGrouped[val.GetDataloaderKey()] = entity.MeetingParticipants{val}
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
		results[ix] = &dataloader.Result{Data: entity.MeetingParticipants{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.MeetingParticipants{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
