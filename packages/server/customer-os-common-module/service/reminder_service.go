package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ReminderService interface {
	CreateReminder(ctx context.Context, tenant, userId, orgId, content string, dueDate time.Time) (string, error)
	UpdateReminder(ctx context.Context, tenant, id string, content *string, dueDate *time.Time, dismissed *bool) error
	GetReminderById(ctx context.Context, id string) (*neo4jentity.ReminderEntity, error)
	RemindersForOrganization(ctx context.Context, organizationID string, dismissed *bool) ([]*neo4jentity.ReminderEntity, error)

	SendNotification(ctx context.Context, reminderId, fronteraPublicPath string) error
}

type reminderService struct {
	services *Services
}

func NewReminderService(services *Services) ReminderService {
	return &reminderService{
		services: services,
	}
}

func (s *reminderService) CreateReminder(ctx context.Context, tenant, userId, organizationId, content string, dueDate time.Time) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderService.CreateReminder")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("userId", userId), log.String("organizationId", organizationId), log.String("content", content), log.String("dueDate", dueDate.String()))

	tenant = common.GetTenantFromContext(ctx)

	reminderId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelReminder)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	err = s.services.Neo4jRepositories.ReminderWriteRepository.CreateReminder(ctx, tenant, reminderId, userId, organizationId, content, utils.Now(), dueDate)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return reminderId, nil
}

func (s *reminderService) UpdateReminder(ctx context.Context, tenant, reminderId string, content *string, dueDate *time.Time, dismissed *bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderService.UpdateReminder")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, reminderId)

	tenant = common.GetTenantFromContext(ctx)

	if content == nil && dueDate == nil && dismissed == nil {
		return nil
	}

	updateData := neo4jrepo.ReminderUpdateFields{}

	if content != nil {
		updateData.Content = content
		updateData.UpdateContent = true
	}
	if dueDate != nil {
		updateData.DueDate = dueDate
		updateData.UpdateDueDate = true
	}
	if dismissed != nil {
		updateData.Dismissed = dismissed
		updateData.UpdateDismissed = true
	}

	err := s.services.Neo4jRepositories.ReminderWriteRepository.UpdateReminder(ctx, tenant, reminderId, updateData)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *reminderService) GetReminderById(ctx context.Context, id string) (*neo4jentity.ReminderEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderService.GetReminderById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, id)

	if reminderDbNode, err := s.services.Neo4jRepositories.ReminderReadRepository.GetReminderById(ctx, id); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Reminder with id {%s} not found", id))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToReminderEntity(reminderDbNode), nil
	}
}

func (s *reminderService) RemindersForOrganization(ctx context.Context, organizationID string, dismissed *bool) ([]*neo4jentity.ReminderEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderService.RemindersForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, organizationID)

	reminderDbNodes, err := s.services.Neo4jRepositories.ReminderReadRepository.GetRemindersOrderByDueDateAsc(ctx, organizationID, dismissed)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	reminderEntities := make([]*neo4jentity.ReminderEntity, 0, len(reminderDbNodes))
	if len(reminderDbNodes) == 0 {
		span.LogFields(log.String("Warning", fmt.Sprintf("Reminders for organization with id {%s} not found", organizationID)))
		return reminderEntities, nil
	}
	for _, v := range reminderDbNodes {
		reminderEntities = append(reminderEntities, neo4jmapper.MapDbNodeToReminderEntity(v))
	}
	return reminderEntities, nil
}

func (s *reminderService) SendNotification(ctx context.Context, reminderId, fronteraPublicPath string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderService.SendNotification")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, reminderId)

	tenant := common.GetTenantFromContext(ctx)

	reminder, err := s.GetReminderById(ctx, reminderId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.notificationProviderSendEmail(ctx, span, fronteraPublicPath, WorkflowReminderNotificationEmail, reminder.UserId, reminder.Content, reminder.OrganizationId, tenant, reminder.CreatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.notificationProviderSendInAppNotification(ctx, span, WorkflowReminderInAppNotification, reminder.UserId, reminder.Content, reminder.OrganizationId, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////
// ///////////////////// Send Email Notification //////////////////////////
// ////////////////////////////////////////////////////////////////////////

func (s *reminderService) notificationProviderSendEmail(
	ctx context.Context,
	span opentracing.Span,
	fronteraPublicPath string,
	workflowId string,
	userId string,
	content string,
	organizationId string,
	tenant string,
	createdAt time.Time,
) error {
	// target user email
	emailDbNode, err := s.services.Neo4jRepositories.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)

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
	userDbNode, err := s.services.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.UserRepository.GetUser")
	}
	var user neo4jentity.UserEntity
	if userDbNode != nil {
		user = *neo4jmapper.MapDbNodeToUserEntity(userDbNode)
	}
	// Organization
	orgDbNode, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)

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
	subject := fmt.Sprintf(WorkflowReminderNotificationSubject, orgName)
	payload := map[string]interface{}{
		"subject": subject,
		"email":   email.Email,
		"orgName": orgName,
		"orgLink": fmt.Sprintf("%s/organization/%s", fronteraPublicPath, organizationId),
	}

	notification := &NovuNotification{
		WorkflowId: workflowId,
		TemplateData: map[string]string{
			"{{reminderContent}}":   content,
			"{{reminderCreatedAt}}": createdAt.Format("Monday 02 Jan 2006"),
			"{{orgName}}":           orgName,
			"{{orgLink}}":           fmt.Sprintf("%s/organization/%s", fronteraPublicPath, organizationId),
		},
		To: &NotifiableUser{
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        email.Email,
			SubscriberID: userId,
		},
		Subject: subject,
		Payload: payload,
	}

	// call notification service
	err = s.services.NovuService.SendNotification(ctx, notification)

	return err
}

// ////////////////////////////////////////////////////////////////////////
// //////////////////// Send In App Notification //////////////////////////
// ////////////////////////////////////////////////////////////////////////
func (h *reminderService) notificationProviderSendInAppNotification(
	ctx context.Context,
	span opentracing.Span,
	workflowId string,
	userId string,
	content string,
	organizationId string,
	tenant string,
) error {
	// target user email
	emailDbNode, err := h.services.Neo4jRepositories.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)

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
	userDbNode, err := h.services.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.services.CommonServices.UserRepository.GetUser")
	}
	var user neo4jentity.UserEntity
	if userDbNode != nil {
		user = *neo4jmapper.MapDbNodeToUserEntity(userDbNode)
	}
	// Organization
	orgDbNode, err := h.services.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)

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
	subject := fmt.Sprintf(WorkflowReminderNotificationSubject, orgName)
	payload := map[string]interface{}{
		"notificationText": fmt.Sprintf("%s: %s", subject, content),
		"orgId":            organizationId,
	}

	notification := &NovuNotification{
		WorkflowId:   workflowId,
		TemplateData: map[string]string{},
		To: &NotifiableUser{
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        email.Email,
			SubscriberID: userId,
		},
		Subject: subject,
		Payload: payload,
	}

	// call notification service
	err = h.services.NovuService.SendNotification(ctx, notification)

	return err
}
