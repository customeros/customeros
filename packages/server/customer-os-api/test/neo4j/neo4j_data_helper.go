package neo4j

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"log"
	"time"
)

func CleanupAllData(ctx context.Context, driver *neo4j.DriverWithContext) {
	ExecuteWriteQuery(ctx, driver, `MATCH (n) DETACH DELETE n`, map[string]any{})
}

func CreateFullTextBasicSearchIndexes(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := fmt.Sprintf("DROP INDEX basicSearchStandard_location_terms IF EXISTS")
	ExecuteWriteQuery(ctx, driver, query, map[string]any{})

	query = fmt.Sprintf("CREATE FULLTEXT INDEX basicSearchStandard_location_terms IF NOT EXISTS FOR (n:State) ON EACH [n.name, n.code] " +
		"OPTIONS {  indexConfig: { `fulltext.analyzer`: 'standard', `fulltext.eventually_consistent`: true } }")
	ExecuteWriteQuery(ctx, driver, query, map[string]any{})

	ExecuteWriteQuery(ctx, driver, query, map[string]any{})
}

func CreateTenant(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MERGE (t:Tenant {name:$tenant})`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant": tenant,
	})
}

func CreateWorkspace(ctx context.Context, driver *neo4j.DriverWithContext, workspace string, provider string, tenant string) {
	query := `MATCH (t:Tenant {name: $tenant})
			  MERGE (t)-[:HAS_WORKSPACE]->(w:Workspace {name:$workspace, provider:$provider})`

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":    tenant,
		"provider":  provider,
		"workspace": workspace,
	})
}

func CreateHubspotExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":           tenant,
		"externalSystemId": "hubspot",
	})
}

func CreateCalComExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":           tenant,
		"externalSystemId": "calcom",
	})
}

func LinkWithHubspotExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, entityId, externalId string, externalUrl, externalSource *string, syncDate time.Time) {
	query := `MATCH (e:ExternalSystem {id:$externalSystemId}), (n {id:$entityId})
			MERGE (n)-[rel:IS_LINKED_WITH {externalId:$externalId}]->(e)
			ON CREATE SET rel.externalUrl=$externalUrl, rel.syncDate=$syncDate, rel.externalSource=$externalSource`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"externalSystemId": "hubspot",
		"entityId":         entityId,
		"externalId":       externalId,
		"externalUrl":      externalUrl,
		"syncDate":         syncDate,
		"externalSource":   externalSource,
	})
}

func CreateDefaultUser(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) string {
	return CreateUser(ctx, driver, tenant, entity.UserEntity{
		FirstName:     "first",
		LastName:      "last",
		Source:        "openline",
		SourceOfTruth: "openline",
	})
}

func CreateDefaultUserWithId(ctx context.Context, driver *neo4j.DriverWithContext, tenant, userId string) string {
	return CreateUserWithId(ctx, driver, tenant, userId, entity.UserEntity{
		FirstName:     "first",
		LastName:      "last",
		Source:        "openline",
		SourceOfTruth: "openline",
	})
}

func CreateUser(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, user entity.UserEntity) string {
	return CreateUserWithId(ctx, driver, tenant, "", user)
}

func CreateAttachment(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, attachment entity.AttachmentEntity) string {
	if len(attachment.Id) == 0 {
		attachmentUuid, _ := uuid.NewRandom()
		attachment.Id = attachmentUuid.String()
	}
	query := "MERGE (a:Attachment_%s {id:randomUUID()}) ON CREATE SET " +
		" a:Attachment, " +
		" a.id=$id, " +
		" a.source=$source, " +
		" a.createdAt=datetime({timezone: 'UTC'}), " +
		" a.name=$name, " +
		" a.mimeType=$mimeType, " +
		" a.extension=$extension, " +
		" a.size=$size, " +
		" a.sourceOfTruth=$sourceOfTruth, " +
		" a.appSource=$appSource " +
		" RETURN a"
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"tenant":        tenant,
		"id":            attachment.Id,
		"name":          attachment.Name,
		"mimeType":      attachment.MimeType,
		"size":          attachment.Size,
		"extension":     attachment.Extension,
		"sourceOfTruth": attachment.SourceOfTruth,
		"source":        attachment.Source,
		"appSource":     attachment.AppSource,
	})
	return attachment.Id
}

func CreateUserWithId(ctx context.Context, driver *neo4j.DriverWithContext, tenant, userId string, user entity.UserEntity) string {
	if userId == "" {
		userUuid, _ := uuid.NewRandom()
		userId = userUuid.String()
	}
	query := `
		MATCH (t:Tenant {name:$tenant})
			MERGE (u:User {
				  	id: $userId,
				  	firstName: $firstName,
				  	lastName: $lastName,
					createdAt :datetime({timezone: 'UTC'}),
					source: $source,
					sourceOfTruth: $sourceOfTruth
				})-[:USER_BELONGS_TO_TENANT]->(t)
			SET u:User_%s, 
				u.roles=$roles,
				u.internal=$internal`
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"tenant":        tenant,
		"userId":        userId,
		"firstName":     user.FirstName,
		"lastName":      user.LastName,
		"source":        user.Source,
		"sourceOfTruth": user.SourceOfTruth,
		"roles":         user.Roles,
		"internal":      user.Internal,
	})
	return userId
}

func CreateDefaultPlayer(ctx context.Context, driver *neo4j.DriverWithContext, authId, provider string) string {
	return CreatePlayerWithId(ctx, driver, "", entity.PlayerEntity{
		AuthId:     authId,
		Provider:   provider,
		IdentityId: utils.StringPtr("test-player-id"),
	})
}

func CreatePlayerWithId(ctx context.Context, driver *neo4j.DriverWithContext, playerId string, player entity.PlayerEntity) string {
	if len(playerId) == 0 {
		playerUuid, _ := uuid.NewRandom()
		playerId = playerUuid.String()
	}
	query := `
			MERGE (p:Player {
				  	id: $playerId,
					authId: $authId,
					provider: $provider
				})
			SET     p.identityId = $identityId,
					p.createdAt = datetime({timezone: 'UTC'}),
					p.updatedAt = datetime({timezone: 'UTC'}),
					p.source =  $source,
					p.sourceOfTruth = $sourceOfTruth,
			        p.appSource = $appSource`

	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query), map[string]any{
		"playerId":      playerId,
		"authId":        player.AuthId,
		"provider":      player.Provider,
		"source":        player.Source,
		"sourceOfTruth": player.SourceOfTruth,
		"appSource":     player.AppSource,
		"identityId":    player.IdentityId,
	})

	return playerId
}

func LinkPlayerToUser(ctx context.Context, driver *neo4j.DriverWithContext, playerId string, userId string, isDefault bool) {
	query := `
			MATCH (p:Player {id:$playerId})
			MATCH (u:User {id:$userId})
			MERGE (p)-[:IDENTIFIES {default: $default}]->(u)
			`
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query), map[string]any{
		"playerId": playerId,
		"userId":   userId,
		"default":  isDefault,
	})

}

func CreateDefaultContact(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) string {
	return CreateContact(ctx, driver, tenant, entity.ContactEntity{Prefix: "MR", FirstName: "first", LastName: "last"})
}

func CreateContactWith(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, firstName string, lastName string) string {
	return CreateContact(ctx, driver, tenant, entity.ContactEntity{Prefix: "MR", FirstName: firstName, LastName: lastName})
}

func CreateContact(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, contact entity.ContactEntity) string {
	var contactId = contact.Id
	if contactId == "" {
		contactUuid, _ := uuid.NewRandom()
		contactId = contactUuid.String()
	}
	query := `MATCH (t:Tenant {name: $tenant}) 
		 		MERGE (c:Contact {id: $contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t) 
			 	ON CREATE SET c.prefix=$prefix, 
						c.firstName=$firstName, 
						c.lastName=$lastName, 
						c.name=$name, 
						c.description=$description,
						c.appSource=$appSource, 
						c.source=$source, 
						c.sourceOfTruth=$sourceOfTruth, 
						c.createdAt=$now, 
		 				c:Contact_%s`

	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"tenant":        tenant,
		"contactId":     contactId,
		"prefix":        contact.Prefix,
		"firstName":     contact.FirstName,
		"lastName":      contact.LastName,
		"name":          contact.Name,
		"description":   contact.Description,
		"now":           utils.Now(),
		"source":        contact.Source,
		"sourceOfTruth": contact.SourceOfTruth,
		"appSource":     utils.StringFirstNonEmpty(contact.AppSource, "test"),
	})
	return contactId
}

func CreateContactWithId(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, contactId string, contact entity.ContactEntity) string {
	contact.Id = contactId
	return CreateContact(ctx, driver, tenant, contact)
}

func CreateDefaultFieldSet(ctx context.Context, driver *neo4j.DriverWithContext, contactId string) string {
	return CreateFieldSet(ctx, driver, contactId, entity.FieldSetEntity{Name: "name", Source: entity.DataSourceOpenline, SourceOfTruth: entity.DataSourceOpenline})
}

func CreateFieldSet(ctx context.Context, driver *neo4j.DriverWithContext, contactId string, fieldSet entity.FieldSetEntity) string {
	var fieldSetId, _ = uuid.NewRandom()
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (s:FieldSet {
				  	id: $fieldSetId,
				  	name: $name,
					source: $source,
					sourceOfTruth: $sourceOfTruth,
					createdAt :datetime({timezone: 'UTC'})
				})<-[:HAS_COMPLEX_PROPERTY]-(c)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"contactId":     contactId,
		"fieldSetId":    fieldSetId.String(),
		"name":          fieldSet.Name,
		"source":        fieldSet.Source,
		"sourceOfTruth": fieldSet.SourceOfTruth,
	})
	return fieldSetId.String()
}

func CreateDefaultCustomFieldInSet(ctx context.Context, driver *neo4j.DriverWithContext, fieldSetId string) string {
	return createCustomFieldInSet(ctx, driver, fieldSetId,
		entity.CustomFieldEntity{
			Name:          "name",
			Source:        entity.DataSourceOpenline,
			SourceOfTruth: entity.DataSourceOpenline,
			DataType:      model.CustomFieldDataTypeText.String(),
			Value:         model.AnyTypeValue{Str: utils.StringPtr("value")}})
}

func createCustomFieldInSet(ctx context.Context, driver *neo4j.DriverWithContext, fieldSetId string, customField entity.CustomFieldEntity) string {
	var fieldId, _ = uuid.NewRandom()
	customField.AdjustValueByDatatype()
	query := fmt.Sprintf(
		"MATCH (s:FieldSet {id:$fieldSetId}) "+
			" MERGE (:%s:CustomField { "+
			"	  id: $fieldId, "+
			"	  %s: $value, "+
			"	  datatype: $datatype, "+
			"	  name: $name, "+
			"	  source: $source, "+
			"	  sourceOfTruth: $sourceOfTruth "+
			"	})<-[:HAS_PROPERTY]-(s)", customField.NodeLabel(), customField.PropertyName())
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"fieldSetId":    fieldSetId,
		"fieldId":       fieldId.String(),
		"name":          customField.Name,
		"datatype":      customField.DataType,
		"value":         customField.Value.RealValue(),
		"source":        customField.Source,
		"sourceOfTruth": customField.SourceOfTruth,
	})
	return fieldId.String()
}

func CreateDefaultCustomFieldInContact(ctx context.Context, driver *neo4j.DriverWithContext, contactId string) string {
	return createCustomFieldInContact(ctx, driver, contactId,
		entity.CustomFieldEntity{
			Name:          "name",
			DataType:      model.CustomFieldDataTypeText.String(),
			Source:        entity.DataSourceOpenline,
			SourceOfTruth: entity.DataSourceOpenline,
			Value:         model.AnyTypeValue{Str: utils.StringPtr("value")}})
}

func createCustomFieldInContact(ctx context.Context, driver *neo4j.DriverWithContext, contactId string, customField entity.CustomFieldEntity) string {
	var fieldId, _ = uuid.NewRandom()
	customField.AdjustValueByDatatype()
	query := fmt.Sprintf(
		"MATCH (c:Contact {id:$contactId}) "+
			" MERGE (:%s:CustomField { "+
			"	  id: $fieldId, "+
			"	  %s: $value, "+
			"	  datatype: $datatype, "+
			"	  name: $name, "+
			"	  source: $source, "+
			"	  sourceOfTruth: $sourceOfTruth "+
			"	})<-[:HAS_PROPERTY]-(c)", customField.NodeLabel(), customField.PropertyName())
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"contactId":     contactId,
		"fieldId":       fieldId.String(),
		"name":          customField.Name,
		"datatype":      customField.DataType,
		"value":         customField.Value.RealValue(),
		"source":        customField.Source,
		"sourceOfTruth": customField.SourceOfTruth,
	})
	return fieldId.String()
}

func CreateEmail(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, entity entity.EmailEntity) string {
	if entity.Email == "" && entity.RawEmail == "" {
		log.Fatalf("Missing email address")
	}

	var emailId, _ = uuid.NewRandom()
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
								MERGE (e:Email {id:$emailId})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
								SET e:Email_%s,
									e.email=$email,
									e.rawEmail=$rawEmail,
									e.isReachable=$isReachable,
									e.createdAt=$createdAt,
									e.updatedAt=$updatedAt
							`, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":      tenant,
		"emailId":     emailId.String(),
		"email":       entity.Email,
		"rawEmail":    entity.RawEmail,
		"isReachable": entity.IsReachable,
		"createdAt":   entity.CreatedAt,
		"updatedAt":   entity.UpdatedAt,
	})
	return emailId.String()
}

