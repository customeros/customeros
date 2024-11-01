package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
	"time"
)

type TrackingRepository interface {
	GetById(ctx context.Context, id string) (*entity.Tracking, error)
	GetNewRecords(ctx context.Context) ([]*entity.Tracking, error)
	GetForPrefilter(ctx context.Context) ([]*entity.Tracking, error)
	GetReadyForIdentification(ctx context.Context) ([]*entity.Tracking, error)
	GetIdentifiedWithDistinctIP(ctx context.Context) ([]*entity.Tracking, error)
	GetForSlackNotifications(ctx context.Context, limit int) ([]*entity.Tracking, error)

	Store(ctx context.Context, tracking entity.Tracking) (string, error)

	SetStateById(ctx context.Context, id string, newState entity.TrackingIdentificationState) error

	MarkAsNotified(ctx context.Context, id string) error
	MarkAsOrganizationCreated(ctx context.Context, id, organizationId string, organizationName, organizationDomain, organizationWebsite *string) error
	MarkAllWithState(ctx context.Context, ip string, state entity.TrackingIdentificationState) error
	MarkAllExcludeIdWithState(ctx context.Context, excludeId, ip string, state entity.TrackingIdentificationState) error
	IncrementNotificationTry(ctx context.Context, id string) error
	WasNotifiedRecently(ctx context.Context, organizationDomain string, lookBackWindowHours int) (bool, error)
}

type trackingRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewTrackingRepository(gormDb *gorm.DB) TrackingRepository {
	return &trackingRepositoryImpl{gormDb: gormDb}
}

func (r *trackingRepositoryImpl) GetById(ctx context.Context, id string) (*entity.Tracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	span.LogFields(tracingLog.String("id", id))

	var result entity.Tracking
	err := r.gormDb.
		Where("id = ?", id).
		First(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			span.LogFields(tracingLog.Bool("result.found", false))
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Bool("result.found", true))

	return &result, nil
}

func (r *trackingRepositoryImpl) GetNewRecords(ctx context.Context) ([]*entity.Tracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.GetNewRecords")
	defer span.Finish()

	var entities []*entity.Tracking
	err := r.gormDb.
		Where("event_type = 'page_exit'").
		Where("state = ?", entity.TrackingIdentificationStateNew).
		Order("created_at asc").
		Limit(500).
		Find(&entities).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(entities)))

	return entities, nil
}

func (r *trackingRepositoryImpl) GetForPrefilter(ctx context.Context) ([]*entity.Tracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.GetForPrefilter")
	defer span.Finish()

	var entities []*entity.Tracking
	err := r.gormDb.Raw("select t.* from tracking t left join enrich_details_prefilter_tracking e on t.ip = e.ip where t.state = 'PREFILTER_ASKED' and e.response is not null limit 500").Scan(&entities).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(entities)))

	return entities, nil
}

func (r *trackingRepositoryImpl) GetReadyForIdentification(ctx context.Context) ([]*entity.Tracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.GetReadyForIdentification")
	defer span.Finish()

	var entities []*entity.Tracking
	err := r.gormDb.
		Where("event_type = 'page_exit'").
		Where("state = ?", entity.TrackingIdentificationStatePrefilteredPass).
		Distinct("ip", "id", "created_at").
		Order("created_at asc").
		Limit(250).
		Find(&entities).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(entities)))

	return entities, nil
}

func (r *trackingRepositoryImpl) GetIdentifiedWithDistinctIP(ctx context.Context) ([]*entity.Tracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.GetIdentifiedWithDistinctIP")
	defer span.Finish()

	var entities []*entity.Tracking
	err := r.gormDb.
		Where("event_type = 'page_exit'").
		Where("state = ?", entity.TrackingIdentificationStateIdentified).
		Distinct("ip", "id", "tenant", "created_at").
		Order("created_at asc").
		Limit(100).
		Find(&entities).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(entities)))

	return entities, nil
}

func (r *trackingRepositoryImpl) GetForSlackNotifications(ctx context.Context, limit int) ([]*entity.Tracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.GetForSlackNotifications")
	defer span.Finish()

	stateIn := []entity.TrackingIdentificationState{
		entity.TrackingIdentificationStateOrganizationCreated,
		entity.TrackingIdentificationStateOrganizationExists,
	}

	var entities []*entity.Tracking
	err := r.gormDb.
		Where("notified = false").
		Where("organization_id is not null").
		Where("event_type = 'page_exit'").
		Where("state in ?", stateIn).
		Where("notification_try < 5").
		Distinct("ip", "id", "tenant", "created_at").
		Order("created_at asc").
		Limit(limit).
		Find(&entities).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(entities)))

	return entities, nil
}

