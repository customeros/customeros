package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"reflect"
)

func (i *Loaders) GetInvoiceLinesForInvoice(ctx context.Context, invoiceId string) (*neo4jentity.InvoiceLineEntities, error) {
	thunk := i.InvoiceLinesForInvoice.Load(ctx, dataloader.StringKey(invoiceId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.InvoiceLineEntities)
	return &resultObj, nil
}

func (b *invoiceBatcher) getInvoiceLinesForInvoice(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceDataLoader.getInvoiceLinesForInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	invoiceLinesEntitiesPtr, err := b.invoiceService.GetInvoiceLinesForInvoices(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get invoice lines for invoices")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	invoiceLinesByInvoiceId := make(map[string]neo4jentity.InvoiceLineEntities)
	for _, val := range *invoiceLinesEntitiesPtr {
		if list, ok := invoiceLinesByInvoiceId[val.DataloaderKey]; ok {
			invoiceLinesByInvoiceId[val.DataloaderKey] = append(list, val)
		} else {
			invoiceLinesByInvoiceId[val.DataloaderKey] = neo4jentity.InvoiceLineEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for invoiceId, record := range invoiceLinesByInvoiceId {
		if ix, ok := keyOrder[invoiceId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, invoiceId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.InvoiceLineEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.InvoiceLineEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
