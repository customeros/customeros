package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

const (
	SearchParamStage                 = "STAGE"
	searchParamIndustry              = "INDUSTRY"
	searchParamEmployee              = "EMPLOYEE_COUNT"
	searchParamCountryA2             = "COUNTRY_A2"
	searchParamTags                  = "TAGS"
	searchParamLinkedInFollowerCount = "LINKEDIN_FOLLOWER_COUNT"
	searchParamIsPublic              = "ORGANIZATIONS_IS_PUBLIC"
	searchParamYearFounded           = "YEAR_FOUNDED"
)

var searchParamsMap = map[string]string{
	"STAGE":                                 SearchParamStage,
	"ORGANIZATIONS_STAGE":                   SearchParamStage,
	"INDUSTRY":                              searchParamIndustry,
	"ORGANIZATIONS_INDUSTRY":                searchParamIndustry,
	"EMPLOYEE_COUNT":                        searchParamEmployee,
	"ORGANIZATIONS_EMPLOYEE_COUNT":          searchParamEmployee,
	"COUNTRY_A2":                            searchParamCountryA2,
	"ORGANIZATIONS_HEADQUARTERS":            searchParamCountryA2,
	"TAGS":                                  searchParamTags,
	"ORGANIZATIONS_TAGS":                    searchParamTags,
	"LINKEDIN_FOLLOWER_COUNT":               searchParamLinkedInFollowerCount,
	"ORGANIZATIONS_LINKEDIN_FOLLOWER_COUNT": searchParamLinkedInFollowerCount,
	"IS_PUBLIC":                             searchParamIsPublic,
	"ORGANIZATIONS_IS_PUBLIC":               searchParamIsPublic,
	"YEAR_FOUNDED":                          searchParamYearFounded,
	"ORGANIZATIONS_YEAR_FOUNDED":            searchParamYearFounded,
}

func getSearchParam(input string) string {
	if searchParam, ok := searchParamsMap[input]; ok {
		return searchParam
	}
	return ""
}

type OrganizationWithFiltersReadRepository interface {
	GetFilteredOrganizationIds(ctx context.Context, tenant string, filter *model.Filter) ([]string, error)
}

type organizationWithFiltersReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOrganizationWithFiltersReadRepository(driver *neo4j.DriverWithContext, database string) OrganizationWithFiltersReadRepository {
	return &organizationWithFiltersReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *organizationWithFiltersReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *organizationWithFiltersReadRepository) GetFilteredOrganizationIds(ctx context.Context, tenant string, filter *model.Filter) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetFilteredOrganizationIds")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	tracing.LogObjectAsJson(span, "filter", filter)

	params := map[string]any{
		"tenant": tenant,
	}
	organizationFilterCypher, organizationFilterParams := "", make(map[string]interface{})
	tagFilterCypher, tagFilterParams := "", make(map[string]interface{})
	locationFilterCypher, locationFilterParams := "", make(map[string]interface{})
	socialFilterCypher, socialFilterParams := "", make(map[string]interface{})

	if filter != nil {
		organizationFilter := new(utils.CypherFilter)
		organizationFilter.Negate = false
		organizationFilter.LogicalOperator = utils.AND
		organizationFilter.Filters = make([]*utils.CypherFilter, 0)

		tagFilter := new(utils.CypherFilter)
		tagFilter.Negate = false
		tagFilter.LogicalOperator = utils.AND
		tagFilter.Filters = make([]*utils.CypherFilter, 0)

		locationFilter := new(utils.CypherFilter)
		locationFilter.Negate = false
		locationFilter.LogicalOperator = utils.AND
		locationFilter.Filters = make([]*utils.CypherFilter, 0)

		socialFilter := new(utils.CypherFilter)
		socialFilter.Negate = false
		socialFilter.LogicalOperator = utils.AND
		socialFilter.Filters = make([]*utils.CypherFilter, 0)

		for _, filterPart := range filter.And {
			if getSearchParam(filterPart.Filter.Property) == SearchParamStage {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterEq(string(neo4jentity.OrganizationPropertyStage), *filterPart.Filter.Value.Str))
			} else if getSearchParam(filterPart.Filter.Property) == searchParamIndustry {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterIn(string(neo4jentity.OrganizationPropertyIndustry), *filterPart.Filter.Value.ArrayStr))
			} else if getSearchParam(filterPart.Filter.Property) == searchParamEmployee {
				if filterPart.Filter.Operation == model.ComparisonOperatorBetween {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter(string(neo4jentity.OrganizationPropertyEmployees), *filterPart.Filter.Value.ArrayInt, utils.BETWEEN))
				} else {
					// expecting only LTE / LT / GTE / GT
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter(string(neo4jentity.OrganizationPropertyEmployees), (*filterPart.Filter.Value.ArrayInt)[0], filterPart.Filter.Operation.GetOperator()))
				}
			} else if getSearchParam(filterPart.Filter.Property) == searchParamYearFounded {
				if filterPart.Filter.Operation == model.ComparisonOperatorBetween {
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter(string(neo4jentity.OrganizationPropertyYearFounded), *filterPart.Filter.Value.ArrayInt, utils.BETWEEN))
				} else {
					// expecting only LTE / LT / GTE / GT
					organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilter(string(neo4jentity.OrganizationPropertyYearFounded), (*filterPart.Filter.Value.ArrayInt)[0], filterPart.Filter.Operation.GetOperator()))
				}
			} else if getSearchParam(filterPart.Filter.Property) == searchParamTags {
				tagFilter.Filters = append(tagFilter.Filters, utils.CreateCypherFilterIn(string(neo4jentity.TagPropertyId), *filterPart.Filter.Value.ArrayStr))
			} else if getSearchParam(filterPart.Filter.Property) == searchParamCountryA2 {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateCypherFilterIn(string(neo4jentity.LocationPropertyCountryCodeA2), *filterPart.Filter.Value.ArrayStr))
			} else if getSearchParam(filterPart.Filter.Property) == searchParamIsPublic {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterEq(string(neo4jentity.OrganizationPropertyIsPublic), *filterPart.Filter.Value.Bool))
			} else if getSearchParam(filterPart.Filter.Property) == searchParamLinkedInFollowerCount {
				socialFilter.Filters = append(socialFilter.Filters, utils.CreateCypherFilter(string(neo4jentity.SocialPropertyUrl), "linkedin.", utils.CONTAINS))
				if filterPart.Filter.Operation == model.ComparisonOperatorBetween {
					socialFilter.Filters = append(socialFilter.Filters, utils.CreateCypherFilter(string(neo4jentity.SocialPropertyFollowersCount), *filterPart.Filter.Value.ArrayInt, utils.BETWEEN))
				} else {
					// expecting only LTE / LT / GTE / GT
					socialFilter.Filters = append(socialFilter.Filters, utils.CreateCypherFilter(string(neo4jentity.SocialPropertyFollowersCount), (*filterPart.Filter.Value.ArrayInt)[0], filterPart.Filter.Operation.GetOperator()))
				}
			}
		}

		if len(organizationFilter.Filters) > 0 {
			organizationFilterCypher, organizationFilterParams = organizationFilter.BuildCypherFilterFragmentWithParamName("o", "o_param_")
		}
		if len(tagFilter.Filters) > 0 {
			tagFilterCypher, tagFilterParams = tagFilter.BuildCypherFilterFragmentWithParamName("t", "t_param_")
		}
		if len(locationFilter.Filters) > 0 {
			locationFilterCypher, locationFilterParams = locationFilter.BuildCypherFilterFragmentWithParamName("l", "l_param_")
		}
		if len(socialFilter.Filters) > 0 {
			socialFilterCypher, socialFilterParams = socialFilter.BuildCypherFilterFragmentWithParamName("s", "s_param_")
		}
	}

	cypher := `MATCH (o:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			WHERE o.hide = false WITH * `
	if tagFilterCypher != "" {
		cypher += ` MATCH (o)-[:TAGGED]->(t:Tag) WITH *`
	}
	if locationFilterCypher != "" {
		cypher += ` MATCH (o)--(l:Location) WITH *`
	}
	if socialFilterCypher != "" {
		cypher += ` MATCH (o)-[:HAS]->(s:Social) WITH *`
	}
	if organizationFilterCypher != "" || tagFilterCypher != "" || locationFilterCypher != "" || socialFilterCypher != "" {
		cypher += " WHERE "
	}
	cypherParts := []string{}
	if organizationFilterCypher != "" {
		cypherParts = append(cypherParts, organizationFilterCypher)
	}
	if tagFilterCypher != "" {
		cypherParts = append(cypherParts, tagFilterCypher)
	}
	if locationFilterCypher != "" {
		cypherParts = append(cypherParts, locationFilterCypher)
	}
	if socialFilterCypher != "" {
		cypherParts = append(cypherParts, socialFilterCypher)
	}
	cypher = cypher + strings.Join(cypherParts, " AND ")
	cypher += " RETURN DISTINCT o.id"

	params = utils.MergeMaps(params, organizationFilterParams)
	params = utils.MergeMaps(params, tagFilterParams)
	params = utils.MergeMaps(params, locationFilterParams)
	params = utils.MergeMaps(params, socialFilterParams)

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		span.LogFields(log.Int("result.count", 0))
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]string))))
	return result.([]string), err
}
