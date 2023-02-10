package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type CustomFieldTemplateRepository interface {
	createCustomFieldTemplateForEntityInTx(tx neo4j.Transaction, tenant, entityTemplateId string, entity *entity.CustomFieldTemplateEntity) error
	createCustomFieldTemplateForFieldSetInTx(tx neo4j.Transaction, tenant, fieldSetTemplateId string, entity *entity.CustomFieldTemplateEntity) error
	FindAllByEntityTemplateId(entityTemplateId string) (any, error)
	FindAllByEntityFieldSetTemplateId(fieldSetTemplateId string) (any, error)
	FindByCustomFieldId(fieldSetId string) (any, error)
}

type customFieldTemplateRepository struct {
	driver *neo4j.Driver
}

func NewCustomFieldTemplateRepository(driver *neo4j.Driver) CustomFieldTemplateRepository {
	return &customFieldTemplateRepository{
		driver: driver,
	}
}

func (r *customFieldTemplateRepository) createCustomFieldTemplateForEntityInTx(tx neo4j.Transaction, tenant, entityTemplateId string, entity *entity.CustomFieldTemplateEntity) error {
	query := "MATCH (e:EntityTemplate {id:$entityTemplateId}) " +
		" MERGE (e)-[:CONTAINS]->(f:CustomFieldTemplate {id:randomUUID(), name:$name}) " +
		" ON CREATE SET f:%s, " +
		"				f.createdAt=$now, " +
		"				f.updated=$now, " +
		"  				f.order=$order, " +
		"				f.mandatory=$mandatory, " +
		"				f.type=$type, " +
		"				f.length=$length, " +
		"  				f.min=$min, " +
		"				f.max=$max"

	_, err := tx.Run(fmt.Sprintf(query, "CustomFieldTemplate_"+tenant),
		map[string]any{
			"entityTemplateId": entityTemplateId,
			"name":             entity.Name,
			"order":            entity.Order,
			"mandatory":        entity.Mandatory,
			"type":             entity.Type,
			"length":           entity.Length,
			"min":              entity.Min,
			"max":              entity.Max,
			"now":              utils.Now(),
		})

	return err
}

func (r *customFieldTemplateRepository) createCustomFieldTemplateForFieldSetInTx(tx neo4j.Transaction, tenant, fieldSetTemplateId string, entity *entity.CustomFieldTemplateEntity) error {
	query := "MATCH (d:FieldSetTemplate {id:$fieldSetTemplateId}) " +
		" MERGE (d)-[:CONTAINS]->(f:CustomFieldTemplate {id:randomUUID(), name:$name}) " +
		" ON CREATE SET f:%s, " +
		"				f.createdAt=$now, " +
		"				f.updatedAt=$now, " +
		"				f.order=$order, " +
		"				f.mandatory=$mandatory, " +
		"				f.type=$type, " +
		"				f.length=$length, " +
		"				f.min=$min, " +
		"				f.max=$max"
	_, err := tx.Run(fmt.Sprintf(query, "CustomFieldTemplate_"+tenant),
		map[string]any{
			"fieldSetTemplateId": fieldSetTemplateId,
			"name":               entity.Name,
			"order":              entity.Order,
			"mandatory":          entity.Mandatory,
			"type":               entity.Type,
			"length":             entity.Length,
			"min":                entity.Min,
			"max":                entity.Max,
			"now":                utils.Now(),
		})

	return err
}

func (r *customFieldTemplateRepository) FindAllByEntityTemplateId(entityTemplateId string) (any, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (:EntityTemplate {id:$entityTemplateId})-[:CONTAINS]->(f:CustomFieldTemplate) RETURN f ORDER BY f.order`,
			map[string]any{
				"entityTemplateId": entityTemplateId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
}

func (r *customFieldTemplateRepository) FindAllByEntityFieldSetTemplateId(fieldSetTemplateId string) (any, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (:FieldSetTemplate {id:$fieldSetTemplateId})-[:CONTAINS]->(f:CustomFieldTemplate) RETURN f ORDER BY f.order`,
			map[string]any{
				"fieldSetTemplateId": fieldSetTemplateId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
}

func (r *customFieldTemplateRepository) FindByCustomFieldId(customFieldId string) (any, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (:CustomField {id:$customFieldId})-[:IS_DEFINED_BY]->(d:CustomFieldTemplate)
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
