package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com.openline-ai.customer-os-analytics-api/graph/generated"
	"github.com.openline-ai.customer-os-analytics-api/graph/model"
	"github.com.openline-ai.customer-os-analytics-api/mapper"
	"github.com.openline-ai.customer-os-analytics-api/repository/entity"
)

// Sessions is the resolver for the sessions field.
func (r *applicationResolver) Sessions(ctx context.Context, obj *model.Application) ([]*model.AppSession, error) {
	operationResult := r.RepositoryHandler.SessionsRepo.FindAllByApplication(entity.ApplicationUniqueIdentifier{
		Tenant:      obj.Tenant,
		AppId:       obj.Name,
		TrackerName: obj.TrackerName,
	})
	return mapper.MapSessions(operationResult.Result.(*entity.SessionEntities)), nil
}

// Application is the resolver for the application field.
func (r *queryResolver) Application(ctx context.Context, id *string) (*model.Application, error) {
	operationResult := r.RepositoryHandler.AppInfoRepo.FindOneById(*id)
	if operationResult.Result == nil {
		return nil, nil
	}

	return mapper.MapApplication(operationResult.Result.(*entity.ApplicationEntity)), nil
}

// Applications is the resolver for the applications field.
func (r *queryResolver) Applications(ctx context.Context) ([]*model.Application, error) {
	operationResult := r.RepositoryHandler.AppInfoRepo.FindAll()
	return mapper.MapApplications(operationResult.Result.(*entity.ApplicationEntities)), nil
}

// Application returns generated.ApplicationResolver implementation.
func (r *Resolver) Application() generated.ApplicationResolver { return &applicationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type applicationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
