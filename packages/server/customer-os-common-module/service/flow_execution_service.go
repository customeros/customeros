package service

import (
	"context"
	"encoding/json"
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

type FlowExecutionService interface {
	GetFlowActionExecutionById(ctx context.Context, flowActionExecution string) (*entity.FlowActionExecutionEntity, error)
	GetFlowExecutionSettingsForEntity(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId string, entityType model.EntityType) (*entity.FlowExecutionSettingsEntity, error)
	GetFlowRequirements(ctx context.Context, flowId string) (*FlowComputeParticipantsRequirementsInput, error)
	GetFlowActionExecutionForParticipantWithActionType(ctx context.Context, entityId string, entityType model.EntityType, actionType entity.FlowActionType) ([]*entity.FlowActionExecutionEntity, error)
	UpdateParticipantFlowRequirements(ctx context.Context, tx *neo4j.ManagedTransaction, participant *entity.FlowParticipantEntity, requirements *FlowComputeParticipantsRequirementsInput) (bool, error)
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
	PrimaryEmailRequired      bool
	LinkedInSocialUrlRequired bool
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

func (s *flowExecutionService) GetFlowExecutionSettingsForEntity(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId string, entityType model.EntityType) (*entity.FlowExecutionSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetFlowExecutionSettingsForEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	node, err := s.services.Neo4jRepositories.FlowExecutionSettingsReadRepository.GetForEntity(ctx, tx, flowId, entityId, entityType.String())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowExecutionSettingsEntity(node), nil
}

func (s *flowExecutionService) GetFlowRequirements(ctx context.Context, flowId string) (*FlowComputeParticipantsRequirementsInput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetFlowRequirements")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowId", flowId))

	requirements := FlowComputeParticipantsRequirementsInput{
		PrimaryEmailRequired:      false,
		LinkedInSocialUrlRequired: false,
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
		if v.Data.Action == entity.FlowActionTypeLinkedinConnectionRequest || v.Data.Action == entity.FlowActionTypeLinkedinMessage {
			requirements.LinkedInSocialUrlRequired = true
		}
	}

	return &requirements, nil
}

