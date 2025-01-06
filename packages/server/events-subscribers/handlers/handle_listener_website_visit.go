package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"
)

const MinHoursBetweenNotifications int = 12 // on same domain for a tenant

func HandleWebsiteVisitorEvent(c context.Context, s *service.Services, sourceEvent commonEnum.FlowListenerEvent, eventData *data_fields.WebsiteVisitEvent) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.HandleMeetingSummaryEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant: eventData.Tenant,
	})

	flows, err := s.WorkflowService.GetFlowsForListenerEvent(ctx, sourceEvent, eventData.Type(), eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flows == nil || len(*flows) == 0 {
		return nil
	}

	var errs error
	for _, flow := range *flows {
		// get next action on the flow
		nextStep, err := s.WorkflowService.GetNextStepInFlow(ctx, flow.ID, &flow.TriggerNodeID)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
			continue
		}

		// send event to agent executioner
		switch nextStep.ToNodeAgent.Agent() {
		case "slack":
			err := publishSlackNotifyEvent(ctx, s, &flow, eventData, sourceEvent)
			if err != nil {
				tracing.TraceErr(span, err)
				errs = multierr.Append(errs, err)
			}

		default:
			err = fmt.Errorf("next flow action not handled for flow %s", flow.ID)
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func publishSlackNotifyEvent(ctx context.Context, s *service.Services, flow *entity.Flow, eventData *data_fields.WebsiteVisitEvent, sourceEvent commonEnum.FlowListenerEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.publishSlackNotifyEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	domains := make([]string, 1)
	domains = append(domains, eventData.Domain)

	// create the org if it does not exist
	orgId, err := s.OrganizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
		Domains: domains,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	eventData.OrganizationID = &orgId

	// build the message
	message, err := buildWebVisitorSlackNotification(ctx, s, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// build the slack.notify event
	slackNotifyEvent, err := buildSlackNotifyEvent(ctx, s, eventData, message)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if slackNotifyEvent == nil {
		err = fmt.Errorf("failed to build slack notify event for %s", eventData.Domain)
		tracing.TraceErr(span, err)
		return err
	}

	// create flow execution record
	flowExecutionRecord, err := createFlowExecutionRecordForWebsiteVisit(ctx, s, flow, eventData, slackNotifyEvent.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// do not publish if flow is configured but not running
	flowExecutionStatus, err := commonEnum.GetFlowExecutionStatus(flowExecutionRecord.Status)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowExecutionStatus != commonEnum.FlowExecutionRunning {
		return nil
	}

	// publish
	flowAgentEvent := dto.FlowAgentEvent{
		FlowExecutionId:  flowExecutionRecord.ID,
		Tenant:           eventData.Tenant,
		ExternalSystemId: enum.SourceReveal,
		SourceEvent:      sourceEvent,
		Name:             commonEnum.AgentSlackNotify,
		DataType:         data_fields.SlackNotifyEventFields{}.Type(),
		Data:             slackNotifyEvent,
	}

	err = s.RabbitMQService.PublishFlowAgentEvent(ctx, flowAgentEvent)
	if err != nil {
		err = fmt.Errorf("failed to publish slack notify event for %s: %w", eventData.Domain, err)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func createFlowExecutionRecordForWebsiteVisit(ctx context.Context, s *service.Services, flow *entity.Flow, eventData *data_fields.WebsiteVisitEvent, eventId string) (*entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.createFlowExecutionRecordForWebsiteVisit")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	flowExecutionRecord, err := s.WorkflowService.BuildAndSaveFlowExecutionRecord(
		ctx, flow.Status, flow.ID, flow.TriggerNodeID, eventId, eventData.Type(),
		commonEnum.AgentSlackNotify.String(), eventData,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return flowExecutionRecord, nil
}

func buildSlackNotifyEvent(ctx context.Context, s *service.Services, eventData *data_fields.WebsiteVisitEvent, message string) (*data_fields.SlackNotifyEventFields, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.buildSlackNotifyEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	channelIds, err := s.PostgresRepositories.SlackChannelNotificationRepository.GetSlackChannels(ctx, eventData.Tenant, "REVEAL-AI")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if len(channelIds) == 0 {
		err = fmt.Errorf("no slack channels found for tenant %s", eventData.Tenant)
		tracing.TraceErr(span, err)
		return nil, err
	}

	channelId := channelIds[0]

	ratelimit := MinHoursBetweenNotifications

	domainContext := data_fields.DomainContext{
		Domain:                       eventData.Domain,
		DomainRateLimit:              true,
		MinHoursBetweenNotifications: &ratelimit,
	}

	event := data_fields.SlackNotifyEventFields{
		Tenant:        eventData.Tenant,
		ChannelID:     channelId.ChannelId,
		Message:       message,
		DomainContext: &domainContext,
	}

	event.ID = event.GenerateID()

	return &event, nil
}

func buildWebVisitorSlackNotification(ctx context.Context, s *service.Services, eventData *data_fields.WebsiteVisitEvent) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.buildWebVisitorSlackNotification")
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
