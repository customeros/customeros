package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ContactRepository interface {
	CreateContact(ctx context.Context, contactId string, event event.ContactCreateEvent) error
	CreateContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, contactId string, event event.ContactCreateEvent) error
	UpdateContact(ctx context.Context, contactId string, event event.ContactUpdateEvent) error
}

type contactRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContactRepository(driver *neo4j.DriverWithContext) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) CreateContact(ctx context.Context, contactId string, event event.ContactCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.CreateContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("contactId", contactId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.CreateContactInTx(ctx, tx, contactId, event)
	})
	return err
}

func (r *contactRepository) CreateContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, contactId string, event event.ContactCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.CreateContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("contactId", contactId))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
				MERGE (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact:Contact_%s {id:$id}) 
		 		ON CREATE SET
						c.firstName = $firstName,
						c.lastName = $lastName,	
						c.prefix = $prefix,
						c.description = $description,
						c.timezone = $timezone,
						c.profilePhotoUrl = $profilePhotoUrl,
						c.name = $name,
						c.source = $source,
						c.sourceOfTruth = $sourceOfTruth,
						c.appSource = $appSource,
						c.createdAt = $createdAt,
						c.updatedAt = $updatedAt,
						c.syncedWithEventStore = true
				ON MATCH SET
						c.name = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.name is null OR c.name = '' THEN $name ELSE c.name END,
						c.firstName = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.firstName is null OR c.firstName = '' THEN $firstName ELSE c.firstName END,
						c.lastName = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.lastName is null OR c.lastName = '' THEN $lastName ELSE c.lastName END,
						c.timezone = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.timezone is null OR c.timezone = '' THEN $timezone ELSE c.timezone END,
						c.profilePhotoUrl = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.profilePhotoUrl is null OR c.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE c.profilePhotoUrl END,
						c.prefix = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.prefix is null OR c.prefix = '' THEN $prefix ELSE c.prefix END,
						c.description = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.description is null OR c.description = '' THEN $description ELSE c.description END,
						c.updatedAt = $updatedAt,
						c.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE c.sourceOfTruth END,
						c.syncedWithEventStore = true
				`, event.Tenant)
	span.LogFields(log.String("query", query))

	return utils.ExecuteQueryInTx(ctx, tx, query, map[string]any{
		"id":              contactId,
		"firstName":       event.FirstName,
		"lastName":        event.LastName,
		"prefix":          event.Prefix,
		"description":     event.Description,
		"timezone":        event.Timezone,
		"profilePhotoUrl": event.ProfilePhotoUrl,
		"name":            event.Name,
		"tenant":          event.Tenant,
		"source":          helper.GetSource(event.Source),
		"sourceOfTruth":   helper.GetSourceOfTruth(event.SourceOfTruth),
		"appSource":       helper.GetAppSource(event.AppSource),
		"createdAt":       event.CreatedAt,
		"updatedAt":       event.UpdatedAt,
		"overwrite":       helper.GetSourceOfTruth(event.SourceOfTruth) == constants.SourceOpenline,
	})
}

func (r *contactRepository) UpdateContact(ctx context.Context, contactId string, event event.ContactUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.UpdateContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("contactId", contactId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact:Contact_%s {id:$id})
		 SET	c.name = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.name is null OR c.name = '' THEN $name ELSE c.name END,
				c.firstName = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.firstName is null OR c.firstName = '' THEN $firstName ELSE c.firstName END,
				c.lastName = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.lastName is null OR c.lastName = '' THEN $lastName ELSE c.lastName END,
				c.timezone = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.timezone is null OR c.timezone = '' THEN $timezone ELSE c.timezone END,
				c.profilePhotoUrl = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.profilePhotoUrl is null OR c.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE c.profilePhotoUrl END,
				c.prefix = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.prefix is null OR c.prefix = '' THEN $prefix ELSE c.prefix END,
				c.description = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.description is null OR c.description = '' THEN $description ELSE c.description END,
				c.updatedAt = $updatedAt,
				c.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE c.sourceOfTruth END,
				c.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":              contactId,
				"tenant":          event.Tenant,
				"firstName":       event.FirstName,
				"lastName":        event.LastName,
				"prefix":          event.Prefix,
				"description":     event.Description,
				"timezone":        event.Timezone,
				"profilePhotoUrl": event.ProfilePhotoUrl,
				"name":            event.Name,
				"updatedAt":       event.UpdatedAt,
				"sourceOfTruth":   helper.GetSource(event.Source),
				"overwrite":       helper.GetSource(event.Source) == constants.SourceOpenline,
			})
		return nil, err
	})
	return err
}
