package webscraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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

var (
	ErrPaymentRequired = errors.New("Jina balance requires topup")
	ErrUnprocessable   = errors.New("Jina cannot process webpage")
)

func (s *webscraperService) Scrape(ctx context.Context, url string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebscraperService.Scrape")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("url", url))

	// fetch page contents
	contents, err := s.fetchPage(ctx, url)
	if err != nil {
		switch err {
		case ErrPaymentRequired:
			tracing.TraceErr(span, err)
			return "", nil

		case ErrUnprocessable:
			span.LogKV("error", ErrUnprocessable)
			return "", err

		default:
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	if strings.Contains(contents, "403 Forbidden") {
		return "", ErrUnprocessable
	}

	return contents, nil
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
		switch resp.StatusCode {
		case 402:
			return "", ErrPaymentRequired

		case 422:
			return "", ErrUnprocessable

		default:
			err = fmt.Errorf("error code: %d", resp.StatusCode)
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error reading response body: %v", err)
		tracing.TraceErr(span, err)
		return "", err
	}

	return string(body), nil
}
