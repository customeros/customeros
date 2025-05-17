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

func (i *Loaders) GetSubmitterParticipantsForIssue(ctx context.Context, issueId string) (*neo4jentity.IssueParticipants, error) {
	thunk := i.SubmitterParticipantsForIssue.Load(ctx, dataloader.StringKey(issueId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.IssueParticipants)
	return &resultObj, nil
}

func (i *Loaders) GetReporterParticipantsForIssue(ctx context.Context, issueId string) (*neo4jentity.IssueParticipants, error) {
	thunk := i.ReporterParticipantsForIssue.Load(ctx, dataloader.StringKey(issueId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.IssueParticipants)
	return &resultObj, nil
}

func (i *Loaders) GetAssigneeParticipantsForIssue(ctx context.Context, issueId string) (*neo4jentity.IssueParticipants, error) {
	thunk := i.AssigneeParticipantsForIssue.Load(ctx, dataloader.StringKey(issueId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.IssueParticipants)
	return &resultObj, nil
}

func (i *Loaders) GetFollowerParticipantsForIssue(ctx context.Context, issueId string) (*neo4jentity.IssueParticipants, error) {
	thunk := i.FollowerParticipantsForIssue.Load(ctx, dataloader.StringKey(issueId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.IssueParticipants)
	return &resultObj, nil
}

func (b *issueParticipantBatcher) getSubmitterParticipantsForIssues(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IssueParticipantDataLoader.getSubmitterParticipantsForIssues")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	participantEntitiesPtr, err := b.issueService.GetSubmitterParticipantsForIssues(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	participantEntitiesGrouped := make(map[string]neo4jentity.IssueParticipants)
	for _, val := range *participantEntitiesPtr {
		if list, ok := participantEntitiesGrouped[val.GetDataloaderKey()]; ok {
			participantEntitiesGrouped[val.GetDataloaderKey()] = append(list, val)
		} else {
			participantEntitiesGrouped[val.GetDataloaderKey()] = neo4jentity.IssueParticipants{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range participantEntitiesGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.IssueParticipants{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.IssueParticipants{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *issueParticipantBatcher) getReporterParticipantsForIssues(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IssueParticipantDataLoader.getReporterParticipantsForIssues")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	participantEntitiesPtr, err := b.issueService.GetReporterParticipantsForIssues(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get issue participants")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	participantEntitiesGrouped := make(map[string]neo4jentity.IssueParticipants)
	for _, val := range *participantEntitiesPtr {
		if list, ok := participantEntitiesGrouped[val.GetDataloaderKey()]; ok {
			participantEntitiesGrouped[val.GetDataloaderKey()] = append(list, val)
		} else {
			participantEntitiesGrouped[val.GetDataloaderKey()] = neo4jentity.IssueParticipants{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range participantEntitiesGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.IssueParticipants{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.IssueParticipants{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *issueParticipantBatcher) getAssigneeParticipantsForIssues(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IssueParticipantDataLoader.getAssigneeParticipantsForIssues")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	participantEntitiesPtr, err := b.issueService.GetAssigneeParticipantsForIssues(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get issue participants")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	participantEntitiesGrouped := make(map[string]neo4jentity.IssueParticipants)
	for _, val := range *participantEntitiesPtr {
		if list, ok := participantEntitiesGrouped[val.GetDataloaderKey()]; ok {
			participantEntitiesGrouped[val.GetDataloaderKey()] = append(list, val)
		} else {
			participantEntitiesGrouped[val.GetDataloaderKey()] = neo4jentity.IssueParticipants{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range participantEntitiesGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.IssueParticipants{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.IssueParticipants{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *issueParticipantBatcher) getFollowerParticipantsForIssues(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IssueParticipantDataLoader.getFollowerParticipantsForIssues")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	participantEntitiesPtr, err := b.issueService.GetFollowerParticipantsForIssues(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get issue participants")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	participantEntitiesGrouped := make(map[string]neo4jentity.IssueParticipants)
	for _, val := range *participantEntitiesPtr {
		if list, ok := participantEntitiesGrouped[val.GetDataloaderKey()]; ok {
			participantEntitiesGrouped[val.GetDataloaderKey()] = append(list, val)
		} else {
			participantEntitiesGrouped[val.GetDataloaderKey()] = neo4jentity.IssueParticipants{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationId, record := range participantEntitiesGrouped {
		ix, ok := keyOrder[organizationId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, organizationId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.IssueParticipants{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.IssueParticipants{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
