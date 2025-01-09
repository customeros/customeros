package service

import (
	"context"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type SkipReason string

const (
	SkipTenantDomain SkipReason = "tenant domain"
	SkipNoData       SkipReason = "no data"
	SkipTooRecent    SkipReason = "too recent"
	SkipNA           SkipReason = ""
)

func (a *agentService) NotifySlackWebsiteVisitor(ctx context.Context, event *data_fields.WebsiteVisitEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilitiesService.notifySlackWebsiteVisitor")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	skip, _ := a.skipNotification(ctx, event)
	if skip {
		return nil
	}

	err := a.services.SlackService.Notify(ctx, event.Tenant, event.ChannelID, &event.Message)

	return nil
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
