package flows

import (
	"net/http"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"
)

func FathomZapier(c *rest.HTTPContext) {
	ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.GinContext.Request.Context(), "Fathom", c.GinContext.Request.Header)
	c.ServiceContext = &ctx
	c.Span = span
	defer span.Finish()
	commontracing.TagComponentRest(span)
    

    //cannot do this as it's an unauthenticated endpoint
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

    // create Fathom Events in DB if they don't already exist
    err := createFathomFlowEvents()
    if err != nil {
        tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to setup Fathom flow events in db"))
        return
    }


	go func() {
		if err := createEventFromFathomAISummaryZapier(ctx, aiSummaryData); err != nil {
			tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to process Fathom AI summary from zapier"))
		}
	}()
	return
}

func createEventFromFathomAISummaryZapier(ctx *rest.HTTPContext, aiSummaryData *FathomZapierPayload) error {

	var meetingSummary data_fields.MeetingSummaryFields

	content, err := aiSummaryData.toMarkdownContent()
	if err != nil {
		tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to convert Fathom AI summary to markdown"))
		return err
	}

	// build meeting summary
	meetingSummary.Content = &content
	participants := aiSummaryData.Meeting.participantEmails()
	meetingSummary.ParticipantEmails = &participants
	meetingSummary.Source = neo4jenum.

	if aiSummaryData.Meeting.ScheduledStartTime.IsZero() {
		meetingSummary.Timestamp = utils.NowPtr()
	} else {
		meetingSummary.Timestamp = utils.TimePtr(aiSummaryData.Meeting.ScheduledStartTime.UTC())
	}

	err = ctx.Services.CommonServices.RabbitMQService.PublishEvent(*ctx.ServiceContext, "", model.FLOW_EVENT, dto.MeetingSummaryCreated{meetingSummary})
	if err != nil {
		tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to publish event"))
	}

	return err
}

func createFathomFlowEvents(ctx *rest.HTTPContext) error {
    events, err := ctx.Services.Repositories.PostgresRepositories.FlowEventsRepository.GetFlowEventsByExternalSystem(enum.Fathom)
    
    if err != nil {
        return err
    }

    for _, event := range events {
        if event.EventName != "fathom.meeting_summary.created" 
    }

    if "fathom.meeting_summary.created" 
     
    return nil
}

func event
