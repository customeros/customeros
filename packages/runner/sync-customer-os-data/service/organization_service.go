package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
)

type OrganizationService interface {
	GetIdForReferencedOrganization(ctx context.Context, tenant, externalSystem string, org entity.ReferencedOrganization) (string, error)
	UpdateLastTouchpointByOrganizationId(ctx context.Context, tenant, organizationID string)
	UpdateLastTouchpointByOrganizationExternalId(ctx context.Context, tenant, organizationExternalId, externalSystem string)
	UpdateLastTouchpointByContactId(ctx context.Context, tenant, contactID string)
	UpdateLastTouchpointByContactIdExternalId(ctx context.Context, tenant, contactExternalId, externalSystem string)
	UpdateLastTouchpointByContactEmailId(ctx context.Context, tenant, emailId string)
}

type organizationService struct {
	repositories *repository.Repositories
}

func NewOrganizationService(repositories *repository.Repositories) OrganizationService {
	return &organizationService{
		repositories: repositories,
	}
}

func (s *organizationService) UpdateLastTouchpointByOrganizationId(ctx context.Context, tenant, organizationID string) {
	if organizationID == "" {
		return
	}

	lastTouchpointAt, lastTouchpointId, err := s.repositories.OrganizationRepository.CalculateAndGetLastTouchpoint(ctx, tenant, organizationID)

	if err != nil {
		return
	}

	if lastTouchpointAt == nil {
		return
	}

	_ = s.repositories.OrganizationRepository.UpdateLastTouchpoint(ctx, tenant, organizationID, *lastTouchpointAt, lastTouchpointId)
}

func (s *organizationService) UpdateLastTouchpointByOrganizationExternalId(ctx context.Context, tenant, organizationExternalId, externalSystem string) {
	orgId, err := s.repositories.OrganizationRepository.GetMatchedOrganizationId(ctx, tenant, entity.OrganizationData{
		BaseData: entity.BaseData{
			ExternalId:     organizationExternalId,
			ExternalSystem: externalSystem,
		},
	})
	if err != nil {
		return
	}
	s.UpdateLastTouchpointByOrganizationId(ctx, tenant, orgId)
}

func (s *organizationService) UpdateLastTouchpointByContactId(ctx context.Context, tenant, contactID string) {
	orgIds, err := s.repositories.OrganizationRepository.GetOrganizationIdsForContact(ctx, tenant, contactID)
	if err != nil {
		return
	}
	for _, orgId := range orgIds {
		s.UpdateLastTouchpointByOrganizationId(ctx, tenant, orgId)
	}
}

func (s *organizationService) UpdateLastTouchpointByContactIdExternalId(ctx context.Context, tenant, contactExternalId, externalSystem string) {
	orgIds, err := s.repositories.OrganizationRepository.GetOrganizationIdsForContactByExternalId(ctx, tenant, contactExternalId, externalSystem)
	if err != nil {
		return
	}
	for _, orgId := range orgIds {
		s.UpdateLastTouchpointByOrganizationId(ctx, tenant, orgId)
	}
}

func (s *organizationService) UpdateLastTouchpointByContactEmailId(ctx context.Context, tenant, emailId string) {
	contactIds, err := s.repositories.ContactRepository.GetContactIdsForEmail(ctx, tenant, emailId)
	if err != nil {
		return
	}
	for _, contactId := range contactIds {
		s.UpdateLastTouchpointByContactId(ctx, tenant, contactId)
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
