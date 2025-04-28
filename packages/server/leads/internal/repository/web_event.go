package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/leads/internal/database"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
)

type WebTrackerEventRepository interface {
	Create(ctx context.Context, event *models.WebTrackerEvent) error
	FindBySessionID(ctx context.Context, sessionID string, limit int) ([]*models.WebTrackerEvent, error)
	FindByTenant(ctx context.Context, tenant string, from, to time.Time, limit int) ([]*models.WebTrackerEvent, error)
}

type webTrackerEventRepository struct {
	read  *gorm.DB
	write *gorm.DB
}

func NewWebTrackerEventRepository(db *database.DbConnections) WebTrackerEventRepository {
	return &webTrackerEventRepository{
		read:  db.ReadDB,
		write: db.WriteDB,
	}
}

func (r *webTrackerEventRepository) Create(ctx context.Context, event *models.WebTrackerEvent) error {
	span, ctx := telemetry.StartPostgresSpan(ctx, "webTrackerEventRepository.Create")
	defer span.Finish()

	// Generate ID if not provided
	if event.ID == "" {
		event.ID = utils.GenerateNanoIDWithPrefix("wevt", 16)
	}

	return r.write.WithContext(ctx).Create(event).Error
}

func (r *webTrackerEventRepository) FindBySessionID(ctx context.Context, sessionID string, limit int) ([]*models.WebTrackerEvent, error) {
	span, ctx := telemetry.StartPostgresSpan(ctx, "webTrackerEventRepository.FindBySessionID")
	defer span.Finish()

	var events []*models.WebTrackerEvent
	err := r.read.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("timestamp ASC").
		Limit(limit).
		Find(&events).Error

	return events, err
}

func (r *webTrackerEventRepository) FindByTenant(ctx context.Context, tenant string, from, to time.Time, limit int) ([]*models.WebTrackerEvent, error) {
	span, ctx := telemetry.StartPostgresSpan(ctx, "webTrackerEventRepository.FindByTenant")
	defer span.Finish()

	var events []*models.WebTrackerEvent
	err := r.read.WithContext(ctx).
		Where("tenant = ? AND timestamp BETWEEN ? AND ?", tenant, from, to).
		Order("timestamp DESC").
		Limit(limit).
		Find(&events).Error

	return events, err
}
