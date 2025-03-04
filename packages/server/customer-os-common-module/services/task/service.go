package task

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neoRepo "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type taskService struct {
	log                logger.Logger
	neo4j              *neoRepo.Repositories
	events             *events.EventsService
	userService        interfaces.UserService
	opportunityService interfaces.OpportunityService
}

func NewTaskService(log logger.Logger, neo4j *neoRepo.Repositories, events *events.EventsService, userSrv interfaces.UserService, opportunitySrv interfaces.OpportunityService) interfaces.TaskService {
	return &taskService{
		log:                log,
		neo4j:              neo4j,
		events:             events,
		userService:        userSrv,
		opportunityService: opportunitySrv,
	}
}

func (s *taskService) IsInitialized() bool {
	return utils.IsInitialized(s)
}

func (s *taskService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, taskFields data_fields.TaskFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TaskService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "taskFields", taskFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	taskId := ""

	if utils.IfNotNilString(id) == "" {
		createFlow = true
		span.LogKV("flow", "create")

		// prepare missing fields
		if utils.IfNotNilString(taskFields.Source) == "" {
			taskFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
		}
		if utils.IfNotNilString(taskFields.AppSource) == "" {
			taskFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}
		if utils.IfNotNilString(taskFields.CreatedByUserId) == "" {
			taskFields.CreatedByUserId = utils.StringPtr(common.GetUserIdFromContext(ctx))
		}
		if taskFields.Status == nil {
			taskFields.Status = utils.ToPtr(enum.TaskStatusTodo)
		}

		taskId, err = s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelTask)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	} else {
		span.LogKV("flow", "update")
		taskId = *id

		// validate task exists
		exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, taskId, model.NodeLabelTask)
		if err != nil || !exists {
			err = errors.New("task not found")
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, taskId)

	// validate assignees
	if taskFields.AssigneeUserIds != nil {
		for _, assigneeId := range *taskFields.AssigneeUserIds {
			// check if user exists
			exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, assigneeId, model.NodeLabelUser)
			if err != nil {
				tracing.TraceErr(span, err)
				return "", err
			}
			if !exists {
				err = errors.Errorf("asignee user %s not found", assigneeId)
				tracing.TraceErr(span, err)
				return "", err
			}
		}
	}
	// validate opportunities
	if taskFields.OpportunityIds != nil {
		for _, opportunityId := range *taskFields.OpportunityIds {
			// check if opportunity exists
			exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, opportunityId, model.NodeLabelOpportunity)
			if err != nil {
				tracing.TraceErr(span, err)
				return "", err
			}
			if !exists {
				err = errors.Errorf("opportunity %s not found", opportunityId)
				tracing.TraceErr(span, err)
				return "", err
			}
		}
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		if createFlow {
			err := s.neo4j.TaskWriteRepository.Create(ctx, txWithPostCommit.Tx, tenant, taskId, taskFields)
			if err != nil {
				s.log.Errorf("Error while saving task %s: %s", taskId, err.Error())
				return nil, err
			}
		} else {
			err := s.neo4j.TaskWriteRepository.Update(ctx, txWithPostCommit.Tx, tenant, taskId, taskFields)
			if err != nil {
				s.log.Errorf("Error while updating task %s: %s", taskId, err.Error())
				return nil, err
			}
		}

		if taskFields.OpportunityIds != nil {
			err := s.neo4j.TaskWriteRepository.SetOpportunities(ctx, txWithPostCommit.Tx, tenant, taskId, *taskFields.OpportunityIds)
			if err != nil {
				s.log.Errorf("Error while setting task opportunities %s: %s", taskId, err.Error())
				return nil, err
			}
		}
		if taskFields.AssigneeUserIds != nil {
			err := s.neo4j.TaskWriteRepository.SetUserAssignees(ctx, txWithPostCommit.Tx, tenant, taskId, *taskFields.AssigneeUserIds)
			if err != nil {
				s.log.Errorf("Error while setting task user assignees %s: %s", taskId, err.Error())
				return nil, err
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// send events
			if createFlow {
				err := s.events.Publisher.PublishFanoutEvent(ctx, taskId, model.TASK, dto.CreateTask{taskFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateTask"))
				}
				s.events.Publisher.PublishNotification(ctx, tenant, taskId, model.TASK, utils.NewEventCompletedDetails().WithCreate())
			} else {
				err := s.events.Publisher.PublishFanoutEvent(ctx, taskId, model.TASK, dto.UpdateTask{taskFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateTask"))
				}
				if common.GetTenantFromContext(ctx) != constants.AppSourceCustomerOsApi {
					s.events.Publisher.PublishNotification(ctx, tenant, taskId, model.TASK, utils.NewEventCompletedDetails().WithUpdate())
				}
			}
			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return taskId, nil
}

func (s *taskService) GetById(ctx context.Context, id string) (*neo4jentity.TaskEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TaskService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, id)

	dbNode, err := s.neo4j.TaskReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToTaskEntity(dbNode), nil
}

func (s *taskService) GetAllByIds(ctx context.Context, ids []string) (*neo4jentity.TaskEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TaskService.GetAllByIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	dbNodes, err := s.neo4j.TaskReadRepository.GetAllByIds(ctx, tenant, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var taskEntities neo4jentity.TaskEntities
	for _, dbNodePtr := range dbNodes {
		taskEntities = append(taskEntities, *neo4jmapper.MapDbNodeToTaskEntity(dbNodePtr))
	}
	return &taskEntities, nil
}
