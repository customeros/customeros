package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FlowExecutionService interface {
	RampUpMailboxes()
	ExecuteScheduledFlowActions()
	ComputeFlowStatistics()
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

func (s *flowExecutionService) rampUpMailbox(ctx context.Context, mailbox *entity.TenantSettingsMailbox) error {
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

		err := s.commonServices.PostgresRepositories.TenantSettingsMailboxRepository.Merge(ctx, mailbox)
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
			Tenant: tenant,
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

	tenants, err := s.commonServices.TenantService.GetAllTenants(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	for _, tenant := range tenants {
		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant: tenant.Name,
		})

		flows, err := s.commonServices.FlowService.FlowGetList(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		for _, flow := range *flows {
			ctx = common.WithCustomContext(ctx, &common.CustomContext{
				Tenant: tenant.Name,
			})

			flowChanged := false

			onHold, err := s.commonServices.Neo4jRepositories.FlowParticipantReadRepository.CountWithStatus(ctx, flow.Id, neo4jEntity.FlowParticipantStatusOnHold)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			if flow.OnHold != onHold {
				flowChanged = true

				err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateInt64Property(ctx, tenant.Name, model.NodeLabelFlow, flow.Id, "onHold", onHold)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}
			}

			ready, err := s.commonServices.Neo4jRepositories.FlowParticipantReadRepository.CountWithStatus(ctx, flow.Id, neo4jEntity.FlowParticipantStatusReady)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			if flow.Ready != ready {
				flowChanged = true

				err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateInt64Property(ctx, tenant.Name, model.NodeLabelFlow, flow.Id, "ready", ready)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}
			}

			scheduled, err := s.commonServices.Neo4jRepositories.FlowParticipantReadRepository.CountWithStatus(ctx, flow.Id, neo4jEntity.FlowParticipantStatusScheduled)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			if flow.Scheduled != scheduled {
				flowChanged = true

				err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateInt64Property(ctx, tenant.Name, model.NodeLabelFlow, flow.Id, "scheduled", scheduled)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}
			}

			inProgress, err := s.commonServices.Neo4jRepositories.FlowParticipantReadRepository.CountWithStatus(ctx, flow.Id, neo4jEntity.FlowParticipantStatusInProgress)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			if flow.InProgress != inProgress {
				flowChanged = true

				err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateInt64Property(ctx, tenant.Name, model.NodeLabelFlow, flow.Id, "inProgress", inProgress)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}
			}

			completed, err := s.commonServices.Neo4jRepositories.FlowParticipantReadRepository.CountWithStatus(ctx, flow.Id, neo4jEntity.FlowParticipantStatusCompleted)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			if flow.Completed != completed {
				flowChanged = true

				err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateInt64Property(ctx, tenant.Name, model.NodeLabelFlow, flow.Id, "completed", completed)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}
			}

			goalAchieved, err := s.commonServices.Neo4jRepositories.FlowParticipantReadRepository.CountWithStatus(ctx, flow.Id, neo4jEntity.FlowParticipantStatusGoalAchieved)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			if flow.GoalAchieved != goalAchieved {
				flowChanged = true

				err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateInt64Property(ctx, tenant.Name, model.NodeLabelFlow, flow.Id, "goalAchieved", goalAchieved)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}
			}

			if !flowChanged {
				continue
			}

			err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateInt64Property(ctx, tenant.Name, model.NodeLabelFlow, flow.Id, "total", onHold+ready+scheduled+inProgress+completed+goalAchieved)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			s.commonServices.RabbitMQService.PublishEventCompleted(ctx, tenant.Name, flow.Id, model.FLOW, utils.NewEventCompletedDetails().WithUpdate())
		}

	}

}
