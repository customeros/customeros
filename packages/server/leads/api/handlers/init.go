package handlers

import (
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
)

type APIHandlers struct {
	WebEvents *WebsiteEventsHandler
}

func InitHandlers(natsConn *nats_internal.NATSConnections, r *repository.Repository) *APIHandlers {
	return &APIHandlers{
		WebEvents: NewWebsiteEventsHandler(natsConn),
	}
}
