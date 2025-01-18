package service

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jEntity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
)

type FlowExecutionService interface {
	RampUpMailboxes()
	ExecuteScheduledFlowActions()
	ComputeFlowStatistics()
}

type flowExecutionService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonService.CommonServices
}

func NewFlowExecutionService(cfg *config.Config, log logger.Logger, commonServices *commonService.CommonServices) FlowExecutionService {
	return &flowExecutionService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
	}
}

func (s *flowExecutionService) RampUpMailboxes() {
	ctx, cancel := utils.GetContextWithTimeout(context.Background(), utils.HalfOfHourDuration)
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "FlowExecutionService.RampUpMailboxes")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	mailboxes, err := s.commonServices.PostgresRepositories.TenantSettingsMailboxRepository.GetForRampUp(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	span.LogFields(log.Int("mailboxes.count", len(mailboxes)))

	for _, mailbox := range mailboxes {
		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant: mailbox.Tenant,
		})

		err := s.rampUpMailbox(ctx, mailbox)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}
}

func (s *flowExecutionService) rampUpMailbox(ctx context.Context, mailbox *postgres_entity.TenantSettingsMailbox) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionService.rampUpMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	for {
		if mailbox.RampUpCurrent >= mailbox.RampUpMax {
			break
		}

		if mailbox.LastRampUpAt.After(utils.StartOfDayInUTC(utils.Now())) {
			break
		}

		mailbox.RampUpCurrent = mailbox.RampUpCurrent + mailbox.RampUpRate

		if mailbox.RampUpCurrent > mailbox.RampUpMax {
			mailbox.RampUpCurrent = mailbox.RampUpMax
		}

		mailbox.LastRampUpAt = mailbox.LastRampUpAt.AddDate(0, 0, 1)

		err := s.commonServices.PostgresRepositories.TenantSettingsMailboxRepository.Merge(ctx, nil, mailbox)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
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
		actionExecution := neo4jmapper.MapDbNodeToFlowActionExecutionEntity(actionExecutionNode)

		tenant := model.GetTenantFromLabels(actionExecutionNode.Labels, model.NodeLabelFlowActionExecution)
		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})

		err := s.commonServices.FlowExecutionService.ProcessActionExecution(ctx, actionExecution)
		if err != nil {
			tracing.TraceErr(span, err)

			actionExecution.StatusUpdatedAt = utils.Now()
			actionExecution.Status = neo4jEntity.FlowActionExecutionStatusTechError

			_, err = s.commonServices.Neo4jRepositories.FlowActionExecutionWriteRepository.Merge(ctx, nil, actionExecution)
			if err != nil {
				tracing.TraceErr(span, err)
			}

			continue
		}
	}
}

func (s *flowExecutionService) ComputeFlowStatistics() {
	ctx, cancel := utils.GetContextWithTimeout(context.Background(), utils.HalfOfHourDuration)
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "FlowExecutionService.ComputeFlowStatistics")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	flowsUpdated, err := s.commonServices.Neo4jRepositories.FlowWriteRepository.UpdateStatistics(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	if flowsUpdated != nil && len(flowsUpdated) > 0 {
		for _, flowData := range flowsUpdated {
			s.commonServices.Events.Publisher.PublishEventCompletedBulk(ctx, flowData.Tenant, flowData.Strings, model.FLOW, utils.NewEventCompletedDetails().WithUpdate())
		}
	}
}
