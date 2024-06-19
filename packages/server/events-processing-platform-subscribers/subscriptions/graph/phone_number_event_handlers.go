package graph

import (
	"context"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type PhoneNumberEventHandler struct {
	repositories *repository.Repositories
}

func NewPhoneNumberEventHandler(repositories *repository.Repositories) *PhoneNumberEventHandler {
	return &PhoneNumberEventHandler{
		repositories: repositories,
	}
}

func (h *PhoneNumberEventHandler) OnPhoneNumberCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberEventHandler.OnPhoneNumberCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.PhoneNumberCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.PhoneNumberCreateFields{
		RawPhoneNumber: eventData.RawPhoneNumber,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.SourceFields.Source),
			SourceOfTruth: helper.GetSourceOfTruth(eventData.SourceFields.SourceOfTruth),
			AppSource:     helper.GetAppSource(eventData.SourceFields.AppSource),
		},
		CreatedAt: eventData.CreatedAt,
	}
	err := h.repositories.Neo4jRepositories.PhoneNumberWriteRepository.CreatePhoneNumber(ctx, eventData.Tenant, phoneNumberId, data)

	return err
}

func (h *PhoneNumberEventHandler) OnPhoneNumberUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberEventHandler.OnPhoneNumberUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.PhoneNumberUpdatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.PhoneNumberWriteRepository.UpdatePhoneNumber(ctx, eventData.Tenant, phoneNumberId, eventData.Source)

	return err
}

func (e *PhoneNumberEventHandler) OnPhoneNumberValidated(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberEventHandler.OnPhoneNumberValidated")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.PhoneNumberValidatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.PhoneNumberValidateFields{
		E164:          eventData.E164,
		CountryCodeA2: eventData.CountryCodeA2,
		ValidatedAt:   eventData.ValidatedAt,
		Source:        constants.SourceOpenline,
		AppSource:     "validation-api",
	}
	err := e.repositories.Neo4jRepositories.PhoneNumberWriteRepository.PhoneNumberValidated(ctx, eventData.Tenant, phoneNumberId, data)

	return err
}

func (h *PhoneNumberEventHandler) OnPhoneNumberValidationFailed(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberEventHandler.OnPhoneNumberValidationFailed")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.PhoneNumberFailedValidationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.PhoneNumberWriteRepository.FailPhoneNumberValidation(ctx, eventData.Tenant, phoneNumberId, eventData.ValidationError)

	return err
}
