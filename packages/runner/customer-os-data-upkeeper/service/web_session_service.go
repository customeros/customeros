package service

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	commonservice "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
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

func (s *webSessionService) ProcessWebSessions() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	span, ctx := tracing.StartTracerSpan(ctx, "WebSessionService.ProcessWebSessions")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	// find timed out page exit events
	lookback := int(enum.WebSessionTimeoutPageExit)
	query := postgres_entity.WebSession{
		LastEventType: enum.WebTrackerPageExit.String(),
		IsActive:      true,
	}
	pageExitSessions, err := s.commonServices.PostgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, &lookback)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	// close sessions & fire webtracker event
	if pageExitSessions != nil {
		err = s.closeSessions(ctx, pageExitSessions)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	// find timed out page view events
	query = postgres_entity.WebSession{
		LastEventType: enum.WebTrackerPageView.String(),
		IsActive:      true,
	}
	pageViewSessions, err := s.commonServices.PostgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, &lookback)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	// close sessions & fire webtracker event
	if pageViewSessions != nil {
		err = s.closeSessions(ctx, pageViewSessions)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}
	return
}

func (s *webSessionService) closeSessions(ctx context.Context, sessions []postgres_entity.WebSession) error {
	span, ctx := tracing.StartTracerSpan(ctx, "WebSessionService.closeSessions")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	var errs error
	for _, session := range sessions {
		// create web visit event
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

		// publish event
		err := s.commonServices.Events.Publisher.PublishWebhookEvent(ctx, event)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}

		// update websession record
		session.IsActive = false
		session.EndTime = utils.NowPtr()
		session.PublishedEvent = true
		_, err = s.commonServices.PostgresRepositories.WebSessionRepository.Update(ctx, session)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}
