package repository

import (
	"context"
	"errors"

	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowListenerRegistryRepository interface {
	Create(ctx context.Context, event *entity.FlowListenerRegistry) (*entity.FlowListenerRegistry, error)
	FindAll(ctx context.Context) (*[]entity.FlowListenerRegistry, error)
	Find(ctx context.Context, event *entity.FlowListenerRegistry) (*entity.FlowListenerRegistry, error)
	Initialize(ctx context.Context) error
}

type flowListenerRegistryRepository struct {
	gormDb *gorm.DB
}

func NewFlowListenerRegistryRepository(gormDb *gorm.DB) FlowListenerRegistryRepository {
	return &flowListenerRegistryRepository{gormDb: gormDb}
}

func (r *flowListenerRegistryRepository) Create(ctx context.Context, event *entity.FlowListenerRegistry) (*entity.FlowListenerRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowListenerRegistryRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.Create(event).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return event, nil
}

func (r *flowListenerRegistryRepository) FindAll(ctx context.Context) (*[]entity.FlowListenerRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowListenerRegistryRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var events []entity.FlowListenerRegistry
	err := r.gormDb.
		Where("enabled = ?", true).
		Order("external_system DESC").
		Find(&events).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &events, nil
}

func (r *flowListenerRegistryRepository) Find(ctx context.Context, event *entity.FlowListenerRegistry) (*entity.FlowListenerRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowListenerRegistryRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var foundEvent entity.FlowListenerRegistry
	query := r.gormDb.Where("enabled = ?", true)

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

func (r *flowListenerRegistryRepository) Initialize(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowListenerRegistryRepository.Initialize")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	requiredEvents := []entity.FlowListenerRegistry{
		{
			ExternalSystem: "fathom",
			ListenerEvent:  commonenum.EventFathomMeetingSummaryCreated.String(),
			FriendlyName:   "Fathom Meeting Summary Created",
			Description:    "New AI meeting summary created by Fathom",
			Enabled:        true,
		},
		{
			ExternalSystem: "flow",
			ListenerEvent:  commonenum.EventFlowContactAdded.String(),
			FriendlyName:   "Contact added to Flow",
			Description:    "A new Contact has been added to a Flow",
			Enabled:        true,
		},
		{
			ExternalSystem: "grain",
			ListenerEvent:  commonenum.EventGrainMeetingSummaryCreated.String(),
			FriendlyName:   "Grain Meeting Summary Created",
			Description:    "New AI meeting summary created by Grain",
			Enabled:        true,
		},
	}

	for _, event := range requiredEvents {
		existingEvent, err := r.Find(ctx, &entity.FlowListenerRegistry{
			ListenerEvent: event.ListenerEvent,
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
