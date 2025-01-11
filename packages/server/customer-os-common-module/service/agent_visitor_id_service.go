package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type AgentVisitorIDService interface {
	CreateAgent(ctx context.Context) (*entity.Agents, error)
	RunAgent(ctx context.Context, agent *entity.Agents, eventData *data_fields.WebsiteVisitEvent) error
	SaveAgentConfig(ctx context.Context, agent *entity.Agents, agentConfig AgentConfig) (*entity.Agents, error)
	AgentConfig(ctx context.Context, agent *entity.Agents) (*AgentConfig, error)
	AgentCapabilities(ctx context.Context, agent *entity.Agents) (*AgentCapabilities, error)
	SetReferrerQueryParams(ctx context.Context, queryParams []ReferrerQueryParams) (*string, error)
	ReferrerQueryParams(ctx context.Context, session *entity.WebSession) ([]ReferrerQueryParams, error)
}

type agentVisitorIDService struct {
	services *Services
}

func NewAgentVisitorIDService(services *Services) AgentVisitorIDService {
	return &agentVisitorIDService{
		services: services,
	}
}

type AgentConfig struct {
	Websites                    []string `json:"websites" validate:"required"`
	SlackEnabled                bool     `json:"slackEnabled"`
	SlackChannelID              string   `json:"slackChannelId,omitempty"`
	NotificationCooldownInHours int      `json:"notificationCooldownHours"`
}

type AgentCapabilities struct {
	Capabilities []Capability `json:"capabilities"`
}

type Capability struct {
	Name     string `json:"name"`
	Action   string `json:"action"`
	Optional bool   `json:"optional,omitempty"`
}

type ReferrerQueryParams struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

const DefaultNotificationCooldownInHours = 12

func (a *agentVisitorIDService) CreateAgent(ctx context.Context) (*entity.Agents, error) {
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
	masterAgent, err := a.services.PostgresRepositories.AgentRegistryRepository.Find(ctx, enum.AgentVisitorID)
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
	agent := &entity.Agents{
		RegistryID:   masterAgent.ID,
		Tenant:       tenant,
		Name:         masterAgent.Name,
		Capabilities: masterAgent.Capabilities,
		Goal:         masterAgent.Goal,
		IsActive:     false,
		VisibleInUI:  true,
		Icon:         masterAgent.Icon,
		Color:        utils.GetRandomColor(),
	}

	config := AgentConfig{
		SlackEnabled:                false,
		NotificationCooldownInHours: DefaultNotificationCooldownInHours,
	}

	savedAgent, err := a.SaveAgentConfig(ctx, agent, config)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return savedAgent, nil
}

