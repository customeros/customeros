package dataloader

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
	"reflect"
)

func (i *Loaders) GetJobRolesForContact(ctx context.Context, contactId string) (*neo4jentity.JobRoleEntities, error) {
	thunk := i.JobRolesForContact.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.JobRoleEntities)
	return &resultObj, nil
}

func (i *Loaders) GetJobRolesForOrganization(ctx context.Context, organizationId string) (*neo4jentity.JobRoleEntities, error) {
	thunk := i.JobRolesForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.JobRoleEntities)
	return &resultObj, nil
}

func (i *Loaders) GetJobRolesForUser(ctx context.Context, userId string) (*neo4jentity.JobRoleEntities, error) {
	thunk := i.JobRolesForUser.Load(ctx, dataloader.StringKey(userId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.JobRoleEntities)
	return &resultObj, nil
}

func (b *jobRoleBatcher) getJobRolesForContacts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "JobRoleDataLoader.getJobRolesForContacts")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	jobRoleEntitiesPtr, err := b.jobRoleService.GetAllForContacts(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get job roles for contacts")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	jobRoleEntitiesGroupedByContactId := make(map[string]neo4jentity.JobRoleEntities)
	for _, val := range *jobRoleEntitiesPtr {
		if list, ok := jobRoleEntitiesGroupedByContactId[val.DataloaderKey]; ok {
			jobRoleEntitiesGroupedByContactId[val.DataloaderKey] = append(list, val)
		} else {
			jobRoleEntitiesGroupedByContactId[val.DataloaderKey] = neo4jentity.JobRoleEntities{val}
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
		results[ix] = &dataloader.Result{Data: neo4jentity.JobRoleEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.JobRoleEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *jobRoleBatcher) getJobRolesForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "JobRoleDataLoader.getJobRolesForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	jobRoleEntitiesPtr, err := b.jobRoleService.GetAllForOrganizations(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get job roles for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	jobRoleEntitiesGroupedByOrganizationId := make(map[string]neo4jentity.JobRoleEntities)
	for _, val := range *jobRoleEntitiesPtr {
		if list, ok := jobRoleEntitiesGroupedByOrganizationId[val.DataloaderKey]; ok {
			jobRoleEntitiesGroupedByOrganizationId[val.DataloaderKey] = append(list, val)
		} else {
			jobRoleEntitiesGroupedByOrganizationId[val.DataloaderKey] = neo4jentity.JobRoleEntities{val}
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
		results[ix] = &dataloader.Result{Data: neo4jentity.JobRoleEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.JobRoleEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *jobRoleBatcher) getJobRolesForUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "JobRoleDataLoader.getJobRolesForUsers")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	jobRoleEntitiesPtr, err := b.jobRoleService.GetAllForUsers(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get job roles for contacts")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	jobRoleEntitiesGroupedByUserId := make(map[string]neo4jentity.JobRoleEntities)
	for _, val := range *jobRoleEntitiesPtr {
		if list, ok := jobRoleEntitiesGroupedByUserId[val.DataloaderKey]; ok {
			jobRoleEntitiesGroupedByUserId[val.DataloaderKey] = append(list, val)
		} else {
			jobRoleEntitiesGroupedByUserId[val.DataloaderKey] = neo4jentity.JobRoleEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for userId, record := range jobRoleEntitiesGroupedByUserId {
		ix, ok := keyOrder[userId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, userId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.JobRoleEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.JobRoleEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
