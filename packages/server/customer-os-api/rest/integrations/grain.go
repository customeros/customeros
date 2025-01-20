package integrations

import (
	"context"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	commontracing "github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
)

func (h *IntegrationHandler) GrainZapier(c *gin.Context) {
	ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.Grain", c.Request.Header)
	defer span.Finish()
	commontracing.TagComponentRest(span)

	tenant, err := h.services.Repositories.PostgresRepositories.TenantRepository.GetTenant(ctx, c.Param("tenantId"))
	if err != nil {
		err := errors.Wrap(err, "Unable to identify tenant")
		tracing.TraceErr(span, err)
		h.responseHandler.HandleError(c, http.StatusUnauthorized, nil)
		return
	}

	// update context with tenant, pass this where tenant is needed
	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    tenant,
		AppSource: constants.AppSourceCustomerOsApiRest,
	})

	if !strings.HasPrefix(c.ContentType(), "application/json") {
		message := "Unsupported Content-Type"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
	}

	if c.Request.UserAgent() == "" {
		h.responseHandler.HandleError(c, http.StatusForbidden, nil)
	}

	// if !strings.EqualFold(c.Request.UserAgent(), "Zapier") {
	// 	handlers.SendError(c, span, http.StatusForbidden, enum.ErrForbidden)
	// }

	h.handleGrainNewRecordingEventZapier(c, ctx)
}

func (h *IntegrationHandler) handleGrainNewRecordingEventZapier(c *gin.Context, ctx context.Context) {
	span, _ := commontracing.StartTracerSpan(c.Request.Context(), "Flows.handleGrainNewRecorderEventZapier")
	defer span.Finish()
	commontracing.TagComponentRest(span)

	var grainDataPayload GrainRecordingData
	err := c.BindJSON(&grainDataPayload)
	if err != nil {
		message := "Unable to parse payload from Zapier"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return
	}

	grainData := &grainDataPayload
	err = grainData.cleanPayload()
	if err != nil {
		message := "Unable to normalize payload from Zapier"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return
	}

	if grainData.RecordingData.IntelligenceNotesMD == "" {
		message := "No Grain meeting in payload"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
	}

	h.responseHandler.HandleAccepted(c)

	go func() {
		if err := h.publishGrainMeetingSummaryCreatedEvent(c, ctx, grainData); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to process Grain AI summary from zapier"))
		}
	}()
	return
}

func (h *IntegrationHandler) publishGrainMeetingSummaryCreatedEvent(c *gin.Context, ctx context.Context, grainData *GrainRecordingData) error {
	span, _ := commontracing.StartTracerSpan(c.Request.Context(), "Flows.publishGrainMeetingSummaryEvent")
	defer span.Finish()
	commontracing.TagComponentRest(span)

	var meeting data_fields.MeetingSummaryEvent

	content := grainData.MeetingNoteContent()
	meeting.Content = &content
	meeting.Tenant = common.GetTenantFromContext(ctx)
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
		ExternalSystemId: commonenum.SourceGrain,
		Name:             commonenum.EventGrainMeetingSummaryCreated,
		DataType:         "MeetingSummaryEvent",
		Data:             &meeting,
	}

	pubErr := h.services.CommonServices.Events.Publisher.PublishWebhookEvent(ctx, event)

	if pubErr != nil {
		tracing.TraceErr(span, errors.Wrap(pubErr, "failed to publish event"))
	}

	return nil
}
