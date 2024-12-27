package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type SlackNotificationEventsRepository interface {
	Create(ctx context.Context, notificationData entity.SlackNotificationEvents) (*entity.SlackNotificationEvents, error)
	Find(ctx context.Context, notificationData entity.SlackNotificationEvents) (*entity.SlackNotificationEvents, error)
	Update(ctx context.Context, notificationData entity.SlackNotificationEvents) (*entity.SlackNotificationEvents, error)
}

type slackNotificationEventsRepository struct {
	gormDb *gorm.DB
}

func NewSlackNotificationEventsRepository(gormDb *gorm.DB) SlackNotificationEventsRepository {
	return &slackNotificationEventsRepository{gormDb: gormDb}
}

func (r *slackNotificationEventsRepository) Create(ctx context.Context, notificationData entity.SlackNotificationEvents) (*entity.SlackNotificationEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackNotificationEventsRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if notificationData.Tenant == "" || notificationData.Domain == "" {
		err := errors.New("Tenant or Domain is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	notificationData.LastNotified = utils.NowPtr()

	var created entity.SlackNotificationEvents
	err := r.gormDb.Create(&notificationData).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &created, nil
}

func (r *slackNotificationEventsRepository) Find(ctx context.Context, notificationData entity.SlackNotificationEvents) (*entity.SlackNotificationEvents, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "SlackNotificationEventsRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var results entity.SlackNotificationEvents
	err := r.gormDb.
		Where(&notificationData).
		First(&results).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &results, nil
}

func (f *slackNotificationEventsRepository) Update(ctx context.Context, notificationData entity.SlackNotificationEvents) (*entity.SlackNotificationEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackNotificationEventsRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if notificationData.Tenant == "" || notificationData.Domain == "" {
		err := errors.New("Tenant or Domain is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if notificationData.LastNotified == nil {
		notificationData.LastNotified = utils.NowPtr()
	}

	var updatedRecord entity.SlackNotificationEvents
	err := f.gormDb.Model(&notificationData).Updates(&notificationData).First(&updatedRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedRecord, nil
}
