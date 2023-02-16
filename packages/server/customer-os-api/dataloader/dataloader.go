package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"net/http"
)

type loadersString string

const loadersKey = loadersString("dataloaders")

type Loaders struct {
	TagsForOrganizationId *dataloader.Loader
}

type batcher struct {
	tagService service.TagService
}

func (i *Loaders) GetTagsForOrganization(ctx context.Context, organizationId string) (*entity.TagEntities, error) {
	thunk := i.TagsForOrganizationId.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*entity.TagEntities), nil
}

// NewDataLoader returns the instantiated Loaders struct for use in a request
func NewDataLoader(services *service.Services) *Loaders {
	b := &batcher{tagService: services.TagService}
	return &Loaders{
		TagsForOrganizationId: dataloader.NewBatchedLoader(b.getTagsForOrganizations, dataloader.WithClearCacheOnBatch()),
	}
}

// Middleware injects a DataLoader into the request context, so it can be used later in the schema resolvers
func Middleware(loaders *Loaders, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCtx := context.WithValue(r.Context(), loadersKey, loaders)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
