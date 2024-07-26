package repository

import (
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

type enrichDetailsPrefilterTrackingRepository struct {
	gormDb *gorm.DB
}

type EnrichDetailsPrefilterTrackingRepository interface {
	GetForSendingRequests(ctx context.Context) ([]*entity.EnrichDetailsPreFilterTracking, error)
	GetByIP(ctx context.Context, IP string) (*entity.EnrichDetailsPreFilterTracking, error)

	RegisterRequest(ctx context.Context, ip string) error
	RegisterResponse(ctx context.Context, ip string, shouldIdentify bool, response string) error
}

func NewEnrichDetailsPrefilterTrackingRepository(gormDb *gorm.DB) EnrichDetailsPrefilterTrackingRepository {
	return &enrichDetailsPrefilterTrackingRepository{gormDb: gormDb}
}

func (r enrichDetailsPrefilterTrackingRepository) GetForSendingRequests(ctx context.Context) ([]*entity.EnrichDetailsPreFilterTracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsPrefilterTrackingRepository.GetForSendingRequests")
	defer span.Finish()

	var entitites []*entity.EnrichDetailsPreFilterTracking
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

func (r enrichDetailsPrefilterTrackingRepository) RegisterRequest(ctx context.Context, ip string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsPrefilterTrackingRepository.RegisterRequest")
	defer span.Finish()

	span.LogFields(tracingLog.String("ip", ip))

	request := entity.EnrichDetailsPreFilterTracking{
		CreatedAt: time.Now(),
		IP:        ip,
	}

	err := r.gormDb.Create(&request).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r enrichDetailsPrefilterTrackingRepository) RegisterResponse(ctx context.Context, ip string, shouldIdentify bool, response string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsPrefilterTrackingRepository.RegisterResponse")
	defer span.Finish()

	span.LogFields(tracingLog.String("ip", ip), tracingLog.Bool("shouldIdentify", shouldIdentify), tracingLog.String("response", response))

	byId, err := r.GetByIP(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	byId.ShouldIdentify = &shouldIdentify
	byId.Response = &response

	err = r.gormDb.Save(byId).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
