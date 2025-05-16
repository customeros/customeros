package repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type LocationRepository interface {
	// Deprecated
	GetMatchedLocationIdForOrganizationBySource(ctx context.Context, organizationId, externalSystem string) (string, error)
	// Deprecated
	GetMatchedLocationIdForContactBySource(ctx context.Context, contactId, externalSystem string) (string, error)
	// Deprecated
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "LocationRepository.GetMatchedLocationIdForOrganizationBySource")
	defer spans.Finish()
	spans.LogKV("organizationId", organizationId)
	spans.LogKV("externalSystem", externalSystem)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(:Organization {id:$organizationId})-[:ASSOCIATED_WITH]->(l:Location {source:$source})-[:LOCATION_BELONGS_TO_TENANT]->(t)
				RETURN l.id limit 1`
	params := map[string]interface{}{
		"tenant":         common.GetTenantFromContext(ctx),
		"source":         externalSystem,
		"organizationId": organizationId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return "", err
	}
	locationIds := dbRecords.([]*db.Record)
	if len(locationIds) > 0 {
		return locationIds[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *locationRepository) GetMatchedLocationIdForContactBySource(ctx context.Context, contactId, externalSystem string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "LocationRepository.GetMatchedLocationIdForContactBySource")
	defer spans.Finish()
	spans.LogKV("contactId", contactId)
	spans.LogKV("externalSystem", externalSystem)

	query := `MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(:Contact {id:$contactId})-[:ASSOCIATED_WITH]->(l:Location {source:$source})-[:LOCATION_BELONGS_TO_TENANT]->(t)
				RETURN l.id limit 1`
	params := map[string]interface{}{
		"tenant":    common.GetTenantFromContext(ctx),
		"source":    externalSystem,
		"contactId": contactId,
	}
	spans.LogKV("query", query)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return "", err
	}
	locationIds := dbRecords.([]*db.Record)
	if len(locationIds) > 0 {
		return locationIds[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *locationRepository) GetById(ctx context.Context, locationId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "LocationRepository.GetById")
	defer spans.Finish()
	spans.LogKV("locationId", locationId)

	query := "MATCH (:Tenant {name:$tenant})<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location {id:$locationId}) RETURN l"
	spans.LogKV("query", query)

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
		spans.TraceError(err)
		return nil, err
	}
	return result.(*dbtype.Node), nil
}
