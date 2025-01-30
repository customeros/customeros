package events_listeners

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowComputeParticipantsRequirementsListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return l.handle(ctx, event.Event.EntityId)
}

func (l *FlowComputeParticipantsRequirementsListener) handle(ctx context.Context, entityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowComputeParticipantsRequirementsListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	flow, err := l.dependencies.CommonServices.FlowService.FlowGetById(ctx, entityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flow == nil {
		err = errors.New("flow not found")
		tracing.TraceErr(span, err)
		return err
	}

	flowRequirements, err := l.dependencies.CommonServices.FlowExecutionService.GetFlowRequirements(ctx, flow.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowParticipants, err := l.dependencies.CommonServices.FlowService.FlowParticipantGetList(ctx, []string{flow.Id})
	if err != nil {
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
