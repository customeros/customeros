package notifications

import (
	"context"
	"fmt"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/reminder/event"
	"time"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type ReminderEventHandler struct {
	services *service.Services
	log      logger.Logger
	cfg      *config.Config
}

func NewReminderEventHandler(log logger.Logger, services *service.Services, cfg *config.Config) *ReminderEventHandler {
	return &ReminderEventHandler{
		services: services,
		log:      log,
		cfg:      cfg,
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
		commonService.WorkflowReminderNotificationEmail,
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
		commonService.WorkflowReminderInAppNotification,
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
	emailDbNode, err := h.services.CommonServices.Neo4jRepositories.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.EmailRepository.GetEmailForUser")
	}

	var email neo4jentity.EmailEntity
	if emailDbNode == nil {
		tracing.TraceErr(span, err)
		err = errors.New("email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendEmail")
	}
	email = *neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)
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
	// Organization
	orgDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.OrganizationRepository.GetOrganization")
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
	subject := fmt.Sprintf(commonService.WorkflowReminderNotificationSubject, orgName)
	payload := map[string]interface{}{
		"subject": subject,
		"email":   email.Email,
		"orgName": orgName,
		"orgLink": fmt.Sprintf("%s/organization/%s", h.cfg.Subscriptions.NotificationsSubscription.RedirectUrl, organizationId),
	}

	notification := &commonService.NovuNotification{
		WorkflowId: workflowId,
		TemplateData: map[string]string{
			"{{reminderContent}}":   content,
			"{{reminderCreatedAt}}": createdAt.Format("Monday 02 Jan 2006"),
			"{{orgName}}":           orgName,
			"{{orgLink}}":           fmt.Sprintf("%s/organization/%s", h.cfg.Subscriptions.NotificationsSubscription.RedirectUrl, organizationId),
		},
		To: &commonService.NotifiableUser{
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        email.Email,
			SubscriberID: userId,
		},
		Subject: subject,
		Payload: payload,
	}

	// call notification service
	err = h.services.CommonServices.NovuService.SendNotification(ctx, notification)

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
	emailDbNode, err := h.services.CommonServices.Neo4jRepositories.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.EmailRepository.GetEmailForUser")
	}

	var email neo4jentity.EmailEntity
	if emailDbNode == nil {
		tracing.TraceErr(span, err)
		err = errors.New("email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendEmail")
	}
	email = *neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)
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
	// Organization
	orgDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.OrganizationRepository.GetOrganization")
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
	subject := fmt.Sprintf(commonService.WorkflowReminderNotificationSubject, orgName)
	payload := map[string]interface{}{
		"notificationText": fmt.Sprintf("%s: %s", subject, content),
		"orgId":            organizationId,
	}

	notification := &commonService.NovuNotification{
		WorkflowId:   workflowId,
		TemplateData: map[string]string{},
		To: &commonService.NotifiableUser{
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        email.Email,
			SubscriberID: userId,
		},
		Subject: subject,
		Payload: payload,
	}

	// call notification service
	err = h.services.CommonServices.NovuService.SendNotification(ctx, notification)

	return err
}
