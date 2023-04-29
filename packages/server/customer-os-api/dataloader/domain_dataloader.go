package dataloader

import (
	"errors"
	"fmt"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"golang.org/x/net/context"
	"reflect"
	"time"
)

const domainContextTimeout = 10 * time.Second

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
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, domainContextTimeout)
	defer cancel()

	domainEntitiesPtr, err := b.domainService.GetDomainsForOrganizations(ctx, ids)
	if err != nil {
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

	if err = assertDomainEntitiesType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func assertDomainEntitiesType(results []*dataloader.Result) error {
	for _, res := range results {
		if _, ok := res.Data.(entity.DomainEntities); !ok {
			return errors.New(fmt.Sprintf("Not expected type :%v", reflect.TypeOf(res.Data)))
		}
	}
	return nil
}
