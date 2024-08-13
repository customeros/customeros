package neo4j

import (
	"context"
	"fmt"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func CreateFullTextBasicSearchIndexes(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := fmt.Sprintf("DROP INDEX basicSearchStandard_location_terms IF EXISTS")
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{})

	query = fmt.Sprintf("CREATE FULLTEXT INDEX basicSearchStandard_location_terms IF NOT EXISTS FOR (n:State) ON EACH [n.name, n.code] " +
		"OPTIONS {  indexConfig: { `fulltext.analyzer`: 'standard', `fulltext.eventually_consistent`: true } }")
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{})

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{})
}

func CreateHubspotExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":           tenant,
		"externalSystemId": string(neo4jenum.Hubspot),
	})
}

func CreateSlackExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)
			SET e.externalSource=$externalSource`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":           tenant,
		"externalSystemId": string(neo4jenum.Slack),
		"externalSource":   "Slack",
	})
}

func CreateCalComExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":           tenant,
		"externalSystemId": "calcom",
	})
}

func LinkWithHubspotExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, entityId, externalId string, externalUrl, externalSource *string, syncDate time.Time) {
	LinkWithExternalSystem(ctx, driver, entityId, externalId, string(neo4jenum.Hubspot), externalUrl, externalSource, syncDate)
}

func LinkWithSlackExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, entityId, externalId string, externalUrl, externalSource *string, syncDate time.Time) {
	LinkWithExternalSystem(ctx, driver, entityId, externalId, string(neo4jenum.Slack), externalUrl, externalSource, syncDate)
}

func LinkWithExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, entityId, externalId, externalSystemId string, externalUrl, externalSource *string, syncDate time.Time) {
	query := `MATCH (e:ExternalSystem {id:$externalSystemId}), (n {id:$entityId})
			MERGE (n)-[rel:IS_LINKED_WITH {externalId:$externalId}]->(e)
			ON CREATE SET rel.externalUrl=$externalUrl, rel.syncDate=$syncDate, rel.externalSource=$externalSource`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"externalSystemId": externalSystemId,
		"entityId":         entityId,
		"externalId":       externalId,
		"externalUrl":      externalUrl,
		"syncDate":         syncDate,
		"externalSource":   externalSource,
	})
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
		" a.fileName=$fileName, " +
		" a.mimeType=$mimeType, " +
		" a.cdnUrl=$cdnUrl, " +
		" a.basePath=$basePath, " +
		" a.sourceOfTruth=$sourceOfTruth, " +
		" a.appSource=$appSource " +
		" RETURN a"
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"tenant":        tenant,
		"id":            attachment.Id,
		"fileName":      attachment.FileName,
		"mimeType":      attachment.MimeType,
		"cdnUrl":        attachment.CdnUrl,
		"basePath":      attachment.BasePath,
		"sourceOfTruth": attachment.SourceOfTruth,
		"source":        attachment.Source,
		"appSource":     attachment.AppSource,
	})
	return attachment.Id
}

func CreateDefaultPlayer(ctx context.Context, driver *neo4j.DriverWithContext, authId, provider string) string {
	return CreatePlayerWithId(ctx, driver, "", neo4jentity.PlayerEntity{
		AuthId:     authId,
		Provider:   provider,
		IdentityId: "test-player-id",
	})
}

func CreatePlayerWithId(ctx context.Context, driver *neo4j.DriverWithContext, playerId string, player neo4jentity.PlayerEntity) string {
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

	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query), map[string]any{
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query), map[string]any{
		"playerId": playerId,
		"userId":   userId,
		"default":  isDefault,
	})

}

func CreateDefaultContact(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) string {
	return CreateContact(ctx, driver, tenant, neo4jentity.ContactEntity{Prefix: "MR", FirstName: "first", LastName: "last"})
}

func CreateContactWith(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, firstName string, lastName string) string {
	return CreateContact(ctx, driver, tenant, neo4jentity.ContactEntity{Prefix: "MR", FirstName: firstName, LastName: lastName})
}

func CreateContact(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, contact neo4jentity.ContactEntity) string {
	contactId := utils.NewUUIDIfEmpty(contact.Id)
	query := `MATCH (t:Tenant {name: $tenant}) 
		 		MERGE (c:Contact {id: $contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t) 
			 	ON CREATE SET c.prefix=$prefix, 
						c.firstName=$firstName, 
						c.lastName=$lastName, 
						c.name=$name, 
						c.description=$description,
						c.profilePhotoUrl=$profilePhotoUrl,
						c.username=$username,
						c.appSource=$appSource, 
						c.source=$source, 
						c.sourceOfTruth=$sourceOfTruth, 
						c.createdAt=datetime(), 
		 				c:Contact_%s`

	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"tenant":          tenant,
		"contactId":       contactId,
		"prefix":          contact.Prefix,
		"firstName":       contact.FirstName,
		"lastName":        contact.LastName,
		"name":            contact.Name,
		"description":     contact.Description,
		"profilePhotoUrl": contact.ProfilePhotoUrl,
		"username":        contact.Username,
		"source":          contact.Source,
		"sourceOfTruth":   contact.SourceOfTruth,
		"appSource":       utils.StringFirstNonEmpty(contact.AppSource, "test"),
	})
	return contactId
}

func CreateContactWithId(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, contactId string, contact neo4jentity.ContactEntity) string {
	contact.Id = contactId
	return CreateContact(ctx, driver, tenant, contact)
}

func CreateDefaultFieldSet(ctx context.Context, driver *neo4j.DriverWithContext, contactId string) string {
	return CreateFieldSet(ctx, driver, contactId, entity.FieldSetEntity{Name: "name", Source: neo4jentity.DataSourceOpenline, SourceOfTruth: neo4jentity.DataSourceOpenline})
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
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
			Source:        neo4jentity.DataSourceOpenline,
			SourceOfTruth: neo4jentity.DataSourceOpenline,
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
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
			Source:        neo4jentity.DataSourceOpenline,
			SourceOfTruth: neo4jentity.DataSourceOpenline,
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
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

// Deprecated
func AddEmailTo(ctx context.Context, driver *neo4j.DriverWithContext, entityType commonModel.EntityType, tenant, entityId, email string, primary bool, label string) string {
	query := ""

	switch entityType {
	case commonModel.CONTACT:
		query = "MATCH (entity:Contact {id:$entityId})--(t:Tenant) "
	case commonModel.USER:
		query = "MATCH (entity:User {id:$entityId})--(t:Tenant) "
	case commonModel.ORGANIZATION:
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

	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, "Email_"+tenant), map[string]any{
		"entityId": entityId,
		"primary":  primary,
		"email":    email,
		"label":    label,
		"emailId":  emailId.String(),
	})
	return emailId.String()
}

func LinkEmail(ctx context.Context, driver *neo4j.DriverWithContext, entityId, emailId string, primary bool, label string) {
	query :=
		`	MATCH (n {id:$entityId})--(t:Tenant) 
			MATCH (e:Email {id: $emailId})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
		 	MERGE (e)<-[rel:HAS]-(n) 
		 	ON CREATE SET rel.label=$label, rel.primary=$primary `

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"entityId": entityId,
		"primary":  primary,
		"emailId":  emailId,
		"label":    label,
	})
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, "PhoneNumber_"+tenant), map[string]any{
		"phoneNumberId": phoneNumberId.String(),
		"entityId":      id,
		"primary":       primary,
		"phoneNumber":   phoneNumber,
		"label":         label,
	})
	return phoneNumberId.String()
}

func LinkPhoneNumber(ctx context.Context, driver *neo4j.DriverWithContext, id, phoneNumberId string, primary bool, label string) {
	query :=
		` 	MATCH (n {id:$entityId})--(t:Tenant) 
			MERGE (p:PhoneNumber {id:$phoneNumberId})-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t) 
			MERGE (p)<-[rel:HAS]-(n) 
			ON CREATE SET rel.label=$label, rel.primary=$primary `
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"entityId":      id,
		"primary":       primary,
		"phoneNumberId": phoneNumberId,
		"label":         label,
	})
}

func CreateEntityTemplate(ctx context.Context, driver *neo4j.DriverWithContext, tenant, extends string) string {
	var templateId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:EntityTemplate {id:$templateId})-[:ENTITY_TEMPLATE_BELONGS_TO_TENANT]->(t)
			ON CREATE SET e.extends=$extends, e.name=$name`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"TemplateId": entityTemplateId,
		"contactId":  contactId,
	})
}

