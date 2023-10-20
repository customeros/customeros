package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LocationRepository interface {
	GetMatchedLocationIdForOrganizationBySource(ctx context.Context, organizationId, externalSystem string) (string, error)
	GetMatchedLocationIdForContactBySource(ctx context.Context, contactId, externalSystem string) (string, error)
	GetById(ctx context.Context, locationId string) (*dbtype.Node, error)
}

type locationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewLocationRepository(driver *neo4j.DriverWithContext) LocationRepository {
	return &locationRepository{
		driver: driver,
	}
}

func (r *locationRepository) GetMatchedLocationIdForOrganizationBySource(ctx context.Context, organizationId, externalSystem string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.GetMatchedLocationIdForOrganizationBySource")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId), log.String("externalSystem", externalSystem))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(:Organization {id:$organizationId})-[:ASSOCIATED_WITH]->(l:Location {source:$source})-[:LOCATION_BELONGS_TO_TENANT]->(t)
				RETURN l.id limit 1`
	params := map[string]interface{}{
		"tenant":         common.GetTenantFromContext(ctx),
		"source":         externalSystem,
		"organizationId": organizationId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	locationIds := dbRecords.([]*db.Record)
	if len(locationIds) > 0 {
		return locationIds[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *locationRepository) GetMatchedLocationIdForContactBySource(ctx context.Context, contactId, externalSystem string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.GetMatchedLocationIdForContactBySource")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("externalSystem", externalSystem))

	query := `MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(:Contact {id:$contactId})-[:ASSOCIATED_WITH]->(l:Location {source:$source})-[:LOCATION_BELONGS_TO_TENANT]->(t)
				RETURN l.id limit 1`
	params := map[string]interface{}{
		"tenant":    common.GetTenantFromContext(ctx),
		"source":    externalSystem,
		"contactId": contactId,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	locationIds := dbRecords.([]*db.Record)
	if len(locationIds) > 0 {
		return locationIds[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *locationRepository) GetById(ctx context.Context, locationId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("locationId", locationId))

	query := "MATCH (:Tenant {name:$tenant})<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location {id:$locationId}) RETURN l"
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"locationId": locationId,
				"tenant":     common.GetTenantFromContext(ctx),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}
