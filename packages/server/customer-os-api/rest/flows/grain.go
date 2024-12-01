package flows

import (
	"net/http"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
)

func GrainZapier(c *rest.HTTPContext) {
	ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.GinContext.Request.Context(), "Grain", c.GinContext.Request.Header)
	c.ServiceContext = &ctx
	c.Span = span
	defer span.Finish()
	commontracing.TagComponentRest(span)

	tenant := rest.ValidateTenant(c.GinContext, *c.ServiceContext, c.Span)
	if tenant == "" {
		rest.SendError(c.GinContext, c.Span, http.StatusForbidden, rest.ErrForbidden)
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
		if err := createEventFromGrainRecordingZapier(ctx, grainData); err != nil {
			tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to process Grain AI summary from zapier"))
		}
	}()
	return
}

func createEventFromGrainRecordingZapier(ctx *rest.HTTPContext, grainData *GrainRecordingData) error {
	var event data_fields.MarkdownEventFields
	var allErrs error

	content := grainData.MeetingNoteContent()
	event.Content = &content
	if grainData.RecordingData.StartDatetime.IsZero() {
		event.CreatedAt = utils.NowPtr()
	} else {
		event.CreatedAt = utils.TimePtr(grainData.RecordingData.StartDatetime.UTC())
	}
	source := neo4jentity.DataSourceGrain
	event.Source = &source

	orgIds, err := getParticipantOrganizationIds(ctx, grainData.RecordingData.getDomains())
	if err != nil {
		allErrs = multierr.Append(allErrs, errors.Wrap(err, "failed to get organization id for participant"))
		tracing.TraceErr(ctx.Span, err)
	}

	orgIds = utils.RemoveDuplicates(orgIds)
	orgIds = utils.RemoveEmpties(orgIds)

	for _, orgId := range orgIds {
		event.OrganizationId = &orgId
		_, err := ctx.Services.CommonServices.MarkdownEventService.Save(*ctx.ServiceContext, nil, nil, event)
		if err != nil {
			allErrs = multierr.Append(allErrs, errors.Wrap(err, "failed to get organization id for participant"))
			tracing.TraceErr(ctx.Span, err)
		}
	}

	// write event to db
	return nil
}
