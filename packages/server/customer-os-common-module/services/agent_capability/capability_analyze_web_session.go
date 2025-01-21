package agent_capability

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
)

type AnalyzeWebSessionInput struct {
	SessionID      string
	VisitorID      string
	OrganizationID string
	Domain         string
}

type AnalyzeWebSessionResult struct {
	IsNewCompanyVisit bool
	IsNewPersonVisit  bool
	PageViews         []string
	SessionDuration   string
	SlackNotification string
}

func (c *agentCapabilityService) handleAnalyzeWebSessionExecution(ctx context.Context, executionContainer *dto.CapabilityExecutionContainer) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.handleAnalyzeWebSessionExecution")
	defer span.Finish()
	tracing.TagComponentService(span)

	input, ok := executionContainer.InputData.(AnalyzeWebSessionInput)
	if !ok {
		err := fmt.Errorf("expected AnalyzeWebSessionInput, got %T", executionContainer.InputData)
		tracing.TraceErr(span, err)
		return err
	}
	result, err := c.executeWebSessionAnalysis(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	executionContainer.OutputData = result
	span.LogKV(
		"event", "capability_executed",
		"output_type", fmt.Sprintf("%T", result),
		"is_new_company", result.IsNewCompanyVisit,
		"is_new_person", result.IsNewPersonVisit,
		"session_duration", result.SessionDuration,
	)
	return nil
}

func (c *agentCapabilityService) executeWebSessionAnalysis(ctx context.Context, data AnalyzeWebSessionInput) (AnalyzeWebSessionResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.executeWebSessionAnalysis")
	defer span.Finish()
	tracing.TagComponentService(span)

	// get unique pageviews
	pageViews, err := c.getUniquePageViews(ctx, data.SessionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return AnalyzeWebSessionResult{}, err
	}

	// calculate session duration
	sessionDuration, err := c.calculateSessionDuration(ctx, data.SessionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return AnalyzeWebSessionResult{}, err
	}

	// determine if a new company visit
	isNewCompany, err := c.isNewCompanyVisit(ctx, data.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return AnalyzeWebSessionResult{}, err
	}

	// determene if a new person visit
	isNewVisitor, err := c.isNewWebsiteVisitor(ctx, data.VisitorID)
	if err != nil {
		tracing.TraceErr(span, err)
		return AnalyzeWebSessionResult{}, err
	}

	//build results object
	results := AnalyzeWebSessionResult{
		PageViews:         pageViews,
		SessionDuration:   sessionDuration,
		IsNewCompanyVisit: isNewCompany,
		IsNewPersonVisit:  isNewVisitor,
	}

	// build timeline event
	timelineMessage, err := c.buildTimelineMessage(ctx, data.SessionID, results)
	if err != nil {
		tracing.TraceErr(span, err)
		return results, err
	}

	// write event to timeline
	err = c.writeSessionToTimeline(ctx, data.OrganizationID, timelineMessage)
	if err != nil {
		tracing.TraceErr(span, err)
		return results, err
	}

	// build slack notification

	return results, nil
}

