package agent_listeners

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func MarkAgentExecutionCompletedWithGoalAchieved(ctx context.Context, postgresRepositories *postgres_repository.Repositories, agentExecutionId string, goalAchieved bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MarkAgentExecutionWithGoalAchieved")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	span.LogFields(log.String("agentExecutionId", agentExecutionId))

	var agentExecution *postgres_entity.AgentExecution
	var err error
	if agentExecutionId != "" {
		agentExecution, err = postgresRepositories.AgentExecutionRepository.GetById(ctx, agentExecutionId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}
	if agentExecution == nil {
		err = fmt.Errorf("agent execution not found")
		tracing.TraceErr(span, err)
		return err
	}

	// update execution with goal achieved
	_, err = postgresRepositories.AgentExecutionRepository.Completed(ctx, agentExecution.ID, utils.Ptr(goalAchieved))
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
