package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

const (
	workingDayStart = 6
	workingDayEnd   = 18
)

type FlowExecutionService interface {
	GetFlowActionExecutionById(ctx context.Context, flowActionExecution string) (*entity.FlowActionExecutionEntity, error)
	GetFlowRequirements(ctx context.Context, flowId string) (*FlowComputeParticipantsRequirementsInput, error)
	UpdateParticipantFlowRequirements(ctx context.Context, tx *neo4j.ManagedTransaction, participant *entity.FlowParticipantEntity, requirements *FlowComputeParticipantsRequirementsInput) error
	ScheduleFlow(ctx context.Context, tx *neo4j.ManagedTransaction, flowId string, flowParticipant *entity.FlowParticipantEntity) error
	ProcessActionExecution(ctx context.Context, scheduledActionExecution *entity.FlowActionExecutionEntity) error
}

type flowExecutionService struct {
	services *Services
}

func NewFlowExecutionService(services *Services) FlowExecutionService {
	return &flowExecutionService{
		services: services,
	}
}

type FlowComputeParticipantsRequirementsInput struct {
	PrimaryEmailRequired bool
}

func (s *flowExecutionService) GetFlowActionExecutionById(ctx context.Context, flowActionExecution string) (*entity.FlowActionExecutionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetFlowActionExecutionById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowActionExecution", flowActionExecution))

	node, err := s.services.Neo4jRepositories.FlowActionExecutionReadRepository.GetById(ctx, flowActionExecution)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowActionExecutionEntity(node), nil
}

func (s *flowExecutionService) GetFlowRequirements(ctx context.Context, flowId string) (*FlowComputeParticipantsRequirementsInput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetFlowRequirements")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowId", flowId))

	requirements := FlowComputeParticipantsRequirementsInput{
		PrimaryEmailRequired: false,
	}

	flowActions, err := s.services.FlowService.FlowActionGetList(ctx, []string{flowId})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	for _, v := range *flowActions {
		if v.Data.Action == entity.FlowActionTypeEmailNew || v.Data.Action == entity.FlowActionTypeEmailReply {
			requirements.PrimaryEmailRequired = true
		}
	}

	return &requirements, nil
}

func (s *flowExecutionService) UpdateParticipantFlowRequirements(ctx context.Context, tx *neo4j.ManagedTransaction, participant *entity.FlowParticipantEntity, requirements *FlowComputeParticipantsRequirementsInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.UpdateParticipantFlowRequirements")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	status := entity.FlowParticipantStatusReady

	if requirements.PrimaryEmailRequired {
		//identify the primary email
		primaryEmail, err := s.services.EmailService.GetPrimaryEmailForEntityId(ctx, participant.EntityType, participant.EntityId)
		if err != nil {
			return errors.Wrap(err, "failed to get primary email for entity id")
		}

		if primaryEmail == nil {
			status = entity.FlowParticipantStatusOnHold
		}
	}

	err := s.services.Neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, tx, tenant, model.NodeLabelFlowParticipant, participant.Id, "status", string(status))
	if err != nil {
		return errors.Wrap(err, "failed to update string property")
	}

	participant.Status = status

	return nil
}

