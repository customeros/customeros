package agent_capability

import (
	"context"
	"fmt"
	"strings"
	"time"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type SendWebVisitorSlackNotificationCapability struct {
	notificationService             interfaces.NotificationService
	workspaceService                interfaces.WorkspaceService
	postgresRepositories            *postgres_repository.Repositories
	sendSlackNotificationCapability *SendSlackNotificationCapability
}

type SendWebVisitorSlackNotificationInput struct {
	Domain          string   `json:"domain"`
	OrganizationID  string   `json:"organizationId"`
	Referrer        string   `json:"referrer"`
	PageViews       []string `json:"pageViews"`
	SessionDuration string   `json:"sessionDuration"`
}

type SendWebVisitorSlackNotificationResult struct {
	Success bool `json:"success"`
	Skipped bool `json:"skipped"`
}

type SendWebVisitorSlackNotificationConfig struct {
	ChannelID     string `json:"channelId"`
	CooldownHours int    `json:"cooldownHours"`
}

func NewSendWebVisitorSlackNotificationCapability(postgresRepositories *postgres_repository.Repositories,
	notificationService interfaces.NotificationService,
	workspaceService interfaces.WorkspaceService,
) *SendWebVisitorSlackNotificationCapability {
	return &SendWebVisitorSlackNotificationCapability{
		postgresRepositories:            postgresRepositories,
		notificationService:             notificationService,
		workspaceService:                workspaceService,
		sendSlackNotificationCapability: NewSendSlackNotificationCapability(notificationService),
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapabilityExecution[SendWebVisitorSlackNotificationInput, SendWebVisitorSlackNotificationResult, SendWebVisitorSlackNotificationConfig] = (*SendWebVisitorSlackNotificationCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                                                                                       = (*SendWebVisitorSlackNotificationCapability)(nil)
)

func (c *SendWebVisitorSlackNotificationCapability) ValidateConfig(config SendWebVisitorSlackNotificationConfig) error {
	if config.ChannelID == "" {
		return errors.New("ChannelID must be set")
	}
	if config.CooldownHours < 0 {
		return errors.New("CooldownHours must be greater than or equal to 0")
	}
	return nil
}

func (c *SendWebVisitorSlackNotificationCapability) ValidateInput(input SendWebVisitorSlackNotificationInput) error {
	if input.Domain == "" {
		return coserrors.ErrCapabilityDomainMissing
	}
	if input.OrganizationID == "" {
		return errors.New("OrganizationID cannot be empty")
	}
	return nil
}

func (c *SendWebVisitorSlackNotificationCapability) GetInput() any {
	return &SendWebVisitorSlackNotificationInput{}
}

func (c *SendWebVisitorSlackNotificationCapability) GetConfig() any {
	return &SendWebVisitorSlackNotificationConfig{}
}

func (c *SendWebVisitorSlackNotificationCapability) GetOutput() any {
	return &SendWebVisitorSlackNotificationResult{}
}

func (c *SendWebVisitorSlackNotificationCapability) Execute(ctx context.Context, data SendWebVisitorSlackNotificationInput, config SendWebVisitorSlackNotificationConfig) (SendWebVisitorSlackNotificationResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendWebVisitorSlackNotificationCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", data)
	tracing.LogObjectAsJson(span, "config", config)

	result := SendWebVisitorSlackNotificationResult{
		Success: false,
		Skipped: false,
	}

	if err := c.ValidateInput(data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return result, err
	}
	if err := c.ValidateConfig(config); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return result, err
	}

	// check if slack notification enabled
	slackChannel, err := c.postgresRepositories.SlackChannelNotificationRepository.GetSlackChannel(ctx, "REVEAL-AI-WEBSITE-VISIT")
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}
	if slackChannel == nil {
		span.LogKV("result.slackChannel", "not found")
		return result, nil
	}
	span.LogKV("result.slackChannel", slackChannel.ChannelId)

	// check if notification should be suppressed
	skip, err := c.skipNotification(ctx, data.Domain, config.CooldownHours)
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}
	if skip {
		result.Success = true
		result.Skipped = true
		return result, nil
	}

	message, err := c.buildWebVisitorSlackNotification(ctx, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}
	if message == nil {
		err = errors.New("Failed to build slack notification message")
		tracing.TraceErr(span, err)
		return result, err
	}

	sendResult, err := c.sendSlackNotificationCapability.Execute(ctx, SendSlackNotificationInput{Message: message}, SendSlackNotificationConfig{ChannelID: config.ChannelID})
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}

	result.Success = sendResult.Success
	tracing.LogObjectAsJson(span, "result", result)
	return result, nil
}

