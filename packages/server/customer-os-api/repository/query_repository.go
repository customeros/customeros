package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type QueryRepository interface {
	GetOrganizationsAndContacts(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip int, limit int, searchTerm *string) (*utils.PairDbNodesWithTotalCount, error)
}

type queryRepository struct {
	driver *neo4j.DriverWithContext
}

func NewQueryRepository(driver *neo4j.DriverWithContext) QueryRepository {
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

func (r *queryRepository) GetOrganizationsAndContacts(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip int, limit int, searchTerm *string) (*utils.PairDbNodesWithTotalCount, error) {
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
		contactFilter.Filters = append(contactFilter.Filters, createCypherFilter("name", *searchTerm))

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

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}
		utils.MergeMapToMap(contactFilterParams, params)
		utils.MergeMapToMap(organizationFilterParams, params)

		if searchTerm != nil {
			params["email"] = searchTerm
			params["location"] = searchTerm
		}

		//region count query
		countQuery := fmt.Sprintf(`CALL {`)

		//fetch organizations and contacts + filters on their properties
		countQuery = countQuery + fmt.Sprintf(`
          MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  MATCH (t)--(c:Contact)
		  MATCH (o)-[rel]-(c)
		  WHERE (%s) OR (%s)
		  RETURN rel`, organizationFilterCypher, contactFilterCypher)
		//fetch organizations and contacts with filters on their emails
		if searchTerm != nil {
			countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)-[:HAS]->(e:Email)
		  MATCH (t)--(c:Contact)
		  MATCH (o)-[rel]-(c)
		  WHERE (e.email CONTAINS $email OR e.rawEmail CONTAINS $email)
		  RETURN rel

		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email)
		  MATCH (o)-[rel]-(c)
		  WHERE (e.email CONTAINS $email or e.rawEmail CONTAINS $email)
		  RETURN rel`)
		}
		//fetch organizations and contacts with filters on their locations
		if searchTerm != nil {
			countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)--(l:Location)
		  MATCH (t)--(c:Contact)
		  MATCH (o)-[rel]-(c)
		  WHERE l.country CONTAINS $location OR l.region CONTAINS $location OR l.locality CONTAINS $location
		  RETURN rel

		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  MATCH (t)--(c:Contact)--(l:Location)
		  MATCH (o)-[rel]-(c)
		  WHERE l.country CONTAINS $location OR l.region CONTAINS $location OR l.locality CONTAINS $location
		  RETURN rel`)
		}

		//fetch organizations without contacts + filters on their properties
		countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})-[rel:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
		  WHERE NOT (o)--(:Contact) AND (%s)	
		  RETURN rel`, organizationFilterCypher)
		//fetch organizations without contacts with filters on their emails
		if searchTerm != nil {
			countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})-[rel:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS]->(e:Email)
		  WHERE NOT (o)--(:Contact)	AND (e.email CONTAINS $email OR e.rawEmail CONTAINS $email)
		  RETURN rel`)
		}
		//fetch organizations without contacts with filters on their locations
		if searchTerm != nil {
			countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})-[rel:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)--(l:Location)
		  WHERE NOT (o)--(:Contact)	AND (l.country CONTAINS $location OR l.region CONTAINS $location OR l.locality CONTAINS $location)
		  RETURN rel`)
		}

		//fetch contacts without organizations + filters on their properties
		countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[rel:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
		  WHERE NOT (c)--(:Organization) AND (%s)
		  RETURN rel`, contactFilterCypher)
		//fetch contacts without organizations with filters on their emails
		if searchTerm != nil {
			countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[rel:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email)
		  WHERE NOT (c)--(:Organization) AND (e.email CONTAINS $email OR e.rawEmail CONTAINS $email)
		  RETURN rel`)
		}
		//fetch contacts without organizations with filters on their locations
		if searchTerm != nil {
			countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[rel:CONTACT_BELONGS_TO_TENANT]-(c:Contact)--(l:Location)
		  WHERE NOT (c)--(:Organization) AND (l.country CONTAINS $location OR l.region CONTAINS $location OR l.locality CONTAINS $location)
		  RETURN rel`)
		}

		countQuery = countQuery + fmt.Sprintf(`} RETURN count(rel)`)

		countQueryResult, err := tx.Run(ctx, countQuery, params)
		if err != nil {
			return nil, err
		}

		countRecord, err := countQueryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		result.Count = countRecord.Values[0].(int64)

		//endregion

		//region query to fetch data
		query := fmt.Sprintf(`CALL {`)

		//fetch organizations and contacts + filters on their properties
		query = query + fmt.Sprintf(`
          MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  MATCH (t)--(c:Contact)
		  MATCH (o)--(c)
		  WHERE (%s) OR (%s)
		  RETURN o, c`, organizationFilterCypher, contactFilterCypher)
		//fetch organizations and contacts with filters on their emails
		if searchTerm != nil {
			query = query + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)-[:HAS]->(e:Email)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
		  MATCH (o)--(c)
		  WHERE (e.email CONTAINS $email OR e.rawEmail CONTAINS $email)
		  RETURN o, c

		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email)
		  MATCH (o)--(c)
		  WHERE (e.email CONTAINS $email OR e.rawEmail CONTAINS $email)
		  RETURN o, c`)
		}
		//fetch organizations and contacts with filters on their locations
		if searchTerm != nil {
			query = query + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)--(l:Location)
		  MATCH (t)--(c:Contact)
		  MATCH (o)--(c)
		  WHERE l.country CONTAINS $location OR l.region CONTAINS $location OR l.locality CONTAINS $location
		  RETURN o, c

		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  MATCH (t)--(c:Contact)--(l:Location)
		  MATCH (o)--(c)
		  WHERE l.country CONTAINS $location OR l.region CONTAINS $location OR l.locality CONTAINS $location
		  RETURN o, c`)
		}

		//fetch organizations without contacts + filters on their properties
		query = query + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  WHERE NOT (o)--(:Contact) AND (%s)	
		  RETURN o, null as c`, organizationFilterCypher)
		//fetch organizations without contacts with filters on their emails
		if searchTerm != nil {
			query = query + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)-[:HAS]->(e:Email)
		  WHERE NOT (o)--(:Contact) AND (e.email CONTAINS $email AND e.rawEmail CONTAINS $email)
		  RETURN o, null as c`)
		}
		//fetch organizations without contacts with filters on their locations
		if searchTerm != nil {
			query = query + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)--(l:Location)
		  WHERE NOT (o)--(:Contact) AND (l.country CONTAINS $location OR l.region CONTAINS $location OR l.locality CONTAINS $location)
		  RETURN o, null as c`)
		}

		//fetch contacts without organizations + filters on their properties
		query = query + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(c:Contact)
		  WHERE NOT (c)--(:Organization) AND (%s)
		  RETURN null as o, c`, contactFilterCypher)
		//fetch contacts without organizations with filters on their emails
		if searchTerm != nil {
			query = query + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(c:Contact)-[:HAS]->(e:Email)
		  WHERE NOT (c)--(:Organization) AND (e.email CONTAINS $email OR e.rawEmail CONTAINS $email)
		  RETURN null as o, c`)
		}
		//fetch contacts without organizations with filters on their locations
		if searchTerm != nil {
			query = query + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(c:Contact)--(l:Location)
		  WHERE NOT (c)--(:Organization) AND (l.country CONTAINS $location OR l.region CONTAINS $location OR l.locality CONTAINS $location)
		  RETURN null as o, c`)
		}
		//endregion

		query = query + fmt.Sprintf(`} RETURN o, c SKIP $skip LIMIT $limit`)

		queryResult, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
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
