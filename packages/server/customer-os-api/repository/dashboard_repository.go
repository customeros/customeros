package repository

import (
	"context"
	"encoding/json"
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
	"time"
)

type DashboardRepository interface {
	GetDashboardViewOrganizationData(ctx context.Context, tenant string, skip, limit int, where *model.Filter, sort *model.SortBy) (*utils.DbNodesWithTotalCount, error)
	GetDashboardNewCustomersData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardCustomerMapData(ctx context.Context, tenant string) ([]map[string]interface{}, error)
	GetDashboardRevenueAtRiskData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardMRRPerCustomerData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardARRBreakdownData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
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
		whereJSON, err := json.Marshal(where)
		if err == nil {
			span.LogFields(log.String("where", string(whereJSON)))
		}
	}
	if sort != nil {
		sortJSON, err := json.Marshal(sort)
		if err == nil {
			span.LogFields(log.Object("sort", sortJSON))
		}
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
					renewalLikelihoodValues = append(renewalLikelihoodValues, mapper.MapOpportunityRenewalLikelihoodFromString(&v))
				}
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("derivedRenewalLikelihood", renewalLikelihoodValues, utils.IN, false))
			} else if filter.Filter.Property == "RENEWAL_CYCLE_NEXT" && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("billingDetailsRenewalCycleNext", *filter.Filter.Value.Time, utils.LTE, false))
			} else if filter.Filter.Property == "RENEWAL_DATE" && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("derivedNextRenewalAt", *filter.Filter.Value.Time, utils.LTE, false))
			} else if filter.Filter.Property == "FORECAST_ARR" && filter.Filter.Value.ArrayInt != nil && len(*filter.Filter.Value.ArrayInt) == 2 {
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("renewalForecastArr", (*filter.Filter.Value.ArrayInt)[0], utils.GTE, false))
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("renewalForecastArr", (*filter.Filter.Value.ArrayInt)[1], utils.LTE, false))
			} else if filter.Filter.Property == "LAST_TOUCHPOINT_AT" && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, createCypherFilter("lastTouchpointAt", *filter.Filter.Value.Time, utils.GTE, false))
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

		span.LogFields(log.String("countQuery", countQuery))

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
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN owner.firstName <> \"\" and not owner.firstName is null THEN owner.firstName ELSE 'ZZZZZZZZZZZZZZZZZZZ' END as OWNER_FIRST_NAME_FOR_SORTING "
				query += ", CASE WHEN owner.lastName <> \"\" and not owner.lastName is null THEN owner.lastName ELSE 'ZZZZZZZZZZZZZZZZZZZ' END as OWNER_LAST_NAME_FOR_SORTING "
			} else {
				query += ", CASE WHEN owner.firstName <> \"\" and not owner.firstName is null THEN owner.firstName ELSE 'AAAAAAAAAAAAAAAAAAA' END as OWNER_FIRST_NAME_FOR_SORTING "
				query += ", CASE WHEN owner.lastName <> \"\" and not owner.lastName is null THEN owner.lastName ELSE 'AAAAAAAAAAAAAAAAAAA' END as OWNER_LAST_NAME_FOR_SORTING "
			}
			aliases += ", OWNER_FIRST_NAME_FOR_SORTING, OWNER_LAST_NAME_FOR_SORTING "
		}
		if sort != nil && sort.By == "NAME" {
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN o.name <> \"\" and not o.name is null THEN o.name ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as NAME_FOR_SORTING "
			} else {
				query += ", o.name as NAME_FOR_SORTING "
			}
			aliases += ", NAME_FOR_SORTING "
		}
		if sort != nil && sort.By == "RENEWAL_LIKELIHOOD" {
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN o.derivedRenewalLikelihoodOrder IS NOT NULL THEN o.derivedRenewalLikelihoodOrder ELSE 9999 END as RENEWAL_LIKELIHOOD_FOR_SORTING "
			} else {
				query += ", CASE WHEN o.derivedRenewalLikelihoodOrder IS NOT NULL THEN o.derivedRenewalLikelihoodOrder ELSE -1 END as RENEWAL_LIKELIHOOD_FOR_SORTING "
			}
			aliases += ", RENEWAL_LIKELIHOOD_FOR_SORTING "
		}
		if sort != nil && sort.By == "RENEWAL_CYCLE_NEXT" {
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN o.billingDetailsRenewalCycleNext IS NOT NULL THEN date(o.billingDetailsRenewalCycleNext) ELSE date('2100-01-01') END as RENEWAL_CYCLE_NEXT_FOR_SORTING "
			} else {
				query += ", CASE WHEN o.billingDetailsRenewalCycleNext IS NOT NULL THEN date(o.billingDetailsRenewalCycleNext) ELSE date('1900-01-01') END as RENEWAL_CYCLE_NEXT_FOR_SORTING "
			}
			aliases += ", RENEWAL_CYCLE_NEXT_FOR_SORTING "
		}
		if sort != nil && sort.By == "RENEWAL_DATE" {
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN o.derivedNextRenewalAt IS NOT NULL THEN date(o.derivedNextRenewalAt) ELSE date('2100-01-01') END as RENEWAL_DATE_FOR_SORTING "
			} else {
				query += ", CASE WHEN o.derivedNextRenewalAt IS NOT NULL THEN date(o.derivedNextRenewalAt) ELSE date('1900-01-01') END as RENEWAL_DATE_FOR_SORTING "
			}
			aliases += ", RENEWAL_DATE_FOR_SORTING "
		}
		if sort != nil && sort.By == "FORECAST_ARR" {
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN o.renewalForecastArr <> \"\" and o.renewalForecastArr IS NOT NULL THEN o.renewalForecastArr ELSE 9999999999999999 END as FORECAST_ARR_FOR_SORTING "
			} else {
				query += ", CASE WHEN o.renewalForecastArr <> \"\" and o.renewalForecastArr IS NOT NULL THEN o.renewalForecastArr ELSE 0 END as FORECAST_ARR_FOR_SORTING "
			}
			aliases += ", FORECAST_ARR_FOR_SORTING "
		}
		if sort != nil && sort.By == "ORGANIZATION" {
			query += " OPTIONAL MATCH (o)-[:SUBSIDIARY_OF]->(parent:Organization) WITH "
			query += aliases + ", parent "
		}

		cypherSort := utils.CypherSort{}
		if sort != nil {
			if sort.By == "NAME" {
				query += " ORDER BY NAME_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == "ORGANIZATION" {
				cypherSort.NewSortRule("NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.OrganizationEntity{})).WithCoalesce().WithAlias("parent")
				cypherSort.NewSortRule("NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.OrganizationEntity{})).WithCoalesce()
				cypherSort.NewSortRule("NAME", sort.Direction.String(), true, reflect.TypeOf(entity.OrganizationEntity{})).WithAlias("parent").WithDescending()
				cypherSort.NewSortRule("NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(entity.OrganizationEntity{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			} else if sort.By == "FORECAST_ARR" {
				query += " ORDER BY FORECAST_ARR_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == "RENEWAL_LIKELIHOOD" {
				query += " ORDER BY RENEWAL_LIKELIHOOD_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == "RENEWAL_CYCLE_NEXT" {
				query += " ORDER BY RENEWAL_CYCLE_NEXT_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == "RENEWAL_DATE" {
				query += " ORDER BY RENEWAL_DATE_FOR_SORTING " + string(sort.Direction)
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
				query += " ORDER BY OWNER_FIRST_NAME_FOR_SORTING " + string(sort.Direction) + ", OWNER_LAST_NAME_FOR_SORTING " + string(sort.Direction)
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

		span.LogFields(log.Object("query", query))
		paramsJson, err := json.Marshal(params)
		if err == nil {
			span.LogFields(log.String("params", string(paramsJson)))
		}

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

func (r *dashboardRepository) GetDashboardNewCustomersData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardNewCustomersData")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("startDate", startDate), log.Object("endDate", endDate))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			`
					WITH $startDate AS startDate, $endDate AS endDate
					WITH startDate.year AS startYear, startDate.month AS startMonth, endDate.year AS endYear, endDate.month AS endMonth, endDate
					WITH range(startYear * 12 + startMonth - 1, endYear * 12 + endMonth - 1) AS monthsRange, endDate
					UNWIND monthsRange AS monthsSinceEpoch
					
					WITH datetime({year: monthsSinceEpoch / 12, 
								   month: monthsSinceEpoch %s, 
								   day: 1}) AS currentDate, endDate

					WITH currentDate,
						 CASE 
						   WHEN currentDate.month = 12 THEN date({year: currentDate.year + 1, month: 1, day: 1})
						   ELSE date({year: currentDate.year, month: currentDate.month + 1, day: 1})
						 END AS startOfNextMonth
						
					WITH currentDate,
						 startOfNextMonth,
						 CASE 
						   WHEN startOfNextMonth.month = 1 THEN date({year: startOfNextMonth.year, month: 1, day: 1}) - duration({days: 1})
						   ELSE startOfNextMonth - duration({days: 1})
						 END AS endOfMonth

					WITH DISTINCT currentDate.year AS year, currentDate.month AS month, currentDate, datetime({year: endOfMonth.year, month: endOfMonth.month, day: endOfMonth.day, hour: 23, minute: 59, second: 59, nanosecond:999999999}) as endOfMonth
					OPTIONAL MATCH (t:Tenant{name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(i:Contract_%s)
					WHERE 
					  o.hide = false AND
					  i.serviceStartedAt.year = year AND 
					  i.serviceStartedAt.month = month AND 
					  (i.endedAt IS NULL OR i.endedAt > endOfMonth)
					
					WITH o, year, month, MIN(i.serviceStartedAt) AS oldestContractDate
					OPTIONAL MATCH (o)-[:HAS_CONTRACT]->(oldest:Contract_%s)
					WHERE oldest.serviceStartedAt = oldestContractDate
					RETURN year, month, COUNT(oldest) AS totalContracts
				`, "% 12 + 1", tenant, tenant, tenant),
			map[string]any{
				"tenant":    tenant,
				"startDate": startDate,
				"endDate":   endDate,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if dbRecords != nil {
		for _, v := range dbRecords.([]*neo4j.Record) {
			year := v.Values[0].(int64)
			month := v.Values[1].(int64)
			count := v.Values[2].(int64)

			record := map[string]interface{}{
				"year":  year,
				"month": month,
				"count": count,
			}

			results = append(results, record)
		}
	}

	return results, nil
}

func (r *dashboardRepository) GetDashboardCustomerMapData(ctx context.Context, tenant string) ([]map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardCustomerMapData")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			`
					MATCH (t:Tenant{name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:ACTIVE_RENEWAL]->(op:Opportunity_%s)
					WHERE o.hide = false AND c.serviceStartedAt IS NOT NULL and op.maxAmount IS NOT NULL
					WITH o,  
						 COLLECT(DISTINCT CASE
						   WHEN c.status = 'ENDED' THEN 'CHURNED'
						   WHEN c.status = 'LIVE' AND op.internalType = 'RENEWAL' AND op.renewalLikelihood = 'HIGH' THEN 'OK'
						   ELSE 'AT_RISK' END) AS statuses,
						 COLLECT(DISTINCT { id: c.id, serviceStartedAt: c.serviceStartedAt, status: c.status, maxAmount: op.maxAmount }) AS contractDetails
					WITH *, CASE
								WHEN ALL(x IN statuses WHERE x = 'CHURNED') THEN 'CHURNED'
								WHEN ALL(x IN statuses WHERE x IN ['OK', 'CHURNED']) THEN 'OK'
								ELSE 'AT_RISK'
							END AS status

					WITH *, REDUCE(s = null, cd IN contractDetails | 
							 CASE WHEN s IS NULL OR cd.serviceStartedAt < s THEN cd.serviceStartedAt ELSE s END
						 ) AS oldestServiceStartedAt

					WITH *, REDUCE(s = null, cd IN contractDetails | 
							 CASE WHEN s IS NULL OR cd.serviceStartedAt > s THEN cd.serviceStartedAt ELSE s END
						 ) AS latestServiceStartedAt

					WITH *, REDUCE(s = 0, cd IN contractDetails | 
							 CASE WHEN cd.serviceStartedAt = latestServiceStartedAt THEN s + cd.maxAmount ELSE s END 
						 ) AS latestContractLiveArr
					
					WITH *, REDUCE(sum = 0, cd IN contractDetails | CASE WHEN cd.status <> 'ENDED' THEN sum + cd.maxAmount ELSE sum END ) AS arr
					
					RETURN o.id,
						 oldestServiceStartedAt,
						 status,
						 CASE WHEN status = 'CHURNED' THEN latestContractLiveArr ELSE arr END as arr
					ORDER BY oldestServiceStartedAt ASC
				`, tenant, tenant, tenant),
			map[string]any{
				"tenant": tenant,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if dbRecords != nil {
		for _, v := range dbRecords.([]*neo4j.Record) {
			organizationId := v.Values[0].(string)
			oldestServiceStartedAt := v.Values[1].(time.Time)
			state := v.Values[2].(string)
			arr := getCorrectValueType(v.Values[3])

			record := map[string]interface{}{
				"organizationId":         organizationId,
				"oldestServiceStartedAt": oldestServiceStartedAt,
				"state":                  state,
				"arr":                    arr,
			}

			results = append(results, record)
		}
	}

	return results, nil
}

func (r *dashboardRepository) GetDashboardRevenueAtRiskData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardNewCustomersData")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("startDate", startDate), log.Object("endDate", endDate))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			`
					MATCH (t:Tenant{name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:ACTIVE_RENEWAL]->(op:Opportunity_%s)
					WHERE 
						o.hide = false AND c.status = 'LIVE' AND op.internalType = 'RENEWAL'
					
					WITH COLLECT(DISTINCT { renewalLikelihood: op.renewalLikelihood, maxAmount: op.maxAmount }) AS contractDetails
					
					return 
						REDUCE(sumHigh = 0, cd IN contractDetails | CASE WHEN cd.renewalLikelihood = 'HIGH' THEN sumHigh + cd.maxAmount ELSE sumHigh END ) AS high,
						REDUCE(sumAtRisk = 0, cd IN contractDetails | CASE WHEN cd.renewalLikelihood <> 'HIGH' THEN sumAtRisk + cd.maxAmount ELSE sumAtRisk END ) AS atRisk
				`, tenant, tenant, tenant),
			map[string]any{
				"tenant":    tenant,
				"startDate": startDate,
				"endDate":   endDate,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if dbRecords != nil {
		for _, v := range dbRecords.([]*neo4j.Record) {
			high := getCorrectValueType(v.Values[0])
			atRisk := getCorrectValueType(v.Values[1])

			record := map[string]interface{}{
				"high":   high,
				"atRisk": atRisk,
			}

			results = append(results, record)
		}
	}

	return results, nil
}

func (r *dashboardRepository) GetDashboardMRRPerCustomerData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardMRRPerCustomerData")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("startDate", startDate), log.Object("endDate", endDate))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			`
					WITH $startDate AS startDate, $endDate AS endDate
					WITH startDate.year AS startYear, startDate.month AS startMonth, endDate.year AS endYear, endDate.month AS endMonth
					WITH range(startYear * 12 + startMonth - 1, endYear * 12 + endMonth - 1) AS monthsRange
					UNWIND monthsRange AS monthsSinceEpoch
					
					WITH datetime({year: monthsSinceEpoch / 12, 
									month: monthsSinceEpoch %s, 
									day: 1}) AS currentDate
					
					WITH currentDate,
						 datetime({year: currentDate.year, 
									month: currentDate.month, 
									day: 1, hour: 0, minute: 0, second: 0, nanosecond: 0o00000000}) as beginOfMonth,
						 currentDate + duration({months: 1}) - duration({nanoseconds: 1}) as endOfMonth
					
					WITH DISTINCT currentDate.year AS year, currentDate.month AS month, beginOfMonth, endOfMonth
					
					OPTIONAL MATCH (t:Tenant{name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WHERE 
						o.hide = false AND o.isCustomer = true AND sli.startedAt is NOT null AND (sli.billed = 'MONTHLY' or sli.billed = 'QUARTERLY' or sli.billed = 'ANNUALLY') AND
						((sli.startedAt.year = beginOfMonth.year AND sli.startedAt.month = beginOfMonth.month) OR sli.startedAt < beginOfMonth) AND (sli.endedAt IS NULL OR sli.endedAt >= endOfMonth)
					
					WITH beginOfMonth, endOfMonth, year, month, COLLECT(DISTINCT { id: sli.id, startedAt: sli.startedAt, endedAt: sli.endedAt, amountPerMonth: CASE WHEN sli.billed = 'MONTHLY' THEN sli.price * sli.quantity ELSE CASE WHEN sli.billed = 'QUARTERLY' THEN  sli.price * sli.quantity / 4 ELSE CASE WHEN sli.billed = 'ANNUALLY' THEN sli.price * sli.quantity / 12 ELSE 0 END END END }) AS contractDetails
					
					WITH beginOfMonth, endOfMonth, contractDetails, year, month,  REDUCE(sumHigh = 0, cd IN contractDetails | sumHigh + cd.amountPerMonth ) AS mrr
					
					return year, month, mrr
				`, "% 12 + 1", tenant, tenant, tenant),
			map[string]any{
				"tenant":    tenant,
				"startDate": startDate,
				"endDate":   endDate,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if dbRecords != nil {
		for _, v := range dbRecords.([]*neo4j.Record) {
			year := v.Values[0].(int64)
			month := v.Values[1].(int64)
			amountPerMonth := getCorrectValueType(v.Values[2])

			record := map[string]interface{}{
				"year":           year,
				"month":          month,
				"amountPerMonth": amountPerMonth,
			}

			results = append(results, record)
		}
	}

	return results, nil
}

func (r *dashboardRepository) GetDashboardARRBreakdownData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardARRBreakdownData")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("startDate", startDate), log.Object("endDate", endDate))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			`
					WITH $startDate AS startDate, $endDate AS endDate
					WITH startDate.year AS startYear, startDate.month AS startMonth, endDate.year AS endYear, endDate.month AS endMonth
					WITH range(startYear * 12 + startMonth - 1, endYear * 12 + endMonth - 1) AS monthsRange
					UNWIND monthsRange AS monthsSinceEpoch
					
					WITH datetime({year: monthsSinceEpoch / 12, 
									month: monthsSinceEpoch %s, 
									day: 1}) AS currentDate
					
					WITH currentDate,
						 datetime({year: currentDate.year, 
									month: currentDate.month, 
									day: 1, hour: 0, minute: 0, second: 0, nanosecond: 0o00000000}) as beginOfMonth,
						 currentDate + duration({months: 1}) - duration({nanoseconds: 1}) as endOfMonth
					
					WITH DISTINCT currentDate.year AS year, currentDate.month AS month, beginOfMonth, endOfMonth
					
					OPTIONAL MATCH (t:Tenant{name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WHERE 
						o.hide = false AND o.isCustomer = true AND sli.startedAt is NOT null AND (sli.billed = 'MONTHLY' or sli.billed = 'QUARTERLY' or sli.billed = 'ANNUALLY')
					
					WITH year, month, COLLECT(DISTINCT { sliId: sli.id, sliCanceled: CASE WHEN sli.isCanceled IS NOT NULL THEN sli.isCanceled ELSE false END, sliStartedAt: sli.startedAt, sliEndedAt: sli.endedAt, sliAmountPerMonth: CASE WHEN sli.billed = 'MONTHLY' THEN sli.price * sli.quantity ELSE CASE WHEN sli.billed = 'QUARTERLY' THEN  sli.price * sli.quantity / 4 ELSE CASE WHEN sli.billed = 'ANNUALLY' THEN sli.price * sli.quantity / 12 ELSE 0 END END END }) AS contractDetails
					
					WITH year, month, contractDetails, 
						REDUCE(s = 0, cd IN contractDetails | CASE WHEN cd.sliCanceled = true AND cd.sliEndedAt.year = year AND cd.sliEndedAt.month = month THEN s + cd.sliAmountPerMonth ELSE s END ) AS cancellations
			
					return year, month, 0 as newlyContracted, 0 as renewals, 0 as upsells, 0 as downgrades, cancellations, 0 as churned, contractDetails
				`, "% 12 + 1", tenant, tenant, tenant),
			map[string]any{
				"tenant":    tenant,
				"startDate": startDate,
				"endDate":   endDate,
			})
		if err != nil {
			return nil, err
		}

		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}

	//,
	//REDUCE(newlyContracted = 0, cd IN contractDetails | CASE WHEN cd.status = 'ENDED' THEN newlyContracted + cd.sliAmountPerMonth ELSE churned END ) AS newlyContracted,
	//	REDUCE(churned = 0, cd IN contractDetails | CASE WHEN cd.status = 'ENDED' THEN churned + cd.sliAmountPerMonth ELSE churned END ) AS churned,
	//
	//	REDUCE(cancellations = 0, cd IN contractDetails | CASE WHEN cd.sliCanceled = true THEN cancellations + cd.sliAmountPerMonth ELSE cancellations END ) AS downgrades,

	var results []map[string]interface{}
	if dbRecords != nil {
		for _, v := range dbRecords.([]*neo4j.Record) {
			year := v.Values[0].(int64)
			month := v.Values[1].(int64)
			newlyContracted := getCorrectValueType(v.Values[2])
			renewals := getCorrectValueType(v.Values[3])
			upsells := getCorrectValueType(v.Values[4])
			downgrades := getCorrectValueType(v.Values[5])
			cancellations := getCorrectValueType(v.Values[6])
			churned := getCorrectValueType(v.Values[7])

			record := map[string]interface{}{
				"year":            year,
				"month":           month,
				"newlyContracted": newlyContracted,
				"renewals":        renewals,
				"upsells":         upsells,
				"downgrades":      downgrades,
				"cancellations":   cancellations,
				"churned":         churned,
			}

			results = append(results, record)
		}
	}

	return results, nil
}

func getCorrectValueType(valueToExtract any) float64 {
	var v float64

	switch val := valueToExtract.(type) {
	case int64:
		v = float64(val)
	case float64:
		v = val
	default:
		fmt.Errorf("unexpected type %T", val)
		v = 0
	}

	return v
}
