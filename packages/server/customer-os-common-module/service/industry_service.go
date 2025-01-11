package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type IndustryService interface {
	GetAllForOrganizationIds(ctx context.Context, organizationIds []string) (*neo4jentity.IndustryEntities, error)
}

type industryService struct {
	log      logger.Logger
	services *Services
}

func NewIndustryService(log logger.Logger, services *Services) IndustryService {
	return &industryService{
		log:      log,
		services: services,
	}
}

func (s *industryService) GetAllForOrganizationIds(ctx context.Context, organizationIds []string) (*neo4jentity.IndustryEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IndustryService.GetAllForOrganizationIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationIds", fmt.Sprintf("%v", organizationIds)))

	industryDbNodes, err := s.services.Neo4jRepositories.IndustryReadRepository.GetAllForOrganizationIds(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		return nil, err
	}
	industryEntities := make(neo4jentity.IndustryEntities, 0)
	for _, v := range industryDbNodes {
		industryEntity := neo4jmapper.MapDbNodeToIndustryEntity(v.Node)
		industryEntity.DataloaderKey = v.LinkedNodeId
		industryEntities = append(industryEntities, *industryEntity)
	}
	return &industryEntities, nil
}
