package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (s *trackingService) NotifyOnSlack(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "TrackingService.NotifyOnSlack")
	defer span.Finish()

	if s.cfg.SlackBotApiKey == "" {
		span.LogFields(log.String("skip", "no slack bot api key"))
		return
	}

	limit := 100
	notifyOnSlackRecords, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetForSlackNotifications(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get records for slack notifications"))
		s.services.Logger.Errorf("failed to get records for slack notifications: %s", err.Error())
		return
	}

	if len(notifyOnSlackRecords) == 0 {
		span.LogFields(log.String("skip", "no records to notify"))
		return
	}

	for _, r := range notifyOnSlackRecords {
		err := s.notifyOnSlack(ctx, r)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to notify on slack"))
		}
	}

	return
}

func (s *trackingService) notifyOnSlack(c context.Context, r *entity.Tracking) error {
	span, ctx := opentracing.StartSpanFromContext(c, "TrackingService.notifyOnSlack")
	defer span.Finish()

	record, err := s.getAndUpdateTrackingRecord(ctx, span, r.ID)
	if err != nil {
		return err
	}

	if record.Notified || record.OrganizationId == nil {
		return nil
	}

	globalOrg, err := s.getEnrichedOrganizationDetails(ctx, span, record)
	if err != nil {
		return err
	}

	if s.shouldSkipNotification(ctx, span, record) {
		return nil
	}

	slackBlock := s.buildSlackNotification(record, globalOrg)
	return s.sendNotifications(ctx, span, record, slackBlock)
}

func (s *trackingService) getAndUpdateTrackingRecord(ctx context.Context, span opentracing.Span, id string) (*entity.Tracking, error) {
	record, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get tracking record"))
		return nil, err
	}

	if record.Tenant != "" {
		tracing.TagTenant(span, record.Tenant)
	}

	err = s.services.CommonServices.PostgresRepositories.TrackingRepository.IncrementNotificationTry(ctx, record.ID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to increment notification try"))
	}

	return record, nil
}

func (s *trackingService) getEnrichedOrganizationDetails(ctx context.Context, span opentracing.Span, record *entity.Tracking) (*entity.GlobalOrganization, error) {
	enrichDetails, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(ctx, record.IP)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	_, primaryDomain := domaincheck.PrimaryDomainCheck(*enrichDetails.CompanyDomain)
	return s.services.CommonServices.PostgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
}

func (s *trackingService) shouldSkipNotification(ctx context.Context, span opentracing.Span, record *entity.Tracking) bool {
	// Skip if domain matches workspace
	if record.OrganizationDomain != nil && *record.OrganizationDomain != "" {
		if s.isWorkspaceDomain(ctx, span, record) {
			return true
		}

		// Skip if recently notified
		notificationSent, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.WasNotifiedRecently(ctx, *record.OrganizationDomain, 24)
		if err != nil {
			tracing.TraceErr(span, err)
			return true
		}
		return notificationSent
	}
	return false
}

func (s *trackingService) isWorkspaceDomain(ctx context.Context, span opentracing.Span, record *entity.Tracking) bool {
	workspaceNodeList, err := s.services.CommonServices.Neo4jRepositories.WorkspaceReadRepository.GetAllForTenant(ctx, record.Tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get workspace nodes"))
		return true
	}

	for _, workspaceNode := range workspaceNodeList {
		props := utils.GetPropsFromNode(*workspaceNode)
		if utils.GetStringPropOrEmpty(props, "name") == *record.OrganizationDomain {
			err := s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsNotified(ctx, record.ID)
			if err != nil {
				tracing.TraceErr(span, err)
			}
			span.LogFields(log.String("skip", "workspace is the same as organization domain"))
			return true
		}
	}
	return false
}

