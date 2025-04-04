package postgres_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type tenantSettingsOpportunityStageRepository struct {
	gormDb *gorm.DB
}

type TenantSettingsOpportunityStageRepository interface {
	GetById(c context.Context, tenant, id string) (*postgres_entity.TenantSettingsOpportunityStage, error)
	GetOrInitialize(c context.Context, tenant string) ([]*postgres_entity.TenantSettingsOpportunityStage, error)
	Init(c context.Context, tenant string) error
	Store(c context.Context, postgres_entity postgres_entity.TenantSettingsOpportunityStage) (*postgres_entity.TenantSettingsOpportunityStage, error)
	Update(ctx context.Context, tenant, id string, label *string, likelihoodRate *int64, visible *bool) (*postgres_entity.TenantSettingsOpportunityStage, error)
}

func NewTenantSettingsOpportunityStageRepository(db *gorm.DB) TenantSettingsOpportunityStageRepository {
	return &tenantSettingsOpportunityStageRepository{gormDb: db}
}

func (r *tenantSettingsOpportunityStageRepository) GetById(ctx context.Context, tenant, id string) (*postgres_entity.TenantSettingsOpportunityStage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "TenantSettingsOpportunityStageRepository.GetById")
	defer spans.Finish()
	spans.LogKV("id", id)

	var result postgres_entity.TenantSettingsOpportunityStage
	err := r.gormDb.
		Where("id = ?", id).
		First(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	return &result, nil
}

func (r *tenantSettingsOpportunityStageRepository) GetOrInitialize(ctx context.Context, tenant string) ([]*postgres_entity.TenantSettingsOpportunityStage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "TenantSettingsOpportunityStageRepository.Get")
	defer spans.Finish()

	var entities []*postgres_entity.TenantSettingsOpportunityStage
	err := r.gormDb.
		Where("tenant = ?", tenant).
		Order("idx asc").
		Find(&entities).Error

	if err != nil {
		spans.TraceError(err)
		return nil, errors.Wrap(err, "error while getting tenant settings opportunity stages")
	}

	if len(entities) == 0 {
		err = r.Init(ctx, tenant)
		if err != nil {
			spans.TraceError(err)
			return nil, errors.Wrap(err, "error while initializing tenant settings opportunity stages")
		}

		err = r.gormDb.
			Where("tenant = ?", tenant).
			Order("idx asc").
			Find(&entities).Error

		if err != nil {
			spans.TraceError(err)
			return nil, errors.Wrap(err, "error while getting tenant settings opportunity stages")
		}
	}

	return entities, nil
}

func (r *tenantSettingsOpportunityStageRepository) Init(c context.Context, tenant string) error {
	spans, ctx := telemetry.StartPostgresSpan(c, "TenantSettingsOpportunityStageRepository.Init")
	defer spans.Finish()

	r.Store(ctx, postgres_entity.TenantSettingsOpportunityStage{
		Tenant:  tenant,
		Value:   "STAGE1",
		Label:   "Identified",
		Order:   1,
		Visible: true,
	})

	r.Store(ctx, postgres_entity.TenantSettingsOpportunityStage{
		Tenant:  tenant,
		Value:   "STAGE2",
		Label:   "Qualified",
		Order:   2,
		Visible: true,
	})

	r.Store(ctx, postgres_entity.TenantSettingsOpportunityStage{
		Tenant:  tenant,
		Value:   "STAGE3",
		Label:   "Committed",
		Order:   3,
		Visible: true,
	})

	for _, idx := range []int{4, 5, 6, 7, 8, 9, 10} {
		_, err := r.Store(ctx, postgres_entity.TenantSettingsOpportunityStage{
			Tenant:  tenant,
			Value:   fmt.Sprintf("STAGE%d", idx),
			Label:   fmt.Sprintf("Stage %d", idx),
			Order:   idx,
			Visible: false,
		})

		if err != nil {
			spans.TraceError(err)
			return err
		}
	}

	return nil
}

func (r *tenantSettingsOpportunityStageRepository) Store(ctx context.Context, postgres_entity postgres_entity.TenantSettingsOpportunityStage) (*postgres_entity.TenantSettingsOpportunityStage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "TenantSettingsOpportunityStageRepository.Store")
	defer spans.Finish()

	err := r.gormDb.Save(&postgres_entity).Error

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("postgres_entity.id", postgres_entity.ID)

	return &postgres_entity, nil
}

func (r *tenantSettingsOpportunityStageRepository) Update(ctx context.Context, tenant, id string, label *string, likelihoodRate *int64, visible *bool) (*postgres_entity.TenantSettingsOpportunityStage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "TenantSettingsOpportunityStageRepository.Update")
	defer spans.Finish()
	spans.LogKV("id", id)

	// update label, rate and visible if not null
	updateFields := map[string]interface{}{}
	if label != nil {
		updateFields["label"] = *label
	}
	if likelihoodRate != nil {
		updateFields["likelihood_rate"] = *likelihoodRate
	}
	if visible != nil {
		updateFields["visible"] = *visible
	}

	err := r.gormDb.
		Model(&postgres_entity.TenantSettingsOpportunityStage{}).
		Where("tenant = ? AND id = ?", tenant, id).
		Updates(updateFields).
		UpdateColumn("updated_at", utils.Now()).
		Error

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return r.GetById(ctx, tenant, id)
}
