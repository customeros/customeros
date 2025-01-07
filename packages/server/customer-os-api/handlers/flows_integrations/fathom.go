package integrations

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
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func FathomZapier(c *gin.Context, s *service.Services) {
	ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.FathomZapier", c.Request.Header)
	defer span.Finish()
	commontracing.TagComponentRest(span)

	tenant, err := s.CommonServices.PostgresRepositories.TenantRepository.GetTenant(ctx, c.Param("tenantId"))
	if err != nil {
		err := errors.Wrap(err, "Unable to identify tenant")
		tracing.TraceErr(span, err)
		handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
		return
	}

	// update context with tenant, pass this where tenant is needed
	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    tenant,
		AppSource: constants.AppSourceCustomerOsApiRest,
	})

	if !strings.HasPrefix(c.ContentType(), "application/json") {
		handlers.SendError(c, span, http.StatusBadRequest, enum.ErrUnsupportedContentType)
		return
	}

	if c.Request.UserAgent() == "" {
		handlers.SendError(c, span, http.StatusForbidden, enum.ErrForbidden)
		return
	}

	// if !strings.EqualFold(c.Request.UserAgent(), "Zapier") {
	// 	handlers.SendError(c, span, http.StatusForbidden, enum.ErrForbidden)
	// 	return
	// }

	handleFathomAISummaryZapier(c, ctx, s)
}

func handleFathomAISummaryZapier(c *gin.Context, ctx context.Context, s *service.Services) {
	span, _ := commontracing.StartTracerSpan(c.Request.Context(), "Flows.handleFathomAISummaryZapier")
	defer span.Finish()
	commontracing.TagComponentRest(span)

	var aiSummaryDataPayload FathomZapierPayload
	err := c.BindJSON(&aiSummaryDataPayload)
	if err != nil {
		handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Unable to parse payload from Zapier"))
		return
	}

	aiSummaryData := &aiSummaryDataPayload
	err = aiSummaryData.toCleanPayload()
	if err != nil {
		handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Unable to normalize payload from Zapier"))
		return
	}

	if aiSummaryData.AISummary.HTMLFormatted == "" {
		handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("No Fathom summary data"))
		return
	}

	c.JSON(http.StatusAccepted, enum.BuildBaseResponse(enum.StatusProcessing))

	go func() {
		if err := publishFathomMeetingSummaryCreatedEvent(c, ctx, s, aiSummaryData); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to process Fathom AI summary from zapier"))
		}
	}()
	return
}

func publishFathomMeetingSummaryCreatedEvent(c *gin.Context, ctx context.Context, s *service.Services, aiSummaryData *FathomZapierPayload) error {
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

	pubErr := s.CommonServices.RabbitMQService.PublishWebhookEvent(ctx, event)
	if pubErr != nil {
		tracing.TraceErr(span, errors.Wrap(pubErr, "failed to publish event"))
	}

	return nil
}
