package service

import (
	"github.com/machinebox/graphql"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/repository"
	"github.com/redis/go-redis/v9"
)

type Services struct {
	MailService         MailService
	CustomerOsService   CustomerOSService
	RedisService        RedisService
	FileStoreApiService FileStoreApiService
}

func InitServices(graphqlClient *graphql.Client, redisClient *redis.Client, config *c.Config, db *c.StorageDB) *Services {
	cosService := NewCustomerOSService(graphqlClient, config)
	apiKeyRepository := repository.NewApiKeyRepository(db)
	services := Services{
		CustomerOsService:   cosService,
		MailService:         NewMailService(config, cosService, apiKeyRepository),
		RedisService:        NewRedisService(redisClient, config),
		FileStoreApiService: NewFileStoreApiService(config),
	}

	return &services
}
