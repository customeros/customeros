package notifications

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type UserEventHandler struct {
	log                  logger.Logger
	repositories         *repository.Repositories
	notificationProvider NotificationProvider
	cfg                  *config.Config
}

func NewUserEventHandler(log logger.Logger, repositories *repository.Repositories, cfg *config.Config) *UserEventHandler {
	return &UserEventHandler{
		log:                  log,
		repositories:         repositories,
		cfg:                  cfg,
		notificationProvider: NewNotificationProvider(log, cfg.Services.Novu.ApiKey),
	}
}

func (h *UserEventHandler) OnUserUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnUserUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)

	err := h.notificationProviderSendEmail(ctx, span, EventIdUserUpdate, userId, eventData.Tenant)

	return err
}

func (h *UserEventHandler) OnJobRoleLinkedToUser(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnJobRoleLinkedToUser")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserLinkJobRoleEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)

	err := h.notificationProviderSendEmail(ctx, span, EventIdUserUpdate, userId, eventData.Tenant)

	return err
}

func (h *UserEventHandler) OnAddRole(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnAddRole")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserAddRoleEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)

	err := h.notificationProviderSendEmail(ctx, span, EventIdUserUpdate, userId, eventData.Tenant)

	return err
}

func (h *UserEventHandler) notificationProviderSendEmail(ctx context.Context, span opentracing.Span, eventId, userId, tenant string) error {

	userDbNode, err := h.repositories.UserRepository.GetUser(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.UserRepository.GetUser")
	}
	var user *entity.UserEntity
	if userDbNode != nil {
		user = graph_db.MapDbNodeToUserEntity(*userDbNode)
	}

	emailDbNode, err := h.repositories.EmailRepository.GetEmailForUser(ctx, tenant, userId)

	var email *entity.EmailEntity
	if emailDbNode != nil {
		email = graph_db.MapDbNodeToEmailEntity(*emailDbNode)
	}

	// call notification service
	err = h.notificationProvider.SendEmail(ctx, &EmailableUser{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        email.Email,
		Message:      "Welcome to CustomerOS!",
		SubscriberID: userId,
	}, EventIdUserUpdate)

	return err
}
