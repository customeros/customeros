package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"math/rand"
	"strings"
)

type TenantService interface {
	GetAllTenants(ctx context.Context) ([]*neo4jentity.TenantEntity, error)
	GetTenantForWorkspace(ctx context.Context, workspaceEntity neo4jentity.WorkspaceEntity) (*neo4jentity.TenantEntity, error)
	GetTenantForUserEmail(ctx context.Context, email string) (*neo4jentity.TenantEntity, error)
	GetTenantSettings(ctx context.Context) (*neo4jentity.TenantSettingsEntity, error)
	GetTenantSettingsForTenant(ctx context.Context, tenant string) (*neo4jentity.TenantSettingsEntity, error)
	GetTenantBillingProfiles(ctx context.Context) (*neo4jentity.TenantBillingProfileEntities, error)
	GetTenantBillingProfile(ctx context.Context, id string) (*neo4jentity.TenantBillingProfileEntity, error)
	GetDefaultTenantBillingProfile(ctx context.Context) (*neo4jentity.TenantBillingProfileEntity, error)

	Merge(ctx context.Context, tenantEntity neo4jentity.TenantEntity) (*neo4jentity.TenantEntity, error)

	HardDelete(ctx context.Context, tenant string) error
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

func (s *tenantService) GetAllTenants(ctx context.Context) ([]*neo4jentity.TenantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetAllTenants")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	nodes, err := s.services.Neo4jRepositories.TenantReadRepository.GetAll(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	tenants := make([]*neo4jentity.TenantEntity, len(nodes))

	for i, node := range nodes {
		tenants[i] = neo4jmapper.MapDbNodeToTenantEntity(node)
	}

	return tenants, nil
}

func (s *tenantService) GetTenantForWorkspace(ctx context.Context, workspaceEntity neo4jentity.WorkspaceEntity) (*neo4jentity.TenantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantForWorkspace")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("workspace", workspaceEntity))

	tenant, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantForWorkspaceProvider(ctx, workspaceEntity.Name, workspaceEntity.Provider)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantForWorkspace: %w", err)
	}

	return neo4jmapper.MapDbNodeToTenantEntity(tenant), nil
}

func (s *tenantService) GetTenantForUserEmail(ctx context.Context, email string) (*neo4jentity.TenantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantForUserEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email))

	tenant, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantForUserEmail(ctx, email)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantForWorkspace: %w", err)
	}

	return neo4jmapper.MapDbNodeToTenantEntity(tenant), nil
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

func (s *tenantService) GetTenantSettingsForTenant(ctx context.Context, tenant string) (*neo4jentity.TenantSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantSettings")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	tracing.TagTenant(span, tenant)

	dbNode, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode), nil
}

func (s *tenantService) GetTenantBillingProfiles(ctx context.Context) (*neo4jentity.TenantBillingProfileEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantBillingProfiles")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbNodes, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantBillingProfiles(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantBillingProfiles: %w", err)
	}

	tenantBillingProfiles := neo4jentity.TenantBillingProfileEntities{}
	for _, dbNode := range dbNodes {
		tenantBillingProfiles = append(tenantBillingProfiles, *neo4jmapper.MapDbNodeToTenantBillingProfileEntity(dbNode))
	}

	return &tenantBillingProfiles, nil
}

func (s *tenantService) GetTenantBillingProfile(ctx context.Context, id string) (*neo4jentity.TenantBillingProfileEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("id", id))

	dbNode, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantBillingProfileById(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantBillingProfile: %w", err)
	}

	return neo4jmapper.MapDbNodeToTenantBillingProfileEntity(dbNode), nil
}

func (s *tenantService) GetDefaultTenantBillingProfile(ctx context.Context) (*neo4jentity.TenantBillingProfileEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetDefaultTenantBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenantBillingProfiles, err := s.GetTenantBillingProfiles(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetDefaultTenantBillingProfile: %w", err)
	}
	if tenantBillingProfiles == nil || len(*tenantBillingProfiles) == 0 {
		return nil, nil
	} else {
		return &(*tenantBillingProfiles)[0], nil
	}
}

func (s *tenantService) Merge(ctx context.Context, tenantEntity neo4jentity.TenantEntity) (*neo4jentity.TenantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.Merge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "tenantEntity", tenantEntity)

	tenantName := strings.ReplaceAll(tenantEntity.Name, " ", "")
	if tenantName == "" {
		err := fmt.Errorf("tenant name is empty")
		tracing.TraceErr(span, err)
		return nil, err
	}

	for i := 0; i < 10; i++ {
		existNode, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantByName(ctx, tenantName)
		if err != nil {
			return nil, fmt.Errorf("merge: %w", err)
		}
		if existNode == nil {
			break
		}
		tenantName = fmt.Sprintf("%s%d", tenantName, rand.Intn(10))
	}

	span.LogFields(log.Object("tenantName", tenantName))
	tenantEntity.Name = tenantName
	tenant, err := s.services.Neo4jRepositories.TenantWriteRepository.CreateTenantIfNotExistAndReturn(ctx, tenantEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	// save tenant in postgres table
	_, err = s.services.PostgresRepositories.TenantRepository.Create(ctx, postgresentity.Tenant{
		Name: tenantName,
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}
	// create tenant specific api key
	err = s.services.PostgresRepositories.TenantWebhookApiKeyRepository.CreateApiKey(ctx, tenantName)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	err = s.services.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenantEntity.Name, "gmail", "gmail")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	err = s.services.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenantEntity.Name, "slack", "slack")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	err = s.services.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenantEntity.Name, "intercom", "intercom")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	err = s.services.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenantEntity.Name, "gcal", "gcal")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	return neo4jmapper.MapDbNodeToTenantEntity(tenant), nil
}

func (s *tenantService) HardDelete(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.HardDelete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "tenant", tenant)

	err := s.services.Neo4jRepositories.TenantWriteRepository.HardDeleteTenant(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
