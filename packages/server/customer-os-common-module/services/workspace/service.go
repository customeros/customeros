package workspace

import (
	"context"
	"fmt"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type workspaceService struct {
	neo4j *neo4j_repository.Repositories
}

func NewWorkspaceService(neo4j *neo4j_repository.Repositories) interfaces.WorkspaceService {
	return &workspaceService{
		neo4j: neo4j,
	}
}

func (s *workspaceService) MergeToTenant(ctx context.Context, tx *neo4j.ManagedTransaction, workspaceEntity neo4jentity.WorkspaceEntity, tenant string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkspaceService.MergeToTenant")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogKV("workspaceEntity", workspaceEntity)
	span.LogKV("tenant", tenant)

	_, err := s.neo4j.WorkspaceWriteRepository.Merge(ctx, tx, tenant, workspaceEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	return true, err
}

func (s *workspaceService) GetWorkspaceDomainsForTenant(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkspaceService.GetWorkspaceDomainsForTenant")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return []string{}, err
	}
	tenant := common.GetTenantFromContext(ctx)

	dbNodes, err := s.neo4j.WorkspaceReadRepository.GetAllForTenant(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return []string{}, err
	}

	var domains []string
	for _, dbNode := range dbNodes {
		workspaceEntity := neo4jmapper.MapDbNodeToWorkspaceEntity(dbNode)
		domains = append(domains, workspaceEntity.Name)
	}
	domains = utils.RemoveEmpties(domains)
	domains = utils.RemoveDuplicates(domains)

	span.LogFields(log.String("domains", fmt.Sprintf("%v", domains)))
	return domains, nil
}
