package service

import (
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"golang.org/x/net/context"
)

const CALCOM = "calcom"

type PersonalIntegrationsService interface {
	GetPersonalIntegration(tenantName, email, integration string) (*postgresEntity.PersonalIntegration, error)
	SavePersonalIntegration(postgresEntity.PersonalIntegration) (*postgresEntity.PersonalIntegration, error)
	GetPersonalIntegrations(tenantName, email string) ([]*postgresEntity.PersonalIntegration, error)
}

type personalIntegrationsService struct {
	services   *Services
	serviceMap map[string][]keyMapping
	log        logger.Logger
}

func NewPersonalIntegrationsService(services *Services, log logger.Logger) PersonalIntegrationsService {
	return &personalIntegrationsService{
		services: services,
		log:      log,
	}
}

func (s *personalIntegrationsService) GetPersonalIntegrations(tenantName, email string) ([]*postgresEntity.PersonalIntegration, error) {
	res := s.services.CommonServices.PostgresRepositories.PersonalIntegrationRepository.FindIntegrations(context.TODO(), tenantName, email)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.([]*postgresEntity.PersonalIntegration), nil
}
func (s *personalIntegrationsService) GetPersonalIntegration(tenantName, email, integration string) (*postgresEntity.PersonalIntegration, error) {
	res := s.services.CommonServices.PostgresRepositories.PersonalIntegrationRepository.FindIntegration(context.TODO(), tenantName, email, integration)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.(*postgresEntity.PersonalIntegration), nil
}
func (s *personalIntegrationsService) SavePersonalIntegration(integration postgresEntity.PersonalIntegration) (*postgresEntity.PersonalIntegration, error) {
	res := s.services.CommonServices.PostgresRepositories.PersonalIntegrationRepository.SaveIntegration(context.TODO(), integration)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.(*postgresEntity.PersonalIntegration), nil
}
