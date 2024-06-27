package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type ContactCreateFields struct {
	FirstName        string       `json:"firstName"`
	LastName         string       `json:"lastName"`
	Prefix           string       `json:"prefix"`
	Description      string       `json:"description"`
	Timezone         string       `json:"timezone"`
	ProfilePhotoUrl  string       `json:"profilePhotoUrl"`
	Name             string       `json:"name"`
	CreatedAt        time.Time    `json:"createdAt"`
	SourceFields     model.Source `json:"sourceFields"`
	AggregateVersion int64        `json:"aggregateVersion"`
}

type ContactUpdateFields struct {
	AggregateVersion      int64  `json:"aggregateVersion"`
	FirstName             string `json:"firstName"`
	LastName              string `json:"lastName"`
	Prefix                string `json:"prefix"`
	Description           string `json:"description"`
	Timezone              string `json:"timezone"`
	ProfilePhotoUrl       string `json:"profilePhotoUrl"`
	Name                  string `json:"name"`
	Source                string `json:"source"`
	UpdateFirstName       bool   `json:"updateFirstName"`
	UpdateLastName        bool   `json:"updateLastName"`
	UpdateName            bool   `json:"updateName"`
	UpdatePrefix          bool   `json:"updatePrefix"`
	UpdateDescription     bool   `json:"updateDescription"`
	UpdateTimezone        bool   `json:"updateTimezone"`
	UpdateProfilePhotoUrl bool   `json:"updateProfilePhotoUrl"`
}

type ContactWriteRepository interface {
	CreateContact(ctx context.Context, tenant, contactId string, data ContactCreateFields) error
	CreateContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, data ContactCreateFields) error
	UpdateContact(ctx context.Context, tenant, contactId string, data ContactUpdateFields) error
	UpdateTimeProperty(ctx context.Context, tenant, contactId string, property entity.ContactProperty, value *time.Time) error
	UpdateAnyProperty(ctx context.Context, tenant, contactId string, property entity.ContactProperty, value any) error
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

func (r *contactWriteRepository) CreateContact(ctx context.Context, tenant, contactId string, data ContactCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.CreateContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.CreateContactInTx(ctx, tx, tenant, contactId, data)
	})
	return err
}

func (r *contactWriteRepository) CreateContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, data ContactCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.CreateContactInTx")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
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
						c.name = $name,
						c.source = $source,
						c.sourceOfTruth = $sourceOfTruth,
						c.appSource = $appSource,
						c.createdAt = $createdAt,
						c.updatedAt = datetime(),
						c.aggregateVersion = $aggregateVersion,
						c.syncedWithEventStore = true
				ON MATCH SET
						c.name = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.name is null OR c.name = '' THEN $name ELSE c.name END,
						c.firstName = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.firstName is null OR c.firstName = '' THEN $firstName ELSE c.firstName END,
						c.lastName = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.lastName is null OR c.lastName = '' THEN $lastName ELSE c.lastName END,
						c.timezone = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.timezone is null OR c.timezone = '' THEN $timezone ELSE c.timezone END,
						c.profilePhotoUrl = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.profilePhotoUrl is null OR c.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE c.profilePhotoUrl END,
						c.prefix = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.prefix is null OR c.prefix = '' THEN $prefix ELSE c.prefix END,
						c.description = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.description is null OR c.description = '' THEN $description ELSE c.description END,
						c.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE c.sourceOfTruth END,
						c.aggregateVersion = $aggregateVersion,
						c.updatedAt = datetime(),
						c.syncedWithEventStore = true
				`, tenant)
	params := map[string]any{
		"id":               contactId,
		"firstName":        data.FirstName,
		"lastName":         data.LastName,
		"prefix":           data.Prefix,
		"description":      data.Description,
		"timezone":         data.Timezone,
		"profilePhotoUrl":  data.ProfilePhotoUrl,
		"name":             data.Name,
		"tenant":           tenant,
		"source":           data.SourceFields.Source,
		"sourceOfTruth":    data.SourceFields.SourceOfTruth,
		"appSource":        data.SourceFields.AppSource,
		"createdAt":        data.CreatedAt,
		"overwrite":        data.SourceFields.SourceOfTruth == constants.SourceOpenline,
		"aggregateVersion": data.AggregateVersion,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteQueryInTx(ctx, tx, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contactWriteRepository) UpdateContact(ctx context.Context, tenant, contactId string, data ContactUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.UpdateContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact:Contact_%s {id:$id})
				WHERE c.aggregateVersion IS NULL OR c.aggregateVersion < $aggregateVersion
		 SET	c.updatedAt = datetime(),
				c.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE c.sourceOfTruth END,
				c.aggregateVersion = $aggregateVersion,
				c.syncedWithEventStore = true`, tenant)

	params := map[string]any{
		"id":               contactId,
		"tenant":           tenant,
		"sourceOfTruth":    data.Source,
		"overwrite":        data.Source == constants.SourceOpenline,
		"aggregateVersion": data.AggregateVersion,
	}

	if data.UpdateFirstName {
		cypher += ", c.firstName = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.firstName is null OR c.firstName = '' THEN $firstName ELSE c.firstName END"
		params["firstName"] = data.FirstName
	}
	if data.UpdateLastName {
		cypher += ", c.lastName = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.lastName is null OR c.lastName = '' THEN $lastName ELSE c.lastName END"
		params["lastName"] = data.LastName
	}
	if data.UpdateName {
		cypher += ", c.name = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.name is null OR c.name = '' THEN $name ELSE c.name END"
		params["name"] = data.Name
	}
	if data.UpdatePrefix {
		cypher += ", c.prefix = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.prefix is null OR c.prefix = '' THEN $prefix ELSE c.prefix END"
		params["prefix"] = data.Prefix
	}
	if data.UpdateDescription {
		cypher += ", c.description = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.description is null OR c.description = '' THEN $description ELSE c.description END"
		params["description"] = data.Description
	}
	if data.UpdateTimezone {
		cypher += ", c.timezone = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.timezone is null OR c.timezone = '' THEN $timezone ELSE c.timezone END"
		params["timezone"] = data.Timezone
	}
	if data.UpdateProfilePhotoUrl {
		cypher += ", c.profilePhotoUrl = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.profilePhotoUrl is null OR c.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE c.profilePhotoUrl END"
		params["profilePhotoUrl"] = data.ProfilePhotoUrl
	}

	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contactWriteRepository) UpdateTimeProperty(ctx context.Context, tenant, contactId string, property entity.ContactProperty, value *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.UpdateTimeProperty")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)
	span.LogFields(log.String("property", string(property)), log.Object("value", value))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id: $contactId})
			SET c.%s = $value`, string(property))
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
		"value":     utils.TimePtrAsAny(value),
	}
	span.LogFields(log.String("cypher", cypher))
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
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