func (s *flowExecutionService) ScheduleFlow(ctx context.Context, tx *neo4j.ManagedTransaction, flowId string, flowParticipant *entity.FlowParticipantEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.ScheduleFlow")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	now := utils.Now()

	_, err := utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, tx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		flowExecutions, err := s.getFlowActionExecutions(ctx, &tx, flowId, flowParticipant.EntityId, flowParticipant.EntityType)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if len(flowExecutions) == 0 {
			startAction, err := s.services.FlowService.FlowActionGetStart(ctx, flowId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			nextActions, err := s.services.FlowService.FlowActionGetNext(ctx, startAction.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			for _, nextAction := range nextActions {

				scheduleAt := now.Add(time.Duration(nextAction.Data.WaitBefore) * time.Minute)

				err := s.scheduleNextAction(ctx, &tx, flowId, flowParticipant, scheduleAt, *nextAction)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
			}

		} else {
			lastActionExecution := flowExecutions[len(flowExecutions)-1]
			lastActionExecutedAt := lastActionExecution.ScheduledAt

			lastAction, err := s.services.FlowService.FlowActionGetById(ctx, lastActionExecution.ActionId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			nextActions, err := s.services.FlowService.FlowActionGetNext(ctx, lastAction.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			for _, nextAction := range nextActions {

				//marking the flow as completed if the next action is FLOW_END
				if nextAction.Data.Action == entity.FlowActionTypeFlowEnd {
					flowParticipant.Status = entity.FlowParticipantStatusCompleted

					_, err = s.services.Neo4jRepositories.FlowParticipantWriteRepository.Merge(ctx, &tx, flowParticipant)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}

					return nil, nil
				}

				scheduleAt := lastActionExecutedAt.Add(time.Duration(nextAction.Data.WaitBefore) * time.Minute)

				err := s.scheduleNextAction(ctx, &tx, flowId, flowParticipant, scheduleAt, *nextAction)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
			}

		}

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *flowExecutionService) getFlowActionExecutions(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId string, entityType model.EntityType) ([]*entity.FlowActionExecutionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetForContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	//get executions for contact
	nodes, err := s.services.Neo4jRepositories.FlowActionExecutionReadRepository.GetForEntity(ctx, tx, flowId, entityId, entityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make([]*entity.FlowActionExecutionEntity, 0)
	for _, node := range nodes {
		entities = append(entities, mapper.MapDbNodeToFlowActionExecutionEntity(node))
	}

	return entities, nil
}

func (s *flowExecutionService) scheduleNextAction(ctx context.Context, tx *neo4j.ManagedTransaction, flowId string, flowParticipant *entity.FlowParticipantEntity, scheduleAt time.Time, nextAction entity.FlowActionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.scheduleNextAction")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	switch nextAction.Data.Action {
	case entity.FlowActionTypeEmailNew, entity.FlowActionTypeEmailReply:
		return s.scheduleEmailAction(ctx, tx, flowId, flowParticipant, scheduleAt, nextAction)
	default:
		tracing.TraceErr(span, fmt.Errorf("Unsupported action type %s", nextAction.Data.Action))
		return errors.New("Unsupported action type")
	}
}

func (s *flowExecutionService) scheduleEmailAction(ctx context.Context, tx *neo4j.ManagedTransaction, flowId string, flowParticipant *entity.FlowParticipantEntity, scheduleAt time.Time, nextAction entity.FlowActionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.scheduleEmailAction")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	// 1. Get the mailbox for contact or associate the best available mailbox
	flowExecutionSettings, err := s.getFlowExecutionSettings(ctx, flowId, flowParticipant.EntityId, flowParticipant.EntityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowExecutionSettings == nil || flowExecutionSettings.Mailbox == nil {
		//compute the best available mailbox and associate

		// 1. get all available mailboxes
		// 2. select the mailbox with the fastest response time
		flowSenders, err := s.services.FlowService.FlowSenderGetList(ctx, []string{flowId})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		mailboxesScheduledAt := make(map[string]*time.Time)
		mailboxesScheduledAt[""] = utils.TimePtr(time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC))

		for _, flowActionSender := range *flowSenders {
			emailEntitites, err := s.services.EmailService.GetAllEmailsForEntityIds(ctx, tenant, model.USER, []string{*flowActionSender.UserId})
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			for _, emailEntity := range *emailEntitites {
				mailboxes, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetAllByUsername(ctx, tenant, emailEntity.RawEmail)
				if err != nil {
					tracing.TraceErr(span, err)
					return err
				}

				for _, mailbox := range mailboxes {
					scheduledAt, err := s.services.Neo4jRepositories.FlowActionExecutionReadRepository.GetFirstSlotForMailbox(ctx, tx, mailbox.MailboxUsername)
					if err != nil {
						tracing.TraceErr(span, err)
						return err
					}

					mailboxesScheduledAt[mailbox.MailboxUsername] = scheduledAt
				}
			}
		}

		fastestMailbox := ""
		for mailbox, scheduledAt := range mailboxesScheduledAt {
			if scheduledAt == nil {
				fastestMailbox = mailbox
				break
			} else if scheduledAt.Before(*mailboxesScheduledAt[fastestMailbox]) {
				fastestMailbox = mailbox
			}
		}

		if fastestMailbox == "" {
			tracing.TraceErr(span, errors.New("No mailbox available"))
			return errors.New("No mailbox available")
		}

		user, err := s.services.UserService.FindUserByEmail(ctx, fastestMailbox)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if user == nil {
			tracing.TraceErr(span, errors.New("User not found"))
			return errors.New("User not found")
		}

		id, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowExecutionSettings)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		flowExecutionSettings = &entity.FlowExecutionSettingsEntity{
			Id:         id,
			FlowId:     flowId,
			EntityId:   flowParticipant.EntityId,
			EntityType: flowParticipant.EntityType.String(),
			Mailbox:    &fastestMailbox,
			UserId:     &user.Id,
		}

		node, err := s.services.Neo4jRepositories.FlowExecutionSettingsWriteRepository.Merge(ctx, tx, flowExecutionSettings)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		flowExecutionSettings = mapper.MapDbNodeToFlowExecutionSettingsEntity(node)
	}

	workingSchedule, err := s.services.PostgresRepositories.UserWorkingScheduleRepository.GetForUser(ctx, tenant, *flowExecutionSettings.UserId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if len(workingSchedule) == 0 {
		return errors.New("User working schedule not found")
	}

	// 2. Schedule the email action
	actualScheduleAt, err := s.getFirstAvailableSlotForMailbox(ctx, tx, tenant, *flowExecutionSettings.Mailbox, scheduleAt, workingSchedule)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.storeNextActionExecutionEntity(ctx, tx, flowId, nextAction.Id, flowParticipant, flowExecutionSettings.Mailbox, actualScheduleAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowParticipant.Status = entity.FlowParticipantStatusScheduled
	_, err = s.services.Neo4jRepositories.FlowParticipantWriteRepository.Merge(ctx, tx, flowParticipant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *flowExecutionService) getFirstAvailableSlotForMailbox(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, mailbox string, scheduleAt time.Time, workingSchedule []*postgresentity.UserWorkingSchedule) (*time.Time, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.getFirstAvailableSlotForMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	mailboxEntity, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, tenant, mailbox)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	//minTimeBetweenEmails := time.Duration(mailboxEntity.MinMinutesBetweenEmails) * time.Minute

	// Get the last scheduled execution for this mailbox
	lastScheduledExecutionNode, err := s.services.Neo4jRepositories.FlowActionExecutionReadRepository.GetLastScheduledForMailbox(ctx, tx, mailbox)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	possibleScheduledAt := scheduleAt

	if lastScheduledExecutionNode != nil {
		lastScheduledExecution := mapper.MapDbNodeToFlowActionExecutionEntity(lastScheduledExecutionNode)
		possibleScheduledAt = maxTime(possibleScheduledAt, lastScheduledExecution.ScheduledAt)
	}

	//check the number of emails scheduled for the day
	for {
		emailsScheduledInDay, err := s.services.Neo4jRepositories.FlowActionExecutionReadRepository.CountEmailsPerMailboxPerDay(ctx, tx, mailbox, utils.StartOfDayInUTC(possibleScheduledAt), utils.EndOfDayInUTC(possibleScheduledAt))
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if emailsScheduledInDay >= int64(mailboxEntity.RampUpCurrent) {
			possibleScheduledAt = possibleScheduledAt.AddDate(0, 0, 1)
			possibleScheduledAt = time.Date(possibleScheduledAt.Year(), possibleScheduledAt.Month(), possibleScheduledAt.Day(), 0, 0, 0, 0, time.UTC)
			continue
		} else {
			break
		}
	}

	// Ensure possibleScheduledAt is not in the past and within working hours
	possibleScheduledAt = adjustToWorkingTimeWithRandom(maxTime(possibleScheduledAt, utils.Now()), workingSchedule, mailboxEntity)

	//Add random seconds and miliseconds to not have 00:00:00 as the scheduled time
	randomSeconds := time.Duration(utils.GenerateRandomInt(0, 60)) * time.Second
	randomMiliseconds := time.Duration(utils.GenerateRandomInt(0, 1000)) * time.Millisecond
	randomMicroseconds := time.Duration(utils.GenerateRandomInt(0, 1000)) * time.Microsecond
	possibleScheduledAt = possibleScheduledAt.Add(randomSeconds).Add(randomMiliseconds).Add(randomMicroseconds)

	return &possibleScheduledAt, nil

	//TODO V2
	//for {
	//	endTime := possibleScheduledAt.Add(minTimeBetweenEmails)
	//
	//	// Check if there's any scheduled execution within the interval
	//	ee, err := s.services.Neo4jRepositories.FlowActionExecutionReadRepository.GetByMailboxAndTimeInterval(ctx, tx, mailbox, possibleScheduledAt, endTime)
	//	if err != nil {
	//		tracing.TraceErr(span, err)
	//		return nil, err
	//	}
	//
	//	existingExecution := mapper.MapDbNodeToFlowActionExecutionEntity(ee)
	//
	//	if existingExecution == nil {
	//		// No execution found in the interval, so this slot is available
	//		// Add a random duration within the min-max range
	//		randomDuration := time.Duration(utils.GenerateRandomInt(mailboxEntity.MinMinutesBetweenEmails, mailboxEntity.MaxMinutesBetweenEmails)) * time.Minute
	//		scheduledTime := possibleScheduledAt.Add(randomDuration)
	//
	//		// Ensure the scheduled time is within working hours
	//		scheduledTime = nextWorkingTime(scheduledTime)
	//
	//		return &scheduledTime, nil
	//	}
	//
	//	// Move the start time to just after the found execution
	//	possibleScheduledAt = possibleScheduledAt.Add(time.Minute)
	//	possibleScheduledAt = nextWorkingTime(possibleScheduledAt)
	//}
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func (s *flowExecutionService) getFlowExecutionSettings(ctx context.Context, flowId, entityId string, entityType model.EntityType) (*entity.FlowExecutionSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.getFlowExecutionSettings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	node, err := s.services.Neo4jRepositories.FlowExecutionSettingsReadRepository.GetByEntityId(ctx, flowId, entityId, entityType.String())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowExecutionSettingsEntity(node), nil
}

func (s *flowExecutionService) storeNextActionExecutionEntity(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, actionId string, flowParticipant *entity.FlowParticipantEntity, mailbox *string, executionTime *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.storeNextActionExecutionEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	id, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowActionExecution)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = s.services.Neo4jRepositories.FlowActionExecutionWriteRepository.Merge(ctx, tx, &entity.FlowActionExecutionEntity{
		Id:              id,
		FlowId:          flowId,
		ActionId:        actionId,
		EntityId:        flowParticipant.EntityId,
		EntityType:      flowParticipant.EntityType,
		Mailbox:         mailbox,
		ScheduledAt:     *executionTime,
		StatusUpdatedAt: utils.Now(),
		Status:          entity.FlowActionExecutionStatusScheduled,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *flowExecutionService) ProcessActionExecution(ctx context.Context, scheduledActionExecution *entity.FlowActionExecutionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.ProcessActionExecution")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogFields(log.Object("scheduledActionExecution", scheduledActionExecution))

	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	shouldInsertEmailMessage := false
	var emailMessage *postgresentity.EmailMessage
	var currentAction *entity.FlowActionEntity

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var err error
		currentAction, err = s.services.FlowService.FlowActionGetById(ctx, scheduledActionExecution.ActionId)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get action by id")
		}

		if currentAction == nil {
			return nil, errors.New("action not found")
		}

		//check if the participant meets flow requirements
		flowRequirements, err := s.services.FlowExecutionService.GetFlowRequirements(ctx, scheduledActionExecution.FlowId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		participant, err := s.services.FlowService.FlowParticipantByEntity(ctx, scheduledActionExecution.FlowId, scheduledActionExecution.EntityId, scheduledActionExecution.EntityType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get flow participant by entity")
		}

		if participant == nil {
			return nil, errors.New("participant not found")
		}

		err = s.UpdateParticipantFlowRequirements(ctx, &tx, participant, flowRequirements)
		if err != nil {
			return nil, errors.Wrap(err, "failed to update participant flow requirements")
		}

		//if participant is not ready, mark the action as business error
		if participant.Status != entity.FlowParticipantStatusReady {

			scheduledActionExecution.StatusUpdatedAt = utils.Now()
			scheduledActionExecution.Status = entity.FlowActionExecutionStatusBusinessError

			_, err = s.services.Neo4jRepositories.FlowActionExecutionWriteRepository.Merge(ctx, &tx, scheduledActionExecution)
			if err != nil {
				return nil, errors.Wrap(err, "failed to merge flow action execution")
			}

			return nil, nil
		}

		if currentAction.Data.Action == entity.FlowActionTypeEmailNew || currentAction.Data.Action == entity.FlowActionTypeEmailReply {

			//prevent duplicate emails for same action
			existingEmail, err := s.services.PostgresRepositories.EmailMessageRepository.GetByProducer(ctx, tenant, scheduledActionExecution.Id, model.NodeLabelFlowActionExecution)
			if err != nil {
				//todo this is producing an error below when trying to insert the email. need to rewrite how we store the email
				return nil, errors.Wrap(err, "failed to get email by producer")
			}

			if existingEmail == nil {

				mailbox, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, common.GetTenantFromContext(ctx), *scheduledActionExecution.Mailbox)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get mailbox by mailbox")
				}

				if mailbox == nil {
					return nil, errors.New("mailbox not found in database")
				}

				primaryEmail, err := s.services.EmailService.GetPrimaryEmailForEntityId(ctx, participant.EntityType, participant.EntityId)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get primary email for entity id")
				}

				if primaryEmail == nil {
					return nil, errors.New("primary email not found")
				}

				toEmail := primaryEmail.RawEmail

				bodyTemplate := *currentAction.Data.BodyTemplate

				if scheduledActionExecution.EntityType == model.CONTACT {
					contactNode, err := s.services.Neo4jRepositories.ContactReadRepository.GetContact(ctx, tenant, scheduledActionExecution.EntityId)
					if err != nil {
						return nil, errors.Wrap(err, "failed to get contact")
					}

					contact := mapper.MapDbNodeToContactEntity(contactNode)

					firstName, lastName := contact.DeriveFirstAndLastNames()
					bodyTemplate = replacePlaceholders(bodyTemplate, "contact_first_name", firstName)
					bodyTemplate = replacePlaceholders(bodyTemplate, "contact_last_name", lastName)
					bodyTemplate = replacePlaceholders(bodyTemplate, "contact_email", toEmail)

					contactWithOrganizations, err := s.services.OrganizationService.GetLatestOrganizationsWithJobRolesForContacts(ctx, []string{contact.Id})
					if err != nil {
						return nil, errors.Wrap(err, "failed to get latest organizations with job roles for contacts")
					}

					if len(*contactWithOrganizations) > 0 {
						contactWithOrganization := (*contactWithOrganizations)[0]
						bodyTemplate = replacePlaceholders(bodyTemplate, "organization_name", contactWithOrganization.Organization.Name)
					} else {
						bodyTemplate = replacePlaceholders(bodyTemplate, "organization_name", "")
					}
				}

				userNode, err := s.services.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, *scheduledActionExecution.Mailbox)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get first user by email")
				}

				if userNode == nil {
					return nil, errors.New("user not found")
				}

				user := mapper.MapDbNodeToUserEntity(userNode)

				bodyTemplate = replacePlaceholders(bodyTemplate, "sender_first_name", user.FirstName)
				bodyTemplate = replacePlaceholders(bodyTemplate, "sender_last_name", user.LastName)

				shouldInsertEmailMessage = true
				emailMessage = &postgresentity.EmailMessage{
					Status:       postgresentity.EmailMessageStatusScheduled,
					ProducerId:   scheduledActionExecution.Id,
					ProducerType: model.NodeLabelFlowActionExecution,
					FromName:     user.FirstName + " " + user.LastName,
					From:         *scheduledActionExecution.Mailbox,
					To:           []string{toEmail},
					Content:      bodyTemplate,
				}

				if currentAction.Data.Action == entity.FlowActionTypeEmailNew {
					emailMessage.Subject = *currentAction.Data.Subject
				}

				if currentAction.Data.Action == entity.FlowActionTypeEmailReply {
					//walk back the flow and identify the previous email
					//get previous email execution from neo4j
					//get previous email from postgres
					//reply to the previous email

					parentEmailAction, err := s.getEmailActionToReply(ctx, scheduledActionExecution.ActionId)
					if err != nil {
						return nil, errors.Wrap(err, "failed to get email action to reply")
					}

					if parentEmailAction == nil {
						return nil, errors.New("no parent email action found")
					}

					parentEmailExecution, err := s.services.Neo4jRepositories.FlowActionExecutionReadRepository.GetExecution(ctx, scheduledActionExecution.FlowId, parentEmailAction.Id, scheduledActionExecution.EntityId, scheduledActionExecution.EntityType)
					if err != nil {
						return nil, errors.Wrap(err, "failed to get execution")
					}

					if parentEmailExecution == nil {
						return nil, errors.New("no parent email execution found")
					}

					parentEmail := mapper.MapDbNodeToFlowActionExecutionEntity(parentEmailExecution)

					parentEmailSent, err := s.services.PostgresRepositories.EmailMessageRepository.GetByProducer(ctx, tenant, parentEmail.Id, model.NodeLabelFlowActionExecution)
					if err != nil {
						return nil, errors.Wrap(err, "failed to get email by producer")
					}

					if parentEmailSent == nil {
						return nil, errors.New("no parent email sent found")
					}

					emailMessage.Subject = "Re: " + parentEmailSent.Subject
					emailMessage.ProviderInReplyTo = parentEmailSent.ProviderMessageId
					emailMessage.ProviderReferences = parentEmailSent.ProviderReferences + " " + parentEmailSent.ProviderMessageId
				}

			}
		}

		scheduledActionExecution.ExecutedAt = utils.TimePtr(utils.Now())
		scheduledActionExecution.StatusUpdatedAt = utils.Now()
		scheduledActionExecution.Status = entity.FlowActionExecutionStatusSuccess

		_, err = s.services.Neo4jRepositories.FlowActionExecutionWriteRepository.Merge(ctx, &tx, scheduledActionExecution)
		if err != nil {
			return nil, errors.Wrap(err, "failed to merge flow action execution")
		}

		flowParticipant, err := s.services.FlowService.FlowParticipantByEntity(ctx, scheduledActionExecution.FlowId, scheduledActionExecution.EntityId, scheduledActionExecution.EntityType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get flow participant by entity")
		}

		err = s.services.FlowExecutionService.ScheduleFlow(ctx, &tx, scheduledActionExecution.FlowId, flowParticipant)
		if err != nil {
			return nil, errors.Wrap(err, "failed to schedule flow")
		}

		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if shouldInsertEmailMessage {
		if emailMessage == nil {
			tracing.TraceErr(span, errors.New("email message is nil"))
			return errors.New("email message is nil")
		}

		//store in PG after the neo4j transaction is committed
		err = s.services.PostgresRepositories.EmailMessageRepository.Store(ctx, tenant, emailMessage)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "failed to store email message")
		}
	}

	// save billable event
	_, err = s.services.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgresentity.BillableEventFlowActionExecuted,
		postgresrepository.BillableEventDetails{
			Subtype:       string(currentAction.Data.Action),
			ReferenceData: fmt.Sprintf("FlowActionExecutionId: %s, storedEmailMessage: %v", scheduledActionExecution.Id, shouldInsertEmailMessage),
		})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to store billable event"))
	}

	return nil
}

