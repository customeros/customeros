package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	commonservice "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
)

type WebSesssionService interface {
	ProcessIntentSignals()
	ProcessWebSessions()
}

type webSessionService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonservice.CommonServices
}

func NewWebSessionService(cfg *config.Config, log logger.Logger, s *commonservice.CommonServices) WebSesssionService {
	return &webSessionService{
		cfg:            cfg,
		log:            log,
		commonServices: s,
	}
}

const (
	WebSessionTimeoutPageExit int = 5  // mins -- page exit without a following page view
	WebSessionTimeoutPageView int = 30 // mins -- page view without a page exit
)

func (s *webSessionService) ProcessWebSessions() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	span, ctx := tracing.StartTracerSpan(ctx, "WebSessionService.ProcessWebSessions")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	// find timed out page exit events
	pageExitSessions, err := s.findTimedOutPageExitEvents(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	// close sessions & fire webtracker event
	if pageExitSessions != nil {
		err = s.closeSessions(ctx, pageExitSessions)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
	}

	// find timed out page view events
	pageViewSessions, err := s.findTimedOutPageViewEvents(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	// close sessions & fire webtracker event
	if pageViewSessions != nil {
		err = s.closeSessions(ctx, pageViewSessions)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
	}
	return
}

func (s *webSessionService) ProcessIntentSignals() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	span, ctx := tracing.StartTracerSpan(ctx, "WebSessionService.ProcessWebSessions")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	// find closed events that have not been analyzed for intent
	sessions, err := s.commonServices.PostgresRepositories.WebSessionRepository.FindAllSessionsForIntentAnalysis(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	if len(sessions) == 0 {
		return
	}

	for _, session := range sessions {
		err := s.commonServices.WebVisitProcessor.Process(ctx, session)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}
}

func (s *webSessionService) findTimedOutPageExitEvents(ctx context.Context) ([]postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionService.findTimedOutPageExitEvents")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	lookback := WebSessionTimeoutPageExit
	query := postgres_entity.WebSession{
		LastEventType: enum.WebTrackerPageExit.String(),
		IsActive:      true,
	}
	return s.commonServices.PostgresRepositories.WebSessionRepository.FindAllActiveSessions(ctx, query, &lookback)
}

func (s *webSessionService) findTimedOutPageViewEvents(ctx context.Context) ([]postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionService.findTimedOutPageViewEvents")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	lookback := WebSessionTimeoutPageView
	query := postgres_entity.WebSession{
		LastEventType: enum.WebTrackerPageView.String(),
		IsActive:      true,
	}
	return s.commonServices.PostgresRepositories.WebSessionRepository.FindAllActiveSessions(ctx, query, &lookback)
}

func (s *webSessionService) closeSessions(ctx context.Context, sessions []postgres_entity.WebSession) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionService.closeSessions")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	if sessions == nil || len(sessions) == 0 {
		span.LogKV("result", "no_sessions_to_process")
		return nil
	}

	var errs error
	for _, session := range sessions {
		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    session.Tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})
		err := s.processClosedSession(innerCtx, session)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (s *webSessionService) processClosedSession(ctx context.Context, session postgres_entity.WebSession) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionService.processClosedSession")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	if session.ID == "" {
		err := errors.New("SessionID cannot be empty")
		tracing.TraceErr(span, err)
		return err
	}

	// close websession record
	endTime := session.LastActivity
	closedSession, err := s.commonServices.PostgresRepositories.WebSessionRepository.UpdateSessionEnd(ctx, session.ID, endTime)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if closedSession == nil {
		err := errors.New("closedSession is nil")
		tracing.TraceErr(span, err)
		return err
	}

	err = s.processUniquePageViews(ctx, session.Tenant, session.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	event, err := s.createCloseSessionWebhookEvent(ctx, *closedSession)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.commonServices.Events.Publisher.PublishWebhookEvent(ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *webSessionService) createCloseSessionWebhookEvent(ctx context.Context, session postgres_entity.WebSession) (dto.WebhookEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionService.createCloseSessionWebhookEvent")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	if session.ID == "" || session.Tenant == "" || session.IP == "" || session.VisitorID == "" {
		err := errors.New("cannot build web visit event")
		tracing.TraceErr(span, err)
		return dto.WebhookEvent{}, err
	}

	eventData := data_fields.WebsiteVisitEvent{
		SessionID: session.ID,
		Tenant:    session.Tenant,
		IPAddress: session.IP,
		VisitorID: session.VisitorID,
		Hostname:  session.Hostname,
	}

	event := dto.WebhookEvent{
		ExternalSystemId: enum.SourceAgent,
		Name:             enum.EventRevealWebsiteVisit,
		DataType:         eventData.Type(),
		Data:             &eventData,
	}

	return event, nil
}

func (c *webSessionService) processUniquePageViews(ctx context.Context, tenant, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "webSessionService.getUniquePageViews")
	defer span.Finish()
	tracing.TagComponentService(span)

	session, err := c.commonServices.PostgresRepositories.WebTrackerEventsRepository.FindAll(ctx, postgres_entity.WebTrackerEvents{
		SessionID: sessionID,
		Tenant:    tenant,
	}, nil)
	if err != nil {
		return err
	}
	if session == nil {
		return nil
	}

	// Use map to track unique pages
	uniquePageMap := make(map[string]struct{})
	for _, page := range session {
		url := c.extractPageUrl(page)
		uniquePageMap[url] = struct{}{}
	}

	// Convert map keys to slice
	uniquePages := make([]string, 0, len(uniquePageMap))
	for pathname := range uniquePageMap {
		uniquePages = append(uniquePages, pathname)
	}

	uniquePages = c.sortUrlsByLength(uniquePages)

	// strore in db
	_, err = c.commonServices.PostgresRepositories.WebSessionRepository.UpdateSessionPageViews(ctx, sessionID, tenant, uniquePages)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to update websession with unique page views"))
		return err
	}

	return nil
}

func (c *webSessionService) sortUrlsByLength(urls []string) []string {
	sort.Slice(urls, func(i, j int) bool {
		// If lengths are different, sort by length
		if len(urls[i]) != len(urls[j]) {
			return len(urls[i]) < len(urls[j])
		}
		// If lengths are equal, sort alphabetically
		return urls[i] < urls[j]
	})
	return urls
}

func (c *webSessionService) extractPageUrl(page postgres_entity.WebTrackerEvents) string {
	url := utils.CleanUrlBasePath(page.Hostname)
	if page.Pathname != "" {
		url = strings.TrimSuffix(fmt.Sprintf("%s/%s",
			utils.CleanUrlBasePath(page.Hostname),
			utils.CleanUrlBasePath(page.Pathname)), "/")
	}
	return url
}
