package dataloader

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
	"reflect"
)

func (i *Loaders) GetTasksForOpportunity(ctx context.Context, opportunityId string) (*neo4jentity.TaskEntities, error) {
	thunk := i.TasksForOpportunity.Load(ctx, dataloader.StringKey(opportunityId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.TaskEntities)
	return &resultObj, nil
}

func (b *taskBatcher) getTasksForOpportunities(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TaskDataLoader.getTasksForOpportunities")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	taskEntitiesPtr, err := b.taskService.GetTasksForOpportunities(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get tasks for opportunities")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	taskEntitiesByOpportunityId := make(map[string]neo4jentity.TaskEntities)
	for _, val := range *taskEntitiesPtr {
		if list, ok := taskEntitiesByOpportunityId[val.DataloaderKey]; ok {
			taskEntitiesByOpportunityId[val.DataloaderKey] = append(list, val)
		} else {
			taskEntitiesByOpportunityId[val.DataloaderKey] = neo4jentity.TaskEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for opportunityId, record := range taskEntitiesByOpportunityId {
		if ix, ok := keyOrder[opportunityId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, opportunityId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.TaskEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.TaskEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