func AddFieldTemplateToEntity(ctx context.Context, driver *neo4j.DriverWithContext, entityTemplateId string) string {
	var templateId, _ = uuid.NewRandom()
	query := `MATCH (e:EntityTemplate {id:$entityTemplateId})
			MERGE (f:CustomFieldTemplate {id:$templateId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"templateId":       templateId.String(),
		"entityTemplateId": entityTemplateId,
		"type":             "TEXT",
		"order":            1,
		"mandatory":        false,
		"name":             "set name",
	})
	return templateId.String()
}

// Deprecated
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
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

// Deprecated
func IssueReportedBy(ctx context.Context, driver *neo4j.DriverWithContext, issueId, entityId string) {
	query := `MATCH (e:Organization|User|Contact {id:$entityId}), (i:Issue {id:$issueId})
			MERGE (e)<-[:REPORTED_BY]-(i)`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"issueId":  issueId,
		"entityId": entityId,
	})
}

func IssueSubmittedBy(ctx context.Context, driver *neo4j.DriverWithContext, issueId, entityId string) {
	query := `MATCH (e:Organization|User|Contact {id:$entityId}), (i:Issue {id:$issueId})
			MERGE (e)<-[:SUBMITTED_BY]-(i)`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"issueId":  issueId,
		"entityId": entityId,
	})
}

// Deprecated
func IssueFollowedBy(ctx context.Context, driver *neo4j.DriverWithContext, issueId, entityId string) {
	query := `MATCH (e:Organization|User|Contact {id:$entityId}), (i:Issue {id:$issueId})
			MERGE (e)<-[:FOLLOWED_BY]-(i)`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"issueId":  issueId,
		"entityId": entityId,
	})
}

func IssueAssignedTo(ctx context.Context, driver *neo4j.DriverWithContext, issueId, entityId string) {
	query := `MATCH (e:Organization|User|Contact {id:$entityId}), (i:Issue {id:$issueId})
			MERGE (e)<-[:ASSIGNED_TO]-(i)`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"issueId":  issueId,
		"entityId": entityId,
	})
}

func TagIssue(ctx context.Context, driver *neo4j.DriverWithContext, issueId, tagId string) {
	query := `MATCH (i:Issue {id:$issueId}), (tag:Tag {id:$tagId})
			MERGE (i)-[r:TAGGED]->(tag)
			ON CREATE SET r.taggedAt=datetime({timezone: 'UTC'})`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tagId":   tagId,
		"issueId": issueId,
	})
}

func TagContact(ctx context.Context, driver *neo4j.DriverWithContext, contactId, tagId string) {
	query := `MATCH (c:Contact {id:$contactId}), (tag:Tag {id:$tagId})
			MERGE (c)-[r:TAGGED]->(tag)
			ON CREATE SET r.taggedAt=datetime({timezone: 'UTC'})`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tagId":     tagId,
		"contactId": contactId,
	})
}

func TagLogEntry(ctx context.Context, driver *neo4j.DriverWithContext, logEntryId, tagId string, taggedAt *time.Time) {
	query := `MATCH (l:LogEntry {id:$logEntryId}), (tag:Tag {id:$tagId})
			MERGE (l)-[r:TAGGED]->(tag)
			ON CREATE SET r.taggedAt=$taggedAt`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tagId":      tagId,
		"logEntryId": logEntryId,
		"taggedAt":   utils.TimePtrAsAny(taggedAt, utils.NowPtr()),
	})
}

func TagOrganization(ctx context.Context, driver *neo4j.DriverWithContext, organizationId, tagId string) {
	query := `MATCH (o:Organization {id:$organizationId}), (tag:Tag {id:$tagId})
			MERGE (o)-[r:TAGGED]->(tag)
			ON CREATE SET r.taggedAt=datetime({timezone: 'UTC'})`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tagId":          tagId,
		"organizationId": organizationId,
	})
}

// Deprecated, use CreateOrg
func CreateOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationName string) string {
	return neo4jtest.CreateOrganization(ctx, driver, tenant, neo4jentity.OrganizationEntity{
		Name: organizationName,
	})
}

func CreateTenantOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationName string) string {
	return neo4jtest.CreateOrganization(ctx, driver, tenant, neo4jentity.OrganizationEntity{
		Name: organizationName,
		Hide: true,
	})
}

func LinkOrganizationAsSubsidiary(ctx context.Context, driver *neo4j.DriverWithContext, parentOrganizationId, subOrganizationId, relationType string) {
	query := `MATCH (parent:Organization {id:$parentOrganizationId}),
			(org:Organization {id:$subOrganizationId})
			MERGE (org)-[rel:SUBSIDIARY_OF]->(parent)
			ON CREATE SET rel.type=$type`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"parentOrganizationId": parentOrganizationId,
		"subOrganizationId":    subOrganizationId,
		"type":                 relationType,
	})
}

func RefreshLastTouchpoint(ctx context.Context, driver *neo4j.DriverWithContext, organizationId, timelineEventId string, timelineEventAt time.Time, timelineEventType model.LastTouchpointType) {
	query := `MATCH (org:Organization {id:$organizationId})
			SET org.lastTouchpointId=$timelineEventId, org.lastTouchpointAt = $timelineEventAt, org.lastTouchpointType=$timelineEventType`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId":    organizationId,
		"timelineEventId":   timelineEventId,
		"timelineEventAt":   timelineEventAt,
		"timelineEventType": timelineEventType,
	})
}

func LinkSuggestedMerge(ctx context.Context, driver *neo4j.DriverWithContext, primaryOrgId, orgId, suggestedBy string, suggestedAt time.Time, confidence float64) {
	query := `MATCH (primary:Organization {id:$primaryOrgId}),
					(org:Organization {id:$orgId})
			MERGE (org)-[rel:SUGGESTED_MERGE]->(primary)
			ON CREATE SET rel.suggestedBy=$suggestedBy, rel.suggestedAt=$suggestedAt, rel.confidence=$confidence`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"orgId":        orgId,
		"primaryOrgId": primaryOrgId,
		"suggestedBy":  suggestedBy,
		"suggestedAt":  suggestedAt,
		"confidence":   confidence,
	})
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
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
			ON CREATE SET r.id=$id, r.jobTitle=$jobTitle, r.primary=$primary,
							r.createdAt=datetime({timezone: 'UTC'}), r.appSource=$appSource`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":             roleId.String(),
		"contactId":      contactId,
		"organizationId": organizationId,
		"jobTitle":       jobTitle,
		"primary":        primary,
		"appSource":      "test",
	})
	return roleId.String()
}

