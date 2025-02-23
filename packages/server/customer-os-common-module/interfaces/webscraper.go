package interfaces

import "context"

type WebscraperService interface {
	Scrape(ctx context.Context, url string) (string, error)
}
