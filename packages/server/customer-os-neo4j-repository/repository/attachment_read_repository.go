package neo4j_repository

import (
	"context"
	"fmt"

	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "AttachmentReadRepository.GetById")
	defer spans.Finish()

	spans.LogKV("id", id)

	cypher := fmt.Sprintf("MATCH (a:Attachment_%s {id:$id}) RETURN a", tenant)
	params := map[string]any{
		"id": id,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})

	if err != nil {
		spans.TraceError(err)
		spans.LogKV("result.found", false)
		return nil, err
	}

	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *attachmentReadRepository) GetFor(ctx context.Context, tenant string, entityType neo4jenum.EntityType, entityRelation *neo4jenum.EntityRelation, ids []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "AttachmentReadRepository.GetFor")
	defer spans.Finish()

	spans.LogKV("entityType", entityType.String(),
		"ids.count", len(ids))
	if entityRelation != nil {
		spans.LogKV("entityRelation", entityRelation.String())
	}
	spans.LogObjectAsJson("ids", ids)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
	})

	if err != nil {
		spans.TraceError(err)
		spans.LogKV("result.count", 0)
		return nil, err
	}

	resultArray := result.([]*utils.DbNodeAndId)
	spans.LogKV("result.count", len(resultArray))
	return resultArray, nil
}
