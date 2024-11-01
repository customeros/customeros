package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TenantSettingsEmailExclusionRepository interface {
	GetExclusionList(ctx context.Context) ([]entity.TenantSettingsEmailExclusion, error)
}

type tenantSettingsEmailExclusionRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewEmailExclusionRepository(gormDb *gorm.DB) TenantSettingsEmailExclusionRepository {
	return &tenantSettingsEmailExclusionRepositoryImpl{gormDb: gormDb}
}

func (repo *tenantSettingsEmailExclusionRepositoryImpl) GetExclusionList(ctx context.Context) ([]entity.TenantSettingsEmailExclusion, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsEmailExclusionRepository.GetExclusionList")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	result := []entity.TenantSettingsEmailExclusion{}
	err := repo.gormDb.Find(&result).Limit(5000).Error

	if err != nil {
		logrus.Errorf("error while getting personal email provider list: %v", err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result)))

	return result, nil
}
