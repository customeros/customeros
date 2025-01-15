package listeners

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

func Handle_FlowParticipantGoalAchieved(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.Handle_FlowParticipantGoalAchieved")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message := input.(*dto.Event)
	event := message.Event.Data.(*dto.FlowParticipantGoalAchieved)

	flow, err := services.FlowService.FlowGetById(ctx, message.Event.EntityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flow == nil {
		err = errors.New("flow not found")
		tracing.TraceErr(span, err)
		return err
	}

	flowParticipant, err := services.FlowService.FlowParticipantByEntity(ctx, flow.Id, event.ParticipantId, event.ParticipantType)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	contact, err := services.ContactService.GetContactById(ctx, flowParticipant.EntityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if contact == nil {
		err = errors.New("contact not found")
		tracing.TraceErr(span, err)
		return err
	}

	contactEmail, err := services.EmailService.GetPrimaryEmailForEntityId(ctx, model.CONTACT, flowParticipant.EntityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowActionExecutions, err := services.FlowExecutionService.GetFlowActionExecutionsForParticipant(ctx, nil, flow.Id, flowParticipant.EntityId, flowParticipant.EntityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, flowActionExecution := range flowActionExecutions {
		if flowActionExecution.Status == entity.FlowActionExecutionStatusScheduled {
			err = services.Neo4jRepositories.FlowActionExecutionWriteRepository.Delete(ctx, nil, flowActionExecution.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	executionSettings, err := services.FlowExecutionService.GetFlowExecutionSettingsForEntity(ctx, nil, flow.Id, flowParticipant.EntityId, flowParticipant.EntityType)
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
		user, err := services.UserService.GetById(ctx, *executionSettings.UserId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if user == nil {
			err = errors.New("user not found")
			tracing.TraceErr(span, err)
			return err
		}

		contactName := ""
		if utils.StringFirstNonEmpty(contact.Name) != "" {
			contactName = contact.Name
		} else if utils.StringFirstNonEmpty(contact.FirstName, contact.LastName) != "" {
			contactName = utils.JoinNonEmpty(" ", contact.FirstName, contact.LastName)
		} else if contactEmail != nil && utils.StringFirstNonEmpty(contactEmail.RawEmail) != "" {
			contactName = contactEmail.RawEmail
		}

		organizationName := ""
		organizationPublicLink := ""
		contactWithOrganizations, err := services.OrganizationService.GetPrimaryOrganizationsWithJobRoleForContacts(ctx, []string{flowParticipant.EntityId})
		if err != nil {
			return errors.Wrap(err, "failed to get primary organizations with job roles for contacts")
		}

		if len(*contactWithOrganizations) > 0 {
			contactWithOrganization := (*contactWithOrganizations)[0]
			organizationName = contactWithOrganization.Organization.Name
			organizationPublicLink = fmt.Sprintf("%s/organization/%s", services.GlobalConfig.ExternalServices.NovuConfig.FronteraUrl, contactWithOrganization.Organization.ID)
		}

		//slack notification
		slackChannel, err := services.PostgresRepositories.SlackChannelNotificationRepository.GetSlackChannel(ctx, message.Event.Tenant, postgresEntity.SlackChannelNotificationWorkflowMailstackReply)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if slackChannel != nil && slackChannel.ChannelId != "" {
			slackMessageText := contactName + " has replied to the email. " + organizationName + " has achieved it’s goal!"

			err = services.SlackService.Notify(ctx, message.Event.Tenant, slackChannel.ChannelId, &slackMessageText)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}

		primaryEmail, err := services.EmailService.GetPrimaryEmailForEntityId(ctx, model.USER, user.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if primaryEmail != nil {
			notification := &service.NovuNotification{
				WorkflowId: service.WorkflowId_FlowParticipantGoalAchievedEmail,
				TemplateData: map[string]string{
					"{{orgLink}}": organizationPublicLink,
					"{{orgName}}": organizationName,
				},
				To: &service.NotifiableUser{
					FirstName:    user.FirstName,
					LastName:     user.LastName,
					Email:        primaryEmail.RawEmail,
					SubscriberID: user.Id,
				},
				Subject: fmt.Sprintf("%s has achieved it’s goal!", flow.Name),
			}

			err = services.NovuService.SendNotification(ctx, notification)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	return nil
}
