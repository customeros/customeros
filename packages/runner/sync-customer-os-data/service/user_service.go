package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
)

type UserService interface {
	GetIdForReferencedUser(ctx context.Context, tenant, externalSystem string, user entity.ReferencedUser) (string, error)
}

type userService struct {
	repositories *repository.Repositories
}

func NewUserService(repositories *repository.Repositories) UserService {
	return &userService{
		repositories: repositories,
	}
}

func (s *userService) GetIdForReferencedUser(ctx context.Context, tenant, externalSystemId string, user entity.ReferencedUser) (string, error) {
	if !user.Available() {
		return "", nil
	}

	if user.ReferencedById() {
		return s.repositories.UserRepository.GetUserIdById(ctx, tenant, user.Id)
	} else if user.ReferencedByExternalId() {
		return s.repositories.UserRepository.GetUserIdByExternalId(ctx, tenant, user.ExternalId, externalSystemId)
	} else if user.ReferencedByExternalOwnerId() {
		return s.repositories.UserRepository.GetUserIdByExternalOwnerId(ctx, tenant, user.ExternalOwnerId, externalSystemId)
	}
	return "", nil
}
