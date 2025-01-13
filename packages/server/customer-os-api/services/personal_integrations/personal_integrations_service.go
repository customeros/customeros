package api_personal_integrations

import (
	cosapi_interfaces "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"golang.org/x/net/context"
)

const CALCOM = "calcom"

type personalIntegrationsService struct {
	log      logger.Logger
	postgres *repository.Repositories
}

func NewPersonalIntegrationsService(log logger.Logger, postgres *repository.Repositories) cosapi_interfaces.PersonalIntegrationsService {
	return &personalIntegrationsService{
		log:      log,
		postgres: postgres,
	}
}

func (s *personalIntegrationsService) GetPersonalIntegrations(tenantName, email string) ([]*postgresEntity.PersonalIntegration, error) {
	res := s.postgres.PersonalIntegrationRepository.FindIntegrations(context.TODO(), tenantName, email)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.([]*postgresEntity.PersonalIntegration), nil
}

func (s *personalIntegrationsService) GetPersonalIntegration(tenantName, email, integration string) (*postgresEntity.PersonalIntegration, error) {
	res := s.postgres.PersonalIntegrationRepository.FindIntegration(context.TODO(), tenantName, email, integration)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.(*postgresEntity.PersonalIntegration), nil
}

func (s *personalIntegrationsService) SavePersonalIntegration(integration postgresEntity.PersonalIntegration) (*postgresEntity.PersonalIntegration, error) {
	res := s.postgres.PersonalIntegrationRepository.SaveIntegration(context.TODO(), integration)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.(*postgresEntity.PersonalIntegration), nil
}
