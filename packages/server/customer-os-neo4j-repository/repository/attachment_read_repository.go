package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type AttachmentReadRepository interface {
	GetById(ctx context.Context, tenant string, id string) (*neo4j.Node, error)
	GetFor(ctx context.Context, tenant string, entityType neo4jenum.EntityType, entityRelation *neo4jenum.EntityRelation, ids []string) ([]*utils.DbNodeAndId, error)
}

type attachmentReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewAttachmentReadRepository(driver *neo4j.DriverWithContext, database string) AttachmentReadRepository {
	return &attachmentReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *attachmentReadRepository) GetById(ctx context.Context, tenant, id string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentReadRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf("MATCH (a:Attachment_%s {id:$id}) RETURN a", tenant)
	params := map[string]any{
		"id": id,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})

	if err != nil {
		tracing.TraceErr(span, err)
		span.LogFields(log.Bool("result.found", false))
		return nil, err
	}

	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (r *attachmentReadRepository) GetFor(ctx context.Context, tenant string, entityType neo4jenum.EntityType, entityRelation *neo4jenum.EntityRelation, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentReadRepository.GetFor")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := fmt.Sprintf(`MATCH (n:%s_%s)-`, entityType.Neo4jLabel(), tenant)

	if entityRelation != nil {
		cypher += fmt.Sprintf("[:%s]->", entityRelation.String())
	} else {
		cypher += "[r]->"
	}

	cypher += fmt.Sprintf(`(a:Attachment_%s)
				WHERE n.id IN $ids
				RETURN a, n.id`, tenant)
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
	})

	if err != nil {
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	if result == nil {
		return nil, nil
	}
	return result.([]*utils.DbNodeAndId), nil
}
