package direct_listeners

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type FlowParticipantScheduleListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewFlowParcipantScheduleListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &FlowParticipantScheduleListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.FlowParticipantSchedule](), // subscribed event
			events.QueueFlowParticipantSchedule,                // listening on Direct queue
		),
		dependencies: deps,
	}
}

func (l *FlowParticipantScheduleListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartListenerSpan(ctx, "FlowParticipantScheduleListener.Handle")
	defer spans.Finish()
	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return l.handle(ctx, event.Event.EntityId)
}

func (l *FlowParticipantScheduleListener) handle(ctx context.Context, entityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowParticipantScheduleListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	flowParticipant, err := l.dependencies.CommonServices.FlowService.FlowParticipantById(ctx, entityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowParticipant == nil {
		err = errors.New("flow participant not found")
		tracing.TraceErr(span, err)
		return err
	}

	flow, err := l.dependencies.CommonServices.FlowService.FlowGetByParticipantId(ctx, nil, flowParticipant.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = l.dependencies.CommonServices.FlowExecutionService.ScheduleFlow(ctx, nil, flow.Id, flowParticipant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
