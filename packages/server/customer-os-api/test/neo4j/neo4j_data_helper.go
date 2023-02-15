package neo4j

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"golang.org/x/net/context"
)

func CleanupAllData(ctx context.Context, driver *neo4j.DriverWithContext) {
	ExecuteWriteQuery(ctx, driver, `MATCH (n) DETACH DELETE n`, map[string]any{})
}

func CreateFullTextBasicSearchIndexes(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := fmt.Sprintf("DROP INDEX basicSearchStandard_%s IF EXISTS", tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{})

	query = fmt.Sprintf("CREATE FULLTEXT INDEX basicSearchStandard_%s FOR (n:Contact_%s|Email_%s|Organization_%s) ON EACH [n.firstName, n.lastName, n.name, n.email] "+
		"OPTIONS {  indexConfig: { `fulltext.analyzer`: 'standard', `fulltext.eventually_consistent`: true } }", tenant, tenant, tenant, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{})

	query = fmt.Sprintf("DROP INDEX basicSearchSimple_%s IF EXISTS", tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{})

	query = fmt.Sprintf("CREATE FULLTEXT INDEX basicSearchSimple_%s FOR (n:Contact_%s|Email_%s|Organization_%s) ON EACH [n.firstName, n.lastName, n.email, n.name] "+
		"OPTIONS {  indexConfig: { `fulltext.analyzer`: 'simple', `fulltext.eventually_consistent`: true } }", tenant, tenant, tenant, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{})
}

func CreateTenant(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MERGE (t:Tenant {name:$tenant})`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant": tenant,
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

func CreateUserWithId(ctx context.Context, driver *neo4j.DriverWithContext, tenant, userId string, user entity.UserEntity) string {
	if len(userId) == 0 {
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
				})-[:USER_BELONGS_TO_TENANT]->(t)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":        tenant,
		"userId":        userId,
		"firstName":     user.FirstName,
		"lastName":      user.LastName,
		"source":        user.Source,
		"sourceOfTruth": user.SourceOfTruth,
	})
	return userId
}

func CreateDefaultContact(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) string {
	return CreateContact(ctx, driver, tenant, entity.ContactEntity{Title: "MR", FirstName: "first", LastName: "last"})
}

func CreateContactWith(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, firstName string, lastName string) string {
	return CreateContact(ctx, driver, tenant, entity.ContactEntity{Title: "MR", FirstName: firstName, LastName: lastName})
}

func CreateContact(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, contact entity.ContactEntity) string {
	var contactId, _ = uuid.NewRandom()
	query := "MATCH (t:Tenant {name:$tenant}) MERGE (c:Contact {id: $contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET c.title=$title, c.firstName=$firstName, c.lastName=$lastName, c.label=$label, c.createdAt=datetime({timezone: 'UTC'}), " +
		" c:%s"

	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, "Contact_"+tenant), map[string]any{
		"tenant":    tenant,
		"contactId": contactId.String(),
		"title":     contact.Title,
		"firstName": contact.FirstName,
		"lastName":  contact.LastName,
		"label":     contact.Label,
	})
	return contactId.String()
}

func CreateContactGroup(ctx context.Context, driver *neo4j.DriverWithContext, tenant, name string) string {
	var contactGroupId, _ = uuid.NewRandom()
	query := `
			MATCH (t:Tenant {name:$tenant})
			MERGE (g:ContactGroup {
				  id: $id,
				  name: $name
				})-[:GROUP_BELONGS_TO_TENANT]->(t)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant": tenant,
		"id":     contactGroupId.String(),
		"name":   name,
	})
	return contactGroupId.String()
}

func AddContactToGroup(ctx context.Context, driver *neo4j.DriverWithContext, contactId, groupId string) {
	query := `MATCH (c:Contact {id:$contactId}), (g:ContactGroup {id:$groupId})
				MERGE (c)-[:BELONGS_TO_GROUP]->(g)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"contactId": contactId,
		"groupId":   groupId,
	})
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
		" MERGE (e:Email {email: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)" +
		" ON CREATE SET " +
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

func AddPhoneNumberToContact(ctx context.Context, driver *neo4j.DriverWithContext, contactId string, e164 string, primary bool, label string) {
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (:PhoneNumber {
				  id: randomUUID(),
				  e164: $e164,
				  label: $label
				})<-[:PHONE_ASSOCIATED_WITH {primary:$primary}]-(c)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"contactId": contactId,
		"primary":   primary,
		"e164":      e164,
		"label":     label,
	})
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

func TagContact(ctx context.Context, driver *neo4j.DriverWithContext, contactId, tagId string) {
	query := `MATCH (c:Contact {id:$contactId}), (tag:Tag {id:$tagId})
			MERGE (c)-[r:TAGGED]->(tag)
			ON CREATE SET r.taggedAt=datetime({timezone: 'UTC'})`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tagId":     tagId,
		"contactId": contactId,
	})
}

func CreateOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationName string) string {
	var organizationId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id})
			ON CREATE SET org.name=$name, org:%s`
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, "Organization_"+tenant), map[string]any{
		"id":     organizationId.String(),
		"tenant": tenant,
		"name":   organizationName,
	})
	return organizationId.String()
}

