package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type ExternalSystemWriteRepository interface {
	LinkWithEntity(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, externalSystem model.ExternalSystem) error
	LinkWithEntityInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, linkedEntityId, linkedEntityNodeLabel string, externalSystem model.ExternalSystem) error
}

type externalSystemWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewExternalSystemWriteRepository(driver *neo4j.DriverWithContext, database string) ExternalSystemWriteRepository {
	return &externalSystemWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *externalSystemWriteRepository) prepareWriteSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *externalSystemWriteRepository) LinkWithEntity(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, externalSystem model.ExternalSystem) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemWriteRepository.LinkWithEntity")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel))
	tracing.LogObjectAsJson(span, "externalSystem", externalSystem)

	session := r.prepareWriteSession(ctx)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.LinkWithEntityInTx(ctx, tx, tenant, linkedEntityId, linkedEntityNodeLabel, externalSystem)
	})
	return err
}

func (r *externalSystemWriteRepository) LinkWithEntityInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, linkedEntityId, linkedEntityNodeLabel string, externalSystem model.ExternalSystem) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemWriteRepository.LinkWithEntityInTx")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel))
	tracing.LogObjectAsJson(span, "externalSystem", externalSystem)

	cypher := fmt.Sprintf(`MATCH (n:%s {id:$entityId}),
			(t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystemId})
		MERGE (n)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext)
		ON CREATE SET
			r.syncDate=$syncDate, 
			r.externalIdSecond=$externalIdSecond,
			r.externalUrl=$externalUrl, 
			r.externalSource=$externalSource
		ON MATCH SET
			r.syncDate=$syncDate, 
			r.externalSource=$externalSource`, linkedEntityNodeLabel+"_"+tenant)
	params := map[string]any{
		"tenant":           tenant,
		"externalSystemId": externalSystem.ExternalSystemId,
		"externalId":       externalSystem.ExternalId,
		"externalUrl":      externalSystem.ExternalUrl,
		"externalSource":   externalSystem.ExternalSource,
		"externalIdSecond": externalSystem.ExternalIdSecond,
		"syncDate":         utils.TimePtrFirstNonNilNillableAsAny(externalSystem.SyncDate),
		"entityId":         linkedEntityId,
	}
	span.LogFields(log.String("cypher", cypher))
	span.LogFields(log.Object("params", params))

	return utils.ExecuteQueryInTx(ctx, tx, cypher, params)
}
