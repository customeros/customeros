package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
)

type QueryRepository interface {
	GetDashboardViewContactsData(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip int, limit int, where *model.Filter) (*utils.DbNodesWithTotalCount, error)
	GetDashboardViewOrganizationData(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip int, limit int, where *model.Filter) (*utils.DbNodesWithTotalCount, error)
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

func (r *queryRepository) GetDashboardViewContactsData(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip int, limit int, where *model.Filter) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	contactFilterCypher, contactFilterParams := "()", make(map[string]interface{})
	emailFilterCypher, emailFilterParams := "()", make(map[string]interface{})
	locationFilterCypher, locationFilterParams := "()", make(map[string]interface{})

	//CONTACT, EMAIL, COUNTRY, REGION, LOCALITY
	//region organization filters
	if where != nil {

		contactFilter := new(utils.CypherFilter)
		contactFilter.Negate = false
		contactFilter.LogicalOperator = utils.OR
		contactFilter.Filters = make([]*utils.CypherFilter, 0)

		emailFilter := new(utils.CypherFilter)
		emailFilter.Negate = false
		emailFilter.LogicalOperator = utils.OR
		emailFilter.Filters = make([]*utils.CypherFilter, 0)

		locationFilter := new(utils.CypherFilter)
		locationFilter.Negate = false
		locationFilter.LogicalOperator = utils.OR
		locationFilter.Filters = make([]*utils.CypherFilter, 0)

		for _, filter := range where.And {
			if filter.Filter.Property == "CONTACT" {
				contactFilter.Filters = append(contactFilter.Filters, createCypherFilter("name", *filter.Filter.Value.Str))
				contactFilter.Filters = append(contactFilter.Filters, createCypherFilter("firstName", *filter.Filter.Value.Str))
				contactFilter.Filters = append(contactFilter.Filters, createCypherFilter("lastName", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "EMAIL" {
				emailFilter.Filters = append(emailFilter.Filters, createCypherFilter("email", *filter.Filter.Value.Str))
				emailFilter.Filters = append(emailFilter.Filters, createCypherFilter("rawEmail", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "COUNTRY" {
				locationFilter.Filters = append(locationFilter.Filters, createCypherFilter("country", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "REGION" {
				locationFilter.Filters = append(locationFilter.Filters, createCypherFilter("region", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "LOCALITY" {
				locationFilter.Filters = append(locationFilter.Filters, createCypherFilter("locality", *filter.Filter.Value.Str))
			}
		}

		contactFilterCypher, contactFilterParams = contactFilter.BuildCypherFilterFragmentWithParamName("c", "c_param_")
		emailFilterCypher, emailFilterParams = emailFilter.BuildCypherFilterFragmentWithParamName("e", "e_param_")
		locationFilterCypher, locationFilterParams = locationFilter.BuildCypherFilterFragmentWithParamName("l", "l_param_")
	}

	//endregion

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}
		utils.MergeMapToMap(contactFilterParams, params)
		utils.MergeMapToMap(emailFilterParams, params)
		utils.MergeMapToMap(locationFilterParams, params)

		//region count query
		countQuery := fmt.Sprintf(`
			MATCH (c:Contact_%s) WITH c
			OPTIONAL MATCH (c)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH c
		`, tenant, tenant)

		if contactFilterCypher != "()" || emailFilterCypher != "()" || locationFilterCypher != "()" {
			countQuery = countQuery + "WHERE "
		}

		countQueryParts := []string{}
		if contactFilterCypher != "()" {
			countQueryParts = append(countQueryParts, contactFilterCypher)
		}
		if emailFilterCypher != "()" {
			countQueryParts = append(countQueryParts, emailFilterCypher)
		}
		if locationFilterCypher != "()" {
			countQueryParts = append(countQueryParts, locationFilterCypher)
		}

		countQuery = countQuery + strings.Join(countQueryParts, " AND ") + fmt.Sprintf(` RETURN count(c)`)

		countQueryResult, err := tx.Run(ctx, countQuery, params)
		if err != nil {
			return nil, err
		}

		countRecord, err := countQueryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		dbNodesWithTotalCount.Count = countRecord.Values[0].(int64)
		//endregion

		//region query to fetch data
		query := fmt.Sprintf(`
			MATCH (c:Contact_%s) WITH c
			OPTIONAL MATCH (c)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH c
		`, tenant, tenant)

		if contactFilterCypher != "()" || emailFilterCypher != "()" || locationFilterCypher != "()" {
			query = query + "WHERE "
		}

		queryParts := []string{}
		if contactFilterCypher != "()" {
			queryParts = append(queryParts, contactFilterCypher)
		}
		if emailFilterCypher != "()" {
			queryParts = append(queryParts, emailFilterCypher)
		}
		if locationFilterCypher != "()" {
			queryParts = append(queryParts, locationFilterCypher)
		}

		//endregion
		query = query + strings.Join(queryParts, " AND ") + fmt.Sprintf(` RETURN c order by c.updatedAt desc SKIP $skip LIMIT $limit`)

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
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *queryRepository) GetDashboardViewOrganizationData(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip int, limit int, where *model.Filter) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	organizationfilterCypher, organizationFilterParams := "()", make(map[string]interface{})
	emailFilterCypher, emailFilterParams := "()", make(map[string]interface{})
	locationFilterCypher, locationFilterParams := "()", make(map[string]interface{})

	//ORGANIZATION, EMAIL, COUNTRY, REGION, LOCALITY
	//region organization filters
	if where != nil {

		organizationFilter := new(utils.CypherFilter)
		organizationFilter.Negate = false
		organizationFilter.LogicalOperator = utils.OR
		organizationFilter.Filters = make([]*utils.CypherFilter, 0)

		emailFilter := new(utils.CypherFilter)
		emailFilter.Negate = false
		emailFilter.LogicalOperator = utils.OR
		emailFilter.Filters = make([]*utils.CypherFilter, 0)

		locationFilter := new(utils.CypherFilter)
		locationFilter.Negate = false
		locationFilter.LogicalOperator = utils.OR
		locationFilter.Filters = make([]*utils.CypherFilter, 0)

		for _, filter := range where.And {
			if filter.Filter.Property == "ORGANIZATION" {
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("name", *filter.Filter.Value.Str))
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("website", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "EMAIL" {
				emailFilter.Filters = append(emailFilter.Filters, createCypherFilter("email", *filter.Filter.Value.Str))
				emailFilter.Filters = append(emailFilter.Filters, createCypherFilter("rawEmail", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "COUNTRY" {
				locationFilter.Filters = append(locationFilter.Filters, createCypherFilter("country", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "REGION" {
				locationFilter.Filters = append(locationFilter.Filters, createCypherFilter("region", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "LOCALITY" {
				locationFilter.Filters = append(locationFilter.Filters, createCypherFilter("locality", *filter.Filter.Value.Str))
			}
		}

		organizationfilterCypher, organizationFilterParams = organizationFilter.BuildCypherFilterFragmentWithParamName("o", "o_param_")
		emailFilterCypher, emailFilterParams = emailFilter.BuildCypherFilterFragmentWithParamName("e", "e_param_")
		locationFilterCypher, locationFilterParams = locationFilter.BuildCypherFilterFragmentWithParamName("l", "l_param_")
	}

	//endregion

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}
		utils.MergeMapToMap(organizationFilterParams, params)
		utils.MergeMapToMap(emailFilterParams, params)
		utils.MergeMapToMap(locationFilterParams, params)

		//region count query
		countQuery := fmt.Sprintf(`
			MATCH (o:Organization_%s) WITH o
			OPTIONAL MATCH (o)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH o
			WHERE (o.tenantOrganization = false OR o.tenantOrganization is null)
		`, tenant, tenant)

		if organizationfilterCypher != "()" || emailFilterCypher != "()" || locationFilterCypher != "()" {
			countQuery = countQuery + " AND "
		}

		countQueryParts := []string{}
		if organizationfilterCypher != "()" {
			countQueryParts = append(countQueryParts, organizationfilterCypher)
		}
		if emailFilterCypher != "()" {
			countQueryParts = append(countQueryParts, emailFilterCypher)
		}
		if locationFilterCypher != "()" {
			countQueryParts = append(countQueryParts, locationFilterCypher)
		}

		countQuery = countQuery + strings.Join(countQueryParts, " AND ") + fmt.Sprintf(` RETURN count(o)`)

		countQueryResult, err := tx.Run(ctx, countQuery, params)
		if err != nil {
			return nil, err
		}

		countRecord, err := countQueryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		dbNodesWithTotalCount.Count = countRecord.Values[0].(int64)
		//endregion

		//region query to fetch data
		query := fmt.Sprintf(`
			MATCH (o:Organization_%s) WITH o
			OPTIONAL MATCH (o)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH o
			WHERE (o.tenantOrganization = false OR o.tenantOrganization is null)
		`, tenant, tenant)

		if organizationfilterCypher != "()" || emailFilterCypher != "()" || locationFilterCypher != "()" {
			query = query + " AND "
		}

		queryParts := []string{}
		if organizationfilterCypher != "()" {
			queryParts = append(queryParts, organizationfilterCypher)
		}
		if emailFilterCypher != "()" {
			queryParts = append(queryParts, emailFilterCypher)
		}
		if locationFilterCypher != "()" {
			queryParts = append(queryParts, locationFilterCypher)
		}

		//endregion
		query = query + strings.Join(queryParts, " AND ") + fmt.Sprintf(` RETURN o order by o.updatedAt desc SKIP $skip LIMIT $limit`)

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
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}
