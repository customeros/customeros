package dataloader

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"
	"reflect"
)

func (i *Loaders) GetCommentsForIssue(ctx context.Context, issueId string) (*neo4jentity.CommentEntities, error) {
	thunk := i.CommentsForIssue.Load(ctx, dataloader.StringKey(issueId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.CommentEntities)
	return &resultObj, nil
}

func (b *commentBatcher) getCommentsForIssues(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "CommentDataLoader.getCommentsForIssues")
	defer spans.Finish()
	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	commentEntitiesPtr, err := b.commentService.GetCommentsForIssues(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	commentEntitiesByMeetingId := make(map[string]neo4jentity.CommentEntities)
	for _, val := range *commentEntitiesPtr {
		if list, ok := commentEntitiesByMeetingId[val.DataloaderKey]; ok {
			commentEntitiesByMeetingId[val.DataloaderKey] = append(list, val)
		} else {
			commentEntitiesByMeetingId[val.DataloaderKey] = neo4jentity.CommentEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for meetingsId, record := range commentEntitiesByMeetingId {
		if ix, ok := keyOrder[meetingsId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, meetingsId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.CommentEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.CommentEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
