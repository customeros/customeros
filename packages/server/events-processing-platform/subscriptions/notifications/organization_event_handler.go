package notifications

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type OrganizationEventHandler struct {
	repositories         *repository.Repositories
	log                  logger.Logger
	notificationProvider NotificationProvider
	cfg                  *config.Config
}

func NewOrganizationEventHandler(log logger.Logger, repositories *repository.Repositories, cfg *config.Config) *OrganizationEventHandler {
	return &OrganizationEventHandler{
		repositories:         repositories,
		log:                  log,
		notificationProvider: NewNotificationProvider(log, cfg.Services.Novu.ApiKey),
	}
}

type eventMetadata struct {
	UserId string `json:"user-id"`
}

func (h *OrganizationEventHandler) OnOrganizationUpdateOwner(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Notifications.OrganizationEventHandler.OnOrganizationUpdateOwner")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationUpdateOwnerEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := eventData.OwnerUserId
	actorUserId := eventData.ActorUserId

	err := h.notificationProviderSendEmail(ctx, span, EventIdUserUpdate, userId, actorUserId, eventData.Tenant)

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (h *OrganizationEventHandler) notificationProviderSendEmail(ctx context.Context, span opentracing.Span, eventId, userId, actorUserId, tenant string) error {
	// target user
	userDbNode, err := h.repositories.UserRepository.GetUser(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.UserRepository.GetUser")
	}
	var user *entity.UserEntity
	if userDbNode != nil {
		user = graph_db.MapDbNodeToUserEntity(*userDbNode)
	}

	// actor user
	actorDbNode, err := h.repositories.UserRepository.GetUser(ctx, tenant, actorUserId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.UserRepository.GetUser")
	}
	var actor *entity.UserEntity
	if userDbNode != nil {
		actor = graph_db.MapDbNodeToUserEntity(*actorDbNode)
	}

	emailDbNode, err := h.repositories.EmailRepository.GetEmailForUser(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.EmailRepository.GetEmailForUser")
	}

	var email *entity.EmailEntity
	if emailDbNode == nil {
		tracing.TraceErr(span, err)
		err = errors.New("email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendEmail")
	}
	email = graph_db.MapDbNodeToEmailEntity(*emailDbNode)

	payload := map[string]interface{}{
		"actorFirstName": actor.FirstName,
		"actorLastName":  actor.LastName,
		"subject":        fmt.Sprintf("%s %s added you as an owner", actor.FirstName, actor.LastName),
	}

	// call notification service
	err = h.notificationProvider.SendEmail(ctx, &EmailableUser{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        email.Email,
		SubscriberID: userId,
	}, payload, EventIdUserUpdate)

	return err
}
