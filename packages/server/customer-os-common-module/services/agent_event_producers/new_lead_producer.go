package agent_producers

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type NewLeadProducer struct {
	postgresRepository  *postgres_repository.Repositories
	neo4jRepository     *neo4j_repository.Repositories
	organizationService interfaces.OrganizationService
	events              *events.EventsService
}

func NewNewLeadProducer(
	postgresRepository *postgres_repository.Repositories,
	neo4jRepository *neo4j_repository.Repositories,
	organizationService interfaces.OrganizationService,
	events *events.EventsService,
) *NewLeadProducer {
	return &NewLeadProducer{
		postgresRepository:  postgresRepository,
		neo4jRepository:     neo4jRepository,
		organizationService: organizationService,
		events:              events,
	}
}

// Add all Agent types subscribed to this event here
func (p *NewLeadProducer) subscribedAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentICPQualifier,
	}
}

func (p *NewLeadProducer) Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	limit := 100
	delayFromPreviousCheckRequestInMinutes := 24 * 60 // 24 hours

	span, ctx := tracing.StartTracerSpan(ctx, "NewLeadProducer.NewLeads")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	// get active icp agents
	icpAgents, err := p.postgresRepository.AgentRepository.GetActiveConfiguredAgentsByTypesCrossTenant(ctx, p.subscribedAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
	var tenants []string
	for _, agent := range icpAgents {
		tenants = append(tenants, agent.Tenant)
	}

	if len(tenants) == 0 {
		span.LogKV("message", "No active icp agents found")
		return
	}

	records, err := p.neo4jRepository.OrganizationReadRepository.GetOrganizationsForIcpCheck(ctx, tenants, limit, delayFromPreviousCheckRequestInMinutes)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	// process organizations
	for _, record := range records {
		p.processLeads(ctx, record)
	}
}

func (p *NewLeadProducer) processLeads(ctx context.Context, record neo4j_repository.TenantAndOrganizationId) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewLeadProducer.processsLeads")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    record.Tenant,
		AppSource: constants.AppSourceUpkeeper,
	})
	recordSpan, innerCtx := tracing.StartTracerSpan(innerCtx, "OrganizationService.IcpCheck.Record")
	tracing.TagEntity(recordSpan, record.OrganizationId)
	tracing.TagTenant(recordSpan, record.Tenant)

	err := p.neo4jRepository.CommonWriteRepository.UpdateTimeProperty(innerCtx, record.Tenant, model.NodeLabelOrganization, record.OrganizationId, string(neo4j_entity.OrganizationPropertyIcpCheckRequestedAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(recordSpan, err)
		return
	}

	// check if any organization domain is known global org
	globalOrgs, err := p.organizationService.GetGlobalOrganizationsByTenantOrganizationId(innerCtx, record.OrganizationId)
	if err != nil {
		tracing.TraceErr(recordSpan, err)
		return
	}
	if len(globalOrgs) == 0 {
		span.LogKV("message", "Organization is not a global organization")
		return
	}

	err = p.events.Publisher.PublishFanoutEvent(innerCtx, record.OrganizationId, model.ORGANIZATION, dto.NewLead{})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error publishing new lead event"))
		return
	}
}
