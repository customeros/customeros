package neo4j

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func CleanupAllData(driver *neo4j.Driver) {
	ExecuteWriteQuery(driver, `MATCH (n) DETACH DELETE n`, map[string]any{})
}

func CreateFullTextBasicSearchIndexes(driver *neo4j.Driver, tenant string) {
	query := fmt.Sprintf("DROP INDEX basicSearchStandard_%s IF EXISTS", tenant)
	ExecuteWriteQuery(driver, query, map[string]any{})

	query = fmt.Sprintf("CREATE FULLTEXT INDEX basicSearchStandard_%s FOR (n:Contact_%s|Email_%s|Organization_%s) ON EACH [n.firstName, n.lastName, n.name, n.email] "+
		"OPTIONS {  indexConfig: { `fulltext.analyzer`: 'standard', `fulltext.eventually_consistent`: true } }", tenant, tenant, tenant, tenant)
	ExecuteWriteQuery(driver, query, map[string]any{})

	query = fmt.Sprintf("DROP INDEX basicSearchSimple_%s IF EXISTS", tenant)
	ExecuteWriteQuery(driver, query, map[string]any{})

	query = fmt.Sprintf("CREATE FULLTEXT INDEX basicSearchSimple_%s FOR (n:Contact_%s|Email_%s|Organization_%s) ON EACH [n.firstName, n.lastName, n.email, n.name] "+
		"OPTIONS {  indexConfig: { `fulltext.analyzer`: 'simple', `fulltext.eventually_consistent`: true } }", tenant, tenant, tenant, tenant)
	ExecuteWriteQuery(driver, query, map[string]any{})
}

func CreateTenant(driver *neo4j.Driver, tenant string) {
	query := `MERGE (t:Tenant {name:$tenant})`
	ExecuteWriteQuery(driver, query, map[string]any{
		"tenant": tenant,
	})
}

func CreateHubspotExternalSystem(driver *neo4j.Driver, tenant string) {
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"tenant":           tenant,
		"externalSystemId": "hubspot",
	})
}

func CreateDefaultUser(driver *neo4j.Driver, tenant string) string {
	return CreateUser(driver, tenant, entity.UserEntity{
		FirstName:     "first",
		LastName:      "last",
		Source:        "openline",
		SourceOfTruth: "openline",
	})
}

func CreateDefaultUserWithId(driver *neo4j.Driver, tenant, userId string) string {
	return CreateUserWithId(driver, tenant, userId, entity.UserEntity{
		FirstName:     "first",
		LastName:      "last",
		Source:        "openline",
		SourceOfTruth: "openline",
	})
}

func CreateUser(driver *neo4j.Driver, tenant string, user entity.UserEntity) string {
	return CreateUserWithId(driver, tenant, "", user)
}

func CreateUserWithId(driver *neo4j.Driver, tenant, userId string, user entity.UserEntity) string {
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
	ExecuteWriteQuery(driver, query, map[string]any{
		"tenant":        tenant,
		"userId":        userId,
		"firstName":     user.FirstName,
		"lastName":      user.LastName,
		"source":        user.Source,
		"sourceOfTruth": user.SourceOfTruth,
	})
	return userId
}

func CreateDefaultContact(driver *neo4j.Driver, tenant string) string {
	return CreateContact(driver, tenant, entity.ContactEntity{Title: "MR", FirstName: "first", LastName: "last"})
}

func CreateContactWith(driver *neo4j.Driver, tenant string, firstName string, lastName string) string {
	return CreateContact(driver, tenant, entity.ContactEntity{Title: "MR", FirstName: firstName, LastName: lastName})
}

func CreateContact(driver *neo4j.Driver, tenant string, contact entity.ContactEntity) string {
	var contactId, _ = uuid.NewRandom()
	query := "MATCH (t:Tenant {name:$tenant}) MERGE (c:Contact {id: $contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET c.title=$title, c.firstName=$firstName, c.lastName=$lastName, c.label=$label, c.createdAt=datetime({timezone: 'UTC'}), " +
		" c:%s"

	ExecuteWriteQuery(driver, fmt.Sprintf(query, "Contact_"+tenant), map[string]any{
		"tenant":    tenant,
		"contactId": contactId.String(),
		"title":     contact.Title,
		"firstName": contact.FirstName,
		"lastName":  contact.LastName,
		"label":     contact.Label,
	})
	return contactId.String()
}

