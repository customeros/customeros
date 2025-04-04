package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
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
	spans, _ := telemetry.StartPostgresSpan(ctx, "TenantSettingsEmailExclusionRepository.GetExclusionList")
	defer spans.Finish()

	result := []postgres_entity.TenantSettingsEmailExclusion{}
	err := repo.gormDb.Find(&result).Limit(5000).Error

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(result))

	return result, nil
}
