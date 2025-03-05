package neo4j_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type AttachmentWriteRepository interface {
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant, id, cdnUrl, basePath, fileName, mimeType string, size int64, createdAt *time.Time, source neo4jentity.DataSource, appSource string) (*dbtype.Node, error)
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

func (r *attachmentWriteRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant, id, cdnUrl, basePath, fileName, mimeType string, size int64, createdAt *time.Time, source neo4jentity.DataSource, appSource string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentWriteRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(
		log.String("id", id),
		log.String("fileName", fileName),
		log.String("mimeType", mimeType),
		log.Int64("size", size))

	id = utils.NewUUIDIfEmpty(id)

	if createdAt == nil {
		createdAt = utils.NowPtr()
	}

	cypher := "MERGE (a:Attachment_%s {id:$id}) ON CREATE SET " +
		" a:Attachment, " +
		" a.source=$source, " +
		" a.createdAt=$createdAt, " +
		" a.cdnUrl=$cdnUrl, " +
		" a.basePath=$basePath, " +
		" a.fileName=$fileName, " +
		" a.mimeType=$mimeType, " +
		" a.size=$size, " +
		" a.appSource=$appSource " +
		" RETURN a"

	params := map[string]interface{}{
		"tenant":    tenant,
		"source":    source,
		"createdAt": *createdAt,
		"id":        id,
		"cdnUrl":    cdnUrl,
		"basePath":  basePath,
		"fileName":  fileName,
		"mimeType":  mimeType,
		"size":      size,
		"appSource": appSource,
	}

	span.LogFields(log.String("cypher", fmt.Sprintf(cypher, tenant)))
	tracing.LogObjectAsJson(span, "params", params)

	queryResult, err := tx.Run(ctx, fmt.Sprintf(cypher, tenant), params)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	result, err := utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	if err != nil {
		tracing.TraceErr(span, err)
		span.LogFields(log.Bool("result.found", false))
		return nil, err
	}

	span.LogFields(log.Bool("result.found", result != nil))
	return result, nil
}
