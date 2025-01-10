package service

import (
	"context"
	"encoding/json"
	"fmt"

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
	GetAgentConfig(ctx context.Context, agent entity.Agents) (*AgentConfig, error)
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

	agent := &entity.Agents{
		RegistryID:  enum.AgentVisitorID.String(),
		Tenant:      tenant,
		Name:        "Identify website visitors",
		IsActive:    false,
		VisibleInUI: true,
	}

	config := AgentConfig{
		SlackEnabled:                false,
		NotificationCooldownInHours: 12,
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

	if agent.RegistryID != enum.AgentVisitorID.String() {
		return errors.New("agent does not match VisitorID agent")
	}

	// try to deanonymize the IP address
	domain, linkedinSlug, err := a.identifyIP(ctx, eventData.IPAddress)
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

	_, err = a.services.PostgresRepositories.WebSessionRepository.Update(ctx, query)

	// determine if new org
	tenant := common.GetTenantFromContext(ctx)
	isNewOrg, err := a.isNewCompanyVisit(ctx, &tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed isNewCompany lookup"))
	}

	// determine if new person
	isNewVisitor, err := a.isNewWebsiteVisitor(ctx, tenant, eventData.VisitorID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed isNewVisitor lookup"))
	}

	// log visit on timeline
	timelineMessage, err := a.buildTimelineMessage(ctx, eventData, isNewOrg, isNewVisitor)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	if timelineMessage == nil {
		err := errors.New("unable to build timeline message")
		tracing.TraceErr(span, err)
	}
	source := enum.SourceReveal.String()
	actionType := enum.ActionGeneric

	action := data_fields.ActionFields{
		Source:     &source,
		CreatedAt:  utils.NowPtr(),
		ActionType: &actionType,
		Content:    timelineMessage,
	}

	// determine if slack notification is configured
	sendNotificationEnabled, slackChannelID, err := a.isSlackNotificationEnabled(ctx, agent)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	if sendNotificationEnabled == false || (sendNotificationEnabled && slackChannelID == "") {
		return
	}

	// handle slack notification
	message, err := a.buildWebVisitorSlackNotification(ctx, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	err = a.services.SlackService.Notify(ctx, tenant, slackChannelID, message)
	if err != nil {
		tracing.TraceErr(span, err)
	}

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
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.buildTimelineMessage")
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
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.calculateSessionDuration")
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
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.getUniquePageViews")
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
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.identifyIP")
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

func (a *agentService) isNewCompanyVisit(ctx context.Context, tenant, domain *string) (bool, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.isNewCompanyVisit")
	defer span.Finish()

	query := entity.WebSession{
		Tenant:    *tenant,
		Domain:    domain,
		EventType: string(EventPageExit),
	}

	results, err := a.services.PostgresRepositories.TrackerEventsRepository.FindAll(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	if len(*results) == 0 {
		return true, nil
	}
	return false, nil
}

func (a *agentVisitorIDService) isNewWebsiteVisitor(ctx context.Context, tenant, visitorId string) (bool, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.isNewWebsiteVisitor")
	defer span.Finish()

	query := entity.TrackerEvents{
		Tenant:    tenant,
		VisitorID: visitorId,
		EventType: string(EventPageExit),
	}

	results, err := a.services.PostgresRepositories.TrackerEventsRepository.FindAll(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	if len(*results) == 0 {
		return true, nil
	}
	return false, nil
}

func (a *agentVisitorIDService) buildWebVisitorSlackNotification(ctx context.Context, eventData *data_fields.WebsiteVisitEvent) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.buildWebVisitorSlackNotification")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	// get org data from global org table
	globalOrg, err := s.PostgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, eventData.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if globalOrg == nil {
		err = s.PostgresRepositories.GlobalOrganizationWebsiteToProcessRepository.AddWebsiteToProcess(ctx, eventData.Domain)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	// Build the text content for the section based on available data
	var contentLines []string
	primaryDomain := eventData.Domain
	if globalOrg != nil {
		primaryDomain = globalOrg.PrimaryDomain
	}
	website := "https://" + primaryDomain

	name := eventData.Domain
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
	if eventData.Referrer != "" {
		referrer := strings.TrimPrefix(eventData.Referrer, "https://")
		referrer = strings.TrimPrefix(referrer, "http://")
		referrer = strings.TrimPrefix(referrer, "www.")
		referrer = strings.Trim(referrer, "/")
		contentLines = append(contentLines, fmt.Sprintf("*Source:* <%s|%s> ", eventData.Referrer, referrer))
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
	]`, name, sectionBlock, *eventData.OrganizationID)

	return layoutBlocks, nil
}

func parseURLParams(queryString string) []data_fields.URLParams {
	// Remove leading ? if present
	queryString = strings.TrimPrefix(queryString, "?")

	// Split the string by & to get individual param-value pairs
	pairs := strings.Split(queryString, "&")

	// Create slice to hold results
	params := make([]data_fields.URLParams, 0, len(pairs))

	// Parse each pair into the struct
	for _, pair := range pairs {
		// Split pair by = to separate param and value
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			params = append(params, data_fields.URLParams{
				Param: parts[0],
				Value: parts[1],
			})
		}
	}

	return params
}

func (a *agentService) skipNotification(ctx context.Context, event *data_fields.WebsiteVisitEvent) (bool, SkipReason) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilitiesService.skipNotifications")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	if event.DomainContext == nil {
		return true, SkipNoData, nil
	}

	// don't send if from workspace domain
	isWorkspaceDomain := a.isWorkspaceDomain(ctx, event.DomainContext.Domain)
	if isWorkspaceDomain {
		return true, SkipTenantDomain, nil
	}

	if event.DomainContext.DomainRateLimit == true {
		// determine if notified from this domain
		query := entity.SlackNotificationEvents{
			Tenant:    event.Tenant,
			ChannelID: event.ChannelID,
			Domain:    event.DomainContext.Domain,
		}

		priorNotifications, err := a.services.PostgresRepositories.SlackNotificationEventsRepository.Find(ctx, query)
		if err != nil {
			tracing.TraceErr(span, err)
			return false, SkipNA, nil
		}

		if priorNotifications == nil {
			return false, SkipNA, nil
		}

		// determine how long since last notification
		hoursSinceLastNotification := time.Now().Sub(*priorNotifications.LastNotified).Hours()
		if hoursSinceLastNotification < float64(*event.DomainContext.MinHoursBetweenNotifications) {
			return true, SkipTooRecent, nil
		}

		return false, SkipNA, priorNotifications
	}

	return false, SkipNA, nil
}

func (a *agentService) isWorkspaceDomain(ctx context.Context, domain string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilitiesService.isWorkspaceDomain")
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
