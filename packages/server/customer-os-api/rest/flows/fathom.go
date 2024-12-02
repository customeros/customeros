package flows

import (
	"net/http"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
)

func FathomZapier(c *rest.HTTPContext) {
	ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.GinContext.Request.Context(), "Fathom", c.GinContext.Request.Header)
	c.ServiceContext = &ctx
	c.Span = span
	defer span.Finish()
	commontracing.TagComponentRest(span)

	c.Tenant = rest.ValidateTenant(c.GinContext, *c.ServiceContext, c.Span)
	if c.Tenant == "" {
		rest.SendError(c.GinContext, c.Span, http.StatusForbidden, rest.ErrForbidden)
		return
	}

	if !strings.HasPrefix(c.GinContext.ContentType(), "application/json") {
		rest.SendError(c.GinContext, c.Span, http.StatusBadRequest, rest.ErrUnsupportedContentType)
		return
	}

	if c.GinContext.Request.UserAgent() == "" {
		rest.SendError(c.GinContext, c.Span, http.StatusForbidden, rest.ErrForbidden)
		return
	}

	if !strings.EqualFold(c.GinContext.Request.UserAgent(), "Zapier") {
		rest.SendError(c.GinContext, c.Span, http.StatusForbidden, rest.ErrForbidden)
		return
	}

	handleFathomAISummaryZapier(c)
}

func handleFathomAISummaryZapier(ctx *rest.HTTPContext) {
	var aiSummaryDataPayload FathomZapierPayload
	err := ctx.GinContext.BindJSON(&aiSummaryDataPayload)
	if err != nil {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to parse payload from Zapier"))
		return
	}

	aiSummaryData := &aiSummaryDataPayload
	err = aiSummaryData.toCleanPayload()
	if err != nil {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to normalize payload from Zapier"))
		return
	}

	if aiSummaryData.AISummary.HTMLFormatted == "" {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("No Fathom summary data"))
		return
	}

	ctx.GinContext.JSON(http.StatusAccepted, rest.BuildBaseResponse(rest.StatusProcessing))

	go func() {
		if err := createEventFromFathomAISummaryZapier(ctx, aiSummaryData); err != nil {
			tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to process Fathom AI summary from zapier"))
		}
	}()
	return
}

func createEventFromFathomAISummaryZapier(ctx *rest.HTTPContext, aiSummaryData *FathomZapierPayload) error {
	var meetingSummary data_fields.MeetingSummaryCreatedEvent
	var event data_fields.FlowEventObject[data_fields.MeetingSummaryCreatedEvent]
	var allErrs error

	content, err := aiSummaryData.toMarkdownContent()
	if err != nil {
		return err
	}

	// build meeting summary
	meetingSummary.Content = &content
	participants := aiSummaryData.Meeting.participantEmails()
	meetingSummary.ParticipantEmails = &participants

	// build event object
	event.Tenant = &ctx.Tenant
	event.Event = data_fields.FlowEventFathomMeetingSummaryCreated
	event.Integration = data_fields.IntegrationFathom
	event.Payload = &meetingSummary

	if aiSummaryData.Meeting.ScheduledStartTime.IsZero() {
		event.Timestamp = utils.NowPtr()
	} else {
		event.Timestamp = utils.TimePtr(aiSummaryData.Meeting.ScheduledStartTime.UTC())
	}

	ctx.Services.CommonServices.RabbitMQService.PublishEvent(
		*ctx.ServiceContext, id, model.FLOW_EVENT, event,
	)

	return allErrs
}
