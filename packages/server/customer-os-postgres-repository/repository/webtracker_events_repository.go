package postgres_repository

import (
	"context"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebTrackerEventsRepository.Create")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "webTrackerData", webTrackerData)

	var created postgres_entity.WebTrackerEvents
	err := r.gormDb.Create(&webTrackerData).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &created, nil
}

func (r *webTrackerEventsRepository) FindAll(ctx context.Context, webTrackerData postgres_entity.WebTrackerEvents, cacheLookbackInDays *int) ([]postgres_entity.WebTrackerEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebTrackerEventsRepository.FindAll")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	query := r.gormDb.Where(&webTrackerData).Order("created_at DESC")

	// Add lookback period if provided
	if cacheLookbackInDays != nil {
		lookbackDate := time.Now().AddDate(0, 0, -*cacheLookbackInDays)
		query = query.Where("created_at > ?", lookbackDate)
	}

	var results []postgres_entity.WebTrackerEvents
	err := query.Find(&results).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return results, nil
}

func (r *webTrackerEventsRepository) FindEventsForPageVisit(ctx context.Context, sessionId, page string) ([]postgres_entity.WebTrackerEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebTrackerEventsRepository.FindEventsForPageVisit")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogKV("sessionId", sessionId, "page", page)

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
	span.LogKV("pathname", pathname)

	var results []postgres_entity.WebTrackerEvents
	err := r.gormDb.
		Where("session_id = ?", sessionId).
		Where("pathname = ?", pathname).
		Order("timestamp ASC").
		Find(&results).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(results)))
	return results, nil
}