func (s *flowExecutionService) GetFlowActionExecutionForParticipantWithActionType(ctx context.Context, entityId string, entityType model.EntityType, actionType entity.FlowActionType) ([]*entity.FlowActionExecutionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetFlowActionExecutionForParticipantWithActionType")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	nodes, err := s.services.Neo4jRepositories.FlowActionExecutionReadRepository.GetForEntityWithActionType(ctx, entityId, entityType, actionType)
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

func (s *flowExecutionService) UpdateParticipantFlowRequirements(ctx context.Context, tx *neo4j.ManagedTransaction, participant *entity.FlowParticipantEntity, requirements *FlowComputeParticipantsRequirementsInput) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.UpdateParticipantFlowRequirements")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	status := entity.FlowParticipantStatusReady

	if requirements.PrimaryEmailRequired {
		//identify the primary email
		primaryEmail, err := s.services.EmailService.GetPrimaryEmailForEntityId(ctx, participant.EntityType, participant.EntityId)
		if err != nil {
			return false, errors.Wrap(err, "failed to get primary email for entity id")
		}

		if primaryEmail == nil {
			status = entity.FlowParticipantStatusOnHold
		}
	}

	if requirements.LinkedInSocialUrlRequired {
		socials, err := s.services.SocialService.GetAllForEntities(ctx, tenant, participant.EntityType, []string{participant.EntityId})
		if err != nil {
			return false, errors.Wrap(err, "failed to get socials for entities")
		}
		found := false
		for _, social := range *socials {
			if strings.Contains(social.Url, "linkedin.com") {
				found = true
				break
			}
		}

		if !found {
			status = entity.FlowParticipantStatusOnHold
		}
	}

	if participant.Status == status {
		return false, nil
	}

	err := s.services.Neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, tx, tenant, model.NodeLabelFlowParticipant, participant.Id, "status", string(status))
	if err != nil {
		return false, errors.Wrap(err, "failed to update string property")
	}

	participant.Status = status

	return true, nil
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

	//check if the participant meets flow requirements
	flowRequirements, err := s.services.FlowExecutionService.GetFlowRequirements(ctx, flowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "failed to get flow requirements")
	}

	_, err = s.UpdateParticipantFlowRequirements(ctx, tx, flowParticipant, flowRequirements)
	if err != nil {
		return errors.Wrap(err, "failed to update participant flow requirements")
	}

	if flowParticipant.Status != entity.FlowParticipantStatusReady {
		//todo return business error
		return nil
	}

	switch nextAction.Data.Action {
	case entity.FlowActionTypeEmailNew, entity.FlowActionTypeEmailReply:
		return s.scheduleEmailAction(ctx, tx, flowId, flowParticipant, scheduleAt, nextAction)
	case entity.FlowActionTypeLinkedinConnectionRequest:
		return s.scheduleSendLinkedInConnection(ctx, tx, flowId, flowParticipant, scheduleAt, nextAction)
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
	flowExecutionSettings, err := s.GetFlowExecutionSettingsForEntity(ctx, tx, flowId, flowParticipant.EntityId, flowParticipant.EntityType)
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
				mailboxes, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetAllByUsername(ctx, emailEntity.RawEmail)
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

		flowExecutionSettings, err = s.upsertFlowExecutionSettings(ctx, tx, tenant, flowId, flowParticipant, &fastestMailbox, &user.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
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
	actualScheduleAt, err := s.getFirstAvailableSlotForMailbox(ctx, tx, *flowExecutionSettings.Mailbox, scheduleAt, workingSchedule)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = s.storeNextActionExecutionEntity(ctx, tx, flowId, nextAction.Id, flowParticipant, flowExecutionSettings.Mailbox, actualScheduleAt)
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

func (s *flowExecutionService) scheduleSendLinkedInConnection(ctx context.Context, tx *neo4j.ManagedTransaction, flowId string, flowParticipant *entity.FlowParticipantEntity, scheduleAt time.Time, nextAction entity.FlowActionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.scheduleSendLinkedInConnection")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	if flowParticipant.EntityType != model.CONTACT {
		return errors.New("Only contacts are supported for LinkedIn connection requests")
	}

	flowSenders, err := s.services.FlowService.FlowSenderGetList(ctx, []string{flowId})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	senderIds := make([]string, 0)
	for _, flowSender := range *flowSenders {
		if flowSender.UserId == nil {
			continue
		}

		activeLinkedinToken, err := s.services.PostgresRepositories.BrowserConfigRepository.GetForUser(ctx, *flowSender.UserId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if activeLinkedinToken != nil {
			senderIds = append(senderIds, *flowSender.UserId)
		}
	}

	span.LogFields(log.String("senderIds", strings.Join(senderIds, ",")))

	if len(senderIds) == 0 {
		return errors.New("No LinkedIn sender available")
	}

	// 1 - Sender is already connected with the contact
	// - create the scheduled execution as EXECUTED
	// - set the UserId as the user to be used in the flow
	// - schedule the next action

	for _, senderId := range senderIds {
		isLinkedWith, err := s.services.Neo4jRepositories.CommonReadRepository.IsLinkedWith(ctx, tenant, flowParticipant.EntityId, model.CONTACT, model.CONNECTED_WITH.String(), senderId, model.USER)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "CommonReadRepository.IsLinkedWith"))
			return err
		}

		if isLinkedWith {
			span.LogFields(log.String("process", senderId+" is already connected with the contact"))
			id, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowActionExecution)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			now := utils.Now()
			_, err = s.services.Neo4jRepositories.FlowActionExecutionWriteRepository.Merge(ctx, tx, &entity.FlowActionExecutionEntity{
				Id:              id,
				FlowId:          flowId,
				ActionId:        nextAction.Id,
				EntityId:        flowParticipant.EntityId,
				EntityType:      flowParticipant.EntityType,
				UserId:          &senderId,
				ExecutedAt:      &now,
				ScheduledAt:     now,
				StatusUpdatedAt: now,
				Status:          entity.FlowActionExecutionStatusSkipped,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			_, err = s.upsertFlowExecutionSettings(ctx, tx, tenant, flowId, flowParticipant, nil, &senderId)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			err = s.ScheduleFlow(ctx, tx, flowId, flowParticipant)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			return nil
		}
	}

	socials, err := s.services.SocialService.GetAllForEntities(ctx, tenant, model.CONTACT, []string{flowParticipant.EntityId})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	socialUrl := ""
	requestSentAlready := false

	// 2 - No flow sender is connected with the contact
	// - identify the fastest user to connect with
	// - store flow execution settings
	// - create the scheduled execution as SCHEDULED

	// - check LinkedinConnectionRequest if there is already a request scheduled in the last 30 days. if there is one, associate it with the flow action execution
	for _, senderId := range senderIds {
		for _, social := range *socials {
			if strings.Contains(social.Url, "linkedin.com") {
				requestSent, err := s.services.Neo4jRepositories.LinkedinConnectionRequestReadRepository.GetPendingRequestByUserForSocialUrl(ctx, tx, tenant, senderId, social.Url)
				if err != nil {
					tracing.TraceErr(span, err)
					return err
				}

				if requestSent != nil {
					socialUrl = social.Url
					requestSentAlready = true
					break
				}
			}
		}

		//if there is a linkedin request sent already to one of the socials for the contact
		if requestSentAlready {
			span.LogFields(log.String("process", "linkedin request sent already by user "+senderId+" to: "+socialUrl))
			id, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowActionExecution)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			now := utils.Now()
			_, err = s.services.Neo4jRepositories.FlowActionExecutionWriteRepository.Merge(ctx, tx, &entity.FlowActionExecutionEntity{
				Id:              id,
				FlowId:          flowId,
				ActionId:        nextAction.Id,
				EntityId:        flowParticipant.EntityId,
				EntityType:      flowParticipant.EntityType,
				UserId:          &senderId,
				SocialUrl:       &socialUrl,
				ExecutedAt:      &now,
				ScheduledAt:     now,
				StatusUpdatedAt: now,
				Status:          entity.FlowActionExecutionStatusInProgress,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			_, err = s.upsertFlowExecutionSettings(ctx, tx, tenant, flowId, flowParticipant, nil, &senderId)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			return nil
		}
	}

	for _, social := range *socials {
		if strings.Contains(social.Url, "linkedin.com") {
			socialUrl = social.Url
			break
		}
	}

	if socialUrl == "" {
		return errors.New("No linkedin social found")
	}

	//if there is no linkedin request sent already to one of the socials for the contact
	span.LogFields(log.String("process", "no linkedin request sent already"))
	fastestUserId := ""
	fastestUserAt := time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, senderId := range senderIds {
		lastScheduledNode, err := s.services.Neo4jRepositories.LinkedinConnectionRequestReadRepository.GetLastScheduledForUser(ctx, tx, senderId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if lastScheduledNode == nil {
			fastestUserId = senderId
			break
		}

		lastScheduled := mapper.MapDbNodeToFlowActionExecutionEntity(lastScheduledNode)
		if lastScheduled.ScheduledAt.Before(fastestUserAt) {
			fastestUserId = senderId
			fastestUserAt = lastScheduled.ScheduledAt
		}
	}

	if fastestUserId == "" {
		return errors.New("No fastest user found")
	}

	flowExecutionSettings, err := s.upsertFlowExecutionSettings(ctx, tx, tenant, flowId, flowParticipant, nil, &fastestUserId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	workingSchedule, err := s.services.PostgresRepositories.UserWorkingScheduleRepository.GetForUser(ctx, tenant, fastestUserId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if len(workingSchedule) == 0 {
		return errors.New("User working schedule not found")
	}

	// 2. Schedule the email action
	actualScheduleAt, err := s.getFirstAvailableSlotForLinkedinConnection(ctx, tx, fastestUserId, scheduleAt, workingSchedule)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowActionExecutionId, err := s.storeNextActionExecutionEntity(ctx, tx, flowId, nextAction.Id, flowParticipant, flowExecutionSettings.Mailbox, actualScheduleAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	id, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelLinkedinConnectionRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	err = s.services.Neo4jRepositories.LinkedinConnectionRequestWriteRepository.Save(ctx, tx, &entity.LinkedinConnectionRequest{
		Id:           id,
		ProducerId:   flowActionExecutionId,
		ProducerType: model.NodeLabelFlowActionExecution,
		SocialUrl:    socialUrl,
		UserId:       fastestUserId,
		ScheduledAt:  *actualScheduleAt,
		Status:       entity.LinkedinConnectionRequestStatusPending,
	})
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

func (s *flowExecutionService) upsertFlowExecutionSettings(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, flowId string, participant *entity.FlowParticipantEntity, mailbox, userId *string) (*entity.FlowExecutionSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.upsertFlowExecutionSettings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	flowExecutionSettings, err := s.GetFlowExecutionSettingsForEntity(ctx, tx, flowId, participant.EntityId, participant.EntityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if flowExecutionSettings == nil {
		id, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowExecutionSettings)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		flowExecutionSettings = &entity.FlowExecutionSettingsEntity{
			Id:         id,
			FlowId:     flowId,
			EntityId:   participant.EntityId,
			EntityType: participant.EntityType,
			Mailbox:    mailbox,
			UserId:     userId,
		}
	}

	if mailbox != nil {
		flowExecutionSettings.Mailbox = mailbox
	}
	if userId != nil {
		flowExecutionSettings.UserId = userId
	}

	node, err := s.services.Neo4jRepositories.FlowExecutionSettingsWriteRepository.Merge(ctx, tx, flowExecutionSettings)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowExecutionSettingsEntity(node), nil
}

func (s *flowExecutionService) getFirstAvailableSlotForMailbox(ctx context.Context, tx *neo4j.ManagedTransaction, mailbox string, scheduleAt time.Time, workingSchedule []*postgresentity.UserWorkingSchedule) (*time.Time, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.getFirstAvailableSlotForMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	mailboxEntity, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, mailbox)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

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
	possibleScheduledAt = adjustToWorkingTimeWithRandom(maxTime(possibleScheduledAt, utils.Now()), workingSchedule, mailboxEntity.MinMinutesBetweenEmails, mailboxEntity.MaxMinutesBetweenEmails)

	//Add random seconds and miliseconds to not have 00:00:00 as the scheduled time
	randomSeconds := time.Duration(utils.GenerateRandomInt(0, 60)) * time.Second
	randomMiliseconds := time.Duration(utils.GenerateRandomInt(0, 1000)) * time.Millisecond
	randomMicroseconds := time.Duration(utils.GenerateRandomInt(0, 1000)) * time.Microsecond
	possibleScheduledAt = possibleScheduledAt.Add(randomSeconds).Add(randomMiliseconds).Add(randomMicroseconds)

	return &possibleScheduledAt, nil
}

func (s *flowExecutionService) getFirstAvailableSlotForLinkedinConnection(ctx context.Context, tx *neo4j.ManagedTransaction, userId string, scheduleAt time.Time, workingSchedule []*postgresentity.UserWorkingSchedule) (*time.Time, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.getFirstAvailableSlotForLinkedinConnection")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Get the last scheduled execution for this mailbox
	lastScheduledNode, err := s.services.Neo4jRepositories.LinkedinConnectionRequestReadRepository.GetLastScheduledForUser(ctx, tx, userId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	possibleScheduledAt := scheduleAt

	if lastScheduledNode != nil {
		lastScheduledExecution := mapper.MapDbNodeToFlowActionExecutionEntity(lastScheduledNode)
		possibleScheduledAt = maxTime(possibleScheduledAt, lastScheduledExecution.ScheduledAt)
	}

	//check the number of requests scheduled for the day
	for {
		requestsScheduledInDay, err := s.services.Neo4jRepositories.LinkedinConnectionRequestReadRepository.CountRequestsPerUserPerDay(ctx, tx, userId, utils.StartOfDayInUTC(possibleScheduledAt), utils.EndOfDayInUTC(possibleScheduledAt))
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if requestsScheduledInDay >= 20 {
			possibleScheduledAt = possibleScheduledAt.AddDate(0, 0, 1)
			possibleScheduledAt = time.Date(possibleScheduledAt.Year(), possibleScheduledAt.Month(), possibleScheduledAt.Day(), 0, 0, 0, 0, time.UTC)
			continue
		} else {
			break
		}
	}

	// Ensure possibleScheduledAt is not in the past and within working hours
	possibleScheduledAt = adjustToWorkingTimeWithRandom(maxTime(possibleScheduledAt, utils.Now()), workingSchedule, 5, 10)

	//Add random seconds and miliseconds to not have 00:00:00 as the scheduled time
	randomSeconds := time.Duration(utils.GenerateRandomInt(0, 60)) * time.Second
	randomMiliseconds := time.Duration(utils.GenerateRandomInt(0, 1000)) * time.Millisecond
	randomMicroseconds := time.Duration(utils.GenerateRandomInt(0, 1000)) * time.Microsecond
	possibleScheduledAt = possibleScheduledAt.Add(randomSeconds).Add(randomMiliseconds).Add(randomMicroseconds)

	return &possibleScheduledAt, nil
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func (s *flowExecutionService) storeNextActionExecutionEntity(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, actionId string, flowParticipant *entity.FlowParticipantEntity, mailbox *string, executionTime *time.Time) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.storeNextActionExecutionEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	id, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowActionExecution)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
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
		return "", err
	}

	return id, nil
}

func (s *flowExecutionService) ProcessActionExecution(ctx context.Context, scheduledActionExecution *entity.FlowActionExecutionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.ProcessActionExecution")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogFields(log.Object("scheduledActionExecution", scheduledActionExecution))

	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	participantUpdated := false

	var emailMessage *postgresentity.EmailMessage
	var sendLinkedInConnection *postgresentity.BrowserAutomationsRun

	var currentAction *entity.FlowActionEntity

	participant, err := s.services.FlowService.FlowParticipantByEntity(ctx, scheduledActionExecution.FlowId, scheduledActionExecution.EntityId, scheduledActionExecution.EntityType)
	if err != nil {
		return errors.Wrap(err, "failed to get flow participant by entity")
	}

	if participant == nil {
		return errors.New("participant not found")
	}

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
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

		participantUpdated, err = s.UpdateParticipantFlowRequirements(ctx, &tx, participant, flowRequirements)
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

				mailbox, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, *scheduledActionExecution.Mailbox)
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
		} else if currentAction.Data.Action == entity.FlowActionTypeLinkedinConnectionRequest {
			if scheduledActionExecution.SocialUrl == nil {
				return nil, errors.New("social url not found")
			}

			if scheduledActionExecution.UserId == nil {
				return nil, errors.New("user id not found")
			}

			linkedinTokens, err := s.services.PostgresRepositories.BrowserConfigRepository.GetForUser(ctx, *scheduledActionExecution.UserId)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get linkedin tokens for user")
			}

			if linkedinTokens == nil {
				return nil, errors.New("linkedin tokens not found")
			}

			payload := []map[string]interface{}{{"url": *scheduledActionExecution.SocialUrl}}
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				return nil, errors.Wrap(err, "failed to marshal payload")
			}
			sendLinkedInConnection = &postgresentity.BrowserAutomationsRun{
				BrowserConfigId: linkedinTokens.Id,
				UserId:          linkedinTokens.UserId,
				Tenant:          linkedinTokens.Tenant,
				Type:            "SEND_CONNECTION_REQUEST",
				Status:          "SCHEDULED",
				Payload:         string(payloadBytes),
			}
		}

		scheduledActionExecution.ExecutedAt = utils.TimePtr(utils.Now())
		scheduledActionExecution.StatusUpdatedAt = utils.Now()
		scheduledActionExecution.Status = entity.FlowActionExecutionStatusSuccess

		_, err = s.services.Neo4jRepositories.FlowActionExecutionWriteRepository.Merge(ctx, &tx, scheduledActionExecution)
		if err != nil {
			return nil, errors.Wrap(err, "failed to merge flow action execution")
		}

		err = s.services.FlowExecutionService.ScheduleFlow(ctx, &tx, scheduledActionExecution.FlowId, participant)
		if err != nil {
			return nil, errors.Wrap(err, "failed to schedule flow")
		}

		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if emailMessage != nil {
		//store in PG after the neo4j transaction is committed
		err = s.services.PostgresRepositories.EmailMessageRepository.Store(ctx, tenant, emailMessage)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "failed to store email message")
		}
	}

	if sendLinkedInConnection != nil {
		//store in PG after the neo4j transaction is committed
		err = s.services.PostgresRepositories.BrowserAutomationRunRepository.Add(ctx, sendLinkedInConnection)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "failed to store email message")
		}
	}

	if emailMessage != nil || sendLinkedInConnection != nil {
		_, err = s.services.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgresentity.BillableEventFlowActionExecuted,
			postgresrepository.BillableEventDetails{
				Subtype:       string(currentAction.Data.Action),
				ReferenceData: fmt.Sprintf("FlowActionExecutionId: %s", scheduledActionExecution.Id),
			})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to store billable event"))
		}
	}

	if participantUpdated {
		s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, participant.Id, model.FLOW_PARTICIPANT, utils.NewEventCompletedDetails().WithUpdate())
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
func adjustToWorkingTimeWithRandom(t time.Time, schedules []*postgresentity.UserWorkingSchedule, minRandom, maxRandom int) time.Time {
	for {
		// Get working hours for the current day
		start, end := getWorkingHoursForDay(t, schedules)

		randomMinutes := time.Duration(utils.GenerateRandomInt(minRandom, maxRandom)) * time.Minute
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