func AddEmailTo(ctx context.Context, driver *neo4j.DriverWithContext, entityType entity.EntityType, tenant, entityId, email string, primary bool, label string) string {
	query := ""

	switch entityType {
	case entity.CONTACT:
		query = "MATCH (entity:Contact {id:$entityId})--(t:Tenant) "
	case entity.USER:
		query = "MATCH (entity:User {id:$entityId})--(t:Tenant) "
	case entity.ORGANIZATION:
		query = "MATCH (entity:Organization {id:$entityId})--(t:Tenant) "
	}

	var emailId, _ = uuid.NewRandom()
	query = query +
		" MERGE (e:Email {rawEmail: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)" +
		" ON CREATE SET " +
		"	e.rawEmail=$email, " +
		"	e.email=$email, " +
		"	e.id=$emailId, " +
		"	e:%s " +
		" WITH e, entity MERGE (e)<-[rel:HAS]-(entity) " +
		" ON CREATE SET rel.label=$label, rel.primary=$primary "

	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, "Email_"+tenant), map[string]any{
		"entityId": entityId,
		"primary":  primary,
		"email":    email,
		"label":    label,
		"emailId":  emailId.String(),
	})
	return emailId.String()
}

func AddPhoneNumberTo(ctx context.Context, driver *neo4j.DriverWithContext, tenant, id, phoneNumber string, primary bool, label string) string {
	var phoneNumberId, _ = uuid.NewRandom()
	query :=
		" MATCH (n {id:$entityId})--(t:Tenant) " +
			" MERGE (p:PhoneNumber {rawPhoneNumber:$phoneNumber})-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t) " +
			" ON CREATE SET " +
			" 	p.e164=$phoneNumber," +
			" 	p.validated=true," +
			"	p.id=$phoneNumberId, " +
			"	p:%s " +
			" WITH p, n MERGE (p)<-[rel:HAS]-(n) " +
			" ON CREATE SET rel.label=$label, rel.primary=$primary "
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, "PhoneNumber_"+tenant), map[string]any{
		"phoneNumberId": phoneNumberId.String(),
		"entityId":      id,
		"primary":       primary,
		"phoneNumber":   phoneNumber,
		"label":         label,
	})
	return phoneNumberId.String()
}

