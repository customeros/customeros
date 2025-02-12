package agent_listeners

import (
	"fmt"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
)

type AgentListeners struct {
	listeners map[enum.AgentListenerEvent]interfaces.AgentListenerUntyped
}

func InitAgentListeners(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *AgentListeners {
	var listeners []interfaces.AgentListenerUntyped
	listeners = append(listeners, NewCompanyIdentifiedListener(logger, postgresRepositories, agentRunnerService))
	listeners = append(listeners, NewCompanyNeedsHelpListener(logger, postgresRepositories))
	listeners = append(listeners, NewContactAddedToCampaignListener(logger, postgresRepositories, agentRunnerService))
	listeners = append(listeners, NewEmailBouncedListener(logger, postgresRepositories, agentRunnerService))
	listeners = append(listeners, NewEmailReplyReceivedListener(logger, postgresRepositories, agentRunnerService))
	listeners = append(listeners, NewIcpFitListener(logger, postgresRepositories))
	listeners = append(listeners, NewIcpNotAFitListener(logger, postgresRepositories))
	listeners = append(listeners, NewMeetingLoggedListener(logger, postgresRepositories))
	listeners = append(listeners, NewNewLeadListener(logger, postgresRepositories, agentRunnerService))
	listeners = append(listeners, NewNewMeetingRecordingListener(logger, postgresRepositories, agentRunnerService))
	listeners = append(listeners, NewNewWebSessionListener(logger, postgresRepositories, agentRunnerService))
	listeners = append(listeners, NewRunIcpQualifierAgent(logger, postgresRepositories, agentRunnerService))
	listeners = append(listeners, NewWebVisitorIdentifiedListener(logger, postgresRepositories))
	listeners = append(listeners, NewWebVisitorNotIdentifiedListener(logger, postgresRepositories))
	listeners = append(listeners, NewStartInvoiceRun(logger, postgresRepositories, agentRunnerService))
	listeners = append(listeners, NewStartInvoiceRunWithAutopayment(logger, postgresRepositories, agentRunnerService))

	agentListeners := AgentListeners{
		listeners: make(map[enum.AgentListenerEvent]interfaces.AgentListenerUntyped),
	}

	for _, listener := range listeners {
		agentListeners.listeners[listener.Type()] = listener
	}

	return &agentListeners
}

// GetListener retrieves the untyped listener based on the listener event.
func (l *AgentListeners) GetListener(listenerEvent enum.AgentListenerEvent) (interfaces.AgentListenerUntyped, error) {
	executor, exists := l.listeners[listenerEvent]
	if !exists {
		return nil, fmt.Errorf("unsupported listener event: %s", listenerEvent)
	}
	return executor, nil
}
