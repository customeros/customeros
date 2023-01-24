package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"time"
)

type ExternalSystemRepository interface {
	Merge(tenant, externalSystem string) error
}

type externalSystemRepository struct {
	driver *neo4j.Driver
}

func NewExternalSystemRepository(driver *neo4j.Driver) ExternalSystemRepository {
	return &externalSystemRepository{
		driver: driver,
	}
}

func (r *externalSystemRepository) Merge(tenant, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
				MATCH (t:Tenant {name:$tenant})
				MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				ON CREATE SET e.name=$externalSystem, e.createdAt=$now
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"now":            time.Now().UTC(),
			})
		return nil, err
	})
	return err
}
