package agent_producers

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

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

	spans, ctx := telemetry.StartCronSpan(ctx, "NewLeadProducer.NewLeads")
	defer spans.Finish()

	// get active icp agents
	icpAgents, err := p.postgresRepository.AgentRepository.GetActiveConfiguredAgentsByTypesCrossTenant(ctx, p.subscribedAgents())
	if err != nil {
		spans.TraceError(err)
		return
	}
	var tenants []string
	for _, agent := range icpAgents {
		tenants = append(tenants, agent.Tenant)
	}

	if len(tenants) == 0 {
		spans.LogKV("message", "No active icp agents found")
		return
	}

	records, err := p.neo4jRepository.OrganizationReadRepository.GetOrganizationsForIcpCheck(ctx, tenants, limit, delayFromPreviousCheckRequestInMinutes)
	if err != nil {
		spans.TraceError(err)
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
	spans, ctx := telemetry.StartCronSpan(ctx, "NewLeadProducer.processsLeads")
	defer spans.Finish()

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
		spans.LogKV("message", "Organization is not a global organization")
		return
	}
	// precheck if global organziation has required fields
	for _, globalOrg := range globalOrgs {
		if globalOrg.Description == "" {
			spans.LogKV("message", "Global organization has no description")
			return
		}
		if globalOrg.EmployeeCount == 0 {
			spans.LogKV("message", "Global organization has no employee count")
			return
		}
		if globalOrg.IndustryNaicsName == "" {
			spans.LogKV("message", "Global organization has no industry naics name")
			return
		}
		if globalOrg.CountryA2 == "" {
			spans.LogKV("message", "Global organization has no country a2")
			return
		}
	}

	err = p.events.Publisher.PublishFanoutEvent(innerCtx, record.OrganizationId, model.ORGANIZATION, dto.NewLead{})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "error publishing new lead event"))
		return
	}
}
