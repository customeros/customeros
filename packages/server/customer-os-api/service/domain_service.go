package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type DomainService interface {
	GetDomainsForOrganizations(ctx context.Context, organizationIds []string) (*entity.DomainEntities, error)
}

type domainService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewDomainService(log logger.Logger, repository *repository.Repositories) DomainService {
	return &domainService{
		log:          log,
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
