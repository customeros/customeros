package notifications

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/aws_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/notifications"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type OrganizationEventHandler struct {
	services             *service.Services
	log                  logger.Logger
	notificationProvider notifications.NotificationProvider
	cfg                  *config.Config
}

func NewOrganizationEventHandler(log logger.Logger, services *service.Services, cfg *config.Config) *OrganizationEventHandler {
	s3 := aws_client.NewS3Client(&aws.Config{Region: aws.String("eu-west-1")})
	return &OrganizationEventHandler{
		services:             services,
		log:                  log,
		notificationProvider: notifications.NewNovuNotificationProvider(log, cfg.Services.Novu.ApiKey, s3),
		cfg:                  cfg,
	}
}

func (h *OrganizationEventHandler) OnOrganizationUpdateOwner(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Notifications.OrganizationEventHandler.OnOrganizationUpdateOwner")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationOwnerUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	if eventData.ActorUserId == eventData.OwnerUserId {
		// do not send notification if the actor is the same as the target
		h.log.Info("actor is the same as the target, skipping notification")
		return nil
	}

	err := h.notificationProviderSendEmail(
		ctx,
		span,
		notifications.WorkflowIdOrgOwnerUpdateEmail,
		eventData.OwnerUserId,
		eventData.ActorUserId,
		eventData.OrganizationId,
		eventData.Tenant,
	)

	if err != nil {
		tracing.TraceErr(span, err)
	}

	err = h.notificationProviderSendInAppNotification(
		ctx,
		span,
		notifications.WorkflowIdOrgOwnerUpdateAppNotification,
		eventData.OwnerUserId,
		eventData.ActorUserId,
		eventData.OrganizationId,
		eventData.Tenant,
	)

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (h *OrganizationEventHandler) notificationProviderSendEmail(ctx context.Context, span opentracing.Span, workflowId, userId, actorUserId, orgId, tenant string) error {
	///////////////////////////////////       Get User, Actor, Org Content       ///////////////////////////////////
	// target user email
	emailDbNode, err := h.services.CommonServices.Neo4jRepositories.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.EmailRepository.GetEmailForUser")
	}

	var email *neo4jentity.EmailEntity
	if emailDbNode == nil {
		tracing.TraceErr(span, err)
		err = errors.New("email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendEmail")
	}
	email = neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)

	// actor user email
	actorEmailDbNode, err := h.services.CommonServices.Neo4jRepositories.EmailReadRepository.GetEmailForUser(ctx, tenant, actorUserId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.EmailRepository.GetEmailForUser")
	}

	var actorEmail *neo4jentity.EmailEntity
	if actorEmailDbNode == nil {
		tracing.TraceErr(span, err)
		err = errors.New("actor email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendEmail")
	}
	actorEmail = neo4jmapper.MapDbNodeToEmailEntity(actorEmailDbNode)

	// target user
	userDbNode, err := h.services.CommonServices.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.UserRepository.GetUser")
	}
	var user neo4jentity.UserEntity
	if userDbNode != nil {
		user = *neo4jmapper.MapDbNodeToUserEntity(userDbNode)
	}

	// actor user
	actorDbNode, err := h.services.CommonServices.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, actorUserId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.UserRepository.GetUser")
	}
	var actor neo4jentity.UserEntity
	if userDbNode != nil {
		actor = *neo4jmapper.MapDbNodeToUserEntity(actorDbNode)
	}

	// Organization
	orgDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, orgId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.OrganizationRepository.GetOrganization")
	}
	var org neo4jentity.OrganizationEntity
	if orgDbNode != nil {
		org = *neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	}

	/////////////////////////////////// Notification Provider Payload And Call ///////////////////////////////////
	orgName := org.Name
	if orgName == "" {
		orgName = "Unnamed"
	}
	subject := fmt.Sprintf(notifications.WorkflowIdOrgOwnerUpdateEmailSubject, actor.FirstName, actor.LastName)
	payload := map[string]interface{}{
		// "html":           html, fill during send notification call
		"subject":        subject,
		"email":          email.Email,
		"orgName":        orgName,
		"userFirstName":  user.FirstName,
		"actorFirstName": actor.FirstName,
		"actorLastName":  actor.LastName,
		"orgLink":        fmt.Sprintf("%s/organization/%s", h.cfg.Subscriptions.NotificationsSubscription.RedirectUrl, orgId),
	}

	overrides := map[string]interface{}{
		"email": map[string]string{
			"replyTo": actorEmail.Email,
		},
	}
	payload["overrides"] = overrides

	notification := &notifications.NovuNotification{
		WorkflowId: workflowId,
		TemplateData: map[string]string{
			"{{userFirstName}}":  user.FirstName,
			"{{actorFirstName}}": actor.FirstName,
			"{{actorLastName}}":  actor.LastName,
			"{{orgName}}":        orgName,
			"{{orgLink}}":        fmt.Sprintf("%s/organization/%s", h.cfg.Subscriptions.NotificationsSubscription.RedirectUrl, orgId),
		},
		To: &notifications.NotifiableUser{
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        email.Email,
			SubscriberID: userId,
		},
		Subject: subject,
		Payload: payload,
	}

	// call notification service
	err = h.notificationProvider.SendNotification(ctx, notification, span)

	return err
}

func (h *OrganizationEventHandler) notificationProviderSendInAppNotification(ctx context.Context, span opentracing.Span, workflowId, userId, actorUserId, orgId, tenant string) error {
	///////////////////////////////////       Get User, Actor, Org Content       ///////////////////////////////////
	// target user email
	emailDbNode, err := h.services.CommonServices.Neo4jRepositories.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.EmailRepository.GetEmailForUser")
	}

	var email *neo4jentity.EmailEntity
	if emailDbNode == nil {
		tracing.TraceErr(span, err)
		err = errors.New("email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendInAppNotification")
	}
	email = neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)

	// target user
	userDbNode, err := h.services.CommonServices.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.UserRepository.GetUser")
	}
	var user neo4jentity.UserEntity
	if userDbNode != nil {
		user = *neo4jmapper.MapDbNodeToUserEntity(userDbNode)
	}

	// actor user
	actorDbNode, err := h.services.CommonServices.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, actorUserId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.UserRepository.GetUser")
	}
	var actor neo4jentity.UserEntity
	if userDbNode != nil {
		actor = *neo4jmapper.MapDbNodeToUserEntity(actorDbNode)
	}

	// Organization
	orgDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, orgId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.OrganizationRepository.GetOrganization")
	}
	var org neo4jentity.OrganizationEntity
	if orgDbNode != nil {
		org = *neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	}
	/////////////////////////////////// Notification Provider Payload And Call ///////////////////////////////////
	subject := fmt.Sprintf(notifications.WorkflowIdOrgOwnerUpdateAppNotificationSubject, actor.FirstName, actor.LastName)
	notification := &notifications.NovuNotification{
		WorkflowId:   workflowId,
		TemplateData: map[string]string{},
		To: &notifications.NotifiableUser{
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        email.Email,
			SubscriberID: userId,
		},
		Subject: subject,
		Payload: map[string]interface{}{
			"notificationText": fmt.Sprintf("%s %s made you the owner of %s", actor.FirstName, actor.LastName, org.Name),
			"orgId":            orgId,
			"isArchived":       org.Hide,
		},
	}

	// call notification service
	err = h.notificationProvider.SendNotification(ctx, notification, span)

	return err
}
