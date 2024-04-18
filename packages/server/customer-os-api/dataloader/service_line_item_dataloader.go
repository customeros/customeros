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

func (i *Loaders) GetServiceLineItemsForContract(ctx context.Context, contractId string) (*neo4jentity.ServiceLineItemEntities, error) {
	thunk := i.ServiceLineItemsForContract.Load(ctx, dataloader.StringKey(contractId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.ServiceLineItemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetServiceLineItemForInvoiceLine(ctx context.Context, invoiceLineId string) (*neo4jentity.ServiceLineItemEntity, error) {
	thunk := i.ServiceLineItemForInvoiceLine.Load(ctx, dataloader.StringKey(invoiceLineId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.ServiceLineItemEntity), nil
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

	serviceLineItemEntitiesByContractId := make(map[string]neo4jentity.ServiceLineItemEntities)
	for _, val := range *serviceLineItemEntitiesPtr {
		if list, ok := serviceLineItemEntitiesByContractId[val.DataloaderKey]; ok {
			serviceLineItemEntitiesByContractId[val.DataloaderKey] = append(list, val)
		} else {
			serviceLineItemEntitiesByContractId[val.DataloaderKey] = neo4jentity.ServiceLineItemEntities{val}
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
		results[ix] = &dataloader.Result{Data: neo4jentity.ServiceLineItemEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.ServiceLineItemEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *serviceLineItemBatcher) getServiceLineItemForInvoiceLine(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemDataLoader.getServiceLineItemsForInvoiceLines")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	sliEntities, err := b.serviceLineItemService.GetServiceLineItemsForInvoiceLines(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get service line items for invoice lines")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	sliByInvoiceLineId := make(map[string]neo4jentity.ServiceLineItemEntity)
	for _, val := range *sliEntities {
		sliByInvoiceLineId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for jobRoleId := range sliByInvoiceLineId {
		if ix, ok := keyOrder[jobRoleId]; ok {
			val := sliByInvoiceLineId[jobRoleId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, jobRoleId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.ServiceLineItemEntity{}), true); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Object("output - results_length", len(results)))

	return results
}
