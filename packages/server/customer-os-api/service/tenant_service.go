package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"math/rand"
)

type TenantService interface {
	GetTenantForWorkspace(ctx context.Context, workspaceEntity entity.WorkspaceEntity) (*entity.TenantEntity, error)
	GetTenantForUserEmail(ctx context.Context, email string) (*entity.TenantEntity, error)
	Merge(ctx context.Context, tenantEntity entity.TenantEntity) (*entity.TenantEntity, error)
}

type tenantService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewTenantService(log logger.Logger, repository *repository.Repositories) TenantService {
	return &tenantService{
		log:          log,
		repositories: repository,
	}
}

func (s *tenantService) Merge(ctx context.Context, tenantEntity entity.TenantEntity) (*entity.TenantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.Merge")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	for {
		existNode, err := s.repositories.TenantRepository.GetByName(ctx, tenantEntity.Name)
		if err != nil {
			return nil, fmt.Errorf("Merge: %w", err)
		}
		if existNode == nil {
			break
		}
		newTenantName := fmt.Sprintf("%s%d", tenantEntity.Name, rand.Intn(10))
		tenantEntity.Name = newTenantName
	}
	span.LogFields(log.Object("tenantName", tenantEntity.Name))
	tenant, err := s.repositories.TenantRepository.Merge(ctx, tenantEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("Merge: %w", err)
	}
	return s.mapDbNodeToTenantEntity(tenant), nil
}

func (s *tenantService) GetTenantForWorkspace(ctx context.Context, workspaceEntity entity.WorkspaceEntity) (*entity.TenantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantForWorkspace")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.Object("workspace", workspaceEntity))

	tenant, err := s.repositories.TenantRepository.GetForWorkspace(ctx, workspaceEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantForWorkspace: %w", err)
	}

	return s.mapDbNodeToTenantEntity(tenant), nil
}

func (s *tenantService) GetTenantForUserEmail(ctx context.Context, email string) (*entity.TenantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantForUserEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.String("email", email))

	tenant, err := s.repositories.TenantRepository.GetForUserEmail(ctx, email)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantForWorkspace: %w", err)
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
		Source:    neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
	}
	return &tenant
}
