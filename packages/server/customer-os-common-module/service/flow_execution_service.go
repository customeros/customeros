package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

const numberOfEmailsPerDay = 2

type FlowExecutionService interface {
	ScheduleFlow(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, contactId string, entityType model.EntityType) error
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

func (s *flowExecutionService) ScheduleFlow(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId string, entityType model.EntityType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.ScheduleFlow")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	now := utils.Now()

	flowExecutions, err := s.getFlowActionExecutions(ctx, flowId, entityId, entityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if len(flowExecutions) == 0 {
		startAction, err := s.services.FlowService.FlowActionGetStart(ctx, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		nextActions, err := s.services.FlowService.FlowActionGetNext(ctx, startAction.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		for _, nextAction := range nextActions {

			scheduleAt := now.Add(time.Duration(nextAction.Data.WaitBefore) * time.Minute)

			err := s.scheduleNextAction(ctx, tx, flowId, entityId, entityType, scheduleAt, *nextAction)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}

	} else {
		lastActionExecution := flowExecutions[len(flowExecutions)-1]
		lastActionExecutedAt := lastActionExecution.ScheduledAt

		lastAction, err := s.services.FlowService.FlowActionGetById(ctx, lastActionExecution.ActionId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		nextActions, err := s.services.FlowService.FlowActionGetNext(ctx, lastAction.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		for _, nextAction := range nextActions {

			//marking the flow as completed if the next action is FLOW_END
			if nextAction.Data.Action == entity.FlowActionTypeFlowEnd {

				var flowParticipant *entity.FlowParticipantEntity

				//TODO support multiple entities
				if entityType == model.CONTACT {
					flowParticipant, err = s.services.FlowService.FlowParticipantByContactId(ctx, flowId, entityId)
					if err != nil {
						tracing.TraceErr(span, err)
						return err
					}
				}

				if flowParticipant == nil {
					tracing.TraceErr(span, errors.New("Flow participant not found"))
					return errors.New("Flow participant not found")
				}

				flowParticipant.Status = entity.FlowParticipantStatusCompleted

				_, err = s.services.Neo4jRepositories.FlowParticipantWriteRepository.Merge(ctx, tx, flowParticipant)
				if err != nil {
					tracing.TraceErr(span, err)
					return err
				}

				return nil
			}

			scheduleAt := lastActionExecutedAt.Add(time.Duration(nextAction.Data.WaitBefore) * time.Minute)

			err := s.scheduleNextAction(ctx, tx, flowId, entityId, entityType, scheduleAt, *nextAction)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
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

func (s *flowExecutionService) scheduleNextAction(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId string, entityType model.EntityType, scheduleAt time.Time, nextAction entity.FlowActionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.scheduleNextAction")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	switch nextAction.Data.Action {
	case entity.FlowActionTypeEmailNew, entity.FlowActionTypeEmailReply:
		return s.scheduleEmailAction(ctx, tx, flowId, entityId, entityType, scheduleAt, nextAction)
	default:
		tracing.TraceErr(span, fmt.Errorf("Unsupported action type %s", nextAction.Data.Action))
		return errors.New("Unsupported action type")
	}
}

func (s *flowExecutionService) scheduleEmailAction(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId string, entityType model.EntityType, scheduleAt time.Time, nextAction entity.FlowActionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.scheduleEmailAction")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	// 1. Get the mailbox for contact or associate the best available mailbox
	flowExecutionSettings, err := s.getFlowExecutionSettings(ctx, flowId, entityId, entityType)
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
			EntityId:   entityId,
			EntityType: entityType.String(),
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

	err = s.storeNextActionExecutionEntity(ctx, tx, flowId, nextAction.Id, entityId, entityType, flowExecutionSettings.Mailbox, *actualScheduleAt)
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

	startDate := utils.StartOfDayInUTC(scheduleAt)
	endDate := utils.EndOfDayInUTC(scheduleAt)

	emailsScheduledAlready, err := s.services.Neo4jRepositories.FlowActionExecutionReadRepository.CountEmailsPerMailboxPerDay(ctx, tx, mailbox, startDate, endDate)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if emailsScheduledAlready >= numberOfEmailsPerDay {
		// If the mailbox has already reached the daily limit, schedule the next email for the next day
		scheduleAt = startDate.Add(24 * time.Hour)
		// todo break recursive based on index or smth
		return s.getFirstAvailableSlotForMailbox(ctx, tx, tenant, mailbox, scheduleAt)
	} else {
		lastScheduledExecutionNode, err := s.services.Neo4jRepositories.FlowActionExecutionReadRepository.GetLastScheduledForMailbox(ctx, tx, mailbox)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		lastScheduledExecution := mapper.MapDbNodeToFlowActionExecutionEntity(lastScheduledExecutionNode)

		mailboxEntity, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, tenant, mailbox)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		//this is a bit wrong as will ocupy all the space in between the last scheduled email and the current schedule time
		if lastScheduledExecution != nil {
			if lastScheduledExecution.ScheduledAt.After(scheduleAt) {
				// If the last scheduled email is after the current schedule time, use the last scheduled time
				scheduleAt = lastScheduledExecution.ScheduledAt.Add(time.Duration(utils.GenerateRandomInt(mailboxEntity.MinMinutesBetweenEmails, mailboxEntity.MaxMinutesBetweenEmails)) * time.Minute)
				return &scheduleAt, nil
			} else {
				// If the last scheduled email is before the current schedule time, use the current schedule time
				return &scheduleAt, nil
			}
		} else {
			scheduleAt = scheduleAt.Add(time.Duration(utils.GenerateRandomInt(mailboxEntity.MinMinutesBetweenEmails, mailboxEntity.MaxMinutesBetweenEmails)) * time.Minute)
			return &scheduleAt, nil
		}
	}
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

func (s *flowExecutionService) storeNextActionExecutionEntity(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, actionId, entityId string, entityType model.EntityType, mailbox *string, executionTime time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.storeNextActionExecutionEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	id, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowActionExecution)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// After the wait duration, the next action will be executed
	actionExecution := entity.FlowActionExecutionEntity{
		Id:          id,
		FlowId:      flowId,
		ActionId:    actionId,
		EntityId:    entityId,
		EntityType:  entityType.String(),
		Mailbox:     mailbox,
		ScheduledAt: executionTime,
		Status:      entity.FlowActionExecutionStatusScheduled,
	}

	_, err = s.services.Neo4jRepositories.FlowActionExecutionWriteRepository.Merge(ctx, tx, &actionExecution)
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

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		currentAction, err := s.services.FlowService.FlowActionGetById(ctx, scheduledActionExecution.ActionId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if currentAction == nil {
			tracing.TraceErr(span, errors.New("Action not found"))
			return nil, errors.New("Action not found")
		}

		if currentAction.Data.Action == entity.FlowActionTypeEmailNew {

			mailbox, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, common.GetTenantFromContext(ctx), *scheduledActionExecution.Mailbox)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			if mailbox == nil {
				//mark the execution as failed
				tracing.TraceErr(span, errors.New("Mailbox not found"))
				return nil, errors.New("Mailbox not found")
			}

			toEmail := ""

			if model.GetEntityType(scheduledActionExecution.EntityType) == model.CONTACT {
				emailNodes, err := s.services.Neo4jRepositories.EmailReadRepository.GetAllEmailNodesForLinkedEntityIds(ctx, tenant, model.CONTACT, []string{scheduledActionExecution.EntityId})
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}

				if emailNodes == nil || len(emailNodes) == 0 {
					tracing.TraceErr(span, errors.New("Email not found"))
					return nil, errors.New("Email not found")
				}

				for _, emailNode := range emailNodes {
					emailEntity := mapper.MapDbNodeToEmailEntity(emailNode.Node)
					if emailEntity != nil && emailEntity.Work != nil && *emailEntity.Work {
						toEmail = emailEntity.RawEmail // TODO we should look for verified emails?
						break
					}
				}
			}

			if toEmail == "" {
				tracing.TraceErr(span, errors.New("Email not found"))
				return nil, errors.New("Email not found")
			}

			_, err = s.services.MailService.SendMail(ctx, &tx, tenant, dto.MailRequest{
				From:    *scheduledActionExecution.Mailbox,
				To:      []string{toEmail},
				Subject: currentAction.Data.Subject,
				Content: *currentAction.Data.BodyTemplate,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}
		//else if currentAction.ActionType == neo4jEntity.FlowActionTypeEmailReply {
		//	TODO reply to previous email
		//}

		scheduledActionExecution.Status = entity.FlowActionExecutionStatusSuccess

		_, err = s.services.Neo4jRepositories.FlowActionExecutionWriteRepository.Merge(ctx, &tx, scheduledActionExecution)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = s.services.FlowExecutionService.ScheduleFlow(ctx, &tx, scheduledActionExecution.FlowId, scheduledActionExecution.EntityId, model.GetEntityType(scheduledActionExecution.EntityType))
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
