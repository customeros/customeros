package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

const MinHoursBetweenNotifications int = 12 // on same domain for a tenant

type URLParams struct {
	Param string `json:"param"`
	Value string `json:"value"`
}

func (a *agentService) VisitorIDAgent(ctx context.Context, eventData *data_fields.WebsiteVisitEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.VisitorIDAgent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

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

	// determine if slack notification is configured

	// handle slack notification
}

func (a *agentService) buildTimelineMessage(ctx context.Context, eventData *data_fields.WebsiteVisitEvent, isNewOrg, isNewVisitor bool) (*string, error) {
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

func (a *agentService) calculateSessionDuration(ctx context.Context, sessionID string) (string, error) {
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

func (a *agentService) getUniquePageViews(ctx context.Context, pageViews []entity.WebTrackerEvents) ([]string, error) {
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

func (a *agentService) identifyIP(ctx context.Context, ipAddress string) (domain, linkedinSlug *string, err error) {
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

	query := entity.WebTrackerEvents{
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

func (a *agentService) isNewWebsiteVisitor(ctx context.Context, tenant, visitorId string) (bool, error) {
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

func (a *agentService) buildWebVisitorSlackNotification(ctx context.Context, eventData *data_fields.WebsiteVisitEvent) (string, error) {
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
