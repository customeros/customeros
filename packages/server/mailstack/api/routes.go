package api

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/mailstack/api/handlers"
	"github.com/customeros/customeros/packages/server/mailstack/internal/repository"
	"github.com/customeros/customeros/packages/server/mailstack/services"
)

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(ctx context.Context, r *gin.Engine, s *services.Services, repos *repository.Repositories) {
	// Health check and status endpoints
	r.GET("/health", handlers.HealthCheck)
	r.GET("/status", handlers.Status(s.IMAPService))

	// API group with version
	api := r.Group("/v1",
		security.ApiKeyCheckerHTTP(
			repos.TenantWebhookAPIKeyRepository,
			repos.AppKeyRepository,
			security.PLATFORM_ADMIN_API,
			security.WithCache(s.Cache),
		),
	)
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
