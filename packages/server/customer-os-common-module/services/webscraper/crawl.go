package webscraper

import (
	"bufio"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

const (
	MaxCrawlDepth   = 5
	MaxPagesToCrawl = 100
)

func (s *webscraperService) Crawl(ctx context.Context, startUrl string) ([]string, error) {
	// timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "webscraperService.Crawl")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	_, primaryDomain := domaincheck.PrimaryDomainCheck(utils.ExtractDomain(startUrl))
	if primaryDomain == "" {
		err := errors.New("Not a valid domain")
		tracing.TraceErr(span, err)
		return nil, err
	}

	workspaceDomains := []string{primaryDomain}
	results := make(chan string, MaxPagesToCrawl)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(1)

	// Start recursive crawl
	go s.crawlRecursive(ctx, startUrl, workspaceDomains, 0, &wg, results, errChan)

	// Wait for completion in separate goroutine
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var allUrls []string
	for url := range results {
		allUrls = append(allUrls, url)
	}

	return allUrls, nil
}

func (s *webscraperService) crawlRecursive(
	ctx context.Context,
	url string,
	workspaceDomains []string,
	depth int,
	wg *sync.WaitGroup,
	results chan<- string,
	errChan chan<- error,
) {
	defer wg.Done()

	// Check for already visited URLs
	if _, visited := s.visitedURLs.LoadOrStore(url, true); visited {
		return
	}

	// Check context and depth
	select {
	case <-ctx.Done():
		return
	default:
		if depth >= MaxCrawlDepth {
			return
		}
	}

	// Create span for this URL
	span, ctx := opentracing.StartSpanFromContext(ctx, "crawlRecursive")
	defer span.Finish()
	span.SetTag("url", url)
	span.SetTag("depth", depth)

	// Scrape current URL
	content, err := s.ScrapeAndClean(ctx, url)
	if err != nil {
		select {
		case errChan <- errors.Wrapf(err, "failed to scrape %s", url):
		default:
		}
		return
	}

	// Send URL to results
	results <- url

	// Get new links
	links := s.linksToCrawl(content, workspaceDomains)

	// Rate limiting channel
	limiter := make(chan struct{}, 5)

	// Crawl new links
	for _, link := range links {
		select {
		case <-ctx.Done():
			return
		case limiter <- struct{}{}:
			wg.Add(1)
			go func(link string) {
				defer func() { <-limiter }()
				s.crawlRecursive(ctx, link, workspaceDomains, depth+1, wg, results, errChan)
			}(link)
		}
	}
}

func (s *webscraperService) extractLinks(content string) []string {
	var urls []string

	// Find the Links/Buttons section
	sections := strings.Split(content, "Links/Buttons:")
	if len(sections) < 2 {
		return urls
	}

	// Process the Links/Buttons section
	linksSection := sections[1]
	scanner := bufio.NewScanner(strings.NewReader(linksSection))

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Extract URL from markdown link format [text](url)
		if strings.Contains(line, "](") {
			start := strings.Index(line, "](") + 2
			end := strings.Index(line[start:], ")")
			if end != -1 {
				url := line[start : start+end]
				// Skip empty URLs and fragment-only URLs
				if url != "" && url != "#" {
					urls = append(urls, url)
				}
			}
		}
	}

	return urls
}

func (s *webscraperService) linksToCrawl(content string, workspaceDomains []string) []string {
	var urls []string
	links := s.extractLinks(content)
	for _, link := range links {
		link = strings.TrimSuffix(link, "#")
		if s.shouldCrawl(link, workspaceDomains) {
			urls = append(urls, link)
		}
	}
	return urls
}

func (s *webscraperService) shouldCrawl(url string, workspaceDomains []string) bool {
	if shouldSkipURL(url) {
		return false
	}

	urlDomain := utils.ExtractDomain(url)
	for _, domain := range workspaceDomains {
		if strings.Contains(urlDomain, domain) {
			return true
		}
	}
	return false
}

func shouldSkipURL(url string) bool {
	skipKeywords := []string{
		"login", "logout", "signin", "signout", "register", "signup",
		"account", "profile", "dashboard", "settings", "preferences",
		"cart", "checkout", "order", "payment", "billing",
		"mailto", "tel", "api", "webhook", "feed", "rss",
		"staging", "stage", "dev", "test", "beta", "sandbox",
		"session", "token", "search", "filter", "sort", "page",
		"share", "support", "help", "faq", "ticket", "contact",
		"calendar", "date", "archive", "tag", "admin", "manage",
		"status", "password",
	}

	lowercaseURL := strings.ToLower(url)
	for _, keyword := range skipKeywords {
		if strings.Contains(lowercaseURL, keyword) {
			return true
		}
	}
	return false
}
