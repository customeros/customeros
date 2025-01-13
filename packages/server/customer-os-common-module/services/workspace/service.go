package workspace

import (
	"context"
	"errors"
	"fmt"

	"github.com/customeros/mailsherpa/mailvalidate"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type workspaceService struct {
	neo4j *neoRepo.Repositories
}

func NewWorkspaceService(neo4j *neoRepo.Repositories) interfaces.WorkspaceService {
	return &workspaceService{
		neo4j: neo4j,
	}
}

func (s *workspaceService) CheckEmailBelongsToTenant(ctx context.Context, email string) (bool, error) {
	tenantDomains, err := s.GetWorkspaceDomainsForTenant(ctx)
	if err != nil {
		return false, err
	}

	validation := mailvalidate.ValidateEmailSyntax(email)
	if !validation.IsValid {
		return false, errors.New("Email is invalid")
	}

	for _, domain := range tenantDomains {
		if domain == validation.Domain {
			return true, nil
		}
	}
	return false, nil
}

func (s *workspaceService) MergeToTenant(ctx context.Context, workspaceEntity neo4jentity.WorkspaceEntity, tenant string) (bool, error) {
	_, err := s.neo4j.WorkspaceWriteRepository.Merge(ctx, workspaceEntity)
	if err != nil {
		return false, fmt.Errorf("MergeToTenant: %w", err)
	}
	result, err := s.neo4j.TenantWriteRepository.LinkWithWorkspace(ctx, tenant, workspaceEntity)
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

	dbNodes, err := s.neo4j.WorkspaceReadRepository.GetAllForTenant(ctx, tenant)
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
