package service

import (
	"github.com/machinebox/graphql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	service "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/client"
	authService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/service"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Services struct {
	AuthServices *authService.Services

	CustomerOsService   CustomerOSService
	CustomerOSApiClient service.CustomerOSApiClient

	MailService MailService

	RedisService        RedisService
	FileStoreApiService fsc.FileStoreApiService
	CommonServices      *commonService.Services
}

func InitServices(graphqlClient *graphql.Client, redisClient *redis.Client, cfg *c.Config, db *gorm.DB, driver *neo4j.DriverWithContext, neo4jDatabase string) *Services {
	cosService := NewCustomerOSService(graphqlClient, cfg)
	customerOSApiClient := service.NewCustomerOsClient(cfg.Service.CustomerOsAPI, cfg.Service.CustomerOsAPIKey)

	services := Services{
		CustomerOsService:   cosService,
		CustomerOSApiClient: customerOSApiClient,
		RedisService:        NewRedisService(redisClient, cfg),
		FileStoreApiService: fsc.NewFileStoreApiService(&cfg.FileStoreApiConfig),
		CommonServices:      commonService.InitServices(db, driver, neo4jDatabase, nil),
	}

	services.MailService = NewMailService(cfg, &services)
	services.AuthServices = authService.InitServices(&cfg.AuthConfig, services.CommonServices, db)

	return &services
}
