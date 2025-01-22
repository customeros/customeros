package service

import (
	"context"
	"errors"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	commonservice "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
)

type WebSesssionService interface {
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

func (s *webSessionService) findTimedOutPageExitEvents(ctx context.Context) ([]postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionService.findTimedOutPageExitEvents")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	lookback := WebSessionTimeoutPageExit
	query := postgres_entity.WebSession{
		LastEventType: enum.WebTrackerPageExit.String(),
		IsActive:      true,
	}
	return s.commonServices.PostgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, &lookback)
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
	return s.commonServices.PostgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, &lookback)
}

func (s *webSessionService) closeSessions(ctx context.Context, sessions []postgres_entity.WebSession) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionService.closeSessions")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	if sessions == nil || len(sessions) == 0 {
		err := errors.New("There are no sessions to process")
		tracing.TraceErr(span, err)
		return err
	}

	var errs error
	for _, session := range sessions {
		err := s.processClosedSession(ctx, session)
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
	endTime := s.calculateSessionEnd(ctx, session)
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

func (s *webSessionService) calculateSessionEnd(ctx context.Context, session postgres_entity.WebSession) time.Time {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionService.calculateSessionEnd")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	return session.LastActivity.Add(time.Minute)
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
	}

	event := dto.WebhookEvent{
		ExternalSystemId: enum.SourceAgent,
		Name:             enum.EventRevealWebsiteVisit,
		DataType:         eventData.Type(),
		Data:             &eventData,
	}

	return event, nil
}
