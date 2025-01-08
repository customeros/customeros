package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
)

type InteractionSessionService interface {
	GetIdForReferencedInteractionSession(ctx context.Context, tenant, externalSystem string, user model.ReferencedInteractionSession) (string, error)
}

type interactionSessionService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewInteractionSessionService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) InteractionSessionService {
	return &interactionSessionService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (s *interactionSessionService) GetIdForReferencedInteractionSession(ctx context.Context, tenant, externalSystemId string, interactionSession model.ReferencedInteractionSession) (string, error) {
	if !interactionSession.Available() {
		return "", nil
	}

	if interactionSession.ReferencedByExternalId() {
		return s.repositories.InteractionSessionRepository.GetInteractionSessionIdByExternalId(ctx, tenant, interactionSession.ExternalId, externalSystemId)
	}
	return "", nil
}
