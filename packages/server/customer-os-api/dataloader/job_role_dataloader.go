package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"reflect"
)

func (i *Loaders) GetJobRolesForContact(ctx context.Context, contactId string) (*entity.JobRoleEntities, error) {
	thunk := i.JobRolesForContact.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.JobRoleEntities)
	return &resultObj, nil
}

func (i *Loaders) GetJobRolesForOrganization(ctx context.Context, organizationId string) (*entity.JobRoleEntities, error) {
	thunk := i.JobRolesForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.JobRoleEntities)
	return &resultObj, nil
}

func (b *jobRoleBatcher) getJobRolesForContacts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	jobRoleEntitiesPtr, err := b.jobRoleService.GetAllForContacts(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get job roles for contacts")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	jobRoleEntitiesGroupedByContactId := make(map[string]entity.JobRoleEntities)
	for _, val := range *jobRoleEntitiesPtr {
		if list, ok := jobRoleEntitiesGroupedByContactId[val.DataloaderKey]; ok {
			jobRoleEntitiesGroupedByContactId[val.DataloaderKey] = append(list, val)
		} else {
			jobRoleEntitiesGroupedByContactId[val.DataloaderKey] = entity.JobRoleEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range jobRoleEntitiesGroupedByContactId {
		ix, ok := keyOrder[contactId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.JobRoleEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.JobRoleEntities{})); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func (b *jobRoleBatcher) getJobRolesForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	jobRoleEntitiesPtr, err := b.jobRoleService.GetAllForOrganizations(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get job roles for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	jobRoleEntitiesGroupedByOrganizationId := make(map[string]entity.JobRoleEntities)
	for _, val := range *jobRoleEntitiesPtr {
		if list, ok := jobRoleEntitiesGroupedByOrganizationId[val.DataloaderKey]; ok {
			jobRoleEntitiesGroupedByOrganizationId[val.DataloaderKey] = append(list, val)
		} else {
			jobRoleEntitiesGroupedByOrganizationId[val.DataloaderKey] = entity.JobRoleEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range jobRoleEntitiesGroupedByOrganizationId {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.JobRoleEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.JobRoleEntities{})); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}
