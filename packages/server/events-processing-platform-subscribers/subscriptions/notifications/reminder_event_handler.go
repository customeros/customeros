package notifications

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/reminder/event"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/aws_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/notifications"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type ReminderEventHandler struct {
	repositories         *repository.Repositories
	log                  logger.Logger
	notificationProvider notifications.NotificationProvider
	cfg                  *config.Config
}

func NewReminderEventHandler(log logger.Logger, repositories *repository.Repositories, cfg *config.Config) *ReminderEventHandler {
	s3 := aws_client.NewS3Client(&aws.Config{Region: aws.String("eu-west-1")})
	return &ReminderEventHandler{
		repositories:         repositories,
		log:                  log,
		notificationProvider: notifications.NewNovuNotificationProvider(log, cfg.Services.Novu.ApiKey, s3),
		cfg:                  cfg,
	}
}

func (h *ReminderEventHandler) OnReminderNotification(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Notifications.ReminderEventHandler.OnReminderNotification")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ReminderNotificationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	err := h.notificationProviderSendEmail(
		ctx,
		span,
		notifications.WorkflowReminderNotificationEmail,
		eventData.UserId,
		eventData.Content,
		eventData.OrganizationId,
		eventData.Tenant,
		eventData.CreatedAt,
	)

	if err != nil {
		tracing.TraceErr(span, err)
	}

	err = h.notificationProviderSendInAppNotification(
		ctx,
		span,
		notifications.WorkflowReminderInAppNotification,
		eventData.UserId,
		eventData.Content,
		eventData.OrganizationId,
		eventData.Tenant,
	)

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////
// ///////////////////// Send Email Notification //////////////////////////
// ////////////////////////////////////////////////////////////////////////

func (h *ReminderEventHandler) notificationProviderSendEmail(
	ctx context.Context,
	span opentracing.Span,
	workflowId string,
	userId string,
	content string,
	organizationId string,
	tenant string,
	createdAt time.Time,
) error {
	// target user email
	emailDbNode, err := h.repositories.Neo4jRepositories.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.EmailRepository.GetEmailForUser")
	}

	var email neo4jentity.EmailEntity
	if emailDbNode == nil {
		tracing.TraceErr(span, err)
		err = errors.New("email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendEmail")
	}
	email = *neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)
	// target user
	userDbNode, err := h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.UserRepository.GetUser")
	}
	var user neo4jentity.UserEntity
	if userDbNode != nil {
		user = *neo4jmapper.MapDbNodeToUserEntity(userDbNode)
	}
	// Organization
	orgDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.OrganizationRepository.GetOrganization")
	}
	var org neo4jentity.OrganizationEntity
	if orgDbNode != nil {
		org = *neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	}
	////////////////////////////////////////////////////
	// ////////// Format email and send it ////////////
	//////////////////////////////////////////////////
	orgName := org.Name
	if orgName == "" {
		orgName = "Unnamed"
	}
	subject := fmt.Sprintf(notifications.WorkflowReminderNotificationSubject, orgName)
	payload := map[string]interface{}{
		"subject": subject,
		"email":   email.Email,
		"orgName": orgName,
		"orgLink": fmt.Sprintf("%s/organization/%s", h.cfg.Subscriptions.NotificationsSubscription.RedirectUrl, organizationId),
	}

	notification := &notifications.NovuNotification{
		WorkflowId: workflowId,
		TemplateData: map[string]string{
			"{{reminderContent}}":   content,
			"{{reminderCreatedAt}}": createdAt.Format("Monday 02 Jan 2006"),
			"{{orgName}}":           orgName,
			"{{orgLink}}":           fmt.Sprintf("%s/organization/%s", h.cfg.Subscriptions.NotificationsSubscription.RedirectUrl, organizationId),
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

// ////////////////////////////////////////////////////////////////////////
// //////////////////// Send In App Notification //////////////////////////
// ////////////////////////////////////////////////////////////////////////
func (h *ReminderEventHandler) notificationProviderSendInAppNotification(
	ctx context.Context,
	span opentracing.Span,
	workflowId string,
	userId string,
	content string,
	organizationId string,
	tenant string,
) error {
	// target user email
	emailDbNode, err := h.repositories.Neo4jRepositories.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.EmailRepository.GetEmailForUser")
	}

	var email neo4jentity.EmailEntity
	if emailDbNode == nil {
		tracing.TraceErr(span, err)
		err = errors.New("email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendEmail")
	}
	email = *neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)
	// target user
	userDbNode, err := h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.UserRepository.GetUser")
	}
	var user neo4jentity.UserEntity
	if userDbNode != nil {
		user = *neo4jmapper.MapDbNodeToUserEntity(userDbNode)
	}
	// Organization
	orgDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.OrganizationRepository.GetOrganization")
	}
	var org neo4jentity.OrganizationEntity
	if orgDbNode != nil {
		org = *neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	}
	////////////////////////////////////////////////////
	// //////// Format Notification and send it ///////
	//////////////////////////////////////////////////
	orgName := org.Name
	if orgName == "" {
		orgName = "Unnamed"
	}
	subject := fmt.Sprintf(notifications.WorkflowReminderNotificationSubject, orgName)
	payload := map[string]interface{}{
		"notificationText": fmt.Sprintf("%s: %s", subject, content),
		"orgId":            organizationId,
	}

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
		Payload: payload,
	}

	// call notification service
	err = h.notificationProvider.SendNotification(ctx, notification, span)

	return err
}
