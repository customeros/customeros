package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type RevenueReportRepository interface {
	GetMonthlyRevenueReportPerCustomerData(ctx context.Context, tenant string, year, month int) ([]map[string]interface{}, error)
}

type revenueReportRepository struct {
	driver *neo4j.DriverWithContext
}

func NewRevenueReportRepository(driver *neo4j.DriverWithContext) RevenueReportRepository {
	return &revenueReportRepository{
		driver: driver,
	}
}

func (r revenueReportRepository) GetMonthlyRevenueReportPerCustomerData(ctx context.Context, tenant string, year, month int) ([]map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}
