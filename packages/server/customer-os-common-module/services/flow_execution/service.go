package flow_execution

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type flowExecutionService struct {
	neo4j    *neo4j_repository.Repositories
	postgres *postgres_repository.Repositories
	events   *events.EventsService
	email    interfaces.EmailService
	flow     interfaces.FlowService
	org      interfaces.OrganizationService
	social   interfaces.SocialService
}

func NewFlowExecutionService(neo4j *neo4j_repository.Repositories, postgres *postgres_repository.Repositories, events *events.EventsService, email interfaces.EmailService, flow interfaces.FlowService, org interfaces.OrganizationService, social interfaces.SocialService) interfaces.FlowExecutionService {
	return &flowExecutionService{
		neo4j:    neo4j,
		postgres: postgres,
		events:   events,
		email:    email,
		flow:     flow,
		org:      org,
		social:   social,
	}
}

func (s *flowExecutionService) SetEmailService(email interfaces.EmailService) {
	s.email = email
}

func (s *flowExecutionService) SetFlowService(flow interfaces.FlowService) {
	s.flow = flow
}

func (s *flowExecutionService) SetOrganizationService(org interfaces.OrganizationService) {
	s.org = org
}

func (s *flowExecutionService) SetSocialService(social interfaces.SocialService) {
	s.social = social
}

func (s *flowExecutionService) IsInitialized() bool {
	return utils.IsInitialized(s)
}

func (s *flowExecutionService) GetFlowActionExecutionById(ctx context.Context, flowActionExecution string) (*neo4j_entity.FlowActionExecutionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetFlowActionExecutionById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowActionExecution", flowActionExecution))

	node, err := s.neo4j.FlowActionExecutionReadRepository.GetById(ctx, flowActionExecution)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowActionExecutionEntity(node), nil
}

func (s *flowExecutionService) GetFlowExecutionSettingsForEntity(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId string, entityType model.EntityType) (*neo4j_entity.FlowExecutionSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetFlowExecutionSettingsForEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	node, err := s.neo4j.FlowExecutionSettingsReadRepository.GetForEntity(ctx, tx, flowId, entityId, entityType.String())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowExecutionSettingsEntity(node), nil
}

func (s *flowExecutionService) GetFlowRequirements(ctx context.Context, flowId string) (*interfaces.FlowComputeParticipantsRequirementsInput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetFlowRequirements")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowId", flowId))

	requirements := interfaces.FlowComputeParticipantsRequirementsInput{
		PrimaryEmailRequired:      false,
		LinkedInSocialUrlRequired: false,
	}

	flowActions, err := s.flow.FlowActionGetList(ctx, []string{flowId})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	for _, v := range *flowActions {
		if v.Data.Action == neo4j_entity.FlowActionTypeEmailNew || v.Data.Action == neo4j_entity.FlowActionTypeEmailReply {
			requirements.PrimaryEmailRequired = true
		}
		if v.Data.Action == neo4j_entity.FlowActionTypeLinkedinConnectionRequest || v.Data.Action == neo4j_entity.FlowActionTypeLinkedinMessage {
			requirements.LinkedInSocialUrlRequired = true
		}
	}

	return &requirements, nil
}

func (s *flowExecutionService) GetFlowActionExecutionsForParticipantWithActionType(ctx context.Context, entityId string, entityType model.EntityType, actionType neo4j_entity.FlowActionType) ([]*neo4j_entity.FlowActionExecutionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetFlowActionExecutionsForParticipantWithActionType")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	nodes, err := s.neo4j.FlowActionExecutionReadRepository.GetForEntityWithActionType(ctx, entityId, entityType, actionType)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make([]*neo4j_entity.FlowActionExecutionEntity, 0)
	for _, node := range nodes {
		entities = append(entities, mapper.MapDbNodeToFlowActionExecutionEntity(node))
	}

	return entities, nil
}

