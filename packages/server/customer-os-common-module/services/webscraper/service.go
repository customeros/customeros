package webscraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type webscraperService struct {
	config               *config.JinaConfig
	postgresRepositories *postgres_repository.Repositories
}

func NewWebscraperService(config *config.JinaConfig, postgres *postgres_repository.Repositories) interfaces.WebscraperService {
	return &webscraperService{
		config:               config,
		postgresRepositories: postgres,
	}
}

func (s *webscraperService) Scrape(ctx context.Context, url, primaryDomain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebscraperService.Scrape")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("url", url), log.String("primaryDomain", primaryDomain))

	// fetch page contents
	contents, err := s.fetchPage(ctx, url)
	if err != nil {
		if isPaymentRequiredError(err) {
			// Log the 402 error but don't treat it as a failure
			span.LogFields(log.String("event", "payment_required"))
			return nil
		}
		tracing.TraceErr(span, err)
		return err
	}

	_, err = s.postgresRepositories.GlobalOrganizationWebpageRepository.Save(ctx, postgres_entity.GlobalOrganizationWebpages{
		Url:           url,
		PrimaryDomain: primaryDomain,
		Content:       contents,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// todo parse page urls and scrape them

	return nil
}

func (s *webscraperService) fetchPage(ctx context.Context, url string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebscraperService.fetchPage")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	requestUrl := s.config.Url + url
	span.LogKV("requestUrl", requestUrl)

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	// Align headers with working curl command
	req.Header.Set("Authorization", "Bearer "+s.config.ApiKey)
	req.Header.Set("X-Return-Format", "markdown")  // Get markdown output
	req.Header.Set("X-Retain-Images", "none")      // Don't retain images
	req.Header.Set("X-With-Links-Summary", "true") // Include links summary

	// Add a timeout
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("error code: %d", resp.StatusCode)
		tracing.TraceErr(span, err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error reading response body: %v", err)
		tracing.TraceErr(span, err)
		return "", err
	}

	return string(body), nil
}

func isPaymentRequiredError(err error) bool {
	// Check for HTTP errors with status code 402
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		// For standard http client errors
		if resp, ok := urlErr.Unwrap().(interface{ StatusCode() int }); ok {
			return resp.StatusCode() == 402
		}
	}

	// Try to find error types that embed an HTTP response
	type statusCoder interface {
		StatusCode() int
	}

	var scErr statusCoder
	if errors.As(err, &scErr) {
		return scErr.StatusCode() == 402
	}

	// Fallback to string checking for other HTTP client implementations
	return strings.Contains(err.Error(), "402") ||
		strings.Contains(strings.ToLower(err.Error()), "payment required")
}
