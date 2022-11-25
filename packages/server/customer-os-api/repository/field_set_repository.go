package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type FieldSetRepository interface {
	LinkWithFieldSetDefinitionInTx(tx neo4j.Transaction, tenant, fieldSetId, definitionId string) error
	MergeFieldSetToContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, *dbtype.Relationship, error)

	UpdateForContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, *dbtype.Relationship, error)
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

func (r *fieldSetRepository) LinkWithFieldSetDefinitionInTx(tx neo4j.Transaction, tenant, fieldSetId, definitionId string) error {
	txResult, err := tx.Run(`
			MATCH (f:FieldSet {id:$fieldSetId})<-[:HAS_COMPLEX_PROPERTY]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[:IS_DEFINED_BY]->(e:EntityDefinition)-[:CONTAINS]->(d:FieldSetDefinition {id:$definitionId})
			MERGE (f)-[r:IS_DEFINED_BY]->(d)
			RETURN r`,
		map[string]any{
			"definitionId": definitionId,
			"fieldSetId":   fieldSetId,
			"tenant":       tenant,
		})
	if err != nil {
		return err
	}
	_, err = txResult.Single()
	return err
}

func (r *fieldSetRepository) MergeFieldSetToContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (f:FieldSet {name: $name})<-[r:HAS_COMPLEX_PROPERTY]-(c)
            ON CREATE SET f.id=randomUUID(), r.added=datetime({timezone: 'UTC'})
			RETURN f, r`,
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
			"name":      entity.Name,
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}

func (r *fieldSetRepository) UpdateForContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[r:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId})
            SET s.name=$name
			RETURN s, r`,
		map[string]interface{}{
			"tenant":     tenant,
			"contactId":  contactId,
			"fieldSetId": entity.Id,
			"name":       entity.Name,
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}
