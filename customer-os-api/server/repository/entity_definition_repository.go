package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type EntityDefinitionRepository interface {
	Create(tenant string, entity *entity.EntityDefinitionEntity) (any, error)
	FindAllByTenant(tenant string) (any, error)
}

type entityDefinitionRepository struct {
	driver *neo4j.Driver
	repos  *RepositoryContainer
}

func NewEntityDefinitionRepository(driver *neo4j.Driver, repos *RepositoryContainer) EntityDefinitionRepository {
	return &entityDefinitionRepository{
		driver: driver,
		repos:  repos,
	}
}

func (r *entityDefinitionRepository) Create(tenant string, entity *entity.EntityDefinitionEntity) (any, error) {
	session := (*r.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(r.createFullEntityDefinitionInTxWork(tenant, entity))
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}

func (r *entityDefinitionRepository) FindAllByTenant(tenant string) (any, error) {
	session := (*r.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (:Tenant {name:$tenant})-[r:USES_ENTITY_DEFINITION]->(e:EntityDefinition) RETURN e, r`,
			map[string]any{
				"tenant": tenant,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
}

func (r *entityDefinitionRepository) createFullEntityDefinitionInTxWork(tenant string, entity *entity.EntityDefinitionEntity) func(tx neo4j.Transaction) (any, error) {
	return func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)-[:USES_ENTITY_DEFINITION {added:datetime({timezone: 'UTC'})}]->(e:EntityDefinition {
				  id: randomUUID(),
				  name: $name,
				  version: $version
			}) ON CREATE SET e.extends=$extends
			RETURN e`,
			map[string]any{
				"tenant":  tenant,
				"name":    entity.Name,
				"version": entity.Version,
				"extends": entity.Extends,
			})
		if err != nil {
			return nil, err
		}
		record, err := queryResult.Single()
		if err != nil {
			return nil, err
		}
		entityDefinitionId := utils.GetPropsFromNode(record.Values[0].(dbtype.Node))["id"].(string)
		for _, v := range entity.FieldSets {
			err := r.repos.FieldSetDefinitionRepository.createFieldSetDefinitionInTx(entityDefinitionId, v, tx)
			if err != nil {
				return nil, err
			}
		}
		for _, v := range entity.CustomFields {
			err := r.repos.CustomFieldDefinitionRepository.createCustomFieldDefinitionForEntityInTx(entityDefinitionId, v, tx)
			if err != nil {
				return nil, err
			}
		}
		return record.Values[0], nil
	}
}
