package workspace

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/pkg/errors"
	"strings"

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

	// send event
	eventData := dto.AddWorkspaceDomainToTenant{
		Domain:   workspaceEntity.Name,
		Provider: workspaceEntity.Provider,
	}
	err = s.events.Publisher.PublishFanoutEvent(ctx, tenant, model.TENANT, eventData)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddWorkspaceDomainToTenant"))
	}

	return true, err
}

func (s *workspaceService) AddDomainAsWorkspace(ctx context.Context, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkspaceService.AddDomainAsWorkspace")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// check if domain is valid
	if !utils.IsValidDomain(domain) {
		err = errors.New("invalid domain")
		tracing.TraceErr(span, err)
		return err
	}

	// check if domain not registered with other workspace
	isUsed, err := s.IsAnyTenantWorkspaceDomain(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if isUsed {
		err = errors.New("domain already used as workspace")
		tracing.TraceErr(span, err)
		return err
	}

	workspaceEntity := neo4jentity.WorkspaceEntity{
		Name: domain,
	}

	_, err = s.MergeToTenant(ctx, nil, workspaceEntity, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
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

func (s *workspaceService) IsWorkspaceDomain(ctx context.Context, domain string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkspaceService.IsWorkspaceDomain")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	workspaceDomains, err := s.GetWorkspaceDomainsForTenant(ctx)
	tracing.LogObjectAsJson(span, "workspaceDomains", workspaceDomains)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	if len(workspaceDomains) == 0 {
		span.LogFields(log.Bool("result", false))
		return false, nil
	}

	for _, d := range workspaceDomains {
		if strings.EqualFold(d, domain) {
			span.LogFields(log.Bool("result", true))
			return true, nil
		}
	}

	span.LogFields(log.Bool("result", false))
	return false, nil
}

func (s *workspaceService) IsAnyTenantWorkspaceDomain(ctx context.Context, domain string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkspaceService.IsAnyTenantWorkspaceDomain")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	workspaceDbNodes, err := s.neo4j.WorkspaceReadRepository.GetByNameCrossTenant(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	if len(workspaceDbNodes) > 0 {
		span.LogFields(log.Bool("result", true))
		return true, nil
	}

	span.LogFields(log.Bool("result", false))
	return false, nil
}
