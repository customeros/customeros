package repository

import (
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type enrichDetailsPrefilterTrackingRepository struct {
	gormDb *gorm.DB
}

type EnrichDetailsPrefilterTrackingRepository interface {
	RegisterRequest(ctx context.Context, request entity.EnrichDetailsPreFilterTracking) error
	GetByIP(ctx context.Context, IP string) (*entity.EnrichDetailsPreFilterTracking, error)
}

func NewEnrichDetailsPrefilterTrackingRepository(gormDb *gorm.DB) EnrichDetailsPrefilterTrackingRepository {
	return &enrichDetailsPrefilterTrackingRepository{gormDb: gormDb}
}

func (r enrichDetailsPrefilterTrackingRepository) RegisterRequest(ctx context.Context, request entity.EnrichDetailsPreFilterTracking) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsPrefilterTrackingRepository.RegisterRequest")
	defer span.Finish()

	err := r.gormDb.Create(&request).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r enrichDetailsPrefilterTrackingRepository) GetByIP(ctx context.Context, ip string) (*entity.EnrichDetailsPreFilterTracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsPrefilterTrackingRepository.GetByIP")
	defer span.Finish()
	span.LogFields(tracingLog.String("ip", ip))

	var entity *entity.EnrichDetailsPreFilterTracking
	err := r.gormDb.
		Where("ip = ?", ip).
		First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		span.LogFields(tracingLog.Bool("result.found", false))
		return nil, nil
	}

	if err != nil {
		tracing.TraceErr(span, err)
	}

	span.LogFields(tracingLog.Bool("result.found", true))

	return entity, err
}
