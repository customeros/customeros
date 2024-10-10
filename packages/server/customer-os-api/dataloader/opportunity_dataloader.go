package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"reflect"
)

func (i *Loaders) GetOpportunitiesForContract(ctx context.Context, tenant, contractId string) (*neo4jentity.OpportunityEntities, error) {
	thunk := i.OpportunitiesForContract.Load(ctx, dataloader.StringKey(contractId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.OpportunityEntities)
	return &resultObj, nil
}

func (i *Loaders) GetOpportunitiesForOrganization(ctx context.Context, orgId string) (*neo4jentity.OpportunityEntities, error) {
	thunk := i.OpportunitiesForOrganization.Load(ctx, dataloader.StringKey(orgId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.OpportunityEntities)
	return &resultObj, nil
}

func (b *opportunityBatcher) getOpportunitiesForContracts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityDataLoader.getOpportunitiesForContracts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	tenant := common.GetTenantFromContext(ctx)

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	opportunityEntitiesPtr, err := b.opportunityService.GetOpportunitiesForContracts(ctx, tenant, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get opportunities for contracts")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	opportunityEntitiesByContractId := make(map[string]neo4jentity.OpportunityEntities)
	for _, val := range *opportunityEntitiesPtr {
		if list, ok := opportunityEntitiesByContractId[val.DataloaderKey]; ok {
			opportunityEntitiesByContractId[val.DataloaderKey] = append(list, val)
		} else {
			opportunityEntitiesByContractId[val.DataloaderKey] = neo4jentity.OpportunityEntities{val}
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
		results[ix] = &dataloader.Result{Data: neo4jentity.OpportunityEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.OpportunityEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *opportunityBatcher) getOpportunitiesForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityDataLoader.getOpportunitiesForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	tenant := common.GetTenantFromContext(ctx)

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	opportunityEntitiesPtr, err := b.opportunityService.GetOpportunitiesForOrganizations(ctx, tenant, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get opportunities for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	opportunityEntitiesByOrganizationId := make(map[string]neo4jentity.OpportunityEntities)
	for _, val := range *opportunityEntitiesPtr {
		if list, ok := opportunityEntitiesByOrganizationId[val.DataloaderKey]; ok {
			opportunityEntitiesByOrganizationId[val.DataloaderKey] = append(list, val)
		} else {
			opportunityEntitiesByOrganizationId[val.DataloaderKey] = neo4jentity.OpportunityEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for orgId, record := range opportunityEntitiesByOrganizationId {
		if ix, ok := keyOrder[orgId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, orgId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.OpportunityEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.OpportunityEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
