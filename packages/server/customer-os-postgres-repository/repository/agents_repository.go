package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type AgentsRepository interface {
	Create(ctx context.Context, automationRecord entity.Agents) (*entity.Agents, error)
	Find(ctx context.Context, automationRecord entity.Agents) (*entity.Agents, error)
	FindAll(ctx context.Context, automationRecord entity.Agents) (*[]entity.Agents, error)
	Update(ctx context.Context, automationRecord entity.Agents) (*entity.Agents, error)
}

type agentsRepository struct {
	gormDb *gorm.DB
}

func NewAgentsRepository(gormDb *gorm.DB) AgentsRepository {
	return &agentsRepository{gormDb: gormDb}
}

func (f *agentsRepository) Create(ctx context.Context, automationRecord entity.Agents) (*entity.Agents, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentsRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	automationRecord.ID = utils.GenerateNanoIdWithPrefix("auto", 16)

	var created entity.Agents
	err := f.gormDb.Create(&automationRecord).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &automationRecord, nil
}

func (f *agentsRepository) FindAll(ctx context.Context, automationRecord entity.Agents) (*[]entity.Agents, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentsRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var agents []entity.Agents
	query := f.gormDb.Where("is_active = true")

	// Add additional filters based on non-zero fields in automationRecord
	if automationRecord != (entity.Agents{}) {
		query = query.Where(&automationRecord)
	}

	err := query.Find(&agents).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &agents, nil
}

func (f *agentsRepository) Find(ctx context.Context, automationRecord entity.Agents) (*entity.Agents, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentsRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var automation entity.Agents
	err := f.gormDb.
		Where(&automationRecord).
		Where("is_active = true").
		First(&automation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &automation, nil
}

func (f *agentsRepository) Update(ctx context.Context, automationRecord entity.Agents) (*entity.Agents, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentsRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if automationRecord.ID == "" {
		err := errors.New("automationID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedAgents entity.Agents
	err := f.gormDb.
		Model(&automationRecord).
		Updates(&automationRecord).
		First(&updatedAgents).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedAgents, nil
}
