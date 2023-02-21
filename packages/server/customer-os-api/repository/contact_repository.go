package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
	"time"
)

type ContactRepository interface {
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newContact entity.ContactEntity, source, sourceOfTruth entity.DataSource) (*dbtype.Node, error)
	Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, contactDtls *entity.ContactEntity) (*dbtype.Node, error)
	Delete(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) error
	SetOwner(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, userId string) error
	RemoveOwner(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string) error
	LinkWithEntityTemplateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, entityTemplateId string) error
	GetPaginatedContacts(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetPaginatedContactsForContactGroup(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort, contactGroupId string) (*utils.DbNodesWithTotalCount, error)
	GetPaginatedContactsForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetAllForConversation(ctx context.Context, session neo4j.SessionWithContext, tenant, conversationId string) ([]*dbtype.Node, error)
	GetContactForRole(ctx context.Context, session neo4j.SessionWithContext, tenant, roleId string) (*dbtype.Node, error)
	AddTag(ctx context.Context, tenant, contactId, tagId string) (*dbtype.Node, error)
	RemoveTag(ctx context.Context, tenant, contactId, tagId string) (*dbtype.Node, error)
}

type contactRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContactRepository(driver *neo4j.DriverWithContext) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) SetOwner(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, userId string) error {
	queryResult, err := tx.Run(ctx, `
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
	_, err = queryResult.Single(ctx)
	return err
}

func (r *contactRepository) RemoveOwner(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string) error {
	_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(c)<-[r:OWNS]-()
			DELETE r`,
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
		})
	return err
}

func (r *contactRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newContact entity.ContactEntity, source, sourceOfTruth entity.DataSource) (*dbtype.Node, error) {
	var createdAt time.Time
	createdAt = utils.Now()
	if newContact.CreatedAt != nil {
		createdAt = *newContact.CreatedAt
	}

	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (c:Contact {id:randomUUID()})-[:CONTACT_BELONGS_TO_TENANT]->(t) ON CREATE SET " +
		" c.title=$title, " +
		" c.firstName=$firstName, " +
		" c.lastName=$lastName, " +
		" c.createdAt=$createdAt, " +
		" c.updatedAt=$createdAt, " +
		" c.source=$source, " +
		" c.sourceOfTruth=$sourceOfTruth, " +
		" c:Contact_%s " +
		" RETURN c"

	if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"title":         newContact.Title,
			"firstName":     newContact.FirstName,
			"lastName":      newContact.LastName,
			"source":        source,
			"sourceOfTruth": sourceOfTruth,
			"createdAt":     createdAt,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}

func (r *contactRepository) Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, contactDtls *entity.ContactEntity) (*dbtype.Node, error) {
	if queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			SET c.firstName=$firstName,
				c.lastName=$lastName,
				c.title=$title,
				c.updatedAt=datetime({timezone: 'UTC'})
			RETURN c`,
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
			"firstName": contactDtls.FirstName,
			"lastName":  contactDtls.LastName,
			"title":     contactDtls.Title,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}

func (r *contactRepository) LinkWithEntityTemplateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, entityTemplateId string) error {
	queryResult, err := tx.Run(ctx, `
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
	_, err = queryResult.Single(ctx)
	return err
}

