package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type WorkspaceService interface {
	MergeToTenant(ctx context.Context, workspaceEntity neo4jentity.WorkspaceEntity, tenant string) (bool, error)
	GetWorkspaceDomainsForTenant(ctx context.Context) ([]string, error)
}

type workspaceService struct {
	services *Services
}

func NewWorkspaceService(services *Services) WorkspaceService {
	return &workspaceService{
		services: services,
	}
}

func (s *workspaceService) MergeToTenant(ctx context.Context, workspaceEntity neo4jentity.WorkspaceEntity, tenant string) (bool, error) {
	_, err := s.services.Neo4jRepositories.WorkspaceWriteRepository.Merge(ctx, workspaceEntity)
	if err != nil {
		return false, fmt.Errorf("MergeToTenant: %w", err)
	}
	result, err := s.services.Neo4jRepositories.TenantWriteRepository.LinkWithWorkspace(ctx, tenant, workspaceEntity)
	return result, err
}

func (s *workspaceService) GetWorkspaceDomainsForTenant(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkspaceService.GetWorkspaceDomainsForTenant")
	defer span.Finish()

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return []string{}, err
	}
	tenant := common.GetTenantFromContext(ctx)

	dbNodes, err := s.services.Neo4jRepositories.WorkspaceReadRepository.GetAllForTenant(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return []string{}, err
	}

	domains := []string{}
	for _, dbNode := range dbNodes {
		workspaceEntity := neo4jmapper.MapDbNodeToWorkspaceEntity(dbNode)
		domains = append(domains, workspaceEntity.Name)
	}
	domains = utils.RemoveEmpties(domains)
	domains = utils.RemoveDuplicates(domains)

	span.LogFields(log.String("domains", fmt.Sprintf("%v", domains)))
	return domains, nil
}