func CreatePhoneNumber(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, phoneNumber entity.PhoneNumberEntity) string {
	var phoneNumberId, _ = uuid.NewRandom()
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
								MERGE (p:PhoneNumber {id:$phoneNumberId})-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t)
								SET p:PhoneNumber_%s,
									p.e164=$e164,
									p.rawPhoneNumber=$rawPhoneNumber,
									p.createdAt=$createdAt,
									p.updatedAt=$updatedAt`, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":         tenant,
		"phoneNumberId":  phoneNumberId.String(),
		"rawPhoneNumber": phoneNumber.RawPhoneNumber,
		"e164":           phoneNumber.E164,
		"createdAt":      phoneNumber.CreatedAt,
		"updatedAt":      phoneNumber.UpdatedAt,
	})
	return phoneNumberId.String()
}

func CreateEntityTemplate(ctx context.Context, driver *neo4j.DriverWithContext, tenant, extends string) string {
	var templateId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:EntityTemplate {id:$templateId})-[:ENTITY_TEMPLATE_BELONGS_TO_TENANT]->(t)
			ON CREATE SET e.extends=$extends, e.name=$name`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"templateId": templateId.String(),
		"tenant":     tenant,
		"extends":    extends,
		"name":       "template name",
	})
	return templateId.String()
}

func LinkEntityTemplateToContact(ctx context.Context, driver *neo4j.DriverWithContext, entityTemplateId, contactId string) {
	query := `MATCH (c:Contact {id:$contactId}),
			(e:EntityTemplate {id:$TemplateId})
			MERGE (c)-[:IS_DEFINED_BY]->(e)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"TemplateId": entityTemplateId,
		"contactId":  contactId,
	})
}

func AddFieldTemplateToEntity(ctx context.Context, driver *neo4j.DriverWithContext, entityTemplateId string) string {
	var templateId, _ = uuid.NewRandom()
	query := `MATCH (e:EntityTemplate {id:$entityTemplateId})
			MERGE (f:CustomFieldTemplate {id:$templateId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"templateId":       templateId.String(),
		"entityTemplateId": entityTemplateId,
		"type":             "TEXT",
		"order":            1,
		"mandatory":        false,
		"name":             "template name",
	})
	return templateId.String()
}

func AddFieldTemplateToSet(ctx context.Context, driver *neo4j.DriverWithContext, setTemplateId string) string {
	var templateId, _ = uuid.NewRandom()
	query := `MATCH (e:FieldSetTemplate {id:$setTemplateId})
			MERGE (f:CustomFieldTemplate {id:$templateId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"templateId":    templateId.String(),
		"setTemplateId": setTemplateId,
		"type":          "TEXT",
		"order":         1,
		"mandatory":     false,
		"name":          "template name",
	})
	return templateId.String()
}

func AddSetTemplateToEntity(ctx context.Context, driver *neo4j.DriverWithContext, entityTemplateId string) string {
	var templateId, _ = uuid.NewRandom()
	query := `MATCH (e:EntityTemplate {id:$entityTemplateId})
			MERGE (f:FieldSetTemplate {id:$templateId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"templateId":       templateId.String(),
		"entityTemplateId": entityTemplateId,
		"type":             "TEXT",
		"order":            1,
		"mandatory":        false,
		"name":             "set name",
	})
	return templateId.String()
}

func CreateTag(ctx context.Context, driver *neo4j.DriverWithContext, tenant, tagName string) string {
	var tagId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})
			ON CREATE SET tag.name=$name, tag.source=$source, tag.appSource=$appSource, tag.createdAt=$now, tag.updatedAt=$now`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":        tagId.String(),
		"tenant":    tenant,
		"name":      tagName,
		"source":    "openline",
		"appSource": "test",
		"now":       utils.Now(),
	})
	return tagId.String()
}

func CreateIssue(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, issue entity.IssueEntity) string {
	var issueId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$id})
			ON CREATE SET 
				i.subject=$subject, 
				i.createdAt=$createdAt,  
				i.updatedAt=$createdAt,
				i.description=$description,
				i.status=$status,
				i.priority=$priority,
				i.source=$source,
				i.appSource=$appSource,
				i.sourceOfTruth=$sourceOfTruth,
				i:TimelineEvent,
				i:Issue_%s,
				i:TimelineEvent_%s`
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":            issueId.String(),
		"tenant":        tenant,
		"subject":       issue.Subject,
		"createdAt":     issue.CreatedAt,
		"description":   issue.Description,
		"status":        issue.Status,
		"priority":      issue.Priority,
		"source":        "openline",
		"sourceOfTruth": "openline",
		"appSource":     "test",
	})
	return issueId.String()
}

