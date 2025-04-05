package postgres_repository

import (
	"context"
	"strings"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type WebTrackerEventsRepository interface {
	Create(ctx context.Context, webWebTrackerData postgres_entity.WebTrackerEvents) (*postgres_entity.WebTrackerEvents, error)
	FindAll(ctx context.Context, webWebTrackerData postgres_entity.WebTrackerEvents, cacheLookbackInDays *int) ([]postgres_entity.WebTrackerEvents, error)
	FindEventsForPageVisit(ctx context.Context, sessionId, page string) ([]postgres_entity.WebTrackerEvents, error)
}

type webTrackerEventsRepository struct {
	gormDb *gorm.DB
}

func NewWebTrackerEventsRepository(gormDb *gorm.DB) WebTrackerEventsRepository {
	return &webTrackerEventsRepository{gormDb: gormDb}
}

func (r *webTrackerEventsRepository) Create(ctx context.Context, webTrackerData postgres_entity.WebTrackerEvents) (*postgres_entity.WebTrackerEvents, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "WebTrackerEventsRepository.Create")
	defer spans.Finish()
	spans.LogKV("webTrackerData", webTrackerData)

	var created postgres_entity.WebTrackerEvents
	err := r.gormDb.Create(&webTrackerData).Scan(&created).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &created, nil
}

func (r *webTrackerEventsRepository) FindAll(ctx context.Context, webTrackerData postgres_entity.WebTrackerEvents, cacheLookbackInDays *int) ([]postgres_entity.WebTrackerEvents, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "WebTrackerEventsRepository.FindAll")
	defer spans.Finish()

	query := r.gormDb.Where(&webTrackerData).Order("created_at DESC")

	// Add lookback period if provided
	if cacheLookbackInDays != nil {
		lookbackDate := time.Now().AddDate(0, 0, -*cacheLookbackInDays)
		query = query.Where("created_at > ?", lookbackDate)
	}

	var results []postgres_entity.WebTrackerEvents
	err := query.Find(&results).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return results, nil
}

func (r *webTrackerEventsRepository) FindEventsForPageVisit(ctx context.Context, sessionId, page string) ([]postgres_entity.WebTrackerEvents, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "WebTrackerEventsRepository.FindEventsForPageVisit")
	defer spans.Finish()
	spans.LogKV("sessionId", sessionId, "page", page)

	if sessionId == "" {
		return nil, errors.New("sessionId cannot be empty")
	}
	if page == "" {
		return nil, errors.New("page cannot be empty")
	}

	pathname := "/"
	domain := utils.ExtractDomain(page)
	pageParts := strings.SplitAfter(page, domain)
	if len(pageParts) > 1 {
		pathname = pageParts[1]
		if !strings.HasSuffix(pathname, "/") {
			pathname += "/"
		}
	}
	spans.LogKV("pathname", pathname)

	var results []postgres_entity.WebTrackerEvents
	err := r.gormDb.
		Where("session_id = ?", sessionId).
		Where("pathname = ?", pathname).
		Order("timestamp ASC").
		Find(&results).
		Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(results))
	return results, nil
}
