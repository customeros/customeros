package repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	enummapper "github.com/customeros/customeros/packages/server/customer-os-api/mapper/enum"
)

const (
	SearchSortParamOrganization           = "ORGANIZATION"
	SearchSortParamWebsite                = "WEBSITE"
	SearchSortParamEmail                  = "EMAIL"
	SearchSortParamCountry                = "COUNTRY"
	SearchSortParamOnboardingStatus       = "ONBOARDING_STATUS"
	SearchSortParamIsCustomer             = "IS_CUSTOMER"
	SearchSortParamStage                  = "STAGE"
	SearchSortParamIndustry               = "INDUSTRY"
	SearchSortParamName                   = "NAME"
	SearchSortParamRenewalLikelihood      = "RENEWAL_LIKELIHOOD"
	SearchSortParamRenewalCycleNext       = "RENEWAL_CYCLE_NEXT"
	SearchSortParamRenewalDate            = "RENEWAL_DATE"
	SearchSortParamChurnDate              = "CHURN_DATE"
	SearchSortParamForecastArr            = "FORECAST_ARR"
	SearchSortParamRegion                 = "REGION"
	SearchSortParamLocality               = "LOCALITY"
	SearchSortParamOwnerId                = "OWNER_ID"
	SearchSortParamLocation               = "LOCATION"
	SearchSortParamOwner                  = "OWNER"
	SearchSortParamLastTouchpoint         = "LAST_TOUCHPOINT"
	SearchSortParamLastTouchpointAt       = "LAST_TOUCHPOINT_AT"
	SearchSortParamLastTouchpointType     = "LAST_TOUCHPOINT_TYPE"
	SearchSortParamRenewalCycle           = "RENEWAL_CYCLE"
	SearchSortParamContractLengthInMonths = "CONTRACT_LENGTH_IN_MONTHS"
	SearchParamExternalId                 = "EXTERNAL_ID"
	SearchSortParamUpdatedAt              = "UPDATED_AT"
)

type DashboardRepository interface {
	// deprecated
	GetDashboardViewOrganizationData(ctx context.Context, tenant string, skip, limit int, where *model.Filter, sort *commonmodel.SortBy) (*utils.DbNodesWithTotalCount, error)
	GetDashboardViewRenewalData(ctx context.Context, tenant string, skip, limit int, where *model.Filter, sort *commonmodel.SortBy) (*utils.RecordsWithTotalCount, error)
}

type dashboardRepository struct {
	driver *neo4j.DriverWithContext
}

func NewDashboardRepository(driver *neo4j.DriverWithContext) DashboardRepository {
	return &dashboardRepository{
		driver: driver,
	}
}

func createStringCypherFilterWithValueOrEmpty(filter *model.FilterItem, propertyName string) *utils.CypherFilter {
	if filter.IncludeEmpty != nil && *filter.IncludeEmpty {
		orFilter := utils.CypherFilter{}
		orFilter.LogicalOperator = utils.OR
		orFilter.Details = new(utils.CypherFilterItem)

		orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter(propertyName, *filter.Value.Str, commonmodel.ComparisonOperatorContains))
		orFilter.Filters = append(orFilter.Filters, utils.CreateCypherFilterEq(propertyName, ""))
		orFilter.Filters = append(orFilter.Filters, utils.CreateCypherFilterIsNull(propertyName))
		return &orFilter
	} else {
		return utils.CreateStringCypherFilter(propertyName, *filter.Value.Str, commonmodel.ComparisonOperatorContains)
	}
}

