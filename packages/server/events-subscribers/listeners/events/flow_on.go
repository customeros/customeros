package events_listeners

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	common_model "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type FlowOnListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewFlowOnListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &FlowOnListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.FlowOn](), // subscribed event
			events.QueueEvents,                // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *FlowOnListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartListenerSpan(ctx, "FlowOnListener.Handle")
	defer spans.Finish()
	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return l.handle(ctx, event.Event.EntityId)
}

func (l *FlowOnListener) handle(ctx context.Context, entityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowOnListener.handle")
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

	if flow.Status != neo4j_entity.FlowStatusOn {
		return nil
	}

	flowParticipants, err := l.dependencies.CommonServices.FlowService.FlowParticipantGetList(ctx, []string{flow.Id})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, v := range *flowParticipants {
		err := l.dependencies.CommonServices.Events.Publisher.PublishDirectEvent(ctx, v.Id, common_model.FLOW_PARTICIPANT, dto.FlowParticipantSchedule{})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
