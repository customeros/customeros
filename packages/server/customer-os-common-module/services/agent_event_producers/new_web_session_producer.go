package agent_producers

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go/log"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type NewWebSessionProducer struct {
	events               *events.EventsService
	postgresRepositories *postgres_repository.Repositories
	webscraperService    interfaces.WebscraperService
}

func NewNewWebSessionProducer(
	events *events.EventsService,
	postgresRepositories *postgres_repository.Repositories,
	webscraperService interfaces.WebscraperService,
) *NewWebSessionProducer {
	return &NewWebSessionProducer{
		events:               events,
		postgresRepositories: postgresRepositories,
		webscraperService:    webscraperService,
	}
}

const (
	WebSessionTimeoutPageExit int = 5  // mins -- page exit without a following page view
	WebSessionTimeoutPageView int = 30 // mins -- page view without a page exit
	WebSessionTimeoutClick    int = 30 // mins -- click without a page exit
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

	// find timed out click events
	clickSessions, err := s.findTimedOutClickEvents(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	// close sessions & fire webtracker event
	if clickSessions != nil {
		err = s.closeSessions(ctx, clickSessions)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
	}
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
	activeSessions, err := s.postgresRepositories.WebSessionRepository.FindAllActiveSessions(ctx, query, &lookback)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(activeSessions)))
	return activeSessions, nil
}

func (s *NewWebSessionProducer) findTimedOutClickEvents(ctx context.Context) ([]postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionProducer.findTimedOutClickEvents")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	lookback := WebSessionTimeoutClick
	query := postgres_entity.WebSession{
		LastEventType: enum.WebTrackerClick.String(),
		IsActive:      true,
	}

	activeSessions, err := s.postgresRepositories.WebSessionRepository.FindAllActiveSessions(ctx, query, &lookback)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(activeSessions)))
	return activeSessions, nil
}

