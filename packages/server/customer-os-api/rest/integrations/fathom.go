package integrations

import (
	"context"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	commontracing "github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
)

func (h *IntegrationHandler) FathomZapier(c *gin.Context, tenant string) {
	ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.FathomZapier", c.Request.Header)
	defer span.Finish()
	commontracing.TagComponentRest(span)
	tracing.TagTenant(span, tenant)

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

	if !strings.EqualFold(c.Request.UserAgent(), "Zapier") {
		h.responseHandler.HandleError(c, http.StatusForbidden, nil)
		return
	}

	h.handleFathomAISummaryZapier(c, ctx)
}

func (h *IntegrationHandler) handleFathomAISummaryZapier(c *gin.Context, ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Integrations.handleFathomAISummaryZapier")
	defer span.Finish()
	commontracing.TagComponentRest(span)

	var aiSummaryDataPayload FathomZapierPayload
	err := c.BindJSON(&aiSummaryDataPayload)
	if err != nil {
		message := "Unable to parse payload from Zapier"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return
	}

	userEmailValidation := mailvalidate.ValidateEmailSyntax(aiSummaryDataPayload.FathomUser.Email)
	email := userEmailValidation.CleanEmail
	if email == "" {
		message := "Cannot find user"
		h.responseHandler.HandleError(c, http.StatusNotFound, &message)
		return
	}
	userId, err := h.GetCustomerOSUser(ctx, []string{email})
	if err != nil {
		message := "User not found"
		h.responseHandler.HandleError(c, http.StatusNotFound, &message)
		return
	}
	if userId == "" {
		h.responseHandler.HandleError(c, http.StatusUnauthorized, nil)
		return
	}
	span.LogKV("userId", userId)
	ctx = common.SetUserIdInContext(ctx, userId)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "Integrations.publishFathomMeetingSummaryCreatedEvent")
	defer span.Finish()
	commontracing.TagComponentRest(span)

	meetingID := aiSummaryData.ID
	participants := aiSummaryData.Meeting.participantEmails()
	content, err := aiSummaryData.toMarkdownContent()
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to convert Fathom AI summary to markdown"))
		return err
	}

	event := dto.NewMeetingRecording{
		MeetingTitle:        aiSummaryData.Meeting.Title,
		Source:              enum.SourceFathom,
		Content:             content,
		ParticipantEmails:   participants,
		MeetingRecordingUrl: aiSummaryData.Recording.ShareURL,
	}
	if !utils.IsStringInSlice(aiSummaryData.FathomUser.Email, event.ParticipantEmails) {
		event.ParticipantEmails = append(event.ParticipantEmails, aiSummaryData.FathomUser.Email)
	}

	if aiSummaryData.Meeting.ScheduledStartTime.IsZero() {
		event.Timestamp = utils.Now()
	} else {
		event.Timestamp = aiSummaryData.Meeting.ScheduledStartTime.UTC()
	}

	pubErr := h.services.CommonServices.Events.Publisher.PublishFanoutEvent(ctx, meetingID, model.MEETING, event)
	if pubErr != nil {
		tracing.TraceErr(span, errors.Wrap(pubErr, "failed to publish event"))
	}

	return nil
}