func CreateContactGroup(driver *neo4j.Driver, tenant, name string) string {
	var contactGroupId, _ = uuid.NewRandom()
	query := `
			MATCH (t:Tenant {name:$tenant})
			MERGE (g:ContactGroup {
				  id: $id,
				  name: $name
				})-[:GROUP_BELONGS_TO_TENANT]->(t)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"tenant": tenant,
		"id":     contactGroupId.String(),
		"name":   name,
	})
	return contactGroupId.String()
}

func AddContactToGroup(driver *neo4j.Driver, contactId, groupId string) {
	query := `MATCH (c:Contact {id:$contactId}), (g:ContactGroup {id:$groupId})
				MERGE (c)-[:BELONGS_TO_GROUP]->(g)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"groupId":   groupId,
	})
}

func CreateDefaultFieldSet(driver *neo4j.Driver, contactId string) string {
	return CreateFieldSet(driver, contactId, entity.FieldSetEntity{Name: "name", Source: entity.DataSourceOpenline, SourceOfTruth: entity.DataSourceOpenline})
}

func CreateFieldSet(driver *neo4j.Driver, contactId string, fieldSet entity.FieldSetEntity) string {
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
	ExecuteWriteQuery(driver, query, map[string]any{
		"contactId":     contactId,
		"fieldSetId":    fieldSetId.String(),
		"name":          fieldSet.Name,
		"source":        fieldSet.Source,
		"sourceOfTruth": fieldSet.SourceOfTruth,
	})
	return fieldSetId.String()
}

func CreateDefaultCustomFieldInSet(driver *neo4j.Driver, fieldSetId string) string {
	return createCustomFieldInSet(driver, fieldSetId,
		entity.CustomFieldEntity{
			Name:          "name",
			Source:        entity.DataSourceOpenline,
			SourceOfTruth: entity.DataSourceOpenline,
			DataType:      model.CustomFieldDataTypeText.String(),
			Value:         model.AnyTypeValue{Str: utils.StringPtr("value")}})
}

