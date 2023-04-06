package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"golang.org/x/net/context"
	"time"
)

type PhoneNumberRepository interface {
	GetIdIfExists(ctx context.Context, tenant, phoneNumber string) (string, error)
	CreatePhoneNumber(ctx context.Context, aggregateId string, event events.PhoneNumberCreatedEvent) error
	UpdatePhoneNumber(ctx context.Context, aggregateId string, event events.PhoneNumberUpdatedEvent) error
	LinkWithContact(ctx context.Context, tenant, contactId, phoneNumberId, label string, primary bool, updatedAt time.Time) error
}

type phoneNumberRepository struct {
	driver *neo4j.DriverWithContext
}

func NewPhoneNumberRepository(driver *neo4j.DriverWithContext) PhoneNumberRepository {
	return &phoneNumberRepository{
		driver: driver,
	}
}

func (r *phoneNumberRepository) GetIdIfExists(ctx context.Context, tenant string, phoneNumber string) (string, error) {
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

func (r *phoneNumberRepository) CreatePhoneNumber(ctx context.Context, aggregateId string, event events.PhoneNumberCreatedEvent) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber:PhoneNumber_%s {id:$id}) 
		 ON CREATE SET p.rawPhoneNumber = $rawPhoneNumber, 
						p.validated = null,
						p.source = $source,
						p.sourceOfTruth = $sourceOfTruth,
						p.appSource = $appSource,
						p.createdAt = $createdAt,
						p.updatedAt = $updatedAt,
						p.syncedWithEventStore = true 
		 ON MATCH SET 	p.syncedWithEventStore = true
`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":             aggregateId,
				"rawPhoneNumber": event.RawPhoneNumber,
				"tenant":         event.Tenant,
				"source":         event.Source,
				"sourceOfTruth":  event.SourceOfTruth,
				"appSource":      event.AppSource,
				"createdAt":      event.CreatedAt,
				"updatedAt":      event.UpdatedAt,
			})
		return nil, err
	})
	return err
}

func (r *phoneNumberRepository) UpdatePhoneNumber(ctx context.Context, aggregateId string, event events.PhoneNumberUpdatedEvent) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber:PhoneNumber_%s {id:$id})
		 SET p.sourceOfTruth = $sourceOfTruth,
			p.updatedAt = $updatedAt,
			p.syncedWithEventStore = true`

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

func (r *phoneNumberRepository) LinkWithContact(ctx context.Context, tenant, contactId, phoneNumberId, label string, primary bool, updatedAt time.Time) error {
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
