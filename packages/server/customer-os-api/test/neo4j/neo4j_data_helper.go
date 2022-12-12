package neo4j

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"time"
)

func CleanupAllData(driver *neo4j.Driver) {
	ExecuteWriteQuery(driver, `MATCH (n) DETACH DELETE n`, map[string]any{})
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
		FirstName: "first",
		LastName:  "last",
		Email:     "user@openline.ai",
	})
}

func CreateUser(driver *neo4j.Driver, tenant string, user entity.UserEntity) string {
	var userId, _ = uuid.NewRandom()
	query := `
		MATCH (t:Tenant {name:$tenant})
			MERGE (u:User {
				  id: $userId,
				  firstName: $firstName,
				  lastName: $lastName,
				  email: $email,
				  createdAt :datetime({timezone: 'UTC'})
				})-[:USER_BELONGS_TO_TENANT]->(t)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"tenant":    tenant,
		"userId":    userId.String(),
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	})
	return userId.String()
}

func CreateDefaultContact(driver *neo4j.Driver, tenant string) string {
	return CreateContact(driver, tenant, entity.ContactEntity{Title: "MR", FirstName: "first", LastName: "last"})
}

func CreateContact(driver *neo4j.Driver, tenant string, contact entity.ContactEntity) string {
	var contactId, _ = uuid.NewRandom()
	query := `
			MATCH (t:Tenant {name:$tenant})
			MERGE (c:Contact {
				  id: $contactId,
				  title: $title,
				  firstName: $firstName,
				  lastName: $lastName,
				  label: $label,
				  notes: $notes,
				  createdAt :datetime({timezone: 'UTC'})
				})-[:CONTACT_BELONGS_TO_TENANT]->(t)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"tenant":    tenant,
		"contactId": contactId.String(),
		"title":     contact.Title,
		"firstName": contact.FirstName,
		"lastName":  contact.LastName,
		"notes":     contact.Notes,
		"label":     contact.Label,
	})
	return contactId.String()
}

func SetContactTypeForContact(driver *neo4j.Driver, contactId, contactTypeId string) {
	query := `
			MATCH (c:Contact {id:$contactId}),
				  (o:ContactType {id:$contactTypeId})
			MERGE (c)-[:IS_OF_TYPE]->(o)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"contactId":     contactId,
		"contactTypeId": contactTypeId,
	})
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
	return CreateFieldSet(driver, contactId, entity.FieldSetEntity{Name: "name"})
}

func CreateFieldSet(driver *neo4j.Driver, contactId string, fieldSet entity.FieldSetEntity) string {
	var fieldSetId, _ = uuid.NewRandom()
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (s:FieldSet {
				  id: $fieldSetId,
				  name: $name
				})<-[:HAS_COMPLEX_PROPERTY {added:datetime({timezone: 'UTC'})}]-(c)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"contactId":  contactId,
		"fieldSetId": fieldSetId.String(),
		"name":       fieldSet.Name,
	})
	return fieldSetId.String()
}

func CreateDefaultCustomFieldInSet(driver *neo4j.Driver, fieldSetId string) string {
	return createCustomFieldInSet(driver, fieldSetId,
		entity.CustomFieldEntity{
			Name:     "name",
			DataType: model.CustomFieldDataTypeText.String(),
			Value:    model.AnyTypeValue{Str: utils.StringPtr("value")}})
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
			"	  name: $name "+
			"	})<-[:HAS_PROPERTY]-(s)", customField.NodeLabel(), customField.PropertyName())
	ExecuteWriteQuery(driver, query, map[string]any{
		"fieldSetId": fieldSetId,
		"fieldId":    fieldId.String(),
		"name":       customField.Name,
		"datatype":   customField.DataType,
		"value":      customField.Value.RealValue(),
	})
	return fieldId.String()
}

func CreateDefaultCustomFieldInContact(driver *neo4j.Driver, contactId string) string {
	return createCustomFieldInContact(driver, contactId,
		entity.CustomFieldEntity{
			Name:     "name",
			DataType: model.CustomFieldDataTypeText.String(),
			Value:    model.AnyTypeValue{Str: utils.StringPtr("value")}})
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
			"	  name: $name "+
			"	})<-[:HAS_PROPERTY]-(c)", customField.NodeLabel(), customField.PropertyName())
	ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"fieldId":   fieldId.String(),
		"name":      customField.Name,
		"datatype":  customField.DataType,
		"value":     customField.Value.RealValue(),
	})
	return fieldId.String()
}

func AddEmailToContact(driver *neo4j.Driver, contactId string, email string, primary bool, label string) {
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (:Email {
				  id: randomUUID(),
				  email: $email,
				  label: $label
				})<-[:EMAILED_AT {primary:$primary}]-(c)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"primary":   primary,
		"email":     email,
		"label":     label,
	})
}

func AddPhoneNumberToContact(driver *neo4j.Driver, contactId string, e164 string, primary bool, label string) {
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (:PhoneNumber {
				  id: randomUUID(),
				  e164: $e164,
				  label: $label
				})<-[:CALLED_AT {primary:$primary}]-(c)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"primary":   primary,
		"e164":      e164,
		"label":     label,
	})
}

func CreateEntityDefinition(driver *neo4j.Driver, tenant, extends string) string {
	var definitionId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:EntityDefinition {id:$definitionId})<-[:USES_ENTITY_DEFINITION]-(t)
			ON CREATE SET e.extends=$extends, e.name=$name`
	ExecuteWriteQuery(driver, query, map[string]any{
		"definitionId": definitionId.String(),
		"tenant":       tenant,
		"extends":      extends,
		"name":         "definition name",
	})
	return definitionId.String()
}

