package agent_producers

import service "github.com/customeros/customeros/packages/server/customer-os-common-module/services"

type AgentProducers struct {
	NewLeadProducer         *NewLeadProducer
	NewSupportVisitProducer *NewSupportVisitProducer
	NewWebSessionProducer   *NewWebSessionProducer
}

func InitAgentProducers(services *service.CommonServices) *AgentProducers {
	return &AgentProducers{
		NewLeadProducer: NewNewLeadProducer(
			services.PostgresRepositories,
			services.Neo4jRepositories,
			services.OrganizationService,
			services.Events,
		),
		NewSupportVisitProducer: NewNewSupportVisitProducer(
			services.Events,
			services.PostgresRepositories,
		),
		NewWebSessionProducer: NewNewWebSessionProducer(
			services.Events,
			services.PostgresRepositories,
		),
	}
}
