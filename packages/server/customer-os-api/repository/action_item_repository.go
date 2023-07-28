package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ActionItemRepository interface {
	LinkWithInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, linkedWith LinkedWith, entityId, includedById string) (*dbtype.Node, error)
	UnlinkWithTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, linkedWith LinkedWith, entityId, includedById string) (*dbtype.Node, error)
	GetFor(ctx context.Context, tenant string, linkedWith LinkedWith, entityIds []string) ([]*utils.DbNodeAndId, error)
}

type actionItemRepository struct {
	driver *neo4j.DriverWithContext
}

func NewActionItemRepository(driver *neo4j.DriverWithContext) ActionItemRepository {
	return &actionItemRepository{
		driver: driver,
	}
}

func (r *actionItemRepository) LinkWithInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, linkedWith LinkedWith, entityId, includedById string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionItemRepository.LinkWithInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (i:%s_%s {id:$includedById}) `, linkedWith, tenant)
	query += fmt.Sprintf(`MATCH (a:ActionItem_%s {id:$entityId}) `, tenant)
	query += `MERGE (i)-[r:INCLUDES]->(a) `
	query += `return i `

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"includedById": includedById,
			"entityId":     entityId,
		})
	span.LogFields(log.String("query", query))
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *actionItemRepository) UnlinkWithTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, linkedWith LinkedWith, entityId, includedById string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionItemRepository.UnlinkWithTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (i:%s_%s {id:$includedById})`, linkedWith, tenant)
	query += `-[r:INCLUDES]->`

	query += fmt.Sprintf(`(a:ActionItem_%s {id:$entityId}) `, tenant)
	query += ` DELETE r `
	query += ` return i `

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"includedById": includedById,
			"entityId":     entityId,
		})
	span.LogFields(log.String("query", query))
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *actionItemRepository) GetFor(ctx context.Context, tenant string, linkedWith LinkedWith, entityIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentRepository.GetAttachmentsForXX")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	var query = "MATCH (n:%s_%s)-[r:INCLUDES]->(a:ActionItem_%s)"
	query += " WHERE n.id IN $entityIds "
	query += " RETURN a, n.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, linkedWith, tenant, tenant),
			map[string]any{
				"tenant":    tenant,
				"entityIds": entityIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	span.LogFields(log.String("query", query))
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}
