package service

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
)

type WebhookService interface {
	LookupIntegration(integrationID string) (Integration, error)
}

type Integration string

// Add all supported integrations here, and also in NewWebhookService below
const (
	IntegrationCalCom   Integration = "calcom"
	IntegrationFathom   Integration = "fathom"
	IntegrationGrain    Integration = "grain"
	IntegrationPostmark Integration = "postmark"
)

func (i Integration) String() string {
	return string(i)
}

func (i Integration) IntegrationID() string {
	return utils.GenerateHashId(i.String(), 12)
}

type webhookService struct {
	log            logger.Logger
	repositories   *repository.Repositories
	services       *Services
	integrationIDs map[string]Integration
}

func NewWebhookService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) WebhookService {
	integrationHashMap := make(map[string]Integration)

	// Add all supported integrations here
	integrations := []Integration{
		IntegrationCalCom,
		IntegrationFathom,
		IntegrationGrain,
	}

	for _, i := range integrations {
		integrationHashMap[i.IntegrationID()] = i
	}

	return &webhookService{
		log:            log,
		repositories:   repositories,
		services:       services,
		integrationIDs: integrationHashMap,
	}
}

func (w *webhookService) LookupIntegration(integrationId string) (Integration, error) {
	if i, ok := w.integrationIDs[integrationId]; ok {
		return i, nil
	}
	return "", fmt.Errorf("invalid integrationID: %s", integrationId)
}
