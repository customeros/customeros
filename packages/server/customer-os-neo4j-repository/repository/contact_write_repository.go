package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type ContactFields struct {
	FirstName       string             `json:"firstName"`
	LastName        string             `json:"lastName"`
	Prefix          string             `json:"prefix"`
	Description     string             `json:"description"`
	Timezone        string             `json:"timezone"`
	ProfilePhotoUrl string             `json:"profilePhotoUrl"`
	Username        string             `json:"username"`
	Name            string             `json:"name"`
	SourceFields    model.SourceFields `json:"sourceFields"`
	CreatedAt       time.Time          `json:"createdAt"`

	UpdateFirstName       bool `json:"updateFirstName"`
	UpdateLastName        bool `json:"updateLastName"`
	UpdateName            bool `json:"updateName"`
	UpdatePrefix          bool `json:"updatePrefix"`
	UpdateDescription     bool `json:"updateDescription"`
	UpdateTimezone        bool `json:"updateTimezone"`
	UpdateProfilePhotoUrl bool `json:"updateProfilePhotoUrl"`
	UpdateUsername        bool `json:"updateUsername"`

	UpdateOnlyIfEmpty bool `json:"updateIfEmpty"`
}

type ContactWriteRepository interface {
	// Deprecated, use SaveContact instead
	CreateContact(ctx context.Context, tenant, contactId string, data ContactFields) error
	// Deprecated, use SaveContact instead
	CreateContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, data ContactFields) error
	// Deprecated, use SaveContact instead
	UpdateContact(ctx context.Context, tenant, contactId string, data ContactFields) error
	SaveContactInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, contactId string, data ContactFields) error
	// Deprecated, Use CommonRepository.UpdateAnyProperty instead
	UpdateAnyProperty(ctx context.Context, tenant, contactId string, property entity.ContactProperty, value any) error
}

type contactWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func (r *contactWriteRepository) SaveContactInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, contactId string, data ContactFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.SaveContactInTx")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		cypher := fmt.Sprintf(`
				MATCH (t:Tenant {name:$tenant})
				MERGE (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})
				ON CREATE SET
					c:Contact_%s,
					c.createdAt = $createdAt,
					c.hide = false,
					c.source = $source
				WITH c
				SET
					c.updatedAt = datetime()`, tenant)
		params := map[string]any{
			"tenant":            tenant,
			"contactId":         contactId,
			"createdAt":         data.CreatedAt,
			"source":            data.SourceFields.GetSource(),
			"updateOnlyIfEmpty": data.UpdateOnlyIfEmpty,
		}

		if data.UpdateFirstName {
			cypher += ", c.firstName = CASE WHEN $updateOnlyIfEmpty = false OR c.firstName is null OR c.firstName = '' THEN $firstName ELSE c.firstName END"
			params["firstName"] = data.FirstName
		}
		if data.UpdateLastName {
			cypher += ", c.lastName = CASE WHEN $updateOnlyIfEmpty = false OR c.lastName is null OR c.lastName = '' THEN $lastName ELSE c.lastName END"
			params["lastName"] = data.LastName
		}
		if data.UpdateName {
			cypher += ", c.name = CASE WHEN $updateOnlyIfEmpty = false OR c.name is null OR c.name = '' THEN $name ELSE c.name END"
			params["name"] = data.Name
		}
		if data.UpdatePrefix {
			cypher += ", c.prefix = CASE WHEN $updateOnlyIfEmpty = false OR c.prefix is null OR c.prefix = '' THEN $prefix ELSE c.prefix END"
			params["prefix"] = data.Prefix
		}
		if data.UpdateDescription {
			cypher += ", c.description = CASE WHEN $updateOnlyIfEmpty = false OR c.description is null OR c.description = '' THEN $description ELSE c.description END"
			params["description"] = data.Description
		}
		if data.UpdateTimezone {
			cypher += ", c.timezone = CASE WHEN $updateOnlyIfEmpty = false OR c.timezone is null OR c.timezone = '' THEN $timezone ELSE c.timezone END"
			params["timezone"] = data.Timezone
		}
		if data.UpdateProfilePhotoUrl {
			cypher += ", c.profilePhotoUrl = CASE WHEN $updateOnlyIfEmpty = false OR c.profilePhotoUrl is null OR c.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE c.profilePhotoUrl END"
			params["profilePhotoUrl"] = data.ProfilePhotoUrl
		}
		if data.UpdateUsername {
			cypher += ", c.username = CASE WHEN $updateOnlyIfEmpty = false OR c.username is null OR c.username = '' THEN $username ELSE c.username END"
			params["username"] = data.Username
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

func NewContactWriteRepository(driver *neo4j.DriverWithContext, database string) ContactWriteRepository {
	return &contactWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *contactWriteRepository) CreateContact(ctx context.Context, tenant, contactId string, data ContactFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.CreateContact")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.CreateContactInTx(ctx, tx, tenant, contactId, data)
	})
	return err
}

func (r *contactWriteRepository) CreateContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, data ContactFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.CreateContactInTx")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
	MERGE (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact:Contact_%s {id:$id})
	ON CREATE SET
	c.firstName = $firstName,
	c.lastName = $lastName,
	c.prefix = $prefix,
	c.description = $description,
	c.timezone = $timezone,
	c.profilePhotoUrl = $profilePhotoUrl,
	c.username = $username,
	c.name = $name,
	c.source = $source,
	c.appSource = $appSource,
	c.createdAt = $createdAt,
	c.updatedAt = datetime(),
	c.hide = $hide
	ON MATCH SET
	c.name = CASE WHEN $overwrite = true OR c.name is null OR c.name = '' THEN $name ELSE c.name END,
	c.firstName = CASE WHEN $overwrite =true OR c.firstName is null OR c.firstName = '' THEN $firstName ELSE c.firstName END,
	c.lastName = CASE WHEN $overwrite = true OR c.lastName is null OR c.lastName = '' THEN $lastName ELSE c.lastName END,
	c.timezone = CASE WHEN $overwrite = true OR c.timezone is null OR c.timezone = '' THEN $timezone ELSE c.timezone END,
	c.profilePhotoUrl = CASE WHEN OR $overwrite = true OR c.profilePhotoUrl is null OR c.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE c.profilePhotoUrl END,
	c.username = CASE WHEN OR $overwrite = true OR c.username is null OR c.username = '' THEN $username ELSE c.username END,
	c.prefix = CASE WHEN OR $overwrite =true OR c.prefix is null OR c.prefix = '' THEN $prefix ELSE c.prefix END,
	c.description = CASE WHEN OR $overwrite = true OR c.description is null OR c.description = '' THEN $description ELSE c.description END,
	c.updatedAt = datetime()
	`, tenant)
	params := map[string]any{
		"id":              contactId,
		"firstName":       data.FirstName,
		"lastName":        data.LastName,
		"prefix":          data.Prefix,
		"description":     data.Description,
		"timezone":        data.Timezone,
		"profilePhotoUrl": data.ProfilePhotoUrl,
		"username":        data.Username,
		"name":            data.Name,
		"tenant":          tenant,
		"source":          data.SourceFields.Source,
		"appSource":       data.SourceFields.AppSource,
		"createdAt":       data.CreatedAt,
		"overwrite":       true,
		"hide":            false,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteQueryInTx(ctx, tx, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contactWriteRepository) UpdateContact(ctx context.Context, tenant, contactId string, data ContactFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.UpdateContact")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact:Contact_%s {id:$id})
	SET    c.updatedAt = datetime()`, tenant)

	params := map[string]any{
		"id":        contactId,
		"tenant":    tenant,
		"overwrite": data.SourceFields.Source == constants.SourceOpenline,
	}

	if data.UpdateFirstName {
		cypher += ", c.firstName = CASE WHEN $overwrite=true OR c.firstName is null OR c.firstName = '' THEN $firstName ELSE c.firstName END"
		params["firstName"] = data.FirstName
	}
	if data.UpdateLastName {
		cypher += ", c.lastName = CASE WHEN $overwrite=true OR c.lastName is null OR c.lastName = '' THEN $lastName ELSE c.lastName END"
		params["lastName"] = data.LastName
	}
	if data.UpdateName {
		cypher += ", c.name = CASE WHEN $overwrite=true OR c.name is null OR c.name = '' THEN $name ELSE c.name END"
		params["name"] = data.Name
	}
	if data.UpdatePrefix {
		cypher += ", c.prefix = CASE WHEN $overwrite=true OR c.prefix is null OR c.prefix = '' THEN $prefix ELSE c.prefix END"
		params["prefix"] = data.Prefix
	}
	if data.UpdateDescription {
		cypher += ", c.description = CASE WHEN $overwrite=true OR c.description is null OR c.description = '' THEN $description ELSE c.description END"
		params["description"] = data.Description
	}
	if data.UpdateTimezone {
		cypher += ", c.timezone = CASE WHEN $overwrite=true OR c.timezone is null OR c.timezone = '' THEN $timezone ELSE c.timezone END"
		params["timezone"] = data.Timezone
	}
	if data.UpdateProfilePhotoUrl {
		cypher += ", c.profilePhotoUrl = CASE WHEN $overwrite=true OR c.profilePhotoUrl is null OR c.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE c.profilePhotoUrl END"
		params["profilePhotoUrl"] = data.ProfilePhotoUrl
	}
	if data.UpdateUsername {
		cypher += ", c.username = CASE WHEN $overwrite=true OR c.username is null OR c.username = '' THEN $username ELSE c.username END"
		params["username"] = data.Username
	}

	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contactWriteRepository) UpdateAnyProperty(ctx context.Context, tenant, contactId string, property entity.ContactProperty, value any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.UpdateTimeProperty")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)
	span.LogFields(log.String("property", string(property)), log.Object("value", value))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id: $contactId})
	SET c.%s = $value`, string(property))
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
		"value":     value,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
