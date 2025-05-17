package events_listeners

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type FlowComputeParticipantsRequirementsListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewFlowComputeParticipantsRequirementsListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &FlowComputeParticipantsRequirementsListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.FlowComputeParticipantsRequirements](), // subscribed event
			events.QueueEvents, // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *FlowComputeParticipantsRequirementsListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartListenerSpan(ctx, "FlowComputeParticipantsRequirementsListener.Handle")
	defer spans.Finish()
	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return l.handle(ctx, event.Event.EntityId)
}

func (l *FlowComputeParticipantsRequirementsListener) handle(ctx context.Context, entityId string) error {
	spans, ctx := telemetry.StartListenerSpan(ctx, "FlowComputeParticipantsRequirementsListener.handle")
	defer spans.Finish()
	spans.LogKV("entityId", entityId)

	flow, err := l.dependencies.CommonServices.FlowService.FlowGetById(ctx, entityId)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if flow == nil {
		err = errors.New("flow not found")
		spans.TraceError(err)
		return err
	}

	flowRequirements, err := l.dependencies.CommonServices.FlowExecutionService.GetFlowRequirements(ctx, flow.Id)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	flowParticipants, err := l.dependencies.CommonServices.FlowService.FlowParticipantGetList(ctx, []string{flow.Id})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, l.dependencies.Neo4jRepositories.Neo4jDriver, l.dependencies.Neo4jRepositories.Database, nil, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		for _, v := range *flowParticipants {
			err := l.dependencies.CommonServices.FlowExecutionService.UpdateParticipantFlowRequirements(ctx, txWithPostCommit, &v, flowRequirements)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
