package postgres_repository

import (
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
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
	spans, ctx := telemetry.StartPostgresSpan(ctx, "EnrichDetailsTrackingRepository.RegisterRequest")
	defer spans.Finish()

	record, err := r.GetByIP(ctx, request.IP)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if record == nil {
		// Create new record
		if err := r.gormDb.WithContext(ctx).Create(&request).Error; err != nil {
			spans.TraceError(err)
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
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r enrichDetailsTrackingRepository) GetByIP(ctx context.Context, ip string) (*postgres_entity.EnrichDetailsTracking, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EnrichDetailsTrackingRepository.GetByIP")
	defer spans.Finish()

	spans.LogKV("ip", ip)

	var postgres_entity *postgres_entity.EnrichDetailsTracking
	err := r.gormDb.
		Where("ip = ?", ip).
		First(&postgres_entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		spans.LogKV("result.found", false)
		return nil, nil
	}

	if err != nil {
		spans.TraceError(err)
	}

	spans.LogKV("result.found", true)

	return postgres_entity, err
}
