package task

import (
	"context"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
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
	// validate opportunities

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
		//
		//	txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
		//		// send events
		//		if createFlow {
		//			if utils.IfNotNilString(issueFields.ReportedByOrganizationId) != "" {
		//				err := s.org.RequestRefreshLastTouchpoint(ctx, *issueFields.ReportedByOrganizationId)
		//				if err != nil {
		//					tracing.TraceErr(span, errors.Wrap(err, "unable to request refresh last touchpoint"))
		//				}
		//			}
		//			err := s.events.Publisher.PublishFanoutEvent(ctx, issueId, model.ISSUE, dto.CreateIssue{issueFields})
		//			if err != nil {
		//				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateIssue"))
		//			}
		//			s.events.Publisher.PublishNotification(ctx, tenant, issueId, model.ISSUE, utils.NewEventCompletedDetails().WithCreate())
		//		} else {
		//			err := s.events.Publisher.PublishFanoutEvent(ctx, issueId, model.ISSUE, dto.UpdateIssue{issueFields})
		//			if err != nil {
		//				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateIssue"))
		//			}
		//			if common.GetTenantFromContext(ctx) != constants.AppSourceCustomerOsApi {
		//				s.events.Publisher.PublishNotification(ctx, tenant, issueId, model.ISSUE, utils.NewEventCompletedDetails().WithUpdate())
		//			}
		//		}
		//		return nil
		//	})
		//
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return taskId, nil
}
