package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CommonWriteRepository interface {
	LinkEntityWithEntity(ctx context.Context, tenant, entityId string, entityType model.EntityType, relationship string, withEntityId string, withEntityType model.EntityType) error
}

type commonWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCommonWriteRepository(driver *neo4j.DriverWithContext, database string) CommonWriteRepository {
	return &commonWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *commonWriteRepository) LinkEntityWithEntity(ctx context.Context, tenant, entityId string, entityType model.EntityType, relationship string, withEntityId string, withEntityType model.EntityType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.LinkEntityWithEntity")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := fmt.Sprintf(`MATCH (parent:%s_%s {id:$entityId}) `, entityType.Neo4jLabel(), tenant)
	cypher += fmt.Sprintf(`MATCH (child:%s_%s {id:$withEntityId}) `, withEntityType.Neo4jLabel(), tenant)
	cypher += fmt.Sprintf(`MERGE (parent)-[r:%s]->(child)`, relationship)

	params := map[string]any{
		"tenant":       tenant,
		"entityId":     entityId,
		"withEntityId": withEntityId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if _, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return nil, nil
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}
