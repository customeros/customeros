package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowEventsRepository interface {
	CreateFlowEvent(ctx context.Context, event *entity.FlowEvent) error
	GetAllFlowEvents(ctx context.Context) ([]entity.FlowEvent, error)
	FindFlowEventByName(ctx context.Context, eventName string) (entity.FlowEvent, error)
	GetFlowEventsByExternalSystem(ctx context.Context, externalSystem enum.ExternalSystemId) ([]entity.FlowEvent, error)
}

type flowEventsRepository struct {
	gormDb *gorm.DB
}

func NewFlowEventsRepository(gormDb *gorm.DB) (FlowEventsRepository, error) {
	r := &flowEventsRepository{gormDb: gormDb}

	if err := r.InitializeEvents(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize flow events: %w", err)
	}
	return r, nil
}

func (r *flowEventsRepository) InitializeEvents(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowEventRepository.InitializeEvents")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	requiredEvents := []entity.FlowEvent{
		{
			ExternalSystem: enum.Fathom.String(),
			Resource:       "meeting_summary",
			Action:         "created",
			Description:    "New AI meeting summary created by Fathom",
			Enabled:        true,
		},
		// Add other required events...
	}

	for _, event := range requiredEvents {
		event.EventName = fmt.Sprintf("%s.%s.%s", event.ExternalSystem, event.Resource, event.Action)

		existingEvent, err := r.FindFlowEventByName(ctx, event.EventName)
		if err != nil {
			return err
		}

		if existingEvent.EventName != "" {
			continue
		}

		if err := r.CreateFlowEvent(ctx, &event); err != nil {
			return err
		}
	}

	return nil
}

func (r *flowEventsRepository) FindFlowEventByName(ctx context.Context, eventName string) (entity.FlowEvent, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowEventRepository.FindFlowEventByName")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var event entity.FlowEvent
	err := r.gormDb.
		Where("event_name = ? AND enabled = true", eventName).
		First(&event).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return event, err
}

func (r *flowEventsRepository) CreateFlowEvent(ctx context.Context, event *entity.FlowEvent) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowEventRepository.CreateFlowEvent")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if err := r.gormDb.Create(event).Error; err != nil {
		return err
	}
	return nil
}

func (r *flowEventsRepository) GetAllFlowEvents(ctx context.Context) ([]entity.FlowEvent, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowEventRepository.GetFlowEvents")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var events []entity.FlowEvent
	err := r.gormDb.
		Where("enabled = true").
		Order("external_system DESC").
		Find(&events).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return events, err
}

func (r *flowEventsRepository) GetFlowEventsByExternalSystem(ctx context.Context, externalSystem enum.ExternalSystemId) ([]entity.FlowEvent, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowEventRepository.GetFlowEvents")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var events []entity.FlowEvent
	err := r.gormDb.
		Where("enabled = true AND external_system = ?", externalSystem.String()).
		Order("external_system DESC").
		Find(&events).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return events, err
}
