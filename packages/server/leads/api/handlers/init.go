package handlers

import (
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/services"
)

type APIHandlers struct {
	WebEvents *WebsiteEventsHandler
}

func InitHandlers(services *services.Services, r *repository.Repositories) *APIHandlers {
	return &APIHandlers{
		WebEvents: NewWebsiteEventsHandler(services),
	}
}
