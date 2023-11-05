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

func (i *Loaders) GetDomainsForOrganization(ctx context.Context, organizationId string) (*entity.DomainEntities, error) {
	thunk := i.DomainsForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.DomainEntities)
	return &resultObj, nil
}

func (b *domainBatcher) getDomainsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainDataLoader.getDomainsForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	domainEntitiesPtr, err := b.domainService.GetDomainsForOrganizations(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get domains for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	domainEntitiesByOrganizationId := make(map[string]entity.DomainEntities)
	for _, val := range *domainEntitiesPtr {
		if list, ok := domainEntitiesByOrganizationId[val.DataloaderKey]; ok {
			domainEntitiesByOrganizationId[val.DataloaderKey] = append(list, val)
		} else {
			domainEntitiesByOrganizationId[val.DataloaderKey] = entity.DomainEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range domainEntitiesByOrganizationId {
		if ix, ok := keyOrder[organizationId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.DomainEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.DomainEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("output - results_length", len(results)))

	return results
}
