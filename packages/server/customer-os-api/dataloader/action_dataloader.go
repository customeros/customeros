package dataloader

import (
	"context"
	"errors"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
	"reflect"
)

func (i *Loaders) GetActionsForInteractionEvent(ctx context.Context, interactionEventId string) (*neo4jentity.ActionEntities, error) {
	thunk := i.ActionsForInteractionEvent.Load(ctx, dataloader.StringKey(interactionEventId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.ActionEntities)
	return &resultObj, nil
}

func (b *actionBatcher) getActionsForInteractionEvents(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ActionDataLoader.getActionsForInteractionEvents")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	actionsForNodes, err := b.actionService.GetActionsForNodes(ctx, neo4jenum.INTERACTION_EVENT, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get attachments for interaction events")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	entitiesByInteractionEventId := make(map[string]neo4jentity.ActionEntities)
	for _, val := range *actionsForNodes {
		if list, ok := entitiesByInteractionEventId[val.DataloaderKey]; ok {
			entitiesByInteractionEventId[val.DataloaderKey] = append(list, val)
		} else {
			entitiesByInteractionEventId[val.DataloaderKey] = neo4jentity.ActionEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for id, record := range entitiesByInteractionEventId {
		if ix, ok := keyOrder[id]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, id)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.ActionEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.ActionEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
