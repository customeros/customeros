package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
)

type DomainService interface {
	GetDomainsForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.DomainEntities, error)
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

func (s *domainService) GetDomainsForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.DomainEntities, error) {
	domains, err := s.repositories.DomainRepository.GetForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		return nil, err
	}
	domainEntities := neo4jentity.DomainEntities{}
	for _, v := range domains {
		domainEntity := neo4jmapper.MapDbNodeToDomainEntity(v.Node)
		domainEntity.DataloaderKey = v.LinkedNodeId
		domainEntities = append(domainEntities, *domainEntity)
	}
	return &domainEntities, nil
}
