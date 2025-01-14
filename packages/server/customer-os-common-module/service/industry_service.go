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
	GetInUseIndustries(ctx context.Context) (*neo4jentity.IndustryEntities, error)
	GetByCode(ctx context.Context, code string) (*neo4jentity.IndustryEntity, error)
	GetClosestByCode(ctx context.Context, code string) (*neo4jentity.IndustryEntity, error)
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

func (s *industryService) GetInUseIndustries(ctx context.Context) (*neo4jentity.IndustryEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IndustryService.GetInUseIndustries")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	industryDbNodes, err := s.services.Neo4jRepositories.IndustryReadRepository.GetInUseIndustries(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		return nil, err
	}
	industryEntities := make(neo4jentity.IndustryEntities, 0)
	for _, v := range industryDbNodes {
		industryEntity := neo4jmapper.MapDbNodeToIndustryEntity(v)
		industryEntities = append(industryEntities, *industryEntity)
	}
	return &industryEntities, nil
}

// Returns the industry entity by code, nil if not found
func (s *industryService) GetByCode(ctx context.Context, code string) (*neo4jentity.IndustryEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IndustryService.GetByCode")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("code", code))

	industryDbNode, err := s.services.Neo4jRepositories.IndustryReadRepository.GetByCode(ctx, code)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if industryDbNode == nil {
		return nil, nil
	}
	industryEntity := neo4jmapper.MapDbNodeToIndustryEntity(industryDbNode)
	return industryEntity, nil
}

func (s *industryService) GetClosestByCode(ctx context.Context, code string) (*neo4jentity.IndustryEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IndustryService.GetClosestByCode")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("code", code))

	queryCode := code
	// If the code is not found, strip the last character and try again until found or empty
	for {
		if len(queryCode) == 0 {
			break
		}
		industryEntity, err := s.GetByCode(ctx, queryCode)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if industryEntity != nil {
			return industryEntity, nil
		}
		queryCode = queryCode[:len(queryCode)-1]
	}

	return nil, nil
}
