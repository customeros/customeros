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

type DomainService interface {
	GetDomainsForOrganizations(ctx context.Context, organizationIds []string) (*entity.DomainEntities, error)
	MergeToTenant(ctx context.Context, domainEntity entity.DomainEntity, tenant string) error
}

type domainService struct {
	repositories *repository.Repositories
}

func NewDomainService(repository *repository.Repositories) DomainService {
	return &domainService{
		repositories: repository,
	}
}

func (s *domainService) GetDomainsForOrganizations(ctx context.Context, organizationIds []string) (*entity.DomainEntities, error) {
	tags, err := s.repositories.DomainRepository.GetForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		return nil, err
	}
	domainEntities := entity.DomainEntities{}
	for _, v := range tags {
		domainEntity := s.mapDbNodeToDomainEntity(*v.Node)
		domainEntity.DataloaderKey = v.LinkedNodeId
		domainEntities = append(domainEntities, *domainEntity)
	}
	return &domainEntities, nil
}

func (s *domainService) MergeToTenant(ctx context.Context, domainEntity entity.DomainEntity, tenant string) error {
	_, err := s.repositories.DomainRepository.Merge(ctx, domainEntity)
	if err != nil {
		return fmt.Errorf("MergeToTenant: %w", err)
	}
	err = s.repositories.TenantRepository.LinkWithDomain(ctx, tenant, domainEntity.Domain)
	return err
}

func (s *domainService) mapDbNodeToDomainEntity(dbNode dbtype.Node) *entity.DomainEntity {
	props := utils.GetPropsFromNode(dbNode)
	domain := entity.DomainEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Domain:    utils.GetStringPropOrEmpty(props, "domain"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource: utils.GetStringPropOrEmpty(props, "appSource"),
		Source:    entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
	}
	return &domain
}
