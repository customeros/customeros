package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type WebscraperService interface {
	Crawl(ctx context.Context, startUrl string) ([]string, error)
	Scrape(ctx context.Context, url string) (string, error)
	ClassifyContentStage(ctx context.Context, url string) (enum.CustomerJourneyStage, error)
	ClassifyWebpageCategory(ctx context.Context, url string) (enum.WebpageCategory, error)
	ClassifyWebpageTopics(ctx context.Context, url string) ([]string, error)
	ProcessWebContent(ctx context.Context, content string) (string, []string)
}
