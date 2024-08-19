package email_validation

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
)

type emailEventHandler struct {
	log         logger.Logger
	cfg         *config.Config
	grpcClients *grpc_client.Clients
	services    *service.Services
}

func NewEmailEventHandler(log logger.Logger, cfg *config.Config, grpcClients *grpc_client.Clients, services *service.Services) *emailEventHandler {
	return &emailEventHandler{
		log:         log,
		cfg:         cfg,
		grpcClients: grpcClients,
		services:    services,
	}
}
