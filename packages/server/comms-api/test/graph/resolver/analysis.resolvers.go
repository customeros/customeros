package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/model"
)

// Describes is the resolver for the describes field.
func (r *analysisResolver) Describes(ctx context.Context, obj *model.Analysis) ([]model.DescriptionNode, error) {
	panic(fmt.Errorf("not implemented: Describes - describes"))
}

// AnalysisCreate is the resolver for the analysis_Create field.
func (r *mutationResolver) AnalysisCreate(ctx context.Context, analysis model.AnalysisInput) (*model.Analysis, error) {
	if r.Resolver.AnalysisCreate != nil {
		return r.Resolver.AnalysisCreate(ctx, analysis)
	}
	panic(fmt.Errorf("not implemented: AnalysisCreate - analysis_Create"))
}

// Analysis is the resolver for the analysis field.
func (r *queryResolver) Analysis(ctx context.Context, id string) (*model.Analysis, error) {
	panic(fmt.Errorf("not implemented: Analysis - analysis"))
}

// Analysis returns generated.AnalysisResolver implementation.
func (r *Resolver) Analysis() generated.AnalysisResolver { return &analysisResolver{r} }

type analysisResolver struct{ *Resolver }
