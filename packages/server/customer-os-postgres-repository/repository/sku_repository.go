package postgres_repository

import (
	"errors"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type SkuRepository interface {
	Get(ctx context.Context, tenant, id string) (*postgresentity.SkuEntity, error)
	GetAll(ctx context.Context, tenant string, archived *bool) ([]*postgresentity.SkuEntity, error)
	Save(ctx context.Context, sku *postgresentity.SkuEntity) (*postgresentity.SkuEntity, error)
	Archive(ctx context.Context, tenant, id string) error
}

type skuRepository struct {
	db *gorm.DB
}

func NewSkuRepository(db *gorm.DB) SkuRepository {
	return &skuRepository{
		db: db,
	}
}

func (repo *skuRepository) Get(ctx context.Context, tenant, skuId string) (*postgresentity.SkuEntity, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "SkuRepository.Get")
	defer spans.Finish()
	spans.LogKV("skuId", skuId)

	if skuId == "" {
		return nil, nil
	}

	var existing *postgresentity.SkuEntity
	err := repo.db.First(&existing, "tenant = ? and id = ?", tenant, skuId).Error
	if err != nil {
		spans.LogKV("result.found", false)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	if existing == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}

	spans.LogKV("result.found", true)
	return existing, nil
}

func (repo *skuRepository) GetAll(ctx context.Context, tenant string, archived *bool) ([]*postgresentity.SkuEntity, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "SkuRepository.GetAll")
	defer spans.Finish()

	var err error
	var existing []*postgresentity.SkuEntity

	if archived == nil {
		err = repo.db.Find(&existing, "tenant = ?", tenant).Error
	} else {
		err = repo.db.Find(&existing, "tenant = ? and archived = ?", tenant, archived).Error
	}

	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return nil, err
	}

	return existing, nil
}

func (repo *skuRepository) Save(ctx context.Context, sku *postgresentity.SkuEntity) (*postgresentity.SkuEntity, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "SkuRepository.Save")
	defer spans.Finish()

	result := repo.db.Save(sku)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, fmt.Errorf("saving slack settings failed: %w", result.Error)
	}
	return sku, nil
}

func (repo *skuRepository) Archive(ctx context.Context, tenant, id string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "SkuRepository.Archive")
	defer spans.Finish()
	spans.LogKV("id", id)

	existing, err := repo.Get(ctx, tenant, id)
	if err != nil {
		return err
	}

	existing.Archived = true

	result := repo.db.Save(existing)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return fmt.Errorf("archiving sku failed: %w", result.Error)
	}

	return nil
}