func UserWorksAs(ctx context.Context, driver *neo4j.DriverWithContext, userId, jobTitle string, description string, primary bool) string {
	var roleId, _ = uuid.NewRandom()
	query := `MATCH (u:User {id:$userId})
			MERGE (u)-[:WORKS_AS]->(r:JobRole)
			ON CREATE SET r.id=$id, r.description=$description, r.jobTitle=$jobTitle, r.primary=$primary,
							r.createdAt=datetime({timezone: 'UTC'}), r.appSource=$appSource`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":          roleId.String(),
		"userId":      userId,
		"jobTitle":    jobTitle,
		"description": description,
		"primary":     primary,
		"appSource":   "test",
	})
	return roleId.String()
}

func UserOwnsContact(ctx context.Context, driver *neo4j.DriverWithContext, userId, contactId string) {
	query := `MATCH (c:Contact {id:$contactId}),
			        (u:User {id:$userId})
			MERGE (u)-[:OWNS]->(c)`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"contactId": contactId,
		"userId":    userId,
	})
}

func DeleteUserOwnsOrganization(ctx context.Context, driver *neo4j.DriverWithContext, userId, organizationId string) {
	query := `MATCH (u:User {id:$userId})-[r:OWNS]->(o:Organization {id:$organizationId})     
			DELETE r`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"userId":         userId,
	})
}

