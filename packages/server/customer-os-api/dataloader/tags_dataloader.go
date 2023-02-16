package dataloader

import (
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"golang.org/x/net/context"
)

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

func (b *batcher) getTagsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var ids []string
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	for ix, key := range keys {
		ids = append(ids, key.String())
		keyOrder[key.String()] = ix
	}

	tagEntitiesPtr, err := b.tagService.GetTagsForOrganizations(ctx, ids)
	if err != nil {
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
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.TagEntities{}, Error: nil}
	}

	return results
}

func (b *batcher) getTagsForContacts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var ids []string
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	for ix, key := range keys {
		ids = append(ids, key.String())
		keyOrder[key.String()] = ix
	}

	tagEntitiesPtr, err := b.tagService.GetTagsForContacts(ctx, ids)
	if err != nil {
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
		ix, ok := keyOrder[contactId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.TagEntities{}, Error: nil}
	}

	return results
}
