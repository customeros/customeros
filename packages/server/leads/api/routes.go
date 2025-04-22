package api

import (
	"context"

	"github.com/customeros/mailstack/api/middleware"
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/leads/api/handlers"
	"github.com/customeros/customeros/packages/server/leads/internal/config"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/services"
)

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(ctx context.Context, r *gin.Engine, s *services.Services, repos *repository.Repositories, cfg *config.Config) {
	if s == nil {
		panic("Services cannot be nil")
	}
	if repos == nil {
		panic("Repositories cannot be nil")
	}

	// Add recovery middlewares
	r.Use(gin.Recovery()) // Gin's built-in recovery

	// setup handlers
	apiHandlers := handlers.InitHandlers(repos, cfg, s)

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