func UserHasCalendar(ctx context.Context, driver *neo4j.DriverWithContext, userId, link, calType string, primary bool) string {
	var calId, _ = uuid.NewRandom()
	query := `MATCH (u:User {id:$userId})
			MERGE (u)-[:HAS_CALENDAR]->(c:Calendar)
			ON CREATE SET c.id=$id, c.link=$link, c.calType=$calType, c.primary=$primary, c.createdAt=datetime({timezone: 'UTC'}), c.appSource=$appSource`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
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

func CreateLocation(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, location neo4jentity.LocationEntity) string {
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

	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"contactId":  contactId,
		"locationId": locationId,
	})
}

func OrganizationAssociatedWithLocation(ctx context.Context, driver *neo4j.DriverWithContext, organizationId, locationId string) {
	query := `MATCH (org:Organization {id:$organizationId}),
			        (l:Location {id:$locationId})
			MERGE (org)-[:ASSOCIATED_WITH]->(l)`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"locationId":     locationId,
	})
}

func CreateNoteForContact(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contactId, content, contentType string, createdAt time.Time) string {
	var noteId, _ = uuid.NewRandom()

	query := "MATCH (c:Contact {id:$contactId}) " +
		"		MERGE (c)-[:NOTED]->(n:Note {id:$id}) " +
		"		ON CREATE SET 	n.html=$content, " +
		"						n.content=$content, " +
		"						n.contentType=$contentType, " +
		"						n.createdAt=$createdAt, " +
		"						n.updatedAt=$createdAt, " +
		"						n.source=$source, " +
		"						n.sourceOfSource=$source, " +
		"						n.appSource=$appSource, " +
		"						n:Note_%s, " +
		"						n:TimelineEvent, " +
		"						n:TimelineEvent_%s"
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":          noteId.String(),
		"contactId":   contactId,
		"content":     content,
		"contentType": contentType,
		"createdAt":   createdAt,
		"source":      "openline",
		"appSource":   "test",
	})
	return noteId.String()
}

func CreateNoteForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId, content string, createdAt time.Time) string {
	var noteId, _ = uuid.NewRandom()

	query := "MATCH (org:Organization {id:$organizationId}) " +
		"		MERGE (org)-[:NOTED]->(n:Note {id:$id}) " +
		"		ON CREATE SET 	n.html=$content, " +
		"						n.content=$content, " +
		"						n.createdAt=$createdAt, " +
		"						n:Note_%s, " +
		"						n:TimelineEvent, " +
		"						n:TimelineEvent_%s"
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":             noteId.String(),
		"organizationId": organizationId,
		"content":        content,
		"createdAt":      createdAt,
		"source":         "openline",
		"appSource":      "test",
	})
	return noteId.String()
}

func LogEntryCreatedByUser(ctx context.Context, driver *neo4j.DriverWithContext, logEntryId, userId string) {
	query := `MATCH (l:LogEntry {id:$logEntryId}),
					(u:User {id:$userId})
			  MERGE (l)-[:CREATED_BY]->(u)
				`

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"logEntryId": logEntryId,
		"userId":     userId,
	})
}

func LinkNoteWithOrganization(ctx context.Context, driver *neo4j.DriverWithContext, noteId, organizationId string) {
	query := `MATCH (n:Note {id:$noteId}),
			(org:Organization {id:$organizationId})
			MERGE (n)<-[:NOTED]-(org)`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"noteId":         noteId,
	})
}

