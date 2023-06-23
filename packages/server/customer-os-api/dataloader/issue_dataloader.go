package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"reflect"
)

func (i *Loaders) GetIssueForInteractionEvent(ctx context.Context, interactionEventId string) (*entity.IssueEntity, error) {
	thunk := i.IssueForInteractionEvent.Load(ctx, dataloader.StringKey(interactionEventId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	resultObj := result.(*entity.IssueEntity)
	return resultObj, nil
}

func (b *issueBatcher) getIssuesForInteractionEvents(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	issueEntities, err := b.issueService.GetIssuesForInteractionEvents(ctx, ids)
	if err != nil {
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get issues for interaction events")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	issueEntityByInteractionEventId := make(map[string]entity.IssueEntity)
	for _, val := range *issueEntities {
		issueEntityByInteractionEventId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for interactionEventId, _ := range issueEntityByInteractionEventId {
		if ix, ok := keyOrder[interactionEventId]; ok {
			val := issueEntityByInteractionEventId[interactionEventId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, interactionEventId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(entity.IssueEntity{}), true); err != nil {
		return []*dataloader.Result{{nil, err}}
	}

	return results
}
