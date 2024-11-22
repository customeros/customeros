package service

import (
	"context"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
)

type FlowService interface {
	ProcessEmailFlows(ctx context.Context, tenant string, data *model.PostmarkEmailWebhookData, participants []string) error
	ProcessMailstackReply(ctx context.Context, tenant string, data *model.PostmarkEmailWebhookData, slackChannel string) error
}

type flowService struct {
	services *Services
	logger   logger.Logger
}

func NewFlowService(services *service.Services) FlowService {
	return &flowService{
		services: services,
		logger:   services.Logger,
	}
}

func (s *flowService) ProcessEmailFlows(ctx context.Context, tenant string, data *model.PostmarkEmailWebhookData, participants []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.ProcessEmailFlows")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Check if sender is a system user
	senderUser, err := s.services.CommonServices.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, data.FromFull.Email)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Only process flows for outbound emails (from system users)
	if senderUser == nil {
		return nil
	}

	// Process each participant
	for _, participantEmail := range participants {
		if err := s.processParticipantFlows(ctx, tenant, participantEmail, data.Subject); err != nil {
			return err
		}
	}

	return nil
}

func (s *flowService) processParticipantFlows(ctx context.Context, tenant, participantEmail, subject string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.processParticipantFlows")
	defer span.Finish()

	// Get contacts with this email
	contacts, err := s.services.CommonServices.Neo4jRepositories.ContactReadRepository.GetContactsWithEmail(ctx, tenant, participantEmail)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Process each contact's flows
	for _, contactNode := range contacts {
		contactEntity := mapper.MapDbNodeToContactEntity(contactNode)
		if err := s.processContactFlows(ctx, tenant, contactEntity, subject); err != nil {
			return err
		}
	}

	return nil
}

func (s *flowService) processContactFlows(ctx context.Context, tenant string, contact *neo4jentity.ContactEntity, subject string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.processContactFlows")
	defer span.Finish()

	// Get flows that include this contact
	flows, err := s.services.CommonServices.FlowService.FlowsGetListWithParticipant(ctx, []string{contact.Id}, commonModel.CONTACT)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Process each flow
	for _, flow := range *flows {
		if err := s.updateFlowParticipantStatus(ctx, tenant, flow.Id, contact, subject); err != nil {
			return err
		}
	}

	return nil
}

func (s *flowService) updateFlowParticipantStatus(ctx context.Context, tenant, flowId string, contact *neo4jentity.ContactEntity, subject string) error {
	flowParticipant, err := s.services.CommonServices.FlowService.FlowParticipantByEntity(ctx, flowId, contact.Id, commonModel.CONTACT)
	if err != nil {
		return err
	}

	if flowParticipant == nil || s.isCompletedStatus(flowParticipant.Status) {
		return nil
	}

	newStatus := s.determineNewStatus(subject)
	flowParticipant.Status = newStatus

	_, err = s.services.CommonServices.Neo4jRepositories.FlowParticipantWriteRepository.Merge(ctx, nil, flowParticipant)
	return err
}

func (s *flowService) isCompletedStatus(status neo4jentity.FlowParticipantStatus) bool {
	return status == neo4jentity.FlowParticipantStatusCompleted || status == neo4jentity.FlowParticipantStatusGoalAchieved
}

func (s *flowService) determineNewStatus(subject string) neo4jentity.FlowParticipantStatus {
	if subject == "Welcome to Embedd - Product Tips" {
		return neo4jentity.FlowParticipantStatusGoalAchieved
	}
	return neo4jentity.FlowParticipantStatusCompleted
}

func (s *flowService) ProcessMailstackReply(ctx context.Context, tenant string, data *model.PostmarkEmailWebhookData, slackChannel string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.ProcessMailstackReply")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Get In-Reply-To header
	inReplyTo := s.getInReplyTo(data)
	if inReplyTo == "" {
		return nil
	}

	// Get original mailstack email
	mailstackEmail, err := s.services.CommonServices.PostgresRepositories.EmailMessageRepository.GetByProviderMessageId(ctx, tenant, inReplyTo)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Verify this is a valid reply to a flow action email
	if !s.isValidFlowReply(mailstackEmail, data.FromFull.Email) {
		return nil
	}

	// Process the flow action reply
	if err := s.processFlowActionReply(ctx, tenant, mailstackEmail, slackChannel); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *flowService) getInReplyTo(data *model.PostmarkEmailWebhookData) string {
	for _, header := range data.Headers {
		if header.Name == "In-Reply-To" {
			return header.Value
		}
	}
	return ""
}