func NoteCreatedByUser(ctx context.Context, driver *neo4j.DriverWithContext, noteId, userId string) {
	query := `MATCH (u:User {id:$userId})
				MATCH (n:Note {id:$noteId})
			MERGE (u)-[:CREATED]->(n)
`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"noteId": noteId,
		"userId": userId,
	})
}

func LinkContactWithOrganization(ctx context.Context, driver *neo4j.DriverWithContext, contactId, organizationId string) string {
	var jobId, _ = uuid.NewRandom()
	query := `MATCH (c:Contact {id:$contactId}),
			(org:Organization {id:$organizationId})
			MERGE (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(org)
			ON CREATE SET 	j.id=$jobId,
							j.createdAt=$now,
							j.updatedAt=$now,
							j.source=$source,
							j.sourceOfSource=$source,
							j.appSource=$source`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"contactId":      contactId,
		"source":         "test",
		"now":            utils.Now(),
		"jobId":          jobId.String(),
	})
	return jobId.String()
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, describesType, tenant), map[string]any{
		"actionId": actionId,
		"nodeId":   nodeId,
	})
}

func CreateInteractionEvent(ctx context.Context, driver *neo4j.DriverWithContext, tenant, identifier, content, contentType string, channel *string, createdAt time.Time) string {
	return CreateInteractionEventFromEntity(ctx, driver, tenant, entity.InteractionEventEntity{
		EventIdentifier: identifier,
		Content:         content,
		ContentType:     contentType,
		Channel:         channel,
		CreatedAt:       &createdAt,
	})
}

func CreateInteractionEventFromEntity(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, ie entity.InteractionEventEntity) string {
	var interactionEventId, _ = uuid.NewRandom()

	query := "MERGE (ie:InteractionEvent {id:$id})" +
		" ON CREATE SET " +
		"	ie.content=$content, " +
		"	ie.createdAt=$createdAt, " +
		"	ie.channel=$channel, " +
		"	ie.contentType=$contentType, " +
		"	ie.source=$source, " +
		"   ie.hide=$hide, " +
		"	ie.sourceOfTruth=$sourceOfTruth, " +
		"	ie.appSource=$appSource," +
		"	ie:InteractionEvent_%s, ie:TimelineEvent, ie:TimelineEvent_%s," +
		"   ie.identifier=$identifier"
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":            interactionEventId.String(),
		"content":       ie.Content,
		"contentType":   ie.ContentType,
		"channel":       ie.Channel,
		"createdAt":     *ie.CreatedAt,
		"source":        "openline",
		"sourceOfTruth": "openline",
		"appSource":     "test",
		"identifier":    ie.EventIdentifier,
		"hide":          ie.Hide,
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, resolvedQuery, map[string]any{
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

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
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

	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"interactionSessionId": interactionSessionId,
		"nodeId":               nodeId,
		"interactionType":      interactionType,
	})
}

func InteractionEventSentBy(ctx context.Context, driver *neo4j.DriverWithContext, interactionEventId, nodeId, interactionType string) {
	query := "MATCH (ie:InteractionEvent {id:$interactionEventId}), " +
		"(n {id:$nodeId}) " +
		" MERGE (ie)-[:SENT_BY {type:$interactionType}]->(n) "
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"interactionEventId": interactionEventId,
		"nodeId":             nodeId,
		"interactionType":    interactionType,
	})
}

func MeetingCreatedBy(ctx context.Context, driver *neo4j.DriverWithContext, meetingId, nodeId string) {
	query := "MATCH (m:Meeting {id:$meetingId}), " +
		"(n {id:$nodeId}) " +
		" MERGE (m)-[:CREATED_BY]->(n) "
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"meetingId": meetingId,
		"nodeId":    nodeId,
	})
}

func MeetingAttendedBy(ctx context.Context, driver *neo4j.DriverWithContext, meetingId, nodeId string) {
	query := "MATCH (m:Meeting {id:$meetingId}), " +
		"(n {id:$nodeId}) " +
		" MERGE (m)-[:ATTENDED_BY]->(n) "
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"meetingId": meetingId,
		"nodeId":    nodeId,
	})
}

