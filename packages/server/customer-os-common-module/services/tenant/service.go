package tenant

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"math/rand"
	"strings"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
)

type tenantService struct {
	log      logger.Logger
	neo4j    *neo4j_repository.Repositories
	postgres *postgres_repository.Repositories
}

func NewTenantService(log logger.Logger, neo4j *neo4j_repository.Repositories, postgres *postgres_repository.Repositories) interfaces.TenantService {
	return &tenantService{
		log:      log,
		neo4j:    neo4j,
		postgres: postgres,
	}
}

func (s *tenantService) GetAllTenants(ctx context.Context) ([]*neo4jentity.TenantEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantService.GetAllTenants")
	defer spans.Finish()

	nodes, err := s.neo4j.TenantReadRepository.GetAll(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	tenants := make([]*neo4jentity.TenantEntity, len(nodes))

	for i, node := range nodes {
		tenants[i] = neo4jmapper.MapDbNodeToTenantEntity(node)
	}

	return tenants, nil
}

func (s *tenantService) GetTenantForUserEmail(ctx context.Context, email string) (*neo4jentity.TenantEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantService.GetTenantForUserEmail")
	defer spans.Finish()

	spans.LogKV("email", email)

	tenant, err := s.neo4j.TenantReadRepository.GetTenantForUserEmail(ctx, email)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("GetTenantForWorkspace: %w", err)
	}

	return neo4jmapper.MapDbNodeToTenantEntity(tenant), nil
}

func (s *tenantService) Merge(ctx context.Context, tx neo4j.ManagedTransaction, tenantEntity neo4jentity.TenantEntity) (*neo4jentity.TenantEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantService.Merge")
	defer spans.Finish()

	spans.LogObjectAsJson("tenantEntity", tenantEntity)

	tenantName := strings.ReplaceAll(tenantEntity.Name, " ", "")
	if tenantName == "" {
		err := fmt.Errorf("tenant name is empty")
		spans.TraceError(err)
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

	spans.LogKV("tenantName", tenantName)
	tenantEntity.Name = tenantName
	tenant, err := s.neo4j.TenantWriteRepository.CreateTenantIfNotExistAndReturn(ctx, tx, tenantEntity)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	// save tenant in postgres table
	_, err = s.postgres.TenantRepository.Create(ctx, postgresentity.Tenant{
		Name: tenantName,
	})
	if err != nil {
		spans.TraceError(err)
	}
	// create tenant specific api key
	err = s.postgres.TenantWebhookApiKeyRepository.CreateApiKey(ctx, tenantName)
	if err != nil {
		spans.TraceError(err)
	}

	err = s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, &tx, tenantEntity.Name, enum.SourceGmail.String(), enum.SourceGmail.String())
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	err = s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, &tx, tenantEntity.Name, enum.SourceOutlook.String(), enum.SourceOutlook.String())
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	err = s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, &tx, tenantEntity.Name, enum.SourceMailstack.String(), enum.SourceMailstack.String())
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	err = s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, &tx, tenantEntity.Name, enum.SourceSlack.String(), enum.SourceSlack.String())
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	err = s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, &tx, tenantEntity.Name, enum.SourceIntercom.String(), enum.SourceIntercom.String())
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	err = s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, &tx, tenantEntity.Name, enum.SourceGrain.String(), enum.SourceGrain.String())
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	err = s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, &tx, tenantEntity.Name, enum.SourceFathom.String(), enum.SourceFathom.String())
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToTenantEntity(tenant), nil
}

func (s *tenantService) HardDelete(ctx context.Context, tenant string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantService.HardDelete")
	defer spans.Finish()

	// Step 1: Permanently delete all tenant data from postgres
	err := s.postgres.CommonRepository.PermanentlyDelete(ctx, tenant)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	err = s.postgres.TenantRepository.PermanentlyDelete(ctx, tenant)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	// Step 2: Permanently delete all tenant data from neo4j
	// TODO implement incremental delete to avoid crashing neo4j DB
	err = s.neo4j.TenantWriteRepository.HardDeleteTenant(ctx, tenant)
	if err != nil {
		spans.TraceError(err)
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
