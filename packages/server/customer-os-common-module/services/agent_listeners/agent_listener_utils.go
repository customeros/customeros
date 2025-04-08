package agent_listeners

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

func MarkAgentExecutionCompletedWithGoalAchieved(ctx context.Context, postgresRepositories *postgres_repository.Repositories, agentExecutionId string, goalAchieved bool) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MarkAgentExecutionWithGoalAchieved")
	defer spans.Finish()

	spans.LogKV("agentExecutionId", agentExecutionId)

	var agentExecution *postgres_entity.AgentExecution
	var err error
	if agentExecutionId != "" {
		agentExecution, err = postgresRepositories.AgentExecutionRepository.GetById(ctx, agentExecutionId)
		if err != nil {
			spans.TraceError(err)
			return err
		}
	}
	if agentExecution == nil {
		err = fmt.Errorf("agent execution not found")
		spans.TraceError(err)
		return err
	}

	// update execution with goal achieved
	_, err = postgresRepositories.AgentExecutionRepository.Completed(ctx, agentExecution.ID, utils.Ptr(goalAchieved))
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