func InteractionEventSentTo(ctx context.Context, driver *neo4j.DriverWithContext, interactionEventId, nodeId, interactionType string) {
	query := "MATCH (ie:InteractionEvent {id:$interactionEventId}), " +
		"(n {id:$nodeId}) " +
		" MERGE (ie)-[:SENT_TO {type:$interactionType}]->(n) "
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"interactionEventId": interactionEventId,
		"nodeId":             nodeId,
		"interactionType":    interactionType,
	})
}

func InteractionEventPartOfInteractionSession(ctx context.Context, driver *neo4j.DriverWithContext, interactionEventId, interactionSessionId string) {
	query := "MATCH (ie:InteractionEvent {id:$interactionEventId}), " +
		"(is:InteractionSession {id:$interactionSessionId}) " +
		" MERGE (ie)-[:PART_OF]->(is) "
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"interactionEventId":   interactionEventId,
		"interactionSessionId": interactionSessionId,
	})
}

func InteractionEventPartOfMeeting(ctx context.Context, driver *neo4j.DriverWithContext, interactionEventId, meetingId string) {
	query := "MATCH (ie:InteractionEvent {id:$interactionEventId}), " +
		"(m:Meeting {id:$meetingId}) " +
		" MERGE (ie)-[:PART_OF]->(m) "
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"interactionEventId": interactionEventId,
		"meetingId":          meetingId,
	})
}

func InteractionEventPartOfIssue(ctx context.Context, driver *neo4j.DriverWithContext, interactionEventId, issueId string) {
	query := "MATCH (ie:InteractionEvent {id:$interactionEventId}), " +
		"(i:Issue {id:$issueId}) " +
		" MERGE (ie)-[:PART_OF]->(i) "
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"interactionEventId": interactionEventId,
		"issueId":            issueId,
	})
}

func InteractionEventRepliesToInteractionEvent(ctx context.Context, driver *neo4j.DriverWithContext, tenant, interactionEventId, repliesToInteractionEventId string) {
	query := "MATCH (ie:InteractionEvent_%s {id:$interactionEventId}), " +
		"(rie:InteractionEvent_%s {id:$repliesToInteractionEventId}) " +
		" MERGE (ie)-[:REPLIES_TO]->(rie) "
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"interactionEventId":          interactionEventId,
		"repliesToInteractionEventId": repliesToInteractionEventId,
	})
}

func CreateState(ctx context.Context, driver *neo4j.DriverWithContext, countryCodeA3, name, code string) {
	query := "MATCH (c:Country{codeA3: $countryCodeA3}) MERGE (c)<-[:BELONGS_TO_COUNTRY]-(az:State { code: $code }) ON CREATE SET az.id = randomUUID(), az.name = $name, az.createdAt = datetime({timezone: 'UTC'}), az.updatedAt = datetime({timezone: 'UTC'})"
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"countryCodeA3": countryCodeA3,
		"name":          name,
		"code":          code,
	})
}

func CreateSocial(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, social neo4jentity.SocialEntity) string {
	var socialId, _ = uuid.NewRandom()
	query := " MERGE (s:Social {id:$id}) " +
		" ON CREATE SET " +
		"				s.url=$url, " +
		"				s.source=$source, " +
		"				s.sourceOfTruth=$source, " +
		"				s.appSource=$appSource, " +
		"				s.createdAt=$now, " +
		"				s.updatedAt=$now, " +
		"				s:Social_%s"

	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"tenant":    tenant,
		"id":        socialId.String(),
		"source":    neo4jentity.DataSourceOpenline,
		"appSource": "test",
		"url":       social.Url,
		"now":       utils.Now(),
	})
	return socialId.String()
}

func LinkSocialWithEntity(ctx context.Context, driver *neo4j.DriverWithContext, entityId, socialId string) {
	query := `MATCH (e {id:$entityId}), (s:Social {id:$socialId}) MERGE (e)-[:HAS]->(s)`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"entityId": entityId,
		"socialId": socialId,
	})
}