func (s *flowExecutionService) UpdateAllParticipantsFlowRequirements(ctx context.Context, flowId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.UpdateAllParticipantsFlowRequirements")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowId", flowId))

	flowParticipants, err := s.flow.FlowParticipantGetList(ctx, []string{flowId})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowRequirements, err := s.GetFlowRequirements(ctx, flowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, v := range *flowParticipants {
		err := s.UpdateParticipantFlowRequirements(ctx, nil, &v, flowRequirements)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	return nil
}

func (s *flowExecutionService) UpdateParticipantFlowRequirements(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, participant *neo4j_entity.FlowParticipantEntity, requirements *interfaces.FlowComputeParticipantsRequirementsInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.UpdateParticipantFlowRequirements")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	if requirements != nil {
		tracing.LogObjectAsJson(span, "requirements", requirements)
	} else {
		span.LogFields(log.String("requirements", "nil"))
	}

	tenant := common.GetTenantFromContext(ctx)

	if participant.Status == neo4j_entity.FlowParticipantStatusCompleted || participant.Status == neo4j_entity.FlowParticipantStatusGoalAchieved {
		return nil
	}

	status := neo4j_entity.FlowParticipantStatusReady
	newUnmeets := make([]neo4j_entity.FlowParticipantRequirementsUnmeet, 0)

	if requirements.PrimaryEmailRequired {
		// identify the primary email
		primaryEmail, err := s.email.GetPrimaryEmailForEntityId(ctx, participant.EntityType, participant.EntityId)
		if err != nil {
			return errors.Wrap(err, "failed to get primary email for entity id")
		}

		if primaryEmail == nil {
			status = neo4j_entity.FlowParticipantStatusOnHold
			newUnmeets = append(newUnmeets, neo4j_entity.FlowParticipantRequirementsUnmeetMissingPrimaryEmail)
		}
	}

	if requirements.LinkedInSocialUrlRequired {
		socials, err := s.social.GetAllForEntities(ctx, tenant, participant.EntityType, []string{participant.EntityId})
		if err != nil {
			return errors.Wrap(err, "failed to get socials for entities")
		}
		found := false
		for _, social := range *socials {
			if strings.Contains(social.Url, "linkedin.com") {
				found = true
				break
			}
		}

		if !found {
			status = neo4j_entity.FlowParticipantStatusOnHold
			newUnmeets = append(newUnmeets, neo4j_entity.FlowParticipantRequirementsUnmeetMissingLinkedinUrl)
		}
	}

	// no new meets, return
	unmeetsChanged := false
	for _, v := range newUnmeets {
		if !utils.ContainsElement(participant.RequirementsUnmeet, v) {
			unmeetsChanged = true
			break
		}
	}

	if !unmeetsChanged && participant.Status == status {
		return nil
	}

	participant.Status = status
	participant.RequirementsUnmeet = newUnmeets

	_, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		_, err := s.neo4j.FlowParticipantWriteRepository.Merge(ctx, txWithPostCommit.Tx, participant)
		if err != nil {
			return nil, errors.Wrap(err, "failed to merge participant")
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.UpdateParticipantFlowRequirements.postCommit")
			defer span.Finish()

			s.events.Publisher.PublishEventCompleted(ctx, tenant, participant.Id, model.FLOW_PARTICIPANT, utils.NewEventCompletedDetails().WithUpdate())

			return nil
		})

		return nil, nil
	})
	if err != nil {
		err := errors.Wrap(err, "failed to execute write in transaction with post commit actions")
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(log.String("result.participant.status", string(participant.Status)))
	participant.Status = status

	return nil
}

