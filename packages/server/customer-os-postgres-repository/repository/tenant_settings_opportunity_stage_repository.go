package repository

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type tenantSettingsOpportunityStageRepository struct {
	gormDb *gorm.DB
}

type TenantSettingsOpportunityStageRepository interface {
	GetById(c context.Context, tenant, id string) (*entity.TenantSettingsOpportunityStage, error)
	GetOrInitialize(c context.Context, tenant string) ([]*entity.TenantSettingsOpportunityStage, error)
	Init(c context.Context, tenant string) error
	Store(c context.Context, entity entity.TenantSettingsOpportunityStage) (*entity.TenantSettingsOpportunityStage, error)
	Update(ctx context.Context, tenant, id string, label *string, likelihoodRate *int64, visible *bool) (*entity.TenantSettingsOpportunityStage, error)
}

func NewTenantSettingsOpportunityStageRepository(db *gorm.DB) TenantSettingsOpportunityStageRepository {
	return &tenantSettingsOpportunityStageRepository{gormDb: db}
}

func (r *tenantSettingsOpportunityStageRepository) GetById(c context.Context, tenant, id string) (*entity.TenantSettingsOpportunityStage, error) {
	span, _ := opentracing.StartSpanFromContext(c, "TenantSettingsOpportunityStageRepository.GetById")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(tracingLog.String("id", id))

	var result entity.TenantSettingsOpportunityStage
	err := r.gormDb.
		Where("id = ?", id).
		First(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &result, nil
}

func (r *tenantSettingsOpportunityStageRepository) GetOrInitialize(c context.Context, tenant string) ([]*entity.TenantSettingsOpportunityStage, error) {
	span, _ := opentracing.StartSpanFromContext(c, "TenantSettingsOpportunityStageRepository.Get")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)

	var entities []*entity.TenantSettingsOpportunityStage
	err := r.gormDb.
		Where("tenant = ?", tenant).
		Order("idx asc").
		Find(&entities).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "error while getting tenant settings opportunity stages")
	}

	if len(entities) == 0 {
		err = r.Init(c, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, errors.Wrap(err, "error while initializing tenant settings opportunity stages")
		}

		err = r.gormDb.
			Where("tenant = ?", tenant).
			Order("idx asc").
			Find(&entities).Error

		if err != nil {
			tracing.TraceErr(span, err)
			return nil, errors.Wrap(err, "error while getting tenant settings opportunity stages")
		}
	}

	return entities, nil
}

func (r *tenantSettingsOpportunityStageRepository) Init(c context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(c, "TenantSettingsOpportunityStageRepository.Init")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)

	r.Store(ctx, entity.TenantSettingsOpportunityStage{
		Tenant:  tenant,
		Value:   "STAGE1",
		Label:   "Identified",
		Order:   1,
		Visible: true,
	})

	r.Store(ctx, entity.TenantSettingsOpportunityStage{
		Tenant:  tenant,
		Value:   "STAGE2",
		Label:   "Qualified",
		Order:   2,
		Visible: true,
	})

	r.Store(ctx, entity.TenantSettingsOpportunityStage{
		Tenant:  tenant,
		Value:   "STAGE3",
		Label:   "Committed",
		Order:   3,
		Visible: true,
	})

	for _, idx := range []int{4, 5, 6, 7, 8, 9, 10} {
		_, err := r.Store(ctx, entity.TenantSettingsOpportunityStage{
			Tenant:  tenant,
			Value:   fmt.Sprintf("STAGE%d", idx),
			Label:   fmt.Sprintf("Stage %d", idx),
			Order:   idx,
			Visible: false,
		})

		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (r *tenantSettingsOpportunityStageRepository) Store(c context.Context, entity entity.TenantSettingsOpportunityStage) (*entity.TenantSettingsOpportunityStage, error) {
	span, _ := opentracing.StartSpanFromContext(c, "TenantSettingsOpportunityStageRepository.Store")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, entity.Tenant)

	err := r.gormDb.Save(&entity).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.String("entity.id", entity.ID))

	return &entity, nil
}

func (r *tenantSettingsOpportunityStageRepository) Update(ctx context.Context, tenant, id string, label *string, likelihoodRate *int64, visible *bool) (*entity.TenantSettingsOpportunityStage, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsOpportunityStageRepository.Update")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(tracingLog.String("id", id))

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
		Model(&entity.TenantSettingsOpportunityStage{}).
		Where("tenant = ? AND id = ?", tenant, id).
		Updates(updateFields).
		UpdateColumn("updated_at", utils.Now()).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return r.GetById(ctx, tenant, id)
}
