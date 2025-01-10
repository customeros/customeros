package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/net/context"
)

type IndustryWriteRepository interface {
	ReplaceForOrganization(ctx context.Context, tenant string, industryCode, organizationId string) error
}

type industryWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewIndustryWriteRepository(driver *neo4j.DriverWithContext, database string) IndustryWriteRepository {
	return &industryWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *industryWriteRepository) ReplaceForOrganization(ctx context.Context, tenant string, industryCode, organizationId string) error {
	//TODO implement me
	panic("implement me")
}
