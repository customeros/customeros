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

const contactContextTimeout = 10 * time.Second

func (i *Loaders) GetContactsForEmail(ctx context.Context, emailId string) (*entity.ContactEntities, error) {
	thunk := i.ContactsForEmail.Load(ctx, dataloader.StringKey(emailId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ContactEntities)
	return &resultObj, nil
}

func (i *Loaders) GetContactsForPhoneNumber(ctx context.Context, phoneNumberId string) (*entity.ContactEntities, error) {
	thunk := i.ContactsForPhoneNumber.Load(ctx, dataloader.StringKey(phoneNumberId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ContactEntities)
	return &resultObj, nil
}

func (b *contactBatcher) getContactsForEmails(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, contactContextTimeout)
	defer cancel()

	contactEntitiesPtr, err := b.contactService.GetContactsForEmails(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get contacts for emails")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	contactEntitiesByEmailId := make(map[string]entity.ContactEntities)
	for _, val := range *contactEntitiesPtr {
		if list, ok := contactEntitiesByEmailId[val.DataloaderKey]; ok {
			contactEntitiesByEmailId[val.DataloaderKey] = append(list, val)
		} else {
			contactEntitiesByEmailId[val.DataloaderKey] = entity.ContactEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for emailId, record := range contactEntitiesByEmailId {
		if ix, ok := keyOrder[emailId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, emailId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.ContactEntities{}, Error: nil}
	}

	if err = assertContactEntitiesType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func (b *contactBatcher) getContactsForPhoneNumbers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, contactContextTimeout)
	defer cancel()

	contactEntitiesPtr, err := b.contactService.GetContactsForPhoneNumbers(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get contacts for phone numbers")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	contactEntitiesByPhoneNumberId := make(map[string]entity.ContactEntities)
	for _, val := range *contactEntitiesPtr {
		if list, ok := contactEntitiesByPhoneNumberId[val.DataloaderKey]; ok {
			contactEntitiesByPhoneNumberId[val.DataloaderKey] = append(list, val)
		} else {
			contactEntitiesByPhoneNumberId[val.DataloaderKey] = entity.ContactEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for phoneNumberId, record := range contactEntitiesByPhoneNumberId {
		if ix, ok := keyOrder[phoneNumberId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, phoneNumberId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.ContactEntities{}, Error: nil}
	}

	if err = assertContactEntitiesType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func assertContactEntitiesType(results []*dataloader.Result) error {
	for _, res := range results {
		if _, ok := res.Data.(entity.ContactEntities); !ok {
			return errors.New(fmt.Sprintf("Not expected type :%v", reflect.TypeOf(res.Data)))
		}
	}
	return nil
}
