package repository

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
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
	userFilterCypher, userFilterParams := "", make(map[string]interface{})
	domainFilterCypher, domainFilterParams := "", make(map[string]interface{})
	parentOrganizationFilterCypher, parentOrganizationFilterParams := "", make(map[string]interface{})

	// ORGANIZATION, EMAIL, COUNTRY, REGION, LOCALITY
	// region organization filters
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
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsName.String() {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("name", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsWebsite.String() {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("website", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsPrimaryDomains.String() {
				domainFilter.Filters = append(domainFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.DomainPropertyDomain), filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsRelationship.String() {
				createInOrEmptyStringFilter(filter, organizationFilter, "relationship")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsOnboardingStatus.String() {
				createInOrEmptyStringFilter(filter, organizationFilter, "onboardingStatus")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsRenewalLikelihood.String() {
				createInOrEmptyStringFilter(filter, organizationFilter, "derivedRenewalLikelihood")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsRenewalDate.String() {
				createTimeFilter(filter, organizationFilter, "derivedNextRenewalAt")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsForecastArr.String() {
				createNumberCypherFilter(filter, organizationFilter, "renewalForecastArr")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsOwner.String() {
				createInOrEmptyStringFilter(filter, userFilter, "id")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsLastTouchpoint.String() {
				createInOrEmptyStringFilter(filter, organizationFilter, string(neo4jentity.OrganizationPropertyLastTouchpointType))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsLastTouchpointDate.String() {
				createTimeFilter(filter, organizationFilter, string(neo4jentity.OrganizationPropertyLastTouchpointAt))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsStage.String() {
				createInOrEmptyStringFilter(filter, organizationFilter, string(neo4jentity.OrganizationPropertyStage))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsSocials.String() {
				socialFilter.Filters = append(socialFilter.Filters, utils.CreateStringCypherFilter("url", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsLeadSource.String() {
				createInOrEmptyStringFilter(filter, organizationFilter, string(neo4jentity.OrganizationPropertyLeadSource))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsCreatedDate.String() {
				createTimeFilter(filter, organizationFilter, "createdAt")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsEmployeeCount.String() {
				createNumberCypherFilter(filter, organizationFilter, "employees")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsContactCount.String() {
				createNumberCypherFilter(filter, organizationFilter, "derivedContactCount")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsYearFounded.String() {
				createNumberCypherFilter(filter, organizationFilter, "yearFounded")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsIndustry.String() {
				createInOrEmptyStringFilter(filter, organizationFilter, string(neo4jentity.OrganizationPropertyIndustry))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsChurnDate.String() {
				createTimeFilter(filter, organizationFilter, "derivedChurnedAt")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsLtv.String() {
				createNumberCypherFilter(filter, organizationFilter, "derivedLtv")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsCountry.String() {
				createInOrEmptyStringFilter(filter, locationFilter, "countryCodeA2")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsCity.String() {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("locality", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsIsPublic.String() {
				createBooleanFilter(filter, organizationFilter, "isPublic")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsTags.String() {
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
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsHeadquarters.String() {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("headquarters", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsParentOrganization.String() {
				parentOrganizationFilter.Filters = append(parentOrganizationFilter.Filters, utils.CreateStringCypherFilter("name", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsUpdatedDate.String() {
				createTimeFilter(filter, organizationFilter, "updatedAt")
			}
		}

		if len(where.And) == 0 {
			for _, filter := range where.Or {
				organizationFilter.LogicalOperator = utils.OR
				if filter.Filter.Property == model.ColumnViewTypeOrganizationsName.String() {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("name", filter.Filter.Value.Str, filter.Filter.Operation))
				}
				if filter.Filter.Property == model.ColumnViewTypeOrganizationsWebsite.String() {
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

	// endregion

	params := map[string]any{
		"tenant": tenant,
		"limit":  limit,
	}

	utils.MergeMapToMap(organizationFilterParams, params)
	utils.MergeMapToMap(socialFilterParams, params)
	utils.MergeMapToMap(tagFilterParams, params)
	utils.MergeMapToMap(locationFilterParams, params)
	utils.MergeMapToMap(userFilterParams, params)
	utils.MergeMapToMap(domainFilterParams, params)
	utils.MergeMapToMap(parentOrganizationFilterParams, params)

	// region count selectQuery
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
		if locationFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (o)-[:ASSOCIATED_WITH]->(l:Location) WITH *`
		}
		if parentOrganizationFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (o)-[:SUBSIDIARY_OF]->(po:Organization) WITH *`
		}

		countQuery += ` WHERE (o.hide = false OR o.hide IS NULL) `

		if organizationFilterCypher != "" || domainFilterCypher != "" || socialFilterCypher != "" || tagFilterCypher != "" || locationFilterCypher != "" || parentOrganizationFilterCypher != "" || userFilterCypher != "" {
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
		if locationFilterCypher != "" {
			countQueryParts = append(countQueryParts, locationFilterCypher)
		}
		if parentOrganizationFilterCypher != "" {
			countQueryParts = append(countQueryParts, parentOrganizationFilterCypher)
		}

		countQuery = countQuery + strings.Join(countQueryParts, " AND ") + fmt.Sprintf(` RETURN count(distinct(o))`)
	}
	// end count region

	selectQuery := ""
	// region selectQuery to fetch data
	{
		selectQuery += fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s) `, tenant)
		if userFilterCypher != "" || (sort != nil && (sort.By == model.ColumnViewTypeOrganizationsOwner.String())) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)<-[:OWNS]-(u:User) WITH *`)
		}
		if domainFilterCypher != "" || (sort != nil && (sort.By == model.ColumnViewTypeOrganizationsPrimaryDomains.String())) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[:HAS_DOMAIN]->(d:Domain{primary: true}) WITH *`)
		}
		if socialFilterCypher != "" {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[:HAS]->(s:Social_%s) WITH *`, tenant)
		}
		if tagFilterCypher != "" {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[:TAGGED]->(t:Tag_%s) WITH *`, tenant)
		}
		if locationFilterCypher != "" || (sort != nil && (sort.By == model.ColumnViewTypeOrganizationsCountry.String() || sort.By == model.ColumnViewTypeOrganizationsCity.String())) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH *`, tenant)
		}
		if parentOrganizationFilterCypher != "" || sort != nil && sort.By == model.ColumnViewTypeOrganizationsParentOrganization.String() {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[:SUBSIDIARY_OF]->(po:Organization_%s) WITH *`, tenant)
		}
		selectQuery += ` WHERE (o.hide = false OR o.hide IS NULL) `

		if organizationFilterCypher != "" || domainFilterCypher != "" || socialFilterCypher != "" || tagFilterCypher != "" || parentOrganizationFilterCypher != "" || locationFilterCypher != "" || userFilterCypher != "" {
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
		if locationFilterCypher != "" {
			queryParts = append(queryParts, locationFilterCypher)
		}
		if parentOrganizationFilterCypher != "" {
			queryParts = append(queryParts, parentOrganizationFilterCypher)
		}
		selectQuery = selectQuery + strings.Join(queryParts, " AND ")
	}
	// endregion

	// sort region
	sortingCypher := ""
	aliases := ""

	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsName.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.name <> \"\" and not o.name is null THEN toLower(o.name) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.name <> \"\" and not o.name is null THEN toLower(o.name) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsWebsite.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.website <> \"\" and not o.website is null THEN toLower(o.website) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.website <> \"\" and not o.website is null THEN toLower(o.website) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsPrimaryDomains.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN d.domain <> \"\" and not d.domain is null THEN toLower(d.domain) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN d.domain <> \"\" and not d.domain is null THEN toLower(d.domain) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsRelationship.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.relationship <> \"\" and not o.relationship is null THEN toLower(o.relationship) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.relationship <> \"\" and not o.relationship is null THEN toLower(o.relationship) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsOnboardingStatus.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.onboardingStatusOrder <> \"\" and not o.onboardingStatusOrder is null THEN o.onboardingStatusOrder ELSE 9999 END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.onboardingStatusOrder <> \"\" and not o.onboardingStatusOrder is null THEN o.onboardingStatusOrder ELSE -1 END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsRenewalLikelihood.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.derivedRenewalLikelihoodOrder <> \"\" and not o.derivedRenewalLikelihoodOrder is null THEN o.derivedRenewalLikelihoodOrder ELSE 9999 END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.derivedRenewalLikelihoodOrder <> \"\" and not o.derivedRenewalLikelihoodOrder is null THEN o.derivedRenewalLikelihoodOrder ELSE -1 END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsRenewalDate.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.derivedNextRenewalAt <> \"\" and not o.derivedNextRenewalAt is null THEN o.derivedNextRenewalAt ELSE datetime({year:2100}) END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.derivedNextRenewalAt <> \"\" and not o.derivedNextRenewalAt is null THEN o.derivedNextRenewalAt ELSE datetime({year:1900}) END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsForecastArr.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.renewalForecastArr <> \"\" and not o.renewalForecastArr is null THEN o.renewalForecastArr ELSE 9999 END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.renewalForecastArr <> \"\" and not o.renewalForecastArr is null THEN o.renewalForecastArr ELSE -1 END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsOwner.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN (COALESCE(u.name, '') + COALESCE(u.firstName, '') + COALESCE(u.lastName, '')) <> '' THEN toLower(COALESCE(u.name, '') + COALESCE(u.firstName, '') + COALESCE(u.lastName, '')) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN (COALESCE(u.name, '') + COALESCE(u.firstName, '') + COALESCE(u.lastName, '')) <> '' THEN toLower(COALESCE(u.name, '') + COALESCE(u.firstName, '') + COALESCE(u.lastName, '')) ELSE '' END as SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsLastTouchpoint.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.lastTouchpointAt <> \"\" and not o.lastTouchpointAt is null THEN o.lastTouchpointAt ELSE datetime({year:2100}) END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.lastTouchpointAt <> \"\" and not o.lastTouchpointAt is null THEN o.lastTouchpointAt ELSE datetime({year:1900}) END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsLastTouchpointDate.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.lastTouchpointAt <> \"\" and not o.lastTouchpointAt is null THEN o.lastTouchpointAt ELSE datetime({year:2100}) END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.lastTouchpointAt <> \"\" and not o.lastTouchpointAt is null THEN o.lastTouchpointAt ELSE datetime({year:1900}) END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsStage.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.stage <> \"\" and not o.stage is null THEN toLower(o.stage) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.stage <> \"\" and not o.stage is null THEN toLower(o.stage) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsLeadSource.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.leadSource <> \"\" and not o.leadSource is null THEN toLower(o.leadSource) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.leadSource <> \"\" and not o.leadSource is null THEN toLower(o.leadSource) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsCreatedDate.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN o.createdAt IS NOT NULL THEN o.createdAt ELSE datetime({year:2100}) END as SORT_BY `
		} else {
			aliases += `CASE WHEN o.createdAt IS NOT NULL THEN o.createdAt ELSE datetime({year:1900}) END as SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsEmployeeCount.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.employees <> \"\" and not o.employees is null THEN o.employees ELSE 999999999 END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.employees <> \"\" and not o.employees is null THEN o.employees ELSE -1 END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsContactCount.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN o.derivedContactCount IS NOT NULL THEN o.derivedContactCount ELSE 999999999 END as SORT_BY `
		} else {
			aliases += `CASE WHEN o.derivedContactCount IS NOT NULL THEN o.derivedContactCount ELSE -1 END as SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsYearFounded.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.yearFounded <> \"\" and not o.yearFounded is null THEN o.yearFounded ELSE 999999999 END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.yearFounded <> \"\" and not o.yearFounded is null THEN o.yearFounded ELSE -1 END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsIndustry.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.industry <> \"\" and not o.industry is null THEN toLower(o.industry) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.industry <> \"\" and not o.industry is null THEN toLower(o.industry) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsChurnDate.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.derivedChurnedAt <> \"\" and not o.derivedChurnedAt is null THEN o.derivedChurnedAt ELSE datetime({year:2100}) END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.derivedChurnedAt <> \"\" and not o.derivedChurnedAt is null THEN o.derivedChurnedAt ELSE datetime({year:1900}) END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsLtv.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.derivedLtv <> \"\" and not o.derivedLtv is null THEN o.derivedLtv ELSE 9999999999999999 END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.derivedLtv <> \"\" and not o.derivedLtv is null THEN o.derivedLtv ELSE -9999999999999999 END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsCountry.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN l.country <> \"\" and not l.country is null THEN toLower(l.country) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN l.country <> \"\" and not l.country is null THEN toLower(l.country) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsCity.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN l.locality <> \"\" and not l.locality is null THEN toLower(l.locality) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN l.locality <> \"\" and not l.locality is null THEN toLower(l.locality) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsIsPublic.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN o.isPublic = true THEN 0 ELSE CASE WHEN o.isPublic = false THEN 1 ELSE 2 END END as SORT_BY "
		} else {
			aliases += "CASE WHEN o.isPublic = false THEN 2 ELSE CASE WHEN o.isPublic = true THEN 1 ELSE 0 END END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsParentOrganization.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += "CASE WHEN po.name <> \"\" and not po.name is null THEN toLower(po.name) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY "
		} else {
			aliases += "CASE WHEN po.name <> \"\" and not po.name is null THEN toLower(po.name) ELSE '' END as SORT_BY "
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeOrganizationsUpdatedDate.String() {
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

		span.LogFields(log.Object("params", params))
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

		span.LogFields(log.Object("params", params))
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

		for _, filter := range where.And {
			// TODO
			// ColumnViewTypeContactsEmails                     ColumnViewType = "CONTACTS_EMAILS"
			// ColumnViewTypeContactsPersonalEmails             ColumnViewType = "CONTACTS_PERSONAL_EMAILS"
			// ColumnViewTypeContactsPhoneNumbers               ColumnViewType = "CONTACTS_PHONE_NUMBERS"
			// ColumnViewTypeContactsPersona                    ColumnViewType = "CONTACTS_PERSONA"
			// ColumnViewTypeContactsLastInteraction            ColumnViewType = "CONTACTS_LAST_INTERACTION"
			// ColumnViewTypeContactsSkills                     ColumnViewType = "CONTACTS_SKILLS"
			// ColumnViewTypeContactsSchools                    ColumnViewType = "CONTACTS_SCHOOLS"
			// ColumnViewTypeContactsLanguages                  ColumnViewType = "CONTACTS_LANGUAGES"
			// ColumnViewTypeContactsTimeInCurrentRole          ColumnViewType = "CONTACTS_TIME_IN_CURRENT_ROLE"
			// ColumnViewTypeContactsExperience                 ColumnViewType = "CONTACTS_EXPERIENCE"
			// ColumnViewTypeContactsConnections                ColumnViewType = "CONTACTS_CONNECTIONS"
			// ColumnViewTypeContactsFlows                      ColumnViewType = "CONTACTS_FLOWS"
			// ColumnViewTypeContactsFlowStatus                 ColumnViewType = "CONTACTS_FLOW_STATUS"
			// ColumnViewTypeContactsFlowNextAction             ColumnViewType = "CONTACTS_FLOW_NEXT_ACTION"
			if filter.Filter.Property == model.ColumnViewTypeContactsName.String() {
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
			if filter.Filter.Property == model.ColumnViewTypeContactsPrimaryEmail.String() {
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
			if filter.Filter.Property == model.ColumnViewTypeContactsCountry.String() {
				createInOrEmptyStringFilter(filter, locationFilter, "countryCodeA2")
			}
			if filter.Filter.Property == model.ColumnViewTypeContactsCity.String() {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.LocationPropertyLocality), filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeContactsRegion.String() {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.LocationPropertyRegion), filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeContactsTags.String() {
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
			if filter.Filter.Property == model.ColumnViewTypeContactsLinkedin.String() {
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
			if filter.Filter.Property == model.ColumnViewTypeContactsLinkedinFollowerCount.String() {
				createNumberCypherFilter(filter, linkedInFilter, string(neo4jentity.SocialPropertyFollowersCount))
			}
			if filter.Filter.Property == model.ColumnViewTypeContactsOrganization.String() {
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
			if filter.Filter.Property == model.ColumnViewTypeContactsJobTitle.String() {
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
	}

	// endregion

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

	// region count selectQuery
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

		countQuery += ` WHERE (c.hide = false OR c.hide IS NULL) `

		if contactFilterCypher != "" || linkedInFilterCypher != "" || tagFilterCypher != "" || locationFilterCypher != "" || primaryEmailFilterCypher != "" || primaryOrganizationFilterCypher != "" || primaryJobRoleFilterCypher != "" {
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

		countQuery = countQuery + strings.Join(countQueryParts, " AND ") + fmt.Sprintf(` RETURN count(distinct(c))`)
	}
	// end count region

	selectQuery := ""
	// region selectQuery to fetch data
	{
		selectQuery += fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact_%s) `, tenant)
		if linkedInFilterCypher != "" || (sort != nil && (sort.By == model.ColumnViewTypeContactsLinkedin.String() || sort.By == model.ColumnViewTypeContactsLinkedinFollowerCount.String())) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)-[:HAS]->(sl:Social_%s) WHERE sl.url CONTAINS 'linkedin.com/in'  WITH *`, tenant)
		}
		if tagFilterCypher != "" {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)-[:TAGGED]->(t:Tag_%s) WITH *`, tenant)
		}
		if locationFilterCypher != "" || (sort != nil && (sort.By == model.ColumnViewTypeContactsCountry.String() || sort.By == model.ColumnViewTypeContactsCity.String() || sort.By == model.ColumnViewTypeContactsRegion.String())) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH *`, tenant)
		}
		if primaryEmailFilterCypher != "" || (sort != nil && (sort.By == model.ColumnViewTypeContactsPrimaryEmail.String())) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)-[:HAS {primary:true}]->(pe:Email_%s) WITH *`, tenant)
		}
		if primaryOrganizationFilterCypher != "" || primaryJobRoleFilterCypher != "" || (sort != nil && (sort.By == model.ColumnViewTypeContactsOrganization.String() || sort.By == model.ColumnViewTypeContactsJobTitle.String())) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (c)--(pj:JobRole_%s {primary:true})--(po:Organization_%s {hide:false}) WITH *`, tenant, tenant)
		}
		selectQuery += ` WHERE (c.hide = false OR c.hide IS NULL) `

		if contactFilterCypher != "" || linkedInFilterCypher != "" || tagFilterCypher != "" || locationFilterCypher != "" || primaryEmailFilterCypher != "" || primaryOrganizationFilterCypher != "" || primaryJobRoleFilterCypher != "" {
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
		selectQuery = selectQuery + strings.Join(queryParts, " AND ")
	}
	// endregion

	// sort region
	aliases := ""

	if sort != nil && sort.By == model.ColumnViewTypeContactsName.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN (COALESCE(c.name, '') + COALESCE(c.firstName, '') + COALESCE(c.lastName, '')) <> '' THEN toLower(COALESCE(c.name, '') + COALESCE(c.firstName, '') + COALESCE(c.lastName, '')) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN (COALESCE(c.name, '') + COALESCE(c.firstName, '') + COALESCE(c.lastName, '')) <> '' THEN toLower(COALESCE(c.name, '') + COALESCE(c.firstName, '') + COALESCE(c.lastName, '')) ELSE '' END as SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeContactsPrimaryEmail.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN (COALESCE(pe.email, '') + COALESCE(pe.rawEmail, '')) <> '' THEN toLower(COALESCE(pe.email, '') + COALESCE(pe.rawEmail, '')) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN (COALESCE(pe.email, '') + COALESCE(pe.rawEmail, '')) <> '' THEN toLower(COALESCE(pe.email, '') + COALESCE(pe.rawEmail, '')) ELSE '' END as SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeContactsCountry.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN l.country <> '' AND NOT l.country IS NULL THEN toLower(l.country) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN l.country <> '' AND NOT l.country IS NULL THEN toLower(l.country) ELSE '' END AS SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeContactsRegion.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN l.region <> '' AND NOT l.region IS NULL THEN toLower(l.region) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN l.region <> '' AND NOT l.region IS NULL THEN toLower(l.region) ELSE '' END AS SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeContactsCity.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN l.locality <> '' AND NOT l.locality IS NULL THEN toLower(l.locality) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN l.locality <> '' AND NOT l.locality IS NULL THEN toLower(l.locality) ELSE '' END AS SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeContactsCreatedAt.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN c.createdAt IS NOT NULL THEN c.createdAt ELSE datetime({year:2100}) END as SORT_BY `
		} else {
			aliases += `CASE WHEN c.createdAt IS NOT NULL THEN c.createdAt ELSE datetime({year:1900}) END as SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeContactsUpdatedAt.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN c.updatedAt IS NOT NULL THEN c.updatedAt ELSE datetime({year:2100}) END as SORT_BY `
		} else {
			aliases += `CASE WHEN c.updatedAt IS NOT NULL THEN c.updatedAt ELSE datetime({year:1900}) END as SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeContactsLinkedin.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN (COALESCE(sl.alias, '') + COALESCE(sl.url, '')) <> '' THEN toLower(COALESCE(sl.alias, '') + COALESCE(sl.url, '')) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN (COALESCE(sl.alias, '') + COALESCE(sl.url, '')) <> '' THEN toLower(COALESCE(sl.alias, '') + COALESCE(sl.url, '')) ELSE '' END as SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeContactsLinkedinFollowerCount.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN sl.followersCount IS NULL THEN 999999999 ELSE sl.followersCount END as SORT_BY `
		} else {
			aliases += `CASE WHEN sl.followersCount IS NULL THEN -999999999 ELSE sl.followersCount END as SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeContactsOrganization.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN po.name <> '' AND NOT po.name IS NULL THEN toLower(po.name) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN po.name <> '' AND NOT po.name IS NULL THEN toLower(po.name) ELSE '' END AS SORT_BY `
		}
	}
	if sort != nil && sort.By == model.ColumnViewTypeContactsJobTitle.String() {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			aliases += `CASE WHEN pj.jobTitle <> '' AND NOT pj.jobTitle IS NULL THEN toLower(pj.jobTitle) ELSE 'zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz' END as SORT_BY `
		} else {
			aliases += `CASE WHEN pj.jobTitle <> '' AND NOT pj.jobTitle IS NULL THEN toLower(pj.jobTitle) ELSE '' END AS SORT_BY `
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

		innerSpan.LogFields(log.Object("params", params))
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

		innerSpan.LogFields(log.Object("params", params))
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
