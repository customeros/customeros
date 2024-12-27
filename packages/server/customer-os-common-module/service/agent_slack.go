package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
)

type SkipReason string

const (
	SkipTenantDomain SkipReason = "tenant domain"
	SkipNoData       SkipReason = "no data"
	SkipTooRecent    SkipReason = "too recent"
	SkipNA           SkipReason = ""
)

func (a *agentService) SlackAgent(ctx context.Context, event *dto.FlowAgentEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.SlackAgent")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	switch event.DataType {
	case "SlackNotifyEventFields":
		eventData, ok := event.Data.(*data_fields.SlackNotifyEventFields)
		if !ok {
			return fmt.Errorf("failed to cast to SlackNotifyEventFields, got type: %T", event.Data)
		}
		if eventData == nil {
			return fmt.Errorf("SlackNotifyEventFields is nil")
		}

		return a.slackBotNotification(ctx, eventData, event.FlowExecutionId)

	default:
		return errors.New("Unsupported event")
	}
}

func (a *agentService) slackBotNotification(ctx context.Context, event *data_fields.SlackNotifyEventFields, flowExecutionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.slackBotNotification")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	skip, _, existingNotificationEvent := a.skipNotification(ctx, event)
	if skip {
		return nil
	}

	err := a.services.SlackService.Notify(ctx, event.Tenant, event.ChannelID, &event.Message)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// save notification record
	record := entity.SlackNotificationEvents{
		Tenant:    event.Tenant,
		ChannelID: event.ChannelID,
		Domain:    event.DomainContext.Domain,
	}

	if existingNotificationEvent == nil {
		_, err := a.services.PostgresRepositories.SlackNotificationEventsRepository.Create(ctx, record)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return nil
	}

	record.ID = existingNotificationEvent.ID
	_, err = a.services.PostgresRepositories.SlackNotificationEventsRepository.Update(ctx, record)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (a *agentService) skipNotification(ctx context.Context, event *data_fields.SlackNotifyEventFields) (bool, SkipReason, *entity.SlackNotificationEvents) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.skipNotifications")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.isWorkspaceDomain")
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
