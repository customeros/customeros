package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type UserRepository interface {
	FindOwnerForContact(tx neo4j.Transaction, tenant, contactId string) (*dbtype.Node, error)
	GetPaginatedUsers(session neo4j.Session, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetById(session neo4j.Session, tenant, userId string) (*dbtype.Node, error)
}

type userRepository struct {
	driver *neo4j.Driver
	repos  *RepositoryContainer
}

func NewUserRepository(driver *neo4j.Driver, repos *RepositoryContainer) UserRepository {
	return &userRepository{
		driver: driver,
		repos:  repos,
	}
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