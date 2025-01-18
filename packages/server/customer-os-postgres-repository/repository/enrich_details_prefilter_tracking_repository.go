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

type enrichDetailsPrefilterTrackingRepository struct {
	gormDb *gorm.DB
}

type EnrichDetailsPrefilterTrackingRepository interface {
	GetForSendingRequests(ctx context.Context) ([]*postgres_entity.EnrichDetailsPreFilterTracking, error)
	GetByIP(ctx context.Context, IP string) (*postgres_entity.EnrichDetailsPreFilterTracking, error)

	RegisterRequest(ctx context.Context, ip string) error
	RegisterResponse(ctx context.Context, ip string, shouldIdentify bool, skipIdenitifyReason, response string) error
}

func NewEnrichDetailsPrefilterTrackingRepository(gormDb *gorm.DB) EnrichDetailsPrefilterTrackingRepository {
	return &enrichDetailsPrefilterTrackingRepository{gormDb: gormDb}
}

func (r enrichDetailsPrefilterTrackingRepository) GetForSendingRequests(ctx context.Context) ([]*postgres_entity.EnrichDetailsPreFilterTracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsPrefilterTrackingRepository.GetForSendingRequests")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var entitites []*postgres_entity.EnrichDetailsPreFilterTracking
	err := r.gormDb.
		Where("response is null").
		Limit(500).
		Find(&entitites).Error

	if err != nil {
		tracing.TraceErr(span, err)
	}

	span.LogFields(tracingLog.Int("result.count", len(entitites)))

	return entitites, err
}

func (r enrichDetailsPrefilterTrackingRepository) GetByIP(ctx context.Context, ip string) (*postgres_entity.EnrichDetailsPreFilterTracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsPrefilterTrackingRepository.GetByIP")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("ip", ip))

	var postgres_entity *postgres_entity.EnrichDetailsPreFilterTracking
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

func (r enrichDetailsPrefilterTrackingRepository) RegisterRequest(ctx context.Context, ip string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsPrefilterTrackingRepository.RegisterRequest")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("ip", ip))

	request := postgres_entity.EnrichDetailsPreFilterTracking{
		CreatedAt: utils.Now(),
		IP:        ip,
	}

	err := r.gormDb.Create(&request).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r enrichDetailsPrefilterTrackingRepository) RegisterResponse(ctx context.Context, ip string, shouldIdentify bool, skipIdenitifyReason, response string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsPrefilterTrackingRepository.RegisterResponse")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("ip", ip), tracingLog.Bool("shouldIdentify", shouldIdentify), tracingLog.String("response", response), tracingLog.String("skipIdenitifyReason", skipIdenitifyReason))

	byId, err := r.GetByIP(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	byId.ShouldIdentify = &shouldIdentify
	byId.Response = &response
	byId.SkipIdentifyReason = skipIdenitifyReason

	err = r.gormDb.Save(byId).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
