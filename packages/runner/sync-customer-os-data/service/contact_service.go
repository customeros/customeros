package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
)

type ContactService interface {
	GetIdForReferencedContact(ctx context.Context, tenant, externalSystem string, contact entity.ReferencedContact) (string, error)
}

type contactService struct {
	repositories *repository.Repositories
}

func NewContactService(repositories *repository.Repositories) ContactService {
	return &contactService{
		repositories: repositories,
	}
}

func (s *contactService) GetIdForReferencedContact(ctx context.Context, tenant, externalSystemId string, contact entity.ReferencedContact) (string, error) {
	if !contact.Available() {
		return "", nil
	}

	if contact.ReferencedById() {
		return s.repositories.ContactRepository.GetContactIdById(ctx, tenant, contact.Id)
	} else if contact.ReferencedByExternalId() {
		return s.repositories.ContactRepository.GetContactIdByExternalId(ctx, tenant, contact.ExternalId, externalSystemId)
	}
	return "", nil
}
