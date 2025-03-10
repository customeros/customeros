package api

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/mailstack/api/handlers"
	"github.com/customeros/customeros/packages/server/mailstack/api/middleware"
	"github.com/customeros/customeros/packages/server/mailstack/internal/repository"
	"github.com/customeros/customeros/packages/server/mailstack/services"
)

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(ctx context.Context, r *gin.Engine, s *services.Services, repos *repository.Repositories, apikey string) {
	// Health check and status endpoints
	r.GET("/health", handlers.HealthCheck)
	r.GET("/status", handlers.Status(s.IMAPService))

	apiKeyMiddleware := middleware.APIKeyMiddleware(middleware.APIKeyConfig{
		HeaderName:  "X-CUSTOMER-OS-API-KEY",
		ValidAPIKey: apikey,
	})

	// API group with version
	api := r.Group("/v1")
	api.Use(apiKeyMiddleware)
	{
		// Mailbox endpoints
		mailboxes := api.Group("/mailboxes")
		{
			mailboxes.GET("", handlers.ListMailboxes(s.IMAPService))
			mailboxes.POST("", handlers.AddMailbox(s.IMAPService, repos.MailboxRepository))
			mailboxes.DELETE("/:id", handlers.RemoveMailbox(s.IMAPService))
		}

		// Other API endpoints can go here
	}
}
