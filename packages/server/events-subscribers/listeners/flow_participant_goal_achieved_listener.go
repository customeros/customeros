package listeners

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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

	executionSettings, err := services.FlowExecutionService.GetFlowExecutionSettingsForEntity(ctx, flow.Id, flowParticipant.EntityId, flowParticipant.EntityType)
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

		primaryEmail, err := services.EmailService.GetPrimaryEmailForEntityId(ctx, model.USER, user.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if primaryEmail == nil {
			span.LogFields(log.String("msg", "primary email not found"))
			return nil
		}

		organizationName := ""
		organizationPublicLink := ""
		contactWithOrganizations, err := services.OrganizationService.GetLatestOrganizationsWithJobRolesForContacts(ctx, []string{flowParticipant.EntityId})
		if err != nil {
			return errors.Wrap(err, "failed to get latest organizations with job roles for contacts")
		}

		if len(*contactWithOrganizations) > 0 {
			contactWithOrganization := (*contactWithOrganizations)[0]
			organizationName = contactWithOrganization.Organization.Name
			organizationPublicLink = fmt.Sprintf("%s/organization/%s", services.GlobalConfig.NovuConfig.FronteraUrl, contactWithOrganization.Organization.ID)
		}

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

	return nil
}
