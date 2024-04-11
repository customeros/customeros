package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
)

type TenantService interface {
	GetTenantSettings(ctx context.Context) (*neo4jentity.TenantSettingsEntity, error)
}

type tenantService struct {
	log      logger.Logger
	services *Services
}

func NewTenantService(log logger.Logger, services *Services) TenantService {
	return &tenantService{
		log:      log,
		services: services,
	}
}

func (s *tenantService) GetTenantSettings(ctx context.Context) (*neo4jentity.TenantSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantSettings")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	dbNode, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode), nil
}