func IssueReportedByOrganization(ctx context.Context, driver *neo4j.DriverWithContext, organizationId, issueId string) {
	query := `MATCH (o:Organization {id:$organizationId}), (i:Issue {id:$issueId})
			MERGE (o)<-[:REPORTED_BY]-(i)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"issueId":        issueId,
		"organizationId": organizationId,
	})
}

func TagIssue(ctx context.Context, driver *neo4j.DriverWithContext, issueId, tagId string) {
	query := `MATCH (i:Issue {id:$issueId}), (tag:Tag {id:$tagId})
			MERGE (i)-[r:TAGGED]->(tag)
			ON CREATE SET r.taggedAt=datetime({timezone: 'UTC'})`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tagId":   tagId,
		"issueId": issueId,
	})
}

func NoteMentionsTag(ctx context.Context, driver *neo4j.DriverWithContext, noteId, tagId string) {
	query := `MATCH (n:Note {id:$noteId}), (tag:Tag {id:$tagId})
			MERGE (n)-[:MENTIONED]->(tag)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tagId":  tagId,
		"noteId": noteId,
	})
}

func TagContact(ctx context.Context, driver *neo4j.DriverWithContext, contactId, tagId string) {
	query := `MATCH (c:Contact {id:$contactId}), (tag:Tag {id:$tagId})
			MERGE (c)-[r:TAGGED]->(tag)
			ON CREATE SET r.taggedAt=datetime({timezone: 'UTC'})`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tagId":     tagId,
		"contactId": contactId,
	})
}

func TagOrganization(ctx context.Context, driver *neo4j.DriverWithContext, organizationId, tagId string) {
	query := `MATCH (o:Organization {id:$organizationId}), (tag:Tag {id:$tagId})
			MERGE (o)-[r:TAGGED]->(tag)
			ON CREATE SET r.taggedAt=datetime({timezone: 'UTC'})`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tagId":          tagId,
		"organizationId": organizationId,
	})
}

func CreateDefaultOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) string {
	return CreateOrg(ctx, driver, tenant, entity.OrganizationEntity{
		Name: "Default org",
	})
}

func CreateOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationName string) string {
	return CreateOrg(ctx, driver, tenant, entity.OrganizationEntity{
		Name: organizationName,
	})
}

func CreateTenantOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationName string) string {
	return CreateOrg(ctx, driver, tenant, entity.OrganizationEntity{
		Name:               organizationName,
		TenantOrganization: true,
	})
}

func LinkOrganizationAsSubsidiary(ctx context.Context, driver *neo4j.DriverWithContext, parentOrganizationId, subOrganizationId, relationType string) {
	query := `MATCH (parent:Organization {id:$parentOrganizationId}),
			(org:Organization {id:$subOrganizationId})
			MERGE (org)-[rel:SUBSIDIARY_OF]->(parent)
			ON CREATE SET rel.type=$type`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"parentOrganizationId": parentOrganizationId,
		"subOrganizationId":    subOrganizationId,
		"type":                 relationType,
	})
}

func LinkSuggestedMerge(ctx context.Context, driver *neo4j.DriverWithContext, orgId, primaryOrgId, suggestedBy string, suggestedAt time.Time, confidence float64) {
	query := `MATCH (primary:Organization {id:$primaryOrgId}),
					(org:Organization {id:$orgId})
			MERGE (org)-[rel:SUGGESTED_MERGE]->(primary)
			ON CREATE SET rel.suggestedBy=$suggestedBy, rel.suggestedAt=$suggestedAt, rel.confidence=$confidence`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"orgId":        orgId,
		"primaryOrgId": primaryOrgId,
		"suggestedBy":  suggestedBy,
		"suggestedAt":  suggestedAt,
		"confidence":   confidence,
	})
}

func CreateOrg(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, organization entity.OrganizationEntity) string {
	var organizationId, _ = uuid.NewRandom()
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$id})
			ON CREATE SET 	org.name=$name, 
							org.description=$description, 
							org.website=$website,
							org.industry=$industry, 
							org.subIndustry=$subIndustry,
							org.industryGroup=$industryGroup,
							org.targetAudience=$targetAudience,	
							org.valueProposition=$valueProposition,
							org.lastFundingRound=$lastFundingRound,
							org.lastFundingAmount=$lastFundingAmount,
							org.isPublic=$isPublic, 
							org.tenantOrganization=$tenantOrganization,
							org.createdAt=$now,
							org.updatedAt=$now`, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":                 organizationId.String(),
		"tenant":             tenant,
		"name":               organization.Name,
		"description":        organization.Description,
		"website":            organization.Website,
		"industry":           organization.Industry,
		"isPublic":           organization.IsPublic,
		"subIndustry":        organization.SubIndustry,
		"industryGroup":      organization.IndustryGroup,
		"targetAudience":     organization.TargetAudience,
		"valueProposition":   organization.ValueProposition,
		"tenantOrganization": organization.TenantOrganization,
		"lastFundingRound":   organization.LastFundingRound,
		"lastFundingAmount":  organization.LastFundingAmount,
		"now":                utils.Now(),
	})
	return organizationId.String()
}

func AddDomainToOrg(ctx context.Context, driver *neo4j.DriverWithContext, organizationId, domain string) {
	query := ` MERGE (d:Domain {domain:$domain})
			ON CREATE SET
				d.id=randomUUID(),
				d.source="test",
				d.appSource="test",
				d.createdAt=$now,
				d.updatedAt=$now
			WITH d
			MATCH (o:Organization {id:$organizationId})
			MERGE (o)-[:HAS_DOMAIN]->(d)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"domain":         domain,
		"now":            utils.Now(),
	})
}

func ContactWorksForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, contactId, organizationId, jobTitle string, primary bool) string {
	var roleId, _ = uuid.NewRandom()
	query := `MATCH (c:Contact {id:$contactId}),
			        (org:Organization {id:$organizationId})
			MERGE (c)-[:WORKS_AS]->(r:JobRole)-[:ROLE_IN]->(org)
			ON CREATE SET r.id=$id, r.jobTitle=$jobTitle, r.primary=$primary, r.responsibilityLevel=$responsibilityLevel,
							r.createdAt=datetime({timezone: 'UTC'}), r.appSource=$appSource`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":                  roleId.String(),
		"contactId":           contactId,
		"organizationId":      organizationId,
		"jobTitle":            jobTitle,
		"primary":             primary,
		"responsibilityLevel": 1,
		"appSource":           "test",
	})
	return roleId.String()
}

func UserWorksAs(ctx context.Context, driver *neo4j.DriverWithContext, userId, jobTitle string, description string, primary bool) string {
	var roleId, _ = uuid.NewRandom()
	query := `MATCH (u:User {id:$userId})
			MERGE (u)-[:WORKS_AS]->(r:JobRole)
			ON CREATE SET r.id=$id, r.description=$description, r.jobTitle=$jobTitle, r.primary=$primary, r.responsibilityLevel=$responsibilityLevel,
							r.createdAt=datetime({timezone: 'UTC'}), r.appSource=$appSource`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":                  roleId.String(),
		"userId":              userId,
		"jobTitle":            jobTitle,
		"description":         description,
		"primary":             primary,
		"responsibilityLevel": 1,
		"appSource":           "test",
	})
	return roleId.String()
}

