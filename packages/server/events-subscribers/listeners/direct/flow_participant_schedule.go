package direct_listeners

import (
	"context"
	"fmt"

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowParticipantScheduleListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = l.validateMessage(ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
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

func (l *FlowParticipantScheduleListener) validateMessage(ctx context.Context, event *dto.Event) (*dto.FlowParticipantSchedule, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowParticipantScheduleListener.validateMessage")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message, ok := event.Event.Data.(*dto.FlowParticipantSchedule)
	if !ok {
		err := fmt.Errorf("expected FlowParticipantSchedule, got %T", event.Event.Data)
		tracing.TraceErr(span, err)
		return nil, err
	}

	if event.Event.EntityId == "" {
		err := errors.New("EntityId not set on event")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return message, nil
}
