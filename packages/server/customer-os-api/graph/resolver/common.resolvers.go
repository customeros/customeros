package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"

	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go/log"
)

// AddTag is the resolver for the addTag field.
func (r *mutationResolver) AddTag(ctx context.Context, input model.AddTagInput) (string, error) {
	ctx, span := tracing.StartGraphQLTracerSpan(ctx, "CommonResolver.AddTag", graphql.GetOperationContext(ctx))
	defer span.Finish()
	tracing.SetDefaultResolverSpanTags(ctx, span)
	span.LogFields(log.Object("request", input))

	tenant := common.GetTenantFromContext(ctx)

	tagId, err := r.Services.CommonServices.TagService.AddTag(ctx, nil, tenant, input.EntityID, commonModel.GetEntityType(input.EntityType.String()), utils.StringOrEmpty(input.Tag.ID), utils.StringOrEmpty(input.Tag.Name), constants.AppSourceCustomerOsApi)
	if err != nil {
		tracing.TraceErr(span, err)
		graphql.AddErrorf(ctx, "Error adding tag to entity")
		return "", err
	}

	return tagId, nil
}

// RemoveTag is the resolver for the removeTag field.
func (r *mutationResolver) RemoveTag(ctx context.Context, input model.RemoveTagInput) (*model.Result, error) {
	ctx, span := tracing.StartGraphQLTracerSpan(ctx, "CommonResolver.RemoveTag", graphql.GetOperationContext(ctx))
	defer span.Finish()
	tracing.SetDefaultResolverSpanTags(ctx, span)
	span.LogFields(log.Object("request", input))

	tenant := common.GetTenantFromContext(ctx)

	err := r.Services.CommonServices.TagService.RemoveTag(ctx, nil, tenant, input.EntityID, commonModel.GetEntityType(input.EntityType.String()), input.TagID, constants.AppSourceCustomerOsApi)
	if err != nil {
		tracing.TraceErr(span, err)
		graphql.AddErrorf(ctx, "Error adding tag to entity")
		return nil, nil
	}

	return &model.Result{Result: true}, nil
}