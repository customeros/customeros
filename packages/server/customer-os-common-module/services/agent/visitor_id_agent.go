package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go/log"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type AgentVisitorIDService struct {
	postgresRepositories *postgres_repository.Repositories
	agentCapabilities    *agent_capability.AgentCapabilities
	agentService         interfaces.AgentService
	workspaceService     interfaces.WorkspaceService
}

func NewAgentVisitorIDService(
	postgresRepositories *postgres_repository.Repositories,
	agentCapabilities *agent_capability.AgentCapabilities,
	agentService interfaces.AgentService,
	workspaceService interfaces.WorkspaceService,
) *AgentVisitorIDService {
	return &AgentVisitorIDService{
		postgresRepositories: postgresRepositories,
		agentCapabilities:    agentCapabilities,
		agentService:         agentService,
		workspaceService:     workspaceService,
	}
}

func (a *AgentVisitorIDService) Create(ctx context.Context) (*postgres_entity.Agents, error) {
	return a.agentService.CreateAgent(ctx, enum.AgentVisitorID)
}

func (a *AgentVisitorIDService) Run(ctx context.Context, agentID string, event *data_fields.WebsiteVisitEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.Run")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("agentID", agentID)
	tracing.LogObjectAsJson(span, "event", event)

	err := a.validateWebsiteVisitEvent(event)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid event"))
		return err
	}

	// create execution record
	executionID, err := a.createAgentExecutionRecord(ctx, agentID, event.Type())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to create agent execution record"))
		return err
	}

	// lookup capabilities

	// execute identify visitor capability
	visitorIDResults, err := a.executeVisitorIDCapability(ctx, agentID, executionID, event)
	if err != nil {
		tracing.TraceErr(span, err)
		errMessage := "unable to lookup visitor ID"
		_, err := a.postgresRepositories.AgentExecutionRepository.Update(
			ctx, executionID, nil, &errMessage, false)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to update agent execution record"))
			return err
		}
		return err
	}

	// no result, return early
	if visitorIDResults.Domain == "" {
		_, err := a.postgresRepositories.AgentExecutionRepository.Update(
			ctx,
			executionID,
			utils.NowPtr(),
			nil, // errorMessage
			false,
		)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to update agent execution record"))
			return err
		}
		return nil
	}

	// execute org creation capability
	orgCreationResults, err := a.executeOrgCreationCapability(ctx, agentID, executionID, visitorIDResults.Domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to create organization"))
		errMessage := "unable to create new organization"
		_, err := a.postgresRepositories.AgentExecutionRepository.Update(
			ctx, executionID, nil, &errMessage, false)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to update agent execution record"))
			return err
		}
		return err
	}

	// execute web session analysis capability
	sessionAnalytics, err := a.executeWebSessionAnalysisCapability(
		ctx,
		agentID,
		executionID,
		visitorIDResults.Domain,
		orgCreationResults.OrganizationID,
		event,
	)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to analyze web session"))
		errMessage := "unable to analyze web session"
		_, err := a.postgresRepositories.AgentExecutionRepository.Update(
			ctx, executionID, nil, &errMessage, false)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to update agent execution record"))
			return err
		}
		return err
	}

	// check if slack notification enabled
	slackChannel, err := a.postgresRepositories.SlackChannelNotificationRepository.GetSlackChannel(ctx, "REVEAL-AI-WEBSITE-VISIT")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if slackChannel == nil {
		span.LogKV("result.slackChannel", "not found")
		return nil
	}
	span.LogKV("result.slackChannel", slackChannel.ChannelId)

	// check if notification should be suppressed
	skip, err := a.skipNotification(ctx, visitorIDResults.Domain, 12)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to check if notification should be suppressed"))
		return err
	}
	if skip {
		return nil
	}

	// build slack message
	message, err := a.buildWebVisitorSlackNotification(ctx, orgCreationResults.OrganizationID, visitorIDResults, sessionAnalytics)
	if err != nil || message == nil {
		err = errors.Wrap(err, "unable to build slack notificatoin")
		tracing.TraceErr(span, err)
		errMessage := "unable to send slack notification"
		_, err := a.postgresRepositories.AgentExecutionRepository.Update(
			ctx, executionID, nil, &errMessage, false)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to update agent execution record"))
			return err
		}
		return err
	}

	// execute send notification capability
	err = a.executeSendSlackNotificationCapability(ctx, agentID, executionID, slackChannel.ChannelId, message)
	if err != nil {
		tracing.TraceErr(span, err)
		errMessage := "unable to send slack notification"
		_, err := a.postgresRepositories.AgentExecutionRepository.Update(
			ctx, executionID, nil, &errMessage, false)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return err
	}

	// update agentExecutionRecord
	_, err = a.postgresRepositories.AgentExecutionRepository.Update(
		ctx,
		executionID,
		utils.NowPtr(),
		nil,
		true,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

// TODO check where used
func (a *AgentVisitorIDService) isSlackNotificationEnabled(ctx context.Context) (bool, string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.isSlackNotificationEnabled")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	return true, ""
}

func (a *AgentVisitorIDService) validateWebsiteVisitEvent(event *data_fields.WebsiteVisitEvent) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	if event.SessionID == "" {
		return fmt.Errorf("sessionId cannot be empty")
	}
	if event.Tenant == "" {
		return fmt.Errorf("tenant cannot be empty")
	}
	if event.IPAddress == "" {
		return fmt.Errorf("ipAddress cannot be empty")
	}
	if event.VisitorID == "" {
		return fmt.Errorf("visitorId cannot be empty")
	}

	return nil
}

