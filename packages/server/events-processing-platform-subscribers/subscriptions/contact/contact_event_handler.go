package contact

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type ContactEventHandler struct {
	repositories *repository.Repositories
	log          logger.Logger
	cfg          *config.Config
	caches       caches.Cache
	grpcClients  *grpc_client.Clients
}

func NewContactEventHandler(repositories *repository.Repositories, log logger.Logger, cfg *config.Config, caches caches.Cache, grpcClients *grpc_client.Clients) *ContactEventHandler {
	return &ContactEventHandler{
		repositories: repositories,
		log:          log,
		cfg:          cfg,
		caches:       caches,
		grpcClients:  grpcClients,
	}
}

func (h *ContactEventHandler) OnEnrichContactRequested(ctx context.Context, evt eventstore.Event) error {
	return nil
}