func (s *trackingService) buildSlackNotification(record *entity.Tracking, globalOrg *entity.GlobalOrganization) string {
	// Build the text content for the section based on available data
	var contentLines []string
	contentLines = append(contentLines, fmt.Sprintf("<%s|*%s*> ", globalOrg.Website, globalOrg.Name))

	if globalOrg.Description != "" {
		contentLines = append(contentLines, fmt.Sprintf("%s \n", globalOrg.Description))
	}

	// Add optional fields only if they're not empty
	if globalOrg.PrimaryDomain != "" && globalOrg.Website != "" {
		contentLines = append(contentLines, fmt.Sprintf("*Website:* <%s|%s> ", globalOrg.Website, globalOrg.PrimaryDomain))
	}

	if globalOrg.LinkedInUrl != "" && globalOrg.LinkedInAlias != "" {
		contentLines = append(contentLines, fmt.Sprintf("*LinkedIn:* <%s|/%s> ", globalOrg.LinkedInUrl, globalOrg.LinkedInAlias))
	}

	// Only add location if both city and country are available
	if globalOrg.City != "" && globalOrg.CountryA2 != "" {
		contentLines = append(contentLines, fmt.Sprintf("*Location:* %s, %s ", globalOrg.City, globalOrg.CountryA2))
	}
	// Add source/referrer only if it exists
	if record.Referrer != "" {
		referrer := strings.TrimPrefix(record.Referrer, "https://")
		referrer = strings.TrimPrefix(referrer, "http://")
		referrer = strings.TrimPrefix(referrer, "www.")
		referrer = strings.Trim(referrer, "/")
		contentLines = append(contentLines, fmt.Sprintf("*Source:* <%s|%s> ", record.Referrer, referrer))
	} else {
		contentLines = append(contentLines, "*Source:* Direct ")
	}

	// Join the lines with newlines
	sectionContent := strings.Join(contentLines, "\n")

	// Create the notification blocks based on logo availability
	var layoutBlocks string
	if globalOrg.LogoUrl != "" {
		// With logo - in same section as content
		layoutBlocks = fmt.Sprintf(`[
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
				},
				"accessory": {
					"type": "image",
					"image_url": "%s",
					"alt_text": "%s logo"
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
		]`, globalOrg.Name, sectionContent, globalOrg.LogoUrl, globalOrg.Name, *record.OrganizationId)
	} else {
		// Without logo - use simple layout
		layoutBlocks = fmt.Sprintf(`[
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
		]`, globalOrg.Name, sectionContent, *record.OrganizationId)
	}

	return layoutBlocks
}

func (s *trackingService) sendNotifications(ctx context.Context, span opentracing.Span, record *entity.Tracking, slackBlock string) error {
	slackChannels, err := s.services.CommonServices.PostgresRepositories.SlackChannelNotificationRepository.GetSlackChannels(ctx, record.Tenant, "REVEAL-AI")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, slackChannel := range slackChannels {
		if slackChannel.CreatedAt.After(record.CreatedAt) {
			err := s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsNotified(ctx, record.ID)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
			continue
		}

		if err := s.sendSlackMessage(ctx, slackChannel.Tenant, slackChannel.ChannelId, slackBlock); err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if err := s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsNotified(ctx, record.ID); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}
	return nil
}

func (s *trackingService) sendSlackMessage(ctx context.Context, tenant, channel, blocks string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingService.sendSlackMessage")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("channel", channel))

	// Create HTTP client
	client := &http.Client{}

	requestBody := map[string]interface{}{
		"channel":      channel,
		"unfurl_links": false,
		"unfurl_media": false,
		"blocks":       blocks,
	}

	// Marshal the request body
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request body"))
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	span.LogFields(log.String("request.body", string(requestBodyBytes)))

	// Create POST request
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create POST request"))
		return fmt.Errorf("failed to create POST request: %v", err)
	}

	botApiKey := ""
	// prepare bot key
	slackSettings, err := s.services.CommonServices.PostgresRepositories.SlackSettingsRepository.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get slack settings"))
	}
	if slackSettings == nil {
		span.LogFields(log.String("skip", "slack settings not found"))
		s.services.Logger.Warnf("slack settings not found for tenant %s", tenant)
		return nil
	} else {
		botApiKey = slackSettings.AccessToken
	}

	// display last first 8 and last 3 chars
	maskedBotApiKey := ""
	if len(botApiKey) > 11 {
		maskedBotApiKey = botApiKey[:8] + "..." + botApiKey[len(botApiKey)-3:]
	}
	span.LogFields(log.String("bot.api.key", maskedBotApiKey))

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+botApiKey)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform POST request"))
		return fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		return fmt.Errorf("failed to read response body: %v", err)
	}

	span.LogFields(log.String("response.body", string(responseBody)))

	return nil
}
