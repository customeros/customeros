package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"net/http"
)

type loadersString string

const loadersKey = loadersString("dataloaders")

type Loaders struct {
	TagsForOrganization *dataloader.Loader
	TagsForContact      *dataloader.Loader
}

type batcher struct {
	tagService service.TagService
}

// NewDataLoader returns the instantiated Loaders struct for use in a request
func NewDataLoader(services *service.Services) *Loaders {
	b := &batcher{tagService: services.TagService}
	return &Loaders{
		TagsForOrganization: dataloader.NewBatchedLoader(b.getTagsForOrganizations, dataloader.WithClearCacheOnBatch()),
		TagsForContact:      dataloader.NewBatchedLoader(b.getTagsForContacts, dataloader.WithClearCacheOnBatch()),
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
