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

func (i *Loaders) GetOpportunitiesForContract(ctx context.Context, contractId string) (*entity.OpportunityEntities, error) {
	thunk := i.OpportunitiesForContract.Load(ctx, dataloader.StringKey(contractId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.OpportunityEntities)
	return &resultObj, nil
}

func (b *opportunityBatcher) getOpportunitiesForContracts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityDataLoader.getOpportunitiesForContracts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	opportunityEntitiesPtr, err := b.opportunityService.GetOpportunitiesForContracts(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get opportunities for contracts")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	opportunityEntitiesByContractId := make(map[string]entity.OpportunityEntities)
	for _, val := range *opportunityEntitiesPtr {
		if list, ok := opportunityEntitiesByContractId[val.DataloaderKey]; ok {
			opportunityEntitiesByContractId[val.DataloaderKey] = append(list, val)
		} else {
			opportunityEntitiesByContractId[val.DataloaderKey] = entity.OpportunityEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contractId, record := range opportunityEntitiesByContractId {
		if ix, ok := keyOrder[contractId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contractId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.OpportunityEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.OpportunityEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
