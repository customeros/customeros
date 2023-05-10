package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type TenantService interface {
	GetTenantForDomain(ctx context.Context, domain string) (*entity.TenantEntity, error)
	Merge(ctx context.Context, tenantEntity entity.TenantEntity) (*entity.TenantEntity, error)
}

type tenantService struct {
	repositories *repository.Repositories
}

func NewTenantService(repository *repository.Repositories) TenantService {
	return &tenantService{
		repositories: repository,
	}
}

func (s *tenantService) Merge(ctx context.Context, tenantEntity entity.TenantEntity) (*entity.TenantEntity, error) {
	tenant, err := s.repositories.TenantRepository.Merge(ctx, tenantEntity)
	if err != nil {
		return nil, fmt.Errorf("Merge: %w", err)
	}
	return s.mapDbNodeToTenantEntity(tenant), nil
}

func (s *tenantService) GetTenantForDomain(ctx context.Context, domain string) (*entity.TenantEntity, error) {
	tenant, err := s.repositories.TenantRepository.GetForDomain(ctx, common.GetTenantFromContext(ctx), domain)
	if err != nil {
		return nil, fmt.Errorf("GetTenantForDomain: %w", err)
	}

	return s.mapDbNodeToTenantEntity(tenant), nil
}

func (s *tenantService) mapDbNodeToTenantEntity(dbNode *dbtype.Node) *entity.TenantEntity {

	if dbNode == nil {
		return nil
	}

	props := utils.GetPropsFromNode(*dbNode)
	tenant := entity.TenantEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource: utils.GetStringPropOrEmpty(props, "appSource"),
		Source:    entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
	}
	return &tenant
}