func (r *trackingRepositoryImpl) MarkAllWithState(ctx context.Context, ip string, state entity.TrackingIdentificationState) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.MarkAllWithState")
	defer span.Finish()

	err := r.gormDb.
		Model(&entity.Tracking{}).
		Where("ip = ?", ip).
		Update("state", state).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *trackingRepositoryImpl) MarkAllExcludeIdWithState(ctx context.Context, excludeId, ip string, state entity.TrackingIdentificationState) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.MarkAllExcludeIdWithState")
	defer span.Finish()

	err := r.gormDb.
		Model(&entity.Tracking{}).
		Where("id != ?", excludeId).
		Where("ip = ?", ip).
		Update("state", state).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *trackingRepositoryImpl) Store(ctx context.Context, tracking entity.Tracking) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.Store")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	if tracking.Tenant != "" {
		tracing.TagTenant(span, tracking.Tenant)
	}
	tracing.LogObjectAsJson(span, "tracking", tracking)

	if tracking.ID == "" {
		tracking.CreatedAt = utils.Now()
	}

	err := r.gormDb.Save(&tracking).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	span.LogFields(tracingLog.String("tracking.id", tracking.ID))

	return tracking.ID, nil
}

func (r *trackingRepositoryImpl) SetStateById(c context.Context, id string, newState entity.TrackingIdentificationState) error {
	span, _ := opentracing.StartSpanFromContext(c, "TrackingRepository.SetStateById")
	defer span.Finish()
	span.LogFields(tracingLog.String("ip", id), tracingLog.String("newState", string(newState)))

	err := r.gormDb.
		Model(&entity.Tracking{}).
		Where("id = ?", id).
		Update("state", newState).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *trackingRepositoryImpl) MarkAsNotified(c context.Context, id string) error {
	span, _ := opentracing.StartSpanFromContext(c, "TrackingRepository.MarkAsNotified")
	defer span.Finish()
	span.LogFields(tracingLog.String("ip", id))

	err := r.gormDb.
		Model(&entity.Tracking{}).
		Where("id = ?", id).
		Update("notified", true).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *trackingRepositoryImpl) MarkAsOrganizationCreated(c context.Context, id, organizationId string, organizationName, organizationDomain, organizationWebsite *string) error {
	span, _ := opentracing.StartSpanFromContext(c, "TrackingRepository.MarkAsOrganizationCreated")
	defer span.Finish()
	span.LogFields(tracingLog.String("ip", id), tracingLog.String("organizationId", organizationId))

	err := r.gormDb.
		Model(&entity.Tracking{}).
		Where("id = ?", id).
		Update("state", entity.TrackingIdentificationStateOrganizationCreated).
		Update("organization_id", organizationId).
		Update("organization_name", organizationName).
		Update("organization_domain", organizationDomain).
		Update("organization_website", organizationWebsite).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *trackingRepositoryImpl) IncrementNotificationTry(ctx context.Context, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.IncrementNotificationTry")
	defer span.Finish()
	span.LogFields(tracingLog.String("id", id))

	err := r.gormDb.
		Model(&entity.Tracking{}).
		Where("id = ?", id).
		Update("notification_try", gorm.Expr("notification_try + 1")).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *trackingRepositoryImpl) WasNotifiedRecently(ctx context.Context, organizationDomain string, lookBackWindowHours int) (bool, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingRepository.WasNotifiedRecently")
	defer span.Finish()
	span.LogFields(tracingLog.String("organizationDomain", organizationDomain), tracingLog.Int("lookBackWindowHours", lookBackWindowHours))

	var count int64
	err := r.gormDb.
		Model(&entity.Tracking{}).
		Where("organization_domain = ?", organizationDomain).
		Where("notified = true").
		Where("created_at > ?", utils.Now().Add(-time.Duration(lookBackWindowHours)*time.Hour)).
		Count(&count).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	span.LogFields(tracingLog.Int("result.count", int(count)))

	return count > 0, nil
}
