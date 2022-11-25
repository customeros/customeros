package resolver

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/integration_tests"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func cleanupAllData(driver *neo4j.Driver) {
	integration_tests.ExecuteWriteQuery(driver, `MATCH (n) DETACH DELETE n`, map[string]any{})
}

func createTenant(driver *neo4j.Driver, tenant string) {
	query := `MERGE (t:Tenant {name:$tenant})`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"tenant": tenant,
	})
}

func createDefaultUser(driver *neo4j.Driver, tenant string) string {
	return createUser(driver, tenant, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Email:     "user@openline.ai",
	})
}

func createUser(driver *neo4j.Driver, tenant string, user entity.UserEntity) string {
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
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"tenant":    tenant,
		"userId":    userId.String(),
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	})
	return userId.String()
}

func createDefaultContact(driver *neo4j.Driver, tenant string) string {
	return createContact(driver, tenant, entity.ContactEntity{Title: "MR", FirstName: "first", LastName: "last"})
}

func createContact(driver *neo4j.Driver, tenant string, contact entity.ContactEntity) string {
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
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
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

func setContactTypeForContact(driver *neo4j.Driver, contactId, contactTypeId string) {
	query := `
			MATCH (c:Contact {id:$contactId}),
				  (o:ContactType {id:$contactTypeId})
			MERGE (c)-[:IS_OF_TYPE]->(o)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"contactId":     contactId,
		"contactTypeId": contactTypeId,
	})
}

func createContactGroup(driver *neo4j.Driver, tenant, name string) string {
	var contactGroupId, _ = uuid.NewRandom()
	query := `
			MATCH (t:Tenant {name:$tenant})
			MERGE (g:ContactGroup {
				  id: $id,
				  name: $name
				})-[:GROUP_BELONGS_TO_TENANT]->(t)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"tenant": tenant,
		"id":     contactGroupId.String(),
		"name":   name,
	})
	return contactGroupId.String()
}

func addContactToGroup(driver *neo4j.Driver, contactId, groupId string) {
	query := `MATCH (c:Contact {id:$contactId}), (g:ContactGroup {id:$groupId})
				MERGE (c)-[:BELONGS_TO_GROUP]->(g)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"groupId":   groupId,
	})
}

func createDefaultFieldSet(driver *neo4j.Driver, contactId string) string {
	return createFieldSet(driver, contactId, entity.FieldSetEntity{Name: "name"})
}

func createFieldSet(driver *neo4j.Driver, contactId string, fieldSet entity.FieldSetEntity) string {
	var fieldSetId, _ = uuid.NewRandom()
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (s:FieldSet {
				  id: $fieldSetId,
				  name: $name
				})<-[:HAS_COMPLEX_PROPERTY {added:datetime({timezone: 'UTC'})}]-(c)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"contactId":  contactId,
		"fieldSetId": fieldSetId.String(),
		"name":       fieldSet.Name,
	})
	return fieldSetId.String()
}

func createDefaultCustomFieldInSet(driver *neo4j.Driver, fieldSetId string) string {
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
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"fieldSetId": fieldSetId,
		"fieldId":    fieldId.String(),
		"name":       customField.Name,
		"datatype":   customField.DataType,
		"value":      customField.Value.RealValue(),
	})
	return fieldId.String()
}

func createDefaultCustomFieldInContact(driver *neo4j.Driver, contactId string) string {
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
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"fieldId":   fieldId.String(),
		"name":      customField.Name,
		"datatype":  customField.DataType,
		"value":     customField.Value.RealValue(),
	})
	return fieldId.String()
}

func addEmailToContact(driver *neo4j.Driver, contactId string, email string, primary bool, label string) {
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (:Email {
				  id: randomUUID(),
				  email: $email,
				  label: $label
				})<-[:EMAILED_AT {primary:$primary}]-(c)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"primary":   primary,
		"email":     email,
		"label":     label,
	})
}

func addPhoneNumberToContact(driver *neo4j.Driver, contactId string, e164 string, primary bool, label string) {
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (:PhoneNumber {
				  id: randomUUID(),
				  e164: $e164,
				  label: $label
				})<-[:CALLED_AT {primary:$primary}]-(c)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"primary":   primary,
		"e164":      e164,
		"label":     label,
	})
}

func createEntityDefinition(driver *neo4j.Driver, tenant, extends string) string {
	var definitionId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:EntityDefinition {id:$definitionId})<-[:USES_ENTITY_DEFINITION]-(t)
			ON CREATE SET e.extends=$extends, e.name=$name`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"definitionId": definitionId.String(),
		"tenant":       tenant,
		"extends":      extends,
		"name":         "definition name",
	})
	return definitionId.String()
}

