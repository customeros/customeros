package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

type CountryReadRepository interface {
	GetDefaultCountryCodeA3(ctx context.Context, tenant string) (string, error)
	GetCountryByCodeIfExists(ctx context.Context, code string) (*dbtype.Node, error)
	GetCountryByCodeA3IfExists(ctx context.Context, codeA3 string) (*dbtype.Node, error)
	GetCountryByCodeA2IfExists(ctx context.Context, codeA2 string) (*dbtype.Node, error)
	GetCountriesPaginated(ctx context.Context, skip, limit int) (*utils.DbNodesWithTotalCount, error)
	GetCountries(ctx context.Context) ([]*dbtype.Node, error)
	GetAllForPhoneNumbers(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
}

type countryReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCountryReadRepository(driver *neo4j.DriverWithContext, database string) CountryReadRepository {
	return &countryReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *countryReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *countryReadRepository) GetDefaultCountryCodeA3(ctx context.Context, tenant string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryReadRepository.GetDefaultCountryCodeA3")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (t:Tenant {name:$tenant})
				OPTIONAL MATCH (tenant)-[:DEFAULT_COUNTRY]->(dc:Country)
				RETURN COALESCE(dc.codeA3, "") AS countryCodeA3`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsString(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	span.LogFields(log.String("result", result.(string)))
	return result.(string), nil
}

func (r *countryReadRepository) GetCountryByCodeIfExists(ctx context.Context, code string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryReadRepository.GetCountryByCodeIfExists")
	defer span.Finish()
	span.LogFields(log.String("code", code))

	if len(strings.TrimSpace(code)) == 2 {
		return r.GetCountryByCodeA2IfExists(ctx, strings.TrimSpace(code))
	} else if len(strings.TrimSpace(code)) == 3 {
		return r.GetCountryByCodeA3IfExists(ctx, strings.TrimSpace(code))
	} else {
		return nil, nil
	}
}

func (r *countryReadRepository) GetCountryByCodeA3IfExists(ctx context.Context, codeA3 string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryReadRepository.GetCountryByCodeA3IfExists")
	defer span.Finish()
	span.LogFields(log.String("codeA3", codeA3))

	cypher := `MATCH (c:Country {codeA3:$codeA3}) WHERE $codeA3 <> "" RETURN c`
	params := map[string]any{
		"codeA3": codeA3,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if len(dbRecords.([]*dbtype.Node)) == 0 {
		span.LogFields(log.Bool("result.found", false))
		return nil, err
	}
	span.LogFields(log.Bool("result.found", true))
	return dbRecords.([]*dbtype.Node)[0], err
}

func (r *countryReadRepository) GetCountryByCodeA2IfExists(ctx context.Context, codeA2 string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryReadRepository.GetCountryByCodeA2IfExists")
	defer span.Finish()
	span.LogFields(log.String("codeA2", codeA2))

	cypher := `MATCH (c:Country {codeA2:$codeA2}) WHERE $codeA2 <> "" RETURN c`
	params := map[string]any{
		"codeA2": codeA2,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if len(dbRecords.([]*dbtype.Node)) == 0 {
		span.LogFields(log.Bool("result.found", false))
		return nil, err
	}
	span.LogFields(log.Bool("result.found", true))
	return dbRecords.([]*dbtype.Node)[0], err
}

func (r *countryReadRepository) GetCountriesPaginated(ctx context.Context, skip, limit int) (*utils.DbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryReadRepository.GetCountryByCodeA3")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "skip", skip)
	tracing.LogObjectAsJson(span, "limit", limit)

	cypherCount := `MATCH (c:Country) RETURN count(c) as count`
	paramsCount := map[string]any{}
	span.LogFields(log.String("cypherCount", cypherCount))
	tracing.LogObjectAsJson(span, "paramsCount", paramsCount)

	cypher := `MATCH (c:Country) RETURN c ORDER BY c.name SKIP $skip LIMIT $limit`
	params := map[string]any{
		"skip":  skip,
		"limit": limit,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypherCount, paramsCount)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(ctx, cypher, params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *countryReadRepository) GetCountries(ctx context.Context) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryReadRepository.GetCountries")
	defer span.Finish()

	cypher := `MATCH (c:Country) RETURN c ORDER BY c.name`
	span.LogFields(log.String("cypher", cypher))

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, map[string]any{})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})

	if err != nil {
		return nil, err
	}
	dbNodes := []*dbtype.Node{}
	for _, v := range dbRecords.([]*neo4j.Record) {
		if v.Values[0] != nil {
			dbNodes = append(dbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
		}
	}
	return dbNodes, err
}

func (r *countryReadRepository) GetAllForPhoneNumbers(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryReadRepository.GetAllForPhoneNumbers")
	defer span.Finish()
	span.LogFields(log.Object("ids", ids))

	query := `MATCH (:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber)-[:LINKED_TO]->(c:Country)
		 		WHERE p.id IN $ids 
		 		RETURN c, p.id`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}
