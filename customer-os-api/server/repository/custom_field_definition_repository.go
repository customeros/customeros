package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
)

type CustomFieldDefinitionRepository interface {
	createCustomFieldDefinitionForEntityInTx(entityDefId string, entity *entity.CustomFieldDefinitionEntity, tx neo4j.Transaction) error
	createCustomFieldDefinitionForFieldSetInTx(fieldSetDefId string, entity *entity.CustomFieldDefinitionEntity, tx neo4j.Transaction) error
}

type customFieldDefinitionRepository struct {
	driver *neo4j.Driver
	repos  *RepositoryContainer
}

func NewCustomFieldDefinitionRepository(driver *neo4j.Driver, repos *RepositoryContainer) CustomFieldDefinitionRepository {
	return &customFieldDefinitionRepository{
		driver: driver,
		repos:  repos,
	}
}

func (r *customFieldDefinitionRepository) createCustomFieldDefinitionForEntityInTx(entityDefId string, entity *entity.CustomFieldDefinitionEntity, tx neo4j.Transaction) error {
	_, err := tx.Run(`
			MATCH (e:EntityDefinition {id:$entityDefinitionId})
			MERGE (e)-[:CONTAINS]->(f:CustomFieldDefinition {
				id: randomUUID(),
				name: $name
			}) ON CREATE SET 
				f.order=$order, 
				f.mandatory=$mandatory,
				f.type=$type,
				f.length=$length,
				f.min=$min,
				f.max=$max`,
		map[string]any{
			"entityDefinitionId": entityDefId,
			"name":               entity.Name,
			"order":              entity.Order,
			"mandatory":          entity.Mandatory,
			"type":               entity.Type,
			"length":             entity.Length,
			"min":                entity.Min,
			"max":                entity.Max,
		})

	return err
}

func (r *customFieldDefinitionRepository) createCustomFieldDefinitionForFieldSetInTx(fieldSetDefId string, entity *entity.CustomFieldDefinitionEntity, tx neo4j.Transaction) error {
	_, err := tx.Run(`
			MATCH (d:FieldSetDefinition {id:$fieldSetDefinitionId})
			MERGE (d)-[:CONTAINS]->(f:CustomFieldDefinition {
				id: randomUUID(),
				name: $name
			}) ON CREATE SET 
				f.order=$order, 
				f.mandatory=$mandatory,
				f.type=$type,
				f.length=$length,
				f.min=$min,
				f.max=$max`,
		map[string]any{
			"fieldSetDefinitionId": fieldSetDefId,
			"name":                 entity.Name,
			"order":                entity.Order,
			"mandatory":            entity.Mandatory,
			"type":                 entity.Type,
			"length":               entity.Length,
			"min":                  entity.Min,
			"max":                  entity.Max,
		})

	return err
}
