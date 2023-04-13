package service

import (
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
)

type Services struct {
	FileService FileService
}

func InitServices(cfg *config.Config, graphqlClient *graphql.Client) *Services {
	return &Services{
		FileService: NewFileService(cfg, graphqlClient),
	}
}
