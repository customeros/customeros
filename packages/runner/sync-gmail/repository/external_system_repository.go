package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type ExternalSystemRepository interface {
	Merge(ctx context.Context, tenant, externalSystem string) (string, error)
}

type externalSystemRepository struct {
	driver *neo4j.DriverWithContext
}

func NewExternalSystemRepository(driver *neo4j.DriverWithContext) ExternalSystemRepository {
	return &externalSystemRepository{
		driver: driver,
	}
}

func (r *externalSystemRepository) Merge(ctx context.Context, tenant, externalSystem string) (string, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})" +
		" ON CREATE SET e.name=$externalSystem, e.createdAt=$now, e.updatedAt=$now, e:%s return e.id"

	records, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "ExternalSystem_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"now":            time.Now().UTC(),
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]*db.Record)) == 0 {
		return "", nil
	}
	return records.([]*db.Record)[0].Values[0].(string), nil
}
