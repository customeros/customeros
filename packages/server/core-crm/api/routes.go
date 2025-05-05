package api

import (
	"context"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	service "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/customeros/customeros/packages/server/core-crm/api/graphql/generated"
	"github.com/customeros/customeros/packages/server/core-crm/api/graphql/resolver"
	"github.com/customeros/customeros/packages/server/core-crm/api/middleware"
	"github.com/customeros/customeros/packages/server/core-crm/internal/config"
	"github.com/customeros/customeros/packages/server/core-crm/internal/utils"
)

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(ctx context.Context, r *gin.Engine, services *service.CommonServices, config *config.AppConfig) {
	if services == nil {
		panic("Services cannot be nil")
	}
	if config == nil {
		panic("Config cannot be nil")
	}

	// Add recovery middlewares
	r.Use(gin.Recovery()) // Gin's built-in recovery

	// GraphQL API
	graphqlHandler, playgroundHandler := SetupGraphQLServer(services)

	graphql := r.Group("/")
	{
		graphql.GET("/", playgroundHandler) // playground
	}

	query := r.Group("/query")
	query.Use(middleware.TenantValidationMiddleware()) // Tenant header validation
	query.Use(middleware.UserIdMiddleware())           // UserId header parsing
	query.Use(middleware.CustomContextMiddleware())    // Add custom context
	query.Use(middleware.TracingMiddleware(ctx))       // Add tracing with parent context
	{
		query.POST("", graphqlHandler) // query
	}

	return
}

// SetupGraphQLServer configures and returns the GraphQL server and playground handlers
func SetupGraphQLServer(services *service.CommonServices) (graphqlHandler, playgroundHandler gin.HandlerFunc) {
	// Create the resolver with dependencies
	resolver := resolver.NewResolver(services)

	// Create a new schema with your resolvers
	schema := generated.NewExecutableSchema(generated.Config{
		Resolvers:  resolver,
		Directives: generated.DirectiveRoot{},
		Complexity: generated.ComplexityRoot{},
	})

	// Create the GraphQL server with custom options
	srv := handler.New(schema)

	// Configure server options
	srv.AddTransport(transport.POST{})          // Support POST requests
	srv.AddTransport(transport.GET{})           // Support GET requests
	srv.AddTransport(transport.MultipartForm{}) // Support multipart form
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	// Add extensions
	srv.Use(extension.Introspection{}) // Enable introspection
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Create playground handler
	playground := playground.Handler("GraphQL", "/query")

	// Return handlers wrapped for Gin with our custom middleware
	return func(c *gin.Context) {
		// Explicitly add the Gin context to the request context
		ginCtx := middleware.GinContextToContextMiddleware()
		ginCtx(c)

		// Add custom middleware to extract tenant from Gin context
		c.Request = c.Request.WithContext(utils.WithTenantContext(c.Request.Context(), c.GetString("Tenant")))

		// Call the GraphQL handler
		srv.ServeHTTP(c.Writer, c.Request)
	}, gin.WrapH(playground)
}
