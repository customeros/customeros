package agent_producers

import (
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
)

type NewSupportVisitProducer struct {
	events             *events.EventsService
	postgresRepository *postgres_repository.Repositories
}

func NewNewSupportVisitProducer(
	events *events.EventsService,
	postgresRepository *postgres_repository.Repositories,
)
