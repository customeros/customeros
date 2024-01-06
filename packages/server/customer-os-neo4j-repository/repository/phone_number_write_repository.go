package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"time"
)

type PhoneNumberCreateFields struct {
	RawPhoneNumber string       `json:"rawPhoneNumber"`
	SourceFields   model.Source `json:"sourceFields"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}

type PhoneNumberWriteRepository interface {
	CreatePhoneNumber(ctx context.Context, tenant, phoneNumberId string, data PhoneNumberCreateFields) error
	UpdatePhoneNumber(ctx context.Context, tenant, phoneNumberId, source string, updatedAt time.Time) error
}

type phoneNumberWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewPhoneNumberWriteRepository(driver *neo4j.DriverWithContext, database string) PhoneNumberWriteRepository {
	return &phoneNumberWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *phoneNumberWriteRepository) CreatePhoneNumber(ctx context.Context, tenant, phoneNumberId string, data PhoneNumberCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberWriteRepository.CreatePhoneNumber")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber:PhoneNumber_%s {id:$id}) 
		 ON CREATE SET p.rawPhoneNumber = $rawPhoneNumber, 
						p.validated = null,
						p.source = $source,
						p.sourceOfTruth = $sourceOfTruth,
						p.appSource = $appSource,
						p.createdAt = $createdAt,
						p.updatedAt = $updatedAt,
						p.syncedWithEventStore = true 
		 ON MATCH SET 	p.syncedWithEventStore = true`, tenant)
	params := map[string]any{
		"id":             phoneNumberId,
		"rawPhoneNumber": data.RawPhoneNumber,
		"tenant":         tenant,
		"source":         data.SourceFields.Source,
		"sourceOfTruth":  data.SourceFields.SourceOfTruth,
		"appSource":      data.SourceFields.AppSource,
		"createdAt":      data.CreatedAt,
		"updatedAt":      data.UpdatedAt,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *phoneNumberWriteRepository) UpdatePhoneNumber(ctx context.Context, tenant, phoneNumberId, source string, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberWriteRepository.UpdatePhoneNumber")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$id})
				WHERE p:PhoneNumber_%s
		 SET 	p.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE p.sourceOfTruth END,
				p.updatedAt = $updatedAt,
				p.syncedWithEventStore = true`, tenant)
	params := map[string]any{
		"id":            phoneNumberId,
		"tenant":        tenant,
		"sourceOfTruth": source,
		"updatedAt":     updatedAt,
		"overwrite":     source == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