func createCustomFieldInSet(driver *neo4j.Driver, fieldSetId string, customField entity.CustomFieldEntity) string {
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
	ExecuteWriteQuery(driver, query, map[string]any{
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

func CreateDefaultCustomFieldInContact(driver *neo4j.Driver, contactId string) string {
	return createCustomFieldInContact(driver, contactId,
		entity.CustomFieldEntity{
			Name:          "name",
			DataType:      model.CustomFieldDataTypeText.String(),
			Source:        entity.DataSourceOpenline,
			SourceOfTruth: entity.DataSourceOpenline,
			Value:         model.AnyTypeValue{Str: utils.StringPtr("value")}})
}

func createCustomFieldInContact(driver *neo4j.Driver, contactId string, customField entity.CustomFieldEntity) string {
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
	ExecuteWriteQuery(driver, query, map[string]any{
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

func AddEmailTo(driver *neo4j.Driver, entityType repository.EntityType, tenant, entityId, email string, primary bool, label string) string {
	query := ""

	switch entityType {
	case repository.CONTACT:
		query = "MATCH (entity:Contact {id:$entityId}) "
	case repository.USER:
		query = "MATCH (entity:User {id:$entityId}) "
	case repository.ORGANIZATION:
		query = "MATCH (entity:Organization {id:$entityId}) "
	}

	var emailId, _ = uuid.NewRandom()
	query = query + "MERGE (e:Email {id: $emailId,email: $email,label: $label})<-[:HAS {primary:$primary}]-(entity) ON CREATE SET e:%s"

	ExecuteWriteQuery(driver, fmt.Sprintf(query, "Email_"+tenant), map[string]any{
		"entityId": entityId,
		"primary":  primary,
		"email":    email,
		"label":    label,
		"emailId":  emailId.String(),
	})
	return emailId.String()
}

func AddPhoneNumberToContact(driver *neo4j.Driver, contactId string, e164 string, primary bool, label string) {
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (:PhoneNumber {
				  id: randomUUID(),
				  e164: $e164,
				  label: $label
				})<-[:PHONE_ASSOCIATED_WITH {primary:$primary}]-(c)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"primary":   primary,
		"e164":      e164,
		"label":     label,
	})
}

func CreateEntityTemplate(driver *neo4j.Driver, tenant, extends string) string {
	var templateId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:EntityTemplate {id:$templateId})-[:ENTITY_TEMPLATE_BELONGS_TO_TENANT]->(t)
			ON CREATE SET e.extends=$extends, e.name=$name`
	ExecuteWriteQuery(driver, query, map[string]any{
		"templateId": templateId.String(),
		"tenant":     tenant,
		"extends":    extends,
		"name":       "template name",
	})
	return templateId.String()
}

func LinkEntityTemplateToContact(driver *neo4j.Driver, entityTemplateId, contactId string) {
	query := `MATCH (c:Contact {id:$contactId}),
			(e:EntityTemplate {id:$TemplateId})
			MERGE (c)-[:IS_DEFINED_BY]->(e)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"TemplateId": entityTemplateId,
		"contactId":  contactId,
	})
}

func AddFieldTemplateToEntity(driver *neo4j.Driver, entityTemplateId string) string {
	var templateId, _ = uuid.NewRandom()
	query := `MATCH (e:EntityTemplate {id:$entityTemplateId})
			MERGE (f:CustomFieldTemplate {id:$templateId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	ExecuteWriteQuery(driver, query, map[string]any{
		"templateId":       templateId.String(),
		"entityTemplateId": entityTemplateId,
		"type":             "TEXT",
		"order":            1,
		"mandatory":        false,
		"name":             "template name",
	})
	return templateId.String()
}

func AddFieldTemplateToSet(driver *neo4j.Driver, setTemplateId string) string {
	var templateId, _ = uuid.NewRandom()
	query := `MATCH (e:FieldSetTemplate {id:$setTemplateId})
			MERGE (f:CustomFieldTemplate {id:$templateId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	ExecuteWriteQuery(driver, query, map[string]any{
		"templateId":    templateId.String(),
		"setTemplateId": setTemplateId,
		"type":          "TEXT",
		"order":         1,
		"mandatory":     false,
		"name":          "template name",
	})
	return templateId.String()
}

func AddSetTemplateToEntity(driver *neo4j.Driver, entityTemplateId string) string {
	var templateId, _ = uuid.NewRandom()
	query := `MATCH (e:EntityTemplate {id:$entityTemplateId})
			MERGE (f:FieldSetTemplate {id:$templateId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	ExecuteWriteQuery(driver, query, map[string]any{
		"templateId":       templateId.String(),
		"entityTemplateId": entityTemplateId,
		"type":             "TEXT",
		"order":            1,
		"mandatory":        false,
		"name":             "set name",
	})
	return templateId.String()
}

func CreateTag(driver *neo4j.Driver, tenant, tagName string) string {
	var tagId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})
			ON CREATE SET tag.name=$name, tag.source=$source, tag.appSource=$appSource, tag.createdAt=$now, tag.updatedAt=$now`
	ExecuteWriteQuery(driver, query, map[string]any{
		"id":        tagId.String(),
		"tenant":    tenant,
		"name":      tagName,
		"source":    "openline",
		"appSource": "test",
		"now":       utils.Now(),
	})
	return tagId.String()
}

func TagContact(driver *neo4j.Driver, contactId, tagId string) {
	query := `MATCH (c:Contact {id:$contactId}), (tag:Tag {id:$tagId})
			MERGE (c)-[r:TAGGED]->(tag)
			ON CREATE SET r.taggedAt=datetime({timezone: 'UTC'})`
	ExecuteWriteQuery(driver, query, map[string]any{
		"tagId":     tagId,
		"contactId": contactId,
	})
}

func CreateOrganization(driver *neo4j.Driver, tenant, organizationName string) string {
	var organizationId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id})
			ON CREATE SET org.name=$name, org:%s`
	ExecuteWriteQuery(driver, fmt.Sprintf(query, "Organization_"+tenant), map[string]any{
		"id":     organizationId.String(),
		"tenant": tenant,
		"name":   organizationName,
	})
	return organizationId.String()
}

func CreateOrganizationType(driver *neo4j.Driver, tenant, organizationTypeName string) string {
	var organizationTypeId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType {id:$id})
			ON CREATE SET ot.name=$name`
	ExecuteWriteQuery(driver, query, map[string]any{
		"id":     organizationTypeId.String(),
		"tenant": tenant,
		"name":   organizationTypeName,
	})
	return organizationTypeId.String()
}

func CreateFullOrganization(driver *neo4j.Driver, tenant string, organization entity.OrganizationEntity) string {
	var organizationId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id})
			ON CREATE SET org.name=$name, org.description=$description, org.domain=$domain, org.website=$website,
							org.industry=$industry, org.isPublic=$isPublic, org.createdAt=datetime({timezone: 'UTC'})
