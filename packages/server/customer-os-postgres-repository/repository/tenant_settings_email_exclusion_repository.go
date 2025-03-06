package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TenantSettingsEmailExclusionRepository interface {
	GetExclusionList(ctx context.Context) ([]postgres_entity.TenantSettingsEmailExclusion, error)
}

type tenantSettingsEmailExclusionRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewEmailExclusionRepository(gormDb *gorm.DB) TenantSettingsEmailExclusionRepository {
	return &tenantSettingsEmailExclusionRepositoryImpl{gormDb: gormDb}
}

func (repo *tenantSettingsEmailExclusionRepositoryImpl) GetExclusionList(ctx context.Context) ([]postgres_entity.TenantSettingsEmailExclusion, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsEmailExclusionRepository.GetExclusionList")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	result := []postgres_entity.TenantSettingsEmailExclusion{}
	err := repo.gormDb.Find(&result).Limit(5000).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result)))

	return result, nil
}
