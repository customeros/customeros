package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type FieldSetDefinitionRepository interface {
	createFieldSetDefinitionInTx(entityDefId string, entity *entity.FieldSetDefinitionEntity, tx neo4j.Transaction) error
	FindAllByEntityDefinitionId(entityDefinitionId string) (any, error)
	FindByFieldSetId(fieldSetId string) (any, error)
}

type fieldSetDefinitionRepository struct {
	driver *neo4j.Driver
	repos  *RepositoryContainer
}

func NewFieldSetDefinitionRepository(driver *neo4j.Driver, repos *RepositoryContainer) FieldSetDefinitionRepository {
	return &fieldSetDefinitionRepository{
		driver: driver,
		repos:  repos,
	}
}

func (r *fieldSetDefinitionRepository) createFieldSetDefinitionInTx(entityDefId string, entity *entity.FieldSetDefinitionEntity, tx neo4j.Transaction) error {
	queryResult, err := tx.Run(`
			MATCH (e:EntityDefinition {id:$entityDefinitionId})
			MERGE (e)-[:CONTAINS]->(f:FieldSetDefinition {
				id: randomUUID(),
				name: $name
			}) ON CREATE SET f.order=$order
			RETURN f`,
		map[string]any{
			"entityDefinitionId": entityDefId,
			"name":               entity.Name,
			"order":              entity.Order,
		})

	record, err := queryResult.Single()
	if err != nil {
		return err
	}
	fieldSetDefinitionId := utils.GetPropsFromNode(record.Values[0].(dbtype.Node))["id"].(string)
	for _, v := range entity.CustomFields {
		err := r.repos.CustomFieldDefinitionRepository.createCustomFieldDefinitionForFieldSetInTx(fieldSetDefinitionId, v, tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *fieldSetDefinitionRepository) FindAllByEntityDefinitionId(entityDefinitionId string) (any, error) {
	session := (*r.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (:EntityDefinition {id:$entityDefinitionId})-[:CONTAINS]->(f:FieldSetDefinition) RETURN f ORDER BY f.order`,
			map[string]any{
				"entityDefinitionId": entityDefinitionId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
}

func (r *fieldSetDefinitionRepository) FindByFieldSetId(fieldSetId string) (any, error) {
	session := (*r.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (:FieldSet {id:$fieldSetId})-[:IS_DEFINED_BY]->(d:FieldSetDefinition)
					RETURN d`,
			map[string]any{
				"fieldSetId": fieldSetId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
}
