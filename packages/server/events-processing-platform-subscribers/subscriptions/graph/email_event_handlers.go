package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type EmailEventHandler struct {
	Repositories *repository.Repositories
}

func NewEmailEventHandler(repositories *repository.Repositories) *EmailEventHandler {
	return &EmailEventHandler{
		Repositories: repositories,
	}
}

func (h *EmailEventHandler) OnEmailCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.EmailCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailId := aggregate.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.EmailCreateFields{
		RawEmail: eventData.RawEmail,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(utils.StringFirstNonEmpty(eventData.SourceFields.Source, eventData.Source)),
			SourceOfTruth: helper.GetSourceOfTruth(utils.StringFirstNonEmpty(eventData.SourceFields.SourceOfTruth, eventData.SourceOfTruth)),
			AppSource:     helper.GetAppSource(utils.StringFirstNonEmpty(eventData.SourceFields.AppSource, eventData.AppSource)),
		},
		CreatedAt: eventData.CreatedAt,
	}
	err := h.Repositories.Neo4jRepositories.EmailWriteRepository.CreateEmail(ctx, eventData.Tenant, emailId, data)

	return err
}

func (h *EmailEventHandler) OnEmailUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.EmailUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailId := aggregate.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.Neo4jRepositories.EmailWriteRepository.UpdateEmail(ctx, eventData.Tenant, emailId, eventData.Source)

	return err
}

func (h *EmailEventHandler) OnEmailValidationFailed(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailValidationFailed")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.EmailFailedValidationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailId := aggregate.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.Neo4jRepositories.EmailWriteRepository.FailEmailValidation(ctx, eventData.Tenant, emailId, eventData.ValidationError)

	return err
}

func (h *EmailEventHandler) OnEmailValidated(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailValidated")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.EmailValidatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailId := aggregate.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.EmailValidatedFields{
		ValidationError: eventData.ValidationError,
		EmailAddress:    eventData.EmailAddress,
		Domain:          eventData.Domain,
		AcceptsMail:     eventData.AcceptsMail,
		CanConnectSmtp:  eventData.CanConnectSmtp,
		HasFullInbox:    eventData.HasFullInbox,
		IsCatchAll:      eventData.IsCatchAll,
		IsDeliverable:   eventData.IsDeliverable,
		IsDisabled:      eventData.IsDisabled,
		IsValidSyntax:   eventData.IsValidSyntax,
		Username:        eventData.Username,
		ValidatedAt:     eventData.ValidatedAt,
		IsReachable:     eventData.IsReachable,
	}
	err := h.Repositories.Neo4jRepositories.EmailWriteRepository.EmailValidated(ctx, eventData.Tenant, emailId, data)

	return err
}
