package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
)

type QueryRepository interface {
	GetOrganizationsAndContacts(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip int, limit int, searchTerm *string) (*utils.PairDbNodesWithTotalCount, error)
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
			params["email"] = strings.ToLower(*searchTerm)
			params["location"] = strings.ToLower(*searchTerm)
		}

		//region count query
		countQuery := fmt.Sprintf(`CALL {`)

		//fetch organizations and contacts + filters on their properties
		countQuery = countQuery + fmt.Sprintf(`
          MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
		  MATCH (o)<-[:ROLE_IN]-(:JobRole)<-[rel:WORKS_AS]-(c)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND ((%s) OR (%s))
		  RETURN rel`, organizationFilterCypher, contactFilterCypher)
		//fetch organizations and contacts with filters on their emails
		if searchTerm != nil {
			countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS]->(e:Email)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
		  MATCH (o)<-[:ROLE_IN]-(:JobRole)<-[rel:WORKS_AS]-(c)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND (toLower(e.email) CONTAINS $email OR toLower(e.rawEmail) CONTAINS $email)
		  RETURN rel

		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email)
		  MATCH (o)<-[:ROLE_IN]-(:JobRole)<-[rel:WORKS_AS]-(c)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND (toLower(e.email) CONTAINS $email or toLower(e.rawEmail) CONTAINS $email)
		  RETURN rel`)
		}
		//fetch organizations and contacts with filters on their locations
		if searchTerm != nil {
			countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:ASSOCIATED_WITH]->(l:Location)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
		  MATCH (o)<-[:ROLE_IN]-(:JobRole)<-[rel:WORKS_AS]-(c)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND (toLower(l.name) CONTAINS $location OR toLower(l.country) CONTAINS $location OR toLower(l.region) CONTAINS $location OR toLower(l.locality) CONTAINS $location)
		  RETURN rel

		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:ASSOCIATED_WITH]->(l:Location)
		  MATCH (o)<-[:ROLE_IN]-(:JobRole)<-[rel:WORKS_AS]-(c)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND (toLower(l.name) CONTAINS $location OR toLower(l.country) CONTAINS $location OR toLower(l.region) CONTAINS $location OR toLower(l.locality) CONTAINS $location)
		  RETURN rel`)
		}

		//fetch organizations without contacts + filters on their properties
		countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})-[rel:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND NOT (o)<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(:Contact) AND (%s)
		  RETURN rel`, organizationFilterCypher)
		//fetch organizations without contacts with filters on their emails
		if searchTerm != nil {
			countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})-[rel:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS]->(e:Email)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND NOT (o)<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(:Contact) AND (toLower(e.email) CONTAINS $email OR toLower(e.rawEmail) CONTAINS $email)
		  RETURN rel`)
		}
		//fetch organizations without contacts with filters on their locations
		if searchTerm != nil {
			countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})-[rel:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:ASSOCIATED_WITH]->(l:Location)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND NOT (o)<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(:Contact) AND (toLower(l.name) CONTAINS $location OR toLower(l.country) CONTAINS $location OR toLower(l.region) CONTAINS $location OR toLower(l.locality) CONTAINS $location)
		  RETURN rel`)
		}

		//fetch contacts without organizations + filters on their properties
		countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[rel:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
		  WHERE NOT (c)-[:WORKS_AS]->(:JobRole)-[:ROLE_IN]->(:Organization) AND (%s)
		  RETURN rel`, contactFilterCypher)
		//fetch contacts without organizations with filters on their emails
		if searchTerm != nil {
			countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[rel:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email)
		  WHERE NOT (c)-[:WORKS_AS]->(:JobRole)-[:ROLE_IN]->(:Organization) AND (toLower(e.email) CONTAINS $email OR toLower(e.rawEmail) CONTAINS $email)
		  RETURN rel`)
		}
		//fetch contacts without organizations with filters on their locations
		if searchTerm != nil {
			countQuery = countQuery + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[rel:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:ASSOCIATED_WITH]->(l:Location)
		  WHERE NOT (c)-[:WORKS_AS]->(:JobRole)-[:ROLE_IN]->(:Organization) AND (toLower(l.name) CONTAINS $location OR toLower(l.country) CONTAINS $location OR toLower(l.region) CONTAINS $location OR toLower(l.locality) CONTAINS $location)
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
          MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
		  MATCH (o)<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(c)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND ((%s) OR (%s))
		  RETURN DISTINCT o, c order by c.updatedAt desc`, organizationFilterCypher, contactFilterCypher)
		//fetch organizations and contacts with filters on their emails
		if searchTerm != nil {
			query = query + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS]->(e:Email)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
		  MATCH (o)<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(c)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND (toLower(e.email) CONTAINS $email OR toLower(e.rawEmail) CONTAINS $email)
		  RETURN DISTINCT o, c order by c.updatedAt desc

		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email)
		  MATCH (o)<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(c)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND (toLower(e.email) CONTAINS $email OR toLower(e.rawEmail) CONTAINS $email)
		  RETURN DISTINCT o, c order by c.updatedAt desc`)
		}
		//fetch organizations and contacts with filters on their locations
		if searchTerm != nil {
			query = query + fmt.Sprintf(`
		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:ASSOCIATED_WITH]->(l:Location)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
		  MATCH (o)<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(c)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND (toLower(l.name) CONTAINS $location OR toLower(l.country) CONTAINS $location OR toLower(l.region) CONTAINS $location OR toLower(l.locality) CONTAINS $location)
		  RETURN o, c order by c.updatedAt desc

		  UNION
		  MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
		  MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:ASSOCIATED_WITH]->(l:Location)
		  MATCH (o)<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(c)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND (toLower(l.name) CONTAINS $location OR toLower(l.country) CONTAINS $location OR toLower(l.region) CONTAINS $location OR toLower(l.locality) CONTAINS $location)
		  RETURN o, c order by c.updatedAt desc`)
		}

		//fetch contacts without organizations + filters on their properties
		query = query + fmt.Sprintf(`
		  UNION
		  MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
		  WHERE NOT (c)-[:WORKS_AS]->(:JobRole)-[:ROLE_IN]->(:Organization) AND (%s)
		  RETURN null as o, c order by c.updatedAt desc`, contactFilterCypher)
		//fetch contacts without organizations with filters on their emails
		if searchTerm != nil {
			query = query + fmt.Sprintf(`
		  UNION
		  MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email)
		  WHERE NOT (c)-[:WORKS_AS]->(:JobRole)-[:ROLE_IN]->(:Organization) AND (toLower(e.email) CONTAINS $email OR toLower(e.rawEmail) CONTAINS $email)
		  RETURN null as o, c order by c.updatedAt desc`)
		}
		//fetch contacts without organizations with filters on their locations
		if searchTerm != nil {
			query = query + fmt.Sprintf(`
		  UNION
		  MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:ASSOCIATED_WITH]->(l:Location)
		  WHERE NOT (c)-[:WORKS_AS]->(:JobRole)-[:ROLE_IN]->(:Organization) AND (toLower(l.name) CONTAINS $location OR toLower(l.country) CONTAINS $location OR toLower(l.region) CONTAINS $location OR toLower(l.locality) CONTAINS $location)
		  RETURN null as o, c order by c.updatedAt desc`)
		}

		//fetch organizations without contacts + filters on their properties
		query = query + fmt.Sprintf(`
		  UNION
		  MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND NOT (o)<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(:Contact) AND (%s)	
		  RETURN o, null as c order by o.updatedAt desc`, organizationFilterCypher)
		//fetch organizations without contacts with filters on their emails
		if searchTerm != nil {
			query = query + fmt.Sprintf(`
		  UNION
		  MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS]->(e:Email)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND NOT (o)<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(:Contact) AND (toLower(e.email) CONTAINS $email AND toLower(e.rawEmail) CONTAINS $email)
		  RETURN o, null as c order by o.updatedAt desc`)
		}
		//fetch organizations without contacts with filters on their locations
		if searchTerm != nil {
			query = query + fmt.Sprintf(`
		  UNION
		  MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:ASSOCIATED_WITH]->(l:Location)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND NOT (o)<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(:Contact) AND (toLower(l.name) CONTAINS $location OR toLower(l.country) CONTAINS $location OR toLower(l.region) CONTAINS $location OR toLower(l.locality) CONTAINS $location)
		  RETURN o, null as c order by o.updatedAt desc`)
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

func (r *queryRepository) GetDashboardViewContactsData(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip int, limit int, where *model.Filter) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	contactFilterCypher, contactFilterParams := "()", make(map[string]interface{})
	emailFilterCypher, emailFilterParams := "()", make(map[string]interface{})
	countryFilterCypher, countryFilterParams := "()", make(map[string]interface{})
	regionFilterCypher, regionFilterParams := "()", make(map[string]interface{})
	localityFilterCypher, localityFilterParams := "()", make(map[string]interface{})

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

		countryFilter := new(utils.CypherFilter)
		countryFilter.Negate = false
		countryFilter.LogicalOperator = utils.OR
		countryFilter.Filters = make([]*utils.CypherFilter, 0)

		regionFilter := new(utils.CypherFilter)
		regionFilter.Negate = false
		regionFilter.LogicalOperator = utils.OR
		regionFilter.Filters = make([]*utils.CypherFilter, 0)

		localityFilter := new(utils.CypherFilter)
		localityFilter.Negate = false
		localityFilter.LogicalOperator = utils.OR
		localityFilter.Filters = make([]*utils.CypherFilter, 0)

		for _, filter := range where.And {
			if filter.Filter.Property == "CONTACT" {
				contactFilter.Filters = append(contactFilter.Filters, createCypherFilter("name", *filter.Filter.Value.Str))
				contactFilter.Filters = append(contactFilter.Filters, createCypherFilter("firstName", *filter.Filter.Value.Str))
				contactFilter.Filters = append(contactFilter.Filters, createCypherFilter("lastName", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "EMAIL" {
				emailFilter.Filters = append(emailFilter.Filters, createCypherFilter("email", *filter.Filter.Value.Str))
				emailFilter.Filters = append(emailFilter.Filters, createCypherFilter("rawEmail", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "COUNTRY" {
				countryFilter.Filters = append(countryFilter.Filters, createCypherFilter("country", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "REGION" {
				regionFilter.Filters = append(regionFilter.Filters, createCypherFilter("region", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "LOCALITY" {
				localityFilter.Filters = append(localityFilter.Filters, createCypherFilter("locality", *filter.Filter.Value.Str))
			}
		}

		contactFilterCypher, contactFilterParams = contactFilter.BuildCypherFilterFragmentWithParamName("c", "c_param_")
		emailFilterCypher, emailFilterParams = emailFilter.BuildCypherFilterFragmentWithParamName("e", "e_param_")
		countryFilterCypher, countryFilterParams = countryFilter.BuildCypherFilterFragmentWithParamName("cr", "cr_param_")
		regionFilterCypher, regionFilterParams = regionFilter.BuildCypherFilterFragmentWithParamName("r", "r_param_")
		localityFilterCypher, localityFilterParams = localityFilter.BuildCypherFilterFragmentWithParamName("l", "l_param_")
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
		utils.MergeMapToMap(countryFilterParams, params)
		utils.MergeMapToMap(regionFilterParams, params)
		utils.MergeMapToMap(localityFilterParams, params)

		//region count query
		countQuery := fmt.Sprintf(`CALL {`)

		countQueryParts := []string{}

		//if no filter, we apply a straight query
		if contactFilterCypher == "()" && emailFilterCypher == "()" && countryFilterCypher == "()" && regionFilterCypher == "()" && localityFilterCypher == "()" {
			countQueryParts = append(countQueryParts, fmt.Sprintf(`
		 MATCH (c:Contact_%s)
		 RETURN c`, tenant))
		}

		//fetch organizations with filters on their properties
		if contactFilterCypher != "()" {
			countQueryParts = append(countQueryParts, fmt.Sprintf(`
		 MATCH (c:Contact_%s)
		 WHERE %s
		 RETURN c`, tenant, contactFilterCypher))
		}

		//fetch organizations with filters on their emails
		if emailFilterCypher != "" && emailFilterCypher != "()" {
			countQueryParts = append(countQueryParts, fmt.Sprintf(`
		  MATCH (c:Contact_%s)-[:HAS]->(e:Email_%s)
		  WHERE %s
		  RETURN c`, tenant, tenant, emailFilterCypher))
		}

		//fetch organizations with filters on their country locations
		if countryFilterCypher != "" && countryFilterCypher != "()" {
			countQueryParts = append(countQueryParts, fmt.Sprintf(`
		 MATCH (c:Contact_%s)-[:ASSOCIATED_WITH]->(cr:Location_%s)
		 WHERE %s
		 RETURN c`, tenant, tenant, countryFilterCypher))
		}

		//fetch organizations with filters on their region locations
		if regionFilterCypher != "" && regionFilterCypher != "()" {
			countQueryParts = append(countQueryParts, fmt.Sprintf(`
		 MATCH (c:Contact_%s)-[:ASSOCIATED_WITH]->(r:Location_%s)
		 WHERE %s
		 RETURN c`, tenant, tenant, regionFilterCypher))
		}

		//fetch organizations with filters on their region locations
		if localityFilterCypher != "" && localityFilterCypher != "()" {
			countQueryParts = append(countQueryParts, fmt.Sprintf(`
		 MATCH (c:Contact_%s)-[:ASSOCIATED_WITH]->(l:Location_%s)
		 WHERE %s
		 RETURN c`, tenant, tenant, localityFilterCypher))
		}

		countQuery = countQuery + strings.Join(countQueryParts, " UNION ") + fmt.Sprintf(`} RETURN count(c)`)

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
		query := fmt.Sprintf(`CALL {`)

		queryParts := []string{}

		//if no filter, we apply a straight query
		if contactFilterCypher == "()" && emailFilterCypher == "()" && countryFilterCypher == "()" && regionFilterCypher == "()" && localityFilterCypher == "()" {
			queryParts = append(queryParts, fmt.Sprintf(`
		 MATCH (c:Contact_%s)
		 RETURN c`, tenant))
		}

		//fetch organizations with filters on their properties
		if contactFilterCypher != "()" {
			queryParts = append(queryParts, fmt.Sprintf(`
		 MATCH (c:Contact_%s)
		 WHERE %s
		 RETURN c`, tenant, contactFilterCypher))
		}

		//fetch organizations with filters on their emails
		if emailFilterCypher != "" && emailFilterCypher != "()" {
			queryParts = append(queryParts, fmt.Sprintf(`
		  MATCH (c:Contact_%s)-[:HAS]->(e:Email_%s)
		  WHERE %s
		  RETURN c`, tenant, tenant, emailFilterCypher))
		}

		//fetch organizations with filters on their country locations
		if countryFilterCypher != "" && countryFilterCypher != "()" {
			queryParts = append(queryParts, fmt.Sprintf(`
		 MATCH (c:Contact_%s)-[:ASSOCIATED_WITH]->(cr:Location_%s)
		 WHERE %s
		 RETURN c`, tenant, tenant, countryFilterCypher))
		}

		//fetch organizations with filters on their region locations
		if regionFilterCypher != "" && regionFilterCypher != "()" {
			queryParts = append(queryParts, fmt.Sprintf(`
		 MATCH (c:Contact_%s)-[:ASSOCIATED_WITH]->(r:Location_%s)
		 WHERE %s
		 RETURN c`, tenant, tenant, regionFilterCypher))
		}

		//fetch organizations with filters on their region locations
		if localityFilterCypher != "" && localityFilterCypher != "()" {
			queryParts = append(queryParts, fmt.Sprintf(`
		 MATCH (c:Contact_%s)-[:ASSOCIATED_WITH]->(l:Location_%s)
		 WHERE %s
		 RETURN c`, tenant, tenant, localityFilterCypher))
		}

		//endregion
		query = query + strings.Join(queryParts, " UNION ") + fmt.Sprintf(`} RETURN c order by c.updatedAt desc SKIP $skip LIMIT $limit`)

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
	countryFilterCypher, countryFilterParams := "()", make(map[string]interface{})
	regionFilterCypher, regionFilterParams := "()", make(map[string]interface{})
	localityFilterCypher, localityFilterParams := "()", make(map[string]interface{})

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

		countryFilter := new(utils.CypherFilter)
		countryFilter.Negate = false
		countryFilter.LogicalOperator = utils.OR
		countryFilter.Filters = make([]*utils.CypherFilter, 0)

		regionFilter := new(utils.CypherFilter)
		regionFilter.Negate = false
		regionFilter.LogicalOperator = utils.OR
		regionFilter.Filters = make([]*utils.CypherFilter, 0)

		localityFilter := new(utils.CypherFilter)
		localityFilter.Negate = false
		localityFilter.LogicalOperator = utils.OR
		localityFilter.Filters = make([]*utils.CypherFilter, 0)

		for _, filter := range where.And {
			if filter.Filter.Property == "ORGANIZATION" {
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("name", *filter.Filter.Value.Str))
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("website", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "EMAIL" {
				emailFilter.Filters = append(emailFilter.Filters, createCypherFilter("email", *filter.Filter.Value.Str))
				emailFilter.Filters = append(emailFilter.Filters, createCypherFilter("rawEmail", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "COUNTRY" {
				countryFilter.Filters = append(countryFilter.Filters, createCypherFilter("country", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "REGION" {
				regionFilter.Filters = append(regionFilter.Filters, createCypherFilter("region", *filter.Filter.Value.Str))
			} else if filter.Filter.Property == "LOCALITY" {
				localityFilter.Filters = append(localityFilter.Filters, createCypherFilter("locality", *filter.Filter.Value.Str))
			}
		}

		organizationfilterCypher, organizationFilterParams = organizationFilter.BuildCypherFilterFragmentWithParamName("o", "o_param_")
		emailFilterCypher, emailFilterParams = emailFilter.BuildCypherFilterFragmentWithParamName("e", "e_param_")
		countryFilterCypher, countryFilterParams = countryFilter.BuildCypherFilterFragmentWithParamName("c", "c_param_")
		regionFilterCypher, regionFilterParams = regionFilter.BuildCypherFilterFragmentWithParamName("r", "r_param_")
		localityFilterCypher, localityFilterParams = localityFilter.BuildCypherFilterFragmentWithParamName("l", "l_param_")
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
		utils.MergeMapToMap(countryFilterParams, params)
		utils.MergeMapToMap(regionFilterParams, params)
		utils.MergeMapToMap(localityFilterParams, params)

		//if searchTerm != nil {
		//	params["email"] = strings.ToLower(*searchTerm)
		//	params["location"] = strings.ToLower(*searchTerm)
		//}

		//region count query
		countQuery := fmt.Sprintf(`CALL {`)

		countQueryParts := []string{}

		//if no filter, we apply a straight query
		if organizationfilterCypher == "()" && emailFilterCypher == "()" && countryFilterCypher == "()" && regionFilterCypher == "()" && localityFilterCypher == "()" {
			countQueryParts = append(countQueryParts, fmt.Sprintf(`
		 MATCH (o:Organization_%s)
		 WHERE (o.tenantOrganization = false OR o.tenantOrganization is null)
		 RETURN o`, tenant))
		}

		//fetch organizations with filters on their properties
		if organizationfilterCypher != "()" {
			countQueryParts = append(countQueryParts, fmt.Sprintf(`
		 MATCH (o:Organization_%s)
		 WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND %s
		 RETURN o`, tenant, organizationfilterCypher))
		}

		//fetch organizations with filters on their emails
		if emailFilterCypher != "" && emailFilterCypher != "()" {
			countQueryParts = append(countQueryParts, fmt.Sprintf(`
		  MATCH (o:Organization_%s)-[:HAS]->(e:Email_%s)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND %s
		  RETURN o`, tenant, tenant, emailFilterCypher))
		}

		//fetch organizations with filters on their country locations
		if countryFilterCypher != "" && countryFilterCypher != "()" {
			countQueryParts = append(countQueryParts, fmt.Sprintf(`
		 MATCH (o:Organization_%s)-[:ASSOCIATED_WITH]->(c:Location_%s)
		 WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND %s
		 RETURN o`, tenant, tenant, countryFilterCypher))
		}

		//fetch organizations with filters on their region locations
		if regionFilterCypher != "" && regionFilterCypher != "()" {
			countQueryParts = append(countQueryParts, fmt.Sprintf(`
		 MATCH (o:Organization_%s)-[:ASSOCIATED_WITH]->(r:Location_%s)
		 WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND %s
		 RETURN o`, tenant, tenant, regionFilterCypher))
		}

		//fetch organizations with filters on their region locations
		if localityFilterCypher != "" && localityFilterCypher != "()" {
			countQueryParts = append(countQueryParts, fmt.Sprintf(`
		 MATCH (o:Organization_%s)-[:ASSOCIATED_WITH]->(l:Location_%s)
		 WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND %s
		 RETURN o`, tenant, tenant, localityFilterCypher))
		}

		countQuery = countQuery + strings.Join(countQueryParts, " UNION ") + fmt.Sprintf(`} RETURN count(o)`)

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
		query := fmt.Sprintf(`CALL {`)

		queryParts := []string{}

		//if no filter, we apply a straight query
		if organizationfilterCypher == "()" && emailFilterCypher == "()" && countryFilterCypher == "()" && regionFilterCypher == "()" && localityFilterCypher == "()" {
			queryParts = append(queryParts, fmt.Sprintf(`
		 MATCH (o:Organization_%s)
		 WHERE (o.tenantOrganization = false OR o.tenantOrganization is null)
		 RETURN o`, tenant))
		}

		//fetch organizations with filters on their properties
		if organizationfilterCypher != "()" {
			queryParts = append(queryParts, fmt.Sprintf(`
		 MATCH (o:Organization_%s)
		 WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND %s
		 RETURN o`, tenant, organizationfilterCypher))
		}

		//fetch organizations with filters on their emails
		if emailFilterCypher != "" && emailFilterCypher != "()" {
			queryParts = append(queryParts, fmt.Sprintf(`
		  MATCH (o:Organization_%s)-[:HAS]->(e:Email_%s)
		  WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND %s
		  RETURN o`, tenant, tenant, emailFilterCypher))
		}

		//fetch organizations with filters on their country locations
		if countryFilterCypher != "" && countryFilterCypher != "()" {
			queryParts = append(queryParts, fmt.Sprintf(`
		 MATCH (o:Organization_%s)-[:ASSOCIATED_WITH]->(c:Location_%s)
		 WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND %s
		 RETURN o`, tenant, tenant, countryFilterCypher))
		}

		//fetch organizations with filters on their region locations
		if regionFilterCypher != "" && regionFilterCypher != "()" {
			queryParts = append(queryParts, fmt.Sprintf(`
		 MATCH (o:Organization_%s)-[:ASSOCIATED_WITH]->(r:Location_%s)
		 WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND %s
		 RETURN o`, tenant, tenant, regionFilterCypher))
		}

		//fetch organizations with filters on their region locations
		if localityFilterCypher != "" && localityFilterCypher != "()" {
			queryParts = append(queryParts, fmt.Sprintf(`
		 MATCH (o:Organization_%s)-[:ASSOCIATED_WITH]->(l:Location_%s)
		 WHERE (o.tenantOrganization = false OR o.tenantOrganization is null) AND %s
		 RETURN o`, tenant, tenant, localityFilterCypher))
		}

		//endregion
		query = query + strings.Join(queryParts, " UNION ") + fmt.Sprintf(`} RETURN o order by o.updatedAt desc SKIP $skip LIMIT $limit`)

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