func (s *flowService) isValidFlowReply(mailstackEmail *entity.EmailMessage, fromEmail string) bool {
	return mailstackEmail != nil &&
		strings.Contains(mailstackEmail.ToString, fromEmail) &&
		mailstackEmail.ProducerType == commonModel.NodeLabelFlowActionExecution
}

func (s *flowService) processFlowActionReply(ctx context.Context, tenant string, mailstackEmail *entity.EmailMessage, slackChannel string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.processFlowActionReply")
	defer span.Finish()

	// Get flow action execution
	flowActionExecution, err := s.services.CommonServices.FlowExecutionService.GetFlowActionExecutionById(ctx, mailstackEmail.ProducerId)
	if err != nil {
		return err
	}

	// Get flow participant
	flowParticipant, err := s.services.CommonServices.FlowService.FlowParticipantByEntity(
		ctx,
		flowActionExecution.FlowId,
		flowActionExecution.EntityId,
		flowActionExecution.EntityType,
	)
	if err != nil {
		return err
	}

	// Update participant status
	if err := s.updateParticipantStatus(ctx, tenant, flowParticipant); err != nil {
		return err
	}

	// Send notifications
	if err := s.sendNotifications(ctx, tenant, flowParticipant, slackChannel); err != nil {
		s.logger.Errorf("Error sending notifications: %v", err)
	}

	return nil
}

func (s *flowService) updateParticipantStatus(ctx context.Context, tenant string, participant *neo4jentity.FlowParticipantEntity) error {
	return s.services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateStringProperty(
		ctx,
		nil,
		tenant,
		commonModel.NodeLabelFlowParticipant,
		participant.Id,
		"status",
		string(neo4jentity.FlowParticipantStatusGoalAchieved),
	)
}

func (s *flowService) sendNotifications(ctx context.Context, tenant string, participant *neo4jentity.FlowParticipantEntity, slackChannel string) error {
	// Get primary email for participant
	primaryEmail, err := s.services.CommonServices.EmailService.GetPrimaryEmailForEntityId(
		ctx,
		participant.EntityType,
		participant.EntityId,
	)
	if err != nil {
		return err
	}

	if primaryEmail == nil {
		primaryEmail = &neo4jentity.EmailEntity{RawEmail: "primary email missing"}
	}

	// Send Slack notification if channel configured
	if slackChannel != "" {
		slackMessage := s.formatSlackMessage(tenant, primaryEmail.RawEmail)
		if err := utils.SendSlackMessage(ctx, slackChannel, slackMessage); err != nil {
			return err
		}
	}

	// Publish events
	if err := s.publishEvents(ctx, tenant, participant); err != nil {
		return err
	}

	return nil
}

func (s *flowService) formatSlackMessage(tenant, email string) string {
	return "*Tenant:* " + tenant + "\n*Goal achieved for:* " + email + "\n"
}

func (s *flowService) publishEvents(ctx context.Context, tenant string, participant *neo4jentity.FlowParticipantEntity) error {
	// Publish goal achieved event
	if err := s.services.CommonServices.RabbitMQService.PublishEvent(
		ctx,
		participant.FlowId,
		commonModel.FLOW,
		dto.FlowParticipantGoalAchieved{
			ParticipantId:   participant.EntityId,
			ParticipantType: participant.EntityType,
		},
	); err != nil {
		return err
	}

	// Publish event completed
	s.services.CommonServices.RabbitMQService.PublishEventCompleted(
		ctx,
		tenant,
		participant.Id,
		commonModel.FLOW_PARTICIPANT,
		utils.NewEventCompletedDetails().WithUpdate(),
	)

	return nil
}
