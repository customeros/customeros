package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type EmailEventHandler struct {
	log         logger.Logger
	services    *service.Services
	grpcClients *grpc_client.Clients
}

func NewEmailEventHandler(log logger.Logger, services *service.Services, grpcClients *grpc_client.Clients) *EmailEventHandler {
	return &EmailEventHandler{
		log:         log,
		services:    services,
		grpcClients: grpcClients,
	}
}

func (h *EmailEventHandler) OnEmailValidatedV2(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailValidatedV2")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.EmailValidatedEventV2
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	emailId := email.GetEmailObjectID(evt.AggregateID, eventData.Tenant)

	data := neo4jrepository.EmailValidatedFields{
		EmailAddress:    eventData.Email,
		Domain:          eventData.Domain,
		Username:        eventData.Username,
		IsValidSyntax:   eventData.IsValidSyntax,
		IsRisky:         eventData.IsRisky,
		IsFirewalled:    eventData.IsFirewalled,
		Provider:        eventData.Provider,
		Firewall:        eventData.Firewall,
		IsCatchAll:      eventData.IsCatchAll,
		Deliverable:     eventData.Deliverable,
		IsMailboxFull:   eventData.IsMailboxFull,
		IsRoleAccount:   eventData.IsRoleAccount,
		IsFreeAccount:   eventData.IsFreeAccount,
		SmtpSuccess:     eventData.SmtpSuccess,
		ResponseCode:    eventData.ResponseCode,
		ErrorCode:       eventData.ErrorCode,
		Description:     eventData.Description,
		ValidatedAt:     eventData.ValidatedAt,
		IsPrimaryDomain: eventData.IsPrimaryDomain,
		PrimaryDomain:   eventData.PrimaryDomain,
		AlternateEmail:  eventData.AlternateEmail,
		RetryValidation: eventData.RetryValidation,
	}

	err := h.services.CommonServices.Neo4jRepositories.EmailWriteRepository.EmailValidated(ctx, eventData.Tenant, emailId, data)

	return err
}

func (h *EmailEventHandler) OnEmailDelete(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailDelete")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.EmailDeleteEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	emailId := email.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	tracing.TagTenant(span, eventData.Tenant)
	tracing.TagEntity(span, emailId)

	err := h.services.CommonServices.Neo4jRepositories.EmailWriteRepository.DeleteEmail(ctx, eventData.Tenant, emailId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("failed to delete email: %v", err)
	}

	return nil
}
