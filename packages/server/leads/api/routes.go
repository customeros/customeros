package api

import (
	"context"

	"github.com/customeros/mailstack/api/middleware"
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/leads/api/handlers"
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
)

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(ctx context.Context, r *gin.Engine, natsConn *nats_internal.NATSConnections, repos *repository.Repositories) {
	if natsConn == nil {
		panic("Services cannot be nil")
	}
	if repos == nil {
		panic("Repositories cannot be nil")
	}

	// Add recovery middlewares
	r.Use(gin.Recovery()) // Gin's built-in recovery

	// setup handlers
	apiHandlers := handlers.InitHandlers(natsConn, repos)

	// Health check and status endpoints (no custom context needed)
	r.GET("/health", handlers.HealthCheck)

	// Rest API
	api := r.Group("/v1")
	{
		// Domain endpoints
		events := api.Group("/events")
		events.Use(middleware.TenantValidationMiddleware()) // Tenant validation for domains
		events.Use(middleware.CustomContextMiddleware())    // Add custom context
		events.Use(middleware.TracingMiddleware(ctx))       // Add tracing with parent context
		{
			events.POST("", apiHandlers.WebEvents.Handle())
		}
	}
}
