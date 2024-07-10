package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/country"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type CountryEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewCountryEventHandler(log logger.Logger, repositories *repository.Repositories) *CountryEventHandler {
	return &CountryEventHandler{
		log:          log,
		repositories: repositories,
	}
}

func (h *CountryEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryEventHandler.OnCountryNew")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData country.CountryCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	id := country.GetCountryObjectID(evt.GetAggregateID())
	span.SetTag(tracing.SpanTagEntityId, id)

	err := h.repositories.Neo4jRepositories.CountryWriteRepository.CreateCountry(ctx, id, eventData.Name, eventData.CodeA2, eventData.CodeA3, eventData.PhoneCode, eventData.CreatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving country %s: %s", id, err.Error())
		return err
	}
	return err
}