func (c *agentCapabilityService) writeSessionToTimeline(ctx context.Context, orgID, timelineMessage string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.writeSessionToTimeline")
	defer span.Finish()
	tracing.TagComponentService(span)

	actionType := enum.ActionGeneric
	metadata := "web_session"

	action := data_fields.ActionFields{
		ActionType: &actionType,
		CreatedAt:  utils.NowPtr(),
		Content:    &timelineMessage,
		Metadata:   &metadata,
	}

	_, err := c.actionService.CreateActionForOrganization(ctx, nil, orgID, action)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (c *agentCapabilityService) getUniquePageViews(ctx context.Context, sessionID string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.getUniquePageViews")
	defer span.Finish()
	tracing.TagComponentService(span)

	session, err := c.postgresRepositories.WebTrackerEventsRepository.FindAll(ctx, postgres_entity.WebTrackerEvents{
		SessionID: sessionID,
		Tenant:    common.GetTenantFromContext(ctx),
	}, nil)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}

	// Use map to track unique pages
	uniquePageMap := make(map[string]struct{})
	for _, page := range session {
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

func (c *agentCapabilityService) calculateSessionDuration(ctx context.Context, sessionID string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.calculateSessionDuration")
	defer span.Finish()
	tracing.TagComponentService(span)

	session, err := c.postgresRepositories.WebSessionRepository.FindSession(ctx, postgres_entity.WebSession{
		ID:     sessionID,
		Tenant: common.GetTenantFromContext(ctx),
	}, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if session == nil {
		return "", err
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

func (c *agentCapabilityService) isNewCompanyVisit(ctx context.Context, domain string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.isNewCompanyVisit")
	defer span.Finish()
	tracing.TagComponentService(span)

	tenant := common.GetTenantFromContext(ctx)

	if tenant == "" || domain == "" {
		err := errors.New("neither tenant or domain can be nil")
		tracing.TraceErr(span, err)
		span.LogKV("tenant", tenant)
		span.LogKV("domain", domain)
		return false, err
	}

	query := postgres_entity.WebSession{
		Tenant: tenant,
		Domain: &domain,
	}

	results, err := c.postgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	if len(results) == 0 {
		return true, nil
	}
	return false, nil
}

func (c *agentCapabilityService) isNewWebsiteVisitor(ctx context.Context, visitorId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.isNewCompanyVisit")
	defer span.Finish()
	tracing.TagComponentService(span)

	query := postgres_entity.WebSession{
		Tenant:    common.GetTenantFromContext(ctx),
		VisitorID: visitorId,
	}

	results, err := c.postgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	if len(results) == 0 {
		return true, nil
	}
	return false, nil
}

func (c *agentCapabilityService) buildTimelineMessage(ctx context.Context, sessionID string, analysis AnalyzeWebSessionResult) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.buildTimelineMessage")
	defer span.Finish()
	tracing.TagComponentService(span)

	var baseMessage string
	switch {
	case analysis.IsNewCompanyVisit:
		baseMessage = fmt.Sprintf("First Visit: %s", analysis.SessionDuration)
	case analysis.IsNewPersonVisit:
		baseMessage = fmt.Sprintf("New Visitor: %s", analysis.SessionDuration)
	case !analysis.IsNewPersonVisit:
		baseMessage = fmt.Sprintf("Repeat Visitor: %s", analysis.SessionDuration)
	default:
		baseMessage = fmt.Sprintf("Web Visitor: %s", analysis.SessionDuration)
	}

	// Append pages with indentation
	var fullMessage strings.Builder
	fullMessage.WriteString(baseMessage)
	for _, page := range analysis.PageViews {
		fullMessage.WriteString(fmt.Sprintf("\n    • %s", page))
	}

	timelineMessage := fullMessage.String()
	return timelineMessage, nil
}

//
// func (a *agentVisitorIDService) buildWebVisitorSlackNotification(ctx context.Context, session *postgres_entity.WebSession, orgID string) (*string, error) {
// 	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.buildWebVisitorSlackNotification")
// 	defer span.Finish()
// 	tracing.SetDefaultListenerSpanTags(ctx, span)
//
// 	// get org data from global org table
// 	globalOrg, err := a.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, *session.Domain)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return nil, err
// 	}
//
// 	if globalOrg == nil {
// 		err = a.postgresRepositories.GlobalOrganizationWebsiteToProcessRepository.AddWebsiteToProcess(ctx, *session.Domain)
// 		if err != nil {
// 			tracing.TraceErr(span, err)
// 		}
// 	}
//
// 	// Build company info section
// 	var companyLines []string
// 	primaryDomain := *session.Domain
// 	if globalOrg != nil {
// 		primaryDomain = globalOrg.PrimaryDomain
// 	}
// 	website := "https://" + primaryDomain
// 	name := *session.Domain
// 	if globalOrg != nil {
// 		name = globalOrg.Name
// 	}
//
// 	companyLines = append(companyLines, fmt.Sprintf("<%s|*%s*>", website, name))
//
// 	if globalOrg != nil && globalOrg.Description != "" {
// 		companyLines = append(companyLines, fmt.Sprintf("%s", globalOrg.Description))
// 	}
//
// 	if website != "" {
// 		companyLines = append(companyLines, fmt.Sprintf("*Website:* <%s|%s>", website, primaryDomain))
// 	}
//
// 	if globalOrg != nil && globalOrg.LinkedInUrl != "" && globalOrg.LinkedInAlias != "" {
// 		companyLines = append(companyLines, fmt.Sprintf("*LinkedIn:* <%s|/%s>", globalOrg.LinkedInUrl, globalOrg.LinkedInAlias))
// 	}
//
// 	if globalOrg != nil && globalOrg.City != "" && globalOrg.CountryA2 != "" {
// 		companyLines = append(companyLines, fmt.Sprintf("*Location:* %s, %s", globalOrg.City, globalOrg.CountryA2))
// 	}
//
// 	if session.Referrer != nil {
// 		referrer := strings.TrimPrefix(*session.Referrer, "https://")
// 		referrer = strings.TrimPrefix(referrer, "http://")
// 		referrer = strings.TrimPrefix(referrer, "www.")
// 		referrer = strings.Trim(referrer, "/")
// 		companyLines = append(companyLines, fmt.Sprintf("*Source:* <%s|%s>", *session.Referrer, referrer))
// 	} else {
// 		companyLines = append(companyLines, "*Source:* Direct")
// 	}
//
// 	companyContent := strings.Join(companyLines, "\n")
//
// 	// Build session info section
// 	var sessionLines []string
//
// 	if session.EndTime != nil && !session.EndTime.IsZero() {
// 		duration, err := a.calculateSessionDuration(ctx, session)
// 		if err != nil {
// 			tracing.TraceErr(span, err)
// 			return nil, err
// 		}
//
// 		sessionLines = append(sessionLines, fmt.Sprintf("*Session Duration:* %s minutes", duration))
// 	}
//
// 	// Get page views
// 	query := postgres_entity.WebTrackerEvents{
// 		SessionID: session.ID,
// 		EventType: enum.WebTrackerPageView.String(),
// 	}
// 	pageViews, err := a.postgresRepositories.WebTrackerEventsRepository.FindAll(ctx, query, nil)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return nil, err
// 	}
//
// 	uniquePages, err := a.getUniquePageViews(ctx, pageViews)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return nil, err
// 	}
//
// 	sessionLines = append(sessionLines, fmt.Sprintf("*Pages Viewed:* %d", len(uniquePages)))
// 	for _, page := range uniquePages {
// 		sessionLines = append(sessionLines, fmt.Sprintf("• <%s%s|%s>", website, page, page))
// 	}
//
// 	sessionContent := strings.Join(sessionLines, "\n")
//
// 	// Handle logo accessory
// 	var logoAccessory string
// 	if globalOrg != nil && globalOrg.LogoUrl != "" {
// 		logoAccessory = fmt.Sprintf(`,
//            "accessory": {
//                "type": "image",
//                "image_url": "%s",
//                "alt_text": "%s logo"
//            }`, globalOrg.LogoUrl, name)
// 	}
//
// 	// Build the final layout
// 	layoutBlocks := fmt.Sprintf(`[
//        {
//            "type": "header",
//            "text": {
//                "type": "plain_text",
//                "text": "A visitor from %s is on your website",
//                "emoji": true
//            }
//        },
//        {
//            "type": "divider"
//        },
//        {
//            "type": "section",
//            "text": {
//                "type": "mrkdwn",
//                "text": "%s"
//            }%s
//        },
//        {
//            "type": "divider"
//        },
//        {
//            "type": "section",
//            "text": {
//                "type": "mrkdwn",
//                "text": "%s"
//            }
//        },
//        {
//            "type": "divider"
//        },
//        {
//            "type": "actions",
//            "elements": [
//                {
//                    "type": "button",
//                    "text": {
//                        "type": "plain_text",
//                        "text": "View in CustomerOS"
//                    },
//                    "url": "https://app.customeros.ai/organization/%s?tab=about",
//                    "value": "click_me_123",
//                    "action_id": "actionId-0"
//                }
//            ]
//        }
//    ]`,
// 		name,
// 		companyContent,
// 		logoAccessory,
// 		sessionContent,
// 		orgID)
//
// 	return &layoutBlocks, nil
// }