func (a *AgentVisitorIDService) executeSendSlackNotificationCapability(
	ctx context.Context,
	agentID, executionID, slackChannelID string,
	message *string,
) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.executeSendSlackNotificationCapability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		return err
	}

	executionContainer := dto.CapabilityExecutionContainer{
		AgentID:          agentID,
		AgentExecutionID: executionID,
		Capability:       enum.CapabilitySendSlackNotification,
		InputData: agent_capability.SendSlackNotificationInput{
			Message:   message,
			ChannelID: slackChannelID,
		},
	}

	executionContainer.AgentID = ""

	return nil
}

func (a *AgentVisitorIDService) executeWebSessionAnalysisCapability(
	ctx context.Context,
	agentID, executionID, domain, organizationID string,
	event *data_fields.WebsiteVisitEvent,
) (agent_capability.AnalyzeWebSessionResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.executeWebSessionAnalysisCapability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	executionContainer := dto.CapabilityExecutionContainer{
		AgentID:          agentID,
		AgentExecutionID: executionID,
		Capability:       enum.CapabilityAnalyzeWebSessionIntent,
		InputData: agent_capability.AnalyzeWebSessionInput{
			SessionID:      event.SessionID,
			VisitorID:      event.VisitorID,
			OrganizationID: organizationID,
			Domain:         domain,
		},
	}

	output, ok := executionContainer.OutputData.(agent_capability.AnalyzeWebSessionResult)
	if !ok {
		err := fmt.Errorf("expected agent_capability.AnalyzeWebSessionResult, got %T", executionContainer.OutputData)
		tracing.TraceErr(span, err)
		return agent_capability.AnalyzeWebSessionResult{}, err
	}

	return output, nil
}

func (a *AgentVisitorIDService) executeOrgCreationCapability(ctx context.Context, agentID, executionID, domain string) (agent_capability.CreateOrganizationResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.executeOrgCreationCapability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	source := "WebVisitor ID Agent"
	executionContainer := dto.CapabilityExecutionContainer{
		AgentID:          agentID,
		AgentExecutionID: executionID,
		Capability:       enum.CapabilityCreateOrganization,
		InputData: data_fields.OrganizationFields{
			Source:  &source,
			Domains: []string{domain},
		},
	}

	output, ok := executionContainer.OutputData.(agent_capability.CreateOrganizationResult)
	if !ok {
		err := fmt.Errorf("expected agent_capability.CreateOrganizationResult, got %T", executionContainer.OutputData)
		tracing.TraceErr(span, err)
		return agent_capability.CreateOrganizationResult{}, err
	}

	return output, nil
}

