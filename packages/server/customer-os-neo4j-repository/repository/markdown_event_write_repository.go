package neo4j_repository

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type MarkdownEventWriteRepository interface {
	CreateInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, markdownEventId string, data data_fields.MarkdownEventFields) error
}

type markdownEventWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewMarkdownEventWriteRepository(driver *neo4j.DriverWithContext, database string) MarkdownEventWriteRepository {
	return &markdownEventWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *markdownEventWriteRepository) CreateInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, markdownEventId string, data data_fields.MarkdownEventFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MarkdownEventWriteRepository.CreateInTx")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, markdownEventId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$orgId})
							MERGE (m:MarkdownEvent {id:$markdownEventId})<-[:HAS_MARKDOWN_EVENT]-(o)
							ON CREATE SET 
								m:MarkdownEvent_%s,
								m:TimelineEvent,
								m:TimelineEvent_%s,
								m.createdAt=$createdAt,
								m.updatedAt=datetime(),
								m.source=$source,
								m.appSource=$appSource,
								m.content=$content,
								o.updatedAt=datetime()
							`, tenant, tenant)
	params := map[string]any{
		"tenant":          tenant,
		"markdownEventId": markdownEventId,
		"orgId":           utils.IfNotNilString(data.OrganizationId),
		"createdAt":       utils.IfNotNilTimeWithDefault(data.CreatedAt, utils.Now()),
		"appSource":       utils.IfNotNilString(data.AppSource),
		"content":         utils.IfNotNilString(data.Content),
	}
	if data.Source != nil {
		params["source"] = data.Source.String()
	} else {
		params["source"] = neo4j_entity.DataSourceOpenline.String()
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}
