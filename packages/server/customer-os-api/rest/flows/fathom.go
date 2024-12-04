package flows

import (
	"net/http"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
)

func FathomZapier(c *rest.HTTPContext) {
	ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.GinContext.Request.Context(), "Fathom", c.GinContext.Request.Header)
	c.ServiceContext = &ctx
	c.Span = span
	defer span.Finish()
	commontracing.TagComponentRest(span)

	tenant, err := c.Services.CommonServices.PostgresRepositories.TenantRepository.GetTenant(ctx, c.GinContext.Param("tenantId"))
	if err != nil {
		err := errors.Wrap(err, "Unable to identify tenant")
		tracing.TraceErr(span, err)
		rest.SendError(c.GinContext, c.Span, http.StatusUnauthorized, rest.ErrUnauthorized)
		return
	}
	c.Tenant = tenant

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
		if err := publishMeetingSummaryCreatedEvent(ctx, aiSummaryData); err != nil {
			tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to process Fathom AI summary from zapier"))
		}
	}()
	return
}

func publishMeetingSummaryCreatedEvent(ctx *rest.HTTPContext, aiSummaryData *FathomZapierPayload) error {
	var meetingSummary data_fields.MeetingSummaryEvent

	content, err := aiSummaryData.toMarkdownContent()
	if err != nil {
		tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to convert Fathom AI summary to markdown"))
		return err
	}

	// build meeting summary
	meetingSummary.Content = &content
	participants := aiSummaryData.Meeting.participantEmails()
	meetingSummary.ParticipantEmails = &participants

	if aiSummaryData.Meeting.ScheduledStartTime.IsZero() {
		meetingSummary.Timestamp = utils.NowPtr()
	} else {
		meetingSummary.Timestamp = utils.TimePtr(aiSummaryData.Meeting.ScheduledStartTime.UTC())
	}

	event, err := dto.NewWebhookEvent(enum.Fathom, "meeting_summary", "created", &meetingSummary)
	if err != nil {
		tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to build webhook event"))
		return err
	}

	err = ctx.Services.CommonServices.RabbitMQService.PublishWebhookEvent(*ctx.ServiceContext, event)
	if err != nil {
		tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to publish event"))
	}

	return err
}
