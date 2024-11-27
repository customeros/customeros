package webhooks

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/customeros/mailwatcher/dmarkstats"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

const EXTERNAL_SYSTEM = "mailstack"

func PostmarkInboundEmail(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "PostmarkInboundEmail", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		// Validate Postmark User-Agent
		if c.Request.UserAgent() == "" || !strings.EqualFold(c.Request.UserAgent(), "Postmark") {
			rest.SendError(c, span, http.StatusForbidden, rest.ErrForbidden)
			return
		}

		httpContext := rest.HTTPContext{
			GinContext:     c,
			ServiceContext: &ctx,
			Span:           span,
			Services:       s,
		}

		// Parse email data
		emailData, err := parseInboundEmail(httpContext)
		if err != nil {
			tracing.LogObjectAsJson(span, "body", c.Request.Body)
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest)
			return
		}

		// Return accepted response immediately
		c.JSON(http.StatusAccepted, rest.BuildBaseResponse(rest.StatusProcessing))

		// Process email asynchronously
		go func() {
			defer func() {
				if r := recover(); r != nil {
					stack := debug.Stack()
					err := fmt.Errorf("panic recovered in email processing: %v\n%s", r, stack)
					tracing.TraceErr(span, err)
				}
			}()

			var err error
			if emailData.IsMonitorEmail() {
				err = processDmarcMonitoringReport(httpContext, &emailData)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to process DMARC report"))
				}
			} else {
				err = processInboundEmail(httpContext, &emailData)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to process inbound email from Postmark"))
				}
			}
		}()
	}
}

func parseInboundEmail(ctx rest.HTTPContext) (PostmarkInboundEmailData, error) {
	var emailData PostmarkInboundEmailData
	err := ctx.GinContext.BindJSON(&emailData)
	if err != nil {
		return emailData, err
	}

	return emailData, nil
}

func processDmarcMonitoringReport(ctx rest.HTTPContext, emailData *PostmarkInboundEmailData) error {
	// Get attachment, unzip, feed file to dmark analyzer service
	attachment := emailData.Attachments[0]
	if attachment.ContentType != "application/zip" && attachment.ContentType != "application/gzip" {
		return fmt.Errorf("attachment %s is not a zip file", attachment.Name)
	}

	provider := emailData.DMARCReportProvider()

	reports, err := decodeAndReadDMARCReportFile(attachment.Content, attachment.ContentType)
	if err != nil {
		return fmt.Errorf("cannot parse dmarc report %s from attachment: %v", attachment.Name, err)
	}
	for _, report := range reports {
		dbReport := buildDMARCReport(ctx, report, provider)
		ctx.Services.Repositories.PostgresRepositories.MailStackDomainRepository.CreateDMARCReport(
			*ctx.ServiceContext, ctx.Tenant, &dbReport)
	}
	return nil
}

func buildDMARCReport(ctx rest.HTTPContext, report dmarcstats.Report, provider string) entity.DMARCMonitoring {
	// todo - need new query to do tenant lookup based on mailstack domain
	tenant := ""

	jsonReport, _ := json.Marshal(report)
	return entity.DMARCMonitoring{
		Tenant:        tenant,
		EmailProvider: provider,
		Domain:        report.Domain,
		ReportStart:   report.ReportPeriod.Start,
		ReportEnd:     report.ReportPeriod.End,
		MessageCount:  report.TotalMessages,
		SPFPass:       report.AuthResults.SPFPassCount,
		DKIMPass:      report.AuthResults.DKIMPassCount,
		DMARCPass:     report.AuthResults.DMARCPassCount,
		Data:          string(jsonReport),
	}
}