func (a *AgentVisitorIDService) executeVisitorIDCapability(ctx context.Context, agentID, executionID string, event *data_fields.WebsiteVisitEvent) (agent_capability.IdentifyWebsiteVisitorResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.executeVisitorIDCapability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	executionContainer := dto.CapabilityExecutionContainer{
		AgentID:          agentID,
		AgentExecutionID: executionID,
		Capability:       enum.CapabilityIdentifyWebVisitor,
		InputData: agent_capability.IdentifyWebsiteVisitorInput{
			IPAddress: event.IPAddress,
			SessionID: event.SessionID,
		},
	}

	output, ok := executionContainer.OutputData.(agent_capability.IdentifyWebsiteVisitorResult)
	if !ok {
		err := fmt.Errorf("expected agent_capability.IdentifyWebsiteVisitorResult, got %T", executionContainer.OutputData)
		tracing.TraceErr(span, err)
		return agent_capability.IdentifyWebsiteVisitorResult{}, err
	}

	return output, nil
}

func (a *AgentVisitorIDService) createAgentExecutionRecord(ctx context.Context, agentID, triggerEventType string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.createAgentExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("agentID", agentID, "triggerEventType", triggerEventType)

	agent, err := a.postgresRepositories.AgentsRepository.Find(ctx, postgres_entity.Agents{
		ID: agentID,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if agent == nil {
		err := errors.New("agent not found")
		tracing.TraceErr(span, err)
		return "", err
	}

	return a.agentService.CreateAgentExecutionRecord(ctx, *agent, triggerEventType)
}

func (a *AgentVisitorIDService) skipNotification(ctx context.Context, domain string, cooldownInHrs int) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.skipNotification")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()
	span.LogFields(log.String("domain", domain), log.Int("cooldownInHrs", cooldownInHrs))

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		return true, err
	}

	if domain == "" {
		span.LogFields(log.Bool("result.skip", true))
		return true, nil
	}

	// don't send if from workspace domain
	isWorkspaceDomain := a.isWorkspaceDomain(ctx, domain)
	if isWorkspaceDomain {
		span.LogFields(log.Bool("result.skip", true))
		return true, nil
	}

	// determine last notification from this domain
	lastNotification, err := a.postgresRepositories.WebSessionRepository.FindLastNotification(ctx, tenant, domain)
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

func (a *AgentVisitorIDService) isWorkspaceDomain(ctx context.Context, domain string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.isWorkspaceDomain")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	workspaceDomains, err := a.workspaceService.GetWorkspaceDomainsForTenant(ctx)
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

func (a *AgentVisitorIDService) buildWebVisitorSlackNotification(
	ctx context.Context,
	orgID string,
	visitorID agent_capability.IdentifyWebsiteVisitorResult,
	sessionAnalytics agent_capability.AnalyzeWebSessionResult,
) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.buildWebVisitorSlackNotification")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	// get org data from global org table
	globalOrg, err := a.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, visitorID.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if globalOrg == nil {
		err = a.postgresRepositories.GlobalOrganizationWebsiteToProcessRepository.AddWebsiteToProcess(ctx, visitorID.Domain)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	// Build company info section
	var companyLines []string
	primaryDomain := visitorID.Domain
	website := "https://" + primaryDomain
	name := visitorID.Domain
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

	if sessionAnalytics.Referrer != "" {
		companyLines = append(companyLines, fmt.Sprintf("*Source:* <https://%s|%s>", sessionAnalytics.Referrer, sessionAnalytics.Referrer))
	} else {
		companyLines = append(companyLines, "*Source:* Direct")
	}

	companyContent := strings.Join(companyLines, "\n")

	// Build session info section
	var sessionLines []string
	sessionLines = append(sessionLines, fmt.Sprintf("*Session Duration:* %s minutes", sessionAnalytics.SessionDuration))
	sessionLines = append(sessionLines, fmt.Sprintf("*Pages Viewed:* %d", len(sessionAnalytics.PageViews)))
	for _, page := range sessionAnalytics.PageViews {
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
		orgID)

	return &layoutBlocks, nil
}
