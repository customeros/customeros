package dataloader

import (
	"context"
	"errors"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
	"reflect"
)

func (i *Loaders) GetExternalSystemsForComment(ctx context.Context, commentId string) (*neo4jentity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForComment.Load(ctx, dataloader.StringKey(commentId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForIssue(ctx context.Context, issueId string) (*neo4jentity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForIssue.Load(ctx, dataloader.StringKey(issueId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForOrganization(ctx context.Context, organizationId string) (*neo4jentity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForLogEntry(ctx context.Context, logEntryId string) (*neo4jentity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForLogEntry.Load(ctx, dataloader.StringKey(logEntryId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForMeeting(ctx context.Context, meetingId string) (*neo4jentity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForMeeting.Load(ctx, dataloader.StringKey(meetingId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForInteractionEvent(ctx context.Context, ieId string) (*neo4jentity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForInteractionEvent.Load(ctx, dataloader.StringKey(ieId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.ExternalSystemEntities)
	return &resultObj, nil
}

func (b *externalSystemBatcher) getExternalSystemsForComments(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemDataLoader.getExternalSystemsForComments")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	return b.getExternalSystemsFor(ctx, keys, commonModel.COMMENT, spans)
}

func (b *externalSystemBatcher) getExternalSystemsForIssues(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemDataLoader.getExternalSystemsForIssues")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	return b.getExternalSystemsFor(ctx, keys, commonModel.ISSUE, spans)
}

func (b *externalSystemBatcher) getExternalSystemsForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemDataLoader.getExternalSystemsForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	return b.getExternalSystemsFor(ctx, keys, commonModel.ORGANIZATION, spans)
}

func (b *externalSystemBatcher) getExternalSystemsForContracts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemDataLoader.getExternalSystemsForContracts")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	return b.getExternalSystemsFor(ctx, keys, commonModel.CONTRACT, spans)
}

func (b *externalSystemBatcher) getExternalSystemsForOpportunities(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemDataLoader.getExternalSystemsForOpportunities")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	return b.getExternalSystemsFor(ctx, keys, commonModel.OPPORTUNITY, spans)
}

func (b *externalSystemBatcher) getExternalSystemsForServiceLineItems(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemDataLoader.getExternalSystemsForServiceLineItems")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	return b.getExternalSystemsFor(ctx, keys, commonModel.SERVICE_LINE_ITEM, spans)
}

func (b *externalSystemBatcher) getExternalSystemsForLogEntries(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemDataLoader.getExternalSystemsForLogEntries")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	return b.getExternalSystemsFor(ctx, keys, commonModel.LOG_ENTRY, spans)
}

func (b *externalSystemBatcher) getExternalSystemsForMeetings(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemDataLoader.getExternalSystemsForMeetings")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	return b.getExternalSystemsFor(ctx, keys, commonModel.MEETING, spans)
}

func (b *externalSystemBatcher) getExternalSystemsForInteractionEvents(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemDataLoader.getExternalSystemsForInteractionEvents")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	return b.getExternalSystemsFor(ctx, keys, commonModel.INTERACTION_EVENT, spans)
}

func (b *externalSystemBatcher) getExternalSystemsFor(ctx context.Context, keys dataloader.Keys, entityType commonModel.EntityType, spans *telemetry.Spans) []*dataloader.Result {
	ids, keyOrder := sortKeys(keys)

	ExternalSystemsPtr, err := b.externalSystemService.GetExternalSystemsForEntities(ctx, ids, entityType)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get external systems for entities")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	ExternalSystemsByEntityId := make(map[string]neo4jentity.ExternalSystemEntities)
	for _, val := range *ExternalSystemsPtr {
		if list, ok := ExternalSystemsByEntityId[val.DataloaderKey]; ok {
			ExternalSystemsByEntityId[val.DataloaderKey] = append(list, val)
		} else {
			ExternalSystemsByEntityId[val.DataloaderKey] = neo4jentity.ExternalSystemEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for entityId, record := range ExternalSystemsByEntityId {
		if ix, ok := keyOrder[entityId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, entityId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.ExternalSystemEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.ExternalSystemEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (i *Loaders) GetExternalSystemsForContract(ctx context.Context, contractId string) (*neo4jentity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForContract.Load(ctx, dataloader.StringKey(contractId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForOpportunity(ctx context.Context, opportunityId string) (*neo4jentity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForOpportunity.Load(ctx, dataloader.StringKey(opportunityId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.ExternalSystemEntities)
	return &resultObj, nil
}

func (i *Loaders) GetExternalSystemsForServiceLineItem(ctx context.Context, serviceLineItemId string) (*neo4jentity.ExternalSystemEntities, error) {
	thunk := i.ExternalSystemsForServiceLineItem.Load(ctx, dataloader.StringKey(serviceLineItemId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.ExternalSystemEntities)
	return &resultObj, nil
}
