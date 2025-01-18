package neo4j_repository

import (
	"context"
	"fmt"
	"time"

	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type InteractionSessionWriteRepository interface {
	CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionSessionId string, data neo4j_entity.InteractionSessionEntity) error
	MergeByIdentifierAndChannel(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, identifier string, syncDate time.Time, message commonmodel.SaveEmailMessage, sessionType commonenum.InteractionSessionType, channel commonenum.InteractionSessionChannel, source, appSource string) (string, error)
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

func (r *interactionSessionWriteRepository) CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionSessionId string, data neo4j_entity.InteractionSessionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionWriteRepository.CreateInTx")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionSessionId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MERGE (i:InteractionSession:InteractionSession_%s {id:$interactionSessionId}) 
							ON CREATE SET 
								i.createdAt=$createdAt,
								i.updatedAt=datetime(),
								i.source=$source,
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
		"createdAt":            utils.NowIfZero(data.CreatedAt),
		"source":               data.Source,
		"appSource":            data.AppSource,
		"channel":              data.Channel.String(),
		"channelData":          data.ChannelData,
		"identifier":           data.Identifier,
		"type":                 data.Type.String(),
		"status":               data.Status.String(),
		"name":                 data.Name,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := tx.Run(ctx, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *interactionSessionWriteRepository) MergeByIdentifierAndChannel(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, identifier string, syncDate time.Time, message commonmodel.SaveEmailMessage, sessionType commonenum.InteractionSessionType, channel commonenum.InteractionSessionChannel, source, appSource string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionWriteRepository.MergeByIdentifierAndChannel")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogKV("identifier", identifier, "channel", channel, "sessionType", sessionType, "source", source, "appSource", appSource)

	cypher := ""
	if identifier == "" {
		cypher += `MATCH (:Tenant {name:$tenant}) 
			 	CREATE (is:InteractionSession_%s {identifier:$identifier, channel:$channel})
				SET `
	} else {
		cypher += `MATCH (:Tenant {name:$tenant}) 
			 	MERGE (is:InteractionSession_%s {identifier:$identifier, channel:$channel}) 
			 	ON CREATE SET `
	}
	cypher += ` is:InteractionSession,
			is.id=randomUUID(),
			is.syncDate=$syncDate,
			is.createdAt=$createdAt,
			is.updatedAt=datetime(),
			is.name=$name,
			is.status=$status,
			is.type=$type,
			is.source=$source,
			is.appSource=$appSource
		WITH is
		RETURN is.id`

	cypher = fmt.Sprintf(cypher, tenant)

	params := map[string]interface{}{
		"tenant":     tenant,
		"source":     source,
		"appSource":  appSource,
		"identifier": identifier,
		"name":       message.Subject,
		"syncDate":   syncDate,
		"createdAt":  message.CreatedAt,
		"status":     commonenum.InteractionSessionStatusActive.String(),
		"type":       sessionType,
		"channel":    channel,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	queryResult, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		qr, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsString(ctx, qr, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return queryResult.(string), nil
}
