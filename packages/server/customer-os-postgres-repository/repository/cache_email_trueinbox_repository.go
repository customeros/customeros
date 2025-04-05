package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type CacheEmailTrueinboxRepository interface {
	Create(ctx context.Context, record postgres_entity.CacheEmailTrueinbox) (*postgres_entity.CacheEmailTrueinbox, error)
	GetLatestByEmail(ctx context.Context, email string) (*postgres_entity.CacheEmailTrueinbox, error)
	GetAllByEmail(ctx context.Context, email string) ([]postgres_entity.CacheEmailTrueinbox, error)
}

type cacheEmailTrueinboxRepository struct {
	db *gorm.DB
}

func NewCacheEmailTrueinboxRepository(gormDb *gorm.DB) CacheEmailTrueinboxRepository {
	return &cacheEmailTrueinboxRepository{db: gormDb}
}

func (r cacheEmailTrueinboxRepository) GetAllByEmail(ctx context.Context, email string) ([]postgres_entity.CacheEmailTrueinbox, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheEmailTrueinboxRepository.GetAllByEmail")
	defer spans.Finish()
	spans.LogKV("email", email)

	var records []postgres_entity.CacheEmailTrueinbox
	err := r.db.Where("email = ?", email).Order("created_at desc").Find(&records).Error
	if err != nil {
		return nil, err
	}

	spans.LogKV("result.count", len(records))

	return records, nil
}

func (r cacheEmailTrueinboxRepository) GetLatestByEmail(ctx context.Context, email string) (*postgres_entity.CacheEmailTrueinbox, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheEmailTrueinboxRepository.GetLatestByEmail")
	defer spans.Finish()
	spans.LogKV("email", email)

	var data postgres_entity.CacheEmailTrueinbox
	err := r.db.Where("email = ?", email).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		spans.TraceError(err)
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r cacheEmailTrueinboxRepository) Create(ctx context.Context, record postgres_entity.CacheEmailTrueinbox) (*postgres_entity.CacheEmailTrueinbox, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheEmailTrueinboxRepository.Create")
	defer spans.Finish()
	spans.LogObjectAsJson("record", record)

	record.CreatedAt = utils.Now()
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &record, nil
}
