package service

import (
	"github.com/machinebox/graphql"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	authService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/service"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/redis/go-redis/v9"
)

type Services struct {
	AuthServices *authService.Services

	MailService         MailService
	CustomerOsService   CustomerOSService
	RedisService        RedisService
	FileStoreApiService FileStoreApiService
	CommonServices      *commonService.Services
}

func InitServices(graphqlClient *graphql.Client, redisClient *redis.Client, cfg *c.Config, db *c.StorageDB) *Services {
	cosService := NewCustomerOSService(graphqlClient, cfg)

	services := Services{
		CustomerOsService:   cosService,
		RedisService:        NewRedisService(redisClient, cfg),
		FileStoreApiService: NewFileStoreApiService(cfg),
		CommonServices:      commonService.InitServices(db.GormDB, nil),
	}

	services.MailService = NewMailService(cfg, &services)
	services.AuthServices = authService.InitServices(&cfg.AuthConfig, db.GormDB)

	return &services
}
