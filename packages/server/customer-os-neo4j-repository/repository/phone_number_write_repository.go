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
}

type PhoneNumberValidateFields struct {
	E164          string    `json:"e164"`
	CountryCodeA2 string    `json:"countryCodeA2"`
	ValidatedAt   time.Time `json:"validatedAt"`
	Source        string    `json:"source"`
	AppSource     string    `json:"appSource"`
}

type PhoneNumberWriteRepository interface {
	CreatePhoneNumber(ctx context.Context, tenant, phoneNumberId string, data PhoneNumberCreateFields) error
	UpdatePhoneNumber(ctx context.Context, tenant, phoneNumberId, rawPhoneNumber, source string) error
	FailPhoneNumberValidation(ctx context.Context, tenant, phoneNumberId, validationError string) error
	PhoneNumberValidated(ctx context.Context, tenant, phoneNumberId string, data PhoneNumberValidateFields) error
	LinkWithContact(ctx context.Context, tenant, contactId, phoneNumberId, label string, primary bool) error
	LinkWithOrganization(ctx context.Context, tenant, organizationId, phoneNumberId, label string, primary bool) error
	LinkWithUser(ctx context.Context, tenant, userId, phoneNumberId, label string, primary bool) error
	CleanPhoneNumberValidation(ctx context.Context, tenant, phoneNumberId string) error
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
						p.updatedAt = datetime(),
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
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *phoneNumberWriteRepository) UpdatePhoneNumber(ctx context.Context, tenant, phoneNumberId, rawPhoneNumber, source string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberWriteRepository.UpdatePhoneNumber")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$id})
				WHERE p:PhoneNumber_%s
		 SET 	p.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE p.sourceOfTruth END,
				p.updatedAt = datetime(),
				p.rawPhoneNumber = $rawPhoneNumber,
				p.syncedWithEventStore = true`, tenant)
	params := map[string]any{
		"id":             phoneNumberId,
		"tenant":         tenant,
		"sourceOfTruth":  source,
		"rawPhoneNumber": rawPhoneNumber,
		"overwrite":      source == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *phoneNumberWriteRepository) FailPhoneNumberValidation(ctx context.Context, tenant, phoneNumberId, validationError string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberWriteRepository.FailPhoneNumberValidation")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$id})
				WHERE p:PhoneNumber_%s
		 		SET p.validationError = $validationError,
		     		p.validated = false,
					p.updatedAt = datetime()`, tenant)
	params := map[string]any{
		"id":              phoneNumberId,
		"tenant":          tenant,
		"validationError": validationError,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *phoneNumberWriteRepository) PhoneNumberValidated(ctx context.Context, tenant, phoneNumberId string, data PhoneNumberValidateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberWriteRepository.PhoneNumberValidated")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$id})
				WHERE p:PhoneNumber_%s
		 		SET p.validationError = $validationError,
					p.e164 = $e164,
		     		p.validated = true,
					p.updatedAt = datetime()
				WITH p
				WHERE $countryCodeA2 <> ''
				WITH p
				CALL {
					WITH p
    				OPTIONAL MATCH (p)-[r:LINKED_TO]->(oldCountry:Country)
    				WHERE oldCountry.codeA2 <> $countryCodeA2
    				DELETE r
				}
				MERGE (c:Country {codeA2: $countryCodeA2})
					ON CREATE SET 	c.createdAt = $now, 
									c.updatedAt = datetime(), 
									c.appSource = $appSource,
									c.source = $source,
									c.sourceOfTruth = $source
				MERGE (p)-[:LINKED_TO]->(c)
				`, tenant)
	params := map[string]any{
		"id":              phoneNumberId,
		"tenant":          tenant,
		"validationError": "",
		"e164":            data.E164,
		"validatedAt":     data.ValidatedAt,
		"countryCodeA2":   data.CountryCodeA2,
		"now":             utils.Now(),
		"appSource":       data.AppSource,
		"source":          data.Source,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *phoneNumberWriteRepository) LinkWithContact(ctx context.Context, tenant, contactId, phoneNumberId, label string, primary bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberWriteRepository.LinkWithContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)

	cypher := `
		MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}),
				(t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId})
		MERGE (c)-[rel:HAS]->(p)
		SET	rel.primary = $primary,
			rel.label = $label,	
			c.updatedAt = datetime(),
			rel.syncedWithEventStore = true`
	params := map[string]any{
		"tenant":        tenant,
		"contactId":     contactId,
		"phoneNumberId": phoneNumberId,
		"label":         label,
		"primary":       primary,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *phoneNumberWriteRepository) LinkWithOrganization(ctx context.Context, tenant, organizationId, phoneNumberId, label string, primary bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberWriteRepository.LinkWithOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)

	cypher := `
		MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId}),
				(t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId})
		MERGE (org)-[rel:HAS]->(p)
		SET	rel.primary = $primary,
			rel.label = $label,	
			org.updatedAt = datetime(),
			rel.syncedWithEventStore = true`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"phoneNumberId":  phoneNumberId,
		"label":          label,
		"primary":        primary,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *phoneNumberWriteRepository) LinkWithUser(ctx context.Context, tenant, userId, phoneNumberId, label string, primary bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberWriteRepository.LinkWithUser")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)

	cypher := `
		MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId}),
				(t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId})
		MERGE (u)-[rel:HAS]->(p)
		SET	rel.primary = $primary,
			rel.label = $label,	
			u.updatedAt = datetime(),
			rel.syncedWithEventStore = true`
	params := map[string]any{
		"tenant":        tenant,
		"userId":        userId,
		"phoneNumberId": phoneNumberId,
		"label":         label,
		"primary":       primary,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *phoneNumberWriteRepository) CleanPhoneNumberValidation(ctx context.Context, tenant, phoneNumberId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberWriteRepository.CleanPhoneNumberValidation")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$id})
				WHERE p:PhoneNumber_%s
		 		SET p.validationError = null,
		     		p.validated = null,
					p.e164 = "",
					p.updatedAt = datetime()`, tenant)
	params := map[string]any{
		"id":     phoneNumberId,
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
