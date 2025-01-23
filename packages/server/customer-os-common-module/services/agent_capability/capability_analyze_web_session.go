package agent_capability

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type AnalyzeWebSessionInput struct {
	SessionID      string
	VisitorID      string
	OrganizationID string
	Domain         string
}

type AnalyzeWebSessionResult struct {
	SessionID         string
	IsNewCompanyVisit bool
	IsNewPersonVisit  bool
	PageViews         []string
	SessionDuration   string
	SlackNotification string
	Hostname          string
	Referrer          string
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

	// validation
	if data.SessionID == "" || data.Domain == "" {
		err := errors.New("Missing required input data")
		tracing.TraceErr(span, err)
		return AnalyzeWebSessionResult{}, err
	}

	// get unique pageviews
	pageViews, err := c.getUniquePageViews(ctx, data.SessionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return AnalyzeWebSessionResult{}, err
	}

	// calculate session duration
	hostname, referrer, sessionDuration, err := c.sessionAnalytics(ctx, data.SessionID)
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

	// build results object
	results := AnalyzeWebSessionResult{
		PageViews:         pageViews,
		SessionDuration:   sessionDuration,
		IsNewCompanyVisit: isNewCompany,
		IsNewPersonVisit:  isNewVisitor,
		Hostname:          hostname,
		Referrer:          referrer,
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
		return []string{}, err
	}
	if session == nil {
		return []string{}, nil
	}

	// Use map to track unique pages
	uniquePageMap := make(map[string]struct{})
	for _, page := range session {
		url := c.cleanUrl(page.Hostname)
		if page.Pathname != "" {
			url = strings.TrimSuffix(fmt.Sprintf("%s/%s", c.cleanUrl(page.Hostname), c.cleanUrl(page.Pathname)), "/")
		}
		uniquePageMap[url] = struct{}{}
	}

	// Convert map keys to slice
	uniquePages := make([]string, 0, len(uniquePageMap))
	for pathname := range uniquePageMap {
		uniquePages = append(uniquePages, pathname)
	}

	sort.Slice(uniquePages, func(i, j int) bool {
		// If lengths are different, sort by length
		if len(uniquePages[i]) != len(uniquePages[j]) {
			return len(uniquePages[i]) < len(uniquePages[j])
		}
		// If lengths are equal, sort alphabetically
		return uniquePages[i] < uniquePages[j]
	})

	return uniquePages, nil
}

func (c *agentCapabilityService) sessionAnalytics(ctx context.Context, sessionID string) (string, string, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.sessionAnalytics")
	defer span.Finish()
	tracing.TagComponentService(span)

	session, err := c.postgresRepositories.WebSessionRepository.FindSession(ctx, postgres_entity.WebSession{
		ID:     sessionID,
		Tenant: common.GetTenantFromContext(ctx),
	}, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", "", "", err
	}
	if session == nil {
		return "", "", "", nil
	}

	hostname := c.cleanUrl(session.Hostname)

	var referrer string
	if session.Referrer == nil {
		referrer = ""
	} else {
		referrer = *session.Referrer
	}
	referrer = c.cleanUrl(referrer)
	if strings.Contains(referrer, "syndicatedsearch.goog") {
		referrer = "google.com"
	}

	sessionDuration, err := c.calculateSessionDuration(ctx, session)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", "", "", err
	}

	return hostname, referrer, sessionDuration, nil
}

func (c *agentCapabilityService) calculateSessionDuration(ctx context.Context, session *postgres_entity.WebSession) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.calculateSessionDuration")
	defer span.Finish()
	tracing.TagComponentService(span)

	if session.EndTime == nil || session.EndTime.IsZero() {
		err := errors.New("Session EndTime not set")
		tracing.TraceErr(span, err)
		return "", err
	}

	duration := session.EndTime.Sub(session.StartTime)
	minutes := duration.Minutes()

	switch {
	case minutes < 1:
		return "less than 1 minute", nil
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
		Tenant:   tenant,
		Domain:   &domain,
		IsActive: false,
	}
	span.LogKV("isActive", "false")

	results, err := c.postgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}
	span.LogKV("recordsReturned", len(results))

	if results == nil || len(results) == 0 {
		return true, nil
	}
	return false, nil
}

func (c *agentCapabilityService) isNewWebsiteVisitor(ctx context.Context, visitorId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.isNewCompanyVisit")
	defer span.Finish()
	tracing.TagComponentService(span)

	tenant := common.GetTenantFromContext(ctx)

	if tenant == "" || visitorId == "" {
		err := errors.New("neither tenant or visitorID can be nil")
		tracing.TraceErr(span, err)
		span.LogKV("tenant", tenant)
		span.LogKV("visitorID", visitorId)
		return false, err
	}

	query := postgres_entity.WebSession{
		Tenant:    tenant,
		VisitorID: visitorId,
		IsActive:  false,
	}

	results, err := c.postgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	if results == nil || len(results) == 0 {
		return true, nil
	}
	return false, nil
}

func (c *agentCapabilityService) buildTimelineMessage(ctx context.Context, sessionID string, analysis AnalyzeWebSessionResult) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.buildTimelineMessage")
	defer span.Finish()
	tracing.TagComponentService(span)

	// Build base message
	var baseMessage string
	switch {
	case analysis.IsNewCompanyVisit && analysis.Referrer != "":
		baseMessage = fmt.Sprintf("A **New visitor** referred by %s browsed for %s", analysis.Referrer, analysis.SessionDuration)
	case analysis.IsNewCompanyVisit && analysis.Referrer == "":
		baseMessage = fmt.Sprintf("A **New visitor** browsed for %s", analysis.SessionDuration)
	case !analysis.IsNewCompanyVisit && analysis.Referrer != "":
		baseMessage = fmt.Sprintf("A **Repeat visitor** referred by %s browsed for %s", analysis.Referrer, analysis.SessionDuration)
	default:
		baseMessage = fmt.Sprintf("A **Repeat visitor** browsed for %s", analysis.SessionDuration)
	}

	// Sort pages by length and alphabetically
	if len(analysis.PageViews) > 0 {
		sort.Slice(analysis.PageViews, func(i, j int) bool {
			if len(analysis.PageViews[i]) != len(analysis.PageViews[j]) {
				return len(analysis.PageViews[i]) < len(analysis.PageViews[j])
			}
			return analysis.PageViews[i] < analysis.PageViews[j]
		})
	}

	// Build full message
	var fullMessage strings.Builder
	fullMessage.WriteString(baseMessage)

	// Only add page views if hostname is present and there are pages to show
	if analysis.Hostname != "" && len(analysis.PageViews) > 0 {
		for _, page := range analysis.PageViews {
			fullUrl := fmt.Sprintf("https://%s", page)
			fullMessage.WriteString(fmt.Sprintf("\n* [%s](%s)", page, fullUrl))
		}
	}

	return fullMessage.String(), nil
}

func (c *agentCapabilityService) cleanUrl(s string) string {
	if s == "" {
		return ""
	}
	clean := strings.Split(s, "?")[0]
	clean = strings.TrimPrefix(clean, "https://")
	clean = strings.TrimPrefix(clean, "http://")
	clean = strings.TrimPrefix(clean, "www.")
	return strings.Trim(clean, "/")
}
