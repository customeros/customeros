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
	contactEmailFilterCypher, contactEmailFilterParams := "1=1", make(map[string]interface{})
	organizationEmailFilterCypher, organizationEmailFilterParams := "1=1", make(map[string]interface{})

	//region contact filters
	if searchTerm != nil {
		contactFilter := new(utils.CypherFilter)
		contactFilter.Negate = false
		contactFilter.LogicalOperator = utils.OR
		contactFilter.Filters = make([]*utils.CypherFilter, 0)

		contactFilter.Filters = append(contactFilter.Filters, createCypherFilter("firstName", *searchTerm))
		contactFilter.Filters = append(contactFilter.Filters, createCypherFilter("lastName", *searchTerm))

		contactFilterCypher, contactFilterParams = contactFilter.BuildCypherFilterFragmentWithParamName("c", "c_param_")

		contactEmailFilter := new(utils.CypherFilter)
		contactEmailFilter.Negate = false
		contactEmailFilter.LogicalOperator = utils.OR
		contactEmailFilter.Filters = make([]*utils.CypherFilter, 0)

		contactEmailFilter.Filters = append(contactEmailFilter.Filters, createCypherFilter("email", *searchTerm))
		contactEmailFilterCypher, contactEmailFilterParams = contactEmailFilter.BuildCypherFilterFragmentWithParamName("ce", "ce_param_")

		fmt.Sprintf(contactFilterCypher)
		fmt.Sprintf(contactEmailFilterCypher)

		contactFilterCypher = fmt.Sprintf("(%s) OR (%s)", contactFilterCypher, contactEmailFilterCypher)
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

		organizationEmailFilter := new(utils.CypherFilter)
		organizationEmailFilter.Negate = false
		organizationEmailFilter.LogicalOperator = utils.OR
		organizationEmailFilter.Filters = make([]*utils.CypherFilter, 0)

		organizationEmailFilter.Filters = append(organizationEmailFilter.Filters, createCypherFilter("email", *searchTerm))
		organizationEmailFilterCypher, organizationEmailFilterParams = organizationEmailFilter.BuildCypherFilterFragmentWithParamName("oe", "oe_param_")

		fmt.Sprintf(organizationFilterCypher)
		fmt.Sprintf(organizationEmailFilterCypher)

		organizationFilterCypher = fmt.Sprintf("(%s) OR (%s)", organizationFilterCypher, organizationEmailFilterCypher)
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
		utils.MergeMapToMap(contactEmailFilterParams, params)
		utils.MergeMapToMap(organizationEmailFilterParams, params)

		queryResult, err := tx.Run(fmt.Sprintf(`
		 CALL {
		  MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[r1:HAS]->(oe:Email) 
		  MATCH (t)--(c:Contact)-[r2:HAS]->(ce:Email) 
		  MATCH (o)-[rel]-(c)
		  WHERE (%s) OR (%s)
		  RETURN rel
		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[rel:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[r1:HAS]->(oe:Email) 
		  WHERE NOT (o)--(:Contact) AND (%s)
		  RETURN rel
		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[rel:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[r1:HAS]->(ce:Email) 
		  WHERE NOT (c)--(:Organization) AND (%s)
		  RETURN rel
		}
		RETURN count(rel)`, contactFilterCypher, organizationFilterCypher, organizationFilterCypher, contactFilterCypher),
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
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)-[r1:HAS]->(oe:Email) 
		  MATCH (t)--(c:Contact)-[r2:HAS]->(ce:Email) 
		  MATCH (o)--(c)
		  WHERE (%s) OR (%s)
		  RETURN o, c
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)-[r1:HAS]->(oe:Email) 
		  WHERE NOT (o)--(:Contact) AND (%s)
		  RETURN o, null as c
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(c:Contact)-[r1:HAS]->(ce:Email) 
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
