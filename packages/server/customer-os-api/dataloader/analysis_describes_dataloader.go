package dataloader

import (
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"golang.org/x/net/context"
	"time"
)

const describesContextTimeout = 10 * time.Second

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
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, describesContextTimeout)
	defer cancel()

	participantEntitiesPtr, err := b.analysisService.GetDescribesForAnalysis(ctx, ids)
	if err != nil {
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

	return results
}
