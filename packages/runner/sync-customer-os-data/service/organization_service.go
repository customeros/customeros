package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
)

type OrganizationService interface {
	GetIdForReferencedOrganization(ctx context.Context, tenant, externalSystem string, org entity.ReferencedOrganization) (string, error)
}

type organizationService struct {
	repositories *repository.Repositories
}

func NewOrganizationService(repositories *repository.Repositories) OrganizationService {
	return &organizationService{
		repositories: repositories,
	}
}

func (s *organizationService) GetIdForReferencedOrganization(ctx context.Context, tenant, externalSystemId string, org entity.ReferencedOrganization) (string, error) {
	if !org.Available() {
		return "", nil
	}

	if org.ReferencedById() {
		return s.repositories.OrganizationRepository.GetOrganizationIdById(ctx, tenant, org.Id)
	} else if org.ReferencedByExternalId() {
		return s.repositories.OrganizationRepository.GetOrganizationIdByExternalId(ctx, tenant, org.ExternalId, externalSystemId)
	} else if org.ReferencedByDomain() {
		return s.repositories.OrganizationRepository.GetOrganizationIdByDomain(ctx, tenant, org.Domain)
	}
	return "", nil
}
