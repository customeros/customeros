package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CommentRepository interface {
	GetById(ctx context.Context, commentId string) (*dbtype.Node, error)
	GetMatchedCommentId(ctx context.Context, externalSystem, externalId string) (string, error)
}

type commentRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCommentRepository(driver *neo4j.DriverWithContext, database string) CommentRepository {
	return &commentRepository{
		driver:   driver,
		database: database,
	}
}

func (r *commentRepository) GetById(ctx context.Context, commentId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("commentId", commentId))

	cypher := fmt.Sprintf(`MATCH (c:Comment_%s {id:$commentId}) RETURN c`, common.GetTenantFromContext(ctx))
	params := map[string]any{
		"commentId": commentId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}

func (r *commentRepository) GetMatchedCommentId(ctx context.Context, externalSystem, externalId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentRepository.GetMatchedCommentId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", externalSystem), log.String("externalId", externalId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (e)<-[:IS_LINKED_WITH {externalId:$commentExternalId}]-(c:Comment)
				WITH c WHERE c IS NOT null
				RETURN c.id ORDER BY c.createdAt limit 1`
	params := map[string]interface{}{
		"tenant":            common.GetTenantFromContext(ctx),
		"externalSystem":    externalSystem,
		"commentExternalId": externalId,
	}
	span.LogFields(log.String("cypher", cypher))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	noteIDs := dbRecords.([]*db.Record)
	if len(noteIDs) > 0 {
		return noteIDs[0].Values[0].(string), nil
	}
	return "", nil
}
