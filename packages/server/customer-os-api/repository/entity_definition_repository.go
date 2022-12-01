package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type EntityDefinitionRepository interface {
	Create(tenant string, entity *entity.EntityDefinitionEntity) (any, error)
	FindAllByTenant(session neo4j.Session, tenant string) ([]*db.Record, error)
	FindAllByTenantAndExtends(session neo4j.Session, tenant, extends string) ([]*db.Record, error)
	FindByContactId(tenant string, contactId string) (any, error)
}

type entityDefinitionRepository struct {
	driver *neo4j.Driver
	repos  *Repositories
}

func NewEntityDefinitionRepository(driver *neo4j.Driver, repos *Repositories) EntityDefinitionRepository {
	return &entityDefinitionRepository{
		driver: driver,
		repos:  repos,
	}
}

func (r *entityDefinitionRepository) Create(tenant string, entity *entity.EntityDefinitionEntity) (any, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	queryResult, err := session.WriteTransaction(r.createFullEntityDefinitionInTxWork(tenant, entity))
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}

func (r *entityDefinitionRepository) FindAllByTenant(session neo4j.Session, tenant string) ([]*db.Record, error) {
	if result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (:Tenant {name:$tenant})-[r:USES_ENTITY_DEFINITION]->(e:EntityDefinition) RETURN e, r`,
			map[string]any{
				"tenant": tenant,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	}); err != nil {
		return nil, err
	} else {
		return result.([]*db.Record), nil
	}
}

func (r *entityDefinitionRepository) FindAllByTenantAndExtends(session neo4j.Session, tenant, extends string) ([]*db.Record, error) {
	if result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (:Tenant {name:$tenant})-[r:USES_ENTITY_DEFINITION]->(e:EntityDefinition) 
					WHERE e.extends=$extends
				RETURN e, r`,
			map[string]any{
				"tenant":  tenant,
				"extends": extends,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	}); err != nil {
		return nil, err
	} else {
		return result.([]*db.Record), nil
	}
}

func (r *entityDefinitionRepository) FindByContactId(tenant string, contactId string) (any, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})-[r:USES_ENTITY_DEFINITION]->(e:EntityDefinition),
					(c)-[:IS_DEFINED_BY]->(e)
					RETURN e, r`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
}

func (r *entityDefinitionRepository) createFullEntityDefinitionInTxWork(tenant string, entity *entity.EntityDefinitionEntity) func(tx neo4j.Transaction) (any, error) {
	return func(tx neo4j.Transaction) (any, error) {
		txResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)-[r:USES_ENTITY_DEFINITION {added:datetime({timezone: 'UTC'})}]->(e:EntityDefinition {
				  id: randomUUID(),
				  name: $name,
				  version: $version
			}) ON CREATE SET e.extends=$extends
			RETURN e, r`,
			map[string]any{
				"tenant":  tenant,
				"name":    entity.Name,
				"version": entity.Version,
				"extends": entity.Extends,
			})
		if err != nil {
			return nil, err
		}
		records, err := txResult.Collect()
		if err != nil {
			return nil, err
		}
		entityDefinitionId := utils.GetPropsFromNode(records[0].Values[0].(dbtype.Node))["id"].(string)
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
		return records, nil
	}
}
