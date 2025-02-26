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
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
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
	WebSessionID    string   `json:"webSessionId"`
}

type SendWebVisitorSlackNotificationConfig struct {
	ChannelID     SlackChannelIdConfig     `json:"channelId" toml:"slack_channel_id"`
	CooldownHours SlackCooldownHoursConfig `json:"cooldownHours" toml:"notification_cooldown_in_hrs"`
}

type SlackCooldownHoursConfig struct {
	Value int64  `json:"value"`
	Error string `json:"error"`
}

func (c *SendWebVisitorSlackNotificationConfig) Validate() bool {
	isValid := true

	// validate channel ID has slack channel format
	if utils.IsBlank(c.ChannelID.Value) {
		c.ChannelID.Error = "Please add a Slack channel"
		isValid = false
	} else {
		c.ChannelID.Error = ""
	}

	if c.CooldownHours.Value < 0 {
		c.CooldownHours.Error = "Cooldown period cannot be a negative number"
		isValid = false
	} else {
		c.CooldownHours.Error = ""
	}

	return isValid
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
	_ interfaces.AgentCapability[SendWebVisitorSlackNotificationInput, NoOutput, SendWebVisitorSlackNotificationConfig] = (*SendWebVisitorSlackNotificationCapability)(nil)
)

func (c *SendWebVisitorSlackNotificationCapability) Type() enum.AgentCapability {
	return enum.CapabilitySendWebVisitorSlackNotification
}

func (c *SendWebVisitorSlackNotificationCapability) Name() string {
	return "Send Slack notification (optional)"
}

func (c *SendWebVisitorSlackNotificationCapability) NewInput() SendWebVisitorSlackNotificationInput {
	return SendWebVisitorSlackNotificationInput{}
}

func (c *SendWebVisitorSlackNotificationCapability) NewConfig() SendWebVisitorSlackNotificationConfig {
	return SendWebVisitorSlackNotificationConfig{}
}

func (c *SendWebVisitorSlackNotificationCapability) DefaultConfig() any {
	config := c.NewConfig()
	config.CooldownHours.Value = 12
	return &config
}

func (c *SendWebVisitorSlackNotificationCapability) DefaultActive() bool {
	return false
}

func (c *SendWebVisitorSlackNotificationCapability) ValidateConfig(config SendWebVisitorSlackNotificationConfig) error {
	if config.ChannelID.Value == "" {
		return errors.New("ChannelID must be set")
	}
	if config.CooldownHours.Value < 0 {
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
	if input.WebSessionID == "" {
		return errors.New("WebSessionID cannot be empty")
	}
	return nil
}

func (c *SendWebVisitorSlackNotificationCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[SendWebVisitorSlackNotificationInput, SendWebVisitorSlackNotificationConfig]) (enum.CapabilityExecutionStatus, NoOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendWebVisitorSlackNotificationCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "executionContainer", executionContainer)

	result := NoOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	// check if notification should be suppressed
	skip, err := c.skipNotification(ctx, executionContainer.InputData.Domain, executionContainer.ConfigData.CooldownHours.Value)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionRetry, result, err
	}
	if skip {
		return enum.CapabilityExecutionCompleted, result, nil
	}

	message, err := c.buildWebVisitorSlackNotification(ctx, executionContainer.InputData)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionRetry, result, err
	}
	if message == nil || *message == "" {
		err = errors.New("Failed to build slack notification message")
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionRetry, result, err
	}

	newExecutionContainer := interfaces.TypedExecutionContainer[SendSlackNotificationInput, SendSlackNotificationConfig]{
		AgentExecutionID: executionContainer.AgentExecutionID,
		InputData: SendSlackNotificationInput{
			Message: *message,
		},
		ConfigData: SendSlackNotificationConfig{
			ChannelID: executionContainer.ConfigData.ChannelID,
		},
	}

	status, output, err := c.sendSlackNotificationCapability.Execute(ctx, newExecutionContainer)
	if err != nil {
		tracing.TraceErr(span, err)
		return status, output, err
	}
	if status == enum.CapabilityExecutionCompleted {
		err = c.postgresRepositories.WebSessionRepository.UpdateSessionWithSlackSentAt(ctx, executionContainer.InputData.WebSessionID)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	return status, output, nil
}

func (c *SendWebVisitorSlackNotificationCapability) skipNotification(ctx context.Context, domain string, cooldownInHrs int64) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendWebVisitorSlackNotificationCapability.skipNotification")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()
	span.LogFields(log.String("domain", domain), log.Int64("cooldownInHrs", cooldownInHrs))

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		return true, err
	}

	// don't send if from workspace domain
	isWorkspaceDomain, err := c.workspaceService.IsWorkspaceDomain(ctx, domain)
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
