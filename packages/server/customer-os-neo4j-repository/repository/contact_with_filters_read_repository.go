package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

const (
	ContactSearchParamStage                 = "STAGE"
	contactSearchParamCity                  = "CITY"
	contactSearchParamCountryA2             = "COUNTRY_A2"
	contactSearchParamTags                  = "TAGS"
	contactSearchParamLinkedInFollowerCount = "LINKEDIN_FOLLOWER_COUNT"
)

var contactSearchParamsMap = map[string]string{
	ContactSearchParamStage:                 ContactSearchParamStage,
	contactSearchParamCountryA2:             contactSearchParamCountryA2,
	"CONTACTS_COUNTRY":                      contactSearchParamCountryA2,
	contactSearchParamCity:                  contactSearchParamCity,
	"CONTACTS_CITY":                         contactSearchParamCity,
	contactSearchParamTags:                  contactSearchParamTags,
	"CONTACTS_TAGS":                         contactSearchParamTags,
	contactSearchParamLinkedInFollowerCount: contactSearchParamLinkedInFollowerCount,
	"CONTACTS_LINKEDIN_FOLLOWER_COUNT":      contactSearchParamLinkedInFollowerCount,
}

func getContactSearchParam(input string) string {
	if searchParam, ok := contactSearchParamsMap[input]; ok {
		return searchParam
	}
	return ""
}

type ContactWithFiltersReadRepository interface {
	GetFilteredContactIds(ctx context.Context, tenant string, filter *model.Filter) ([]string, error)
}

type contactWithFiltersReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewContactWithFiltersReadRepository(driver *neo4j.DriverWithContext, database string) ContactWithFiltersReadRepository {
	return &contactWithFiltersReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *contactWithFiltersReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *contactWithFiltersReadRepository) GetFilteredContactIds(ctx context.Context, tenant string, filter *model.Filter) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWithFiltersReadRepository.GetFilteredContactIds")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.LogObjectAsJson(span, "filter", filter)

	params := map[string]any{
		"tenant": tenant,
	}
	contactFilterCypher, contactFilterParams := "", make(map[string]interface{})
	organizationFilterCypher, organizationFilterParams := "", make(map[string]interface{})
	tagFilterCypher, tagFilterParams := "", make(map[string]interface{})
	locationFilterCypher, locationFilterParams := "", make(map[string]interface{})
	socialFilterCypher, socialFilterParams := "", make(map[string]interface{})
	jobRoleFilterCypher, jobRoleFilterParams := "", make(map[string]interface{})

	if filter != nil {
		contactFilter := new(utils.CypherFilter)
		contactFilter.Negate = false
		contactFilter.LogicalOperator = utils.AND
		contactFilter.Filters = make([]*utils.CypherFilter, 0)

		organizationFilter := new(utils.CypherFilter)
		organizationFilter.Negate = false
		organizationFilter.LogicalOperator = utils.AND
		organizationFilter.Filters = make([]*utils.CypherFilter, 0)

		jobRoleFilter := new(utils.CypherFilter)
		jobRoleFilter.Negate = false
		jobRoleFilter.LogicalOperator = utils.AND
		jobRoleFilter.Filters = make([]*utils.CypherFilter, 0)

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
			if getContactSearchParam(filterPart.Filter.Property) == contactSearchParamTags {
				tagFilter.Filters = append(tagFilter.Filters, utils.CreateCypherFilterIn(string(neo4jentity.TagPropertyId), *filterPart.Filter.Value.ArrayStr))
			} else if getContactSearchParam(filterPart.Filter.Property) == contactSearchParamLinkedInFollowerCount {
				socialFilter.Filters = append(socialFilter.Filters, utils.CreateCypherFilter(string(neo4jentity.SocialPropertyUrl), "linkedin.", utils.CONTAINS))
				if filterPart.Filter.Operation == model.ComparisonOperatorBetween {
					socialFilter.Filters = append(socialFilter.Filters, utils.CreateCypherFilter(string(neo4jentity.SocialPropertyFollowersCount), *filterPart.Filter.Value.ArrayInt, utils.BETWEEN))
				} else {
					// expecting only LTE / LT / GTE / GT
					socialFilter.Filters = append(socialFilter.Filters, utils.CreateCypherFilter(string(neo4jentity.SocialPropertyFollowersCount), (*filterPart.Filter.Value.ArrayInt)[0], filterPart.Filter.Operation.GetOperator()))
				}
			} else if getContactSearchParam(filterPart.Filter.Property) == ContactSearchParamStage {
				organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateCypherFilterEq(string(neo4jentity.OrganizationPropertyStage), *filterPart.Filter.Value.Str))
			} else if getContactSearchParam(filterPart.Filter.Property) == contactSearchParamCountryA2 {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateCypherFilterIn(string(neo4jentity.LocationPropertyCountryCodeA2), *filterPart.Filter.Value.ArrayStr))
			} else if getContactSearchParam(filterPart.Filter.Property) == contactSearchParamCity {
				locationFilter.Filters = append(locationFilter.Filters, utils.CreateCypherFilterIn(string(neo4jentity.LocationPropertyLocality), *filterPart.Filter.Value.ArrayStr))
			}
		}

		if len(contactFilter.Filters) > 0 {
			contactFilterCypher, contactFilterParams = contactFilter.BuildCypherFilterFragmentWithParamName("c", "c_param_")
		}
		if len(organizationFilter.Filters) > 0 {
			organizationFilterCypher, organizationFilterParams = organizationFilter.BuildCypherFilterFragmentWithParamName("o", "o_param_")
		}
		if len(jobRoleFilter.Filters) > 0 {
			jobRoleFilterCypher, jobRoleFilterParams = jobRoleFilter.BuildCypherFilterFragmentWithParamName("j", "j_param_")
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

	cypher := `MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			WITH * `
	if organizationFilterCypher != "" || jobRoleFilterCypher != "" {
		cypher += ` MATCH (c)--(j:JobRole)--(o:Organization) WITH *`
	}
	if tagFilterCypher != "" {
		cypher += ` MATCH (c)-[:TAGGED]->(t:Tag) WITH *`
	}
	if locationFilterCypher != "" {
		cypher += ` MATCH (c)--(l:Location) WITH *`
	}
	if socialFilterCypher != "" {
		cypher += ` MATCH (c)-[:HAS]->(s:Social) WITH *`
	}
	if contactFilterCypher != "" || organizationFilterCypher != "" || jobRoleFilterCypher != "" || tagFilterCypher != "" || locationFilterCypher != "" || socialFilterCypher != "" {
		cypher += " WHERE "
	}
	cypherParts := []string{}
	if contactFilterCypher != "" {
		cypherParts = append(cypherParts, contactFilterCypher)
	}
	if organizationFilterCypher != "" {
		cypherParts = append(cypherParts, organizationFilterCypher)
	}
	if jobRoleFilterCypher != "" {
		cypherParts = append(cypherParts, jobRoleFilterCypher)
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
	cypher += " RETURN DISTINCT c.id"

	params = utils.MergeMaps(params, contactFilterParams)
	params = utils.MergeMaps(params, organizationFilterParams)
	params = utils.MergeMaps(params, jobRoleFilterParams)
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
