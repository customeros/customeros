package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	"golang.org/x/net/context"
)

type ContactRepository interface {
	CreateContact(ctx context.Context, aggregateId string, event events.ContactCreatedEvent) error
	UpdateContact(ctx context.Context, aggregateId string, event events.ContactUpdatedEvent) error
}

type contactRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContactRepository(driver *neo4j.DriverWithContext) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) CreateContact(ctx context.Context, aggregateId string, event events.ContactCreatedEvent) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:CONTACT_BELONGS_TO_TENANT]-(p:Contact:Contact_%s {id:$id}) 
		 ON CREATE SET 	p.firstname = $firstName,
						p.lastname = $lastName,	
						p.name = $name,	
						p.prefix = $prefix,
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
				"id":            aggregateId,
				"firstName":     event.FirstName,
				"lastName":      event.LastName,
				"name":          event.Name,
				"prefix":        event.Prefix,
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

func (r *contactRepository) UpdateContact(ctx context.Context, aggregateId string, event events.ContactUpdatedEvent) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(p:Contact:Contact_%s {id:$id})
		 SET	p.firstName = $firstName,
				p.lastName = $lastName,
				p.name = $name,
				p.prefix = $prefix,
				p.sourceOfTruth = $sourceOfTruth,
				p.updatedAt = $updatedAt,
				p.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":            aggregateId,
				"tenant":        event.Tenant,
				"firstName":     event.FirstName,
				"lastName":      event.LastName,
				"name":          event.Name,
				"prefix":        event.Prefix,
				"sourceOfTruth": event.SourceOfTruth,
				"updatedAt":     event.UpdatedAt,
			})
		return nil, err
	})
	return err
}
