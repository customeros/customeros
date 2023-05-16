package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"reflect"
	"time"
)

const userContextTimeout = 10 * time.Second

func (i *Loaders) GetUsersForEmail(ctx context.Context, emailId string) (*entity.UserEntities, error) {
	thunk := i.UsersForEmail.Load(ctx, dataloader.StringKey(emailId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.UserEntities)
	return &resultObj, nil
}

func (i *Loaders) GetUsersForPhoneNumber(ctx context.Context, phoneNumberId string) (*entity.UserEntities, error) {
	thunk := i.UsersForPhoneNumber.Load(ctx, dataloader.StringKey(phoneNumberId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.UserEntities)
	return &resultObj, nil
}

func (i *Loaders) GetUsersForPerson(ctx context.Context, personId string) (*entity.UserEntities, error) {
	thunk := i.UsersForPerson.Load(ctx, dataloader.StringKey(personId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.UserEntities)
	return &resultObj, nil
}

func (b *userBatcher) getUsersForEmails(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, userContextTimeout)
	defer cancel()

	userEntitiesPtr, err := b.userService.GetUsersForEmails(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get users for emails")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntitiesByEmailId := make(map[string]entity.UserEntities)
	for _, val := range *userEntitiesPtr {
		if list, ok := userEntitiesByEmailId[val.DataloaderKey]; ok {
			userEntitiesByEmailId[val.DataloaderKey] = append(list, val)
		} else {
			userEntitiesByEmailId[val.DataloaderKey] = entity.UserEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for emailId, record := range userEntitiesByEmailId {
		if ix, ok := keyOrder[emailId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, emailId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.UserEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.UserEntities{})); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func (b *userBatcher) getUsersForPhoneNumbers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, userContextTimeout)
	defer cancel()

	userEntitiesPtr, err := b.userService.GetUsersForPhoneNumbers(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get users for phone numbers")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntitiesByPhoneNumberId := make(map[string]entity.UserEntities)
	for _, val := range *userEntitiesPtr {
		if list, ok := userEntitiesByPhoneNumberId[val.DataloaderKey]; ok {
			userEntitiesByPhoneNumberId[val.DataloaderKey] = append(list, val)
		} else {
			userEntitiesByPhoneNumberId[val.DataloaderKey] = entity.UserEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for phoneNumberId, record := range userEntitiesByPhoneNumberId {
		if ix, ok := keyOrder[phoneNumberId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, phoneNumberId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.UserEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.UserEntities{})); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func (b *userBatcher) getUsersForPersons(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, userContextTimeout)
	defer cancel()

	userEntitiesPtr, err := b.userService.GetUsersForPersons(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get users for persons")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntitiesByPersonId := make(map[string]entity.UserEntities)
	for _, val := range *userEntitiesPtr {
		if list, ok := userEntitiesByPersonId[val.DataloaderKey]; ok {
			userEntitiesByPersonId[val.DataloaderKey] = append(list, val)
		} else {
			userEntitiesByPersonId[val.DataloaderKey] = entity.UserEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for phoneNumberId, record := range userEntitiesByPersonId {
		if ix, ok := keyOrder[phoneNumberId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, phoneNumberId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.UserEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.UserEntities{})); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}
