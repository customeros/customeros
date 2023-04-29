package dataloader

import (
	"errors"
	"fmt"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"golang.org/x/net/context"
	"reflect"
	"time"
)

const describedByContextTimeout = 10 * time.Second

func (i *Loaders) GetDescribedByForMeeting(ctx context.Context, meetingId string) (*entity.AnalysisEntities, error) {
	thunk := i.DescribedByForMeeting.Load(ctx, dataloader.StringKey(meetingId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.AnalysisEntities)
	return &resultObj, nil
}

func (b *analysisBatcher) getDescribedByForMeeting(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, describesContextTimeout)
	defer cancel()

	analysisEntitiesPtr, err := b.analysisService.GetDescribedByForXX(ctx, ids, repository.DESCRIBES_TYPE_MEETING)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get analysis for meeting")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	analysisGrouped := make(map[string]entity.AnalysisEntities)
	for _, val := range *analysisEntitiesPtr {
		if list, ok := analysisGrouped[val.GetDataloaderKey()]; ok {
			analysisGrouped[val.GetDataloaderKey()] = append(list, val)
		} else {
			analysisGrouped[val.GetDataloaderKey()] = entity.AnalysisEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range analysisGrouped {
		ix, ok := keyOrder[contactId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.AnalysisEntities{}, Error: nil}
	}

	if err = assertAnalysisEntitiesType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func (i *Loaders) GetDescribedByForInteractionSession(ctx context.Context, meetingId string) (*entity.AnalysisEntities, error) {
	thunk := i.DescribedByForInteractionSession.Load(ctx, dataloader.StringKey(meetingId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.AnalysisEntities)
	return &resultObj, nil
}

func (b *analysisBatcher) getDescribedByForInteractionSession(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := context.WithTimeout(ctx, describesContextTimeout)
	defer cancel()

	analysisEntitiesPtr, err := b.analysisService.GetDescribedByForXX(ctx, ids, repository.DESCRIBES_TYPE_INTERACTION_SESSION)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get analysis for interaction session")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	analysisGrouped := make(map[string]entity.AnalysisEntities)
	for _, val := range *analysisEntitiesPtr {
		if list, ok := analysisGrouped[val.GetDataloaderKey()]; ok {
			analysisGrouped[val.GetDataloaderKey()] = append(list, val)
		} else {
			analysisGrouped[val.GetDataloaderKey()] = entity.AnalysisEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contactId, record := range analysisGrouped {
		ix, ok := keyOrder[contactId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, contactId)
		}
	}

	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: entity.AnalysisEntities{}, Error: nil}
	}

	if err = assertAnalysisEntitiesType(results); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}

func assertAnalysisEntitiesType(results []*dataloader.Result) error {
	for _, res := range results {
		if _, ok := res.Data.(entity.AnalysisEntities); !ok {
			return errors.New(fmt.Sprintf("Not expected type :%v", reflect.TypeOf(res.Data)))
		}
	}
	return nil
}