`
	ExecuteWriteQuery(driver, query, map[string]any{
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

func ContactWorksForOrganization(driver *neo4j.Driver, contactId, organizationId, jobTitle string, primary bool) string {
	var roleId, _ = uuid.NewRandom()
	query := `MATCH (c:Contact {id:$contactId}),
			        (org:Organization {id:$organizationId})
			MERGE (c)-[:WORKS_AS]->(r:JobRole)-[:ROLE_IN]->(org)
			ON CREATE SET r.id=$id, r.jobTitle=$jobTitle, r.primary=$primary, r.responsibilityLevel=$responsibilityLevel,
							r.createdAt=datetime({timezone: 'UTC'}), r.appSource=$appSource`
	ExecuteWriteQuery(driver, query, map[string]any{
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

func SetOrganizationTypeForOrganization(driver *neo4j.Driver, organizationId, organizationTypeId string) {
	query := `
			MATCH (org:Organization {id:$organizationId}),
				  (ot:OrganizationType {id:$organizationTypeId})
			MERGE (org)-[:IS_OF_TYPE]->(ot)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"organizationId":     organizationId,
		"organizationTypeId": organizationTypeId,
	})
}

func UserOwnsContact(driver *neo4j.Driver, userId, contactId string) {
	query := `MATCH (c:Contact {id:$contactId}),
			        (u:User {id:$userId})
			MERGE (u)-[:OWNS]->(c)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"userId":    userId,
	})
}

func CreateConversation(driver *neo4j.Driver, userId, contactId string) string {
	var conversationId, _ = uuid.NewRandom()
	query := `MATCH (c:Contact {id:$contactId}),
			        (u:User {id:$userId})
			MERGE (u)-[:PARTICIPATES]->(o:Conversation {id:$conversationId})<-[:PARTICIPATES]-(c)
			ON CREATE SET o.startedAt=datetime({timezone: 'UTC'}), o.status="ACTIVE", o.channel="VOICE", o.messageCount=0 `
	ExecuteWriteQuery(driver, query, map[string]any{
		"contactId":      contactId,
		"userId":         userId,
		"conversationId": conversationId.String(),
	})
	return conversationId.String()
}

func CreatePageView(driver *neo4j.Driver, contactId string, actionEntity entity.PageViewEntity) string {
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
	ExecuteWriteQuery(driver, query, map[string]any{
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

func CreateLocation(driver *neo4j.Driver, tenant string, location entity.LocationEntity) string {
	var locationId, _ = uuid.NewRandom()
	query := "MERGE (l:Location {id:$locationId}) " +
		" ON CREATE SET l.name=$name, " +
		"				l.source=$source, " +
		"				l.appSource=$appSource, " +
		"				l.createdAt=$now, " +
		"				l.updatedAt=$now, " +
		"				l:Location_%s"

	ExecuteWriteQuery(driver, fmt.Sprintf(query, tenant), map[string]any{
		"locationId": locationId.String(),
		"source":     location.Source,
		"appSource":  location.AppSource,
		"name":       location.Name,
		"now":        utils.Now(),
	})
	return locationId.String()
}

func CreatePlaceForLocation(driver *neo4j.Driver, place entity.PlaceEntity, locationId string) string {
	var placeId, _ = uuid.NewRandom()
	query := `MATCH (l:Location {id:$locationId})
	MERGE (l)-[:LOCATED_AT]->(a:Place {id:$placeId})
			ON CREATE SET  a.country=$country, 
							a.state=$state, 
							a.city=$city, 
							a.address=$address,
							a.address2=$address2, 
							a.zip=$zip, 
							a.fax=$fax, 
							a.phone=$phone,
							a.source=$source, 
							a.sourceOfTruth=$sourceOfTruth, 
							a.appSource=$appSource,
							a.createdAt=datetime({timezone: 'UTC'}), 
							a.updatedAt=datetime({timezone: 'UTC'})`
	ExecuteWriteQuery(driver, query, map[string]any{
		"placeId":       placeId.String(),
		"locationId":    locationId,
		"source":        place.Source,
		"appSource":     place.AppSource,
		"sourceOfTruth": place.Source,
		"country":       place.Country,
		"state":         place.State,
		"city":          place.City,
		"address":       place.Address,
		"address2":      place.Address2,
		"zip":           place.Zip,
		"phone":         place.Phone,
		"fax":           place.Fax,
	})
	return placeId.String()
}

func ContactAssociatedWithLocation(driver *neo4j.Driver, contactId, locationId string) {
	query := `MATCH (c:Contact {id:$contactId}),
			        (l:Location {id:$locationId})
			MERGE (c)-[:ASSOCIATED_WITH]->(l)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"contactId":  contactId,
		"locationId": locationId,
	})
}

