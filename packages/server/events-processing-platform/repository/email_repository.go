package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"golang.org/x/net/context"
	"time"
)

type EmailRepository interface {
	GetIdIfExists(ctx context.Context, tenant, email string) (string, error)
	CreateEmail(ctx context.Context, aggregateId string, event events.EmailCreatedEvent) error
	UpdateEmail(ctx context.Context, aggregateId string, event events.EmailUpdatedEvent) error
	LinkWithContact(ctx context.Context, tenant, contactId, emailId, label string, primary bool, updatedAt time.Time) error
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

func (r *emailRepository) CreateEmail(ctx context.Context, aggregateId string, event events.EmailCreatedEvent) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email:Email_%s {id:$id}) 
		 ON CREATE SET e.rawEmail = $rawEmail, 
						e.validated = null,
						e.source = $source,
						e.sourceOfTruth = $sourceOfTruth,
						e.appSource = $appSource,
						e.createdAt = $createdAt,
						e.updatedAt = $updatedAt,
						e.syncedWithEventStore = true 
		 ON MATCH SET 	e.syncedWithEventStore = true
`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":            aggregateId,
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

func (r *emailRepository) UpdateEmail(ctx context.Context, aggregateId string, event events.EmailUpdatedEvent) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email:Email_%s {id:$id})
		 SET e.sourceOfTruth = $sourceOfTruth,
			e.updatedAt = $updatedAt,
			e.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":            aggregateId,
				"tenant":        event.Tenant,
				"sourceOfTruth": event.SourceOfTruth,
				"updatedAt":     event.UpdatedAt,
			})
		return nil, err
	})
	return err
}

func (r *emailRepository) LinkWithContact(ctx context.Context, tenant, contactId, emailId, label string, primary bool, updatedAt time.Time) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `
		MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}),
				(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
		MERGE (c)-[rel:HAS]->(p)
		SET	rel.primary = $primary,
			rel.label = $label,	
			c.updatedAt = $updatedAt,
			rel.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
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
