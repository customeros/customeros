package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type EntityTemplateRepository interface {
	Create(tenant string, entity *entity.EntityTemplateEntity) (any, error)
	FindAllByTenant(session neo4j.Session, tenant string) ([]*db.Record, error)
	FindAllByTenantAndExtends(session neo4j.Session, tenant, extends string) ([]*db.Record, error)
	FindByContactId(tenant string, contactId string) (any, error)
}

type entityTemplateRepository struct {
	driver       *neo4j.Driver
	repositories *Repositories
}

func NewEntityTemplateRepository(driver *neo4j.Driver, repositories *Repositories) EntityTemplateRepository {
	return &entityTemplateRepository{
		driver:       driver,
		repositories: repositories,
	}
}

func (r *entityTemplateRepository) Create(tenant string, entity *entity.EntityTemplateEntity) (any, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	queryResult, err := session.WriteTransaction(r.createFullEntityTemplateInTxWork(tenant, entity))
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}

func (r *entityTemplateRepository) FindAllByTenant(session neo4j.Session, tenant string) ([]*db.Record, error) {
	if result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (:Tenant {name:$tenant})<-[r:ENTITY_TEMPLATE_BELONGS_TO_TENANT]-(e:EntityTemplate) RETURN e`,
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

func (r *entityTemplateRepository) FindAllByTenantAndExtends(session neo4j.Session, tenant, extends string) ([]*db.Record, error) {
	if result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (:Tenant {name:$tenant})<-[r:ENTITY_TEMPLATE_BELONGS_TO_TENANT]-(e:EntityTemplate) 
					WHERE e.extends=$extends
				RETURN e`,
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

func (r *entityTemplateRepository) FindByContactId(tenant string, contactId string) (any, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})<-[r:ENTITY_TEMPLATE_BELONGS_TO_TENANT]->(e:EntityTemplate),
					(c)-[:IS_DEFINED_BY]->(e)
					RETURN e`,
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

func (r *entityTemplateRepository) createFullEntityTemplateInTxWork(tenant string, entity *entity.EntityTemplateEntity) func(tx neo4j.Transaction) (any, error) {
	return func(tx neo4j.Transaction) (any, error) {
		txResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[r:ENTITY_TEMPLATE_BELONGS_TO_TENANT]-(e:EntityTemplate {id: randomUUID()}) 
			ON CREATE SET e.extends=$extends, e.createdAt=datetime({timezone: 'UTC'}), e.name=$name, e.version=$version
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
		records, err := txResult.Collect()
		if err != nil {
			return nil, err
		}
		entityTemplateId := utils.GetPropsFromNode(records[0].Values[0].(dbtype.Node))["id"].(string)
		for _, v := range entity.FieldSets {
			err := r.repositories.FieldSetTemplateRepository.createFieldSetTemplateInTx(entityTemplateId, v, tx)
			if err != nil {
				return nil, err
			}
		}
		for _, v := range entity.CustomFields {
			err := r.repositories.CustomFieldTemplateRepository.createCustomFieldTemplateForEntityInTx(entityTemplateId, v, tx)
			if err != nil {
				return nil, err
			}
		}
		return records, nil
	}
}
