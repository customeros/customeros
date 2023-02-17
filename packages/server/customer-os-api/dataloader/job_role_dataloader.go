package dataloader

import (
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"golang.org/x/net/context"
	"time"
)

const jobRoleContextTimeout = 10 * time.Second

func (i *Loaders) GetJobRolesForContact(ctx context.Context, contactId string) (*entity.JobRoleEntities, error) {
	thunk := i.JobRolesForContact.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.JobRoleEntities)
	return &resultObj, nil
}

func (b *jobRoleBatcher) getJobRolesForContacts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, jobRoleContextTimeout)
	defer cancel()

	jobRoleEntitiesPtr, err := b.jobRoleService.GetAllForContacts(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get job roles for contacts")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	jobRoleEntitiesGrouped := make(map[string]entity.JobRoleEntities)
	for _, val := range *jobRoleEntitiesPtr {
		if list, ok := jobRoleEntitiesGrouped[val.DataloaderKey]; ok {
			jobRoleEntitiesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			jobRoleEntitiesGrouped[val.DataloaderKey] = entity.JobRoleEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range jobRoleEntitiesGrouped {
		ix, ok := keyOrder[contactId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.JobRoleEntities{}, Error: nil}
	}

	return results
}
