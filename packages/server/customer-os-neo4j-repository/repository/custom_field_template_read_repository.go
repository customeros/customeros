package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type CustomFieldTemplateReadRepository interface {
	GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetById(ctx context.Context, tenant, id string) (*dbtype.Node, error)
}

type customFieldTemplateReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCustomFieldTemplateReadRepository(driver *neo4j.DriverWithContext, database string) CustomFieldTemplateReadRepository {
	return &customFieldTemplateReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *customFieldTemplateReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *customFieldTemplateReadRepository) GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateReadRepository.GetAllForTenant")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CUSTOM_FIELD_TEMPLATE_BELONGS_TO_TENANT]-(cft:CustomFieldTemplate) RETURN cft ORDER BY cft.entityType, cft.order`
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
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), err
}

func (r *customFieldTemplateReadRepository) GetById(ctx context.Context, tenant, id string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateReadRepository.GetById")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, id)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CUSTOM_FIELD_TEMPLATE_BELONGS_TO_TENANT]-(cft:CustomFieldTemplate {id:$id}) RETURN cft`
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Object("result", result))
	return result.(*dbtype.Node), err
}
