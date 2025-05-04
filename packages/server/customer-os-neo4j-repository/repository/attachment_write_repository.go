package neo4j_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "AttachmentWriteRepository.Create")
	defer spans.Finish()

	spans.LogKV("id", id,
		"fileName", fileName,
		"mimeType", mimeType,
		"size", size)

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

	spans.LogKV("cypher", fmt.Sprintf(cypher, tenant))
	spans.LogObjectAsJson("params", params)

	queryResult, err := tx.Run(ctx, fmt.Sprintf(cypher, tenant), params)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	result, err := utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	if err != nil {
		spans.TraceError(err)
		spans.LogKV("result.found", false)
		return nil, err
	}

	spans.LogKV("result.found", result != nil)
	return result, nil
}
