package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type FieldSetRepository interface {
	LinkWithFieldSetDefinitionInTx(tenant string, fieldSetId string, fieldSetDefinitionId string, tx neo4j.Transaction) error
}

type fieldSetRepository struct {
	driver *neo4j.Driver
	repos  *RepositoryContainer
}

func NewFieldSetRepository(driver *neo4j.Driver, repos *RepositoryContainer) FieldSetRepository {
	return &fieldSetRepository{
		driver: driver,
		repos:  repos,
	}
}

func (r *fieldSetRepository) LinkWithFieldSetDefinitionInTx(tenant string, fieldSetId string, fieldSetDefinitionId string, tx neo4j.Transaction) error {
	txResult, err := tx.Run(`
			MATCH (f:FieldSet {id:$fieldSetId})<-[:HAS_COMPLEX_PROPERTY]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[:IS_DEFINED_BY]->(e:EntityDefinition)-[:CONTAINS]->(d:FieldSetDefinition {id:$fieldSetDefinitionId})
			MERGE (f)-[r:IS_DEFINED_BY]->(d)
			RETURN r`,
		map[string]any{
			"fieldSetDefinitionId": fieldSetDefinitionId,
			"fieldSetId":           fieldSetId,
			"tenant":               tenant,
		})
	if err != nil {
		return err
	}
	_, err = txResult.Single()
	return err
}
