package webscraper

import (
	"context"
	"fmt"
	"io"
	"net/http"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"

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

	// fetch page contents
	contents, err := s.fetchPage(ctx, url)
	if err != nil {
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

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.ApiKey)
	req.Header.Set("X-Return-Format", "markdown")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}