func OrganizationAssociatedWithLocation(driver *neo4j.Driver, organizationId, locationId string) {
	query := `MATCH (org:Organization {id:$organizationId}),
			        (l:Location {id:$locationId})
			MERGE (org)-[:ASSOCIATED_WITH]->(l)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"organizationId": organizationId,
		"locationId":     locationId,
	})
}

func CreateNoteForContact(driver *neo4j.Driver, tenant, contactId, html string) string {
	var noteId, _ = uuid.NewRandom()
	query := "MATCH (c:Contact {id:$contactId}) " +
		"		MERGE (c)-[:NOTED]->(n:Note {id:$id}) " +
		"		ON CREATE SET n.html=$html, n.createdAt=datetime({timezone: 'UTC'}), n:%s"
	ExecuteWriteQuery(driver, fmt.Sprintf(query, "Note_"+tenant), map[string]any{
		"id":        noteId.String(),
		"contactId": contactId,
		"html":      html,
	})
	return noteId.String()
}

func CreateNoteForOrganization(driver *neo4j.Driver, tenant, organizationId, html string) string {
	var noteId, _ = uuid.NewRandom()
	query := "MATCH (org:Organization {id:$organizationId}) " +
		"		MERGE (org)-[:NOTED]->(n:Note {id:$id}) " +
		"		ON CREATE SET n.html=$html, n.createdAt=datetime({timezone: 'UTC'}), n:%s"
	ExecuteWriteQuery(driver, fmt.Sprintf(query, "Note_"+tenant), map[string]any{
		"id":             noteId.String(),
		"organizationId": organizationId,
		"html":           html,
	})
	return noteId.String()
}

func NoteCreatedByUser(driver *neo4j.Driver, noteId, userId string) {
	query := `MATCH (u:User {id:$userId})
				MATCH (n:Note {id:$noteId})
			MERGE (u)-[:CREATED]->(n)
`
	ExecuteWriteQuery(driver, query, map[string]any{
		"noteId": noteId,
		"userId": userId,
	})
}

func LinkContactWithOrganization(driver *neo4j.Driver, contactId, organizationId string) {
	query := `MATCH (c:Contact {id:$contactId}),
			(org:Organization {id:$organizationId})
			MERGE (c)-[:CONTACT_OF]->(org)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"organizationId": organizationId,
		"contactId":      contactId,
	})
}

func Q1(driver *neo4j.Driver, tenant string) int64 {
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  MATCH (t)--(c:Contact)
		  MATCH (o)--(c) 
RETURN count(t)`)
	result := ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{
		"tenant": tenant,
	})
	return int64(result.(*db.Record).Values[0].(int64))
}

func Q2(driver *neo4j.Driver, tenant string) int64 {
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  WHERE NOT (o)--(:Contact)
RETURN count(t)`)
	result := ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{
		"tenant": tenant,
	})
	return int64(result.(*db.Record).Values[0].(int64))
}

func Q3(driver *neo4j.Driver, tenant string) int64 {
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})--(c:Contact)
		  WHERE NOT (c)--(:Organization)
RETURN count(t)`)
	result := ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{
		"tenant": tenant,
	})
	return int64(result.(*db.Record).Values[0].(int64))
}

func Q4(driver *neo4j.Driver, tenant string) int64 {
	query := fmt.Sprintf(`CALL {
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  MATCH (t)--(c:Contact)
		  MATCH (o)--(c)
		  RETURN count(o) as t
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(o:Organization)
		  WHERE NOT (o)--(:Contact)
		  RETURN count(o) as t
		  UNION
		  MATCH (t:Tenant {name:$tenant})--(c:Contact)
		  WHERE NOT (c)--(:Organization)
		  RETURN count(c) as t
		} RETURN sum(t)`)
	result := ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{
		"tenant": tenant,
	})
	return int64(result.(*db.Record).Values[0].(int64))
}

func GetCountOfNodes(driver *neo4j.Driver, nodeLabel string) int {
	query := fmt.Sprintf(`MATCH (n:%s) RETURN count(n)`, nodeLabel)
	result := ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetCountOfRelationships(driver *neo4j.Driver, relationship string) int {
	query := fmt.Sprintf(`MATCH (a)-[r:%s]-(b) RETURN count(distinct r)`, relationship)
	result := ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetTotalCountOfNodes(driver *neo4j.Driver) int {
	query := `MATCH (n) RETURN count(n)`
	result := ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetAllLabels(driver *neo4j.Driver) []string {
	query := `MATCH (n) RETURN DISTINCT labels(n)`
	dbRecords := ExecuteReadQueryWithCollectionReturn(driver, query, map[string]any{})
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
