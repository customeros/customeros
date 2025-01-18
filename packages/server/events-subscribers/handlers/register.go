package handlers

import (
	"context"
	"reflect"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/listeners"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/model"
)

func InitHandlerRegistration(eventsService *events.EventsService, dependencies *model.DependencyContainer) {

	// Register Flow handlers
	eventsService.RegisterHandler(dto.FlowOn{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.Handle_FlowOn(ctx, dependencies.CommonServices, event)
		},
		EventType: reflect.TypeOf(dto.FlowOn{}).Name(),
		DataType:  reflect.TypeOf(dto.FlowOn{}),
	})

	eventsService.RegisterHandler(dto.FlowParticipantSchedule{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.Handle_FlowParticipantSchedule(ctx, dependencies.CommonServices, event)
		},
		EventType: reflect.TypeOf(dto.FlowParticipantSchedule{}).Name(),
		DataType:  reflect.TypeOf(dto.FlowParticipantSchedule{}),
	})

	eventsService.RegisterHandler(dto.FlowComputeParticipantsRequirements{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.Handle_FlowComputeParticipantsRequirements(ctx, dependencies, event)
		},
		EventType: reflect.TypeOf(dto.FlowComputeParticipantsRequirements{}).Name(),
		DataType:  reflect.TypeOf(dto.FlowComputeParticipantsRequirements{}),
	})

	eventsService.RegisterHandler(dto.FlowParticipantGoalAchieved{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.Handle_FlowParticipantGoalAchieved(ctx, dependencies, event)
		},
		EventType: reflect.TypeOf(dto.FlowParticipantGoalAchieved{}).Name(),
		DataType:  reflect.TypeOf(dto.FlowParticipantGoalAchieved{}),
	})

	// Mailstack handlers
	eventsService.RegisterHandler(dto.MailstackProvisionBuyRequest{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.Handle_MailstackProvisionBuyRequest(ctx, dependencies, event)
		},
		EventType: reflect.TypeOf(dto.MailstackProvisionBuyRequest{}).Name(),
		DataType:  reflect.TypeOf(dto.MailstackProvisionBuyRequest{}),
	})

	eventsService.RegisterHandler(dto.MailstackProvisionMailbox{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.Handle_MailstackProvisionMailbox(ctx, dependencies, event)
		},
		EventType: reflect.TypeOf(dto.MailstackProvisionMailbox{}).Name(),
		DataType:  reflect.TypeOf(dto.MailstackProvisionMailbox{}),
	})

	// Contact handlers
	eventsService.RegisterHandler(dto.AddSocialToContact{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.OnSocialAddedToContact(ctx, dependencies, event)
		},
		EventType: reflect.TypeOf(dto.AddSocialToContact{}).Name(),
		DataType:  reflect.TypeOf(dto.AddSocialToContact{}),
	})

	eventsService.RegisterHandler(dto.RequestEnrichContact{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.OnRequestedEnrichContact(ctx, dependencies, event)
		},
		EventType: reflect.TypeOf(dto.RequestEnrichContact{}).Name(),
		DataType:  reflect.TypeOf(dto.RequestEnrichContact{}),
	})

	eventsService.RegisterHandler(dto.HideContact{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.OnContactHidden(ctx, dependencies, event)
		},
		EventType: reflect.TypeOf(dto.HideContact{}).Name(),
		DataType:  reflect.TypeOf(dto.HideContact{}),
	})

	// Organization handlers
	eventsService.RegisterHandler(dto.RequestRefreshLastTouchpoint{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.OnRequestLastTouchpointRefresh(ctx, dependencies.CommonServices, event)
		},
		EventType: reflect.TypeOf(dto.RequestRefreshLastTouchpoint{}).Name(),
		DataType:  reflect.TypeOf(dto.RequestRefreshLastTouchpoint{}),
	})

	eventsService.RegisterHandler(dto.RequestEnrichOrganization{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.OnRequestedEnrichOrganization(ctx, dependencies, event)
		},
		EventType: reflect.TypeOf(dto.RequestEnrichOrganization{}).Name(),
		DataType:  reflect.TypeOf(dto.RequestEnrichOrganization{}),
	})

	// Automation Engine
	eventsService.RegisterHandler(dto.WebhookEvent{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.OnWebhookEventCreated(ctx, dependencies, event)
		},
		EventType: reflect.TypeOf(dto.WebhookEvent{}).Name(),
		DataType:  reflect.TypeOf(dto.WebhookEvent{}),
	})

	// Flow Engine

	// Email handlers
	eventsService.RegisterHandler(dto.RequestValidateEmail{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.OnRequestedValidateEmail(ctx, dependencies, event)
		},
		EventType: reflect.TypeOf(dto.RequestValidateEmail{}).Name(),
		DataType:  reflect.TypeOf(dto.RequestValidateEmail{}),
	})

	// Flow Engine handlers
	eventsService.RegisterHandler(dto.WebhookEvent{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.OnWebhookEventCreated(ctx, dependencies, event)
		},
		EventType: reflect.TypeOf(dto.WebhookEvent{}).Name(),
		DataType:  reflect.TypeOf(dto.WebhookEvent{}),
	})

	eventsService.RegisterHandler(dto.FlowAgentEvent{}, interfaces.EventHandler{
		HandlerFunc: func(ctx context.Context, event any) error {
			return listeners.OnFlowAgentEventCreated(ctx, dependencies.CommonServices, event)
		},
		EventType: reflect.TypeOf(dto.FlowAgentEvent{}).Name(),
		DataType:  reflect.TypeOf(dto.FlowAgentEvent{}),
	})

}
