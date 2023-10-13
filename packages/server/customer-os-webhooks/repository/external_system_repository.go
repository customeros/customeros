package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ExternalSystemRepository interface {
	MergeExternalSystem(ctx context.Context, tenant, externalSystemId, externalSystemName string) error
}

type externalSystemRepository struct {
	driver *neo4j.DriverWithContext
}

func NewExternalSystemRepository(driver *neo4j.DriverWithContext) ExternalSystemRepository {
	return &externalSystemRepository{
		driver: driver,
	}
}

func (r *externalSystemRepository) MergeExternalSystem(ctx context.Context, tenant, externalSystemId, externalSystemName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemRepository.MergeExternalSystem")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("externalSystemId", externalSystemId), log.String("externalSystemName", externalSystemName))

	query := fmt.Sprintf(`MATCH(t:Tenant {name:$tenant})
							MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId}) 
							ON CREATE SET e.name=$externalSystemName, e.createdAt=$now, e.updatedAt=$now, e:ExternalSystem_%s`, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	return r.executeQuery(ctx, query, map[string]any{
		"tenant":             tenant,
		"externalSystemId":   externalSystemId,
		"externalSystemName": utils.FirstNotEmpty(externalSystemName, externalSystemId),
		"now":                utils.Now(),
	})
}

func (r *externalSystemRepository) executeQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteQuery(ctx, *r.driver, query, params)
}
