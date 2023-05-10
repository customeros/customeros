package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type ExternalSystemRepository interface {
	Merge(ctx context.Context, tenant, externalSystem string) error
}

type externalSystemRepository struct {
	driver *neo4j.DriverWithContext
}

func NewExternalSystemRepository(driver *neo4j.DriverWithContext) ExternalSystemRepository {
	return &externalSystemRepository{
		driver: driver,
	}
}

func (r *externalSystemRepository) Merge(ctx context.Context, tenant, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})" +
		" ON CREATE SET e.name=$externalSystem, e.createdAt=$now, e.updatedAt=$now, e:%s "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "ExternalSystem_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"now":            time.Now().UTC(),
			})
		return nil, err
	})
	return err
}
