// @openapi 3.0.0
package events

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/net/html"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func FathomZapier(services *service.Services) gin.HandlerFunc {
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

		handleFathomAISummaryZapier(httpContext)
	}
}

func handleFathomAISummaryZapier(ctx rest.HTTPContext) {
	var aiSummaryData FathomZapierPayload
	err := ctx.GinContext.BindJSON(&aiSummaryData)
	if err != nil {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to parse payload from Zapier"))
		return
	}

	if aiSummaryData.AISummary.HTMLFormatted == "" {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("No Fantom summary data"))
		return
	}

	ctx.GinContext.JSON(http.StatusAccepted, rest.BuildBaseResponse(rest.StatusProcessing))

	go func() {
		if err := createEventFromFathomAISummaryZapier(ctx, &aiSummaryData); err != nil {
			tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to process Fathom AI summary from zapier"))
		}
	}()
	return
}

func createEventFromFathomAISummaryZapier(ctx rest.HTTPContext, aiSummaryData *FathomZapierPayload) error {
	var event data_fields.MarkdownEventFields
	var allErrs error

	content, err := processFathomSummaryFromZapier(aiSummaryData)
	if err != nil {
		return err
	}

	source := neo4jentity.DataSourceFathom
	event.Source = &source
	if aiSummaryData.Meeting.ScheduledStartTime.IsZero() {
		event.CreatedAt = utils.NowPtr()
	} else {
		event.CreatedAt = utils.TimePtr(aiSummaryData.Meeting.ScheduledStartTime.UTC())
	}
	event.Content = &content

	domains := aiSummaryData.ExternalDomains()
	orgIds, err := getParticipantOrganizationIds(ctx, domains)
	if err != nil {
		allErrs = multierr.Append(allErrs, errors.Wrap(err, "failed to get organization id for participant"))
		tracing.TraceErr(ctx.Span, err)
	}

	orgIds = utils.RemoveEmpties(orgIds)
	orgIds = utils.RemoveDuplicates(orgIds)

	for _, org := range orgIds {
		event.OrganizationId = &org
		_, err := ctx.Services.CommonServices.MarkdownEventService.Save(*ctx.ServiceContext, nil, nil, event)
		if err != nil {
			allErrs = multierr.Append(allErrs, errors.Wrap(err, "failed to save markdown event"))
			tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to save markdown event"))
		}
	}

	return allErrs
}

func processFathomSummaryFromZapier(raw *FathomZapierPayload) (string, error) {
	// Convert HTML to clean markdown
	cleanMarkdown, err := convertFathomHTMLToCleanMarkdown(raw.AISummary.HTMLFormatted)
	if err != nil {
		return "", fmt.Errorf("error converting HTML to markdown: %w", err)
	}

	// Build the additional sections
	var builder strings.Builder
	builder.WriteString(cleanMarkdown)

	// Add participants section
	builder.WriteString("\n\n### Meeting Participants\n")
	for _, participant := range raw.Meeting.Invitees {
		builder.WriteString(fmt.Sprintf("- %s (%s)\n", participant.Name, participant.Email))
	}

	// Add duration
	builder.WriteString("\n### Meeting Duration\n")
	mins, err := strconv.ParseFloat(raw.Recording.DurationInMinutes, 64)
	if err == nil {
		builder.WriteString(fmt.Sprintf("%d minutes\n\n", int(mins)))
	}

	// Add recording link
	builder.WriteString(fmt.Sprintf("[View Recording](%s)\n", raw.Recording.ShareURL))

	return builder.String(), nil
}

func convertFathomHTMLToCleanMarkdown(htmlContent string) (string, error) {
	// Parse HTML
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	var process func(*html.Node)

	process = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			builder.WriteString(n.Data)
		case html.ElementNode:
			switch n.Data {
			case "h1":
				builder.WriteString("\n# ")
			case "h2":
				builder.WriteString("\n## ")
			case "h3":
				builder.WriteString("\n### ")
			case "p":
				builder.WriteString("\n\n")
			case "ul":
				builder.WriteString("\n")
			case "li":
				builder.WriteString("\n- ")
			case "br":
				builder.WriteString("\n")
			case "a":
				// Skip the href and just process the text content
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				process(c)
			}

			switch n.Data {
			case "p", "h1", "h2", "h3", "ul", "li":
				builder.WriteString("\n")
			}
		}
	}

	process(doc)

	// Clean up the output
	output := builder.String()

	// Remove multiple newlines
	output = regexp.MustCompile(`\n{3,}`).ReplaceAllString(output, "\n\n")
	// Remove leading/trailing whitespace
	output = strings.TrimSpace(output)
	// Ensure consistent newlines between sections
	output = regexp.MustCompile(`\n## `).ReplaceAllString(output, "\n\n## ")
	output = regexp.MustCompile(`\n### `).ReplaceAllString(output, "\n\n### ")

	return output, nil
}
