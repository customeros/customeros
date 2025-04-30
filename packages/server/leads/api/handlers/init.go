package handlers

import (
	"github.com/customeros/customeros/packages/server/leads/services"
)

type APIHandlers struct {
	WebEvents *WebsiteEventsHandler
}

func InitHandlers(services *services.Services) *APIHandlers {
	return &APIHandlers{
		WebEvents: NewWebsiteEventsHandler(services),
	}
}