func CreateOrganizationType(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationTypeName string) string {
	var organizationTypeId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType {id:$id})
			ON CREATE SET ot.name=$name`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":     organizationTypeId.String(),
		"tenant": tenant,
		"name":   organizationTypeName,
	})
	return organizationTypeId.String()
}

func CreateFullOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, organization entity.OrganizationEntity) string {
	var organizationId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id})
			ON CREATE SET org.name=$name, org.description=$description, org.domain=$domain, org.website=$website,
							org.industry=$industry, org.isPublic=$isPublic, org.createdAt=datetime({timezone: 'UTC'})
`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":          organizationId.String(),
		"tenant":      tenant,
		"name":        organization.Name,
		"description": organization.Description,
		"domain":      organization.Domain,
		"website":     organization.Website,
		"industry":    organization.Industry,
		"isPublic":    organization.IsPublic,
	})
	return organizationId.String()
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

func SetOrganizationTypeForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, organizationId, organizationTypeId string) {
	query := `
			MATCH (org:Organization {id:$organizationId}),
				  (ot:OrganizationType {id:$organizationTypeId})
			MERGE (org)-[:IS_OF_TYPE]->(ot)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId":     organizationId,
		"organizationTypeId": organizationTypeId,
	})
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

func CreateConversation(ctx context.Context, driver *neo4j.DriverWithContext, userId, contactId string) string {
	var conversationId, _ = uuid.NewRandom()
	query := `MATCH (c:Contact {id:$contactId}),
			        (u:User {id:$userId})
			MERGE (u)-[:PARTICIPATES]->(o:Conversation {id:$conversationId})<-[:PARTICIPATES]-(c)
			ON CREATE SET o.startedAt=datetime({timezone: 'UTC'}), o.status="ACTIVE", o.channel="VOICE", o.messageCount=0 `
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"contactId":      contactId,
		"userId":         userId,
		"conversationId": conversationId.String(),
	})
	return conversationId.String()
}

func CreatePageView(ctx context.Context, driver *neo4j.DriverWithContext, contactId string, actionEntity entity.PageViewEntity) string {
	var actionId, _ = uuid.NewRandom()
	query := `MATCH (c:Contact {id:$contactId})
			MERGE (c)-[:HAS_ACTION]->(a:Action:PageView {id:$actionId})
			ON CREATE SET
				a.trackerName=$trackerName,
				a.startedAt=$startedAt,
				a.endedAt=$endedAt,
				a.application=$application,
				a.pageUrl=$pageUrl,
				a.pageTitle=$pageTitle,
				a.sessionId=$sessionId,
				a.orderInSession=$orderInSession,
				a.engagedTime=$engagedTime`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"contactId":      contactId,
		"actionId":       actionId.String(),
		"trackerName":    actionEntity.TrackerName,
		"startedAt":      actionEntity.StartedAt,
		"endedAt":        actionEntity.EndedAt,
		"application":    actionEntity.Application,
		"pageUrl":        actionEntity.PageUrl,
		"pageTitle":      actionEntity.PageTitle,
		"sessionId":      actionEntity.SessionId,
		"orderInSession": actionEntity.OrderInSession,
		"engagedTime":    actionEntity.EngagedTime,
	})
	return actionId.String()
}

func CreateLocation(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, location entity.LocationEntity) string {
	var locationId, _ = uuid.NewRandom()
	query := "MERGE (l:Location {id:$locationId}) " +
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
		"				l:Location_%s"

	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"locationId": locationId.String(),
		"source":     location.Source,
		"appSource":  location.AppSource,
		"name":       location.Name,
		"now":        utils.Now(),
		"country":    location.Country,
		"region":     location.Region,
		"locality":   location.Locality,
		"address":    location.Address,
		"address2":   location.Address2,
		"zip":        location.Zip,
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

func CreateNoteForContact(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contactId, html string) string {
	var noteId, _ = uuid.NewRandom()
	query := "MATCH (c:Contact {id:$contactId}) " +
		"		MERGE (c)-[:NOTED]->(n:Note {id:$id}) " +
		"		ON CREATE SET n.html=$html, n.createdAt=datetime({timezone: 'UTC'}), n:%s"
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, "Note_"+tenant), map[string]any{
		"id":        noteId.String(),
		"contactId": contactId,
		"html":      html,
	})
	return noteId.String()
}

func CreateNoteForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId, html string) string {
	var noteId, _ = uuid.NewRandom()
	query := "MATCH (org:Organization {id:$organizationId}) " +
		"		MERGE (org)-[:NOTED]->(n:Note {id:$id}) " +
		"		ON CREATE SET n.html=$html, n.createdAt=datetime({timezone: 'UTC'}), n:%s"
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, "Note_"+tenant), map[string]any{
		"id":             noteId.String(),
		"organizationId": organizationId,
		"html":           html,
	})
	return noteId.String()
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
			MERGE (c)-[:CONTACT_OF]->(org)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"contactId":      contactId,
	})
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
