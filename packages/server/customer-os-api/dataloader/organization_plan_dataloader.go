package dataloader

import (
	"context"
	"reflect"

	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (i *Loaders) GetOrganizationPlanMilestonesForOrganizationPlan(ctx context.Context, orgPlanId string) (*neo4jentity.OrganizationPlanMilestoneEntities, error) {
	thunk := i.OrganizationPlanMilestonesForOrganizationPlan.Load(ctx, dataloader.StringKey(orgPlanId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.OrganizationPlanMilestoneEntities)
	return &resultObj, nil
}

func (b *organizationPlanBatcher) getOrganizationPlanMilestonesForOrganizationPlans(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "organizationPlanBatcher.getOrganizationPlanMilestonesForOrganizationPlans")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys.length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	opMilestoneEntities, err := b.organizationPlanService.GetOrganizationPlanMilestonesForOrganizationPlans(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get org plan milestones for org plans")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	opMilestoneEntitiesByOrgPlanId := make(map[string]neo4jentity.OrganizationPlanMilestoneEntities)
	for _, val := range *opMilestoneEntities {
		if list, ok := opMilestoneEntitiesByOrgPlanId[val.DataloaderKey]; ok {
			opMilestoneEntitiesByOrgPlanId[val.DataloaderKey] = append(list, val)
		} else {
			opMilestoneEntitiesByOrgPlanId[val.DataloaderKey] = neo4jentity.OrganizationPlanMilestoneEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for orgPlanId, record := range opMilestoneEntitiesByOrgPlanId {
		if ix, ok := keyOrder[orgPlanId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, orgPlanId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.OrganizationPlanMilestoneEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.OrganizationPlanMilestoneEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("result.length", len(results)))

	return results
}
