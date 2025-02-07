package agent_producers

import (
	"context"
	"fmt"
	"sort"
	"strings"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type NewWebSessionProducer struct {
	events               *events.EventsService
	postgresRepositories *postgres_repository.Repositories
}

func NewNewWebSessionProducer(
	events *events.EventsService,
	postgresRepositories *postgres_repository.Repositories,
) *NewWebSessionProducer {
	return &NewWebSessionProducer{
		events:               events,
		postgresRepositories: postgresRepositories,
	}
}

const (
	WebSessionTimeoutPageExit int = 5  // mins -- page exit without a following page view
	WebSessionTimeoutPageView int = 30 // mins -- page view without a page exit
)

// Add all Agent types subscribed to this event here
func (p *NewWebSessionProducer) subscribedAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentWebVisitorIdentifier,
	}
}

func (s *NewWebSessionProducer) Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	span, ctx := tracing.StartTracerSpan(ctx, "NewWebSessionProducer.ProcessWebSessions")
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

func (s *NewWebSessionProducer) findTimedOutPageExitEvents(ctx context.Context) ([]postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionProducer.findTimedOutPageExitEvents")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	lookback := WebSessionTimeoutPageExit
	query := postgres_entity.WebSession{
		LastEventType: enum.WebTrackerPageExit.String(),
		IsActive:      true,
	}
	return s.postgresRepositories.WebSessionRepository.FindAllActiveSessions(ctx, query, &lookback)
}

func (s *NewWebSessionProducer) findTimedOutPageViewEvents(ctx context.Context) ([]postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionProducer.findTimedOutPageViewEvents")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	lookback := WebSessionTimeoutPageView
	query := postgres_entity.WebSession{
		LastEventType: enum.WebTrackerPageView.String(),
		IsActive:      true,
	}
	return s.postgresRepositories.WebSessionRepository.FindAllActiveSessions(ctx, query, &lookback)
}

func (s *NewWebSessionProducer) closeSessions(ctx context.Context, sessions []postgres_entity.WebSession) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionProducer.closeSessions")
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
			AppSource: constants.AppSourceUpkeeper,
		})
		err := s.processClosedSession(innerCtx, session)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (s *NewWebSessionProducer) processClosedSession(ctx context.Context, session postgres_entity.WebSession) error {
	span, ctx := tracing.StartTracerSpan(ctx, "NewWebSessionProducer.processClosedSession")
	defer span.Finish()
	tracing.TagComponentCronJob(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	if session.ID == "" {
		err := errors.New("SessionID cannot be empty")
		tracing.TraceErr(span, err)
		return err
	}

	// close websession record
	endTime := session.LastActivity
	closedSession, err := s.postgresRepositories.WebSessionRepository.UpdateSessionEnd(ctx, session.ID, endTime)
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

	event, err := s.createCloseSessionEvent(ctx, *closedSession)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.events.Publisher.PublishFanoutEvent(ctx, event.WebSessionID, model.WEB_SESSION, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *NewWebSessionProducer) createCloseSessionEvent(ctx context.Context, session postgres_entity.WebSession) (dto.NewWebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionProducer.createCloseSessionWebhookEvent")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	if session.ID == "" || session.Tenant == "" || session.IP == "" || session.VisitorID == "" {
		err := errors.New("cannot build new web session event")
		tracing.TraceErr(span, err)
		return dto.NewWebSession{}, err
	}

	eventData := dto.NewWebSession{
		WebSessionID: session.ID,
		IPAddress:    session.IP,
		VisitorID:    session.VisitorID,
		Hostname:     session.Hostname,
	}

	return eventData, nil
}

func (c *NewWebSessionProducer) processUniquePageViews(ctx context.Context, tenant, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionProducer.getUniquePageViews")
	defer span.Finish()
	tracing.TagComponentService(span)

	session, err := c.postgresRepositories.WebTrackerEventsRepository.FindAll(ctx, postgres_entity.WebTrackerEvents{
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
	_, err = c.postgresRepositories.WebSessionRepository.UpdateSessionPageViews(ctx, sessionID, tenant, uniquePages)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to update websession with unique page views"))
		return err
	}

	return nil
}

func (c *NewWebSessionProducer) sortUrlsByLength(urls []string) []string {
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

func (c *NewWebSessionProducer) extractPageUrl(page postgres_entity.WebTrackerEvents) string {
	url := utils.CleanUrlBasePath(page.Hostname)
	if page.Pathname != "" {
		url = strings.TrimSuffix(fmt.Sprintf("%s/%s",
			utils.CleanUrlBasePath(page.Hostname),
			utils.CleanUrlBasePath(page.Pathname)), "/")
	}
	return url
}
