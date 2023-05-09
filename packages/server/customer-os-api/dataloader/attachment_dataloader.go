package dataloader

import (
	"errors"
	"fmt"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"golang.org/x/net/context"
	"reflect"
	"time"
)

const attachmentContextTimeout = 10 * time.Second

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
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, attachmentContextTimeout)
	defer cancel()

	attachmentEntitiesPtr, err := b.attachmentService.GetAttachmentsForNode(ctx, repository.INCLUDED_BY_INTERACTION_EVENT, nil, ids)
	if err != nil {
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

	if err = assertAttachmentEntitiesType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func (b *attachmentBatcher) getAttachmentsForInteractionSessions(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, attachmentContextTimeout)
	defer cancel()

	attachmentEntitiesPtr, err := b.attachmentService.GetAttachmentsForNode(ctx, repository.INCLUDED_BY_INTERACTION_SESSION, nil, ids)
	if err != nil {
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

	if err = assertAttachmentEntitiesType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func (b *attachmentBatcher) getAttachmentsForMeetings(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, attachmentContextTimeout)
	defer cancel()

	attachmentEntitiesPtr, err := b.attachmentService.GetAttachmentsForNode(ctx, repository.INCLUDED_BY_MEETING, nil, ids)
	if err != nil {
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

	if err = assertAttachmentEntitiesType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func assertAttachmentEntitiesType(results []*dataloader.Result) error {
	for _, res := range results {
		if _, ok := res.Data.(entity.AttachmentEntities); !ok {
			return errors.New(fmt.Sprintf("Not expected type :%v", reflect.TypeOf(res.Data)))
		}
	}
	return nil
}
