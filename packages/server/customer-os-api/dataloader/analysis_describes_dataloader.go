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

func (i *Loaders) GetDescribesForAnalysis(ctx context.Context, analysisId string) (*entity.AnalysisDescribes, error) {
	thunk := i.DescribesForAnalysis.Load(ctx, dataloader.StringKey(analysisId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.AnalysisDescribes)
	return &resultObj, nil
}

func (b *analysisBatcher) getDescribesForAnalysis(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalysisDataLoader.getDescribesForAnalysis", opentracing.ChildOf(tracing.ExtractSpanCtx(ctx)))
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	participantEntitiesPtr, err := b.analysisService.GetDescribesForAnalysis(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get interaction event participants")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	analysisDescribeGrouped := make(map[string]entity.AnalysisDescribes)
	for _, val := range *participantEntitiesPtr {
		if list, ok := analysisDescribeGrouped[val.GetDataloaderKey()]; ok {
			analysisDescribeGrouped[val.GetDataloaderKey()] = append(list, val)
		} else {
			analysisDescribeGrouped[val.GetDataloaderKey()] = entity.AnalysisDescribes{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range analysisDescribeGrouped {
		ix, ok := keyOrder[contactId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.AnalysisDescribes{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.AnalysisDescribes{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("output - results_length", len(results)))

	return results
}
