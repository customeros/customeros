package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
)

func (i *Loaders) GetRelationshipsForOrganization(ctx context.Context, organizationID string) (entity.OrganizationRelationships, error) {
	thunk := i.RelationshipsForOrganization.Load(ctx, dataloader.StringKey(organizationID))
	result, err := thunk()
	if err != nil {
		return entity.OrganizationRelationships{}, err
	}
	return result.(entity.OrganizationRelationships), nil
}

func (i *Loaders) GetRelationshipStagesForOrganization(ctx context.Context, organizationID string) (*entity.OrganizationRelationshipsWithStages, error) {
	thunk := i.RelationshipStagesForOrganization.Load(ctx, dataloader.StringKey(organizationID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.OrganizationRelationshipsWithStages)
	return &resultObj, nil
}

func (b *relationshipBatcher) getRelationshipsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRelationshipDataLoader.getRelationshipsForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	entitiesPtr, err := b.organizationRelationshipService.GetRelationshipsForOrganizations(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{entity.OrganizationRelationships{}, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *relationshipBatcher) getRelationshipStagesForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRelationshipDataLoader.getRelationshipStagesForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	entitiesPtr, err := b.organizationRelationshipService.GetRelationshipsWithStagesForOrganizations(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get relationships with stages for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	organizationRelationshipsWithStagesGrouped := make(map[string]entity.OrganizationRelationshipsWithStages)
	for _, val := range *entitiesPtr {
		if list, ok := organizationRelationshipsWithStagesGrouped[val.DataloaderKey]; ok {
			organizationRelationshipsWithStagesGrouped[val.DataloaderKey] = append(list, val)
		} else {
			organizationRelationshipsWithStagesGrouped[val.DataloaderKey] = entity.OrganizationRelationshipsWithStages{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range organizationRelationshipsWithStagesGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.OrganizationRelationshipsWithStages{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.OrganizationRelationshipsWithStages{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
