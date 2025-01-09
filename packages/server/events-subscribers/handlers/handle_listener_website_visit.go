package handlers

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"
)

// add all subsribed agents for this handler here
var SubscribedAgents = [1]enum.Agent{
	enum.AgentVisitorID,
}

func HandleWebsiteVisitorEvent(c context.Context, s *service.Services, sourceEvent commonEnum.FlowListenerEvent, eventData *data_fields.WebsiteVisitEvent) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.HandleWebsiteVisitorEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant: eventData.Tenant,
	})

	// find list of active agents that listen on event
	agents, err := s.PostgresRepositories.AgentsRepository.FindAllFromAgentsList(ctx, SubscribedAgents[:])
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if agents == nil {
		return nil
	}

	var errs error
	for _, agent := range agents {
		var loopErr error
		// create start agent execution record
		agentExecutionRecord := entity.AgentExecution{
			AgentID:      &agent.ID,
			TriggerEvent: eventData.Type(),
			Status:       enum.AgentExecutionRunning.String(),
			StartedAt:    utils.NowPtr(),
		}
		execution, err := s.PostgresRepositories.AgentExecutionRepository.Create(ctx, agentExecutionRecord)
		if err != nil {
			tracing.TraceErr(span, err)
			loopErr = multierr.Append(loopErr, err)
		}
		if execution == nil {
			err := errors.New("unable to create agent execution record")
			tracing.TraceErr(span, err)
			loopErr = multierr.Append(loopErr, err)
		}

		// call agent service
		if execution != nil {
			err = s.AgentService.VisitorIDAgent(ctx, eventData)
			if err != nil {
				tracing.TraceErr(span, err)
				loopErr = multierr.Append(loopErr, err)
			}
		}

		// update automation execution record with results
		if loopErr != nil {
			errorMessage := loopErr.Error()
			agentExecutionRecord.ErrorMessage = &errorMessage
			agentExecutionRecord.Status = enum.AgentExecutionFail.String()
		} else {
			agentExecutionRecord.Status = enum.AgentExecutionSuccess.String()
			agentExecutionRecord.CompletedAt = utils.NowPtr()
		}

		execution, err = s.PostgresRepositories.AgentExecutionRepository.Update(ctx, agentExecutionRecord)
		if err != nil {
			tracing.TraceErr(span, err)
			loopErr = multierr.Append(loopErr, err)
		}

		errs = multierr.Append(errs, loopErr)
	}

	return errs
}
