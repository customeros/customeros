package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type slackChannelNotificationRepository struct {
	db *gorm.DB
}

type SlackChannelNotificationRepository interface {
	GetSlackChannel(c context.Context, workflow postgres_entity.SlackChannelNotificationWorkflow) (*postgres_entity.SlackChannelNotification, error)
}

func NewSlackChannelNotificationRepository(db *gorm.DB) SlackChannelNotificationRepository {
	return &slackChannelNotificationRepository{db: db}
}

func (r *slackChannelNotificationRepository) GetSlackChannel(c context.Context, workflow postgres_entity.SlackChannelNotificationWorkflow) (*postgres_entity.SlackChannelNotification, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "SlackChannelNotificationRepository.GetSlackChannels")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogKV("workflow", workflow)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set on context")
		tracing.TraceErr(span, err)
		return nil, err
	}
	tracing.TagTenant(span, tenant)

	var e *postgres_entity.SlackChannelNotification
	err := r.db.
		Where("tenant = ?", tenant).
		Where("workflow = ?", workflow).
		First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return e, nil
}
