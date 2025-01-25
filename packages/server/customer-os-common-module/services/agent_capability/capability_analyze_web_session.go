package agent_capability

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type AnalyzeWebSessionCapability struct {
	postgresRepositories *postgres_repository.Repositories
	actionService        interfaces.ActionService
}

func (c *AnalyzeWebSessionCapability) GetInput() any {
	return &AnalyzeWebSessionInput{}
}

func (c *AnalyzeWebSessionCapability) GetConfig() any {
	return &NoConfig{}
}

func (c *AnalyzeWebSessionCapability) GetOutput() any {
	return &AnalyzeWebSessionResult{}
}

func NewAnalyzeWebSessionCapability(
	postgresRepositories *postgres_repository.Repositories,
	actionService interfaces.ActionService,
) *AnalyzeWebSessionCapability {
	return &AnalyzeWebSessionCapability{
		postgresRepositories: postgresRepositories,
		actionService:        actionService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapabilityExecution[AnalyzeWebSessionInput, AnalyzeWebSessionResult, NoConfig] = (*AnalyzeWebSessionCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                              = (*AnalyzeWebSessionCapability)(nil)
)

type AnalyzeWebSessionInput struct {
	SessionID      string `json:"session_id"`
	VisitorID      string `json:"visitor_id"`
	OrganizationID string `json:"organization_id"`
	Domain         string `json:"domain"`
}

type AnalyzeWebSessionResult struct {
	SessionID         string   `json:"session_id"`
	IsNewCompanyVisit bool     `json:"is_new_company_visit"`
	IsNewPersonVisit  bool     `json:"is_new_person_visit"`
	PageViews         []string `json:"page_views"`
	SessionDuration   string   `json:"session_duration"`
	SlackNotification string   `json:"slack_notification"`
	Hostname          string   `json:"hostname"`
	Referrer          string   `json:"referrer"`
}

func (c *AnalyzeWebSessionCapability) Execute(ctx context.Context, data AnalyzeWebSessionInput, config NoConfig) (AnalyzeWebSessionResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalyzeWebSessionCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)

	// validation
	if data.SessionID == "" || data.Domain == "" {
		err := errors.New("Missing required input data")
		tracing.TraceErr(span, err)
		return AnalyzeWebSessionResult{}, err
	}

	// analyze session
	results, err := c.sessionAnalytics(ctx, data.SessionID)
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
	results.IsNewCompanyVisit = isNewCompany

	// determene if a new person visit
	isNewVisitor, err := c.isNewWebsiteVisitor(ctx, data.VisitorID)
	if err != nil {
		tracing.TraceErr(span, err)
		return AnalyzeWebSessionResult{}, err
	}
	results.IsNewPersonVisit = isNewVisitor

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

func (c *AnalyzeWebSessionCapability) writeSessionToTimeline(ctx context.Context, orgID, timelineMessage string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalyzeWebSessionCapability.writeSessionToTimeline")
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

func (c *AnalyzeWebSessionCapability) sessionAnalytics(ctx context.Context, sessionID string) (AnalyzeWebSessionResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalyzeWebSessionCapability.sessionAnalytics")
	defer span.Finish()
	tracing.TagComponentService(span)

	session, err := c.postgresRepositories.WebSessionRepository.FindSession(ctx, postgres_entity.WebSession{
		ID:     sessionID,
		Tenant: common.GetTenantFromContext(ctx),
	}, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return AnalyzeWebSessionResult{}, err
	}
	if session == nil {
		return AnalyzeWebSessionResult{}, nil
	}

	hostname := utils.CleanUrlBasePath(session.Hostname)

	var referrer string
	if session.Referrer == nil {
		referrer = ""
	} else {
		referrer = *session.Referrer
	}
	referrer = utils.CleanUrlBasePath(referrer)
	if strings.Contains(referrer, "syndicatedsearch.goog") {
		referrer = "google.com"
	}

	sessionDuration, err := c.calculateSessionDuration(ctx, session)
	if err != nil {
		tracing.TraceErr(span, err)
		return AnalyzeWebSessionResult{}, err
	}

	results := AnalyzeWebSessionResult{
		SessionID:       sessionID,
		PageViews:       session.UniquePageViews,
		SessionDuration: sessionDuration,
		Hostname:        hostname,
		Referrer:        referrer,
	}

	return results, nil
}

func (c *AnalyzeWebSessionCapability) calculateSessionDuration(ctx context.Context, session *postgres_entity.WebSession) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalyzeWebSessionCapability.calculateSessionDuration")
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

func (c *AnalyzeWebSessionCapability) isNewCompanyVisit(ctx context.Context, domain string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalyzeWebSessionCapability.isNewCompanyVisit")
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

func (c *AnalyzeWebSessionCapability) isNewWebsiteVisitor(ctx context.Context, visitorId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalyzeWebSessionCapability.isNewCompanyVisit")
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

func (c *AnalyzeWebSessionCapability) buildTimelineMessage(ctx context.Context, sessionID string, analysis AnalyzeWebSessionResult) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalyzeWebSessionCapability.buildTimelineMessage")
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

// ExecuteUntyped implements the AgentCapabilityUntyped interface.
// It casts the generic input and config to the specific types and delegates to the typed Execute method.
func (c *AnalyzeWebSessionCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
	typedInput, ok := input.(*AnalyzeWebSessionInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for AnalyzeWebSessionCapability: expected AnalyzeWebSessionInput")
	}

	typedConfig, ok := config.(*NoConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type for AnalyzeWebSessionCapability: expected NoCOnfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
