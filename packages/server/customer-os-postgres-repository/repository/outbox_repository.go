package postgres_repository

import (
	"context"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"gorm.io/gorm"
)

type OutboxRepository interface {
	Create(ctx context.Context, event *postgres_entity.OutboxEvent) (*postgres_entity.OutboxEvent, error)
	CreateWithTxn(ctx context.Context, txn *gorm.DB, event *postgres_entity.OutboxEvent) error
	GetPendingEvents(ctx context.Context, limit int) ([]*postgres_entity.OutboxEvent, error)
	MarkAsProcessing(ctx context.Context, id string, lockDuration time.Duration) error
	MarkAsCompleted(ctx context.Context, id string) error
	MarkAsFailed(ctx context.Context, id string, errorMessage string) error
	IncrementRetryCount(ctx context.Context, id string) error
	DeleteProcessedEvents(ctx context.Context, olderThan time.Duration, limit int64) (int64, error)
}

type outboxRepository struct {
	gormDb *gorm.DB
}

func NewOutboxRepository(gormDb *gorm.DB) OutboxRepository {
	return &outboxRepository{
		gormDb: gormDb,
	}
}

func (r *outboxRepository) Create(ctx context.Context, event *postgres_entity.OutboxEvent) (*postgres_entity.OutboxEvent, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "outboxRepository.Create")
	defer spans.Finish()

	err := r.gormDb.WithContext(ctx).Create(event).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return event, nil
}

func (r *outboxRepository) CreateWithTxn(ctx context.Context, txn *gorm.DB, event *postgres_entity.OutboxEvent) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "outboxRepository.CreateWithTxn")
	defer spans.Finish()

	err := txn.WithContext(ctx).Create(event).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *outboxRepository) GetPendingEvents(ctx context.Context, limit int) ([]*postgres_entity.OutboxEvent, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "outboxRepository.GetPendingEvents")
	defer spans.Finish()
	spans.LogKV("limit", limit)

	var events []*postgres_entity.OutboxEvent

	err := r.gormDb.WithContext(ctx).
		Where("status = ? AND (lock_until IS NULL OR lock_until < ?)",
			postgres_entity.OutboxPending, time.Now()).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(events))
	return events, err
}

func (r *outboxRepository) MarkAsProcessing(ctx context.Context, id string, lockDuration time.Duration) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "outboxRepository.MarkAsProcessing")
	defer spans.Finish()
	spans.LogKV("id", id)

	lockUntil := time.Now().Add(lockDuration)

	result := r.gormDb.WithContext(ctx).
		Model(&postgres_entity.OutboxEvent{}).
		Where("id = ? AND (lock_until IS NULL OR lock_until < ?)", id, time.Now()).
		Updates(map[string]interface{}{
			"status":     postgres_entity.OutboxProcessing,
			"lock_until": lockUntil,
		})

	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *outboxRepository) MarkAsCompleted(ctx context.Context, id string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "outboxRepository.MarkAsCompleted")
	defer spans.Finish()
	spans.LogKV("id", id)

	now := utils.Now()

	err := r.gormDb.WithContext(ctx).
		Model(&postgres_entity.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       postgres_entity.OutboxCompleted,
			"processed_at": now,
			"lock_until":   nil,
		}).Error

	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *outboxRepository) MarkAsFailed(ctx context.Context, id string, errorMessage string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "outboxRepository.MarkAsFailed")
	defer spans.Finish()
	spans.LogKV("errorMessage", errorMessage, "id", id)

	err := r.gormDb.WithContext(ctx).
		Model(&postgres_entity.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        postgres_entity.OutboxFailed,
			"error_message": errorMessage,
			"lock_until":    nil,
		}).Error

	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}

func (r *outboxRepository) IncrementRetryCount(ctx context.Context, id string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "outboxRepository.IncrementRetryCount")
	defer spans.Finish()

	err := r.gormDb.WithContext(ctx).
		Model(&postgres_entity.OutboxEvent{}).
		Where("id = ?", id).
		UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *outboxRepository) DeleteProcessedEvents(ctx context.Context, olderThan time.Duration, limit int64) (int64, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "outboxRepository.DeleteProcessedEvents")
	defer spans.Finish()

	cutoffTime := time.Now().Add(-olderThan)

	result := r.gormDb.WithContext(ctx).
		Where("status = ? AND processed_at < ?",
			postgres_entity.OutboxCompleted, cutoffTime).
		Limit(int(limit)).
		Delete(&postgres_entity.OutboxEvent{})

	if result.Error != nil {
		spans.TraceError(result.Error)
		return 0, result.Error
	}

	spans.LogKV("result.rowsDeleted", result.RowsAffected)

	return result.RowsAffected, nil
}
