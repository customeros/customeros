package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type TenantService interface {
	GetAllTenants(ctx context.Context) ([]*entity.TenantEntity, error)
}

type tenantService struct {
	repositories *repository.Repositories
}

func NewTenantService(repository *repository.Repositories) TenantService {
	return &tenantService{
		repositories: repository,
	}
}

func (s *tenantService) GetAllTenants(ctx context.Context) ([]*entity.TenantEntity, error) {
	nodes, err := s.repositories.TenantRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllTenants: %w", err)
	}

	tenants := make([]*entity.TenantEntity, len(nodes))

	for i, node := range nodes {
		tenants[i] = s.mapDbNodeToTenantEntity(node)
	}

	return tenants, nil
}

func (s *tenantService) mapDbNodeToTenantEntity(dbNode *dbtype.Node) *entity.TenantEntity {

	if dbNode == nil {
		return nil
	}

	props := utils.GetPropsFromNode(*dbNode)
	tenant := entity.TenantEntity{
		Id:   utils.GetStringPropOrEmpty(props, "id"),
		Name: utils.GetStringPropOrEmpty(props, "name"),
	}
	return &tenant
}
