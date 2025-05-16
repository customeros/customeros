package icp

import (
	"context"
	"log"

	nats_common "github.com/customeros/customeros/packages/server/customer-os-common-module/nats"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/customeros/customeros/packages/server/enums"

	"github.com/customeros/customeros/packages/server/core-crm/interfaces"
	"github.com/customeros/customeros/packages/server/core-crm/internal/telemetry"
)

type icpService struct {
	icpFitConsumer       *nats_common.AsyncEventsConsumer
	organizationConsumer *nats_common.AsyncEventsConsumer
	tenantConsumer       *nats_common.AsyncEventsConsumer
	repositories         *postgres_repository.Repositories
}

func NewICPService(
	natsConn *nats_common.NATSConnections,
	repository *postgres_repository.Repositories,
) interfaces.NatsService {
	// configure nats consumer
	organizationConsumer, err := setupAsyncConsumer(natsConn, enums.StreamOrganization, enums.ServicesICP, enums.EventOrganizationCreated)
	if err != nil {
		log.Fatalf("Cannot setup Organization consumer")
	}
	tenantConsumer, err := setupAsyncConsumer(natsConn, enums.StreamTenant, enums.ServicesICP, enums.EventTenantCreated)
	if err != nil {
		log.Fatalf("Cannot setup Tenant consumer")
	}
	icpFitConsumer, err := setupAsyncConsumer(natsConn, enums.StreamICP, enums.ServicesICP, enums.EventICPFitDetermined)
	if err != nil {
		log.Fatalf("Cannot setup Content consumer")
	}

	s := &icpService{
		icpFitConsumer:       icpFitConsumer,
		organizationConsumer: organizationConsumer,
		tenantConsumer:       tenantConsumer,
		repositories:         repository,
	}

	// register event handlers
	icpFitConsumer.RegisterHandler(enums.EventICPFitDetermined.String(), s.handleICPFitResponse)
	organizationConsumer.RegisterHandler(enums.EventOrganizationCreated.String(), s.handleCompanyCreated)
	tenantConsumer.RegisterHandler(enums.EventTenantCreated.String(), s.handleTenantCreated)
	return s
}

func setupAsyncConsumer(natsConn *nats_common.NATSConnections, stream enums.NatsStream, service enums.Services, subject enums.NatsEventType) (*nats_common.AsyncEventsConsumer, error) {
	asyncEventConfig := &nats_common.AsyncConsumerConfig{
		StreamName:        stream,
		ServiceName:       service,
		SubscribedSubject: subject.String(),
	}

	return nats_common.NewAsyncEventsConsumer(natsConn, asyncEventConfig)
}

// Start begins listening for organization events and processing them
func (s *icpService) Start(ctx context.Context) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "icpService.Start")
	defer span.Finish()

	err := s.organizationConsumer.Start(ctx)
	if err != nil {
		span.TraceError(err)
		return err
	}

	return s.tenantConsumer.Start(ctx)
}

func (s *icpService) Stop() {
	s.organizationConsumer.Stop()
	s.tenantConsumer.Stop()
}
