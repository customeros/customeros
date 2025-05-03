package workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type workspaceService struct {
	neo4j  *neo4j_repository.Repositories
	events *events.EventsService
}

func NewWorkspaceService(neo4j *neo4j_repository.Repositories, events *events.EventsService) interfaces.WorkspaceService {
	return &workspaceService{
		neo4j:  neo4j,
		events: events,
	}
}

func (s *workspaceService) MergeToTenant(ctx context.Context, tx *neo4j.ManagedTransaction, workspaceEntity neo4jentity.WorkspaceEntity, tenant string) (bool, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WorkspaceService.MergeToTenant")
	defer spans.Finish()

	spans.LogKV("workspaceEntity", workspaceEntity)
	spans.LogKV("tenant", tenant)

	_, err := s.neo4j.WorkspaceWriteRepository.Merge(ctx, tx, tenant, workspaceEntity)
	if err != nil {
		spans.TraceError(err)
		return false, err
	}

	// send event
	eventData := dto.AddWorkspaceDomainToTenant{
		Domain:   workspaceEntity.Name,
		Provider: workspaceEntity.Provider,
	}
	err = s.events.Publisher.PublishFanoutEvent(ctx, tenant, model.TENANT, eventData)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "unable to publish message AddWorkspaceDomainToTenant"))
	}

	return true, err
}

func (s *workspaceService) AddDomainAsWorkspace(ctx context.Context, domain string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WorkspaceService.AddDomainAsWorkspace")
	defer spans.Finish()

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// check if domain is valid
	if !utils.IsValidDomain(domain) {
		err = errors.New("invalid domain")
		spans.TraceError(err)
		return err
	}

	// check if domain not registered with other workspace
	isUsed, err := s.IsAnyTenantWorkspaceDomain(ctx, domain)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	if isUsed {
		err = errors.New("domain already used as workspace")
		spans.TraceError(err)
		return err
	}

	workspaceEntity := neo4jentity.WorkspaceEntity{
		Name: domain,
	}

	_, err = s.MergeToTenant(ctx, nil, workspaceEntity, tenant)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (s *workspaceService) GetWorkspaceDomainsForTenant(ctx context.Context) ([]string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WorkspaceService.GetWorkspaceDomainsForTenant")
	defer spans.Finish()

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return []string{}, err
	}
	tenant := common.GetTenantFromContext(ctx)

	dbNodes, err := s.neo4j.WorkspaceReadRepository.GetAllForTenant(ctx, tenant)
	if err != nil {
		spans.TraceError(err)
		return []string{}, err
	}

	var domains []string
	for _, dbNode := range dbNodes {
		workspaceEntity := neo4jmapper.MapDbNodeToWorkspaceEntity(dbNode)
		domains = append(domains, workspaceEntity.Name)
	}
	domains = utils.RemoveEmpties(domains)
	domains = utils.RemoveDuplicates(domains)

	spans.LogKV("domains", fmt.Sprintf("%v", domains))
	return domains, nil
}

func (s *workspaceService) IsWorkspaceDomain(ctx context.Context, domain string) (bool, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WorkspaceService.IsWorkspaceDomain")

	defer spans.Finish()

	workspaceDomains, err := s.GetWorkspaceDomainsForTenant(ctx)
	spans.LogObjectAsJson("workspaceDomains", workspaceDomains)
	if err != nil {
		spans.TraceError(err)
		return false, err
	}

	if len(workspaceDomains) == 0 {
		spans.LogFields(log.Bool("result", false))
		return false, nil
	}

	for _, d := range workspaceDomains {
		if strings.EqualFold(d, domain) {
			spans.LogFields(log.Bool("result", true))
			return true, nil
		}
	}

	spans.LogFields(log.Bool("result", false))
	return false, nil
}

func (s *workspaceService) IsAnyTenantWorkspaceDomain(ctx context.Context, domain string) (bool, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WorkspaceService.IsAnyTenantWorkspaceDomain")

	defer spans.Finish()

	workspaceDbNodes, err := s.neo4j.WorkspaceReadRepository.GetByNameCrossTenant(ctx, domain)
	if err != nil {
		spans.TraceError(err)
		return false, err
	}

	if len(workspaceDbNodes) > 0 {
		spans.LogFields(log.Bool("result", true))
		return true, nil
	}

	spans.LogFields(log.Bool("result", false))
	return false, nil
}
