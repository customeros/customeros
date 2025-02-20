package interfaces

import "context"

type WebscraperService interface {
	Scrape(ctx context.Context, url, primaryDomain string) error
}