func (s *flowExecutionService) ScheduleFlow(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, flowId string, flowParticipant *neo4j_entity.FlowParticipantEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.ScheduleFlow")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	now := utils.Now()

	_, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (interface{}, error) {
		flow, err := s.flow.FlowGetByParticipantId(ctx, txWithPostCommit.Tx, flowParticipant.Id)
		if err != nil {
			return nil, err
		}

		if flow.Status != neo4j_entity.FlowStatusOn {
			return nil, nil
		}

		// check if the participant meets flow requirements
		flowRequirements, err := s.GetFlowRequirements(ctx, flowId)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get flow requirements")
		}

		err = s.UpdateParticipantFlowRequirements(ctx, txWithPostCommit, flowParticipant, flowRequirements)
		if err != nil {
			return nil, errors.Wrap(err, "failed to update participant flow requirements")
		}

		if flowParticipant.Status != neo4j_entity.FlowParticipantStatusReady {
			return nil, nil
		}

		flowExecutions, err := s.GetFlowActionExecutionsForParticipant(ctx, txWithPostCommit.Tx, flowId, flowParticipant.EntityId, flowParticipant.EntityType)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		span.LogFields(log.Int("flowExecutionsCount", len(flowExecutions)))

		if len(flowExecutions) == 0 {
			startAction, err := s.flow.FlowActionGetStart(ctx, flowId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			nextActions, err := s.flow.FlowActionGetNext(ctx, startAction.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			for _, nextAction := range nextActions {

				scheduleAt := now.Add(time.Duration(nextAction.Data.WaitBefore) * time.Minute)

				err := s.scheduleNextAction(ctx, txWithPostCommit, flowId, flowParticipant, scheduleAt, *nextAction)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
			}

		} else {
			lastActionExecution := flowExecutions[len(flowExecutions)-1]
			lastActionExecutedAt := lastActionExecution.ScheduledAt

			span.LogFields(log.String("lastActionExecution.Status", string(lastActionExecution.Status)))

			if lastActionExecution.Status != neo4j_entity.FlowActionExecutionStatusSuccess {
				return nil, nil
			}

			lastAction, err := s.flow.FlowActionGetById(ctx, lastActionExecution.ActionId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			nextActions, err := s.flow.FlowActionGetNext(ctx, lastAction.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			for _, nextAction := range nextActions {

				// marking the flow as completed if the next action is FLOW_END
				if nextAction.Data.Action == neo4j_entity.FlowActionTypeFlowEnd {
					flowParticipant.Status = neo4j_entity.FlowParticipantStatusCompleted

					_, err = s.neo4j.FlowParticipantWriteRepository.Merge(ctx, txWithPostCommit.Tx, flowParticipant)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}

					return nil, nil
				}

				scheduleAt := lastActionExecutedAt.Add(time.Duration(nextAction.Data.WaitBefore) * time.Minute)

				err := s.scheduleNextAction(ctx, txWithPostCommit, flowId, flowParticipant, scheduleAt, *nextAction)
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

func (s *flowExecutionService) GetFlowActionExecutionsForParticipants(ctx context.Context, flowParticipantIds []string) (*neo4j_entity.FlowActionExecutionEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetFlowActionExecutionsForParticipants")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// get executions for contact
	nodes, err := s.neo4j.FlowActionExecutionReadRepository.GetForParticipants(ctx, flowParticipantIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4j_entity.FlowActionExecutionEntities, 0)
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowActionExecutionEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowExecutionService) GetFlowActionExecutionsForParticipant(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId string, entityType model.EntityType) ([]*neo4j_entity.FlowActionExecutionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.GetFlowActionExecutionsForParticipant")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// get executions for contact
	nodes, err := s.neo4j.FlowActionExecutionReadRepository.GetForEntity(ctx, tx, flowId, entityId, entityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make([]*neo4j_entity.FlowActionExecutionEntity, 0)
	for _, node := range nodes {
		entities = append(entities, mapper.MapDbNodeToFlowActionExecutionEntity(node))
	}

	return entities, nil
}

func (s *flowExecutionService) scheduleNextAction(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, flowId string, flowParticipant *neo4j_entity.FlowParticipantEntity, scheduleAt time.Time, nextAction neo4j_entity.FlowActionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.scheduleNextAction")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if flowParticipant.Status != neo4j_entity.FlowParticipantStatusReady {
		// todo return business error
		return nil
	}

	switch nextAction.Data.Action {
	case neo4j_entity.FlowActionTypeEmailNew, neo4j_entity.FlowActionTypeEmailReply:
		return s.scheduleEmailAction(ctx, txWithPostCommit, flowId, flowParticipant, scheduleAt, nextAction)
	case neo4j_entity.FlowActionTypeLinkedinConnectionRequest:
		return s.scheduleSendLinkedInConnection(ctx, txWithPostCommit, flowId, flowParticipant, scheduleAt, nextAction)
	default:
		tracing.TraceErr(span, fmt.Errorf("Unsupported action type %s", nextAction.Data.Action))
		return errors.New("Unsupported action type")
	}
}

func (s *flowExecutionService) scheduleEmailAction(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, flowId string, flowParticipant *neo4j_entity.FlowParticipantEntity, scheduleAt time.Time, nextAction neo4j_entity.FlowActionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.scheduleEmailAction")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	// 1. Get the mailbox for contact or associate the best available mailbox
	flowExecutionSettings, err := s.GetFlowExecutionSettingsForEntity(ctx, txWithPostCommit.Tx, flowId, flowParticipant.EntityId, flowParticipant.EntityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowExecutionSettings == nil || flowExecutionSettings.Mailbox == nil {
		// compute the best available mailbox and associate

		// 1. get all available mailboxes
		// 2. select the mailbox with the fastest response time
		flowSenders, err := s.flow.FlowSenderGetList(ctx, []string{flowId})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		mailboxesScheduledAt := make(map[string]*time.Time)
		mailboxesScheduledAt[""] = utils.TimePtr(time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC))

		for _, flowActionSender := range *flowSenders {
			if flowActionSender.UserId == nil {
				continue
			}

			mailboxes, err := s.postgres.TenantSettingsMailboxRepository.GetAllByUserId(ctx, *flowActionSender.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			for _, mailbox := range mailboxes {
				scheduledAt, err := s.neo4j.FlowActionExecutionReadRepository.GetFirstSlotForMailbox(ctx, txWithPostCommit.Tx, mailbox.MailboxUsername)
				if err != nil {
					tracing.TraceErr(span, err)
					return err
				}

				mailboxesScheduledAt[mailbox.MailboxUsername] = scheduledAt
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

		mailbox, err := s.postgres.TenantSettingsMailboxRepository.GetByMailbox(ctx, fastestMailbox)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		flowExecutionSettings, err = s.upsertFlowExecutionSettings(ctx, txWithPostCommit.Tx, tenant, flowId, flowParticipant, &fastestMailbox, &mailbox.UserId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	workingSchedule, err := s.postgres.UserWorkingScheduleRepository.GetForUser(ctx, tenant, *flowExecutionSettings.UserId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if len(workingSchedule) == 0 {
		return errors.New("User working schedule not found")
	}

	// 2. Schedule the email action
	actualScheduleAt, err := s.getFirstAvailableSlotForMailbox(ctx, txWithPostCommit.Tx, *flowExecutionSettings.Mailbox, scheduleAt, workingSchedule)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = s.storeNextActionExecutionEntity(ctx, txWithPostCommit.Tx, &neo4j_entity.FlowActionExecutionEntity{
		FlowId:          flowId,
		ActionId:        nextAction.Id,
		ParticipantId:   flowParticipant.Id,
		EntityId:        flowParticipant.EntityId,
		EntityType:      flowParticipant.EntityType,
		Mailbox:         flowExecutionSettings.Mailbox,
		UserId:          flowExecutionSettings.UserId,
		ScheduledAt:     *actualScheduleAt,
		StatusUpdatedAt: utils.Now(),
		Status:          neo4j_entity.FlowActionExecutionStatusScheduled,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowParticipant.Status = neo4j_entity.FlowParticipantStatusScheduled
	_, err = s.neo4j.FlowParticipantWriteRepository.Merge(ctx, txWithPostCommit.Tx, flowParticipant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *flowExecutionService) scheduleSendLinkedInConnection(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, flowId string, flowParticipant *neo4j_entity.FlowParticipantEntity, scheduleAt time.Time, nextAction neo4j_entity.FlowActionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.scheduleSendLinkedInConnection")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	if flowParticipant.EntityType != model.CONTACT {
		return errors.New("Only contacts are supported for LinkedIn connection requests")
	}

	flowSenders, err := s.flow.FlowSenderGetList(ctx, []string{flowId})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	senderIds := make([]string, 0)
	for _, flowSender := range *flowSenders {
		if flowSender.UserId == nil {
			continue
		}

		activeLinkedinToken, err := s.postgres.BrowserConfigRepository.GetForUser(ctx, *flowSender.UserId)
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
		isLinkedWith, err := s.neo4j.CommonReadRepository.IsLinkedWith(ctx, tenant, flowParticipant.EntityId, model.CONTACT, model.CONNECTED_WITH.String(), senderId, model.USER)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "CommonReadRepository.IsLinkedWith"))
			return err
		}

		if isLinkedWith {
			span.LogFields(log.String("process", senderId+" is already connected with the contact"))
			id, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowActionExecution)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			now := utils.Now()
			_, err = s.neo4j.FlowActionExecutionWriteRepository.Merge(ctx, txWithPostCommit.Tx, &neo4j_entity.FlowActionExecutionEntity{
				Id:              id,
				FlowId:          flowId,
				ActionId:        nextAction.Id,
				ParticipantId:   flowParticipant.Id,
				EntityId:        flowParticipant.EntityId,
				EntityType:      flowParticipant.EntityType,
				UserId:          &senderId,
				ExecutedAt:      &now,
				ScheduledAt:     now,
				StatusUpdatedAt: now,
				Status:          neo4j_entity.FlowActionExecutionStatusSkipped,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			_, err = s.upsertFlowExecutionSettings(ctx, txWithPostCommit.Tx, tenant, flowId, flowParticipant, nil, &senderId)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			err = s.ScheduleFlow(ctx, txWithPostCommit, flowId, flowParticipant)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			return nil
		}
	}

	socials, err := s.social.GetAllForEntities(ctx, tenant, model.CONTACT, []string{flowParticipant.EntityId})
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
				requestSent, err := s.neo4j.LinkedinConnectionRequestReadRepository.GetPendingRequestByUserForSocialUrl(ctx, txWithPostCommit.Tx, tenant, senderId, social.Url)
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

		// if there is a linkedin request sent already to one of the socials for the contact
		if requestSentAlready {
			span.LogFields(log.String("process", "linkedin request sent already by user "+senderId+" to: "+socialUrl))
			id, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowActionExecution)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			now := utils.Now()
			_, err = s.neo4j.FlowActionExecutionWriteRepository.Merge(ctx, txWithPostCommit.Tx, &neo4j_entity.FlowActionExecutionEntity{
				Id:              id,
				FlowId:          flowId,
				ActionId:        nextAction.Id,
				ParticipantId:   flowParticipant.Id,
				EntityId:        flowParticipant.EntityId,
				EntityType:      flowParticipant.EntityType,
				UserId:          &senderId,
				SocialUrl:       &socialUrl,
				ExecutedAt:      &now,
				ScheduledAt:     now,
				StatusUpdatedAt: now,
				Status:          neo4j_entity.FlowActionExecutionStatusInProgress,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			_, err = s.upsertFlowExecutionSettings(ctx, txWithPostCommit.Tx, tenant, flowId, flowParticipant, nil, &senderId)
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

	// if there is no linkedin request sent already to one of the socials for the contact
	span.LogFields(log.String("process", "no linkedin request sent already"))
	fastestUserId := ""
	fastestUserAt := time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, senderId := range senderIds {
		lastScheduledNode, err := s.neo4j.LinkedinConnectionRequestReadRepository.GetLastScheduledForUser(ctx, txWithPostCommit.Tx, senderId)
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

	_, err = s.upsertFlowExecutionSettings(ctx, txWithPostCommit.Tx, tenant, flowId, flowParticipant, nil, &fastestUserId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	workingSchedule, err := s.postgres.UserWorkingScheduleRepository.GetForUser(ctx, tenant, fastestUserId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if len(workingSchedule) == 0 {
		return errors.New("User working schedule not found")
	}

	// 2. Schedule the email action
	actualScheduleAt, err := s.getFirstAvailableSlotForLinkedinConnection(ctx, txWithPostCommit.Tx, fastestUserId, scheduleAt, workingSchedule)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowActionExecutionId, err := s.storeNextActionExecutionEntity(ctx, txWithPostCommit.Tx, &neo4j_entity.FlowActionExecutionEntity{
		FlowId:          flowId,
		ActionId:        nextAction.Id,
		ParticipantId:   flowParticipant.Id,
		EntityId:        flowParticipant.EntityId,
		EntityType:      flowParticipant.EntityType,
		UserId:          &fastestUserId,
		SocialUrl:       &socialUrl,
		ScheduledAt:     *actualScheduleAt,
		StatusUpdatedAt: utils.Now(),
		Status:          neo4j_entity.FlowActionExecutionStatusScheduled,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	id, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelLinkedinConnectionRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	err = s.neo4j.LinkedinConnectionRequestWriteRepository.Save(ctx, txWithPostCommit.Tx, &neo4j_entity.LinkedinConnectionRequest{
		Id:           id,
		ProducerId:   flowActionExecutionId,
		ProducerType: model.NodeLabelFlowActionExecution,
		SocialUrl:    socialUrl,
		UserId:       fastestUserId,
		ScheduledAt:  *actualScheduleAt,
		Status:       neo4j_entity.LinkedinConnectionRequestStatusPending,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	flowParticipant.Status = neo4j_entity.FlowParticipantStatusScheduled
	_, err = s.neo4j.FlowParticipantWriteRepository.Merge(ctx, txWithPostCommit.Tx, flowParticipant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *flowExecutionService) upsertFlowExecutionSettings(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, flowId string, participant *neo4j_entity.FlowParticipantEntity, mailbox, userId *string) (*neo4j_entity.FlowExecutionSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.upsertFlowExecutionSettings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	flowExecutionSettings, err := s.GetFlowExecutionSettingsForEntity(ctx, tx, flowId, participant.EntityId, participant.EntityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if flowExecutionSettings == nil {
		id, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowExecutionSettings)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		flowExecutionSettings = &neo4j_entity.FlowExecutionSettingsEntity{
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

	node, err := s.neo4j.FlowExecutionSettingsWriteRepository.Merge(ctx, tx, flowExecutionSettings)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowExecutionSettingsEntity(node), nil
}

func (s *flowExecutionService) getFirstAvailableSlotForMailbox(ctx context.Context, tx *neo4j.ManagedTransaction, mailbox string, scheduleAt time.Time, workingSchedule []*postgres_entity.UserWorkingSchedule) (*time.Time, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.getFirstAvailableSlotForMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	mailboxEntity, err := s.postgres.TenantSettingsMailboxRepository.GetByMailbox(ctx, mailbox)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Get the last scheduled execution for this mailbox
	lastScheduledExecutionNode, err := s.neo4j.FlowActionExecutionReadRepository.GetLastScheduledForMailbox(ctx, tx, mailbox)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	possibleScheduledAt := scheduleAt

	if lastScheduledExecutionNode != nil {
		lastScheduledExecution := mapper.MapDbNodeToFlowActionExecutionEntity(lastScheduledExecutionNode)
		possibleScheduledAt = maxTime(possibleScheduledAt, lastScheduledExecution.ScheduledAt)
	}

	// check the number of emails scheduled for the day
	for {
		emailsScheduledInDay, err := s.neo4j.FlowActionExecutionReadRepository.CountEmailsPerMailboxPerDay(ctx, tx, mailbox, utils.StartOfDayInUTC(possibleScheduledAt), utils.EndOfDayInUTC(possibleScheduledAt))
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

	// Add random seconds and miliseconds to not have 00:00:00 as the scheduled time
	randomSeconds := time.Duration(utils.GenerateRandomInt(0, 60)) * time.Second
	randomMiliseconds := time.Duration(utils.GenerateRandomInt(0, 1000)) * time.Millisecond
	randomMicroseconds := time.Duration(utils.GenerateRandomInt(0, 1000)) * time.Microsecond
	possibleScheduledAt = possibleScheduledAt.Add(randomSeconds).Add(randomMiliseconds).Add(randomMicroseconds)

	return &possibleScheduledAt, nil
}

func (s *flowExecutionService) getFirstAvailableSlotForLinkedinConnection(ctx context.Context, tx *neo4j.ManagedTransaction, userId string, scheduleAt time.Time, workingSchedule []*postgres_entity.UserWorkingSchedule) (*time.Time, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.getFirstAvailableSlotForLinkedinConnection")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Get the last scheduled execution for this mailbox
	lastScheduledNode, err := s.neo4j.LinkedinConnectionRequestReadRepository.GetLastScheduledForUser(ctx, tx, userId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	possibleScheduledAt := scheduleAt

	if lastScheduledNode != nil {
		lastScheduledExecution := mapper.MapDbNodeToFlowActionExecutionEntity(lastScheduledNode)
		possibleScheduledAt = maxTime(possibleScheduledAt, lastScheduledExecution.ScheduledAt)
	}

	// check the number of requests scheduled for the day
	for {
		requestsScheduledInDay, err := s.neo4j.LinkedinConnectionRequestReadRepository.CountRequestsPerUserPerDay(ctx, tx, userId, utils.StartOfDayInUTC(possibleScheduledAt), utils.EndOfDayInUTC(possibleScheduledAt))
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

	// Add random seconds and miliseconds to not have 00:00:00 as the scheduled time
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

func (s *flowExecutionService) storeNextActionExecutionEntity(ctx context.Context, tx *neo4j.ManagedTransaction, input *neo4j_entity.FlowActionExecutionEntity) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.storeNextActionExecutionEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	id, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowActionExecution)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	input.Id = id

	_, err = s.neo4j.FlowActionExecutionWriteRepository.Merge(ctx, tx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return id, nil
}

func (s *flowExecutionService) ProcessActionExecution(ctx context.Context, scheduledActionExecution *neo4j_entity.FlowActionExecutionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.ProcessActionExecution")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tenant := common.GetTenantFromContext(ctx)

	span.LogFields(log.Object("scheduledActionExecution", scheduledActionExecution))

	var currentAction *neo4j_entity.FlowActionEntity

	participant, err := s.flow.FlowParticipantByEntity(ctx, scheduledActionExecution.FlowId, scheduledActionExecution.EntityId, scheduledActionExecution.EntityType)
	if err != nil {
		return errors.Wrap(err, "failed to get flow participant by entity")
	}

	if participant == nil {
		return errors.New("participant not found")
	}

	addBillableEvent := false

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, nil, func(txWithPostCommit *utils.TxWithPostCommit) (interface{}, error) {
		currentAction, err = s.flow.FlowActionGetById(ctx, scheduledActionExecution.ActionId)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get action by id")
		}

		if currentAction == nil {
			return nil, errors.New("action not found")
		}

		// check if the participant meets flow requirements
		flowRequirements, err := s.GetFlowRequirements(ctx, scheduledActionExecution.FlowId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = s.UpdateParticipantFlowRequirements(ctx, txWithPostCommit, participant, flowRequirements)
		if err != nil {
			return nil, errors.Wrap(err, "failed to update participant flow requirements")
		}

		// if participant is not ready, mark the action as business error
		if participant.Status != neo4j_entity.FlowParticipantStatusReady {

			scheduledActionExecution.StatusUpdatedAt = utils.Now()
			scheduledActionExecution.Status = neo4j_entity.FlowActionExecutionStatusBusinessError

			_, err = s.neo4j.FlowActionExecutionWriteRepository.Merge(ctx, txWithPostCommit.Tx, scheduledActionExecution)
			if err != nil {
				return nil, errors.Wrap(err, "failed to merge flow action execution")
			}

			return nil, nil
		}

		if currentAction.Data.Action == neo4j_entity.FlowActionTypeEmailNew || currentAction.Data.Action == neo4j_entity.FlowActionTypeEmailReply {

			// prevent duplicate emails for same action
			existingEmail, err := s.postgres.EmailMessageRepository.GetByProducer(ctx, tenant, scheduledActionExecution.Id, model.NodeLabelFlowActionExecution)
			if err != nil {
				// todo this is producing an error below when trying to insert the email. need to rewrite how we store the email
				return nil, errors.Wrap(err, "failed to get email by producer")
			}

			if existingEmail == nil {
				span.LogFields(log.Bool("process.existingEmail", true))

				mailbox, err := s.postgres.TenantSettingsMailboxRepository.GetByMailbox(ctx, *scheduledActionExecution.Mailbox)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to get mailbox by mailbox"))
					return nil, errors.Wrap(err, "failed to get mailbox by mailbox")
				}

				if mailbox == nil {
					tracing.TraceErr(span, errors.New("Mailbox not found in database"))
					return nil, errors.New("mailbox not found in database")
				}

				primaryEmail, err := s.email.GetPrimaryEmailForEntityId(ctx, participant.EntityType, participant.EntityId)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to get primary email for entity id"))
					return nil, errors.Wrap(err, "failed to get primary email for entity id")
				}

				if primaryEmail == nil {
					tracing.TraceErr(span, errors.New("Primary email not found"))
					return nil, errors.New("primary email not found")
				}

				toEmail := primaryEmail.RawEmail
				span.LogFields(log.String("process.toEmail", toEmail))

				subjectTemplate := utils.IfNotNilString(currentAction.Data.Subject)
				bodyTemplate := utils.IfNotNilString(currentAction.Data.BodyTemplate)

				span.LogFields(log.Bool("process.bodyTemplate.available", bodyTemplate != ""))
				span.LogFields(log.Bool("process.subjectTemplate.available", subjectTemplate != ""))

				span.LogFields(log.String("process.scheduledActionExecution.EntityType", scheduledActionExecution.EntityType.String()))
				if scheduledActionExecution.EntityType == model.CONTACT {
					contactNode, err := s.neo4j.ContactReadRepository.GetContact(ctx, tenant, scheduledActionExecution.EntityId)
					if err != nil {
						return nil, errors.Wrap(err, "failed to get contact")
					}

					contact := mapper.MapDbNodeToContactEntity(contactNode)

					firstName, lastName := contact.DeriveFirstAndLastNames()
					bodyTemplate = s.ReplacePlaceholder(bodyTemplate, "contact_first_name", firstName)
					bodyTemplate = s.ReplacePlaceholder(bodyTemplate, "contact_last_name", lastName)
					bodyTemplate = s.ReplacePlaceholder(bodyTemplate, "contact_email", toEmail)

					contactWithOrganizations, err := s.org.GetPrimaryOrganizationsWithJobRoleForContacts(ctx, []string{contact.Id})
					if err != nil {
						return nil, errors.Wrap(err, "failed to get latest organizations with job roles for contacts")
					}

					if len(*contactWithOrganizations) > 0 {
						contactWithOrganization := (*contactWithOrganizations)[0]
						bodyTemplate = s.ReplacePlaceholder(bodyTemplate, "organization_name", contactWithOrganization.Organization.Name)
						subjectTemplate = strings.ReplaceAll(subjectTemplate, "{{organization_name}}", contactWithOrganization.Organization.Name)
					} else {
						bodyTemplate = s.ReplacePlaceholder(bodyTemplate, "organization_name", "")
						subjectTemplate = strings.ReplaceAll(subjectTemplate, "{{organization_name}}", "")
					}

				}

				mailbox, err = s.postgres.TenantSettingsMailboxRepository.GetByMailbox(ctx, *scheduledActionExecution.Mailbox)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get mailbox by mailbox")
				}

				userNode, err := s.neo4j.UserReadRepository.GetUserById(ctx, tenant, mailbox.UserId)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get first user by email")
				}

				if userNode == nil {
					return nil, errors.New("user not found")
				}

				user := mapper.MapDbNodeToUserEntity(userNode)

				bodyTemplate = s.ReplacePlaceholder(bodyTemplate, "sender_first_name", user.FirstName)
				bodyTemplate = s.ReplacePlaceholder(bodyTemplate, "sender_last_name", user.LastName)

				addBillableEvent = true
				emailMessage := &postgres_entity.EmailMessage{
					Status:       postgres_entity.EmailMessageStatusScheduled,
					ProducerId:   scheduledActionExecution.Id,
					ProducerType: model.NodeLabelFlowActionExecution,
					FromName:     user.FirstName + " " + user.LastName,
					From:         *scheduledActionExecution.Mailbox,
					To:           []string{toEmail},
					Content:      bodyTemplate,
				}

				if currentAction.Data.Action == neo4j_entity.FlowActionTypeEmailNew {
					emailMessage.Subject = subjectTemplate
				}

				if currentAction.Data.Action == neo4j_entity.FlowActionTypeEmailReply {
					// walk back the flow and identify the previous email
					// get previous email execution from neo4j
					// get previous email from postgres
					// reply to the previous email

					parentEmailAction, err := s.getEmailActionToReply(ctx, scheduledActionExecution.ActionId)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, errors.Wrap(err, "failed to get email action to reply")
					}

					if parentEmailAction == nil {
						tracing.TraceErr(span, err)
						return nil, errors.New("no parent email action found")
					}

					parentEmailExecution, err := s.neo4j.FlowActionExecutionReadRepository.GetExecution(ctx, scheduledActionExecution.FlowId, parentEmailAction.Id, scheduledActionExecution.EntityId, scheduledActionExecution.EntityType)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, errors.Wrap(err, "failed to get execution")
					}

					if parentEmailExecution == nil {
						tracing.TraceErr(span, err)
						return nil, errors.New("no parent email execution found")
					}

					parentEmail := mapper.MapDbNodeToFlowActionExecutionEntity(parentEmailExecution)

					parentEmailSent, err := s.postgres.EmailMessageRepository.GetByProducer(ctx, tenant, parentEmail.Id, model.NodeLabelFlowActionExecution)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, errors.Wrap(err, "failed to get email by producer")
					}

					if parentEmailSent == nil {
						tracing.TraceErr(span, err)
						return nil, errors.New("no parent email sent found")
					}

					emailMessage.Subject = "Re: " + parentEmailSent.Subject
					emailMessage.ProviderInReplyTo = parentEmailSent.ProviderMessageId
					emailMessage.ProviderReferences = parentEmailSent.ProviderReferences + " " + parentEmailSent.ProviderMessageId
				}

				txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
					span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.ProcessActionExecution.PostCommitAction.StoreEmailMessage")
					defer span.Finish()

					err = s.postgres.EmailMessageRepository.Store(ctx, tenant, emailMessage)
					if err != nil {
						tracing.TraceErr(span, err)
						return errors.Wrap(err, "failed to store email message")
					}

					return nil
				})
			}
		} else if currentAction.Data.Action == neo4j_entity.FlowActionTypeLinkedinConnectionRequest {
			if scheduledActionExecution.SocialUrl == nil {
				tracing.TraceErr(span, err)
				return nil, errors.New("social url not found")
			}

			if scheduledActionExecution.UserId == nil {
				tracing.TraceErr(span, err)
				return nil, errors.New("user id not found")
			}

			linkedinTokens, err := s.postgres.BrowserConfigRepository.GetForUser(ctx, *scheduledActionExecution.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, errors.Wrap(err, "failed to get linkedin tokens for user")
			}

			if linkedinTokens == nil {
				tracing.TraceErr(span, err)
				return nil, errors.New("linkedin tokens not found")
			}

			payload := map[string]interface{}{"profileUrl": *scheduledActionExecution.SocialUrl}
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, errors.Wrap(err, "failed to marshal payload")
			}

			addBillableEvent = true
			sendLinkedInConnection := &postgres_entity.BrowserAutomationsRun{
				BrowserConfigId: linkedinTokens.Id,
				UserId:          linkedinTokens.UserId,
				Tenant:          linkedinTokens.Tenant,
				Type:            "SEND_CONNECTION_REQUEST",
				Status:          "SCHEDULED",
				Payload:         string(payloadBytes),
			}

			txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
				span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.ProcessActionExecution.PostCommitAction.SendLinkedInConnection")
				defer span.Finish()

				err = s.postgres.BrowserAutomationRunRepository.Add(ctx, sendLinkedInConnection)
				if err != nil {
					tracing.TraceErr(span, err)
					return errors.Wrap(err, "failed to store email message")
				}

				return nil
			})
		}

		scheduledActionExecution.ExecutedAt = utils.TimePtr(utils.Now())
		scheduledActionExecution.StatusUpdatedAt = utils.Now()
		scheduledActionExecution.Status = neo4j_entity.FlowActionExecutionStatusSuccess

		_, err = s.neo4j.FlowActionExecutionWriteRepository.Merge(ctx, txWithPostCommit.Tx, scheduledActionExecution)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, errors.Wrap(err, "failed to merge flow action execution")
		}

		err = s.ScheduleFlow(ctx, txWithPostCommit, scheduledActionExecution.FlowId, participant)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, errors.Wrap(err, "failed to schedule flow")
		}

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)

		participant.Status = neo4j_entity.FlowParticipantStatusError
		_, err = s.neo4j.FlowParticipantWriteRepository.Merge(ctx, nil, participant)

		s.events.Publisher.PublishEventCompleted(ctx, tenant, participant.Id, model.FLOW_PARTICIPANT, utils.NewEventCompletedDetails().WithUpdate())

		return err
	}

	span.LogFields(log.Bool("addBillableEvent", addBillableEvent))
	if addBillableEvent {
		_, err = s.postgres.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgres_entity.BillableEventFlowActionExecuted,
			postgres_repository.BillableEventDetails{
				Subtype:       string(currentAction.Data.Action),
				ReferenceData: fmt.Sprintf("FlowActionExecutionId: %s", scheduledActionExecution.Id),
			})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to store billable event"))
		}
	}

	return nil
}

func (s *flowExecutionService) getEmailActionToReply(ctx context.Context, actionId string) (*neo4j_entity.FlowActionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.getEmailActionToReply")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	var previous *neo4j_entity.FlowActionEntity

	previousNodes, err := s.neo4j.FlowActionReadRepository.GetPrevious(ctx, actionId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get previous nodes for action")
	}

	if previousNodes == nil || len(previousNodes) == 0 {
		return nil, errors.New("no previous nodes found for action")
	}

	if len(previousNodes) == 1 {
		previousEntity := mapper.MapDbNodeToFlowActionEntity(previousNodes[0])

		if previousEntity.Data.Action == neo4j_entity.FlowActionTypeFlowStart {
			return nil, errors.New("reached start of flow")
		} else if previousEntity.Data.Action == neo4j_entity.FlowActionTypeEmailNew {
			previous = previousEntity
		} else {
			return s.getEmailActionToReply(ctx, previousEntity.Id)
		}
	}

	// Handle multiple previous nodes
	for _, previousNode := range previousNodes {
		previousEntity := mapper.MapDbNodeToFlowActionEntity(previousNode)

		// Base case: If we find an email action, return it
		if previousEntity.Data.Action == neo4j_entity.FlowActionTypeEmailNew {
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

func (s *flowExecutionService) ReplacePlaceholder(input, variableName, value string) string {
	return strings.Replace(input, "{{"+variableName+"}}", value, -1)
}

func adjustToWorkingTimeWithRandom(t time.Time, schedules []*postgres_entity.UserWorkingSchedule, minRandom, maxRandom int) time.Time {
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
func getWorkingHoursForDay(day time.Time, schedules []*postgres_entity.UserWorkingSchedule) (time.Time, time.Time) {
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
