package neo4j_repository

import (
	"fmt"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"golang.org/x/net/context"
	"time"
)

type PhoneNumberCreateFields struct {
	RawPhoneNumber string             `json:"rawPhoneNumber"`
	SourceFields   model.SourceFields `json:"sourceFields"`
	CreatedAt      time.Time          `json:"createdAt"`
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
	RemoveRelationship(ctx context.Context, entityType commonModel.EntityType, tenant, entityId, phoneNumber string) error
	RemoveRelationshipById(ctx context.Context, entityType commonModel.EntityType, tenant, entityId, phoneNumberId string) error
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberWriteRepository.CreatePhoneNumber")
	defer spans.Finish()

	spans.TagEntity(phoneNumberId)
	spans.LogObjectAsJson("data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber:PhoneNumber_%s {id:$id}) 
		 ON CREATE SET p.rawPhoneNumber = $rawPhoneNumber, 
						p.validated = null,
						p.source = $source,
						p.appSource = $appSource,
						p.createdAt = $createdAt,
						p.updatedAt = datetime()`, tenant)
	params := map[string]any{
		"id":             phoneNumberId,
		"rawPhoneNumber": data.RawPhoneNumber,
		"tenant":         tenant,
		"source":         data.SourceFields.Source,
		"appSource":      data.SourceFields.AppSource,
		"createdAt":      data.CreatedAt,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *phoneNumberWriteRepository) UpdatePhoneNumber(ctx context.Context, tenant, phoneNumberId, rawPhoneNumber, source string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberWriteRepository.UpdatePhoneNumber")
	defer spans.Finish()

	spans.TagEntity(phoneNumberId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$id})
				WHERE p:PhoneNumber_%s
		 SET 	p.updatedAt = datetime(),
				p.rawPhoneNumber = $rawPhoneNumber`, tenant)
	params := map[string]any{
		"id":             phoneNumberId,
		"tenant":         tenant,
		"rawPhoneNumber": rawPhoneNumber,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *phoneNumberWriteRepository) FailPhoneNumberValidation(ctx context.Context, tenant, phoneNumberId, validationError string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberWriteRepository.FailPhoneNumberValidation")
	defer spans.Finish()

	spans.TagEntity(phoneNumberId)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *phoneNumberWriteRepository) PhoneNumberValidated(ctx context.Context, tenant, phoneNumberId string, data PhoneNumberValidateFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberWriteRepository.PhoneNumberValidated")
	defer spans.Finish()

	spans.TagEntity(phoneNumberId)
	spans.LogObjectAsJson("data", data)

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
									c.source = $source
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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *phoneNumberWriteRepository) LinkWithContact(ctx context.Context, tenant, contactId, phoneNumberId, label string, primary bool) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberWriteRepository.LinkWithContact")
	defer spans.Finish()

	spans.TagEntity(phoneNumberId)

	cypher := `
		MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}),
				(t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId})
		MERGE (c)-[rel:HAS]->(p)
		SET	rel.primary = $primary,
			rel.label = $label,	
			c.updatedAt = datetime()`
	params := map[string]any{
		"tenant":        tenant,
		"contactId":     contactId,
		"phoneNumberId": phoneNumberId,
		"label":         label,
		"primary":       primary,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *phoneNumberWriteRepository) LinkWithOrganization(ctx context.Context, tenant, organizationId, phoneNumberId, label string, primary bool) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberWriteRepository.LinkWithOrganization")
	defer spans.Finish()

	spans.TagEntity(phoneNumberId)

	cypher := `
		MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId}),
				(t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId})
		MERGE (org)-[rel:HAS]->(p)
		SET	rel.primary = $primary,
			rel.label = $label,	
			org.updatedAt = datetime()`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"phoneNumberId":  phoneNumberId,
		"label":          label,
		"primary":        primary,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *phoneNumberWriteRepository) LinkWithUser(ctx context.Context, tenant, userId, phoneNumberId, label string, primary bool) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberWriteRepository.LinkWithUser")
	defer spans.Finish()

	spans.TagEntity(phoneNumberId)

	cypher := `
		MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId}),
				(t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId})
		MERGE (u)-[rel:HAS]->(p)
		SET	rel.primary = $primary,
			rel.label = $label,	
			u.updatedAt = datetime()`
	params := map[string]any{
		"tenant":        tenant,
		"userId":        userId,
		"phoneNumberId": phoneNumberId,
		"label":         label,
		"primary":       primary,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *phoneNumberWriteRepository) CleanPhoneNumberValidation(ctx context.Context, tenant, phoneNumberId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberWriteRepository.CleanPhoneNumberValidation")
	defer spans.Finish()

	spans.TagEntity(phoneNumberId)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *phoneNumberWriteRepository) RemoveRelationship(ctx context.Context, entityType commonModel.EntityType, tenant, entityId, phoneNumber string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberWriteRepository.RemoveRelationship")
	defer spans.Finish()

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := ""
	switch entityType {
	case commonModel.CONTACT:
		cypher = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case commonModel.USER:
		cypher = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case commonModel.ORGANIZATION:
		cypher = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	}
	cypher += `MATCH (entity)-[rel:HAS]->(p:PhoneNumber)
			WHERE p.e164 = $phoneNumber OR p.rawPhoneNumber = $phoneNumber
            DELETE rel`
	params := map[string]any{
		"entityId":    entityId,
		"phoneNumber": phoneNumber,
		"tenant":      tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *phoneNumberWriteRepository) RemoveRelationshipById(ctx context.Context, entityType commonModel.EntityType, tenant, entityId, phoneNumberId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberWriteRepository.RemoveRelationshipById")
	defer spans.Finish()

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := ""
	switch entityType {
	case commonModel.CONTACT:
		cypher = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case commonModel.USER:
		cypher = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case commonModel.ORGANIZATION:
		cypher = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	}
	cypher += `MATCH (entity)-[rel:HAS]->(p:PhoneNumber {id:$phoneNumberId})
            DELETE rel`
	params := map[string]any{
		"entityId":      entityId,
		"phoneNumberId": phoneNumberId,
		"tenant":        tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}
