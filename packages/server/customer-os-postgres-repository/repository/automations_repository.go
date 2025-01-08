package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type AutomationsRepository interface {
	Create(ctx context.Context, automationRecord entity.Automations) (*entity.Automations, error)
	Find(ctx context.Context, automationRecord entity.Automations) (*entity.Automations, error)
	FindAll(ctx context.Context, automationRecord entity.Automations) (*[]entity.Automations, error)
	Update(ctx context.Context, automationRecord entity.Automations) (*entity.Automations, error)
}

type automationsRepository struct {
	gormDb *gorm.DB
}

func NewAutomationsRepository(gormDb *gorm.DB) AutomationsRepository {
	return &automationsRepository{gormDb: gormDb}
}

func (f *automationsRepository) Create(ctx context.Context, automationRecord entity.Automations) (*entity.Automations, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AutomationsRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	automationRecord.ID = utils.GenerateNanoIdWithPrefix("auto", 16)

	var created entity.Automations
	err := f.gormDb.Create(&automationRecord).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &automationRecord, nil
}

func (f *automationsRepository) FindAll(ctx context.Context, automationRecord entity.Automations) (*[]entity.Automations, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AutomationsRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var automations []entity.Automations
	query := f.gormDb.Where("is_active = true")

	// Add additional filters based on non-zero fields in automationRecord
	if automationRecord != (entity.Automations{}) {
		query = query.Where(&automationRecord)
	}

	err := query.Find(&automations).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &automations, nil
}

func (f *automationsRepository) Find(ctx context.Context, automationRecord entity.Automations) (*entity.Automations, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AutomationsRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var automation entity.Automations
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

func (f *automationsRepository) Update(ctx context.Context, automationRecord entity.Automations) (*entity.Automations, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AutomationsRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if automationRecord.ID == "" {
		err := errors.New("automationID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedAutomations entity.Automations
	err := f.gormDb.
		Model(&automationRecord).
		Updates(&automationRecord).
		First(&updatedAutomations).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedAutomations, nil
}
