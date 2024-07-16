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

type enrichDetailsTrackingRepository struct {
	gormDb *gorm.DB
}

type EnrichDetailsTrackingRepository interface {
	RegisterRequest(ctx context.Context, request entity.EnrichDetailsTracking) error
	GetByIP(ctx context.Context, IP string) (*entity.EnrichDetailsTracking, error)
}

func NewEnrichDetailsTrackingRepository(gormDb *gorm.DB) EnrichDetailsTrackingRepository {
	return &enrichDetailsTrackingRepository{gormDb: gormDb}
}

func (r enrichDetailsTrackingRepository) RegisterRequest(ctx context.Context, request entity.EnrichDetailsTracking) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsTrackingRepository.RegisterRequest")
	defer span.Finish()

	err := r.gormDb.Create(&request).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r enrichDetailsTrackingRepository) GetByIP(ctx context.Context, ip string) (*entity.EnrichDetailsTracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsTrackingRepository.GetByIP")
	defer span.Finish()
	span.LogFields(tracingLog.String("ip", ip))

	var entity *entity.EnrichDetailsTracking
	err := r.gormDb.
		Where("ip = ?", ip).
		First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return entity, err
}
