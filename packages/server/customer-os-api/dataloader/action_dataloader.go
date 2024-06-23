package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionDataLoader.getActionsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	actionsForNodes, err := b.actionService.GetActionsForNodes(ctx, neo4jenum.INTERACTION_EVENT, ids)
	if err != nil {
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Object("output - results_length", len(results)))

	return results
}
