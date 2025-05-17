package dataloader

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
	"reflect"
)

func (i *Loaders) GetInteractionSessionForInteractionEvent(ctx context.Context, interactionEventId string) (*neo4jentity.InteractionSessionEntity, error) {
	thunk := i.InteractionSessionForInteractionEvent.Load(ctx, dataloader.StringKey(interactionEventId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	resultObj := result.(*neo4jentity.InteractionSessionEntity)
	return resultObj, nil
}

func (b *interactionSessionBatcher) getInteractionSessionsForInteractionEvents(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "InteractionSessionDataLoader.getInteractionSessionsForInteractionEvents")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	interactionSessionEntities, err := b.interactionSessionService.GetInteractionSessionsForInteractionEvents(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get interaction sessions for interaction events")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	interactionSessionEntityByInteractionEventId := make(map[string]neo4jentity.InteractionSessionEntity)
	for _, val := range *interactionSessionEntities {
		interactionSessionEntityByInteractionEventId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for interactionEventId, _ := range interactionSessionEntityByInteractionEventId {
		if ix, ok := keyOrder[interactionEventId]; ok {
			val := interactionSessionEntityByInteractionEventId[interactionEventId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, interactionEventId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.InteractionSessionEntity{}), true); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
