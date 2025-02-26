package webscraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type webscraperService struct {
	config               *config.JinaConfig
	postgresRepositories *postgres_repository.Repositories
	aiService            interfaces.AIService
	visitedURLs          sync.Map
	limiter              chan struct{}
}

func NewWebscraperService(config *config.JinaConfig, postgres *postgres_repository.Repositories, aiService interfaces.AIService) interfaces.WebscraperService {
	return &webscraperService{
		config:               config,
		postgresRepositories: postgres,
		aiService:            aiService,
		limiter:              make(chan struct{}, 5),
	}
}

const (
	WebpageScrapeTTLInDays = 180
	CleanPageModel         = enum.AIModelGemini
)

var (
	ErrPaymentRequired = errors.New("Jina balance requires topup")
	ErrUnprocessable   = errors.New("Jina cannot process webpage")
)

func (s *webscraperService) Scrape(ctx context.Context, url string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebscraperService.Scrape")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("url", url))

	url = strings.TrimSuffix(url, "/")
	url = strings.TrimSuffix(url, "#")
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	// check cache
	cachedData, err := s.checkCache(ctx, url, WebpageScrapeTTLInDays)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if cachedData != "" {
		return cachedData, nil
	}

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

	if strings.Contains(contents, "Robot Challenge") {
		return "", ErrUnprocessable
	}

	if contents == "" {
		return "", nil
	}

	content, links := s.processWebContent(ctx, contents)

	_, primaryDomain := domaincheck.PrimaryDomainCheck(utils.ExtractDomain(url))

	_, err = s.postgresRepositories.GlobalOrganizationWebpageRepository.Save(ctx, postgres_entity.GlobalOrganizationWebpages{
		PrimaryDomain: primaryDomain,
		Url:           url,
		Content:       content,
		Links:         links,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return contents, nil
}

func (s *webscraperService) processWebContent(ctx context.Context, content string) (string, []string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebscraperService.postProcessWebContent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// extract links and save
	sections := strings.Split(content, "Links/Buttons:")
	if len(sections) < 2 {
		cleanContent := s.processMarkdownWebpage(content)
		return cleanContent, nil
	}

	content = s.processMarkdownWebpage(sections[0])
	links := s.extractLinks(sections[1])

	return content, links
}

func (s *webscraperService) fetchPage(ctx context.Context, url string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebscraperService.fetchPage")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if s.config.ApiKey == "" {
		return "", errors.New("Jina API key not set")
	}

	requestUrl := s.config.Url + url
	span.LogKV("requestUrl", requestUrl)

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	// Align headers with working curl command
	req.Header.Set("Authorization", "Bearer "+s.config.ApiKey)
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

func (s *webscraperService) checkCache(ctx context.Context, url string, cacheTTL int) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "webscraperService.checkCache")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	record, err := s.postgresRepositories.GlobalOrganizationWebpageRepository.GetWebpage(ctx, url, cacheTTL)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if record == nil || record.Content == "" {
		return "", nil
	}

	return record.Content, nil
}
