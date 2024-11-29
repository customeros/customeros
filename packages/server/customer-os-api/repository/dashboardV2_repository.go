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
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
	"strings"
	"sync"
)

type DashboardV2Repository interface {
	GetDashboardViewOrganizationDataV2(ctx context.Context, tenant string, limit int, where *model.Filter, sort *commonmodel.SortBy) (*utils.StringsWithTotalCount, error)
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

	ownerId := []string{}
	ownerIncludeEmpty := false

	//ORGANIZATION, EMAIL, COUNTRY, REGION, LOCALITY
	//region organization filters
	if where != nil {
		organizationFilter := new(utils.CypherFilter)
		organizationFilter.Negate = false
		organizationFilter.LogicalOperator = utils.AND
		organizationFilter.Details = &utils.CypherFilterItem{
			SupportCaseSensitive: true,
		}
		organizationFilter.Filters = make([]*utils.CypherFilter, 0)

		socialFilter := new(utils.CypherFilter)
		socialFilter.Negate = false
		socialFilter.LogicalOperator = utils.AND
		socialFilter.Details = &utils.CypherFilterItem{
			SupportCaseSensitive: true,
		}
		socialFilter.Filters = make([]*utils.CypherFilter, 0)

		tagFilter := new(utils.CypherFilter)
		tagFilter.Negate = false
		tagFilter.LogicalOperator = utils.AND
		tagFilter.Details = &utils.CypherFilterItem{
			SupportCaseSensitive: true,
		}
		tagFilter.Filters = make([]*utils.CypherFilter, 0)

		emailFilter := new(utils.CypherFilter)
		emailFilter.Negate = false
		emailFilter.LogicalOperator = utils.OR
		emailFilter.Details = &utils.CypherFilterItem{
			SupportCaseSensitive: true,
		}
		emailFilter.Filters = make([]*utils.CypherFilter, 0)

		locationFilter := new(utils.CypherFilter)
		locationFilter.Negate = false
		locationFilter.LogicalOperator = utils.OR
		locationFilter.Details = &utils.CypherFilterItem{
			SupportCaseSensitive: true,
		}
		locationFilter.Filters = make([]*utils.CypherFilter, 0)

		for _, filter := range where.And {
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsName.String() {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("name", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsWebsite.String() {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("website", filter.Filter.Value.Str, filter.Filter.Operation))
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
				createBetweenOrEmptyTimeFilter(filter, organizationFilter, "derivedNextRenewalAt")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsLastTouchpointDate.String() {
				createBetweenOrEmptyTimeFilter(filter, organizationFilter, "lastTouchpointAt")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsStage.String() {
				createInOrEmptyStringFilter(filter, organizationFilter, "stage")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsSocials.String() {
				socialFilter.Filters = append(socialFilter.Filters, utils.CreateStringCypherFilter("url", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsLeadSource.String() {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("leadSource", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsCreatedDate.String() {
				createBetweenOrEmptyTimeFilter(filter, organizationFilter, "createdAt")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsEmployeeCount.String() {
				createNumberCypherFilter(filter, organizationFilter, "employees")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsYearFounded.String() {
				createNumberCypherFilter(filter, organizationFilter, "yearFounded")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsIndustry.String() {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("industry", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsChurnDate.String() {
				createBetweenOrEmptyTimeFilter(filter, organizationFilter, "derivedChurnedAt")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsLtv.String() {
				createNumberCypherFilter(filter, organizationFilter, "derivedLtv")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsCountry.String() {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("country", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsCity.String() {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("locality", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsIsPublic.String() {
				createBooleanFilter(filter, organizationFilter, "isPublic")
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsTags.String() {

				cf := utils.CypherFilter{}

				if filter.Filter.Operation == commonmodel.ComparisonOperatorIsEmpty {
					cf.Details = new(utils.CypherFilterItem)
					cf.Details.NodeProperty = "trCount = 0"
					cf.Details.ComparisonOperator = commonmodel.ComparisonOperatorCountRelation
				} else if filter.Filter.Operation == commonmodel.ComparisonOperatorIsNotEmpty {
					cf.Details = new(utils.CypherFilterItem)
					cf.Details.NodeProperty = "trCount > 0"
					cf.Details.ComparisonOperator = commonmodel.ComparisonOperatorCountRelation
				} else {
					cf.Details = new(utils.CypherFilterItem)
					cf.Details.NodeProperty = "name"
					cf.Details.Value = filter.Filter.Value.Str
					cf.Details.ComparisonOperator = filter.Filter.Operation
				}

				tagFilter.Filters = append(tagFilter.Filters, &cf)
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsHeadquarters.String() {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("headquarters", filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == model.ColumnViewTypeOrganizationsUpdatedDate.String() {
				createBetweenOrEmptyTimeFilter(filter, organizationFilter, "updatedAt")
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
	}

	//endregion

	params := map[string]any{
		"tenant":  tenant,
		"ownerId": ownerId,
		"limit":   limit,
	}

	utils.MergeMapToMap(organizationFilterParams, params)
	utils.MergeMapToMap(socialFilterParams, params)
	utils.MergeMapToMap(tagFilterParams, params)
	utils.MergeMapToMap(locationFilterParams, params)

	//region count selectQuery
	countQuery := ""
	{
		countQuery += fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s) `, tenant)
		if len(ownerId) > 0 || ownerIncludeEmpty {
			countQuery += ` OPTIONAL MATCH (o)<-[:OWNS]-(owner:User) WITH *`
		}
		if socialFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (o)-[:HAS]->(s:Social) WITH *`
		}
		if tagFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (o)-[tr:TAGGED]->(t:Tag) WITH *, count(tr) as trCount`
		}
		if locationFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (o)-[:ASSOCIATED_WITH]->(l:Location) WITH *`
		}

		countQuery += ` WHERE (o.hide = false OR o.hide IS NULL) `

		if organizationFilterCypher != "" || socialFilterCypher != "" || tagFilterCypher != "" || locationFilterCypher != "" || len(ownerId) > 0 || ownerIncludeEmpty {
			countQuery += " AND "
		}

		countQueryParts := []string{}
		if organizationFilterCypher != "" {
			countQueryParts = append(countQueryParts, organizationFilterCypher)
		}
		if len(ownerId) > 0 || ownerIncludeEmpty {
			if len(ownerId) == 0 {
				countQueryParts = append(countQueryParts, fmt.Sprintf(` owner.id IS NULL `))
			} else if ownerIncludeEmpty {
				countQueryParts = append(countQueryParts, fmt.Sprintf(` (owner.id IN $ownerId OR owner.id IS NULL) `))
			} else {
				countQueryParts = append(countQueryParts, fmt.Sprintf(` owner.id IN $ownerId `))
			}
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

		countQuery = countQuery + strings.Join(countQueryParts, " AND ") + fmt.Sprintf(` RETURN count(distinct(o))`)
	}
	//end count region

	selectQuery := ""
	//region selectQuery to fetch data
	{
		selectQuery += fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s) `, tenant)
		if len(ownerId) > 0 || ownerIncludeEmpty {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)<-[:OWNS]-(owner:User) WITH *`)
		}
		if socialFilterCypher != "" {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[:HAS]->(s:Social_%s) WITH *`, tenant)
		}
		if tagFilterCypher != "" {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[tr:TAGGED]->(t:Tag_%s) WITH *, count(tr) as trCount`, tenant)
		}
		if locationFilterCypher != "" {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH *`, tenant)
		}
		if sort != nil && sort.By == SearchSortParamOwner {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (o)<-[:OWNS]-(owner:User_%s) WITH *`, tenant)
		}
		selectQuery += ` WHERE (o.hide = false OR o.hide IS NULL) `

		if organizationFilterCypher != "" || socialFilterCypher != "" || tagFilterCypher != "" || locationFilterCypher != "" || len(ownerId) > 0 || ownerIncludeEmpty {
			selectQuery += " AND "
		}

		queryParts := []string{}
		if organizationFilterCypher != "" {
			queryParts = append(queryParts, organizationFilterCypher)
		}
		if len(ownerId) > 0 || ownerIncludeEmpty {
			if len(ownerId) == 0 {
				queryParts = append(queryParts, fmt.Sprintf(` owner.id IS NULL `))
			} else if ownerIncludeEmpty {
				queryParts = append(queryParts, fmt.Sprintf(` (owner.id IN $ownerId OR owner.id IS NULL) `))
			} else {
				queryParts = append(queryParts, fmt.Sprintf(` owner.id IN $ownerId `))
			}
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
		selectQuery = selectQuery + strings.Join(queryParts, " AND ")
	}
	//endregion

	// sort region
	//aliases := " o, d, l"
	selectQuery += " WITH o "
	//if sort != nil && sort.By == SearchSortParamOwner {
	//	if sort.Direction == commonmodel.SortingDirectionAsc {
	//		selectQuery += ", CASE WHEN owner.firstName <> \"\" and not owner.firstName is null THEN owner.firstName ELSE 'ZZZZZZZZZZZZZZZZZZZ' END as OWNER_FIRST_NAME_FOR_SORTING "
	//		selectQuery += ", CASE WHEN owner.lastName <> \"\" and not owner.lastName is null THEN owner.lastName ELSE 'ZZZZZZZZZZZZZZZZZZZ' END as OWNER_LAST_NAME_FOR_SORTING "
	//	} else {
	//		selectQuery += ", CASE WHEN owner.firstName <> \"\" and not owner.firstName is null THEN owner.firstName ELSE 'AAAAAAAAAAAAAAAAAAA' END as OWNER_FIRST_NAME_FOR_SORTING "
	//		selectQuery += ", CASE WHEN owner.lastName <> \"\" and not owner.lastName is null THEN owner.lastName ELSE 'AAAAAAAAAAAAAAAAAAA' END as OWNER_LAST_NAME_FOR_SORTING "
	//	}
	//	aliases += ", OWNER_FIRST_NAME_FOR_SORTING, OWNER_LAST_NAME_FOR_SORTING "
	//}
	//if sort != nil && sort.By == SearchSortParamName {
	//	if sort.Direction == commonmodel.SortingDirectionAsc {
	//		selectQuery += ", CASE WHEN o.name <> \"\" and not o.name is null THEN o.name ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as NAME_FOR_SORTING "
	//	} else {
	//		selectQuery += ", o.name as NAME_FOR_SORTING "
	//	}
	//	aliases += ", NAME_FOR_SORTING "
	//}
	//if sort != nil && sort.By == SearchSortParamRenewalLikelihood {
	//	if sort.Direction == commonmodel.SortingDirectionAsc {
	//		selectQuery += ", CASE WHEN o.derivedRenewalLikelihoodOrder IS NOT NULL THEN o.derivedRenewalLikelihoodOrder ELSE 9999 END as RENEWAL_LIKELIHOOD_FOR_SORTING "
	//	} else {
	//		selectQuery += ", CASE WHEN o.derivedRenewalLikelihoodOrder IS NOT NULL THEN o.derivedRenewalLikelihoodOrder ELSE -1 END as RENEWAL_LIKELIHOOD_FOR_SORTING "
	//	}
	//	aliases += ", RENEWAL_LIKELIHOOD_FOR_SORTING "
	//}
	//if sort != nil && sort.By == SearchSortParamRelationship {
	//	if sort.Direction == commonmodel.SortingDirectionAsc {
	//		selectQuery += ", CASE WHEN o.relationship <> '' AND NOT o.relationship IS NULL THEN o.relationship ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as RELATIONSHIP_FOR_SORTING "
	//	} else {
	//		selectQuery += ", o.relationship as RELATIONSHIP_FOR_SORTING "
	//	}
	//	aliases += ", RELATIONSHIP_FOR_SORTING "
	//}
	//if sort != nil && sort.By == SearchSortParamStage {
	//	if sort.Direction == commonmodel.SortingDirectionAsc {
	//		selectQuery += ", CASE WHEN o.stage <> '' AND NOT o.stage IS NULL THEN o.stage ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as STAGE_FOR_SORTING "
	//	} else {
	//		selectQuery += ", o.stage as STAGE_FOR_SORTING "
	//	}
	//	aliases += ", STAGE_FOR_SORTING "
	//}
	//if sort != nil && sort.By == SearchSortParamIndustry {
	//	if sort.Direction == commonmodel.SortingDirectionAsc {
	//		selectQuery += ", CASE WHEN o.industry <> '' AND NOT o.industry IS NULL THEN o.industry ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as STAGE_FOR_SORTING "
	//	} else {
	//		selectQuery += ", o.industry as INDUSTRY_FOR_SORTING "
	//	}
	//	aliases += ", INDUSTRY_FOR_SORTING "
	//}
	//if sort != nil && sort.By == SearchSortParamRenewalCycleNext {
	//	if sort.Direction == commonmodel.SortingDirectionAsc {
	//		selectQuery += ", CASE WHEN o.billingDetailsRenewalCycleNext IS NOT NULL THEN date(o.billingDetailsRenewalCycleNext) ELSE date('2100-01-01') END as RENEWAL_CYCLE_NEXT_FOR_SORTING "
	//	} else {
	//		selectQuery += ", CASE WHEN o.billingDetailsRenewalCycleNext IS NOT NULL THEN date(o.billingDetailsRenewalCycleNext) ELSE date('1900-01-01') END as RENEWAL_CYCLE_NEXT_FOR_SORTING "
	//	}
	//	aliases += ", RENEWAL_CYCLE_NEXT_FOR_SORTING "
	//}
	//if sort != nil && sort.By == SearchSortParamRenewalDate {
	//	if sort.Direction == commonmodel.SortingDirectionAsc {
	//		selectQuery += ", CASE WHEN o.derivedNextRenewalAt IS NOT NULL THEN date(o.derivedNextRenewalAt) ELSE date('2100-01-01') END as RENEWAL_DATE_FOR_SORTING "
	//	} else {
	//		selectQuery += ", CASE WHEN o.derivedNextRenewalAt IS NOT NULL THEN date(o.derivedNextRenewalAt) ELSE date('1900-01-01') END as RENEWAL_DATE_FOR_SORTING "
	//	}
	//	aliases += ", RENEWAL_DATE_FOR_SORTING "
	//}
	//if sort != nil && sort.By == SearchSortParamChurnDate {
	//	if sort.Direction == commonmodel.SortingDirectionAsc {
	//		selectQuery += ", CASE WHEN o.derivedChurnedAt IS NOT NULL THEN date(o.derivedChurnedAt) ELSE date('2100-01-01') END as CHURN_DATE_FOR_SORTING "
	//	} else {
	//		selectQuery += ", CASE WHEN o.derivedChurnedAt IS NOT NULL THEN date(o.derivedChurnedAt) ELSE date('1900-01-01') END as CHURN_DATE_FOR_SORTING "
	//	}
	//	aliases += ", RENEWAL_DATE_FOR_SORTING "
	//}
	//if sort != nil && sort.By == SearchSortParamOnboardingStatus {
	//	if sort.Direction == commonmodel.SortingDirectionAsc {
	//		selectQuery += ", CASE WHEN o.onboardingStatusOrder IS NOT NULL THEN o.onboardingStatusOrder ELSE 9999 END as ONBOARDING_STATUS_FOR_SORTING "
	//		selectQuery += ", o.onboardingUpdatedAt AS ONBOARDING_UPDATED_AT_FOR_SORTING "
	//	} else {
	//		selectQuery += ", CASE WHEN o.onboardingStatusOrder IS NOT NULL THEN o.onboardingStatusOrder ELSE -1 END as ONBOARDING_STATUS_FOR_SORTING "
	//		selectQuery += ", o.onboardingUpdatedAt AS ONBOARDING_UPDATED_AT_FOR_SORTING "
	//	}
	//	aliases += ", ONBOARDING_STATUS_FOR_SORTING, ONBOARDING_UPDATED_AT_FOR_SORTING "
	//}
	//if sort != nil && sort.By == SearchSortParamForecastArr {
	//	if sort.Direction == commonmodel.SortingDirectionAsc {
	//		selectQuery += ", CASE WHEN o.renewalForecastArr <> \"\" and o.renewalForecastArr IS NOT NULL THEN o.renewalForecastArr ELSE 9999999999999999 END as FORECAST_ARR_FOR_SORTING "
	//	} else {
	//		selectQuery += ", CASE WHEN o.renewalForecastArr <> \"\" and o.renewalForecastArr IS NOT NULL THEN o.renewalForecastArr ELSE 0 END as FORECAST_ARR_FOR_SORTING "
	//	}
	//	aliases += ", FORECAST_ARR_FOR_SORTING "
	//}
	//if sort != nil && sort.By == SearchSortParamOrganization {
	//	selectQuery += " OPTIONAL MATCH (o)-[:SUBSIDIARY_OF]->(parent:Organization) WITH "
	//	selectQuery += aliases + ", parent "
	//}

	cypherSort := utils.CypherSort{}
	//if sort != nil {
	//	if sort.By == SearchSortParamName {
	//		if sort.CaseSensitive != nil && *sort.CaseSensitive {
	//			selectQuery += " ORDER BY NAME_FOR_SORTING " + string(sort.Direction)
	//		} else {
	//			selectQuery += " ORDER BY toLower(NAME_FOR_SORTING) " + string(sort.Direction)
	//		}
	//	} else if sort.By == SearchSortParamRelationship {
	//		selectQuery += " ORDER BY RELATIONSHIP_FOR_SORTING " + string(sort.Direction)
	//	} else if sort.By == SearchSortParamStage {
	//		selectQuery += " ORDER BY STAGE_FOR_SORTING " + string(sort.Direction)
	//	} else if sort.By == SearchSortParamIndustry {
	//		selectQuery += " ORDER BY INDUSTRY_FOR_SORTING " + string(sort.Direction)
	//	} else if sort.By == SearchSortParamOrganization {
	//		cypherSort.NewSortRule("NAME", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithCoalesce().WithAlias("parent")
	//		cypherSort.NewSortRule("NAME", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithCoalesce()
	//		cypherSort.NewSortRule("NAME", string(sort.Direction), true, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithAlias("parent").WithDescending()
	//		cypherSort.NewSortRule("NAME", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
	//		selectQuery += string(cypherSort.SortingCypherFragment("o"))
	//	} else if sort.By == SearchSortParamForecastArr {
	//		selectQuery += " ORDER BY FORECAST_ARR_FOR_SORTING " + string(sort.Direction)
	//	} else if sort.By == SearchSortParamRenewalLikelihood {
	//		selectQuery += " ORDER BY RENEWAL_LIKELIHOOD_FOR_SORTING " + string(sort.Direction)
	//	} else if sort.By == SearchSortParamOnboardingStatus {
	//		selectQuery += " ORDER BY ONBOARDING_STATUS_FOR_SORTING " + string(sort.Direction) +
	//			", ONBOARDING_UPDATED_AT_FOR_SORTING " + string(sort.Direction)
	//	} else if sort.By == SearchSortParamRenewalCycleNext {
	//		selectQuery += " ORDER BY RENEWAL_CYCLE_NEXT_FOR_SORTING " + string(sort.Direction)
	//	} else if sort.By == SearchSortParamRenewalDate {
	//		selectQuery += " ORDER BY RENEWAL_DATE_FOR_SORTING " + string(sort.Direction)
	//	} else if sort.By == SearchSortParamChurnDate {
	//		selectQuery += " ORDER BY CHURN_DATE_FOR_SORTING " + string(sort.Direction)
	//	} else if sort.By == "DOMAIN" {
	//		cypherSort.NewSortRule("DOMAIN", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.DomainEntity{}))
	//		selectQuery += string(cypherSort.SortingCypherFragment("d"))
	//	} else if sort.By == SearchSortParamLocation {
	//		cypherSort.NewSortRule("COUNTRY", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.LocationEntity{}))
	//		cypherSort.NewSortRule("REGION", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.LocationEntity{}))
	//		cypherSort.NewSortRule("LOCALITY", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.LocationEntity{}))
	//		selectQuery += string(cypherSort.SortingCypherFragment("l"))
	//	} else if sort.By == "OWNER" {
	//		if sort.CaseSensitive != nil && *sort.CaseSensitive {
	//			selectQuery += " ORDER BY OWNER_FIRST_NAME_FOR_SORTING " + string(sort.Direction) + ", OWNER_LAST_NAME_FOR_SORTING " + string(sort.Direction)
	//		} else {
	//			selectQuery += " ORDER BY toLower(OWNER_FIRST_NAME_FOR_SORTING) " + string(sort.Direction) + ", toLower(OWNER_LAST_NAME_FOR_SORTING) " + string(sort.Direction)
	//		}
	//	} else if sort.By == SearchSortParamLastTouchpointAt || sort.By == SearchSortParamLastTouchpoint {
	//		cypherSort.NewSortRule("LAST_TOUCHPOINT_AT", string(sort.Direction), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
	//		selectQuery += string(cypherSort.SortingCypherFragment("o"))
	//	} else if sort.By == SearchSortParamLastTouchpointType {
	//		cypherSort.NewSortRule("LAST_TOUCHPOINT_TYPE", string(sort.Direction), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
	//		selectQuery += string(cypherSort.SortingCypherFragment("o"))
	//	} else if sort.By == SearchSortParamUpdatedAt {
	//		cypherSort.NewSortRule("UPDATED_AT", string(sort.Direction), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
	//		selectQuery += string(cypherSort.SortingCypherFragment("o"))
	//	}
	//} else
	//{
	cypherSort.NewSortRule("UPDATED_AT", string(commonmodel.SortingDirectionDesc), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
	selectQuery += string(cypherSort.SortingCypherFragment("o"))
	//}

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
	if filter.Filter.Operation == commonmodel.ComparisonOperatorIsEmpty && filter.Filter.Value.Str != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, nil, commonmodel.ComparisonOperatorIsEmpty))
	} else if filter.Filter.Operation == commonmodel.ComparisonOperatorIn && filter.Filter.Value.ArrayStr != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, filter.Filter.Value.ArrayStr, commonmodel.ComparisonOperatorIn))
	}
}

func createBetweenOrEmptyTimeFilter(filter *model.Filter, cypherFilter *utils.CypherFilter, neo4jProperty string) {
	if filter.Filter.Operation == commonmodel.ComparisonOperatorBetween && filter.Filter.Value.ArrayTime != nil && len(*filter.Filter.Value.ArrayTime) == 2 {
		times := *filter.Filter.Value.ArrayTime
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, times[0], commonmodel.ComparisonOperatorGte))
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, times[1], commonmodel.ComparisonOperatorLte))
	} else if filter.Filter.Operation == commonmodel.ComparisonOperatorIsEmpty {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, nil, commonmodel.ComparisonOperatorIsEmpty))
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
