package events_listeners

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	common_model "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/novu"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type FlowParticipantGoalAchievedListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewFlowParticipantGoalAchievedListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &FlowParticipantGoalAchievedListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.FlowParticipantGoalAchieved](), // subscribed event
			events.QueueEvents, // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *FlowParticipantGoalAchievedListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowParticipantGoalAchievedListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	data, err := events.DecodeEventData[dto.FlowParticipantGoalAchieved](ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if data.ParticipantId == "" {
		err := errors.New("ParticipantID is empty")
		tracing.TraceErr(span, err)
		return err
	}

	if data.ParticipantType == "" {
		err := errors.New("ParticipantType is empty")
		tracing.TraceErr(span, err)
		return err
	}

	return l.handle(ctx, event.Event.EntityId, data)
}

func (l *FlowParticipantGoalAchievedListener) handle(ctx context.Context, entityId string, message dto.FlowParticipantGoalAchieved) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowParticipantGoalAchievedListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	flow, err := l.dependencies.CommonServices.FlowService.FlowGetById(ctx, entityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flow == nil {
		err = errors.New("flow not found")
		tracing.TraceErr(span, err)
		return err
	}

	flowParticipant, err := l.dependencies.CommonServices.FlowService.FlowParticipantByEntity(ctx, flow.Id, message.ParticipantId, message.ParticipantType)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	contact, err := l.dependencies.CommonServices.ContactService.GetContactById(ctx, flowParticipant.EntityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if contact == nil {
		err = errors.New("contact not found")
		tracing.TraceErr(span, err)
		return err
	}

	contactEmail, err := l.dependencies.CommonServices.EmailService.GetPrimaryEmailForEntityId(ctx, common_model.CONTACT, flowParticipant.EntityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowActionExecutions, err := l.dependencies.CommonServices.FlowExecutionService.GetFlowActionExecutionsForParticipant(ctx, nil, flow.Id, flowParticipant.EntityId, flowParticipant.EntityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, flowActionExecution := range flowActionExecutions {
		if flowActionExecution.Status == neo4j_entity.FlowActionExecutionStatusScheduled {
			err = l.dependencies.Neo4jRepositories.FlowActionExecutionWriteRepository.Delete(ctx, nil, flowActionExecution.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	executionSettings, err := l.dependencies.CommonServices.FlowExecutionService.GetFlowExecutionSettingsForEntity(ctx, nil, flow.Id, flowParticipant.EntityId, flowParticipant.EntityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if executionSettings == nil {
		err = errors.New("execution settings not found")
		tracing.TraceErr(span, err)
		return err
	}

	if executionSettings.UserId != nil {
		user, err := l.dependencies.CommonServices.UserService.GetById(ctx, *executionSettings.UserId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if user == nil {
			err = errors.New("user not found")
			tracing.TraceErr(span, err)
			return err
		}

		contactName := contact.FullName()
		if contactName == "" && contactEmail != nil && utils.StringFirstNonEmpty(contactEmail.RawEmail) != "" {
			contactName = contactEmail.RawEmail
		}

		organizationName := ""
		organizationPublicLink := ""
		contactWithOrganizations, err := l.dependencies.CommonServices.OrganizationService.GetPrimaryOrganizationsWithJobRoleForContacts(ctx, []string{flowParticipant.EntityId})
		if err != nil {
			return errors.Wrap(err, "failed to get primary organizations with job roles for contacts")
		}

		if len(*contactWithOrganizations) > 0 {
			contactWithOrganization := (*contactWithOrganizations)[0]
			organizationName = contactWithOrganization.Organization.Name
			organizationPublicLink = fmt.Sprintf("%s/organization/%s", l.dependencies.CommonConfig.External.NovuConfig.FronteraUrl, contactWithOrganization.Organization.ID)
		}

		// slack notification
		slackChannel, err := l.dependencies.PostgresRepositories.SlackChannelNotificationRepository.GetSlackChannel(ctx, postgres_entity.SlackChannelNotificationWorkflowMailstackReply)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if slackChannel != nil && slackChannel.ChannelId != "" {
			slackMessageText := contactName + " has replied to the email. " + organizationName + " has achieved it’s goal!"

			err = l.dependencies.CommonServices.NotificationService.NotifySlackChannel(ctx, slackChannel.ChannelId, slackMessageText)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}

		primaryEmail, err := l.dependencies.CommonServices.EmailService.GetPrimaryEmailForEntityId(ctx, common_model.USER, user.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if primaryEmail != nil {
			notification := &interfaces.NovuNotification{
				WorkflowId: novu.WorkflowId_FlowParticipantGoalAchievedEmail,
				TemplateData: map[string]string{
					"{{orgLink}}": organizationPublicLink,
					"{{orgName}}": organizationName,
				},
				To: &interfaces.NotifiableUser{
					FirstName:    user.FirstName,
					LastName:     user.LastName,
					Email:        primaryEmail.RawEmail,
					SubscriberID: user.Id,
				},
				Subject: fmt.Sprintf("%s has achieved it’s goal!", flow.Name),
			}

			err = l.dependencies.CommonServices.NovuService.SendNotification(ctx, notification)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	return nil
}
