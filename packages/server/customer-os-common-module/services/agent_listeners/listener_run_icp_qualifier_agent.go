package agent_listeners

import (
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
)

type RunIcpQualifierAgent struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*RunIcpQualifierAgent)(nil)
)

func NewRunIcpQualifierAgent(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *RunIcpQualifierAgent {
	return &RunIcpQualifierAgent{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.NewLead](), // subscribed event
			events.QueueAgents,                 // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *RunIcpQualifierAgent) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentICPQualifier,
	}
}

func (l *RunIcpQualifierAgent) Type() enum.AgentListenerEvent {
	return enum.EventRunICPQualifierAgent
}

func (l *RunIcpQualifierAgent) Name() string {
	return "Run ICP Qualifier Agent"
}

func (l *RunIcpQualifierAgent) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}
