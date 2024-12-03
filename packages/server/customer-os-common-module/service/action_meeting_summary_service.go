package service

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type ActionMeetingSummaryService interface {
	HandleMeetingSummaryEvent(ctx context.Context, eventData *data_fields.MeetingSummaryEvent, source enum.ExternalSystemId) error
	PublishEventCreateTimelineEvent(ctx context.Context, mdEvent *data_fields.MarkdownEventFields) error
	PublishEventCreateContact(ctx context.Context) error
	PublishEventCreateOrganizations(ctx context.Context) error
}

type actionMeetingSummaryService struct {
	services *Services
}

func NewActionMeetingSummaryService(services *Services) ActionMeetingSummaryService {
	return &actionMeetingSummaryService{
		services: services,
	}
}

func (a *actionMeetingSummaryService) HandleMeetingSummaryEvent(ctx context.Context, eventData *data_fields.MeetingSummaryEvent, source enum.ExternalSystemId) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionMeetingSummaryService.HandleMeetingSummaryEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	mdEvent, err := a.buildMarkdownEvent(ctx, span, eventData, source)
	if err != nil {
		err = fmt.Errorf("failed to handle meeting summary event: %v", err)
		tracing.TraceErr(span, err)
		return err
	}

	err = a.services.RabbitMQService.PublishEvent(ctx, "", model.MARKDOWN_EVENT, mdEvent)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to publish markdown event"))
		return err
	}

	return nil
}

func (a *actionMeetingSummaryService) PublishEventCreateTimelineEvent(ctx context.Context, mdEvent *data_fields.MarkdownEventFields) error {
	return nil
}

func (a *actionMeetingSummaryService) PublishEventCreateContact(ctx context.Context) error {
	return nil
}

func (a *actionMeetingSummaryService) PublishEventCreateOrganizations(ctx context.Context) error {
	return nil
}

func (a *actionMeetingSummaryService) buildMarkdownEvent(ctx context.Context, span opentracing.Span, event *data_fields.MeetingSummaryEvent, source enum.ExternalSystemId) (*data_fields.MarkdownEventFields, error) {
	var mdEvent data_fields.MarkdownEventFields
	var sourceId entity.DataSource

	switch source {
	case enum.Fathom:
		sourceId = entity.DataSourceFathom
	case enum.Grain:
		sourceId = entity.DataSourceGrain
	default:
		return nil, fmt.Errorf("Unuspported source: %v", source)
	}

	mdEvent.Source = &sourceId
	mdEvent.Content = event.Content
	mdEvent.CreatedAt = event.Timestamp

	return &mdEvent, nil
}
