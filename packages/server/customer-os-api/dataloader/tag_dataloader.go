package dataloader

import (
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"golang.org/x/net/context"
	"time"
)

const tagContextTimeout = 10 * time.Second

func (i *Loaders) GetTagsForOrganization(ctx context.Context, organizationId string) (*entity.TagEntities, error) {
	thunk := i.TagsForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.TagEntities)
	return &resultObj, nil
}

func (i *Loaders) GetTagsForContact(ctx context.Context, contactId string) (*entity.TagEntities, error) {
	thunk := i.TagsForContact.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.TagEntities)
	return &resultObj, nil
}

func (i *Loaders) GetTagsForIssue(ctx context.Context, issueId string) (*entity.TagEntities, error) {
	thunk := i.TagsForIssue.Load(ctx, dataloader.StringKey(issueId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.TagEntities)
	return &resultObj, nil
}

func (b *tagBatcher) getTagsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, tagContextTimeout)
	defer cancel()

	tagEntitiesPtr, err := b.tagService.GetTagsForOrganizations(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get tags for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	tagEntitiesByOrganizationId := make(map[string]entity.TagEntities)
	for _, val := range *tagEntitiesPtr {
		if list, ok := tagEntitiesByOrganizationId[val.DataloaderKey]; ok {
			tagEntitiesByOrganizationId[val.DataloaderKey] = append(list, val)
		} else {
			tagEntitiesByOrganizationId[val.DataloaderKey] = entity.TagEntities{val}
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
		results[ix] = &dataloader.Result{Data: entity.TagEntities{}, Error: nil}
	}

	return results
}

func (b *tagBatcher) getTagsForContacts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, tagContextTimeout)
	defer cancel()

	tagEntitiesPtr, err := b.tagService.GetTagsForContacts(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get tags for contacts")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	tagEntitiesByContactId := make(map[string]entity.TagEntities)
	for _, val := range *tagEntitiesPtr {
		if list, ok := tagEntitiesByContactId[val.DataloaderKey]; ok {
			tagEntitiesByContactId[val.DataloaderKey] = append(list, val)
		} else {
			tagEntitiesByContactId[val.DataloaderKey] = entity.TagEntities{val}
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
		results[ix] = &dataloader.Result{Data: entity.TagEntities{}, Error: nil}
	}

	return results
}

func (b *tagBatcher) getTagsForIssues(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, tagContextTimeout)
	defer cancel()

	tagEntitiesPtr, err := b.tagService.GetTagsForIssues(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get tags for issues")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	tagEntitiesByIssueId := make(map[string]entity.TagEntities)
	for _, val := range *tagEntitiesPtr {
		if list, ok := tagEntitiesByIssueId[val.DataloaderKey]; ok {
			tagEntitiesByIssueId[val.DataloaderKey] = append(list, val)
		} else {
			tagEntitiesByIssueId[val.DataloaderKey] = entity.TagEntities{val}
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
		results[ix] = &dataloader.Result{Data: entity.TagEntities{}, Error: nil}
	}

	return results
}
