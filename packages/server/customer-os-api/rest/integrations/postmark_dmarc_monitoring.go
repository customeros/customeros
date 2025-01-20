package integrations

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

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/mailwatcher/dmarkstats"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (h *IntegrationHandler) PostmarkDMARCMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.PostmarkDMARCMonitor", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		// Validate Postmark User-Agent
		if c.Request.UserAgent() == "" || !strings.EqualFold(c.Request.UserAgent(), "Postmark") {
			tracing.TraceErr(span, fmt.Errorf("Invalid user agent %s", c.Request.UserAgent()))
			h.responseHandler.HandleError(c, http.StatusForbidden, nil)
			return
		}

		// Parse email data
		emailData, err := h.parseInboundEmail(c)
		if err != nil {
			tracing.LogObjectAsJson(span, "body", c.Request.Body)
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		// Return accepted response immediately
		h.responseHandler.HandleAccepted(c)

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
				err = h.processDmarcMonitoringReport(c, &emailData)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to process DMARC report"))
				}
			}
		}()
	}
}

func (h *IntegrationHandler) processDmarcMonitoringReport(c *gin.Context, emailData *PostmarkInboundEmailData) error {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Flows.processDmarcMonitoringReport")
	defer span.Finish()
	tracing.TagComponentRest(span)

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
		dbReport := h.buildDMARCReport(c, report, provider)
		// todo - if tenant is empty, don't send report to database
		// leaving this in for now to verify everything is working as expected
		h.services.Repositories.PostgresRepositories.MailStackDomainRepository.CreateDMARCReport(
			ctx, dbReport.Tenant, &dbReport)
	}
	return nil
}

func (h *IntegrationHandler) buildDMARCReport(c *gin.Context, report dmarcstats.Report, provider string) postgres_entity.DMARCMonitoring {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Flows.buildDMARCReport")
	defer span.Finish()
	tracing.TagComponentRest(span)

	tenant, err := h.services.CommonServices.MailstackService.GetTenantForMailstackDomain(ctx, report.Domain)
	if err != nil {
		tracing.TraceErr(span, fmt.Errorf("Unable to get tenant for domain %s: %v", report.Domain, err))
	}

	jsonReport, _ := json.Marshal(report)
	return postgres_entity.DMARCMonitoring{
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
