package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/customer-os-neo4j-repository/model"
	"github.com/openline-ai/customer-os-neo4j-repository/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type InteractionSessionCreateFields struct {
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
	SourceFields model.Source `json:"sourceFields"`
	Channel      string       `json:"channel"`
	ChannelData  string       `json:"channelData"`
	Identifier   string       `json:"identifier"`
	Type         string       `json:"type"`
	Status       string       `json:"status"`
	Name         string       `json:"name"`
}

type InteractionSessionWriteRepository interface {
	Create(ctx context.Context, tenant, interactionSessionId string, data InteractionSessionCreateFields) error
}

type interactionSessionWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInteractionSessionWriteRepository(driver *neo4j.DriverWithContext, database string) InteractionSessionWriteRepository {
	return &interactionSessionWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *interactionSessionWriteRepository) Create(ctx context.Context, tenant, interactionSessionId string, data InteractionSessionCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionWriteRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionSessionId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MERGE (i:InteractionSession:InteractionSession_%s {id:$interactionSessionId}) 
							ON CREATE SET 
								i.createdAt=$createdAt,
								i.updatedAt=$updatedAt,
								i.source=$source,
								i.sourceOfTruth=$sourceOfTruth,
								i.appSource=$appSource,
								i.status=$status,
								i.channel=$channel,
								i.channelData=$channelData,
								i.identifier=$identifier,
								i.type=$type,
								i.name=$name
							`, tenant)
	params := map[string]any{
		"tenant":               tenant,
		"interactionSessionId": interactionSessionId,
		"createdAt":            data.CreatedAt,
		"updatedAt":            data.UpdatedAt,
		"source":               data.SourceFields.Source,
		"sourceOfTruth":        data.SourceFields.Source,
		"appSource":            data.SourceFields.AppSource,
		"channel":              data.Channel,
		"channelData":          data.ChannelData,
		"identifier":           data.Identifier,
		"type":                 data.Type,
		"status":               data.Status,
		"name":                 data.Name,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *interactionSessionWriteRepository) executeWriteQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteWriteQueryOnDb(ctx, *r.driver, r.database, query, params)
}
