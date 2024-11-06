package service

import (
	"context"
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type RegistrationService interface {
	PrepareDefaultTenantSetup(ctx context.Context, loggedInUserEmail string) error
}

type registrationService struct {
	services *Services
}

func NewRegistrationService(services *Services) RegistrationService {
	return &registrationService{
		services: services,
	}
}

func (s *registrationService) MergeToTenant(ctx context.Context, registrationEntity neo4jentity.RegistrationEntity, tenant string) (bool, error) {
	_, err := s.services.Neo4jRepositories.RegistrationWriteRepository.Merge(ctx, registrationEntity)
	if err != nil {
		return false, fmt.Errorf("MergeToTenant: %w", err)
	}
	result, err := s.services.Neo4jRepositories.TenantWriteRepository.LinkWithRegistration(ctx, tenant, registrationEntity)
	return result, err
}
