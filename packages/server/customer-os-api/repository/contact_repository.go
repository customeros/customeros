package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"time"
)

type ContactRepository interface {
	Create(tx neo4j.Transaction, tenant string, newContact entity.ContactEntity) (*dbtype.Node, error)
	Update(tx neo4j.Transaction, tenant, contactId string, contactDtls *entity.ContactEntity) (*dbtype.Node, error)
	Delete(session neo4j.Session, tenant, contactId string) error
	SetOwner(tx neo4j.Transaction, tenant, contactId, userId string) error
	RemoveOwner(tx neo4j.Transaction, tenant, contactId string) error
	LinkWithEntityTemplateInTx(tx neo4j.Transaction, tenant, contactId, entityTemplateId string) error
	LinkWithContactTypeInTx(tx neo4j.Transaction, tenant, contactId, contactTypeId string) error
	UnlinkFromContactTypesInTx(tx neo4j.Transaction, tenant, contactId string) error
	GetPaginatedContacts(session neo4j.Session, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetPaginatedContactsForContactGroup(session neo4j.Session, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort, contactGroupId string) (*utils.DbNodesWithTotalCount, error)
	GetAllForConversation(session neo4j.Session, tenant, conversationId string) ([]*dbtype.Node, error)
}

type contactRepository struct {
	driver *neo4j.Driver
}

func NewContactRepository(driver *neo4j.Driver) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
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

func (r *contactRepository) Create(tx neo4j.Transaction, tenant string, newContact entity.ContactEntity) (*dbtype.Node, error) {
	var createdAt time.Time
	createdAt = time.Now().UTC()
	if newContact.CreatedAt != nil {
		createdAt = *newContact.CreatedAt
	}

	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (c:Contact {id:randomUUID()})-[:CONTACT_BELONGS_TO_TENANT]->(t) ON CREATE SET" +
		" c.title=$title, " +
		" c.firstName=$firstName, " +
		" c.lastName=$lastName, " +
		" c.readonly=$readonly, " +
		" c.label=$label, " +
		" c.createdAt=$createdAt, " +
		" c:%s " +
		" RETURN c"

	if queryResult, err := tx.Run(fmt.Sprintf(query, "Contact_"+tenant),
		map[string]interface{}{
			"tenant":    tenant,
			"title":     newContact.Title,
			"firstName": newContact.FirstName,
			"lastName":  newContact.LastName,
			"readonly":  newContact.Readonly,
			"label":     newContact.Label,
			"createdAt": createdAt,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}
}

func (r *contactRepository) Update(tx neo4j.Transaction, tenant, contactId string, contactDtls *entity.ContactEntity) (*dbtype.Node, error) {
	if queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			SET c.firstName=$firstName,
				c.lastName=$lastName,
				c.label=$label,
				c.title=$title,
			    c.readonly=$readonly
			RETURN c`,
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
			"firstName": contactDtls.FirstName,
			"lastName":  contactDtls.LastName,
			"label":     contactDtls.Label,
			"title":     contactDtls.Title,
			"readonly":  contactDtls.Readonly,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}
}

func (r *contactRepository) LinkWithEntityTemplateInTx(tx neo4j.Transaction, tenant, contactId, entityTemplateId string) error {
	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})<-[:ENTITY_TEMPLATE_BELONGS_TO_TENANT]-(e:EntityTemplate {id:$entityTemplateId})
			WHERE e.extends=$extends
			MERGE (c)-[r:IS_DEFINED_BY]->(e)
			RETURN r`,
		map[string]any{
			"entityTemplateId": entityTemplateId,
			"contactId":        contactId,
			"tenant":           tenant,
			"extends":          model.EntityTemplateExtensionContact,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single()
	return err
}

func (r *contactRepository) LinkWithContactTypeInTx(tx neo4j.Transaction, tenant, contactId, contactTypeId string) error {
	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})<-[:CONTACT_TYPE_BELONGS_TO_TENANT]-(e:ContactType {id:$contactTypeId})
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

func (r *contactRepository) GetPaginatedContacts(session neo4j.Session, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("c")
		countParams := map[string]any{
			"tenant": tenant,
		}
		utils.MergeMapToMap(filterParams, countParams)
		queryResult, err := tx.Run(fmt.Sprintf("MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) %s RETURN count(c) as count", filterCypherStr),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single()
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}
		utils.MergeMapToMap(filterParams, params)

		queryResult, err = tx.Run(fmt.Sprintf(
			"MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) "+
				" %s "+
				" RETURN c "+
				" %s "+
				" SKIP $skip LIMIT $limit", filterCypherStr, sort.SortingCypherFragment("c")),
			params)
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

func (r *contactRepository) GetPaginatedContactsForContactGroup(session neo4j.Session, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort, contactGroupId string) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("c")
		countParams := map[string]any{
			"tenant":         tenant,
			"contactGroupId": contactGroupId,
		}
		utils.MergeMapToMap(filterParams, countParams)
		queryResult, err := tx.Run(fmt.Sprintf("MATCH (t:Tenant {name:$tenant})<-[:GROUP_BELONGS_TO_TENANT]-(g:ContactGroup {id:$contactGroupId}), "+
			" (g)<-[:BELONGS_TO_GROUP]-(c:Contact) %s "+
			" RETURN count(c) as count", filterCypherStr),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single()
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant":         tenant,
			"skip":           skip,
			"limit":          limit,
			"contactGroupId": contactGroupId,
		}
		utils.MergeMapToMap(filterParams, params)

		queryResult, err = tx.Run(fmt.Sprintf(
			"MATCH (t:Tenant {name:$tenant})<-[:GROUP_BELONGS_TO_TENANT]-(g:ContactGroup {id:$contactGroupId}), "+
				" (g)<-[:BELONGS_TO_GROUP]-(c:Contact) "+
				" %s "+
				" RETURN c "+
				" %s "+
				" SKIP $skip LIMIT $limit", filterCypherStr, sort.SortingCypherFragment("c")),
			params)
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

func (r *contactRepository) Delete(session neo4j.Session, tenant, contactId string) error {
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (c)-[:HAS_PROPERTY]->(f:CustomField)
			OPTIONAL MATCH (c)-[:CALLED_AT]->(p:PhoneNumber)
			OPTIONAL MATCH (c)-[:EMAILED_AT]->(e:Email)
			OPTIONAL MATCH (c)-[:LOCATED_AT]->(a:Address)
			OPTIONAL MATCH (c)-[:HAS_COMPLEX_PROPERTY]->(fs:FieldSet)
			OPTIONAL MATCH (c)-[:HAS_ROLE]->(r:Role)
            DETACH DELETE p, e, f, fs, a, r, c`,
			map[string]interface{}{
				"contactId": contactId,
				"tenant":    tenant,
			})
		return nil, err
	})
	return err
}

// FIXME alexb test
func (r *contactRepository) GetAllForConversation(session neo4j.Session, tenant, conversationId string) ([]*dbtype.Node, error) {
	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:PARTICIPATES]->(o:Conversation {id:$conversationId})
			RETURN c`,
			map[string]any{
				"tenant":         tenant,
				"conversationId": conversationId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	if err != nil {
		return nil, err
	}
	dbNodes := []*dbtype.Node{}
	for _, v := range dbRecords.([]*neo4j.Record) {
		if v.Values[0] != nil {
			dbNodes = append(dbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
		}
	}
	return dbNodes, err
}
