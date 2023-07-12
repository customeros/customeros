package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type EmailRepository interface {
	GetIdIfExists(ctx context.Context, tenant, email string) (string, error)
	CreateEmail(ctx context.Context, emailId string, event events.EmailCreateEvent) error
	UpdateEmail(ctx context.Context, emailId string, event events.EmailUpdateEvent) error
	FailEmailValidation(ctx context.Context, emailId string, event events.EmailFailedValidationEvent) error
	EmailValidated(ctx context.Context, emailId string, event events.EmailValidatedEvent) error
	LinkWithContact(ctx context.Context, tenant, contactId, emailId, label string, primary bool, updatedAt time.Time) error
	LinkWithOrganization(ctx context.Context, tenant, organizationId, emailId, label string, primary bool, updatedAt time.Time) error
	LinkWithUser(ctx context.Context, tenant, userId, emailId, label string, primary bool, updatedAt time.Time) error
}

type emailRepository struct {
	driver *neo4j.DriverWithContext
}

func NewEmailRepository(driver *neo4j.DriverWithContext) EmailRepository {
	return &emailRepository{
		driver: driver,
	}
}

func (r *emailRepository) GetIdIfExists(ctx context.Context, tenant string, email string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.GetIdIfExists")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("email", email))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (e:Email_%s) WHERE e.email = $email OR e.rawEmail = $email RETURN e.id LIMIT 1"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"email": email,
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

func (r *emailRepository) CreateEmail(ctx context.Context, emailId string, event events.EmailCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.CreateEmail")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("emailId", emailId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant}) 
              MERGE (e:Email:Email_%s {id:$id})
				 SET e.rawEmail = $rawEmail, 
					e.validated = null,
					e.source = $source,
					e.sourceOfTruth = $sourceOfTruth,
					e.appSource = $appSource,
					e.createdAt = $createdAt,
					e.updatedAt = $updatedAt,
					e.syncedWithEventStore = true 
		 MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e)

`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":            emailId,
				"rawEmail":      event.RawEmail,
				"tenant":        event.Tenant,
				"source":        event.Source,
				"sourceOfTruth": event.SourceOfTruth,
				"appSource":     event.AppSource,
				"createdAt":     event.CreatedAt,
				"updatedAt":     event.UpdatedAt,
			})
		return nil, err
	})
	return err
}

func (r *emailRepository) UpdateEmail(ctx context.Context, emailId string, event events.EmailUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.UpdateEmail")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("emailId", emailId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email:Email_%s {id:$id})
		 SET e.sourceOfTruth = $sourceOfTruth,
			e.updatedAt = $updatedAt,
			e.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":            emailId,
				"tenant":        event.Tenant,
				"sourceOfTruth": event.SourceOfTruth,
				"updatedAt":     event.UpdatedAt,
			})
		return nil, err
	})
	return err
}

func (r *emailRepository) FailEmailValidation(ctx context.Context, emailId string, event events.EmailFailedValidationEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.FailEmailValidation")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("emailId", emailId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email:Email_%s {id:$id})
		 		SET e.validationError = $validationError,
		     		e.validated = false,
					e.updatedAt = $validatedAt`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":              emailId,
				"tenant":          event.Tenant,
				"validationError": event.ValidationError,
				"validatedAt":     event.ValidatedAt,
			})
		return nil, err
	})
	return err
}

func (r *emailRepository) EmailValidated(ctx context.Context, emailId string, event events.EmailValidatedEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.EmailValidated")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("emailId", emailId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email:Email_%s {id:$id})
		 		SET e.validationError = $validationError,
					e.email = $email,
		     		e.validated = true,
					e.acceptsMail = $acceptsMail,
					e.canConnectSmtp = $canConnectSmtp,
					e.hasFullInbox = $hasFullInbox,
					e.isCatchAll = $isCatchAll,
					e.isDeliverable = $isDeliverable,
					e.isDisabled = $isDisabled,
					e.isValidSyntax = $isValidSyntax,
					e.username = $username,
					e.updatedAt = $validatedAt,
					e.isReachable = $isReachable
				WITH e
				MERGE (d:Domain {domain:$domain})
				ON CREATE SET 	d.id=randomUUID(), 
								d.createdAt=$now, 
								d.updatedAt=$now,
								d.appSource=$source,
								d.source=$appSource
				WITH d, e
				MERGE (e)-[:HAS_DOMAIN]->(d)`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":              emailId,
				"tenant":          event.Tenant,
				"validationError": event.ValidationError,
				"email":           event.EmailAddress,
				"domain":          event.Domain,
				"acceptsMail":     event.AcceptsMail,
				"canConnectSmtp":  event.CanConnectSmtp,
				"hasFullInbox":    event.HasFullInbox,
				"isCatchAll":      event.IsCatchAll,
				"isDeliverable":   event.IsDeliverable,
				"isDisabled":      event.IsDisabled,
				"isValidSyntax":   event.IsValidSyntax,
				"username":        event.Username,
				"validatedAt":     event.ValidatedAt,
				"isReachable":     event.IsReachable,
				"now":             utils.Now(),
				"source":          constants.SourceOpenline,
				"appSource":       constants.SourceEventProcessingPlatform,
			})
		return nil, err
	})
	return err
}

func (r *emailRepository) LinkWithContact(ctx context.Context, tenant, contactId, emailId, label string, primary bool, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.LinkWithContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("emailId", emailId), log.String("contactId", contactId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `
		MATCH (t:Tenant {name:$tenant})
		MERGE (c:Contact:Contact_%s {id:$contactId})
		ON CREATE SET c.syncedWithEventStore = true
		MERGE (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c)
		MERGE (e:Email:Email_%s {id:$emailId})
		ON CREATE SET e.syncedWithEventStore = true
		MERGE (e)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
		MERGE (c)-[rel:HAS]->(e)
		SET	rel.primary = $primary,
			rel.label = $label,	
			c.updatedAt = $updatedAt,
			rel.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
				"emailId":   emailId,
				"label":     label,
				"primary":   primary,
				"updatedAt": updatedAt,
			})
		if err != nil {
			return nil, err
		}
		return nil, err
	})
	return err
}

func (r *emailRepository) LinkWithOrganization(ctx context.Context, tenant, organizationId, emailId, label string, primary bool, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.LinkWithOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("emailId", emailId), log.String("organizationId", organizationId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `
		MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId}),
				(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
		MERGE (org)-[rel:HAS]->(e)
		SET	rel.primary = $primary,
			rel.label = $label,	
			org.updatedAt = $updatedAt,
			rel.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
				"emailId":        emailId,
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

func (r *emailRepository) LinkWithUser(ctx context.Context, tenant, userId, emailId, label string, primary bool, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.LinkWithUser")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("emailId", emailId), log.String("userId", userId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `
		MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId}),
				(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
		MERGE (u)-[rel:HAS]->(e)
		SET	rel.primary = $primary,
			rel.label = $label,	
			u.updatedAt = $updatedAt,
			rel.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":    tenant,
				"userId":    userId,
				"emailId":   emailId,
				"label":     label,
				"primary":   primary,
				"updatedAt": updatedAt,
			})
		if err != nil {
			return nil, err
		}
		return nil, err
	})
	return err
}
