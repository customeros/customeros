package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TrackingRepository interface {
	GetById(ctx context.Context, id string) (*entity.Tracking, error)
	GetNotIdentifiedRecords(ctx context.Context) ([]*entity.Tracking, error)

	Store(ctx context.Context, tracking entity.Tracking) (string, error)
	MarkAsIdentified(ctx context.Context, id string) error
}

type trackingRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewTrackingRepository(gormDb *gorm.DB) TrackingRepository {
	return &trackingRepositoryImpl{gormDb: gormDb}
}

func (r trackingRepositoryImpl) GetById(ctx context.Context, id string) (*entity.Tracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.GetNotIdentifiedTrackingRecords")
	defer span.Finish()

	var result entity.Tracking
	err := r.gormDb.Model(&entity.Tracking{}).Find(&result, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &result, nil
}

func (r trackingRepositoryImpl) GetNotIdentifiedRecords(ctx context.Context) ([]*entity.Tracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.GetNotIdentifiedTrackingRecords")
	defer span.Finish()

	var entities []*entity.Tracking
	err := r.gormDb.
		Where("identified = ?", false).
		Distinct("ip", "id").
		Limit(50).
		Find(&entities).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return entities, nil
}

func (repo *trackingRepositoryImpl) Store(ctx context.Context, tracking entity.Tracking) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.Store")
	defer span.Finish()

	err := repo.gormDb.Save(&tracking).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return tracking.ID, nil
}

func (r trackingRepositoryImpl) MarkAsIdentified(ctx context.Context, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.MarkAsIdentified")
	defer span.Finish()
	span.LogFields(tracingLog.String("id", id))

	byId, err := r.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if byId == nil {
		tracing.TraceErr(span, errors.New("record not found"))
		return nil
	}

	byId.Identified = true

	err = r.gormDb.Save(byId).Error
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}
