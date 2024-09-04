package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type AttachmentWriteRepository interface {
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant, id, cdnUrl, basePath, fileName, mimeType string, size int64, createdAt *time.Time, source, sourceOfTruth neo4jentity.DataSource, appSource string) (*dbtype.Node, error)
}

type attachmentWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewAttachmentWriteRepository(driver *neo4j.DriverWithContext, database string) AttachmentWriteRepository {
	return &attachmentWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *attachmentWriteRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant, id, cdnUrl, basePath, fileName, mimeType string, size int64, createdAt *time.Time, source, sourceOfTruth neo4jentity.DataSource, appSource string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentWriteRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	id = utils.NewUUIDIfEmpty(id)

	if createdAt == nil {
		createdAt = utils.NowPtr()
	}

	query := "MERGE (a:Attachment_%s {id:$id}) ON CREATE SET " +
		" a:Attachment, " +
		" a.source=$source, " +
		" a.createdAt=$createdAt, " +
		" a.cdnUrl=$cdnUrl, " +
		" a.basePath=$basePath, " +
		" a.fileName=$fileName, " +
		" a.mimeType=$mimeType, " +
		" a.size=$size, " +
		" a.sourceOfTruth=$sourceOfTruth, " +
		" a.appSource=$appSource " +
		" RETURN a"

	span.LogFields(log.String("query", query))

	if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"source":        source,
			"createdAt":     *createdAt,
			"id":            id,
			"cdnUrl":        cdnUrl,
			"basePath":      basePath,
			"fileName":      fileName,
			"mimeType":      mimeType,
			"size":          size,
			"sourceOfTruth": sourceOfTruth,
			"appSource":     appSource,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}
