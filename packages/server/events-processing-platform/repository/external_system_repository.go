package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ExternalSystemRepository interface {
	LinkWithEntity(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, externalSystem cmnmod.ExternalSystem) error
}

type externalSystemRepository struct {
	driver *neo4j.DriverWithContext
}

func NewExternalSystemRepository(driver *neo4j.DriverWithContext) ExternalSystemRepository {
	return &externalSystemRepository{
		driver: driver,
	}
}

func (r *externalSystemRepository) LinkWithEntity(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, externalSystem cmnmod.ExternalSystem) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemRepository.LinkWithEntity")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.Object("externalSystem", externalSystem), log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel))

	query := fmt.Sprintf(`MATCH (n:%s {id:$entityId}),
			(t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystemId})
		MERGE (n)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext)
		ON CREATE SET
			r.syncDate=$syncDate, 
			r.externalUrl=$externalUrl, 
			r.externalSource=$externalSource
		ON MATCH SET
			r.syncDate=$syncDate, 
			r.externalSource=$externalSource`, linkedEntityNodeLabel+"_"+tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	return r.executeQuery(ctx, query, map[string]any{
		"tenant":           tenant,
		"externalSystemId": externalSystem.ExternalSystemId,
		"externalId":       externalSystem.ExternalId,
		"externalUrl":      externalSystem.ExternalUrl,
		"externalSource":   externalSystem.ExternalSource,
		"syncDate":         utils.TimePtrFirstNonNilNillableAsAny(externalSystem.SyncDate),
		"entityId":         linkedEntityId,
	})
}

// Common database interaction method
func (r *externalSystemRepository) executeQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteQuery(ctx, *r.driver, query, params)
}
