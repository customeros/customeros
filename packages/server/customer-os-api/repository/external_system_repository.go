package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
)

type ExternalSystemRepository interface {
	LinkNodeWithExternalSystemInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, nodeId, nodeType string, relationship neo4jentity.ExternalSystemEntity) error
}

type externalSystemRepository struct {
	driver *neo4j.DriverWithContext
}

func NewExternalSystemRepository(driver *neo4j.DriverWithContext) ExternalSystemRepository {
	return &externalSystemRepository{
		driver: driver,
	}
}

func (r *externalSystemRepository) LinkNodeWithExternalSystemInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, nodeId, nodeLabel string, externalSystem neo4jentity.ExternalSystemEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemRepository.LinkContactWithExternalSystemInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})," +
		" (n:%s {id:$nodeId}) " +
		" MERGE (n)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e) " +
		" ON CREATE SET e:%s, " +
		"				r.syncDate=$syncDate, " +
		"				r.externalUrl=$externalUrl, " +
		"				r.externalSource=$externalSource, " +
		"				e.createdAt=datetime({timezone: 'UTC'}) " +
		" ON MATCH SET r.syncDate=$syncDate " +
		" RETURN r"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, nodeLabel, "ExternalSystem_"+tenant),
		map[string]any{
			"nodeId":           nodeId,
			"tenant":           tenant,
			"syncDate":         *externalSystem.Relationship.SyncDate,
			"externalId":       externalSystem.Relationship.ExternalId,
			"externalSystemId": externalSystem.ExternalSystemId,
			"externalUrl":      externalSystem.Relationship.ExternalUrl,
			"externalSource":   externalSystem.Relationship.ExternalSource,
		})

	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}
