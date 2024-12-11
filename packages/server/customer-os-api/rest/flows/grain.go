package flows

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func GrainZapier(c *gin.Context, s *service.Services) {
	ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.Grain", c.Request.Header)
	defer span.Finish()
	commontracing.TagComponentRest(span)

	tenant, err := s.CommonServices.PostgresRepositories.TenantRepository.GetTenant(ctx, c.Param("tenantId"))
	if err != nil {
		err := errors.Wrap(err, "Unable to identify tenant")
		tracing.TraceErr(span, err)
		rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
		return
	}

	// update context with tenant, pass this where tenant is needed
	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    tenant,
		AppSource: constants.AppSourceCustomerOsApiRest,
	})

	if !strings.HasPrefix(c.ContentType(), "application/json") {
		rest.SendError(c, span, http.StatusBadRequest, rest.ErrUnsupportedContentType)
	}

	if c.Request.UserAgent() == "" {
		rest.SendError(c, span, http.StatusForbidden, rest.ErrForbidden)
	}

	if !strings.EqualFold(c.Request.UserAgent(), "Zapier") {
		rest.SendError(c, span, http.StatusForbidden, rest.ErrForbidden)
	}

	handleGrainNewRecordingEventZapier(c, ctx, s)
}

func handleGrainNewRecordingEventZapier(c *gin.Context, ctx context.Context, s *service.Services) {
	span, _ := commontracing.StartTracerSpan(c.Request.Context(), "Flows.handleGrainNewRecorderEventZapier")
	defer span.Finish()
	commontracing.TagComponentRest(span)

	var grainDataPayload GrainRecordingData
	err := c.BindJSON(&grainDataPayload)
	if err != nil {
		rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to parse payload from Zapier"))
		return
	}

	grainData := &grainDataPayload
	err = grainData.cleanPayload()
	if err != nil {
		rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to normalize payload from Zapier"))
		return
	}

	if grainData.RecordingData.IntelligenceNotesMD == "" {
		rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("No Grain meeting in payload"))
	}

	c.JSON(http.StatusAccepted, rest.BuildBaseResponse(rest.StatusProcessing))

	go func() {
		if err := publishGrainMeetingSummaryCreatedEvent(c, ctx, s, grainData); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to process Grain AI summary from zapier"))
		}
	}()
	return
}

func publishGrainMeetingSummaryCreatedEvent(c *gin.Context, ctx context.Context, s *service.Services, grainData *GrainRecordingData) error {
	span, _ := commontracing.StartTracerSpan(c.Request.Context(), "Flows.publishGrainMeetingSummaryEvent")
	defer span.Finish()
	commontracing.TagComponentRest(span)

	var meeting data_fields.MeetingSummaryEvent

	content := grainData.MeetingNoteContent()
	meeting.Content = &content
	meeting.MeetingID = grainData.RecordingData.ID
	participants := grainData.RecordingData.participantEmails()
	meeting.ParticipantEmails = &participants

	if grainData.RecordingData.StartDatetime.IsZero() {
		meeting.Timestamp = utils.NowPtr()
	} else {
		meeting.Timestamp = utils.TimePtr(grainData.RecordingData.StartDatetime.UTC())
	}

	// build webhook event
	event := dto.WebhookEvent{
		ExternalSystemId: enum.Grain,
		Name:             commonenum.EventGrainMeetingSummaryCreated,
		DataType:         "MeetingSummaryEvent",
		Data:             &meeting,
	}

	pubErr := s.CommonServices.RabbitMQService.PublishWebhookEvent(ctx, event)

	if pubErr != nil {
		tracing.TraceErr(span, errors.Wrap(pubErr, "failed to publish event"))
	}

	return nil
}
