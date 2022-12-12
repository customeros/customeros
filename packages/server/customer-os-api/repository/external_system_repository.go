package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
)

type ExternalSystemRepository interface {
	LinkContactWithExternalSystemInTx(tx neo4j.Transaction, tenant, contactId string, relationship entity.ExternalReferenceRelationship) error
}

type externalSystemRepository struct {
	driver *neo4j.Driver
}

func NewExternalSystemRepository(driver *neo4j.Driver) ExternalSystemRepository {
	return &externalSystemRepository{
		driver: driver,
	}
}

func (e *externalSystemRepository) LinkContactWithExternalSystemInTx(tx neo4j.Transaction, tenant, contactId string, relationship entity.ExternalReferenceRelationship) error {
	queryResult, err := tx.Run(`
			MATCH (e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				  (c:Contact {id:$contactId})
			MERGE (c)-[r:IS_LINKED_WITH {id:$referenceId}]->(e)
			ON CREATE SET r.syncDate=$syncDate
			ON MATCH SET r.syncDate=$syncDate	
			RETURN r`,
		map[string]any{
			"contactId":        contactId,
			"tenant":           tenant,
			"syncDate":         relationship.SyncDate,
			"referenceId":      relationship.Id,
			"externalSystemId": relationship.ExternalSystemId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single()
	return err
}
