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
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

const (
	workingDayStart = 9
	workingDayEnd   = 18
)

type FlowExecutionService interface {
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

func (s *flowExecutionService) ScheduleFlow(ctx context.Context, tx *neo4j.ManagedTransaction, flowId string, flowParticipant *entity.FlowParticipantEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.ScheduleFlow")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	now := utils.Now()

	_, err := utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, tx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		flowExecutions, err := s.getFlowActionExecutions(ctx, flowId, flowParticipant.EntityId, flowParticipant.EntityType)
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

func (s *flowExecutionService) getFlowActionExecutions(ctx context.Context, flowId, entityId string, entityType model.EntityType) ([]*entity.FlowActionExecutionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetForContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	//get executions for contact
	nodes, err := s.services.Neo4jRepositories.FlowActionExecutionReadRepository.GetForEntity(ctx, flowId, entityId, entityType.String())
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
			UserId:     nil,
		}

		node, err := s.services.Neo4jRepositories.FlowExecutionSettingsWriteRepository.Merge(ctx, tx, flowExecutionSettings)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		flowExecutionSettings = mapper.MapDbNodeToFlowExecutionSettingsEntity(node)
	}

	// 2. Schedule the email action
	actualScheduleAt, err := s.getFirstAvailableSlotForMailbox(ctx, tx, tenant, *flowExecutionSettings.Mailbox, scheduleAt)
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

func (s *flowExecutionService) getFirstAvailableSlotForMailbox(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, mailbox string, scheduleAt time.Time) (*time.Time, error) {
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

	var possibleScheduledAt time.Time
	if lastScheduledExecutionNode != nil {
		lastScheduledExecution := mapper.MapDbNodeToFlowActionExecutionEntity(lastScheduledExecutionNode)
		possibleScheduledAt = lastScheduledExecution.ScheduledAt
	} else {
		possibleScheduledAt = scheduleAt
	}

	//check the number of emails scheduled for the day
	for {
		emailsScheduledInDay, err := s.services.Neo4jRepositories.FlowActionExecutionReadRepository.CountEmailsPerMailboxPerDay(ctx, tx, mailbox, utils.StartOfDayInUTC(possibleScheduledAt), utils.EndOfDayInUTC(possibleScheduledAt))
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if emailsScheduledInDay >= int64(mailboxEntity.MaxEmailsPerDay) {
			possibleScheduledAt = possibleScheduledAt.AddDate(0, 0, 1)
			possibleScheduledAt = time.Date(possibleScheduledAt.Year(), possibleScheduledAt.Month(), possibleScheduledAt.Day(), 0, 0, 0, 0, time.UTC)
			continue
		} else {
			break
		}
	}

	// Ensure possibleScheduledAt is not in the past and within working hours
	possibleScheduledAt = nextWorkingTime(maxTime(possibleScheduledAt, time.Now().UTC()))

	randomMinutes := time.Duration(utils.GenerateRandomInt(mailboxEntity.MinMinutesBetweenEmails, mailboxEntity.MaxMinutesBetweenEmails)) * time.Minute
	possibleScheduledAt = possibleScheduledAt.Add(randomMinutes)

	// Ensure the scheduled time is within working hours
	possibleScheduledAt = nextWorkingTime(possibleScheduledAt)

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
		Id:          id,
		FlowId:      flowId,
		ActionId:    actionId,
		EntityId:    flowParticipant.EntityId,
		EntityType:  flowParticipant.EntityType,
		Mailbox:     mailbox,
		ScheduledAt: *executionTime,
		Status:      entity.FlowActionExecutionStatusScheduled,
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

	var emailMessage *postgresEntity.EmailMessage

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		currentAction, err := s.services.FlowService.FlowActionGetById(ctx, scheduledActionExecution.ActionId)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get action by id")
		}

		if currentAction == nil {
			return nil, errors.New("action not found")
		}

		//TODO verify if the action execution hasn't produced an email in the email table by producerId and producerType

		if currentAction.Data.Action == entity.FlowActionTypeEmailNew {

			mailbox, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, common.GetTenantFromContext(ctx), *scheduledActionExecution.Mailbox)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get mailbox by mailbox")
			}

			if mailbox == nil {
				return nil, errors.New("mailbox not found in database")
			}

			toEmail := ""

			//identify the primary work email associated with the entity
			emailNodes, err := s.services.Neo4jRepositories.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, tenant, scheduledActionExecution.EntityType, []string{scheduledActionExecution.EntityId})
			if err != nil {
				return nil, errors.Wrap(err, "failed to get all email nodes for linked entity ids")
			}

			if emailNodes == nil || len(emailNodes) == 0 {
				return nil, errors.New("no email nodes found for linked entity ids")
			}

			for _, emailNode := range emailNodes {
				emailEntity := mapper.MapDbNodeToEmailEntity(emailNode.Node)
				if emailEntity != nil && emailEntity.Work != nil && *emailEntity.Work {
					toEmail = emailEntity.RawEmail // TODO we should look for verified emails?
					break
				}
			}

			if toEmail == "" {
				return nil, errors.New("no work email found for entity")
			}

			emailMessage = &postgresEntity.EmailMessage{
				Status:       postgresEntity.EmailMessageStatusScheduled,
				ProducerId:   scheduledActionExecution.Id,
				ProducerType: model.NodeLabelFlowActionExecution,
				From:         *scheduledActionExecution.Mailbox,
				To:           []string{toEmail},
				Subject:      *currentAction.Data.Subject,
				Content:      *currentAction.Data.BodyTemplate,
			}
		}
		//else if currentAction.ActionType == neo4jEntity.FlowActionTypeEmailReply {
		//	TODO reply to previous email
		//}

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

	return nil
}

func isWorkingDay(t time.Time) bool {
	weekday := t.Weekday()
	return weekday >= time.Monday && weekday <= time.Friday
}

func isWithinWorkingHours(t time.Time) bool {
	if !isWorkingDay(t) {
		return false
	}
	hour := t.UTC().Hour()
	return hour >= workingDayStart && hour < workingDayEnd
}

func nextWorkingTime(t time.Time) time.Time {
	t = t.UTC()
	for !isWithinWorkingHours(t) {
		if !isWorkingDay(t) {
			// Move to next day at 9:00 UTC
			t = time.Date(t.Year(), t.Month(), t.Day()+1, workingDayStart, 0, 0, 0, time.UTC)
		} else if t.Hour() < workingDayStart {
			// Move to 9:00 UTC same day
			t = time.Date(t.Year(), t.Month(), t.Day(), workingDayStart, 0, 0, 0, time.UTC)
		} else {
			// Move to 9:00 UTC next day
			t = time.Date(t.Year(), t.Month(), t.Day()+1, workingDayStart, 0, 0, 0, time.UTC)
		}
	}
	return t
}
