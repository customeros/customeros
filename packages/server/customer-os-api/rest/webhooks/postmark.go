package webhooks

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RedirectToPayInvoice", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		if c.Request.UserAgent() == "" {
			rest.SendError(c, span, http.StatusForbidden, rest.ErrForbidden)
		}

		if !strings.EqualFold(c.Request.UserAgent(), "Postmark") {
			rest.SendError(c, span, http.StatusForbidden, rest.ErrForbidden)
		}

		httpContext := rest.HTTPContext{
			GinContext:     c,
			ServiceContext: &ctx,
			Span:           span,
			Services:       s,
		}

		emailData, err := parseInboundEmail(httpContext)
		if err != nil {
			tracing.LogObjectAsJson(span, "body", c.Request.Body)
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest)
			return
		}

		httpContext.GinContext.JSON(http.StatusAccepted, rest.BuildBaseResponse(rest.StatusProcessing))

		go func() {
			if emailData.IsMonitorEmail() {
				// Process dmarc monitoring report
				err := processDmarcMonitoringReport(httpContext, &emailData)
				if err != nil {
					tracing.TraceErr(httpContext.Span, errors.Wrap(err, "failed to process DMARC report"))
				}
			} else {
				// Process normal mailstack email
				err := processInboundEmail(httpContext, &emailData)
				if err != nil {
					tracing.TraceErr(httpContext.Span, errors.Wrap(err, "failed to process inbound email from Postmark"))
				}
			}
		}()
		return
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
	if attachment.ContentType != "application/zip" {
		return fmt.Errorf("attachment %s is not a zip file", attachment.Name)
	}

	provider := emailData.DMARCReportProvider()

	reports, err := decodeAndReadDMARCReportFile(emailData.Attachments[0].Content)
	if err != nil {
		return fmt.Errorf("cannot parse dmarc report %s from attachment: %v", attachment.Name, err)
	}
	for _, report := range reports {
		dbReport := buildDMARCReport(report, provider, ctx.Tenant)
		ctx.Services.Repositories.PostgresRepositories.MailStackDomainRepository.CreateDMARCReport(
			*ctx.ServiceContext, ctx.Tenant, &dbReport)
	}
	return nil
}

func buildDMARCReport(report dmarcstats.Report, provider, tenant string) entity.DMARCMonitoring {
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

func decodeAndReadDMARCReportFile(attachment string) ([]dmarcstats.Report, error) {
	var reports []dmarcstats.Report

	// Decode base64 string to bytes
	decoded, err := base64.StdEncoding.DecodeString(attachment)
	if err != nil {
		return reports, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create a reader from the decoded bytes
	zipReader, err := zip.NewReader(bytes.NewReader(decoded), int64(len(decoded)))
	if err != nil {
		return reports, fmt.Errorf("failed to create zip reader: %w", err)
	}

	// Read each file in the zip
	for _, file := range zipReader.File {
		rc, err := file.Open()
		if err != nil {
			return reports, fmt.Errorf("failed to open zip file %s: %w", file.Name, err)
		}
		defer rc.Close()

		// Read the file contents
		content, err := io.ReadAll(rc)
		if err != nil {
			return reports, fmt.Errorf("failed to read zip file %s: %w", file.Name, err)
		}
		reader := bytes.NewReader(content)
		report, err := dmarcstats.AnalyzeDMARCReport(reader)
		if err != nil {
			return reports, fmt.Errorf("failed to read get dmarc report for file %s: %w", file.Name, err)
		}
		reports = append(reports, *report)

	}

	return reports, nil
}
