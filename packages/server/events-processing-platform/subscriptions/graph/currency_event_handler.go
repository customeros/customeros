package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/currency"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type CurrencyEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewCurrencyEventHandler(log logger.Logger, repositories *repository.Repositories) *CurrencyEventHandler {
	return &CurrencyEventHandler{
		log:          log,
		repositories: repositories,
	}
}

func (h *CurrencyEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CurrencyEventHandler.OnCurrencyNew")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData currency.CurrencyCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	id := currency.GetCurrencyObjectID(evt.GetAggregateID())
	span.SetTag(tracing.SpanTagEntityId, id)

	source := helper.GetSource(eventData.SourceFields.Source)
	appSource := helper.GetAppSource(eventData.SourceFields.AppSource)
	err := h.repositories.Neo4jRepositories.CurrencyWriteRepository.CreateCurrency(ctx, id, eventData.Name, eventData.Symbol, source, appSource, eventData.CreatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving currency %s: %s", id, err.Error())
		return err
	}
	return err
}
