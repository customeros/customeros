package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
	"log"
	"time"
)

type IncludesType string

const (
	INCUDED_BY_INTERACTION_SESSION IncludesType = "InteractionSession"
	INCLUDED_BY_INTERACTION_EVENT  IncludesType = "InteractionEvent"
	INCLUDED_BY_NOTE               IncludesType = "Note"
)

type AttachmentRepository interface {
	LinkWithXXIncludesAttachmentInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, includesType IncludesType, attachmentId, includedById string) (*dbtype.Node, error)
	GetAttachmentsForXX(ctx context.Context, tenant string, includesType IncludesType, ids []string) ([]*utils.DbNodeAndId, error)
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

func (r *attachmentRepository) LinkWithXXIncludesAttachmentInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, includesType IncludesType, attachmentId, includedById string) (*dbtype.Node, error) {

	query := fmt.Sprintf(`MATCH (i:%s_%s {id:$includedById}) `, includesType, tenant)
	query += fmt.Sprintf(`MATCH (a:Attachment_%s {id:$attachmentId}) `, tenant)
	query += `MERGE (i)-[r:INCLUDES]->(a) `
	query += `return i `

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"includedById": includedById,
			"attachmentId": attachmentId,
		})
	log.Printf("*************Result: %v", queryResult)
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *attachmentRepository) GetAttachmentsForXX(ctx context.Context, tenant string, includesType IncludesType, ids []string) ([]*utils.DbNodeAndId, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (n:%s_%s)-[DESCRIBES]->(a:Attachment_%s) " +
		" WHERE n.id IN $ids " +
		" RETURN a, n.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, includesType, tenant, tenant),
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
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
