package service

import (
	"github.com/machinebox/graphql"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/repository"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/redis/go-redis/v9"
)

type Services struct {
	MailService         MailService
	CustomerOsService   CustomerOSService
	RedisService        RedisService
	FileStoreApiService FileStoreApiService
	CommonServices      *commonService.Services
}

func InitServices(graphqlClient *graphql.Client, redisClient *redis.Client, config *c.Config, db *c.StorageDB) *Services {
	cosService := NewCustomerOSService(graphqlClient, config)
	apiKeyRepository := repository.NewApiKeyRepository(db)
	services := Services{
		CustomerOsService:   cosService,
		MailService:         NewMailService(config, cosService, apiKeyRepository),
		RedisService:        NewRedisService(redisClient, config),
		FileStoreApiService: NewFileStoreApiService(config),
		CommonServices:      commonService.InitServices(db.GormDB, nil),
	}

	return &services
}
