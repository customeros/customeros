package dataloader

import (
	"errors"
	"fmt"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"golang.org/x/net/context"
	"reflect"
	"time"
)

const organizationContextTimeout = 10 * time.Second

func (i *Loaders) GetOrganizationsForEmail(ctx context.Context, emailId string) (*entity.OrganizationEntities, error) {
	thunk := i.OrganizationsForEmail.Load(ctx, dataloader.StringKey(emailId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.OrganizationEntities)
	return &resultObj, nil
}

func (i *Loaders) GetOrganizationsForPhoneNumber(ctx context.Context, phoneNumberId string) (*entity.OrganizationEntities, error) {
	thunk := i.OrganizationsForPhoneNumber.Load(ctx, dataloader.StringKey(phoneNumberId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.OrganizationEntities)
	return &resultObj, nil
}

func (b *organizationBatcher) getOrganizationsForEmails(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, organizationContextTimeout)
	defer cancel()

	organizationEntitiesPtr, err := b.organizationService.GetOrganizationsForEmails(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get organizations for emails")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	organizationEntitiesByEmailId := make(map[string]entity.OrganizationEntities)
	for _, val := range *organizationEntitiesPtr {
		if list, ok := organizationEntitiesByEmailId[val.DataloaderKey]; ok {
			organizationEntitiesByEmailId[val.DataloaderKey] = append(list, val)
		} else {
			organizationEntitiesByEmailId[val.DataloaderKey] = entity.OrganizationEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for emailId, record := range organizationEntitiesByEmailId {
		if ix, ok := keyOrder[emailId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, emailId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.OrganizationEntities{}, Error: nil}
	}

	if err = assertOrganizationEntitiesType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func (b *organizationBatcher) getOrganizationsForPhoneNumbers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, organizationContextTimeout)
	defer cancel()

	organizationEntitiesPtr, err := b.organizationService.GetOrganizationsForPhoneNumbers(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get organizations for phone numbers")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	organizationEntitiesByPhoneNumberId := make(map[string]entity.OrganizationEntities)
	for _, val := range *organizationEntitiesPtr {
		if list, ok := organizationEntitiesByPhoneNumberId[val.DataloaderKey]; ok {
			organizationEntitiesByPhoneNumberId[val.DataloaderKey] = append(list, val)
		} else {
			organizationEntitiesByPhoneNumberId[val.DataloaderKey] = entity.OrganizationEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for phoneNumberId, record := range organizationEntitiesByPhoneNumberId {
		if ix, ok := keyOrder[phoneNumberId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, phoneNumberId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.OrganizationEntities{}, Error: nil}
	}

	if err = assertOrganizationEntitiesType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func assertOrganizationEntitiesType(results []*dataloader.Result) error {
	for _, res := range results {
		if _, ok := res.Data.(entity.OrganizationEntities); !ok {
			return errors.New(fmt.Sprintf("Not expected type :%v", reflect.TypeOf(res.Data)))
		}
	}
	return nil
}
