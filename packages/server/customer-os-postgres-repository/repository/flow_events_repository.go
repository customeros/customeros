package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowEventsRepository interface {
	CreateFlowEvent(ctx context.Context, event *entity.FlowEvent) error
	GetAllFlowEvents(ctx context.Context) ([]entity.FlowEvent, error)
	GetFlowEventsByExternalSystem(ctx context.Context, externalSystem enum.ExternalSystemId) ([]entity.FlowEvent, error)
}

type flowEventsRepository struct {
	gormDb *gorm.DB
}

func NewFlowEventsRepository(gormDb *gorm.DB) FlowEventsRepository {
	return &flowEventsRepository{gormDb: gormDb}
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
