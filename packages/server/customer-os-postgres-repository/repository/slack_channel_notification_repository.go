package repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type slackChannelNotificationRepository struct {
	db *gorm.DB
}

type SlackChannelNotificationRepository interface {
	GetSlackChannel(c context.Context, tenant string, workflow entity.SlackChannelNotificationWorkflow) (*entity.SlackChannelNotification, error)
}

func NewSlackChannelNotificationRepository(db *gorm.DB) SlackChannelNotificationRepository {
	return &slackChannelNotificationRepository{db: db}
}

func (r *slackChannelNotificationRepository) GetSlackChannel(c context.Context, tenant string, workflow entity.SlackChannelNotificationWorkflow) (*entity.SlackChannelNotification, error) {
	span, _ := opentracing.StartSpanFromContext(c, "SlackChannelNotificationRepository.GetSlackChannels")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(tracingLog.Object("workflow", workflow))

	var e *entity.SlackChannelNotification
	err := r.db.
		Where("tenant = ?", tenant).
		Where("workflow = ?", workflow).
		First(&e).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return e, nil
}