func UserOwnsContact(ctx context.Context, driver *neo4j.DriverWithContext, userId, contactId string) {
	query := `MATCH (c:Contact {id:$contactId}),
			        (u:User {id:$userId})
			MERGE (u)-[:OWNS]->(c)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"contactId": contactId,
		"userId":    userId,
	})
}

func UserOwnsOrganization(ctx context.Context, driver *neo4j.DriverWithContext, userId, organizationId string) {
	query := `MATCH (o:Organization {id:$organizationId}),
			        (u:User {id:$userId})
			MERGE (u)-[:OWNS]->(o)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"userId":         userId,
	})
}

func CreateConversation(ctx context.Context, driver *neo4j.DriverWithContext, tenant, userId, contactId, subject string, startedAt time.Time) string {
	var conversationId, _ = uuid.NewRandom()
	query := `MATCH (c:Contact {id:$contactId}),
			        (u:User {id:$userId})
			MERGE (u)-[:PARTICIPATES]->(o:Conversation:Conversation_%s:TimelineEvent:TimelineEvent_%s {id:$conversationId})<-[:PARTICIPATES]-(c)
			ON CREATE SET 	o.startedAt=$startedAt, 
							o.status="ACTIVE", 
							o.channel="VOICE",
							o.subject=$subject,
							o.messageCount=0 `
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"contactId":      contactId,
		"userId":         userId,
		"subject":        subject,
		"startedAt":      startedAt,
		"conversationId": conversationId.String(),
	})
	return conversationId.String()
}

func UserHasCalendar(ctx context.Context, driver *neo4j.DriverWithContext, userId, link, calType string, primary bool) string {
	var calId, _ = uuid.NewRandom()
	query := `MATCH (u:User {id:$userId})
			MERGE (u)-[:HAS_CALENDAR]->(c:Calendar)
			ON CREATE SET c.id=$id, c.link=$link, c.calType=$calType, c.primary=$primary, c.createdAt=datetime({timezone: 'UTC'}), c.appSource=$appSource`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":        calId.String(),
		"calType":   calType,
		"userId":    userId,
		"link":      link,
		"primary":   primary,
		"appSource": "test",
	})
	return calId.String()
}

func CreatePageView(ctx context.Context, driver *neo4j.DriverWithContext, contactId string, pageViewEntity entity.PageViewEntity) string {
	var actionId, _ = uuid.NewRandom()
	query := `MATCH (c:Contact {id:$contactId})
			MERGE (c)-[:HAS_ACTION]->(a:TimelineEvent:PageView {id:$actionId})
			ON CREATE SET
				a.trackerName=$trackerName,
				a.startedAt=$startedAt,
				a.endedAt=$endedAt,
				a.application=$application,
				a.pageUrl=$pageUrl,
				a.pageTitle=$pageTitle,
				a.sessionId=$sessionId,
				a.orderInSession=$orderInSession,
				a.engagedTime=$engagedTime,
				a.source=$source,	
				a.sourceOfTruth=$sourceOfTruth,	
				a.appSource=$appSource`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"contactId":      contactId,
		"actionId":       actionId.String(),
		"trackerName":    pageViewEntity.TrackerName,
		"startedAt":      pageViewEntity.StartedAt,
		"endedAt":        pageViewEntity.EndedAt,
		"application":    pageViewEntity.Application,
		"pageUrl":        pageViewEntity.PageUrl,
		"pageTitle":      pageViewEntity.PageTitle,
		"sessionId":      pageViewEntity.SessionId,
		"orderInSession": pageViewEntity.OrderInSession,
		"engagedTime":    pageViewEntity.EngagedTime,
		"source":         "openline",
		"sourceOfTruth":  "openline",
		"appSource":      "test",
	})
	return actionId.String()
}

func CreateLocation(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, location entity.LocationEntity) string {
	var locationId, _ = uuid.NewRandom()
	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (l:Location {id:$locationId})-[:LOCATION_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET l.name=$name, " +
		"				l.source=$source, " +
		"				l.appSource=$appSource, " +
		"				l.createdAt=$now, " +
		"				l.updatedAt=$now, " +
		"				l.country=$country, " +
		"				l.region=$region, " +
		"				l.locality=$locality, " +
		"				l.address=$address, " +
		"				l.address2=$address2, " +
		"				l.zip=$zip, " +
		"				l.addressType=$addressType, " +
		"				l.houseNumber=$houseNumber, " +
		"				l.postalCode=$postalCode, " +
		"				l.plusFour=$plusFour, " +
		"				l.commercial=$commercial, " +
		"				l.predirection=$predirection, " +
		"				l.district=$district, " +
		"				l.street=$street, " +
		"				l.rawAddress=$rawAddress, " +
		"				l.latitude=$latitude, " +
		"				l.longitude=$longitude, " +
		"				l.utcOffset=$utcOffset, " +
		"				l.timeZone=$timeZone, " +
		"				l:Location_%s"

	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"tenant":       tenant,
		"locationId":   locationId.String(),
		"source":       location.Source,
		"appSource":    location.AppSource,
		"name":         location.Name,
		"now":          utils.Now(),
		"country":      location.Country,
		"region":       location.Region,
		"locality":     location.Locality,
		"address":      location.Address,
		"address2":     location.Address2,
		"zip":          location.Zip,
		"addressType":  location.AddressType,
		"houseNumber":  location.HouseNumber,
		"postalCode":   location.PostalCode,
		"plusFour":     location.PlusFour,
		"commercial":   location.Commercial,
		"predirection": location.Predirection,
		"district":     location.District,
		"street":       location.Street,
		"rawAddress":   location.RawAddress,
		"latitude":     location.Latitude,
		"longitude":    location.Longitude,
		"utcOffset":    location.UtcOffset,
		"timeZone":     location.TimeZone,
	})
	return locationId.String()
}

func ContactAssociatedWithLocation(ctx context.Context, driver *neo4j.DriverWithContext, contactId, locationId string) {
	query := `MATCH (c:Contact {id:$contactId}),
			        (l:Location {id:$locationId})
			MERGE (c)-[:ASSOCIATED_WITH]->(l)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"contactId":  contactId,
		"locationId": locationId,
	})
}

func OrganizationAssociatedWithLocation(ctx context.Context, driver *neo4j.DriverWithContext, organizationId, locationId string) {
	query := `MATCH (org:Organization {id:$organizationId}),
			        (l:Location {id:$locationId})
			MERGE (org)-[:ASSOCIATED_WITH]->(l)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"locationId":     locationId,
	})
}

func CreateNoteForContact(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contactId, html string, createdAt time.Time) string {
	var noteId, _ = uuid.NewRandom()

	query := "MATCH (c:Contact {id:$contactId}) " +
		"		MERGE (c)-[:NOTED]->(n:Note {id:$id}) " +
		"		ON CREATE SET 	n.html=$html, " +
		"						n.createdAt=$createdAt, " +
		"						n.updatedAt=$createdAt, " +
		"						n.source=$source, " +
		"						n.sourceOfSource=$source, " +
		"						n.appSource=$appSource, " +
		"						n:Note_%s, " +
		"						n:TimelineEvent, " +
		"						n:TimelineEvent_%s"
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":        noteId.String(),
		"contactId": contactId,
		"html":      html,
		"createdAt": createdAt,
		"source":    "openline",
		"appSource": "test",
	})
	return noteId.String()
}

func CreateNoteForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId, html string, createdAt time.Time) string {
	var noteId, _ = uuid.NewRandom()

	query := "MATCH (org:Organization {id:$organizationId}) " +
		"		MERGE (org)-[:NOTED]->(n:Note {id:$id}) " +
		"		ON CREATE SET 	n.html=$html, " +
		"						n.createdAt=$createdAt, " +
		"						n:Note_%s, " +
		"						n:TimelineEvent, " +
		"						n:TimelineEvent_%s"
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":             noteId.String(),
		"organizationId": organizationId,
		"html":           html,
		"createdAt":      createdAt,
		"source":         "openline",
		"appSource":      "test",
	})
	return noteId.String()
}

func LinkNoteWithOrganization(ctx context.Context, driver *neo4j.DriverWithContext, noteId, organizationId string) {
	query := `MATCH (n:Note {id:$noteId}),
			(org:Organization {id:$organizationId})
			MERGE (n)<-[:NOTED]-(org)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"noteId":         noteId,
	})
}

func NoteCreatedByUser(ctx context.Context, driver *neo4j.DriverWithContext, noteId, userId string) {
	query := `MATCH (u:User {id:$userId})
				MATCH (n:Note {id:$noteId})
			MERGE (u)-[:CREATED]->(n)
`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"noteId": noteId,
		"userId": userId,
	})
}

func LinkContactWithOrganization(ctx context.Context, driver *neo4j.DriverWithContext, contactId, organizationId string) {
	query := `MATCH (c:Contact {id:$contactId}),
			(org:Organization {id:$organizationId})
			MERGE (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(org)
			ON CREATE SET 	j.id=randomUUID(), 
							j.createdAt=$now,
							j.updatedAt=$now,
							j.source=$source,
							j.sourceOfSource=$source,
							j.appSource=$source`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"contactId":      contactId,
		"source":         "test",
		"now":            utils.Now(),
	})
}

func CreateAnalysis(ctx context.Context, driver *neo4j.DriverWithContext, tenant, content, contentType, analysisType string, createdAt time.Time) string {
	var analysisId, _ = uuid.NewRandom()

	query := "MERGE (a:Analysis {id:$id})" +
		" ON CREATE SET " +
		"	a.content=$content, " +
		"	a.createdAt=$createdAt, " +
		"	a.analysisType=$analysisType, " +
		"	a.contentType=$contentType, " +
		"	a.source=$source, " +
		"	a.sourceOfTruth=$sourceOfTruth, " +
		"	a.appSource=$appSource," +
		"	a:Analysis_%s"
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"id":            analysisId.String(),
		"content":       content,
		"contentType":   contentType,
		"analysisType":  analysisType,
		"createdAt":     createdAt,
		"source":        "openline",
		"sourceOfTruth": "openline",
		"appSource":     "test",
	})
	return analysisId.String()
}

func AnalysisDescribes(ctx context.Context, driver *neo4j.DriverWithContext, tenant, actionId, nodeId string, describesType string) {
	query := "MATCH (a:Analysis_%s {id:$actionId}), " +
		"(n:%s_%s {id:$nodeId}) " +
		" MERGE (a)-[:DESCRIBES]->(n) "
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, describesType, tenant), map[string]any{
		"actionId": actionId,
		"nodeId":   nodeId,
	})
}

func CreateInteractionEvent(ctx context.Context, driver *neo4j.DriverWithContext, tenant, identifier, content, contentType string, channel *string, createdAt time.Time) string {
	var interactionEventId, _ = uuid.NewRandom()

	query := "MERGE (ie:InteractionEvent {id:$id})" +
		" ON CREATE SET " +
		"	ie.content=$content, " +
		"	ie.createdAt=$createdAt, " +
		"	ie.channel=$channel, " +
		"	ie.contentType=$contentType, " +
		"	ie.source=$source, " +
		"	ie.sourceOfTruth=$sourceOfTruth, " +
		"	ie.appSource=$appSource," +
		"	ie:InteractionEvent_%s, ie:TimelineEvent, ie:TimelineEvent_%s," +
		"   ie.identifier=$identifier"
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":            interactionEventId.String(),
		"content":       content,
		"contentType":   contentType,
		"channel":       channel,
		"createdAt":     createdAt,
		"source":        "openline",
		"sourceOfTruth": "openline",
		"appSource":     "test",
		"identifier":    identifier,
	})
	return interactionEventId.String()
}

func CreateInteractionSession(ctx context.Context, driver *neo4j.DriverWithContext, tenant, identifier, name, sessionType, status, channel string, createdAt time.Time, inTimeline bool) string {
	var interactionSessionId, _ = uuid.NewRandom()

	query := "MERGE (is:InteractionSession {id:$id})" +
		" ON CREATE SET " +
		"	is.createdAt=$createdAt, " +
		"	is.updatedAt=$updatedAt, " +
		"	is.name=$name, " +
		"	is.type=$type, " +
		"	is.channel=$channel, " +
		"	is.status=$status, " +
		"	is.source=$source, " +
		"	is.sourceOfTruth=$sourceOfTruth, " +
		"	is.appSource=$appSource," +
		"   is.identifier=$identifier, " +
		"	is:InteractionSession_%s"

	resolvedQuery := ""
	if inTimeline {
		query += ", is:TimelineEvent, is:TimelineEvent_%s"

		resolvedQuery = fmt.Sprintf(query, tenant, tenant)
	} else {
		resolvedQuery = fmt.Sprintf(query, tenant)
	}
	ExecuteWriteQuery(ctx, driver, resolvedQuery, map[string]any{
		"id":            interactionSessionId.String(),
		"name":          name,
		"type":          sessionType,
		"channel":       channel,
		"status":        status,
		"createdAt":     createdAt,
		"updatedAt":     createdAt.Add(time.Duration(10) * time.Minute),
		"source":        "openline",
		"sourceOfTruth": "openline",
		"appSource":     "test",
		"identifier":    identifier,
	})
	return interactionSessionId.String()
}

func CreateActionItemLinkedWith(ctx context.Context, driver *neo4j.DriverWithContext, tenant, linkedWith string, linkedWithId, content string, createdAt time.Time) string {
	var actionItemId, _ = uuid.NewRandom()

	session := utils.NewNeo4jWriteSession(ctx, *driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (i:%s_%s{id:$linkedWithId}) `, linkedWith, tenant)
	query += fmt.Sprintf(`MERGE (i)-[r:INCLUDES]->(a:ActionItem_%s{id:$actionItemId}) `, tenant)
	query += fmt.Sprintf("ON CREATE SET " +
		" a:ActionItem, " +
		" a.createdAt=$createdAt, " +
		" a.content=$content, " +
		" a.source=$source, " +
		" a.sourceOfTruth=$sourceOfTruth, " +
		" a.appSource=$appSource ")

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"linkedWithId":  linkedWithId,
		"actionItemId":  actionItemId.String(),
		"content":       content,
		"createdAt":     createdAt,
		"source":        "openline",
		"sourceOfTruth": "openline",
		"appSource":     "test",
	})
	return actionItemId.String()
}

