package dataloader

import (
	"context"
	"fmt"
	"github.com.openline-ai.customer-os-analytics-api/graph/model"
	"github.com.openline-ai.customer-os-analytics-api/mapper"
	"github.com.openline-ai.customer-os-analytics-api/repository"
	"github.com.openline-ai.customer-os-analytics-api/repository/entity"
	"github.com/graph-gophers/dataloader"
	gopher_dataloader "github.com/graph-gophers/dataloader"
	"net/http"
)

type loadersString string

const loadersKey = loadersString("dataloaders")

type Loaders struct {
	PageViewsBySessionId *dataloader.Loader
}

type pageViewsBatcher struct {
	repo repository.PageViewRepository
}

func (i *Loaders) GetPageViewsForSession(ctx context.Context, sessionId string) ([]*model.PageView, error) {
	thunk := i.PageViewsBySessionId.Load(ctx, gopher_dataloader.StringKey(sessionId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.([]*model.PageView), nil
}

// NewDataLoader returns the instantiated Loaders struct for use in a request
func NewDataLoader(repositoryContainer *repository.RepositoryContainer) *Loaders {
	pageViews := &pageViewsBatcher{repo: repositoryContainer.PageViewRepo}
	return &Loaders{
		PageViewsBySessionId: dataloader.NewBatchedLoader(pageViews.get),
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

func (b *pageViewsBatcher) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var sessionIDs []string
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	for ix, key := range keys {
		sessionIDs = append(sessionIDs, key.String())
		keyOrder[key.String()] = ix
	}

	queryResult := b.repo.GetAllBySessionIds(sessionIDs)

	if queryResult.Error != nil {
		return []*dataloader.Result{{Data: nil, Error: queryResult.Error}}
	}

	pageViews := mapper.MapPageViews(queryResult.Result.(*entity.PageViewEntities))
	pageViesBySessionId := map[string][]*model.PageView{}
	for _, val := range pageViews {
		pageViesBySessionId[val.SessionId] = append(pageViesBySessionId[val.SessionId], val)
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for sessionId, record := range pageViesBySessionId {
		ix, ok := keyOrder[sessionId]
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, sessionId)
		}
	}
	// fill array positions with errors where not found in DB
	for sessionId, ix := range keyOrder {
		err := fmt.Errorf("page views not found %s", sessionId)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}

	return results
}