func processInboundEmail(ctx rest.HTTPContext, emailData *PostmarkInboundEmailData) error {
	tenant, err := getTenant(ctx, emailData)
	if err != nil {
		return err
	}
	ctx.Tenant = tenant

	participants := emailData.AllParticipantEmails()
	username, err := getUsername(ctx, participants)
	if err != nil || username == "" {
		ctx.Span.LogFields(tracingLog.Bool("mailbox.found", false))
		tracing.TraceErr(ctx.Span, err)
		return err
	}

	messageId := emailData.GetHeaderValue("Message-Id")
	emailExistsInDb, err := ctx.Services.CommonServices.PostgresRepositories.RawEmailRepository.EmailExistsByMessageId(
		*ctx.ServiceContext, EXTERNAL_SYSTEM, ctx.Tenant, username, messageId,
	)
	if err != nil {
		tracing.TraceErr(ctx.Span, err)
		return fmt.Errorf("Unable to determine if email exists in db for messageId %s: %v", messageId, err)
	}

	if emailExistsInDb {
		return nil
	}

	dbEmailEntity := emailData.ToRawDbObject()
	jsonEmailEntity, err := json.Marshal(dbEmailEntity)
	if err != nil {
		return fmt.Errorf("Unable to produce JSON email object for db: %v", err)
	}

	dbErr := ctx.Services.CommonServices.PostgresRepositories.RawEmailRepository.Store(*ctx.ServiceContext, EXTERNAL_SYSTEM, ctx.Tenant, username, dbEmailEntity.ProviderMessageId, messageId, string(jsonEmailEntity), dbEmailEntity.Sent, entity.REAL_TIME)
	if dbErr != nil {
		ctx.Span.LogFields(tracingLog.Object("raw_email", jsonEmailEntity))
		tracing.TraceErr(ctx.Span, err)
		return fmt.Errorf("Error writing email to db: %v", err)
	}

	// Check to see if email is a reply to a flow.  If so, mark as complete.
	// This should be handled in the email processor common service, not here.
	// Same with goal achieved.
	// Move slack notifications to processing service
	// Handle attachments

	return nil
}

func getTenant(ctx rest.HTTPContext, emailData *PostmarkInboundEmailData) (string, error) {
	nameFromBcc := emailData.TenantFromBcc()

	n, err := ctx.Services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantByNameIgnoreCase(*ctx.ServiceContext, nameFromBcc)
	if err != nil {
		ctx.Span.LogFields(tracingLog.Bool("tenant.found", false))
		tracing.TraceErr(ctx.Span, err)
		return "", fmt.Errorf("Cannot identify tenant %s: %v", nameFromBcc, err)
	}

	if n == nil {
		ctx.Span.LogFields(tracingLog.Bool("tenant.found", false))
		return "", fmt.Errorf("No valid tenant %s: %v", nameFromBcc, err)
	}

	tenant := mapper.MapDbNodeToTenantEntity(n)
	ctx.Span.LogFields(tracingLog.Bool("tenant.found", true))
	ctx.Span.LogFields(tracingLog.String("tenant.name", tenant.Name))

	return tenant.Name, nil
}

func getUsername(ctx rest.HTTPContext, EmailParticipants []string) (string, error) {
	for _, p := range EmailParticipants {
		userByEmail, err := ctx.Services.CommonServices.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(
			*ctx.ServiceContext, ctx.Tenant, p)
		if err != nil {
			return "", fmt.Errorf("Error getting username from email %s: %v", p, err)
		}
		if userByEmail != nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("Unable to find user amongst email participants")
}

func decodeAndReadDMARCReportFile(attachment, contentType string) ([]dmarcstats.Report, error) {
	var reports []dmarcstats.Report

	// Decode base64 string to bytes
	decoded, err := base64.StdEncoding.DecodeString(attachment)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	switch contentType {
	case "application/zip":
		reports, err = handleZipFile(decoded)
	case "application/gzip":
		reports, err = handleGzipFile(decoded)
	default:
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	if err != nil {
		return nil, err
	}

	return reports, nil
}

func handleZipFile(decoded []byte) ([]dmarcstats.Report, error) {
	var reports []dmarcstats.Report

	zipReader, err := zip.NewReader(bytes.NewReader(decoded), int64(len(decoded)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	for _, file := range zipReader.File {
		report, err := processZipFile(file)
		if err != nil {
			return nil, err
		}
		reports = append(reports, *report)
	}

	return reports, nil
}

func handleGzipFile(decoded []byte) ([]dmarcstats.Report, error) {
	var reports []dmarcstats.Report

	gzReader, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	report, err := dmarcstats.AnalyzeDMARCReport(gzReader)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze DMARC report from gzip: %w", err)
	}

	reports = append(reports, *report)
	return reports, nil
}

func processZipFile(file *zip.File) (*dmarcstats.Report, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file %s: %w", file.Name, err)
	}
	defer rc.Close()

	report, err := dmarcstats.AnalyzeDMARCReport(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze DMARC report for file %s: %w", file.Name, err)
	}

	return report, nil
}
