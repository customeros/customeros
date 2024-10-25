package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
	"strings"
	"time"
)

const (
	SearchSortParamOrganization           = "ORGANIZATION"
	SearchSortParamWebsite                = "WEBSITE"
	SearchSortParamEmail                  = "EMAIL"
	SearchSortParamCountry                = "COUNTRY"
	SearchSortParamOnboardingStatus       = "ONBOARDING_STATUS"
	SearchSortParamIsCustomer             = "IS_CUSTOMER"
	SearchSortParamRelationship           = "RELATIONSHIP"
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
	GetDashboardViewOrganizationData(ctx context.Context, tenant string, skip, limit int, where *model.Filter, sort *model.SortBy) (*utils.DbNodesWithTotalCount, error)
	GetDashboardViewRenewalData(ctx context.Context, tenant string, skip, limit int, where *model.Filter, sort *model.SortBy) (*utils.RecordsWithTotalCount, error)
	GetDashboardNewCustomersData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardCustomerMapData(ctx context.Context, tenant string) ([]map[string]interface{}, error)
	GetDashboardRevenueAtRiskData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardMRRPerCustomerData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardARRBreakdownData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardARRBreakdownUpsellsAndDowngradesData(ctx context.Context, tenant, queryType string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardARRBreakdownRenewalsData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardARRBreakdownValueData(ctx context.Context, tenant string, date time.Time) (float64, error)
	GetDashboardRetentionRateContractsRenewalsData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardRetentionRateContractsChurnedData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardAverageTimeToOnboardPerMonth(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardOnboardingCompletionPerMonth(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetDashboardGRRData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error)
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

		orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter(propertyName, *filter.Value.Str, utils.CONTAINS))
		orFilter.Filters = append(orFilter.Filters, utils.CreateCypherFilterEq(propertyName, ""))
		orFilter.Filters = append(orFilter.Filters, utils.CreateCypherFilterIsNull(propertyName))
		return &orFilter
	} else {
		return utils.CreateStringCypherFilter(propertyName, *filter.Value.Str, utils.CONTAINS)
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
	ownerIncludeEmpty := false
	externalId := ""

	//ORGANIZATION, EMAIL, COUNTRY, REGION, LOCALITY
	//region organization filters
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

				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("name", *filter.Filter.Value.Str, utils.CONTAINS))
				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("website", *filter.Filter.Value.Str, utils.CONTAINS))
				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("customerOsId", *filter.Filter.Value.Str, utils.CONTAINS))
				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("referenceId", *filter.Filter.Value.Str, utils.CONTAINS))

				organizationFilter.Filters = append(organizationFilter.Filters, &orFilter)
			} else if filter.Filter.Property == SearchSortParamName {
				organizationFilter.Filters = append(organizationFilter.Filters, createStringCypherFilterWithValueOrEmpty(filter.Filter, "name"))
			} else if filter.Filter.Property == SearchSortParamWebsite {
				organizationFilter.Filters = append(organizationFilter.Filters, createStringCypherFilterWithValueOrEmpty(filter.Filter, "website"))
			} else if filter.Filter.Property == SearchSortParamRelationship && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("relationship", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamStage && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("stage", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamIndustry && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("industry", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamEmail {
				emailFilter.Filters = append(emailFilter.Filters, utils.CreateStringCypherFilter("email", *filter.Filter.Value.Str, utils.CONTAINS))
				emailFilter.Filters = append(emailFilter.Filters, utils.CreateStringCypherFilter("rawEmail", *filter.Filter.Value.Str, utils.CONTAINS))
			} else if filter.Filter.Property == SearchSortParamCountry {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("country", *filter.Filter.Value.Str, utils.EQUALS))
			} else if filter.Filter.Property == SearchSortParamRegion {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("region", *filter.Filter.Value.Str, utils.EQUALS))
			} else if filter.Filter.Property == SearchSortParamLocality {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("locality", *filter.Filter.Value.Str, utils.EQUALS))
			} else if filter.Filter.Property == SearchSortParamOwnerId {
				if filter.Filter.Value.ArrayStr != nil {
					ownerId = *filter.Filter.Value.ArrayStr
				}
				ownerIncludeEmpty = *filter.Filter.IncludeEmpty
			} else if filter.Filter.Property == SearchParamExternalId {
				externalId = *filter.Filter.Value.Str
			} else if filter.Filter.Property == SearchSortParamIsCustomer && filter.Filter.Value.ArrayBool != nil && len(*filter.Filter.Value.ArrayBool) >= 1 {
				if (*filter.Filter.Value.ArrayBool)[0] {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterEq("relationship", neo4jenum.Customer.String()))
				} else {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterNotEq("relationship", neo4jenum.Customer.String()))
				}
			} else if filter.Filter.Property == SearchSortParamRenewalLikelihood && filter.Filter.Value.ArrayStr != nil && len(*filter.Filter.Value.ArrayStr) >= 1 {
				renewalLikelihoodValues := make([]string, 0)
				for _, v := range *filter.Filter.Value.ArrayStr {
					renewalLikelihoodValues = append(renewalLikelihoodValues, mapper.MapOpportunityRenewalLikelihoodFromString(&v))
				}
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("derivedRenewalLikelihood", renewalLikelihoodValues))
			} else if filter.Filter.Property == SearchSortParamOnboardingStatus && filter.Filter.Value.ArrayStr != nil && len(*filter.Filter.Value.ArrayStr) >= 1 {
				onboardingStatusValues := make([]string, 0)
				for _, v := range *filter.Filter.Value.ArrayStr {
					onboardingStatusValues = append(onboardingStatusValues, mapper.MapOnboardingStatusFromString(&v))
				}
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("onboardingStatus", onboardingStatusValues))
			} else if filter.Filter.Property == SearchSortParamRenewalCycleNext && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("billingDetailsRenewalCycleNext", *filter.Filter.Value.Time, utils.LTE))
			} else if filter.Filter.Property == SearchSortParamRenewalDate && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("derivedNextRenewalAt", *filter.Filter.Value.Time, utils.LTE))
			} else if filter.Filter.Property == SearchSortParamForecastArr && filter.Filter.Value.ArrayInt != nil && len(*filter.Filter.Value.ArrayInt) == 2 {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("renewalForecastArr", (*filter.Filter.Value.ArrayInt)[0], utils.GTE))
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("renewalForecastArr", (*filter.Filter.Value.ArrayInt)[1], utils.LTE))
			} else if (filter.Filter.Property == SearchSortParamLastTouchpointAt || filter.Filter.Property == SearchSortParamLastTouchpoint) && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("lastTouchpointAt", *filter.Filter.Value.Time, utils.GTE))
			} else if filter.Filter.Property == SearchSortParamLastTouchpointType && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("lastTouchpointType", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamUpdatedAt && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("updatedAt", *filter.Filter.Value.Time, utils.GTE))
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

	//region count query
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

	span.LogFields(log.String("countQuery", countQuery))

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
	//end count region

	//region query to fetch data
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

	//endregion
	query = query + strings.Join(queryParts, " AND ")

	// sort region
	aliases := " o, d, l"
	query += " WITH o, d, l "
	if sort != nil && sort.By == SearchSortParamOwner {
		if sort.Direction == model.SortingDirectionAsc {
			query += ", CASE WHEN owner.firstName <> \"\" and not owner.firstName is null THEN owner.firstName ELSE 'ZZZZZZZZZZZZZZZZZZZ' END as OWNER_FIRST_NAME_FOR_SORTING "
			query += ", CASE WHEN owner.lastName <> \"\" and not owner.lastName is null THEN owner.lastName ELSE 'ZZZZZZZZZZZZZZZZZZZ' END as OWNER_LAST_NAME_FOR_SORTING "
		} else {
			query += ", CASE WHEN owner.firstName <> \"\" and not owner.firstName is null THEN owner.firstName ELSE 'AAAAAAAAAAAAAAAAAAA' END as OWNER_FIRST_NAME_FOR_SORTING "
			query += ", CASE WHEN owner.lastName <> \"\" and not owner.lastName is null THEN owner.lastName ELSE 'AAAAAAAAAAAAAAAAAAA' END as OWNER_LAST_NAME_FOR_SORTING "
		}
		aliases += ", OWNER_FIRST_NAME_FOR_SORTING, OWNER_LAST_NAME_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamName {
		if sort.Direction == model.SortingDirectionAsc {
			query += ", CASE WHEN o.name <> \"\" and not o.name is null THEN o.name ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as NAME_FOR_SORTING "
		} else {
			query += ", o.name as NAME_FOR_SORTING "
		}
		aliases += ", NAME_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamRenewalLikelihood {
		if sort.Direction == model.SortingDirectionAsc {
			query += ", CASE WHEN o.derivedRenewalLikelihoodOrder IS NOT NULL THEN o.derivedRenewalLikelihoodOrder ELSE 9999 END as RENEWAL_LIKELIHOOD_FOR_SORTING "
		} else {
			query += ", CASE WHEN o.derivedRenewalLikelihoodOrder IS NOT NULL THEN o.derivedRenewalLikelihoodOrder ELSE -1 END as RENEWAL_LIKELIHOOD_FOR_SORTING "
		}
		aliases += ", RENEWAL_LIKELIHOOD_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamRelationship {
		if sort.Direction == model.SortingDirectionAsc {
			query += ", CASE WHEN o.relationship <> '' AND NOT o.relationship IS NULL THEN o.relationship ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as RELATIONSHIP_FOR_SORTING "
		} else {
			query += ", o.relationship as RELATIONSHIP_FOR_SORTING "
		}
		aliases += ", RELATIONSHIP_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamStage {
		if sort.Direction == model.SortingDirectionAsc {
			query += ", CASE WHEN o.stage <> '' AND NOT o.stage IS NULL THEN o.stage ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as STAGE_FOR_SORTING "
		} else {
			query += ", o.stage as STAGE_FOR_SORTING "
		}
		aliases += ", STAGE_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamIndustry {
		if sort.Direction == model.SortingDirectionAsc {
			query += ", CASE WHEN o.industry <> '' AND NOT o.industry IS NULL THEN o.industry ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as STAGE_FOR_SORTING "
		} else {
			query += ", o.industry as INDUSTRY_FOR_SORTING "
		}
		aliases += ", INDUSTRY_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamRenewalCycleNext {
		if sort.Direction == model.SortingDirectionAsc {
			query += ", CASE WHEN o.billingDetailsRenewalCycleNext IS NOT NULL THEN date(o.billingDetailsRenewalCycleNext) ELSE date('2100-01-01') END as RENEWAL_CYCLE_NEXT_FOR_SORTING "
		} else {
			query += ", CASE WHEN o.billingDetailsRenewalCycleNext IS NOT NULL THEN date(o.billingDetailsRenewalCycleNext) ELSE date('1900-01-01') END as RENEWAL_CYCLE_NEXT_FOR_SORTING "
		}
		aliases += ", RENEWAL_CYCLE_NEXT_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamRenewalDate {
		if sort.Direction == model.SortingDirectionAsc {
			query += ", CASE WHEN o.derivedNextRenewalAt IS NOT NULL THEN date(o.derivedNextRenewalAt) ELSE date('2100-01-01') END as RENEWAL_DATE_FOR_SORTING "
		} else {
			query += ", CASE WHEN o.derivedNextRenewalAt IS NOT NULL THEN date(o.derivedNextRenewalAt) ELSE date('1900-01-01') END as RENEWAL_DATE_FOR_SORTING "
		}
		aliases += ", RENEWAL_DATE_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamChurnDate {
		if sort.Direction == model.SortingDirectionAsc {
			query += ", CASE WHEN o.derivedChurnedAt IS NOT NULL THEN date(o.derivedChurnedAt) ELSE date('2100-01-01') END as CHURN_DATE_FOR_SORTING "
		} else {
			query += ", CASE WHEN o.derivedChurnedAt IS NOT NULL THEN date(o.derivedChurnedAt) ELSE date('1900-01-01') END as CHURN_DATE_FOR_SORTING "
		}
		aliases += ", RENEWAL_DATE_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamOnboardingStatus {
		if sort.Direction == model.SortingDirectionAsc {
			query += ", CASE WHEN o.onboardingStatusOrder IS NOT NULL THEN o.onboardingStatusOrder ELSE 9999 END as ONBOARDING_STATUS_FOR_SORTING "
			query += ", o.onboardingUpdatedAt AS ONBOARDING_UPDATED_AT_FOR_SORTING "
		} else {
			query += ", CASE WHEN o.onboardingStatusOrder IS NOT NULL THEN o.onboardingStatusOrder ELSE -1 END as ONBOARDING_STATUS_FOR_SORTING "
			query += ", o.onboardingUpdatedAt AS ONBOARDING_UPDATED_AT_FOR_SORTING "
		}
		aliases += ", ONBOARDING_STATUS_FOR_SORTING, ONBOARDING_UPDATED_AT_FOR_SORTING "
	}
	if sort != nil && sort.By == SearchSortParamForecastArr {
		if sort.Direction == model.SortingDirectionAsc {
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
		} else if sort.By == SearchSortParamRelationship {
			query += " ORDER BY RELATIONSHIP_FOR_SORTING " + string(sort.Direction)
		} else if sort.By == SearchSortParamStage {
			query += " ORDER BY STAGE_FOR_SORTING " + string(sort.Direction)
		} else if sort.By == SearchSortParamIndustry {
			query += " ORDER BY INDUSTRY_FOR_SORTING " + string(sort.Direction)
		} else if sort.By == SearchSortParamOrganization {
			cypherSort.NewSortRule("NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithCoalesce().WithAlias("parent")
			cypherSort.NewSortRule("NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithCoalesce()
			cypherSort.NewSortRule("NAME", sort.Direction.String(), true, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithAlias("parent").WithDescending()
			cypherSort.NewSortRule("NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
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
			cypherSort.NewSortRule("DOMAIN", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.DomainEntity{}))
			query += string(cypherSort.SortingCypherFragment("d"))
		} else if sort.By == SearchSortParamLocation {
			cypherSort.NewSortRule("COUNTRY", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.LocationEntity{}))
			cypherSort.NewSortRule("REGION", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.LocationEntity{}))
			cypherSort.NewSortRule("LOCALITY", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.LocationEntity{}))
			query += string(cypherSort.SortingCypherFragment("l"))
		} else if sort.By == "OWNER" {
			if sort.CaseSensitive != nil && *sort.CaseSensitive {
				query += " ORDER BY OWNER_FIRST_NAME_FOR_SORTING " + string(sort.Direction) + ", OWNER_LAST_NAME_FOR_SORTING " + string(sort.Direction)
			} else {
				query += " ORDER BY toLower(OWNER_FIRST_NAME_FOR_SORTING) " + string(sort.Direction) + ", toLower(OWNER_LAST_NAME_FOR_SORTING) " + string(sort.Direction)
			}
		} else if sort.By == SearchSortParamLastTouchpointAt || sort.By == SearchSortParamLastTouchpoint {
			cypherSort.NewSortRule("LAST_TOUCHPOINT_AT", sort.Direction.String(), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
			query += string(cypherSort.SortingCypherFragment("o"))
		} else if sort.By == SearchSortParamLastTouchpointType {
			cypherSort.NewSortRule("LAST_TOUCHPOINT_TYPE", sort.Direction.String(), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
			query += string(cypherSort.SortingCypherFragment("o"))
		} else if sort.By == SearchSortParamUpdatedAt {
			cypherSort.NewSortRule("UPDATED_AT", sort.Direction.String(), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
			query += string(cypherSort.SortingCypherFragment("o"))
		}
	} else {
		cypherSort.NewSortRule("UPDATED_AT", string(model.SortingDirectionDesc), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
		query += string(cypherSort.SortingCypherFragment("o"))
	}
	// end sort region
	query += fmt.Sprintf(` RETURN distinct(o) `)
	query += fmt.Sprintf(` SKIP $skip LIMIT $limit`)

	span.LogFields(log.String("query", query))
	tracing.LogObjectAsJson(span, "params", params)

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

func (r *dashboardRepository) GetDashboardViewRenewalData(ctx context.Context, tenant string, skip, limit int, where *model.Filter, sort *model.SortBy) (*utils.RecordsWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardViewRenewalData")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Int("skip", skip), log.Int("limit", limit))
	tracing.LogObjectAsJson(span, "where", where)
	tracing.LogObjectAsJson(span, "sort", sort)

	dbRecordsWithTotalCount := new(utils.RecordsWithTotalCount)

	organizationFilterCypher, organizationFilterParams := "", make(map[string]interface{})
	contractFilterCypher, contractFilterParams := "", make(map[string]interface{})
	opportunityFilterCypher, opportunityFilterParams := "", make(map[string]interface{})
	emailFilterCypher, emailFilterParams := "", make(map[string]interface{})
	locationFilterCypher, locationFilterParams := "", make(map[string]interface{})

	ownerId := []string{}
	ownerIncludeEmpty := false

	//ORGANIZATION, EMAIL, COUNTRY, REGION, LOCALITY
	//region organization & contract filters
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

				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("name", *filter.Filter.Value.Str, utils.CONTAINS))
				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("website", *filter.Filter.Value.Str, utils.CONTAINS))
				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("customerOsId", *filter.Filter.Value.Str, utils.CONTAINS))
				orFilter.Filters = append(orFilter.Filters, utils.CreateStringCypherFilter("referenceId", *filter.Filter.Value.Str, utils.CONTAINS))

				organizationFilter.Filters = append(organizationFilter.Filters, &orFilter)
			} else if filter.Filter.Property == SearchSortParamName {
				organizationFilter.Filters = append(organizationFilter.Filters, createStringCypherFilterWithValueOrEmpty(filter.Filter, "name"))
			} else if filter.Filter.Property == SearchSortParamWebsite {
				organizationFilter.Filters = append(organizationFilter.Filters, createStringCypherFilterWithValueOrEmpty(filter.Filter, "website"))
			} else if filter.Filter.Property == SearchSortParamRelationship && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("relationship", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamStage && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("stage", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamIndustry && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("industry", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamEmail {
				emailFilter.Filters = append(emailFilter.Filters, utils.CreateStringCypherFilter("email", *filter.Filter.Value.Str, utils.CONTAINS))
				emailFilter.Filters = append(emailFilter.Filters, utils.CreateStringCypherFilter("rawEmail", *filter.Filter.Value.Str, utils.CONTAINS))
			} else if filter.Filter.Property == SearchSortParamCountry {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("country", *filter.Filter.Value.Str, utils.EQUALS))
			} else if filter.Filter.Property == SearchSortParamRegion {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("region", *filter.Filter.Value.Str, utils.EQUALS))
			} else if filter.Filter.Property == SearchSortParamLocality {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateStringCypherFilter("locality", *filter.Filter.Value.Str, utils.EQUALS))
			} else if filter.Filter.Property == SearchSortParamOwnerId {
				if filter.Filter.Value.ArrayStr != nil {
					ownerId = *filter.Filter.Value.ArrayStr
				}
				ownerIncludeEmpty = *filter.Filter.IncludeEmpty
			} else if filter.Filter.Property == SearchSortParamIsCustomer && filter.Filter.Value.ArrayBool != nil && len(*filter.Filter.Value.ArrayBool) >= 1 {
				if (*filter.Filter.Value.ArrayBool)[0] {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterEq("relationship", neo4jenum.Customer.String()))
				} else {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterNotEq("relationship", neo4jenum.Customer.String()))
				}
			} else if filter.Filter.Property == SearchSortParamRenewalLikelihood && filter.Filter.Value.ArrayStr != nil && len(*filter.Filter.Value.ArrayStr) >= 1 {
				renewalLikelihoodValues := make([]string, 0)
				for _, v := range *filter.Filter.Value.ArrayStr {
					renewalLikelihoodValues = append(renewalLikelihoodValues, mapper.MapOpportunityRenewalLikelihoodFromString(&v))
				}
				opportunityFilter.Filters = append(opportunityFilter.Filters, utils.CreateCypherFilterIn("renewalLikelihood", renewalLikelihoodValues))
			} else if filter.Filter.Property == SearchSortParamOnboardingStatus && filter.Filter.Value.ArrayStr != nil && len(*filter.Filter.Value.ArrayStr) >= 1 {
				onboardingStatusValues := make([]string, 0)
				for _, v := range *filter.Filter.Value.ArrayStr {
					onboardingStatusValues = append(onboardingStatusValues, mapper.MapOnboardingStatusFromString(&v))
				}
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("onboardingStatus", onboardingStatusValues))
			} else if filter.Filter.Property == SearchSortParamRenewalCycleNext && filter.Filter.Value.Time != nil {
				opportunityFilter.Filters = append(opportunityFilter.Filters, utils.CreateCypherFilter("renewedAt", *filter.Filter.Value.Time, utils.LTE))
			} else if filter.Filter.Property == SearchSortParamRenewalDate && filter.Filter.Value.Time != nil {
				opportunityFilter.Filters = append(opportunityFilter.Filters, utils.CreateCypherFilter("renewedAt", *filter.Filter.Value.Time, utils.LTE))
			} else if filter.Filter.Property == SearchSortParamForecastArr && filter.Filter.Value.ArrayInt != nil && len(*filter.Filter.Value.ArrayInt) == 2 {
				opportunityFilter.Filters = append(opportunityFilter.Filters, utils.CreateCypherFilter("maxAmount", (*filter.Filter.Value.ArrayInt)[0], utils.GTE))
				opportunityFilter.Filters = append(opportunityFilter.Filters, utils.CreateCypherFilter("maxAmount", (*filter.Filter.Value.ArrayInt)[1], utils.LTE))
			} else if (filter.Filter.Property == SearchSortParamLastTouchpointAt || filter.Filter.Property == SearchSortParamLastTouchpoint) && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("lastTouchpointAt", *filter.Filter.Value.Time, utils.GTE))
			} else if filter.Filter.Property == SearchSortParamLastTouchpointType && filter.Filter.Value.ArrayStr != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn("lastTouchpointType", *filter.Filter.Value.ArrayStr))
			} else if filter.Filter.Property == SearchSortParamUpdatedAt && filter.Filter.Value.Time != nil {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter("updatedAt", *filter.Filter.Value.Time, utils.GTE))
			} else if filter.Filter.Property == SearchSortParamRenewalCycle {
				if filter.Filter.Value.Str != nil {
					switch *filter.Filter.Value.Str {
					case "MONTHLY":
						contractFilter.Filters = append(contractFilter.Filters, utils.CreateCypherFilter("lengthInMonths", 1, utils.EQUALS))
					case "QUARTERLY":
						contractFilter.Filters = append(contractFilter.Filters, utils.CreateCypherFilter("lengthInMonths", 3, utils.EQUALS))
					case "ANNUALLY":
						contractFilter.Filters = append(contractFilter.Filters, utils.CreateCypherFilter("lengthInMonths", 12, utils.GTE))
					}
				}
			} else if filter.Filter.Property == SearchSortParamContractLengthInMonths {
				if filter.Filter.Value.Int != nil {
					contractFilter.Filters = append(contractFilter.Filters, utils.CreateCypherFilter("lengthInMonths", *filter.Filter.Value.Int, utils.EQUALS))
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
		utils.MergeMapToMap(contractFilterParams, params)
		utils.MergeMapToMap(opportunityFilterParams, params)

		//region count query
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

		span.LogFields(log.String("countQuery", countQuery))

		countQueryResult, err := tx.Run(ctx, countQuery, params)
		if err != nil {
			return nil, err
		}

		countRecord, err := countQueryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		dbRecordsWithTotalCount.Count = countRecord.Values[0].(int64)
		//endregion

		//region query to fetch data
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

		//endregion
		query = query + strings.Join(queryParts, " AND ")

		// sort region
		aliases := " o, contract, op, owner"
		query += " WITH o, contract, op, owner "
		if sort != nil && sort.By == SearchSortParamOwner {
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN owner.firstName <> \"\" and not owner.firstName is null THEN owner.firstName ELSE 'ZZZZZZZZZZZZZZZZZZZ' END as OWNER_FIRST_NAME_FOR_SORTING "
				query += ", CASE WHEN owner.lastName <> \"\" and not owner.lastName is null THEN owner.lastName ELSE 'ZZZZZZZZZZZZZZZZZZZ' END as OWNER_LAST_NAME_FOR_SORTING "
			} else {
				query += ", CASE WHEN owner.firstName <> \"\" and not owner.firstName is null THEN owner.firstName ELSE 'AAAAAAAAAAAAAAAAAAA' END as OWNER_FIRST_NAME_FOR_SORTING "
				query += ", CASE WHEN owner.lastName <> \"\" and not owner.lastName is null THEN owner.lastName ELSE 'AAAAAAAAAAAAAAAAAAA' END as OWNER_LAST_NAME_FOR_SORTING "
			}
			aliases += ", OWNER_FIRST_NAME_FOR_SORTING, OWNER_LAST_NAME_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamName {
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN o.name <> \"\" and not o.name is null THEN o.name ELSE 'ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ' END as NAME_FOR_SORTING "
			} else {
				query += ", o.name as NAME_FOR_SORTING "
			}
			aliases += ", NAME_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamRenewalLikelihood {
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN op.renewalLikelihood IS NOT NULL THEN op.renewalLikelihood ELSE 9999 END as RENEWAL_LIKELIHOOD_FOR_SORTING "
			} else {
				query += ", CASE WHEN op.renewalLikelihood IS NOT NULL THEN op.renewalLikelihood ELSE -1 END as RENEWAL_LIKELIHOOD_FOR_SORTING "
			}
			aliases += ", RENEWAL_LIKELIHOOD_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamRenewalCycleNext {
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN op.renewedAt IS NOT NULL THEN date(op.renewedAt) ELSE date('2100-01-01') END as RENEWAL_CYCLE_NEXT_FOR_SORTING "
			} else {
				query += ", CASE WHEN op.renewedAt IS NOT NULL THEN date(op.renewedAt) ELSE date('1900-01-01') END as RENEWAL_CYCLE_NEXT_FOR_SORTING "
			}
			aliases += ", RENEWAL_CYCLE_NEXT_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamRenewalDate {
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN op.renewedAt IS NOT NULL THEN date(op.renewedAt) ELSE date('2100-01-01') END as RENEWAL_DATE_FOR_SORTING "
			} else {
				query += ", CASE WHEN op.renewedAt IS NOT NULL THEN date(op.renewedAt) ELSE date('1900-01-01') END as RENEWAL_DATE_FOR_SORTING "
			}
			aliases += ", RENEWAL_DATE_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamContractLengthInMonths {
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN contract.lengthInMonths IS NOT NULL AND contract.lengthInMonths > 0 THEN contract.lengthInMonths ELSE 9999 END as CONTRACT_LENGTH_FOR_SORTING "
			} else {
				query += ", CASE WHEN contract.lengthInMonths IS NOT NULL THEN contract.lengthInMonths ELSE -1 END as CONTRACT_LENGTH_FOR_SORTING "
			}
			aliases += ", RENEWAL_LIKELIHOOD_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamOnboardingStatus {
			if sort.Direction == model.SortingDirectionAsc {
				query += ", CASE WHEN o.onboardingStatusOrder IS NOT NULL THEN o.onboardingStatusOrder ELSE 9999 END as ONBOARDING_STATUS_FOR_SORTING "
				query += ", o.onboardingUpdatedAt AS ONBOARDING_UPDATED_AT_FOR_SORTING "
			} else {
				query += ", CASE WHEN o.onboardingStatusOrder IS NOT NULL THEN o.onboardingStatusOrder ELSE -1 END as ONBOARDING_STATUS_FOR_SORTING "
				query += ", o.onboardingUpdatedAt AS ONBOARDING_UPDATED_AT_FOR_SORTING "
			}
			aliases += ", ONBOARDING_STATUS_FOR_SORTING, ONBOARDING_UPDATED_AT_FOR_SORTING "
		}
		if sort != nil && sort.By == SearchSortParamForecastArr {
			if sort.Direction == model.SortingDirectionAsc {
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
				cypherSort.NewSortRule("NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithCoalesce().WithAlias("parent")
				cypherSort.NewSortRule("NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithCoalesce()
				cypherSort.NewSortRule("NAME", sort.Direction.String(), true, reflect.TypeOf(neo4jentity.OrganizationEntity{})).WithAlias("parent").WithDescending()
				cypherSort.NewSortRule("NAME", sort.Direction.String(), *sort.CaseSensitive, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
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
				cypherSort.NewSortRule("LAST_TOUCHPOINT_AT", sort.Direction.String(), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			} else if sort.By == SearchSortParamLastTouchpointType {
				cypherSort.NewSortRule("LAST_TOUCHPOINT_TYPE", sort.Direction.String(), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			} else if sort.By == SearchSortParamContractLengthInMonths {
				query += " ORDER BY CONTRACT_LENGTH_FOR_SORTING " + string(sort.Direction)
			} else if sort.By == SearchSortParamUpdatedAt {
				cypherSort.NewSortRule("UPDATED_AT", sort.Direction.String(), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
				query += string(cypherSort.SortingCypherFragment("o"))
			}
		} else {
			cypherSort.NewSortRule("UPDATED_AT", string(model.SortingDirectionDesc), false, reflect.TypeOf(neo4jentity.OrganizationEntity{}))
			query += string(cypherSort.SortingCypherFragment("o"))
		}
		// end sort region
		query += fmt.Sprintf(` RETURN o, contract, op `)
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
	dbRecordsWithTotalCount.Records = dbRecords.([]*db.Record)
	//each record will contain three nodes, organization, contract and opportunity
	return dbRecordsWithTotalCount, nil
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
					  o.hide = false AND o.relationship = $customerRelationship AND
					  i.serviceStartedAt.year = year AND 
					  i.serviceStartedAt.month = month AND 
					  (i.endedAt IS NULL OR i.endedAt > endOfMonth)
					
					WITH o, year, month, MIN(i.serviceStartedAt) AS oldestContractDate
					OPTIONAL MATCH (o)-[:HAS_CONTRACT]->(oldest:Contract_%s)
					WHERE oldest.serviceStartedAt = oldestContractDate
					RETURN year, month, COUNT(oldest) AS totalContracts
				`, "% 12 + 1", tenant, tenant, tenant),
			map[string]any{
				"tenant":               tenant,
				"startDate":            startDate,
				"endDate":              endDate,
				"customerRelationship": neo4jenum.Customer,
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
					MATCH (t:Tenant{name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[r]->(op:Opportunity_%s)
					WHERE 	o.hide = false AND 
							o.relationship = $customerRelationship AND 
							c.serviceStartedAt IS NOT NULL AND 
							NOT c.status IN [$contractStatusDraft] AND 
							op.internalType = $opportunityInternalTypeRenewal AND 
							op.internalStage in [$opportunityInternalStageOpen]

					WITH o, c, op, COLLECT(type(r)) AS relTypes

					WITH o.id AS oid,
						COLLECT(DISTINCT CASE
							WHEN c.status = $contractStatusEnded THEN 'CHURNED'
							WHEN c.status IN [$contractStatusLive,$contractStatusScheduled] AND 'ACTIVE_RENEWAL' in relTypes AND op.renewalLikelihood = $likelihoodHigh THEN 'OK'
							WHEN op.renewalLikelihood IN [$likelihoodLow, $likelihoodZero] THEN 'HIGH_RISK'
							ELSE 'MEDIUM_RISK'
						END) AS statuses,
						COLLECT(DISTINCT { serviceStartedAt: c.serviceStartedAt }) AS contractsStartedAt
					
					WITH oid, CASE
						WHEN ALL(x IN statuses WHERE x = 'CHURNED') THEN 'CHURNED'
						WHEN ALL(x IN statuses WHERE x IN ['OK', 'CHURNED']) THEN 'OK'
						WHEN ALL(x IN statuses WHERE x IN ['OK', 'CHURNED', 'MEDIUM_RISK']) THEN 'MEDIUM_RISK'
						ELSE 'HIGH_RISK'
					END AS status,
					REDUCE(s = null, cs IN contractsStartedAt | 
						CASE WHEN s IS NULL OR cs.serviceStartedAt < s THEN cs.serviceStartedAt ELSE s END
					) AS oldestServiceStartedAt
					
					MATCH (o:Organization_%s{id:oid})-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WHERE sli.endedAt IS NULL AND (sli.isCanceled IS NULL OR sli.isCanceled = false)
						AND (sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
					
					WITH oid, oldestServiceStartedAt, status, c, sli
					ORDER BY c.serviceStartedAt DESC
					
					WITH oid, oldestServiceStartedAt, status, c, COLLECT(sli) AS sliListPerContract
					ORDER BY CASE WHEN status = 'CHURNED' THEN c.endedAt ELSE c.serviceStartedAt END DESC
					
					WITH oid, oldestServiceStartedAt, status, COLLECT({cId: c.id, cStatus: c.status, sliList: sliListPerContract}) AS contractsDetails
					
					WITH oid, oldestServiceStartedAt, status, CASE WHEN status = 'CHURNED' THEN [contractsDetails[0]]
						ELSE REDUCE(a = [], c IN contractsDetails | 
							CASE WHEN c.cStatus IN ['LIVE','DRAFT'] THEN a + c ELSE a END
						) END AS contracts
					
					WITH oid, oldestServiceStartedAt, status, REDUCE(s = [], c IN contracts | 
						s + REDUCE(a = 0, sli IN c.sliList | 
							a + CASE WHEN sli.billed = 'MONTHLY' THEN sli.price * sli.quantity * 12
							ELSE CASE WHEN sli.billed = 'QUARTERLY' THEN sli.price * sli.quantity * 4
							ELSE CASE WHEN sli.billed = 'ANNUALLY' THEN sli.price * sli.quantity
							ELSE 0 END END END
						)
					) AS arrList
					
					RETURN oid, oldestServiceStartedAt, status, REDUCE(a = TOFLOAT(0), arr IN arrList | a + arr) AS arr
					ORDER BY oldestServiceStartedAt ASC
				`, tenant, tenant, tenant, tenant, tenant, tenant),
			map[string]any{
				"tenant":                         tenant,
				"contractStatusLive":             neo4jenum.ContractStatusLive.String(),
				"contractStatusDraft":            neo4jenum.ContractStatusDraft.String(),
				"contractStatusEnded":            neo4jenum.ContractStatusEnded.String(),
				"contractStatusScheduled":        neo4jenum.ContractStatusScheduled.String(),
				"contractStatusOutOfContract":    neo4jenum.ContractStatusOutOfContract.String(),
				"opportunityInternalTypeRenewal": neo4jenum.OpportunityInternalTypeRenewal.String(),
				"opportunityInternalStageOpen":   neo4jenum.OpportunityInternalStageOpen.String(),
				"likelihoodHigh":                 neo4jenum.RenewalLikelihoodHigh.String(),
				"likelihoodMedium":               neo4jenum.RenewalLikelihoodMedium.String(),
				"likelihoodLow":                  neo4jenum.RenewalLikelihoodLow.String(),
				"likelihoodZero":                 neo4jenum.RenewalLikelihoodZero.String(),
				"customerRelationship":           neo4jenum.Customer.String(),
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
			oldestServiceStartedAt := utils.GetTimePropFromNeo4jOrZeroTime(v.Values[1])
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
						o.hide = false AND o.relationship = $customerRelationship AND c.status = 'LIVE' AND op.internalType = 'RENEWAL' AND op.internalStage in ['OPEN']
					
					WITH COLLECT(DISTINCT { renewalLikelihood: op.renewalLikelihood, maxAmount: op.maxAmount }) AS contractDetails
					
					return 
						REDUCE(sumHigh = 0, cd IN contractDetails | CASE WHEN cd.renewalLikelihood = 'HIGH' THEN sumHigh + cd.maxAmount ELSE sumHigh END ) AS high,
						REDUCE(sumAtRisk = 0, cd IN contractDetails | CASE WHEN cd.renewalLikelihood <> 'HIGH' THEN sumAtRisk + cd.maxAmount ELSE sumAtRisk END ) AS atRisk
				`, tenant, tenant, tenant),
			map[string]any{
				"tenant":               tenant,
				"startDate":            startDate,
				"endDate":              endDate,
				"customerRelationship": neo4jenum.Customer.String(),
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
									day: 1, hour: 0, minute: 0, second: 0, nanosecond: 0o00000000}) AS currentDate
					
					WITH currentDate as beginOfMonth,
						 currentDate + duration({months: 1}) as startOfNextMonth
					
					WITH DISTINCT beginOfMonth.year AS year, beginOfMonth.month AS month, beginOfMonth, startOfNextMonth
					
					OPTIONAL MATCH (t:Tenant{name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WHERE 
						o.hide = false AND o.relationship = $customerRelationship AND (c.endedAt IS NULL or c.endedAt >= startOfNextMonth) AND (sli.billed = 'MONTHLY' or sli.billed = 'QUARTERLY' or sli.billed = 'ANNUALLY') AND
						sli.startedAt < startOfNextMonth AND (sli.endedAt IS NULL OR sli.endedAt >= startOfNextMonth)
					
					WITH year, month, startOfNextMonth, COLLECT(DISTINCT { id: sli.id, startedAt: sli.startedAt, endedAt: sli.endedAt, amountPerMonth: CASE WHEN sli.billed = 'MONTHLY' THEN sli.price * sli.quantity ELSE CASE WHEN sli.billed = 'QUARTERLY' THEN  sli.price * sli.quantity / 3 ELSE CASE WHEN sli.billed = 'ANNUALLY' THEN sli.price * sli.quantity / 12 ELSE 0 END END END }) AS contractDetails
					
					WITH year, month,  startOfNextMonth, contractDetails, REDUCE(sumHigh = 0, cd IN contractDetails | sumHigh + cd.amountPerMonth ) AS mrr
					
					return year, month, mrr
				`, "% 12 + 1", tenant, tenant, tenant),
			map[string]any{
				"tenant":               tenant,
				"startDate":            startDate,
				"endDate":              endDate,
				"customerRelationship": neo4jenum.Customer.String(),
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
					WHERE o.hide = false AND o.relationship = $customerRelationship AND sli.startedAt IS NOT NULL AND (sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
					
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
				"tenant":               tenant,
				"startDate":            startDate,
				"endDate":              endDate,
				"customerRelationship": neo4jenum.Customer.String(),
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
					
					OPTIONAL MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WHERE o.hide = false AND o.relationship = $customerRelationship AND (sli.endedAt IS NULL OR sli.endedAt > beginOfMonth) AND sli.startedAt < startOfNextMonth AND (sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
					
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
				`, "% 12 + 1", tenant, tenant, tenant, q),
			map[string]any{
				"tenant":               tenant,
				"startDate":            startDate,
				"endDate":              endDate,
				"customerRelationship": neo4jenum.Customer.String(),
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
					
					OPTIONAL MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WITH year, month, beginOfMonth, endOfMonth, c.serviceStartedAt AS cssa, c.lengthInMonths AS clim, sli ORDER BY sli.startedAt ASC

					WHERE 	o.hide = false AND 
							o.relationship = $customerRelationship AND 
							c.serviceStartedAt IS NOT NULL AND 
							(c.endedAt IS NULL OR c.endedAt > beginOfMonth) AND 
							(sli.endedAt IS NULL OR sli.endedAt > beginOfMonth) AND 
							sli.startedAt < beginOfMonth AND 
							(sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
					WITH year, month, beginOfMonth, endOfMonth, cssa, clim, sli.parentId AS pp, COLLECT(sli) AS versions
					WITH year, month, beginOfMonth, endOfMonth, cssa, clim, pp, LAST(versions) AS lastSliVersion
					WITH year, month, beginOfMonth, endOfMonth, cssa, clim, pp, lastSliVersion
					WHERE
						CASE WHEN clim >= 12 THEN cssa.YEAR < beginOfMonth.YEAR AND cssa.MONTH = beginOfMonth.MONTH AND (beginOfMonth.YEAR - cssa.YEAR) %s = 0
						ELSE 1 = 1 END AND
						CASE WHEN clim = 3 THEN
							(lastSliVersion.billed IN ['MONTHLY', 'QUARTERLY'] AND beginOfMonth.MONTH IN [cssa.MONTH - 9, cssa.MONTH - 6, cssa.MONTH - 3, cssa.MONTH, cssa.MONTH + 3, cssa.MONTH + 6, cssa.MONTH + 9]) OR
							(lastSliVersion.billed = 'ANNUALLY' AND beginOfMonth.MONTH = cssa.MONTH)
						ELSE 1 = 1 END AND
						CASE WHEN clim = 1 THEN
							lastSliVersion.billed = 'MONTHLY' OR
							(lastSliVersion.billed = 'QUARTERLY' AND beginOfMonth.MONTH IN [cssa.MONTH - 9, cssa.MONTH - 6, cssa.MONTH - 3, cssa.MONTH, cssa.MONTH + 3, cssa.MONTH + 6, cssa.MONTH + 9]) OR
							(lastSliVersion.billed = 'ANNUALLY' AND beginOfMonth.MONTH = cssa.MONTH)
						ELSE 1 = 1 END
					
					WITH year, month, clim, COLLECT(lastSliVersion) AS lasts
					WITH year, month, 
						REDUCE(s = 0.0, a IN lasts | s + 
							TOFLOAT(
								CASE WHEN clim = 0 THEN 0.0 ELSE
									CASE a.billed 
										WHEN 'MONTHLY' THEN 
                    						clim
										WHEN 'QUARTERLY' THEN 
                    						CASE WHEN clim / 3 = 0 THEN 1 ELSE clim / 3 END
										WHEN 'ANNUALLY' THEN 
                    						CASE WHEN clim / 12 = 0 THEN 1 ELSE clim / 12 END
										ELSE 0.0
									END
								END * a.price * a.quantity)
							) AS amount

					RETURN year, month, SUM(amount)
				`, "% 12 + 1", tenant, tenant, tenant, "% (clim/12)"),
			map[string]any{
				"tenant":               tenant,
				"startDate":            startDate,
				"endDate":              endDate,
				"customerRelationship": neo4jenum.Customer.String(),
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

func (r *dashboardRepository) GetDashboardARRBreakdownValueData(ctx context.Context, tenant string, date time.Time) (float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardARRBreakdownValueData")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("date", date))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			`
					WITH $date AS date

					WITH datetime({
						YEAR: date.YEAR,
						MONTH: date.MONTH,
						 DAY: 1,
							 HOUR: 0,
							 MINUTE: 0,
							 SECOND: 0,
							 NANOSECOND: 0o00000000
					}) AS beginOfMonth
					
					WITH beginOfMonth,
						 beginOfMonth + duration({MONTHS: 1}) - duration({NANOSECONDS: 1}) AS endOfMonth,
						 beginOfMonth + duration({MONTHS: 1}) AS startOfNextMonth
					
					OPTIONAL MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WITH beginOfMonth, endOfMonth, c.id AS cid, sli ORDER BY sli.startedAt ASC

					WHERE 
					o.hide = false AND o.relationship = $customerRelationship
					AND c.serviceStartedAt IS NOT NULL 
					AND sli.startedAt < startOfNextMonth 
					AND (c.endedAt IS NULL OR c.endedAt >= startOfNextMonth) 
					AND (sli.endedAt IS NULL OR sli.endedAt >= startOfNextMonth) 
					AND (sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
				
					WITH CASE WHEN sli.billed = 'MONTHLY' THEN 12 ELSE (CASE WHEN sli.billed = 'QUARTERLY' THEN 4 ELSE (CASE WHEN sli.billed = 'ANNUALLY' THEN 1 ELSE 0 END) END) END * sli.price * sli.quantity AS amount
						
					RETURN SUM(amount)
				`, tenant, tenant, tenant),
			map[string]any{
				"tenant":               tenant,
				"date":                 date,
				"customerRelationship": neo4jenum.Customer.String(),
			})
		if err != nil {
			return nil, err
		}

		return queryResult.Collect(ctx)
	})
	if err != nil {
		return 0.0, err
	}

	if dbRecords != nil {
		for _, v := range dbRecords.([]*neo4j.Record) {
			return getCorrectValueType(v.Values[0]), nil
		}
	}

	return 0.0, nil
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
						 currentDate + duration({MONTHS: 1}) - duration({NANOSECONDS: 1}) AS endOfMonth
					
					WITH DISTINCT currentDate.YEAR AS year, currentDate.MONTH AS month, beginOfMonth, endOfMonth
					
					OPTIONAL MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WITH year, month, beginOfMonth, endOfMonth, c.id as cid, c.serviceStartedAt AS cssa, c.lengthInMonths AS clim, sli ORDER BY sli.startedAt ASC

					WHERE o.hide = false AND o.relationship = $customerRelationship AND c.serviceStartedAt IS NOT NULL AND (c.endedAt IS NULL OR c.endedAt > endOfMonth) AND (sli.endedAt IS NULL OR sli.endedAt > beginOfMonth) AND sli.startedAt < beginOfMonth AND (sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
					WITH year, month, beginOfMonth, endOfMonth, cid, cssa, clim, sli.parentId AS pp, COLLECT(sli) AS versions
					WITH year, month, beginOfMonth, endOfMonth, cid, cssa, clim, pp, LAST(versions) AS lastSliVersion
					WITH year, month, beginOfMonth, endOfMonth, cid, cssa, clim, pp, lastSliVersion
					WHERE
						CASE WHEN clim >= 12 THEN cssa.YEAR < beginOfMonth.YEAR AND cssa.MONTH = beginOfMonth.MONTH AND (beginOfMonth.YEAR - cssa.YEAR) %s = 0
						ELSE 1 = 1 END AND
						CASE WHEN clim = 3 THEN
							(lastSliVersion.billed IN ['MONTHLY', 'QUARTERLY'] AND beginOfMonth.MONTH IN [cssa.MONTH - 9, cssa.MONTH - 6, cssa.MONTH - 3, cssa.MONTH, cssa.MONTH + 3, cssa.MONTH + 6, cssa.MONTH + 9]) OR
							(lastSliVersion.billed = 'ANNUALLY' AND beginOfMonth.MONTH = cssa.MONTH)
						ELSE 1 = 1 END AND
						CASE WHEN clim = 1 THEN
							lastSliVersion.billed = 'MONTHLY' OR
							(lastSliVersion.billed = 'QUARTERLY' AND beginOfMonth.MONTH IN [cssa.MONTH - 9, cssa.MONTH - 6, cssa.MONTH - 3, cssa.MONTH, cssa.MONTH + 3, cssa.MONTH + 6, cssa.MONTH + 9]) OR
							(lastSliVersion.billed = 'ANNUALLY' AND beginOfMonth.MONTH = cssa.MONTH)
						ELSE 1 = 1 END
					
					WITH year, month, cid, lastSliVersion
					return year, month, COUNT(DISTINCT(cid)) AS contractsWithRenewals
				`, "% 12 + 1", tenant, tenant, tenant, "% (clim/12)"),
			map[string]any{
				"tenant":               tenant,
				"startDate":            startDate,
				"endDate":              endDate,
				"customerRelationship": neo4jenum.Customer.String(),
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
					
					OPTIONAL MATCH (t:Tenant {name: $tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WITH year, month, beginOfMonth, startOfNextMonth, c.id as id, c.serviceStartedAt as serviceStartedAt, c.endedAt AS contractEndedAt, sli ORDER BY sli.startedAt ASC

					WHERE o.hide = false AND o.relationship = $customerRelationship AND c.serviceStartedAt IS NOT NULL AND contractEndedAt >= beginOfMonth AND contractEndedAt < startOfNextMonth AND sli.startedAt IS NOT NULL AND (sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
					return year, month, COUNT(DISTINCT id) AS value
				`, "% 12 + 1", tenant, tenant, tenant),
			map[string]any{
				"tenant":               tenant,
				"startDate":            startDate,
				"endDate":              endDate,
				"customerRelationship": neo4jenum.Customer.String(),
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

func (r *dashboardRepository) GetDashboardAverageTimeToOnboardPerMonth(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardAverageTimeToOnboardPerMonth")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("startDate", startDate), log.Object("endDate", endDate))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := fmt.Sprintf(`
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
						 currentDate + duration({MONTHS: 1}) - duration({NANOSECONDS: 1}) AS endOfMonth
					WITH DISTINCT currentDate.YEAR AS year, currentDate.MONTH AS month, beginOfMonth, endOfMonth
						MATCH (o:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
						OPTIONAL MATCH (o)<-[:ACTION_ON]-(action:Action {type:"ONBOARDING_STATUS_CHANGED"})
							WHERE action.status IN ['DONE','SUCCESSFUL'] 
							AND action.createdAt >= beginOfMonth
  							AND action.createdAt <= endOfMonth 
					WITH year, month, o, action
						OPTIONAL MATCH (o)<-[:ACTION_ON]-(previousDone:Action {type:"ONBOARDING_STATUS_CHANGED"})
							WHERE previousDone.createdAt < action.createdAt
  							AND previousDone.status IN ['DONE','SUCCESSFUL']
					WITH year, month, o, action, max(previousDone.createdAt) as previousDoneCreatedAt
						OPTIONAL MATCH (o)<-[:ACTION_ON]-(startAction:Action {type:"ONBOARDING_STATUS_CHANGED"})
							WHERE startAction.createdAt < action.createdAt
							AND NOT startAction.status IN ['DONE','SUCCESSFUL']
							AND (previousDoneCreatedAt IS NULL OR startAction.createdAt>previousDoneCreatedAt)
					WITH year, month, action.createdAt as endDate, coalesce(min(startAction.createdAt), action.createdAt) as startDate
					RETURN year, month, avg(duration.inSeconds(startDate, endDate)) AS durationInSeconds ORDER BY year, month
				`, "% 12 + 1")
	params := map[string]any{
		"tenant":    tenant,
		"startDate": startDate,
		"endDate":   endDate,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
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
			record := map[string]interface{}{
				"year":  v.Values[0].(int64),
				"month": v.Values[1].(int64),
			}
			if v.Values[2] != nil {
				record["duration"] = v.Values[2].(neo4j.Duration)
			}
			results = append(results, record)
		}
	}

	return results, nil
}

func (r *dashboardRepository) GetDashboardOnboardingCompletionPerMonth(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardRepository.GetDashboardOnboardingCompletionPerMonth")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("startDate", startDate), log.Object("endDate", endDate))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := fmt.Sprintf(`
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
						 currentDate + duration({MONTHS: 1}) - duration({NANOSECONDS: 1}) AS endOfMonth
					WITH DISTINCT currentDate.YEAR AS year, currentDate.MONTH AS month, beginOfMonth, endOfMonth
						MATCH (o:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
						OPTIONAL MATCH (o)<-[:ACTION_ON]-(action:Action {type:"ONBOARDING_STATUS_CHANGED"})
							WHERE action.status IN ['DONE','SUCCESSFUL'] 
  							AND action.createdAt >= beginOfMonth
  							AND action.createdAt <= endOfMonth 
					WITH year, month, endOfMonth, o, action
						OPTIONAL MATCH (o)<-[:ACTION_ON]-(previousAction:Action {type:"ONBOARDING_STATUS_CHANGED"})
							WHERE previousAction.createdAt < action.createdAt 
					WITH year, month, endOfMonth, o, action, previousAction
						ORDER BY previousAction.createdAt DESC
					WITH year, month, endOfMonth, o, action, head(collect(previousAction)) AS previousAction  
						OPTIONAL MATCH (o)<-[:ACTION_ON]-(lastActionInPeriod:Action {type:"ONBOARDING_STATUS_CHANGED"})
							WHERE lastActionInPeriod.createdAt < endOfMonth
					WITH year, month, o, action, previousAction, lastActionInPeriod
						ORDER BY lastActionInPeriod.createdAt DESC
					WITH year, month, action, previousAction, head(collect(lastActionInPeriod)) AS lastAction
					WITH year, month, 
     					CASE WHEN action IS NOT NULL AND (previousAction IS NULL OR NOT previousAction.status IN ['DONE', 'SUCCESSFUL']) THEN action.id ELSE null END AS actionId,
     					CASE WHEN lastAction IS NOT NULL AND NOT lastAction.status IN ['DONE', 'SUCCESSFUL'] THEN lastAction.id ELSE null END AS lastActionId
					RETURN year, month, count(distinct(actionId)) AS completedOnboardings, count(distinct(lastActionId)) AS notCompletedOnboardings  ORDER BY year, month
				`, "% 12 + 1")
	params := map[string]any{
		"tenant":    tenant,
		"startDate": startDate,
		"endDate":   endDate,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
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
			record := map[string]interface{}{
				"year":                    v.Values[0].(int64),
				"month":                   v.Values[1].(int64),
				"completedOnboardings":    v.Values[2].(int64),
				"notCompletedOnboardings": v.Values[3].(int64),
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

func (r *dashboardRepository) GetDashboardGRRData(ctx context.Context, tenant string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
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
					WITH startDate.year AS startYear, startDate.month AS startMonth, endDate.year AS endYear, endDate.month AS endMonth
					WITH RANGE(startYear * 12 + startMonth - 1, endYear * 12 + endMonth - 1) AS monthsRange
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
						 currentDate + duration({months: 1}) AS startOfNextMonth
					WITH currentDate.year AS year, currentDate.month AS month, beginOfMonth, endOfMonth, startOfNextMonth
										OPTIONAL MATCH (t:Tenant {name: $tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract_%s)-[:HAS_SERVICE]->(sli:ServiceLineItem_%s)
					WITH year, month, beginOfMonth, endOfMonth, startOfNextMonth, c.id as cid, c.endedAt as cEndedAt, sli.parentId as sp, sli
					ORDER BY sli.startedAt ASC
					WHERE c.serviceStartedAt IS NOT NULL
					  AND sli.startedAt < startOfNextMonth
					  AND (sli.billed = 'MONTHLY' OR sli.billed = 'QUARTERLY' OR sli.billed = 'ANNUALLY')
					WITH year, month, beginOfMonth, endOfMonth, startOfNextMonth, cid, cEndedAt,
						 COLLECT({
							 sliId: sli.id,
							 sliV1: CASE WHEN sli.id = sli.parentId THEN TRUE ELSE FALSE END,
							 sliStartedAt: sli.startedAt,
							 sliEndedAt: sli.endedAt,
							 amount: CASE 
											   WHEN sli.billed = 'MONTHLY' THEN sli.price * sli.quantity 
											   ELSE CASE 
													  WHEN sli.billed = 'QUARTERLY' THEN  sli.price * sli.quantity / 3 
													  ELSE CASE 
															 WHEN sli.billed = 'ANNUALLY' THEN sli.price * sli.quantity / 12 
															 ELSE 0 
														   END 
													END 
											 END
						 }) as sliPerContract
					WITH year, month, beginOfMonth, endOfMonth, cid, cEndedAt, startOfNextMonth,
						 REDUCE(s = [], sliItem in sliPerContract | [
											CASE WHEN sliItem.sliV1 = true AND sliItem.sliStartedAt < beginOfMonth - duration({months: 10}) THEN sliItem.amount ELSE 0 END,
											CASE WHEN sliItem.sliV1 = true AND sliItem.sliStartedAt < beginOfMonth - duration({months: 9}) THEN sliItem.amount ELSE 0 END,
											CASE WHEN sliItem.sliV1 = true AND sliItem.sliStartedAt < beginOfMonth - duration({months: 8}) THEN sliItem.amount ELSE 0 END,
											CASE WHEN sliItem.sliV1 = true AND sliItem.sliStartedAt < beginOfMonth - duration({months: 7}) THEN sliItem.amount ELSE 0 END,
											CASE WHEN sliItem.sliV1 = true AND sliItem.sliStartedAt < beginOfMonth - duration({months: 6}) THEN sliItem.amount ELSE 0 END,
											CASE WHEN sliItem.sliV1 = true AND sliItem.sliStartedAt < beginOfMonth - duration({months: 5}) THEN sliItem.amount ELSE 0 END,
											CASE WHEN sliItem.sliV1 = true AND sliItem.sliStartedAt < beginOfMonth - duration({months: 4}) THEN sliItem.amount ELSE 0 END,
											CASE WHEN sliItem.sliV1 = true AND sliItem.sliStartedAt < beginOfMonth - duration({months: 3}) THEN sliItem.amount ELSE 0 END,
											CASE WHEN sliItem.sliV1 = true AND sliItem.sliStartedAt < beginOfMonth - duration({months: 2}) THEN sliItem.amount ELSE 0 END,
											CASE WHEN sliItem.sliV1 = true AND sliItem.sliStartedAt < beginOfMonth - duration({months: 1}) THEN sliItem.amount ELSE 0 END,
											CASE WHEN sliItem.sliV1 = true AND sliItem.sliStartedAt < beginOfMonth - duration({months: 0}) THEN sliItem.amount ELSE 0 END,
											CASE WHEN sliItem.sliV1 = true AND sliItem.sliStartedAt < beginOfMonth + duration({months: 1}) - duration({nanoseconds: 1}) THEN sliItem.amount ELSE 0 END
						 ]) as sliActivePerContract,
						 REDUCE(s = [], sliItem in sliPerContract | 
							 CASE WHEN sliItem.sliStartedAt < startOfNextMonth  AND (sliItem.sliEndedAt IS NULL or sliItem.sliEndedAt > endOfMonth) AND (cEndedAt IS NULL OR cEndedAt >= beginOfMonth ) THEN s + sliItem ELSE s END
						 ) as sliActiveInCurrentMonth
					WITH year, month, beginOfMonth, endOfMonth, cid, sliActivePerContract, sliActiveInCurrentMonth, 
					REDUCE(acc = null, i IN RANGE(0, 11) |
						CASE WHEN sliActivePerContract[i] > 0 AND acc IS NULL THEN sliActivePerContract[i] ELSE acc END
					) AS baseline, 
					REDUCE(acc = 0, i IN sliActiveInCurrentMonth | acc + i.amount ) AS currentAmount
					WITH year, month, cid, baseline, currentAmount
					with year, month, SUM(baseline) as baselineValue, SUM(currentAmount) as currentValue
					return year, month, CASE WHEN baselineValue <> 0 THEN (currentValue / baselineValue) * 100 ELSE 0 END AS grossRevenueRetentionRate
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
