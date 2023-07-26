package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ContactRepository interface {
	CreateContact(ctx context.Context, aggregateId string, event events.ContactCreateEvent) error
	UpdateContact(ctx context.Context, aggregateId string, event events.ContactUpdateEvent) error
}

type contactRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContactRepository(driver *neo4j.DriverWithContext) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) CreateContact(ctx context.Context, contactId string, event events.ContactCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.CreateContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("contactId", contactId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
				MERGE (p:Contact:Contact_%s {id:$id}) 
		 		SET 	p.firstName = $firstName,
						p.lastName = $lastName,	
						p.prefix = $prefix,
						p.description = $description,
						p.timezone = $timezone,
						p.name = $name,
						p.source = $source,
						p.sourceOfTruth = $sourceOfTruth,
						p.appSource = $appSource,
						p.createdAt = $createdAt,
						p.updatedAt = $updatedAt,
						p.syncedWithEventStore = true
				MERGE (t)<-[:CONTACT_BELONGS_TO_TENANT]-(p)`, event.Tenant)

	span.LogFields(log.String("query", query))

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"id":            contactId,
				"firstName":     event.FirstName,
				"lastName":      event.LastName,
				"prefix":        event.Prefix,
				"description":   event.Description,
				"timezone":      event.Timezone,
				"name":          event.Name,
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

func (r *contactRepository) UpdateContact(ctx context.Context, contactId string, event events.ContactUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.UpdateContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("contactId", contactId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(p:Contact:Contact_%s {id:$id})
		 SET	p.firstName = $firstName,
				p.lastName = $lastName,
				p.prefix = $prefix,
				p.description = $description,
				p.timezone = $timezone,
				p.name = $name,
				p.sourceOfTruth = $sourceOfTruth,
				p.updatedAt = $updatedAt,
				p.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":            contactId,
				"tenant":        event.Tenant,
				"firstName":     event.FirstName,
				"lastName":      event.LastName,
				"prefix":        event.Prefix,
				"description":   event.Description,
				"timezone":      event.Timezone,
				"name":          event.Name,
				"sourceOfTruth": event.SourceOfTruth,
				"updatedAt":     event.UpdatedAt,
			})
		return nil, err
	})
	return err
}
