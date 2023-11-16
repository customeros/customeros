package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
	"strings"
)

type DashboardRepository interface {
	GetDashboardViewOrganizationData(ctx context.Context, tenant string, skip, limit int, where *model.Filter, sort *model.SortBy) (*utils.DbNodesWithTotalCount, error)
}

type dashboardRepository struct {
	driver *neo4j.DriverWithContext
}

func NewDashboardRepository(driver *neo4j.DriverWithContext) DashboardRepository {
	return &dashboardRepository{
		driver: driver,
	}
}

func createStringCypherFilter(propertyName string, searchTerm any, comparator utils.ComparisonOperator) *utils.CypherFilter {
	filter := utils.CypherFilter{}
	filter.Details = new(utils.CypherFilterItem)
	filter.Details.NodeProperty = propertyName
	filter.Details.Value = &searchTerm
	filter.Details.ComparisonOperator = comparator
	filter.Details.SupportCaseSensitive = true
	return &filter
}

func createCypherFilter(propertyName string, searchTerm any, comparator utils.ComparisonOperator, caseSensitive bool) *utils.CypherFilter {
	filter := utils.CypherFilter{}
	filter.Details = new(utils.CypherFilterItem)
	filter.Details.NodeProperty = propertyName
	filter.Details.Value = &searchTerm
	filter.Details.ComparisonOperator = comparator
	filter.Details.SupportCaseSensitive = caseSensitive
	return &filter
}

func createStringCypherFilterWithValueOrEmpty(filter *model.FilterItem, propertyName string) *utils.CypherFilter {
	if filter.IncludeEmpty != nil && *filter.IncludeEmpty {
		orFilter := utils.CypherFilter{}
		orFilter.LogicalOperator = utils.OR
		orFilter.Details = new(utils.CypherFilterItem)

		orFilter.Filters = append(orFilter.Filters, createStringCypherFilter(propertyName, *filter.Value.Str, utils.CONTAINS))
		orFilter.Filters = append(orFilter.Filters, createCypherFilter(propertyName, "", utils.EQUALS, false))
		orFilter.Filters = append(orFilter.Filters, createCypherFilter(propertyName, nil, utils.IS_NULL, false))
		return &orFilter
	} else {
		return createStringCypherFilter(propertyName, *filter.Value.Str, utils.CONTAINS)
	}
}

