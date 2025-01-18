package postgres_repository

import (
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type enrichDetailsTrackingRepository struct {
	gormDb *gorm.DB
}

type EnrichDetailsTrackingRepository interface {
	Save(ctx context.Context, request postgres_entity.EnrichDetailsTracking) error
	GetByIP(ctx context.Context, IP string) (*postgres_entity.EnrichDetailsTracking, error)
}

func NewEnrichDetailsTrackingRepository(gormDb *gorm.DB) EnrichDetailsTrackingRepository {
	return &enrichDetailsTrackingRepository{gormDb: gormDb}
}

func (r enrichDetailsTrackingRepository) Save(ctx context.Context, request postgres_entity.EnrichDetailsTracking) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EnrichDetailsTrackingRepository.RegisterRequest")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	record, err := r.GetByIP(ctx, request.IP)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if record == nil {
		// Create new record
		if err := r.gormDb.WithContext(ctx).Create(&request).Error; err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return nil
	}

	// Update existing record
	request.UpdatedAt = utils.Now()
	if err := r.gormDb.WithContext(ctx).
		Model(&postgres_entity.EnrichDetailsTracking{}).
		Where("ip = ?", request.IP).
		Updates(request).Error; err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r enrichDetailsTrackingRepository) GetByIP(ctx context.Context, ip string) (*postgres_entity.EnrichDetailsTracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsTrackingRepository.GetByIP")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("ip", ip))

	var postgres_entity *postgres_entity.EnrichDetailsTracking
	err := r.gormDb.
		Where("ip = ?", ip).
		First(&postgres_entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		span.LogFields(tracingLog.Bool("result.found", false))
		return nil, nil
	}

	if err != nil {
		tracing.TraceErr(span, err)
	}

	span.LogFields(tracingLog.Bool("result.found", true))

	return postgres_entity, err
}
