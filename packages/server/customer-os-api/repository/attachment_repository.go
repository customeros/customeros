package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
	"time"
)

type IncludesType string
type IncludesNature string

const (
	INCLUDED_BY_INTERACTION_SESSION IncludesType = "InteractionSession"
	INCLUDED_BY_INTERACTION_EVENT   IncludesType = "InteractionEvent"
	INCLUDED_BY_MEETING             IncludesType = "Meeting"
	INCLUDED_BY_NOTE                IncludesType = "Note"

	INCLUDE_NATURE_RECORDING IncludesNature = "Recording"
)

type AttachmentRepository interface {
	LinkWithXXIncludesAttachmentInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, includesType IncludesType, includesNature *IncludesNature, attachmentId, includedById string) (*dbtype.Node, error)
	UnlinkWithXXIncludesAttachmentInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, includesType IncludesType, includesNature *IncludesNature, attachmentId, includedById string) (*dbtype.Node, error)
	GetAttachmentsForXX(ctx context.Context, tenant string, includesType IncludesType, includesNature *IncludesNature, ids []string) ([]*utils.DbNodeAndId, error)
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newAttachment entity.AttachmentEntity, source, sourceOfTruth entity.DataSource) (*dbtype.Node, error)
}

type attachmentRepository struct {
	driver *neo4j.DriverWithContext
}

func NewAttachmentRepository(driver *neo4j.DriverWithContext) AttachmentRepository {
	return &attachmentRepository{
		driver: driver,
	}
}

func (r *attachmentRepository) LinkWithXXIncludesAttachmentInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, includesType IncludesType, includesNature *IncludesNature, attachmentId, includedById string) (*dbtype.Node, error) {

	query := fmt.Sprintf(`MATCH (i:%s_%s {id:$includedById}) `, includesType, tenant)
	query += fmt.Sprintf(`MATCH (a:Attachment_%s {id:$attachmentId}) `, tenant)
	if includesNature != nil {
		query += `MERGE (i)-[r:INCLUDES {nature: $includesNature}]->(a) `
	} else {
		query += `MERGE (i)-[r:INCLUDES]->(a) `
	}
	query += `return i `

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"includedById":   includedById,
			"attachmentId":   attachmentId,
			"includesNature": includesNature,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *attachmentRepository) UnlinkWithXXIncludesAttachmentInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, includesType IncludesType, includesNature *IncludesNature, attachmentId, includedById string) (*dbtype.Node, error) {

	query := fmt.Sprintf(`MATCH (i:%s_%s {id:$includedById})`, includesType, tenant)
	if includesNature != nil {
		query += `MERGE (i)-[r:INCLUDES {nature: $includesNature}]->(a) `
	} else {
		query += `-[r:INCLUDES]->`
	}

	query += fmt.Sprintf(`(a:Attachment_%s {id:$attachmentId}) `, tenant)
	query += ` DELETE r `
	query += ` return i `

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"includedById":   includedById,
			"attachmentId":   attachmentId,
			"includesNature": includesNature,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *attachmentRepository) GetAttachmentsForXX(ctx context.Context, tenant string, includesType IncludesType, includesNature *IncludesNature, ids []string) ([]*utils.DbNodeAndId, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	var query string
	query = "MATCH (n:%s_%s)-[:INCLUDES {nature: $includesNature}]->(a:Attachment_%s) "
	query += " WHERE n.id IN $ids "
	query += " RETURN a, n.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, includesType, tenant, tenant),
			map[string]any{
				"tenant":         tenant,
				"ids":            ids,
				"includesNature": includesNature,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *attachmentRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newAttachment entity.AttachmentEntity, source, sourceOfTruth entity.DataSource) (*dbtype.Node, error) {
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
