package intent_signals

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"strings"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type WebVisitProcessor struct {
	postgresRepositories *postgres_repository.Repositories
	events               *events.EventsService
	organizationService  interfaces.OrganizationService
}

func NewWebVisitProcessor(postgresRepos *postgres_repository.Repositories, events *events.EventsService, org interfaces.OrganizationService) *WebVisitProcessor {
	return &WebVisitProcessor{
		postgresRepositories: postgresRepos,
		events:               events,
		organizationService:  org,
	}
}

const (
	IntentDetected   int8 = 1
	NoIntentDetected int8 = 2
)

func (p *WebVisitProcessor) Process(ctx context.Context, webSession postgres_entity.WebSession) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebVisitProcessor.Process")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	return p.processSupportSignal(ctx, webSession)
}

func (p *WebVisitProcessor) processSupportSignal(ctx context.Context, webSession postgres_entity.WebSession) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebVisitProcessor.processSupportSignal")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = p.validateWebSession(ctx, webSession)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if !p.isSignal(ctx, webSession.UniquePageViews, "support") {
		_, err := p.postgresRepositories.WebSessionRepository.UpdateIntentSignal(ctx, webSession.ID, webSession.Tenant, NoIntentDetected)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return nil
	}

	event := dto.IntentEvent{
		EventName:      enum.EventIntentSignal,
		Source:         enum.SourceWebtracker,
		SourceID:       webSession.ID,
		Tenant:         webSession.Tenant,
		IntentType:     enum.IntentSupportRequired,
		OrganizationID: *webSession.OrganizationId,
	}

	err = p.events.Publisher.PublishEvent(ctx, webSession.ID, model.INTENT_SIGNAL, &event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// update webSession record in database
	_, err = p.postgresRepositories.WebSessionRepository.UpdateIntentSignal(ctx, webSession.ID, webSession.Tenant, IntentDetected)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (p *WebVisitProcessor) isSignal(ctx context.Context, pageViews []string, pattern string) bool {
	for _, page := range pageViews {
		if strings.Contains(page, pattern) {
			return true
		}
	}
	return false
}

func (p *WebVisitProcessor) validateWebSession(ctx context.Context, webSession postgres_entity.WebSession) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebVisitProcessor.validateWebSession")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	switch {
	case webSession.ID == "":
		err := errors.New("session ID cannot be empty")
		return err
	case webSession.Tenant == "":
		err := errors.New("tenant cannot be empty")
		return err
	case len(webSession.UniquePageViews) == 0:
		err := errors.New("unique page views cannot be empty")
		return err
	case utils.IfNotNilString(webSession.OrganizationId) == "":
		err := errors.New("organization id cannot be empty")
		return err
	default:
		return nil
	}
}
