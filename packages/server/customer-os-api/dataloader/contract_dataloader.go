package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"reflect"
)

func (i *Loaders) GetContractsForOrganization(ctx context.Context, organizationId string) (*neo4jentity.ContractEntities, error) {
	thunk := i.ContractsForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.ContractEntities)
	return &resultObj, nil
}

func (i *Loaders) GetContractForInvoice(ctx context.Context, invoiceId string) (*neo4jentity.ContractEntity, error) {
	thunk := i.ContractForInvoice.Load(ctx, dataloader.StringKey(invoiceId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.ContractEntity), nil
}

func (b *contractBatcher) getContractsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractDataLoader.getContractsForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	contractEntitiesPtr, err := b.contractService.GetContractsForOrganizations(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get contracts for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	contractEntitiesByOrganizationId := make(map[string]neo4jentity.ContractEntities)
	for _, val := range *contractEntitiesPtr {
		if list, ok := contractEntitiesByOrganizationId[val.DataloaderKey]; ok {
			contractEntitiesByOrganizationId[val.DataloaderKey] = append(list, val)
		} else {
			contractEntitiesByOrganizationId[val.DataloaderKey] = neo4jentity.ContractEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range contractEntitiesByOrganizationId {
		if ix, ok := keyOrder[organizationId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.ContractEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.ContractEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *contractBatcher) getContractsForInvoices(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractDataLoader.getContractsForInvoices")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	contractEntities, err := b.contractService.GetContractsForInvoices(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get contracts for invoices")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	contractEntityByInvoiceId := make(map[string]neo4jentity.ContractEntity)
	for _, val := range *contractEntities {
		contractEntityByInvoiceId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for invoiceId, _ := range contractEntityByInvoiceId {
		if ix, ok := keyOrder[invoiceId]; ok {
			val := contractEntityByInvoiceId[invoiceId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, invoiceId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.ContractEntity{}), true); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("result.length", len(results)))

	return results
}
