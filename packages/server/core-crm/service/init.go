package service

import (
	"github.com/customeros/customeros/packages/server/core-crm/service/outbox_processor"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	nats_common "github.com/customeros/customeros/packages/server/customer-os-common-module/nats"
	services_common "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type Services struct {
	CommonServices  *services_common.CommonServices
	OutboxProcessor *outbox_processor.OutboxProcessor
}

func InitServices(
	log logger.Logger,
	neo4jRepositories *neo4j_repository.Repositories,
	postgresRepositories *postgres_repository.Repositories,
	warehouseRepositories *postgres_repository.WarehouseRepositories,
	cfg *config.CommonConfig,
	natsConn *nats_common.NATSConnections,
) *Services {

	outboxProcessor := outbox_processor.NewOutboxProcessor(
		natsConn,
		postgresRepositories,
	)

	return &Services{
		OutboxProcessor: outboxProcessor,
		CommonServices: services_common.InitCommonServices(
			log,
			neo4jRepositories,
			postgresRepositories,
			warehouseRepositories,
			cfg,
			natsConn,
			nil,
		),
	}
}
