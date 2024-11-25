package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ContactWriteRepository interface {
	SaveContactInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, contactId string, data data_fields.ContactFields, updateOnlyIfEmpty bool) error
	ResetEnrichAttempts(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, contactId string) error
}

type contactWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewContactWriteRepository(driver *neo4j.DriverWithContext, database string) ContactWriteRepository {
	return &contactWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *contactWriteRepository) SaveContactInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, contactId string, data data_fields.ContactFields, updateOnlyIfEmpty bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.SaveContactInTx")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)
	tracing.LogObjectAsJson(span, "data", data)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		cypher := fmt.Sprintf(`
				MATCH (t:Tenant {name:$tenant})
				MERGE (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})
				ON CREATE SET
					c:Contact_%s,
					c.createdAt = $createdAt,
					c.hide = false,
					c.source = $source,
					c.appSource = $appSource
				WITH c
				SET
					c.updatedAt = datetime()`, tenant)
		params := map[string]any{
			"tenant":            tenant,
			"contactId":         contactId,
			"createdAt":         utils.TimePtrAsAny(data.CreatedAt),
			"source":            utils.IfNotNilString(data.Source),
			"appSource":         utils.IfNotNilString(data.AppSource),
			"updateOnlyIfEmpty": updateOnlyIfEmpty,
		}

		if data.FirstName != nil {
			cypher += ", c.firstName = CASE WHEN $updateOnlyIfEmpty = false OR c.firstName is null OR c.firstName = '' THEN $firstName ELSE c.firstName END"
			params["firstName"] = *data.FirstName
		}
		if data.LastName != nil {
			cypher += ", c.lastName = CASE WHEN $updateOnlyIfEmpty = false OR c.lastName is null OR c.lastName = '' THEN $lastName ELSE c.lastName END"
			params["lastName"] = *data.LastName
		}
		if data.Name != nil {
			cypher += ", c.name = CASE WHEN $updateOnlyIfEmpty = false OR c.name is null OR c.name = '' THEN $name ELSE c.name END"
			params["name"] = *data.Name
		}
		if data.Prefix != nil {
			cypher += ", c.prefix = CASE WHEN $updateOnlyIfEmpty = false OR c.prefix is null OR c.prefix = '' THEN $prefix ELSE c.prefix END"
			params["prefix"] = *data.Prefix
		}
		if data.Description != nil {
			cypher += ", c.description = CASE WHEN $updateOnlyIfEmpty = false OR c.description is null OR c.description = '' THEN $description ELSE c.description END"
			params["description"] = *data.Description
		}
		if data.Timezone != nil {
			cypher += ", c.timezone = CASE WHEN $updateOnlyIfEmpty = false OR c.timezone is null OR c.timezone = '' THEN $timezone ELSE c.timezone END"
			params["timezone"] = *data.Timezone
		}
		if data.ProfilePhotoUrl != nil {
			cypher += ", c.profilePhotoUrl = CASE WHEN $updateOnlyIfEmpty = false OR c.profilePhotoUrl is null OR c.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE c.profilePhotoUrl END"
			params["profilePhotoUrl"] = *data.ProfilePhotoUrl
		}
		if data.Username != nil {
			cypher += ", c.username = CASE WHEN $updateOnlyIfEmpty = false OR c.username is null OR c.username = '' THEN $username ELSE c.username END"
			params["username"] = *data.Username
		}
		if data.Hide != nil {
			cypher += ", c.hide = $hide"
			cypher += `, c.hiddenAt = CASE WHEN $hide = true THEN datetime() ELSE null END `
			params["hide"] = *data.Hide
		}

		span.LogFields(log.String("cypher", cypher))
		tracing.LogObjectAsJson(span, "params", params)

		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *contactWriteRepository) ResetEnrichAttempts(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, contactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.ResetEnrichAttempts")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, contactId)

	cypher := `MATCH (t:Tenant {name: $tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id: $contactId})
	WHERE c.enrichedAt IS NULL
	REMOVE c.techEnrichAttempts, c.techEnrichRequestedAt`
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}
