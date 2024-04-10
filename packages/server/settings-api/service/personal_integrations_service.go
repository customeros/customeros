package service

import (
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
)

const CALCOM = "calcom"

type PersonalIntegrationsService interface {
	GetPersonalIntegration(tenantName, email, integration string) (*postgresEntity.PersonalIntegration, error)
	SavePersonalIntegration(postgresEntity.PersonalIntegration) (*postgresEntity.PersonalIntegration, error)
	GetPersonalIntegrations(tenantName, email string) ([]*postgresEntity.PersonalIntegration, error)
}

type personalIntegrationsService struct {
	repositories *repository.PostgresRepositories
	serviceMap   map[string][]keyMapping
	log          logger.Logger
}

func NewPersonalIntegrationsService(repositories *repository.PostgresRepositories, log logger.Logger) PersonalIntegrationsService {
	return &personalIntegrationsService{
		repositories: repositories,
		log:          log,
	}
}

func (s *personalIntegrationsService) GetPersonalIntegrations(tenantName, email string) ([]*postgresEntity.PersonalIntegration, error) {
	res := s.repositories.PostgresRepositories.PersonalIntegrationRepository.FindIntegrations(tenantName, email)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.([]*postgresEntity.PersonalIntegration), nil
}
func (s *personalIntegrationsService) GetPersonalIntegration(tenantName, email, integration string) (*postgresEntity.PersonalIntegration, error) {
	res := s.repositories.PostgresRepositories.PersonalIntegrationRepository.FindIntegration(tenantName, email, integration)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.(*postgresEntity.PersonalIntegration), nil
}
func (s *personalIntegrationsService) SavePersonalIntegration(integration postgresEntity.PersonalIntegration) (*postgresEntity.PersonalIntegration, error) {
	res := s.repositories.PostgresRepositories.PersonalIntegrationRepository.SaveIntegration(integration)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.(*postgresEntity.PersonalIntegration), nil
}
