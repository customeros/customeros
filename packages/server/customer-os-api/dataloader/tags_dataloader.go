package dataloader

import (
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"golang.org/x/net/context"
)

func (b *batcher) getTagsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var orgIds []string
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	for ix, key := range keys {
		orgIds = append(orgIds, key.String())
		keyOrder[key.String()] = ix
	}

	tagEntitiesPtr, err := b.tagService.GetTagsForOrganizations(ctx, orgIds)
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
