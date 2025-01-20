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

func (h *IntegrationHandler) FathomZapier(c *gin.Context) {
	ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.FathomZapier", c.Request.Header)
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
		return
	}

	if c.Request.UserAgent() == "" {
		h.responseHandler.HandleError(c, http.StatusForbidden, nil)
		return
	}

	// if !strings.EqualFold(c.Request.UserAgent(), "Zapier") {
	// 	handlers.SendError(c, span, http.StatusForbidden, enum.ErrForbidden)
	// 	return
	// }

	h.handleFathomAISummaryZapier(c, ctx)
}

func (h *IntegrationHandler) handleFathomAISummaryZapier(c *gin.Context, ctx context.Context) {
	span, _ := commontracing.StartTracerSpan(c.Request.Context(), "Flows.handleFathomAISummaryZapier")
	defer span.Finish()
	commontracing.TagComponentRest(span)

	var aiSummaryDataPayload FathomZapierPayload
	err := c.BindJSON(&aiSummaryDataPayload)
	if err != nil {
		message := "Unable to parse payload from Zapier"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return
	}

	aiSummaryData := &aiSummaryDataPayload
	err = aiSummaryData.toCleanPayload()
	if err != nil {
		message := "Unable to normalize payload from Zapier"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return
	}

	if aiSummaryData.AISummary.HTMLFormatted == "" {
		message := "No Fathom data"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return
	}

	h.responseHandler.HandleAccepted(c)

	go func() {
		if err := h.publishFathomMeetingSummaryCreatedEvent(c, ctx, aiSummaryData); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to process Fathom AI summary from zapier"))
		}
	}()
	return
}

func (h *IntegrationHandler) publishFathomMeetingSummaryCreatedEvent(c *gin.Context, ctx context.Context, aiSummaryData *FathomZapierPayload) error {
	span, _ := commontracing.StartTracerSpan(c.Request.Context(), "Flows.publishFathomMeetingSummaryCreatedEvent")
	defer span.Finish()
	commontracing.TagComponentRest(span)

	var meetingSummary data_fields.MeetingSummaryEvent

	content, err := aiSummaryData.toMarkdownContent()
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to convert Fathom AI summary to markdown"))
		return err
	}

	// build meeting summary
	meetingSummary.Tenant = common.GetTenantFromContext(ctx)
	meetingSummary.MeetingID = aiSummaryData.ID
	meetingSummary.Content = &content
	participants := aiSummaryData.Meeting.participantEmails()
	meetingSummary.ParticipantEmails = &participants

	if aiSummaryData.Meeting.ScheduledStartTime.IsZero() {
		meetingSummary.Timestamp = utils.NowPtr()
	} else {
		meetingSummary.Timestamp = utils.TimePtr(aiSummaryData.Meeting.ScheduledStartTime.UTC())
	}

	event := dto.WebhookEvent{
		ExternalSystemId: commonenum.SourceFathom,
		Name:             commonenum.EventFathomMeetingSummaryCreated,
		DataType:         data_fields.MeetingSummaryEvent{}.Type(),
		Data:             meetingSummary,
	}

	pubErr := h.services.CommonServices.Events.Publisher.PublishWebhookEvent(ctx, event)
	if pubErr != nil {
		tracing.TraceErr(span, errors.Wrap(pubErr, "failed to publish event"))
	}

	return nil
}
