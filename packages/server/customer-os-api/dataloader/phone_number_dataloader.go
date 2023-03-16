package dataloader

import (
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"golang.org/x/net/context"
	"time"
)

const phoneNumberContextTimeout = 10 * time.Second

func (i *Loaders) GetPhoneNumbersForOrganization(ctx context.Context, organizationId string) (*entity.PhoneNumberEntities, error) {
	thunk := i.PhoneNumbersForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.PhoneNumberEntities)
	return &resultObj, nil
}

func (b *phoneNumberBatcher) getPhoneNumbersForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, phoneNumberContextTimeout)
	defer cancel()

	phoneNumberEntitiesPtr, err := b.phoneNumberService.GetAllForEntityTypeByIds(ctx, entity.ORGANIZATION, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get phone numbers for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	phoneNumberEntitiesGrouped := make(map[string]entity.PhoneNumberEntities)
	for _, val := range *phoneNumberEntitiesPtr {
		if list, ok := phoneNumberEntitiesGrouped[val.DataloaderKey]; ok {
			phoneNumberEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			phoneNumberEntitiesGrouped[val.DataloaderKey] = entity.PhoneNumberEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range phoneNumberEntitiesGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.PhoneNumberEntities{}, Error: nil}
	}

	return results
}
