package services

import (
	"context"

	"github.com/customeros/customeros/packages/server/leads/internal/config"
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/services/ipdata"
	"github.com/customeros/customeros/packages/server/leads/services/proxy_manager"
	"github.com/customeros/customeros/packages/server/leads/services/session_manager"
	"github.com/customeros/customeros/packages/server/leads/services/snitcher"
	"github.com/customeros/customeros/packages/server/leads/services/visitor_identity"
	"github.com/customeros/customeros/packages/server/leads/services/web_event_processor"
	"github.com/customeros/customeros/packages/server/leads/services/webtracker"
)

type Services struct {
	IPDataService          *ipdata.IPDataService
	ProxyManager           *proxy_manager.ProxyManager
	SessionManager         *session_manager.SessionManager
	SnitcherService        *snitcher.SnitcherService
	VisitorIdentityService *visitor_identity.VisitorIdentityService
	WebEventProcessor      web_event_processor.WebEventProcessor
	WebtrackerService      webtracker.WebtrackerService
}

func (s *Services) Start(ctx context.Context) error {
	return nil
}

func (s *Services) Stop(ctx context.Context) {
	return
}

func InitServices(natsConn *nats_internal.NATSConnections, repositories *repository.Repositories, config *config.Config) *Services {
	return nil
}