func (a *agentVisitorIDService) RunAgent(ctx context.Context, agent *entity.Agents, eventData *data_fields.WebsiteVisitEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.RunAgent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validate agent
	if agent.RegistryID != enum.AgentVisitorID.String() {
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
	domain, _, err := a.identifyIP(ctx, eventData.IPAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if domain == nil {
		return nil
	}

	// Create org if doesn't exist
	orgID, err := a.services.OrganizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
		Domains: []string{*domain},
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// update session table with ID data
	query := entity.WebSession{
		ID:     eventData.SessionID,
		IP:     eventData.IPAddress,
		Domain: domain,
	}

	session, err := a.services.PostgresRepositories.WebSessionRepository.Update(ctx, query)
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
	isNewVisitor, err := a.isNewWebsiteVisitor(ctx, tenant, eventData.VisitorID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed isNewVisitor lookup"))
		return err
	}

	// log visit on timeline
	timelineMessage, err := a.buildTimelineMessage(ctx, eventData, isNewOrg, isNewVisitor)
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
	_, err = a.services.ActionService.CreateActionForOrganization(ctx, nil, orgID, action)
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

	err = a.services.SlackService.Notify(ctx, tenant, agentConfig.SlackChannelID, message)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return nil
}

func (a *agentVisitorIDService) SaveAgentConfig(ctx context.Context, agent *entity.Agents, agentConfig AgentConfig) (*entity.Agents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.SetConfig")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if agent.IsActive == true && len(agentConfig.Websites) == 0 {
		return nil, errors.New("Cannot turn on agent without a website")
	}

	jsonBytes, err := json.Marshal(agentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal slack config: %v", err)
	}

	jsonStr := string(jsonBytes)
	agent.Config = &jsonStr

	// Save to db
	savedAgent, err := a.services.PostgresRepositories.AgentsRepository.Create(ctx, agent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return savedAgent, nil
}

func (a *agentVisitorIDService) AgentConfig(ctx context.Context, agent *entity.Agents) (*AgentConfig, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.GetConfig")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if agent.RegistryID != enum.AgentVisitorID.String() {
		return nil, fmt.Errorf("agent invalid for this method")
	}

	if agent.Config == nil {
		return nil, nil
	}

	var config AgentConfig
	err := json.Unmarshal([]byte(*agent.Config), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse agent config: %v", err)
	}

	return &config, nil
}

func (a *agentVisitorIDService) AgentCapabilities(ctx context.Context, agent *entity.Agents) (*AgentCapabilities, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.AgentCapabilities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	var caps AgentCapabilities

	if err := json.Unmarshal([]byte(agent.Capabilities), &caps); err != nil {
		return nil, fmt.Errorf("failed to parse capabilities: %v", err)
	}
	return &caps, nil
}

func (a *agentVisitorIDService) buildTimelineMessage(ctx context.Context, eventData *data_fields.WebsiteVisitEvent, isNewOrg, isNewVisitor bool) (*string, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.buildTimelineMessage")
	defer span.Finish()

	query := entity.WebTrackerEvents{
		Tenant:    eventData.Tenant,
		SessionID: eventData.SessionID,
		EventType: enum.WebTrackerPageView.String(),
	}
	pageViews, err := a.services.PostgresRepositories.WebTrackerEventsRepository.FindAll(ctx, query, nil)
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

	sessionDuration, err := a.calculateSessionDuration(ctx, eventData.SessionID)
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

func (a *agentVisitorIDService) calculateSessionDuration(ctx context.Context, sessionID string) (string, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.calculateSessionDuration")
	defer span.Finish()

	query := entity.WebSession{
		ID:       sessionID,
		IsActive: false,
	}
	session, err := a.services.PostgresRepositories.WebSessionRepository.FindSession(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if session == nil {
		return "", nil
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

func (a *agentVisitorIDService) getUniquePageViews(ctx context.Context, pageViews []entity.WebTrackerEvents) ([]string, error) {
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

	snitcherData, err := a.services.EnrichmentService.Snitcher(ctx, ipAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, nil, err
	}

	if snitcherData == nil || snitcherData.Data == nil {
		return nil, nil, nil
	}

	_, primaryDomain := domaincheck.PrimaryDomainCheck(snitcherData.Data.Domain)

	return &primaryDomain, &snitcherData.Data.Profiles.LinkedIn.Handle, nil
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

	query := entity.WebSession{
		Tenant: *tenant,
		Domain: domain,
	}

	results, err := a.services.PostgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, nil)
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

	query := entity.WebSession{
		Tenant:    tenant,
		VisitorID: visitorId,
	}

	results, err := a.services.PostgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	if len(results) == 0 {
		return true, nil
	}
	return false, nil
}

func (a *agentVisitorIDService) isSlackNotificationEnabled(ctx context.Context, agent *entity.Agents) (bool, *AgentConfig, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.isSlackNotificationEnabled")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agentConfig, err := a.AgentConfig(ctx, agent)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, nil, err
	}

	if agentConfig == nil {
		return false, nil, nil
	}

	return agentConfig.SlackEnabled, agentConfig, nil

}

func (a *agentVisitorIDService) buildWebVisitorSlackNotification(ctx context.Context, session *entity.WebSession, orgID string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.buildWebVisitorSlackNotification")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	// get org data from global org table
	globalOrg, err := a.services.PostgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, *session.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if globalOrg == nil {
		err = a.services.PostgresRepositories.GlobalOrganizationWebsiteToProcessRepository.AddWebsiteToProcess(ctx, *session.Domain)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	// Build the text content for the section based on available data
	var contentLines []string
	primaryDomain := *session.Domain
	if globalOrg != nil {
		primaryDomain = globalOrg.PrimaryDomain
	}
	website := "https://" + primaryDomain

	name := *session.Domain
	if globalOrg != nil {
		name = globalOrg.Name
	}
	contentLines = append(contentLines, fmt.Sprintf("<%s|*%s*> ", website, name))
	if globalOrg != nil && globalOrg.Description != "" {
		contentLines = append(contentLines, fmt.Sprintf("%s \n", globalOrg.Description))
	}
	// Add optional fields only if they're not empty
	if website != "" {
		contentLines = append(contentLines, fmt.Sprintf("*Website:* <%s|%s> ", website, primaryDomain))
	}
	if globalOrg != nil && globalOrg.LinkedInUrl != "" && globalOrg.LinkedInAlias != "" {
		contentLines = append(contentLines, fmt.Sprintf("*LinkedIn:* <%s|/%s> ", globalOrg.LinkedInUrl, globalOrg.LinkedInAlias))
	}
	// Only add location if both city and country are available
	if globalOrg != nil && globalOrg.City != "" && globalOrg.CountryA2 != "" {
		contentLines = append(contentLines, fmt.Sprintf("*Location:* %s, %s ", globalOrg.City, globalOrg.CountryA2))
	}
	// Add source/referrer only if it exists
	if session.Referrer != nil {
		referrer := strings.TrimPrefix(*session.Referrer, "https://")
		referrer = strings.TrimPrefix(referrer, "http://")
		referrer = strings.TrimPrefix(referrer, "www.")
		referrer = strings.Trim(referrer, "/")
		contentLines = append(contentLines, fmt.Sprintf("*Source:* <%s|%s> ", session.Referrer, referrer))
	} else {
		contentLines = append(contentLines, "*Source:* Direct ")
	}
	// Join the lines with newlines
	sectionContent := strings.Join(contentLines, "\n")

	// Create the section block without the accessory first
	sectionBlock := fmt.Sprintf(`{
		"type": "section",
		"text": {
			"type": "mrkdwn",
			"text": "%s"
		}
	}`, sectionContent)

	// If logo exists, add the accessory field
	if globalOrg != nil && globalOrg.LogoUrl != "" {
		sectionBlock = fmt.Sprintf(`{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "%s"
			},
			"accessory": {
				"type": "image",
				"image_url": "%s",
				"alt_text": "%s logo"
			}
		}`, sectionContent, globalOrg.LogoUrl, name)
	}

	// Create the final layout using the section block
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
		%s,
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
	]`, name, sectionBlock, orgID)

	return &layoutBlocks, nil
}

func (a *agentVisitorIDService) skipNotification(ctx context.Context, agentConfig *AgentConfig, session *entity.WebSession) (bool, error) {
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
	lastNotification, err := a.services.PostgresRepositories.WebSessionRepository.FindLastNotification(ctx, session.Tenant, *session.Domain)
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

	workspaceDomains, err := a.services.WorkspaceService.GetWorkspaceDomainsForTenant(ctx)
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

func (a *agentVisitorIDService) SetReferrerQueryParams(ctx context.Context, queryParams []ReferrerQueryParams) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.SetReferrerQueryParams")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	bytes, err := json.Marshal(queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal QueryParam: %w", err)
	}
	results := string(bytes)
	return &results, nil
}

func (a *agentVisitorIDService) ReferrerQueryParams(ctx context.Context, session *entity.WebSession) ([]ReferrerQueryParams, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.ReferrerQueryParams")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	var param []ReferrerQueryParams
	if err := json.Unmarshal([]byte(*session.QueryParams), &param); err != nil {
		return nil, fmt.Errorf("failed to unmarshal QueryParam: %w", err)
	}
	return param, nil
}
