package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type ContactGroupRepository interface {
	GetPaginatedContactGroups(session neo4j.Session, tenant string, skip, limit int, sorting *utils.Sorts) (*utils.DbNodesWithTotalCount, error)
}

type contactGroupRepository struct {
	driver *neo4j.Driver
	repos  *RepositoryContainer
}

func NewContactGroupRepository(driver *neo4j.Driver, repos *RepositoryContainer) ContactGroupRepository {
	return &contactGroupRepository{
		driver: driver,
		repos:  repos,
	}
}

func (r *contactGroupRepository) GetPaginatedContactGroups(session neo4j.Session, tenant string, skip, limit int, sorting *utils.Sorts) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (:Tenant {name:$tenant})<-[:GROUP_BELONGS_TO_TENANT]-(cg:ContactGroup) RETURN count(cg) as count`,
			map[string]any{
				"tenant": tenant,
			})
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single()
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(fmt.Sprintf(
			"MATCH (:Tenant {name:$tenant})<-[:GROUP_BELONGS_TO_TENANT]-(g:ContactGroup) "+
				" RETURN g "+
				" %s "+
				" SKIP $skip LIMIT $limit", sorting.SortingCypherFragment("g")),
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
