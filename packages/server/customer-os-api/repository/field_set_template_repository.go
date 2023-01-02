package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type FieldSetTemplateRepository interface {
	createFieldSetTemplateInTx(entityTemplateId string, entity *entity.FieldSetTemplateEntity, tx neo4j.Transaction) error
	FindAllByEntityTemplateId(entityTemplateId string) (any, error)
	FindByFieldSetId(fieldSetId string) (any, error)
}

type fieldSetTemplateRepository struct {
	driver       *neo4j.Driver
	repositories *Repositories
}

func NewFieldSetTemplateRepository(driver *neo4j.Driver, repositories *Repositories) FieldSetTemplateRepository {
	return &fieldSetTemplateRepository{
		driver:       driver,
		repositories: repositories,
	}
}

func (r *fieldSetTemplateRepository) createFieldSetTemplateInTx(entityTemplateId string, entity *entity.FieldSetTemplateEntity, tx neo4j.Transaction) error {
	queryResult, err := tx.Run(`
			MATCH (e:EntityTemplate {id:$entityTemplateId})
			MERGE (e)-[:CONTAINS]->(f:FieldSetTemplate {
				id: randomUUID(),
				name: $name
			}) ON CREATE SET f.order=$order
			RETURN f`,
		map[string]any{
			"entityTemplateId": entityTemplateId,
			"name":             entity.Name,
			"order":            entity.Order,
		})

	record, err := queryResult.Single()
	if err != nil {
		return err
	}
	fieldSetTemplateId := utils.GetPropsFromNode(record.Values[0].(dbtype.Node))["id"].(string)
	for _, v := range entity.CustomFields {
		err := r.repositories.CustomFieldTemplateRepository.createCustomFieldTemplateForFieldSetInTx(fieldSetTemplateId, v, tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *fieldSetTemplateRepository) FindAllByEntityTemplateId(entityTemplateId string) (any, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (:EntityTemplate {id:$entityTemplateId})-[:CONTAINS]->(f:FieldSetTemplate) RETURN f ORDER BY f.order`,
			map[string]any{
				"entityTemplateId": entityTemplateId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
}

func (r *fieldSetTemplateRepository) FindByFieldSetId(fieldSetId string) (any, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (:FieldSet {id:$fieldSetId})-[:IS_DEFINED_BY]->(d:FieldSetTemplate)
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
