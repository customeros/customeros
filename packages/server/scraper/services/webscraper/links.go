package webscraper

import (
	"bufio"
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

func (s *webscraperService) linksToCrawl(ctx context.Context, content string, workspaceDomains []string) ([]string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "webscraperService.linksToCrawl")
	defer spans.Finish()

	var urls []string

	webpages, err := s.postgresRepositories.ScrapedWebpageRepository.GetAllWebpagesByPrimaryDomains(ctx, workspaceDomains)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if webpages == nil {
		return nil, nil
	}

	for _, webpage := range webpages {
		for _, link := range webpage.Links {
			if s.shouldCrawl(link, workspaceDomains) {
				urls = append(urls, link)
			}
		}
	}

	return urls, nil
}

func (s *webscraperService) extractLinks(linksSection string) []string {
	var urls []string

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
					url = strings.TrimSuffix(url, "#")
					url = strings.TrimSuffix(url, "/")
					if strings.HasPrefix(url, "http") && !utils.IsStringInSlice(url, urls) {
						urls = append(urls, url)
					}
				}
			}
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