func (r *dashboardRepository) GetDashboardViewOrganizationData(ctx context.Context, tenant string, skip, limit int, where *model.Filter, sort *model.SortBy) (*utils.DbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardViewOrganizationData")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Int("skip", skip), log.Int("limit", limit))
	if where != nil {
		span.LogFields(log.Object("where", where))
	}
	if sort != nil {
		span.LogFields(log.Object("sort", sort))
	}

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	organizationfilterCypher, organizationFilterParams := "", make(map[string]interface{})
	emailFilterCypher, emailFilterParams := "", make(map[string]interface{})
	locationFilterCypher, locationFilterParams := "", make(map[string]interface{})

	ownerId := []string{}
	ownerIncludeEmpty := true

	//ORGANIZATION, EMAIL, COUNTRY, REGION, LOCALITY
	//region organization filters
	if where != nil {
		organizationFilter := new(utils.CypherFilter)
		organizationFilter.Negate = false
		organizationFilter.LogicalOperator = utils.AND
		organizationFilter.Filters = make([]*utils.CypherFilter, 0)

		emailFilter := new(utils.CypherFilter)
		emailFilter.Negate = false
		emailFilter.LogicalOperator = utils.AND
		emailFilter.Filters = make([]*utils.CypherFilter, 0)

		locationFilter := new(utils.CypherFilter)
		locationFilter.Negate = false
		locationFilter.LogicalOperator = utils.OR
		locationFilter.Filters = make([]*utils.CypherFilter, 0)

		for _, filter := range where.And {
			if filter.Filter.Property == "ORGANIZATION" {
				orFilter := utils.CypherFilter{}
				orFilter.LogicalOperator = utils.OR
				orFilter.Details = new(utils.CypherFilterItem)

				orFilter.Filters = append(orFilter.Filters, createStringCypherFilter("name", *filter.Filter.Value.Str, utils.CONTAINS))
				orFilter.Filters = append(orFilter.Filters, createStringCypherFilter("website", *filter.Filter.Value.Str, utils.CONTAINS))
				orFilter.Filters = append(orFilter.Filters, createStringCypherFilter("customerOsId", *filter.Filter.Value.Str, utils.CONTAINS))
				orFilter.Filters = append(orFilter.Filters, createStringCypherFilter("referenceId", *filter.Filter.Value.Str, utils.CONTAINS))

				organizationFilter.Filters = append(organizationFilter.Filters, &orFilter)
			} else if filter.Filter.Property == "NAME" {
				organizationFilter.Filters = append(organizationFilter.Filters, createStringCypherFilterWithValueOrEmpty(filter.Filter, "name"))
			} else if filter.Filter.Property == "WEBSITE" {
				organizationFilter.Filters = append(organizationFilter.Filters, createStringCypherFilterWithValueOrEmpty(filter.Filter, "website"))
			} else if filter.Filter.Property == "EMAIL" {
				emailFilter.Filters = append(emailFilter.Filters, createStringCypherFilter("email", *filter.Filter.Value.Str, utils.CONTAINS))
				emailFilter.Filters = append(emailFilter.Filters, createStringCypherFilter("rawEmail", *filter.Filter.Value.Str, utils.CONTAINS))
			} else if filter.Filter.Property == "COUNTRY" {
				locationFilter.Filters = append(locationFilter.Filters, createStringCypherFilter("country", *filter.Filter.Value.Str, utils.EQUALS))
			} else if filter.Filter.Property == "REGION" {
				locationFilter.Filters = append(locationFilter.Filters, createStringCypherFilter("region", *filter.Filter.Value.Str, utils.EQUALS))
			} else if filter.Filter.Property == "LOCALITY" {
				locationFilter.Filters = append(locationFilter.Filters, createStringCypherFilter("locality", *filter.Filter.Value.Str, utils.EQUALS))
			} else if filter.Filter.Property == "OWNER_ID" {
				ownerId = *filter.Filter.Value.ArrayStr
				ownerIncludeEmpty = *filter.Filter.IncludeEmpty
			} else if filter.Filter.Property == "IS_CUSTOMER" && filter.Filter.Value.ArrayBool != nil && len(*filter.Filter.Value.ArrayBool) >= 1 {
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("isCustomer", *filter.Filter.Value.ArrayBool, utils.IN, false))
			} else if filter.Filter.Property == "RENEWAL_LIKELIHOOD" && filter.Filter.Value.ArrayStr != nil && len(*filter.Filter.Value.ArrayStr) >= 1 {
				renewalLikelihoodValues := make([]string, 0)
				for _, v := range *filter.Filter.Value.ArrayStr {
					renewalLikelihoodValues = append(renewalLikelihoodValues, mapper.MapRenewalLikelihoodFromString(&v))
				}
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("renewalLikelihood", renewalLikelihoodValues, utils.IN, false))
			} else if filter.Filter.Property == "RENEWAL_CYCLE_NEXT" && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("billingDetailsRenewalCycleNext", *filter.Filter.Value.Time, utils.LTE, false))
			} else if filter.Filter.Property == "FORECAST_AMOUNT" && filter.Filter.Value.ArrayInt != nil && len(*filter.Filter.Value.ArrayInt) == 2 {
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("renewalForecastAmount", (*filter.Filter.Value.ArrayInt)[0], utils.GTE, false))
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("renewalForecastAmount", (*filter.Filter.Value.ArrayInt)[1], utils.LTE, false))
			} else if filter.Filter.Property == "LAST_TOUCHPOINT_AT" && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("lastTouchpointAt", *filter.Filter.Value.Time, utils.LTE, false))
			} else if filter.Filter.Property == "LAST_TOUCHPOINT_TYPE" && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("lastTouchpointType", *filter.Filter.Value.ArrayStr, utils.IN, false))
			}
		}

		if len(organizationFilter.Filters) > 0 {
			organizationfilterCypher, organizationFilterParams = organizationFilter.BuildCypherFilterFragmentWithParamName("o", "o_param_")
		}
		if len(emailFilter.Filters) > 0 {
			emailFilterCypher, emailFilterParams = emailFilter.BuildCypherFilterFragmentWithParamName("e", "e_param_")
		}
		if len(locationFilter.Filters) > 0 {
			locationFilterCypher, locationFilterParams = locationFilter.BuildCypherFilterFragmentWithParamName("l", "l_param_")
		}
	}

	//endregion
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"tenant":  tenant,
			"ownerId": ownerId,
			"skip":    skip,
			"limit":   limit,
		}

		utils.MergeMapToMap(organizationFilterParams, params)
		utils.MergeMapToMap(emailFilterParams, params)
		utils.MergeMapToMap(locationFilterParams, params)

		//region count query
		countQuery := fmt.Sprintf(`MATCH (o:Organization_%s)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) WITH * `, tenant)
		if len(ownerId) > 0 {
			countQuery += fmt.Sprintf(` OPTIONAL MATCH (o)<-[:OWNS]-(owner:User) WITH *`)
		}
		if emailFilterCypher != "" {
			countQuery += fmt.Sprintf(` MATCH (o)-[:HAS]->(e:Email_%s) WITH *`, tenant)
		}
		if locationFilterCypher != "" {
			countQuery += fmt.Sprintf(` MATCH (o)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH *`, tenant)
		}
		countQuery += fmt.Sprintf(` WHERE (o.hide = false OR o.hide is null)`)

		if organizationfilterCypher != "" || emailFilterCypher != "" || locationFilterCypher != "" || len(ownerId) > 0 {
			countQuery += " AND "
		}

		countQueryParts := []string{}
		if organizationfilterCypher != "" {
			countQueryParts = append(countQueryParts, organizationfilterCypher)
		}
		if len(ownerId) > 0 {
			if ownerIncludeEmpty {
				countQueryParts = append(countQueryParts, fmt.Sprintf(` (owner.id IN $ownerId OR owner.id IS NULL) `))
			} else {
				countQueryParts = append(countQueryParts, fmt.Sprintf(` owner.id IN $ownerId `))
			}
		}
		if emailFilterCypher != "" {
			countQueryParts = append(countQueryParts, emailFilterCypher)
		}
		if locationFilterCypher != "" {
			countQueryParts = append(countQueryParts, locationFilterCypher)
		}

		countQuery = countQuery + strings.Join(countQueryParts, " AND ") + fmt.Sprintf(` RETURN count(distinct(o))`)

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
		query := fmt.Sprintf(` MATCH (o:Organization_%s)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) WITH * `, tenant)
		if len(ownerId) > 0 {
			query += fmt.Sprintf(` OPTIONAL MATCH (o)<-[:OWNS]-(owner:User) WITH *`)
		}
		query += fmt.Sprintf(` OPTIONAL MATCH (o)-[:HAS_DOMAIN]->(d:Domain) WITH *`)
		query += fmt.Sprintf(` OPTIONAL MATCH (o)-[:HAS]->(e:Email_%s) WITH *`, tenant)
		query += fmt.Sprintf(` OPTIONAL MATCH (o)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH *`, tenant)
		if sort != nil && sort.By == "OWNER" {
			query += fmt.Sprintf(` OPTIONAL MATCH (o)<-[:OWNS]-(owner:User_%s) WITH *`, tenant)
		}
		query += ` WHERE (o.hide = false OR o.hide is null) `

		if organizationfilterCypher != "" || emailFilterCypher != "" || locationFilterCypher != "" || len(ownerId) > 0 {
			query += " AND "
		}

		queryParts := []string{}
		if organizationfilterCypher != "" {
			queryParts = append(queryParts, organizationfilterCypher)
		}
		if len(ownerId) > 0 {
			if ownerIncludeEmpty {
				queryParts = append(queryParts, fmt.Sprintf(` (owner.id IN $ownerId OR owner.id IS NULL) `))
			} else {
				queryParts = append(queryParts, fmt.Sprintf(` owner.id IN $ownerId `))
			}
		}
		if emailFilterCypher != "" {
			queryParts = append(queryParts, emailFilterCypher)
		}
		if locationFilterCypher != "" {
			queryParts = append(queryParts, locationFilterCypher)
		}

		//endregion
		query = query + strings.Join(queryParts, " AND ")

		// sort region
		aliases := " o, d, l"
		query += " WITH o, d, l "
		if sort != nil && sort.By == "OWNER" {
			query += ", owner"
			aliases += ", owner "
		}
		if sort != nil && sort.By == "ORGANIZATION" {
			query += " OPTIONAL MATCH (o)-[:SUBSIDIARY_OF]->(parent:Organization) WITH "
			query += aliases + ", parent "
		}
		cypherSort := utils.CypherSort{}
		if sort != nil {
			if sort.By == "NAME" || sort.By == "ORGANIZATION" {
				cypherSort.NewSortRule("NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.OrganizationEntity{})).WithCoalesce().WithAlias("parent")
				cypherSort.NewSortRule("NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.OrganizationEntity{})).WithCoalesce()
				cypherSort.NewSortRule("NAME", sort.Direction.String(), true, reflect.TypeOf(entity.OrganizationEntity{})).WithAlias("parent").WithDescending()
				cypherSort.NewSortRule("NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.OrganizationEntity{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			} else if sort.By == "FORECAST_AMOUNT" {
				cypherSort.NewSortRule("FORECAST_AMOUNT", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.RenewalForecast{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			} else if sort.By == "RENEWAL_LIKELIHOOD" {
				cypherSort.NewSortRule("RENEWAL_LIKELIHOOD", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.RenewalLikelihood{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			} else if sort.By == "RENEWAL_CYCLE_NEXT" {
				cypherSort.NewSortRule("RENEWAL_CYCLE_NEXT", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.BillingDetails{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			} else if sort.By == "LAST_TOUCHPOINT" {
				cypherSort.NewSortRule("LAST_TOUCHPOINT_AT", sort.Direction.String(), false, reflect.TypeOf(entity.OrganizationEntity{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			} else if sort.By == "DOMAIN" {
				cypherSort.NewSortRule("DOMAIN", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.DomainEntity{}))
				query += string(cypherSort.SortingCypherFragment("d"))
			} else if sort.By == "LOCATION" {
				cypherSort.NewSortRule("COUNTRY", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.LocationEntity{}))
				cypherSort.NewSortRule("REGION", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.LocationEntity{}))
				cypherSort.NewSortRule("LOCALITY", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.LocationEntity{}))
				query += string(cypherSort.SortingCypherFragment("l"))
			} else if sort.By == "OWNER" {
				cypherSort.NewSortRule("FIRST_NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.UserEntity{}))
				cypherSort.NewSortRule("LAST_NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.UserEntity{}))
				query += string(cypherSort.SortingCypherFragment("owner"))
			} else if sort.By == "LAST_TOUCHPOINT_AT" {
				cypherSort.NewSortRule("LAST_TOUCHPOINT_AT", sort.Direction.String(), false, reflect.TypeOf(entity.OrganizationEntity{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			} else if sort.By == "LAST_TOUCHPOINT_TYPE" {
				cypherSort.NewSortRule("LAST_TOUCHPOINT_TYPE", sort.Direction.String(), false, reflect.TypeOf(entity.OrganizationEntity{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			}
		} else {
			cypherSort.NewSortRule("UPDATED_AT", string(model.SortingDirectionDesc), false, reflect.TypeOf(entity.OrganizationEntity{}))
			query += string(cypherSort.SortingCypherFragment("o"))
		}
		// end sort region
		query += fmt.Sprintf(` RETURN distinct(o) `)
		query += fmt.Sprintf(` SKIP $skip LIMIT $limit`)

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