func CreateMeeting(ctx context.Context, driver *neo4j.DriverWithContext, tenant, name string, createdAt time.Time) string {
	var meetingId, _ = uuid.NewRandom()

	query := "MERGE (m:Meeting_%s {id:$id}) " +
		" ON CREATE SET m:Meeting, " +
		"				m.name=$name, " +
		"				m.createdAt=$createdAt, " +
		"				m.updatedAt=$updatedAt, " +
		"				m.start=$createdAt, " +
		"				m.end=$updatedAt, " +
		"				m.appSource=$appSource, " +
		"				m.source=$source, " +
		"				m.sourceOfTruth=$sourceOfTruth, " +
		"				m:TimelineEvent, " +
		"				m:TimelineEvent_%s " +
		" RETURN m"

	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":            meetingId.String(),
		"name":          name,
		"createdAt":     createdAt,
		"updatedAt":     createdAt,
		"source":        "openline",
		"sourceOfTruth": "openline",
		"appSource":     "test",
	})
	return meetingId.String()
}

func InteractionSessionAttendedBy(ctx context.Context, driver *neo4j.DriverWithContext, tenant, interactionSessionId, nodeId, interactionType string) {
	query := "MATCH (is:InteractionSession_%s {id:$interactionSessionId}), " +
		"(n {id:$nodeId}) " +
		" MERGE (is)-[:ATTENDED_BY {type:$interactionType}]->(n) "
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"interactionSessionId": interactionSessionId,
		"nodeId":               nodeId,
		"interactionType":      interactionType,
	})
}

func InteractionEventSentBy(ctx context.Context, driver *neo4j.DriverWithContext, interactionEventId, nodeId, interactionType string) {
	query := "MATCH (ie:InteractionEvent {id:$interactionEventId}), " +
		"(n {id:$nodeId}) " +
		" MERGE (ie)-[:SENT_BY {type:$interactionType}]->(n) "
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"interactionEventId": interactionEventId,
		"nodeId":             nodeId,
		"interactionType":    interactionType,
	})
}

func MeetingCreatedBy(ctx context.Context, driver *neo4j.DriverWithContext, meetingId, nodeId string) {
	query := "MATCH (m:Meeting {id:$meetingId}), " +
		"(n {id:$nodeId}) " +
		" MERGE (m)-[:CREATED_BY]->(n) "
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"meetingId": meetingId,
		"nodeId":    nodeId,
	})
}

func MeetingAttendedBy(ctx context.Context, driver *neo4j.DriverWithContext, meetingId, nodeId string) {
	query := "MATCH (m:Meeting {id:$meetingId}), " +
		"(n {id:$nodeId}) " +
		" MERGE (m)-[:ATTENDED_BY]->(n) "
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"meetingId": meetingId,
		"nodeId":    nodeId,
	})
}

func InteractionEventSentTo(ctx context.Context, driver *neo4j.DriverWithContext, interactionEventId, nodeId, interactionType string) {
	query := "MATCH (ie:InteractionEvent {id:$interactionEventId}), " +
		"(n {id:$nodeId}) " +
		" MERGE (ie)-[:SENT_TO {type:$interactionType}]->(n) "
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"interactionEventId": interactionEventId,
		"nodeId":             nodeId,
		"interactionType":    interactionType,
	})
}

func InteractionEventPartOfInteractionSession(ctx context.Context, driver *neo4j.DriverWithContext, interactionEventId, interactionSessionId string) {
	query := "MATCH (ie:InteractionEvent {id:$interactionEventId}), " +
		"(is:InteractionSession {id:$interactionSessionId}) " +
		" MERGE (ie)-[:PART_OF]->(is) "
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"interactionEventId":   interactionEventId,
		"interactionSessionId": interactionSessionId,
	})
}

func InteractionEventPartOfMeeting(ctx context.Context, driver *neo4j.DriverWithContext, interactionEventId, meetingId string) {
	query := "MATCH (ie:InteractionEvent {id:$interactionEventId}), " +
		"(m:Meeting {id:$meetingId}) " +
		" MERGE (ie)-[:PART_OF]->(m) "
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"interactionEventId": interactionEventId,
		"meetingId":          meetingId,
	})
}

func InteractionEventPartOfIssue(ctx context.Context, driver *neo4j.DriverWithContext, interactionEventId, issueId string) {
	query := "MATCH (ie:InteractionEvent {id:$interactionEventId}), " +
		"(i:Issue {id:$issueId}) " +
		" MERGE (ie)-[:PART_OF]->(i) "
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"interactionEventId": interactionEventId,
		"issueId":            issueId,
	})
}

func InteractionEventRepliesToInteractionEvent(ctx context.Context, driver *neo4j.DriverWithContext, tenant, interactionEventId, repliesToInteractionEventId string) {
	query := "MATCH (ie:InteractionEvent_%s {id:$interactionEventId}), " +
		"(rie:InteractionEvent_%s {id:$repliesToInteractionEventId}) " +
		" MERGE (ie)-[:REPLIES_TO]->(rie) "
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"interactionEventId":          interactionEventId,
		"repliesToInteractionEventId": repliesToInteractionEventId,
	})
}

func CreateCountry(ctx context.Context, driver *neo4j.DriverWithContext, codeA2, codeA3, name, phoneCode string) {
	query := `MERGE (c:Country{codeA3: $codeA3}) 
				ON CREATE SET 
					c.phoneCode = $phoneCode,
					c.codeA2 = $codeA2,
					c.name = $name, 
					c.createdAt = $now, 
					c.updatedAt = $now`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"codeA2":    codeA2,
		"codeA3":    codeA3,
		"phoneCode": phoneCode,
		"name":      name,
		"now":       utils.Now(),
	})
}

func CreateCountryWith(ctx context.Context, driver *neo4j.DriverWithContext, id, countryCodeA3, name string) {
	var countryId = id
	if countryId == "" {
		countryUuid, _ := uuid.NewRandom()
		countryId = countryUuid.String()
	}
	query := "MERGE (c:Country{codeA3: $countryCodeA3}) ON CREATE SET c.id = $countryId, c.name = $name, c.createdAt = datetime({timezone: 'UTC'}), c.updatedAt = datetime({timezone: 'UTC'})"
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"countryId":     countryId,
		"countryCodeA3": countryCodeA3,
		"name":          name,
	})
}

func CreateState(ctx context.Context, driver *neo4j.DriverWithContext, countryCodeA3, name, code string) {
	query := "MATCH (c:Country{codeA3: $countryCodeA3}) MERGE (c)<-[:BELONGS_TO_COUNTRY]-(az:State { code: $code }) ON CREATE SET az.id = randomUUID(), az.name = $name, az.createdAt = datetime({timezone: 'UTC'}), az.updatedAt = datetime({timezone: 'UTC'})"
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"countryCodeA3": countryCodeA3,
		"name":          name,
		"code":          code,
	})
}

