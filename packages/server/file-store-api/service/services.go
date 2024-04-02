package service

import (
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/logger"
)

type Services struct {
	FileService FileService
}

func InitServices(cfg *config.Config, graphqlClient *graphql.Client, log logger.Logger) *Services {
	return &Services{
		FileService: NewFileService(cfg, graphqlClient, log),
	}
}
