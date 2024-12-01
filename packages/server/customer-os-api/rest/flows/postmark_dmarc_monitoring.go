package flows

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
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func PostmarkDMARCMonitor(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "PostmarkInboundEmail", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		// Validate Postmark User-Agent
		if c.Request.UserAgent() == "" || !strings.EqualFold(c.Request.UserAgent(), "Postmark") {
			tracing.TraceErr(span, fmt.Errorf("Invalid user agent %s", c.Request.UserAgent()))
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
			}
		}()
	}
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
		// todo - if tenant is empty, don't send report to database
		// leaving this in for now to verify everything is working as expected
		ctx.Services.Repositories.PostgresRepositories.MailStackDomainRepository.CreateDMARCReport(
			*ctx.ServiceContext, ctx.Tenant, &dbReport)
	}
	return nil
}

func buildDMARCReport(ctx rest.HTTPContext, report dmarcstats.Report, provider string) entity.DMARCMonitoring {
	tenant, err := ctx.Services.CommonServices.MailstackService.GetTenantForMailstackDomain(*ctx.ServiceContext, report.Domain)
	if err != nil {
		tracing.TraceErr(ctx.Span, fmt.Errorf("Unable to get tenant for domain %s: %v", report.Domain, err))
	}

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
