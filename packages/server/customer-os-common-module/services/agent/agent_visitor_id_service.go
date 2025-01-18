package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type agentVisitorIDService struct {
	postgresRepositories *postgres_repository.Repositories
	actionService        interfaces.ActionService
	enrichmentService    interfaces.EnrichmentService
	organizationService  interfaces.OrganizationService
	notificationService  interfaces.NotificationService
	workspaceService     interfaces.WorkspaceService
}

func NewAgentVisitorIDService(
	postgres *postgres_repository.Repositories,
	action interfaces.ActionService,
	enrichment interfaces.EnrichmentService,
	org interfaces.OrganizationService,
	notification interfaces.NotificationService,
	workspace interfaces.WorkspaceService,
) interfaces.AgentService {
	return &agentVisitorIDService{
		postgresRepositories: postgres,
		actionService:        action,
		enrichmentService:    enrichment,
		organizationService:  org,
		notificationService:  notification,
		workspaceService:     workspace,
	}
}

type AgentConfig struct {
	Websites                    []string `json:"websites" validate:"required"`
	SlackEnabled                bool     `json:"slackEnabled"`
	SlackChannelID              string   `json:"slackChannelId,omitempty"`
	NotificationCooldownInHours int      `json:"notificationCooldownHours"`
}

const DefaultNotificationCooldownInHours = 12

func (a *agentVisitorIDService) SetActionService(action interfaces.ActionService) {
	a.actionService = action
}

func (a *agentVisitorIDService) SetOrganizationService(org interfaces.OrganizationService) {
	a.organizationService = org
}

func (a *agentVisitorIDService) IsInitialized() bool {
	return utils.IsInitialized(a)
}

