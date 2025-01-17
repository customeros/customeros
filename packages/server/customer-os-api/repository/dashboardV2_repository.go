package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
	"strings"
	"sync"
)

type DashboardV2Repository interface {
	GetDashboardViewOrganizationDataV2(ctx context.Context, tenant string, limit int, where *model.Filter, sort *commonmodel.SortBy) (*utils.StringsWithTotalCount, error)
	GetDashboardViewContactDataV2(ctx context.Context, tenant string, limit int, where *model.Filter, sort *commonmodel.SortBy) (*utils.StringsWithTotalCount, error)
}

type dashboardV2Repository struct {
	driver *neo4j.DriverWithContext
}

func NewDashboardV2Repository(driver *neo4j.DriverWithContext) DashboardV2Repository {
	return &dashboardV2Repository{
		driver: driver,
	}
}

func (r *dashboardV2Repository) GetDashboardViewOrganizationDataV2(ctx context.Context, tenant string, limit int, where *model.Filter, sort *commonmodel.SortBy) (*utils.StringsWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardViewOrganizationDataV2")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Int("limit", limit))
	tracing.LogObjectAsJson(span, "where", where)
	tracing.LogObjectAsJson(span, "sort", sort)

	organizationFilterCypher, organizationFilterParams := "", make(map[string]interface{})
	socialFilterCypher, socialFilterParams := "", make(map[string]interface{})
	tagFilterCypher, tagFilterParams := "", make(map[string]interface{})
	locationFilterCypher, locationFilterParams := "", make(map[string]interface{})
	industryFilterCypher, industryFilterParams := "", make(map[string]interface{})
	userFilterCypher, userFilterParams := "", make(map[string]interface{})
	domainFilterCypher, domainFilterParams := "", make(map[string]interface{})
	parentOrganizationFilterCypher, parentOrganizationFilterParams := "", make(map[string]interface{})

	//ORGANIZATION, EMAIL, COUNTRY, REGION, LOCALITY
	//region organization filters
	if where != nil {
		organizationFilter := new(utils.CypherFilter)
		organizationFilter.Negate = false
		organizationFilter.LogicalOperator = utils.AND
		organizationFilter.Filters = make([]*utils.CypherFilter, 0)

		socialFilter := new(utils.CypherFilter)
		socialFilter.Negate = false
		socialFilter.LogicalOperator = utils.AND
		socialFilter.Filters = make([]*utils.CypherFilter, 0)

		tagFilter := new(utils.CypherFilter)
		tagFilter.Negate = false
		tagFilter.LogicalOperator = utils.AND
		tagFilter.Filters = make([]*utils.CypherFilter, 0)

		emailFilter := new(utils.CypherFilter)
		emailFilter.Negate = false
		emailFilter.LogicalOperator = utils.AND
		emailFilter.Filters = make([]*utils.CypherFilter, 0)

		locationFilter := new(utils.CypherFilter)
		locationFilter.Negate = false
		locationFilter.LogicalOperator = utils.AND
		locationFilter.Filters = make([]*utils.CypherFilter, 0)

		industryFilter := new(utils.CypherFilter)
		industryFilter.Negate = false
		industryFilter.LogicalOperator = utils.AND
		industryFilter.Filters = make([]*utils.CypherFilter, 0)

		parentOrganizationFilter := new(utils.CypherFilter)
		parentOrganizationFilter.Negate = false
		parentOrganizationFilter.LogicalOperator = utils.AND
		parentOrganizationFilter.Filters = make([]*utils.CypherFilter, 0)

		userFilter := new(utils.CypherFilter)
		userFilter.Negate = false
		userFilter.LogicalOperator = utils.AND
		userFilter.Filters = make([]*utils.CypherFilter, 0)

		domainFilter := new(utils.CypherFilter)
		domainFilter.Negate = false
		domainFilter.LogicalOperator = utils.AND
		domainFilter.Filters = make([]*utils.CypherFilter, 0)

		for _, filter := range where.And {
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsName) {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("name", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsWebsite) {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("website", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsPrimaryDomains) {
				domainFilter.Filters = append(domainFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.DomainPropertyDomain), filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsRelationship) {
				createInOrEmptyStringFilter(filter, organizationFilter, "relationship")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsOnboardingStatus) {
				createInOrEmptyStringFilter(filter, organizationFilter, "onboardingStatus")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsRenewalLikelihood) {
				createInOrEmptyStringFilter(filter, organizationFilter, "derivedRenewalLikelihood")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsRenewalDate) {
				createTimeFilter(filter, organizationFilter, "derivedNextRenewalAt")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsForecastArr) {
				createNumberCypherFilter(filter, organizationFilter, "renewalForecastArr")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsOwner) {
				createInOrEmptyStringFilter(filter, userFilter, "id")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsLastTouchpoint) {
				createInOrEmptyStringFilter(filter, organizationFilter, string(neo4jentity.OrganizationPropertyLastTouchpointType))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsLastTouchpointDate) {
				createTimeFilter(filter, organizationFilter, string(neo4jentity.OrganizationPropertyLastTouchpointAt))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsStage) {
				createInOrEmptyStringFilter(filter, organizationFilter, string(neo4jentity.OrganizationPropertyStage))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsSocials) {
				socialFilter.Filters = append(socialFilter.Filters, utils.CreateStringCypherFilter("url", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsLeadSource) {
				createInOrEmptyStringFilter(filter, organizationFilter, string(neo4jentity.OrganizationPropertyLeadSource))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsCreatedDate) {
				createTimeFilter(filter, organizationFilter, "createdAt")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsEmployeeCount) {
				createNumberCypherFilter(filter, organizationFilter, "employees")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsContactCount) {
				createNumberCypherFilter(filter, organizationFilter, "derivedContactCount")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsYearFounded) {
				createNumberCypherFilter(filter, organizationFilter, "yearFounded")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsIndustry) {
				createInOrEmptyStringFilter(filter, industryFilter, string(neo4jentity.IndustryPropertyCode))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsChurnDate) {
				createTimeFilter(filter, organizationFilter, "derivedChurnedAt")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsLtv) {
				createNumberCypherFilter(filter, organizationFilter, "derivedLtv")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsCountry) {
				createInOrEmptyStringFilter(filter, locationFilter, "countryCodeA2")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsCity) {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("locality", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsIsPublic) {
				createBooleanFilter(filter, organizationFilter, "isPublic")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsTags) {
				// special case for not in tags
				if filter.Filter.Operation == commonmodel.ComparisonOperatorNotIn && filter.Filter.Value.ArrayStr != nil {
					rawCypher := ""
					for _, v := range *filter.Filter.Value.ArrayStr {
						if rawCypher != "" {
							rawCypher += " AND "
						}
						rawCypher += fmt.Sprintf(` NOT (o)-[:TAGGED]->(:Tag {name:"%s"}) `, v)
					}
					tagFilter.Filters = append(tagFilter.Filters, utils.CreateRawCypherFilter(rawCypher))
				} else {
					createInOrEmptyStringFilter(filter, tagFilter, string(neo4jentity.TagPropertyName))
				}
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsHeadquarters) {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("headquarters", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsParentOrganization) {
				parentOrganizationFilter.Filters = append(parentOrganizationFilter.Filters, utils.CreateStringCypherFilter("name", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsUpdatedDate) {
				createTimeFilter(filter, organizationFilter, "updatedAt")
			}
		}

		if len(where.And) == 0 {
			for _, filter := range where.Or {
				organizationFilter.LogicalOperator = utils.OR
				if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsName) {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("name", filter.Filter.Value.Str, filter.Filter.Operation))
				}
				if filter.Filter.Property == string(postgresEntity.ColumnViewTypeOrganizationsWebsite) {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("website", filter.Filter.Value.Str, filter.Filter.Operation))
				}
			}
		}

		if len(organizationFilter.Filters) > 0 {
			organizationFilterCypher, organizationFilterParams = organizationFilter.BuildCypherFilterFragmentWithParamName("o", "o_param_")
		}
		if len(socialFilter.Filters) > 0 {
			socialFilterCypher, socialFilterParams = socialFilter.BuildCypherFilterFragmentWithParamName("s", "s_param_")
		}
		if len(tagFilter.Filters) > 0 {
			tagFilterCypher, tagFilterParams = tagFilter.BuildCypherFilterFragmentWithParamName("t", "t_param_")
		}
		if len(locationFilter.Filters) > 0 {
			locationFilterCypher, locationFilterParams = locationFilter.BuildCypherFilterFragmentWithParamName("l", "l_param_")
		}
		if len(industryFilter.Filters) > 0 {
			industryFilterCypher, industryFilterParams = industryFilter.BuildCypherFilterFragmentWithParamName("i", "i_param_")
		}
		if len(userFilter.Filters) > 0 {
			userFilterCypher, userFilterParams = userFilter.BuildCypherFilterFragmentWithParamName("u", "u_param_")
		}
		if len(domainFilter.Filters) > 0 {
			domainFilterCypher, domainFilterParams = domainFilter.BuildCypherFilterFragmentWithParamName("d", "d_param_")
		}
		if len(parentOrganizationFilter.Filters) > 0 {
			parentOrganizationFilterCypher, parentOrganizationFilterParams = parentOrganizationFilter.BuildCypherFilterFragmentWithParamName("po", "po_param_")
		}
	}

	//endregion

	params := map[string]any{
		"tenant": tenant,
		"limit":  limit,
	}

	utils.MergeMapToMap(organizationFilterParams, params)
	utils.MergeMapToMap(socialFilterParams, params)
	utils.MergeMapToMap(tagFilterParams, params)
	utils.MergeMapToMap(locationFilterParams, params)
	utils.MergeMapToMap(industryFilterParams, params)
	utils.MergeMapToMap(userFilterParams, params)
	utils.MergeMapToMap(domainFilterParams, params)
	utils.MergeMapToMap(parentOrganizationFilterParams, params)

	//region count selectQuery
	countQuery := ""
	{
		countQuery += fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s) `, tenant)
		if domainFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (o)-[:HAS_DOMAIN]->(d:Domain{primary: true}) WITH *`
		}
		if userFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (o)<-[:OWNS]-(u:User) WITH *`
		}
		if socialFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (o)-[:HAS]->(s:Social) WITH *`
		}
		if tagFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (o)-[:TAGGED]->(t:Tag) WITH *`
		}
		if industryFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (o)-[:HAS_INDUSTRY]->(i:Industry) WITH *`
		}
		if locationFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (o)-[:ASSOCIATED_WITH]->(l:Location) WITH *`
		}
		if parentOrganizationFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (o)-[:SUBSIDIARY_OF]->(po:Organization) WITH *`
		}

		countQuery += ` WHERE (o.hide = false OR o.hide IS NULL) `

		if organizationFilterCypher != "" || domainFilterCypher != "" || socialFilterCypher != "" || tagFilterCypher != "" || locationFilterCypher != "" || parentOrganizationFilterCypher != "" || userFilterCypher != "" || industryFilterCypher != "" {
			countQuery += " AND "
		}

		countQueryParts := []string{}
		if organizationFilterCypher != "" {
			countQueryParts = append(countQueryParts, organizationFilterCypher)
		}
		if domainFilterCypher != "" {
			countQueryParts = append(countQueryParts, domainFilterCypher)
		}
		if userFilterCypher != "" {
			countQueryParts = append(countQueryParts, userFilterCypher)
		}
		if socialFilterCypher != "" {
			countQueryParts = append(countQueryParts, socialFilterCypher)
		}
		if tagFilterCypher != "" {
			countQueryParts = append(countQueryParts, tagFilterCypher)
		}
		if industryFilterCypher != "" {
			countQueryParts = append(countQueryParts, industryFilterCypher)
		}
		if locationFilterCypher != "" {
			countQueryParts = append(countQueryParts, locationFilterCypher)
		}
		if parentOrganizationFilterCypher != "" {
			countQueryParts = append(countQueryParts, parentOrganizationFilterCypher)
		}

		countQuery = countQuery + strings.Join(countQueryParts, " AND ") + fmt.Sprintf(` RETURN count(distinct(o))`)
	}
	//end count region

	selectQuery := ""
	//region selectQuery to fetch data
	{
		selectQuery += fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s) `, tenant)
		if userFilterCypher != "" || (sort != nil && (sort.By == string(postgresEntity.ColumnViewTypeOrganizationsOwner))) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)<-[:OWNS]-(u:User) WITH *`)
		}
		if domainFilterCypher != "" || (sort != nil && (sort.By == string(postgresEntity.ColumnViewTypeOrganizationsPrimaryDomains))) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[:HAS_DOMAIN]->(d:Domain{primary: true}) WITH *`)
		}
		if socialFilterCypher != "" {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[:HAS]->(s:Social_%s) WITH *`, tenant)
		}
		if tagFilterCypher != "" {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[:TAGGED]->(t:Tag_%s) WITH *`, tenant)
		}
		if industryFilterCypher != "" || (sort != nil && (sort.By == string(postgresEntity.ColumnViewTypeOrganizationsIndustry))) {
			selectQuery += ` OPTIONAL MATCH (o)-[:HAS_INDUSTRY]->(i:Industry) WITH *`
		}
		if locationFilterCypher != "" || (sort != nil && (sort.By == string(postgresEntity.ColumnViewTypeOrganizationsCountry) || sort.By == string(postgresEntity.ColumnViewTypeOrganizationsCity))) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH *`, tenant)
		}
		if parentOrganizationFilterCypher != "" || sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsParentOrganization) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[:SUBSIDIARY_OF]->(po:Organization_%s) WITH *`, tenant)
		}
		selectQuery += ` WHERE (o.hide = false OR o.hide IS NULL) `

		if organizationFilterCypher != "" || domainFilterCypher != "" || socialFilterCypher != "" || tagFilterCypher != "" || parentOrganizationFilterCypher != "" || locationFilterCypher != "" || userFilterCypher != "" || industryFilterCypher != "" {
			selectQuery += " AND "
		}

		queryParts := []string{}
		if organizationFilterCypher != "" {
			queryParts = append(queryParts, organizationFilterCypher)
		}
		if domainFilterCypher != "" {
			queryParts = append(queryParts, domainFilterCypher)
		}
		if userFilterCypher != "" {
			queryParts = append(queryParts, userFilterCypher)
		}
		if socialFilterCypher != "" {
			queryParts = append(queryParts, socialFilterCypher)
		}
		if tagFilterCypher != "" {
			queryParts = append(queryParts, tagFilterCypher)
		}
		if industryFilterCypher != "" {
			queryParts = append(queryParts, industryFilterCypher)
		}
		if locationFilterCypher != "" {
			queryParts = append(queryParts, locationFilterCypher)
		}
		if parentOrganizationFilterCypher != "" {
			queryParts = append(queryParts, parentOrganizationFilterCypher)
		}
		selectQuery = selectQuery + strings.Join(queryParts, " AND ")
	}
	//endregion

	// sort region
	sortingCypher := ""
	aliases := ""

	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsName) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.name <> \"\" and not o.name is null THEN toLower(o.name) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.name <> \"\" and not o.name is null THEN toLower(o.name) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsWebsite) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.website <> \"\" and not o.website is null THEN toLower(o.website) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.website <> \"\" and not o.website is null THEN toLower(o.website) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsPrimaryDomains) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN d.domain <> \"\" and not d.domain is null THEN toLower(d.domain) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN d.domain <> \"\" and not d.domain is null THEN toLower(d.domain) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsRelationship) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.relationship <> \"\" and not o.relationship is null THEN toLower(o.relationship) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.relationship <> \"\" and not o.relationship is null THEN toLower(o.relationship) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsOnboardingStatus) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.onboardingStatusOrder <> \"\" and not o.onboardingStatusOrder is null THEN o.onboardingStatusOrder ELSE 9999 END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.onboardingStatusOrder <> \"\" and not o.onboardingStatusOrder is null THEN o.onboardingStatusOrder ELSE -1 END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsRenewalLikelihood) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.derivedRenewalLikelihoodOrder <> \"\" and not o.derivedRenewalLikelihoodOrder is null THEN o.derivedRenewalLikelihoodOrder ELSE 9999 END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.derivedRenewalLikelihoodOrder <> \"\" and not o.derivedRenewalLikelihoodOrder is null THEN o.derivedRenewalLikelihoodOrder ELSE -1 END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsRenewalDate) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.derivedNextRenewalAt <> \"\" and not o.derivedNextRenewalAt is null THEN o.derivedNextRenewalAt ELSE datetime({year:2100}) END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.derivedNextRenewalAt <> \"\" and not o.derivedNextRenewalAt is null THEN o.derivedNextRenewalAt ELSE datetime({year:1900}) END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsForecastArr) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.renewalForecastArr <> \"\" and not o.renewalForecastArr is null THEN o.renewalForecastArr ELSE 9999 END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.renewalForecastArr <> \"\" and not o.renewalForecastArr is null THEN o.renewalForecastArr ELSE -1 END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsOwner) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN (COALESCE(u.name, '') + COALESCE(u.firstName, '') + COALESCE(u.lastName, '')) <> '' THEN toLower(COALESCE(u.name, '') + COALESCE(u.firstName, '') + COALESCE(u.lastName, '')) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN (COALESCE(u.name, '') + COALESCE(u.firstName, '') + COALESCE(u.lastName, '')) <> '' THEN toLower(COALESCE(u.name, '') + COALESCE(u.firstName, '') + COALESCE(u.lastName, '')) ELSE '' END as SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsLastTouchpoint) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.lastTouchpointAt <> \"\" and not o.lastTouchpointAt is null THEN o.lastTouchpointAt ELSE datetime({year:2100}) END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.lastTouchpointAt <> \"\" and not o.lastTouchpointAt is null THEN o.lastTouchpointAt ELSE datetime({year:1900}) END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsLastTouchpointDate) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.lastTouchpointAt <> \"\" and not o.lastTouchpointAt is null THEN o.lastTouchpointAt ELSE datetime({year:2100}) END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.lastTouchpointAt <> \"\" and not o.lastTouchpointAt is null THEN o.lastTouchpointAt ELSE datetime({year:1900}) END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsStage) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.stage <> \"\" and not o.stage is null THEN toLower(o.stage) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.stage <> \"\" and not o.stage is null THEN toLower(o.stage) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsLeadSource) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.leadSource <> \"\" and not o.leadSource is null THEN toLower(o.leadSource) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.leadSource <> \"\" and not o.leadSource is null THEN toLower(o.leadSource) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsCreatedDate) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN o.createdAt IS NOT NULL THEN o.createdAt ELSE datetime({year:2100}) END as SORT_BY `
		} else {
			aliases += `CASE WHEN o.createdAt IS NOT NULL THEN o.createdAt ELSE datetime({year:1900}) END as SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsEmployeeCount) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.employees <> \"\" and not o.employees is null THEN o.employees ELSE 999999999 END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.employees <> \"\" and not o.employees is null THEN o.employees ELSE -1 END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsContactCount) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN o.derivedContactCount IS NOT NULL THEN o.derivedContactCount ELSE 999999999 END as SORT_BY `
		} else {
			aliases += `CASE WHEN o.derivedContactCount IS NOT NULL THEN o.derivedContactCount ELSE -1 END as SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsYearFounded) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.yearFounded <> \"\" and not o.yearFounded is null THEN o.yearFounded ELSE 999999999 END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.yearFounded <> \"\" and not o.yearFounded is null THEN o.yearFounded ELSE -1 END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsIndustry) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN i.name <> "" and not i.name is null THEN toLower(i.name) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN i.name <> "" and not i.name is null THEN toLower(i.name) ELSE '' END as SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsChurnDate) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.derivedChurnedAt <> \"\" and not o.derivedChurnedAt is null THEN o.derivedChurnedAt ELSE datetime({year:2100}) END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.derivedChurnedAt <> \"\" and not o.derivedChurnedAt is null THEN o.derivedChurnedAt ELSE datetime({year:1900}) END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsLtv) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.derivedLtv <> \"\" and not o.derivedLtv is null THEN o.derivedLtv ELSE 9999999999999999 END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.derivedLtv <> \"\" and not o.derivedLtv is null THEN o.derivedLtv ELSE -9999999999999999 END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsCountry) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN l.country <> \"\" and not l.country is null THEN toLower(l.country) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN l.country <> \"\" and not l.country is null THEN toLower(l.country) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsCity) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN l.locality <> \"\" and not l.locality is null THEN toLower(l.locality) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN l.locality <> \"\" and not l.locality is null THEN toLower(l.locality) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsIsPublic) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.isPublic = true THEN 0 ELSE CASE WHEN o.isPublic = false THEN 1 ELSE 2 END END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.isPublic = false THEN 2 ELSE CASE WHEN o.isPublic = true THEN 1 ELSE 0 END END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsParentOrganization) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN po.name <> \"\" and not po.name is null THEN toLower(po.name) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN po.name <> \"\" and not po.name is null THEN toLower(po.name) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeOrganizationsUpdatedDate) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.updatedAt <> \"\" and not o.updatedAt is null THEN o.updatedAt ELSE datetime({year:2100}) END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.updatedAt <> \"\" and not o.updatedAt is null THEN o.updatedAt ELSE datetime({year:1900}) END as SORT_BY "
		}
	}

	if sort != nil {
		sortingCypher += " ORDER BY SORT_BY " + string(sort.Direction)
	}

	if len(aliases) > 0 {
		selectQuery += " WITH *, " + aliases
	} else {
		selectQuery += " WITH * "
	}

	cypherSort := utils.CypherSort{}
	if sort != nil {
		selectQuery += " " + sortingCypher
	} else {
		cypherSort.NewSortRule("UPDATED_AT", string(commonmodel.SortingDirectionDesc), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
		selectQuery += string(cypherSort.SortingCypherFragment("o"))
	}

	// end sort region
	selectQuery += fmt.Sprintf(` RETURN distinct(o.id) `)
	selectQuery += fmt.Sprintf(` LIMIT $limit`)

	stringsWithTotalCount := &utils.StringsWithTotalCount{}

	var wg sync.WaitGroup
	var firstErr error
	var mu sync.Mutex

	setError := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
	}

	wg.Add(1)
	go func(ctx context.Context, result *utils.StringsWithTotalCount) {
		span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardViewOrganizationDataV2.CountQuery")
		defer span.Finish()
		defer wg.Done()
		tracing.SetDefaultServiceSpanTags(ctx, span)

		tracing.LogObjectAsJson(span, "params", params)
		span.LogFields(log.String("countQuery", countQuery))

		session := utils.NewNeo4jReadSession(ctx, *r.driver)
		defer session.Close(ctx)

		countRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			countQueryResult, err := tx.Run(ctx, countQuery, params)
			if err != nil {
				return nil, err
			} else {
				return utils.ExtractSingleRecordFirstValueAsType[int64](ctx, countQueryResult, err)
			}
		})

		if err != nil {
			tracing.TraceErr(span, err)
			setError(err)
			return
		}

		result.Count = countRecord.(int64)
	}(ctx, stringsWithTotalCount)

	wg.Add(1)
	go func(ctx context.Context, result *utils.StringsWithTotalCount) {
		span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardViewOrganizationDataV2.SelectQuery")
		defer span.Finish()
		defer wg.Done()
		tracing.SetDefaultServiceSpanTags(ctx, span)

		tracing.LogObjectAsJson(span, "params", params)
		span.LogFields(log.String("selectQuery", selectQuery))

		session := utils.NewNeo4jReadSession(ctx, *r.driver)
		defer session.Close(ctx)

		dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			queryResult, err := tx.Run(ctx, selectQuery, params)
			if err != nil {
				return nil, err
			} else {
				return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
			}
		})

		if err != nil {
			tracing.TraceErr(span, err)
			setError(err)
			return
		}

		result.Strings = dbRecords.([]string)
	}(ctx, stringsWithTotalCount)

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return stringsWithTotalCount, nil
}

func createInOrEmptyStringFilter(filter *model.Filter, cypherFilter *utils.CypherFilter, neo4jProperty string) {
	if filter.Filter.Operation == commonmodel.ComparisonOperatorIsEmpty {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, nil, commonmodel.ComparisonOperatorIsEmpty))
	} else if filter.Filter.Operation == commonmodel.ComparisonOperatorIsNotEmpty {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, nil, commonmodel.ComparisonOperatorIsNotEmpty))
	} else if filter.Filter.Operation == commonmodel.ComparisonOperatorIn && filter.Filter.Value.ArrayStr != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, filter.Filter.Value.ArrayStr, commonmodel.ComparisonOperatorIn))
	} else if filter.Filter.Operation == commonmodel.ComparisonOperatorNotIn && filter.Filter.Value.ArrayStr != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, filter.Filter.Value.ArrayStr, commonmodel.ComparisonOperatorNotIn))
	}
}

func createTimeFilter(filter *model.Filter, cypherFilter *utils.CypherFilter, neo4jProperty string) {
	if filter.Filter.Operation == commonmodel.ComparisonOperatorBetween && filter.Filter.Value.ArrayTime != nil && len(*filter.Filter.Value.ArrayTime) == 2 {
		times := *filter.Filter.Value.ArrayTime
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, times[0], commonmodel.ComparisonOperatorGte))
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, times[1], commonmodel.ComparisonOperatorLte))
	} else if filter.Filter.Operation == commonmodel.ComparisonOperatorIsEmpty {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, nil, commonmodel.ComparisonOperatorIsEmpty))
	} else if filter.Filter.Operation == commonmodel.ComparisonOperatorGte && filter.Filter.Value.Time != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, *filter.Filter.Value.Time, commonmodel.ComparisonOperatorGte))
	} else if filter.Filter.Operation == commonmodel.ComparisonOperatorGt && filter.Filter.Value.Time != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, *filter.Filter.Value.Time, commonmodel.ComparisonOperatorGt))
	} else if filter.Filter.Operation == commonmodel.ComparisonOperatorLte && filter.Filter.Value.Time != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, *filter.Filter.Value.Time, commonmodel.ComparisonOperatorLte))
	} else if filter.Filter.Operation == commonmodel.ComparisonOperatorLt && filter.Filter.Value.Time != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, *filter.Filter.Value.Time, commonmodel.ComparisonOperatorLt))
	}
}

func createNumberCypherFilter(filter *model.Filter, cypherFilter *utils.CypherFilter, neo4jProperty string) {
	if filter.Filter.Value.Int == nil && filter.Filter.Value.Float == nil && filter.Filter.Operation != commonmodel.ComparisonOperatorIsEmpty && filter.Filter.Operation != commonmodel.ComparisonOperatorIsNotEmpty {
		return
	}

	if filter.Filter.Operation == commonmodel.ComparisonOperatorIsEmpty {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, nil, commonmodel.ComparisonOperatorIsNull))
	} else if filter.Filter.Operation == commonmodel.ComparisonOperatorIsEmpty {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, nil, commonmodel.ComparisonOperatorIsNotNull))
	} else {
		if filter.Filter.Value.Int != nil {
			createInternalNumberCypherFilter(filter.Filter.Value.Int, filter, cypherFilter, neo4jProperty)
		} else {
			createInternalNumberCypherFilter(filter.Filter.Value.Float, filter, cypherFilter, neo4jProperty)
		}
	}
}

func createInternalNumberCypherFilter(val any, filter *model.Filter, cypherFilter *utils.CypherFilter, neo4jProperty string) {
	if val == nil && filter.Filter.Operation != commonmodel.ComparisonOperatorIsEmpty && filter.Filter.Operation != commonmodel.ComparisonOperatorIsNotEmpty {
		return
	}

	if filter.Filter.Operation != commonmodel.ComparisonOperatorIsEmpty &&
		filter.Filter.Operation != commonmodel.ComparisonOperatorIsNotEmpty &&
		filter.Filter.Operation != commonmodel.ComparisonOperatorLt &&
		filter.Filter.Operation != commonmodel.ComparisonOperatorLte &&
		filter.Filter.Operation != commonmodel.ComparisonOperatorGt &&
		filter.Filter.Operation != commonmodel.ComparisonOperatorGte &&
		filter.Filter.Operation != commonmodel.ComparisonOperatorEquals &&
		filter.Filter.Operation != commonmodel.ComparisonOperatorNotEquals {
		return
	}

	// not equals should also show empty values
	if filter.Filter.Operation == commonmodel.ComparisonOperatorNotEquals {
		orFilter := utils.CypherFilter{}
		orFilter.LogicalOperator = utils.OR
		orFilter.Details = new(utils.CypherFilterItem)

		orFilter.Filters = append(orFilter.Filters, utils.CreateCypherFilter(neo4jProperty, "", commonmodel.ComparisonOperatorIsEmpty))
		orFilter.Filters = append(orFilter.Filters, utils.CreateCypherFilter(neo4jProperty, val, filter.Filter.Operation))

		cypherFilter.Filters = append(cypherFilter.Filters, &orFilter)
	} else {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, val, filter.Filter.Operation))
	}
}

func createBooleanFilter(filter *model.Filter, cypherFilter *utils.CypherFilter, neo4jProperty string) {
	if (filter.Filter.Operation == commonmodel.ComparisonOperatorEquals || filter.Filter.Operation == commonmodel.ComparisonOperatorNotEquals) && filter.Filter.Value.Bool != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, filter.Filter.Value.Bool, filter.Filter.Operation))
	}
}

func (r *dashboardV2Repository) GetDashboardViewContactDataV2(ctx context.Context, tenant string, limit int, where *model.Filter, sort *commonmodel.SortBy) (*utils.StringsWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardViewContactDataV2")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Int("limit", limit))
	tracing.LogObjectAsJson(span, "where", where)
	tracing.LogObjectAsJson(span, "sort", sort)

	contactFilterCypher, contactFilterParams := "", make(map[string]interface{})
	linkedInFilterCypher, linkedInFilterParams := "", make(map[string]interface{})
	tagFilterCypher, tagFilterParams := "", make(map[string]interface{})
	locationFilterCypher, locationFilterParams := "", make(map[string]interface{})
	primaryEmailFilterCypher, primaryEmailFilterParams := "", make(map[string]interface{})
	primaryOrganizationFilterCypher, primaryOrganizationFilterParams := "", make(map[string]interface{})
	primaryJobRoleFilterCypher, primaryJobRoleFilterParams := "", make(map[string]interface{})
	phoneNumberFilterCypher, phoneNumberFilterParams := "", make(map[string]interface{})
	flowFilterCypher, flowFilterParams := "", make(map[string]interface{})
	flowParticipantFilterCypher, flowParticipantFilterParams := "", make(map[string]interface{})
	userConnectedFilterCypher, userConnectedFilterParams := "", make(map[string]interface{})

	if where != nil {
		contactFilter := new(utils.CypherFilter)
		contactFilter.Negate = false
		contactFilter.LogicalOperator = utils.AND
		contactFilter.Filters = make([]*utils.CypherFilter, 0)

		linkedInFilter := new(utils.CypherFilter)
		linkedInFilter.Negate = false
		linkedInFilter.LogicalOperator = utils.AND
		linkedInFilter.Filters = make([]*utils.CypherFilter, 0)

		tagFilter := new(utils.CypherFilter)
		tagFilter.Negate = false
		tagFilter.LogicalOperator = utils.AND
		tagFilter.Filters = make([]*utils.CypherFilter, 0)

		primaryEmailFilter := new(utils.CypherFilter)
		primaryEmailFilter.Negate = false
		primaryEmailFilter.LogicalOperator = utils.AND
		primaryEmailFilter.Filters = make([]*utils.CypherFilter, 0)

		primaryOrganizationFilter := new(utils.CypherFilter)
		primaryOrganizationFilter.Negate = false
		primaryOrganizationFilter.LogicalOperator = utils.AND
		primaryOrganizationFilter.Filters = make([]*utils.CypherFilter, 0)

		primaryJobRoleFilter := new(utils.CypherFilter)
		primaryJobRoleFilter.Negate = false
		primaryJobRoleFilter.LogicalOperator = utils.AND
		primaryJobRoleFilter.Filters = make([]*utils.CypherFilter, 0)

		locationFilter := new(utils.CypherFilter)
		locationFilter.Negate = false
		locationFilter.LogicalOperator = utils.AND
		locationFilter.Filters = make([]*utils.CypherFilter, 0)

		phoneNumberFilter := new(utils.CypherFilter)
		phoneNumberFilter.Negate = false
		phoneNumberFilter.LogicalOperator = utils.AND
		phoneNumberFilter.Filters = make([]*utils.CypherFilter, 0)

		flowFilter := new(utils.CypherFilter)
		flowFilter.Negate = false
		flowFilter.LogicalOperator = utils.AND
		flowFilter.Filters = make([]*utils.CypherFilter, 0)

		flowParticipantFilter := new(utils.CypherFilter)
		flowParticipantFilter.Negate = false
		flowParticipantFilter.LogicalOperator = utils.AND
		flowParticipantFilter.Filters = make([]*utils.CypherFilter, 0)

		userConnectedFilter := new(utils.CypherFilter)
		userConnectedFilter.Negate = false
		userConnectedFilter.LogicalOperator = utils.AND
		userConnectedFilter.Filters = make([]*utils.CypherFilter, 0)

		for _, filter := range where.And {
			// TODO
			//ColumnViewTypeContactsEmails                     ColumnViewType = "CONTACTS_EMAILS"
			//ColumnViewTypeContactsPersonalEmails             ColumnViewType = "CONTACTS_PERSONAL_EMAILS"
			//ColumnViewTypeContactsPersona                    ColumnViewType = "CONTACTS_PERSONA"
			//ColumnViewTypeContactsLastInteraction            ColumnViewType = "CONTACTS_LAST_INTERACTION"
			//ColumnViewTypeContactsSkills                     ColumnViewType = "CONTACTS_SKILLS"
			//ColumnViewTypeContactsSchools                    ColumnViewType = "CONTACTS_SCHOOLS"
			//ColumnViewTypeContactsLanguages                  ColumnViewType = "CONTACTS_LANGUAGES"
			//ColumnViewTypeContactsTimeInCurrentRole          ColumnViewType = "CONTACTS_TIME_IN_CURRENT_ROLE"
			//ColumnViewTypeContactsExperience                 ColumnViewType = "CONTACTS_EXPERIENCE"
			//ColumnViewTypeContactsFlowNextAction             ColumnViewType = "CONTACTS_FLOW_NEXT_ACTION"
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsName) {
				logicalOperator := utils.AND
				if filter.Filter.Operation == commonmodel.ComparisonOperatorContains || filter.Filter.Operation == commonmodel.ComparisonOperatorIsNotEmpty {
					logicalOperator = utils.OR
				} else if filter.Filter.Operation == commonmodel.ComparisonOperatorNotContains || filter.Filter.Operation == commonmodel.ComparisonOperatorIsEmpty {
					logicalOperator = utils.AND
				} else {
					continue
				}
				innerGroupFilter := new(utils.CypherFilter)
				innerGroupFilter.Negate = false
				innerGroupFilter.LogicalOperator = logicalOperator
				innerGroupFilter.Filters = make([]*utils.CypherFilter, 0)
				innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.ContactPropertyFirstName), filter.Filter.Value.Str, filter.Filter.Operation))
				innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.ContactPropertyLastName), filter.Filter.Value.Str, filter.Filter.Operation))
				innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.ContactPropertyName), filter.Filter.Value.Str, filter.Filter.Operation))
				contactFilter.Filters = append(contactFilter.Filters, innerGroupFilter)
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsPrimaryEmail) {
				logicalOperator := utils.AND
				if filter.Filter.Operation == commonmodel.ComparisonOperatorContains || filter.Filter.Operation == commonmodel.ComparisonOperatorIsNotEmpty {
					logicalOperator = utils.OR
				} else if filter.Filter.Operation == commonmodel.ComparisonOperatorNotContains || filter.Filter.Operation == commonmodel.ComparisonOperatorIsEmpty {
					logicalOperator = utils.AND
				} else {
					continue
				}
				innerGroupFilter := new(utils.CypherFilter)
				innerGroupFilter.Negate = false
				innerGroupFilter.LogicalOperator = logicalOperator
				innerGroupFilter.Filters = make([]*utils.CypherFilter, 0)
				innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.EmailPropertyRawEmail), filter.Filter.Value.Str, filter.Filter.Operation))
				innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.EmailPropertyEmail), filter.Filter.Value.Str, filter.Filter.Operation))
				primaryEmailFilter.Filters = append(contactFilter.Filters, innerGroupFilter)
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeEmailVerificationPrimaryEmail) {

				if filter.Filter.Operation == commonmodel.ComparisonOperatorIsEmpty {
					innerGroupFilter := new(utils.CypherFilter)
					innerGroupFilter.Negate = false
					innerGroupFilter.LogicalOperator = utils.AND
					innerGroupFilter.Filters = make([]*utils.CypherFilter, 0)

					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNull(string(neo4jentity.EmailPropertyIsFirewalled)))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNull(string(neo4jentity.EmailPropertyIsFreeAccount)))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNull(string(neo4jentity.EmailPropertyIsRisky)))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNull(string(neo4jentity.EmailPropertyIsValidSyntax)))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNull(string(neo4jentity.EmailPropertyIsMailboxFull)))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNull(string(neo4jentity.EmailPropertyIsCatchAll)))

					primaryEmailFilter.Filters = append(contactFilter.Filters, innerGroupFilter)
				} else if filter.Filter.Operation == commonmodel.ComparisonOperatorIsNotEmpty {
					innerGroupFilter := new(utils.CypherFilter)
					innerGroupFilter.Negate = false
					innerGroupFilter.LogicalOperator = utils.OR
					innerGroupFilter.Filters = make([]*utils.CypherFilter, 0)

					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNotNull(string(neo4jentity.EmailPropertyIsFirewalled)))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNotNull(string(neo4jentity.EmailPropertyIsFreeAccount)))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNotNull(string(neo4jentity.EmailPropertyIsRisky)))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNotNull(string(neo4jentity.EmailPropertyIsValidSyntax)))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNotNull(string(neo4jentity.EmailPropertyIsMailboxFull)))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNotNull(string(neo4jentity.EmailPropertyIsCatchAll)))

					primaryEmailFilter.Filters = append(contactFilter.Filters, innerGroupFilter)
				} else {
					innerGroupFilter := new(utils.CypherFilter)
					innerGroupFilter.Filters = make([]*utils.CypherFilter, 0)

					if filter.Filter.Operation == commonmodel.ComparisonOperatorIn {
						innerGroupFilter.Negate = false
						innerGroupFilter.LogicalOperator = utils.OR
					} else if filter.Filter.Operation == commonmodel.ComparisonOperatorNotIn {
						innerGroupFilter.Negate = true
						innerGroupFilter.LogicalOperator = utils.AND
					}

					if filter.Filter.Value.ArrayStr != nil {
						for _, v := range *filter.Filter.Value.ArrayStr {

							if v == "firewall_protected" {
								innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterEq(string(neo4jentity.EmailPropertyIsFirewalled), true))
							}
							if v == "free_account" {
								innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterEq(string(neo4jentity.EmailPropertyIsFreeAccount), true))
							}
							if v == "no_risk" {
								innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterEq(string(neo4jentity.EmailPropertyIsRisky), false))
							}
							if v == "incorrect_format" {
								innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterEq(string(neo4jentity.EmailPropertyIsValidSyntax), false))
							}
							if v == "invalid_mailbox" { // todo WTF is this

								andFilter := new(utils.CypherFilter)
								andFilter.Negate = false
								andFilter.LogicalOperator = utils.AND
								andFilter.Filters = make([]*utils.CypherFilter, 0)

								andFilter.Filters = append(andFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.EmailPropertyDeliverable), "UNDELIVERABLE", commonmodel.ComparisonOperatorEquals))
								andFilter.Filters = append(andFilter.Filters, utils.CreateCypherFilterEq(string(neo4jentity.EmailPropertyIsMailboxFull), true))

								innerGroupFilter.Filters = append(innerGroupFilter.Filters, andFilter)
							}
							if v == "mailbox_full" {
								innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterEq(string(neo4jentity.EmailPropertyIsMailboxFull), true))
							}
							if v == "catch_all" {
								innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterEq(string(neo4jentity.EmailPropertyIsCatchAll), true))
							}
							if v == "not_verified" {
								innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNull(string(neo4jentity.EmailPropertyValidatedAt)))
							}
							if v == "verification_in_progress" {
								andFilter := new(utils.CypherFilter)
								andFilter.Negate = false
								andFilter.LogicalOperator = utils.AND
								andFilter.Filters = make([]*utils.CypherFilter, 0)

								andFilter.Filters = append(andFilter.Filters, utils.CreateCypherFilterIsNull(string(neo4jentity.EmailPropertyValidatedAt)))
								andFilter.Filters = append(andFilter.Filters, utils.CreateCypherFilterIsNotNull(string(neo4jentity.EmailPropertyValidationRequestedAt)))

								innerGroupFilter.Filters = append(innerGroupFilter.Filters, andFilter)
							}
						}
					}

					primaryEmailFilter.Filters = append(contactFilter.Filters, innerGroupFilter)
				}
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsCountry) {
				createInOrEmptyStringFilter(filter, locationFilter, "countryCodeA2")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsCity) {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.LocationPropertyLocality), filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsRegion) {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.LocationPropertyRegion), filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsTags) {
				// special case for not in tags
				if filter.Filter.Operation == commonmodel.ComparisonOperatorNotIn && filter.Filter.Value.ArrayStr != nil {
					rawCypher := ""
					for _, v := range *filter.Filter.Value.ArrayStr {
						if rawCypher != "" {
							rawCypher += " AND "
						}
						rawCypher += fmt.Sprintf(` NOT (c)-[:TAGGED]->(:Tag {name:"%s"}) `, v)
					}
					tagFilter.Filters = append(tagFilter.Filters, utils.CreateRawCypherFilter(rawCypher))
				} else {
					createInOrEmptyStringFilter(filter, tagFilter, string(neo4jentity.TagPropertyName))
				}
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsLinkedin) {
				logicalOperator := utils.AND
				if filter.Filter.Operation == commonmodel.ComparisonOperatorContains || filter.Filter.Operation == commonmodel.ComparisonOperatorIsNotEmpty {
					logicalOperator = utils.OR
				} else if filter.Filter.Operation == commonmodel.ComparisonOperatorNotContains || filter.Filter.Operation == commonmodel.ComparisonOperatorIsEmpty {
					logicalOperator = utils.AND
				} else {
					continue
				}
				parentInnerGroupFilter := new(utils.CypherFilter)
				parentInnerGroupFilter.Negate = false
				parentInnerGroupFilter.LogicalOperator = utils.OR
				parentInnerGroupFilter.Filters = make([]*utils.CypherFilter, 0)

				innerGroupFilter1 := new(utils.CypherFilter)
				innerGroupFilter1.Negate = false
				innerGroupFilter1.LogicalOperator = logicalOperator
				innerGroupFilter1.Filters = make([]*utils.CypherFilter, 0)
				innerGroupFilter1.Filters = append(innerGroupFilter1.Filters, utils.CreateStringCypherFilter(string(neo4jentity.SocialPropertyAlias), filter.Filter.Value.Str, filter.Filter.Operation))
				innerGroupFilter1.Filters = append(innerGroupFilter1.Filters, utils.CreateStringCypherFilter(string(neo4jentity.SocialPropertyUrl), filter.Filter.Value.Str, filter.Filter.Operation))

				parentInnerGroupFilter.Filters = append(parentInnerGroupFilter.Filters, innerGroupFilter1)

				if filter.Filter.Operation == commonmodel.ComparisonOperatorNotContains {
					innerGroupFilter2 := new(utils.CypherFilter)
					innerGroupFilter2.Negate = false
					innerGroupFilter2.LogicalOperator = utils.AND
					innerGroupFilter2.Filters = make([]*utils.CypherFilter, 0)
					innerGroupFilter2.Filters = append(innerGroupFilter2.Filters, utils.CreateCypherFilterIsNull(string(neo4jentity.SocialPropertyAlias)))
					parentInnerGroupFilter.Filters = append(parentInnerGroupFilter.Filters, innerGroupFilter2)
				}

				linkedInFilter.Filters = append(linkedInFilter.Filters, parentInnerGroupFilter)
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsLinkedinFollowerCount) {
				createNumberCypherFilter(filter, linkedInFilter, string(neo4jentity.SocialPropertyFollowersCount))
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsOrganization) {
				if filter.Filter.Operation == commonmodel.ComparisonOperatorNotContains {
					innerGroupFilter := new(utils.CypherFilter)
					innerGroupFilter.Negate = false
					innerGroupFilter.LogicalOperator = utils.OR
					innerGroupFilter.Filters = make([]*utils.CypherFilter, 0)
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.OrganizationPropertyName), filter.Filter.Value.Str, filter.Filter.Operation))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNull(string(neo4jentity.OrganizationPropertyName)))
					primaryOrganizationFilter.Filters = append(primaryOrganizationFilter.Filters, innerGroupFilter)
				} else {
					primaryOrganizationFilter.Filters = append(primaryOrganizationFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.OrganizationPropertyName), filter.Filter.Value.Str, filter.Filter.Operation))
				}
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsJobTitle) {
				if filter.Filter.Operation == commonmodel.ComparisonOperatorNotContains {
					innerGroupFilter := new(utils.CypherFilter)
					innerGroupFilter.Negate = false
					innerGroupFilter.LogicalOperator = utils.OR
					innerGroupFilter.Filters = make([]*utils.CypherFilter, 0)
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateStringCypherFilter("jobTitle", filter.Filter.Value.Str, filter.Filter.Operation))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNull("jobTitle"))
					primaryJobRoleFilter.Filters = append(primaryJobRoleFilter.Filters, innerGroupFilter)
				} else {
					primaryJobRoleFilter.Filters = append(primaryJobRoleFilter.Filters, utils.CreateStringCypherFilter("jobTitle", filter.Filter.Value.Str, filter.Filter.Operation))
				}
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsTimeInCurrentRole) {
				createTimeFilter(filter, primaryJobRoleFilter, "startedAt")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsPhoneNumbers) {
				if filter.Filter.Operation == commonmodel.ComparisonOperatorNotContains {
					innerGroupFilter := new(utils.CypherFilter)
					innerGroupFilter.Negate = false
					innerGroupFilter.LogicalOperator = utils.OR
					innerGroupFilter.Filters = make([]*utils.CypherFilter, 0)
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateStringCypherFilter("rawPhoneNumber", filter.Filter.Value.Str, filter.Filter.Operation))
					innerGroupFilter.Filters = append(innerGroupFilter.Filters, utils.CreateCypherFilterIsNull("rawPhoneNumber"))
					phoneNumberFilter.Filters = append(phoneNumberFilter.Filters, innerGroupFilter)
				} else {
					phoneNumberFilter.Filters = append(phoneNumberFilter.Filters, utils.CreateStringCypherFilter("rawPhoneNumber", filter.Filter.Value.Str, filter.Filter.Operation))
				}
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsFlows) {
				// special case for not in flows
				if filter.Filter.Operation == commonmodel.ComparisonOperatorNotIn && filter.Filter.Value.ArrayStr != nil {
					rawCypher := ""
					for _, v := range *filter.Filter.Value.ArrayStr {
						if rawCypher != "" {
							rawCypher += " AND "
						}
						rawCypher += fmt.Sprintf(` NOT (c)<-[:HAS]-(:FlowParticipant)<-[:HAS]-(:Flow {id:"%s"}) `, v)
					}
					flowFilter.Filters = append(flowFilter.Filters, utils.CreateRawCypherFilter(rawCypher))
				} else {
					createInOrEmptyStringFilter(filter, flowFilter, "id")
				}
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsFlowStatus) {
				createInOrEmptyStringFilter(filter, flowParticipantFilter, "status")
			}
			if filter.Filter.Property == string(postgresEntity.ColumnViewTypeContactsConnections) {
				// special case for not connected with
				if filter.Filter.Operation == commonmodel.ComparisonOperatorNotIn && filter.Filter.Value.ArrayStr != nil {
					rawCypher := ""
					for _, v := range *filter.Filter.Value.ArrayStr {
						if rawCypher != "" {
							rawCypher += " AND "
						}
						rawCypher += fmt.Sprintf(` NOT (c)-[:CONNECTED_WITH]->(:User{id:"%s"}) `, v)
					}
					flowFilter.Filters = append(flowFilter.Filters, utils.CreateRawCypherFilter(rawCypher))
				} else {
					createInOrEmptyStringFilter(filter, userConnectedFilter, "id")
				}
			}
		}

		if len(contactFilter.Filters) > 0 {
			contactFilterCypher, contactFilterParams = contactFilter.BuildCypherFilterFragmentWithParamName("c", "c_param_")
		}
		if len(linkedInFilter.Filters) > 0 {
			linkedInFilterCypher, linkedInFilterParams = linkedInFilter.BuildCypherFilterFragmentWithParamName("sl", "sl_param_")
		}
		if len(tagFilter.Filters) > 0 {
			tagFilterCypher, tagFilterParams = tagFilter.BuildCypherFilterFragmentWithParamName("t", "t_param_")
		}
		if len(locationFilter.Filters) > 0 {
			locationFilterCypher, locationFilterParams = locationFilter.BuildCypherFilterFragmentWithParamName("l", "l_param_")
		}
		if len(primaryEmailFilter.Filters) > 0 {
			primaryEmailFilterCypher, primaryEmailFilterParams = primaryEmailFilter.BuildCypherFilterFragmentWithParamName("pe", "pe_param_")
		}
		if len(primaryOrganizationFilter.Filters) > 0 {
			primaryOrganizationFilterCypher, primaryOrganizationFilterParams = primaryOrganizationFilter.BuildCypherFilterFragmentWithParamName("po", "po_param_")
		}
		if len(primaryJobRoleFilter.Filters) > 0 {
			primaryJobRoleFilterCypher, primaryJobRoleFilterParams = primaryJobRoleFilter.BuildCypherFilterFragmentWithParamName("pj", "pj_param_")
		}
		if len(phoneNumberFilter.Filters) > 0 {
			phoneNumberFilterCypher, phoneNumberFilterParams = phoneNumberFilter.BuildCypherFilterFragmentWithParamName("pn", "pn_param_")
		}
		if len(flowFilter.Filters) > 0 {
			flowFilterCypher, flowFilterParams = flowFilter.BuildCypherFilterFragmentWithParamName("f", "f_param_")
		}
		if len(flowParticipantFilter.Filters) > 0 {
			flowParticipantFilterCypher, flowParticipantFilterParams = flowParticipantFilter.BuildCypherFilterFragmentWithParamName("fc", "fc_param_")
		}
		if len(userConnectedFilter.Filters) > 0 {
			userConnectedFilterCypher, userConnectedFilterParams = userConnectedFilter.BuildCypherFilterFragmentWithParamName("uc", "uc_param_")
		}
	}

	//endregion

	params := map[string]any{
		"tenant": tenant,
		"limit":  limit,
	}

	utils.MergeMapToMap(contactFilterParams, params)
	utils.MergeMapToMap(linkedInFilterParams, params)
	utils.MergeMapToMap(tagFilterParams, params)
	utils.MergeMapToMap(locationFilterParams, params)
	utils.MergeMapToMap(primaryEmailFilterParams, params)
	utils.MergeMapToMap(primaryOrganizationFilterParams, params)
	utils.MergeMapToMap(primaryJobRoleFilterParams, params)
	utils.MergeMapToMap(phoneNumberFilterParams, params)
	utils.MergeMapToMap(flowFilterParams, params)
	utils.MergeMapToMap(flowParticipantFilterParams, params)
	utils.MergeMapToMap(userConnectedFilterParams, params)

	//region count selectQuery
	countQuery := ""
	{
		countQuery += fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact_%s) `, tenant)
		if linkedInFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (c)-[:HAS]->(sl:Social) WHERE sl.url CONTAINS 'linkedin.com/in' WITH *`
		}
		if tagFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (c)-[:TAGGED]->(t:Tag) WITH *`
		}
		if locationFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (c)-[:ASSOCIATED_WITH]->(l:Location) WITH *`
		}
		if primaryEmailFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (c)-[:HAS {primary:true}]->(pe:Email) WITH *`
		}
		if primaryOrganizationFilterCypher != "" || primaryJobRoleFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (c)--(pj:JobRole {primary:true})--(po:Organization {hide:false}) WITH *`
		}
		if phoneNumberFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (c)-[:HAS]->(pn:PhoneNumber) WITH *`
		}
		if flowFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (c)<-[:HAS]-(:FlowParticipant)<-[:HAS]-(f:Flow) WITH *`
		}
		if flowParticipantFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (c)<-[:HAS]-(fc:FlowParticipant)<-[:HAS]-(:Flow) WITH *`
		}
		if userConnectedFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (c)-[:CONNECTED_WITH]->(uc:User) WITH *`
		}

		countQuery += ` WHERE (c.hide = false OR c.hide IS NULL) `

		if contactFilterCypher != "" || linkedInFilterCypher != "" || tagFilterCypher != "" || locationFilterCypher != "" || primaryEmailFilterCypher != "" || primaryOrganizationFilterCypher != "" || primaryJobRoleFilterCypher != "" || phoneNumberFilterCypher != "" || flowFilterCypher != "" || flowParticipantFilterCypher != "" || userConnectedFilterCypher != "" {
			countQuery += " AND "
		}

		countQueryParts := []string{}
		if contactFilterCypher != "" {
			countQueryParts = append(countQueryParts, contactFilterCypher)
		}
		if linkedInFilterCypher != "" {
			countQueryParts = append(countQueryParts, linkedInFilterCypher)
		}
		if tagFilterCypher != "" {
			countQueryParts = append(countQueryParts, tagFilterCypher)
		}
		if locationFilterCypher != "" {
			countQueryParts = append(countQueryParts, locationFilterCypher)
		}
		if primaryEmailFilterCypher != "" {
			countQueryParts = append(countQueryParts, primaryEmailFilterCypher)
		}
		if primaryOrganizationFilterCypher != "" {
			countQueryParts = append(countQueryParts, primaryOrganizationFilterCypher)
		}
		if primaryJobRoleFilterCypher != "" {
			countQueryParts = append(countQueryParts, primaryJobRoleFilterCypher)
		}
		if phoneNumberFilterCypher != "" {
			countQueryParts = append(countQueryParts, phoneNumberFilterCypher)
		}
		if flowFilterCypher != "" {
			countQueryParts = append(countQueryParts, flowFilterCypher)
		}
		if flowParticipantFilterCypher != "" {
			countQueryParts = append(countQueryParts, flowParticipantFilterCypher)
		}
		if userConnectedFilterCypher != "" {
			countQueryParts = append(countQueryParts, userConnectedFilterCypher)
		}

		countQuery = countQuery + strings.Join(countQueryParts, " AND ") + fmt.Sprintf(` RETURN count(distinct(c))`)
	}
	//end count region

	selectQuery := ""
	//region selectQuery to fetch data
	{
		selectQuery += fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact_%s) `, tenant)
		if linkedInFilterCypher != "" || (sort != nil && (sort.By == string(postgresEntity.ColumnViewTypeContactsLinkedin) || sort.By == string(postgresEntity.ColumnViewTypeContactsLinkedinFollowerCount))) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)-[:HAS]->(sl:Social_%s) WHERE sl.url CONTAINS 'linkedin.com/in'  WITH *`, tenant)
		}
		if tagFilterCypher != "" {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)-[:TAGGED]->(t:Tag_%s) WITH *`, tenant)
		}
		if locationFilterCypher != "" || (sort != nil && (sort.By == string(postgresEntity.ColumnViewTypeContactsCountry) || sort.By == string(postgresEntity.ColumnViewTypeContactsCity) || sort.By == string(postgresEntity.ColumnViewTypeContactsRegion))) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH *`, tenant)
		}
		if primaryEmailFilterCypher != "" || (sort != nil && (sort.By == string(postgresEntity.ColumnViewTypeContactsPrimaryEmail))) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)-[:HAS {primary:true}]->(pe:Email_%s) WITH *`, tenant)
		}
		if primaryOrganizationFilterCypher != "" || primaryJobRoleFilterCypher != "" || (sort != nil && (sort.By == string(postgresEntity.ColumnViewTypeContactsOrganization) || sort.By == string(postgresEntity.ColumnViewTypeContactsJobTitle) || sort.By == string(postgresEntity.ColumnViewTypeContactsTimeInCurrentRole))) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)--(pj:JobRole_%s {primary:true})--(po:Organization_%s {hide:false}) WITH *`, tenant, tenant)
		}
		if phoneNumberFilterCypher != "" || (sort != nil && (sort.By == string(postgresEntity.ColumnViewTypeContactsPhoneNumbers))) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)-[:HAS]->(pn:PhoneNumber_%s) WITH *`, tenant)
		}
		if flowFilterCypher != "" || (sort != nil && (sort.By == string(postgresEntity.ColumnViewTypeContactsFlows))) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)<-[:HAS]-(:FlowParticipant_%s)<-[:HAS]-(f:Flow_%s) WITH *`, tenant, tenant)
		}
		if flowParticipantFilterCypher != "" || (sort != nil && (sort.By == string(postgresEntity.ColumnViewTypeContactsFlowStatus))) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)<-[:HAS]-(fc:FlowParticipant_%s)<-[:HAS]-(:Flow_%s) WITH *`, tenant, tenant)
		}
		if userConnectedFilterCypher != "" || (sort != nil && (sort.By == string(postgresEntity.ColumnViewTypeContactsConnections))) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)-[:CONNECTED_WITH]->(uc:User_%s) WITH *`, tenant)
		}

		selectQuery += ` WHERE (c.hide = false OR c.hide IS NULL) `

		if contactFilterCypher != "" || linkedInFilterCypher != "" || tagFilterCypher != "" || locationFilterCypher != "" || primaryEmailFilterCypher != "" || primaryOrganizationFilterCypher != "" || primaryJobRoleFilterCypher != "" || phoneNumberFilterCypher != "" || flowFilterCypher != "" || flowParticipantFilterCypher != "" || userConnectedFilterCypher != "" {
			selectQuery += " AND "
		}

		queryParts := []string{}
		if contactFilterCypher != "" {
			queryParts = append(queryParts, contactFilterCypher)
		}
		if linkedInFilterCypher != "" {
			queryParts = append(queryParts, linkedInFilterCypher)
		}
		if tagFilterCypher != "" {
			queryParts = append(queryParts, tagFilterCypher)
		}
		if locationFilterCypher != "" {
			queryParts = append(queryParts, locationFilterCypher)
		}
		if primaryEmailFilterCypher != "" {
			queryParts = append(queryParts, primaryEmailFilterCypher)
		}
		if primaryOrganizationFilterCypher != "" {
			queryParts = append(queryParts, primaryOrganizationFilterCypher)
		}
		if primaryJobRoleFilterCypher != "" {
			queryParts = append(queryParts, primaryJobRoleFilterCypher)
		}
		if phoneNumberFilterCypher != "" {
			queryParts = append(queryParts, phoneNumberFilterCypher)
		}
		if flowFilterCypher != "" {
			queryParts = append(queryParts, flowFilterCypher)
		}
		if flowParticipantFilterCypher != "" {
			queryParts = append(queryParts, flowParticipantFilterCypher)
		}
		if userConnectedFilterCypher != "" {
			queryParts = append(queryParts, userConnectedFilterCypher)
		}
		selectQuery = selectQuery + strings.Join(queryParts, " AND ")
	}
	//endregion

	// sort region
	aliases := ""

	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsName) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN (COALESCE(c.name, '') + COALESCE(c.firstName, '') + COALESCE(c.lastName, '')) <> '' THEN toLower(COALESCE(c.name, '') + COALESCE(c.firstName, '') + COALESCE(c.lastName, '')) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN (COALESCE(c.name, '') + COALESCE(c.firstName, '') + COALESCE(c.lastName, '')) <> '' THEN toLower(COALESCE(c.name, '') + COALESCE(c.firstName, '') + COALESCE(c.lastName, '')) ELSE '' END as SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsPrimaryEmail) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN (COALESCE(pe.email, '') + COALESCE(pe.rawEmail, '')) <> '' THEN toLower(COALESCE(pe.email, '') + COALESCE(pe.rawEmail, '')) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN (COALESCE(pe.email, '') + COALESCE(pe.rawEmail, '')) <> '' THEN toLower(COALESCE(pe.email, '') + COALESCE(pe.rawEmail, '')) ELSE '' END as SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsCountry) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN l.country <> '' AND NOT l.country IS NULL THEN toLower(l.country) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN l.country <> '' AND NOT l.country IS NULL THEN toLower(l.country) ELSE '' END AS SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsRegion) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN l.region <> '' AND NOT l.region IS NULL THEN toLower(l.region) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN l.region <> '' AND NOT l.region IS NULL THEN toLower(l.region) ELSE '' END AS SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsCity) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN l.locality <> '' AND NOT l.locality IS NULL THEN toLower(l.locality) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN l.locality <> '' AND NOT l.locality IS NULL THEN toLower(l.locality) ELSE '' END AS SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsCreatedAt) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN c.createdAt IS NOT NULL THEN c.createdAt ELSE datetime({year:2100}) END as SORT_BY `
		} else {
			aliases += `CASE WHEN c.createdAt IS NOT NULL THEN c.createdAt ELSE datetime({year:1900}) END as SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsUpdatedAt) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN c.updatedAt IS NOT NULL THEN c.updatedAt ELSE datetime({year:2100}) END as SORT_BY `
		} else {
			aliases += `CASE WHEN c.updatedAt IS NOT NULL THEN c.updatedAt ELSE datetime({year:1900}) END as SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsLinkedin) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN (COALESCE(sl.alias, '') + COALESCE(sl.url, '')) <> '' THEN toLower(COALESCE(sl.alias, '') + COALESCE(sl.url, '')) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN (COALESCE(sl.alias, '') + COALESCE(sl.url, '')) <> '' THEN toLower(COALESCE(sl.alias, '') + COALESCE(sl.url, '')) ELSE '' END as SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsLinkedinFollowerCount) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN sl.followersCount IS NULL THEN 999999999 ELSE sl.followersCount END as SORT_BY `
		} else {
			aliases += `CASE WHEN sl.followersCount IS NULL THEN -999999999 ELSE sl.followersCount END as SORT_BY `
		}

	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsOrganization) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN po.name <> '' AND NOT po.name IS NULL THEN toLower(po.name) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN po.name <> '' AND NOT po.name IS NULL THEN toLower(po.name) ELSE '' END AS SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsJobTitle) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN pj.jobTitle <> '' AND NOT pj.jobTitle IS NULL THEN toLower(pj.jobTitle) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN pj.jobTitle <> '' AND NOT pj.jobTitle IS NULL THEN toLower(pj.jobTitle) ELSE '' END AS SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsTimeInCurrentRole) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN pj.startedAt IS NOT NULL THEN pj.startedAt ELSE datetime({year:2100}) END as SORT_BY `
		} else {
			aliases += `CASE WHEN pj.startedAt IS NOT NULL THEN pj.startedAt ELSE datetime({year:1900}) END as SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsPhoneNumbers) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN pn.rawPhoneNumber <> '' AND NOT pn.rawPhoneNumber IS NULL THEN toLower(pn.rawPhoneNumber) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN pn.rawPhoneNumber <> '' AND NOT pn.rawPhoneNumber IS NULL THEN toLower(pn.rawPhoneNumber) ELSE '' END AS SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsFlows) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN f.name <> '' AND NOT f.name IS NULL THEN toLower(f.name) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN f.name <> '' AND NOT f.name IS NULL THEN toLower(f.name) ELSE '' END AS SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsFlowStatus) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN fc.status <> '' AND NOT fc.status IS NULL THEN toLower(fc.status) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN fc.status <> '' AND NOT fc.status IS NULL THEN toLower(fc.status) ELSE '' END AS SORT_BY `
		}
	}
	if sort != nil && sort.By == string(postgresEntity.ColumnViewTypeContactsConnections) {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN uc.name <> '' AND NOT uc.name IS NULL THEN toLower(uc.name) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN uc.name <> '' AND NOT uc.name IS NULL THEN toLower(uc.name) ELSE '' END AS SORT_BY `
		}
	}

	if len(aliases) > 0 {
		selectQuery += " WITH *, " + aliases
	} else {
		selectQuery += " WITH * "
	}

	cypherSort := utils.CypherSort{}
	if sort != nil && len(aliases) > 0 {
		selectQuery += " ORDER BY SORT_BY " + string(sort.Direction)
	} else {
		cypherSort.NewSortRule("UPDATED_AT", string(commonmodel.SortingDirectionDesc), false, reflect.TypeOf(neo4jentity.ContactEntity{}))
		selectQuery += string(cypherSort.SortingCypherFragment("c"))
	}

	// end sort region
	selectQuery += fmt.Sprintf(` RETURN distinct(c.id) `)
	selectQuery += fmt.Sprintf(` LIMIT $limit`)

	stringsWithTotalCount := &utils.StringsWithTotalCount{}

	var wg sync.WaitGroup
	var firstErr error
	var mu sync.Mutex

	setError := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
	}

	wg.Add(1)
	go func(ctx context.Context, result *utils.StringsWithTotalCount) {
		innerSpan, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardViewContactDataV2.CountQuery")
		defer innerSpan.Finish()
		defer wg.Done()
		tracing.SetDefaultServiceSpanTags(ctx, innerSpan)

		tracing.LogObjectAsJson(innerSpan, "params", params)
		innerSpan.LogFields(log.String("countQuery", countQuery))

		session := utils.NewNeo4jReadSession(ctx, *r.driver)
		defer session.Close(ctx)

		countRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			countQueryResult, err := tx.Run(ctx, countQuery, params)
			if err != nil {
				return nil, err
			} else {
				return utils.ExtractSingleRecordFirstValueAsType[int64](ctx, countQueryResult, err)
			}
		})

		if err != nil {
			tracing.TraceErr(innerSpan, err)
			setError(err)
			return
		}

		result.Count = countRecord.(int64)
	}(ctx, stringsWithTotalCount)

	wg.Add(1)
	go func(ctx context.Context, result *utils.StringsWithTotalCount) {
		innerSpan, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardViewContactDataV2.SelectQuery")
		defer span.Finish()
		defer wg.Done()
		tracing.SetDefaultServiceSpanTags(ctx, innerSpan)

		tracing.LogObjectAsJson(innerSpan, "params", params)
		innerSpan.LogFields(log.String("selectQuery", selectQuery))

		session := utils.NewNeo4jReadSession(ctx, *r.driver)
		defer session.Close(ctx)

		dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			queryResult, err := tx.Run(ctx, selectQuery, params)
			if err != nil {
				return nil, err
			} else {
				return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
			}
		})

		if err != nil {
			tracing.TraceErr(innerSpan, err)
			setError(err)
			return
		}

		result.Strings = dbRecords.([]string)
	}(ctx, stringsWithTotalCount)

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return stringsWithTotalCount, nil
}
