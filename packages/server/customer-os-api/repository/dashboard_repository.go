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
	"time"
)

type DashboardRepository interface {
	GetDashboardViewOrganizationData(ctx context.Context, tenant string, skip, limit int, where *model.Filter, sort *model.SortBy) (*utils.DbNodesWithTotalCount, error)
	GetDashboardNewCustomersData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardCustomerMapData(ctx context.Context, tenant string) ([]map[string]interface{}, error)
	GetDashboardRevenueAtRiskData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardMRRPerCustomerData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardARRBreakdownData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardARRBreakdownUpsellsAndDowngradesData(ctx context.Context, tenant, queryType string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardARRBreakdownRenewalsData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardRetentionRateContractsRenewalsData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardRetentionRateContractsChurnedData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
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
	tracing.LogObjectAsJson(span, "where", where)
	tracing.LogObjectAsJson(span, "sort", sort)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	organizationFilterCypher, organizationFilterParams := "", make(map[string]interface{})
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
			organizationFilterCypher, organizationFilterParams = organizationFilter.BuildCypherFilterFragmentWithParamName("o", "o_param_")
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
		countQuery := `MATCH (o:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
		if len(ownerId) > 0 {
			countQuery += ` OPTIONAL MATCH (o)<-[:OWNS]-(owner:User) WITH *`
		}
		if emailFilterCypher != "" {
			countQuery += ` MATCH (o)-[:HAS]->(e:Email) WITH *`
		}
		if locationFilterCypher != "" {
			countQuery += ` MATCH (o)-[:ASSOCIATED_WITH]->(l:Location) WITH *`
		}
		countQuery += ` WHERE o.hide = false `

		if organizationFilterCypher != "" || emailFilterCypher != "" || locationFilterCypher != "" || len(ownerId) > 0 {
			countQuery += " AND "
		}

		countQueryParts := []string{}
		if organizationFilterCypher != "" {
			countQueryParts = append(countQueryParts, organizationFilterCypher)
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
		query := `MATCH (o:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
		if len(ownerId) > 0 {
			query += fmt.Sprintf(` OPTIONAL MATCH (o)<-[:OWNS]-(owner:User) WITH *`)
		}
		query += fmt.Sprintf(` OPTIONAL MATCH (o)-[:HAS_DOMAIN]->(d:Domain) WITH *`)
		query += fmt.Sprintf(` OPTIONAL MATCH (o)-[:HAS]->(e:Email_%s) WITH *`, tenant)
		query += fmt.Sprintf(` OPTIONAL MATCH (o)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH *`, tenant)
		if sort != nil && sort.By == "OWNER" {
			query += fmt.Sprintf(` OPTIONAL MATCH (o)<-[:OWNS]-(owner:User_%s) WITH *`, tenant)
		}
		query += ` WHERE (o.hide = false) `

		if organizationFilterCypher != "" || emailFilterCypher != "" || locationFilterCypher != "" || len(ownerId) > 0 {
			query += " AND "
		}

		queryParts := []string{}
		if organizationFilterCypher != "" {
			queryParts = append(queryParts, organizationFilterCypher)
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
		tracing.LogObjectAsJson(span, "params", params)

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
					
					WITH beginOfMonth, endOfMonth, year, month, COLLECT(DISTINCT { id: sli.id, startedAt: sli.startedAt, endedAt: sli.endedAt, amountPerMonth: CASE WHEN sli.billed = 'MONTHLY' THEN sli.price * sli.quantity ELSE CASE WHEN sli.billed = 'QUARTERLY' THEN  sli.price * sli.quantity / 3 ELSE CASE WHEN sli.billed = 'ANNUALLY' THEN sli.price * sli.quantity / 12 ELSE 0 END END END }) AS contractDetails
					
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
					WITH datetime({
						year: monthsSinceEpoch / 12,
						month: monthsSinceEpoch %s,
						day: 1
					}) AS currentDate
					
					WITH currentDate,
						 datetime({
							 year: currentDate.year,
							 month: currentDate.month,
							 day: 1,
							 hour: 0,
							 minute: 0,
							 second: 0,
							 nanosecond: 0o00000000
						 }) as beginOfMonth,
						 currentDate + duration({months: 1}) - duration({nanoseconds: 1}) as endOfMonth,
						 currentDate + duration({months: 1}) as startOfNextMonth
					
					WITH DISTINCT currentDate.year AS year, currentDate.month AS month, beginOfMonth, endOfMonth, startOfNextMonth
					
					OPTIONAL MATCH (t:Tenant{name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WHERE o.hide = false AND o.isCustomer = true AND sli.startedAt IS NOT NULL AND (sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
					
					WITH year, month, endOfMonth, startOfNextMonth, c, sli ORDER BY sli.parentId ASC
					
					WITH year, month, endOfMonth, startOfNextMonth, c, 
						COLLECT({
							id: sli.id, 
							parentId: sli.parentId,
							sliCanceled: CASE WHEN sli.isCanceled IS NOT NULL THEN sli.isCanceled ELSE false END,
							sliStartedAt: sli.startedAt,
							sliEndedAt: sli.endedAt,
							sliAmountPerMonth: toFloat(CASE WHEN sli.billed = 'MONTHLY' THEN sli.price * sli.quantity * 12 ELSE CASE WHEN sli.billed = 'QUARTERLY' THEN sli.price * sli.quantity * 4 ELSE CASE WHEN sli.billed = 'ANNUALLY' THEN sli.price * sli.quantity ELSE 0 END END END),
							sliPreviousAmountPerMonth: toFloat(CASE WHEN sli.previousBilled = 'MONTHLY' THEN sli.previousPrice * sli.previousQuantity * 12 ELSE CASE WHEN sli.previousBilled = 'QUARTERLY' THEN sli.previousPrice * sli.previousQuantity * 4 ELSE CASE WHEN sli.previousBilled = 'ANNUALLY' THEN sli.previousPrice * sli.previousQuantity ELSE 0 END END END)
						}) AS sliDetails
					
					WITH year, month, endOfMonth, startOfNextMonth, c.id AS contractId, c.status AS contractStatus, c.serviceStartedAt AS contractStartedAt, c.endedAt AS contractEndedAt, COLLECT(sliDetails) AS sliList
					
					WITH year, month, endOfMonth, startOfNextMonth, COLLECT(DISTINCT {contractId: contractId, contractStatus: contractStatus, contractStartedAt: contractStartedAt, contractEndedAt: contractEndedAt, sliList: sliList}) AS contractDetails
					
					WITH year, month, endOfMonth, startOfNextMonth, contractDetails, 

						REDUCE(newlyContracted = 0, cd IN REDUCE (f = [], cd in contractDetails | CASE WHEN cd.contractStartedAt.year = year AND cd.contractStartedAt.month = month AND (cd.contractEndedAt IS NULL OR cd.contractEndedAt > endOfMonth) THEN f + [cd] ELSE f END) |
							newlyContracted + REDUCE(totalSli = 0, sliItem IN cd.sliList |
								REDUCE(cc = 0, arr IN sliItem |
									cc + REDUCE(innerTotal = 0, innerSliItem IN arr |
										CASE WHEN innerSliItem.sliCanceled = false AND innerSliItem.sliEndedAt IS NULL 
											 THEN innerTotal + innerSliItem.sliAmountPerMonth 
											 ELSE innerTotal 
										END
									)
								)
							)
						) AS newlyContracted,

						REDUCE(cancellations = 0, cd IN contractDetails |
							cancellations + REDUCE(totalSli = 0, sliItem IN cd.sliList |
								REDUCE(cc = 0, arr IN sliItem |
									cc + REDUCE(innerTotal = 0, innerSliItem IN arr |
										CASE WHEN innerSliItem.sliCanceled = true AND innerSliItem.sliEndedAt.year = year AND innerSliItem.sliEndedAt.month = month 
											 THEN innerTotal + innerSliItem.sliAmountPerMonth 
											 ELSE innerTotal 
										END
									)
								)
							)
						) AS cancellations,

						REDUCE(churned = 0, cd IN REDUCE (f = [], cd in contractDetails | CASE WHEN cd.contractEndedAt.year = year AND cd.contractEndedAt.month = month AND cd.contractStatus = 'ENDED' THEN f + [cd] ELSE f END) |
							churned + REDUCE(totalSli = 0, sliItem IN cd.sliList |
								REDUCE(cc = 0, arr IN sliItem |
									cc + REDUCE(innerTotal = 0, innerSliItem IN arr |
										CASE WHEN innerSliItem.sliCanceled = false AND innerSliItem.sliEndedAt IS NULL 
											 THEN innerTotal + innerSliItem.sliAmountPerMonth 
											 ELSE innerTotal 
										END
									)
								)
							)
						) AS churned
					
					RETURN year, month, newlyContracted, 0 as renewals, cancellations, churned, contractDetails

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
			newlyContracted := getCorrectValueType(v.Values[2])
			renewals := getCorrectValueType(v.Values[3])
			cancellations := getCorrectValueType(v.Values[4])
			churned := getCorrectValueType(v.Values[5])

			record := map[string]interface{}{
				"year":            year,
				"month":           month,
				"newlyContracted": newlyContracted,
				"renewals":        renewals,
				"cancellations":   cancellations,
				"churned":         churned,
			}

			results = append(results, record)
		}
	}

	return results, nil
}

func (r *dashboardRepository) GetDashboardARRBreakdownUpsellsAndDowngradesData(ctx context.Context, tenant, queryType string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardARRBreakdownUpsellsAndDowngradesData")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("startDate", startDate), log.Object("endDate", endDate))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	var q string
	if queryType == "UPSELLS" {
		q = "<"
	} else {
		q = ">"
	}

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			`
					WITH $startDate AS startDate, $endDate AS endDate
					
					WITH startDate.YEAR AS startYear, startDate.MONTH AS startMonth, endDate.YEAR AS endYear, endDate.MONTH AS endMonth
					WITH RANGE(startYear * 12 + startMonth - 1, endYear * 12 + endMonth - 1) AS monthsRange
					UNWIND monthsRange AS monthsSinceEpoch
					
					WITH datetime({
						YEAR: monthsSinceEpoch / 12,
						MONTH: monthsSinceEpoch %s,
						DAY: 1
					}) AS currentDate
					
					WITH currentDate,
						 datetime({
							 YEAR: currentDate.YEAR,
							 MONTH: currentDate.MONTH,
							 DAY: 1,
							 HOUR: 0,
							 MINUTE: 0,
							 SECOND: 0,
							 NANOSECOND: 0o00000000
						 }) AS beginOfMonth,
						 currentDate + duration({MONTHS: 1}) - duration({NANOSECONDS: 1}) AS endOfMonth,
						 currentDate + duration({MONTHS: 1}) AS startOfNextMonth
					
					WITH DISTINCT currentDate.YEAR AS year, currentDate.MONTH AS month, beginOfMonth, endOfMonth, startOfNextMonth
					
					OPTIONAL MATCH (t:Tenant {name: $tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WHERE (sli.endedAt IS NULL OR sli.endedAt > beginOfMonth) AND sli.startedAt < startOfNextMonth AND (sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
					
					WITH year, month, beginOfMonth, endOfMonth, startOfNextMonth, sli ORDER BY sli.startedAt ASC
					
					WITH year, month, beginOfMonth, endOfMonth, startOfNextMonth, sli.parentId AS pp, COLLECT(sli) AS versions
					WHERE SIZE(versions) > 1
					
					WITH year, month, beginOfMonth, endOfMonth, startOfNextMonth, pp, HEAD(versions) AS baseSliVersion, LAST(versions) AS lastSliVersion

					WITH year, month, beginOfMonth, endOfMonth, startOfNextMonth, pp, baseSliVersion, lastSliVersion,
						 TOFLOAT(CASE WHEN baseSliVersion.billed = 'MONTHLY' THEN baseSliVersion.price * baseSliVersion.quantity * 12 ELSE CASE WHEN baseSliVersion.billed = 'QUARTERLY' THEN baseSliVersion.price * baseSliVersion.quantity * 4 ELSE CASE WHEN baseSliVersion.billed = 'ANNUALLY' THEN baseSliVersion.price * baseSliVersion.quantity ELSE 0 END END END) AS headAmount,
						 TOFLOAT(CASE WHEN lastSliVersion.billed = 'MONTHLY' THEN lastSliVersion.price * lastSliVersion.quantity * 12 ELSE CASE WHEN lastSliVersion.billed = 'QUARTERLY' THEN lastSliVersion.price * lastSliVersion.quantity * 4 ELSE CASE WHEN lastSliVersion.billed = 'ANNUALLY' THEN lastSliVersion.price * lastSliVersion.quantity ELSE 0 END END END) AS lastAmount
					WHERE 
						baseSliVersion.startedAt < beginOfMonth 
						AND (lastSliVersion.endedAt IS NULL OR lastSliVersion.endedAt > endOfMonth)
						AND lastSliVersion.startedAt < startOfNextMonth 
						AND lastSliVersion.startedAt >= beginOfMonth
						AND headAmount %s lastAmount
					
					WITH year, month, SUM(lastAmount) AS totalHigh, SUM(headAmount) AS totalLow
					
					WITH year, month, totalHigh - totalLow AS total
					
					RETURN year, month, CASE WHEN total < 0 THEN -total ELSE total END AS total
				`, "% 12 + 1", tenant, tenant, q),
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
			value := getCorrectValueType(v.Values[2])
			record := map[string]interface{}{
				"year":  year,
				"month": month,
				"value": value,
			}

			results = append(results, record)
		}
	}

	return results, nil
}

func (r *dashboardRepository) GetDashboardARRBreakdownRenewalsData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardARRBreakdownRenewalsData")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("startDate", startDate), log.Object("endDate", endDate))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			`
					WITH $startDate AS startDate, $endDate AS endDate
					
					WITH startDate.YEAR AS startYear, startDate.MONTH AS startMonth, endDate.YEAR AS endYear, endDate.MONTH AS endMonth
					WITH RANGE(startYear * 12 + startMonth - 1, endYear * 12 + endMonth - 1) AS monthsRange
					UNWIND monthsRange AS monthsSinceEpoch
					
					WITH datetime({
						YEAR: monthsSinceEpoch / 12,
						MONTH: monthsSinceEpoch %s,
						DAY: 1
					}) AS currentDate
					
					WITH currentDate,
						 datetime({
							 YEAR: currentDate.YEAR,
							 MONTH: currentDate.MONTH,
							 DAY: 1,
							 HOUR: 0,
							 MINUTE: 0,
							 SECOND: 0,
							 NANOSECOND: 0o00000000
						 }) AS beginOfMonth,
						 currentDate + duration({MONTHS: 1}) - duration({NANOSECONDS: 1}) AS endOfMonth,
						 currentDate + duration({MONTHS: 1}) AS startOfNextMonth
					
					WITH DISTINCT currentDate.YEAR AS year, currentDate.MONTH AS month, beginOfMonth, endOfMonth, startOfNextMonth
					
					OPTIONAL MATCH (t:Tenant {name: $tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WITH year, month, beginOfMonth, endOfMonth, c.serviceStartedAt AS cssa, c.renewalCycle AS crc, c.renewalPeriods as crp, sli ORDER BY sli.startedAt ASC

					WHERE c.serviceStartedAt IS NOT NULL AND (c.endedAt IS NULL OR c.endedAt > beginOfMonth) AND (sli.endedAt IS NULL OR sli.endedAt > beginOfMonth) AND sli.startedAt < beginOfMonth AND (sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
					WITH year, month, beginOfMonth, endOfMonth, cssa, crc, crp, sli.parentId AS pp, COLLECT(sli) AS versions
					WITH year, month, beginOfMonth, endOfMonth, cssa, crc, crp, pp, LAST(versions) AS lastSliVersion
					WITH year, month, beginOfMonth, endOfMonth, cssa, crc, crp, pp, lastSliVersion
					WHERE
						CASE WHEN crc = 'ANNUALLY' THEN (CASE WHEN crp IS NULL THEN cssa.MONTH = beginOfMonth.MONTH ELSE cssa.MONTH = beginOfMonth.MONTH AND (beginOfMonth.YEAR - cssa.YEAR) %s = 0 END) ELSE 1 = 1 END AND
						CASE WHEN crc = 'QUARTERLY' THEN
							(lastSliVersion.billed IN ['MONTHLY', 'QUARTERLY'] AND beginOfMonth.MONTH IN [cssa.MONTH - 9, cssa.MONTH - 6, cssa.MONTH - 3, cssa.MONTH, cssa.MONTH + 3, cssa.MONTH + 6, cssa.MONTH + 9]) OR
							(lastSliVersion.billed = 'ANNUALLY' AND beginOfMonth.MONTH = cssa.MONTH)
						ELSE 1 = 1 END AND
						CASE WHEN crc = 'MONTHLY' THEN
							lastSliVersion.billed = 'MONTHLY' OR
							(lastSliVersion.billed = 'QUARTERLY' AND beginOfMonth.MONTH IN [cssa.MONTH - 9, cssa.MONTH - 6, cssa.MONTH - 3, cssa.MONTH, cssa.MONTH + 3, cssa.MONTH + 6, cssa.MONTH + 9]) OR
							(lastSliVersion.billed = 'ANNUALLY' AND beginOfMonth.MONTH = cssa.MONTH)
						ELSE 1 = 1 END
					
					WITH year, month, crc, crp, COLLECT(lastSliVersion) AS lasts
					WITH year, month, 
						REDUCE(s = 0.0, a IN lasts | s + 
							TOFLOAT(
							CASE WHEN crc = 'ANNUALLY' AND crp IS NULL THEN (CASE WHEN a.billed = 'MONTHLY' THEN 12 ELSE (CASE WHEN a.billed = 'QUARTERLY' THEN 4 ELSE (CASE WHEN a.billed = 'ANNUALLY' THEN 1 ELSE 0 END) END) END) ELSE 
							CASE WHEN crc = 'ANNUALLY' AND crp IS NOT NULL THEN (CASE WHEN a.billed = 'MONTHLY' THEN crp * 12 ELSE (CASE WHEN a.billed = 'QUARTERLY' THEN crp * 4 ELSE (CASE WHEN a.billed = 'ANNUALLY' THEN crp * 1 ELSE 0 END) END) END) ELSE 
							CASE WHEN crc = 'QUARTERLY' THEN (CASE WHEN a.billed = 'MONTHLY' THEN 3 ELSE (CASE WHEN a.billed = 'QUARTERLY' THEN 1 ELSE (CASE WHEN a.billed = 'ANNUALLY' THEN 1 ELSE 0 END) END) END) ELSE
							CASE WHEN crc = 'MONTHLY' THEN 1 ELSE 0 END END END END * a.price * a.quantity)) AS amount

					RETURN year, month, SUM(amount)
				`, "% 12 + 1", tenant, tenant, "% crp"),
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
			value := getCorrectValueType(v.Values[2])
			record := map[string]interface{}{
				"year":  year,
				"month": month,
				"value": value,
			}

			results = append(results, record)
		}
	}

	return results, nil
}

func (r *dashboardRepository) GetDashboardRetentionRateContractsRenewalsData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardRetentionRateContractsRenewalsData")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("startDate", startDate), log.Object("endDate", endDate))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			`
					WITH $startDate AS startDate, $endDate AS endDate
					
					WITH startDate.YEAR AS startYear, startDate.MONTH AS startMonth, endDate.YEAR AS endYear, endDate.MONTH AS endMonth
					WITH RANGE(startYear * 12 + startMonth - 1, endYear * 12 + endMonth - 1) AS monthsRange
					UNWIND monthsRange AS monthsSinceEpoch
					
					WITH datetime({
						YEAR: monthsSinceEpoch / 12,
						MONTH: monthsSinceEpoch %s,
						DAY: 1
					}) AS currentDate
					
					WITH currentDate,
						 datetime({
							 YEAR: currentDate.YEAR,
							 MONTH: currentDate.MONTH,
							 DAY: 1,
							 HOUR: 0,
							 MINUTE: 0,
							 SECOND: 0,
							 NANOSECOND: 0o00000000
						 }) AS beginOfMonth,
						 currentDate + duration({MONTHS: 1}) - duration({NANOSECONDS: 1}) AS endOfMonth,
						 currentDate + duration({MONTHS: 1}) AS startOfNextMonth
					
					WITH DISTINCT currentDate.YEAR AS year, currentDate.MONTH AS month, beginOfMonth, endOfMonth, startOfNextMonth
					
					OPTIONAL MATCH (t:Tenant {name: $tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WITH year, month, beginOfMonth, endOfMonth, c.serviceStartedAt AS cssa, c.renewalCycle AS crc, c.renewalPeriods as crp, sli ORDER BY sli.startedAt ASC

					WHERE c.serviceStartedAt IS NOT NULL AND (c.endedAt IS NULL OR c.endedAt > beginOfMonth) AND (sli.endedAt IS NULL OR sli.endedAt > beginOfMonth) AND sli.startedAt < beginOfMonth AND (sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
					WITH year, month, beginOfMonth, endOfMonth, cssa, crc, crp, sli.parentId AS pp, COLLECT(sli) AS versions
					WITH year, month, beginOfMonth, endOfMonth, cssa, crc, crp, pp, LAST(versions) AS lastSliVersion
					WITH year, month, beginOfMonth, endOfMonth, cssa, crc, crp, pp, lastSliVersion
					WHERE
						CASE WHEN crc = 'ANNUALLY' THEN (CASE WHEN crp IS NULL THEN cssa.MONTH = beginOfMonth.MONTH ELSE cssa.MONTH = beginOfMonth.MONTH AND (beginOfMonth.YEAR - cssa.YEAR) %s = 0 END) ELSE 1 = 1 END AND
						CASE WHEN crc = 'QUARTERLY' THEN
							(lastSliVersion.billed IN ['MONTHLY', 'QUARTERLY'] AND beginOfMonth.MONTH IN [cssa.MONTH - 9, cssa.MONTH - 6, cssa.MONTH - 3, cssa.MONTH, cssa.MONTH + 3, cssa.MONTH + 6, cssa.MONTH + 9]) OR
							(lastSliVersion.billed = 'ANNUALLY' AND beginOfMonth.MONTH = cssa.MONTH)
						ELSE 1 = 1 END AND
						CASE WHEN crc = 'MONTHLY' THEN
							lastSliVersion.billed = 'MONTHLY' OR
							(lastSliVersion.billed = 'QUARTERLY' AND beginOfMonth.MONTH IN [cssa.MONTH - 9, cssa.MONTH - 6, cssa.MONTH - 3, cssa.MONTH, cssa.MONTH + 3, cssa.MONTH + 6, cssa.MONTH + 9]) OR
							(lastSliVersion.billed = 'ANNUALLY' AND beginOfMonth.MONTH = cssa.MONTH)
						ELSE 1 = 1 END
					
					WITH year, month, crc, COUNT(lastSliVersion) AS numberOfSliRenewedPerContract
					RETURN year, month, COUNT(numberOfSliRenewedPerContract) AS numberOfContractsRenewed
				`, "% 12 + 1", tenant, tenant, "% crp"),
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
			value := getCorrectValueType(v.Values[2])
			record := map[string]interface{}{
				"year":  year,
				"month": month,
				"value": value,
			}

			results = append(results, record)
		}
	}

	return results, nil
}

func (r *dashboardRepository) GetDashboardRetentionRateContractsChurnedData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardRetentionRateContractsChurnedData")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("startDate", startDate), log.Object("endDate", endDate))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			`
					WITH $startDate AS startDate, $endDate AS endDate
					
					WITH startDate.YEAR AS startYear, startDate.MONTH AS startMonth, endDate.YEAR AS endYear, endDate.MONTH AS endMonth
					WITH RANGE(startYear * 12 + startMonth - 1, endYear * 12 + endMonth - 1) AS monthsRange
					UNWIND monthsRange AS monthsSinceEpoch
					
					WITH datetime({
						YEAR: monthsSinceEpoch / 12,
						MONTH: monthsSinceEpoch %s,
						DAY: 1
					}) AS currentDate
					
					WITH currentDate,
						 datetime({
							 YEAR: currentDate.YEAR,
							 MONTH: currentDate.MONTH,
							 DAY: 1,
							 HOUR: 0,
							 MINUTE: 0,
							 SECOND: 0,
							 NANOSECOND: 0o00000000
						 }) AS beginOfMonth,
						 currentDate + duration({MONTHS: 1}) AS startOfNextMonth
					
					WITH DISTINCT currentDate.YEAR AS year, currentDate.MONTH AS month, beginOfMonth, startOfNextMonth
					
					OPTIONAL MATCH (t:Tenant {name: $tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WITH year, month, beginOfMonth, startOfNextMonth, c.id as id, c.serviceStartedAt as serviceStartedAt, c.endedAt AS contractEndedAt sli ORDER BY sli.startedAt ASC

					WHERE c.serviceStartedAt IS NOT NULL AND contractEndedAt >= beginOfMonth AND contractEndedAt < startOfNextMonth AND sli.startedAt IS NOT NULL AND (sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
					return year, month, COUNT(DISTINCT id) AS value
				`, "% 12 + 1", tenant, tenant),
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
			value := getCorrectValueType(v.Values[2])
			record := map[string]interface{}{
				"year":  year,
				"month": month,
				"value": value,
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
