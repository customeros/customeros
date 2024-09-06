package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
)

func (i *Loaders) GetFlowSequenceContactsForFlowSequence(ctx context.Context, flowSequenceId string) (*neo4jentity.FlowSequenceContactEntities, error) {
	thunk := i.FlowSequenceContactsForFlowSequence.Load(ctx, dataloader.StringKey(flowSequenceId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.FlowSequenceContactEntities)
	return &resultObj, nil
}

func (i *Loaders) GetFlowSequenceSendersForFlowSequence(ctx context.Context, flowSequenceId string) (*neo4jentity.FlowSequenceSenderEntities, error) {
	thunk := i.FlowSequenceSendersForFlowSequence.Load(ctx, dataloader.StringKey(flowSequenceId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.FlowSequenceSenderEntities)
	return &resultObj, nil
}

func (b *flowBatcher) getFlowSequenceContactsForFlowSequence(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowDataLoader.getFlowSequenceContactsForFlowSequence")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	flowSequenceContactEntitiesPtr, err := b.flowService.FlowSequenceContactGetList(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get tags for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	contactEntitiesBySequenceId := make(map[string]neo4jentity.FlowSequenceContactEntities)
	for _, val := range *flowSequenceContactEntitiesPtr {
		if list, ok := contactEntitiesBySequenceId[val.DataloaderKey]; ok {
			contactEntitiesBySequenceId[val.DataloaderKey] = append(list, val)
		} else {
			contactEntitiesBySequenceId[val.DataloaderKey] = neo4jentity.FlowSequenceContactEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for id, record := range contactEntitiesBySequenceId {
		if ix, ok := keyOrder[id]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, id)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.FlowSequenceContactEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.FlowSequenceContactEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *flowBatcher) getFlowSequenceSendersForFlowSequence(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowDataLoader.getFlowSequenceSendersForFlowSequence")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	flowSequenceSenderEntitiesPtr, err := b.flowService.FlowSequenceSenderGetList(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get tags for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	senderEntitiesBySequenceId := make(map[string]neo4jentity.FlowSequenceSenderEntities)
	for _, val := range *flowSequenceSenderEntitiesPtr {
		if list, ok := senderEntitiesBySequenceId[val.DataloaderKey]; ok {
			senderEntitiesBySequenceId[val.DataloaderKey] = append(list, val)
		} else {
			senderEntitiesBySequenceId[val.DataloaderKey] = neo4jentity.FlowSequenceSenderEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for id, record := range senderEntitiesBySequenceId {
		if ix, ok := keyOrder[id]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, id)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.FlowSequenceSenderEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.FlowSequenceSenderEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
