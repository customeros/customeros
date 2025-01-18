package cosapi_interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
)

type SearchService interface {
	GCliSearch(ctx context.Context, keyword string, limit *int) (*entity.SearchResultEntities, error)
}
