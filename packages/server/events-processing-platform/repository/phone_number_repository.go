package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type PhoneNumberRepository interface {
	GetIdIfExists(ctx context.Context, tenant, phoneNumber string) (string, error)
	CreatePhoneNumber(ctx context.Context, phoneNumberId string, event events.PhoneNumberCreateEvent) error
	UpdatePhoneNumber(ctx context.Context, phoneNumberId string, event events.PhoneNumberUpdatedEvent) error
	FailPhoneNumberValidation(ctx context.Context, phoneNumberId string, event events.PhoneNumberFailedValidationEvent) error
	PhoneNumberValidated(ctx context.Context, phoneNumberId string, event events.PhoneNumberValidatedEvent) error
	LinkWithContact(ctx context.Context, tenant, contactId, phoneNumberId, label string, primary bool, updatedAt time.Time) error
	LinkWithOrganization(ctx context.Context, tenant, organizationId, phoneNumberId, label string, primary bool, updatedAt time.Time) error
	LinkWithUser(ctx context.Context, tenant, userId, phoneNumberId, label string, primary bool, updatedAt time.Time) error
	GetCountryCodeA2ForPhoneNumber(ctx context.Context, tenant, phoneNumberId string) (string, error)
}

type phoneNumberRepository struct {
	driver *neo4j.DriverWithContext
}

func NewPhoneNumberRepository(driver *neo4j.DriverWithContext) PhoneNumberRepository {
	return &phoneNumberRepository{
		driver: driver,
	}
}

func (r *phoneNumberRepository) GetIdIfExists(ctx context.Context, tenant, phoneNumber string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.GetIdIfExists")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("phoneNumber", phoneNumber))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (p:PhoneNumber_%s) WHERE p.e164 = $phoneNumber OR p.rawPhoneNumber = $phoneNumber RETURN p.id LIMIT 1"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"phoneNumber": phoneNumber,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return "", err
	}
	if len(result.([]*db.Record)) == 0 {
		return "", nil
	}
	return result.([]*db.Record)[0].Values[0].(string), err
}

func (r *phoneNumberRepository) CreatePhoneNumber(ctx context.Context, phoneNumberId string, event events.PhoneNumberCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.CreatePhoneNumber")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("phoneNumberId", phoneNumberId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber:PhoneNumber_%s {id:$id}) 
		 ON CREATE SET p.rawPhoneNumber = $rawPhoneNumber, 
						p.validated = null,
						p.source = $source,
						p.sourceOfTruth = $sourceOfTruth,
						p.appSource = $appSource,
						p.createdAt = $createdAt,
						p.updatedAt = $updatedAt,
						p.syncedWithEventStore = true 
		 ON MATCH SET 	p.syncedWithEventStore = true`, event.Tenant)

	return r.executeQuery(ctx, query, map[string]any{
		"id":             phoneNumberId,
		"rawPhoneNumber": event.RawPhoneNumber,
		"tenant":         event.Tenant,
		"source":         helper.GetSource(utils.StringFirstNonEmpty(event.SourceFields.Source, event.Source)),
		"sourceOfTruth":  helper.GetSourceOfTruth(utils.StringFirstNonEmpty(event.SourceFields.SourceOfTruth, event.SourceOfTruth)),
		"appSource":      helper.GetAppSource(utils.StringFirstNonEmpty(event.SourceFields.AppSource, event.AppSource)),
		"createdAt":      event.CreatedAt,
		"updatedAt":      event.UpdatedAt,
	})
}

func (r *phoneNumberRepository) UpdatePhoneNumber(ctx context.Context, phoneNumberId string, event events.PhoneNumberUpdatedEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.UpdatePhoneNumber")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("phoneNumberId", phoneNumberId))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber:PhoneNumber_%s {id:$id})
		 SET 	p.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE p.sourceOfTruth END,
				p.updatedAt = $updatedAt,
				p.syncedWithEventStore = true`, event.Tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	return r.executeQuery(ctx, query, map[string]any{
		"id":            phoneNumberId,
		"tenant":        event.Tenant,
		"sourceOfTruth": event.Source,
		"updatedAt":     event.UpdatedAt,
		"overwrite":     event.Source == constants.SourceOpenline,
	})
}