func (r *dashboardRepository) GetDashboardViewOrganizationData(ctx context.Context, tenant string, skip, limit int, where *model.Filter, sort *commonmodel.SortBy) (*utils.DbNodesWithTotalCount, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "DashboardRepository.GetDashboardViewOrganizationData")
	defer spans.Finish()
	spans.LogKV(log.Int("skip", skip), log.Int("limit", limit))
	spans.LogObjectAsJson("where", where)
	spans.LogObjectAsJson("sort", sort)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	organizationFilterCypher, organizationFilterParams := "", make(map[string]interface{})
	emailFilterCypher, emailFilterParams := "", make(map[string]interface{})
	locationFilterCypher, locationFilterParams := "", make(map[string]interface{})

	ownerId := []string{}
	ownerIncludeEmpty := false
	externalId := ""

	// ORGANIZATION, EMAIL, COUNTRY, REGION, LOCALITY
	// region organization filters
	if where != nil {
		organizationFilter := new(utils.CypherFilter)
		organizationFilter.Negate = false
		organizationFilter.LogicalOperator = utils.AND
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
			if filter.Filter.Property == SearchSortParamOrganization {
				orFilter := utils.CypherFilter{}
				orFilter.LogicalOperator = utils.OR
				orFilter.Details = new(utils.CypherFilterItem)

				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("name", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorContains))
				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("website", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorContains))
				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("customerOsId", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorContains))
				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("referenceId", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorContains))

				organizationFilter.Filters = append(organizationFilter.Filters, &orFilter)
			} else if filter.Filter.Property == SearchSortParamName {
				organizationFilter.Filters = append(organizationFilter.Filters, createStringCypherFilterWithValueOrEmpty(filter.Filter, "name"))
			} else if filter.Filter.Property == SearchSortParamWebsite {
				organizationFilter.Filters = append(organizationFilter.Filters, createStringCypherFilterWithValueOrEmpty(filter.Filter, "website"))
			} else if filter.Filter.Property == SearchSortParamStage && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("stage", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamIndustry && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("industry", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamEmail {
				emailFilter.Filters = append(emailFilter.Filters, utils.CreateStringCypherFilter("email", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorContains))
				emailFilter.Filters = append(emailFilter.Filters, utils.CreateStringCypherFilter("rawEmail", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorContains))
			} else if filter.Filter.Property == SearchSortParamCountry {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("country", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorEq))
			} else if filter.Filter.Property == SearchSortParamRegion {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("region", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorEq))
			} else if filter.Filter.Property == SearchSortParamLocality {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("locality", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorEq))
			} else if filter.Filter.Property == SearchSortParamOwnerId {
				if filter.Filter.Value.ArrayStr != nil {
					ownerId = *filter.Filter.Value.ArrayStr
				}
				ownerIncludeEmpty = *filter.Filter.IncludeEmpty
			} else if filter.Filter.Property == SearchParamExternalId {
				externalId = *filter.Filter.Value.Str
			} else if filter.Filter.Property == SearchSortParamIsCustomer && filter.Filter.Value.ArrayBool != nil && len(*filter.Filter.Value.ArrayBool) >= 1 {
				if (*filter.Filter.Value.ArrayBool)[0] {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterEq("stage", enum.Customer.String()))
				} else {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterNotEq("stage", enum.Customer.String()))
				}
			} else if filter.Filter.Property == SearchSortParamRenewalLikelihood && filter.Filter.Value.ArrayStr != nil && len(*filter.Filter.Value.ArrayStr) >= 1 {
				renewalLikelihoodValues := make([]string, 0)
				for _, v := range *filter.Filter.Value.ArrayStr {
					renewalLikelihoodValues = append(renewalLikelihoodValues, enummapper.MapOpportunityRenewalLikelihoodFromString(&v))
				}
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("derivedRenewalLikelihood", renewalLikelihoodValues))
			} else if filter.Filter.Property == SearchSortParamOnboardingStatus && filter.Filter.Value.ArrayStr != nil && len(*filter.Filter.Value.ArrayStr) >= 1 {
				onboardingStatusValues := make([]string, 0)
				for _, v := range *filter.Filter.Value.ArrayStr {
					onboardingStatusValues = append(onboardingStatusValues, enummapper.MapOnboardingStatusFromString(&v))
				}
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("onboardingStatus", onboardingStatusValues))
			} else if filter.Filter.Property == SearchSortParamRenewalCycleNext && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("billingDetailsRenewalCycleNext", *filter.Filter.Value.Time, commonmodel.ComparisonOperatorLte))
			} else if filter.Filter.Property == SearchSortParamRenewalDate && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("derivedNextRenewalAt", *filter.Filter.Value.Time, commonmodel.ComparisonOperatorLte))
			} else if filter.Filter.Property == SearchSortParamForecastArr && filter.Filter.Value.ArrayInt != nil && len(*filter.Filter.Value.ArrayInt) == 2 {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("renewalForecastArr", (*filter.Filter.Value.ArrayInt)[0], commonmodel.ComparisonOperatorGte))
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("renewalForecastArr", (*filter.Filter.Value.ArrayInt)[1], commonmodel.ComparisonOperatorLte))
			} else if (filter.Filter.Property == SearchSortParamLastTouchpointAt || filter.Filter.Property == SearchSortParamLastTouchpoint) && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("lastTouchpointAt", *filter.Filter.Value.Time, commonmodel.ComparisonOperatorGte))
			} else if filter.Filter.Property == SearchSortParamLastTouchpointType && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("lastTouchpointType", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamUpdatedAt && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("updatedAt", *filter.Filter.Value.Time, commonmodel.ComparisonOperatorGte))
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

	// endregion
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	params := map[string]any{
		"tenant":     tenant,
		"ownerId":    ownerId,
		"externalId": externalId,
		"skip":       skip,
		"limit":      limit,
	}

	utils.MergeMapToMap(organizationFilterParams, params)
	utils.MergeMapToMap(emailFilterParams, params)
	utils.MergeMapToMap(locationFilterParams, params)

	// region count query
	countQuery := `MATCH (o:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	if len(ownerId) > 0 || ownerIncludeEmpty {
		countQuery += ` OPTIONAL MATCH (o)<-[:OWNS]-(owner:User) WITH *`
	}
	if emailFilterCypher != "" {
		countQuery += ` MATCH (o)-[:HAS]->(e:Email) WITH *`
	}
	if locationFilterCypher != "" {
		countQuery += ` MATCH (o)-[:ASSOCIATED_WITH]->(l:Location) WITH *`
	}
	if externalId != "" {
		countQuery += ` MATCH (o)-[:IS_LINKED_WITH {externalId:$externalId}]->(ext:ExternalSystem) WITH *`
	}
	countQuery += ` WHERE o.hide = false `

	if organizationFilterCypher != "" || emailFilterCypher != "" || locationFilterCypher != "" || len(ownerId) > 0 || ownerIncludeEmpty {
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
	if emailFilterCypher != "" {
		countQueryParts = append(countQueryParts, emailFilterCypher)
	}
	if locationFilterCypher != "" {
		countQueryParts = append(countQueryParts, locationFilterCypher)
	}

	countQuery = countQuery + strings.Join(countQueryParts, " AND ") + fmt.Sprintf(` RETURN count(distinct(o))`)

	spans.LogKV("countQuery", countQuery)

	countRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		countQueryResult, err := tx.Run(ctx, countQuery, params)
		if err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsType[int64](ctx, countQueryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	dbNodesWithTotalCount.Count = countRecord.(int64)
	// end count region

	// region query to fetch data
	query := `MATCH (o:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	if len(ownerId) > 0 || ownerIncludeEmpty {
		query += fmt.Sprintf(` OPTIONAL MATCH (o)<-[:OWNS]-(owner:User) WITH *`)
	}
	if externalId != "" {
		query += ` MATCH (o)-[:IS_LINKED_WITH {externalId:$externalId}]->(ext:ExternalSystem) WITH *`
	}
	query += fmt.Sprintf(` OPTIONAL MATCH (o)-[:HAS_DOMAIN]->(d:Domain) WITH *`)
	query += fmt.Sprintf(` OPTIONAL MATCH (o)-[:HAS]->(e:Email_%s) WITH *`, tenant)
	query += fmt.Sprintf(` OPTIONAL MATCH (o)-[:ASSOCIATED_WITH]->(l:Location_%s) WITH *`, tenant)
	if sort != nil && sort.By == SearchSortParamOwner {
		query += fmt.Sprintf(` OPTIONAL MATCH (o)<-[:OWNS]-(owner:User_%s) WITH *`, tenant)
	}
	query += ` WHERE (o.hide = false) `

	if organizationFilterCypher != "" || emailFilterCypher != "" || locationFilterCypher != "" || len(ownerId) > 0 || ownerIncludeEmpty {
		query += " AND "
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
	if emailFilterCypher != "" {
		queryParts = append(queryParts, emailFilterCypher)
	}
	if locationFilterCypher != "" {
		queryParts = append(queryParts, locationFilterCypher)
	}

	// endregion
	query = query + strings.Join(queryParts, " AND ")

	// sort region
	aliases := " o, d, l"
	query += " WITH o, d, l "
	if sort != nil && sort.By == SearchSortParamOwner {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			query += ", CASE WHEN owner.firstName <> \"\" and not owner.firstName is null THEN owner.firstName ELSE 'ZZZZZZZZZZZZZZZZZZZ' END as OWNER_FIRST_NAME_FOR_SORTING "
			query += ", CASE WHEN owner.lastName <> \"\" and not owner.lastName is null THEN owner.lastName ELSE 'ZZZZZZZZZZZZZZZZZZZ' END as OWNER_LAST_NAME_FOR_SORTING "
		} else {
			query += ", CASE WHEN owner.firstName <> \"\" and not owner.firstName is null THEN owner.firstName ELSE 'AAAAAAAAAAAAAAAAAAA' END as OWNER_FIRST_NAME_FOR_SORTING "
			query += ", CASE WHEN owner.lastName <> \"\" and not owner.lastName is null THEN owner.lastName ELSE 'AAAAAAAAAAAAAAAAAAA' END as OWNER_LAST_NAME_FOR_SORTING "
		}
		aliases += ", OWNER_FIRST_NAME_FOR_SORTING, OWNER_LAST_NAME_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamName {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			query += ", CASE WHEN o.name <> \"\" and not o.name is null THEN o.name ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as NAME_FOR_SORTING "
		} else {
			query += ", o.name as NAME_FOR_SORTING "
		}
		aliases += ", NAME_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamRenewalLikelihood {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			query += ", CASE WHEN o.derivedRenewalLikelihoodOrder IS NOT NULL THEN o.derivedRenewalLikelihoodOrder ELSE 9999 END as RENEWAL_LIKELIHOOD_FOR_SORTING "
		} else {
			query += ", CASE WHEN o.derivedRenewalLikelihoodOrder IS NOT NULL THEN o.derivedRenewalLikelihoodOrder ELSE -1 END as RENEWAL_LIKELIHOOD_FOR_SORTING "
		}
		aliases += ", RENEWAL_LIKELIHOOD_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamStage {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			query += ", CASE WHEN o.stage <> '' AND NOT o.stage IS NULL THEN o.stage ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as STAGE_FOR_SORTING "
		} else {
			query += ", o.stage as STAGE_FOR_SORTING "
		}
		aliases += ", STAGE_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamIndustry {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			query += ", CASE WHEN o.industry <> '' AND NOT o.industry IS NULL THEN o.industry ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as STAGE_FOR_SORTING "
		} else {
			query += ", o.industry as INDUSTRY_FOR_SORTING "
		}
		aliases += ", INDUSTRY_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamRenewalCycleNext {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			query += ", CASE WHEN o.billingDetailsRenewalCycleNext IS NOT NULL THEN date(o.billingDetailsRenewalCycleNext) ELSE date('2100-01-01') END as RENEWAL_CYCLE_NEXT_FOR_SORTING "
		} else {
			query += ", CASE WHEN o.billingDetailsRenewalCycleNext IS NOT NULL THEN date(o.billingDetailsRenewalCycleNext) ELSE date('1900-01-01') END as RENEWAL_CYCLE_NEXT_FOR_SORTING "
		}
		aliases += ", RENEWAL_CYCLE_NEXT_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamRenewalDate {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			query += ", CASE WHEN o.derivedNextRenewalAt IS NOT NULL THEN date(o.derivedNextRenewalAt) ELSE date('2100-01-01') END as RENEWAL_DATE_FOR_SORTING "
		} else {
			query += ", CASE WHEN o.derivedNextRenewalAt IS NOT NULL THEN date(o.derivedNextRenewalAt) ELSE date('1900-01-01') END as RENEWAL_DATE_FOR_SORTING "
		}
		aliases += ", RENEWAL_DATE_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamChurnDate {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			query += ", CASE WHEN o.derivedChurnedAt IS NOT NULL THEN date(o.derivedChurnedAt) ELSE date('2100-01-01') END as CHURN_DATE_FOR_SORTING "
		} else {
			query += ", CASE WHEN o.derivedChurnedAt IS NOT NULL THEN date(o.derivedChurnedAt) ELSE date('1900-01-01') END as CHURN_DATE_FOR_SORTING "
		}
		aliases += ", RENEWAL_DATE_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamOnboardingStatus {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			query += ", CASE WHEN o.onboardingStatusOrder IS NOT NULL THEN o.onboardingStatusOrder ELSE 9999 END as ONBOARDING_STATUS_FOR_SORTING "
			query += ", o.onboardingUpdatedAt AS ONBOARDING_UPDATED_AT_FOR_SORTING "
		} else {
			query += ", CASE WHEN o.onboardingStatusOrder IS NOT NULL THEN o.onboardingStatusOrder ELSE -1 END as ONBOARDING_STATUS_FOR_SORTING "
			query += ", o.onboardingUpdatedAt AS ONBOARDING_UPDATED_AT_FOR_SORTING "
		}
		aliases += ", ONBOARDING_STATUS_FOR_SORTING, ONBOARDING_UPDATED_AT_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamForecastArr {
		if sort.Direction == commonmodel.SortingDirectionAsc {
			query += ", CASE WHEN o.renewalForecastArr <> \"\" and o.renewalForecastArr IS NOT NULL THEN o.renewalForecastArr ELSE 9999999999999999 END as FORECAST_ARR_FOR_SORTING "
		} else {
			query += ", CASE WHEN o.renewalForecastArr <> \"\" and o.renewalForecastArr IS NOT NULL THEN o.renewalForecastArr ELSE 0 END as FORECAST_ARR_FOR_SORTING "
		}
		aliases += ", FORECAST_ARR_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamOrganization {
		query += " OPTIONAL MATCH (o)-[:SUBSIDIARY_OF]->(parent:Organization) WITH "
		query += aliases + ", parent "
	}

	cypherSort := utils.CypherSort{}
	if sort != nil {
		if sort.By == SearchSortParamName {
			if sort.CaseSensitive != nil && *sort.CaseSensitive {
				query += " ORDER BY NAME_FOR_SORTING " + string(sort.Direction)
			} else {
				query += " ORDER BY toLower(NAME_FOR_SORTING) " + string(sort.Direction)
			}
		} else if sort.By == SearchSortParamStage {
			query += " ORDER BY STAGE_FOR_SORTING " + string(sort.Direction)
		} else if sort.By == SearchSortParamIndustry {
			query += " ORDER BY INDUSTRY_FOR_SORTING " + string(sort.Direction)
		} else if sort.By == SearchSortParamOrganization {
			cypherSort.NewSortRule("NAME", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithCoalesce().WithAlias("parent")
			cypherSort.NewSortRule("NAME", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithCoalesce()
			cypherSort.NewSortRule("NAME", string(sort.Direction), true, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithAlias("parent").WithDescending()
			cypherSort.NewSortRule("NAME", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
			query += string(cypherSort.SortingCypherFragment("o"))
		} else if sort.By == SearchSortParamForecastArr {
			query += " ORDER BY FORECAST_ARR_FOR_SORTING " + string(sort.Direction)
		} else if sort.By == SearchSortParamRenewalLikelihood {
			query += " ORDER BY RENEWAL_LIKELIHOOD_FOR_SORTING " + string(sort.Direction)
		} else if sort.By == SearchSortParamOnboardingStatus {
			query += " ORDER BY ONBOARDING_STATUS_FOR_SORTING " + string(sort.Direction) +
				", ONBOARDING_UPDATED_AT_FOR_SORTING " + string(sort.Direction)
		} else if sort.By == SearchSortParamRenewalCycleNext {
			query += " ORDER BY RENEWAL_CYCLE_NEXT_FOR_SORTING " + string(sort.Direction)
		} else if sort.By == SearchSortParamRenewalDate {
			query += " ORDER BY RENEWAL_DATE_FOR_SORTING " + string(sort.Direction)
		} else if sort.By == SearchSortParamChurnDate {
			query += " ORDER BY CHURN_DATE_FOR_SORTING " + string(sort.Direction)
		} else if sort.By == "DOMAIN" {
			cypherSort.NewSortRule("DOMAIN", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.DomainEntity{}))
			query += string(cypherSort.SortingCypherFragment("d"))
		} else if sort.By == SearchSortParamLocation {
			cypherSort.NewSortRule("COUNTRY", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.LocationEntity{}))
			cypherSort.NewSortRule("REGION", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.LocationEntity{}))
			cypherSort.NewSortRule("LOCALITY", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.LocationEntity{}))
			query += string(cypherSort.SortingCypherFragment("l"))
		} else if sort.By == "OWNER" {
			if sort.CaseSensitive != nil && *sort.CaseSensitive {
				query += " ORDER BY OWNER_FIRST_NAME_FOR_SORTING " + string(sort.Direction) + ", OWNER_LAST_NAME_FOR_SORTING " + string(sort.Direction)
			} else {
				query += " ORDER BY toLower(OWNER_FIRST_NAME_FOR_SORTING) " + string(sort.Direction) + ", toLower(OWNER_LAST_NAME_FOR_SORTING) " + string(sort.Direction)
			}
		} else if sort.By == SearchSortParamLastTouchpointAt || sort.By == SearchSortParamLastTouchpoint {
			cypherSort.NewSortRule("LAST_TOUCHPOINT_AT", string(sort.Direction), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
			query += string(cypherSort.SortingCypherFragment("o"))
		} else if sort.By == SearchSortParamLastTouchpointType {
			cypherSort.NewSortRule("LAST_TOUCHPOINT_TYPE", string(sort.Direction), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
			query += string(cypherSort.SortingCypherFragment("o"))
		} else if sort.By == SearchSortParamUpdatedAt {
			cypherSort.NewSortRule("UPDATED_AT", string(sort.Direction), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
			query += string(cypherSort.SortingCypherFragment("o"))
		}
	} else {
		cypherSort.NewSortRule("UPDATED_AT", string(commonmodel.SortingDirectionDesc), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
		query += string(cypherSort.SortingCypherFragment("o"))
	}
	// end sort region
	query += fmt.Sprintf(` RETURN distinct(o) `)
	query += fmt.Sprintf(` SKIP $skip LIMIT $limit`)

	spans.LogKV(log.String("query", query))
	spans.LogObjectAsJson("params", params)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
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

func (r *dashboardRepository) GetDashboardViewRenewalData(ctx context.Context, tenant string, skip, limit int, where *model.Filter, sort *commonmodel.SortBy) (*utils.RecordsWithTotalCount, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "DashboardRepository.GetDashboardViewRenewalData")
	defer spans.Finish()

	spans.LogKV(log.Int("skip", skip), log.Int("limit", limit))
	spans.LogObjectAsJson("where", where)
	spans.LogObjectAsJson("sort", sort)

	dbRecordsWithTotalCount := new(utils.RecordsWithTotalCount)

	organizationFilterCypher, organizationFilterParams := "", make(map[string]interface{})
	contractFilterCypher, contractFilterParams := "", make(map[string]interface{})
	opportunityFilterCypher, opportunityFilterParams := "", make(map[string]interface{})
	emailFilterCypher, emailFilterParams := "", make(map[string]interface{})
	locationFilterCypher, locationFilterParams := "", make(map[string]interface{})

	ownerId := []string{}
	ownerIncludeEmpty := false

	// ORGANIZATION, EMAIL, COUNTRY, REGION, LOCALITY
	// region organization & contract filters
	if where != nil {
		organizationFilter := new(utils.CypherFilter)
		organizationFilter.Negate = false
		organizationFilter.LogicalOperator = utils.AND
		organizationFilter.Filters = make([]*utils.CypherFilter, 0)

		emailFilter := new(utils.CypherFilter)
		emailFilter.Negate = false
		emailFilter.LogicalOperator = utils.OR
		emailFilter.Filters = make([]*utils.CypherFilter, 0)

		locationFilter := new(utils.CypherFilter)
		locationFilter.Negate = false
		locationFilter.LogicalOperator = utils.OR
		locationFilter.Filters = make([]*utils.CypherFilter, 0)

		contractFilter := new(utils.CypherFilter)
		contractFilter.Negate = false
		contractFilter.LogicalOperator = utils.AND
		contractFilter.Filters = make([]*utils.CypherFilter, 0)

		opportunityFilter := new(utils.CypherFilter)
		opportunityFilter.Negate = false
		opportunityFilter.LogicalOperator = utils.AND
		opportunityFilter.Filters = make([]*utils.CypherFilter, 0)

		for _, filter := range where.And {
			if filter.Filter.Property == SearchSortParamOrganization {
				orFilter := utils.CypherFilter{}
				orFilter.LogicalOperator = utils.OR
				orFilter.Details = new(utils.CypherFilterItem)

				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("name", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorContains))
				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("website", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorContains))
				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("customerOsId", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorContains))
				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("referenceId", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorContains))

				organizationFilter.Filters = append(organizationFilter.Filters, &orFilter)
			} else if filter.Filter.Property == SearchSortParamName {
				organizationFilter.Filters = append(organizationFilter.Filters, createStringCypherFilterWithValueOrEmpty(filter.Filter, "name"))
			} else if filter.Filter.Property == SearchSortParamWebsite {
				organizationFilter.Filters = append(organizationFilter.Filters, createStringCypherFilterWithValueOrEmpty(filter.Filter, "website"))
			} else if filter.Filter.Property == SearchSortParamStage && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("stage", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamIndustry && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("industry", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamEmail {
				emailFilter.Filters = append(emailFilter.Filters, utils.CreateStringCypherFilter("email", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorContains))
				emailFilter.Filters = append(emailFilter.Filters, utils.CreateStringCypherFilter("rawEmail", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorContains))
			} else if filter.Filter.Property == SearchSortParamCountry {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("country", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorEq))
			} else if filter.Filter.Property == SearchSortParamRegion {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("region", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorEq))
			} else if filter.Filter.Property == SearchSortParamLocality {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("locality", *filter.Filter.Value.Str, commonmodel.ComparisonOperatorEq))
			} else if filter.Filter.Property == SearchSortParamOwnerId {
				if filter.Filter.Value.ArrayStr != nil {
					ownerId = *filter.Filter.Value.ArrayStr
				}
				ownerIncludeEmpty = *filter.Filter.IncludeEmpty
			} else if filter.Filter.Property == SearchSortParamIsCustomer && filter.Filter.Value.ArrayBool != nil && len(*filter.Filter.Value.ArrayBool) >= 1 {
				if (*filter.Filter.Value.ArrayBool)[0] {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterEq("stage", enum.Customer.String()))
				} else {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterNotEq("stage", enum.Customer.String()))
				}
			} else if filter.Filter.Property == SearchSortParamRenewalLikelihood && filter.Filter.Value.ArrayStr != nil && len(*filter.Filter.Value.ArrayStr) >= 1 {
				renewalLikelihoodValues := make([]string, 0)
				for _, v := range *filter.Filter.Value.ArrayStr {
					renewalLikelihoodValues = append(renewalLikelihoodValues, enummapper.MapOpportunityRenewalLikelihoodFromString(&v))
				}
				opportunityFilter.Filters = append(opportunityFilter.Filters, utils.CreateCypherFilterIn("renewalLikelihood", renewalLikelihoodValues))
			} else if filter.Filter.Property == SearchSortParamOnboardingStatus && filter.Filter.Value.ArrayStr != nil && len(*filter.Filter.Value.ArrayStr) >= 1 {
				onboardingStatusValues := make([]string, 0)
				for _, v := range *filter.Filter.Value.ArrayStr {
					onboardingStatusValues = append(onboardingStatusValues, enummapper.MapOnboardingStatusFromString(&v))
				}
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("onboardingStatus", onboardingStatusValues))
			} else if filter.Filter.Property == SearchSortParamRenewalCycleNext && filter.Filter.Value.Time != nil {
				opportunityFilter.Filters = append(opportunityFilter.Filters, utils.CreateCypherFilter("renewedAt", *filter.Filter.Value.Time, commonmodel.ComparisonOperatorLte))
			} else if filter.Filter.Property == SearchSortParamRenewalDate && filter.Filter.Value.Time != nil {
				opportunityFilter.Filters = append(opportunityFilter.Filters, utils.CreateCypherFilter("renewedAt", *filter.Filter.Value.Time, commonmodel.ComparisonOperatorLte))
			} else if filter.Filter.Property == SearchSortParamForecastArr && filter.Filter.Value.ArrayInt != nil && len(*filter.Filter.Value.ArrayInt) == 2 {
				opportunityFilter.Filters = append(opportunityFilter.Filters, utils.CreateCypherFilter("maxAmount", (*filter.Filter.Value.ArrayInt)[0], commonmodel.ComparisonOperatorGte))
				opportunityFilter.Filters = append(opportunityFilter.Filters, utils.CreateCypherFilter("maxAmount", (*filter.Filter.Value.ArrayInt)[1], commonmodel.ComparisonOperatorLte))
			} else if (filter.Filter.Property == SearchSortParamLastTouchpointAt || filter.Filter.Property == SearchSortParamLastTouchpoint) && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("lastTouchpointAt", *filter.Filter.Value.Time, commonmodel.ComparisonOperatorGte))
			} else if filter.Filter.Property == SearchSortParamLastTouchpointType && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("lastTouchpointType", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamUpdatedAt && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("updatedAt", *filter.Filter.Value.Time, commonmodel.ComparisonOperatorGte))
			} else if filter.Filter.Property == SearchSortParamRenewalCycle {
				if filter.Filter.Value.Str != nil {
					switch *filter.Filter.Value.Str {
					case "MONTHLY":
						contractFilter.Filters = append(contractFilter.Filters, utils.CreateCypherFilter("lengthInMonths", 1, commonmodel.ComparisonOperatorEq))
					case "QUARTERLY":
						contractFilter.Filters = append(contractFilter.Filters, utils.CreateCypherFilter("lengthInMonths", 3, commonmodel.ComparisonOperatorEq))
					case "ANNUALLY":
						contractFilter.Filters = append(contractFilter.Filters, utils.CreateCypherFilter("lengthInMonths", 12, commonmodel.ComparisonOperatorGte))
					}
				}
			} else if filter.Filter.Property == SearchSortParamContractLengthInMonths {
				if filter.Filter.Value.Int != nil {
					contractFilter.Filters = append(contractFilter.Filters, utils.CreateCypherFilter("lengthInMonths", *filter.Filter.Value.Int, commonmodel.ComparisonOperatorEq))
				}
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
		if len(contractFilter.Filters) > 0 {
			contractFilterCypher, contractFilterParams = contractFilter.BuildCypherFilterFragmentWithParamName("contract", "contract_param_")
		}
		if len(opportunityFilter.Filters) > 0 {
			opportunityFilterCypher, opportunityFilterParams = opportunityFilter.BuildCypherFilterFragmentWithParamName("op", "opportunity_param_")
		}
	}

	// endregion
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
		utils.MergeMapToMap(contractFilterParams, params)
		utils.MergeMapToMap(opportunityFilterParams, params)

		// region count query
		countQuery := `MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
					 MATCH (o)-[:HAS_CONTRACT]->(contract:Contract)-[:CONTRACT_BELONGS_TO_TENANT]->(t)
					 MATCH (contract)-[:ACTIVE_RENEWAL]->(op:Opportunity) `
		if len(ownerId) > 0 || ownerIncludeEmpty {
			countQuery += ` OPTIONAL MATCH (op)<-[:OWNS]-(owner:User) WITH *`
		}
		if emailFilterCypher != "" {
			countQuery += ` MATCH (o)-[:HAS]->(e:Email) WITH *`
		}
		if locationFilterCypher != "" {
			countQuery += ` MATCH (o)-[:ASSOCIATED_WITH]->(l:Location) WITH *`
		}
		if contractFilterCypher != "" {
			countQuery += ` MATCH (o)-[:HAS_CONTRACT]->(contract:Contract) WITH *`
		}
		if opportunityFilterCypher != "" {
			countQuery += ` MATCH (contract)-[:ACTIVE_RENEWAL]->(op:Opportunity) WITH *`
		}
		countQuery += ` WHERE o.hide = false `

		if organizationFilterCypher != "" || emailFilterCypher != "" || locationFilterCypher != "" || contractFilterCypher != "" || opportunityFilterCypher != "" || len(ownerId) > 0 || ownerIncludeEmpty {
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
		if emailFilterCypher != "" {
			countQueryParts = append(countQueryParts, emailFilterCypher)
		}
		if locationFilterCypher != "" {
			countQueryParts = append(countQueryParts, locationFilterCypher)
		}
		if contractFilterCypher != "" {
			countQueryParts = append(countQueryParts, contractFilterCypher)
		}
		if opportunityFilterCypher != "" {
			countQueryParts = append(countQueryParts, opportunityFilterCypher)
		}

		countQuery = countQuery + strings.Join(countQueryParts, " AND ") + fmt.Sprintf(` RETURN count(distinct(contract))`)

		spans.LogKV(log.String("countQuery", countQuery))

		countQueryResult, err := tx.Run(ctx, countQuery, params)
		if err != nil {
			return nil, err
		}

		countRecord, err := countQueryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		dbRecordsWithTotalCount.Count = countRecord.Values[0].(int64)
		// endregion

		// region query to fetch data
		query := `MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
					 MATCH (o)-[:HAS_CONTRACT]->(contract:Contract)-[:CONTRACT_BELONGS_TO_TENANT]->(t)
					 MATCH (contract)-[:ACTIVE_RENEWAL]->(op:Opportunity)
					 `

		query += fmt.Sprintf(` OPTIONAL MATCH (op)<-[:OWNS]-(owner:User) WITH *`)

		query += ` WHERE (o.hide = false) `

		if organizationFilterCypher != "" || contractFilterCypher != "" || emailFilterCypher != "" || locationFilterCypher != "" || len(ownerId) > 0 || ownerIncludeEmpty {
			query += " AND "
		}

		queryParts := []string{}
		if organizationFilterCypher != "" {
			queryParts = append(queryParts, organizationFilterCypher)
		}
		if len(ownerId) > 0 || ownerIncludeEmpty {
			if len(ownerId) == 0 {
				countQueryParts = append(countQueryParts, fmt.Sprintf(` owner.id IS NULL `))
			} else if ownerIncludeEmpty {
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
		if contractFilterCypher != "" {
			queryParts = append(queryParts, contractFilterCypher)
		}

		// endregion
		query = query + strings.Join(queryParts, " AND ")

		// sort region
		aliases := " o, contract, op, owner"
		query += " WITH o, contract, op, owner "
		if sort != nil && sort.By == SearchSortParamOwner {
			if sort.Direction == commonmodel.SortingDirectionAsc {
				query += ", CASE WHEN owner.firstName <> \"\" and not owner.firstName is null THEN owner.firstName ELSE 'ZZZZZZZZZZZZZZZZZZZ' END as OWNER_FIRST_NAME_FOR_SORTING "
				query += ", CASE WHEN owner.lastName <> \"\" and not owner.lastName is null THEN owner.lastName ELSE 'ZZZZZZZZZZZZZZZZZZZ' END as OWNER_LAST_NAME_FOR_SORTING "
			} else {
				query += ", CASE WHEN owner.firstName <> \"\" and not owner.firstName is null THEN owner.firstName ELSE 'AAAAAAAAAAAAAAAAAAA' END as OWNER_FIRST_NAME_FOR_SORTING "
				query += ", CASE WHEN owner.lastName <> \"\" and not owner.lastName is null THEN owner.lastName ELSE 'AAAAAAAAAAAAAAAAAAA' END as OWNER_LAST_NAME_FOR_SORTING "
			}
			aliases += ", OWNER_FIRST_NAME_FOR_SORTING, OWNER_LAST_NAME_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamName {
			if sort.Direction == commonmodel.SortingDirectionAsc {
				query += ", CASE WHEN o.name <> \"\" and not o.name is null THEN o.name ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as NAME_FOR_SORTING "
			} else {
				query += ", o.name as NAME_FOR_SORTING "
			}
			aliases += ", NAME_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamRenewalLikelihood {
			if sort.Direction == commonmodel.SortingDirectionAsc {
				query += ", CASE WHEN op.renewalLikelihood IS NOT NULL THEN op.renewalLikelihood ELSE 9999 END as RENEWAL_LIKELIHOOD_FOR_SORTING "
			} else {
				query += ", CASE WHEN op.renewalLikelihood IS NOT NULL THEN op.renewalLikelihood ELSE -1 END as RENEWAL_LIKELIHOOD_FOR_SORTING "
			}
			aliases += ", RENEWAL_LIKELIHOOD_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamRenewalCycleNext {
			if sort.Direction == commonmodel.SortingDirectionAsc {
				query += ", CASE WHEN op.renewedAt IS NOT NULL THEN date(op.renewedAt) ELSE date('2100-01-01') END as RENEWAL_CYCLE_NEXT_FOR_SORTING "
			} else {
				query += ", CASE WHEN op.renewedAt IS NOT NULL THEN date(op.renewedAt) ELSE date('1900-01-01') END as RENEWAL_CYCLE_NEXT_FOR_SORTING "
			}
			aliases += ", RENEWAL_CYCLE_NEXT_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamRenewalDate {
			if sort.Direction == commonmodel.SortingDirectionAsc {
				query += ", CASE WHEN op.renewedAt IS NOT NULL THEN date(op.renewedAt) ELSE date('2100-01-01') END as RENEWAL_DATE_FOR_SORTING "
			} else {
				query += ", CASE WHEN op.renewedAt IS NOT NULL THEN date(op.renewedAt) ELSE date('1900-01-01') END as RENEWAL_DATE_FOR_SORTING "
			}
			aliases += ", RENEWAL_DATE_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamContractLengthInMonths {
			if sort.Direction == commonmodel.SortingDirectionAsc {
				query += ", CASE WHEN contract.lengthInMonths IS NOT NULL AND contract.lengthInMonths > 0 THEN contract.lengthInMonths ELSE 9999 END as CONTRACT_LENGTH_FOR_SORTING "
			} else {
				query += ", CASE WHEN contract.lengthInMonths IS NOT NULL THEN contract.lengthInMonths ELSE -1 END as CONTRACT_LENGTH_FOR_SORTING "
			}
			aliases += ", RENEWAL_LIKELIHOOD_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamOnboardingStatus {
			if sort.Direction == commonmodel.SortingDirectionAsc {
				query += ", CASE WHEN o.onboardingStatusOrder IS NOT NULL THEN o.onboardingStatusOrder ELSE 9999 END as ONBOARDING_STATUS_FOR_SORTING "
				query += ", o.onboardingUpdatedAt AS ONBOARDING_UPDATED_AT_FOR_SORTING "
			} else {
				query += ", CASE WHEN o.onboardingStatusOrder IS NOT NULL THEN o.onboardingStatusOrder ELSE -1 END as ONBOARDING_STATUS_FOR_SORTING "
				query += ", o.onboardingUpdatedAt AS ONBOARDING_UPDATED_AT_FOR_SORTING "
			}
			aliases += ", ONBOARDING_STATUS_FOR_SORTING, ONBOARDING_UPDATED_AT_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamForecastArr {
			if sort.Direction == commonmodel.SortingDirectionAsc {
				query += ", CASE WHEN op.maxAmount <> \"\" and op.maxAmount IS NOT NULL THEN op.maxAmount ELSE 9999999999999999 END as FORECAST_ARR_FOR_SORTING "
			} else {
				query += ", CASE WHEN op.maxAmount <> \"\" and op.maxAmount IS NOT NULL THEN op.maxAmount ELSE 0 END as FORECAST_ARR_FOR_SORTING "
			}
			aliases += ", FORECAST_ARR_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamOrganization {
			query += " OPTIONAL MATCH (o)-[:SUBSIDIARY_OF]->(parent:Organization) WITH "
			query += aliases + ", parent "
		}

		cypherSort := utils.CypherSort{}
		if sort != nil {
			if sort.By == SearchSortParamName {
				query += " ORDER BY NAME_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == SearchSortParamOrganization {
				cypherSort.NewSortRule("NAME", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithCoalesce().WithAlias("parent")
				cypherSort.NewSortRule("NAME", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithCoalesce()
				cypherSort.NewSortRule("NAME", string(sort.Direction), true, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithAlias("parent").WithDescending()
				cypherSort.NewSortRule("NAME", string(sort.Direction), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			} else if sort.By == SearchSortParamForecastArr {
				query += " ORDER BY FORECAST_ARR_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == SearchSortParamRenewalLikelihood {
				query += " ORDER BY RENEWAL_LIKELIHOOD_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == SearchSortParamOnboardingStatus {
				query += " ORDER BY ONBOARDING_STATUS_FOR_SORTING " + string(sort.Direction) +
					", ONBOARDING_UPDATED_AT_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == SearchSortParamRenewalCycleNext {
				query += " ORDER BY RENEWAL_CYCLE_NEXT_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == SearchSortParamRenewalDate {
				query += " ORDER BY RENEWAL_DATE_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == "OWNER" {
				query += " ORDER BY OWNER_FIRST_NAME_FOR_SORTING " + string(sort.Direction) + ", OWNER_LAST_NAME_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == SearchSortParamLastTouchpointAt || sort.By == SearchSortParamLastTouchpoint {
				cypherSort.NewSortRule("LAST_TOUCHPOINT_AT", string(sort.Direction), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			} else if sort.By == SearchSortParamLastTouchpointType {
				cypherSort.NewSortRule("LAST_TOUCHPOINT_TYPE", string(sort.Direction), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			} else if sort.By == SearchSortParamContractLengthInMonths {
				query += " ORDER BY CONTRACT_LENGTH_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == SearchSortParamUpdatedAt {
				cypherSort.NewSortRule("UPDATED_AT", string(sort.Direction), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			}
		} else {
			cypherSort.NewSortRule("UPDATED_AT", string(commonmodel.SortingDirectionDesc), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
			query += string(cypherSort.SortingCypherFragment("o"))
		}
		// end sort region
		query += fmt.Sprintf(` RETURN o, contract, op `)
		query += fmt.Sprintf(` SKIP $skip LIMIT $limit`)

		spans.LogKV(log.Object("query", query))
		spans.LogObjectAsJson("params", params)

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
	dbRecordsWithTotalCount.Records = dbRecords.([]*db.Record)
	// each record will contain three nodes, organization, contract and opportunity
	return dbRecordsWithTotalCount, nil
}
