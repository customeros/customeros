package reminders

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"time"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/novu"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type reminderService struct {
	neo4j *neo4j_repository.Repositories
	novu  interfaces.NovuService
}

func NewReminderService(neo4j *neo4j_repository.Repositories, novu interfaces.NovuService) interfaces.ReminderService {
	return &reminderService{
		neo4j: neo4j,
		novu:  novu,
	}
}

func (s *reminderService) CreateReminder(ctx context.Context, tenant, userId, organizationId, content string, dueDate time.Time) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ReminderService.CreateReminder")
	defer spans.Finish()

	tenant = common.GetTenantFromContext(ctx)

	reminderId, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelReminder)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	err = s.neo4j.ReminderWriteRepository.CreateReminder(ctx, tenant, reminderId, userId, organizationId, content, utils.Now(), dueDate)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	return reminderId, nil
}

func (s *reminderService) UpdateReminder(ctx context.Context, tenant, reminderId string, content *string, dueDate *time.Time, dismissed, sent *bool) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ReminderService.UpdateReminder")
	defer spans.Finish()

	spans.TagEntity(reminderId)

	tenant = common.GetTenantFromContext(ctx)

	if content == nil && dueDate == nil && dismissed == nil {
		return nil
	}

	updateData := neo4j_repository.ReminderUpdateFields{}

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
	if sent != nil {
		updateData.Sent = sent
		updateData.UpdateSent = true
	}

	err := s.neo4j.ReminderWriteRepository.UpdateReminder(ctx, tenant, reminderId, updateData)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (s *reminderService) GetReminderById(ctx context.Context, id string) (*neo4jentity.ReminderEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ReminderService.GetReminderById")
	defer spans.Finish()

	spans.TagEntity(id)

	if reminderDbNode, err := s.neo4j.ReminderReadRepository.GetReminderById(ctx, id); err != nil {
		spans.TraceError(err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Reminder with id {%s} not found", id))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToReminderEntity(reminderDbNode), nil
	}
}

func (s *reminderService) RemindersForOrganization(ctx context.Context, organizationID string, dismissed *bool) ([]*neo4jentity.ReminderEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ReminderService.RemindersForOrganization")
	defer spans.Finish()

	spans.TagEntity(organizationID)

	reminderDbNodes, err := s.neo4j.ReminderReadRepository.GetRemindersOrderByDueDateAsc(ctx, organizationID, dismissed)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	reminderEntities := make([]*neo4jentity.ReminderEntity, 0, len(reminderDbNodes))
	if len(reminderDbNodes) == 0 {
		spans.LogKV("Warning", fmt.Sprintf("Reminders for organization with id {%s} not found", organizationID))
		return reminderEntities, nil
	}
	for _, v := range reminderDbNodes {
		reminderEntities = append(reminderEntities, neo4jmapper.MapDbNodeToReminderEntity(v))
	}
	return reminderEntities, nil
}

func (s *reminderService) SendNotification(ctx context.Context, reminderId, fronteraPublicPath string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ReminderService.SendNotification")
	defer spans.Finish()

	spans.TagEntity(reminderId)

	tenant := common.GetTenantFromContext(ctx)

	reminder, err := s.GetReminderById(ctx, reminderId)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	err = s.notificationProviderSendEmail(ctx, *spans, fronteraPublicPath, novu.WorkflowReminderNotificationEmail, reminder.UserId, reminder.Content, reminder.OrganizationId, tenant, reminder.CreatedAt)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	// TODO NOVU FAILS TO RECEIVE THIS NOTIFICATION

	//err = s.notificationProviderSendInAppNotification(ctx, span, WorkflowReminderInAppNotification, reminder.UserId, reminder.Content, reminder.OrganizationId, tenant)
	//if err != nil {
	//	spans.TraceError( err)
	//	return err
	//}

	return nil
}

// ////////////////////////////////////////////////////////////////////////
// ///////////////////// Send Email Notification //////////////////////////
// ////////////////////////////////////////////////////////////////////////

func (s *reminderService) notificationProviderSendEmail(
	ctx context.Context,
	spans telemetry.Spans,
	fronteraPublicPath string,
	workflowId string,
	userId string,
	content string,
	organizationId string,
	tenant string,
	createdAt time.Time,
) error {
	// target user email
	emailDbNode, err := s.neo4j.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)
	if err != nil {
		spans.TraceError(err)
		return errors.Wrap(err, "h.services.CommonServices.EmailRepository.GetEmailForUser")
	}

	var email neo4jentity.EmailEntity
	if emailDbNode == nil {
		spans.TraceError(err)
		err = errors.New("email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendEmail")
	}
	email = *neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)
	// target user
	userDbNode, err := s.neo4j.UserReadRepository.GetUserById(ctx, tenant, userId)
	if err != nil {
		spans.TraceError(err)
		return errors.Wrap(err, "h.services.CommonServices.UserRepository.GetUser")
	}
	var user neo4jentity.UserEntity
	if userDbNode != nil {
		user = *neo4jmapper.MapDbNodeToUserEntity(userDbNode)
	}
	// Organization
	orgDbNode, err := s.neo4j.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)
	if err != nil {
		spans.TraceError(err)
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
	subject := fmt.Sprintf(novu.WorkflowReminderNotificationSubject, orgName)
	payload := map[string]interface{}{
		"subject": subject,
		"email":   email.Email,
		"orgName": orgName,
		"orgLink": fmt.Sprintf("%s/organization/%s", fronteraPublicPath, organizationId),
	}

	notification := &interfaces.NovuNotification{
		WorkflowId: workflowId,
		TemplateData: map[string]string{
			"{{reminderContent}}":   content,
			"{{reminderCreatedAt}}": createdAt.Format("Monday 02 Jan 2006"),
			"{{orgName}}":           orgName,
			"{{orgLink}}":           fmt.Sprintf("%s/organization/%s", fronteraPublicPath, organizationId),
		},
		To: &interfaces.NotifiableUser{
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        email.Email,
			SubscriberID: userId,
		},
		Subject: subject,
		Payload: payload,
	}

	// call notification service
	err = s.novu.SendNotification(ctx, notification)

	return err
}

