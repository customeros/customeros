package dataloader

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
	"reflect"
)

func (i *Loaders) GetDomainsForOrganization(ctx context.Context, organizationId string) (*neo4jentity.DomainEntities, error) {
	thunk := i.DomainsForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.DomainEntities)
	return &resultObj, nil
}

func (b *domainBatcher) getDomainsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "DomainDataLoader.getDomainsForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	domainEntitiesPtr, err := b.domainService.GetAllDomainsForOrganizations(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get domains for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	domainEntitiesByOrganizationId := make(map[string]neo4jentity.DomainEntities)
	for _, val := range *domainEntitiesPtr {
		if list, ok := domainEntitiesByOrganizationId[val.DataloaderKey]; ok {
			domainEntitiesByOrganizationId[val.DataloaderKey] = append(list, val)
		} else {
			domainEntitiesByOrganizationId[val.DataloaderKey] = neo4jentity.DomainEntities{val}
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
		results[ix] = &dataloader.Result{Data: neo4jentity.DomainEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.DomainEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
