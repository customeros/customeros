package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type ExternalSystemWriteRepository interface {
	CreateIfNotExists(ctx context.Context, tenant, externalSystemId, externalSystemName string) error
	LinkWithEntity(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, externalSystem model.ExternalSystem) error
	LinkWithEntityInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, linkedEntityId, linkedEntityNodeLabel string, externalSystem model.ExternalSystem) error
	SetProperty(ctx context.Context, tenant, externalSystemId, propertyName string, propertyValue any) error
	SetPrimaryExternalId(ctx context.Context, tenant, externalSystemId string, externalId, linkedEntityNodeLabel, linkedEntityId string) error
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

func (r *externalSystemWriteRepository) CreateIfNotExists(ctx context.Context, tenant, externalSystemId, externalSystemName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemWriteRepository.CreateIfNotExists")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("externalSystemId", externalSystemId), log.String("externalSystemName", externalSystemName))

	cypher := fmt.Sprintf(`MATCH(t:Tenant {name:$tenant})
							MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId}) 
							ON CREATE SET e.name=$externalSystemName, e.createdAt=$now, e.updatedAt=datetime(), e:ExternalSystem_%s`, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"externalSystemId":   externalSystemId,
		"externalSystemName": utils.FirstNotEmpty(externalSystemName, externalSystemId),
		"now":                utils.Now(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *externalSystemWriteRepository) LinkWithEntity(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, externalSystem model.ExternalSystem) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemWriteRepository.LinkWithEntity")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
		"syncDate":         utils.TimePtrAsAny(externalSystem.SyncDate),
		"entityId":         linkedEntityId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	return utils.ExecuteQueryInTx(ctx, tx, cypher, params)
}

func (r *externalSystemWriteRepository) SetProperty(ctx context.Context, tenant, externalSystemId, propertyName string, propertyValue any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemWriteRepository.SetProperty")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("externalSystemId", externalSystemId), log.String("propertyName", propertyName), log.Object("propertyValue", propertyValue))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
		SET e.%s=$propertyValue, e.updatedAt=datetime()`, propertyName)
	params := map[string]any{
		"tenant":           tenant,
		"externalSystemId": externalSystemId,
		"propertyValue":    propertyValue,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *externalSystemWriteRepository) SetPrimaryExternalId(ctx context.Context, tenant, externalSystemId string, externalId, linkedEntityNodeLabel, linkedEntityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemWriteRepository.SetPrimaryExternalId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(
		log.String("externalSystemId", externalSystemId),
		log.String("externalId", externalId),
		log.String("linkedEntityNodeLabel", linkedEntityNodeLabel),
		log.String("linkedEntityId", linkedEntityId))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
		MATCH (n:%s {id:$linkedEntityId})
		WITH n, e
		MERGE (n)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e)
		SET r.primary=true
		WITH n, e
		MATCH (n)-[r:IS_LINKED_WITH]->(e)
		WHERE r.externalId <> $externalId
		SET r.primary=false`, linkedEntityNodeLabel+"_"+tenant)
	params := map[string]any{
		"tenant":           tenant,
		"externalSystemId": externalSystemId,
		"externalId":       externalId,
		"linkedEntityId":   linkedEntityId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
