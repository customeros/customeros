package tenant

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type tenantService struct {
	log      logger.Logger
	neo4j    *neoRepo.Repositories
	postgres *repository.Repositories
}

func NewTenantService(log logger.Logger, neo4j *neoRepo.Repositories, postgres *repository.Repositories) interfaces.TenantService {
	return &tenantService{
		log:      log,
		neo4j:    neo4j,
		postgres: postgres,
	}
}

func (s *tenantService) GetAllTenants(ctx context.Context) ([]*neo4jentity.TenantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetAllTenants")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	nodes, err := s.neo4j.TenantReadRepository.GetAll(ctx)
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

func (s *tenantService) GetTenantForUserEmail(ctx context.Context, email string) (*neo4jentity.TenantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantForUserEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email))

	tenant, err := s.neo4j.TenantReadRepository.GetTenantForUserEmail(ctx, email)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantForWorkspace: %w", err)
	}

	return neo4jmapper.MapDbNodeToTenantEntity(tenant), nil
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
		existNode, err := s.neo4j.TenantReadRepository.GetTenantByName(ctx, tenantName)
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
	tenant, err := s.neo4j.TenantWriteRepository.CreateTenantIfNotExistAndReturn(ctx, tenantEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	// save tenant in postgres table
	_, err = s.postgres.TenantRepository.Create(ctx, postgresentity.Tenant{
		Name: tenantName,
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}
	// create tenant specific api key
	err = s.postgres.TenantWebhookApiKeyRepository.CreateApiKey(ctx, tenantName)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	err = s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenantEntity.Name, "gmail", "gmail")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	err = s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenantEntity.Name, "slack", "slack")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	err = s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenantEntity.Name, "intercom", "intercom")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	err = s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenantEntity.Name, "gcal", "gcal")
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

	// Step 1: Permanently delete all tenant data from postgres
	err := s.postgres.CommonRepository.PermanentlyDelete(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.postgres.TenantRepository.PermanentlyDelete(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Step 2: Permanently delete all tenant data from neo4j
	// TODO implement incremental delete to avoid crashing neo4j DB
	err = s.neo4j.TenantWriteRepository.HardDeleteTenant(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Step 3: Permanently delete tenant bucket from S3
	// TODO implement this

	// Step 4: Permanently delete tenant data integration app
	// TODO implement this

	// Step 5: Permanently delete tenant mailboxes data OpenSRS
	// TODO implement this

	return nil
}
