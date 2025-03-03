package agent_capability

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type AnalyzeWebSessionCapability struct {
	events               *events.EventsService
	postgresRepositories *postgres_repository.Repositories
	actionService        interfaces.ActionService
	webscraperService    interfaces.WebscraperService
}

func NewAnalyzeWebSessionCapability(
	events *events.EventsService,
	postgresRepositories *postgres_repository.Repositories,
	actionService interfaces.ActionService,
	webscraperService interfaces.WebscraperService,
) *AnalyzeWebSessionCapability {
	return &AnalyzeWebSessionCapability{
		events:               events,
		postgresRepositories: postgresRepositories,
		actionService:        actionService,
		webscraperService:    webscraperService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[AnalyzeWebSessionInput, AnalyzeWebSessionOutput, postgres_entity.NoConfig] = (*AnalyzeWebSessionCapability)(nil)
)

func (c *AnalyzeWebSessionCapability) Type() enum.AgentCapability {
	return enum.CapabilityAnalyzeWebSessionIntent
}

func (c *AnalyzeWebSessionCapability) Name() string {
	return "Analyze behavior for intent signals"
}

func (c *AnalyzeWebSessionCapability) NewInput() AnalyzeWebSessionInput {
	return AnalyzeWebSessionInput{}
}

func (c *AnalyzeWebSessionCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *AnalyzeWebSessionCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *AnalyzeWebSessionCapability) DefaultActive() bool {
	return true
}

func (c *AnalyzeWebSessionCapability) ValidateInput(data AnalyzeWebSessionInput) error {
	if data.WebSessionID == "" {
		return errors.New("missing required input data: SessionID")
	}
	if data.Domain == "" {
		return coserrors.ErrCapabilityDomainMissing
	}
	if data.OrganizationID == "" {
		return errors.New("missing required input data: OrganizationID")
	}
	return nil
}

func (c *AnalyzeWebSessionCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

type AnalyzeWebSessionInput struct {
	WebSessionID   string `json:"webSessionId"`
	OrganizationID string `json:"organizationId"`
	Domain         string `json:"domain"`
}

type AnalyzeWebSessionOutput struct {
	WebSessionID      string   `json:"webSessionId"`
	IsNewCompanyVisit bool     `json:"isNewCompanyVisit"`
	IsNewPersonVisit  bool     `json:"isNewPersonVisit"`
	PageViews         []string `json:"pageViews"`
	SessionDuration   string   `json:"sessionDuration"`
	SlackNotification string   `json:"slackNotification"`
	Hostname          string   `json:"hostname"`
	Referrer          string   `json:"referrer"`
}

func (c *AnalyzeWebSessionCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[AnalyzeWebSessionInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, AnalyzeWebSessionOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalyzeWebSessionCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := AnalyzeWebSessionOutput{}

	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}

	// update session with organization id
	err := c.postgresRepositories.WebSessionRepository.SetOrganizationId(ctx, executionContainer.InputData.WebSessionID, executionContainer.InputData.OrganizationID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	// get session details
	session, err := c.postgresRepositories.WebSessionRepository.FindSession(ctx, postgres_entity.WebSession{
		ID: executionContainer.InputData.WebSessionID,
	}, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}
	if session == nil {
		err := errors.New("cannot identify session")
		span.LogKV("sessionId", executionContainer.InputData.WebSessionID)
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}
	if len(session.UniquePageViews) == 0 {
		err := errors.New("no page views to analyze")
		span.LogKV("sessionId", executionContainer.InputData.WebSessionID)
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	// analyze session
	result, err = c.sessionAnalytics(ctx, executionContainer.InputData.WebSessionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionRetry, result, err
	}

	// determine if a new company visit
	isNewCompany, err := c.isNewCompanyVisit(ctx, executionContainer.InputData.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionRetry, result, err
	}
	result.IsNewCompanyVisit = isNewCompany

	// build timeline event
	timelineMessage, err := c.buildTimelineMessage(ctx, executionContainer.InputData.WebSessionID, result)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionRetry, result, err
	}

	// write event to timeline
	err = c.writeSessionToTimeline(ctx, executionContainer.InputData.OrganizationID, timelineMessage)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}

type Session struct {
	PageVisits []PageVisit `json:"pageVisits"`
}

type PageVisit struct {
	SessionID    string                    `json:"sessionId"`
	IPAddress    string                    `json:"ipAddress"`
	Url          string                    `json:"url"`
	StageSignal  enum.CustomerJourneyStage `json:"stageSignal"`
	PageCategory enum.WebpageCategory      `json:"pageCategory"`
	Topics       []string                  `json:"topics"`

	EntryTimestamp time.Time `json:"entryTimestamp"`
	ExitTimestamp  time.Time `json:"exitTimestamp"`
	ClickCount     int64     `json:"clickCount"`
	Referer        string    `json:"referer"`
	SecondsOnPage  int64     `json:"secondsOnPage"`

	VisitorEmail     string `json:"visitorEmail"`
	VisitorEmailType string `json:"visitorEmailType"`
	VisitorDomain    string `json:"visitorDomain"`
}

const (
	BusinessEmail = "business"
	PersonalEmail = "personal"
)

func (c *AnalyzeWebSessionCapability) processPageVisit(ctx context.Context, sessisonId, page, domain string) (PageVisit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalyzeWebSessionCapability.processPageVisit")
	defer span.Finish()
	tracing.TagComponentService(span)

	url := page
	if !strings.HasPrefix(page, "http") {
		url = fmt.Sprintf("https://%s", page)
	}

	visit := PageVisit{
		SessionID:     sessisonId,
		Url:           url,
		VisitorDomain: domain,
	}

	// get session events
	events, err := c.postgresRepositories.WebTrackerEventsRepository.FindEventsForPageVisit(ctx, sessisonId, page)
	if err != nil {
		tracing.TraceErr(span, err)
		return visit, err
	}

	for _, event := range events {
		if visit.IPAddress == "" {
			visit.IPAddress = event.IP
		}

		if visit.Referer == "" && event.Referrer != "" {
			visit.Referer = event.Referer()
		}

		switch event.EventType {
		case "page_view":
			if visit.EntryTimestamp.IsZero() {
				visit.EntryTimestamp = event.Timestamp
			} else if event.Timestamp.Before(visit.EntryTimestamp) {
				visit.EntryTimestamp = event.Timestamp
			}

		case "page_exit":
			if visit.ExitTimestamp.IsZero() {
				visit.ExitTimestamp = event.Timestamp
			} else if event.Timestamp.After(visit.ExitTimestamp) {
				visit.ExitTimestamp = event.Timestamp
			}

		case "click":
			visit.ClickCount++

		case "identify":
			email, err := event.VisitorEmail()
			if err != nil {
				tracing.TraceErr(span, err)
				continue
			}
			emailValidate := mailvalidate.ValidateEmailSyntax(email)
			if emailValidate.Error != "" || emailValidate.IsSystemGenerated {
				continue
			}

			switch {
			case !emailValidate.IsFreeAccount:
				visit.VisitorDomain = emailValidate.Domain
				if !emailValidate.IsRoleAccount {
					visit.VisitorEmail = emailValidate.CleanEmail
					visit.VisitorEmailType = BusinessEmail
				}

			case emailValidate.IsFreeAccount:
				if visit.VisitorEmailType != BusinessEmail && !emailValidate.IsRoleAccount {
					visit.VisitorEmail = emailValidate.CleanEmail
					visit.VisitorEmailType = PersonalEmail
				}
			}

		}
	}

	if !visit.EntryTimestamp.IsZero() && !visit.ExitTimestamp.IsZero() {
		duration := visit.ExitTimestamp.Sub(visit.EntryTimestamp)
		visit.SecondsOnPage = int64(duration.Seconds())
	}

	content, err := c.webscraperService.Scrape(ctx, page)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	pageCategory, err := c.webscraperService.ClassifyWebpageCategory(ctx, page, &content)
	if err != nil {
		tracing.TraceErr(span, err)
		return visit, err
	}
	visit.PageCategory = pageCategory

	// return early is there is no webscrape content to analyze
	if content == "" {
		return visit, nil
	}

	contentStage, err := c.webscraperService.ClassifyContentStage(ctx, page, &content)
	if err != nil {
		tracing.TraceErr(span, err)
		return visit, err
	}
	visit.StageSignal = contentStage

	pageTopics, err := c.webscraperService.ClassifyWebpageTopics(ctx, page, &content)
	if err != nil {
		tracing.TraceErr(span, err)
		return visit, err
	}
	visit.Topics = pageTopics

	return visit, nil
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

func (c *AnalyzeWebSessionCapability) sessionAnalytics(ctx context.Context, sessionID string) (AnalyzeWebSessionOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalyzeWebSessionCapability.sessionAnalytics")
	defer span.Finish()
	tracing.TagComponentService(span)

	session, err := c.postgresRepositories.WebSessionRepository.FindSession(ctx, postgres_entity.WebSession{
		ID:     sessionID,
		Tenant: common.GetTenantFromContext(ctx),
	}, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return AnalyzeWebSessionOutput{}, err
	}
	if session == nil {
		return AnalyzeWebSessionOutput{}, nil
	}

	hostname := utils.StripUrlToBasePath(session.Hostname)

	var referrer string
	if session.Referrer == nil {
		referrer = ""
	} else {
		referrer = *session.Referrer
	}
	referrer = utils.StripUrlToBasePath(referrer)
	if strings.Contains(referrer, "syndicatedsearch.goog") {
		referrer = "google.com"
	}

	sessionDuration, err := c.calculateSessionDuration(ctx, session)
	if err != nil {
		tracing.TraceErr(span, err)
		return AnalyzeWebSessionOutput{}, err
	}

	results := AnalyzeWebSessionOutput{
		WebSessionID:    sessionID,
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

	results, err := c.postgresRepositories.WebSessionRepository.FindAllActiveSessions(ctx, query, nil)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalyzeWebSessionCapability.isNewWebsiteVisitor")
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
		Tenant:   tenant,
		IsActive: false,
	}

	results, err := c.postgresRepositories.WebSessionRepository.FindAllActiveSessions(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	if results == nil || len(results) == 0 {
		return true, nil
	}
	return false, nil
}

func (c *AnalyzeWebSessionCapability) buildTimelineMessage(ctx context.Context, sessionID string, analysis AnalyzeWebSessionOutput) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalyzeWebSessionCapability.buildTimelineMessage")
	defer span.Finish()
	tracing.TagComponentService(span)

	// Build base message
	var baseMessage string
	switch {
	case analysis.Referrer == "":
		baseMessage = fmt.Sprintf("A visitor browsed for %s", analysis.SessionDuration)
	default:
		baseMessage = fmt.Sprintf("A visitor referred by **%s** browsed for %s", analysis.Referrer, analysis.SessionDuration)
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