func CreateOrganizationRelationship(ctx context.Context, driver *neo4j.DriverWithContext, name string) {
	query := `MERGE (r:OrganizationRelationship {name:$name}) ON CREATE SET r.id=randomUUID()`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"name": name,
	})
}

func CreateActionForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId string, actionType neo4jenum.ActionType, createdAt time.Time) string {
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
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":             actionId.String(),
		"organizationId": organizationId,
		"type":           actionType,
		"createdAt":      createdAt,
		"source":         "openline",
		"appSource":      "test",
	})
	return actionId.String()
}

func CreateActionForInteractionEvent(ctx context.Context, driver *neo4j.DriverWithContext, tenant, interactionEventId string, actionType neo4jenum.ActionType, createdAt time.Time) string {
	var actionId, _ = uuid.NewRandom()

	query := "MATCH (i:InteractionEvent {id:$interactionEventId}) " +
		"		MERGE (i)<-[:ACTION_ON]-(a:Action {id:$id}) " +
		"		ON CREATE SET 	a.type=$type, " +
		"						a.createdAt=$createdAt, " +
		"						a.updatedAt=$createdAt, " +
		"						a.source=$source, " +
		"						a.appSource=$appSource, " +
		"						a:Action_%s, " +
		"						a:TimelineEvent, " +
		"						a:TimelineEvent_%s"
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":                 actionId.String(),
		"interactionEventId": interactionEventId,
		"type":               actionType,
		"createdAt":          createdAt,
		"source":             "openline",
		"appSource":          "test",
	})
	return actionId.String()
}

func CreateActionForOrganizationWithProperties(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId string, actionType neo4jenum.ActionType, createdAt time.Time, extraProperties map[string]string) string {
	var actionId, _ = uuid.NewRandom()

	query := `MATCH (o:Organization {id:$organizationId}) 
				MERGE (o)<-[:ACTION_ON]-(a:Action {id:$id}) 
				ON CREATE SET 	a.type=$type, 
								a.createdAt=$createdAt, 
								a.source=$source, 
								a.appSource=$appSource, 
								a:Action_%s, 
								a:TimelineEvent, 
								a:TimelineEvent_%s,
								a += $extraProperties`
	neo4jtest.ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":              actionId.String(),
		"organizationId":  organizationId,
		"type":            actionType,
		"createdAt":       createdAt,
		"source":          "openline",
		"appSource":       "test",
		"extraProperties": extraProperties,
	})
	return actionId.String()
}

// Deprecated
func ActiveRenewalOpportunityForContract(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId, opportunityId string) string {
	query := fmt.Sprintf(`
				MATCH (c:Contract_%s {id:$contractId}), (op:Opportunity_%s {id:$opportunityId})
				MERGE (c)-[:ACTIVE_RENEWAL]->(op)
				`, tenant, tenant)

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"opportunityId": opportunityId,
		"contractId":    contractId,
	})
	return opportunityId
}

func OpportunityCreatedBy(ctx context.Context, driver *neo4j.DriverWithContext, opportunityId, entityId string) {
	query := `MATCH (e:User {id:$entityId}), (op:Opportunity {id:$opportunityId})
			MERGE (e)<-[:CREATED_BY]-(op)`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"opportunityId": opportunityId,
		"entityId":      entityId,
	})
}

// Deprecated
func OpportunityOwnedBy(ctx context.Context, driver *neo4j.DriverWithContext, opportunityId, entityId string) {
	query := `MATCH (e:User {id:$entityId}), (op:Opportunity {id:$opportunityId})
			MERGE (e)-[:OWNS]->(op)`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"opportunityId": opportunityId,
		"entityId":      entityId,
	})
}

func GetTotalCountOfNodes(ctx context.Context, driver *neo4j.DriverWithContext) int {
	query := `MATCH (n) RETURN count(n)`
	result := neo4jtest.ExecuteReadQueryWithSingleReturn(ctx, driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}
