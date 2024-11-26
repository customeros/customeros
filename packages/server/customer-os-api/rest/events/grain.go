// @openapi 3.0.0
package events

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"net/http"
	"regexp"
	"strings"

	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

func GrainZapier(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Grain", c.Request.Header)
		defer span.Finish()
		commontracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			return
		}

		if !strings.HasPrefix(c.ContentType(), "application/json") {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrUnsupportedContentType)
		}

		httpContext := rest.HTTPContext{
			GinContext:     c,
			ServiceContext: &ctx,
			Span:           span,
			Services:       services,
			Tenant:         tenant,
		}

		if c.Request.UserAgent() == "" {
			rest.SendError(c, span, http.StatusForbidden, rest.ErrForbidden)
		}

		if !strings.EqualFold(c.Request.UserAgent(), "Zapier") {
			rest.SendError(c, span, http.StatusForbidden, rest.ErrForbidden)
		}

		handleGrainNewRecordingEventZapier(httpContext)
	}
}

func handleGrainNewRecordingEventZapier(ctx rest.HTTPContext) {
	var grainData GrainRecordingData
	err := ctx.GinContext.BindJSON(&grainData)
	if err != nil {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to parse payload from Zapier"))
		return
	}

	err = cleanJsonPayload(&grainData)
	if err != nil {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to normalize payload from Zapier"))
		return
	}

	if grainData.Data.IntelligenceNotesMD == "" {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("No Grain meeting in payload"))
	}

	ctx.GinContext.JSON(http.StatusAccepted, rest.BuildBaseResponse(rest.StatusProcessing))

	go func() {
		if err := createEventFromGrainRecordingZapier(ctx, &grainData); err != nil {
			tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to process Grain AI summary from zapier"))
		}
	}()
	return
}

func cleanJsonPayload(data *GrainRecordingData) error {
	ownersJson := utils.ReplaceSingleQuotesWithDoubleQuotes(data.Data.OwnersStr)
	var owners []string
	err := json.Unmarshal([]byte(ownersJson), &owners)
	if err != nil {
		return err
	}
	data.Data.Owners = owners

	participantsJson := utils.ReplaceSingleQuotesWithDoubleQuotes(data.Data.ParticipantsStr)

	participantsJson = strings.Replace(participantsJson, " True", " true", -1)
	participantsJson = strings.Replace(participantsJson, " False", " false", -1)
	participantsJson = strings.Replace(participantsJson, " None", " \"\"", -1)

	var participants []GrainParticipant
	err = json.Unmarshal([]byte(participantsJson), &participants)
	if err != nil {
		return err
	}
	data.Data.Participants = participants

	tagsJson := utils.ReplaceSingleQuotesWithDoubleQuotes(data.Data.TagsStr)
	var tags []string
	err = json.Unmarshal([]byte(tagsJson), &tags)
	if err != nil {
		return err
	}

	return nil
}

func createEventFromGrainRecordingZapier(ctx rest.HTTPContext, grainData *GrainRecordingData) error {
	var event data_fields.MarkdownEventFields
	var allErrs error

	content := extractGrainMeetingNotes(grainData)
	event.Content = &content
	if grainData.Data.StartDatetime.IsZero() {
		event.CreatedAt = utils.NowPtr()
	} else {
		event.CreatedAt = utils.TimePtr(grainData.Data.StartDatetime.UTC())
	}
	source := neo4jentity.DataSourceGrain
	event.Source = &source

	orgIds, err := getParticipantOrganizationIds(ctx, getParticipantDomains(&grainData.Data.Participants))
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

func getParticipantDomains(participants *[]GrainParticipant) []string {
	var domains []string

	for _, participant := range *participants {
		emailData := mailvalidate.ValidateEmailSyntax(*participant.Email)
		if !emailData.IsValid {
			continue
		}
		domains = append(domains, emailData.Domain)
	}
	return domains
}

func extractGrainMeetingNotes(grainData *GrainRecordingData) string {
	if grainData == nil || grainData.Data.IntelligenceNotesMD == "" {
		return ""
	}

	// Remove hyperlinks using regex
	// Matches markdown links like [(time)](url)
	linkPattern := regexp.MustCompile(`\[\([^)]+\)\]\([^)]+\)`)
	cleanNotes := linkPattern.ReplaceAllStringFunc(grainData.Data.IntelligenceNotesMD, func(match string) string {
		// Extract just the text between the brackets
		textStart := strings.Index(match, "(") + 1
		textEnd := strings.Index(match, ")")
		return match[textStart:textEnd]
	})

	// Calculate meeting duration
	duration := grainData.Data.EndDatetime.Sub(grainData.Data.StartDatetime)
	durationMinutes := int(duration.Minutes())

	// Build participants section
	var attendees, declines []string
	for _, p := range grainData.Data.Participants {
		email := "No email"
		if p.Email != nil {
			email = *p.Email
		}
		var participant string
		if email == "No email" {

			participant = p.Name
		} else {
			participant = fmt.Sprintf("- %s (%s)", p.Name, email)
		}
		if p.ConfirmedAttendee {
			attendees = append(attendees, participant)
		} else {
			declines = append(declines, participant)
		}
	}

	// Build the final markdown
	var sb strings.Builder

	// Add meeting details
	sb.WriteString(fmt.Sprintf("# %s\n\n", grainData.Data.Title))
	sb.WriteString(fmt.Sprintf("**Meeting Date:** %s\n", grainData.Data.StartDatetime.Format("January 2, 2006 3:04 PM MST")))
	sb.WriteString(fmt.Sprintf("**Duration:** %d minutes\n\n", durationMinutes))

	// Add participants
	sb.WriteString("**Attended**\n")
	sb.WriteString(strings.Join(attendees, "\n"))
	sb.WriteString("\n\n")
	if len(declines) > 0 {
		sb.WriteString("**Declined/No Show**\n")
		sb.WriteString(strings.Join(declines, "\n"))
		sb.WriteString("\n\n")
	}

	// Add meeting notes
	sb.WriteString("## Meeting Notes\n")
	sb.WriteString(cleanNotes)
	sb.WriteString("\n\n")

	// Add recording link
	sb.WriteString(fmt.Sprintf("[View Recording](%s)\n", grainData.Data.PublicURL))

	return sb.String()
}
