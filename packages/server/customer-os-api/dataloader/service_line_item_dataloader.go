package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"reflect"
)

func (i *Loaders) GetServiceLineItemsForContract(ctx context.Context, contractId string) (*entity.ServiceLineItemEntities, error) {
	thunk := i.ServiceLineItemsForContract.Load(ctx, dataloader.StringKey(contractId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.ServiceLineItemEntities)
	return &resultObj, nil
}

func (b *serviceLineItemBatcher) getServiceLineItemsForContracts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemDataLoader.getServiceLineItemsForContracts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	serviceLineItemEntitiesPtr, err := b.serviceLineItemService.GetServiceLineItemsForContracts(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get service line items for contracts")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	serviceLineItemEntitiesByContractId := make(map[string]entity.ServiceLineItemEntities)
	for _, val := range *serviceLineItemEntitiesPtr {
		if list, ok := serviceLineItemEntitiesByContractId[val.DataloaderKey]; ok {
			serviceLineItemEntitiesByContractId[val.DataloaderKey] = append(list, val)
		} else {
			serviceLineItemEntitiesByContractId[val.DataloaderKey] = entity.ServiceLineItemEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contractId, record := range serviceLineItemEntitiesByContractId {
		if ix, ok := keyOrder[contractId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contractId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.ServiceLineItemEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.ServiceLineItemEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
