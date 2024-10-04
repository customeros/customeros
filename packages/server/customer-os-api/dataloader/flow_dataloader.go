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

func (i *Loaders) GetFlowContactsForFlow(ctx context.Context, flowId string) (*neo4jentity.FlowParticipantEntities, error) {
	thunk := i.FlowContactsForFlow.Load(ctx, dataloader.StringKey(flowId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.FlowParticipantEntities)
	return &resultObj, nil
}

func (i *Loaders) GetFlowActionSendersForFlowAction(ctx context.Context, actionId string) (*neo4jentity.FlowActionSenderEntities, error) {
	thunk := i.FlowActionSendersForFlowAction.Load(ctx, dataloader.StringKey(actionId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.FlowActionSenderEntities)
	return &resultObj, nil
}

func (i *Loaders) GetFlowActionsForFlow(ctx context.Context, flowId string) (*neo4jentity.FlowActionEntities, error) {
	thunk := i.FlowActionsForFlow.Load(ctx, dataloader.StringKey(flowId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.FlowActionEntities)
	return &resultObj, nil
}

func (i *Loaders) GetFlowsWithContact(ctx context.Context, contactId string) (*neo4jentity.FlowEntities, error) {
	thunk := i.FlowsWithContact.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.FlowEntities)
	return &resultObj, nil
}

func (b *flowBatcher) getFlowContactsForFlow(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowDataLoader.getFlowContactsForFlow")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	flowSequenceContactEntitiesPtr, err := b.flowService.FlowParticipantGetList(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get tags for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	contactEntitiesBySequenceId := make(map[string]neo4jentity.FlowParticipantEntities)
	for _, val := range *flowSequenceContactEntitiesPtr {
		if list, ok := contactEntitiesBySequenceId[val.DataloaderKey]; ok {
			contactEntitiesBySequenceId[val.DataloaderKey] = append(list, val)
		} else {
			contactEntitiesBySequenceId[val.DataloaderKey] = neo4jentity.FlowParticipantEntities{val}
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
		results[ix] = &dataloader.Result{Data: neo4jentity.FlowParticipantEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.FlowParticipantEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *flowBatcher) getFlowActionSendersForFlowAction(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowDataLoader.getFlowActionSendersForFlowAction")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	flowSequenceSenderEntitiesPtr, err := b.flowService.FlowActionSenderGetList(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get tags for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	senderEntitiesBySequenceId := make(map[string]neo4jentity.FlowActionSenderEntities)
	for _, val := range *flowSequenceSenderEntitiesPtr {
		if list, ok := senderEntitiesBySequenceId[val.DataloaderKey]; ok {
			senderEntitiesBySequenceId[val.DataloaderKey] = append(list, val)
		} else {
			senderEntitiesBySequenceId[val.DataloaderKey] = neo4jentity.FlowActionSenderEntities{val}
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
		results[ix] = &dataloader.Result{Data: neo4jentity.FlowActionSenderEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.FlowActionSenderEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *flowBatcher) getFlowActionsForFlow(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowDataLoader.getFlowActionsForFlow")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	flowSequenceStepEntitiesPtr, err := b.flowService.FlowActionGetList(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get flow sequece steps for flow sequence")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	entitiesBySequenceId := make(map[string]neo4jentity.FlowActionEntities)
	for _, val := range *flowSequenceStepEntitiesPtr {
		if list, ok := entitiesBySequenceId[val.DataloaderKey]; ok {
			entitiesBySequenceId[val.DataloaderKey] = append(list, val)
		} else {
			entitiesBySequenceId[val.DataloaderKey] = neo4jentity.FlowActionEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for id, record := range entitiesBySequenceId {
		if ix, ok := keyOrder[id]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, id)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.FlowActionEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.FlowActionEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *flowBatcher) getFlowsWithContact(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowDataLoader.getFlowsWithContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	flowSequenceEntitiesPtr, err := b.flowService.FlowsGetListWithContact(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get flow sequece steps for contact")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	entitiesByContactId := make(map[string]neo4jentity.FlowEntities)
	for _, val := range *flowSequenceEntitiesPtr {
		if list, ok := entitiesByContactId[val.DataloaderKey]; ok {
			entitiesByContactId[val.DataloaderKey] = append(list, val)
		} else {
			entitiesByContactId[val.DataloaderKey] = neo4jentity.FlowEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for id, record := range entitiesByContactId {
		if ix, ok := keyOrder[id]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, id)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.FlowEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.FlowEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
