package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowTransitionsRegistryRepository interface {
	Initialize(ctx context.Context) error
	GetAll(ctx context.Context, transition *entity.FlowTransitionsRegistry) ([]entity.FlowTransitionsRegistry, error)
	GetAllActive(ctx context.Context, transition *entity.FlowTransitionsRegistry) ([]entity.FlowTransitionsRegistry, error)
	Find(ctx context.Context, transition entity.FlowTransitionsRegistry) (*entity.FlowTransitionsRegistry, error)
	Create(ctx context.Context, transition *entity.FlowTransitionsRegistry) (*entity.FlowTransitionsRegistry, error)
}

type flowTransitionsRegistryRepository struct {
	gormDb *gorm.DB
}

type FromNodeRecord struct {
	FromNode     string `gorm:"column:from_node"`
	FromNodeType string `gorm:"column:from_node_type"`
}

func NewFlowTransitionsRegistryRepository(gormDb *gorm.DB) FlowTransitionsRegistryRepository {
	return &flowTransitionsRegistryRepository{gormDb: gormDb}
}

func (r *flowTransitionsRegistryRepository) Create(ctx context.Context, transition *entity.FlowTransitionsRegistry) (*entity.FlowTransitionsRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowTransitionsRegistryRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.WithContext(ctx).Create(transition).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return transition, nil
}

func (r *flowTransitionsRegistryRepository) GetAll(ctx context.Context, transition *entity.FlowTransitionsRegistry) ([]entity.FlowTransitionsRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowTransitionsRegistryRepository.GetAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var transitions []entity.FlowTransitionsRegistry
	query := r.gormDb.WithContext(ctx)

	// Add additional filters if transition is not nil
	if transition != nil {
		query = query.Where(transition)
	}

	err := query.
		Order("from_node DESC").
		Find(&transitions).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return transitions, nil
}

func (r *flowTransitionsRegistryRepository) GetAllActive(ctx context.Context, transition *entity.FlowTransitionsRegistry) ([]entity.FlowTransitionsRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowTransitionsRegistryRepository.GetAllActive")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var transitions []entity.FlowTransitionsRegistry
	query := r.gormDb.WithContext(ctx).
		Where("status = ?", enum.FlowNodeEdgeStatusActive.String())

	// Add additional filters if transition is not nil
	if transition != nil {
		query = query.Where(transition)
	}

	err := query.
		Order("from_node DESC").
		Find(&transitions).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return transitions, nil
}

func (r *flowTransitionsRegistryRepository) Find(ctx context.Context, transition entity.FlowTransitionsRegistry) (*entity.FlowTransitionsRegistry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowTransitionsRegistryRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var foundTransition entity.FlowTransitionsRegistry
	query := r.gormDb.WithContext(ctx).
		Where("status = ?", enum.FlowNodeEdgeStatusActive.String())

	// Add additional filters based on non-zero fields in transition
	if transition != (entity.FlowTransitionsRegistry{}) {
		query = query.Where(&transition)
	}

	err := query.First(&foundTransition).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &foundTransition, nil
}

func (r *flowTransitionsRegistryRepository) Initialize(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowTransitionsRegistryRepository.Initialize")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	return nil
}