func (s *flowExecutionService) getEmailActionToReply(ctx context.Context, actionId string) (*entity.FlowActionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.getEmailActionToReply")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	var previous *entity.FlowActionEntity

	previousNodes, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetPrevious(ctx, actionId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get previous nodes for action")
	}

	if previousNodes == nil || len(previousNodes) == 0 {
		return nil, errors.New("no previous nodes found for action")
	}

	if len(previousNodes) == 1 {
		previousEntity := mapper.MapDbNodeToFlowActionEntity(previousNodes[0])

		if previousEntity.Data.Action == entity.FlowActionTypeFlowStart {
			return nil, errors.New("reached start of flow")
		} else if previousEntity.Data.Action == entity.FlowActionTypeEmailNew {
			previous = previousEntity
		} else {
			return s.getEmailActionToReply(ctx, previousEntity.Id)
		}
	}

	// Handle multiple previous nodes
	for _, previousNode := range previousNodes {
		previousEntity := mapper.MapDbNodeToFlowActionEntity(previousNode)

		// Base case: If we find an email action, return it
		if previousEntity.Data.Action == entity.FlowActionTypeEmailNew {
			return previousEntity, nil
		}

		// Recursively search for an email action
		result, err := s.getEmailActionToReply(ctx, previousEntity.Id)
		if err == nil && result != nil {
			return result, nil
		}
	}

	return previous, nil
}

