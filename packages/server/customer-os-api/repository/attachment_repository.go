package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type AttachmentRepository interface {
	LinkWithXXIncludesAttachmentInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, linkedWith LinkedWith, linkedNature *LinkedNature, attachmentId, includedById string) (*dbtype.Node, error)
	UnlinkWithXXIncludesAttachmentInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, linkedWith LinkedWith, linkedNature *LinkedNature, attachmentId, includedById string) (*dbtype.Node, error)
	GetAttachmentsForXX(ctx context.Context, tenant string, linkedWith LinkedWith, linkedNature *LinkedNature, ids []string) ([]*utils.DbNodeAndId, error)
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newAttachment entity.AttachmentEntity, source, sourceOfTruth neo4jentity.DataSource) (*dbtype.Node, error)
}

type attachmentRepository struct {
	driver *neo4j.DriverWithContext
}

func NewAttachmentRepository(driver *neo4j.DriverWithContext) AttachmentRepository {
	return &attachmentRepository{
		driver: driver,
	}
}

func (r *attachmentRepository) LinkWithXXIncludesAttachmentInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, linkedWith LinkedWith, linkedNature *LinkedNature, attachmentId, includedById string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentRepository.LinkWithXXIncludesAttachmentInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (i:%s_%s {id:$includedById}) `, linkedWith, tenant)
	query += fmt.Sprintf(`MATCH (a:Attachment_%s {id:$attachmentId}) `, tenant)
	if linkedNature != nil {
		query += `MERGE (i)-[r:INCLUDES {nature: $linkedNature}]->(a) `
	} else {
		query += `MERGE (i)-[r:INCLUDES]->(a) `
	}
	query += `return i `

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"includedById": includedById,
			"attachmentId": attachmentId,
			"linkedNature": linkedNature,
		})
	span.LogFields(log.String("query", query))
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *attachmentRepository) UnlinkWithXXIncludesAttachmentInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, linkedWith LinkedWith, linkedNature *LinkedNature, attachmentId, includedById string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentRepository.UnlinkWithXXIncludesAttachmentInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (i:%s_%s {id:$includedById})`, linkedWith, tenant)
	if linkedNature != nil {
		query += `-[r:INCLUDES {nature: $linkedNature}]->`
	} else {
		query += `-[r:INCLUDES]->`
	}

	query += fmt.Sprintf(`(a:Attachment_%s {id:$attachmentId}) `, tenant)
	query += ` DELETE r `
	query += ` return i `

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"includedById": includedById,
			"attachmentId": attachmentId,
			"linkedNature": linkedNature,
		})
	span.LogFields(log.String("query", query))
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *attachmentRepository) GetAttachmentsForXX(ctx context.Context, tenant string, linkedWith LinkedWith, linkedNature *LinkedNature, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentRepository.GetAttachmentsForXX")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	var query string
	if linkedNature == nil {
		query = "MATCH (n:%s_%s)-[r:INCLUDES]->(a:Attachment_%s)"
	} else {
		query = "MATCH (n:%s_%s)-[:INCLUDES {nature: $linkedNature}]->(a:Attachment_%s) "
	}
	query += " WHERE n.id IN $ids "

	if linkedNature == nil {
		query += " AND r.nature IS NULL "
	}
	query += " RETURN a, n.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, linkedWith, tenant, tenant),
			map[string]any{
				"tenant":       tenant,
				"ids":          ids,
				"linkedNature": linkedNature,
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

func (r *attachmentRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newAttachment entity.AttachmentEntity, source, sourceOfTruth neo4jentity.DataSource) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	var createdAt time.Time
	createdAt = utils.Now()
	if newAttachment.CreatedAt != nil {
		createdAt = *newAttachment.CreatedAt
	}

	query := "MERGE (a:Attachment_%s {id:randomUUID()}) ON CREATE SET " +
		" a:Attachment, " +
		" a.source=$source, " +
		" a.createdAt=$createdAt, " +
		" a.name=$name, " +
		" a.mimeType=$mimeType, " +
		" a.extension=$extension, " +
		" a.size=$size, " +
		" a.sourceOfTruth=$sourceOfTruth, " +
		" a.appSource=$appSource " +
		" RETURN a"

	span.LogFields(log.String("query", query))

	if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"source":        source,
			"createdAt":     createdAt,
			"name":          newAttachment.Name,
			"mimeType":      newAttachment.MimeType,
			"extension":     newAttachment.Extension,
			"size":          newAttachment.Size,
			"sourceOfTruth": sourceOfTruth,
			"appSource":     newAttachment.AppSource,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}
