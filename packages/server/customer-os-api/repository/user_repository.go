package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type UserRepository interface {
	Create(tx neo4j.Transaction, tenant string, entity entity.UserEntity) (*dbtype.Node, error)
	Update(session neo4j.Session, tenant string, entity entity.UserEntity) (*dbtype.Node, error)
	FindUserByEmail(session neo4j.Session, tenant string, email string) (*dbtype.Node, error)
	FindOwnerForContact(tx neo4j.Transaction, tenant, contactId string) (*dbtype.Node, error)
	FindCreatorForNote(tx neo4j.Transaction, tenant, noteId string) (*dbtype.Node, error)
	GetPaginatedUsers(session neo4j.Session, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetById(session neo4j.Session, tenant, userId string) (*dbtype.Node, error)
	GetAllForConversation(session neo4j.Session, tenant, conversationId string) ([]*dbtype.Node, error)
}

type userRepository struct {
	driver *neo4j.Driver
}

func NewUserRepository(driver *neo4j.Driver) UserRepository {
	return &userRepository{
		driver: driver,
	}
}

func (r *userRepository) Create(tx neo4j.Transaction, tenant string, entity entity.UserEntity) (*dbtype.Node, error) {
	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (u:User {id: randomUUID()})-[:USER_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET u.firstName=$firstName, " +
		"				u.lastName=$lastName, " +
		"				u.createdAt=$now, " +
		"				u.updatedAt=$now, " +
		" 				u.source=$source, " +
		"				u.sourceOfTruth=$sourceOfTruth, " +
		"				u:%s" +
		" RETURN u"

	queryResult, err := tx.Run(fmt.Sprintf(query, "User_"+tenant),
		map[string]any{
			"tenant":        tenant,
			"firstName":     entity.FirstName,
			"lastName":      entity.LastName,
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
}

func (r *userRepository) Update(session neo4j.Session, tenant string, entity entity.UserEntity) (*dbtype.Node, error) {
	query := "MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" SET 	u.firstName=$firstName, " +
		"		u.lastName=$lastName, " +
		"		u.updatedAt=datetime({timezone: 'UTC'}), " +
		"		u.sourceOfTruth=$sourceOfTruth " +
		" RETURN u"

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(query,
			map[string]any{
				"userId":        entity.Id,
				"tenant":        tenant,
				"firstName":     entity.FirstName,
				"lastName":      entity.LastName,
				"sourceOfTruth": entity.SourceOfTruth,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *userRepository) FindUserByEmail(session neo4j.Session, tenant string, email string) (*dbtype.Node, error) {
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
			MATCH (:Email {email:$email})<-[:HAS]-(u:User),
			(u)-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
			RETURN u`,
			map[string]any{
				"tenant": tenant,
				"email":  email,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *userRepository) FindOwnerForContact(tx neo4j.Transaction, tenant, contactId string) (*dbtype.Node, error) {
	if queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})<-[:OWNS]-(u:User)
			RETURN u`,
		map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
		}); err != nil {
		return nil, err
	} else {
		dbRecords, err := queryResult.Collect()
		if err != nil {
			return nil, err
		} else if len(dbRecords) == 0 {
			return nil, nil
		} else {
			return utils.NodePtr(dbRecords[0].Values[0].(dbtype.Node)), nil
		}
	}
}

func (r *userRepository) FindCreatorForNote(tx neo4j.Transaction, tenant, noteId string) (*dbtype.Node, error) {
	if queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:CREATED]->(n:Note {id:$noteId})
			RETURN u`,
		map[string]any{
			"tenant": tenant,
			"noteId": noteId,
		}); err != nil {
		return nil, err
	} else {
		dbRecords, err := queryResult.Collect()
		if err != nil {
			return nil, err
		} else if len(dbRecords) == 0 {
			return nil, nil
		} else {
			return utils.NodePtr(dbRecords[0].Values[0].(dbtype.Node)), nil
		}
	}
}

func (r *userRepository) GetPaginatedUsers(session neo4j.Session, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("u")
		countParams := map[string]any{
			"tenant": tenant,
		}
		utils.MergeMapToMap(filterParams, countParams)
		queryResult, err := tx.Run(fmt.Sprintf("MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) %s RETURN count(u) as count", filterCypherStr),
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
			"MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) "+
				" %s "+
				" RETURN u "+
				" %s "+
				" SKIP $skip LIMIT $limit", filterCypherStr, sort.SortingCypherFragment("u")),
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

func (r *userRepository) GetById(session neo4j.Session, tenant, userId string) (*dbtype.Node, error) {
	dbRecord, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
			RETURN u`,
			map[string]any{
				"tenant": tenant,
				"userId": userId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Single()
	})
	return utils.NodePtr(dbRecord.(*db.Record).Values[0].(dbtype.Node)), err
}

func (r *userRepository) GetAllForConversation(session neo4j.Session, tenant, conversationId string) ([]*dbtype.Node, error) {
	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:PARTICIPATES]->(o:Conversation {id:$conversationId})
			RETURN u`,
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
