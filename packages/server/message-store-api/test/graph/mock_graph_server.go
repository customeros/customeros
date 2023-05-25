package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/resolver"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func GraphqlHandler() (gin.HandlerFunc, *resolver.Resolver) {
	graphResolver := &resolver.Resolver{}

	schemaConfig := generated.Config{Resolvers: graphResolver}
	schemaConfig.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.Role) (res interface{}, err error) {
		return next(ctx)
	}
	schemaConfig.Directives.HasTenant = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		return next(ctx)
	}
	schemaConfig.Directives.HasIdentityId = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		return next(ctx)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(schemaConfig))

	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		return gqlerror.Errorf("Internal server error! %v ", err)
	})
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		// Error hook place, Returned error can be customized. Check https://gqlgen.com/reference/errors/
		return err
	})

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}, graphResolver
}
