package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
)

type ContactRoleService interface {
}

type contactRoleService struct {
	repositories *repository.Repositories
}

func NewContactRoleService(repositories *repository.Repositories) ContactRoleService {
	return &contactRoleService{
		repositories: repositories,
	}
}
