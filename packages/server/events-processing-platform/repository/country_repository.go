package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type CountryRepository interface {
	GetDefaultCountryCodeA3(ctx context.Context, tenant string) (string, error)
}

type countryRepository struct {
	driver *neo4j.DriverWithContext
}

func NewCountryRepository(driver *neo4j.DriverWithContext) CountryRepository {
	return &countryRepository{
		driver: driver,
	}
}

func (r *countryRepository) GetDefaultCountryCodeA3(ctx context.Context, tenant string) (string, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})
				OPTIONAL MATCH (tenant)-[:DEFAULT_COUNTRY]->(dc:Country)
				RETURN COALESCE(dc.codeA3, "") AS countryCodeA3`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsString(ctx, queryResult, err)
		}
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
