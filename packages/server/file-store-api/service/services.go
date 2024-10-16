package service

import (
	"github.com/machinebox/graphql"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/logger"
)

type Services struct {
	FileService FileService
}

func InitServices(cfg *config.Config, commonServices *commonService.Services, graphqlClient *graphql.Client, log logger.Logger) *Services {
	return &Services{
		FileService: NewFileService(cfg, commonServices, graphqlClient, log),
	}
}
