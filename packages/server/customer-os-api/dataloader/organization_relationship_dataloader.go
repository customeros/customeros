package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"reflect"
	"time"
)

const relationshipContextTimeout = 10 * time.Second

func (i *Loaders) GetRelationshipsForOrganization(ctx context.Context, organizationID string) (entity.OrganizationRelationships, error) {
	thunk := i.RelationshipsForOrganization.Load(ctx, dataloader.StringKey(organizationID))
	result, err := thunk()
	if err != nil {
		return entity.OrganizationRelationships{}, err
	}
	return result.(entity.OrganizationRelationships), nil
}

func (b *relationshipBatcher) getRelationshipsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, relationshipContextTimeout)
	defer cancel()

	entitiesPtr, err := b.organizationRelationshipService.GetRelationshipsForOrganizations(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: entity.OrganizationRelationships{}, Error: errors.New("deadline exceeded to get relationships for organizations")}}
		}
		return []*dataloader.Result{{Data: entity.OrganizationRelationships{}, Error: err}}
	}

	relationsGrouped := map[string]entity.OrganizationRelationships{}
	for _, val := range *entitiesPtr {
		if list, ok := relationsGrouped[val.DataloaderKey]; ok {
			relationsGrouped[val.DataloaderKey] = append(list, val.GetOrganizationRelationship())
		} else {
			relationsGrouped[val.DataloaderKey] = entity.OrganizationRelationships{val.GetOrganizationRelationship()}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range relationsGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.OrganizationRelationships{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.OrganizationRelationships{})); err != nil {
		return []*dataloader.Result{{entity.OrganizationRelationships{}, err}}
	}

	return results
}
