package postgres_repository

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
	"time"
)

type CacheEmailScrubbyRepository interface {
	Save(ctx context.Context, cacheEmailScrubby postgres_entity.CacheEmailScrubby) (*postgres_entity.CacheEmailScrubby, error)
	GetAllByEmail(ctx context.Context, email string) ([]postgres_entity.CacheEmailScrubby, error)
	GetLatestByEmail(ctx context.Context, email string) (*postgres_entity.CacheEmailScrubby, error)
	SetStatus(ctx context.Context, email, status string) error
	SetJustChecked(ctx context.Context, id string) (*postgres_entity.CacheEmailScrubby, error)
	GetToCheck(ctx context.Context, delayFromPreviousCheckInHours, limit int) ([]postgres_entity.CacheEmailScrubby, error)
}

type cacheEmailScrubbyRepository struct {
	db *gorm.DB
}

func NewCacheEmailScrubbyRepository(gormDb *gorm.DB) CacheEmailScrubbyRepository {
	return &cacheEmailScrubbyRepository{db: gormDb}
}

func (r *cacheEmailScrubbyRepository) Save(ctx context.Context, cacheEmailScrubby postgres_entity.CacheEmailScrubby) (*postgres_entity.CacheEmailScrubby, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheEmailScrubbyRepository.Save")
	defer spans.Finish()
	spans.LogObjectAsJson("cacheEmailScrubby", cacheEmailScrubby)

	now := utils.Now()
	cacheEmailScrubby.CreatedAt = now
	cacheEmailScrubby.CheckedAt = now

	if err := r.db.WithContext(ctx).Save(&cacheEmailScrubby).Error; err != nil {
		return nil, err
	}

	return &cacheEmailScrubby, nil
}

func (r *cacheEmailScrubbyRepository) GetAllByEmail(ctx context.Context, email string) ([]postgres_entity.CacheEmailScrubby, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheEmailScrubbyRepository.GetAllByEmail")
	defer spans.Finish()
	spans.LogKV("email", email)

	var cacheEmailScrubbys []postgres_entity.CacheEmailScrubby
	result := r.db.WithContext(ctx).Where("email = ?", email).Order("created_at desc").Find(&cacheEmailScrubbys)

	if result.Error != nil {
		return nil, result.Error
	}

	return cacheEmailScrubbys, nil
}

func (r *cacheEmailScrubbyRepository) GetLatestByEmail(ctx context.Context, email string) (*postgres_entity.CacheEmailScrubby, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheEmailScrubbyRepository.GetLatestByEmail")
	defer spans.Finish()
	spans.LogKV("email", email)

	var cacheEmailScrubby postgres_entity.CacheEmailScrubby
	result := r.db.WithContext(ctx).Where("email = ?", email).Order("created_at desc").First(&cacheEmailScrubby)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			spans.TraceError(result.Error)
			return nil, result.Error
		}
	}

	return &cacheEmailScrubby, nil
}

func (r *cacheEmailScrubbyRepository) SetStatus(ctx context.Context, email, status string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheEmailScrubbyRepository.SetStatus")
	defer spans.Finish()
	spans.LogKV("email", email)
	spans.LogKV("status", status)

	result := r.db.WithContext(ctx).
		Model(&postgres_entity.CacheEmailScrubby{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"status":     status,
			"checked_at": utils.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	// If no records were affected, it's not considered an error
	if result.RowsAffected == 0 {
		spans.LogKV("message", "No records found for the given email")
	}

	return nil
}

func (r *cacheEmailScrubbyRepository) SetJustChecked(ctx context.Context, id string) (*postgres_entity.CacheEmailScrubby, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheEmailScrubbyRepository.SetJustChecked")
	defer spans.Finish()
	spans.LogKV("id", id)

	var cacheEmailScrubby postgres_entity.CacheEmailScrubby
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&cacheEmailScrubby)

	if result.Error != nil {
		return nil, result.Error
	}

	cacheEmailScrubby.CheckedAt = utils.Now()

	if err := r.db.WithContext(ctx).Save(&cacheEmailScrubby).Error; err != nil {
		return nil, err
	}

	return &cacheEmailScrubby, nil
}

func (r *cacheEmailScrubbyRepository) GetToCheck(ctx context.Context, delayFromPreviousCheckInHours, limit int) ([]postgres_entity.CacheEmailScrubby, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheEmailScrubbyRepository.GetToCheck")
	defer spans.Finish()
	spans.LogKV("delayFromPreviousCheckInHours", delayFromPreviousCheckInHours)
	spans.LogKV("limit", limit)

	var cacheEmailScrubbys []postgres_entity.CacheEmailScrubby

	// Calculate the cutoff time
	cutoffTime := time.Now().Add(-time.Duration(delayFromPreviousCheckInHours) * time.Hour)

	result := r.db.WithContext(ctx).
		Where("status = ? AND checked_at < ?", string(postgres_entity.ScrubbyStatusPending), cutoffTime).
		Order("checked_at asc").
		Limit(limit).
		Find(&cacheEmailScrubbys)

	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	return cacheEmailScrubbys, nil
}
