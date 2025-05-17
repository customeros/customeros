package dataloader

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
	"reflect"
)

func (i *Loaders) GetTagsForOrganization(ctx context.Context, organizationId string) (*neo4jentity.TagEntities, error) {
	thunk := i.TagsForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.TagEntities)
	return &resultObj, nil
}

func (i *Loaders) GetTagsForContact(ctx context.Context, contactId string) (*neo4jentity.TagEntities, error) {
	thunk := i.TagsForContact.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.TagEntities)
	return &resultObj, nil
}

func (i *Loaders) GetTagsForIssue(ctx context.Context, issueId string) (*neo4jentity.TagEntities, error) {
	thunk := i.TagsForIssue.Load(ctx, dataloader.StringKey(issueId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.TagEntities)
	return &resultObj, nil
}

func (i *Loaders) GetTagsForLogEntry(ctx context.Context, logEntryId string) (*neo4jentity.TagEntities, error) {
	thunk := i.TagsForLogEntry.Load(ctx, dataloader.StringKey(logEntryId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.TagEntities)
	return &resultObj, nil
}

func (b *tagBatcher) getTagsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TagDataLoader.getTagsForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	tagEntitiesPtr, err := b.tagService.GetTagsForOrganizations(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get tags for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	tagEntitiesByOrganizationId := make(map[string]neo4jentity.TagEntities)
	for _, val := range *tagEntitiesPtr {
		if list, ok := tagEntitiesByOrganizationId[val.DataloaderKey]; ok {
			tagEntitiesByOrganizationId[val.DataloaderKey] = append(list, val)
		} else {
			tagEntitiesByOrganizationId[val.DataloaderKey] = neo4jentity.TagEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range tagEntitiesByOrganizationId {
		if ix, ok := keyOrder[organizationId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.TagEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.TagEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *tagBatcher) getTagsForContacts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TagDataLoader.getTagsForContacts")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	tagEntitiesPtr, err := b.tagService.GetTagsForContacts(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get tags for contacts")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	tagEntitiesByContactId := make(map[string]neo4jentity.TagEntities)
	for _, val := range *tagEntitiesPtr {
		if list, ok := tagEntitiesByContactId[val.DataloaderKey]; ok {
			tagEntitiesByContactId[val.DataloaderKey] = append(list, val)
		} else {
			tagEntitiesByContactId[val.DataloaderKey] = neo4jentity.TagEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range tagEntitiesByContactId {
		if ix, ok := keyOrder[contactId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.TagEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.TagEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *tagBatcher) getTagsForIssues(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TagDataLoader.getTagsForIssues")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	tagEntitiesPtr, err := b.tagService.GetTagsForIssues(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get tags for issues")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	tagEntitiesByIssueId := make(map[string]neo4jentity.TagEntities)
	for _, val := range *tagEntitiesPtr {
		if list, ok := tagEntitiesByIssueId[val.DataloaderKey]; ok {
			tagEntitiesByIssueId[val.DataloaderKey] = append(list, val)
		} else {
			tagEntitiesByIssueId[val.DataloaderKey] = neo4jentity.TagEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for issueId, record := range tagEntitiesByIssueId {
		if ix, ok := keyOrder[issueId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, issueId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.TagEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.TagEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *tagBatcher) getTagsForLogEntries(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TagDataLoader.getTagsForLogEntries")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	tagEntitiesPtr, err := b.tagService.GetTagsForLogEntries(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get tags for log entries")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	tagEntitiesByLogEntryId := make(map[string]neo4jentity.TagEntities)
	for _, val := range *tagEntitiesPtr {
		if list, ok := tagEntitiesByLogEntryId[val.DataloaderKey]; ok {
			tagEntitiesByLogEntryId[val.DataloaderKey] = append(list, val)
		} else {
			tagEntitiesByLogEntryId[val.DataloaderKey] = neo4jentity.TagEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for logEntryId, record := range tagEntitiesByLogEntryId {
		if ix, ok := keyOrder[logEntryId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, logEntryId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.TagEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.TagEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
