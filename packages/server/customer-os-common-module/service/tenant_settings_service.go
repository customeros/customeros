package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type TenantSettingsService interface {
	GetTenantSettings(ctx context.Context) (*neo4jentity.TenantSettingsEntity, error)
	GetTenantSettingsForTenant(ctx context.Context, tenant string) (*neo4jentity.TenantSettingsEntity, error)
	UpdateTenantSettings(ctx context.Context, dataFields data_fields.TenantSettingsFields) error
}

type tenantSettingsService struct {
	log      logger.Logger
	services *Services
}

func NewTenantSettingsService(log logger.Logger, services *Services) TenantSettingsService {
	return &tenantSettingsService{
		log:      log,
		services: services,
	}
}

func (s *tenantSettingsService) GetTenantSettings(ctx context.Context) (*neo4jentity.TenantSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsService.GetTenantSettings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbNode, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode), nil
}

func (s *tenantSettingsService) GetTenantSettingsForTenant(ctx context.Context, tenant string) (*neo4jentity.TenantSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsService.GetTenantSettingsForTenant")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, tenant)

	dbNode, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode), nil
}

func (s *tenantSettingsService) UpdateTenantSettings(ctx context.Context, dataFields data_fields.TenantSettingsFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsService.UpdateTenantSettings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	if dataFields.IsEmpty() {
		// nothing to update, return
		return nil
	}

	// update tenant settings in neo4j
	err = s.services.Neo4jRepositories.TenantWriteRepository.UpdateTenantSettings(ctx, tenant, dataFields)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error("Unable to update tenant settings", err)
		return err
	}

	// send event to RabbitMQ
	err = s.services.RabbitMQService.PublishEvent(ctx, tenant, model.TENANT_SETTINGS, dto.UpdateTenantSettings{dataFields})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateTenantSettings"))
	}

	return nil
}
