package listeners

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
)

func Handle_FlowOn(ctx context.Context, services *service.CommonServices, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.Handle_FlowOn")
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

	if flow.Status != neo4jentity.FlowStatusOn {
		return nil
	}

	flowParticipants, err := services.FlowService.FlowParticipantGetList(ctx, []string{flow.Id})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, v := range *flowParticipants {
		err := services.Events.Publisher.PublishEventOnExchange(ctx, v.Id, model.FLOW_PARTICIPANT, dto.FlowParticipantSchedule{}, events.EventsExchangeName, events.EventsFlowParticipantScheduleRoutingKey)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func Handle_FlowParticipantSchedule(ctx context.Context, services *service.CommonServices, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.Handle_FlowParticipantSchedule")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message := input.(*dto.Event)

	flowParticipant, err := services.FlowService.FlowParticipantById(ctx, message.Event.EntityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowParticipant == nil {
		err = errors.New("flow participant not found")
		tracing.TraceErr(span, err)
		return err
	}

	flow, err := services.FlowService.FlowGetByParticipantId(ctx, nil, flowParticipant.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = services.FlowExecutionService.ScheduleFlow(ctx, nil, flow.Id, flowParticipant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
