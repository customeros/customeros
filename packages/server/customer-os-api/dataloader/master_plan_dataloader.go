package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"reflect"
)

func (i *Loaders) GetMasterPlanMilestonesForMasterPlan(ctx context.Context, masterPlanId string) (*neo4jentity.MasterPlanMilestoneEntities, error) {
	thunk := i.MasterPlanMilestonesForMasterPlan.Load(ctx, dataloader.StringKey(masterPlanId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.MasterPlanMilestoneEntities)
	return &resultObj, nil
}

func (b *masterPlanBatcher) getMasterPlanMilestonesForMasterPlans(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanBatcher.getMasterPlanMilestonesForMasterPlans")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys.length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	masterPlanMilestoneEntities, err := b.masterPlanService.GetMasterPlanMilestonesForMasterPlans(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get master plan milestones for master plans")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	masterPlanMilestoneEntitiesByMasterPlanId := make(map[string]neo4jentity.MasterPlanMilestoneEntities)
	for _, val := range *masterPlanMilestoneEntities {
		if list, ok := masterPlanMilestoneEntitiesByMasterPlanId[val.DataloaderKey]; ok {
			masterPlanMilestoneEntitiesByMasterPlanId[val.DataloaderKey] = append(list, val)
		} else {
			masterPlanMilestoneEntitiesByMasterPlanId[val.DataloaderKey] = neo4jentity.MasterPlanMilestoneEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for masterPlanId, record := range masterPlanMilestoneEntitiesByMasterPlanId {
		if ix, ok := keyOrder[masterPlanId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, masterPlanId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.MasterPlanMilestoneEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.MasterPlanMilestoneEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("result.length", len(results)))

	return results
}
