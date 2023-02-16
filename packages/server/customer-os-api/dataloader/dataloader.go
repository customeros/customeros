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
	EmailsForContact    *dataloader.Loader
}

type batcher struct {
	tagService   service.TagService
	emailService service.EmailService
}

// NewDataLoader returns the instantiated Loaders struct for use in a request
func NewDataLoader(services *service.Services) *Loaders {
	tagBatcher := &batcher{
		tagService: services.TagService,
	}
	emailBatcher := &batcher{
		emailService: services.EmailService,
	}
	return &Loaders{
		TagsForOrganization: dataloader.NewBatchedLoader(tagBatcher.getTagsForOrganizations, dataloader.WithClearCacheOnBatch()),
		TagsForContact:      dataloader.NewBatchedLoader(tagBatcher.getTagsForContacts, dataloader.WithClearCacheOnBatch()),
		EmailsForContact:    dataloader.NewBatchedLoader(emailBatcher.getEmailsForContacts, dataloader.WithClearCacheOnBatch()),
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
