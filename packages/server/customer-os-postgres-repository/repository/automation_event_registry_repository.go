package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type AutomationEventRegistryRepository interface {
	Create(ctx context.Context, event *entity.AutomationEventRegistry) (*entity.AutomationEventRegistry, error)
	FindAll(ctx context.Context) (*[]entity.AutomationEventRegistry, error)
	Find(ctx context.Context, event *entity.AutomationEventRegistry) (*entity.AutomationEventRegistry, error)
	Initialize(ctx context.Context) error
}

type flowEventRegistryRepository struct {
	gormDb *gorm.DB
}

func NewAutomationEventRegistryRepository(gormDb *gorm.DB) AutomationEventRegistryRepository {
	return &flowEventRegistryRepository{gormDb: gormDb}
}

func (r *flowEventRegistryRepository) Create(ctx context.Context, event *entity.AutomationEventRegistry) (*entity.AutomationEventRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AutomationEventRegistryRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.Create(event).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return event, nil
}

func (r *flowEventRegistryRepository) FindAll(ctx context.Context) (*[]entity.AutomationEventRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AutomationEventRegistryRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var events []entity.AutomationEventRegistry
	err := r.gormDb.
		Where("status = ?", enum.AutomationNodeEdgeStatusActive.String()).
		Order("external_system DESC").
		Find(&events).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &events, nil
}

func (r *flowEventRegistryRepository) Find(ctx context.Context, event *entity.AutomationEventRegistry) (*entity.AutomationEventRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AutomationEventRegistryRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var foundEvent entity.AutomationEventRegistry
	query := r.gormDb.
		Where("status = ?", enum.AutomationNodeEdgeStatusActive.String())

	if event != nil {
		query = query.Where(event)
	}

	err := query.First(&foundEvent).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &foundEvent, nil
}

func (r *flowEventRegistryRepository) Initialize(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AutomationEventRegistryRepository.Initialize")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	requiredEvents := []entity.AutomationEventRegistry{
		{
			ExternalSystem: "fathom",
			EventEvent:     commonenum.EventFathomMeetingSummaryCreated.String(),
			FriendlyName:   "Fathom Meeting Summary Created",
			Description:    "New AI meeting summary created by Fathom",
			Status:         commonenum.AutomationNodeEdgeStatusActive.String(),
		},
		{
			ExternalSystem: "flow",
			EventEvent:     commonenum.EventAutomationContactAdded.String(),
			FriendlyName:   "Contact added to Automation",
			Description:    "A new Contact has been added to a Automation",
			Status:         commonenum.AutomationNodeEdgeStatusActive.String(),
		},
		{
			ExternalSystem: "grain",
			EventEvent:     commonenum.EventGrainMeetingSummaryCreated.String(),
			FriendlyName:   "Grain Meeting Summary Created",
			Description:    "New AI meeting summary created by Grain",
			Status:         commonenum.AutomationNodeEdgeStatusActive.String(),
		},
		{
			ExternalSystem: "reveal",
			EventEvent:     commonenum.EventRevealWebsiteVisitNew.String(),
			FriendlyName:   "New Website Visitor",
			Description:    "Organization visits your website for the first time",
			Status:         commonenum.AutomationNodeEdgeStatusActive.String(),
		},
		{
			ExternalSystem: "reveal",
			EventEvent:     commonenum.EventRevealWebsiteVisitRepeat.String(),
			FriendlyName:   "Repeat Webpage Visitor",
			Description:    "Organization comes back to your website",
			Status:         commonenum.AutomationNodeEdgeStatusActive.String(),
		},
	}

	for _, event := range requiredEvents {
		existingEvent, err := r.Find(ctx, &entity.AutomationEventRegistry{
			EventEvent: event.EventEvent,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if existingEvent != nil {
			continue
		}

		if _, err := r.Create(ctx, &event); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
