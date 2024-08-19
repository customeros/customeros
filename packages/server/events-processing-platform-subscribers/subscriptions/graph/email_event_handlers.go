package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
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

func (h *EmailEventHandler) OnEmailCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.EmailCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	emailId := email.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	data := neo4jrepository.EmailCreateFields{
		RawEmail: eventData.RawEmail,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(utils.StringFirstNonEmpty(eventData.SourceFields.Source, eventData.Source)),
			SourceOfTruth: helper.GetSourceOfTruth(utils.StringFirstNonEmpty(eventData.SourceFields.SourceOfTruth, eventData.SourceOfTruth)),
			AppSource:     helper.GetAppSource(utils.StringFirstNonEmpty(eventData.SourceFields.AppSource, eventData.AppSource)),
		},
		CreatedAt: eventData.CreatedAt,
	}
	err := h.services.CommonServices.Neo4jRepositories.EmailWriteRepository.CreateEmail(ctx, eventData.Tenant, emailId, data)

	if eventData.LinkWithType != nil && eventData.LinkWithId != nil {
		if *eventData.LinkWithType == "CONTACT" {
			err = h.services.CommonServices.Neo4jRepositories.EmailWriteRepository.LinkWithContact(ctx, eventData.Tenant, *eventData.LinkWithId, emailId, "Work", true)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
		//TODO continue and generify
	}

	return err
}

func (h *EmailEventHandler) OnEmailUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)
	tracing.LogObjectAsJson(span, "eventData", evt)

	var eventData event.EmailUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	emailId := email.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("rawEmail", eventData.RawEmail))

	emailDbNode, err := h.services.CommonServices.Neo4jRepositories.EmailReadRepository.GetById(ctx, eventData.Tenant, emailId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	emailBeforeUpdate := neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)

	err = h.services.CommonServices.Neo4jRepositories.EmailWriteRepository.UpdateEmail(ctx, eventData.Tenant, emailId, eventData.RawEmail, eventData.Source)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// email address updated
	if emailBeforeUpdate.RawEmail != eventData.RawEmail && emailBeforeUpdate.Email != eventData.RawEmail {
		err = h.services.CommonServices.Neo4jRepositories.EmailWriteRepository.CleanEmailValidation(ctx, eventData.Tenant, emailId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil
		}

		if strings.Contains(eventData.RawEmail, "@") {
			ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
			_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
				return h.grpcClients.EmailClient.RequestEmailValidation(ctx, &emailpb.RequestEmailValidationGrpcRequest{
					Tenant:    eventData.Tenant,
					Id:        emailId,
					AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to request email validation for emailId: %s, tenant: %s, err: %s", emailId, eventData.Tenant, err.Error())
				return nil
			}
		}
	}

	return nil
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
		EmailAddress:   eventData.Email,
		Domain:         eventData.Domain,
		Username:       eventData.Username,
		IsValidSyntax:  eventData.IsValidSyntax,
		IsRisky:        eventData.IsRisky,
		IsFirewalled:   eventData.IsFirewalled,
		Provider:       eventData.Provider,
		Firewall:       eventData.Firewall,
		IsCatchAll:     eventData.IsCatchAll,
		CanConnectSMTP: eventData.CanConnectSMTP,
		IsDeliverable:  eventData.IsDeliverable,
		IsMailboxFull:  eventData.IsMailboxFull,
		IsRoleAccount:  eventData.IsRoleAccount,
		IsFreeAccount:  eventData.IsFreeAccount,
		SmtpSuccess:    eventData.SmtpSuccess,
		ResponseCode:   eventData.ResponseCode,
		ErrorCode:      eventData.ErrorCode,
		Description:    eventData.Description,
		SmtpResponse:   eventData.SmtpResponse,
		ValidatedAt:    eventData.ValidatedAt,
	}

	err := h.services.CommonServices.Neo4jRepositories.EmailWriteRepository.EmailValidated(ctx, eventData.Tenant, emailId, data)

	return err
}
