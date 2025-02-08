package integrations

import (
	"context"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	commontracing "github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/mailsherpa/mailvalidate"
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

	if !strings.EqualFold(c.Request.UserAgent(), "Zapier") {
		h.responseHandler.HandleError(c, http.StatusForbidden, nil)
	}

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
		return
	}

	userEmail, err := h.getCustomerOSUser(ctx, grainData.RecordingData.Owners)
	if err != nil {
		message := "User not found"
		h.responseHandler.HandleError(c, http.StatusNotFound, &message)
		return
	}
	ctx = common.SetUserEmailInContext(ctx, userEmail)

	h.responseHandler.HandleAccepted(c)

	go func() {
		if err := h.publishGrainMeetingSummaryCreatedEvent(c, ctx, grainData); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to process Grain AI summary from zapier"))
		}
	}()
	return
}

func (h *IntegrationHandler) getCustomerOSUser(ctx context.Context, meetingOwners []string) (string, error) {
	span, _ := commontracing.StartTracerSpan(ctx, "Flows.getCustomerOSUser")
	defer span.Finish()
	commontracing.TagComponentRest(span)

	if len(meetingOwners) == 0 {
		return "", coserrors.ErrCannotIdentifyUser
	}

	for _, owner := range meetingOwners {
		validation := mailvalidate.ValidateEmailSyntax(owner)
		email := validation.CleanEmail
		if email == "" {
			continue
		}
		user, _ := h.services.CommonServices.UserService.FindUserByEmail(ctx, email)
		if user != nil && user.Id != "" {
			return owner, nil
		}
	}
	return "", coserrors.ErrCannotIdentifyUser
}

func (h *IntegrationHandler) publishGrainMeetingSummaryCreatedEvent(c *gin.Context, ctx context.Context, grainData *GrainRecordingData) error {
	span, _ := commontracing.StartTracerSpan(c.Request.Context(), "Flows.publishGrainMeetingSummaryEvent")
	defer span.Finish()
	commontracing.TagComponentRest(span)

	content := grainData.MeetingNoteContent()
	participants := grainData.RecordingData.participantEmails()
	meetingID := grainData.RecordingData.ID

	event := dto.NewMeetingRecording{
		MeetingTitle:        grainData.RecordingData.Title,
		Source:              enum.SourceGrain,
		Content:             &content,
		ParticipantEmails:   &participants,
		MeetingRecordingUrl: grainData.RecordingData.PublicURL,
	}

	if grainData.RecordingData.StartDatetime.IsZero() {
		event.Timestamp = utils.NowPtr()
	} else {
		event.Timestamp = utils.TimePtr(grainData.RecordingData.StartDatetime.UTC())
	}

	pubErr := h.services.CommonServices.Events.Publisher.PublishFanoutEvent(ctx, meetingID, model.MEETING, event)

	if pubErr != nil {
		tracing.TraceErr(span, errors.Wrap(pubErr, "failed to publish new meeting recording event"))
	}

	return nil
}
