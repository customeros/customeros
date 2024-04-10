package service

import (
	"github.com/machinebox/graphql"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	authService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/service"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/redis/go-redis/v9"
)

type Services struct {
	PostgresRepositories *postgresRepository.Repositories

	AuthServices *authService.Services

	MailService         MailService
	CustomerOsService   CustomerOSService
	RedisService        RedisService
	FileStoreApiService fsc.FileStoreApiService
	CommonServices      *commonService.Services
}

func InitServices(graphqlClient *graphql.Client, redisClient *redis.Client, cfg *c.Config, db *c.StorageDB) *Services {
	postgresRepositories := postgresRepository.InitRepositories(db.GormDB)

	cosService := NewCustomerOSService(graphqlClient, cfg)

	services := Services{
		PostgresRepositories: postgresRepositories,
		CustomerOsService:    cosService,
		RedisService:         NewRedisService(redisClient, cfg),
		FileStoreApiService:  fsc.NewFileStoreApiService(&cfg.FileStoreApiConfig),
		CommonServices:       commonService.InitServices(db.GormDB, nil),
	}

	services.MailService = NewMailService(cfg, &services)
	services.AuthServices = authService.InitServices(&cfg.AuthConfig, db.GormDB)

	return &services
}