func linkEntityDefinitionToContact(driver *neo4j.Driver, entityDefinitionId, contactId string) {
	query := `MATCH (c:Contact {id:$contactId}),
			(e:EntityDefinition {id:$definitionId})
			MERGE (c)-[:IS_DEFINED_BY]->(e)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"definitionId": entityDefinitionId,
		"contactId":    contactId,
	})
}

func addFieldDefinitionToEntity(driver *neo4j.Driver, entityDefinitionId string) string {
	var definitionId, _ = uuid.NewRandom()
	query := `MATCH (e:EntityDefinition {id:$entityDefinitionId})
			MERGE (f:CustomFieldDefinition {id:$definitionId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"definitionId":       definitionId.String(),
		"entityDefinitionId": entityDefinitionId,
		"type":               "TEXT",
		"order":              1,
		"mandatory":          false,
		"name":               "definition name",
	})
	return definitionId.String()
}

func addFieldDefinitionToSet(driver *neo4j.Driver, setDefinitionId string) string {
	var definitionId, _ = uuid.NewRandom()
	query := `MATCH (e:FieldSetDefinition {id:$setDefinitionId})
			MERGE (f:CustomFieldDefinition {id:$definitionId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"definitionId":    definitionId.String(),
		"setDefinitionId": setDefinitionId,
		"type":            "TEXT",
		"order":           1,
		"mandatory":       false,
		"name":            "definition name",
	})
	return definitionId.String()
}

func addSetDefinitionToEntity(driver *neo4j.Driver, entityDefinitionId string) string {
	var definitionId, _ = uuid.NewRandom()
	query := `MATCH (e:EntityDefinition {id:$entityDefinitionId})
			MERGE (f:FieldSetDefinition {id:$definitionId})<-[:CONTAINS]-(e)
			ON CREATE SET f.name=$name, f.type=$type, f.order=$order, f.mandatory=$mandatory`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"definitionId":       definitionId.String(),
		"entityDefinitionId": entityDefinitionId,
		"type":               "TEXT",
		"order":              1,
		"mandatory":          false,
		"name":               "set name",
	})
	return definitionId.String()
}

func createContactType(driver *neo4j.Driver, tenant, contactTypeName string) string {
	var contactTypeId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)-[:USES_CONTACT_TYPE]->(c:ContactType {id:$id})
			ON CREATE SET c.name=$name`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"id":     contactTypeId.String(),
		"tenant": tenant,
		"name":   contactTypeName,
	})
	return contactTypeId.String()
}

func createCompany(driver *neo4j.Driver, tenant, companyName string) string {
	var companyId, _ = uuid.NewRandom()
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:COMPANY_BELONGS_TO_TENANT]-(co:Company {id:$id})
			ON CREATE SET co.name=$name`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"id":     companyId.String(),
		"tenant": tenant,
		"name":   companyName,
	})
	return companyId.String()
}

func contactWorksForCompany(driver *neo4j.Driver, contactId, companyId, jobTitle string) string {
	var positionId, _ = uuid.NewRandom()
	query := `MATCH (c:Contact {id:$contactId}),
			        (co:Company {id:$companyId})
			MERGE (c)-[:WORKS_AT {id:$id, jobTitle:$jobTitle}]->(co)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"id":        positionId.String(),
		"contactId": contactId,
		"companyId": companyId,
		"jobTitle":  jobTitle,
	})
	return positionId.String()
}

func userOwnsContact(driver *neo4j.Driver, userId, contactId string) {
	query := `MATCH (c:Contact {id:$contactId}),
			        (u:User {id:$userId})
			MERGE (u)-[:OWNS]->(c)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"userId":    userId,
	})
}

func createConversation(driver *neo4j.Driver, userId, contactId string) string {
	var conversationId, _ = uuid.NewRandom()
	query := `MATCH (c:Contact {id:$contactId}),
			        (u:User {id:$userId})
			MERGE (u)-[:PARTICIPATES]->(o:Conversation {id:$conversationId, startedAt:datetime({timezone: 'UTC'})})<-[:PARTICIPATES]-(c)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"contactId":      contactId,
		"userId":         userId,
		"conversationId": conversationId.String(),
	})
	return conversationId.String()
}

func getCountOfNodes(driver *neo4j.Driver, nodeLabel string) int {
	query := fmt.Sprintf(`MATCH (n:%s) RETURN count(n)`, nodeLabel)
	result := integration_tests.ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func getCountOfRelationships(driver *neo4j.Driver, relationship string) int {
	query := fmt.Sprintf(`MATCH (a)-[r:%s]-(b) RETURN count(distinct r)`, relationship)
	result := integration_tests.ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}