func (r *contactRepository) GetPaginatedContacts(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("c")
		countParams := map[string]any{
			"tenant": tenant,
		}
		utils.MergeMapToMap(filterParams, countParams)
		queryResult, err := tx.Run(ctx, fmt.Sprintf("MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) %s RETURN count(c) as count", filterCypherStr),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}
		utils.MergeMapToMap(filterParams, params)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
			"MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) "+
				" %s "+
				" RETURN c "+
				" %s "+
				" SKIP $skip LIMIT $limit", filterCypherStr, sort.SortingCypherFragment("c")),
			params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *contactRepository) GetPaginatedContactsForContactGroup(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort, contactGroupId string) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("c")
		countParams := map[string]any{
			"tenant":         tenant,
			"contactGroupId": contactGroupId,
		}
		utils.MergeMapToMap(filterParams, countParams)
		queryResult, err := tx.Run(ctx, fmt.Sprintf("MATCH (t:Tenant {name:$tenant})<-[:GROUP_BELONGS_TO_TENANT]-(g:ContactGroup {id:$contactGroupId}), "+
			" (g)<-[:BELONGS_TO_GROUP]-(c:Contact) %s "+
			" RETURN count(c) as count", filterCypherStr),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant":         tenant,
			"skip":           skip,
			"limit":          limit,
			"contactGroupId": contactGroupId,
		}
		utils.MergeMapToMap(filterParams, params)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
			"MATCH (t:Tenant {name:$tenant})<-[:GROUP_BELONGS_TO_TENANT]-(g:ContactGroup {id:$contactGroupId}), "+
				" (g)<-[:BELONGS_TO_GROUP]-(c:Contact) "+
				" %s "+
				" RETURN c "+
				" %s "+
				" SKIP $skip LIMIT $limit", filterCypherStr, sort.SortingCypherFragment("c")),
			params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *contactRepository) GetPaginatedContactsForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("c")
		countParams := map[string]any{
			"tenant":         tenant,
			"organizationId": organizationId,
		}
		utils.MergeMapToMap(filterParams, countParams)
		queryResult, err := tx.Run(ctx, fmt.Sprintf("MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})--(c:Contact) "+
			" %s "+
			" RETURN count(c) as count", filterCypherStr),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant":         tenant,
			"skip":           skip,
			"limit":          limit,
			"organizationId": organizationId,
		}
		utils.MergeMapToMap(filterParams, params)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
			"MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})--(c:Contact) "+
				" %s "+
				" RETURN c "+
				" %s "+
				" SKIP $skip LIMIT $limit", filterCypherStr, sort.SortingCypherFragment("c")),
			params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *contactRepository) Delete(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) error {
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (c)-[:HAS_PROPERTY]->(f:CustomField)
			OPTIONAL MATCH (c)-[:HAS_COMPLEX_PROPERTY]->(fs:FieldSet)
			OPTIONAL MATCH (c)-[:WORKS_AS]->(j:JobRole)
            DETACH DELETE f, fs, j, c`,
			map[string]interface{}{
				"contactId": contactId,
				"tenant":    tenant,
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) GetAllForConversation(ctx context.Context, session neo4j.SessionWithContext, tenant, conversationId string) ([]*dbtype.Node, error) {
	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:PARTICIPATES]->(o:Conversation {id:$conversationId})
			RETURN c`,
			map[string]any{
				"tenant":         tenant,
				"conversationId": conversationId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
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

func (r *contactRepository) GetContactForRole(ctx context.Context, session neo4j.SessionWithContext, tenant, roleId string) (*dbtype.Node, error) {
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:JobRole {id:$roleId})<-[:WORKS_AS]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			RETURN c`,
			map[string]any{
				"tenant": tenant,
				"roleId": roleId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *contactRepository) AddTag(ctx context.Context, tenant, contactId, tagId string) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant}), " +
		" (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t), " +
		" (tag:Tag {id:$tagId})-[:TAG_BELONGS_TO_TENANT]->(t) " +
		" MERGE (c)-[rel:TAGGED]->(tag) " +
		" ON CREATE SET rel.taggedAt=$now " +
		" RETURN c"

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
				"tagId":     tagId,
				"now":       utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *contactRepository) RemoveTag(ctx context.Context, tenant, contactId, tagId string) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant}), " +
		" (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t), " +
		" (tag:Tag {id:$tagId})-[:TAG_BELONGS_TO_TENANT]->(t) " +
		" OPTIONAL MATCH (c)-[rel:TAGGED]->(tag) " +
		" DELETE rel " +
		" RETURN c"

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
				"tagId":     tagId,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
