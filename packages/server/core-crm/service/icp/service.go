package icp

import (
	"context"
	"log"

	nats_common "github.com/customeros/customeros/packages/server/customer-os-common-module/nats"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/customeros/customeros/packages/server/enums"

	"github.com/customeros/customeros/packages/server/core-crm/interfaces"
)

type icpService struct {
	asyncEvents  *nats_common.AsyncEventsConsumer
	repositories *postgres_repository.Repositories
}

func NewICPService(
	natsConn *nats_common.NATSConnections,
	repository *postgres_repository.Repositories,
) interfaces.NatsService {
	// configure nats consumer
	asyncEventConfig := &nats_common.AsyncConsumerConfig{
		StreamName:         enums.StreamOrganization,
		ServiceName:        "icp",
		SubscribedSubjects: []string{enums.EventOrganizationCreated.String()},
	}

	asyncEvents, err := nats_common.NewAsyncEventsConsumer(natsConn, asyncEventConfig)
	if err != nil {
		log.Fatalf("Unable to start Nats on %s", &asyncEventConfig.ServiceName)
	}

	s := &icpService{
		asyncEvents:  asyncEvents,
		repositories: repository,
	}

	// register event handler
	asyncEvents.RegisterHandler(enums.EventOrganizationCreated.String(), s.handleCompanyCreated)
	return s
}

// Start begins listening for raw email events and processing them
func (s *icpService) Start(ctx context.Context) error {
	return s.asyncEvents.Start(ctx)
}

func (s *icpService) Stop() {
	s.asyncEvents.Stop()
}