func LinkEntityDefinitionToContact(driver *neo4j.Driver, entityDefinitionId, contactId string) {
	query := `MATCH (c:Contact {id:$contactId}),
			(e:EntityDefinition {id:$definitionId})
			MERGE (c)-[:IS_DEFINED_BY]->(e)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"definitionId": entityDefinitionId,
		"contactId":    contactId,
	})
}

func AddFieldDefinitionToEntity(driver *neo4j.Driver, entityDefinitionId string) string {
	var definitionId, _ = uuid.NewRandom()
	query := `MATCH (e:EntityDefinition {id:$entityDefinitionId})
			MERGE (f:CustomFieldDefinition {id:$definitionId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	ExecuteWriteQuery(driver, query, map[string]any{
		"definitionId":       definitionId.String(),
		"entityDefinitionId": entityDefinitionId,
		"type":               "TEXT",
		"order":              1,
		"mandatory":          false,
		"name":               "definition name",
	})
	return definitionId.String()
}

func AddFieldDefinitionToSet(driver *neo4j.Driver, setDefinitionId string) string {
	var definitionId, _ = uuid.NewRandom()
	query := `MATCH (e:FieldSetDefinition {id:$setDefinitionId})
			MERGE (f:CustomFieldDefinition {id:$definitionId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	ExecuteWriteQuery(driver, query, map[string]any{
		"definitionId":    definitionId.String(),
		"setDefinitionId": setDefinitionId,
		"type":            "TEXT",
		"order":           1,
		"mandatory":       false,
		"name":            "definition name",
	})
	return definitionId.String()
}

func AddSetDefinitionToEntity(driver *neo4j.Driver, entityDefinitionId string) string {
	var definitionId, _ = uuid.NewRandom()
	query := `MATCH (e:EntityDefinition {id:$entityDefinitionId})
			MERGE (f:FieldSetDefinition {id:$definitionId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	ExecuteWriteQuery(driver, query, map[string]any{
		"definitionId":       definitionId.String(),
		"entityDefinitionId": entityDefinitionId,
		"type":               "TEXT",
		"order":              1,
		"mandatory":          false,
		"name":               "set name",
	})
	return definitionId.String()
}

func CreateContactType(driver *neo4j.Driver, tenant, contactTypeName string) string {
	var contactTypeId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)-[:USES_CONTACT_TYPE]->(c:ContactType {id:$id})
			ON CREATE SET c.name=$name`
	ExecuteWriteQuery(driver, query, map[string]any{
		"id":     contactTypeId.String(),
		"tenant": tenant,
		"name":   contactTypeName,
	})
	return contactTypeId.String()
}

func CreateCompany(driver *neo4j.Driver, tenant, companyName string) string {
	var companyId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:COMPANY_BELONGS_TO_TENANT]-(co:Company {id:$id})
			ON CREATE SET co.name=$name`
	ExecuteWriteQuery(driver, query, map[string]any{
		"id":     companyId.String(),
		"tenant": tenant,
		"name":   companyName,
	})
	return companyId.String()
}

func ContactWorksForCompany(driver *neo4j.Driver, contactId, companyId, jobTitle string) string {
	var positionId, _ = uuid.NewRandom()
	query := `MATCH (c:Contact {id:$contactId}),
			        (co:Company {id:$companyId})
			MERGE (c)-[:WORKS_AT {id:$id, jobTitle:$jobTitle}]->(co)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"id":        positionId.String(),
		"contactId": contactId,
		"companyId": companyId,
		"jobTitle":  jobTitle,
	})
	return positionId.String()
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
			MERGE (u)-[:PARTICIPATES]->(o:Conversation {id:$conversationId, startedAt:datetime({timezone: 'UTC'})})<-[:PARTICIPATES]-(c)`
	ExecuteWriteQuery(driver, query, map[string]any{
		"contactId":      contactId,
		"userId":         userId,
		"conversationId": conversationId.String(),
	})
	return conversationId.String()
}

func AddMessageToConversation(driver *neo4j.Driver, conversationId, messageChannel string, time time.Time) string {
	var messageId, _ = uuid.NewRandom()
	query := `MATCH (c:Conversation {id:$conversationId})
			MERGE (c)-[:CONSISTS_OF]->(m:Message:Action {id:$messageId})
			ON CREATE SET m.channel=$channel, m.startedAt=$startedAt, m.conversationId=$conversationId`
	ExecuteWriteQuery(driver, query, map[string]any{
		"conversationId": conversationId,
		"messageId":      messageId.String(),
		"channel":        messageChannel,
		"startedAt":      time,
	})
	return messageId.String()
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
