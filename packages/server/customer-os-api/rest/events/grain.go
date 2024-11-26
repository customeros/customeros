// @openapi 3.0.0
package events

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

func GrainZapier(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Fathom", c.Request.Header)
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
	if err != nil && grainData.Data.IntelligenceNotesMD != "" {
		ctx.GinContext.JSON(http.StatusAccepted, rest.BuildBaseResponse(rest.StatusProcessing))

		go func() {
			if err := createEventFromGrainRecordingZapier(ctx, &grainData); err != nil {
				tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to process Grain AI summary from zapier"))
			}
		}()
		return
	}

	return
}

func createEventFromGrainRecordingZapier(ctx rest.HTTPContext, grainData *GrainRecordingData) error {
	var meeting Event

	meeting.Content = extractGrainMeetingNotes(grainData)
	meeting.EventTimestamp = grainData.Data.StartDatetime
	meeting.Tenant = ctx.Tenant
	meeting.Source = string(EventGrain)

	orgIds, err := getParticipantOrganizationIds(ctx, getParticipantDomains(ctx, &grainData.Data.Participants))
	if err != nil {
		tracing.TraceErr(ctx.Span, errors.Wrap(err, "Error getting participant organization Ids"))
	}
	meeting.OrganizationIDs = orgIds

	// write event to db
	return nil
}

func getParticipantDomains(ctx rest.HTTPContext, participants *[]GrainParticipant) []string {
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
