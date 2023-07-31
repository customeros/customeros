package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type UserService interface {
	GetAllUsersForTenant(ctx context.Context, tenant string) ([]*entity.UserEntity, error)
}

type userService struct {
	repositories *repository.Repositories
}

func NewUserService(repository *repository.Repositories) UserService {
	return &userService{
		repositories: repository,
	}
}

func (s *userService) GetAllUsersForTenant(ctx context.Context, tenant string) ([]*entity.UserEntity, error) {
	nodes, err := s.repositories.UserRepository.GetAllForTenant(ctx, tenant)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsersForTenant: %w", err)
	}

	users := make([]*entity.UserEntity, len(nodes))

	for i, node := range nodes {
		users[i] = s.mapDbNodeToUserEntity(node)
	}

	return users, nil
}

func (s *userService) mapDbNodeToUserEntity(dbNode *dbtype.Node) *entity.UserEntity {

	if dbNode == nil {
		return nil
	}

	props := utils.GetPropsFromNode(*dbNode)
	user := entity.UserEntity{
		Id: utils.GetStringPropOrEmpty(props, "id"),
	}
	return &user
}
