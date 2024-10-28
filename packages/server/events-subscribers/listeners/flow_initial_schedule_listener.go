package listeners

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
)

func Handle_FlowInitialSchedule(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.Handle_FlowInitialSchedule")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message := input.(*events.Event)

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

	if flow.Status != neo4jentity.FlowStatusScheduling {
		return nil
	}

	flowParticipants, err := services.FlowService.FlowParticipantGetList(ctx, []string{flow.Id})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = utils.ExecuteWriteInTransaction(ctx, services.Neo4jRepositories.Neo4jDriver, services.Neo4jRepositories.Database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
		for _, v := range *flowParticipants {
			err := services.FlowExecutionService.ScheduleFlow(ctx, &tx, flow.Id, &v)
			if err != nil {
				return nil, err
			}
		}

		flow.Status = neo4jentity.FlowStatusActive

		_, err := services.Neo4jRepositories.FlowWriteRepository.Merge(ctx, &tx, flow)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
