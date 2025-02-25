package interfaces

import "context"

type WebscraperService interface {
	Crawl(ctx context.Context, startUrl string) ([]string, error)
	Scrape(ctx context.Context, url string) (string, error)
	ScrapeAndClean(ctx context.Context, url string) (string, error)
}
