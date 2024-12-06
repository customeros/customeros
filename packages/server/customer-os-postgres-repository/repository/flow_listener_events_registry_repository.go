package repository

import (
	"context"
	"errors"

	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowListenerRegistryRepository interface {
	CreateFlowListenerEvent(ctx context.Context, event *entity.FlowListenerRegistry) error
	GetAllFlowListenerEvents(ctx context.Context) ([]entity.FlowListenerRegistry, error)
	FindFlowListenerEvent(ctx context.Context, eventName string) (entity.FlowListenerRegistry, error)
	GetFlowListenerEventsByExternalSystem(ctx context.Context, externalSystem enum.ExternalSystemId) ([]entity.FlowListenerRegistry, error)
	InitializeFlowListenerEvents(ctx context.Context) error
}

type flowListenerRegistryRepository struct {
	gormDb *gorm.DB
}

func NewFlowListenerRegistryRepository(gormDb *gorm.DB) FlowListenerRegistryRepository {
	return &flowListenerRegistryRepository{gormDb: gormDb}
}

func (r *flowListenerRegistryRepository) InitializeFlowListenerEvents(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowListenerRegistryRepository.InitializeListeners")
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
			ExternalSystem: "grain",
			ListenerEvent:  commonenum.EventGrainMeetingSummaryCreated.String(),
			FriendlyName:   "Grain Meeting Summary Created",
			Description:    "New AI meeting summary created by Grain",
			Enabled:        true,
		},
		// Add other required events...
	}

	for _, event := range requiredEvents {
		existingEvent, err := r.FindFlowListenerEvent(ctx, event.ListenerEvent)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if existingEvent.ListenerEvent != "" {
			continue
		}

		if err := r.CreateFlowListenerEvent(ctx, &event); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (r *flowListenerRegistryRepository) FindFlowListenerEvent(ctx context.Context, eventName string) (entity.FlowListenerRegistry, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowListenerRegistryRepository.FindFlowEventByName")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var event entity.FlowListenerRegistry
	err := r.gormDb.
		Where("listener_event = ? AND enabled = true", eventName).
		First(&event).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return event, err
}

func (r *flowListenerRegistryRepository) CreateFlowListenerEvent(ctx context.Context, event *entity.FlowListenerRegistry) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowListenerRegistryRepository.CreateFlowListenerEvent")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if err := r.gormDb.Create(event).Error; err != nil {
		return err
	}
	return nil
}

func (r *flowListenerRegistryRepository) GetAllFlowListenerEvents(ctx context.Context) ([]entity.FlowListenerRegistry, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowListenerRegistryRepository.GetAllFlowListenerEvents")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var events []entity.FlowListenerRegistry
	err := r.gormDb.
		Where("enabled = true").
		Order("external_system DESC").
		Find(&events).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return events, err
}

func (r *flowListenerRegistryRepository) GetFlowListenerEventsByExternalSystem(ctx context.Context, externalSystem enum.ExternalSystemId) ([]entity.FlowListenerRegistry, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowListenerRegistryRepository.GetFlowListenerEventsByExternalSystem")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var events []entity.FlowListenerRegistry
	err := r.gormDb.
		Where("enabled = true AND external_system = ?", externalSystem.String()).
		Order("external_system DESC").
		Find(&events).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return events, err
}
