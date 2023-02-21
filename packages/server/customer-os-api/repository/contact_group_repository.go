package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type ContactGroupRepository interface {
	Create(ctx context.Context, tenant string, entity entity.ContactGroupEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tenant string, entity entity.ContactGroupEntity) (*dbtype.Node, error)
	GetPaginatedContactGroups(ctx context.Context, tenant string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
}

type contactGroupRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContactGroupRepository(driver *neo4j.DriverWithContext) ContactGroupRepository {
	return &contactGroupRepository{
		driver: driver,
	}
}

func (r *contactGroupRepository) Create(ctx context.Context, tenant string, entity entity.ContactGroupEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (g:ContactGroup {id: randomUUID()})-[:GROUP_BELONGS_TO_TENANT]->(t)" +
		" ON CREATE SET g.name=$name, " +
		"				g.source=$source, " +
		"				g.sourceOfTruth=$sourceOfTruth, " +
		" 				g.createdAt=datetime({timezone: 'UTC'}), " +
		"				g:%s " +
		" RETURN g"

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "ContactGroup_"+tenant),
			map[string]any{
				"tenant":        tenant,
				"name":          entity.Name,
				"source":        entity.Source,
				"sourceOfTruth": entity.SourceOfTruth,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *contactGroupRepository) Update(ctx context.Context, tenant string, entity entity.ContactGroupEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbNode, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (g:ContactGroup {id:$groupId})-[:GROUP_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			SET g.name=$name, 
				g.sourceOfTruth=$sourceOfTruth
			RETURN g`,
			map[string]any{
				"tenant":        tenant,
				"groupId":       entity.Id,
				"name":          entity.Name,
				"sourceOfTruth": entity.SourceOfTruth,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbNode.(*dbtype.Node), nil
}

func (r *contactGroupRepository) GetPaginatedContactGroups(ctx context.Context, tenant string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("g")
		countParams := map[string]any{
			"tenant": tenant,
		}
		utils.MergeMapToMap(filterParams, countParams)

		queryResult, err := tx.Run(ctx, fmt.Sprintf("MATCH (:Tenant {name:$tenant})<-[:GROUP_BELONGS_TO_TENANT]-(g:ContactGroup) %s RETURN count(g) as count", filterCypherStr),
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
			"MATCH (:Tenant {name:$tenant})<-[:GROUP_BELONGS_TO_TENANT]-(g:ContactGroup) "+
				" %s "+
				" RETURN g "+
				" %s "+
				" SKIP $skip LIMIT $limit", filterCypherStr, sorting.SortingCypherFragment("g")),
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