func (a *agentVisitorIDService) CreateAgent(ctx context.Context) (*postgres_entity.Agents, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.CreateAgent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not found in context")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// get config from registry
	masterAgent, err := a.postgresRepositories.AgentRegistryRepository.Find(ctx, enum.AgentVisitorID)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if masterAgent == nil {
		err := errors.New("no agent found in agent repository")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// build new agent
	agent := &postgres_entity.Agents{
		Type:         masterAgent.ID,
		Tenant:       tenant,
		Name:         masterAgent.Name,
		Capabilities: masterAgent.Capabilities,
		Goal:         masterAgent.Goal,
		IsActive:     false,
		VisibleInUI:  true,
		Icon:         masterAgent.Icon,
		Color:        utils.GetRandomColor(),
	}

	//config := AgentConfig{
	//	SlackEnabled:                false,
	//	NotificationCooldownInHours: DefaultNotificationCooldownInHours,
	//}

	return agent, nil
}

func (a *agentVisitorIDService) RunAgent(ctx context.Context, agent *postgres_entity.Agents, eventData any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.RunAgent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validate type
	event, ok := eventData.(data_fields.WebsiteVisitEvent)
	if !ok {
		return errors.New("invalid trigger event")
	}

	// validate agent
	if agent.Type != enum.AgentVisitorID.String() {
		return errors.New("agent does not match VisitorID agent")
	}

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		return errors.New("tenant not set on context")
	}

	if tenant != agent.Tenant {
		return errors.New("agent does not belong to tenant")
	}

	// try to deanonymize the IP address
	domain, _, err := a.identifyIP(ctx, event.IPAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if domain == nil {
		return nil
	}

	// Create org if doesn't exist
	orgID, err := a.organizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
		Domains: []string{*domain},
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// update session table with ID data
	query := postgres_entity.WebSession{
		ID:     event.SessionID,
		IP:     event.IPAddress,
		Domain: domain,
	}

	session, err := a.postgresRepositories.WebSessionRepository.Update(ctx, query)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// determine if new org
	isNewOrg, err := a.isNewCompanyVisit(ctx, &tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed isNewCompany lookup"))
		return err
	}

	// determine if new person
	isNewVisitor, err := a.isNewWebsiteVisitor(ctx, tenant, event.VisitorID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed isNewVisitor lookup"))
		return err
	}

	// log visit on timeline
	timelineMessage, err := a.buildTimelineMessage(ctx, &event, isNewOrg, isNewVisitor)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if timelineMessage == nil {
		err := errors.New("unable to build timeline message")
		tracing.TraceErr(span, err)
		return err
	}
	source := enum.SourceAgent.String()
	actionType := enum.ActionGeneric

	action := data_fields.ActionFields{
		Source:     &source,
		CreatedAt:  utils.NowPtr(),
		ActionType: &actionType,
		Content:    timelineMessage,
	}
	_, err = a.actionService.CreateActionForOrganization(ctx, nil, orgID, action)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// determine if slack notification is configured
	slackEnabled, agentConfig, err := a.isSlackNotificationEnabled(ctx, agent)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	if slackEnabled == false || agentConfig.SlackChannelID == "" {
		return nil
	}

	// determine if notification should be skipped
	skip, err := a.skipNotification(ctx, agentConfig, session)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	if skip {
		return nil
	}

	// handle slack notification
	message, err := a.buildWebVisitorSlackNotification(ctx, session, orgID)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	err = a.notificationService.NotifySlackChannel(ctx, tenant, agentConfig.SlackChannelID, message)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return nil
}

func (a *agentVisitorIDService) AgentCapabilities(ctx context.Context, agent *postgres_entity.Agents) (*interfaces.AgentCapabilities, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.AgentCapabilities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	var caps interfaces.AgentCapabilities

	if err := json.Unmarshal([]byte(agent.Capabilities), &caps); err != nil {
		return nil, fmt.Errorf("failed to parse capabilities: %v", err)
	}
	return &caps, nil
}

func (a *agentVisitorIDService) buildTimelineMessage(ctx context.Context, eventData *data_fields.WebsiteVisitEvent, isNewOrg, isNewVisitor bool) (*string, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.buildTimelineMessage")
	defer span.Finish()

	query := postgres_entity.WebTrackerEvents{
		Tenant:    eventData.Tenant,
		SessionID: eventData.SessionID,
		EventType: enum.WebTrackerPageView.String(),
	}
	pageViews, err := a.postgresRepositories.WebTrackerEventsRepository.FindAll(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if pageViews == nil {
		err := errors.New("could not locate any records for session id: " + eventData.SessionID)
		return nil, err
	}

	uniquePageViews, err := a.getUniquePageViews(ctx, pageViews)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	sessionDuration, err := a.calculateSessionDuration(ctx, &postgres_entity.WebSession{
		ID: eventData.SessionID,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Build base message
	var baseMessage string
	switch {
	case isNewOrg:
		baseMessage = fmt.Sprintf("First Visit: %s", sessionDuration)
	case isNewVisitor:
		baseMessage = fmt.Sprintf("New Visitor: %s", sessionDuration)
	case !isNewVisitor:
		baseMessage = fmt.Sprintf("Repeat Visitor: %s", sessionDuration)
	default:
		baseMessage = fmt.Sprintf("Web Visitor: %s", sessionDuration)
	}

	// Append pages with indentation
	var fullMessage strings.Builder
	fullMessage.WriteString(baseMessage)
	for _, page := range uniquePageViews {
		fullMessage.WriteString(fmt.Sprintf("\n    • %s", page))
	}

	timelineMessage := fullMessage.String()
	return &timelineMessage, nil
}

func (a *agentVisitorIDService) calculateSessionDuration(ctx context.Context, session *postgres_entity.WebSession) (string, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.calculateSessionDuration")
	defer span.Finish()

	if session.EndTime.IsZero() {
		session, err := a.postgresRepositories.WebSessionRepository.FindSession(ctx, *session, nil)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		if session == nil {
			return "", nil
		}
	}

	duration := session.EndTime.Sub(session.StartTime)
	minutes := duration.Minutes()

	switch {
	case minutes < 1:
		return "<1 minute", nil
	case minutes == 1:
		return "1 minute", nil
	default:
		return fmt.Sprintf("%.0f minutes", minutes), nil
	}
}

func (a *agentVisitorIDService) getUniquePageViews(ctx context.Context, pageViews []postgres_entity.WebTrackerEvents) ([]string, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.getUniquePageViews")
	defer span.Finish()

	// Use map to track unique pages
	uniquePageMap := make(map[string]struct{})
	for _, page := range pageViews {
		if page.Pathname != "" {
			uniquePageMap[page.Pathname] = struct{}{}
		}
	}

	// Convert map keys to slice
	uniquePages := make([]string, 0, len(uniquePageMap))
	for pathname := range uniquePageMap {
		uniquePages = append(uniquePages, pathname)
	}

	return uniquePages, nil
}

func (a *agentVisitorIDService) identifyIP(ctx context.Context, ipAddress string) (domain, linkedinSlug *string, err error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.identifyIP")
	defer span.Finish()

	snitcherData, err := a.enrichmentService.IPIdentity(ctx, ipAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, nil, err
	}

	if snitcherData == nil {
		return nil, nil, nil
	}

	_, primaryDomain := domaincheck.PrimaryDomainCheck(snitcherData.Company.Domain)

	return &primaryDomain, &snitcherData.Company.Profiles.LinkedIn.Handle, nil
}

func (a *agentVisitorIDService) isNewCompanyVisit(ctx context.Context, tenant, domain *string) (bool, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.isNewCompanyVisit")
	defer span.Finish()

	if tenant == nil || domain == nil {
		err := errors.New("neither tenant or domain can be nil")
		tracing.TraceErr(span, err)
		span.LogKV("tenant", tenant)
		span.LogKV("domain", domain)
		return false, err
	}

	query := postgres_entity.WebSession{
		Tenant: *tenant,
		Domain: domain,
	}

	results, err := a.postgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	if len(results) == 0 {
		return true, nil
	}
	return false, nil
}

func (a *agentVisitorIDService) isNewWebsiteVisitor(ctx context.Context, tenant, visitorId string) (bool, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.isNewWebsiteVisitor")
	defer span.Finish()

	query := postgres_entity.WebSession{
		Tenant:    tenant,
		VisitorID: visitorId,
	}

	results, err := a.postgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	if len(results) == 0 {
		return true, nil
	}
	return false, nil
}

func (a *agentVisitorIDService) isSlackNotificationEnabled(ctx context.Context, agent *postgres_entity.Agents) (bool, *AgentConfig, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.isSlackNotificationEnabled")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	return true, nil, nil

	// TODO implement
	//agentConfig, err := a.AgentConfig(ctx, agent)
	//if err != nil {
	//	tracing.TraceErr(span, err)
	//	return false, nil, err
	//}
	//
	//if agentConfig == nil {
	//	return false, nil, nil
	//}
	//
	//config, ok := agentConfig.(AgentConfig)
	//if !ok {
	//	err = errors.New("Agent does not support this config type")
	//	tracing.TraceErr(span, err)
	//	return false, nil, err
	//}
	//
	//return config.SlackEnabled, &config, nil
}

func (a *agentVisitorIDService) buildWebVisitorSlackNotification(ctx context.Context, session *postgres_entity.WebSession, orgID string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.buildWebVisitorSlackNotification")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	// get org data from global org table
	globalOrg, err := a.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, *session.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if globalOrg == nil {
		err = a.postgresRepositories.GlobalOrganizationWebsiteToProcessRepository.AddWebsiteToProcess(ctx, *session.Domain)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	// Build company info section
	var companyLines []string
	primaryDomain := *session.Domain
	if globalOrg != nil {
		primaryDomain = globalOrg.PrimaryDomain
	}
	website := "https://" + primaryDomain
	name := *session.Domain
	if globalOrg != nil {
		name = globalOrg.Name
	}

	companyLines = append(companyLines, fmt.Sprintf("<%s|*%s*>", website, name))

	if globalOrg != nil && globalOrg.Description != "" {
		companyLines = append(companyLines, fmt.Sprintf("%s", globalOrg.Description))
	}

	if website != "" {
		companyLines = append(companyLines, fmt.Sprintf("*Website:* <%s|%s>", website, primaryDomain))
	}

	if globalOrg != nil && globalOrg.LinkedInUrl != "" && globalOrg.LinkedInAlias != "" {
		companyLines = append(companyLines, fmt.Sprintf("*LinkedIn:* <%s|/%s>", globalOrg.LinkedInUrl, globalOrg.LinkedInAlias))
	}

	if globalOrg != nil && globalOrg.City != "" && globalOrg.CountryA2 != "" {
		companyLines = append(companyLines, fmt.Sprintf("*Location:* %s, %s", globalOrg.City, globalOrg.CountryA2))
	}

	if session.Referrer != nil {
		referrer := strings.TrimPrefix(*session.Referrer, "https://")
		referrer = strings.TrimPrefix(referrer, "http://")
		referrer = strings.TrimPrefix(referrer, "www.")
		referrer = strings.Trim(referrer, "/")
		companyLines = append(companyLines, fmt.Sprintf("*Source:* <%s|%s>", *session.Referrer, referrer))
	} else {
		companyLines = append(companyLines, "*Source:* Direct")
	}

	companyContent := strings.Join(companyLines, "\n")

	// Build session info section
	var sessionLines []string

	if session.EndTime != nil && !session.EndTime.IsZero() {
		duration, err := a.calculateSessionDuration(ctx, session)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		sessionLines = append(sessionLines, fmt.Sprintf("*Session Duration:* %s minutes", duration))
	}

	// Get page views
	query := postgres_entity.WebTrackerEvents{
		SessionID: session.ID,
		EventType: enum.WebTrackerPageView.String(),
	}
	pageViews, err := a.postgresRepositories.WebTrackerEventsRepository.FindAll(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	uniquePages, err := a.getUniquePageViews(ctx, pageViews)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	sessionLines = append(sessionLines, fmt.Sprintf("*Pages Viewed:* %d", len(uniquePages)))
	for _, page := range uniquePages {
		sessionLines = append(sessionLines, fmt.Sprintf("• <%s%s|%s>", website, page, page))
	}

	sessionContent := strings.Join(sessionLines, "\n")

	// Handle logo accessory
	var logoAccessory string
	if globalOrg != nil && globalOrg.LogoUrl != "" {
		logoAccessory = fmt.Sprintf(`,
           "accessory": {
               "type": "image",
               "image_url": "%s",
               "alt_text": "%s logo"
           }`, globalOrg.LogoUrl, name)
	}

	// Build the final layout
	layoutBlocks := fmt.Sprintf(`[
       {
           "type": "header",
           "text": {
               "type": "plain_text",
               "text": "A visitor from %s is on your website",
               "emoji": true
           }
       },
       {
           "type": "divider"
       },
       {
           "type": "section",
           "text": {
               "type": "mrkdwn",
               "text": "%s"
           }%s
       },
       {
           "type": "divider"
       },
       {
           "type": "section",
           "text": {
               "type": "mrkdwn",
               "text": "%s"
           }
       },
       {
           "type": "divider"
       },
       {
           "type": "actions",
           "elements": [
               {
                   "type": "button",
                   "text": {
                       "type": "plain_text",
                       "text": "View in CustomerOS"
                   },
                   "url": "https://app.customeros.ai/organization/%s?tab=about",
                   "value": "click_me_123",
                   "action_id": "actionId-0"
               }
           ]
       }
   ]`,
		name,
		companyContent,
		logoAccessory,
		sessionContent,
		orgID)

	return &layoutBlocks, nil
}

func (a *agentVisitorIDService) skipNotification(ctx context.Context, agentConfig *AgentConfig, session *postgres_entity.WebSession) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.skipNotifications")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	if session.Domain == nil {
		return true, nil
	}

	// don't send if from workspace domain
	isWorkspaceDomain := a.isWorkspaceDomain(ctx, *session.Domain)
	if isWorkspaceDomain {
		return true, nil
	}

	if agentConfig == nil {
		err := errors.New("agent config not set")
		tracing.TraceErr(span, err)
		return true, err
	}

	// determine last notification from this domain
	lastNotification, err := a.postgresRepositories.WebSessionRepository.FindLastNotification(ctx, session.Tenant, *session.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, nil
	}

	if lastNotification == nil {
		return false, nil
	}

	// determine how long since last notification
	hoursSinceLastNotification := time.Now().Sub(*lastNotification.SentSlackNotification).Hours()
	if hoursSinceLastNotification < float64(agentConfig.NotificationCooldownInHours) {
		return true, nil
	}

	return false, nil
}

func (a *agentVisitorIDService) isWorkspaceDomain(ctx context.Context, domain string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.isWorkspaceDomain")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	workspaceDomains, err := a.workspaceService.GetWorkspaceDomainsForTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return false
	}

	for _, d := range workspaceDomains {
		if strings.EqualFold(d, domain) {
			return true
		}
	}

	return false
}
