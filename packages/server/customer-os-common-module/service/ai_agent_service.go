package service

import (
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
)

type AIAgentService interface{}

type aiAgentService struct {
	repositories *postgresRepository.Repositories
	services     *Services
}

func NewAIAgentService(repositories *postgresRepository.Repositories, services *Services) AIAgentService {
	return &aiAgentService{
		repositories: repositories,
		services:     services,
	}
}
