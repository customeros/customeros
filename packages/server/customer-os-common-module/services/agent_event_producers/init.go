package agent_producers

import service "github.com/customeros/customeros/packages/server/customer-os-common-module/services"

type AgentProducers struct {
	NewLeadProducer            *NewLeadProducer
	NewWebSessionProducer      *NewWebSessionProducer
	InvoiceProducer            *InvoiceProducer
	SendInvoiceProducer        *SendInvoiceProducer
	SendOverdueInvoiceProducer *SendOverdueInvoiceProducer
}

func InitAgentProducers(services *service.CommonServices) *AgentProducers {
	return &AgentProducers{
		NewLeadProducer: NewNewLeadProducer(
			services.PostgresRepositories,
			services.Neo4jRepositories,
			services.OrganizationService,
			services.Events,
		),
		NewWebSessionProducer: NewNewWebSessionProducer(
			services.Events,
			services.PostgresRepositories,
			services.WebscraperService,
		),
		InvoiceProducer: NewInvoiceProducer(
			services.PostgresRepositories,
			services.Neo4jRepositories,
			services.Events,
			services.Logger,
		),
		SendInvoiceProducer: NewSendInvoiceProducer(
			services.PostgresRepositories,
			services.Neo4jRepositories,
			services.Events,
			services.Logger,
		),
		SendOverdueInvoiceProducer: NewSendOverdueInvoiceProducer(
			services.PostgresRepositories,
			services.Neo4jRepositories,
			services.Events,
			services.Logger,
		),
	}
}
