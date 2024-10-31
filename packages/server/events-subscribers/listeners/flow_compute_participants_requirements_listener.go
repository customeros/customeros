package listeners

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

func Handle_FlowComputeParticipantsRequirements(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.FlowComputeParticipantsRequirements")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message := input.(*dto.Event)

	flow, err := services.FlowService.FlowGetById(ctx, message.Event.EntityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flow == nil {
		err = errors.New("flow not found")
		tracing.TraceErr(span, err)
		return err
	}

	flowRequirements, err := services.FlowExecutionService.GetFlowRequirements(ctx, flow.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowParticipants, err := services.FlowService.FlowParticipantGetList(ctx, []string{flow.Id})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowParticipantsMap := make(map[model.EntityType][]string)

	_, err = utils.ExecuteWriteInTransaction(ctx, services.Neo4jRepositories.Neo4jDriver, services.Neo4jRepositories.Database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
		for _, v := range *flowParticipants {
			updated, err := services.FlowExecutionService.UpdateParticipantFlowRequirements(ctx, &tx, &v, flowRequirements)
			if err != nil {
				return nil, err
			}

			if updated {
				flowParticipantsMap[v.EntityType] = append(flowParticipantsMap[v.EntityType], v.Id)
			}
		}

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for entityType, flowParticipantsIds := range flowParticipantsMap {
		services.RabbitMQService.PublishEventCompletedBulk(ctx, message.Event.Tenant, flowParticipantsIds, entityType, utils.NewEventCompletedDetails().WithUpdate())
	}

	return nil
}
