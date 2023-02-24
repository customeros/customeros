package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type tenantRepository struct {
	driver *neo4j.DriverWithContext
}

type TenantRepository interface {
	TenantExists(ctx context.Context, name string) (bool, error)
}

func NewTenantRepository(driver *neo4j.DriverWithContext) TenantRepository {
	return &tenantRepository{
		driver: driver,
	}
}

func (u *tenantRepository) TenantExists(ctx context.Context, name string) (bool, error) {
	session := (*u.driver).NewSession(
		ctx,
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeRead,
			BoltLogger: neo4j.ConsoleBoltLogger(),
		},
	)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$name})
			RETURN t.id`,
			map[string]interface{}{
				"name": name,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return false, err
	}
	if len(records.([]*neo4j.Record)) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
