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

func GrainZapier(c *rest.HTTPContext) {
	ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.GinContext.Request.Context(), "Grain", c.GinContext.Request.Header)
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
	}

	if c.GinContext.Request.UserAgent() == "" {
		rest.SendError(c.GinContext, c.Span, http.StatusForbidden, rest.ErrForbidden)
	}

	if !strings.EqualFold(c.GinContext.Request.UserAgent(), "Zapier") {
		rest.SendError(c.GinContext, c.Span, http.StatusForbidden, rest.ErrForbidden)
	}

	handleGrainNewRecordingEventZapier(c)
}

func handleGrainNewRecordingEventZapier(ctx *rest.HTTPContext) {
	var grainDataPayload GrainRecordingData
	err := ctx.GinContext.BindJSON(&grainDataPayload)
	if err != nil {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to parse payload from Zapier"))
		return
	}

	grainData := &grainDataPayload
	err = grainData.cleanPayload()
	if err != nil {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to normalize payload from Zapier"))
		return
	}

	if grainData.RecordingData.IntelligenceNotesMD == "" {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("No Grain meeting in payload"))
	}

	ctx.GinContext.JSON(http.StatusAccepted, rest.BuildBaseResponse(rest.StatusProcessing))

	go func() {
		if err := publishGrainMeetingSummaryCreatedEvent(ctx, grainData); err != nil {
			tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to process Grain AI summary from zapier"))
		}
	}()
	return
}

func publishGrainMeetingSummaryCreatedEvent(ctx *rest.HTTPContext, grainData *GrainRecordingData) error {
	var meeting data_fields.MeetingSummaryEvent

	content := grainData.MeetingNoteContent()
	meeting.Content = &content
	participants := grainData.RecordingData.participantEmails()
	meeting.ParticipantEmails = &participants

	if grainData.RecordingData.StartDatetime.IsZero() {
		meeting.Timestamp = utils.NowPtr()
	} else {
		meeting.Timestamp = utils.TimePtr(grainData.RecordingData.StartDatetime.UTC())
	}

	event, err := dto.NewWebhookEvent(enum.Grain, "meeting_summary", "created", &meeting)
	if err != nil {
		tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to build webhook event"))
		return err
	}

	pubErr := ctx.Services.CommonServices.RabbitMQService.PublishWebhookEvent(*ctx.ServiceContext, event)

	if pubErr != nil {
		tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to publish event"))
	}

	return nil
}
