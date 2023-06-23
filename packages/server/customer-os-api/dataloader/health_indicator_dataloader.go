package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"reflect"
)

func (i *Loaders) GetHealthIndicatorForOrganization(ctx context.Context, organizationId string) (*entity.HealthIndicatorEntity, error) {
	thunk := i.HealthIndicatorForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	resultObj := result.(*entity.HealthIndicatorEntity)
	return resultObj, nil
}

func (b *healthIndicatorBatcher) getHealthIndicatorsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	healthIndicatorEntities, err := b.healthIndicatorService.GetHealthIndicatorsForOrganizations(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get health indicators for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	healthIndicatorEntityByOrganizationId := make(map[string]entity.HealthIndicatorEntity)
	for _, val := range *healthIndicatorEntities {
		healthIndicatorEntityByOrganizationId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, _ := range healthIndicatorEntityByOrganizationId {
		if ix, ok := keyOrder[organizationId]; ok {
			val := healthIndicatorEntityByOrganizationId[organizationId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(entity.HealthIndicatorEntity{}), true); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}