func (s *NewWebSessionProducer) closeSessions(ctx context.Context, sessions []postgres_entity.WebSession) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionProducer.closeSessions")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	if len(sessions) == 0 {
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
	closedSession, err := s.postgresRepositories.WebSessionRepository.SetSessionEnd(ctx, session.ID, endTime)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if closedSession == nil {
		err := errors.New("closedSession is nil")
		tracing.TraceErr(span, err)
		return err
	}

	pageViews, err := s.processUniquePageViews(ctx, session.Tenant, session.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// process page visits
	identifiedVisitor := IdentifiedVisitor{}
	for _, page := range pageViews {
		visitor, err := s.processPageView(ctx, session.ID, page, endTime)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}
		if identifiedVisitor.Domain == "" && visitor.Domain != "" {
			identifiedVisitor.Domain = visitor.Domain
		}
		if identifiedVisitor.Email == "" && visitor.Email != "" {
			identifiedVisitor.Email = visitor.Email
			identifiedVisitor.EmailType = visitor.EmailType
			identifiedVisitor.Domain = visitor.Domain
		}
		if identifiedVisitor.EmailType == enum.EmailPersonal && visitor.EmailType == enum.EmailBusiness {
			identifiedVisitor.Email = visitor.Email
			identifiedVisitor.EmailType = visitor.EmailType
			identifiedVisitor.Domain = visitor.Domain
		}
	}

	// update websession with identity
	if identifiedVisitor.Domain != "" || identifiedVisitor.Email != "" {
		err := s.postgresRepositories.WebSessionRepository.SetVisitorIdentity(ctx, session.ID, &identifiedVisitor.Domain, &identifiedVisitor.Email, &identifiedVisitor.EmailType)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	err = s.events.Publisher.PublishFanoutEvent(ctx, session.ID, model.WEB_SESSION, dto.NewWebSession{})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

type IdentifiedVisitor struct {
	Domain    string
	Email     string
	EmailType enum.EmailType
}

func (s *NewWebSessionProducer) processPageView(ctx context.Context, sessionId, page string, endTime time.Time) (IdentifiedVisitor, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionProducer.processPageView")
	defer span.Finish()
	tracing.TagComponentCronJob(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	span.LogFields(log.String("sessionId", sessionId), log.String("page", page), log.String("endTime", endTime.String()))

	visitor := IdentifiedVisitor{}

	// get session events
	events, err := s.postgresRepositories.WebTrackerEventsRepository.FindEventsForPageVisit(ctx, sessionId, page)
	if err != nil {
		tracing.TraceErr(span, err)
		return visitor, err
	}

	url := page
	if !strings.HasPrefix(page, "http") {
		url = fmt.Sprintf("https://%s", page)
	}

	visit := postgres_entity.WebSessionPageVisit{
		SessionID: sessionId,
		Url:       url,
	}

	for _, event := range events {
		if visit.Referer == "" && event.Referrer != "" {
			visit.Referer = event.Referer()
		}
		if visit.Hostname == "" && event.Hostname != "" {
			visit.Hostname = event.Hostname
		}
		if visit.Pathname == "" && event.Pathname != "" {
			visit.Pathname = event.Pathname
		}

		switch event.EventType {
		case "page_view":
			if visit.EntryTimestamp.IsZero() {
				visit.EntryTimestamp = event.Timestamp
			} else if event.Timestamp.Before(visit.EntryTimestamp) {
				visit.EntryTimestamp = event.Timestamp
			}

		case "page_exit":
			if event.Timestamp.IsZero() {
				event.Timestamp = endTime
			}

			if visit.ExitTimestamp.IsZero() {
				visit.ExitTimestamp = event.Timestamp
			} else if event.Timestamp.After(visit.ExitTimestamp) {
				visit.ExitTimestamp = event.Timestamp
			}

		case "identify":
			email, err := event.VisitorEmail()
			if err != nil {
				tracing.TraceErr(span, err)
				continue
			}
			if email == "" {
				continue
			}

			emailValidate := mailvalidate.ValidateEmailSyntax(email)
			if emailValidate.Error != "" || emailValidate.IsSystemGenerated {
				continue
			}

			switch {
			case !emailValidate.IsFreeAccount:
				visitor.Domain = emailValidate.Domain
				if !emailValidate.IsRoleAccount {
					visitor.Email = emailValidate.CleanEmail
					visitor.EmailType = enum.EmailBusiness
				}

			case emailValidate.IsFreeAccount:
				if visitor.EmailType != enum.EmailBusiness && !emailValidate.IsRoleAccount {
					visitor.Email = emailValidate.CleanEmail
					visitor.EmailType = enum.EmailPersonal
				}
			}

		}
	}

	// save pagevisit to db
	_, err = s.postgresRepositories.WebSessionPageVisitRepository.Create(ctx, visit)
	if err != nil {
		tracing.TraceErr(span, err)
		return visitor, err
	}

	// check to see if webpage has been scraped
	scrapedPage, err := s.postgresRepositories.ScrapedWebpageRepository.GetWebpage(ctx, url, 365)
	if err != nil {
		tracing.TraceErr(span, err)
		return visitor, err
	}

	// handle all permutations of scraped data
	switch {
	case scrapedPage == nil:
		err := s.scrapeAndClassifyWebpage(ctx, url)
		if err != nil {
			tracing.TraceErr(span, err)
			return visitor, err
		}
		return visitor, nil

	case scrapedPage.Error != "":
		if scrapedPage.Category == "" {
			_, err = s.webscraperService.ClassifyWebpageCategory(ctx, url, &scrapedPage.Content)
			if err != nil {
				tracing.TraceErr(span, err)
				return visitor, err
			}
		}
		return visitor, nil

	case scrapedPage.Content == "":
		err := s.scrapeAndClassifyWebpage(ctx, url)
		if err != nil {
			tracing.TraceErr(span, err)
			return visitor, err
		}
		return visitor, nil

	default:
		if scrapedPage.Category == "" {
			_, err = s.webscraperService.ClassifyWebpageCategory(ctx, url, &scrapedPage.Content)
			if err != nil {
				tracing.TraceErr(span, err)
				return visitor, err
			}
		}

		if scrapedPage.ContentStage == "" {
			_, err = s.webscraperService.ClassifyContentStage(ctx, url, &scrapedPage.Content)
			if err != nil {
				tracing.TraceErr(span, err)
				return visitor, err
			}
		}

		if len(scrapedPage.Topics) == 0 {
			_, err = s.webscraperService.ClassifyWebpageTopics(ctx, page, &scrapedPage.Content)
			if err != nil {
				tracing.TraceErr(span, err)
				return visitor, err
			}
		}
		return visitor, nil
	}
}

func (s *NewWebSessionProducer) scrapeAndClassifyWebpage(ctx context.Context, url string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionProducer.scrapeAndClassifyWebpage")
	defer span.Finish()
	tracing.TagComponentService(span)

	content, err := s.webscraperService.Scrape(ctx, url)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = s.webscraperService.ClassifyWebpageCategory(ctx, url, &content)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if content == "" {
		return nil
	}

	_, err = s.webscraperService.ClassifyContentStage(ctx, url, &content)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = s.webscraperService.ClassifyWebpageTopics(ctx, url, &content)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (s *NewWebSessionProducer) processUniquePageViews(ctx context.Context, tenant, sessionID string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionProducer.getUniquePageViews")
	defer span.Finish()
	tracing.TagComponentService(span)

	session, err := s.postgresRepositories.WebTrackerEventsRepository.FindAll(ctx, postgres_entity.WebTrackerEvents{
		SessionID: sessionID,
		Tenant:    tenant,
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
		url := s.extractPageUrl(page)
		uniquePageMap[url] = struct{}{}
	}

	// Convert map keys to slice
	uniquePages := make([]string, 0, len(uniquePageMap))
	for pathname := range uniquePageMap {
		uniquePages = append(uniquePages, pathname)
	}

	uniquePages = s.sortUrlsByLength(uniquePages)

	// strore in db
	_, err = s.postgresRepositories.WebSessionRepository.SetSessionPageViews(ctx, sessionID, tenant, uniquePages)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to update websession with unique page views"))
	}

	return uniquePages, nil
}

func (s *NewWebSessionProducer) sortUrlsByLength(urls []string) []string {
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

func (s *NewWebSessionProducer) extractPageUrl(page postgres_entity.WebTrackerEvents) string {
	url := utils.StripUrlToBasePath(page.Hostname)
	if page.Pathname != "" {
		url = strings.TrimSuffix(fmt.Sprintf("%s/%s",
			utils.StripUrlToBasePath(page.Hostname),
			utils.StripUrlToBasePath(page.Pathname)), "/")
	}
	return url
}
