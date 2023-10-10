package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
)

const CALCOM = "calcom"

type PersonalIntegrationsService interface {
	GetPersonalIntegration(tenantName, email, integration string) (*entity.PersonalIntegration, error)
	SavePersonalIntegration(entity.PersonalIntegration) (*entity.PersonalIntegration, error)
	GetPersonalIntegrations(tenantName, email string) ([]*entity.PersonalIntegration, error)
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

func (s *personalIntegrationsService) GetPersonalIntegrations(tenantName, email string) ([]*entity.PersonalIntegration, error) {
	res := s.repositories.CommonRepositories.PersonalIntegrationRepository.FindIntegrations(tenantName, email)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.([]*entity.PersonalIntegration), nil
}
func (s *personalIntegrationsService) GetPersonalIntegration(tenantName, email, integration string) (*entity.PersonalIntegration, error) {
	res := s.repositories.CommonRepositories.PersonalIntegrationRepository.FindIntegration(tenantName, email, integration)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.(*entity.PersonalIntegration), nil
}
func (s *personalIntegrationsService) SavePersonalIntegration(integration entity.PersonalIntegration) (*entity.PersonalIntegration, error) {
	res := s.repositories.CommonRepositories.PersonalIntegrationRepository.SaveIntegration(integration)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.(*entity.PersonalIntegration), nil
}