func (c *SendWebVisitorSlackNotificationCapability) skipNotification(ctx context.Context, domain string, cooldownInHrs int) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendWebVisitorSlackNotificationCapability.skipNotification")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()
	span.LogFields(log.String("domain", domain), log.Int("cooldownInHrs", cooldownInHrs))

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		return true, err
	}

	// don't send if from workspace domain
	isWorkspaceDomain := c.isWorkspaceDomain(ctx, domain)
	if isWorkspaceDomain {
		span.LogFields(log.Bool("result.skip", true))
		return true, nil
	}

	// determine last notification from this domain
	lastNotification, err := c.postgresRepositories.WebSessionRepository.FindLastNotification(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, err)
		span.LogFields(log.Bool("result.skip", false))
		return false, nil
	}

	if lastNotification == nil || lastNotification.SentSlackNotification == nil {
		span.LogFields(log.Bool("result.skip", false))
		return false, nil
	}

	// determine how long since last notification
	hoursSinceLastNotification := time.Now().Sub(*lastNotification.SentSlackNotification).Hours()
	if hoursSinceLastNotification < float64(cooldownInHrs) {
		span.LogFields(log.Bool("result.skip", true))
		return true, nil
	}

	span.LogFields(log.Bool("result.skip", false))
	return false, nil
}

func (c *SendWebVisitorSlackNotificationCapability) isWorkspaceDomain(ctx context.Context, domain string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendWebVisitorSlackNotificationCapability.isWorkspaceDomain")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	workspaceDomains, err := c.workspaceService.GetWorkspaceDomainsForTenant(ctx)
	tracing.LogObjectAsJson(span, "workspaceDomains", workspaceDomains)
	if err != nil {
		tracing.TraceErr(span, err)
		return false
	}

	if len(workspaceDomains) == 0 {
		err := errors.New("no workspace domains found for tenant")
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

func (c *SendWebVisitorSlackNotificationCapability) buildWebVisitorSlackNotification(
	ctx context.Context, data SendWebVisitorSlackNotificationInput,
) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendWebVisitorSlackNotificationCapability.buildWebVisitorSlackNotification")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	// get org data from global org table
	globalOrg, err := c.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, data.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if globalOrg == nil {
		err = c.postgresRepositories.GlobalOrganizationWebsiteToProcessRepository.AddWebsiteToProcess(ctx, data.Domain)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	// Build company info section
	var companyLines []string
	primaryDomain := data.Domain
	website := "https://" + primaryDomain
	name := data.Domain
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

	if data.Referrer != "" {
		companyLines = append(companyLines, fmt.Sprintf("*Source:* <https://%s|%s>", data.Referrer, data.Referrer))
	} else {
		companyLines = append(companyLines, "*Source:* Direct")
	}

	companyContent := strings.Join(companyLines, "\n")

	// Build session info section
	var sessionLines []string
	sessionLines = append(sessionLines, fmt.Sprintf("*Session Duration:* %s", data.SessionDuration))
	sessionLines = append(sessionLines, "*Pages Viewed:*")
	for _, page := range data.PageViews {
		sessionLines = append(sessionLines, fmt.Sprintf("• <https://%s|%s>", page, page))
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
		data.OrganizationID)

	return &layoutBlocks, nil
}

// ExecuteUntyped implements the AgentCapabilityUntyped interface.
// It casts the generic input and config to the specific types and delegates to the typed Execute method.
func (c *SendWebVisitorSlackNotificationCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
	typedInput, ok := input.(*SendWebVisitorSlackNotificationInput)
	if !ok || typedInput == nil {
		return nil, fmt.Errorf("invalid input type: expected SendWebVisitorSlackNotificationInput")
	}

	typedConfig, ok := config.(*SendWebVisitorSlackNotificationConfig)
	if !ok || typedConfig == nil {
		return nil, fmt.Errorf("invalid config type: expected SendWebVisitorSlackNotificationConfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
