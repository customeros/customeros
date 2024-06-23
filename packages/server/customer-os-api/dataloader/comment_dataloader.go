package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentDataLoader.getCommentsForIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	commentEntitiesPtr, err := b.commentService.GetCommentsForIssues(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