// ////////////////////////////////////////////////////////////////////////
// //////////////////// Send In App Notification //////////////////////////
// ////////////////////////////////////////////////////////////////////////
func (h *reminderService) notificationProviderSendInAppNotification(
	ctx context.Context,
	spans telemetry.Spans,
	workflowId string,
	userId string,
	content string,
	organizationId string,
	tenant string,
) error {
	// target user email
	emailDbNode, err := h.neo4j.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)
	if err != nil {
		spans.TraceError(err)
		return errors.Wrap(err, "h.services.CommonServices.EmailRepository.GetEmailForUser")
	}

	var email neo4jentity.EmailEntity
	if emailDbNode == nil {
		spans.TraceError(err)
		err = errors.New("email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendEmail")
	}
	email = *neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)
	// target user
	userDbNode, err := h.neo4j.UserReadRepository.GetUserById(ctx, tenant, userId)
	if err != nil {
		spans.TraceError(err)
		return errors.Wrap(err, "h.services.CommonServices.UserRepository.GetUser")
	}
	var user neo4jentity.UserEntity
	if userDbNode != nil {
		user = *neo4jmapper.MapDbNodeToUserEntity(userDbNode)
	}
	// Organization
	orgDbNode, err := h.neo4j.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)
	if err != nil {
		spans.TraceError(err)
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
	subject := fmt.Sprintf(novu.WorkflowReminderNotificationSubject, orgName)
	payload := map[string]interface{}{
		"notificationText": fmt.Sprintf("%s: %s", subject, content),
		"orgId":            organizationId,
	}

	notification := &interfaces.NovuNotification{
		WorkflowId:   workflowId,
		TemplateData: map[string]string{},
		To: &interfaces.NotifiableUser{
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        email.Email,
			SubscriberID: userId,
		},
		Subject: subject,
		Payload: payload,
	}

	// call notification service
	err = h.novu.SendNotification(ctx, notification)

	return err
}
