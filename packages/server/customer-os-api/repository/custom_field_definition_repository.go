package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type CustomFieldDefinitionRepository interface {
	createCustomFieldDefinitionForEntityInTx(entityDefinitionId string, entity *entity.CustomFieldDefinitionEntity, tx neo4j.Transaction) error
	createCustomFieldDefinitionForFieldSetInTx(fieldSetDefinitionId string, entity *entity.CustomFieldDefinitionEntity, tx neo4j.Transaction) error
	FindAllByEntityDefinitionId(entityDefinitionId string) (any, error)
	FindAllByEntityFieldSetDefinitionId(fieldSetDefinitionId string) (any, error)
	FindByCustomFieldId(fieldSetId string) (any, error)
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

func (r *customFieldDefinitionRepository) createCustomFieldDefinitionForEntityInTx(entityDefinitionId string, entity *entity.CustomFieldDefinitionEntity, tx neo4j.Transaction) error {
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
			"entityDefinitionId": entityDefinitionId,
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

func (r *customFieldDefinitionRepository) createCustomFieldDefinitionForFieldSetInTx(fieldSetDefinition string, entity *entity.CustomFieldDefinitionEntity, tx neo4j.Transaction) error {
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
			"fieldSetDefinitionId": fieldSetDefinition,
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

func (r *customFieldDefinitionRepository) FindAllByEntityDefinitionId(entityDefinitionId string) (any, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (:EntityDefinition {id:$entityDefinitionId})-[:CONTAINS]->(f:CustomFieldDefinition) RETURN f ORDER BY f.order`,
			map[string]any{
				"entityDefinitionId": entityDefinitionId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
}

func (r *customFieldDefinitionRepository) FindAllByEntityFieldSetDefinitionId(fieldSetDefinitionId string) (any, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (:FieldSetDefinition {id:$fieldSetDefinitionId})-[:CONTAINS]->(f:CustomFieldDefinition) RETURN f ORDER BY f.order`,
			map[string]any{
				"fieldSetDefinitionId": fieldSetDefinitionId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
}

func (r *customFieldDefinitionRepository) FindByCustomFieldId(customFieldId string) (any, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (:CustomField {id:$customFieldId})-[:IS_DEFINED_BY]->(d:CustomFieldDefinition)
					RETURN d`,
			map[string]any{
				"customFieldId": customFieldId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
}
