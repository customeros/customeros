package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type ContactRepository interface {
	Create(tx neo4j.Transaction, tenant string, newContact entity.ContactEntity) (*dbtype.Node, error)
	SetOwner(tx neo4j.Transaction, tenant, contactId, userId string) error
	RemoveOwner(tx neo4j.Transaction, tenant, contactId string) error
	LinkWithEntityDefinitionInTx(tx neo4j.Transaction, tenant, contactId, entityDefinitionId string) error
	LinkWithContactTypeInTx(tx neo4j.Transaction, tenant, contactId, contactTypeId string) error
	UnlinkFromContactTypesInTx(tx neo4j.Transaction, tenant, contactId string) error
	GetPaginatedContacts(session neo4j.Session, tenant string, skip, limit int, sorting *utils.Sortings) (*utils.DbNodesWithTotalCount, error)
}

type contactRepository struct {
	driver *neo4j.Driver
	repos  *RepositoryContainer
}

func (r *contactRepository) SetOwner(tx neo4j.Transaction, tenant, contactId, userId string) error {
	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t)
			MERGE (u)-[r:OWNS]->(c)
			RETURN r`,
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
			"userId":    userId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single()
	return err
}

func (r *contactRepository) RemoveOwner(tx neo4j.Transaction, tenant, contactId string) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(c)<-[r:OWNS]-()
			DELETE r`,
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
		})
	return err
}

func NewContactRepository(driver *neo4j.Driver, repos *RepositoryContainer) ContactRepository {
	return &contactRepository{
		driver: driver,
		repos:  repos,
	}
}

func (r *contactRepository) Create(tx neo4j.Transaction, tenant string, newContact entity.ContactEntity) (*dbtype.Node, error) {
	if queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})
			CREATE (c:Contact {
				  id: randomUUID(),
				  title: $title,
				  firstName: $firstName,
				  lastName: $lastName,
				  label: $label,
				  notes: $notes,
                  createdAt :datetime({timezone: 'UTC'})
			})-[:CONTACT_BELONGS_TO_TENANT]->(t)
			RETURN c`,
		map[string]interface{}{
			"tenant":    tenant,
			"firstName": newContact.FirstName,
			"lastName":  newContact.LastName,
			"label":     newContact.Label,
			"title":     newContact.Title,
			"notes":     newContact.Notes,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNodePtr(queryResult, err)
	}
}

func (r *contactRepository) LinkWithEntityDefinitionInTx(tx neo4j.Transaction, tenant, contactId, entityDefinitionId string) error {
	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})-[:USES_ENTITY_DEFINITION]->(e:EntityDefinition {id:$entityDefinitionId})
			WHERE e.extends=$extends
			MERGE (c)-[r:IS_DEFINED_BY]->(e)
			RETURN r`,
		map[string]any{
			"entityDefinitionId": entityDefinitionId,
			"contactId":          contactId,
			"tenant":             tenant,
			"extends":            model.EntityDefinitionExtensionContact,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single()
	return err
}

func (r *contactRepository) LinkWithContactTypeInTx(tx neo4j.Transaction, tenant, contactId, contactTypeId string) error {
	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})-[:USES_CONTACT_TYPE]->(e:ContactType {id:$contactTypeId})
			MERGE (c)-[r:IS_OF_TYPE]->(e)
			RETURN r`,
		map[string]any{
			"tenant":        tenant,
			"contactId":     contactId,
			"contactTypeId": contactTypeId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single()
	return err
}

func (r *contactRepository) UnlinkFromContactTypesInTx(tx neo4j.Transaction, tenant, contactId string) error {
	if _, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				(c)-[r:IS_OF_TYPE]->(o:ContactType)
			DELETE r`,
		map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
		}); err != nil {
		return err
	}
	return nil
}

func (r *contactRepository) GetPaginatedContacts(session neo4j.Session, tenant string, skip, limit int, sorting *utils.Sortings) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) RETURN count(c) as count`,
			map[string]any{
				"tenant": tenant,
			})
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single()
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(fmt.Sprintf(
			"MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) "+
				" RETURN c "+
				" %s "+
				" SKIP $skip LIMIT $limit", sorting.CypherFragment("c")),
			map[string]any{
				"tenant": tenant,
				"skip":   skip,
				"limit":  limit,
			})
		return queryResult.Collect()
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}
