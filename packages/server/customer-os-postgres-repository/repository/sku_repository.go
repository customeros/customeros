package postgres_repository

import (
	"errors"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
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
	span, _ := opentracing.StartSpanFromContext(ctx, "SkuRepository.Get")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(tracingLog.String("skuId", skuId))

	if skuId == "" {
		return nil, nil
	}

	var existing *postgresentity.SkuEntity
	err := repo.db.First(&existing, "tenant = ? and skuId = ?", tenant, skuId).Error
	if err != nil {
		span.LogFields(tracingLog.Bool("result.found", false))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	if existing == nil {
		span.LogFields(tracingLog.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(tracingLog.Bool("result.found", true))
	return existing, nil
}

func (repo *skuRepository) GetAll(ctx context.Context, tenant string, archived *bool) ([]*postgresentity.SkuEntity, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SkuRepository.GetAll")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	var err error
	var existing []*postgresentity.SkuEntity

	if archived == nil {
		err = repo.db.Find(&existing, "tenant = ?", tenant).Error
	} else {
		err = repo.db.Find(&existing, "tenant = ? and archived = ?", tenant, archived).Error
	}

	if err != nil {
		span.LogFields(tracingLog.Bool("result.found", false))
		tracing.TraceErr(span, err)
		return nil, err
	}

	return existing, nil
}

func (repo *skuRepository) Save(ctx context.Context, sku *postgresentity.SkuEntity) (*postgresentity.SkuEntity, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SkuRepository.Save")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	result := repo.db.Save(sku)
	if result.Error != nil {
		return nil, fmt.Errorf("saving slack settings failed: %w", result.Error)
	}
	return sku, nil
}

func (repo *skuRepository) Archive(ctx context.Context, tenant, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SkuRepository.Archive")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.String("id", id))

	existing, err := repo.Get(ctx, tenant, id)
	if err != nil {
		return err
	}

	existing.Archived = true

	result := repo.db.Save(existing)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return fmt.Errorf("archiving sku failed: %w", result.Error)
	}

	return nil
}
