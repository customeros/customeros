package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
)

type RevenueReportService interface {
	GetMonthlyRevenueReportPerCustomerData(ctx context.Context, year, month int) (*entity.FileGeneratorResponseData, error)
}

type revenueReportService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewRevenueReportService(log logger.Logger, repositories *repository.Repositories, services *Services) RevenueReportService {
	return &revenueReportService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *revenueReportService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *revenueReportService) GetMonthlyRevenueReportPerCustomerData(ctx context.Context, year, month int) (*entity.FileGeneratorResponseData, error) {
	//TODO implement me
	panic("implement me")
}
