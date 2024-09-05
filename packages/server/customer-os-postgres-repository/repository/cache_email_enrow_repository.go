package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

type CacheEmailEnrowRepository interface {
	RegisterRequest(ctx context.Context, record entity.CacheEmailEnrow) (*entity.CacheEmailEnrow, error)
	AddResponse(ctx context.Context, requestId, qualification, response string) error
	GetLatestByEmail(ctx context.Context, email string) (*entity.CacheEmailEnrow, error)
	GetAllByEmail(ctx context.Context, email string) ([]entity.CacheEmailEnrow, error)
	GetWithoutResponses(ctx context.Context) ([]*entity.CacheEmailEnrow, error)
}

type cacheEmailEnrowRepository struct {
	db *gorm.DB
}

func NewCacheEmailEnrowRepository(gormDb *gorm.DB) CacheEmailEnrowRepository {
	return &cacheEmailEnrowRepository{db: gormDb}
}

func (r cacheEmailEnrowRepository) RegisterRequest(ctx context.Context, record entity.CacheEmailEnrow) (*entity.CacheEmailEnrow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailEnrowRepository.RegisterRequest")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "record", record)

	now := utils.Now()
	record.CreatedAt = now
	record.UpdatedAt = now
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (r cacheEmailEnrowRepository) AddResponse(ctx context.Context, requestId, qualification, response string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailEnrowRepository.AddResponse")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("requestId", requestId)

	// Add response to the request with the given requestId, empty response and latest by created_at
	err := r.db.
		Model(&entity.CacheEmailEnrow{}).
		Where("request_id = ?", requestId).
		Where("response = ?", "").
		Order("created_at desc").
		Limit(1).
		UpdateColumn("data", response).
		UpdateColumn("qualification", qualification).
		UpdateColumn("updated_at", utils.Now()).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (r cacheEmailEnrowRepository) GetAllByEmail(ctx context.Context, email string) ([]entity.CacheEmailEnrow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailEnrowRepository.GetAllByEmail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("email", email)

	var records []entity.CacheEmailEnrow
	err := r.db.Where("email = ?", email).Order("created_at desc").Find(&records).Error
	if err != nil {
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(records)))

	return records, nil
}

func (r cacheEmailEnrowRepository) GetLatestByEmail(ctx context.Context, email string) (*entity.CacheEmailEnrow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailEnrowRepository.GetLatestByEmail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("email", email)

	var data entity.CacheEmailEnrow
	err := r.db.Where("email = ?", email).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r cacheEmailEnrowRepository) GetWithoutResponses(ctx context.Context) ([]*entity.CacheEmailEnrow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.GetWithoutResponses")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var entity []*entity.CacheEmailEnrow
	err := r.db.
		Where("response = ?", "").
		Where("created_at < ?", utils.Now().Add(-10*time.Minute)).
		Limit(50).
		Find(&entity).Error

	return entity, err
}
