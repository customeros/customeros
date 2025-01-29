package oldlisteners

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

func Handle_FlowComputeParticipantsRequirements(ctx context.Context, dependencies *model.DependencyContainer, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.FlowComputeParticipantsRequirements")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message := input.(*dto.Event)

	flow, err := dependencies.CommonServices.FlowService.FlowGetById(ctx, message.Event.EntityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flow == nil {
		err = errors.New("flow not found")
		tracing.TraceErr(span, err)
		return err
	}

	flowRequirements, err := dependencies.CommonServices.FlowExecutionService.GetFlowRequirements(ctx, flow.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowParticipants, err := dependencies.CommonServices.FlowService.FlowParticipantGetList(ctx, []string{flow.Id})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, dependencies.Neo4jRepositories.Neo4jDriver, dependencies.Neo4jRepositories.Database, nil, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		for _, v := range *flowParticipants {
			err := dependencies.CommonServices.FlowExecutionService.UpdateParticipantFlowRequirements(ctx, txWithPostCommit, &v, flowRequirements)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
