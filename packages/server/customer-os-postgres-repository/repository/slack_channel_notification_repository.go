package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type slackChannelNotificationRepository struct {
	db *gorm.DB
}

type SlackChannelNotificationRepository interface {
	GetSlackChannel(c context.Context, tenant, workflow string) ([]*entity.SlackChannelNotification, error)
}

func NewSlackChannelNotificationRepository(db *gorm.DB) SlackChannelNotificationRepository {
	return &slackChannelNotificationRepository{db: db}
}

func (r *slackChannelNotificationRepository) GetSlackChannel(c context.Context, tenant, workflow string) ([]*entity.SlackChannelNotification, error) {
	span, _ := opentracing.StartSpanFromContext(c, "SlackChannelNotificationRepository.MarkAsNotified")
	defer span.Finish()
	span.LogFields(tracingLog.String("tenant", tenant), tracingLog.String("workflow", workflow))

	var entities []*entity.SlackChannelNotification
	err := r.db.
		Where("tenant = ?", tenant).
		Where("workflow = ?", workflow).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}
