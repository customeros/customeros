package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/99designs/gqlgen/graphql"

	"github.com.openline-ai.customer-os-analytics-api/dataloader"
	"github.com.openline-ai.customer-os-analytics-api/graph/generated"
	"github.com.openline-ai.customer-os-analytics-api/graph/model"
	"github.com.openline-ai.customer-os-analytics-api/mapper"
	"github.com.openline-ai.customer-os-analytics-api/repository/entity"
	"github.com.openline-ai.customer-os-analytics-api/repository/helper"
)

// PageViews is the resolver for the pageViews field.
func (r *appSessionResolver) PageViews(ctx context.Context, obj *model.AppSession) ([]*model.PageView, error) {
	return dataloader.For(ctx).GetPageViewsForSession(ctx, obj.ID)
}

// Sessions is the resolver for the sessions field.
func (r *applicationResolver) Sessions(ctx context.Context, obj *model.Application, timeFilter model.TimeFilter, dataFilter []*model.AppSessionsDataFilter, paginationFilter *model.PaginationFilter) (*model.AppSessionsPage, error) {
	operationResult := r.RepositoryContainer.SessionsRepo.FindAllByApplication(entity.ApplicationUniqueIdentifier{
		Tenant:      obj.Tenant,
		AppId:       obj.Name,
		TrackerName: obj.TrackerName,
	}, timeFilter, dataFilter, paginationFilter.GetPage(), paginationFilter.GetLimit())

	paginatedResult := operationResult.Result.(*helper.Pagination)

	return &model.AppSessionsPage{
		Content:       mapper.MapSessions(paginatedResult.Rows.(*entity.SessionEntities)),
		TotalPages:    paginatedResult.TotalPages,
		TotalElements: paginatedResult.TotalRows,
	}, nil
}

// Application is the resolver for the application field.
func (r *queryResolver) Application(ctx context.Context, id *string) (*model.Application, error) {
	operationResult := r.RepositoryContainer.AppInfoRepo.FindOneById(ctx, *id)
	if operationResult.Result == nil {
		graphql.AddErrorf(ctx, "Application with id %s not found", *id)
		return nil, nil
	}

	return mapper.MapApplication(operationResult.Result.(*entity.ApplicationEntity)), nil
}

// Applications is the resolver for the applications field.
func (r *queryResolver) Applications(ctx context.Context) ([]*model.Application, error) {
	operationResult := r.RepositoryContainer.AppInfoRepo.FindAll(ctx)
	return mapper.MapApplications(operationResult.Result.(*entity.ApplicationEntities)), nil
}

// AppSession returns generated.AppSessionResolver implementation.
func (r *Resolver) AppSession() generated.AppSessionResolver { return &appSessionResolver{r} }

// Application returns generated.ApplicationResolver implementation.
func (r *Resolver) Application() generated.ApplicationResolver { return &applicationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type appSessionResolver struct{ *Resolver }
type applicationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
