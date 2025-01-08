package service

import (
	"context"

	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
)

const MinHoursBetweenNotifications int = 12 // on same domain for a tenant

func (a *agentService) VisitorIDAgent(ctx context.Context, eventData *data_fields.WebsiteVisitEvent) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.VisitorIDAgent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	// try to deanonymize the IP address
	domain, linkedinSlug, err := a.identifyIP(ctx, *eventData.IPAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	if domain == nil {
		return
	}

	// Create org

	// update tracker table with ID data
	query := entity.TrackerEvents{
		IP:           *eventData.IPAddress,
		Domain:       domain,
		LinkedinSlug: linkedinSlug,
	}

	_, err = a.services.PostgresRepositories.TrackerEventsRepository.Update(ctx, query)

	// determine if new org
	isNewOrg, err := a.isNewCompanyVisit(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed isNewCompany lookup"))
	}

	// determine if new person
	isNewVisitor, err := a.isNewWebsiteVisitor(ctx, services, *tenant, trackerData.VisitorID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed isNewVisitor lookup"))
	}

	// determine if slack notification is configured

	// handle slack notification

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

	query := entity.TrackerEvents{
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
