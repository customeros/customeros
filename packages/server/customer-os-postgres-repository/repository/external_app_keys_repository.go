package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository/helper"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type externalAppKeysRepository struct {
	gormDb *gorm.DB
}

type ExternalAppKeysRepository interface {
	IncrementUsageCount(ctx context.Context, id uint64) helper.QueryResult
	GetAppKeys(ctx context.Context, app, group string, usageLimit int) helper.QueryResult
}

func NewExternalAppKeysRepository(gormDb *gorm.DB) ExternalAppKeysRepository {
	return &externalAppKeysRepository{gormDb: gormDb}
}

func (e externalAppKeysRepository) GetAppKeys(ctx context.Context, app, group string, usageLimit int) helper.QueryResult {
	spans, _ := telemetry.StartPostgresSpan(ctx, "ExternalAppKeysRepository.GetAppKeys")
	defer spans.Finish()

	spans.LogKV("app", app, "group1", group)

	var appKeys []postgres_entity.ExternalAppKeys
	err := e.gormDb.
		Where("app = ? AND group1 = ? AND usage_count < ?", app, group, usageLimit).
		Find(&appKeys).Limit(10).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: appKeys}
}

func (e externalAppKeysRepository) IncrementUsageCount(ctx context.Context, id uint64) helper.QueryResult {
	spans, _ := telemetry.StartPostgresSpan(ctx, "ExternalAppKeysRepository.IncrementUsageCount")
	defer spans.Finish()

	spans.LogKV("id", id)

	// create entry if not exists
	appKey := postgres_entity.ExternalAppKeys{
		ID: id,
	}
	err := e.gormDb.
		Where("id = ?", id).
		First(&appKey).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	// increment usage_count
	err = e.gormDb.
		Model(&postgres_entity.ExternalAppKeys{}).
		Where("id = ?", id).
		UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).
		UpdateColumn("updated_at", utils.Now()).
		Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: appKey}
}