func replacePlaceholders(input, variableName, value string) string {
	return strings.Replace(input, "{{"+variableName+"}}", value, -1)
}
func adjustToWorkingTimeWithRandom(t time.Time, schedules []*postgresentity.UserWorkingSchedule, mailbox *postgresentity.TenantSettingsMailbox) time.Time {
	for {
		// Get working hours for the current day
		start, end := getWorkingHoursForDay(t, schedules)

		randomMinutes := time.Duration(utils.GenerateRandomInt(mailbox.MinMinutesBetweenEmails, mailbox.MaxMinutesBetweenEmails)) * time.Minute
		t = t.Add(randomMinutes)

		if !start.IsZero() && !end.IsZero() {
			// Check if time is within working hours
			if t.After(start) && t.Before(end) {
				return t // It's within working hours
			}
			if t.Before(start) {
				// Move to the start of today's working hours
				return start
			}
		}

		// Move to the next day's start time
		t = time.Date(t.Year(), t.Month(), t.Day(), start.Hour(), start.Minute(), 0, 0, time.UTC).AddDate(0, 0, 1)
	}
}

// Helper to get the start and end times for the current weekday based on schedules
func getWorkingHoursForDay(day time.Time, schedules []*postgresentity.UserWorkingSchedule) (time.Time, time.Time) {
	dayStr := day.Weekday().String()[:3] // Get day abbreviation, e.g., "Mon"
	for _, schedule := range schedules {
		if IsDayInRange(dayStr, schedule.DayRange) {

			startParts := strings.Split(schedule.StartHour, ":")
			endParts := strings.Split(schedule.EndHour, ":")

			startHour, err := strconv.Atoi(startParts[0])
			if err != nil {
				return time.Time{}, time.Time{}
			}

			startMinute, err := strconv.Atoi(startParts[1])
			if err != nil {
				return time.Time{}, time.Time{}
			}

			endHour, err := strconv.Atoi(endParts[0])
			if err != nil {
				return time.Time{}, time.Time{}
			}

			endMinute, err := strconv.Atoi(endParts[1])
			if err != nil {
				return time.Time{}, time.Time{}
			}

			start := time.Date(day.Year(), day.Month(), day.Day(), startHour, startMinute, 0, 0, time.UTC)
			end := time.Date(day.Year(), day.Month(), day.Day(), endHour, endMinute, 0, 0, time.UTC)

			if day.After(end) {
				return time.Time{}, time.Time{} // No working hours for this day
			}

			return start, end
		}
	}
	return time.Time{}, time.Time{} // No working hours for this day
}

// Helper to get the earliest working start time of the next working day
func startOfNextWorkingDay(t time.Time, schedules []*postgresentity.UserWorkingSchedule) time.Time {
	for {
		start, _ := getWorkingHoursForDay(t, schedules)
		if !start.IsZero() {
			return start
		}
		t = t.AddDate(0, 0, 1) // Move to the next day
	}
}

// Helper function to check if a day is within a day range like "Mon-Wed"
func IsDayInRange(day, dayRange string) bool {
	daysOfWeek := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

	if dayRange == day+"-"+day { // For single day entries like "Mon-Mon"
		return true
	}

	rangeParts := strings.Split(dayRange, "-")
	if len(rangeParts) != 2 {
		return false
	}

	startIdx, endIdx := indexOf(daysOfWeek, rangeParts[0]), indexOf(daysOfWeek, rangeParts[1])
	dayIdx := indexOf(daysOfWeek, day)

	return dayIdx >= startIdx && dayIdx <= endIdx
}

// Helper to find the index of a day in the daysOfWeek slice
func indexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}
