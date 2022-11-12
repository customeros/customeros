package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

type ContactRepository interface {
	LinkWithEntityDefinitionInTx(tenant string, contactId string, entityDefinitionId string, tx neo4j.Transaction) error
}

type contactRepository struct {
	driver *neo4j.Driver
	repos  *RepositoryContainer
}

func NewContactRepository(driver *neo4j.Driver, repos *RepositoryContainer) ContactRepository {
	return &contactRepository{
		driver: driver,
		repos:  repos,
	}
}

func (r *contactRepository) LinkWithEntityDefinitionInTx(tenant string, contactId string, entityDefinitionId string, tx neo4j.Transaction) error {
	txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})-[:USES_ENTITY_DEFINITION]->(e:EntityDefinition {id:$entityDefinitionId})
			WHERE e.extends=$extends
			MERGE (c)-[r:IS_DEFINED_BY]->(e)
			RETURN r`,
		map[string]any{
			"entityDefinitionId": entityDefinitionId,
			"contactId":          contactId,
			"tenant":             tenant,
			"extends":            model.EntityDefinitionExtensionContact,
		})
	if err != nil {
		return err
	}
	_, err = txResult.Single()
	return err
}
