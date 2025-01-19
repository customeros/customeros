package service

import (
	"github.com/customeros/customeros/packages/runner/sync-gmail-raw/config"
	"github.com/customeros/customeros/packages/runner/sync-gmail-raw/logger"
	"github.com/customeros/customeros/packages/runner/sync-gmail-raw/repository"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Services struct {
	cfg *config.Config

	Repositories *repository.Repositories

	CommonServices *commonService.CommonServices

	EmailService   EmailService
	MeetingService MeetingService
}

func InitServices(driver *neo4j.DriverWithContext, postgresDB *commonConfig.PostgresDB, cfg *config.Config, logger *logger.ExtendedLogger) *Services {
	services := new(Services)
	services.cfg = cfg
	services.Repositories = repository.InitRepos(driver, postgresDB)

	neo4jRepositories := neo4jrepository.InitNeo4jRepositories(driver, cfg.Neo4jDb.Database)
	postgresRepositories := postgresRepository.InitRepositories(postgresDB)

	services.CommonServices = commonService.InitCommonServices(
		logger,
		neo4jRepositories,
		postgresRepositories,
		&cfg.CommonConfig,
		grpc_client.InitClients(nil),
		&commonService.InitOptions{},
	)

	services.EmailService = NewEmailService(cfg, services.Repositories, services)
	services.MeetingService = NewMeetingService(cfg, services.Repositories, services)

	return services
}
