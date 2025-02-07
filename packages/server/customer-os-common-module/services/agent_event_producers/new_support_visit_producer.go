package agent_producers

import (
	"context"
	"strings"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type NewSupportVisitProducer struct {
	events             *events.EventsService
	postgresRepository *postgres_repository.Repositories
}

func NewNewSupportVisitProducer(
	events *events.EventsService,
	postgresRepository *postgres_repository.Repositories,
) *NewSupportVisitProducer {
	return &NewSupportVisitProducer{
		events:             events,
		postgresRepository: postgresRepository,
	}
}

const WebSessionLookbackPeriodInHours = 48

type NewSupporVisitConfig struct {
	Webpages    agent_capability.ConfigMultipleValues `json:"webpages"`
	UrlPatterns agent_capability.ConfigMultipleValues `json:"urlPatterns"`
}

// Add all Agent types subscribed to this event here
func (p *NewSupportVisitProducer) subscribedAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentSupportSpotter,
	}
}

func (p *NewSupportVisitProducer) NewConfig() NewSupporVisitConfig {
	return NewSupporVisitConfig{}
}

func (p *NewSupportVisitProducer) Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "NewSupportVisitProducer.Execute")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	// get active subscribed agents (make common?)
	activeAgents, tenants, err := p.getActiveAgents(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
	if activeAgents == nil || tenants == nil {
		return
	}

	for _, agent := range activeAgents {
		ctx = common.SetTenantInContext(ctx, agent.Tenant)
		err := p.processAgent(ctx, agent)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}
	return
}

func (p *NewSupportVisitProducer) processAgent(ctx context.Context, agent postgres_entity.Agent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewLeadProducer.processAgent")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	// get rules for event production from listener config
	config := p.NewConfig()
	err := agent.GetListenerConfigByType(enum.EventNewSupportVisit, &config)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	if len(config.Webpages.Value) == 0 && len(config.UrlPatterns.Value) == 0 {
		span.LogKV("configEmpty", true)
		return nil
	}

	// get websessions ready for processing
	sessions, err := p.postgresRepository.WebSessionRepository.FindAllSessionsForSupportAnalysis(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var errs error
	for _, session := range sessions {
		isHelpNeeded := postgres_entity.NoSupportNeedDetected

		if p.processWebSessionForHelpNeeded(ctx, session, config) {
			err := p.events.Publisher.PublishFanoutEvent(ctx, *session.OrganizationId, model.ORGANIZATION, dto.NewSupportVisit{})
			if err != nil {
				tracing.TraceErr(span, err)
				errs = multierr.Append(errs, err)
			}
			isHelpNeeded = postgres_entity.SupportNeedDetected
		}

		_, err = p.postgresRepository.WebSessionRepository.UpdateSupportSignals(ctx, session.ID, common.GetTenantFromContext(ctx), isHelpNeeded)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}
	return errs
}

func (p *NewSupportVisitProducer) processWebSessionForHelpNeeded(ctx context.Context, session postgres_entity.WebSession, config NewSupporVisitConfig) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewLeadProducer.processWebSessionForHelpNeeded")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	if p.processPageViewsForWebpageMatch(ctx, session.UniquePageViews, config.Webpages.Value) {
		return true
	}

	if p.processPageViewsForUrlPatternMatch(ctx, session.UniquePageViews, config.UrlPatterns.Value) {
		return true
	}
	return false
}

func (p *NewSupportVisitProducer) processPageViewsForWebpageMatch(ctx context.Context, pages []string, targetPages []string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewLeadProducer.processPageViewsForWebpageMatch")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	for _, page := range pages {
		for _, targetPage := range targetPages {
			if strings.Contains(page, targetPage) {
				return true
			}
		}
	}
	return false
}

func (p *NewSupportVisitProducer) processPageViewsForUrlPatternMatch(ctx context.Context, pages []string, urlPatterns []string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewLeadProducer.processPageViewsForUrlPatternMatch")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	for _, page := range pages {
		for _, urlPattern := range urlPatterns {
			if utils.MatchUrlPattern(urlPattern, page) {
				return true
			}
		}
	}
	return false
}

func (p *NewSupportVisitProducer) getActiveAgents(ctx context.Context) ([]postgres_entity.Agent, []string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewLeadProducer.processsLeads")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	agents, err := p.postgresRepository.AgentRepository.GetActiveConfiguredAgentsByTypesCrossTenant(ctx, p.subscribedAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, nil, err
	}

	var tenants []string
	for _, agent := range agents {
		tenants = append(tenants, agent.Tenant)
	}

	if len(tenants) == 0 {
		span.LogKV("message", "No active icp agents found")
		return nil, nil, nil
	}

	return agents, tenants, nil
}
