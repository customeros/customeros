package agent_producers

import (
	"context"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type NewSupportVisitProducer struct {
	events             *events.EventsService
	postgresRepository *postgres_repository.Repositories
}

func NewNewSupportVisitProducer(
	events *events.EventsService,
	postgresRepository *postgres_repository.Repositories,
) *NewSupportVisitProducer {
	return &NewSupportVisitProducer{
		events:             events,
		postgresRepository: postgresRepository,
	}
}

const WebSessionLookbackPeriodInHours = 48

// Add all Agent types subscribed to this event here
func (p *NewSupportVisitProducer) subscribedAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentSupportSpotter,
	}
}

func (p *NewSupportVisitProducer) Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "NewSupportVisitProducer.Execute")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	// get active subscribed agents (make common?)

	// get rules for event production from listener config

	// get websessions ready for processing

	// process websessions
}
