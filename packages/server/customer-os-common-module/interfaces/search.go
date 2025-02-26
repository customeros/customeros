package interfaces

import "context"

type SearchService interface {
	SearchWebsites(ctx context.Context, primaryDomain, query string) (string, error)
}