func CreateSocial(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, social entity.SocialEntity) string {
	var socialId, _ = uuid.NewRandom()
	query := " MERGE (s:Social {id:$id}) " +
		" ON CREATE SET s.platformName=$platformName, " +
		"				s.url=$url, " +
		"				s.source=$source, " +
		"				s.sourceOfTruth=$source, " +
		"				s.appSource=$appSource, " +
		"				s.createdAt=$now, " +
		"				s.updatedAt=$now, " +
		"				s:Social_%s"

	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"tenant":       tenant,
		"id":           socialId.String(),
		"source":       entity.DataSourceOpenline,
		"appSource":    "test",
		"platformName": social.PlatformName,
		"url":          social.Url,
		"now":          utils.Now(),
	})
	return socialId.String()
}

func LinkSocialWithEntity(ctx context.Context, driver *neo4j.DriverWithContext, entityId, socialId string) {
	query := `MATCH (e {id:$entityId}), (s:Social {id:$socialId}) MERGE (e)-[:HAS]->(s)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"entityId": entityId,
		"socialId": socialId,
	})
}

func CreateOrganizationRelationship(ctx context.Context, driver *neo4j.DriverWithContext, name string) {
	query := `MERGE (r:OrganizationRelationship {name:$name}) ON CREATE SET r.id=randomUUID()`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"name": name,
	})
}

func CreateOrganizationRelationshipStages(ctx context.Context, driver *neo4j.DriverWithContext, tenant, relationship string, stages []string) {
	query := `WITH $stages as stages
			UNWIND stages as stage
			MATCH (t:Tenant {name:$tenant}), (r:OrganizationRelationship {name:$relationship})
			MERGE (t)<-[:STAGE_BELONGS_TO_TENANT]-(s:OrganizationRelationshipStage {name:stage})<-[:HAS_STAGE]-(r)
			ON CREATE SET s.id=randomUUID(), s:OrganizationRelationshipStage_%s`
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"relationship": relationship,
		"tenant":       tenant,
		"stages":       stages,
	})

	for k, stage := range stages {
		query = `MATCH (:Tenant {name:$tenant})<-[:STAGE_BELONGS_TO_TENANT]-(s:OrganizationRelationshipStage {name:$stage})<-[:HAS_STAGE]-(:OrganizationRelationship {name:$relationship})
			SET s.order=$order`
		ExecuteWriteQuery(ctx, driver, query, map[string]any{
			"relationship": relationship,
			"tenant":       tenant,
			"stage":        stage,
			"order":        k,
		})
	}
}

func LinkOrganizationWithRelationship(ctx context.Context, driver *neo4j.DriverWithContext, organizationId, relationship string) {
	query := `MATCH (org {id:$organizationId}), (or:OrganizationRelationship {name:$relationship}) MERGE (org)-[:IS]->(or)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"relationship":   relationship,
	})
}

func LinkOrganizationWithRelationshipAndStage(ctx context.Context, driver *neo4j.DriverWithContext, organizationId, relationship, stage string) {
	query := `MATCH (org {id:$organizationId}), 
				(or:OrganizationRelationship {name:$relationship}),
				(or)-[:HAS_STAGE]->(s:OrganizationRelationshipStage {name:$stage}) 
				MERGE (s)<-[:HAS_STAGE]-(org)-[:IS]->(or)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"relationship":   relationship,
		"stage":          stage,
	})
}

func CreateHealthIndicator(ctx context.Context, driver *neo4j.DriverWithContext, tenant, name string, order int64) string {
	var id, _ = uuid.NewRandom()
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:HEALTH_INDICATOR_BELONGS_TO_TENANT]-(h:HealthIndicator {id:$id})
			ON CREATE SET h:HealthIndicator_%s, h.name=$name, h.order=$order, h.createdAt=$now, h.updatedAt=$now`, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":     id.String(),
		"tenant": tenant,
		"name":   name,
		"order":  order,
		"now":    utils.Now(),
	})
	return id.String()
}

func SetHealthIndicatorForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, organizationId, healthIndicatorId string) {
	query := `MATCH (o:Organization {id:$organizationId}),
			        (h:HealthIndicator {id:$healthIndicatorId})
			MERGE (o)-[:HAS_INDICATOR]->(h)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId":    organizationId,
		"healthIndicatorId": healthIndicatorId,
	})
}

func CreateActionForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId string, actionType entity.ActionType, createdAt time.Time) string {
	var actionId, _ = uuid.NewRandom()

	query := "MATCH (o:Organization {id:$organizationId}) " +
		"		MERGE (o)<-[:ACTION_ON]-(a:Action {id:$id}) " +
		"		ON CREATE SET 	a.type=$type, " +
		"						a.createdAt=$createdAt, " +
		"						a.updatedAt=$createdAt, " +
		"						a.source=$source, " +
		"						a.appSource=$appSource, " +
		"						a:Action_%s, " +
		"						a:TimelineEvent, " +
		"						a:TimelineEvent_%s"
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":             actionId.String(),
		"organizationId": organizationId,
		"type":           actionType,
		"createdAt":      createdAt,
		"source":         "openline",
		"appSource":      "test",
	})
	return actionId.String()
}

func GetCountOfNodes(ctx context.Context, driver *neo4j.DriverWithContext, nodeLabel string) int {
	query := fmt.Sprintf(`MATCH (n:%s) RETURN count(n)`, nodeLabel)
	result := ExecuteReadQueryWithSingleReturn(ctx, driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetCountOfRelationships(ctx context.Context, driver *neo4j.DriverWithContext, relationship string) int {
	query := fmt.Sprintf(`MATCH (a)-[r:%s]-(b) RETURN count(distinct r)`, relationship)
	result := ExecuteReadQueryWithSingleReturn(ctx, driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetCountOfRelationshipsForNodeWithId(ctx context.Context, driver *neo4j.DriverWithContext, relationship, id string) int {
	query := fmt.Sprintf(`MATCH (a {id:$id})-[r:%s]-(b) RETURN count(distinct r)`, relationship)
	result := ExecuteReadQueryWithSingleReturn(ctx, driver, query, map[string]any{
		"id": id,
	})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetTotalCountOfNodes(ctx context.Context, driver *neo4j.DriverWithContext) int {
	query := `MATCH (n) RETURN count(n)`
	result := ExecuteReadQueryWithSingleReturn(ctx, driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetAllLabels(ctx context.Context, driver *neo4j.DriverWithContext) []string {
	query := `MATCH (n) RETURN DISTINCT labels(n)`
	dbRecords := ExecuteReadQueryWithCollectionReturn(ctx, driver, query, map[string]any{})
	labels := []string{}
	for _, v := range dbRecords {
		for _, nodeLabels := range v.Values {
			for _, label := range nodeLabels.([]interface{}) {
				if !contains(labels, label.(string)) {
					labels = append(labels, label.(string))
				}
			}
		}
	}
	return labels
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