func (r *phoneNumberRepository) LinkWithContact(ctx context.Context, tenant, contactId, phoneNumberId, label string, primary bool, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.LinkWithContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("phoneNumberId", phoneNumberId), log.String("contactId", contactId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `
		MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}),
				(t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId})
		MERGE (c)-[rel:HAS]->(p)
		SET	rel.primary = $primary,
			rel.label = $label,	
			c.updatedAt = $updatedAt,
			rel.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":        tenant,
				"contactId":     contactId,
				"phoneNumberId": phoneNumberId,
				"label":         label,
				"primary":       primary,
				"updatedAt":     updatedAt,
			})
		if err != nil {
			return nil, err
		}
		return nil, err
	})
	return err
}

func (r *phoneNumberRepository) LinkWithOrganization(ctx context.Context, tenant, organizationId, phoneNumberId, label string, primary bool, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.LinkWithOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("phoneNumberId", phoneNumberId), log.String("organizationId", organizationId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `
		MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId}),
				(t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId})
		MERGE (org)-[rel:HAS]->(p)
		SET	rel.primary = $primary,
			rel.label = $label,	
			org.updatedAt = $updatedAt,
			rel.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
				"phoneNumberId":  phoneNumberId,
				"label":          label,
				"primary":        primary,
				"updatedAt":      updatedAt,
			})
		if err != nil {
			return nil, err
		}
		return nil, err
	})
	return err
}

func (r *phoneNumberRepository) LinkWithUser(ctx context.Context, tenant, userId, phoneNumberId, label string, primary bool, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.LinkWithUser")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("phoneNumberId", phoneNumberId), log.String("userId", userId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `
		MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId}),
				(t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId})
		MERGE (u)-[rel:HAS]->(p)
		SET	rel.primary = $primary,
			rel.label = $label,	
			u.updatedAt = $updatedAt,
			rel.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":        tenant,
				"userId":        userId,
				"phoneNumberId": phoneNumberId,
				"label":         label,
				"primary":       primary,
				"updatedAt":     updatedAt,
			})
		if err != nil {
			return nil, err
		}
		return nil, err
	})
	return err
}

func (r *phoneNumberRepository) FailPhoneNumberValidation(ctx context.Context, phoneNumberId string, event events.PhoneNumberFailedValidationEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.FailPhoneNumberValidation")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("phoneNumberId", phoneNumberId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber:PhoneNumber_%s {id:$id})
		 		SET p.validationError = $validationError,
		     		p.validated = false,
					p.updatedAt = $validatedAt`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":              phoneNumberId,
				"tenant":          event.Tenant,
				"validationError": event.ValidationError,
				"validatedAt":     event.ValidatedAt,
			})
		return nil, err
	})
	return err
}

func (r *phoneNumberRepository) PhoneNumberValidated(ctx context.Context, phoneNumberId string, event events.PhoneNumberValidatedEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.PhoneNumberValidated")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("phoneNumberId", phoneNumberId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber:PhoneNumber_%s {id:$id})
		 		SET p.validationError = $validationError,
					p.e164 = $e164,
		     		p.validated = true,
					p.updatedAt = $validatedAt
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
									c.updatedAt = $now, 
									c.appSource = $appSource,
									c.source = $source,
									c.sourceOfTruth = $source
				MERGE (p)-[:LINKED_TO]->(c)
				`
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":              phoneNumberId,
				"tenant":          event.Tenant,
				"validationError": "",
				"e164":            event.E164,
				"validatedAt":     event.ValidatedAt,
				"countryCodeA2":   event.CountryCodeA2,
				"now":             utils.Now(),
				"appSource":       "validation-api",
				"source":          "openline",
			})
		return nil, err
	})
	return err
}

func (r *phoneNumberRepository) GetCountryCodeA2ForPhoneNumber(ctx context.Context, tenant, phoneNumberId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.GetCountryCodeA2ForPhoneNumber")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("phoneNumberId", phoneNumberId))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (p:PhoneNumber {id:$phoneNumberId})-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
				OPTIONAL MATCH (p)-[:LINKED_TO]->(c:Country)
				OPTIONAL MATCH (tenant)-[:DEFAULT_COUNTRY]->(dc:Country)
				RETURN COALESCE(c.codeA2, dc.codeA2, '') AS countryCodeA2 LIMIT 1`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":        tenant,
				"phoneNumberId": phoneNumberId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsString(ctx, queryResult, err)
		}
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (r *phoneNumberRepository) executeQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteQuery(ctx, *r.driver, query, params)
}
