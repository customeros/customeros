package api

import (
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/mailstack/internal/repository"
	"github.com/customeros/customeros/packages/server/mailstack/services"
)

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(r *gin.Engine, s *services.Services, repos *repository.Repositories) {
	// Health check and status endpoints
	r.GET("/health", HealthCheck)
	r.GET("/status", Status(s.IMAPService))

	// API group with version
	api := r.Group("/api/v1")
	{
		// Mailbox endpoints
		mailboxes := api.Group("/mailboxes")
		{
			mailboxes.GET("", ListMailboxes(s.IMAPService))
			mailboxes.POST("", AddMailbox(s.IMAPService))
			mailboxes.DELETE("/:id", RemoveMailbox(s.IMAPService))
		}

		// Other API endpoints can go here
	}
}
