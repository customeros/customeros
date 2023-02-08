package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type QueryRepository interface {
	GetOrganizationsAndContacts(session neo4j.Session, tenant string, skip int, limit int, searchTerm *string) (*utils.PairDbNodesWithTotalCount, error)
}

type queryRepository struct {
	driver *neo4j.Driver
}

func NewQueryRepository(driver *neo4j.Driver) QueryRepository {
	return &queryRepository{
		driver: driver,
	}
}

func createCypherFilter(propertyName string, searchTerm string) *utils.CypherFilter {
	filter := utils.CypherFilter{}
	filter.Details = new(utils.CypherFilterItem)
	filter.Details.NodeProperty = propertyName
	filter.Details.Value = &searchTerm
	filter.Details.ComparisonOperator = utils.CONTAINS
	filter.Details.SupportCaseSensitive = true
	return &filter
}

func (r *queryRepository) GetOrganizationsAndContacts(session neo4j.Session, tenant string, skip int, limit int, searchTerm *string) (*utils.PairDbNodesWithTotalCount, error) {
	result := new(utils.PairDbNodesWithTotalCount)

	contactFilterCypher, contactFilterParams := "1=1", make(map[string]interface{})
	organizationFilterCypher, organizationFilterParams := "1=1", make(map[string]interface{})

	//region contact filters
	if searchTerm != nil {
		contactFilter := new(utils.CypherFilter)
		contactFilter.Negate = false
		contactFilter.LogicalOperator = utils.OR
		contactFilter.Filters = make([]*utils.CypherFilter, 0)

		contactFilter.Filters = append(contactFilter.Filters, createCypherFilter("firstName", *searchTerm))
		contactFilter.Filters = append(contactFilter.Filters, createCypherFilter("lastName", *searchTerm))

		contactFilterCypher, contactFilterParams = contactFilter.BuildCypherFilterFragmentWithParamName("c", "c_param_")
	}

	//endregion

	//region organization filters
	if searchTerm != nil {

		organizationFilter := new(utils.CypherFilter)
		organizationFilter.Negate = false
		organizationFilter.LogicalOperator = utils.OR
		organizationFilter.Filters = make([]*utils.CypherFilter, 0)

		organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("name", *searchTerm))

		organizationFilterCypher, organizationFilterParams = organizationFilter.BuildCypherFilterFragmentWithParamName("o", "o_param_")
	}

	//endregion

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		params := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}
		utils.MergeMapToMap(contactFilterParams, params)
		utils.MergeMapToMap(organizationFilterParams, params)

		queryResult, err := tx.Run(fmt.Sprintf(`
		CALL {
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  MATCH (t)--(c:Contact)
		  MATCH (o)--(c)
		  WHERE (%s) OR (%s)
		  RETURN count(o) as t
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  WHERE NOT (o)--(:Contact) AND (%s)
		  RETURN count(o) as t
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(c:Contact)
		  WHERE NOT (c)--(:Organization) AND (%s)
		  RETURN count(c) as t
		}
		RETURN sum(t)`, contactFilterCypher, organizationFilterCypher, organizationFilterCypher, contactFilterCypher),
			params)
		if err != nil {
			return nil, err
		}
		countRecord, err := queryResult.Single()
		if err != nil {
			return nil, err
		}
		result.Count = countRecord.Values[0].(int64)

		if queryResult, err := tx.Run(fmt.Sprintf(`
		CALL {
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  MATCH (t)--(c:Contact)
		  MATCH (o)--(c)
		  WHERE (%s) OR (%s)
		  RETURN o, c
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  WHERE NOT (o)--(:Contact) AND (%s)
		  RETURN o, null as c
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(c:Contact)
		  WHERE NOT (c)--(:Organization) AND (%s)
		  RETURN null as o, c
		}
		RETURN o, c
		SKIP $skip LIMIT $limit`, organizationFilterCypher, contactFilterCypher, organizationFilterCypher, contactFilterCypher),
			params); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	if err != nil {
		return nil, err
	}

	for _, v := range dbRecords.([]*neo4j.Record) {
		pair := new(utils.Pair[*dbtype.Node, *dbtype.Node])
		if v.Values[0] != nil {
			node := v.Values[0].(dbtype.Node)
			pair.First = &node
		}
		if v.Values[1] != nil {
			node := v.Values[1].(dbtype.Node)
			pair.Second = &node
		}

		result.Pairs = append(result.Pairs, pair)
	}
	return result, err
}

func max(a int64, b int64) int64 {
	if a < b {
		return b
	}
	return a
}
