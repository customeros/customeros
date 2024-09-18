package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FlowExecutionService interface {
	ExecuteScheduledFlowActions()
}

type flowExecutionService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonService.Services
}

func NewFlowExecutionService(cfg *config.Config, log logger.Logger, commonServices *commonService.Services) FlowExecutionService {
	return &flowExecutionService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
	}
}

func (s *flowExecutionService) ExecuteScheduledFlowActions() {
	ctx, cancel := utils.GetContextWithTimeout(context.Background(), utils.HalfOfHourDuration)
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "FlowExecutionService.ExecuteScheduledFlowActions")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	actionsToExecute, err := s.commonServices.Neo4jRepositories.FlowActionExecutionReadRepository.GetScheduledBefore(ctx, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	span.LogFields(log.Int("actionsToExecute.count", len(actionsToExecute)))

	for _, actionExecutionNode := range actionsToExecute {
		// Do something with the action
		actionExecution := neo4jmapper.MapDbNodeToFlowActionExecutionEntity(actionExecutionNode)

		err := s.processActionExecution(ctx, actionExecution)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}
	}

}

func (s *flowExecutionService) processActionExecution(ctx context.Context, scheduledActionExecution *neo4jEntity.FlowActionExecutionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.processActionExecution")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("scheduledActionExecution", scheduledActionExecution))

	session := utils.NewNeo4jWriteSession(ctx, *s.commonServices.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		currentAction, err := s.commonServices.FlowService.FlowActionGetById(ctx, scheduledActionExecution.ActionId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if currentAction.ActionType == neo4jEntity.FlowActionTypeEmailNew {
			//TODO send email
		} else if currentAction.ActionType == neo4jEntity.FlowActionTypeEmailReply {
			//TODO reply to previous email
		}

		scheduledActionExecution.Status = neo4jEntity.FlowActionExecutionStatusSuccess

		_, err = s.commonServices.Neo4jRepositories.FlowActionExecutionWriteRepository.Merge(ctx, &tx, scheduledActionExecution)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = s.commonServices.FlowExecutionService.ScheduleFlowForContact(ctx, &tx, scheduledActionExecution.FlowId, scheduledActionExecution.ContactId)
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
